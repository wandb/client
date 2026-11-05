package api_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/require"

	"github.com/wandb/wandb/core/internal/api"
	"github.com/wandb/wandb/core/internal/clients"
)

func TestRequestID_NewIDPerDo(t *testing.T) {
	headers := make(chan http.Header, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers <- r.Header.Clone()
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()
	var logs bytes.Buffer
	client := api.NewClient(
		api.ClientOptions{
			RetryPolicy: clients.RetryMostFailures,
			Logger:      slog.New(slog.NewJSONHandler(&logs, nil)),
		},
	)
	req, err := retryablehttp.NewRequest(http.MethodPost, server.URL, nil)
	require.NoError(t, err)
	req.Header = nil
	var ids []string
	for range 2 {
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
		require.Empty(t, (<-headers).Get("X-Wandb-Request-Id"))
		var entry map[string]any
		require.NoError(t, json.Unmarshal(logs.Bytes(), &entry))
		require.Equal(t, "api: HTTP error", entry["msg"])
		require.NotEmpty(t, entry["request_id"])
		ids = append(ids, entry["request_id"].(string))
		logs.Reset()
	}
	require.NotEqual(t, ids[0], ids[1])
	require.Nil(t, req.Header)
	require.Empty(t, api.RequestID(req.Context()))
}
