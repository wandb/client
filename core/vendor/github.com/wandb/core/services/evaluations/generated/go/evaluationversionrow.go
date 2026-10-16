// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package coreweaveevaluations

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/wandb/core/services/evaluations/generated/go/internal/apijson"
	"github.com/wandb/core/services/evaluations/generated/go/internal/apiquery"
	"github.com/wandb/core/services/evaluations/generated/go/internal/param"
	"github.com/wandb/core/services/evaluations/generated/go/internal/requestconfig"
	"github.com/wandb/core/services/evaluations/generated/go/option"
)

// EvaluationVersionRowService contains methods and other services that help with
// interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvaluationVersionRowService] method instead.
type EvaluationVersionRowService struct {
	Options []option.RequestOption
}

// NewEvaluationVersionRowService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewEvaluationVersionRowService(opts ...option.RequestOption) (r *EvaluationVersionRowService) {
	r = &EvaluationVersionRowService{}
	r.Options = opts
	return
}

// Reads Cases and selected Scores at the EvaluationVersion cutoff. `filter` uses
// an AIP-160-style subset: comparisons, unary `NOT` or `-`, explicit `AND`, and
// `*value*` substring matching; `OR`, parentheses, implicit conjunctions, and `:`
// are rejected. `order_by` uses AIP-132-style syntax with one column and optional
// `desc`. Use bare `subject` for the scalar subject, or `i.<name>`, `o.<name>`,
// and `s.<name>` for input, output, and score columns; the longer `input`,
// `output`, and `score` prefixes are accepted as aliases. Quote and JSON-escape
// non-identifier names. Filtering and ordering apply before pagination.
//
// The encoded request target (path and query string) must be at most 8000 bytes.
// Larger requests return 414 URI Too Long with error code `request_uri_too_long`;
// retry the same `filter`, `order_by`, `limit`, and `cursor` through the
// corresponding POST `/rows/search` operation.
func (r *EvaluationVersionRowService) List(
	ctx context.Context,
	evaluationVersionID string,
	params EvaluationVersionRowListParams,
	opts ...option.RequestOption,
) (res *EvaluationVersionRowListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	precfg, err := requestconfig.PreRequestOptions(opts...)
	if err != nil {
		return nil, err
	}
	requestconfig.UseDefaultParam(&params.Namespace, precfg.Namespace)
	requestconfig.UseDefaultParam(&params.ScopeID, precfg.ScopeID)
	if params.Namespace.Value == "" {
		err = errors.New("missing required namespace parameter")
		return nil, err
	}
	if params.ScopeID.Value == "" {
		err = errors.New("missing required scope_id parameter")
		return nil, err
	}
	if evaluationVersionID == "" {
		err = errors.New("missing required evaluation_version_id parameter")
		return nil, err
	}
	path := fmt.Sprintf(
		"v1/namespaces/%s/scopes/%s/evaluation-versions/%s/rows",
		params.Namespace,
		params.ScopeID,
		evaluationVersionID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Reads Cases and selected Scores at the EvaluationVersion cutoff. `filter` and
// `order_by` use the same AIP-style syntax as the GET route. Filtering and
// ordering apply before pagination. This operation does not create a saved query.
//
// Request bodies must be smaller than 1 MiB; larger bodies are rejected with 413
// Request Entity Too Large.
func (r *EvaluationVersionRowService) Search(
	ctx context.Context,
	evaluationVersionID string,
	params EvaluationVersionRowSearchParams,
	opts ...option.RequestOption,
) (res *EvaluationVersionRowSearchResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	precfg, err := requestconfig.PreRequestOptions(opts...)
	if err != nil {
		return nil, err
	}
	requestconfig.UseDefaultParam(&params.Namespace, precfg.Namespace)
	requestconfig.UseDefaultParam(&params.ScopeID, precfg.ScopeID)
	if params.Namespace.Value == "" {
		err = errors.New("missing required namespace parameter")
		return nil, err
	}
	if params.ScopeID.Value == "" {
		err = errors.New("missing required scope_id parameter")
		return nil, err
	}
	if evaluationVersionID == "" {
		err = errors.New("missing required evaluation_version_id parameter")
		return nil, err
	}
	path := fmt.Sprintf(
		"v1/namespaces/%s/scopes/%s/evaluation-versions/%s/rows/search",
		params.Namespace,
		params.ScopeID,
		evaluationVersionID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type EvaluationVersionRowListResponse struct {
	Rows       []EvaluationVersionRowListResponseRow `json:"rows"        api:"required"`
	NextCursor string                                `json:"next_cursor"`
	JSON       evaluationVersionRowListResponseJSON  `json:"-"`
}

// evaluationVersionRowListResponseJSON contains the JSON metadata for the struct
// [EvaluationVersionRowListResponse]
type evaluationVersionRowListResponseJSON struct {
	Rows        apijson.Field
	NextCursor  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationVersionRowListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionRowListResponseJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionRowListResponseRow struct {
	Attributes map[string]interface{}                      `json:"attributes" api:"required"`
	CaseID     string                                      `json:"case_id"    api:"required"          format:"uuid"`
	Input      interface{}                                 `json:"input"      api:"required"`
	Output     interface{}                                 `json:"output"     api:"required"`
	Scores     []EvaluationVersionRowListResponseRowsScore `json:"scores"     api:"required"`
	Subject    string                                      `json:"subject"    api:"required,nullable"`
	JSON       evaluationVersionRowListResponseRowJSON     `json:"-"`
}

// evaluationVersionRowListResponseRowJSON contains the JSON metadata for the
// struct [EvaluationVersionRowListResponseRow]
type evaluationVersionRowListResponseRowJSON struct {
	Attributes  apijson.Field
	CaseID      apijson.Field
	Input       apijson.Field
	Output      apijson.Field
	Scores      apijson.Field
	Subject     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationVersionRowListResponseRow) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionRowListResponseRowJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionRowListResponseRowsScore struct {
	ScoreID         string                                        `json:"score_id"          api:"required" format:"uuid"`
	ScorerID        string                                        `json:"scorer_id"         api:"required" format:"uuid"`
	ScorerVersionID string                                        `json:"scorer_version_id" api:"required" format:"uuid"`
	Value           interface{}                                   `json:"value"             api:"required"`
	JSON            evaluationVersionRowListResponseRowsScoreJSON `json:"-"`
}

// evaluationVersionRowListResponseRowsScoreJSON contains the JSON metadata for the
// struct [EvaluationVersionRowListResponseRowsScore]
type evaluationVersionRowListResponseRowsScoreJSON struct {
	ScoreID         apijson.Field
	ScorerID        apijson.Field
	ScorerVersionID apijson.Field
	Value           apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *EvaluationVersionRowListResponseRowsScore) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionRowListResponseRowsScoreJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionRowSearchResponse struct {
	Rows       []EvaluationVersionRowSearchResponseRow `json:"rows"        api:"required"`
	NextCursor string                                  `json:"next_cursor"`
	JSON       evaluationVersionRowSearchResponseJSON  `json:"-"`
}

// evaluationVersionRowSearchResponseJSON contains the JSON metadata for the struct
// [EvaluationVersionRowSearchResponse]
type evaluationVersionRowSearchResponseJSON struct {
	Rows        apijson.Field
	NextCursor  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationVersionRowSearchResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionRowSearchResponseJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionRowSearchResponseRow struct {
	Attributes map[string]interface{}                        `json:"attributes" api:"required"`
	CaseID     string                                        `json:"case_id"    api:"required"          format:"uuid"`
	Input      interface{}                                   `json:"input"      api:"required"`
	Output     interface{}                                   `json:"output"     api:"required"`
	Scores     []EvaluationVersionRowSearchResponseRowsScore `json:"scores"     api:"required"`
	Subject    string                                        `json:"subject"    api:"required,nullable"`
	JSON       evaluationVersionRowSearchResponseRowJSON     `json:"-"`
}

// evaluationVersionRowSearchResponseRowJSON contains the JSON metadata for the
// struct [EvaluationVersionRowSearchResponseRow]
type evaluationVersionRowSearchResponseRowJSON struct {
	Attributes  apijson.Field
	CaseID      apijson.Field
	Input       apijson.Field
	Output      apijson.Field
	Scores      apijson.Field
	Subject     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationVersionRowSearchResponseRow) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionRowSearchResponseRowJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionRowSearchResponseRowsScore struct {
	ScoreID         string                                          `json:"score_id"          api:"required" format:"uuid"`
	ScorerID        string                                          `json:"scorer_id"         api:"required" format:"uuid"`
	ScorerVersionID string                                          `json:"scorer_version_id" api:"required" format:"uuid"`
	Value           interface{}                                     `json:"value"             api:"required"`
	JSON            evaluationVersionRowSearchResponseRowsScoreJSON `json:"-"`
}

// evaluationVersionRowSearchResponseRowsScoreJSON contains the JSON metadata for
// the struct [EvaluationVersionRowSearchResponseRowsScore]
type evaluationVersionRowSearchResponseRowsScoreJSON struct {
	ScoreID         apijson.Field
	ScorerID        apijson.Field
	ScorerVersionID apijson.Field
	Value           apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *EvaluationVersionRowSearchResponseRowsScore) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionRowSearchResponseRowsScoreJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionRowListParams struct {
	// Namespace that assigned the scope reference.
	//
	// Defaults to "wandb" if not set.
	Namespace param.Field[string] `path:"namespace" api:"required"`
	// Opaque reference the namespace assigned to this scope. Printable ASCII only;
	// this service compares it byte for byte and never parses it. Percent-encode it as
	// a single path segment, and prefer base64url when minting references that are not
	// already URL-safe.
	//
	// Use [option.WithScopeID] on the client to set a global default for this field.
	ScopeID param.Field[string] `path:"scope_id" api:"required"`
	// Continuation returned by the same EvaluationVersion, filter, and order_by. Page
	// size may change.
	Cursor param.Field[string] `query:"cursor"`
	// Filter expression using explicit AND, unary NOT or -, comparisons (=, !=, <,
	// <=, >, >=), and JSON scalar values. A string surrounded by \* performs a
	// case-sensitive substring match. OR, parentheses, implicit conjunctions,
	// one-sided wildcards, nested fields, functions, the : has operator, and
	// column-to-column comparisons are not supported.
	Filter param.Field[string] `query:"filter"`
	Limit  param.Field[int64]  `query:"limit"`
	// One column reference, optionally followed by desc. Ascending is the default.
	OrderBy param.Field[string] `query:"order_by"`
}

// URLQuery serializes [EvaluationVersionRowListParams]'s query parameters as
// `url.Values`.
func (r EvaluationVersionRowListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type EvaluationVersionRowSearchParams struct {
	// Namespace that assigned the scope reference.
	//
	// Defaults to "wandb" if not set.
	Namespace param.Field[string] `path:"namespace" api:"required"`
	// Opaque reference the namespace assigned to this scope. Printable ASCII only;
	// this service compares it byte for byte and never parses it. Percent-encode it as
	// a single path segment, and prefer base64url when minting references that are not
	// already URL-safe.
	//
	// Use [option.WithScopeID] on the client to set a global default for this field.
	ScopeID param.Field[string] `path:"scope_id" api:"required"`
	// Continuation for the same EvaluationVersion, filter, and order_by.
	Cursor param.Field[string] `json:"cursor"`
	// Same filter expression syntax as the GET filter parameter.
	Filter param.Field[string] `json:"filter"`
	Limit  param.Field[int64]  `json:"limit"`
	// Same ordering expression syntax as GET order_by. Omit for Dataset order.
	OrderBy param.Field[string] `json:"order_by"`
}

func (r EvaluationVersionRowSearchParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
