# Request log correlation

Existing API retry and final HTTP-error logs include a local `request_id`.
`Do` creates an ID unless the caller attached one with `WithRequestID`.
Retries retain that ID. Reusing an unmodified request in another `Do` call
creates a new ID.

Filestream attaches an ID before sending so its existing request summary,
retry/error messages, and success message can be joined. The success message
also includes `duration_seconds`, covering the HTTP call and its retries.
Successful heartbeat requests remain excluded from these INFO logs.

These changes add fields to existing messages, not per-attempt outcome logs.
IDs are not HTTP headers or deduplication keys. A successful HTTP response
does not prove that history is visible through server readback.
