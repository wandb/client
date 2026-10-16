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

// EvalTableVersionRowService contains methods and other services that help with
// interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvalTableVersionRowService] method instead.
type EvalTableVersionRowService struct {
	Options []option.RequestOption
}

// NewEvalTableVersionRowService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewEvalTableVersionRowService(opts ...option.RequestOption) (r *EvalTableVersionRowService) {
	r = &EvalTableVersionRowService{}
	r.Options = opts
	return
}

// Reads Cases and selected Scores at the EvaluationVersion cutoff. Use `i`, `o`,
// and `s` for input, output, and score columns; `input`, `output`, and `score`
// remain aliases. `subject` is not available. Filtering and ordering apply before
// pagination.
//
// The encoded request target (path and query string) must be at most 8000 bytes.
// Larger requests return 414 URI Too Long with error code `request_uri_too_long`;
// retry the same `filter`, `order_by`, `limit`, and `cursor` through the
// corresponding POST `/rows/search` operation.
func (r *EvalTableVersionRowService) List(
	ctx context.Context,
	evaluationVersionID string,
	params EvalTableVersionRowListParams,
	opts ...option.RequestOption,
) (res *EvalTableVersionRowListResponse, err error) {
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
		"v1/namespaces/%s/scopes/%s/eval-table-versions/%s/rows",
		params.Namespace,
		params.ScopeID,
		evaluationVersionID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Reads Cases and selected Scores at the EvaluationVersion cutoff. `filter` and
// `order_by` use the same AIP-style syntax as the GET route. Filtering and
// ordering apply before pagination.
//
// Request bodies must be smaller than 1 MiB; larger bodies are rejected with 413
// Request Entity Too Large.
func (r *EvalTableVersionRowService) Search(
	ctx context.Context,
	evaluationVersionID string,
	params EvalTableVersionRowSearchParams,
	opts ...option.RequestOption,
) (res *EvalTableVersionRowSearchResponse, err error) {
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
		"v1/namespaces/%s/scopes/%s/eval-table-versions/%s/rows/search",
		params.Namespace,
		params.ScopeID,
		evaluationVersionID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type EvalTableVersionRowListResponse struct {
	Rows       []EvalTableVersionRowListResponseRow `json:"rows"        api:"required"`
	NextCursor string                               `json:"next_cursor"`
	JSON       evalTableVersionRowListResponseJSON  `json:"-"`
}

// evalTableVersionRowListResponseJSON contains the JSON metadata for the struct
// [EvalTableVersionRowListResponse]
type evalTableVersionRowListResponseJSON struct {
	Rows        apijson.Field
	NextCursor  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvalTableVersionRowListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evalTableVersionRowListResponseJSON) RawJSON() string {
	return r.raw
}

type EvalTableVersionRowListResponseRow struct {
	Attributes map[string]interface{}                     `json:"attributes" api:"required"`
	CaseID     string                                     `json:"case_id"    api:"required" format:"uuid"`
	Input      interface{}                                `json:"input"      api:"required"`
	Output     interface{}                                `json:"output"     api:"required"`
	Scores     []EvalTableVersionRowListResponseRowsScore `json:"scores"     api:"required"`
	JSON       evalTableVersionRowListResponseRowJSON     `json:"-"`
}

// evalTableVersionRowListResponseRowJSON contains the JSON metadata for the struct
// [EvalTableVersionRowListResponseRow]
type evalTableVersionRowListResponseRowJSON struct {
	Attributes  apijson.Field
	CaseID      apijson.Field
	Input       apijson.Field
	Output      apijson.Field
	Scores      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvalTableVersionRowListResponseRow) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evalTableVersionRowListResponseRowJSON) RawJSON() string {
	return r.raw
}

type EvalTableVersionRowListResponseRowsScore struct {
	ID              string                                       `json:"id"                api:"required" format:"uuid"`
	ScorerID        string                                       `json:"scorer_id"         api:"required" format:"uuid"`
	ScorerVersionID string                                       `json:"scorer_version_id" api:"required" format:"uuid"`
	Value           interface{}                                  `json:"value"             api:"required"`
	JSON            evalTableVersionRowListResponseRowsScoreJSON `json:"-"`
}

// evalTableVersionRowListResponseRowsScoreJSON contains the JSON metadata for the
// struct [EvalTableVersionRowListResponseRowsScore]
type evalTableVersionRowListResponseRowsScoreJSON struct {
	ID              apijson.Field
	ScorerID        apijson.Field
	ScorerVersionID apijson.Field
	Value           apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *EvalTableVersionRowListResponseRowsScore) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evalTableVersionRowListResponseRowsScoreJSON) RawJSON() string {
	return r.raw
}

type EvalTableVersionRowSearchResponse struct {
	Rows       []EvalTableVersionRowSearchResponseRow `json:"rows"        api:"required"`
	NextCursor string                                 `json:"next_cursor"`
	JSON       evalTableVersionRowSearchResponseJSON  `json:"-"`
}

// evalTableVersionRowSearchResponseJSON contains the JSON metadata for the struct
// [EvalTableVersionRowSearchResponse]
type evalTableVersionRowSearchResponseJSON struct {
	Rows        apijson.Field
	NextCursor  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvalTableVersionRowSearchResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evalTableVersionRowSearchResponseJSON) RawJSON() string {
	return r.raw
}

type EvalTableVersionRowSearchResponseRow struct {
	Attributes map[string]interface{}                       `json:"attributes" api:"required"`
	CaseID     string                                       `json:"case_id"    api:"required" format:"uuid"`
	Input      interface{}                                  `json:"input"      api:"required"`
	Output     interface{}                                  `json:"output"     api:"required"`
	Scores     []EvalTableVersionRowSearchResponseRowsScore `json:"scores"     api:"required"`
	JSON       evalTableVersionRowSearchResponseRowJSON     `json:"-"`
}

// evalTableVersionRowSearchResponseRowJSON contains the JSON metadata for the
// struct [EvalTableVersionRowSearchResponseRow]
type evalTableVersionRowSearchResponseRowJSON struct {
	Attributes  apijson.Field
	CaseID      apijson.Field
	Input       apijson.Field
	Output      apijson.Field
	Scores      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvalTableVersionRowSearchResponseRow) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evalTableVersionRowSearchResponseRowJSON) RawJSON() string {
	return r.raw
}

type EvalTableVersionRowSearchResponseRowsScore struct {
	ID              string                                         `json:"id"                api:"required" format:"uuid"`
	ScorerID        string                                         `json:"scorer_id"         api:"required" format:"uuid"`
	ScorerVersionID string                                         `json:"scorer_version_id" api:"required" format:"uuid"`
	Value           interface{}                                    `json:"value"             api:"required"`
	JSON            evalTableVersionRowSearchResponseRowsScoreJSON `json:"-"`
}

// evalTableVersionRowSearchResponseRowsScoreJSON contains the JSON metadata for
// the struct [EvalTableVersionRowSearchResponseRowsScore]
type evalTableVersionRowSearchResponseRowsScoreJSON struct {
	ID              apijson.Field
	ScorerID        apijson.Field
	ScorerVersionID apijson.Field
	Value           apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *EvalTableVersionRowSearchResponseRowsScore) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evalTableVersionRowSearchResponseRowsScoreJSON) RawJSON() string {
	return r.raw
}

type EvalTableVersionRowListParams struct {
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
	// Continuation from the same version, filter, and order_by; page size may change.
	Cursor param.Field[string] `query:"cursor"`
	// AIP-160-style comparisons joined by explicit AND; NOT and - negate. Use i, o,
	// and s for input, output, and score; the long prefixes are aliases. _value_
	// matches a substring. OR, parentheses, implicit conjunctions, one-sided
	// wildcards, nested fields, functions, :, column comparisons, and subject are
	// unsupported.
	Filter param.Field[string] `query:"filter"`
	Limit  param.Field[int64]  `query:"limit"`
	// One i, o, or s column followed by optional desc. The long prefixes are aliases.
	// Ascending is the default.
	OrderBy param.Field[string] `query:"order_by"`
}

// URLQuery serializes [EvalTableVersionRowListParams]'s query parameters as
// `url.Values`.
func (r EvalTableVersionRowListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type EvalTableVersionRowSearchParams struct {
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
	// Continuation for the same EvalTableVersion, filter, and order_by.
	Cursor param.Field[string] `json:"cursor"`
	// Same filter expression syntax as the GET filter parameter.
	Filter param.Field[string] `json:"filter"`
	Limit  param.Field[int64]  `json:"limit"`
	// Same ordering expression syntax as GET order_by. Omit for Dataset order.
	OrderBy param.Field[string] `json:"order_by"`
}

func (r EvalTableVersionRowSearchParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
