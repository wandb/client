from unittest.mock import Mock

import pytest
from wandb.proto import wandb_internal_pb2 as pb
from wandb.sdk.artifacts._generated import (
    ArtifactFragment,
    ArtifactMembershipAliases,
    ArtifactMembershipByName,
)
from wandb.sdk.artifacts._validators import FullArtifactPath, LinkArtifactFields
from wandb.sdk.artifacts.artifact import Artifact

_PATH = FullArtifactPath(
    prefix="test-entity", project="test-project", name="test-collection:latest"
)


def _membership_data(*, flat_aliases=None, version_index=3):
    artifact = ArtifactFragment(
        name="test-collection",
        version_index=version_index,
        artifact_type={"name": "dataset"},
        artifact_sequence={
            "name": "test-collection",
            "project": {
                "id": "UHJvamVjdDox",
                "internalId": "UHJvamVjdEludGVybmFsSWQ6MQ==",
                "name": "source-project",
                "entity": {"name": "source-entity"},
            },
        },
        id="PLACEHOLDER",
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
    ).model_dump()
    membership = {
        "id": "PLACEHOLDER",
        "versionIndex": version_index,
        "createdAt": "PLACEHOLDER",
        "artifactCollection": {
            "__typename": "ArtifactPortfolio",
            "name": "test-collection",
            "project": {
                "id": "UHJvamVjdDox",
                "internalId": "UHJvamVjdEludGVybmFsSWQ6MQ==",
                "name": "canonical-project",
                "entity": {"name": "canonical-entity"},
            },
        },
        "artifact": artifact,
    }
    if flat_aliases is not None:
        membership["aliases"] = [
            {"id": "PLACEHOLDER", "alias": alias} for alias in flat_aliases
        ]
    return ArtifactMembershipByName.model_validate(
        {"project": {"artifactCollectionMembership": membership}}
    )


def _alias_page(aliases, *, has_next_page=False, end_cursor=None):
    return ArtifactMembershipAliases.model_validate(
        {
            "project": {
                "artifactCollectionMembership": {
                    "aliasesConnection": {
                        "pageInfo": {
                            "hasNextPage": has_next_page,
                            "endCursor": end_cursor,
                        },
                        "edges": [
                            {"node": {"id": "PLACEHOLDER", "alias": alias}}
                            for alias in aliases
                        ],
                    }
                }
            }
        }
    )


def _service_api(*responses, pagination_enabled=True):
    service_api = Mock()
    service_api.feature_enabled.side_effect = lambda feature, **_: (
        feature == pb.ARTIFACT_MEMBERSHIP_ALIASES_PAGINATION and pagination_enabled
    )
    pending = iter(responses)
    service_api.variable_snapshots = []

    def execute_graphql(*args, **kwargs):
        service_api.variable_snapshots.append(dict(kwargs.get("variables") or {}))
        result = next(pending)
        if isinstance(result, Exception):
            raise result
        return result

    service_api.execute_graphql.side_effect = execute_graphql
    return service_api


def test_from_name_uses_flat_aliases_on_old_server():
    service_api = _service_api(
        _membership_data(flat_aliases=["v3", "best"]), pagination_enabled=False
    )

    artifact = Artifact._from_name(path=_PATH, service_api=service_api)

    assert artifact.aliases == ["best"]
    service_api.execute_graphql.assert_called_once()
    assert "aliases" not in service_api.execute_graphql.call_args.kwargs["omit_fields"]


def test_from_name_defers_all_alias_pages_and_uses_canonical_identity():
    service_api = _service_api(
        _membership_data(),
        _alias_page(["v3", "first"], has_next_page=True, end_cursor="cursor-1"),
        _alias_page(["best"]),
    )

    artifact = Artifact._from_name(path=_PATH, service_api=service_api)

    assert service_api.execute_graphql.call_count == 1
    assert artifact._aliases is None
    assert artifact._saved_aliases is None
    assert "aliases" in service_api.execute_graphql.call_args.kwargs["omit_fields"]

    assert artifact.aliases == ["first", "best"]
    assert service_api.execute_graphql.call_count == 3
    first_alias_vars = service_api.variable_snapshots[1]
    assert first_alias_vars["name"] == "test-collection:v3"
    assert first_alias_vars["entity"] == "canonical-entity"
    assert first_alias_vars["project"] == "canonical-project"
    assert first_alias_vars["cursor"] is None
    assert first_alias_vars["perPage"] == 1000
    assert artifact._saved_aliases == ["first", "best"]
    assert artifact.aliases == ["first", "best"]
    assert service_api.execute_graphql.call_count == 3


def test_aliases_retry_starts_over_after_failed_later_page():
    service_api = _service_api(
        _membership_data(),
        _alias_page(["first"], has_next_page=True, end_cursor="cursor-1"),
        TimeoutError("transient network error"),
        _alias_page(["first"], has_next_page=True, end_cursor="cursor-1"),
        _alias_page(["best"]),
    )
    artifact = Artifact._from_name(path=_PATH, service_api=service_api)

    with pytest.raises(TimeoutError):
        _ = artifact.aliases

    assert artifact._aliases is None
    assert artifact.aliases == ["first", "best"]
    assert service_api.variable_snapshots[3]["cursor"] is None


def test_aliases_setter_fetches_before_replacing():
    service_api = _service_api(_membership_data(), _alias_page(["old"]))
    artifact = Artifact._from_name(path=_PATH, service_api=service_api)

    artifact.aliases = ["new"]

    assert artifact._saved_aliases == ["old"]
    assert artifact.aliases == ["new"]
    assert service_api.execute_graphql.call_count == 2


def test_copy_as_link_does_not_inherit_source_pending_aliases():
    service_api = _service_api(_membership_data(), _alias_page(["linked"]))
    artifact = Artifact._from_name(path=_PATH, service_api=service_api)
    linked = artifact._create_linked_artifact_using_source_artifact(
        LinkArtifactFields(
            entity_name="linked-entity",
            project_name="linked-project",
            name="linked-collection:v7",
            version="v7",
            aliases=["known"],
            linked_at="PLACEHOLDER",
        )
    )

    assert linked.aliases == ["known"]
    assert service_api.execute_graphql.call_count == 1
    assert artifact._aliases is None


def test_update_skips_alias_fetch_when_aliases_untouched(mocker):
    service_api = _service_api(
        _membership_data(),
        Mock(result=Mock(artifact=Mock(), artifact_membership=None)),
    )
    artifact = Artifact._from_name(path=_PATH, service_api=service_api)
    add_aliases = mocker.patch.object(artifact, "_add_aliases")
    delete_aliases = mocker.patch.object(artifact, "_delete_aliases")
    mocker.patch.object(artifact, "_assign_attrs")

    artifact._update()

    assert service_api.execute_graphql.call_count == 2
    add_aliases.assert_not_called()
    delete_aliases.assert_not_called()
    assert artifact._aliases is None


def test_update_diffs_explicit_alias_edit_against_fetched_aliases(mocker):
    service_api = _service_api(
        _membership_data(),
        _alias_page(["old"]),
        Mock(result=Mock(artifact=Mock(), artifact_membership=None)),
    )
    artifact = Artifact._from_name(path=_PATH, service_api=service_api)
    artifact.aliases = ["new"]
    add_aliases = mocker.patch.object(artifact, "_add_aliases")
    delete_aliases = mocker.patch.object(artifact, "_delete_aliases")
    mocker.patch.object(artifact, "_assign_attrs")

    artifact._update()

    assert add_aliases.call_args.args[0] == {"new"}
    assert delete_aliases.call_args.args[0] == {"old"}
    assert service_api.execute_graphql.call_count == 3
