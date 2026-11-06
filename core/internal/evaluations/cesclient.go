// Package evaluations is a client for the CoreWeave Evaluation Service (CES).
package evaluations

import (
	"log/slog"
	"net/url"

	ces "github.com/wandb/cwevals-go"
	"github.com/wandb/cwevals-go/option"

	"github.com/wandb/wandb/core/internal/api"
	"github.com/wandb/wandb/core/internal/clients"
	"github.com/wandb/wandb/core/internal/httplayers"
	"github.com/wandb/wandb/core/internal/settings"
)

func NewCESHTTPClient(
	cesBaseURL *url.URL,
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
			httplayers.LimitTo(cesBaseURL, credentialProvider),
		),
	})

	return api.AsStandardClient(retryClient)
}

func NewEvalTableService(
	cesBaseURL *url.URL,
	httpClient api.HTTPDoer,
) *ces.EvalTableService {
	return ces.NewClient(
		option.WithBaseURL(cesBaseURL.String()),
		option.WithHTTPClient(httpClient),
		option.WithMaxRetries(0),
	).EvalTables
}
