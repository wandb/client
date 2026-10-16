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

// EvaluationQueryResultService contains methods and other services that help with
// interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvaluationQueryResultService] method instead.
type EvaluationQueryResultService struct {
	Options []option.RequestOption
}

// NewEvaluationQueryResultService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewEvaluationQueryResultService(opts ...option.RequestOption) (r *EvaluationQueryResultService) {
	r = &EvaluationQueryResultService{}
	r.Options = opts
	return
}

// Runs one of the Evaluation's saved queries on its latest data and returns
// aggregated Scores. The first page reads every Score published so far over the
// Evaluation's registered Scorers. A fixed Evaluation reads its DatasetVersion's
// Cases, and a live Evaluation reads every Case published to its Dataset head.
// Later pages keep the first page's Cases and Scores. Results sort by the group_by
// columns in order, ascending with nulls last. For each grouping column, missing
// values and explicit nulls are treated the same. Without grouping, the result
// contains one row with an empty group array, even for an empty dataset. With
// grouping, an empty dataset returns no rows.
func (r *EvaluationQueryResultService) List(
	ctx context.Context,
	evaluationID string,
	evaluationQueryID string,
	params EvaluationQueryResultListParams,
	opts ...option.RequestOption,
) (res *EvaluationQueryResultListResponse, err error) {
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
	if evaluationID == "" {
		err = errors.New("missing required evaluation_id parameter")
		return nil, err
	}
	if evaluationQueryID == "" {
		err = errors.New("missing required evaluation_query_id parameter")
		return nil, err
	}
	path := fmt.Sprintf(
		"v1/namespaces/%s/scopes/%s/evaluations/%s/queries/%s/results",
		params.Namespace,
		params.ScopeID,
		evaluationID,
		evaluationQueryID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

type EvaluationQueryResultListResponse struct {
	// The DatasetVersion whose Cases the results read. For a live Evaluation's latest
	// data, this is the Dataset head when the first page was requested.
	DatasetVersionID  string                                     `json:"dataset_version_id"  api:"required" format:"uuid"`
	EvaluationID      string                                     `json:"evaluation_id"       api:"required" format:"uuid"`
	EvaluationQueryID string                                     `json:"evaluation_query_id" api:"required" format:"uuid"`
	GroupBy           []EvaluationQueryResultListResponseGroupBy `json:"group_by"            api:"required"`
	Metrics           []EvaluationQueryResultListResponseMetric  `json:"metrics"             api:"required"`
	// Namespace that assigned the scope reference.
	Namespace EvaluationQueryResultListResponseNamespace `json:"namespace" api:"required"`
	// The results include Scores published up to this commit. For a live Evaluation's
	// latest data, Cases too. Every page of one traversal uses the first page's
	// commit.
	PublishedCommitID int64                                  `json:"published_commit_id" api:"required"`
	Rows              []EvaluationQueryResultListResponseRow `json:"rows"                api:"required"`
	ScopeRef          string                                 `json:"scope_ref"           api:"required"`
	// The EvaluationVersion whose data the results read. Absent when they read the
	// Evaluation's latest data.
	EvaluationVersionID string `json:"evaluation_version_id" format:"uuid"`
	// Cursor for the next page. Absent on the last page.
	NextCursor string                                `json:"next_cursor"`
	JSON       evaluationQueryResultListResponseJSON `json:"-"`
}

// evaluationQueryResultListResponseJSON contains the JSON metadata for the struct
// [EvaluationQueryResultListResponse]
type evaluationQueryResultListResponseJSON struct {
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

func (r *EvaluationQueryResultListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryResultListResponseJSON) RawJSON() string {
	return r.raw
}

type EvaluationQueryResultListResponseGroupBy struct {
	Kind      EvaluationQueryResultListResponseGroupByKind      `json:"kind"       api:"required"`
	ValueType EvaluationQueryResultListResponseGroupByValueType `json:"value_type" api:"required"`
	Name      string                                            `json:"name"`
	Source    EvaluationQueryResultListResponseGroupBySource    `json:"source"`
	JSON      evaluationQueryResultListResponseGroupByJSON      `json:"-"`
}

// evaluationQueryResultListResponseGroupByJSON contains the JSON metadata for the
// struct [EvaluationQueryResultListResponseGroupBy]
type evaluationQueryResultListResponseGroupByJSON struct {
	Kind        apijson.Field
	ValueType   apijson.Field
	Name        apijson.Field
	Source      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryResultListResponseGroupBy) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryResultListResponseGroupByJSON) RawJSON() string {
	return r.raw
}

type EvaluationQueryResultListResponseGroupByKind string

const (
	EvaluationQueryResultListResponseGroupByKindSubject      EvaluationQueryResultListResponseGroupByKind = "subject"
	EvaluationQueryResultListResponseGroupByKindDatasetField EvaluationQueryResultListResponseGroupByKind = "dataset_field"
	EvaluationQueryResultListResponseGroupByKindScore        EvaluationQueryResultListResponseGroupByKind = "score"
)

func (r EvaluationQueryResultListResponseGroupByKind) IsKnown() bool {
	switch r {
	case EvaluationQueryResultListResponseGroupByKindSubject,
		EvaluationQueryResultListResponseGroupByKindDatasetField,
		EvaluationQueryResultListResponseGroupByKindScore:
		return true
	}
	return false
}

type EvaluationQueryResultListResponseGroupByValueType string

const (
	EvaluationQueryResultListResponseGroupByValueTypeBoolean EvaluationQueryResultListResponseGroupByValueType = "boolean"
	EvaluationQueryResultListResponseGroupByValueTypeInteger EvaluationQueryResultListResponseGroupByValueType = "integer"
	EvaluationQueryResultListResponseGroupByValueTypeNumber  EvaluationQueryResultListResponseGroupByValueType = "number"
	EvaluationQueryResultListResponseGroupByValueTypeString  EvaluationQueryResultListResponseGroupByValueType = "string"
)

func (r EvaluationQueryResultListResponseGroupByValueType) IsKnown() bool {
	switch r {
	case EvaluationQueryResultListResponseGroupByValueTypeBoolean,
		EvaluationQueryResultListResponseGroupByValueTypeInteger,
		EvaluationQueryResultListResponseGroupByValueTypeNumber,
		EvaluationQueryResultListResponseGroupByValueTypeString:
		return true
	}
	return false
}

type EvaluationQueryResultListResponseGroupBySource string

const (
	EvaluationQueryResultListResponseGroupBySourceInput  EvaluationQueryResultListResponseGroupBySource = "input"
	EvaluationQueryResultListResponseGroupBySourceOutput EvaluationQueryResultListResponseGroupBySource = "output"
)

func (r EvaluationQueryResultListResponseGroupBySource) IsKnown() bool {
	switch r {
	case EvaluationQueryResultListResponseGroupBySourceInput, EvaluationQueryResultListResponseGroupBySourceOutput:
		return true
	}
	return false
}

type EvaluationQueryResultListResponseMetric struct {
	Aggregation     EvaluationQueryResultListResponseMetricsAggregation `json:"aggregation"       api:"required"`
	Name            string                                              `json:"name"              api:"required"`
	ScorerID        string                                              `json:"scorer_id"         api:"required" format:"uuid"`
	ScorerVersionID string                                              `json:"scorer_version_id" api:"required" format:"uuid"`
	// Aggregation output type: number for mean, integer for count and boolean sum, and
	// the scorer's type otherwise.
	ValueType EvaluationQueryResultListResponseMetricsValueType `json:"value_type" api:"required"`
	JSON      evaluationQueryResultListResponseMetricJSON       `json:"-"`
}

// evaluationQueryResultListResponseMetricJSON contains the JSON metadata for the
// struct [EvaluationQueryResultListResponseMetric]
type evaluationQueryResultListResponseMetricJSON struct {
	Aggregation     apijson.Field
	Name            apijson.Field
	ScorerID        apijson.Field
	ScorerVersionID apijson.Field
	ValueType       apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *EvaluationQueryResultListResponseMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryResultListResponseMetricJSON) RawJSON() string {
	return r.raw
}

type EvaluationQueryResultListResponseMetricsAggregation string

const (
	EvaluationQueryResultListResponseMetricsAggregationMean  EvaluationQueryResultListResponseMetricsAggregation = "mean"
	EvaluationQueryResultListResponseMetricsAggregationSum   EvaluationQueryResultListResponseMetricsAggregation = "sum"
	EvaluationQueryResultListResponseMetricsAggregationMin   EvaluationQueryResultListResponseMetricsAggregation = "min"
	EvaluationQueryResultListResponseMetricsAggregationMax   EvaluationQueryResultListResponseMetricsAggregation = "max"
	EvaluationQueryResultListResponseMetricsAggregationCount EvaluationQueryResultListResponseMetricsAggregation = "count"
)

func (r EvaluationQueryResultListResponseMetricsAggregation) IsKnown() bool {
	switch r {
	case EvaluationQueryResultListResponseMetricsAggregationMean,
		EvaluationQueryResultListResponseMetricsAggregationSum,
		EvaluationQueryResultListResponseMetricsAggregationMin,
		EvaluationQueryResultListResponseMetricsAggregationMax,
		EvaluationQueryResultListResponseMetricsAggregationCount:
		return true
	}
	return false
}

// Aggregation output type: number for mean, integer for count and boolean sum, and
// the scorer's type otherwise.
type EvaluationQueryResultListResponseMetricsValueType string

const (
	EvaluationQueryResultListResponseMetricsValueTypeBoolean EvaluationQueryResultListResponseMetricsValueType = "boolean"
	EvaluationQueryResultListResponseMetricsValueTypeInteger EvaluationQueryResultListResponseMetricsValueType = "integer"
	EvaluationQueryResultListResponseMetricsValueTypeNumber  EvaluationQueryResultListResponseMetricsValueType = "number"
	EvaluationQueryResultListResponseMetricsValueTypeString  EvaluationQueryResultListResponseMetricsValueType = "string"
)

func (r EvaluationQueryResultListResponseMetricsValueType) IsKnown() bool {
	switch r {
	case EvaluationQueryResultListResponseMetricsValueTypeBoolean,
		EvaluationQueryResultListResponseMetricsValueTypeInteger,
		EvaluationQueryResultListResponseMetricsValueTypeNumber,
		EvaluationQueryResultListResponseMetricsValueTypeString:
		return true
	}
	return false
}

// Namespace that assigned the scope reference.
type EvaluationQueryResultListResponseNamespace string

const (
	EvaluationQueryResultListResponseNamespaceWandb EvaluationQueryResultListResponseNamespace = "wandb"
)

func (r EvaluationQueryResultListResponseNamespace) IsKnown() bool {
	switch r {
	case EvaluationQueryResultListResponseNamespaceWandb:
		return true
	}
	return false
}

type EvaluationQueryResultListResponseRow struct {
	// Values in group_by order, with null for missing or null fields. Empty for an
	// ungrouped query.
	Group   []EvaluationQueryResultListResponseRowsGroupUnion `json:"group"   api:"required"`
	Metrics []EvaluationQueryResultListResponseRowsMetric     `json:"metrics" api:"required"`
	// Number of Cases in this group.
	Targets int64                                    `json:"targets" api:"required"`
	JSON    evaluationQueryResultListResponseRowJSON `json:"-"`
}

// evaluationQueryResultListResponseRowJSON contains the JSON metadata for the
// struct [EvaluationQueryResultListResponseRow]
type evaluationQueryResultListResponseRowJSON struct {
	Group       apijson.Field
	Metrics     apijson.Field
	Targets     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryResultListResponseRow) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryResultListResponseRowJSON) RawJSON() string {
	return r.raw
}

// A JSON scalar (boolean, integer, number, or string) or null.
//
// Union satisfied by [shared.UnionBool], [shared.UnionInt], [shared.UnionFloat] or
// [shared.UnionString].
type EvaluationQueryResultListResponseRowsGroupUnion interface {
	ImplementsEvaluationQueryResultListResponseRowsGroupUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*EvaluationQueryResultListResponseRowsGroupUnion)(nil)).Elem(),
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

type EvaluationQueryResultListResponseRowsMetric struct {
	Name string `json:"name" api:"required"`
	// Aggregate value. When no targets have a non-null latest visible Score, count
	// returns 0 and other aggregations return null.
	Value EvaluationQueryResultListResponseRowsMetricsValueUnion `json:"value" api:"required"`
	JSON  evaluationQueryResultListResponseRowsMetricJSON        `json:"-"`
}

// evaluationQueryResultListResponseRowsMetricJSON contains the JSON metadata for
// the struct [EvaluationQueryResultListResponseRowsMetric]
type evaluationQueryResultListResponseRowsMetricJSON struct {
	Name        apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryResultListResponseRowsMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryResultListResponseRowsMetricJSON) RawJSON() string {
	return r.raw
}

// Aggregate value. When no targets have a non-null latest visible Score, count
// returns 0 and other aggregations return null.
//
// Union satisfied by [shared.UnionBool], [shared.UnionInt], [shared.UnionFloat] or
// [shared.UnionString].
type EvaluationQueryResultListResponseRowsMetricsValueUnion interface {
	ImplementsEvaluationQueryResultListResponseRowsMetricsValueUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*EvaluationQueryResultListResponseRowsMetricsValueUnion)(nil)).Elem(),
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

type EvaluationQueryResultListParams struct {
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

// URLQuery serializes [EvaluationQueryResultListParams]'s query parameters as
// `url.Values`.
func (r EvaluationQueryResultListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
