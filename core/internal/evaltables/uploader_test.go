package evaltables_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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

	// errorStatus, if set, returns an HTTP status to fail a request with,
	// or 0 to respond normally.
	errorStatus func(path string) int
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
			if ces.errorStatus != nil {
				if status := ces.errorStatus(r.URL.Path); status != 0 {
					w.WriteHeader(status)
					return
				}
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

func TestUpload_PassesThroughAllValueTypes(t *testing.T) {
	ces := newFakeCES(t)
	uploader, _ := newTestUploader(t, ces, projectFoundJSON)
	// The large integer and 1.0 would change if decoded as float64.
	row := `{` +
		`"input":{` +
		`"flag":true,` +
		`"count":9007199254740993,` +
		`"temperature":1.0,` +
		`"prompt":"draw a cat",` +
		`"config":{"seed":7,"tags":["a","b"]}` +
		`},` +
		`"output":{"answer":"a cat"},` +
		`"scores":{"passed":false,"votes":3,"accuracy":0.5,"grade":"B"}` +
		`}`

	_, err := uploader.Upload(context.Background(), &spb.UploadEvalTableRequest{
		Name: "my-table",
		DatasetFields: []*spb.EvalTableField{
			{Source: "input", Name: "flag", ValueType: "boolean"},
			{Source: "input", Name: "count", ValueType: "integer"},
			{Source: "input", Name: "temperature", ValueType: "number"},
			{Source: "input", Name: "prompt", ValueType: "string"},
			{Source: "input", Name: "config", ValueType: "json"},
			{Source: "output", Name: "answer", ValueType: "string"},
		},
		Scorers: []*spb.EvalTableScorer{
			{Name: "passed", ValueType: "boolean"},
			{Name: "votes", ValueType: "integer"},
			{Name: "accuracy", ValueType: "number"},
			{Name: "grade", ValueType: "string"},
		},
		RowsJson: []string{row},
	})

	require.NoError(t, err)
	requests := ces.Requests()
	require.Len(t, requests, 4)
	assert.JSONEq(t, `{
		"dataset_fields": [
			{"source": "input", "name": "flag", "value_type": "boolean"},
			{"source": "input", "name": "count", "value_type": "integer"},
			{"source": "input", "name": "temperature", "value_type": "number"},
			{"source": "input", "name": "prompt", "value_type": "string"},
			{"source": "input", "name": "config", "value_type": "json"},
			{"source": "output", "name": "answer", "value_type": "string"}
		],
		"scorers": [
			{"name": "passed", "value_type": "boolean"},
			{"name": "votes", "value_type": "integer"},
			{"name": "accuracy", "value_type": "number"},
			{"name": "grade", "value_type": "string"}
		]
	}`, requests[1].Body)
	assert.Equal(t, `{"rows":[`+row+`]}`, requests[2].Body)
}

func TestUpload_DistinctIdempotencyKeys(t *testing.T) {
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
	keys := make(map[string]bool)
	for _, req := range requests {
		require.NotEmpty(t, req.IdempotencyKey, req.Path)
		keys[req.IdempotencyKey] = true
	}
	// Each request needs its own key, since CES rejects a reused key with a
	// different body. Uploading the same table again, such as from another
	// writer in shared mode, must not reuse the first upload's keys.
	assert.Len(t, keys, len(requests))
}

func TestUpload_RetriesReuseIdempotencyKey(t *testing.T) {
	evaltables.SetRetryWaitMinForTest(t, time.Millisecond)
	ces := newFakeCES(t)
	var failed atomic.Bool
	ces.errorStatus = func(path string) int {
		// Fail only the first rows request, with a status that makes the CES HTTP
		// client retry.
		if strings.HasSuffix(path, "/rows") && failed.CompareAndSwap(false, true) {
			return http.StatusServiceUnavailable
		}
		return 0
	}
	uploader, _ := newTestUploader(t, ces, projectFoundJSON)

	_, err := uploader.Upload(context.Background(), &spb.UploadEvalTableRequest{
		Name:     "my-table",
		RowsJson: []string{`{}`},
	})
	require.NoError(t, err)

	var rowsRequests []cesRequest
	for _, req := range ces.Requests() {
		if strings.HasSuffix(req.Path, "/rows") {
			rowsRequests = append(rowsRequests, req)
		}
	}
	// The failed attempt and its retry. The retry must reuse the key so
	// that CES replays it instead of adding the rows again.
	require.Len(t, rowsRequests, 2)
	assert.Equal(t,
		rowsRequests[0].IdempotencyKey,
		rowsRequests[1].IdempotencyKey)
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

func TestUpload_RejectsUnknownColumnTypes(t *testing.T) {
	testCases := []struct {
		name    string
		request *spb.UploadEvalTableRequest
		wantErr string
	}{
		{
			name: "source",
			request: &spb.UploadEvalTableRequest{
				DatasetFields: []*spb.EvalTableField{
					{Source: "metadata", Name: "a", ValueType: "string"},
				},
			},
			wantErr: `field "a" has unknown source "metadata"`,
		},
		{
			name: "field value type",
			request: &spb.UploadEvalTableRequest{
				DatasetFields: []*spb.EvalTableField{
					{Source: "input", Name: "a", ValueType: "float"},
				},
			},
			wantErr: `field "a" has unknown value type "float"`,
		},
		{
			name: "scorer value type",
			request: &spb.UploadEvalTableRequest{
				Scorers: []*spb.EvalTableScorer{{Name: "a", ValueType: "json"}},
			},
			wantErr: `scorer "a" has unknown value type "json"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ces := newFakeCES(t)
			uploader, _ := newTestUploader(t, ces, projectFoundJSON)

			_, err := uploader.Upload(context.Background(), tc.request)

			assert.ErrorContains(t, err, tc.wantErr)
			assert.Empty(t, ces.Requests())
		})
	}
}

func TestUpload_RejectsOversizedColumns(t *testing.T) {
	ces := newFakeCES(t)
	uploader, _ := newTestUploader(t, ces, projectFoundJSON)

	// Each "<" is escaped to 6 bytes, so the name alone exceeds the limit.
	name := strings.Repeat("<", 3<<20)
	_, err := uploader.Upload(context.Background(), &spb.UploadEvalTableRequest{
		DatasetFields: []*spb.EvalTableField{
			{Source: "input", Name: name, ValueType: "string"},
		},
	})

	assert.ErrorContains(t, err, "columns request is")
	assert.Empty(t, ces.Requests())
}

func TestUpload_ConvertsMedia(t *testing.T) {
	ces := newFakeCES(t)
	uploader, _ := newTestUploader(t, ces, projectFoundJSON)
	// W&B media JSON, as from Python's Media.to_json.
	row := `{` +
		`"input":{"prompt":"draw <a> cat & dog","count":9007199254740993,"temperature":1.0},` +
		`"output":{` +
		`"answer":"<b>a cat</b> & a dog",` +
		`"tokens":9007199254740993,` +
		`"ratio":1.0,` +
		`"image":{` +
		`"_type":"image-file",` +
		`"path":"media/eval_tables/images/eval/cat.png",` +
		`"sha256":"aaa",` +
		`"size":1024,` +
		`"format":"png",` +
		`"width":32,` +
		`"height":24,` +
		`"caption":"a cat",` +
		`"grouping":1,` +
		`"boxes":{"predictions":{` +
		`"_type":"boxes2D",` +
		`"path":"media/eval_tables/boxes/eval/cat.boxes2D.json",` +
		`"sha256":"bbb",` +
		`"size":10` +
		`}},` +
		`"masks":{"truth":{` +
		`"_type":"mask",` +
		`"artifact_path":"wandb-artifact://_id/abc/mask.png",` +
		`"sha256":"ccc",` +
		`"size":20` +
		`}}` +
		`},` +
		`"audio":{` +
		`"_type":"audio-file",` +
		`"artifact_path":"wandb-artifact://_id/def/meow.wav",` +
		`"_latest_artifact_path":"wandb-artifact://entity/project/sounds:latest/meow.wav",` +
		`"sha256":"ddd",` +
		`"size":2048,` +
		`"caption":"meow"` +
		`},` +
		`"video":{` +
		`"_type":"video-file",` +
		`"path":"media/eval tables/videos/cat#1.mp4",` +
		`"sha256":"eee",` +
		`"size":4096,` +
		`"width":64,` +
		`"height":48` +
		`}` +
		`},` +
		`"scores":{}` +
		`}`

	_, err := uploader.Upload(context.Background(), &spb.UploadEvalTableRequest{
		Name: "my-table",
		DatasetFields: []*spb.EvalTableField{
			{Source: "input", Name: "prompt", ValueType: "string"},
			{Source: "input", Name: "count", ValueType: "integer"},
			{Source: "input", Name: "temperature", ValueType: "number"},
			{Source: "output", Name: "answer", ValueType: "string"},
			{Source: "output", Name: "tokens", ValueType: "integer"},
			{Source: "output", Name: "ratio", ValueType: "number"},
			{Source: "output", Name: "image", ValueType: "json", WandbMediaType: "image-file"},
			{Source: "output", Name: "audio", ValueType: "json", WandbMediaType: "audio-file"},
			{Source: "output", Name: "video", ValueType: "json", WandbMediaType: "video-file"},
		},
		RowsJson: []string{row},
	})

	require.NoError(t, err)
	requests := ces.Requests()
	require.Len(t, requests, 4)
	assert.JSONEq(t, `{
		"dataset_fields": [
			{"source": "input", "name": "prompt", "value_type": "string"},
			{"source": "input", "name": "count", "value_type": "integer"},
			{"source": "input", "name": "temperature", "value_type": "number"},
			{"source": "output", "name": "answer", "value_type": "string"},
			{"source": "output", "name": "tokens", "value_type": "integer"},
			{"source": "output", "name": "ratio", "value_type": "number"},
			{
				"source": "output",
				"name": "image",
				"value_type": "json",
				"extension_type": "wandb-image",
				"extension_schema_version": 1
			},
			{
				"source": "output",
				"name": "audio",
				"value_type": "json",
				"extension_type": "wandb-audio",
				"extension_schema_version": 1
			},
			{
				"source": "output",
				"name": "video",
				"value_type": "json",
				"extension_type": "wandb-video",
				"extension_schema_version": 1
			}
		],
		"scorers": []
	}`, requests[1].Body)
	assert.JSONEq(t, `{"rows": [{
		"input": {
			"prompt": "draw <a> cat & dog",
			"count": 9007199254740993,
			"temperature": 1.0
		},
		"output": {
			"answer": "<b>a cat</b> & a dog",
			"tokens": 9007199254740993,
			"ratio": 1.0,
			"image": {
				"extension_type": "wandb-image",
				"schema_version": 1,
				"wb_media_type": "image-file",
				"format": "png",
				"sha256": "aaa",
				"size": 1024,
				"uri": "wandb-run-file://test-entity/test-project/test-run/media/eval_tables/images/eval/cat.png",
				"width": 32,
				"height": 24,
				"caption": "a cat",
				"boxes": {"predictions": {
					"_type": "boxes2D",
					"sha256": "bbb",
					"size": 10,
					"uri": "wandb-run-file://test-entity/test-project/test-run/media/eval_tables/boxes/eval/cat.boxes2D.json"
				}},
				"masks": {"truth": {
					"_type": "mask",
					"sha256": "ccc",
					"size": 20,
					"uri": "wandb-artifact://_id/abc/mask.png"
				}}
			},
			"audio": {
				"extension_type": "wandb-audio",
				"schema_version": 1,
				"wb_media_type": "audio-file",
				"sha256": "ddd",
				"size": 2048,
				"uri": "wandb-artifact://_id/def/meow.wav",
				"caption": "meow"
			},
			"video": {
				"extension_type": "wandb-video",
				"schema_version": 1,
				"wb_media_type": "video-file",
				"sha256": "eee",
				"size": 4096,
				"uri": "wandb-run-file://test-entity/test-project/test-run/media/eval%20tables/videos/cat%231.mp4",
				"width": 64,
				"height": 48
			}
		},
		"scores": {}
	}]}`, requests[2].Body)
	// JSONEq compares numbers as float64, so check that other values keep
	// their exact JSON: unchanged without media, and re-encoded otherwise.
	assert.Contains(t, requests[2].Body,
		`"input":{"prompt":"draw <a> cat & dog","count":9007199254740993,"temperature":1.0}`)
	assert.Contains(t, requests[2].Body, `"answer":"<b>a cat</b> & a dog"`)
	assert.Contains(t, requests[2].Body, `"ratio":1.0`)
	assert.Contains(t, requests[2].Body, `"tokens":9007199254740993`)
}

func TestUpload_KeepsNullMedia(t *testing.T) {
	ces := newFakeCES(t)
	uploader, _ := newTestUploader(t, ces, projectFoundJSON)

	_, err := uploader.Upload(context.Background(), &spb.UploadEvalTableRequest{
		DatasetFields: []*spb.EvalTableField{
			{Source: "input", Name: "image", ValueType: "json", WandbMediaType: "image-file"},
			{Source: "output", Name: "image", ValueType: "json", WandbMediaType: "image-file"},
		},
		RowsJson: []string{
			`{"input":{"image":null},"output":null,"scores":{}}`,
			`{"input":{},"scores":{}}`,
		},
	})

	require.NoError(t, err)
	requests := ces.Requests()
	require.Len(t, requests, 4)
	assert.JSONEq(t,
		`{"rows":[`+
			`{"input":{"image":null},"output":null,"scores":{}},`+
			`{"input":{},"scores":{}}`+
			`]}`,
		requests[2].Body)
}

func TestUpload_RejectsInvalidMedia(t *testing.T) {
	testCases := []struct {
		name      string
		mediaType string
		media     string
		wantErr   string
	}{
		{
			name:      "unknown media type",
			mediaType: "html-file",
			media:     `{"_type":"html-file","path":"a.html","sha256":"a","size":1}`,
			wantErr:   `field "media" has unknown media type "html-file"`,
		},
		{
			name:      "mismatched media type",
			mediaType: "image-file",
			media:     `{"_type":"audio-file","path":"a.wav","sha256":"a","size":1}`,
			wantErr:   `field "media": expected image-file media, got "audio-file"`,
		},
		{
			name:      "no file",
			mediaType: "audio-file",
			media:     `{"_type":"audio-file","sha256":"a","size":1}`,
			wantErr:   `field "media": media has no path or artifact path`,
		},
		{
			name:      "overlay without file",
			mediaType: "image-file",
			media: `{"_type":"image-file","path":"a.png","format":"png",` +
				`"sha256":"a","size":1,"boxes":{"b":{"_type":"boxes2D"}}}`,
			wantErr: `field "media": boxes "b": media has no path or artifact path`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ces := newFakeCES(t)
			uploader, _ := newTestUploader(t, ces, projectFoundJSON)

			_, err := uploader.Upload(context.Background(), &spb.UploadEvalTableRequest{
				DatasetFields: []*spb.EvalTableField{{
					Source:         "output",
					Name:           "media",
					ValueType:      "json",
					WandbMediaType: tc.mediaType,
				}},
				RowsJson: []string{`{"output":{"media":` + tc.media + `},"scores":{}}`},
			})

			assert.ErrorContains(t, err, tc.wantErr)
			assert.Empty(t, ces.Requests())
		})
	}
}
