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

func TestScorerNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Scorers.New(context.TODO(), coreweaveevaluations.ScorerNewParams{
		Namespace: coreweaveevaluations.F("wandb"),
		ScopeID:   coreweaveevaluations.F("scope_id"),
		AllowedTargetKinds: coreweaveevaluations.F(
			[]coreweaveevaluations.ScorerNewParamsAllowedTargetKind{
				coreweaveevaluations.ScorerNewParamsAllowedTargetKindCase,
			},
		),
		Name:        coreweaveevaluations.F("x"),
		ValueType:   coreweaveevaluations.F(coreweaveevaluations.ScorerNewParamsValueTypeNull),
		Description: coreweaveevaluations.F("description"),
	})
	if err != nil {
		var apierr *coreweaveevaluations.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
