# File is generated by: tox -e codemod

import wandb
from wandb.proto.wandb_telemetry_pb2 import Imports as TelemetryImports
from wandb.proto.wandb_telemetry_pb2 import TelemetryRecord

if wandb.TYPE_CHECKING:  # type: ignore
    from typing import ContextManager, Type, Optional
    from types import TracebackType

    # avoid cycle, use string type reference
    from typing import TYPE_CHECKING

    if TYPE_CHECKING:
        from .. import wandb_run


class _TelemetryObject(object):
    # _run: Optional["wandb_run.Run"]

    def __init__(self, run = None):
        self._run = run or wandb.run
        self._obj = TelemetryRecord()

    def __enter__(self):
        return self._obj

    def __exit__(
        self,
        exctype,
        excinst,
        exctb,
    ):
        if not self._run:
            return
        self._run._telemetry_callback(self._obj)


def context(run = None):
    return _TelemetryObject(run=run)


__all__ = [
    "TelemetryImports",
    "TelemetryRecord",
    "context",
]
