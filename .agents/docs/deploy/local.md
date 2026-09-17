---
owner: main
date: 2026-09-17
target: local
---
# Deploy — local (this machine)

"Deploying" the harness means installing the context, skills, and agents into a home directory. Development never touches the real `$HOME` (`REQUIREMENTS.md` R6); deploying to it is the operator's action.

## Prerequisites

Go 1.24+ and git. `~/.local/bin` on `PATH`. Claude Code and/or the Antigravity CLI installed; both are optional — the installer writes for the vendor recorded by the existing install, else both, or selectively with `--vendor`.

## Build and install

```bash
cd ~/workspace/styrman-harness
gofmt -l cmd internal && go vet ./... && go test -race ./...   # before deploying anything; gofmt takes the source dirs, not `.` (it walks .agents/scratch, BUG-048)
go run ./cmd/sh-agent install                        # builds the binary from this checkout, writes ~/.claude and ~/.gemini/config, installs the build to ~/.local/bin/sh-agent
# or for Antigravity only:
# go run ./cmd/sh-agent install --vendor antigravity
sh-agent status                                      # must end with "drift none"
```

## Update

```bash
sh-agent update            # fetch; refuse a dirty checkout; if behind: pull --ff-only; rebuild and reinstall
sh-agent update --check    # only reports; exit 2 when behind
```

The normal way a checkout moves forward, run whenever the operator decides the harness has changed; `remote.md` § The resync loop covers the timer and the dirty-checkout refusal.

## From an unpacked bundle

A copy of the repository with no `.git` — for a major remodel, or a machine that cannot authenticate to GitHub. It needs Go, and nothing else:

```bash
cd <unpacked directory>
go run ./cmd/sh-agent install        # installs the running binary as is; never rebuilds (BUG-015)
sh-agent status                      # first line names a plain copy with no commit; must end with "drift none"
```

`sh-agent update` refuses a plain copy, by design. To move it forward, unpack a newer copy and run `install` from it again; to roll back, install the earlier one.

## Report back

`sh-agent report` sends this machine's experience to the private `styrman-harness-feedback` repository so the harness can be improved from it. What leaves: a COPY of `~/.config/sh-agent/feedback.jsonl`, `retros.jsonl`, `state.json` and `proposals/*.md`, redacted, plus `report.json` — host, OS and arch, vendor, installed file count and time, the harness checkout's commit, drift, skill and agent counts, the Go that built the running binary, and the list of redactions. The originals are never modified. A symlink among them is skipped and named, never followed.

```bash
sh-agent report --dry-run     # prints where it would write, every file, and every redaction; writes nothing, clones nothing
sh-agent report               # clones the feedback repository into ~/.config/sh-agent/feedback on first use, then commits and pushes
```

- **Redaction is structural, and no longer walks into JSON.** In `feedback.jsonl` and `retros.jsonl`, a line whose raw text OR JSON-decoded text carries a needle is replaced whole by `{"redacted":"[redacted by sh-agent report]"}`; a key block (`-----BEGIN` — or a `BEGIN … PRIVATE KEY` header whose dashes were rewritten — to its footer, two concatenated keys kept open, unless the dashes were also stripped and only a space joins them, in which case the second key's body is sent — a stated limit) is found first and replaced by `{"redacted":"[redacted by sh-agent report: N lines]"}`. On the line its header opened, a mentioned footer closes the block only if that mention contains the header's label ("... up to -----END below" does not close it, because "below" is not the label); a footer on any later line closes the block regardless of what it names, so a footer naming a different label, or one merely mentioned on a later line, still closes it early. `state.json` is one document: any needle anywhere replaces the whole file with the plain JSON marker. Text files (`proposals/*.md`) get the same whole-line or whole-block treatment with the plain text marker, never the JSON object. The shapes: vendor prefixes as written (`-----BEGIN`, `PRIVATE KEY`, `AKIA`, `ASIA`, `AIza`, `ya29.`, `GOCSPX-`, `ghp_`, `gho_`, `ghs_`, `ghu_`, `github_pat_`, `glpat-`, `sk-ant-`, `sk-proj-`, `xoxb-`/`xoxp-`/`xoxa-`/`xoxs-`), and case-insensitive: an assignment to a `password`, `passwd`, `secret`, `token` or `api_key`-style name (`GITHUB_TOKEN=`, `DB_PASSWORD:`), `Bearer <value>`, `Authorization: <anything>` (a quoted header name matches too: `"Authorization": …`), and `scheme://user:pass@` (the user may be empty: `redis://:pass@host`). `report.json` records file, line (absent for a whole document), pattern, lines and unterminated — never a JSON key path, and never the text. It over-redacts on purpose — "tokens: 3" or "ASIAN" go too, and an unterminated `-----BEGIN` in a JSONL file redacts to end of file, dropping later records rather than risk shipping the rest of the key (both `--dry-run` and a real run show `unterminated: everything after it in this file was dropped` next to that line, and `report.json` carries `unterminated: true`).
- **Anything that matches no shape IS sent, and `--dry-run` cannot show it:** the dry run lists what was caught, never what was missed. This is a safety net, not a guarantee, so before EVERY report, not just the first, read what was added under `~/.config/sh-agent/` since the last one. Known gaps in the net: a secret on a later line than the needle that located it — a YAML `password: |` or `>-` block, an indented continuation, an HCL heredoc, a shell `\` continuation, pretty-printed JSON — is sent, because redaction only removes lines a needle itself matched. The decoded check reads every JSON value that decodes from the start of the line, numbers kept as written, but this is not exhaustive: a leading BOM, `NaN`, or anything else that stops decoding partway — a leading log prefix, a raw control character, a truncated record, `Infinity`, very deep nesting — is checked only as raw text from that point; a value shadowed by a duplicate key decodes fine but is discarded before the decoded check ever sees it, so only the raw-text check (run on the line regardless) can still catch it. Text files (`proposals/*.md`) are checked only as raw text, never decoded, so a needle reachable only through a JSON escape there is missed.
- **Where it goes:** `reports/<host>/<YYYY-MM-DD-HHMMSS>/` in the feedback checkout, suffixed `-2`, `-3` when that directory exists, so no run overwrites another. Only the paths the run wrote are staged and committed; any other file in the checkout is named in the output and left alone. The push names the feedback URL itself rather than `origin`, so `origin`'s own `pushurl` or an extra `url` cannot redirect it; a git URL rewrite (`url.*.insteadOf`, `pushInsteadOf`) or a remote named after the feedback URL still applies, as it does to every push on the machine — that is the operator's own configuration, trusted rather than refused. `-u` points the branch at that URL so a later plain `git pull --rebase` fetches from the same place.
- **Only the feedback repository.** The checkout — `~/.config/sh-agent/feedback` or `--feedback-repo DIR` — must be the top of its own git checkout with `origin` set exactly to `git@github.com:kesokaj/styrman-harness-feedback.git`, or `report` refuses before writing anything; a harness checkout or any other project fails that. `--feedback-repo` is never cloned into. `report` also refuses to run with `GIT_DIR`, `GIT_WORK_TREE`, `GIT_INDEX_FILE`, `GIT_OBJECT_DIRECTORY` or `GIT_COMMON_DIR` set, and refuses a `reports` or `reports/<host>` path in the checkout that is not a real directory.
- **No harness checkout here** (an unpacked bundle, or nothing): `report` still runs; `harness_commit` is `null`, and with no harness tree at all so are `drift`, `skills` and `agents`. With nothing installed, `vendor`, `files` and `installed_at` are `null`. With nothing under `~/.config/sh-agent/` it still sends `report.json` — "this machine is alive" is the signal.

| What fails | What you see and do |
|---|---|
| First-use clone fails (no SSH key, no network, repository missing) | `report: could not reach <url>: …; clone it by hand or pass --feedback-repo`, exit 1, nothing left behind. Fix access, or clone by hand into `~/.config/sh-agent/feedback`. |
| Push rejected — another place pushed first | `report: push rejected; run git -C <feedback> pull --rebase and re-run — the commit stays local`, exit 1. Pull, then re-run; the rejected commit goes with the next push. |
| `~/.config/sh-agent/feedback` exists but is not a checkout | Refused, exit 1. Move it aside or pass `--feedback-repo`. |
| The checkout's origin is not the feedback repository | `report: <dir> is not the feedback repository: its origin is not "…"`, exit 1, nothing written — the checkout's own origin is never printed. Point it at the right checkout. |
| `GIT_DIR` (or another redirecting variable) is set | `report: GIT_DIR is set; report commits only in the feedback checkout — unset it and re-run`, exit 1. |

## Configure

Nothing to configure. `--vendor VENDOR` selects the target CLI: `antigravity` (or `agy`) or `claude`; the default is the vendor recorded by the existing install, else `all` — passing `--vendor all` explicitly widens a single-vendor install to both through `install`; `update` keeps the recorded vendor instead. `--repo DIR` overrides where the checkout is; `$SH_AGENT_REPO` does the same persistently. `--home DIR` installs elsewhere and is what every development run uses (a scratch home under `.agents/scratch/`, never `/tmp`). From its own git checkout, `install` rebuilds the binary with `go build` first; the running binary is installed as is when `go` is missing (with a warning) or when the tree is a plain copy with no `.git` — an unpacked bundle (BUG-015).

## Verify

- `sh-agent status` shows the current commit, a `binary` line carrying the same revision (a warning there means the installed binary is older than the checkout or unstamped: run `sh-agent update`), and no drift. `install` writes `"attribution": {"commit": "", "pr": ""}` in `~/.claude/settings.json` when the key is absent (creating the file when needed, every other key untouched) and prints a `settings:` line saying so; `status` prints an advisory note when the file is missing or the key is present with a different value.
- `sh-agent render | head -3` prints the assembled context.
- Start each CLI once and check it read the files: Claude Code `/memory` lists `~/.claude/CLAUDE.md` importing `AGENTS.md`, `/skills` and `/agents` list the harness entries; Antigravity `/skills` and `/agents` do the same, and `grep 'Failed to parse' ~/.gemini/antigravity-cli/log/cli-*.log` is empty (BUG-004).

## First session after install

Claude Code asks once per session before the first read or write under `~/.config/sh-agent/` (feedback log, proposals, the retro index), because it is outside the project. Approve it, or add the directory to `permissions.additionalDirectories` in `~/.claude/settings.json` yourself — the installer touches no vendor setting except the attribution key (ADR-025). Antigravity prompts the same way for paths outside the workspace.

## Operator steps

Standing by-hand procedures on this box, never pending items (those live in the newest handoff):

- After a push from another machine: `sh-agent update`, then `sh-agent status` shows a stamped `binary` line and no drift.
- After any change to `context/`, `agents/`, or `skills/` reaches this home: the three smoke checks of `.agents/docs/testing/2026-09-03-verify-real-home.md` on each CLI, from a fresh clone of the demo service.
- When `govulncheck ./...` names a toolchain-fixed advisory: get a newer Go, never raise the `go` directive or add a `toolchain` line in `go.mod` (ADR-029 — a pin makes the module unbuildable on a host that cannot fetch a toolchain). Try the system package manager first; when it has nothing newer, as Debian trixie did on 2026-09-16 (`golang-1.24-go 1.24.4-1` and no later package), install a user-level SDK instead, which needs no root:

```bash
go install golang.org/dl/go1.27.1@latest     # the wrapper lands in ~/.go/bin
go1.27.1 download                            # the SDK lands in ~/sdk/go1.27.1
ln -sfn "$HOME/sdk/go1.27.1/bin/go" "$HOME/.local/bin/go"   # ~/.local/bin is first on PATH here
go install golang.org/x/vuln/cmd/govulncheck@latest         # rebuild it, or it cannot parse the newer sources
go version && govulncheck ./...
```

  Undo by deleting that one symlink; `/usr/bin/go` is untouched and the system package stays where it was. After the switch, run the full gate and `sh-agent install`, so the deployed binary is built by the new toolchain — `go version -m ~/.local/bin/sh-agent` names it.

## Roll back

`sh-agent uninstall` removes exactly the files the manifest lists. To return to an earlier harness commit: `git checkout <commit> && go run ./cmd/sh-agent install`.

## Sandbox instead of the real home

`.agents/docs/testing/2026-09-03-verify-sandbox.md` records how a throwaway `HOME` with both CLIs logged in is used to test a change end to end before it reaches the real home.
