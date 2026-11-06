package evaltables

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"sync"

	"github.com/Khan/genqlient/graphql"
	"github.com/google/uuid"

	ces "github.com/wandb/cwevals-go"
	"github.com/wandb/cwevals-go/option"

	"github.com/wandb/wandb/core/internal/api"
	"github.com/wandb/wandb/core/internal/gql"
	"github.com/wandb/wandb/core/internal/observability"
	"github.com/wandb/wandb/core/internal/runbranch"
	"github.com/wandb/wandb/core/internal/runhandle"
	"github.com/wandb/wandb/core/internal/settings"
	"github.com/wandb/wandb/core/internal/wboperation"
	spb "github.com/wandb/wandb/core/pkg/service_go_proto"
)

const (
	// maxRequestBodyBytes is CES's request body limit.
	maxRequestBodyBytes = 16 << 20

	// targetRowBatchBytes is the size above which a row append request body
	// is split, so we leave headroom below CES's 16 MiB request limit.
	targetRowBatchBytes = 8 << 20

	maxRowsPerBatch = 10_000

	idempotencyKeyVersion = "v3"
)

type Uploader struct {
	graphqlClient graphql.Client
	runHandle     *runhandle.RunHandle

	evalTables    *ces.EvalTableService
	evalTablesErr error

	// scopeMu guards scopeID.
	scopeMu sync.Mutex

	// scopeID is the run's project's internal ID, or empty if not yet fetched.
	scopeID string
}

type UploaderParams struct {
	CredentialProvider api.CredentialProvider
	GraphqlClient      graphql.Client
	Logger             *observability.CoreLogger
	RunHandle          *runhandle.RunHandle
	Settings           *settings.Settings
}

func NewUploader(params UploaderParams) *Uploader {
	u := &Uploader{
		graphqlClient: params.GraphqlClient,
		runHandle:     params.RunHandle,
	}
	u.evalTables, u.evalTablesErr = NewEvalTableService(
		params.CredentialProvider,
		params.Logger.Logger,
		params.Settings,
	)
	return u
}

// Upload creates an EvalTable in CES, adds its columns and rows,
// and cuts a version.
//
// It deletes the request's rows file, if any.
func (u *Uploader) Upload(
	ctx context.Context,
	request *spb.UploadEvalTableRequest,
) (*spb.UploadEvalTableResponse, error) {
	if path := request.GetRowsFilePath(); path != "" {
		defer func() { _ = os.Remove(path) }()
	}

	evalTables := u.evalTables
	if evalTables == nil {
		return nil, u.evalTablesErr
	}

	columns := ces.EvalTableColumnNewParams{
		DatasetFields: ces.F(datasetFieldParams(request.GetDatasetFields())),
		Scorers:       ces.F(scorerParams(request.GetScorers())),
	}
	if err := validateColumns(columns); err != nil {
		return nil, err
	}

	upserter, err := u.runHandle.Upserter()
	if err != nil {
		return nil, fmt.Errorf("evaltables: %v", err)
	}
	runPath := upserter.RunPath()

	scopeID, err := u.projectScopeID(ctx, runPath)
	if err != nil {
		return nil, err
	}
	scope := ces.F(scopeID)
	columns.ScopeID = scope
	keys := newIdempotencyKeys()

	created, err := evalTables.New(
		ctx,
		ces.EvalTableNewParams{ScopeID: scope, Name: ces.F(request.GetName())},
		keys.option("create"),
	)
	if err != nil {
		return nil, fmt.Errorf("evaltables: failed to create EvalTable: %v", err)
	}

	_, err = evalTables.Columns.New(
		ctx,
		created.EvaluationID,
		columns,
		keys.option("columns"),
	)
	if err != nil {
		return nil, fmt.Errorf("evaltables: failed to add columns: %v", err)
	}

	// NewProgress only fails if the operation already has progress.
	progress, _ := wboperation.Get(ctx).NewProgress()
	totalRows := len(request.GetRowsJson())
	rowsDone := 0
	batchIndex := 0
	for batch, err := range rowBatches(requestRows(request)) {
		if err != nil {
			return nil, fmt.Errorf("evaltables: failed to read rows: %v", err)
		}

		if _, err := evalTables.Rows.Add(
			ctx,
			created.EvaluationID,
			ces.EvalTableRowAddParams{ScopeID: scope},
			keys.option(fmt.Sprintf("rows-%d", batchIndex)),
			// The rows are already JSON, so send them without decoding them
			// into the CES SDK's types.
			option.WithRequestBody("application/json", batch.body),
		); err != nil {
			return nil, fmt.Errorf("evaltables: failed to add rows: %v", err)
		}
		batchIndex++

		rowsDone += batch.rows
		if totalRows > 0 {
			progress.SetUnitsOfTotal(rowsDone, totalRows, "rows")
		} else {
			// A rows file's row count isn't known until it's read.
			progress.SetUnitsDone(rowsDone, "rows")
		}
	}

	version, err := evalTables.Versions.New(
		ctx,
		created.EvaluationID,
		ces.EvalTableVersionNewParams{ScopeID: scope},
		keys.option("version"),
	)
	if err != nil {
		return nil, fmt.Errorf("evaltables: failed to cut version: %v", err)
	}

	return &spb.UploadEvalTableResponse{
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
		return "", fmt.Errorf("evaltables: failed to look up project: %v", err)
	}

	project := response.GetProject()
	if project == nil || project.GetInternalId() == "" {
		return "", fmt.Errorf(
			"evaltables: project %s/%s not found",
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
func requestRows(request *spb.UploadEvalTableRequest) iter.Seq2[[]byte, error] {
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

// validateColumns rejects a columns request that CES would, before the
// EvalTable is created.
//
// The CES SDK's types define the valid values, so updating the CES SDK
// updates the check.
func validateColumns(columns ces.EvalTableColumnNewParams) error {
	for _, field := range columns.DatasetFields.Value {
		switch {
		case !field.Source.Value.IsKnown():
			return fmt.Errorf(
				"evaltables: field %q has unknown source %q",
				field.Name.Value, field.Source.Value)
		case !field.ValueType.Value.IsKnown():
			return fmt.Errorf(
				"evaltables: field %q has unknown value type %q",
				field.Name.Value, field.ValueType.Value)
		case field.ExtensionType.Present && !field.ExtensionType.Value.IsKnown():
			return fmt.Errorf(
				"evaltables: field %q has unknown extension type %q",
				field.Name.Value, field.ExtensionType.Value)
		}
	}

	for _, scorer := range columns.Scorers.Value {
		if !scorer.ValueType.Value.IsKnown() {
			return fmt.Errorf(
				"evaltables: scorer %q has unknown value type %q",
				scorer.Name.Value, scorer.ValueType.Value)
		}
	}

	body, err := columns.MarshalJSON()
	if err != nil {
		return fmt.Errorf("evaltables: failed to encode columns: %v", err)
	}
	if len(body) >= maxRequestBodyBytes {
		return fmt.Errorf(
			"evaltables: columns request is %d bytes, must be under %d",
			len(body), maxRequestBodyBytes)
	}

	return nil
}

type rowBatch struct {
	body []byte
	rows int
}

// rowBatches groups JSON-encoded rows into row append request bodies.
//
// Each body has at most maxRowsPerBatch rows and is at most
// targetRowBatchBytes, unless it has a single larger row.
// It stops after the first error from rows.
func rowBatches(rows iter.Seq2[[]byte, error]) iter.Seq2[rowBatch, error] {
	const prefix, suffix = `{"rows":[`, `]}`

	return func(yield func(rowBatch, error) bool) {
		var body bytes.Buffer
		nrows := 0

		flush := func() bool {
			body.WriteString(suffix)
			batch := rowBatch{body: bytes.Clone(body.Bytes()), rows: nrows}
			body.Reset()
			nrows = 0
			return yield(batch, nil)
		}

		for row, err := range rows {
			if err != nil {
				yield(rowBatch{}, err)
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

type idempotencyKeys struct {
	prefix string
}

// newIdempotencyKeys derives CES idempotency keys for one EvalTable upload.
//
// The keys share a random ID, so that retries within the upload reuse the
// CES resources it already created, and separate uploads never collide.
// Uploads can't key on where the table is logged: in shared mode, writers
// to the same run can log the same key at the same step.
//
// CES replays a request with the same key and body, and rejects one with
// the same key but a different body.
func newIdempotencyKeys() idempotencyKeys {
	return idempotencyKeys{
		prefix: fmt.Sprintf(
			"wandb-eval-table-%s-%s",
			idempotencyKeyVersion,
			uuid.Must(uuid.NewV7()),
		),
	}
}

// option returns a request option setting the key for an operation.
func (k idempotencyKeys) option(operation string) option.RequestOption {
	return option.WithHeader("Idempotency-Key", k.prefix+"-"+operation)
}
