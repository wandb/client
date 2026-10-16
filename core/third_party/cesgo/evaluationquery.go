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
)

// EvaluationQueryService contains methods and other services that help with
// interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvaluationQueryService] method instead.
type EvaluationQueryService struct {
	Options []option.RequestOption
	Results *EvaluationQueryResultService
}

// NewEvaluationQueryService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewEvaluationQueryService(opts ...option.RequestOption) (r *EvaluationQueryService) {
	r = &EvaluationQueryService{}
	r.Options = opts
	r.Results = NewEvaluationQueryResultService(opts...)
	return
}

// Saves grouping fields and Score aggregations for an Evaluation of Cases. The
// definition must fit the Evaluation's latest data. This operation does not
// compute the aggregates.
//
// Request bodies must be smaller than 1 MiB; larger bodies are rejected with 413
// Request Entity Too Large.
func (r *EvaluationQueryService) New(
	ctx context.Context,
	evaluationID string,
	params EvaluationQueryNewParams,
	opts ...option.RequestOption,
) (res *EvaluationQueryNewResponse, err error) {
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
	path := fmt.Sprintf(
		"v1/namespaces/%s/scopes/%s/evaluations/%s/queries",
		params.Namespace,
		params.ScopeID,
		evaluationID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Pages one Evaluation's saved queries by query ID.
func (r *EvaluationQueryService) List(
	ctx context.Context,
	evaluationID string,
	params EvaluationQueryListParams,
	opts ...option.RequestOption,
) (res *EvaluationQueryListResponse, err error) {
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
	path := fmt.Sprintf(
		"v1/namespaces/%s/scopes/%s/evaluations/%s/queries",
		params.Namespace,
		params.ScopeID,
		evaluationID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Returns the definition as stored. A query under another Evaluation or scope
// returns 404.
func (r *EvaluationQueryService) Get(
	ctx context.Context,
	evaluationID string,
	evaluationQueryID string,
	query EvaluationQueryGetParams,
	opts ...option.RequestOption,
) (res *EvaluationQueryGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	precfg, err := requestconfig.PreRequestOptions(opts...)
	if err != nil {
		return nil, err
	}
	requestconfig.UseDefaultParam(&query.Namespace, precfg.Namespace)
	requestconfig.UseDefaultParam(&query.ScopeID, precfg.ScopeID)
	if query.Namespace.Value == "" {
		err = errors.New("missing required namespace parameter")
		return nil, err
	}
	if query.ScopeID.Value == "" {
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
		"v1/namespaces/%s/scopes/%s/evaluations/%s/queries/%s",
		query.Namespace,
		query.ScopeID,
		evaluationID,
		evaluationQueryID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type EvaluationQueryNewResponse struct {
	Definition        EvaluationQueryNewResponseDefinition `json:"definition"          api:"required"`
	EvaluationID      string                               `json:"evaluation_id"       api:"required" format:"uuid"`
	EvaluationQueryID string                               `json:"evaluation_query_id" api:"required" format:"uuid"`
	Name              string                               `json:"name"                api:"required"`
	// Namespace that assigned the scope reference.
	Namespace EvaluationQueryNewResponseNamespace `json:"namespace" api:"required"`
	ScopeRef  string                              `json:"scope_ref" api:"required"`
	JSON      evaluationQueryNewResponseJSON      `json:"-"`
}

// evaluationQueryNewResponseJSON contains the JSON metadata for the struct
// [EvaluationQueryNewResponse]
type evaluationQueryNewResponseJSON struct {
	Definition        apijson.Field
	EvaluationID      apijson.Field
	EvaluationQueryID apijson.Field
	Name              apijson.Field
	Namespace         apijson.Field
	ScopeRef          apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *EvaluationQueryNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryNewResponseJSON) RawJSON() string {
	return r.raw
}

type EvaluationQueryNewResponseDefinition struct {
	Metrics []EvaluationQueryNewResponseDefinitionMetric `json:"metrics" api:"required"`
	// Stored definition version; always present in responses.
	V       int64                                         `json:"v"        api:"required"`
	GroupBy []EvaluationQueryNewResponseDefinitionGroupBy `json:"group_by"`
	JSON    evaluationQueryNewResponseDefinitionJSON      `json:"-"`
}

// evaluationQueryNewResponseDefinitionJSON contains the JSON metadata for the
// struct [EvaluationQueryNewResponseDefinition]
type evaluationQueryNewResponseDefinitionJSON struct {
	Metrics     apijson.Field
	V           apijson.Field
	GroupBy     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryNewResponseDefinition) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryNewResponseDefinitionJSON) RawJSON() string {
	return r.raw
}

type EvaluationQueryNewResponseDefinitionMetric struct {
	// Mean accepts integer, number, and boolean (fraction true). Sum accepts integer,
	// number, and boolean (count of true non-null Scores). Min and max accept integer,
	// number, boolean, and string. Strings compare lexicographically. Count accepts
	// these four types and counts non-null latest visible Scores, including zero and
	// false.
	Aggregation     EvaluationQueryNewResponseDefinitionMetricsAggregation `json:"aggregation"       api:"required"`
	Name            string                                                 `json:"name"              api:"required"`
	ScorerVersionID string                                                 `json:"scorer_version_id" api:"required" format:"uuid"`
	JSON            evaluationQueryNewResponseDefinitionMetricJSON         `json:"-"`
}

// evaluationQueryNewResponseDefinitionMetricJSON contains the JSON metadata for
// the struct [EvaluationQueryNewResponseDefinitionMetric]
type evaluationQueryNewResponseDefinitionMetricJSON struct {
	Aggregation     apijson.Field
	Name            apijson.Field
	ScorerVersionID apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *EvaluationQueryNewResponseDefinitionMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryNewResponseDefinitionMetricJSON) RawJSON() string {
	return r.raw
}

// Mean accepts integer, number, and boolean (fraction true). Sum accepts integer,
// number, and boolean (count of true non-null Scores). Min and max accept integer,
// number, boolean, and string. Strings compare lexicographically. Count accepts
// these four types and counts non-null latest visible Scores, including zero and
// false.
type EvaluationQueryNewResponseDefinitionMetricsAggregation string

const (
	EvaluationQueryNewResponseDefinitionMetricsAggregationMean  EvaluationQueryNewResponseDefinitionMetricsAggregation = "mean"
	EvaluationQueryNewResponseDefinitionMetricsAggregationSum   EvaluationQueryNewResponseDefinitionMetricsAggregation = "sum"
	EvaluationQueryNewResponseDefinitionMetricsAggregationMin   EvaluationQueryNewResponseDefinitionMetricsAggregation = "min"
	EvaluationQueryNewResponseDefinitionMetricsAggregationMax   EvaluationQueryNewResponseDefinitionMetricsAggregation = "max"
	EvaluationQueryNewResponseDefinitionMetricsAggregationCount EvaluationQueryNewResponseDefinitionMetricsAggregation = "count"
)

func (r EvaluationQueryNewResponseDefinitionMetricsAggregation) IsKnown() bool {
	switch r {
	case EvaluationQueryNewResponseDefinitionMetricsAggregationMean,
		EvaluationQueryNewResponseDefinitionMetricsAggregationSum,
		EvaluationQueryNewResponseDefinitionMetricsAggregationMin,
		EvaluationQueryNewResponseDefinitionMetricsAggregationMax,
		EvaluationQueryNewResponseDefinitionMetricsAggregationCount:
		return true
	}
	return false
}

type EvaluationQueryNewResponseDefinitionGroupBy struct {
	Kind EvaluationQueryNewResponseDefinitionGroupByKind `json:"kind" api:"required"`
	// Case field name.
	Name string `json:"name"`
	// Case field location.
	Source EvaluationQueryNewResponseDefinitionGroupBySource `json:"source"`
	JSON   evaluationQueryNewResponseDefinitionGroupByJSON   `json:"-"`
	union  EvaluationQueryNewResponseDefinitionGroupByUnion
}

// evaluationQueryNewResponseDefinitionGroupByJSON contains the JSON metadata for
// the struct [EvaluationQueryNewResponseDefinitionGroupBy]
type evaluationQueryNewResponseDefinitionGroupByJSON struct {
	Kind        apijson.Field
	Name        apijson.Field
	Source      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r evaluationQueryNewResponseDefinitionGroupByJSON) RawJSON() string {
	return r.raw
}

func (r *EvaluationQueryNewResponseDefinitionGroupBy) UnmarshalJSON(data []byte) (err error) {
	*r = EvaluationQueryNewResponseDefinitionGroupBy{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [EvaluationQueryNewResponseDefinitionGroupByUnion] interface
// which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumn],
// [EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumn],
// [EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumn].
func (r EvaluationQueryNewResponseDefinitionGroupBy) AsUnion() EvaluationQueryNewResponseDefinitionGroupByUnion {
	return r.union
}

// Union satisfied by
// [EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumn],
// [EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumn] or
// [EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumn].
type EvaluationQueryNewResponseDefinitionGroupByUnion interface {
	implementsEvaluationQueryNewResponseDefinitionGroupBy()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*EvaluationQueryNewResponseDefinitionGroupByUnion)(nil)).Elem(),
		"kind",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumn{}),
			DiscriminatorValue: "subject",
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type: reflect.TypeOf(
				EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumn{},
			),
			DiscriminatorValue: "dataset_field",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumn{}),
			DiscriminatorValue: "score",
		},
	)
}

type EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumn struct {
	Kind EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumnKind `json:"kind" api:"required"`
	JSON evaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumnJSON `json:"-"`
}

// evaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumnJSON contains
// the JSON metadata for the struct
// [EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumn]
type evaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumnJSON struct {
	Kind        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumn) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumnJSON) RawJSON() string {
	return r.raw
}

func (r EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumn) implementsEvaluationQueryNewResponseDefinitionGroupBy() {
}

type EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumnKind string

const (
	EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumnKindSubject EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumnKind = "subject"
)

func (r EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumnKind) IsKnown() bool {
	switch r {
	case EvaluationQueryNewResponseDefinitionGroupBySubjectEvaluationColumnKindSubject:
		return true
	}
	return false
}

type EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumn struct {
	Kind EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnKind `json:"kind" api:"required"`
	// Case field name.
	Name string `json:"name" api:"required"`
	// Case field location.
	Source EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnSource `json:"source" api:"required"`
	JSON   evaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnJSON   `json:"-"`
}

// evaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnJSON
// contains the JSON metadata for the struct
// [EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumn]
type evaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnJSON struct {
	Kind        apijson.Field
	Name        apijson.Field
	Source      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumn) UnmarshalJSON(
	data []byte,
) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnJSON) RawJSON() string {
	return r.raw
}

func (r EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumn) implementsEvaluationQueryNewResponseDefinitionGroupBy() {
}

type EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnKind string

const (
	EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnKindDatasetField EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnKind = "dataset_field"
)

func (r EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnKind) IsKnown() bool {
	switch r {
	case EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnKindDatasetField:
		return true
	}
	return false
}

// Case field location.
type EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnSource string

const (
	EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnSourceInput  EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnSource = "input"
	EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnSourceOutput EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnSource = "output"
)

func (r EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnSource) IsKnown() bool {
	switch r {
	case EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnSourceInput,
		EvaluationQueryNewResponseDefinitionGroupByDatasetFieldEvaluationColumnSourceOutput:
		return true
	}
	return false
}

type EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumn struct {
	Kind EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumnKind `json:"kind" api:"required"`
	// Scorer name selected by this EvaluationVersion.
	Name string                                                               `json:"name" api:"required"`
	JSON evaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumnJSON `json:"-"`
}

// evaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumnJSON contains
// the JSON metadata for the struct
// [EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumn]
type evaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumnJSON struct {
	Kind        apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumn) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumnJSON) RawJSON() string {
	return r.raw
}

func (r EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumn) implementsEvaluationQueryNewResponseDefinitionGroupBy() {
}

type EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumnKind string

const (
	EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumnKindScore EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumnKind = "score"
)

func (r EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumnKind) IsKnown() bool {
	switch r {
	case EvaluationQueryNewResponseDefinitionGroupByScoreEvaluationColumnKindScore:
		return true
	}
	return false
}

type EvaluationQueryNewResponseDefinitionGroupByKind string

const (
	EvaluationQueryNewResponseDefinitionGroupByKindSubject      EvaluationQueryNewResponseDefinitionGroupByKind = "subject"
	EvaluationQueryNewResponseDefinitionGroupByKindDatasetField EvaluationQueryNewResponseDefinitionGroupByKind = "dataset_field"
	EvaluationQueryNewResponseDefinitionGroupByKindScore        EvaluationQueryNewResponseDefinitionGroupByKind = "score"
)

func (r EvaluationQueryNewResponseDefinitionGroupByKind) IsKnown() bool {
	switch r {
	case EvaluationQueryNewResponseDefinitionGroupByKindSubject,
		EvaluationQueryNewResponseDefinitionGroupByKindDatasetField,
		EvaluationQueryNewResponseDefinitionGroupByKindScore:
		return true
	}
	return false
}

// Case field location.
type EvaluationQueryNewResponseDefinitionGroupBySource string

const (
	EvaluationQueryNewResponseDefinitionGroupBySourceInput  EvaluationQueryNewResponseDefinitionGroupBySource = "input"
	EvaluationQueryNewResponseDefinitionGroupBySourceOutput EvaluationQueryNewResponseDefinitionGroupBySource = "output"
)

func (r EvaluationQueryNewResponseDefinitionGroupBySource) IsKnown() bool {
	switch r {
	case EvaluationQueryNewResponseDefinitionGroupBySourceInput,
		EvaluationQueryNewResponseDefinitionGroupBySourceOutput:
		return true
	}
	return false
}

// Namespace that assigned the scope reference.
type EvaluationQueryNewResponseNamespace string

const (
	EvaluationQueryNewResponseNamespaceWandb EvaluationQueryNewResponseNamespace = "wandb"
)

func (r EvaluationQueryNewResponseNamespace) IsKnown() bool {
	switch r {
	case EvaluationQueryNewResponseNamespaceWandb:
		return true
	}
	return false
}

type EvaluationQueryListResponse struct {
	EvaluationQueries []EvaluationQueryListResponseEvaluationQuery `json:"evaluation_queries" api:"required"`
	NextCursor        string                                       `json:"next_cursor"`
	JSON              evaluationQueryListResponseJSON              `json:"-"`
}

// evaluationQueryListResponseJSON contains the JSON metadata for the struct
// [EvaluationQueryListResponse]
type evaluationQueryListResponseJSON struct {
	EvaluationQueries apijson.Field
	NextCursor        apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *EvaluationQueryListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryListResponseJSON) RawJSON() string {
	return r.raw
}

type EvaluationQueryListResponseEvaluationQuery struct {
	Definition        EvaluationQueryListResponseEvaluationQueriesDefinition `json:"definition"          api:"required"`
	EvaluationID      string                                                 `json:"evaluation_id"       api:"required" format:"uuid"`
	EvaluationQueryID string                                                 `json:"evaluation_query_id" api:"required" format:"uuid"`
	Name              string                                                 `json:"name"                api:"required"`
	// Namespace that assigned the scope reference.
	Namespace EvaluationQueryListResponseEvaluationQueriesNamespace `json:"namespace" api:"required"`
	ScopeRef  string                                                `json:"scope_ref" api:"required"`
	JSON      evaluationQueryListResponseEvaluationQueryJSON        `json:"-"`
}

// evaluationQueryListResponseEvaluationQueryJSON contains the JSON metadata for
// the struct [EvaluationQueryListResponseEvaluationQuery]
type evaluationQueryListResponseEvaluationQueryJSON struct {
	Definition        apijson.Field
	EvaluationID      apijson.Field
	EvaluationQueryID apijson.Field
	Name              apijson.Field
	Namespace         apijson.Field
	ScopeRef          apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *EvaluationQueryListResponseEvaluationQuery) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryListResponseEvaluationQueryJSON) RawJSON() string {
	return r.raw
}

type EvaluationQueryListResponseEvaluationQueriesDefinition struct {
	Metrics []EvaluationQueryListResponseEvaluationQueriesDefinitionMetric `json:"metrics" api:"required"`
	// Stored definition version; always present in responses.
	V       int64                                                           `json:"v"        api:"required"`
	GroupBy []EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBy `json:"group_by"`
	JSON    evaluationQueryListResponseEvaluationQueriesDefinitionJSON      `json:"-"`
}

// evaluationQueryListResponseEvaluationQueriesDefinitionJSON contains the JSON
// metadata for the struct [EvaluationQueryListResponseEvaluationQueriesDefinition]
type evaluationQueryListResponseEvaluationQueriesDefinitionJSON struct {
	Metrics     apijson.Field
	V           apijson.Field
	GroupBy     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryListResponseEvaluationQueriesDefinition) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryListResponseEvaluationQueriesDefinitionJSON) RawJSON() string {
	return r.raw
}

type EvaluationQueryListResponseEvaluationQueriesDefinitionMetric struct {
	// Mean accepts integer, number, and boolean (fraction true). Sum accepts integer,
	// number, and boolean (count of true non-null Scores). Min and max accept integer,
	// number, boolean, and string. Strings compare lexicographically. Count accepts
	// these four types and counts non-null latest visible Scores, including zero and
	// false.
	Aggregation     EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregation `json:"aggregation"       api:"required"`
	Name            string                                                                   `json:"name"              api:"required"`
	ScorerVersionID string                                                                   `json:"scorer_version_id" api:"required" format:"uuid"`
	JSON            evaluationQueryListResponseEvaluationQueriesDefinitionMetricJSON         `json:"-"`
}

// evaluationQueryListResponseEvaluationQueriesDefinitionMetricJSON contains the
// JSON metadata for the struct
// [EvaluationQueryListResponseEvaluationQueriesDefinitionMetric]
type evaluationQueryListResponseEvaluationQueriesDefinitionMetricJSON struct {
	Aggregation     apijson.Field
	Name            apijson.Field
	ScorerVersionID apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *EvaluationQueryListResponseEvaluationQueriesDefinitionMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryListResponseEvaluationQueriesDefinitionMetricJSON) RawJSON() string {
	return r.raw
}

// Mean accepts integer, number, and boolean (fraction true). Sum accepts integer,
// number, and boolean (count of true non-null Scores). Min and max accept integer,
// number, boolean, and string. Strings compare lexicographically. Count accepts
// these four types and counts non-null latest visible Scores, including zero and
// false.
type EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregation string

const (
	EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregationMean  EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregation = "mean"
	EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregationSum   EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregation = "sum"
	EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregationMin   EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregation = "min"
	EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregationMax   EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregation = "max"
	EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregationCount EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregation = "count"
)

func (r EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregation) IsKnown() bool {
	switch r {
	case EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregationMean,
		EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregationSum,
		EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregationMin,
		EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregationMax,
		EvaluationQueryListResponseEvaluationQueriesDefinitionMetricsAggregationCount:
		return true
	}
	return false
}

type EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBy struct {
	Kind EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByKind `json:"kind" api:"required"`
	// Case field name.
	Name string `json:"name"`
	// Case field location.
	Source EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySource `json:"source"`
	JSON   evaluationQueryListResponseEvaluationQueriesDefinitionGroupByJSON   `json:"-"`
	union  EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByUnion
}

// evaluationQueryListResponseEvaluationQueriesDefinitionGroupByJSON contains the
// JSON metadata for the struct
// [EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBy]
type evaluationQueryListResponseEvaluationQueriesDefinitionGroupByJSON struct {
	Kind        apijson.Field
	Name        apijson.Field
	Source      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r evaluationQueryListResponseEvaluationQueriesDefinitionGroupByJSON) RawJSON() string {
	return r.raw
}

func (r *EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBy) UnmarshalJSON(data []byte) (err error) {
	*r = EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBy{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByUnion] interface
// which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumn],
// [EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumn],
// [EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumn].
func (r EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBy) AsUnion() EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByUnion {
	return r.union
}

// Union satisfied by
// [EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumn],
// [EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumn]
// or
// [EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumn].
type EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByUnion interface {
	implementsEvaluationQueryListResponseEvaluationQueriesDefinitionGroupBy()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByUnion)(nil)).Elem(),
		"kind",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type: reflect.TypeOf(
				EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumn{},
			),
			DiscriminatorValue: "subject",
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type: reflect.TypeOf(
				EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumn{},
			),
			DiscriminatorValue: "dataset_field",
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type: reflect.TypeOf(
				EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumn{},
			),
			DiscriminatorValue: "score",
		},
	)
}

type EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumn struct {
	Kind EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumnKind `json:"kind" api:"required"`
	JSON evaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumnJSON `json:"-"`
}

// evaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumnJSON
// contains the JSON metadata for the struct
// [EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumn]
type evaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumnJSON struct {
	Kind        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumn) UnmarshalJSON(
	data []byte,
) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumnJSON) RawJSON() string {
	return r.raw
}

func (r EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumn) implementsEvaluationQueryListResponseEvaluationQueriesDefinitionGroupBy() {
}

type EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumnKind string

const (
	EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumnKindSubject EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumnKind = "subject"
)

func (r EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumnKind) IsKnown() bool {
	switch r {
	case EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySubjectEvaluationColumnKindSubject:
		return true
	}
	return false
}

type EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumn struct {
	Kind EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnKind `json:"kind" api:"required"`
	// Case field name.
	Name string `json:"name" api:"required"`
	// Case field location.
	Source EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnSource `json:"source" api:"required"`
	JSON   evaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnJSON   `json:"-"`
}

// evaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnJSON
// contains the JSON metadata for the struct
// [EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumn]
type evaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnJSON struct {
	Kind        apijson.Field
	Name        apijson.Field
	Source      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumn) UnmarshalJSON(
	data []byte,
) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnJSON) RawJSON() string {
	return r.raw
}

func (r EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumn) implementsEvaluationQueryListResponseEvaluationQueriesDefinitionGroupBy() {
}

type EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnKind string

const (
	EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnKindDatasetField EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnKind = "dataset_field"
)

func (r EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnKind) IsKnown() bool {
	switch r {
	case EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnKindDatasetField:
		return true
	}
	return false
}

// Case field location.
type EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnSource string

const (
	EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnSourceInput  EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnSource = "input"
	EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnSourceOutput EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnSource = "output"
)

func (r EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnSource) IsKnown() bool {
	switch r {
	case EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnSourceInput,
		EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByDatasetFieldEvaluationColumnSourceOutput:
		return true
	}
	return false
}

type EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumn struct {
	Kind EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumnKind `json:"kind" api:"required"`
	// Scorer name selected by this EvaluationVersion.
	Name string                                                                                 `json:"name" api:"required"`
	JSON evaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumnJSON `json:"-"`
}

// evaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumnJSON
// contains the JSON metadata for the struct
// [EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumn]
type evaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumnJSON struct {
	Kind        apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumn) UnmarshalJSON(
	data []byte,
) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumnJSON) RawJSON() string {
	return r.raw
}

func (r EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumn) implementsEvaluationQueryListResponseEvaluationQueriesDefinitionGroupBy() {
}

type EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumnKind string

const (
	EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumnKindScore EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumnKind = "score"
)

func (r EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumnKind) IsKnown() bool {
	switch r {
	case EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByScoreEvaluationColumnKindScore:
		return true
	}
	return false
}

type EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByKind string

const (
	EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByKindSubject      EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByKind = "subject"
	EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByKindDatasetField EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByKind = "dataset_field"
	EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByKindScore        EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByKind = "score"
)

func (r EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByKind) IsKnown() bool {
	switch r {
	case EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByKindSubject,
		EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByKindDatasetField,
		EvaluationQueryListResponseEvaluationQueriesDefinitionGroupByKindScore:
		return true
	}
	return false
}

// Case field location.
type EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySource string

const (
	EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySourceInput  EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySource = "input"
	EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySourceOutput EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySource = "output"
)

func (r EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySource) IsKnown() bool {
	switch r {
	case EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySourceInput,
		EvaluationQueryListResponseEvaluationQueriesDefinitionGroupBySourceOutput:
		return true
	}
	return false
}

// Namespace that assigned the scope reference.
type EvaluationQueryListResponseEvaluationQueriesNamespace string

const (
	EvaluationQueryListResponseEvaluationQueriesNamespaceWandb EvaluationQueryListResponseEvaluationQueriesNamespace = "wandb"
)

func (r EvaluationQueryListResponseEvaluationQueriesNamespace) IsKnown() bool {
	switch r {
	case EvaluationQueryListResponseEvaluationQueriesNamespaceWandb:
		return true
	}
	return false
}

type EvaluationQueryGetResponse struct {
	Definition        EvaluationQueryGetResponseDefinition `json:"definition"          api:"required"`
	EvaluationID      string                               `json:"evaluation_id"       api:"required" format:"uuid"`
	EvaluationQueryID string                               `json:"evaluation_query_id" api:"required" format:"uuid"`
	Name              string                               `json:"name"                api:"required"`
	// Namespace that assigned the scope reference.
	Namespace EvaluationQueryGetResponseNamespace `json:"namespace" api:"required"`
	ScopeRef  string                              `json:"scope_ref" api:"required"`
	JSON      evaluationQueryGetResponseJSON      `json:"-"`
}

// evaluationQueryGetResponseJSON contains the JSON metadata for the struct
// [EvaluationQueryGetResponse]
type evaluationQueryGetResponseJSON struct {
	Definition        apijson.Field
	EvaluationID      apijson.Field
	EvaluationQueryID apijson.Field
	Name              apijson.Field
	Namespace         apijson.Field
	ScopeRef          apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *EvaluationQueryGetResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryGetResponseJSON) RawJSON() string {
	return r.raw
}

type EvaluationQueryGetResponseDefinition struct {
	Metrics []EvaluationQueryGetResponseDefinitionMetric `json:"metrics" api:"required"`
	// Stored definition version; always present in responses.
	V       int64                                         `json:"v"        api:"required"`
	GroupBy []EvaluationQueryGetResponseDefinitionGroupBy `json:"group_by"`
	JSON    evaluationQueryGetResponseDefinitionJSON      `json:"-"`
}

// evaluationQueryGetResponseDefinitionJSON contains the JSON metadata for the
// struct [EvaluationQueryGetResponseDefinition]
type evaluationQueryGetResponseDefinitionJSON struct {
	Metrics     apijson.Field
	V           apijson.Field
	GroupBy     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryGetResponseDefinition) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryGetResponseDefinitionJSON) RawJSON() string {
	return r.raw
}

type EvaluationQueryGetResponseDefinitionMetric struct {
	// Mean accepts integer, number, and boolean (fraction true). Sum accepts integer,
	// number, and boolean (count of true non-null Scores). Min and max accept integer,
	// number, boolean, and string. Strings compare lexicographically. Count accepts
	// these four types and counts non-null latest visible Scores, including zero and
	// false.
	Aggregation     EvaluationQueryGetResponseDefinitionMetricsAggregation `json:"aggregation"       api:"required"`
	Name            string                                                 `json:"name"              api:"required"`
	ScorerVersionID string                                                 `json:"scorer_version_id" api:"required" format:"uuid"`
	JSON            evaluationQueryGetResponseDefinitionMetricJSON         `json:"-"`
}

// evaluationQueryGetResponseDefinitionMetricJSON contains the JSON metadata for
// the struct [EvaluationQueryGetResponseDefinitionMetric]
type evaluationQueryGetResponseDefinitionMetricJSON struct {
	Aggregation     apijson.Field
	Name            apijson.Field
	ScorerVersionID apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *EvaluationQueryGetResponseDefinitionMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryGetResponseDefinitionMetricJSON) RawJSON() string {
	return r.raw
}

// Mean accepts integer, number, and boolean (fraction true). Sum accepts integer,
// number, and boolean (count of true non-null Scores). Min and max accept integer,
// number, boolean, and string. Strings compare lexicographically. Count accepts
// these four types and counts non-null latest visible Scores, including zero and
// false.
type EvaluationQueryGetResponseDefinitionMetricsAggregation string

const (
	EvaluationQueryGetResponseDefinitionMetricsAggregationMean  EvaluationQueryGetResponseDefinitionMetricsAggregation = "mean"
	EvaluationQueryGetResponseDefinitionMetricsAggregationSum   EvaluationQueryGetResponseDefinitionMetricsAggregation = "sum"
	EvaluationQueryGetResponseDefinitionMetricsAggregationMin   EvaluationQueryGetResponseDefinitionMetricsAggregation = "min"
	EvaluationQueryGetResponseDefinitionMetricsAggregationMax   EvaluationQueryGetResponseDefinitionMetricsAggregation = "max"
	EvaluationQueryGetResponseDefinitionMetricsAggregationCount EvaluationQueryGetResponseDefinitionMetricsAggregation = "count"
)

func (r EvaluationQueryGetResponseDefinitionMetricsAggregation) IsKnown() bool {
	switch r {
	case EvaluationQueryGetResponseDefinitionMetricsAggregationMean,
		EvaluationQueryGetResponseDefinitionMetricsAggregationSum,
		EvaluationQueryGetResponseDefinitionMetricsAggregationMin,
		EvaluationQueryGetResponseDefinitionMetricsAggregationMax,
		EvaluationQueryGetResponseDefinitionMetricsAggregationCount:
		return true
	}
	return false
}

type EvaluationQueryGetResponseDefinitionGroupBy struct {
	Kind EvaluationQueryGetResponseDefinitionGroupByKind `json:"kind" api:"required"`
	// Case field name.
	Name string `json:"name"`
	// Case field location.
	Source EvaluationQueryGetResponseDefinitionGroupBySource `json:"source"`
	JSON   evaluationQueryGetResponseDefinitionGroupByJSON   `json:"-"`
	union  EvaluationQueryGetResponseDefinitionGroupByUnion
}

// evaluationQueryGetResponseDefinitionGroupByJSON contains the JSON metadata for
// the struct [EvaluationQueryGetResponseDefinitionGroupBy]
type evaluationQueryGetResponseDefinitionGroupByJSON struct {
	Kind        apijson.Field
	Name        apijson.Field
	Source      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r evaluationQueryGetResponseDefinitionGroupByJSON) RawJSON() string {
	return r.raw
}

func (r *EvaluationQueryGetResponseDefinitionGroupBy) UnmarshalJSON(data []byte) (err error) {
	*r = EvaluationQueryGetResponseDefinitionGroupBy{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [EvaluationQueryGetResponseDefinitionGroupByUnion] interface
// which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumn],
// [EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumn],
// [EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumn].
func (r EvaluationQueryGetResponseDefinitionGroupBy) AsUnion() EvaluationQueryGetResponseDefinitionGroupByUnion {
	return r.union
}

// Union satisfied by
// [EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumn],
// [EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumn] or
// [EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumn].
type EvaluationQueryGetResponseDefinitionGroupByUnion interface {
	implementsEvaluationQueryGetResponseDefinitionGroupBy()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*EvaluationQueryGetResponseDefinitionGroupByUnion)(nil)).Elem(),
		"kind",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumn{}),
			DiscriminatorValue: "subject",
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type: reflect.TypeOf(
				EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumn{},
			),
			DiscriminatorValue: "dataset_field",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumn{}),
			DiscriminatorValue: "score",
		},
	)
}

type EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumn struct {
	Kind EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumnKind `json:"kind" api:"required"`
	JSON evaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumnJSON `json:"-"`
}

// evaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumnJSON contains
// the JSON metadata for the struct
// [EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumn]
type evaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumnJSON struct {
	Kind        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumn) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumnJSON) RawJSON() string {
	return r.raw
}

func (r EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumn) implementsEvaluationQueryGetResponseDefinitionGroupBy() {
}

type EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumnKind string

const (
	EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumnKindSubject EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumnKind = "subject"
)

func (r EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumnKind) IsKnown() bool {
	switch r {
	case EvaluationQueryGetResponseDefinitionGroupBySubjectEvaluationColumnKindSubject:
		return true
	}
	return false
}

type EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumn struct {
	Kind EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnKind `json:"kind" api:"required"`
	// Case field name.
	Name string `json:"name" api:"required"`
	// Case field location.
	Source EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnSource `json:"source" api:"required"`
	JSON   evaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnJSON   `json:"-"`
}

// evaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnJSON
// contains the JSON metadata for the struct
// [EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumn]
type evaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnJSON struct {
	Kind        apijson.Field
	Name        apijson.Field
	Source      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumn) UnmarshalJSON(
	data []byte,
) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnJSON) RawJSON() string {
	return r.raw
}

func (r EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumn) implementsEvaluationQueryGetResponseDefinitionGroupBy() {
}

type EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnKind string

const (
	EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnKindDatasetField EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnKind = "dataset_field"
)

func (r EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnKind) IsKnown() bool {
	switch r {
	case EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnKindDatasetField:
		return true
	}
	return false
}

// Case field location.
type EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnSource string

const (
	EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnSourceInput  EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnSource = "input"
	EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnSourceOutput EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnSource = "output"
)

func (r EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnSource) IsKnown() bool {
	switch r {
	case EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnSourceInput,
		EvaluationQueryGetResponseDefinitionGroupByDatasetFieldEvaluationColumnSourceOutput:
		return true
	}
	return false
}

type EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumn struct {
	Kind EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumnKind `json:"kind" api:"required"`
	// Scorer name selected by this EvaluationVersion.
	Name string                                                               `json:"name" api:"required"`
	JSON evaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumnJSON `json:"-"`
}

// evaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumnJSON contains
// the JSON metadata for the struct
// [EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumn]
type evaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumnJSON struct {
	Kind        apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumn) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumnJSON) RawJSON() string {
	return r.raw
}

func (r EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumn) implementsEvaluationQueryGetResponseDefinitionGroupBy() {
}

type EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumnKind string

const (
	EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumnKindScore EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumnKind = "score"
)

func (r EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumnKind) IsKnown() bool {
	switch r {
	case EvaluationQueryGetResponseDefinitionGroupByScoreEvaluationColumnKindScore:
		return true
	}
	return false
}

type EvaluationQueryGetResponseDefinitionGroupByKind string

const (
	EvaluationQueryGetResponseDefinitionGroupByKindSubject      EvaluationQueryGetResponseDefinitionGroupByKind = "subject"
	EvaluationQueryGetResponseDefinitionGroupByKindDatasetField EvaluationQueryGetResponseDefinitionGroupByKind = "dataset_field"
	EvaluationQueryGetResponseDefinitionGroupByKindScore        EvaluationQueryGetResponseDefinitionGroupByKind = "score"
)

func (r EvaluationQueryGetResponseDefinitionGroupByKind) IsKnown() bool {
	switch r {
	case EvaluationQueryGetResponseDefinitionGroupByKindSubject,
		EvaluationQueryGetResponseDefinitionGroupByKindDatasetField,
		EvaluationQueryGetResponseDefinitionGroupByKindScore:
		return true
	}
	return false
}

// Case field location.
type EvaluationQueryGetResponseDefinitionGroupBySource string

const (
	EvaluationQueryGetResponseDefinitionGroupBySourceInput  EvaluationQueryGetResponseDefinitionGroupBySource = "input"
	EvaluationQueryGetResponseDefinitionGroupBySourceOutput EvaluationQueryGetResponseDefinitionGroupBySource = "output"
)

func (r EvaluationQueryGetResponseDefinitionGroupBySource) IsKnown() bool {
	switch r {
	case EvaluationQueryGetResponseDefinitionGroupBySourceInput,
		EvaluationQueryGetResponseDefinitionGroupBySourceOutput:
		return true
	}
	return false
}

// Namespace that assigned the scope reference.
type EvaluationQueryGetResponseNamespace string

const (
	EvaluationQueryGetResponseNamespaceWandb EvaluationQueryGetResponseNamespace = "wandb"
)

func (r EvaluationQueryGetResponseNamespace) IsKnown() bool {
	switch r {
	case EvaluationQueryGetResponseNamespaceWandb:
		return true
	}
	return false
}

type EvaluationQueryNewParams struct {
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
	ScopeID    param.Field[string]                             `path:"scope_id" api:"required"`
	Definition param.Field[EvaluationQueryNewParamsDefinition] `                api:"required" json:"definition"`
	Name       param.Field[string]                             `                api:"required" json:"name"`
}

func (r EvaluationQueryNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type EvaluationQueryNewParamsDefinition struct {
	// Metrics with unique names, using ScorerVersions registered on the Evaluation.
	Metrics param.Field[[]EvaluationQueryNewParamsDefinitionMetric] `json:"metrics" api:"required"`
	// Groups Cases by distinct combinations of chosen Dataset fields or Scores.
	// Results list values in this order. Omit, use null, or provide an empty list to
	// calculate metrics across all Cases.
	GroupBy param.Field[[]EvaluationQueryNewParamsDefinitionGroupByUnion] `json:"group_by"`
	// Definition version; defaults to 2 when omitted.
	V param.Field[int64] `json:"v"`
}

func (r EvaluationQueryNewParamsDefinition) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type EvaluationQueryNewParamsDefinitionMetric struct {
	// Mean accepts integer, number, and boolean (fraction true). Sum accepts integer,
	// number, and boolean (count of true non-null Scores). Min and max accept integer,
	// number, boolean, and string. Strings compare lexicographically. Count accepts
	// these four types and counts non-null latest visible Scores, including zero and
	// false.
	Aggregation     param.Field[EvaluationQueryNewParamsDefinitionMetricsAggregation] `json:"aggregation"       api:"required"`
	Name            param.Field[string]                                               `json:"name"              api:"required"`
	ScorerVersionID param.Field[string]                                               `json:"scorer_version_id" api:"required" format:"uuid"`
}

func (r EvaluationQueryNewParamsDefinitionMetric) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Mean accepts integer, number, and boolean (fraction true). Sum accepts integer,
// number, and boolean (count of true non-null Scores). Min and max accept integer,
// number, boolean, and string. Strings compare lexicographically. Count accepts
// these four types and counts non-null latest visible Scores, including zero and
// false.
type EvaluationQueryNewParamsDefinitionMetricsAggregation string

const (
	EvaluationQueryNewParamsDefinitionMetricsAggregationMean  EvaluationQueryNewParamsDefinitionMetricsAggregation = "mean"
	EvaluationQueryNewParamsDefinitionMetricsAggregationSum   EvaluationQueryNewParamsDefinitionMetricsAggregation = "sum"
	EvaluationQueryNewParamsDefinitionMetricsAggregationMin   EvaluationQueryNewParamsDefinitionMetricsAggregation = "min"
	EvaluationQueryNewParamsDefinitionMetricsAggregationMax   EvaluationQueryNewParamsDefinitionMetricsAggregation = "max"
	EvaluationQueryNewParamsDefinitionMetricsAggregationCount EvaluationQueryNewParamsDefinitionMetricsAggregation = "count"
)

func (r EvaluationQueryNewParamsDefinitionMetricsAggregation) IsKnown() bool {
	switch r {
	case EvaluationQueryNewParamsDefinitionMetricsAggregationMean,
		EvaluationQueryNewParamsDefinitionMetricsAggregationSum,
		EvaluationQueryNewParamsDefinitionMetricsAggregationMin,
		EvaluationQueryNewParamsDefinitionMetricsAggregationMax,
		EvaluationQueryNewParamsDefinitionMetricsAggregationCount:
		return true
	}
	return false
}

type EvaluationQueryNewParamsDefinitionGroupBy struct {
	Kind param.Field[EvaluationQueryNewParamsDefinitionGroupByKind] `json:"kind" api:"required"`
	// Case field name.
	Name param.Field[string] `json:"name"`
	// Case field location.
	Source param.Field[EvaluationQueryNewParamsDefinitionGroupBySource] `json:"source"`
}

func (r EvaluationQueryNewParamsDefinitionGroupBy) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r EvaluationQueryNewParamsDefinitionGroupBy) implementsEvaluationQueryNewParamsDefinitionGroupByUnion() {
}

// Satisfied by [EvaluationQueryNewParamsDefinitionGroupBySubjectEvaluationColumn],
// [EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumn],
// [EvaluationQueryNewParamsDefinitionGroupByScoreEvaluationColumn],
// [EvaluationQueryNewParamsDefinitionGroupBy].
type EvaluationQueryNewParamsDefinitionGroupByUnion interface {
	implementsEvaluationQueryNewParamsDefinitionGroupByUnion()
}

type EvaluationQueryNewParamsDefinitionGroupBySubjectEvaluationColumn struct {
	Kind param.Field[EvaluationQueryNewParamsDefinitionGroupBySubjectEvaluationColumnKind] `json:"kind" api:"required"`
}

func (r EvaluationQueryNewParamsDefinitionGroupBySubjectEvaluationColumn) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r EvaluationQueryNewParamsDefinitionGroupBySubjectEvaluationColumn) implementsEvaluationQueryNewParamsDefinitionGroupByUnion() {
}

type EvaluationQueryNewParamsDefinitionGroupBySubjectEvaluationColumnKind string

const (
	EvaluationQueryNewParamsDefinitionGroupBySubjectEvaluationColumnKindSubject EvaluationQueryNewParamsDefinitionGroupBySubjectEvaluationColumnKind = "subject"
)

func (r EvaluationQueryNewParamsDefinitionGroupBySubjectEvaluationColumnKind) IsKnown() bool {
	switch r {
	case EvaluationQueryNewParamsDefinitionGroupBySubjectEvaluationColumnKindSubject:
		return true
	}
	return false
}

type EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumn struct {
	Kind param.Field[EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnKind] `json:"kind" api:"required"`
	// Case field name.
	Name param.Field[string] `json:"name" api:"required"`
	// Case field location.
	Source param.Field[EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnSource] `json:"source" api:"required"`
}

func (r EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumn) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumn) implementsEvaluationQueryNewParamsDefinitionGroupByUnion() {
}

type EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnKind string

const (
	EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnKindDatasetField EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnKind = "dataset_field"
)

func (r EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnKind) IsKnown() bool {
	switch r {
	case EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnKindDatasetField:
		return true
	}
	return false
}

// Case field location.
type EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnSource string

const (
	EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnSourceInput  EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnSource = "input"
	EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnSourceOutput EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnSource = "output"
)

func (r EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnSource) IsKnown() bool {
	switch r {
	case EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnSourceInput,
		EvaluationQueryNewParamsDefinitionGroupByDatasetFieldEvaluationColumnSourceOutput:
		return true
	}
	return false
}

type EvaluationQueryNewParamsDefinitionGroupByScoreEvaluationColumn struct {
	Kind param.Field[EvaluationQueryNewParamsDefinitionGroupByScoreEvaluationColumnKind] `json:"kind" api:"required"`
	// Scorer name selected by this EvaluationVersion.
	Name param.Field[string] `json:"name" api:"required"`
}

func (r EvaluationQueryNewParamsDefinitionGroupByScoreEvaluationColumn) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r EvaluationQueryNewParamsDefinitionGroupByScoreEvaluationColumn) implementsEvaluationQueryNewParamsDefinitionGroupByUnion() {
}

type EvaluationQueryNewParamsDefinitionGroupByScoreEvaluationColumnKind string

const (
	EvaluationQueryNewParamsDefinitionGroupByScoreEvaluationColumnKindScore EvaluationQueryNewParamsDefinitionGroupByScoreEvaluationColumnKind = "score"
)

func (r EvaluationQueryNewParamsDefinitionGroupByScoreEvaluationColumnKind) IsKnown() bool {
	switch r {
	case EvaluationQueryNewParamsDefinitionGroupByScoreEvaluationColumnKindScore:
		return true
	}
	return false
}

type EvaluationQueryNewParamsDefinitionGroupByKind string

const (
	EvaluationQueryNewParamsDefinitionGroupByKindSubject      EvaluationQueryNewParamsDefinitionGroupByKind = "subject"
	EvaluationQueryNewParamsDefinitionGroupByKindDatasetField EvaluationQueryNewParamsDefinitionGroupByKind = "dataset_field"
	EvaluationQueryNewParamsDefinitionGroupByKindScore        EvaluationQueryNewParamsDefinitionGroupByKind = "score"
)

func (r EvaluationQueryNewParamsDefinitionGroupByKind) IsKnown() bool {
	switch r {
	case EvaluationQueryNewParamsDefinitionGroupByKindSubject,
		EvaluationQueryNewParamsDefinitionGroupByKindDatasetField,
		EvaluationQueryNewParamsDefinitionGroupByKindScore:
		return true
	}
	return false
}

// Case field location.
type EvaluationQueryNewParamsDefinitionGroupBySource string

const (
	EvaluationQueryNewParamsDefinitionGroupBySourceInput  EvaluationQueryNewParamsDefinitionGroupBySource = "input"
	EvaluationQueryNewParamsDefinitionGroupBySourceOutput EvaluationQueryNewParamsDefinitionGroupBySource = "output"
)

func (r EvaluationQueryNewParamsDefinitionGroupBySource) IsKnown() bool {
	switch r {
	case EvaluationQueryNewParamsDefinitionGroupBySourceInput, EvaluationQueryNewParamsDefinitionGroupBySourceOutput:
		return true
	}
	return false
}

type EvaluationQueryListParams struct {
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
	Cursor  param.Field[string] `                               query:"cursor"`
	Limit   param.Field[int64]  `                               query:"limit"`
}

// URLQuery serializes [EvaluationQueryListParams]'s query parameters as
// `url.Values`.
func (r EvaluationQueryListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type EvaluationQueryGetParams struct {
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
}
