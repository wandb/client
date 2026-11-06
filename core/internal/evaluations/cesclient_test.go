package evaluations_test

import (
	"context"
	"net/url"
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

func mustParseURL(t *testing.T, rawURL string) *url.URL {
	t.Helper()
	u, err := url.Parse(rawURL)
	require.NoError(t, err)
	return u
}

func TestCESBaseURL_Default(t *testing.T) {
	t.Setenv("CES_BASE_URL", "")

	cesURL, err := evaluations.CESBaseURL()

	require.NoError(t, err)
	assert.Equal(t, "https://evaluations.wandb.ai/", cesURL.String())
}

func TestCESBaseURL_EnvironmentOverride(t *testing.T) {
	t.Setenv("CES_BASE_URL", "https://ces.example.com")

	cesURL, err := evaluations.CESBaseURL()

	require.NoError(t, err)
	assert.Equal(t, "https://ces.example.com", cesURL.String())
}

func TestEvalTableService_UsesRunCredentialsNotSDKEnvironment(t *testing.T) {
	t.Setenv("WANDB_API_KEY", "env-api-key")
	t.Setenv("CES_CUSTOM_HEADERS", "X-From-Env: yes")
	server := apitest.NewRecordingServer()
	defer server.Close()
	cesURL := mustParseURL(t, server.URL+"/evaluations")
	service := evaluations.NewEvalTableService(
		cesURL,
		evaluations.NewCESHTTPClient(
			cesURL,
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
	assert.Empty(t, req.Header.Get("X-From-Env"))
}
