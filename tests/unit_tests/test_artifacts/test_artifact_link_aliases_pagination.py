"""Tests for lazy membership aliases returned by Artifact.link()."""

from unittest.mock import Mock

from wandb.proto import wandb_internal_pb2 as pb
from wandb.sdk.artifacts._generated import (
    ArtifactFragment,
    ArtifactMembershipAliases,
    LinkArtifact,
)
from wandb.sdk.artifacts._validators import FullArtifactPath
from wandb.sdk.artifacts.artifact import Artifact


def _artifact_fragment():
    return ArtifactFragment(
        name="source-collection",
        version_index=0,
        artifact_type={"name": "dataset"},
        artifact_sequence={
            "name": "source-collection",
            "project": {
                "id": "UHJvamVjdDox",
                "internalId": "UHJvamVjdEludGVybmFsSWQ6MQ==",
                "name": "source-project",
                "entity": {"name": "source-entity"},
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
            "name": "target-collection",
            "project": {
                "id": "UHJvamVjdDox",
                "internalId": "UHJvamVjdEludGVybmFsSWQ6MQ==",
                "name": "target-project",
                "entity": {"name": "target-entity"},
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
        pb.ARTIFACT_MEMBERSHIP_IN_LINK_ARTIFACT_RESPONSE: membership_enabled,
        pb.ARTIFACT_MEMBERSHIP_ALIASES_PAGINATION: pagination_enabled,
    }
    service_api.feature_enabled.side_effect = lambda feature, **_: enabled.get(
        feature, False
    )
    result = {"versionIndex": 3}
    if membership_enabled:
        result["artifactMembership"] = _membership(aliases=aliases)
    service_api.execute_graphql.return_value = LinkArtifact.model_validate(
        {"result": result}
    )
    path = FullArtifactPath(
        prefix="source-entity", project="source-project", name="source-collection:v0"
    )
    artifact = Artifact._from_attrs(path, _artifact_fragment(), service_api)
    return artifact, service_api


def test_link_uses_flat_aliases_when_pagination_disabled():
    artifact, service_api = _setup(
        pagination_enabled=False, membership_enabled=True, aliases=["best"]
    )

    linked = artifact.link("target-entity/target-project/target-collection")

    assert linked.aliases == ["best"]
    service_api.execute_graphql.assert_called_once()
    kwargs = service_api.execute_graphql.call_args.kwargs
    assert "aliases" not in kwargs["omit_fields"]
    assert kwargs["omit_variables"] is None


def test_link_defers_all_aliases_until_first_read():
    artifact, service_api = _setup(pagination_enabled=True, membership_enabled=True)
    link_response = service_api.execute_graphql.return_value
    service_api.execute_graphql.side_effect = [
        link_response,
        _alias_page(["v3", "best"]),
    ]

    linked = artifact.link("target-entity/target-project/target-collection")

    assert service_api.execute_graphql.call_count == 1
    assert linked._aliases is None
    assert "aliases" in service_api.execute_graphql.call_args.kwargs["omit_fields"]
    assert linked.aliases == ["best"]
    assert service_api.execute_graphql.call_count == 2
    alias_variables = service_api.execute_graphql.call_args.kwargs["variables"]
    assert alias_variables["entity"] == "target-entity"
    assert alias_variables["project"] == "target-project"
    assert alias_variables["name"] == "target-collection:v3"


def test_link_without_membership_feature_uses_fallback(mocker):
    artifact, service_api = _setup(pagination_enabled=True, membership_enabled=False)
    fetched = mocker.patch("wandb.Api")

    result = artifact.link("target-entity/target-project/target-collection")

    assert result is fetched.return_value._artifact.return_value
    assert (
        "artifactMembership"
        in service_api.execute_graphql.call_args.kwargs["omit_fields"]
    )
