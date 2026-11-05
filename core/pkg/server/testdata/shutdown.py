import json
import os
import pathlib
import time

import psutil

import wandb
import wandb.sdk.lib.service.service_process as service_process

service_process.get_core_path = lambda: os.environ["WANDB_TEST_CORE_BINARY"]
if os.environ.get("M2_OS_PID_SHUTDOWN") == "1":
    original_popen = service_process.subprocess.Popen

    def popen_with_parent_signal(command, *args, **kwargs):
        if command[0] == os.environ["WANDB_TEST_CORE_BINARY"]:
            command = [*command, "--os-pid-shutdown"]
        return original_popen(command, *args, **kwargs)

    service_process.subprocess.Popen = popen_with_parent_signal

r = wandb.init(
    mode="offline",
    project="m2-shutdown",
    dir=os.environ["M2_DIR"],
    settings=wandb.Settings(
        x_transaction_log_flush_interval=3600, x_disable_stats=True, console="off"
    ),
)
for i in range(25):
    r.log({"m2_row": i})
# A round trip after logging ensures the core has consumed the IPC messages.
r._interface.deliver_get_summary().wait_or(timeout=10)
children = [
    p for p in psutil.Process().children(recursive=True) if "wandb-core" in p.name()
]
pathlib.Path(os.environ["M2_DIR"], "ready.json").write_text(
    json.dumps({"pid": children[0].pid, "file": r.settings.sync_file})
)
while True:
    time.sleep(1)
