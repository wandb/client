// Package evaluations is a client for the CoreWeave Evaluations Service (CES).
package evaluations

import (
	"log/slog"
	"net/url"
	"os"

	ces "github.com/wandb/cwevals-go"
	"github.com/wandb/cwevals-go/option"

	"github.com/wandb/wandb/core/internal/api"
	"github.com/wandb/wandb/core/internal/clients"
	"github.com/wandb/wandb/core/internal/httplayers"
	"github.com/wandb/wandb/core/internal/settings"
)

const (
	cesNamespace  = "wandb"
	cesBaseURLEnv = "CES_BASE_URL"

	// defaultCESBaseURL matches the SDK's option.WithEnvironmentProduction.
	defaultCESBaseURL = "https://evaluations.wandb.ai/"
)

// CESBaseURL returns CES_BASE_URL, or the SDK's default CES URL if it's unset.
func CESBaseURL() (*url.URL, error) {
	if override := os.Getenv(cesBaseURLEnv); override != "" {
		return url.Parse(override)
	}
	return url.Parse(defaultCESBaseURL)
}

func NewCESHTTPClient(
	cesBaseURL *url.URL,
	credentialProvider api.CredentialProvider,
	logger *slog.Logger,
	s *settings.Settings,
) api.HTTPDoer {
	retryClient := api.NewClient(api.ClientOptions{
		// Retries 503 publication_failed and 429, but not 409
		// idempotency_key_reuse or validation errors.
		RetryPolicy:        clients.RetryMostFailures,
		RetryMax:           api.DefaultRetryMax,
		RetryWaitMin:       api.DefaultRetryWaitMin,
		RetryWaitMax:       api.DefaultRetryWaitMax,
		NonRetryTimeout:    api.DefaultNonRetryTimeout,
		Proxy:              s.GetProxyFn(),
		ProxyConnectHeader: s.GetProxyConnectHeader(),
		InsecureDisableSSL: s.IsInsecureDisableSSL(),
		Logger:             logger,
		PreRetryLayers: httplayers.Concat(
			httplayers.DefaultHeaders(s.GetExtraHTTPHeaders()),
			httplayers.LimitTo(cesBaseURL, credentialProvider),
		),
	})

	return api.AsStandardClient(retryClient)
}

// NewEvalTableService returns a client for CES's EvalTable routes.
//
// Unlike ces.NewClient, it ignores the SDK's environment defaults such as
// WANDB_API_KEY and CES_CUSTOM_HEADERS: httpClient authenticates and retries
// requests instead.
func NewEvalTableService(
	cesBaseURL *url.URL,
	httpClient api.HTTPDoer,
) *ces.EvalTableService {
	return ces.NewEvalTableService(
		option.WithBaseURL(cesBaseURL.String()),
		option.WithHTTPClient(httpClient),
		option.WithNamespace(cesNamespace),
		option.WithMaxRetries(0),
	)
}
