"""Public API contracts: recovery must preserve the resolved rule and its scope."""

from __future__ import annotations

import json
from copy import deepcopy

from pytest import fixture, mark, raises
from wandb.apis.public import Api
from wandb.automations import (
    Automation,
    DoNothing,
    OnLinkArtifact,
    OnRunMetric,
    OnRunState,
    RunEvent,
    ScopeType,
)
from wandb.errors import CommError
from wandb.proto.wandb_api_pb2 import ApiErrorResponse
from wandb.sdk.lib.service.service_connection import WandbApiFailedError

from ._helpers import _connection, _service_api, _trigger


@fixture
def api():
    # These methods depend only on _service_api, supplied by each test. Avoid
    # constructing a real service or authenticating to a backend.
    return object.__new__(Api)


def _error(status):
    return WandbApiFailedError("request failed", ApiErrorResponse(http_status=status))


def _saved(event):
    node = _trigger("rule")
    node["scope"] = event.scope.model_dump()
    node["enabled"] = False
    node["description"] = "saved description"
    node["event"] = {
        "__typename": "FilterEventTriggeringCondition",
        "eventType": event.event_type,
        "filter": json.dumps({"filter": event.filter.model_dump()})
        if event.event_type == "LINK_MODEL"
        else event.filter.model_dump_json(),
    }
    node["action"] = {
        "__typename": "GenericWebhookTriggeredAction",
        "integration": {
            "__typename": "GenericWebhookIntegration",
            "id": "integration",
            "name": "test integration",
            "urlEndpoint": "https://example.com/webhook",
        },
        "requestPayload": json.dumps({"nested": {"value": 1}}),
    }
    return node


@mark.parametrize("scope_fixture", ["project", "artifact_collection", "team"])
def test_duplicate_create_fetches_only_the_requested_scope(api, request, scope_fixture):
    scope = request.getfixturevalue(scope_fixture)
    event = OnLinkArtifact(scope=scope)
    target = _trigger("rule")
    target["scope"] = event.scope.model_dump()
    sibling = deepcopy(target)
    sibling["id"] = "sibling-rule"
    sibling["scope"]["id"] = "another-scope"
    sibling["action"] = {}  # Another scope must not even be parsed.
    nodes = [sibling, target]
    if event.scope.scope_type == ScopeType.PROJECT:
        response = {"scope": {"edges": [{"node": {"triggers": nodes}}]}}
    elif event.scope.scope_type == ScopeType.ARTIFACT_COLLECTION:
        response = {"scope": {"project": {"triggers": nodes}}}
    else:
        response = {"scope": {"triggers": _connection(nodes)}}
    api._service_api = service = _service_api(_error(409), response)

    existing = api.create_automation(
        event >> DoNothing(), name="rule", fetch_existing=True
    )

    assert existing.id == target["id"]
    assert existing.scope.id == scope.id
    assert existing.name == "rule"
    calls = service.execute_graphql.call_args_list
    assert len(calls) == 2
    assert "viewer" not in calls[1].args[0]
    variables = calls[1].kwargs["variables"]
    if event.scope.scope_type == ScopeType.ENTITY:
        assert variables["entity"] == scope.name
    else:
        assert variables["scopeID"] == scope.id


@mark.parametrize("case", ["artifact", "metric", "state", "overrides"])
def test_missing_update_recreates_the_resolved_rule(api, project, case):
    events = {
        "artifact": OnLinkArtifact(scope=project),
        "metric": OnRunMetric(
            scope=project, filter=RunEvent.metric("loss").avg(5).gt(2)
        ),
        "state": OnRunState(scope=project, filter=RunEvent.state == "failed"),
    }
    node = _saved(events.get(case, events["artifact"]))
    saved = Automation.model_validate(node)
    overrides = {}
    expected = deepcopy(node)
    if case == "overrides":
        overrides = {
            "name": "replacement",
            "description": "new description",
            "enabled": True,
            "event": events["state"],
            "action": DoNothing(),
        }
        expected.update(
            name="replacement",
            description="new description",
            enabled=True,
            event={
                "__typename": "FilterEventTriggeringCondition",
                "eventType": "RUN_STATE",
                "filter": events["state"].filter.model_dump_json(),
            },
            action={"__typename": "NoOpTriggeredAction", "noOp": True},
        )
    expected["id"] = "new-rule-id"
    api._service_api = service = _service_api(
        _error(404), {"result": {"trigger": expected}}
    )

    recreated = api.update_automation(saved, create_missing=True, **overrides)

    assert recreated == Automation.model_validate(expected)
    calls = service.execute_graphql.call_args_list
    assert len(calls) == 2
    assert "mutation CreateAutomation" in calls[1].args[0]
    update_input = calls[0].kwargs["variables"]["input"]
    create_input = calls[1].kwargs["variables"]["input"]
    assert create_input == {k: v for k, v in update_input.items() if k != "id"}
    assert saved.id == node["id"]
    assert saved.enabled is False


@mark.parametrize(
    "operation,status,recover",
    [
        ("create", 409, False),
        ("update", 404, False),
        ("create", 403, True),
        ("update", 403, True),
    ],
    ids=["conflict", "not-found", "create-forbidden", "update-forbidden"],
)
def test_unrecoverable_errors_do_not_issue_another_request(
    api, project, operation, status, recover
):
    event = OnLinkArtifact(scope=project)
    api._service_api = service = _service_api(_error(status))
    with raises(CommError):
        if operation == "create":
            api.create_automation(
                event >> DoNothing(), name="rule", fetch_existing=recover
            )
        else:
            api.update_automation(
                Automation.model_validate(_saved(event)),
                create_missing=recover,
            )
    assert service.execute_graphql.call_count == 1
