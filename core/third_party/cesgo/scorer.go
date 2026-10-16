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

// ScorerService contains methods and other services that help with interacting
// with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewScorerService] method instead.
type ScorerService struct {
	Options  []option.RequestOption
	Versions *ScorerVersionService
}

// NewScorerService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewScorerService(opts ...option.RequestOption) (r *ScorerService) {
	r = &ScorerService{}
	r.Options = opts
	r.Versions = NewScorerVersionService(opts...)
	return
}

// Request bodies must be smaller than 1 MiB; larger bodies are rejected with 413
// Request Entity Too Large.
func (r *ScorerService) New(
	ctx context.Context,
	params ScorerNewParams,
	opts ...option.RequestOption,
) (res *ScorerNewResponse, err error) {
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
	path := fmt.Sprintf("v1/namespaces/%s/scopes/%s/scorers", params.Namespace, params.ScopeID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type ScorerNewResponse struct {
	AllowedTargetKinds []ScorerNewResponseAllowedTargetKind `json:"allowed_target_kinds" api:"required"`
	Description        string                               `json:"description"          api:"required,nullable"`
	Name               string                               `json:"name"                 api:"required"`
	// Namespace that assigned the scope reference.
	Namespace ScorerNewResponseNamespace `json:"namespace"  api:"required"`
	ScopeRef  string                     `json:"scope_ref"  api:"required"`
	ScorerID  string                     `json:"scorer_id"  api:"required" format:"uuid"`
	ValueType ScorerNewResponseValueType `json:"value_type" api:"required"`
	VersionID string                     `json:"version_id" api:"required" format:"uuid"`
	JSON      scorerNewResponseJSON      `json:"-"`
}

// scorerNewResponseJSON contains the JSON metadata for the struct
// [ScorerNewResponse]
type scorerNewResponseJSON struct {
	AllowedTargetKinds apijson.Field
	Description        apijson.Field
	Name               apijson.Field
	Namespace          apijson.Field
	ScopeRef           apijson.Field
	ScorerID           apijson.Field
	ValueType          apijson.Field
	VersionID          apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *ScorerNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r scorerNewResponseJSON) RawJSON() string {
	return r.raw
}

type ScorerNewResponseAllowedTargetKind string

const (
	ScorerNewResponseAllowedTargetKindCase       ScorerNewResponseAllowedTargetKind = "case"
	ScorerNewResponseAllowedTargetKindComparison ScorerNewResponseAllowedTargetKind = "comparison"
	ScorerNewResponseAllowedTargetKindTrajectory ScorerNewResponseAllowedTargetKind = "trajectory"
)

func (r ScorerNewResponseAllowedTargetKind) IsKnown() bool {
	switch r {
	case ScorerNewResponseAllowedTargetKindCase,
		ScorerNewResponseAllowedTargetKindComparison,
		ScorerNewResponseAllowedTargetKindTrajectory:
		return true
	}
	return false
}

// Namespace that assigned the scope reference.
type ScorerNewResponseNamespace string

const (
	ScorerNewResponseNamespaceWandb ScorerNewResponseNamespace = "wandb"
)

func (r ScorerNewResponseNamespace) IsKnown() bool {
	switch r {
	case ScorerNewResponseNamespaceWandb:
		return true
	}
	return false
}

type ScorerNewResponseValueType string

const (
	ScorerNewResponseValueTypeNull    ScorerNewResponseValueType = "null"
	ScorerNewResponseValueTypeBoolean ScorerNewResponseValueType = "boolean"
	ScorerNewResponseValueTypeInteger ScorerNewResponseValueType = "integer"
	ScorerNewResponseValueTypeNumber  ScorerNewResponseValueType = "number"
	ScorerNewResponseValueTypeString  ScorerNewResponseValueType = "string"
	ScorerNewResponseValueTypeJson    ScorerNewResponseValueType = "json"
)

func (r ScorerNewResponseValueType) IsKnown() bool {
	switch r {
	case ScorerNewResponseValueTypeNull,
		ScorerNewResponseValueTypeBoolean,
		ScorerNewResponseValueTypeInteger,
		ScorerNewResponseValueTypeNumber,
		ScorerNewResponseValueTypeString,
		ScorerNewResponseValueTypeJson:
		return true
	}
	return false
}

type ScorerNewParams struct {
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
	ScopeID            param.Field[string]                             `path:"scope_id" api:"required"`
	AllowedTargetKinds param.Field[[]ScorerNewParamsAllowedTargetKind] `                api:"required" json:"allowed_target_kinds"`
	// Must not have leading or trailing whitespace.
	Name        param.Field[string]                   `json:"name"        api:"required"`
	ValueType   param.Field[ScorerNewParamsValueType] `json:"value_type"  api:"required"`
	Description param.Field[string]                   `json:"description"`
}

func (r ScorerNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ScorerNewParamsAllowedTargetKind string

const (
	ScorerNewParamsAllowedTargetKindCase       ScorerNewParamsAllowedTargetKind = "case"
	ScorerNewParamsAllowedTargetKindComparison ScorerNewParamsAllowedTargetKind = "comparison"
	ScorerNewParamsAllowedTargetKindTrajectory ScorerNewParamsAllowedTargetKind = "trajectory"
)

func (r ScorerNewParamsAllowedTargetKind) IsKnown() bool {
	switch r {
	case ScorerNewParamsAllowedTargetKindCase,
		ScorerNewParamsAllowedTargetKindComparison,
		ScorerNewParamsAllowedTargetKindTrajectory:
		return true
	}
	return false
}

type ScorerNewParamsValueType string

const (
	ScorerNewParamsValueTypeNull    ScorerNewParamsValueType = "null"
	ScorerNewParamsValueTypeBoolean ScorerNewParamsValueType = "boolean"
	ScorerNewParamsValueTypeInteger ScorerNewParamsValueType = "integer"
	ScorerNewParamsValueTypeNumber  ScorerNewParamsValueType = "number"
	ScorerNewParamsValueTypeString  ScorerNewParamsValueType = "string"
	ScorerNewParamsValueTypeJson    ScorerNewParamsValueType = "json"
)

func (r ScorerNewParamsValueType) IsKnown() bool {
	switch r {
	case ScorerNewParamsValueTypeNull,
		ScorerNewParamsValueTypeBoolean,
		ScorerNewParamsValueTypeInteger,
		ScorerNewParamsValueTypeNumber,
		ScorerNewParamsValueTypeString,
		ScorerNewParamsValueTypeJson:
		return true
	}
	return false
}
