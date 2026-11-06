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
			// Unlike other clients, not limited to a URL: only the CES SDK
			// uses this client, and only for its base URL.
			credentialProvider,
		),
	})

	return api.AsStandardClient(retryClient)
}

// NewEvalTableService returns a CES EvalTable client that authenticates and
// retries with httpClient instead of the SDK's environment credentials and
// retries.
func NewEvalTableService(
	cesBaseURL string,
	httpClient api.HTTPDoer,
) *ces.EvalTableService {
	return ces.NewClient(
		option.WithBaseURL(cesBaseURL),
		option.WithHTTPClient(httpClient),
		option.WithMaxRetries(0),
	).EvalTables
}
