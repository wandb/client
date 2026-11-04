package evaltables_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ces "github.com/wandb/cwevals-go"

	"github.com/wandb/wandb/core/internal/api"
	"github.com/wandb/wandb/core/internal/apitest"
	"github.com/wandb/wandb/core/internal/evaltables"
	"github.com/wandb/wandb/core/internal/observabilitytest"
	"github.com/wandb/wandb/core/internal/settings"
)

func mustParseURL(t *testing.T, rawURL string) *url.URL {
	t.Helper()
	u, err := url.Parse(rawURL)
	require.NoError(t, err)
	return u
}

func TestCESBaseURL(t *testing.T) {
	testCases := []struct {
		wandbBaseURL string
		want         string
	}{
		{"https://api.wandb.ai", "https://evaluations.wandb.ai"},
		{"https://api.forge.coreweave.com", "https://evaluations.wandb.ai"},
		{"https://forge.coreweave.com/api/wandb", "https://evaluations.wandb.ai"},
		{
			"https://api.qa.forge.coreweave.com",
			"https://api.qa.forge.coreweave.com/evaluations",
		},
		{"https://wandb.example.com", "https://wandb.example.com/evaluations"},
		{"https://example.com/wandb/", "https://example.com/wandb/evaluations"},
	}

	for _, tc := range testCases {
		t.Run(tc.wandbBaseURL, func(t *testing.T) {
			t.Setenv("CES_BASE_URL", "")

			cesURL, err := evaltables.CESBaseURL(mustParseURL(t, tc.wandbBaseURL))

			require.NoError(t, err)
			assert.Equal(t, tc.want, cesURL.String())
		})
	}
}

func TestCESBaseURL_EnvironmentOverride(t *testing.T) {
	t.Setenv("CES_BASE_URL", "https://ces.example.com")

	cesURL, err := evaltables.CESBaseURL(mustParseURL(t, "https://api.wandb.ai"))

	require.NoError(t, err)
	assert.Equal(t, "https://ces.example.com", cesURL.String())
}

func TestEvalTableService_UsesRunCredentialsNotSDKEnvironment(t *testing.T) {
	t.Setenv("WANDB_API_KEY", "env-api-key")
	t.Setenv("CES_CUSTOM_HEADERS", "X-From-Env: yes")
	server := apitest.NewRecordingServer()
	defer server.Close()
	cesURL := mustParseURL(t, server.URL+"/evaluations")
	service := evaltables.NewEvalTableService(
		cesURL,
		evaltables.NewCESHTTPClient(
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
