// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package coreweaveevaluations

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/wandb/core/services/evaluations/generated/go/internal/apijson"
	"github.com/wandb/core/services/evaluations/generated/go/internal/param"
	"github.com/wandb/core/services/evaluations/generated/go/internal/requestconfig"
	"github.com/wandb/core/services/evaluations/generated/go/option"
)

// DatasetService contains methods and other services that help with interacting
// with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetService] method instead.
type DatasetService struct {
	Options  []option.RequestOption
	Cases    *DatasetCaseService
	Versions *DatasetVersionService
}

// NewDatasetService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewDatasetService(opts ...option.RequestOption) (r *DatasetService) {
	r = &DatasetService{}
	r.Options = opts
	r.Cases = NewDatasetCaseService(opts...)
	r.Versions = NewDatasetVersionService(opts...)
	return
}

// Request bodies must be smaller than 1 MiB; larger bodies are rejected with 413
// Request Entity Too Large.
func (r *DatasetService) New(
	ctx context.Context,
	params DatasetNewParams,
	opts ...option.RequestOption,
) (res *DatasetNewResponse, err error) {
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
	path := fmt.Sprintf("v1/namespaces/%s/scopes/%s/datasets", params.Namespace, params.ScopeID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// A value in a top-level Case input or output field.
type CaseFieldValueParam struct {
	ExtensionType param.Field[CaseFieldValueExtensionType] `json:"extension_type" api:"required"`
	SchemaVersion param.Field[CaseFieldValueSchemaVersion] `json:"schema_version" api:"required"`
	Sha256        param.Field[string]                      `json:"sha256"         api:"required"`
	Size          param.Field[int64]                       `json:"size"           api:"required"`
	Uri           param.Field[string]                      `json:"uri"            api:"required"`
	WbMediaType   param.Field[CaseFieldValueWbMediaType]   `json:"wb_media_type"  api:"required"`
	Boxes         param.Field[interface{}]                 `json:"boxes"`
	Caption       param.Field[string]                      `json:"caption"`
	Format        param.Field[string]                      `json:"format"`
	Height        param.Field[int64]                       `json:"height"`
	Masks         param.Field[interface{}]                 `json:"masks"`
	Width         param.Field[int64]                       `json:"width"`
}

func (r CaseFieldValueParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r CaseFieldValueParam) ImplementsCaseFieldValueUnionParam() {}

// A value in a top-level Case input or output field.
//
// Satisfied by [WandbImageV1Param], [WandbAudioV1Param], [WandbVideoV1Param],
// [shared.UnionBool], [shared.UnionFloat], [shared.UnionString],
// [CaseFieldValueArrayParam], [CaseFieldValueMapParam], [CaseFieldValueParam].
//
// Use [Raw()] to specify an arbitrary value for this param
type CaseFieldValueUnionParam interface {
	ImplementsCaseFieldValueUnionParam()
}

type CaseFieldValueArrayParam []interface{}

func (r CaseFieldValueArrayParam) ImplementsCaseFieldValueUnionParam() {}

type CaseFieldValueMapParam map[string]interface{}

func (r CaseFieldValueMapParam) ImplementsCaseFieldValueUnionParam() {}

type CaseFieldValueExtensionType string

const (
	CaseFieldValueExtensionTypeWandbImage CaseFieldValueExtensionType = "wandb-image"
	CaseFieldValueExtensionTypeWandbAudio CaseFieldValueExtensionType = "wandb-audio"
	CaseFieldValueExtensionTypeWandbVideo CaseFieldValueExtensionType = "wandb-video"
)

func (r CaseFieldValueExtensionType) IsKnown() bool {
	switch r {
	case CaseFieldValueExtensionTypeWandbImage,
		CaseFieldValueExtensionTypeWandbAudio,
		CaseFieldValueExtensionTypeWandbVideo:
		return true
	}
	return false
}

type CaseFieldValueSchemaVersion int64

const (
	CaseFieldValueSchemaVersion1 CaseFieldValueSchemaVersion = 1
)

func (r CaseFieldValueSchemaVersion) IsKnown() bool {
	switch r {
	case CaseFieldValueSchemaVersion1:
		return true
	}
	return false
}

type CaseFieldValueWbMediaType string

const (
	CaseFieldValueWbMediaTypeImageFile CaseFieldValueWbMediaType = "image-file"
	CaseFieldValueWbMediaTypeAudioFile CaseFieldValueWbMediaType = "audio-file"
	CaseFieldValueWbMediaTypeVideoFile CaseFieldValueWbMediaType = "video-file"
)

func (r CaseFieldValueWbMediaType) IsKnown() bool {
	switch r {
	case CaseFieldValueWbMediaTypeImageFile, CaseFieldValueWbMediaTypeAudioFile, CaseFieldValueWbMediaTypeVideoFile:
		return true
	}
	return false
}

// A Case input or output value.
//
// Satisfied by [CaseValueMapParam], [shared.UnionBool], [shared.UnionFloat],
// [shared.UnionString], [CaseValueArrayParam].
type CaseValueUnionParam interface {
	ImplementsCaseValueUnionParam()
}

type CaseValueMapParam map[string]CaseFieldValueUnionParam

func (r CaseValueMapParam) ImplementsCaseValueUnionParam() {}

type CaseValueArrayParam []interface{}

func (r CaseValueArrayParam) ImplementsCaseValueUnionParam() {}

// A V1 W&B audio value in an evaluation input or output field.
type WandbAudioV1Param struct {
	ExtensionType param.Field[WandbAudioV1ExtensionType] `json:"extension_type" api:"required"`
	SchemaVersion param.Field[WandbAudioV1SchemaVersion] `json:"schema_version" api:"required"`
	Sha256        param.Field[string]                    `json:"sha256"         api:"required"`
	Size          param.Field[int64]                     `json:"size"           api:"required"`
	Uri           param.Field[string]                    `json:"uri"            api:"required"`
	WbMediaType   param.Field[WandbAudioV1WbMediaType]   `json:"wb_media_type"  api:"required"`
	Caption       param.Field[string]                    `json:"caption"`
	ExtraFields   map[string]interface{}                 `json:"-,extras"`
}

func (r WandbAudioV1Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r WandbAudioV1Param) ImplementsCaseFieldValueUnionParam() {}

type WandbAudioV1ExtensionType string

const (
	WandbAudioV1ExtensionTypeWandbAudio WandbAudioV1ExtensionType = "wandb-audio"
)

func (r WandbAudioV1ExtensionType) IsKnown() bool {
	switch r {
	case WandbAudioV1ExtensionTypeWandbAudio:
		return true
	}
	return false
}

type WandbAudioV1SchemaVersion int64

const (
	WandbAudioV1SchemaVersion1 WandbAudioV1SchemaVersion = 1
)

func (r WandbAudioV1SchemaVersion) IsKnown() bool {
	switch r {
	case WandbAudioV1SchemaVersion1:
		return true
	}
	return false
}

type WandbAudioV1WbMediaType string

const (
	WandbAudioV1WbMediaTypeAudioFile WandbAudioV1WbMediaType = "audio-file"
)

func (r WandbAudioV1WbMediaType) IsKnown() bool {
	switch r {
	case WandbAudioV1WbMediaTypeAudioFile:
		return true
	}
	return false
}

// A V1 W&B image value in an evaluation input or output field.
type WandbImageV1Param struct {
	ExtensionType param.Field[WandbImageV1ExtensionType] `json:"extension_type" api:"required"`
	Format        param.Field[string]                    `json:"format"         api:"required"`
	SchemaVersion param.Field[WandbImageV1SchemaVersion] `json:"schema_version" api:"required"`
	Sha256        param.Field[string]                    `json:"sha256"         api:"required"`
	Size          param.Field[int64]                     `json:"size"           api:"required"`
	Uri           param.Field[string]                    `json:"uri"            api:"required"`
	WbMediaType   param.Field[WandbImageV1WbMediaType]   `json:"wb_media_type"  api:"required"`
	Boxes         param.Field[map[string]interface{}]    `json:"boxes"`
	Caption       param.Field[string]                    `json:"caption"`
	Height        param.Field[int64]                     `json:"height"`
	Masks         param.Field[map[string]interface{}]    `json:"masks"`
	Width         param.Field[int64]                     `json:"width"`
	ExtraFields   map[string]interface{}                 `json:"-,extras"`
}

func (r WandbImageV1Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r WandbImageV1Param) ImplementsCaseFieldValueUnionParam() {}

type WandbImageV1ExtensionType string

const (
	WandbImageV1ExtensionTypeWandbImage WandbImageV1ExtensionType = "wandb-image"
)

func (r WandbImageV1ExtensionType) IsKnown() bool {
	switch r {
	case WandbImageV1ExtensionTypeWandbImage:
		return true
	}
	return false
}

type WandbImageV1SchemaVersion int64

const (
	WandbImageV1SchemaVersion1 WandbImageV1SchemaVersion = 1
)

func (r WandbImageV1SchemaVersion) IsKnown() bool {
	switch r {
	case WandbImageV1SchemaVersion1:
		return true
	}
	return false
}

type WandbImageV1WbMediaType string

const (
	WandbImageV1WbMediaTypeImageFile WandbImageV1WbMediaType = "image-file"
)

func (r WandbImageV1WbMediaType) IsKnown() bool {
	switch r {
	case WandbImageV1WbMediaTypeImageFile:
		return true
	}
	return false
}

// A V1 W&B video value in an evaluation input or output field.
type WandbVideoV1Param struct {
	ExtensionType param.Field[WandbVideoV1ExtensionType] `json:"extension_type" api:"required"`
	SchemaVersion param.Field[WandbVideoV1SchemaVersion] `json:"schema_version" api:"required"`
	Sha256        param.Field[string]                    `json:"sha256"         api:"required"`
	Size          param.Field[int64]                     `json:"size"           api:"required"`
	Uri           param.Field[string]                    `json:"uri"            api:"required"`
	WbMediaType   param.Field[WandbVideoV1WbMediaType]   `json:"wb_media_type"  api:"required"`
	Caption       param.Field[string]                    `json:"caption"`
	Height        param.Field[int64]                     `json:"height"`
	Width         param.Field[int64]                     `json:"width"`
	ExtraFields   map[string]interface{}                 `json:"-,extras"`
}

func (r WandbVideoV1Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r WandbVideoV1Param) ImplementsCaseFieldValueUnionParam() {}

type WandbVideoV1ExtensionType string

const (
	WandbVideoV1ExtensionTypeWandbVideo WandbVideoV1ExtensionType = "wandb-video"
)

func (r WandbVideoV1ExtensionType) IsKnown() bool {
	switch r {
	case WandbVideoV1ExtensionTypeWandbVideo:
		return true
	}
	return false
}

type WandbVideoV1SchemaVersion int64

const (
	WandbVideoV1SchemaVersion1 WandbVideoV1SchemaVersion = 1
)

func (r WandbVideoV1SchemaVersion) IsKnown() bool {
	switch r {
	case WandbVideoV1SchemaVersion1:
		return true
	}
	return false
}

type WandbVideoV1WbMediaType string

const (
	WandbVideoV1WbMediaTypeVideoFile WandbVideoV1WbMediaType = "video-file"
)

func (r WandbVideoV1WbMediaType) IsKnown() bool {
	switch r {
	case WandbVideoV1WbMediaTypeVideoFile:
		return true
	}
	return false
}

type DatasetNewResponse struct {
	Attributes    map[string]interface{} `json:"attributes"      api:"required"`
	CreatedAt     time.Time              `json:"created_at"      api:"required" format:"date-time"`
	DatasetID     string                 `json:"dataset_id"      api:"required" format:"uuid"`
	HeadVersionID string                 `json:"head_version_id" api:"required" format:"uuid"`
	Name          string                 `json:"name"            api:"required"`
	// Namespace that assigned the scope reference.
	Namespace DatasetNewResponseNamespace `json:"namespace" api:"required"`
	ScopeRef  string                      `json:"scope_ref" api:"required"`
	JSON      datasetNewResponseJSON      `json:"-"`
}

// datasetNewResponseJSON contains the JSON metadata for the struct
// [DatasetNewResponse]
type datasetNewResponseJSON struct {
	Attributes    apijson.Field
	CreatedAt     apijson.Field
	DatasetID     apijson.Field
	HeadVersionID apijson.Field
	Name          apijson.Field
	Namespace     apijson.Field
	ScopeRef      apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *DatasetNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetNewResponseJSON) RawJSON() string {
	return r.raw
}

// Namespace that assigned the scope reference.
type DatasetNewResponseNamespace string

const (
	DatasetNewResponseNamespaceWandb DatasetNewResponseNamespace = "wandb"
)

func (r DatasetNewResponseNamespace) IsKnown() bool {
	switch r {
	case DatasetNewResponseNamespaceWandb:
		return true
	}
	return false
}

type DatasetNewParams struct {
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
	ScopeID    param.Field[string]                 `path:"scope_id" api:"required"`
	Name       param.Field[string]                 `                api:"required" json:"name"`
	Attributes param.Field[map[string]interface{}] `                               json:"attributes"`
}

func (r DatasetNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
