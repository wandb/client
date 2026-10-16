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

// DatasetVersionService contains methods and other services that help with
// interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetVersionService] method instead.
type DatasetVersionService struct {
	Options []option.RequestOption
	Cases   *DatasetVersionCaseService
	Fields  *DatasetVersionFieldService
}

// NewDatasetVersionService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDatasetVersionService(opts ...option.RequestOption) (r *DatasetVersionService) {
	r = &DatasetVersionService{}
	r.Options = opts
	r.Cases = NewDatasetVersionCaseService(opts...)
	r.Fields = NewDatasetVersionFieldService(opts...)
	return
}

// Create an immutable DatasetVersion snapshot
func (r *DatasetVersionService) New(
	ctx context.Context,
	datasetID string,
	body DatasetVersionNewParams,
	opts ...option.RequestOption,
) (res *DatasetVersionNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	precfg, err := requestconfig.PreRequestOptions(opts...)
	if err != nil {
		return nil, err
	}
	requestconfig.UseDefaultParam(&body.Namespace, precfg.Namespace)
	requestconfig.UseDefaultParam(&body.ScopeID, precfg.ScopeID)
	if body.Namespace.Value == "" {
		err = errors.New("missing required namespace parameter")
		return nil, err
	}
	if body.ScopeID.Value == "" {
		err = errors.New("missing required scope_id parameter")
		return nil, err
	}
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/namespaces/%s/scopes/%s/datasets/%s/versions", body.Namespace, body.ScopeID, datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Get a DatasetVersion
func (r *DatasetVersionService) Get(
	ctx context.Context,
	datasetVersionID string,
	query DatasetVersionGetParams,
	opts ...option.RequestOption,
) (res *DatasetVersionGetResponse, err error) {
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
	if datasetVersionID == "" {
		err = errors.New("missing required dataset_version_id parameter")
		return nil, err
	}
	path := fmt.Sprintf(
		"v1/namespaces/%s/scopes/%s/dataset-versions/%s",
		query.Namespace,
		query.ScopeID,
		datasetVersionID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type DatasetVersionNewResponse struct {
	CaseCount        int64     `json:"case_count"         api:"required"`
	CreatedAt        time.Time `json:"created_at"         api:"required" format:"date-time"`
	DatasetID        string    `json:"dataset_id"         api:"required" format:"uuid"`
	DatasetVersionID string    `json:"dataset_version_id" api:"required" format:"uuid"`
	LastCommitID     int64     `json:"last_commit_id"     api:"required"`
	// Namespace that assigned the scope reference.
	Namespace       DatasetVersionNewResponseNamespace `json:"namespace"         api:"required"`
	ScopeRef        string                             `json:"scope_ref"         api:"required"`
	Status          DatasetVersionNewResponseStatus    `json:"status"            api:"required"`
	TrajectoryCount int64                              `json:"trajectory_count"  api:"required"`
	VersionSequence int64                              `json:"version_sequence"  api:"required"`
	ParentVersionID string                             `json:"parent_version_id"                format:"uuid"`
	JSON            datasetVersionNewResponseJSON      `json:"-"`
}

// datasetVersionNewResponseJSON contains the JSON metadata for the struct
// [DatasetVersionNewResponse]
type datasetVersionNewResponseJSON struct {
	CaseCount        apijson.Field
	CreatedAt        apijson.Field
	DatasetID        apijson.Field
	DatasetVersionID apijson.Field
	LastCommitID     apijson.Field
	Namespace        apijson.Field
	ScopeRef         apijson.Field
	Status           apijson.Field
	TrajectoryCount  apijson.Field
	VersionSequence  apijson.Field
	ParentVersionID  apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *DatasetVersionNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetVersionNewResponseJSON) RawJSON() string {
	return r.raw
}

// Namespace that assigned the scope reference.
type DatasetVersionNewResponseNamespace string

const (
	DatasetVersionNewResponseNamespaceWandb DatasetVersionNewResponseNamespace = "wandb"
)

func (r DatasetVersionNewResponseNamespace) IsKnown() bool {
	switch r {
	case DatasetVersionNewResponseNamespaceWandb:
		return true
	}
	return false
}

type DatasetVersionNewResponseStatus string

const (
	DatasetVersionNewResponseStatusPending   DatasetVersionNewResponseStatus = "pending"
	DatasetVersionNewResponseStatusPublished DatasetVersionNewResponseStatus = "published"
)

func (r DatasetVersionNewResponseStatus) IsKnown() bool {
	switch r {
	case DatasetVersionNewResponseStatusPending, DatasetVersionNewResponseStatusPublished:
		return true
	}
	return false
}

type DatasetVersionGetResponse struct {
	CaseCount        int64     `json:"case_count"         api:"required"`
	CreatedAt        time.Time `json:"created_at"         api:"required" format:"date-time"`
	DatasetID        string    `json:"dataset_id"         api:"required" format:"uuid"`
	DatasetVersionID string    `json:"dataset_version_id" api:"required" format:"uuid"`
	LastCommitID     int64     `json:"last_commit_id"     api:"required"`
	// Namespace that assigned the scope reference.
	Namespace       DatasetVersionGetResponseNamespace `json:"namespace"         api:"required"`
	ScopeRef        string                             `json:"scope_ref"         api:"required"`
	Status          DatasetVersionGetResponseStatus    `json:"status"            api:"required"`
	TrajectoryCount int64                              `json:"trajectory_count"  api:"required"`
	VersionSequence int64                              `json:"version_sequence"  api:"required"`
	ParentVersionID string                             `json:"parent_version_id"                format:"uuid"`
	JSON            datasetVersionGetResponseJSON      `json:"-"`
}

// datasetVersionGetResponseJSON contains the JSON metadata for the struct
// [DatasetVersionGetResponse]
type datasetVersionGetResponseJSON struct {
	CaseCount        apijson.Field
	CreatedAt        apijson.Field
	DatasetID        apijson.Field
	DatasetVersionID apijson.Field
	LastCommitID     apijson.Field
	Namespace        apijson.Field
	ScopeRef         apijson.Field
	Status           apijson.Field
	TrajectoryCount  apijson.Field
	VersionSequence  apijson.Field
	ParentVersionID  apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *DatasetVersionGetResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetVersionGetResponseJSON) RawJSON() string {
	return r.raw
}

// Namespace that assigned the scope reference.
type DatasetVersionGetResponseNamespace string

const (
	DatasetVersionGetResponseNamespaceWandb DatasetVersionGetResponseNamespace = "wandb"
)

func (r DatasetVersionGetResponseNamespace) IsKnown() bool {
	switch r {
	case DatasetVersionGetResponseNamespaceWandb:
		return true
	}
	return false
}

type DatasetVersionGetResponseStatus string

const (
	DatasetVersionGetResponseStatusPending   DatasetVersionGetResponseStatus = "pending"
	DatasetVersionGetResponseStatusPublished DatasetVersionGetResponseStatus = "published"
)

func (r DatasetVersionGetResponseStatus) IsKnown() bool {
	switch r {
	case DatasetVersionGetResponseStatusPending, DatasetVersionGetResponseStatusPublished:
		return true
	}
	return false
}

type DatasetVersionNewParams struct {
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

type DatasetVersionGetParams struct {
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
