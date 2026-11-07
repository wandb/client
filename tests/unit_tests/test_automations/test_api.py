"""Parked WB-42558 scope-recovery regression tests."""

from __future__ import annotations

from copy import deepcopy

from pytest import fixture, mark
from wandb.apis.public import Api
from wandb.automations import DoNothing, OnLinkArtifact, ScopeType
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
