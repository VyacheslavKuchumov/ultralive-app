#!/usr/bin/env python3

from __future__ import annotations

import asyncio
import fcntl
import os
import subprocess
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path

import slixmpp


def _require_env(name: str) -> str:
    value = os.getenv(name, "").strip()
    if not value:
        raise ValueError(f"Missing required env var: {name}")
    return value


def _bare_jid(value: str) -> str:
    return value.split("/", 1)[0].strip().lower()


@dataclass(frozen=True)
class BotConfig:
    xmpp_jid: str
    xmpp_password: str
    xmpp_server: str
    xmpp_port: int
    xmpp_resource: str
    allowed_senders: set[str]
    repo_path: Path
    deploy_lock_file: Path
    reply_max_chars: int

    @property
    def xmpp_login_jid(self) -> str:
        bare = _bare_jid(self.xmpp_jid)
        if "/" in self.xmpp_jid:
            return self.xmpp_jid
        return f"{bare}/{self.xmpp_resource}"

    @classmethod
    def from_env(cls) -> "BotConfig":
        xmpp_jid = _require_env("XMPP_JID")
        xmpp_password = _require_env("XMPP_PASSWORD")

        bare = _bare_jid(xmpp_jid)
        if "@" not in bare:
            raise ValueError("XMPP_JID must be a valid JID, e.g. ultralive-release-bot@vyachik-dev.ru")

        default_server = bare.split("@", 1)[1]
        xmpp_server = os.getenv("XMPP_SERVER", default_server).strip() or default_server
        xmpp_port = int(os.getenv("XMPP_PORT", "5222"))
        xmpp_resource = os.getenv("XMPP_RESOURCE", "release-bot").strip() or "release-bot"

        allowed_senders_raw = _require_env("XMPP_ALLOWED_SENDERS")
        allowed_senders = {
            _bare_jid(item)
            for item in allowed_senders_raw.split(",")
            if item.strip()
        }
        if not allowed_senders:
            raise ValueError("XMPP_ALLOWED_SENDERS cannot be empty")

        repo_path = Path(os.getenv("REPO_PATH", "/workspace")).resolve()
        lock_file = Path(os.getenv("DEPLOY_LOCK_FILE", str(repo_path / ".release-bot-deploy.lock"))).resolve()
        reply_max_chars = int(os.getenv("BOT_REPLY_MAX_CHARS", "3000"))

        return cls(
            xmpp_jid=xmpp_jid,
            xmpp_password=xmpp_password,
            xmpp_server=xmpp_server,
            xmpp_port=xmpp_port,
            xmpp_resource=xmpp_resource,
            allowed_senders=allowed_senders,
            repo_path=repo_path,
            deploy_lock_file=lock_file,
            reply_max_chars=reply_max_chars,
        )


def _truncate(text: str, max_chars: int = 3500) -> str:
    value = text.strip()
    if len(value) <= max_chars:
        return value

    head = value[: max_chars // 2]
    tail = value[-(max_chars // 2) :]
    return f"{head}\n... output truncated ...\n{tail}"


def _split_for_messages(text: str, max_chars: int) -> list[str]:
    if not text:
        return [""]

    chunks: list[str] = []
    value = text
    while value:
        part = value[:max_chars]
        chunks.append(part)
        value = value[max_chars:]
    return chunks


def _run_command(cmd: list[str], cwd: Path, timeout_seconds: int = 1800) -> tuple[int, str]:
    result = subprocess.run(
        cmd,
        cwd=str(cwd),
        capture_output=True,
        text=True,
        errors="replace",
        timeout=timeout_seconds,
        check=False,
    )
    combined = "\n".join(filter(None, [result.stdout.strip(), result.stderr.strip()]))
    return result.returncode, combined


def run_deploy(config: BotConfig) -> tuple[bool, str]:
    config.deploy_lock_file.parent.mkdir(parents=True, exist_ok=True)

    with config.deploy_lock_file.open("w", encoding="utf-8") as lock_fd:
        try:
            fcntl.flock(lock_fd.fileno(), fcntl.LOCK_EX | fcntl.LOCK_NB)
        except BlockingIOError:
            return False, "Deploy is already running."

        commands = [["/app/deploy.sh"]]
        log_lines: list[str] = []

        for cmd in commands:
            printable_cmd = " ".join(cmd)
            log_lines.append(f"$ {printable_cmd}")

            try:
                code, output = _run_command(cmd, config.repo_path)
            except Exception as exc:  # noqa: BLE001
                log_lines.append(f"Command failed with exception: {exc}")
                return False, "\n".join(log_lines)

            if output:
                log_lines.append(_truncate(output))
            log_lines.append(f"exit code: {code}")

            if code != 0:
                return False, "\n".join(log_lines)

        return True, "\n".join(log_lines)


def run_status(config: BotConfig) -> tuple[bool, str]:
    commands = [
        ["git", "rev-parse", "--abbrev-ref", "HEAD"],
        ["git", "rev-parse", "HEAD"],
    ]

    values: list[str] = []
    for cmd in commands:
        try:
            code, output = _run_command(cmd, config.repo_path, timeout_seconds=15)
        except Exception as exc:  # noqa: BLE001
            return False, f"Status failed: {exc}"

        if code != 0:
            return False, f"Status failed on: {' '.join(cmd)}\n{_truncate(output)}"
        values.append(output.strip())

    branch = values[0]
    revision = values[1]
    return True, f"Repo: {config.repo_path}\nBranch: {branch}\nCommit: {revision}"


class ReleaseBot(slixmpp.ClientXMPP):
    def __init__(self, config: BotConfig):
        super().__init__(config.xmpp_login_jid, config.xmpp_password)
        self.config = config
        self.deploy_lock = asyncio.Lock()

        self.register_plugin("xep_0030")
        self.register_plugin("xep_0199")

        self.add_event_handler("session_start", self.session_start)
        self.add_event_handler("message", self.on_message)

    async def session_start(self, _event):
        self.send_presence()
        await self.get_roster()
        print(f"[release-bot] connected as {self.boundjid.full}")

    def _reply(self, to_jid: str, message: str):
        for chunk in _split_for_messages(message, self.config.reply_max_chars):
            self.send_message(mto=to_jid, mbody=chunk, mtype="chat")

    async def _handle_update(self, to_jid: str):
        async with self.deploy_lock:
            started_at = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%SZ")
            self._reply(to_jid, f"Starting deploy at {started_at}")

            ok, output = await asyncio.to_thread(run_deploy, self.config)
            if ok:
                self._reply(to_jid, f"Deploy completed successfully.\n{output}")
            else:
                self._reply(to_jid, f"Deploy failed.\n{output}")

    async def on_message(self, msg):
        if msg["type"] not in {"chat", "normal"}:
            return

        sender = _bare_jid(str(msg["from"]))
        to_jid = str(msg["from"].bare)

        if sender not in self.config.allowed_senders:
            self._reply(to_jid, "Unauthorized sender.")
            print(f"[release-bot] blocked message from {sender}")
            return

        body = (msg["body"] or "").strip()
        if not body:
            return

        command = body.split()[0].lower()
        if command == "/help":
            self._reply(to_jid, "Available commands: /help, /ping, /status, /update")
            return

        if command == "/ping":
            self._reply(to_jid, "pong")
            return

        if command == "/status":
            ok, output = await asyncio.to_thread(run_status, self.config)
            self._reply(to_jid, output if ok else f"Status failed.\n{output}")
            return

        if command == "/update":
            if self.deploy_lock.locked():
                self._reply(to_jid, "Deploy is already running.")
                return
            asyncio.create_task(self._handle_update(to_jid))
            return

        self._reply(to_jid, "Unknown command. Use /help")


def main() -> int:
    try:
        config = BotConfig.from_env()
    except Exception as exc:  # noqa: BLE001
        print(f"[release-bot] invalid config: {exc}")
        return 1

    bot = ReleaseBot(config)
    connected = bot.connect((config.xmpp_server, config.xmpp_port))
    if not connected:
        print("[release-bot] failed to connect to XMPP server")
        return 1

    bot.process(forever=True)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
