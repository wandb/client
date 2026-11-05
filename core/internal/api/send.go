package api

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/wandb/wandb/core/internal/wboperation"
)

// Do implements RetryableClient.Do.
func (client *clientImpl) Do(req *retryablehttp.Request) (*http.Response, error) {
	ctx := req.Context()
	if RequestID(ctx) == "" {
		ctx = WithRequestID(ctx)
	}
	// Track retried errors to provide a better error message at the end.
	req = req.WithContext(withRetryObserver(ctx))

	resp, err := client.retryableHTTP.Do(req)
	wboperation.Get(req.Context()).ClearError()

	if err != nil {
		if lastStatus := lastRetriedError(req.Context()); lastStatus != "" {
			return nil, &RetryError{Inner: err, LastStatus: lastStatus}
		} else {
			return nil, err
		}
	}

	// This is a bug that happens with retryablehttp sometimes.
	if resp == nil {
		return nil, fmt.Errorf("api: nil error and nil response")
	}

	client.logFinalResponseOnError(req, resp)
	return resp, nil
}
