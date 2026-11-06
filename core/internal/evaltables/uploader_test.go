package evaltables_test

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/wandb/wandb/core/internal/api"
	"github.com/wandb/wandb/core/internal/evaltables"
	"github.com/wandb/wandb/core/internal/gqlmock"
	"github.com/wandb/wandb/core/internal/observabilitytest"
	"github.com/wandb/wandb/core/internal/runhandle"
	"github.com/wandb/wandb/core/internal/runupsertertest"
	"github.com/wandb/wandb/core/internal/settings"
	"github.com/wandb/wandb/core/internal/wboperation"
	spb "github.com/wandb/wandb/core/pkg/service_go_proto"
)

const tablePath = "/evaluations/v1/namespaces/wandb/scopes/project-internal-id/eval-tables"

type cesRequest struct {
	Path           string
	IdempotencyKey string
	Body           string
}

// fakeCES is a CES server that records requests and returns fixed IDs.
type fakeCES struct {
	*httptest.Server

	mu       sync.Mutex
	requests []cesRequest

	// onRequest, if set, is called with each request's path before
	// responding.
	onRequest func(path string)
}

func newFakeCES(t *testing.T) *fakeCES {
	t.Helper()
	ces := &fakeCES{}

	ces.Server = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)

			ces.mu.Lock()
			ces.requests = append(ces.requests, cesRequest{
				Path:           r.URL.Path,
				IdempotencyKey: r.Header.Get("Idempotency-Key"),
				Body:           string(body),
			})
			ces.mu.Unlock()

			if ces.onRequest != nil {
				ces.onRequest(r.URL.Path)
			}

			w.Header().Set("Content-Type", "application/json")
			switch {
			case r.URL.Path == tablePath:
				_, _ = w.Write([]byte(`{
					"evaluation_id": "eval-id",
					"dataset_id": "dataset-id"
				}`))
			case strings.HasSuffix(r.URL.Path, "/versions"):
				_, _ = w.Write([]byte(`{
					"evaluation_version_id": "eval-version-id",
					"dataset_version_id": "dataset-version-id"
				}`))
			default:
				_, _ = w.Write([]byte(`{}`))
			}
		}))
	t.Cleanup(ces.Close)

	return ces
}

func (c *fakeCES) Requests() []cesRequest {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.requests
}

const projectFoundJSON = `{"project": {"internalId": "project-internal-id"}}`

// newTestUploader returns an Uploader for a run in test-entity/test-project
// that uploads to ces.
//
// projectJSON is the response to the project's internal ID lookup.
func newTestUploader(
	t *testing.T,
	ces *fakeCES,
	projectJSON string,
) (*evaltables.Uploader, *gqlmock.MockClient) {
	t.Helper()
	return newTestUploaderWithSettings(t, projectJSON, settings.From(&spb.Settings{
		XCesBaseUrl: wrapperspb.String(ces.URL + "/evaluations"),
	}))
}

func newTestUploaderWithSettings(
	t *testing.T,
	projectJSON string,
	s *settings.Settings,
) (*evaltables.Uploader, *gqlmock.MockClient) {
	t.Helper()

	mockGQL := gqlmock.NewMockClient()
	mockGQL.StubMatchOnce(gqlmock.WithOpName("ProjectInternalID"), projectJSON)

	runHandle := runhandle.New()
	require.NoError(t, runHandle.Init(runupsertertest.NewOfflineUpserter(t)))

	return evaltables.NewUploader(evaltables.UploaderParams{
		CredentialProvider: api.NewAPIKeyCredentialProvider("test-api-key"),
		GraphqlClient:      mockGQL,
		Logger:             observabilitytest.NewTestLogger(t),
		RunHandle:          runHandle,
		Settings:           s,
	}), mockGQL
}

func TestUpload(t *testing.T) {
	ces := newFakeCES(t)
	uploader, mockGQL := newTestUploader(t, ces, projectFoundJSON)

	response, err := uploader.Upload(context.Background(), &spb.UploadEvalTableRequest{
		Name: "my-table",
		DatasetFields: []*spb.EvalTableField{
			{Source: "input", Name: "prompt", ValueType: "string"},
			{Source: "output", Name: "answer", ValueType: "string"},
		},
		Scorers: []*spb.EvalTableScorer{{Name: "accuracy", ValueType: "number"}},
		RowsJson: []string{
			`{"input":{"prompt":"a"},"output":{"answer":"x"},"scores":{"accuracy":1}}`,
			`{"input":{"prompt":"b"},"output":{"answer":"y"},"scores":{"accuracy":0}}`,
		},
	})

	require.NoError(t, err)
	assert.Equal(t, &spb.UploadEvalTableResponse{
		EvaluationId:        "eval-id",
		EvaluationVersionId: "eval-version-id",
		DatasetId:           "dataset-id",
		DatasetVersionId:    "dataset-version-id",
	}, response)
	gqlmock.AssertVariables(t, mockGQL.AllRequests()[0],
		gqlmock.GQLVar("entityName", gomock.Eq("test-entity")),
		gqlmock.GQLVar("projectName", gomock.Eq("test-project")))

	requests := ces.Requests()
	require.Len(t, requests, 4)
	assert.Equal(t, tablePath, requests[0].Path)
	assert.JSONEq(t, `{"name": "my-table"}`, requests[0].Body)
	assert.Equal(t, tablePath+"/eval-id/columns", requests[1].Path)
	assert.JSONEq(t, `{
		"dataset_fields": [
			{"source": "input", "name": "prompt", "value_type": "string"},
			{"source": "output", "name": "answer", "value_type": "string"}
		],
		"scorers": [{"name": "accuracy", "value_type": "number"}]
	}`, requests[1].Body)
	assert.Equal(t, tablePath+"/eval-id/rows", requests[2].Path)
	assert.JSONEq(t, `{"rows": [
		{"input": {"prompt": "a"}, "output": {"answer": "x"}, "scores": {"accuracy": 1}},
		{"input": {"prompt": "b"}, "output": {"answer": "y"}, "scores": {"accuracy": 0}}
	]}`, requests[2].Body)
	assert.Equal(t, tablePath+"/eval-id/versions", requests[3].Path)
}

func TestUpload_DeclaresExtensionColumns(t *testing.T) {
	ces := newFakeCES(t)
	uploader, _ := newTestUploader(t, ces, projectFoundJSON)

	_, err := uploader.Upload(context.Background(), &spb.UploadEvalTableRequest{
		Name: "my-table",
		DatasetFields: []*spb.EvalTableField{{
			Source:                 "output",
			Name:                   "image",
			ValueType:              "json",
			ExtensionType:          "wandb-image",
			ExtensionSchemaVersion: 1,
		}},
	})

	require.NoError(t, err)
	requests := ces.Requests()
	require.Len(t, requests, 3) // No rows request for an empty table.
	assert.Equal(t, tablePath+"/eval-id/columns", requests[1].Path)
	assert.JSONEq(t, `{
		"dataset_fields": [{
			"source": "output",
			"name": "image",
			"value_type": "json",
			"extension_type": "wandb-image",
			"extension_schema_version": 1
		}],
		"scorers": []
	}`, requests[1].Body)
}

func TestUpload_IdempotencyKeys(t *testing.T) {
	ces := newFakeCES(t)
	uploader, _ := newTestUploader(t, ces, projectFoundJSON)
	request := &spb.UploadEvalTableRequest{
		Name:     "my-table",
		RowsJson: []string{`{}`},
	}

	_, err := uploader.Upload(context.Background(), request)
	require.NoError(t, err)
	_, err = uploader.Upload(context.Background(), request)
	require.NoError(t, err)

	requests := ces.Requests()
	require.Len(t, requests, 8)
	keys := make([]string, len(requests))
	for i, req := range requests {
		keys[i] = req.IdempotencyKey
	}
	uuidPattern := `[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}`
	for _, upload := range [][]string{keys[:4], keys[4:]} {
		prefix, found := strings.CutSuffix(upload[0], "-create")
		require.True(t, found, upload[0])
		assert.Regexp(t, `^wandb-eval-table-v3-`+uuidPattern+`$`, prefix)
		assert.Equal(t,
			[]string{
				prefix + "-create",
				prefix + "-columns",
				prefix + "-rows-0",
				prefix + "-version",
			},
			upload)
	}
	// Uploading the same table again, such as from another writer in
	// shared mode, must not collide with the first upload's keys.
	assert.NotEqual(t, keys[0], keys[4])
}

func TestUpload_SplitsRowsIntoBatches(t *testing.T) {
	ces := newFakeCES(t)
	uploader, _ := newTestUploader(t, ces, projectFoundJSON)
	// Each row is 1 MiB, so at most 7 fit under the 8 MiB target with the
	// surrounding JSON.
	row, err := json.Marshal(map[string]string{"x": strings.Repeat("a", 1<<20)})
	require.NoError(t, err)
	rows := make([]string, 10)
	for i := range rows {
		rows[i] = string(row)
	}

	_, err = uploader.Upload(context.Background(),
		&spb.UploadEvalTableRequest{Name: "t", RowsJson: rows})
	require.NoError(t, err)

	var batchSizes []int
	for _, req := range ces.Requests() {
		if !strings.HasSuffix(req.Path, "/rows") {
			continue
		}
		var body struct{ Rows []json.RawMessage }
		require.NoError(t, json.Unmarshal([]byte(req.Body), &body))
		assert.LessOrEqual(t, len(req.Body), 8<<20)
		batchSizes = append(batchSizes, len(body.Rows))
	}
	assert.Equal(t, []int{7, 3}, batchSizes)
}

func TestUpload_ReportsRowProgress(t *testing.T) {
	ces := newFakeCES(t)
	uploader, _ := newTestUploader(t, ces, projectFoundJSON)
	// Each row is 1 MiB, so the rows are sent in batches of 7 and 3.
	row, err := json.Marshal(map[string]string{"x": strings.Repeat("a", 1<<20)})
	require.NoError(t, err)
	rows := make([]string, 10)
	for i := range rows {
		rows[i] = string(row)
	}
	op := wboperation.NewOperations().New("uploading eval table t")
	var progress []string
	ces.onRequest = func(path string) {
		if strings.HasSuffix(path, "/rows") || strings.HasSuffix(path, "/versions") {
			progress = append(progress, op.ToProto().GetProgress())
		}
	}

	_, err = uploader.Upload(op.Context(context.Background()),
		&spb.UploadEvalTableRequest{Name: "t", RowsJson: rows})

	require.NoError(t, err)
	assert.Equal(t, []string{"", "7/10 rows", "10/10 rows"}, progress)
}

// writeRowsFile writes rows to a gzipped file, one per line.
func writeRowsFile(t *testing.T, rows ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "rows.jl.gz")
	file, err := os.Create(path)
	require.NoError(t, err)
	gzWriter := gzip.NewWriter(file)
	for _, row := range rows {
		_, err = gzWriter.Write([]byte(row + "\n"))
		require.NoError(t, err)
	}
	require.NoError(t, gzWriter.Close())
	require.NoError(t, file.Close())
	return path
}

func TestUpload_ReadsAndDeletesRowsFile(t *testing.T) {
	ces := newFakeCES(t)
	uploader, _ := newTestUploader(t, ces, projectFoundJSON)
	// Longer than bufio.Scanner's default limit.
	longRow := fmt.Sprintf(`{"x":%q}`, strings.Repeat("a", 100_000))
	path := writeRowsFile(t, `{"x":1}`, longRow)

	_, err := uploader.Upload(context.Background(),
		&spb.UploadEvalTableRequest{Name: "t", RowsFilePath: path})

	require.NoError(t, err)
	requests := ces.Requests()
	require.Len(t, requests, 4)
	assert.Equal(t, `{"rows":[{"x":1},`+longRow+`]}`, requests[2].Body)
	assert.NoFileExists(t, path)
}

func TestUpload_ReportsRowsFileProgressWithoutTotal(t *testing.T) {
	ces := newFakeCES(t)
	uploader, _ := newTestUploader(t, ces, projectFoundJSON)
	op := wboperation.NewOperations().New("uploading eval table t")

	_, err := uploader.Upload(op.Context(context.Background()),
		&spb.UploadEvalTableRequest{
			Name:         "t",
			RowsFilePath: writeRowsFile(t, `{"x":1}`, `{"x":2}`),
		})

	require.NoError(t, err)
	assert.Equal(t, "2 rows", op.ToProto().GetProgress())
}

func TestUpload_DeletesRowsFileOnFailure(t *testing.T) {
	ces := newFakeCES(t)
	uploader, _ := newTestUploader(t, ces, `{"project": null}`)
	path := writeRowsFile(t, `{"x":1}`)

	_, err := uploader.Upload(context.Background(),
		&spb.UploadEvalTableRequest{Name: "t", RowsFilePath: path})

	require.Error(t, err)
	assert.NoFileExists(t, path)
}

func TestUpload_NoCESURL(t *testing.T) {
	uploader, _ := newTestUploaderWithSettings(t, projectFoundJSON, settings.New())

	_, err := uploader.Upload(context.Background(),
		&spb.UploadEvalTableRequest{Name: "t"})

	assert.ErrorContains(t, err, "no CES URL in settings")
}

func TestUpload_ProjectNotFound(t *testing.T) {
	ces := newFakeCES(t)
	uploader, _ := newTestUploader(t, ces, `{"project": null}`)

	_, err := uploader.Upload(context.Background(),
		&spb.UploadEvalTableRequest{Name: "t"})

	assert.ErrorContains(t, err, "project test-entity/test-project not found")
	assert.Empty(t, ces.Requests())
}
