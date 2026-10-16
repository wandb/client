// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package coreweaveevaluations

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"slices"

	"github.com/tidwall/gjson"

	"github.com/wandb/core/services/evaluations/generated/go/internal/apijson"
	"github.com/wandb/core/services/evaluations/generated/go/internal/param"
	"github.com/wandb/core/services/evaluations/generated/go/internal/requestconfig"
	"github.com/wandb/core/services/evaluations/generated/go/option"
)

// EvaluationService contains methods and other services that help with interacting
// with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvaluationService] method instead.
type EvaluationService struct {
	Options  []option.RequestOption
	Queries  *EvaluationQueryService
	Versions *EvaluationVersionService
}

// NewEvaluationService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewEvaluationService(opts ...option.RequestOption) (r *EvaluationService) {
	r = &EvaluationService{}
	r.Options = opts
	r.Queries = NewEvaluationQueryService(opts...)
	r.Versions = NewEvaluationVersionService(opts...)
	return
}

// Creates an immutable definition with a fixed DatasetVersion source or a live
// Dataset source. The source contains its kind and exactly one corresponding ID. A
// source outside the scope returns 404. Attributes are stored as canonical JSON;
// an exact idempotent retry returns the original response.
//
// Request bodies must be smaller than 1 MiB; larger bodies are rejected with 413
// Request Entity Too Large.
func (r *EvaluationService) New(
	ctx context.Context,
	params EvaluationNewParams,
	opts ...option.RequestOption,
) (res *EvaluationNewResponse, err error) {
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
	path := fmt.Sprintf("v1/namespaces/%s/scopes/%s/evaluations", params.Namespace, params.ScopeID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type EvaluationNewResponse struct {
	Attributes   map[string]interface{} `json:"attributes"    api:"required"`
	EvaluationID string                 `json:"evaluation_id" api:"required" format:"uuid"`
	Name         string                 `json:"name"          api:"required"`
	// Namespace that assigned the scope reference.
	Namespace         EvaluationNewResponseNamespace         `json:"namespace"           api:"required"`
	PrimaryTargetKind EvaluationNewResponsePrimaryTargetKind `json:"primary_target_kind" api:"required"`
	ScopeRef          string                                 `json:"scope_ref"           api:"required"`
	Source            EvaluationNewResponseSource            `json:"source"              api:"required"`
	JSON              evaluationNewResponseJSON              `json:"-"`
}

// evaluationNewResponseJSON contains the JSON metadata for the struct
// [EvaluationNewResponse]
type evaluationNewResponseJSON struct {
	Attributes        apijson.Field
	EvaluationID      apijson.Field
	Name              apijson.Field
	Namespace         apijson.Field
	PrimaryTargetKind apijson.Field
	ScopeRef          apijson.Field
	Source            apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *EvaluationNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationNewResponseJSON) RawJSON() string {
	return r.raw
}

// Namespace that assigned the scope reference.
type EvaluationNewResponseNamespace string

const (
	EvaluationNewResponseNamespaceWandb EvaluationNewResponseNamespace = "wandb"
)

func (r EvaluationNewResponseNamespace) IsKnown() bool {
	switch r {
	case EvaluationNewResponseNamespaceWandb:
		return true
	}
	return false
}

type EvaluationNewResponsePrimaryTargetKind string

const (
	EvaluationNewResponsePrimaryTargetKindCase       EvaluationNewResponsePrimaryTargetKind = "case"
	EvaluationNewResponsePrimaryTargetKindTrajectory EvaluationNewResponsePrimaryTargetKind = "trajectory"
)

func (r EvaluationNewResponsePrimaryTargetKind) IsKnown() bool {
	switch r {
	case EvaluationNewResponsePrimaryTargetKindCase, EvaluationNewResponsePrimaryTargetKindTrajectory:
		return true
	}
	return false
}

type EvaluationNewResponseSource struct {
	Kind             EvaluationNewResponseSourceKind `json:"kind"               api:"required"`
	DatasetID        string                          `json:"dataset_id"                        format:"uuid"`
	DatasetVersionID string                          `json:"dataset_version_id"                format:"uuid"`
	JSON             evaluationNewResponseSourceJSON `json:"-"`
	union            EvaluationNewResponseSourceUnion
}

// evaluationNewResponseSourceJSON contains the JSON metadata for the struct
// [EvaluationNewResponseSource]
type evaluationNewResponseSourceJSON struct {
	Kind             apijson.Field
	DatasetID        apijson.Field
	DatasetVersionID apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r evaluationNewResponseSourceJSON) RawJSON() string {
	return r.raw
}

func (r *EvaluationNewResponseSource) UnmarshalJSON(data []byte) (err error) {
	*r = EvaluationNewResponseSource{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [EvaluationNewResponseSourceUnion] interface which you can
// cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [EvaluationNewResponseSourceFixedEvaluationSource],
// [EvaluationNewResponseSourceLiveEvaluationSource].
func (r EvaluationNewResponseSource) AsUnion() EvaluationNewResponseSourceUnion {
	return r.union
}

// Union satisfied by [EvaluationNewResponseSourceFixedEvaluationSource] or
// [EvaluationNewResponseSourceLiveEvaluationSource].
type EvaluationNewResponseSourceUnion interface {
	implementsEvaluationNewResponseSource()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*EvaluationNewResponseSourceUnion)(nil)).Elem(),
		"kind",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(EvaluationNewResponseSourceFixedEvaluationSource{}),
			DiscriminatorValue: "fixed",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(EvaluationNewResponseSourceLiveEvaluationSource{}),
			DiscriminatorValue: "live",
		},
	)
}

type EvaluationNewResponseSourceFixedEvaluationSource struct {
	DatasetVersionID string                                               `json:"dataset_version_id" api:"required" format:"uuid"`
	Kind             EvaluationNewResponseSourceFixedEvaluationSourceKind `json:"kind"               api:"required"`
	JSON             evaluationNewResponseSourceFixedEvaluationSourceJSON `json:"-"`
}

// evaluationNewResponseSourceFixedEvaluationSourceJSON contains the JSON metadata
// for the struct [EvaluationNewResponseSourceFixedEvaluationSource]
type evaluationNewResponseSourceFixedEvaluationSourceJSON struct {
	DatasetVersionID apijson.Field
	Kind             apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *EvaluationNewResponseSourceFixedEvaluationSource) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationNewResponseSourceFixedEvaluationSourceJSON) RawJSON() string {
	return r.raw
}

func (r EvaluationNewResponseSourceFixedEvaluationSource) implementsEvaluationNewResponseSource() {}

type EvaluationNewResponseSourceFixedEvaluationSourceKind string

const (
	EvaluationNewResponseSourceFixedEvaluationSourceKindFixed EvaluationNewResponseSourceFixedEvaluationSourceKind = "fixed"
)

func (r EvaluationNewResponseSourceFixedEvaluationSourceKind) IsKnown() bool {
	switch r {
	case EvaluationNewResponseSourceFixedEvaluationSourceKindFixed:
		return true
	}
	return false
}

type EvaluationNewResponseSourceLiveEvaluationSource struct {
	DatasetID string                                              `json:"dataset_id" api:"required" format:"uuid"`
	Kind      EvaluationNewResponseSourceLiveEvaluationSourceKind `json:"kind"       api:"required"`
	JSON      evaluationNewResponseSourceLiveEvaluationSourceJSON `json:"-"`
}

// evaluationNewResponseSourceLiveEvaluationSourceJSON contains the JSON metadata
// for the struct [EvaluationNewResponseSourceLiveEvaluationSource]
type evaluationNewResponseSourceLiveEvaluationSourceJSON struct {
	DatasetID   apijson.Field
	Kind        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluationNewResponseSourceLiveEvaluationSource) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluationNewResponseSourceLiveEvaluationSourceJSON) RawJSON() string {
	return r.raw
}

func (r EvaluationNewResponseSourceLiveEvaluationSource) implementsEvaluationNewResponseSource() {}

type EvaluationNewResponseSourceLiveEvaluationSourceKind string

const (
	EvaluationNewResponseSourceLiveEvaluationSourceKindLive EvaluationNewResponseSourceLiveEvaluationSourceKind = "live"
)

func (r EvaluationNewResponseSourceLiveEvaluationSourceKind) IsKnown() bool {
	switch r {
	case EvaluationNewResponseSourceLiveEvaluationSourceKindLive:
		return true
	}
	return false
}

type EvaluationNewResponseSourceKind string

const (
	EvaluationNewResponseSourceKindFixed EvaluationNewResponseSourceKind = "fixed"
	EvaluationNewResponseSourceKindLive  EvaluationNewResponseSourceKind = "live"
)

func (r EvaluationNewResponseSourceKind) IsKnown() bool {
	switch r {
	case EvaluationNewResponseSourceKindFixed, EvaluationNewResponseSourceKindLive:
		return true
	}
	return false
}

type EvaluationNewParams struct {
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
	ScopeID           param.Field[string]                               `path:"scope_id" api:"required"`
	Name              param.Field[string]                               `                api:"required" json:"name"`
	PrimaryTargetKind param.Field[EvaluationNewParamsPrimaryTargetKind] `                api:"required" json:"primary_target_kind"`
	Source            param.Field[EvaluationNewParamsSourceUnion]       `                api:"required" json:"source"`
	Attributes        param.Field[map[string]interface{}]               `                               json:"attributes"`
}

func (r EvaluationNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type EvaluationNewParamsPrimaryTargetKind string

const (
	EvaluationNewParamsPrimaryTargetKindCase       EvaluationNewParamsPrimaryTargetKind = "case"
	EvaluationNewParamsPrimaryTargetKindTrajectory EvaluationNewParamsPrimaryTargetKind = "trajectory"
)

func (r EvaluationNewParamsPrimaryTargetKind) IsKnown() bool {
	switch r {
	case EvaluationNewParamsPrimaryTargetKindCase, EvaluationNewParamsPrimaryTargetKindTrajectory:
		return true
	}
	return false
}

type EvaluationNewParamsSource struct {
	Kind             param.Field[EvaluationNewParamsSourceKind] `json:"kind"               api:"required"`
	DatasetID        param.Field[string]                        `json:"dataset_id"                        format:"uuid"`
	DatasetVersionID param.Field[string]                        `json:"dataset_version_id"                format:"uuid"`
}

func (r EvaluationNewParamsSource) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r EvaluationNewParamsSource) implementsEvaluationNewParamsSourceUnion() {}

// Satisfied by [EvaluationNewParamsSourceFixedEvaluationSource],
// [EvaluationNewParamsSourceLiveEvaluationSource], [EvaluationNewParamsSource].
type EvaluationNewParamsSourceUnion interface {
	implementsEvaluationNewParamsSourceUnion()
}

type EvaluationNewParamsSourceFixedEvaluationSource struct {
	DatasetVersionID param.Field[string]                                             `json:"dataset_version_id" api:"required" format:"uuid"`
	Kind             param.Field[EvaluationNewParamsSourceFixedEvaluationSourceKind] `json:"kind"               api:"required"`
}

func (r EvaluationNewParamsSourceFixedEvaluationSource) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r EvaluationNewParamsSourceFixedEvaluationSource) implementsEvaluationNewParamsSourceUnion() {}

type EvaluationNewParamsSourceFixedEvaluationSourceKind string

const (
	EvaluationNewParamsSourceFixedEvaluationSourceKindFixed EvaluationNewParamsSourceFixedEvaluationSourceKind = "fixed"
)

func (r EvaluationNewParamsSourceFixedEvaluationSourceKind) IsKnown() bool {
	switch r {
	case EvaluationNewParamsSourceFixedEvaluationSourceKindFixed:
		return true
	}
	return false
}

type EvaluationNewParamsSourceLiveEvaluationSource struct {
	DatasetID param.Field[string]                                            `json:"dataset_id" api:"required" format:"uuid"`
	Kind      param.Field[EvaluationNewParamsSourceLiveEvaluationSourceKind] `json:"kind"       api:"required"`
}

func (r EvaluationNewParamsSourceLiveEvaluationSource) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r EvaluationNewParamsSourceLiveEvaluationSource) implementsEvaluationNewParamsSourceUnion() {}

type EvaluationNewParamsSourceLiveEvaluationSourceKind string

const (
	EvaluationNewParamsSourceLiveEvaluationSourceKindLive EvaluationNewParamsSourceLiveEvaluationSourceKind = "live"
)

func (r EvaluationNewParamsSourceLiveEvaluationSourceKind) IsKnown() bool {
	switch r {
	case EvaluationNewParamsSourceLiveEvaluationSourceKindLive:
		return true
	}
	return false
}

type EvaluationNewParamsSourceKind string

const (
	EvaluationNewParamsSourceKindFixed EvaluationNewParamsSourceKind = "fixed"
	EvaluationNewParamsSourceKindLive  EvaluationNewParamsSourceKind = "live"
)

func (r EvaluationNewParamsSourceKind) IsKnown() bool {
	switch r {
	case EvaluationNewParamsSourceKindFixed, EvaluationNewParamsSourceKindLive:
		return true
	}
	return false
}
