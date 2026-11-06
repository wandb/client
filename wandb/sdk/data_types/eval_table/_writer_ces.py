from __future__ import annotations

import gzip
import logging
import math
import os
import tempfile
from collections.abc import Mapping, Sequence
from dataclasses import dataclass
from typing import TYPE_CHECKING, Any, Literal

from typing_extensions import NotRequired, TypedDict

import wandb
from wandb.analytics import get_telemetry_recorder
from wandb.analytics.opentelemetry.opentelemetry_proxy import LowCardinalityAttributes
from wandb.errors import CommError, UsageError
from wandb.proto import wandb_internal_pb2 as pb
from wandb.sdk.data_types.base_types.media import Media
from wandb.sdk.data_types.base_types.wb_value import WBValue
from wandb.sdk.data_types.eval_table import _media_ces
from wandb.sdk.data_types.eval_table._media_ces import _encode_json
from wandb.sdk.data_types.eval_table._writer import WriteResult, WriteRow
from wandb.sdk.data_types.table import Table

if TYPE_CHECKING:
    from wandb.sdk.data_types.table import ColumnKey, LogMode
    from wandb.sdk.wandb_run import Run as LocalRun


_logger = logging.getLogger(__name__)

# Server request-body limit.
_MAX_REQUEST_BODY_BYTES = 16 << 20
_MAX_REQUEST_BODY_SIZE = "16 MiB"
# Tables with more encoded row characters than this send their rows to
# wandb-core in a file, keeping requests far below the 2 GiB message limit.
_ROWS_FILE_THRESHOLD_CHARS = 64 << 20
# Evaluations request-schema limits.
_MAX_DATASET_FIELDS = 10_000
_MAX_DATASET_FIELD_NAME_LENGTH = 512
_MAX_EVAL_TABLE_NAME_LENGTH = 256
# Allow ten count-limited batches while failing fast on runaway tables.
_MAX_ROWS_PER_TABLE = 10_000
_MAX_SCORERS = 256
_MAX_SCORER_NAME_LENGTH = 256
_MAX_OVERSIZED_MEDIA_LOCATIONS = 5

PrimitiveValueType = Literal["boolean", "integer", "number", "string"]
_CESFieldSource = Literal["input", "output"]


class _CESDatasetField(TypedDict):
    """A CES dataset field declaration, named like `pb.EvalTableField`."""

    source: _CESFieldSource
    name: str
    value_type: PrimitiveValueType | Literal["json"]
    extension_type: NotRequired[_media_ces.CESExtensionType]
    extension_schema_version: NotRequired[int]


class _CESScorer(TypedDict):
    """A CES scorer declaration, named like `pb.EvalTableScorer`."""

    name: str
    value_type: PrimitiveValueType


class _CESRow(TypedDict):
    input: dict[str, Any] | None
    output: dict[str, Any] | None
    scores: dict[str, Any]
    example_key: NotRequired[str]


@dataclass(frozen=True)
class _CESFieldType:
    """Track a column's value schema and build its CES field declaration."""

    value_type: PrimitiveValueType | Literal["json"]
    extension_type: _media_ces.CESExtensionType | None = None
    extension_schema_version: int | None = None

    def declaration(self, source: _CESFieldSource, name: str) -> _CESDatasetField:
        """Build the CES client field declaration for this observed column type."""
        result: _CESDatasetField = {
            "source": source,
            "name": name,
            "value_type": self.value_type,
        }
        if self.extension_type is not None:
            result["extension_type"] = self.extension_type
        if self.extension_schema_version is not None:
            result["extension_schema_version"] = self.extension_schema_version
        return result


@dataclass(frozen=True)
class _BoundRun:
    run: LocalRun
    eval_table_key: str


@dataclass(frozen=True)
class _CESWritePayloads:
    dataset_fields: list[_CESDatasetField]
    scorers: list[_CESScorer]
    rows_json: list[str]
    media_cells_examined: int
    oversized_media_cells: int
    oversized_locations: tuple[str, ...]


# Bytes in a row-add body other than encoded rows and their separating commas.
_ROW_BATCH_ENVELOPE_BYTES = len(_encode_json({"rows": []}))


def _ces_row(
    *,
    row_index: int,
    inputs: dict[str, Any],
    output: dict[str, Any] | None,
    scores: dict[str, Any],
) -> _CESRow:
    """Build a row-add row, keying a row without inputs by position.

    CES derives the example key of a row with inputs from a hash of its whole
    input, so comparisons match rows with equal inputs.
    """
    if inputs:
        return {"input": inputs, "output": output, "scores": scores}
    return {
        "input": None,
        "output": output,
        "scores": scores,
        "example_key": str(row_index + 1),
    }


def _encode_rows(rows: Sequence[_CESRow]) -> list[str]:
    """Encode rows as JSON, rejecting any too large to upload alone."""
    rows_json: list[str] = []

    for row_index, row in enumerate(rows):
        row_json = _encode_json(row)
        if _ROW_BATCH_ENVELOPE_BYTES + len(row_json) >= _MAX_REQUEST_BODY_BYTES:
            raise UsageError(
                "CES EvalTable rows payload contains a row at index "
                f"{row_index} whose encoded request must be smaller than "
                f"{_MAX_REQUEST_BODY_SIZE}."
            )
        rows_json.append(row_json.decode())

    return rows_json


def _write_rows_file(rows_json: Sequence[str]) -> str:
    """Write rows to a gzipped temporary file, one per line.

    wandb-core deletes the file after reading it.

    Returns:
        The file's path.
    """
    file = tempfile.NamedTemporaryFile(
        prefix="wandb-eval-table-",
        suffix=".rows.jl.gz",
        delete=False,
    )
    try:
        with (
            file,
            gzip.open(file, mode="wt", encoding="utf-8", compresslevel=1) as rows,
        ):
            for row_json in rows_json:
                # Encoded JSON has no raw newlines.
                rows.write(row_json + "\n")
    except BaseException:
        # Delete the partial file, then propagate the original error.
        os.remove(file.name)
        raise

    return file.name


def _upload_eval_table(
    run: LocalRun,
    request: pb.UploadEvalTableRequest,
) -> pb.UploadEvalTableResponse:
    """Upload an EvalTable to CES through wandb-core and wait for the result.

    Raises:
        CommError: If the upload fails.
    """
    if not run._interface:
        raise UsageError("Cannot upload an EvalTable without a backend.")

    handle = run._interface.deliver_upload_eval_table(request)
    response = handle.wait_or(timeout=None).response.upload_eval_table_response
    if response.error_message:
        raise CommError(f"Failed to upload EvalTable: {response.error_message}")

    return response


class CESWriter:
    """Write an immutable EvalTable through the Evaluations service."""

    def __init__(
        self,
        *,
        unsupported_media_mode: str = "stub",
    ) -> None:
        self._unsupported_media_mode = unsupported_media_mode
        self._bound: _BoundRun | None = None

    def validate_cell_value(self, value: Any, column: ColumnKey) -> None:
        """Validate W&B values before the table is bound or written."""
        if _media_ces.is_supported_wandb_media(value):
            return
        if isinstance(value, WBValue):
            self._validate_wandb_value(value, column)

    def bind_to_run(self, run: LocalRun, key: str) -> None:
        """Bind the run and the history key the table is logged at."""
        if not run.entity or not run.project:
            raise UsageError("CES EvalTable logging requires a W&B entity and project.")
        self._bound = _BoundRun(run=run, eval_table_key=key)

    def write(
        self,
        *,
        name: str,
        rows: Sequence[WriteRow],
        ncols: int,
        log_mode: LogMode,
    ) -> WriteResult:
        """Prepare and persist the CES resources, then return their history value."""
        bound_run = self._require_bound()

        write_payloads = self._build_write_payloads(name=name, rows=rows)
        self._record_media_telemetry(write_payloads)
        if write_payloads.oversized_media_cells:
            locations = ", ".join(write_payloads.oversized_locations)
            wandb.termwarn(
                "EvalTable replaced "
                f"{write_payloads.oversized_media_cells} of "
                f"{write_payloads.media_cells_examined} media cells with null because "
                "their encoded payload reached the "
                f"{_media_ces.CES_MAX_CELL_BYTES / 1_000_000:g} MB limit. "
                f"First affected: {locations}.",
            )

        request = pb.UploadEvalTableRequest(
            name=name,
            dataset_fields=[
                pb.EvalTableField(**dataset_field)
                for dataset_field in write_payloads.dataset_fields
            ],
            scorers=[pb.EvalTableScorer(**scorer) for scorer in write_payloads.scorers],
        )
        rows_chars = sum(len(row_json) for row_json in write_payloads.rows_json)
        if rows_chars > _ROWS_FILE_THRESHOLD_CHARS:
            request.rows_file_path = _write_rows_file(write_payloads.rows_json)
        else:
            request.rows_json.extend(write_payloads.rows_json)

        response = _upload_eval_table(bound_run.run, request)

        _logger.debug(
            "CES EvalTable recorded evaluation_version_id=%s",
            response.evaluation_version_id,
        )

        return WriteResult(
            history_value={
                # Frontend dispatches on `_type` and validates the schema version.
                "_type": "eval-table-ces",
                "schema_version": 1,
                "ncols": ncols,
                "nrows": len(rows),
                "log_mode": log_mode,
                "evaluation_id": response.evaluation_id,
                "evaluation_version_id": response.evaluation_version_id,
                "dataset_id": response.dataset_id,
                "dataset_version_id": response.dataset_version_id,
            },
            logged_id=response.evaluation_version_id,
        )

    def _require_bound(self) -> _BoundRun:
        if self._bound is None:
            raise UsageError("EvalTable must be logged with run.log().")
        return self._bound

    def _build_write_payloads(
        self,
        *,
        name: str,
        rows: Sequence[WriteRow],
    ) -> _CESWritePayloads:
        """Validate and normalize rows, infer schemas, and encode the rows."""
        self._validate_name(
            "EvalTable",
            name,
            max_length=_MAX_EVAL_TABLE_NAME_LENGTH,
        )
        if not rows:
            raise UsageError("CES EvalTable logging requires at least one row.")
        if len(rows) > _MAX_ROWS_PER_TABLE:
            raise UsageError(
                "CES EvalTable logging currently supports at most "
                f"{_MAX_ROWS_PER_TABLE} rows per table."
            )

        bound_run = self._require_bound()
        dataset_field_types: dict[tuple[_CESFieldSource, str], _CESFieldType] = {}
        dataset_field_order: dict[tuple[_CESFieldSource, str], None] = {}
        scorer_types: dict[str, PrimitiveValueType] = {}
        scorer_order: dict[str, None] = {}
        normalized_rows: list[_CESRow] = []
        media_cells_examined = 0
        oversized_media_cells = 0
        oversized_locations: list[str] = []
        class_label_accumulator = _media_ces.ClassLabelAccumulator()

        for row_index, row in enumerate(rows):
            inputs, input_media, input_oversized = (
                self._normalize_row_input_or_output_values(
                    row.inputs,
                    source="input",
                    types=dataset_field_types,
                    order=dataset_field_order,
                    bound_run=bound_run,
                    row_index=row_index,
                    oversized_locations=oversized_locations,
                    class_label_accumulator=class_label_accumulator,
                )
            )
            media_cells_examined += input_media
            oversized_media_cells += input_oversized
            output = None
            if row.output is not None:
                output, output_media, output_oversized = (
                    self._normalize_row_input_or_output_values(
                        row.output,
                        source="output",
                        types=dataset_field_types,
                        order=dataset_field_order,
                        bound_run=bound_run,
                        row_index=row_index,
                        oversized_locations=oversized_locations,
                        class_label_accumulator=class_label_accumulator,
                    )
                )
                media_cells_examined += output_media
                oversized_media_cells += output_oversized
            scores = self._normalize_row_score_values(
                row.scores,
                types=scorer_types,
                order=scorer_order,
            )
            normalized_rows.append(
                _ces_row(
                    row_index=row_index,
                    inputs=inputs,
                    output=output,
                    scores=scores,
                )
            )

        missing_fields = [
            column_name
            for source, column_name in dataset_field_order
            if (source, column_name) not in dataset_field_types
        ]
        missing_scorers = [
            column_name
            for column_name in scorer_order
            if column_name not in scorer_types
        ]
        if missing_fields or missing_scorers:
            missing = sorted({*missing_fields, *missing_scorers})
            raise UsageError(
                "Cannot infer primitive types for all-null EvalTable column(s): "
                f"{missing}."
            )

        if len(dataset_field_order) > _MAX_DATASET_FIELDS:
            raise UsageError(
                "CES EvalTable logging supports at most "
                f"{_MAX_DATASET_FIELDS} input and output columns combined."
            )
        if len(scorer_order) > _MAX_SCORERS:
            raise UsageError(
                f"CES EvalTable logging supports at most {_MAX_SCORERS} score columns."
            )

        dataset_fields: list[_CESDatasetField] = [
            dataset_field_types[(source, column_name)].declaration(source, column_name)
            for source, column_name in dataset_field_order
        ]
        scorers: list[_CESScorer] = [
            {"name": column_name, "value_type": scorer_types[column_name]}
            for column_name in scorer_order
        ]
        self._validate_body_size(
            "columns", {"dataset_fields": dataset_fields, "scorers": scorers}
        )
        # Encode rows first so every size error precedes config updates and requests.
        rows_json = _encode_rows(normalized_rows)
        class_label_accumulator.flush(bound_run.run)
        return _CESWritePayloads(
            dataset_fields=dataset_fields,
            scorers=scorers,
            rows_json=rows_json,
            media_cells_examined=media_cells_examined,
            oversized_media_cells=oversized_media_cells,
            oversized_locations=tuple(oversized_locations),
        )

    def _normalize_row_input_or_output_values(
        self,
        values: Mapping[ColumnKey, Any],
        *,
        source: _CESFieldSource,
        types: dict[tuple[_CESFieldSource, str], _CESFieldType],
        order: dict[tuple[_CESFieldSource, str], None],
        bound_run: _BoundRun,
        row_index: int,
        oversized_locations: list[str],
        class_label_accumulator: _media_ces.ClassLabelAccumulator,
    ) -> tuple[dict[str, Any], int, int]:
        """Normalize one input/output mapping and accumulate its inferred types."""
        normalized_values: dict[str, Any] = {}
        media_cells = 0
        oversized_cells = 0
        for column, value in values.items():
            name = str(column)
            key = (source, name)
            if key not in order:
                self._validate_name(
                    f"{source} column",
                    name,
                    max_length=_MAX_DATASET_FIELD_NAME_LENGTH,
                )
                order[key] = None
            if _media_ces.is_supported_wandb_media(value):
                media_cells += 1
                normalized, value_type, oversized_size = self._normalize_media(
                    value,
                    bound_run,
                    source=source,
                    column_name=name,
                    class_label_accumulator=class_label_accumulator,
                )
                if oversized_size is not None:
                    oversized_cells += 1
                    if len(oversized_locations) < _MAX_OVERSIZED_MEDIA_LOCATIONS:
                        oversized_locations.append(
                            f"row {row_index}, {source} column {name!r} "
                            f"({oversized_size} bytes)"
                        )
            else:
                normalized, primitive_type = self._normalize_primitive(value, column)
                value_type = (
                    _CESFieldType(primitive_type)
                    if primitive_type is not None
                    else None
                )
            if value_type is not None:
                types[key] = self._merge_field_type(column, types.get(key), value_type)
            normalized_values[name] = normalized
        return normalized_values, media_cells, oversized_cells

    def _normalize_row_score_values(
        self,
        values: Mapping[ColumnKey, Any],
        *,
        types: dict[str, PrimitiveValueType],
        order: dict[str, None],
    ) -> dict[str, Any]:
        """Normalize one score mapping and accumulate its inferred types."""
        normalized_values: dict[str, Any] = {}
        for column, value in values.items():
            name = str(column)
            if name not in order:
                self._validate_name(
                    "score column",
                    name,
                    max_length=_MAX_SCORER_NAME_LENGTH,
                )
                order[name] = None
            if _media_ces.is_supported_wandb_media(value):
                raise UsageError(
                    f"EvalTable score column {column!r} contains unsupported value "
                    f"type {type(value).__name__!r}; only primitive values are supported."
                )
            normalized, value_type = self._normalize_primitive(value, column)
            if value_type is not None:
                types[name] = self._merge_type(column, types.get(name), value_type)
            normalized_values[name] = normalized
        return normalized_values

    def _normalize_media(
        self,
        value: Media,
        bound_run: _BoundRun,
        *,
        source: _CESFieldSource,
        column_name: str,
        class_label_accumulator: _media_ces.ClassLabelAccumulator,
    ) -> tuple[Any, _CESFieldType, int | None]:
        """Return a CES extension value, its field type, and oversized byte count."""
        try:
            media_cell = _media_ces.prepare_media(
                value,
                bound_run.run,
                _media_ces.EvalTableMediaField(
                    eval_table_key=bound_run.eval_table_key,
                    source="inputs" if source == "input" else "outputs",
                    column_name=column_name,
                ),
                class_label_accumulator=class_label_accumulator,
            )
        except _media_ces.UnsupportedMediaVariantError as error:
            if self._unsupported_media_mode == "raise":
                raise
            wandb.termwarn(error.stub_warning, repeat=False)
            return (
                None,
                _CESFieldType(
                    value_type="json",
                    extension_type=_media_ces.media_extension_type(value),
                    extension_schema_version=1,
                ),
                None,
            )
        field_type = _CESFieldType(
            value_type="json",
            extension_type=media_cell.extension_type,
            extension_schema_version=media_cell.extension_schema_version,
        )
        oversized_size = media_cell.encoded_size if media_cell.oversized else None
        return media_cell.value, field_type, oversized_size

    def _normalize_primitive(
        self,
        value: Any,
        column: ColumnKey,
    ) -> tuple[Any, PrimitiveValueType | None]:
        """Return a JSON-safe value and its observed CES primitive type."""
        if isinstance(value, WBValue):
            value = self._normalize_wandb_value(value, column)
        if value is None:
            return None, None

        np = wandb.util.np
        if isinstance(value, bool) or (np is not None and isinstance(value, np.bool_)):
            return bool(value), "boolean"
        if isinstance(value, int) or (np is not None and isinstance(value, np.integer)):
            return int(value), "integer"
        if isinstance(value, float) or (
            np is not None and isinstance(value, np.floating)
        ):
            normalized = float(value)
            if math.isnan(normalized):
                # JSON has no NaN literal; send null while retaining numeric schema.
                return None, "number"
            if not math.isfinite(normalized):
                raise UsageError(
                    f"CES EvalTable column {column!r} contains a non-finite number."
                )
            return normalized, "number"
        if isinstance(value, str) or (np is not None and isinstance(value, np.str_)):
            return str(value), "string"

        raise UsageError(
            f"CES EvalTable column {column!r} contains unsupported value "
            f"type {type(value).__name__!r}; only primitive values are supported."
        )

    def _validate_name(self, kind: str, name: str, *, max_length: int) -> None:
        if not name or len(name) > max_length:
            raise UsageError(
                f"CES {kind} names must contain between 1 and {max_length} "
                f"characters; got {name!r}."
            )

    def _validate_wandb_value(self, value: WBValue, column: ColumnKey) -> None:
        if isinstance(value, Table):
            raise TypeError(
                f"Column {column!r} contains a {type(value).__name__}; "
                "CES EvalTable logging does not support nested Tables."
            )
        if self._unsupported_media_mode == "raise":
            value_kind = "media" if isinstance(value, Media) else "value"
            raise TypeError(
                f"Column {column!r} contains unsupported wandb {value_kind} type "
                f"{type(value).__name__!r}. CES EvalTable logging does not support "
                "wandb media/value types yet. Pass unsupported_media_mode='stub' "
                "to log a placeholder string instead."
            )

    def _normalize_wandb_value(self, value: WBValue, column: ColumnKey) -> str:
        """Validate or replace an unsupported W&B value with a placeholder."""
        # Keep media adaptation behind one hook so CES-native media can replace
        # the unsupported fallback as its schemas and upload paths are added.
        self._validate_wandb_value(value, column)
        wandb.termwarn(
            f"wandb.{type(value).__name__} values are not yet supported by CES "
            "EvalTable logging. They will be logged as placeholder strings.",
            repeat=False,
        )
        return f"[wandb.{type(value).__name__} not yet supported]"

    def _merge_type(
        self,
        column: ColumnKey,
        existing: PrimitiveValueType | None,
        observed: PrimitiveValueType,
    ) -> PrimitiveValueType:
        """Merge observed types, widening integer and number columns to number."""
        if existing is None or existing == observed:
            return observed
        if {existing, observed} == {"integer", "number"}:
            return "number"
        # An explicit permissive dtype can bypass Table's usual type check.
        raise UsageError(
            f"CES EvalTable column {column!r} mixes {existing!r} and "
            f"{observed!r} values."
        )

    def _merge_field_type(
        self,
        column: ColumnKey,
        existing: _CESFieldType | None,
        observed: _CESFieldType,
    ) -> _CESFieldType:
        if existing is None or existing == observed:
            return observed
        if (
            existing.extension_type is None
            and observed.extension_type is None
            and existing.value_type != "json"
            and observed.value_type != "json"
        ):
            return _CESFieldType(
                self._merge_type(
                    column,
                    existing.value_type,
                    observed.value_type,
                )
            )
        existing_type = existing.extension_type or existing.value_type
        observed_type = observed.extension_type or observed.value_type
        raise UsageError(
            f"EvalTable column {column!r} mixes {existing_type!r} and "
            f"{observed_type!r} values."
        )

    def _record_media_telemetry(self, prepared: _CESWritePayloads) -> None:
        # Track whether cells above 3.5 MB justify extended media metadata support.
        if prepared.media_cells_examined == 0:
            return

        recorder = get_telemetry_recorder()
        recorder.increment_counter(
            "eval_table_ces_media_write",
            LowCardinalityAttributes(),
        )
        if prepared.oversized_media_cells:
            recorder.increment_counter(
                "eval_table_ces_media_write_with_oversized_cells",
                LowCardinalityAttributes(),
            )

    def _validate_body_size(self, operation: str, body: dict[str, Any]) -> None:
        """Reject a request body at or above the CES hard limit."""
        if len(_encode_json(body)) >= _MAX_REQUEST_BODY_BYTES:
            raise UsageError(
                f"CES EvalTable {operation} payload must be smaller than "
                f"{_MAX_REQUEST_BODY_SIZE}."
            )
