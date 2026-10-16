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

// DatasetVersionCaseService contains methods and other services that help with
// interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetVersionCaseService] method instead.
type DatasetVersionCaseService struct {
	Options []option.RequestOption
}

// NewDatasetVersionCaseService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewDatasetVersionCaseService(opts ...option.RequestOption) (r *DatasetVersionCaseService) {
	r = &DatasetVersionCaseService{}
	r.Options = opts
	return
}

// List cases in a DatasetVersion
func (r *DatasetVersionCaseService) List(
	ctx context.Context,
	datasetVersionID string,
	params DatasetVersionCaseListParams,
	opts ...option.RequestOption,
) (res *DatasetVersionCaseListResponse, err error) {
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
	if datasetVersionID == "" {
		err = errors.New("missing required dataset_version_id parameter")
		return nil, err
	}
	path := fmt.Sprintf(
		"v1/namespaces/%s/scopes/%s/dataset-versions/%s/cases",
		params.Namespace,
		params.ScopeID,
		datasetVersionID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

type DatasetVersionCaseListResponse struct {
	Cases      []DatasetVersionCaseListResponseCase `json:"cases"       api:"required"`
	NextCursor string                               `json:"next_cursor"`
	JSON       datasetVersionCaseListResponseJSON   `json:"-"`
}

// datasetVersionCaseListResponseJSON contains the JSON metadata for the struct
// [DatasetVersionCaseListResponse]
type datasetVersionCaseListResponseJSON struct {
	Cases       apijson.Field
	NextCursor  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *DatasetVersionCaseListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetVersionCaseListResponseJSON) RawJSON() string {
	return r.raw
}

type DatasetVersionCaseListResponseCase struct {
	Attributes map[string]interface{}                 `json:"attributes" api:"required"`
	CaseID     string                                 `json:"case_id"    api:"required"          format:"uuid"`
	CommitID   int64                                  `json:"commit_id"  api:"required"`
	Input      interface{}                            `json:"input"      api:"required"`
	Output     interface{}                            `json:"output"     api:"required"`
	Subject    string                                 `json:"subject"    api:"required,nullable"`
	JSON       datasetVersionCaseListResponseCaseJSON `json:"-"`
}

// datasetVersionCaseListResponseCaseJSON contains the JSON metadata for the struct
// [DatasetVersionCaseListResponseCase]
type datasetVersionCaseListResponseCaseJSON struct {
	Attributes  apijson.Field
	CaseID      apijson.Field
	CommitID    apijson.Field
	Input       apijson.Field
	Output      apijson.Field
	Subject     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *DatasetVersionCaseListResponseCase) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetVersionCaseListResponseCaseJSON) RawJSON() string {
	return r.raw
}

type DatasetVersionCaseListParams struct {
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

// URLQuery serializes [DatasetVersionCaseListParams]'s query parameters as
// `url.Values`.
func (r DatasetVersionCaseListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
