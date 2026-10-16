// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package coreweaveevaluations

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"slices"

	"github.com/tidwall/gjson"

	"github.com/wandb/core/services/evaluations/generated/go/internal/apijson"
	"github.com/wandb/core/services/evaluations/generated/go/internal/apiquery"
	"github.com/wandb/core/services/evaluations/generated/go/internal/param"
	"github.com/wandb/core/services/evaluations/generated/go/internal/requestconfig"
	"github.com/wandb/core/services/evaluations/generated/go/option"
	"github.com/wandb/core/services/evaluations/generated/go/shared"
)

// EvaluationVersionQueryResultService contains methods and other services that
// help with interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvaluationVersionQueryResultService] method instead.
type EvaluationVersionQueryResultService struct {
	Options []option.RequestOption
}

// NewEvaluationVersionQueryResultService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewEvaluationVersionQueryResultService(opts ...option.RequestOption) (r *EvaluationVersionQueryResultService) {
	r = &EvaluationVersionQueryResultService{}
	r.Options = opts
	return
}

// Runs one of the Evaluation's saved queries on this EvaluationVersion's Cases and
// returns aggregated Scores. Results use the Scores captured when the
// EvaluationVersion was frozen. Returns 400 when the EvaluationVersion's
// DatasetVersion lacks a grouped field or the EvaluationVersion does not select a
// metric's or grouped Score's Scorer. Results sort by the group_by columns in
// order, ascending with nulls last. For each grouping column, missing values and
// explicit nulls are treated the same. Without grouping, the result contains one
// row with an empty group array, even for an empty dataset. With grouping, an
// empty dataset returns no rows.
func (r *EvaluationVersionQueryResultService) List(
	ctx context.Context,
	evaluationVersionID string,
	evaluationQueryID string,
	params EvaluationVersionQueryResultListParams,
	opts ...option.RequestOption,
) (res *EvaluationVersionQueryResultListResponse, err error) {
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
	if evaluationQueryID == "" {
		err = errors.New("missing required evaluation_query_id parameter")
		return nil, err
	}
	path := fmt.Sprintf(
		"v1/namespaces/%s/scopes/%s/evaluation-versions/%s/queries/%s/results",
		params.Namespace,
		params.ScopeID,
		evaluationVersionID,
		evaluationQueryID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

type EvaluationVersionQueryResultListResponse struct {
	// The DatasetVersion whose Cases the results read. For a live Evaluation's latest
	// data, this is the Dataset head when the first page was requested.
	DatasetVersionID  string                                            `json:"dataset_version_id"  api:"required" format:"uuid"`
	EvaluationID      string                                            `json:"evaluation_id"       api:"required" format:"uuid"`
	EvaluationQueryID string                                            `json:"evaluation_query_id" api:"required" format:"uuid"`
	GroupBy           []EvaluationVersionQueryResultListResponseGroupBy `json:"group_by"            api:"required"`
	Metrics           []EvaluationVersionQueryResultListResponseMetric  `json:"metrics"             api:"required"`
	// Namespace that assigned the scope reference.
	Namespace EvaluationVersionQueryResultListResponseNamespace `json:"namespace" api:"required"`
	// The results include Scores published up to this commit. For a live Evaluation's
	// latest data, Cases too. Every page of one traversal uses the first page's
	// commit.
	PublishedCommitID int64                                         `json:"published_commit_id" api:"required"`
	Rows              []EvaluationVersionQueryResultListResponseRow `json:"rows"                api:"required"`
	ScopeRef          string                                        `json:"scope_ref"           api:"required"`
	// The EvaluationVersion whose data the results read. Absent when they read the
	// Evaluation's latest data.
	EvaluationVersionID string `json:"evaluation_version_id" format:"uuid"`
	// Cursor for the next page. Absent on the last page.
	NextCursor string                                       `json:"next_cursor"`
	JSON       evaluationVersionQueryResultListResponseJSON `json:"-"`
}

// evaluationVersionQueryResultListResponseJSON contains the JSON metadata for the
// struct [EvaluationVersionQueryResultListResponse]
type evaluationVersionQueryResultListResponseJSON struct {
	DatasetVersionID    apijson.Field
	EvaluationID        apijson.Field
	EvaluationQueryID   apijson.Field
	GroupBy             apijson.Field
	Metrics             apijson.Field
	Namespace           apijson.Field
	PublishedCommitID   apijson.Field
	Rows                apijson.Field
	ScopeRef            apijson.Field
	EvaluationVersionID apijson.Field
	NextCursor          apijson.Field
	raw                 string
	ExtraFields         map[string]apijson.Field
}

func (r *EvaluationVersionQueryResultListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionQueryResultListResponseJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionQueryResultListResponseGroupBy struct {
	Kind      EvaluationVersionQueryResultListResponseGroupByKind      `json:"kind"       api:"required"`
	ValueType EvaluationVersionQueryResultListResponseGroupByValueType `json:"value_type" api:"required"`
	Name      string                                                   `json:"name"`
	Source    EvaluationVersionQueryResultListResponseGroupBySource    `json:"source"`
	JSON      evaluationVersionQueryResultListResponseGroupByJSON      `json:"-"`
}

// evaluationVersionQueryResultListResponseGroupByJSON contains the JSON metadata
// for the struct [EvaluationVersionQueryResultListResponseGroupBy]
type evaluationVersionQueryResultListResponseGroupByJSON struct {
	Kind        apijson.Field
	ValueType   apijson.Field
	Name        apijson.Field
	Source      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationVersionQueryResultListResponseGroupBy) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionQueryResultListResponseGroupByJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionQueryResultListResponseGroupByKind string

const (
	EvaluationVersionQueryResultListResponseGroupByKindSubject      EvaluationVersionQueryResultListResponseGroupByKind = "subject"
	EvaluationVersionQueryResultListResponseGroupByKindDatasetField EvaluationVersionQueryResultListResponseGroupByKind = "dataset_field"
	EvaluationVersionQueryResultListResponseGroupByKindScore        EvaluationVersionQueryResultListResponseGroupByKind = "score"
)

func (r EvaluationVersionQueryResultListResponseGroupByKind) IsKnown() bool {
	switch r {
	case EvaluationVersionQueryResultListResponseGroupByKindSubject,
		EvaluationVersionQueryResultListResponseGroupByKindDatasetField,
		EvaluationVersionQueryResultListResponseGroupByKindScore:
		return true
	}
	return false
}

type EvaluationVersionQueryResultListResponseGroupByValueType string

const (
	EvaluationVersionQueryResultListResponseGroupByValueTypeBoolean EvaluationVersionQueryResultListResponseGroupByValueType = "boolean"
	EvaluationVersionQueryResultListResponseGroupByValueTypeInteger EvaluationVersionQueryResultListResponseGroupByValueType = "integer"
	EvaluationVersionQueryResultListResponseGroupByValueTypeNumber  EvaluationVersionQueryResultListResponseGroupByValueType = "number"
	EvaluationVersionQueryResultListResponseGroupByValueTypeString  EvaluationVersionQueryResultListResponseGroupByValueType = "string"
)

func (r EvaluationVersionQueryResultListResponseGroupByValueType) IsKnown() bool {
	switch r {
	case EvaluationVersionQueryResultListResponseGroupByValueTypeBoolean,
		EvaluationVersionQueryResultListResponseGroupByValueTypeInteger,
		EvaluationVersionQueryResultListResponseGroupByValueTypeNumber,
		EvaluationVersionQueryResultListResponseGroupByValueTypeString:
		return true
	}
	return false
}

type EvaluationVersionQueryResultListResponseGroupBySource string

const (
	EvaluationVersionQueryResultListResponseGroupBySourceInput  EvaluationVersionQueryResultListResponseGroupBySource = "input"
	EvaluationVersionQueryResultListResponseGroupBySourceOutput EvaluationVersionQueryResultListResponseGroupBySource = "output"
)

func (r EvaluationVersionQueryResultListResponseGroupBySource) IsKnown() bool {
	switch r {
	case EvaluationVersionQueryResultListResponseGroupBySourceInput,
		EvaluationVersionQueryResultListResponseGroupBySourceOutput:
		return true
	}
	return false
}

type EvaluationVersionQueryResultListResponseMetric struct {
	Aggregation     EvaluationVersionQueryResultListResponseMetricsAggregation `json:"aggregation"       api:"required"`
	Name            string                                                     `json:"name"              api:"required"`
	ScorerID        string                                                     `json:"scorer_id"         api:"required" format:"uuid"`
	ScorerVersionID string                                                     `json:"scorer_version_id" api:"required" format:"uuid"`
	// Aggregation output type: number for mean, integer for count and boolean sum, and
	// the scorer's type otherwise.
	ValueType EvaluationVersionQueryResultListResponseMetricsValueType `json:"value_type" api:"required"`
	JSON      evaluationVersionQueryResultListResponseMetricJSON       `json:"-"`
}

// evaluationVersionQueryResultListResponseMetricJSON contains the JSON metadata
// for the struct [EvaluationVersionQueryResultListResponseMetric]
type evaluationVersionQueryResultListResponseMetricJSON struct {
	Aggregation     apijson.Field
	Name            apijson.Field
	ScorerID        apijson.Field
	ScorerVersionID apijson.Field
	ValueType       apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *EvaluationVersionQueryResultListResponseMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionQueryResultListResponseMetricJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionQueryResultListResponseMetricsAggregation string

const (
	EvaluationVersionQueryResultListResponseMetricsAggregationMean  EvaluationVersionQueryResultListResponseMetricsAggregation = "mean"
	EvaluationVersionQueryResultListResponseMetricsAggregationSum   EvaluationVersionQueryResultListResponseMetricsAggregation = "sum"
	EvaluationVersionQueryResultListResponseMetricsAggregationMin   EvaluationVersionQueryResultListResponseMetricsAggregation = "min"
	EvaluationVersionQueryResultListResponseMetricsAggregationMax   EvaluationVersionQueryResultListResponseMetricsAggregation = "max"
	EvaluationVersionQueryResultListResponseMetricsAggregationCount EvaluationVersionQueryResultListResponseMetricsAggregation = "count"
)

func (r EvaluationVersionQueryResultListResponseMetricsAggregation) IsKnown() bool {
	switch r {
	case EvaluationVersionQueryResultListResponseMetricsAggregationMean,
		EvaluationVersionQueryResultListResponseMetricsAggregationSum,
		EvaluationVersionQueryResultListResponseMetricsAggregationMin,
		EvaluationVersionQueryResultListResponseMetricsAggregationMax,
		EvaluationVersionQueryResultListResponseMetricsAggregationCount:
		return true
	}
	return false
}

// Aggregation output type: number for mean, integer for count and boolean sum, and
// the scorer's type otherwise.
type EvaluationVersionQueryResultListResponseMetricsValueType string

const (
	EvaluationVersionQueryResultListResponseMetricsValueTypeBoolean EvaluationVersionQueryResultListResponseMetricsValueType = "boolean"
	EvaluationVersionQueryResultListResponseMetricsValueTypeInteger EvaluationVersionQueryResultListResponseMetricsValueType = "integer"
	EvaluationVersionQueryResultListResponseMetricsValueTypeNumber  EvaluationVersionQueryResultListResponseMetricsValueType = "number"
	EvaluationVersionQueryResultListResponseMetricsValueTypeString  EvaluationVersionQueryResultListResponseMetricsValueType = "string"
)

func (r EvaluationVersionQueryResultListResponseMetricsValueType) IsKnown() bool {
	switch r {
	case EvaluationVersionQueryResultListResponseMetricsValueTypeBoolean,
		EvaluationVersionQueryResultListResponseMetricsValueTypeInteger,
		EvaluationVersionQueryResultListResponseMetricsValueTypeNumber,
		EvaluationVersionQueryResultListResponseMetricsValueTypeString:
		return true
	}
	return false
}

// Namespace that assigned the scope reference.
type EvaluationVersionQueryResultListResponseNamespace string

const (
	EvaluationVersionQueryResultListResponseNamespaceWandb EvaluationVersionQueryResultListResponseNamespace = "wandb"
)

func (r EvaluationVersionQueryResultListResponseNamespace) IsKnown() bool {
	switch r {
	case EvaluationVersionQueryResultListResponseNamespaceWandb:
		return true
	}
	return false
}

type EvaluationVersionQueryResultListResponseRow struct {
	// Values in group_by order, with null for missing or null fields. Empty for an
	// ungrouped query.
	Group   []EvaluationVersionQueryResultListResponseRowsGroupUnion `json:"group"   api:"required"`
	Metrics []EvaluationVersionQueryResultListResponseRowsMetric     `json:"metrics" api:"required"`
	// Number of Cases in this group.
	Targets int64                                           `json:"targets" api:"required"`
	JSON    evaluationVersionQueryResultListResponseRowJSON `json:"-"`
}

// evaluationVersionQueryResultListResponseRowJSON contains the JSON metadata for
// the struct [EvaluationVersionQueryResultListResponseRow]
type evaluationVersionQueryResultListResponseRowJSON struct {
	Group       apijson.Field
	Metrics     apijson.Field
	Targets     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationVersionQueryResultListResponseRow) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionQueryResultListResponseRowJSON) RawJSON() string {
	return r.raw
}

// A JSON scalar (boolean, integer, number, or string) or null.
//
// Union satisfied by [shared.UnionBool], [shared.UnionInt], [shared.UnionFloat] or
// [shared.UnionString].
type EvaluationVersionQueryResultListResponseRowsGroupUnion interface {
	ImplementsEvaluationVersionQueryResultListResponseRowsGroupUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*EvaluationVersionQueryResultListResponseRowsGroupUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.True,
			Type:       reflect.TypeOf(shared.UnionBool(false)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.False,
			Type:       reflect.TypeOf(shared.UnionBool(false)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionInt(0)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(shared.UnionString("")),
		},
	)
}

type EvaluationVersionQueryResultListResponseRowsMetric struct {
	Name string `json:"name" api:"required"`
	// Aggregate value. When no targets have a non-null latest visible Score, count
	// returns 0 and other aggregations return null.
	Value EvaluationVersionQueryResultListResponseRowsMetricsValueUnion `json:"value" api:"required"`
	JSON  evaluationVersionQueryResultListResponseRowsMetricJSON        `json:"-"`
}

// evaluationVersionQueryResultListResponseRowsMetricJSON contains the JSON
// metadata for the struct [EvaluationVersionQueryResultListResponseRowsMetric]
type evaluationVersionQueryResultListResponseRowsMetricJSON struct {
	Name        apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationVersionQueryResultListResponseRowsMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionQueryResultListResponseRowsMetricJSON) RawJSON() string {
	return r.raw
}

// Aggregate value. When no targets have a non-null latest visible Score, count
// returns 0 and other aggregations return null.
//
// Union satisfied by [shared.UnionBool], [shared.UnionInt], [shared.UnionFloat] or
// [shared.UnionString].
type EvaluationVersionQueryResultListResponseRowsMetricsValueUnion interface {
	ImplementsEvaluationVersionQueryResultListResponseRowsMetricsValueUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*EvaluationVersionQueryResultListResponseRowsMetricsValueUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.True,
			Type:       reflect.TypeOf(shared.UnionBool(false)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.False,
			Type:       reflect.TypeOf(shared.UnionBool(false)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionInt(0)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(shared.UnionString("")),
		},
	)
}

type EvaluationVersionQueryResultListParams struct {
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
	// Cursor from the previous page's next_cursor. Omit for the first page.
	Cursor param.Field[string] `query:"cursor"`
	Limit  param.Field[int64]  `query:"limit"`
}

// URLQuery serializes [EvaluationVersionQueryResultListParams]'s query parameters
// as `url.Values`.
func (r EvaluationVersionQueryResultListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
