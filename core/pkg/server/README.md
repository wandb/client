# SDK service shutdown

`wandb-core` gives active runs up to five seconds to flush after SIGINT,
SIGTERM, or parent death. `--shutdown-grace-period` changes this budget; zero
restores immediate shutdown. A second signal forces an immediate stop.

Shutdown stops accepting IPC work, flushes each active transaction log's
existing buffer, waits for connections to drain, then finishes the remaining
streams with nonzero exit code 255. Finishing uses the existing upload order.
The deadline bounds shutdown even when a connection, filesystem, or upload is
stuck. It does not guarantee that all queued work reached disk or the server.

Each run's `.wandb` log is fsynced when the run writes its exit record and once
more at the end of shutdown, including when the grace period expires or a
second signal arrives. That last sync is bounded by one second so a stuck
filesystem cannot block exit.

Shutdown also updates each run's `run_end.json` in the run directory: the end
reason, the shutdown reason (`signal: <name>` or `parent_exited`), whether the
run drained in time, its exit code, and the filestream offsets the server
acknowledged. The sender writes the same file when a fatal error stops uploads
and when the run exits normally, and flow control adds a count of records it
could not read back. A missing file after a run means it never reached any of
these points, for example because of SIGKILL.

On Linux, parent-death notification uses SIGTERM so the handler can run. Parent
polling remains the fallback. This does not add Windows parent-death detection.
SIGKILL, OOM kills, and loss of the machine cannot execute shutdown code.
Between those sync points the log is only flushed, so a machine crash can still
lose its tail, and fsync on network filesystems is only as strong as the
filesystem makes it. Retained files may need `wandb sync --include-online`.

The process test covers Python logging through IPC and reads the transaction
log and `run_end.json` after SIGTERM, parent death, and a SIGKILL control. It runs when
`WANDB_TEST_CORE_BINARY` names a built service and `WANDB_TEST_PYTHON` names a
Python environment with the SDK dependencies and `psutil`.
