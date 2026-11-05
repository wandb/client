# API request outcomes

`NewClient` emits a local `sdk_request_outcome` log for each HTTP attempt,
including retries and transport failures. `Do` assigns a random
`X-Wandb-Request-Id` to a cloned request header map. Retries keep that ID;
another `Do` call gets a new ID. This header is for correlation, not deduplication.

The log includes channel, method, HTTP status, duration, and whether the
transport received a response or returned an error. It excludes URLs, bodies,
and credentials. Storage visibility remains unknown even after 2xx; a lost
response also leaves acceptance unknown. GraphQL application errors within
HTTP 200 are not classified by this transport log.

Coverage is limited to requests sent through `NewClient`. These local logs do
not provide durable row accounting, direct-upload coverage, or a shutdown
report, and are not automatically shipped to the backend.
