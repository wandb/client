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

// EvaluationVersionResultService contains methods and other services that help
// with interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvaluationVersionResultService] method instead.
type EvaluationVersionResultService struct {
	Options []option.RequestOption
}

// NewEvaluationVersionResultService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewEvaluationVersionResultService(opts ...option.RequestOption) (r *EvaluationVersionResultService) {
	r = &EvaluationVersionResultService{}
	r.Options = opts
	return
}

// Pages the pinned DatasetVersion's Cases in position order at its frozen cutoff;
// Cases appended later never appear. Each target includes the latest Score visible
// at the EvaluationVersion's frozen commit per pinned ScorerVersion, in scorer
// order; unscored targets have an empty scores array. published_commit_id reports
// that immutable cutoff. Cursors track target position only because the
// EvaluationVersion fixes Score visibility. next_cursor is absent on the last
// page. Trajectory targets return an empty page until the Trajectory resource is
// implemented.
func (r *EvaluationVersionResultService) List(
	ctx context.Context,
	evaluationVersionID string,
	params EvaluationVersionResultListParams,
	opts ...option.RequestOption,
) (res *EvaluationVersionResultListResponse, err error) {
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
		"v1/namespaces/%s/scopes/%s/evaluation-versions/%s/results",
		params.Namespace,
		params.ScopeID,
		evaluationVersionID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

type EvaluationVersionResultListResponse struct {
	PublishedCommitID int64                                       `json:"published_commit_id" api:"required"`
	Targets           []EvaluationVersionResultListResponseTarget `json:"targets"             api:"required"`
	NextCursor        string                                      `json:"next_cursor"`
	JSON              evaluationVersionResultListResponseJSON     `json:"-"`
}

// evaluationVersionResultListResponseJSON contains the JSON metadata for the
// struct [EvaluationVersionResultListResponse]
type evaluationVersionResultListResponseJSON struct {
	PublishedCommitID apijson.Field
	Targets           apijson.Field
	NextCursor        apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *EvaluationVersionResultListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionResultListResponseJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionResultListResponseTarget struct {
	Scores []EvaluationVersionResultListResponseTargetsScore `json:"scores" api:"required"`
	Target EvaluationVersionResultListResponseTargetsTarget  `json:"target" api:"required"`
	JSON   evaluationVersionResultListResponseTargetJSON     `json:"-"`
}

// evaluationVersionResultListResponseTargetJSON contains the JSON metadata for the
// struct [EvaluationVersionResultListResponseTarget]
type evaluationVersionResultListResponseTargetJSON struct {
	Scores      apijson.Field
	Target      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationVersionResultListResponseTarget) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionResultListResponseTargetJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionResultListResponseTargetsScore struct {
	CommitID        int64                                                     `json:"commit_id"         api:"required"`
	ScoreID         string                                                    `json:"score_id"          api:"required" format:"uuid"`
	ScorerID        string                                                    `json:"scorer_id"         api:"required" format:"uuid"`
	ScorerVersionID string                                                    `json:"scorer_version_id" api:"required" format:"uuid"`
	Value           interface{}                                               `json:"value"             api:"required"`
	ValueType       EvaluationVersionResultListResponseTargetsScoresValueType `json:"value_type"        api:"required"`
	JSON            evaluationVersionResultListResponseTargetsScoreJSON       `json:"-"`
}

// evaluationVersionResultListResponseTargetsScoreJSON contains the JSON metadata
// for the struct [EvaluationVersionResultListResponseTargetsScore]
type evaluationVersionResultListResponseTargetsScoreJSON struct {
	CommitID        apijson.Field
	ScoreID         apijson.Field
	ScorerID        apijson.Field
	ScorerVersionID apijson.Field
	Value           apijson.Field
	ValueType       apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *EvaluationVersionResultListResponseTargetsScore) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionResultListResponseTargetsScoreJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionResultListResponseTargetsScoresValueType string

const (
	EvaluationVersionResultListResponseTargetsScoresValueTypeNull    EvaluationVersionResultListResponseTargetsScoresValueType = "null"
	EvaluationVersionResultListResponseTargetsScoresValueTypeBoolean EvaluationVersionResultListResponseTargetsScoresValueType = "boolean"
	EvaluationVersionResultListResponseTargetsScoresValueTypeInteger EvaluationVersionResultListResponseTargetsScoresValueType = "integer"
	EvaluationVersionResultListResponseTargetsScoresValueTypeNumber  EvaluationVersionResultListResponseTargetsScoresValueType = "number"
	EvaluationVersionResultListResponseTargetsScoresValueTypeString  EvaluationVersionResultListResponseTargetsScoresValueType = "string"
	EvaluationVersionResultListResponseTargetsScoresValueTypeJson    EvaluationVersionResultListResponseTargetsScoresValueType = "json"
)

func (r EvaluationVersionResultListResponseTargetsScoresValueType) IsKnown() bool {
	switch r {
	case EvaluationVersionResultListResponseTargetsScoresValueTypeNull,
		EvaluationVersionResultListResponseTargetsScoresValueTypeBoolean,
		EvaluationVersionResultListResponseTargetsScoresValueTypeInteger,
		EvaluationVersionResultListResponseTargetsScoresValueTypeNumber,
		EvaluationVersionResultListResponseTargetsScoresValueTypeString,
		EvaluationVersionResultListResponseTargetsScoresValueTypeJson:
		return true
	}
	return false
}

type EvaluationVersionResultListResponseTargetsTarget struct {
	ID   string                                               `json:"id"   api:"required" format:"uuid"`
	Kind EvaluationVersionResultListResponseTargetsTargetKind `json:"kind" api:"required"`
	JSON evaluationVersionResultListResponseTargetsTargetJSON `json:"-"`
}

// evaluationVersionResultListResponseTargetsTargetJSON contains the JSON metadata
// for the struct [EvaluationVersionResultListResponseTargetsTarget]
type evaluationVersionResultListResponseTargetsTargetJSON struct {
	ID          apijson.Field
	Kind        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationVersionResultListResponseTargetsTarget) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionResultListResponseTargetsTargetJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionResultListResponseTargetsTargetKind string

const (
	EvaluationVersionResultListResponseTargetsTargetKindCase       EvaluationVersionResultListResponseTargetsTargetKind = "case"
	EvaluationVersionResultListResponseTargetsTargetKindComparison EvaluationVersionResultListResponseTargetsTargetKind = "comparison"
	EvaluationVersionResultListResponseTargetsTargetKindTrajectory EvaluationVersionResultListResponseTargetsTargetKind = "trajectory"
)

func (r EvaluationVersionResultListResponseTargetsTargetKind) IsKnown() bool {
	switch r {
	case EvaluationVersionResultListResponseTargetsTargetKindCase,
		EvaluationVersionResultListResponseTargetsTargetKindComparison,
		EvaluationVersionResultListResponseTargetsTargetKindTrajectory:
		return true
	}
	return false
}

type EvaluationVersionResultListParams struct {
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

// URLQuery serializes [EvaluationVersionResultListParams]'s query parameters as
// `url.Values`.
func (r EvaluationVersionResultListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
