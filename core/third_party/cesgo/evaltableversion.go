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

// EvalTableVersionService contains methods and other services that help with
// interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvalTableVersionService] method instead.
type EvalTableVersionService struct {
	Options []option.RequestOption
	Columns *EvalTableVersionColumnService
	Rows    *EvalTableVersionRowService
}

// NewEvalTableVersionService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewEvalTableVersionService(opts ...option.RequestOption) (r *EvalTableVersionService) {
	r = &EvalTableVersionService{}
	r.Options = opts
	r.Columns = NewEvalTableVersionColumnService(opts...)
	r.Rows = NewEvalTableVersionRowService(opts...)
	return
}

// Cut Dataset and Evaluation versions for an EvalTable
func (r *EvalTableVersionService) New(
	ctx context.Context,
	evaluationID string,
	body EvalTableVersionNewParams,
	opts ...option.RequestOption,
) (res *EvalTableVersionNewResponse, err error) {
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
	if evaluationID == "" {
		err = errors.New("missing required evaluation_id parameter")
		return nil, err
	}
	path := fmt.Sprintf(
		"v1/namespaces/%s/scopes/%s/eval-tables/%s/versions",
		body.Namespace,
		body.ScopeID,
		evaluationID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type EvalTableVersionNewResponse struct {
	DatasetVersionID    string                          `json:"dataset_version_id"    api:"required" format:"uuid"`
	EvaluationVersionID string                          `json:"evaluation_version_id" api:"required" format:"uuid"`
	JSON                evalTableVersionNewResponseJSON `json:"-"`
}

// evalTableVersionNewResponseJSON contains the JSON metadata for the struct
// [EvalTableVersionNewResponse]
type evalTableVersionNewResponseJSON struct {
	DatasetVersionID    apijson.Field
	EvaluationVersionID apijson.Field
	raw                 string
	ExtraFields         map[string]apijson.Field
}

func (r *EvalTableVersionNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evalTableVersionNewResponseJSON) RawJSON() string {
	return r.raw
}

type EvalTableVersionNewParams struct {
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
