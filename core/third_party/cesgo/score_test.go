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

func TestScoreNew(t *testing.T) {
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
	_, err := client.Scores.New(context.TODO(), coreweaveevaluations.ScoreNewParams{
		Namespace: coreweaveevaluations.F("wandb"),
		ScopeID:   coreweaveevaluations.F("scope_id"),
		Scores: coreweaveevaluations.F([]coreweaveevaluations.ScoreNewParamsScore{{
			ScorerID:        coreweaveevaluations.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			ScorerVersionID: coreweaveevaluations.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			Target: coreweaveevaluations.F(coreweaveevaluations.ScoreNewParamsScoresTarget{
				ID:   coreweaveevaluations.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
				Kind: coreweaveevaluations.F(coreweaveevaluations.ScoreNewParamsScoresTargetKindCase),
			}),
			Value: coreweaveevaluations.F[any](map[string]interface{}{}),
			Attributes: coreweaveevaluations.F(map[string]interface{}{
				"foo": "bar",
			}),
		}}),
	})
	if err != nil {
		var apierr *coreweaveevaluations.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestScoreGet(t *testing.T) {
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
	_, err := client.Scores.Get(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		coreweaveevaluations.ScoreGetParams{
			Namespace: coreweaveevaluations.F("wandb"),
			ScopeID:   coreweaveevaluations.F("scope_id"),
		},
	)
	if err != nil {
		var apierr *coreweaveevaluations.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
