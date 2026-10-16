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

// EvalTableRowService contains methods and other services that help with
// interacting with the CoreWeave Evaluations API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvalTableRowService] method instead.
type EvalTableRowService struct {
	Options []option.RequestOption
}

// NewEvalTableRowService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewEvalTableRowService(opts ...option.RequestOption) (r *EvalTableRowService) {
	r = &EvalTableRowService{}
	r.Options = opts
	return
}

// Request bodies must be smaller than 16 MiB; larger bodies are rejected with 413
// Request Entity Too Large.
func (r *EvalTableRowService) Add(
	ctx context.Context,
	evaluationID string,
	params EvalTableRowAddParams,
	opts ...option.RequestOption,
) (res *EvalTableRowAddResponse, err error) {
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
		"v1/namespaces/%s/scopes/%s/eval-tables/%s/rows",
		params.Namespace,
		params.ScopeID,
		evaluationID,
	)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type EvalTableRowAddResponse struct {
	CaseIDs  []string                    `json:"case_ids"  api:"required" format:"uuid"`
	ScoreIDs []string                    `json:"score_ids" api:"required" format:"uuid"`
	JSON     evalTableRowAddResponseJSON `json:"-"`
}

// evalTableRowAddResponseJSON contains the JSON metadata for the struct
// [EvalTableRowAddResponse]
type evalTableRowAddResponseJSON struct {
	CaseIDs     apijson.Field
	ScoreIDs    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvalTableRowAddResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evalTableRowAddResponseJSON) RawJSON() string {
	return r.raw
}

type EvalTableRowAddParams struct {
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
	Rows    param.Field[[]EvalTableRowAddParamsRow] `                api:"required" json:"rows"`
}

func (r EvalTableRowAddParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type EvalTableRowAddParamsRow struct {
	// For objects, top-level field names must contain 1 to 512 Unicode characters and
	// must not have leading or trailing whitespace. Numeric fields infer number,
	// including whole-valued literals. Indexed numbers use float64 and may round
	// integers beyond 2^53; canonical Case JSON remains exact.
	Input param.Field[CaseValueUnionParam] `json:"input" api:"required"`
	// For objects, top-level field names must contain 1 to 512 Unicode characters and
	// must not have leading or trailing whitespace. Numeric fields infer number,
	// including whole-valued literals. Indexed numbers use float64 and may round
	// integers beyond 2^53; canonical Case JSON remains exact.
	Output param.Field[CaseValueUnionParam] `json:"output" api:"required"`
	// Keys must not have leading or trailing whitespace.
	Scores     param.Field[map[string]interface{}] `json:"scores"     api:"required"`
	Attributes param.Field[map[string]interface{}] `json:"attributes"`
}

func (r EvalTableRowAddParamsRow) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
