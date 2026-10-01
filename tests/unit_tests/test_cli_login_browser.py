"""`wandb login --browser`, and how it interacts with API key logins."""

from __future__ import annotations

import datetime
import pathlib

import pytest
from wandb.cli import cli
from wandb.errors import AuthenticationError
from wandb.sdk import wandb_setup
from wandb.sdk.lib import wbauth
from wandb.sdk.lib.wbauth import browser_login

HOST = "https://test-host"


def make_token_set(refresh_token: str = "wb_rt_stored") -> browser_login.TokenSet:
    return browser_login.TokenSet(
        access_token="wb_at_stored",
        refresh_token=refresh_token,
        expires_at=datetime.datetime(2030, 1, 1, tzinfo=datetime.timezone.utc),
    )


@pytest.fixture(autouse=True)
def browser_is_usable(monkeypatch: pytest.MonkeyPatch) -> None:
    """Pin the browser check so it does not depend on this machine."""
    monkeypatch.setattr(browser_login, "can_open_browser", lambda: True)


@pytest.fixture
def credentials_file(
    tmp_path: pathlib.Path,
    monkeypatch: pytest.MonkeyPatch,
) -> pathlib.Path:
    path = tmp_path / "credentials.json"
    monkeypatch.setenv("WANDB_CREDENTIALS_FILE", str(path))
    wandb_setup.singleton().settings.credentials_file = str(path)
    return path


@pytest.fixture
def no_verify(monkeypatch: pytest.MonkeyPatch) -> None:
    """Skip the server round trip that confirms the new credentials."""
    monkeypatch.setattr(wbauth.AuthBrowserLogin, "verify", lambda self: None)


@pytest.fixture
def fake_login(monkeypatch: pytest.MonkeyPatch) -> dict:
    """Stand in for the browser half of login. Empty means it was not used."""
    calls: dict = {}

    # The real parameters, so a call with an argument it lacks fails.
    def fake(
        *,
        host: wbauth.HostUrl,
        timeout: float,
        verify: bool | str,
    ) -> browser_login.TokenSet:
        calls.update(host=host, timeout=timeout, verify=verify)
        return make_token_set("wb_rt_fresh")

    monkeypatch.setattr(browser_login, "login", fake)
    return calls


@pytest.fixture
def api_key_logins(monkeypatch: pytest.MonkeyPatch) -> list[dict]:
    """Record wandb.login() calls, which answer that the login succeeded."""
    calls: list[dict] = []

    def fake(**kwargs: object) -> bool:
        calls.append(kwargs)
        return True

    monkeypatch.setattr("wandb.login", fake)
    return calls


@pytest.fixture
def revocations(monkeypatch: pytest.MonkeyPatch) -> list[str]:
    revoked: list[str] = []
    monkeypatch.setattr(
        browser_login,
        "revoke",
        lambda *, host, refresh_token, **kw: revoked.append(refresh_token),
    )
    return revoked


@pytest.mark.usefixtures("local_settings")
def test_login_asks_for_an_api_key_by_default(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
    api_key_logins: list[dict],
):
    result = runner.invoke(cli.cli, ["login", "--host", HOST])

    assert result.exit_code == 0, result.output
    assert fake_login == {}
    [call] = api_key_logins
    assert call["key"] is None
    assert call["relogin"] is False


@pytest.mark.usefixtures("local_settings")
def test_relogin_asks_for_an_api_key(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
    api_key_logins: list[dict],
):
    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--relogin"])

    assert result.exit_code == 0, result.output
    assert fake_login == {}
    [call] = api_key_logins
    assert call["key"] is None
    assert call["relogin"] is True


@pytest.mark.usefixtures("local_settings")
def test_relogin_replaces_a_browser_login_with_an_api_key(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
    api_key_logins: list[dict],
    revocations: list[str],
):
    browser_login.save_credentials(
        credentials_file, wbauth.HostUrl(HOST), make_token_set()
    )

    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--relogin"])

    assert result.exit_code == 0, result.output
    assert fake_login == {}
    assert api_key_logins[-1]["relogin"] is True
    # The browser login would outrank the new API key, so it goes.
    assert (
        browser_login.load_credentials(credentials_file, wbauth.HostUrl(HOST)) is None
    )
    assert revocations == ["wb_rt_stored"]


@pytest.mark.usefixtures("local_settings", "no_verify")
def test_browser_flag_uses_the_browser(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
    api_key_logins: list[dict],
):
    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--browser"])

    assert result.exit_code == 0, result.output
    assert api_key_logins == []
    assert fake_login["host"].url == HOST
    stored = browser_login.load_credentials(credentials_file, wbauth.HostUrl(HOST))
    assert stored.refresh_token == "wb_rt_fresh"


@pytest.mark.usefixtures("local_settings", "no_verify")
def test_browser_flag_passes_the_timeout(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
):
    result = runner.invoke(
        cli.cli, ["login", "--host", HOST, "--browser", "--timeout", "42"]
    )

    assert result.exit_code == 0, result.output
    assert fake_login["timeout"] == 42


@pytest.mark.usefixtures("local_settings", "no_verify")
def test_browser_flag_ignores_other_credentials(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
    monkeypatch: pytest.MonkeyPatch,
):
    """An API key in the environment does not stand in for a browser login."""
    monkeypatch.setenv("WANDB_API_KEY", "k" * 40)

    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--browser"])

    assert result.exit_code == 0, result.output
    assert fake_login != {}


@pytest.mark.usefixtures("local_settings", "no_verify")
def test_browser_flag_reuses_the_stored_browser_login(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
    revocations: list[str],
):
    browser_login.save_credentials(
        credentials_file, wbauth.HostUrl(HOST), make_token_set("wb_rt_old")
    )

    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--browser"])

    assert result.exit_code == 0, result.output
    assert fake_login == {}
    assert revocations == []
    assert "Already logged in" in result.output
    stored = browser_login.load_credentials(credentials_file, wbauth.HostUrl(HOST))
    assert stored.refresh_token == "wb_rt_old"


@pytest.mark.usefixtures("local_settings")
def test_browser_flag_logs_in_again_when_the_stored_login_is_rejected(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
    revocations: list[str],
    monkeypatch: pytest.MonkeyPatch,
):
    browser_login.save_credentials(
        credentials_file, wbauth.HostUrl(HOST), make_token_set("wb_rt_old")
    )
    verifications: list[str] = []

    def verify(self: wbauth.AuthBrowserLogin) -> None:
        stored = browser_login.load_credentials(credentials_file, self.host)
        verifications.append(stored.refresh_token)
        if stored.refresh_token == "wb_rt_old":
            raise AuthenticationError("revoked")

    monkeypatch.setattr(wbauth.AuthBrowserLogin, "verify", verify)

    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--browser"])

    assert result.exit_code == 0, result.output
    assert verifications == ["wb_rt_old", "wb_rt_fresh"]
    assert "no longer works" in result.output
    assert revocations == ["wb_rt_old"]
    stored = browser_login.load_credentials(credentials_file, wbauth.HostUrl(HOST))
    assert stored.refresh_token == "wb_rt_fresh"


@pytest.mark.usefixtures("local_settings", "no_verify")
def test_browser_relogin_revokes_the_login_it_replaces(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
    revocations: list[str],
):
    browser_login.save_credentials(
        credentials_file, wbauth.HostUrl(HOST), make_token_set("wb_rt_old")
    )

    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--browser", "--relogin"])

    assert result.exit_code == 0, result.output
    assert revocations == ["wb_rt_old"]
    stored = browser_login.load_credentials(credentials_file, wbauth.HostUrl(HOST))
    assert stored.refresh_token == "wb_rt_fresh"


@pytest.mark.usefixtures("local_settings")
def test_browser_flag_fails_without_a_browser(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
    api_key_logins: list[dict],
    monkeypatch: pytest.MonkeyPatch,
):
    monkeypatch.setattr(browser_login, "can_open_browser", lambda: False)

    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--browser"])

    assert result.exit_code != 0
    assert fake_login == {}
    assert api_key_logins == []
    assert "needs a browser on this machine" in result.output


def test_browser_flag_rejects_a_key(runner, credentials_file: pathlib.Path):
    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--browser", "k" * 40])

    assert result.exit_code != 0
    assert "Pass an API key or --browser, not both." in result.output


@pytest.mark.usefixtures("local_settings", "no_verify")
def test_browser_login_persists_the_host(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
):
    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--browser"])

    assert result.exit_code == 0, result.output
    saved_settings = wandb_setup.singleton().settings.read_system_settings()
    assert saved_settings.all().get("base_url") == HOST


@pytest.mark.usefixtures("local_settings", "no_verify")
def test_browser_login_cloud_clears_a_saved_host(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
):
    system_settings = wandb_setup.singleton().settings.read_system_settings()
    system_settings.set("base_url", HOST, globally=True)
    system_settings.save()

    result = runner.invoke(cli.cli, ["login", "--cloud", "--browser"])

    assert result.exit_code == 0, result.output
    saved_settings = wandb_setup.singleton().settings.read_system_settings()
    assert saved_settings.all().get("base_url") is None


def test_login_rejects_host_with_cloud(runner, credentials_file: pathlib.Path):
    result = runner.invoke(cli.cli, ["login", "--cloud", "--host", HOST])

    assert result.exit_code != 0
    assert "Cannot use --host and --cloud together." in result.output


@pytest.mark.usefixtures("local_settings")
def test_login_with_a_key_discards_a_browser_login(
    runner,
    credentials_file: pathlib.Path,
    revocations: list[str],
    api_key_logins: list[dict],
):
    browser_login.save_credentials(
        credentials_file, wbauth.HostUrl(HOST), make_token_set()
    )

    result = runner.invoke(cli.cli, ["login", "--host", HOST, "k" * 40])

    assert result.exit_code == 0, result.output
    assert (
        browser_login.load_credentials(credentials_file, wbauth.HostUrl(HOST)) is None
    )
    assert revocations == ["wb_rt_stored"]


@pytest.mark.usefixtures("local_settings")
def test_rejected_api_key_keeps_the_browser_login(
    runner,
    credentials_file: pathlib.Path,
    revocations: list[str],
    monkeypatch: pytest.MonkeyPatch,
):
    browser_login.save_credentials(
        credentials_file, wbauth.HostUrl(HOST), make_token_set()
    )

    def fail(**kwargs: object) -> bool:
        raise AuthenticationError("rejected")

    monkeypatch.setattr("wandb.login", fail)

    result = runner.invoke(cli.cli, ["login", "--host", HOST, "k" * 40])

    assert result.exit_code != 0
    assert (
        browser_login.load_credentials(credentials_file, wbauth.HostUrl(HOST))
        is not None
    )
    assert revocations == []


@pytest.mark.usefixtures("local_settings", "no_verify")
def test_browser_relogin_keeps_the_previous_login_when_the_new_one_fails(
    runner,
    credentials_file: pathlib.Path,
    revocations: list[str],
    monkeypatch: pytest.MonkeyPatch,
):
    browser_login.save_credentials(
        credentials_file, wbauth.HostUrl(HOST), make_token_set("wb_rt_old")
    )

    def fail(**kwargs: object) -> browser_login.TokenSet:
        raise AuthenticationError("cancelled")

    monkeypatch.setattr(browser_login, "login", fail)

    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--browser", "--relogin"])

    assert result.exit_code != 0
    stored = browser_login.load_credentials(credentials_file, wbauth.HostUrl(HOST))
    assert stored.refresh_token == "wb_rt_old"
    assert revocations == []


@pytest.mark.usefixtures("local_settings")
def test_browser_relogin_restores_the_previous_login_when_verify_fails(
    runner,
    credentials_file: pathlib.Path,
    fake_login: dict,
    revocations: list[str],
    monkeypatch: pytest.MonkeyPatch,
):
    browser_login.save_credentials(
        credentials_file, wbauth.HostUrl(HOST), make_token_set("wb_rt_old")
    )

    def fail(self: object) -> None:
        raise AuthenticationError("rejected")

    monkeypatch.setattr(wbauth.AuthBrowserLogin, "verify", fail)

    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--browser", "--relogin"])

    assert result.exit_code != 0
    stored = browser_login.load_credentials(credentials_file, wbauth.HostUrl(HOST))
    assert stored.refresh_token == "wb_rt_old"
    assert revocations == []


@pytest.mark.usefixtures("local_settings")
def test_login_with_no_key_keeps_a_browser_login(
    runner,
    credentials_file: pathlib.Path,
    revocations: list[str],
    fake_login: dict,
    api_key_logins: list[dict],
):
    browser_login.save_credentials(
        credentials_file, wbauth.HostUrl(HOST), make_token_set()
    )

    result = runner.invoke(cli.cli, ["login", "--host", HOST, "--verify"])

    assert result.exit_code == 0, result.output
    assert (
        browser_login.load_credentials(credentials_file, wbauth.HostUrl(HOST))
        is not None
    )
    assert revocations == []
