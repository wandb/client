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
	// cesNamespace is the CES namespace of scopes that W&B assigns.
	//
	// W&B scopes are projects, identified by their internal ID.
	cesNamespace = "wandb"

	// cesBaseURLEnv overrides the CES URL derived from the W&B base URL.
	cesBaseURLEnv = "CES_BASE_URL"

	mtsaasWandbBaseURL = "https://api.wandb.ai"
	mtsaasCESBaseURL   = "https://evaluations.wandb.ai"
)

// forgeHosts maps each CoreWeave Forge host to the W&B API host behind it.
//
// It must match FORGE_HOSTS in wandb/sdk/lib/urls.py.
var forgeHosts = map[string]string{
	"forge.coreweave.com":    "api.wandb.ai",
	"qa.forge.coreweave.com": "api.qa.wandb.ai",
}

// CESBaseURL returns the CES URL for the W&B server at wandbBaseURL.
//
// The CES_BASE_URL environment variable takes precedence. Multi-tenant SaaS
// uses a separate host, including when reached through Forge, and other
// deployments serve CES under the W&B server's /evaluations path.
func CESBaseURL(wandbBaseURL *url.URL) (*url.URL, error) {
	if override := os.Getenv(cesBaseURLEnv); override != "" {
		return url.Parse(override)
	}

	if forgeUpstreamURL(wandbBaseURL) == mtsaasWandbBaseURL {
		return url.Parse(mtsaasCESBaseURL)
	}

	return wandbBaseURL.JoinPath("evaluations"), nil
}

// forgeUpstreamURL returns the W&B API URL behind a Forge URL, or the URL
// unchanged.
func forgeUpstreamURL(u *url.URL) string {
	hostname := u.Hostname()
	for forge, upstream := range forgeHosts {
		if hostname == forge || hostname == "api."+forge || hostname == upstream {
			return "https://" + upstream
		}
	}
	return u.String()
}

// NewCESHTTPClient returns an HTTP client for requests to CES at cesBaseURL.
//
// It authenticates requests to CES with the run's W&B credentials and retries
// them the way wandb-core retries requests to the W&B server.
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
