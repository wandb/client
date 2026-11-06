package evaluations_test

import (
	"context"
	"net/http"
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

func TestEvalTableService_UsesGivenURLAndRunCredentials(t *testing.T) {
	server := apitest.NewRecordingServer()
	defer server.Close()
	t.Setenv("CES_BASE_URL", "https://ces.example.com")
	t.Setenv("WANDB_API_KEY", "env-api-key")
	t.Setenv("WANDB_ACCESS_TOKEN", "env-access-token")
	t.Setenv("CES_CUSTOM_HEADERS", "X-From-Env: yes")
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
	assert.Equal(t, "yes", req.Header.Get("X-From-Env"))
}

func TestCESHTTPClient_OnlyAuthenticatesCESRequests(t *testing.T) {
	cesServer := apitest.NewRecordingServer()
	defer cesServer.Close()
	otherServer := apitest.NewRecordingServer()
	defer otherServer.Close()
	client := evaluations.NewCESHTTPClient(
		mustParseURL(t, cesServer.URL+"/evaluations"),
		api.NewAPIKeyCredentialProvider("run-api-key"),
		observabilitytest.NewTestLogger(t).Logger,
		settings.New(),
	)

	for _, rawURL := range []string{
		cesServer.URL + "/evaluations/v1",
		cesServer.URL + "/graphql",
		otherServer.URL + "/evaluations/v1",
	} {
		req, err := http.NewRequest(http.MethodGet, rawURL, http.NoBody)
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
	}

	cesRequests := cesServer.Requests()
	require.Len(t, cesRequests, 2)
	assert.NotEmpty(t, cesRequests[0].Header.Get("Authorization"))
	assert.Empty(t, cesRequests[1].Header.Get("Authorization"))
	otherRequests := otherServer.Requests()
	require.Len(t, otherRequests, 1)
	assert.Empty(t, otherRequests[0].Header.Get("Authorization"))
}

func mustParseURL(t *testing.T, rawURL string) *url.URL {
	t.Helper()
	u, err := url.Parse(rawURL)
	require.NoError(t, err)
	return u
}
