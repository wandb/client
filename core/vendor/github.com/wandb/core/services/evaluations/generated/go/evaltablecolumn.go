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

// EvalTableColumnService contains methods and other services that help with
// interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvalTableColumnService] method instead.
type EvalTableColumnService struct {
	Options []option.RequestOption
}

// NewEvalTableColumnService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewEvalTableColumnService(opts ...option.RequestOption) (r *EvalTableColumnService) {
	r = &EvalTableColumnService{}
	r.Options = opts
	return
}

// Request bodies must be smaller than 16 MiB; larger bodies are rejected with 413
// Request Entity Too Large.
func (r *EvalTableColumnService) New(
	ctx context.Context,
	evaluationID string,
	params EvalTableColumnNewParams,
	opts ...option.RequestOption,
) (res *EvalTableColumnNewResponse, err error) {
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
		"v1/namespaces/%s/scopes/%s/eval-tables/%s/columns",
		params.Namespace,
		params.ScopeID,
		evaluationID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type EvalTableColumnNewResponse struct {
	DatasetFields []EvalTableColumnNewResponseDatasetField `json:"dataset_fields" api:"required"`
	Scorers       []EvalTableColumnNewResponseScorer       `json:"scorers"        api:"required"`
	JSON          evalTableColumnNewResponseJSON           `json:"-"`
}

// evalTableColumnNewResponseJSON contains the JSON metadata for the struct
// [EvalTableColumnNewResponse]
type evalTableColumnNewResponseJSON struct {
	DatasetFields apijson.Field
	Scorers       apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *EvalTableColumnNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evalTableColumnNewResponseJSON) RawJSON() string {
	return r.raw
}

type EvalTableColumnNewResponseDatasetField struct {
	Name   string                                        `json:"name"   api:"required"`
	Source EvalTableColumnNewResponseDatasetFieldsSource `json:"source" api:"required"`
	// json may represent an extension-backed field or an ordinary JSON field.
	ValueType EvalTableColumnNewResponseDatasetFieldsValueType `json:"value_type" api:"required"`
	// Extension schema version; omitted for ordinary fields.
	ExtensionSchemaVersion int64 `json:"extension_schema_version"`
	// Extension schema type; omitted for ordinary fields.
	ExtensionType EvalTableColumnNewResponseDatasetFieldsExtensionType `json:"extension_type"`
	JSON          evalTableColumnNewResponseDatasetFieldJSON           `json:"-"`
}

// evalTableColumnNewResponseDatasetFieldJSON contains the JSON metadata for the
// struct [EvalTableColumnNewResponseDatasetField]
type evalTableColumnNewResponseDatasetFieldJSON struct {
	Name                   apijson.Field
	Source                 apijson.Field
	ValueType              apijson.Field
	ExtensionSchemaVersion apijson.Field
	ExtensionType          apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *EvalTableColumnNewResponseDatasetField) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evalTableColumnNewResponseDatasetFieldJSON) RawJSON() string {
	return r.raw
}

type EvalTableColumnNewResponseDatasetFieldsSource string

const (
	EvalTableColumnNewResponseDatasetFieldsSourceInput  EvalTableColumnNewResponseDatasetFieldsSource = "input"
	EvalTableColumnNewResponseDatasetFieldsSourceOutput EvalTableColumnNewResponseDatasetFieldsSource = "output"
)

func (r EvalTableColumnNewResponseDatasetFieldsSource) IsKnown() bool {
	switch r {
	case EvalTableColumnNewResponseDatasetFieldsSourceInput, EvalTableColumnNewResponseDatasetFieldsSourceOutput:
		return true
	}
	return false
}

// json may represent an extension-backed field or an ordinary JSON field.
type EvalTableColumnNewResponseDatasetFieldsValueType string

const (
	EvalTableColumnNewResponseDatasetFieldsValueTypeNull    EvalTableColumnNewResponseDatasetFieldsValueType = "null"
	EvalTableColumnNewResponseDatasetFieldsValueTypeBoolean EvalTableColumnNewResponseDatasetFieldsValueType = "boolean"
	EvalTableColumnNewResponseDatasetFieldsValueTypeInteger EvalTableColumnNewResponseDatasetFieldsValueType = "integer"
	EvalTableColumnNewResponseDatasetFieldsValueTypeNumber  EvalTableColumnNewResponseDatasetFieldsValueType = "number"
	EvalTableColumnNewResponseDatasetFieldsValueTypeString  EvalTableColumnNewResponseDatasetFieldsValueType = "string"
	EvalTableColumnNewResponseDatasetFieldsValueTypeJson    EvalTableColumnNewResponseDatasetFieldsValueType = "json"
)

func (r EvalTableColumnNewResponseDatasetFieldsValueType) IsKnown() bool {
	switch r {
	case EvalTableColumnNewResponseDatasetFieldsValueTypeNull,
		EvalTableColumnNewResponseDatasetFieldsValueTypeBoolean,
		EvalTableColumnNewResponseDatasetFieldsValueTypeInteger,
		EvalTableColumnNewResponseDatasetFieldsValueTypeNumber,
		EvalTableColumnNewResponseDatasetFieldsValueTypeString,
		EvalTableColumnNewResponseDatasetFieldsValueTypeJson:
		return true
	}
	return false
}

// Extension schema type; omitted for ordinary fields.
type EvalTableColumnNewResponseDatasetFieldsExtensionType string

const (
	EvalTableColumnNewResponseDatasetFieldsExtensionTypeWandbImage EvalTableColumnNewResponseDatasetFieldsExtensionType = "wandb-image"
	EvalTableColumnNewResponseDatasetFieldsExtensionTypeWandbAudio EvalTableColumnNewResponseDatasetFieldsExtensionType = "wandb-audio"
	EvalTableColumnNewResponseDatasetFieldsExtensionTypeWandbVideo EvalTableColumnNewResponseDatasetFieldsExtensionType = "wandb-video"
)

func (r EvalTableColumnNewResponseDatasetFieldsExtensionType) IsKnown() bool {
	switch r {
	case EvalTableColumnNewResponseDatasetFieldsExtensionTypeWandbImage,
		EvalTableColumnNewResponseDatasetFieldsExtensionTypeWandbAudio,
		EvalTableColumnNewResponseDatasetFieldsExtensionTypeWandbVideo:
		return true
	}
	return false
}

type EvalTableColumnNewResponseScorer struct {
	Name            string                                     `json:"name"              api:"required"`
	ScorerID        string                                     `json:"scorer_id"         api:"required" format:"uuid"`
	ScorerVersionID string                                     `json:"scorer_version_id" api:"required" format:"uuid"`
	ValueType       EvalTableColumnNewResponseScorersValueType `json:"value_type"        api:"required"`
	JSON            evalTableColumnNewResponseScorerJSON       `json:"-"`
}

// evalTableColumnNewResponseScorerJSON contains the JSON metadata for the struct
// [EvalTableColumnNewResponseScorer]
type evalTableColumnNewResponseScorerJSON struct {
	Name            apijson.Field
	ScorerID        apijson.Field
	ScorerVersionID apijson.Field
	ValueType       apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *EvalTableColumnNewResponseScorer) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evalTableColumnNewResponseScorerJSON) RawJSON() string {
	return r.raw
}

type EvalTableColumnNewResponseScorersValueType string

const (
	EvalTableColumnNewResponseScorersValueTypeBoolean EvalTableColumnNewResponseScorersValueType = "boolean"
	EvalTableColumnNewResponseScorersValueTypeInteger EvalTableColumnNewResponseScorersValueType = "integer"
	EvalTableColumnNewResponseScorersValueTypeNumber  EvalTableColumnNewResponseScorersValueType = "number"
	EvalTableColumnNewResponseScorersValueTypeString  EvalTableColumnNewResponseScorersValueType = "string"
)

func (r EvalTableColumnNewResponseScorersValueType) IsKnown() bool {
	switch r {
	case EvalTableColumnNewResponseScorersValueTypeBoolean,
		EvalTableColumnNewResponseScorersValueTypeInteger,
		EvalTableColumnNewResponseScorersValueTypeNumber,
		EvalTableColumnNewResponseScorersValueTypeString:
		return true
	}
	return false
}

type EvalTableColumnNewParams struct {
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
	ScopeID       param.Field[string]                                 `path:"scope_id" api:"required"`
	DatasetFields param.Field[[]EvalTableColumnNewParamsDatasetField] `                               json:"dataset_fields"`
	Scorers       param.Field[[]EvalTableColumnNewParamsScorer]       `                               json:"scorers"`
}

func (r EvalTableColumnNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type EvalTableColumnNewParamsDatasetField struct {
	Name   param.Field[string]                                      `json:"name"   api:"required"`
	Source param.Field[EvalTableColumnNewParamsDatasetFieldsSource] `json:"source" api:"required"`
	// Use json for extension-backed fields and ordinary JSON fields.
	ValueType param.Field[EvalTableColumnNewParamsDatasetFieldsValueType] `json:"value_type" api:"required"`
	// Extension schema version. Set with extension_type; omit both for ordinary
	// fields.
	ExtensionSchemaVersion param.Field[int64] `json:"extension_schema_version"`
	// Extension schema type. Set with extension_schema_version; omit both for ordinary
	// fields.
	ExtensionType param.Field[EvalTableColumnNewParamsDatasetFieldsExtensionType] `json:"extension_type"`
}

func (r EvalTableColumnNewParamsDatasetField) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type EvalTableColumnNewParamsDatasetFieldsSource string

const (
	EvalTableColumnNewParamsDatasetFieldsSourceInput  EvalTableColumnNewParamsDatasetFieldsSource = "input"
	EvalTableColumnNewParamsDatasetFieldsSourceOutput EvalTableColumnNewParamsDatasetFieldsSource = "output"
)

func (r EvalTableColumnNewParamsDatasetFieldsSource) IsKnown() bool {
	switch r {
	case EvalTableColumnNewParamsDatasetFieldsSourceInput, EvalTableColumnNewParamsDatasetFieldsSourceOutput:
		return true
	}
	return false
}

// Use json for extension-backed fields and ordinary JSON fields.
type EvalTableColumnNewParamsDatasetFieldsValueType string

const (
	EvalTableColumnNewParamsDatasetFieldsValueTypeBoolean EvalTableColumnNewParamsDatasetFieldsValueType = "boolean"
	EvalTableColumnNewParamsDatasetFieldsValueTypeInteger EvalTableColumnNewParamsDatasetFieldsValueType = "integer"
	EvalTableColumnNewParamsDatasetFieldsValueTypeNumber  EvalTableColumnNewParamsDatasetFieldsValueType = "number"
	EvalTableColumnNewParamsDatasetFieldsValueTypeString  EvalTableColumnNewParamsDatasetFieldsValueType = "string"
	EvalTableColumnNewParamsDatasetFieldsValueTypeJson    EvalTableColumnNewParamsDatasetFieldsValueType = "json"
)

func (r EvalTableColumnNewParamsDatasetFieldsValueType) IsKnown() bool {
	switch r {
	case EvalTableColumnNewParamsDatasetFieldsValueTypeBoolean,
		EvalTableColumnNewParamsDatasetFieldsValueTypeInteger,
		EvalTableColumnNewParamsDatasetFieldsValueTypeNumber,
		EvalTableColumnNewParamsDatasetFieldsValueTypeString,
		EvalTableColumnNewParamsDatasetFieldsValueTypeJson:
		return true
	}
	return false
}

// Extension schema type. Set with extension_schema_version; omit both for ordinary
// fields.
type EvalTableColumnNewParamsDatasetFieldsExtensionType string

const (
	EvalTableColumnNewParamsDatasetFieldsExtensionTypeWandbImage EvalTableColumnNewParamsDatasetFieldsExtensionType = "wandb-image"
	EvalTableColumnNewParamsDatasetFieldsExtensionTypeWandbAudio EvalTableColumnNewParamsDatasetFieldsExtensionType = "wandb-audio"
	EvalTableColumnNewParamsDatasetFieldsExtensionTypeWandbVideo EvalTableColumnNewParamsDatasetFieldsExtensionType = "wandb-video"
)

func (r EvalTableColumnNewParamsDatasetFieldsExtensionType) IsKnown() bool {
	switch r {
	case EvalTableColumnNewParamsDatasetFieldsExtensionTypeWandbImage,
		EvalTableColumnNewParamsDatasetFieldsExtensionTypeWandbAudio,
		EvalTableColumnNewParamsDatasetFieldsExtensionTypeWandbVideo:
		return true
	}
	return false
}

type EvalTableColumnNewParamsScorer struct {
	// Must not have leading or trailing whitespace.
	Name      param.Field[string]                                   `json:"name"       api:"required"`
	ValueType param.Field[EvalTableColumnNewParamsScorersValueType] `json:"value_type" api:"required"`
}

func (r EvalTableColumnNewParamsScorer) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type EvalTableColumnNewParamsScorersValueType string

const (
	EvalTableColumnNewParamsScorersValueTypeBoolean EvalTableColumnNewParamsScorersValueType = "boolean"
	EvalTableColumnNewParamsScorersValueTypeInteger EvalTableColumnNewParamsScorersValueType = "integer"
	EvalTableColumnNewParamsScorersValueTypeNumber  EvalTableColumnNewParamsScorersValueType = "number"
	EvalTableColumnNewParamsScorersValueTypeString  EvalTableColumnNewParamsScorersValueType = "string"
)

func (r EvalTableColumnNewParamsScorersValueType) IsKnown() bool {
	switch r {
	case EvalTableColumnNewParamsScorersValueTypeBoolean,
		EvalTableColumnNewParamsScorersValueTypeInteger,
		EvalTableColumnNewParamsScorersValueTypeNumber,
		EvalTableColumnNewParamsScorersValueTypeString:
		return true
	}
	return false
}
