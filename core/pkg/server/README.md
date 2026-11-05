# SDK service shutdown

`wandb-core` gives active runs up to five seconds to flush after SIGINT,
SIGTERM, or parent death. `--shutdown-grace-period` changes this budget; zero
restores immediate shutdown. A second signal forces an immediate stop.

Shutdown stops accepting IPC work, flushes each active transaction log's
existing buffer, waits for connections to drain, then finishes the remaining
streams with nonzero exit code 255. Finishing uses the existing upload order.
The deadline bounds shutdown even when a connection, filesystem, or upload is
stuck. It does not guarantee that all queued work reached disk or the server.

On Linux, parent-death notification uses SIGTERM so the handler can run. Parent
polling remains the fallback. This does not add Windows parent-death detection.
SIGKILL, OOM kills, and loss of the machine cannot execute shutdown code.
Transaction-log flushing does not promise power-loss durability or remote
storage visibility. Retained files may need `wandb sync --include-online`.

The process test covers Python logging through IPC and reads the transaction
log after SIGTERM, parent death, and a SIGKILL control. It runs when
`WANDB_TEST_CORE_BINARY` names a built service and `WANDB_TEST_PYTHON` names a
Python environment with the SDK dependencies and `psutil`.
