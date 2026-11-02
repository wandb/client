"""Tests for lazy membership aliases returned by Artifact._update()."""

from unittest.mock import Mock

from wandb.proto import wandb_internal_pb2 as pb
from wandb.sdk.artifacts._generated import (
    ArtifactFragment,
    ArtifactMembershipAliases,
    UpdateArtifact,
)
from wandb.sdk.artifacts._validators import FullArtifactPath
from wandb.sdk.artifacts.artifact import Artifact


def _artifact_fragment():
    return ArtifactFragment(
        name="test-collection",
        version_index=0,
        artifact_type={"name": "dataset"},
        artifact_sequence={
            "name": "test-collection",
            "project": {
                "id": "UHJvamVjdDox",
                "internalId": "UHJvamVjdEludGVybmFsSWQ6MQ==",
                "name": "test-project",
                "entity": {"name": "test-entity"},
            },
        },
        id="artifact-id",
        description=None,
        tags=[],
        ttl_duration_seconds=-2,
        ttl_is_inherited=False,
        metadata="{}",
        state="COMMITTED",
        size=0,
        digest="FAKE_DIGEST",
        file_count=0,
        commit_hash="PLACEHOLDER",
        created_at="PLACEHOLDER",
        updated_at=None,
        history_step=None,
    )


def _membership(*, aliases=None):
    membership = {
        "id": "membership-id",
        "versionIndex": 3,
        "createdAt": "PLACEHOLDER",
        "artifactCollection": {
            "__typename": "ArtifactPortfolio",
            "name": "test-collection",
            "project": {
                "id": "UHJvamVjdDox",
                "internalId": "UHJvamVjdEludGVybmFsSWQ6MQ==",
                "name": "test-project",
                "entity": {"name": "test-entity"},
            },
        },
        "artifact": _artifact_fragment().model_dump(),
    }
    if aliases is not None:
        membership["aliases"] = [
            {"id": f"alias-{alias}", "alias": alias} for alias in aliases
        ]
    return membership


def _alias_page(aliases):
    return ArtifactMembershipAliases.model_validate(
        {
            "project": {
                "artifactCollectionMembership": {
                    "aliasesConnection": {
                        "pageInfo": {"hasNextPage": False, "endCursor": None},
                        "edges": [
                            {"node": {"id": f"alias-{alias}", "alias": alias}}
                            for alias in aliases
                        ],
                    }
                }
            }
        }
    )


def _setup(*, pagination_enabled, membership_enabled, aliases=None):
    service_api = Mock()
    enabled = {
        pb.ARTIFACT_MEMBERSHIP_ON_ARTIFACT_PAYLOAD: membership_enabled,
        pb.ARTIFACT_MEMBERSHIP_ALIASES_PAGINATION: pagination_enabled,
    }
    service_api.feature_enabled.side_effect = lambda feature, **_: enabled.get(
        feature, False
    )
    result = {"artifact": _artifact_fragment().model_dump()}
    if membership_enabled:
        result["artifactMembership"] = _membership(aliases=aliases)
    service_api.execute_graphql.return_value = UpdateArtifact.model_validate(
        {"result": result}
    )
    path = FullArtifactPath(
        prefix="test-entity", project="test-project", name="test-collection:v0"
    )
    artifact = Artifact._from_attrs(path, _artifact_fragment(), service_api)
    return artifact, service_api


def test_update_uses_flat_aliases_when_pagination_disabled():
    artifact, service_api = _setup(
        pagination_enabled=False, membership_enabled=True, aliases=["best"]
    )

    artifact._update()

    assert artifact.aliases == ["best"]
    service_api.execute_graphql.assert_called_once()
    kwargs = service_api.execute_graphql.call_args.kwargs
    assert "aliases" not in kwargs["omit_fields"]
    assert kwargs["omit_variables"] is None


def test_update_defers_all_aliases_until_first_read():
    artifact, service_api = _setup(pagination_enabled=True, membership_enabled=True)
    response = service_api.execute_graphql.return_value
    service_api.execute_graphql.side_effect = [response, _alias_page(["v3", "best"])]

    artifact._update()

    assert service_api.execute_graphql.call_count == 1
    assert artifact._aliases is None
    assert artifact.version == "v3"
    assert "aliases" in service_api.execute_graphql.call_args.kwargs["omit_fields"]
    assert artifact.aliases == ["best"]
    assert service_api.execute_graphql.call_count == 2
    alias_vars = service_api.execute_graphql.call_args.kwargs["variables"]
    assert alias_vars["name"] == "test-collection:v3"


def test_update_without_membership_feature_uses_artifact_payload():
    artifact, service_api = _setup(pagination_enabled=True, membership_enabled=False)

    artifact._update()

    service_api.execute_graphql.assert_called_once()
    assert artifact.aliases == []
    assert (
        "artifactMembership"
        in service_api.execute_graphql.call_args.kwargs["omit_fields"]
    )
