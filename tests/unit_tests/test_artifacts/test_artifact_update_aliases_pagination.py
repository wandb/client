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
    artifact._aliases = None
    artifact._saved_aliases = None
    response = service_api.execute_graphql.return_value
    service_api.execute_graphql.side_effect = [response, _alias_page(["v3", "best"])]

    artifact._update()

    assert service_api.execute_graphql.call_count == 1
    assert artifact._aliases is None
    assert artifact.version == "v3"
    assert "aliases" in service_api.execute_graphql.call_args.kwargs["omit_fields"]
    assert service_api.execute_graphql.call_args.kwargs["omit_variables"] == {
        "includeAliases"
    }
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
    assert service_api.execute_graphql.call_args.kwargs["omit_variables"] is None


def test_update_skips_alias_diff_when_already_pending_and_untouched():
    """`_update()` must not pay for an alias fetch just to confirm nothing changed.

    Once a prior fetch has left `.aliases` pending (unread), a later
    `_update()` that never touches aliases either should skip resolving
    them entirely -- the diff against `_saved_aliases` is guaranteed empty
    without paying for a single page, let alone every remaining one.
    """
    artifact, service_api = _setup(pagination_enabled=True, membership_enabled=True)
    artifact._aliases = None
    artifact._saved_aliases = None

    artifact._update()

    # Only the `UpdateArtifact` mutation ran -- no alias-pagination
    # request, and no AddAliases/DeleteAliases mutation either.
    service_api.execute_graphql.assert_called_once()


def test_update_keeps_already_loaded_aliases_without_refetching():
    artifact, service_api = _setup(pagination_enabled=True, membership_enabled=True)
    artifact._version = "v3"
    artifact._name = "test-collection:v3"
    artifact._aliases = ["best"]
    artifact._saved_aliases = ["best"]

    artifact._update()
    artifact._update()

    assert artifact.aliases == ["best"]
    assert artifact._saved_aliases == ["best"]
    assert service_api.execute_graphql.call_count == 2
    assert all(
        "mutation UpdateArtifact" in call.args[0]
        for call in service_api.execute_graphql.call_args_list
    )


def test_linked_save_preserves_portfolio_membership_and_diffs_its_aliases(mocker):
    """The update response's source membership must not replace link state."""
    service_api = Mock()
    service_api.feature_enabled.side_effect = lambda feature, **_: (
        feature
        in {
            pb.ARTIFACT_MEMBERSHIP_ON_ARTIFACT_PAYLOAD,
            pb.ARTIFACT_MEMBERSHIP_ALIASES_PAGINATION,
        }
    )
    source = _artifact_fragment()
    source_membership = _membership(aliases=None)
    source_membership["versionIndex"] = 0
    source_membership["createdAt"] = "source-created-at"
    update_response = UpdateArtifact.model_validate(
        {
            "result": {
                "artifact": source.model_dump(),
                "artifactMembership": source_membership,
            }
        }
    )
    calls = []

    def execute_graphql(query, **kwargs):
        calls.append((query, kwargs))
        if "mutation UpdateArtifact" in query:
            return update_response
        if "query ArtifactMembershipAliases" in query:
            return _alias_page(["v7", "portfolio-only"])
        if "mutation AddAliases" in query or "mutation DeleteAliases" in query:
            return None
        raise AssertionError(f"Unexpected GraphQL operation: {query}")

    service_api.execute_graphql.side_effect = execute_graphql
    path = FullArtifactPath(
        prefix="portfolio-entity",
        project="portfolio-project",
        name="portfolio-collection:v7",
    )
    portfolio_membership = _membership(aliases=None)
    portfolio_membership["versionIndex"] = 7
    portfolio_membership["createdAt"] = "portfolio-created-at"
    portfolio_membership["artifactCollection"]["name"] = "portfolio-collection"
    portfolio_membership["artifactCollection"]["project"]["name"] = "portfolio-project"
    portfolio_membership["artifactCollection"]["project"]["entity"]["name"] = (
        "portfolio-entity"
    )
    artifact = Artifact._from_attrs(
        path, source, service_api, membership=portfolio_membership
    )
    mocker.patch.object(artifact, "_is_source_project_read_only", return_value=False)
    assert artifact.is_link is True
    assert artifact.version == "v7"

    artifact._update()

    assert artifact.version == "v7"
    assert artifact.linked_at == "portfolio-created-at"
    assert artifact._aliases is None
    assert artifact.aliases == ["portfolio-only"]
    alias_call = next(
        (query, kwargs)
        for query, kwargs in calls
        if "query ArtifactMembershipAliases" in query
    )
    assert alias_call[1]["variables"]["entity"] == "portfolio-entity"
    assert alias_call[1]["variables"]["project"] == "portfolio-project"
    assert alias_call[1]["variables"]["name"] == "portfolio-collection:v7"

    artifact.aliases = ["prod"]
    artifact._update()

    assert artifact.version == "v7"
    assert artifact.aliases == ["prod"]
    assert sum("query ArtifactMembershipAliases" in query for query, _ in calls) == 1
    add_call = next(
        (query, kwargs) for query, kwargs in calls if "mutation AddAliases" in query
    )
    delete_call = next(
        (query, kwargs) for query, kwargs in calls if "mutation DeleteAliases" in query
    )
    for _, kwargs in (add_call, delete_call):
        alias = kwargs["variables"]["input"]["aliases"][0]
        assert alias["entityName"] == "portfolio-entity"
        assert alias["projectName"] == "portfolio-project"
        assert alias["artifactCollectionName"] == "portfolio-collection"
    assert add_call[1]["variables"]["input"]["aliases"][0]["alias"] == "prod"
    assert (
        delete_call[1]["variables"]["input"]["aliases"][0]["alias"] == "portfolio-only"
    )

    artifact._update()
    assert sum("query ArtifactMembershipAliases" in query for query, _ in calls) == 1
    assert sum("mutation AddAliases" in query for query, _ in calls) == 1
    assert sum("mutation DeleteAliases" in query for query, _ in calls) == 1
