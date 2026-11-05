package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/require"

	"github.com/wandb/wandb/core/internal/api"
	"github.com/wandb/wandb/core/internal/clients"
)

func TestRequestOutcomesRetry(t *testing.T) {
	for _, scenario := range []string{"429", "lost_response", "terminal"} {
		t.Run(scenario, func(t *testing.T) {
			var count atomic.Int32
			var ids []string
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ids = append(ids, r.Header.Get("X-Wandb-Request-Id"))
					if count.Add(1) == 1 {
						switch scenario {
						case "429":
							w.WriteHeader(429)
							return
						case "lost_response":
							conn, _, err := w.(http.Hijacker).Hijack()
							if err == nil {
								_ = conn.Close()
							}
							return
						case "terminal":
							w.WriteHeader(403)
							return
						}
					}
					_, _ = io.WriteString(w, `{}`)
				}),
			)
			defer server.Close()
			var logs bytes.Buffer
			client := api.NewClient(api.ClientOptions{
				RetryMax: 1, RetryPolicy: clients.RetryMostFailures,
				Logger: slog.New(slog.NewJSONHandler(&logs, nil)),
			})
			req, err := retryablehttp.NewRequest(
				http.MethodPost,
				server.URL+"/files/e/p/r/file_stream?secret=hidden",
				strings.NewReader(`{"private":"row"}`),
			)
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())
			want := 2
			if scenario == "terminal" {
				want = 1
			}
			require.EqualValues(t, want, count.Load())
			var outcomes []map[string]any
			for _, line := range strings.Split(strings.TrimSpace(logs.String()), "\n") {
				var entry map[string]any
				require.NoError(t, json.Unmarshal([]byte(line), &entry))
				if entry["log_event_name"] == "sdk_request_outcome" {
					outcomes = append(outcomes, entry)
				}
			}
			require.Len(t, outcomes, want)
			for i, outcome := range outcomes {
				require.NotEmpty(t, outcome["request_id"])
				require.Equal(t, outcomes[0]["request_id"], outcome["request_id"])
				require.Equal(t, "filestream", outcome["channel"])
				require.Equal(t, "unknown", outcome["storage_visibility"])
				require.Equal(t, ids[i], outcome["request_id"])
				encoded, err := json.Marshal(outcome)
				require.NoError(t, err)
				require.NotContains(t, string(encoded), "hidden")
				require.NotContains(t, string(encoded), "private")
			}
			if scenario == "lost_response" {
				require.Equal(t, "transport_error", outcomes[0]["outcome"])
			}
			if scenario != "terminal" {
				require.EqualValues(t, 200, outcomes[1]["status"])
			}
			require.Empty(
				t,
				req.Header.Get("X-Wandb-Request-Id"),
				"do not mutate the caller's headers",
			)
		})
	}
}

func TestRequestOutcomeNewIDPerDo(t *testing.T) {
	ids := make(chan string, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids <- r.Header.Get("X-Wandb-Request-Id")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client := api.NewClient(api.ClientOptions{})
	req, err := retryablehttp.NewRequest(http.MethodGet, server.URL+"/graphql", nil)
	require.NoError(t, err)
	req.Header = nil
	for range 2 {
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
	}
	first, second := <-ids, <-ids
	require.NotEmpty(t, first)
	require.NotEmpty(t, second)
	require.NotEqual(t, first, second)
	require.Nil(t, req.Header)
}

func TestRequestOutcomeChannels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	for path, channel := range map[string]string{
		"/sdk/otel/v1/logs":    "telemetry",
		"/sdk/otel/v1/metrics": "telemetry",
		"/graphql":             "graphql",
		"/unknown":             "other",
	} {
		t.Run(path, func(t *testing.T) {
			var logs bytes.Buffer
			client := api.NewClient(api.ClientOptions{
				Logger: slog.New(slog.NewJSONHandler(&logs, nil)),
			})
			req, err := retryablehttp.NewRequest(http.MethodPost, server.URL+path, nil)
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())
			var outcome map[string]any
			require.NoError(t, json.Unmarshal(logs.Bytes(), &outcome))
			require.Equal(t, "sdk_request_outcome", outcome["log_event_name"])
			require.Equal(t, channel, outcome["channel"])
		})
	}
}
