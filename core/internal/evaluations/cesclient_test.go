package evaluations_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ces "github.com/wandb/cwevals-go"

	"github.com/wandb/wandb/core/internal/api"
	"github.com/wandb/wandb/core/internal/apitest"
	"github.com/wandb/wandb/core/internal/evaluations"
	"github.com/wandb/wandb/core/internal/observabilitytest"
	"github.com/wandb/wandb/core/internal/settings"
)

func TestEvalTableService_UsesGivenURLAndRunCredentials(t *testing.T) {
	server := apitest.NewRecordingServer()
	defer server.Close()
	t.Setenv("CES_BASE_URL", "https://ces.example.com")
	t.Setenv("WANDB_API_KEY", "env-api-key")
	t.Setenv("WANDB_ACCESS_TOKEN", "env-access-token")
	t.Setenv("CES_CUSTOM_HEADERS", "X-From-Env: yes")
	service := evaluations.NewEvalTableService(
		server.URL+"/evaluations",
		evaluations.NewCESHTTPClient(
			api.NewAPIKeyCredentialProvider("run-api-key"),
			observabilitytest.NewTestLogger(t).Logger,
			settings.New(),
		),
	)

	// The recording server's empty response fails to decode; only the
	// request matters here.
	_, _ = service.New(context.Background(), ces.EvalTableNewParams{
		ScopeID: ces.F("scope"),
		Name:    ces.F("table"),
	})

	requests := server.Requests()
	require.Len(t, requests, 1)
	req := requests[0]
	assert.Equal(t,
		"/evaluations/v1/namespaces/wandb/scopes/scope/eval-tables",
		req.URL.Path)
	assert.Equal(t,
		"Basic YXBpOnJ1bi1hcGkta2V5", // "api:run-api-key"
		req.Header.Get("Authorization"))
	assert.Equal(t, "yes", req.Header.Get("X-From-Env"))
}
