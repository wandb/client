// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package coreweaveevaluations

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/wandb/core/services/evaluations/generated/go/internal/apijson"
	"github.com/wandb/core/services/evaluations/generated/go/internal/param"
	"github.com/wandb/core/services/evaluations/generated/go/internal/requestconfig"
	"github.com/wandb/core/services/evaluations/generated/go/option"
)

// EvaluationVersionService contains methods and other services that help with
// interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvaluationVersionService] method instead.
type EvaluationVersionService struct {
	Options []option.RequestOption
	Results *EvaluationVersionResultService
	Queries *EvaluationVersionQueryService
	Rows    *EvaluationVersionRowService
}

// NewEvaluationVersionService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewEvaluationVersionService(opts ...option.RequestOption) (r *EvaluationVersionService) {
	r = &EvaluationVersionService{}
	r.Options = opts
	r.Results = NewEvaluationVersionResultService(opts...)
	r.Queries = NewEvaluationVersionQueryService(opts...)
	r.Rows = NewEvaluationVersionRowService(opts...)
	return
}

// Creates an immutable snapshot of one DatasetVersion and the Evaluation's
// currently registered ScorerVersions.
//
// Request bodies must be smaller than 1 MiB; larger bodies are rejected with 413
// Request Entity Too Large.
func (r *EvaluationVersionService) New(
	ctx context.Context,
	evaluationID string,
	params EvaluationVersionNewParams,
	opts ...option.RequestOption,
) (res *EvaluationVersionNewResponse, err error) {
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
		"v1/namespaces/%s/scopes/%s/evaluations/%s/versions",
		params.Namespace,
		params.ScopeID,
		evaluationID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Returns the immutable EvaluationVersion snapshot exactly as created.
func (r *EvaluationVersionService) Get(
	ctx context.Context,
	evaluationVersionID string,
	query EvaluationVersionGetParams,
	opts ...option.RequestOption,
) (res *EvaluationVersionGetResponse, err error) {
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
	if evaluationVersionID == "" {
		err = errors.New("missing required evaluation_version_id parameter")
		return nil, err
	}
	path := fmt.Sprintf(
		"v1/namespaces/%s/scopes/%s/evaluation-versions/%s",
		query.Namespace,
		query.ScopeID,
		evaluationVersionID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type EvaluationVersionNewResponse struct {
	DatasetID           string `json:"dataset_id"            api:"required" format:"uuid"`
	DatasetVersionID    string `json:"dataset_version_id"    api:"required" format:"uuid"`
	EvaluationID        string `json:"evaluation_id"         api:"required" format:"uuid"`
	EvaluationVersionID string `json:"evaluation_version_id" api:"required" format:"uuid"`
	// Scope publication watermark captured when the snapshot was created.
	LastCommitID int64 `json:"last_commit_id" api:"required"`
	// Namespace that assigned the scope reference.
	Namespace         EvaluationVersionNewResponseNamespace         `json:"namespace"           api:"required"`
	PrimaryTargetKind EvaluationVersionNewResponsePrimaryTargetKind `json:"primary_target_kind" api:"required"`
	ScopeRef          string                                        `json:"scope_ref"           api:"required"`
	// Selected ScorerVersions in snapshot order.
	Scorers []EvaluationVersionNewResponseScorer `json:"scorers" api:"required"`
	JSON    evaluationVersionNewResponseJSON     `json:"-"`
}

// evaluationVersionNewResponseJSON contains the JSON metadata for the struct
// [EvaluationVersionNewResponse]
type evaluationVersionNewResponseJSON struct {
	DatasetID           apijson.Field
	DatasetVersionID    apijson.Field
	EvaluationID        apijson.Field
	EvaluationVersionID apijson.Field
	LastCommitID        apijson.Field
	Namespace           apijson.Field
	PrimaryTargetKind   apijson.Field
	ScopeRef            apijson.Field
	Scorers             apijson.Field
	raw                 string
	ExtraFields         map[string]apijson.Field
}

func (r *EvaluationVersionNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionNewResponseJSON) RawJSON() string {
	return r.raw
}

// Namespace that assigned the scope reference.
type EvaluationVersionNewResponseNamespace string

const (
	EvaluationVersionNewResponseNamespaceWandb EvaluationVersionNewResponseNamespace = "wandb"
)

func (r EvaluationVersionNewResponseNamespace) IsKnown() bool {
	switch r {
	case EvaluationVersionNewResponseNamespaceWandb:
		return true
	}
	return false
}

type EvaluationVersionNewResponsePrimaryTargetKind string

const (
	EvaluationVersionNewResponsePrimaryTargetKindCase       EvaluationVersionNewResponsePrimaryTargetKind = "case"
	EvaluationVersionNewResponsePrimaryTargetKindTrajectory EvaluationVersionNewResponsePrimaryTargetKind = "trajectory"
)

func (r EvaluationVersionNewResponsePrimaryTargetKind) IsKnown() bool {
	switch r {
	case EvaluationVersionNewResponsePrimaryTargetKindCase, EvaluationVersionNewResponsePrimaryTargetKindTrajectory:
		return true
	}
	return false
}

type EvaluationVersionNewResponseScorer struct {
	ScorerID  string                                 `json:"scorer_id"  api:"required" format:"uuid"`
	VersionID string                                 `json:"version_id" api:"required" format:"uuid"`
	JSON      evaluationVersionNewResponseScorerJSON `json:"-"`
}

// evaluationVersionNewResponseScorerJSON contains the JSON metadata for the struct
// [EvaluationVersionNewResponseScorer]
type evaluationVersionNewResponseScorerJSON struct {
	ScorerID    apijson.Field
	VersionID   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationVersionNewResponseScorer) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionNewResponseScorerJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionGetResponse struct {
	DatasetID           string `json:"dataset_id"            api:"required" format:"uuid"`
	DatasetVersionID    string `json:"dataset_version_id"    api:"required" format:"uuid"`
	EvaluationID        string `json:"evaluation_id"         api:"required" format:"uuid"`
	EvaluationVersionID string `json:"evaluation_version_id" api:"required" format:"uuid"`
	// Scope publication watermark captured when the snapshot was created.
	LastCommitID int64 `json:"last_commit_id" api:"required"`
	// Namespace that assigned the scope reference.
	Namespace         EvaluationVersionGetResponseNamespace         `json:"namespace"           api:"required"`
	PrimaryTargetKind EvaluationVersionGetResponsePrimaryTargetKind `json:"primary_target_kind" api:"required"`
	ScopeRef          string                                        `json:"scope_ref"           api:"required"`
	// Selected ScorerVersions in snapshot order.
	Scorers []EvaluationVersionGetResponseScorer `json:"scorers" api:"required"`
	JSON    evaluationVersionGetResponseJSON     `json:"-"`
}

// evaluationVersionGetResponseJSON contains the JSON metadata for the struct
// [EvaluationVersionGetResponse]
type evaluationVersionGetResponseJSON struct {
	DatasetID           apijson.Field
	DatasetVersionID    apijson.Field
	EvaluationID        apijson.Field
	EvaluationVersionID apijson.Field
	LastCommitID        apijson.Field
	Namespace           apijson.Field
	PrimaryTargetKind   apijson.Field
	ScopeRef            apijson.Field
	Scorers             apijson.Field
	raw                 string
	ExtraFields         map[string]apijson.Field
}

func (r *EvaluationVersionGetResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionGetResponseJSON) RawJSON() string {
	return r.raw
}

// Namespace that assigned the scope reference.
type EvaluationVersionGetResponseNamespace string

const (
	EvaluationVersionGetResponseNamespaceWandb EvaluationVersionGetResponseNamespace = "wandb"
)

func (r EvaluationVersionGetResponseNamespace) IsKnown() bool {
	switch r {
	case EvaluationVersionGetResponseNamespaceWandb:
		return true
	}
	return false
}

type EvaluationVersionGetResponsePrimaryTargetKind string

const (
	EvaluationVersionGetResponsePrimaryTargetKindCase       EvaluationVersionGetResponsePrimaryTargetKind = "case"
	EvaluationVersionGetResponsePrimaryTargetKindTrajectory EvaluationVersionGetResponsePrimaryTargetKind = "trajectory"
)

func (r EvaluationVersionGetResponsePrimaryTargetKind) IsKnown() bool {
	switch r {
	case EvaluationVersionGetResponsePrimaryTargetKindCase, EvaluationVersionGetResponsePrimaryTargetKindTrajectory:
		return true
	}
	return false
}

type EvaluationVersionGetResponseScorer struct {
	ScorerID  string                                 `json:"scorer_id"  api:"required" format:"uuid"`
	VersionID string                                 `json:"version_id" api:"required" format:"uuid"`
	JSON      evaluationVersionGetResponseScorerJSON `json:"-"`
}

// evaluationVersionGetResponseScorerJSON contains the JSON metadata for the struct
// [EvaluationVersionGetResponseScorer]
type evaluationVersionGetResponseScorerJSON struct {
	ScorerID    apijson.Field
	VersionID   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationVersionGetResponseScorer) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationVersionGetResponseScorerJSON) RawJSON() string {
	return r.raw
}

type EvaluationVersionNewParams struct {
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
	// Fixed Evaluations require their registered DatasetVersion ID or omission. Live
	// Evaluations default to the current Dataset head.
	DatasetVersionID param.Field[string] `json:"dataset_version_id" format:"uuid"`
}

func (r EvaluationVersionNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type EvaluationVersionGetParams struct {
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
