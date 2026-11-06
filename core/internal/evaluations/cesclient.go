// Package evaluations is a client for the CoreWeave Evaluations Service (CES).
package evaluations

import (
	"log/slog"

	ces "github.com/wandb/cwevals-go"
	"github.com/wandb/cwevals-go/option"

	"github.com/wandb/wandb/core/internal/api"
	"github.com/wandb/wandb/core/internal/clients"
	"github.com/wandb/wandb/core/internal/httplayers"
	"github.com/wandb/wandb/core/internal/settings"
)

func NewCESHTTPClient(
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
			// Not limited to a URL: this client only sends requests to CES.
			credentialProvider,
		),
	})

	return api.AsStandardClient(retryClient)
}

// NewEvalTableService returns a client for CES's EvalTable routes.
//
// It takes the CES URL, namespace and headers from the SDK's environment
// defaults, such as CES_BASE_URL. httpClient authenticates and retries
// requests, overriding the SDK's credentials and retries.
func NewEvalTableService(httpClient api.HTTPDoer) *ces.EvalTableService {
	return ces.NewClient(
		option.WithHTTPClient(httpClient),
		option.WithMaxRetries(0),
	).EvalTables
}
