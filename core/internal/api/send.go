package api

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/wandb/wandb/core/internal/wboperation"
)

// Do implements RetryableClient.Do.
func (client *clientImpl) Do(req *retryablehttp.Request) (*http.Response, error) {
	// A logical request keeps one ID across HTTP retries. A new Do call gets
	// a new ID even when it reuses the caller's request object.
	req = req.WithContext(req.Context())
	req.Header = req.Header.Clone()
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	req.Header.Set("X-Wandb-Request-Id", rand.Text())
	// Track retried errors to provide a better error message at the end.
	req = req.WithContext(withRetryObserver(req.Context()))

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
