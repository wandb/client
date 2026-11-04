"""Unit tests for the CES-backed wandb.EvalTable."""

from __future__ import annotations

import json
import sys
from types import ModuleType
from typing import Any
from unittest.mock import MagicMock

import pytest
import wandb
from google.protobuf.message import Message
from wandb.errors import CommError, UsageError
from wandb.proto import wandb_internal_pb2 as pb
from wandb.sdk.data_types._dtypes import AnyType
from wandb.sdk.data_types.eval_table import _writer_ces as ces
from wandb.sdk.interface.interface import InterfaceBase


@pytest.fixture
def run(mock_run):
    return mock_run(settings={"entity": "e", "project": "p", "mode": "online"})


@pytest.fixture(autouse=True)
def default_eval_table_server_feature_enabled(monkeypatch):
    monkeypatch.setattr(
        "wandb.apis.public.service_api.ServiceApi.feature_enabled",
        lambda self, feature: True,
    )


@pytest.fixture(autouse=True)
def coreweave_evaluations_module(monkeypatch):
    client_module = ModuleType("coreweave_evaluations")
    client_module.__path__ = []
    types_module = ModuleType("coreweave_evaluations.types")
    types_module.__path__ = []
    image_module = ModuleType("coreweave_evaluations.types.wandb_image_v1_param")
    image_module.WandbImageV1Param = dict

    monkeypatch.setitem(sys.modules, "coreweave_evaluations", client_module)
    monkeypatch.setitem(sys.modules, "coreweave_evaluations.types", types_module)
    monkeypatch.setitem(
        sys.modules,
        "coreweave_evaluations.types.wandb_image_v1_param",
        image_module,
    )


def _log_eval_table_result(
    response: pb.LogEvalTableResponse,
) -> pb.Result:
    return pb.Result(response=pb.Response(log_eval_table_response=response))


_UPLOADED = pb.LogEvalTableResponse(
    evaluation_id="evaluation-1",
    evaluation_version_id="evaluation-version-1",
    dataset_id="dataset-1",
    dataset_version_id="dataset-version-1",
)


@pytest.fixture
def deliver_log_eval_table(monkeypatch) -> MagicMock:
    """Replace wandb-core's EvalTable upload with one that succeeds."""
    deliver = MagicMock()
    deliver.return_value.wait_or.return_value = _log_eval_table_result(_UPLOADED)
    monkeypatch.setattr(InterfaceBase, "deliver_log_eval_table", deliver)
    return deliver


def _sent(deliver: MagicMock) -> list[pb.LogEvalTableRequest]:
    """Return the requests sent to wandb-core."""
    return [call.args[0] for call in deliver.call_args_list]


def _rows(request: pb.LogEvalTableRequest) -> list[dict[str, Any]]:
    return [json.loads(row) for row in request.rows_json]


def _as_dicts(messages: list[Message]) -> list[dict[str, Any]]:
    """Convert messages to dicts of their set fields."""
    return [
        {field.name: value for field, value in message.ListFields()}
        for message in messages
    ]


def test_ces_eval_table_writes_columns_rows_and_version(
    deliver_log_eval_table,
    run,
    monkeypatch,
):
    debug = MagicMock()
    monkeypatch.setattr(
        "wandb.sdk.data_types.eval_table._writer_ces._logger.debug",
        debug,
    )
    et = wandb.EvalTable(
        columns=["prompt", "truth", "answer", "confidence", "correct"],
        data=[
            ["2+2", "4", "4", 1, True],
            ["2+3", "5", "4", 0.25, False],
        ],
        input_columns=["prompt", "truth"],
        output_columns=["answer", "confidence"],
        score_columns=["correct"],
        backend="ces",
    )

    run.log({"math_eval": et})

    [request] = _sent(deliver_log_eval_table)
    assert request.name == "math_eval"
    assert request.history_key == "math_eval"
    assert request.step == "0"
    assert _as_dicts(request.dataset_fields) == [
        {"source": "input", "name": "prompt", "value_type": "string"},
        {"source": "input", "name": "truth", "value_type": "string"},
        {"source": "output", "name": "answer", "value_type": "string"},
        {"source": "output", "name": "confidence", "value_type": "number"},
    ]
    assert _as_dicts(request.scorers) == [
        {"name": "correct", "value_type": "boolean"},
    ]
    assert _rows(request) == [
        {
            "input": {"prompt": "2+2", "truth": "4"},
            "output": {"answer": "4", "confidence": 1},
            "scores": {"correct": True},
        },
        {
            "input": {"prompt": "2+3", "truth": "5"},
            "output": {"answer": "4", "confidence": 0.25},
            "scores": {"correct": False},
        },
    ]
    debug.assert_called_once_with(
        "CES EvalTable recorded evaluation_version_id=%s",
        "evaluation-version-1",
    )

    history_value = et.to_json(run)
    assert history_value == {
        "_type": "eval-table-ces",
        "schema_version": 1,
        "ncols": 5,
        "nrows": 2,
        "log_mode": "IMMUTABLE",
        "evaluation_id": "evaluation-1",
        "evaluation_version_id": "evaluation-version-1",
        "dataset_id": "dataset-1",
        "dataset_version_id": "dataset-version-1",
    }
    assert "evaluate_call_id" not in history_value
    assert et._immutable_write_result is not None
    assert et._immutable_write_result.logged_id == "evaluation-version-1"


def test_ces_eval_table_without_inputs_keys_rows_by_position(
    deliver_log_eval_table,
    run,
):
    et = wandb.EvalTable(
        columns=["answer", "correct"],
        data=[["4", True], ["5", False]],
        score_columns=["correct"],
        backend="ces",
    )

    run.log({"eval": et})

    [request] = _sent(deliver_log_eval_table)
    assert _rows(request) == [
        {
            "input": None,
            "output": {"answer": "4"},
            "scores": {"correct": True},
            "example_key": "1",
        },
        {
            "input": None,
            "output": {"answer": "5"},
            "scores": {"correct": False},
            "example_key": "2",
        },
    ]


def test_ces_eval_table_stubs_media_until_native_support_exists(
    deliver_log_eval_table,
    mock_wandb_log,
    run,
):
    histogram = wandb.Histogram([1, 2, 3])
    table = wandb.EvalTable(
        columns=["histogram"],
        data=[[histogram]],
        backend="ces",
    )

    run.log({"media_eval": table})

    [request] = _sent(deliver_log_eval_table)
    rows = _rows(request)
    assert rows[0]["output"]["histogram"] == "[wandb.Histogram not yet supported]"
    mock_wandb_log.assert_warned("wandb.Histogram values are not yet supported")


def test_ces_eval_table_writes_supported_media(
    deliver_log_eval_table,
    run,
):
    from PIL import Image as PILImage

    image = wandb.Image(PILImage.new("RGB", (2, 2), color="red"))
    et = wandb.EvalTable(
        columns=["image"],
        data=[[image]],
        backend="ces",
    )

    run.log({"media_eval": et})

    [request] = _sent(deliver_log_eval_table)
    rows = _rows(request)
    image_value = rows[0]["output"]["image"]
    assert image_value["extension_type"] == "wandb-image"
    assert image_value["uri"].startswith("wandb-run-file://")


@pytest.mark.parametrize(
    ("oversized_media_cells", "expected_metrics"),
    [
        (0, ["eval_table_ces_media_write"]),
        (
            1,
            [
                "eval_table_ces_media_write",
                "eval_table_ces_media_write_with_oversized_cells",
            ],
        ),
    ],
)
def test_ces_eval_table_media_telemetry_counts_affected_writes(
    monkeypatch,
    oversized_media_cells,
    expected_metrics,
):
    recorder = MagicMock()
    monkeypatch.setattr(ces, "get_telemetry_recorder", lambda: recorder)
    prepared = ces._CESWritePayloads(
        dataset_fields=[],
        scorers=[],
        rows_json=[],
        media_cells_examined=1,
        oversized_media_cells=oversized_media_cells,
        oversized_locations=(),
    )

    ces.CESWriter()._record_media_telemetry(prepared)

    assert [
        call.args[0] for call in recorder.increment_counter.call_args_list
    ] == expected_metrics
    recorder.log.assert_not_called()


def test_ces_eval_table_writes_supported_media_in_raise_mode(
    deliver_log_eval_table, run
):
    from PIL import Image as PILImage

    image = wandb.Image(PILImage.new("RGB", (2, 2), color="red"))

    table = wandb.EvalTable(
        columns=["image"],
        data=[[image]],
        backend="ces",
        unsupported_media_mode="raise",
    )
    run.log({"media_eval": table})

    [request] = _sent(deliver_log_eval_table)
    rows = _rows(request)
    assert rows[0]["output"]["image"]["extension_type"] == "wandb-image"


def test_ces_writer_validates_cells_during_construction():
    nested_table = wandb.Table(columns=["value"], data=[[1]])

    with pytest.raises(TypeError, match="does not support nested Tables"):
        wandb.EvalTable(
            columns=["table"],
            data=[[nested_table]],
            backend="ces",
        )


def test_ces_eval_table_rejects_oversized_row_before_network(
    deliver_log_eval_table,
    run,
    monkeypatch,
):
    value = "large" * 100
    row = {
        "input": None,
        "output": {"value": value},
        "scores": {},
        "example_key": "1",
    }
    body_size = len(ces._encode_json({"rows": [row]}))
    monkeypatch.setattr(ces, "_MAX_REQUEST_BODY_BYTES", body_size)
    et = wandb.EvalTable(
        columns=["value"],
        data=[[value]],
        output_columns=["value"],
        backend="ces",
    )

    with pytest.raises(UsageError, match="row at index 0"):
        run.log({"eval": et})

    deliver_log_eval_table.assert_not_called()


def test_ces_eval_table_rejects_total_row_count_before_network(
    deliver_log_eval_table,
    run,
    monkeypatch,
):
    monkeypatch.setattr(ces, "_MAX_ROWS_PER_TABLE", 2)
    et = wandb.EvalTable(
        columns=["value"],
        data=[[1], [2], [3]],
        output_columns=["value"],
        backend="ces",
    )

    with pytest.raises(UsageError, match="at most 2 rows per table"):
        run.log({"eval": et})

    deliver_log_eval_table.assert_not_called()


def test_ces_eval_table_raises_upload_error(deliver_log_eval_table, run):
    deliver_log_eval_table.return_value.wait_or.return_value = _log_eval_table_result(
        pb.LogEvalTableResponse(error_message="upload failed")
    )
    et = wandb.EvalTable(
        columns=["value"],
        data=[[1]],
        output_columns=["value"],
        backend="ces",
    )

    with pytest.raises(CommError, match="upload failed"):
        run.log({"eval": et})

    assert et._immutable_write_result is None


def test_ces_eval_table_infers_python_and_numpy_integers(
    deliver_log_eval_table,
    run,
):
    np = pytest.importorskip("numpy")
    et = wandb.EvalTable(
        columns=["python_int", "numpy_int", "numeric"],
        data=[[1, np.int64(2), np.int64(3)], [4, np.int32(5), np.float64(6.5)]],
        output_columns=["python_int", "numpy_int", "numeric"],
        backend="ces",
    )

    run.log({"typed_eval": et})

    [request] = _sent(deliver_log_eval_table)
    assert _as_dicts(request.dataset_fields) == [
        {"source": "output", "name": "python_int", "value_type": "integer"},
        {"source": "output", "name": "numpy_int", "value_type": "integer"},
        {"source": "output", "name": "numeric", "value_type": "number"},
    ]


def test_ces_eval_table_classifies_integer_valued_floats_as_numbers(
    deliver_log_eval_table,
    run,
):
    et = wandb.EvalTable(
        columns=["value"],
        data=[[1.0], [2.0], [3.0]],
        output_columns=["value"],
        backend="ces",
    )

    run.log({"typed_eval": et})

    [request] = _sent(deliver_log_eval_table)
    assert _as_dicts(request.dataset_fields) == [
        {"source": "output", "name": "value", "value_type": "number"},
    ]


def test_ces_eval_table_serializes_nan_as_typed_null(deliver_log_eval_table, run):
    et = wandb.EvalTable(
        columns=["value"],
        data=[[float("nan")]],
        output_columns=["value"],
        backend="ces",
    )

    run.log({"nan_eval": et})

    [request] = _sent(deliver_log_eval_table)
    assert _as_dicts(request.dataset_fields) == [
        {"source": "output", "name": "value", "value_type": "number"},
    ]
    assert _rows(request)[0]["output"] == {"value": None}


@pytest.mark.parametrize(
    ("data", "message"),
    [
        ([[None], [None]], "all-null"),
        ([[float("inf")]], "non-finite"),
        ([[{"nested": True}]], "only primitive values"),
    ],
)
def test_ces_eval_table_rejects_invalid_columns_before_network(
    deliver_log_eval_table,
    run,
    data,
    message,
):
    et = wandb.EvalTable(
        columns=["value"],
        data=data,
        backend="ces",
    )
    with pytest.raises(UsageError, match=message):
        run.log({"invalid_eval": et})

    deliver_log_eval_table.assert_not_called()


def test_ces_eval_table_rejects_mixed_column_types_on_construction():
    with pytest.raises(TypeError, match="incompatible types"):
        wandb.EvalTable(
            columns=["value"],
            data=[["x"], [1]],
            backend="ces",
        )


def test_ces_eval_table_rejects_mixed_types_with_permissive_dtype_before_network(
    deliver_log_eval_table,
    run,
):
    et = wandb.EvalTable(
        columns=["value"],
        data=[[True], [1]],
        dtype=AnyType,
        backend="ces",
    )

    with pytest.raises(UsageError, match="mixes 'boolean' and 'integer'"):
        run.log({"mixed_eval": et})

    deliver_log_eval_table.assert_not_called()


def test_ces_eval_table_reports_extension_types_in_mixed_column_error():
    writer = ces.CESWriter()

    with pytest.raises(UsageError, match="mixes 'wandb-image' and 'wandb-audio'"):
        writer._merge_field_type(
            "media",
            ces._CESFieldType("json", "wandb-image", 1),
            ces._CESFieldType("json", "wandb-audio", 1),
        )


@pytest.mark.parametrize(
    (
        "columns",
        "input_columns",
        "output_columns",
        "score_columns",
        "message",
    ),
    [
        (["x" * 513], ["x" * 513], None, None, "input column names"),
        (["x" * 513], None, ["x" * 513], None, "output column names"),
        (["x" * 257], None, None, ["x" * 257], "score column names"),
    ],
)
def test_ces_eval_table_rejects_invalid_column_name_lengths_before_network(
    deliver_log_eval_table,
    run,
    columns,
    input_columns,
    output_columns,
    score_columns,
    message,
):
    et = wandb.EvalTable(
        columns=columns,
        data=[[1]],
        input_columns=input_columns,
        output_columns=output_columns,
        score_columns=score_columns,
        backend="ces",
    )

    with pytest.raises(UsageError, match=message):
        run.log({"invalid_eval": et})

    deliver_log_eval_table.assert_not_called()


def test_ces_eval_table_describes_field_limit_as_input_and_output_columns(
    monkeypatch,
    deliver_log_eval_table,
    run,
):
    monkeypatch.setattr(ces, "_MAX_DATASET_FIELDS", 1)
    et = wandb.EvalTable(
        columns=["prompt", "answer"],
        data=[["question", "response"]],
        input_columns=["prompt"],
        output_columns=["answer"],
        backend="ces",
    )

    with pytest.raises(
        UsageError,
        match="at most 1 input and output columns combined",
    ):
        run.log({"invalid_eval": et})

    deliver_log_eval_table.assert_not_called()


def test_ces_eval_table_rejects_invalid_table_name_length_before_network(
    deliver_log_eval_table,
    run,
):
    et = wandb.EvalTable(columns=["value"], data=[[1]], backend="ces")

    with pytest.raises(UsageError, match="EvalTable names"):
        run.log({"x" * 257: et})

    deliver_log_eval_table.assert_not_called()


def test_ces_error_uses_original_integer_column(deliver_log_eval_table, run):
    et = wandb.EvalTable(
        columns=[3],
        data=[[float("inf")]],
        backend="ces",
    )

    with pytest.raises(UsageError, match="column 3") as exc_info:
        run.log({"invalid_eval": et})

    assert "column '3'" not in str(exc_info.value)
    deliver_log_eval_table.assert_not_called()


def test_ces_eval_table_retry_resends_log_location(deliver_log_eval_table, run):
    deliver_log_eval_table.return_value.wait_or.side_effect = [
        _log_eval_table_result(pb.LogEvalTableResponse(error_message="failed")),
        _log_eval_table_result(_UPLOADED),
    ]
    et = wandb.EvalTable(
        columns=["value"],
        data=[[1], [2]],
        output_columns=["value"],
        backend="ces",
    )
    et.bind_to_run(run, "eval", 7)

    with pytest.raises(CommError, match="failed"):
        et.to_json(run)
    et.to_json(run)

    # wandb-core derives idempotency keys from the log location, so a retry
    # reuses the CES resources created by the failed attempt.
    first, retry = _sent(deliver_log_eval_table)
    assert first == retry
    assert (first.history_key, first.step) == ("eval", "7")


def test_ces_eval_table_rejects_mixed_type_mode():
    with pytest.raises(UsageError, match="allow_mixed_types=False"):
        wandb.EvalTable(
            columns=["value"],
            data=[[1]],
            allow_mixed_types=True,
            backend="ces",
        )
