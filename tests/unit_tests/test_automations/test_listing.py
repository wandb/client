"""Tests for listing automations when some can't be parsed by this SDK version."""

from __future__ import annotations

import json
from typing import Any
from unittest.mock import MagicMock

from pytest import fixture, raises
from wandb.apis.public import automations as automations_module
from wandb.apis.public.automations import EntityAutomations, LegacyEntityAutomations


def _trigger(
    name: str,
    *,
    event_type: str = "CREATE_ARTIFACT",
    action: dict[str, Any] | None = None,
) -> dict[str, Any]:
    return {
        "__typename": "Trigger",
        "id": f"id-{name}",
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": None,
        "name": name,
        "description": None,
        "enabled": True,
        "scope": {
            "__typename": "Project",
            "id": "UHJvamVjdDox",
            "name": "my-project",
        },
        "event": {
            "__typename": "FilterEventTriggeringCondition",
            "eventType": event_type,
            "filter": json.dumps({"filter": {"$or": [{"$and": []}]}}),
        },
        "action": action
        if action is not None
        else {"__typename": "NoOpTriggeredAction", "noOp": True},
    }


def _connection(nodes: list[Any]) -> dict[str, Any]:
    return {
        "pageInfo": {"endCursor": None, "hasNextPage": False},
        "edges": [{"node": node} for node in nodes],
    }


# A page mixing automations this SDK can and can't parse.
MIXED_TRIGGERS = [
    _trigger("supported"),
    # An action type this SDK has no fragment for: the server selects no fields.
    _trigger("unknown-action", action={}),
    # An event type this SDK doesn't know about.
    _trigger("unknown-event", event_type="SOME_FUTURE_EVENT"),
    _trigger("also-supported"),
]


@fixture
def termwarn(monkeypatch) -> MagicMock:
    mock = MagicMock()
    monkeypatch.setattr(automations_module.wandb, "termwarn", mock)
    return mock


def _service_api(response: dict[str, Any]) -> MagicMock:
    service_api = MagicMock()
    service_api.execute_graphql.return_value = response
    return service_api


def test_entity_automations_skips_unparseable(termwarn):
    response = {"scope": {"triggers": _connection(MIXED_TRIGGERS)}}

    automations = list(EntityAutomations(_service_api(response), entity="my-entity"))

    assert [a.name for a in automations] == ["supported", "also-supported"]
    termwarn.assert_called()


def test_entity_automations_no_warning_when_all_parse(termwarn):
    triggers = [_trigger("a"), _trigger("b")]
    response = {"scope": {"triggers": _connection(triggers)}}

    automations = list(EntityAutomations(_service_api(response), entity="my-entity"))

    assert [a.name for a in automations] == ["a", "b"]
    termwarn.assert_not_called()


def test_legacy_entity_automations_skips_unparseable(termwarn):
    project = {"__typename": "Project", "triggers": MIXED_TRIGGERS}
    response = {"scope": {"projects": _connection([project])}}

    automations = list(
        LegacyEntityAutomations(_service_api(response), entity="my-entity")
    )

    assert [a.name for a in automations] == ["supported", "also-supported"]
    termwarn.assert_called()


def test_legacy_entity_automations_filters_by_name_before_parsing(termwarn):
    project = {"__typename": "Project", "triggers": MIXED_TRIGGERS}
    response = {"scope": {"projects": _connection([project])}}

    automations = list(
        LegacyEntityAutomations(
            _service_api(response), entity="my-entity", name="also-supported"
        )
    )

    assert [a.name for a in automations] == ["also-supported"]
    # The unparseable automations were filtered out by name, so no warning.
    termwarn.assert_not_called()


def test_entity_automations_raises_on_malformed_response():
    with raises(ValueError, match="Unexpected response data"):
        list(EntityAutomations(_service_api({"scope": None}), entity="my-entity"))


def test_legacy_entity_automations_raises_on_malformed_response():
    with raises(ValueError, match="Unexpected response data"):
        list(LegacyEntityAutomations(_service_api({"scope": None}), entity="x"))
