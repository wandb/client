package filestream

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/wandb/wandb/core/internal/api"
	"github.com/wandb/wandb/core/internal/clients"
	"github.com/wandb/wandb/core/internal/observability"
	"github.com/wandb/wandb/core/internal/settings"
	"github.com/wandb/wandb/core/internal/wboperation"
	spb "github.com/wandb/wandb/core/pkg/service_go_proto"
)

func TestRequestLogging(t *testing.T) {
	const sending = "filestream: sending request"
	const sent = "filestream: request sent"
	for _, scenario := range []struct {
		name      string
		heartbeat bool
		failure   string
		attempts  int
		messages  []string
	}{
		{"data", false, "", 1, []string{sending, sent}},
		{"429", false, "429", 2, []string{sending, "api: retrying HTTP error", sent}},
		{"lost_response", false, "lost", 2, []string{sending, "api: retrying error", sent}},
		{"403", false, "403", 1, []string{sending, "api: HTTP error"}},
		{"heartbeat", true, "", 1, nil},
		{"heartbeat_retry", true, "429", 2, []string{"api: retrying HTTP error"}},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			var count atomic.Int32
			headers := make(chan http.Header, 2)
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					headers <- r.Header.Clone()
					_, _ = io.Copy(io.Discard, r.Body)
					if count.Add(1) == 1 {
						switch scenario.failure {
						case "429":
							w.WriteHeader(429)
							return
						case "403":
							w.WriteHeader(403)
							return
						case "lost":
							conn, _, err := w.(http.Hijacker).Hijack()
							if err == nil {
								_ = conn.Close()
							}
							return
						}
					}
					_, _ = io.WriteString(w, `{}`)
				}),
			)
			defer server.Close()
			var logs bytes.Buffer
			logger := observability.NewCoreLogger(slog.New(slog.NewJSONHandler(&logs, nil)), nil)
			baseURL, err := url.Parse(server.URL)
			require.NoError(t, err)
			fs := fileStream{
				baseURL:         baseURL,
				path:            "files/e/p/r/file_stream",
				beforeRunEndCtx: context.Background(),
				settings: &settings.Settings{
					Proto: &spb.Settings{XFileStreamNoGzip: wrapperspb.Bool(true)},
				},
				logger:     logger,
				operations: wboperation.NewOperations(),
				apiClient: api.NewClient(
					api.ClientOptions{
						Logger:       logger.Logger,
						RetryMax:     1,
						RetryPolicy:  clients.RetryMostFailures,
						RetryWaitMin: time.Millisecond,
						RetryWaitMax: time.Millisecond,
					},
				),
			}
			data := &FileStreamRequestJSON{}
			if !scenario.heartbeat {
				data.Files = map[string]OffsetAndContent{
					HistoryFileName: {Offset: 7, Content: []string{`{"private":"synthetic-row"}`}},
				}
			}
			start := time.Now()
			err = fs.send(data, make(chan map[string]any, 1))
			elapsed := time.Since(start)
			if scenario.failure == "403" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, scenario.attempts, count.Load())
			for range scenario.attempts {
				require.Empty(t, (<-headers).Get("X-Wandb-Request-Id"))
			}
			var entries []map[string]any
			for _, line := range strings.Split(strings.TrimSpace(logs.String()), "\n") {
				if line == "" {
					continue
				}
				var entry map[string]any
				require.NoError(t, json.Unmarshal([]byte(line), &entry))
				require.NotContains(t, line, "synthetic-row")
				entries = append(entries, entry)
			}
			var messages []string
			for _, entry := range entries {
				messages = append(messages, entry["msg"].(string))
				require.NotEmpty(t, entry["request_id"])
				require.Equal(t, entries[0]["request_id"], entry["request_id"])
			}
			require.Equal(t, scenario.messages, messages)
			if !scenario.heartbeat {
				require.EqualValues(t, 7, entries[0]["history_offset"])
				require.EqualValues(t, 1, entries[0]["history_lines"])
				if scenario.failure != "403" {
					success := entries[len(entries)-1]
					require.Equal(t, "200 OK", success["status"])
					duration := success["duration_seconds"].(float64)
					require.Positive(t, duration)
					require.LessOrEqual(t, duration, elapsed.Seconds())
					if scenario.attempts > 1 {
						require.GreaterOrEqual(t, duration, time.Millisecond.Seconds())
					}
				}
			}
		})
	}
}
