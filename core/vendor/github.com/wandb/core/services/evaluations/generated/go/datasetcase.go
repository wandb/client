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

// DatasetCaseService contains methods and other services that help with
// interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetCaseService] method instead.
type DatasetCaseService struct {
	Options []option.RequestOption
}

// NewDatasetCaseService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDatasetCaseService(opts ...option.RequestOption) (r *DatasetCaseService) {
	r = &DatasetCaseService{}
	r.Options = opts
	return
}

// Creates Cases in a Dataset and catalogues top-level fields from input and output
// objects.
//
// Request bodies must be smaller than 16 MiB; larger bodies are rejected with 413
// Request Entity Too Large.
func (r *DatasetCaseService) Add(
	ctx context.Context,
	datasetID string,
	params DatasetCaseAddParams,
	opts ...option.RequestOption,
) (res *DatasetCaseAddResponse, err error) {
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
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/namespaces/%s/scopes/%s/datasets/%s/cases", params.Namespace, params.ScopeID, datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, params, &res, opts...)
	return res, err
}

type DatasetCaseAddResponse struct {
	CaseIDs            []string                   `json:"case_ids"             api:"required" format:"uuid"`
	CommitID           int64                      `json:"commit_id"            api:"required"`
	PublishOperationID string                     `json:"publish_operation_id" api:"required" format:"uuid"`
	JSON               datasetCaseAddResponseJSON `json:"-"`
}

// datasetCaseAddResponseJSON contains the JSON metadata for the struct
// [DatasetCaseAddResponse]
type datasetCaseAddResponseJSON struct {
	CaseIDs            apijson.Field
	CommitID           apijson.Field
	PublishOperationID apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *DatasetCaseAddResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetCaseAddResponseJSON) RawJSON() string {
	return r.raw
}

type DatasetCaseAddParams struct {
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
	ScopeID param.Field[string]                     `path:"scope_id" api:"required"`
	Cases   param.Field[[]DatasetCaseAddParamsCase] `                api:"required" json:"cases"`
}

func (r DatasetCaseAddParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type DatasetCaseAddParamsCase struct {
	// For objects, top-level field names must contain 1 to 512 Unicode characters and
	// must not have leading or trailing whitespace. Numeric fields infer number,
	// including whole-valued literals. Indexed numbers use float64 and may round
	// integers beyond 2^53; canonical Case JSON remains exact.
	Input param.Field[CaseValueUnionParam] `json:"input" api:"required"`
	// For objects, top-level field names must contain 1 to 512 Unicode characters and
	// must not have leading or trailing whitespace. Numeric fields infer number,
	// including whole-valued literals. Indexed numbers use float64 and may round
	// integers beyond 2^53; canonical Case JSON remains exact.
	Output     param.Field[CaseValueUnionParam]    `json:"output"     api:"required"`
	Attributes param.Field[map[string]interface{}] `json:"attributes"`
	// Omit or use null when no subject applies.
	Subject param.Field[string] `json:"subject"`
}

func (r DatasetCaseAddParamsCase) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
