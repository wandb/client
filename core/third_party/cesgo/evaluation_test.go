// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package coreweaveevaluations_test

import (
	"context"
	"errors"
	"os"
	"testing"

	coreweaveevaluations "github.com/wandb/core/services/evaluations/generated/go"
	"github.com/wandb/core/services/evaluations/generated/go/internal/testutil"
	"github.com/wandb/core/services/evaluations/generated/go/option"
)

func TestEvaluationNewWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := coreweaveevaluations.NewClient(
		option.WithBaseURL(baseURL),
		option.WithBearerToken("My Bearer Token"),
		option.WithUsername("My Username"),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Evaluations.New(context.TODO(), coreweaveevaluations.EvaluationNewParams{
		Namespace:         coreweaveevaluations.F("wandb"),
		ScopeID:           coreweaveevaluations.F("scope_id"),
		Name:              coreweaveevaluations.F("x"),
		PrimaryTargetKind: coreweaveevaluations.F(coreweaveevaluations.EvaluationNewParamsPrimaryTargetKindCase),
		Source: coreweaveevaluations.F[coreweaveevaluations.EvaluationNewParamsSourceUnion](
			coreweaveevaluations.EvaluationNewParamsSourceFixedEvaluationSource{
				DatasetVersionID: coreweaveevaluations.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
				Kind: coreweaveevaluations.F(
					coreweaveevaluations.EvaluationNewParamsSourceFixedEvaluationSourceKindFixed,
				),
			},
		),
		Attributes: coreweaveevaluations.F(map[string]interface{}{
			"foo": "bar",
		}),
	})
	if err != nil {
		var apierr *coreweaveevaluations.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
