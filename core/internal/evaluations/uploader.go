package evaluations

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"net/url"
	"os"
	"sync"

	"github.com/Khan/genqlient/graphql"

	ces "github.com/wandb/cwevals-go"
	"github.com/wandb/cwevals-go/option"

	"github.com/wandb/wandb/core/internal/api"
	"github.com/wandb/wandb/core/internal/gql"
	"github.com/wandb/wandb/core/internal/observability"
	"github.com/wandb/wandb/core/internal/runbranch"
	"github.com/wandb/wandb/core/internal/runhandle"
	"github.com/wandb/wandb/core/internal/settings"
	spb "github.com/wandb/wandb/core/pkg/service_go_proto"
)

const (
	// targetRowBatchBytes is the size above which a row append request body
	// is split, leaving headroom below CES's 16 MiB request limit.
	targetRowBatchBytes = 8 << 20

	// maxRowsPerBatch is the most rows in one row append request.
	maxRowsPerBatch = 10_000

	// idempotencyKeyVersion versions the idempotency key format.
	//
	// Bump it if the row batching changes, so that a retry of a partially
	// uploaded table can't reuse keys for differently batched rows.
	idempotencyKeyVersion = "v2"
)

// Uploader uploads a run's EvalTables to CES.
type Uploader struct {
	graphqlClient graphql.Client
	runHandle     *runhandle.RunHandle

	// evalTables returns the CES client, creating it on first use.
	evalTables func() (*ces.EvalTableService, error)

	// scopeMu guards scopeID.
	scopeMu sync.Mutex

	// scopeID is the run's project's internal ID, or empty if not yet known.
	scopeID string
}

type UploaderParams struct {
	BaseURL            api.WBBaseURL
	CredentialProvider api.CredentialProvider
	GraphqlClient      graphql.Client
	Logger             *observability.CoreLogger
	RunHandle          *runhandle.RunHandle
	Settings           *settings.Settings
}

func NewUploader(params UploaderParams) *Uploader {
	return &Uploader{
		graphqlClient: params.GraphqlClient,
		runHandle:     params.RunHandle,
		evalTables: sync.OnceValues(func() (*ces.EvalTableService, error) {
			cesBaseURL, err := CESBaseURL((*url.URL)(params.BaseURL))
			if err != nil {
				return nil, fmt.Errorf("evaluations: invalid CES URL: %v", err)
			}

			httpClient := NewCESHTTPClient(
				cesBaseURL,
				params.CredentialProvider,
				params.Logger.Logger,
				params.Settings,
			)
			return NewEvalTableService(cesBaseURL, httpClient), nil
		}),
	}
}

// Upload creates an EvalTable in CES, adds its columns and rows,
// and cuts a version.
//
// It must be called after the run is initialized. It deletes the request's
// rows file, if any.
func (u *Uploader) Upload(
	ctx context.Context,
	request *spb.LogEvalTableRequest,
) (*spb.LogEvalTableResponse, error) {
	if path := request.GetRowsFilePath(); path != "" {
		defer func() { _ = os.Remove(path) }()
	}

	evalTables, err := u.evalTables()
	if err != nil {
		return nil, err
	}

	upserter, err := u.runHandle.Upserter()
	if err != nil {
		return nil, fmt.Errorf("evaluations: %v", err)
	}
	runPath := upserter.RunPath()

	scopeID, err := u.projectScopeID(ctx, runPath)
	if err != nil {
		return nil, err
	}
	scope := ces.F(scopeID)
	keys := newIdempotencyKeys(runPath, request.GetHistoryKey(), request.GetStep())

	created, err := evalTables.New(
		ctx,
		ces.EvalTableNewParams{ScopeID: scope, Name: ces.F(request.GetName())},
		keys.option("create"),
	)
	if err != nil {
		return nil, fmt.Errorf("evaluations: failed to create EvalTable: %v", err)
	}

	_, err = evalTables.Columns.New(
		ctx,
		created.EvaluationID,
		ces.EvalTableColumnNewParams{
			ScopeID:       scope,
			DatasetFields: ces.F(datasetFieldParams(request.GetDatasetFields())),
			Scorers:       ces.F(scorerParams(request.GetScorers())),
		},
		keys.option("columns"),
	)
	if err != nil {
		return nil, fmt.Errorf("evaluations: failed to add columns: %v", err)
	}

	batchIndex := 0
	for body, err := range rowBatchBodies(requestRows(request)) {
		if err != nil {
			return nil, fmt.Errorf("evaluations: failed to read rows: %v", err)
		}

		if _, err := evalTables.Rows.Add(
			ctx,
			created.EvaluationID,
			ces.EvalTableRowAddParams{ScopeID: scope},
			keys.option(fmt.Sprintf("rows-%d", batchIndex)),
			// The rows are already JSON, so send them without decoding them
			// into the SDK's types.
			option.WithRequestBody("application/json", body),
		); err != nil {
			return nil, fmt.Errorf("evaluations: failed to add rows: %v", err)
		}
		batchIndex++
	}

	version, err := evalTables.Versions.New(
		ctx,
		created.EvaluationID,
		ces.EvalTableVersionNewParams{ScopeID: scope},
		keys.option("version"),
	)
	if err != nil {
		return nil, fmt.Errorf("evaluations: failed to cut version: %v", err)
	}

	return &spb.LogEvalTableResponse{
		EvaluationId:        created.EvaluationID,
		EvaluationVersionId: version.EvaluationVersionID,
		DatasetId:           created.DatasetID,
		DatasetVersionId:    version.DatasetVersionID,
	}, nil
}

// projectScopeID returns the internal ID of the run's project.
func (u *Uploader) projectScopeID(
	ctx context.Context,
	runPath runbranch.RunPath,
) (string, error) {
	u.scopeMu.Lock()
	scopeID := u.scopeID
	u.scopeMu.Unlock()

	if scopeID != "" {
		return scopeID, nil
	}

	// Concurrent uploads may look up the ID more than once; that's harmless.
	response, err := gql.ProjectInternalID(
		ctx,
		u.graphqlClient,
		runPath.Entity,
		runPath.Project,
	)
	if err != nil {
		return "", fmt.Errorf("evaluations: failed to look up project: %v", err)
	}

	project := response.GetProject()
	if project == nil || project.GetInternalId() == "" {
		return "", fmt.Errorf(
			"evaluations: project %s/%s not found",
			runPath.Entity,
			runPath.Project,
		)
	}

	scopeID = project.GetInternalId()

	u.scopeMu.Lock()
	u.scopeID = scopeID
	u.scopeMu.Unlock()

	return scopeID, nil
}

func datasetFieldParams(
	fields []*spb.EvalTableField,
) []ces.EvalTableColumnNewParamsDatasetField {
	params := make([]ces.EvalTableColumnNewParamsDatasetField, len(fields))

	for i, field := range fields {
		params[i] = ces.EvalTableColumnNewParamsDatasetField{
			Name: ces.F(field.GetName()),
			Source: ces.F(
				ces.EvalTableColumnNewParamsDatasetFieldsSource(field.GetSource())),
			ValueType: ces.F(
				ces.EvalTableColumnNewParamsDatasetFieldsValueType(field.GetValueType())),
		}

		if field.GetExtensionType() != "" {
			params[i].ExtensionType = ces.F(
				ces.EvalTableColumnNewParamsDatasetFieldsExtensionType(
					field.GetExtensionType()))
			params[i].ExtensionSchemaVersion = ces.F(
				int64(field.GetExtensionSchemaVersion()))
		}
	}

	return params
}

func scorerParams(scorers []*spb.EvalTableScorer) []ces.EvalTableColumnNewParamsScorer {
	params := make([]ces.EvalTableColumnNewParamsScorer, len(scorers))

	for i, scorer := range scorers {
		params[i] = ces.EvalTableColumnNewParamsScorer{
			Name: ces.F(scorer.GetName()),
			ValueType: ces.F(
				ces.EvalTableColumnNewParamsScorersValueType(scorer.GetValueType())),
		}
	}

	return params
}

// requestRows returns the request's JSON-encoded rows.
func requestRows(request *spb.LogEvalTableRequest) iter.Seq2[[]byte, error] {
	if path := request.GetRowsFilePath(); path != "" {
		return rowsFromFile(path)
	}

	return func(yield func([]byte, error) bool) {
		for _, row := range request.GetRowsJson() {
			if !yield([]byte(row), nil) {
				return
			}
		}
	}
}

// rowsFromFile returns the rows in a gzipped file with one row per line.
func rowsFromFile(path string) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		file, err := os.Open(path)
		if err != nil {
			yield(nil, err)
			return
		}
		defer func() { _ = file.Close() }()

		gzReader, err := gzip.NewReader(file)
		if err != nil {
			yield(nil, err)
			return
		}
		defer func() { _ = gzReader.Close() }()

		// Not a bufio.Scanner, which limits line lengths.
		reader := bufio.NewReader(gzReader)
		for {
			line, err := reader.ReadBytes('\n')
			line = bytes.TrimSuffix(line, []byte{'\n'})

			if len(line) > 0 && !yield(line, nil) {
				return
			}

			switch {
			case errors.Is(err, io.EOF):
				return
			case err != nil:
				yield(nil, err)
				return
			}
		}
	}
}

// rowBatchBodies groups JSON-encoded rows into row append request bodies.
//
// Each body has at most maxRowsPerBatch rows and is at most
// targetRowBatchBytes, unless it has a single larger row.
// It stops after the first error from rows.
func rowBatchBodies(rows iter.Seq2[[]byte, error]) iter.Seq2[[]byte, error] {
	const prefix, suffix = `{"rows":[`, `]}`

	return func(yield func([]byte, error) bool) {
		var body bytes.Buffer
		nrows := 0

		flush := func() bool {
			body.WriteString(suffix)
			batch := bytes.Clone(body.Bytes())
			body.Reset()
			nrows = 0
			return yield(batch, nil)
		}

		for row, err := range rows {
			if err != nil {
				yield(nil, err)
				return
			}

			if nrows > 0 &&
				(nrows >= maxRowsPerBatch ||
					body.Len()+1+len(row)+len(suffix) > targetRowBatchBytes) {
				if !flush() {
					return
				}
			}

			if nrows == 0 {
				body.WriteString(prefix)
			} else {
				body.WriteByte(',')
			}
			body.Write(row)
			nrows++
		}

		if nrows > 0 {
			flush()
		}
	}
}

// idempotencyKeys derives CES idempotency keys for one EvalTable upload.
//
// The keys depend only on where the table is logged, so that retrying the
// upload at the same location reuses the CES resources it already created.
// CES replays a request with the same key and body, and rejects one with
// the same key but a different body.
type idempotencyKeys struct {
	prefix string
}

func newIdempotencyKeys(
	runPath runbranch.RunPath,
	historyKey string,
	step string,
) idempotencyKeys {
	location, err := json.Marshal([]string{
		runPath.Entity,
		runPath.Project,
		runPath.RunID,
		historyKey,
		step,
	})
	if err != nil {
		// Marshaling a slice of strings cannot fail.
		panic(errors.Join(errors.New("evaluations: unexpected error"), err))
	}
	digest := sha256.Sum256(location)

	return idempotencyKeys{
		prefix: fmt.Sprintf(
			"wandb-eval-table-%s-%s",
			idempotencyKeyVersion,
			hex.EncodeToString(digest[:]),
		),
	}
}

// option returns a request option setting the key for an operation.
func (k idempotencyKeys) option(operation string) option.RequestOption {
	return option.WithHeader("Idempotency-Key", k.prefix+"-"+operation)
}
