// Package evaltables is a client for EvalTables in the CoreWeave Evaluation Service (CES).
package evaltables

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	ces "github.com/wandb/cwevals-go"
	"github.com/wandb/cwevals-go/option"

	"github.com/wandb/wandb/core/internal/api"
	"github.com/wandb/wandb/core/internal/clients"
	"github.com/wandb/wandb/core/internal/httplayers"
	"github.com/wandb/wandb/core/internal/settings"
)

// NewEvalTableService returns a client for the EvalTables API at the
// settings' CES URL.
//
// It authenticates with credentialProvider instead of the SDK's environment
// variables, and retries with wandb-core's HTTP client instead of the SDK's.
func NewEvalTableService(
	credentialProvider api.CredentialProvider,
	logger *slog.Logger,
	s *settings.Settings,
) (*ces.EvalTableService, error) {
	cesBaseURL, err := parseCESBaseURL(s)
	if err != nil {
		return nil, err
	}

	return ces.NewClient(
		option.WithBaseURL(cesBaseURL.String()),
		option.WithHTTPClient(
			newCESHTTPClient(cesBaseURL, credentialProvider, logger, s)),
		option.WithMaxRetries(0),
		// Drop the SDK's environment credentials; the HTTP client adds the
		// run's, but only for requests to the CES URL.
		option.WithHeaderDel("Authorization"),
	).EvalTables, nil
}

func parseCESBaseURL(s *settings.Settings) (*url.URL, error) {
	rawURL := s.GetCESBaseURL()
	if rawURL == "" {
		return nil, errors.New("evaltables: no CES URL in settings")
	}

	cesBaseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("evaltables: invalid CES URL: %v", err)
	}

	return cesBaseURL, nil
}

func newCESHTTPClient(
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
