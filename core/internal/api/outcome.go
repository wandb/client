package api

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/wandb/wandb/core/internal/httplayers"
)

// requestOutcomeLogger records each attempt, including retries and transport
// failures. It never logs URLs, credentials, response bodies, or row contents.
type requestOutcomeLogger struct{ logger *slog.Logger }

func (o requestOutcomeLogger) WrapHTTP(send httplayers.HTTPDoFunc) httplayers.HTTPDoFunc {
	return func(req *http.Request) (*http.Response, error) {
		start := time.Now()
		resp, err := send(req)
		status := 0
		outcome := "unknown"
		if resp != nil {
			status = resp.StatusCode
			outcome = "response_received"
		}
		if err != nil {
			outcome = "transport_error"
		}
		o.logger.Info("sdk request outcome",
			"log_event_name", "sdk_request_outcome",
			"request_id", req.Header.Get("X-Wandb-Request-Id"),
			"channel", requestChannel(req.URL.Path),
			"method", req.Method,
			"status", status,
			"outcome", outcome,
			"duration_seconds", time.Since(start).Seconds(),
			"storage_visibility", "unknown",
		)
		return resp, err
	}
}

func requestChannel(path string) string {
	switch {
	case strings.HasSuffix(path, "/file_stream"):
		return "filestream"
	case strings.HasSuffix(path, "/graphql"):
		return "graphql"
	case strings.HasSuffix(path, "/sdk/otel/v1/logs"),
		strings.HasSuffix(path, "/sdk/otel/v1/metrics"):
		return "telemetry"
	default:
		return "other"
	}
}
