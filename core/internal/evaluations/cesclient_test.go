package evaluations_test

import (
	"context"
	"errors"
	"net/http"
	"os"
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

func TestEvalTableService_UsesRunCredentials(t *testing.T) {
	server := apitest.NewRecordingServer()
	defer server.Close()
	t.Setenv("CES_BASE_URL", server.URL+"/evaluations")
	t.Setenv("WANDB_API_KEY", "env-api-key")
	t.Setenv("WANDB_ACCESS_TOKEN", "env-access-token")
	t.Setenv("CES_CUSTOM_HEADERS", "X-From-Env: yes")
	service := evaluations.NewEvalTableService(
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

type recordingDoer struct{ requests []*http.Request }

func (d *recordingDoer) Do(req *http.Request) (*http.Response, error) {
	d.requests = append(d.requests, req)
	return nil, errors.New("not sent")
}

func TestEvalTableService_DefaultsToSDKURL(t *testing.T) {
	t.Setenv("CES_BASE_URL", "")
	require.NoError(t, os.Unsetenv("CES_BASE_URL"))
	doer := &recordingDoer{}
	service := evaluations.NewEvalTableService(doer)

	_, _ = service.New(context.Background(), ces.EvalTableNewParams{
		ScopeID: ces.F("scope"),
		Name:    ces.F("table"),
	})

	require.Len(t, doer.requests, 1)
	assert.Equal(t,
		"https://evaluations.wandb.ai/v1/namespaces/wandb/scopes/scope/eval-tables",
		doer.requests[0].URL.String())
}
