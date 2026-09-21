# AGENTS.md — styrman-harness

Project rules for this repo. They override the global context, which this repo is the source of (`context/`).

## What this repo is

The agent harness: one global context (agents.md standard), skills (Agent Skills `SKILL.md` standard), subagents, and `sh-agent`, the Go binary that installs them into Claude Code and the Antigravity CLI. No hooks, no gates, no CI.

## How the harness works

The map in `README.md` § How it works defines how the harness reaches a machine and how experience comes back. It is the contract, not an illustration, and these rules are strict:

1. **Two ways in, one way back, nothing else.** A place gets the harness from a checkout kept current with `sh-agent update`, or from a release binary with the harness built into it, installed with `sh-agent install` and moved forward by `sh-agent update` (what each kind of place fetches: rule 6's one home). Experience comes back only through `sh-agent report` into `styrman-harness-feedback`, through a throwaway clone that keeps nothing on the machine. No service, no push to a machine, no copy vendored into a project, no timer as the mechanism (R8, R9, ADR-030, ADR-033).
2. **A tag is a release.** `/preflight vX.Y.Z` makes a signed tag and publishes that tag's binaries with signed checksums; what each kind of place does with it is rule 6's one home (D18 of the release-binary plan, ADR-033 superseding ADR-031 D3).
3. **Changing a path is a design question.** Adding, removing or rerouting any line on the map goes to the operator first, gets an ADR, and lands in ONE commit with the code, the docs and the map. A change to install, update, report or release behaviour whose diff does not touch the map is incomplete; `/review --docs` checks for it.
4. **Dashed means a gap, and it stays honest.** A dashed line is done by hand or limited to one machine; it turns solid in the commit that closes the gap, and nothing is drawn solid before it works for every place.
5. **The map is pinned by a test.** `TestHarnessMap` fails when the map is missing, when it names an `sh-agent` verb the binary does not have, or when install, update or report is missing from it. The interactive page of the map (a private claude.ai artifact, not in this repository) is only a view: when the map changes, the page is updated from it, and it is never the source.
6. **The delivery claim has one home.** What a tag delivers, to which kind of place, and how a place moves forward is stated once, in `.agents/docs/specs/sh-agent/delivery.md` § The model; README's prose, this file, the deploy guides and the ADRs point at it in one sentence and never restate it; the README map draws the paths and its labels name them, which is not a restatement, and a guide or README row describing what a COMMAND does at a place (`update` fetches, verifies, replaces, reinstalls) is that command's documentation, not the claim — rule 6 is about which place gets what from a tag, never about how a command works. The v1.3.0 gate spent three review rounds finding the same claim stale in six files after one ADR (ADR-030 → ADR-033); a change to the claim is one file, plus the map's lines when a path moves (rule 3).

## Setup commands

```bash
gofmt -l cmd internal && go vet ./... && go test -race ./...   # before any completion claim; gofmt takes the source dirs, not `.`, because it walks .agents/scratch and the Go tool does not (BUG-048)
# a change to context/*.md or rules/*.md must also pass TestRenderedContextFitsAntigravityBlock (internal/harness/size_test.go,
# part of `go test ./...` above), which pins the rendered context under 22,000 bytes inside Antigravity's 24,000-byte block (ADR-035)
go run ./cmd/sh-agent install --home .agents/scratch/home  # exercise the installer against a scratch home (gitignored)
go run ./cmd/sh-agent list                              # sanity-check skills and agents parse
govulncheck ./...                                       # clean on go1.25.8+ (see below); on an older Go it reports two stdlib advisories — watch for NEW ones
```

`go.mod` carries NO `toolchain` directive, deliberately (ADR-029 withdrew the one ADR-028 added): pinning a patched toolchain made the module unbuildable on a host that has neither that toolchain cached nor a reachable proxy, because `GOTOOLCHAIN=auto` honours a pin by downloading it — measured, including through `go run`, which fails before any of this code executes. On an older Go, `govulncheck` therefore reported two standard-library advisories as affecting (`os` and `os/exec`, reached from the inventory readers and from `Reinstall`'s `LookPath`; the reader path now goes through `os.DirFS`, so the exact call chain a run names will differ). They are fixed in go1.24.6 and go1.25.8, and the fix is upgrading the Go on the machine, never pinning it here. Debian trixie ships neither, so the Linux box runs a user-level go1.27.1 SDK ahead of `/usr/bin/go`, needing no root (`.agents/docs/deploy/local.md` § Operator steps), and reports none since 2026-09-17.

## Layout

```
context/AGENTS.md      the global context — single source, deployed to both vendors
context/dev-setup.md   environment, stack, GCP defaults — appended to the context at install
agents/*.md            subagents in Claude Code format; Antigravity's agent.md is derived at install
skills/<name>/         skills; .agents/skills is a symlink here so `npx skills add` lands in skills/
commands/, rules/      markdown prompts invoked as /name, and extra rules; both deploy to both vendors verbatim
cmd/sh-agent, internal/harness   the binary and its tests
.agents/docs/          this repo's own trail (plans, changes, ADRs, reviews, testing, deploy/{local,remote,gcp}.md, bugs/INDEX.md)
.agents/scratch/       gitignored scratch space; every generated or temporary file goes here, never /tmp
```

Nothing generated is committed. The deployed files are assembled by `sh-agent install`.

## Rules for working here

- The installer must never run against the operator's real `$HOME` during development. Use `--home .agents/scratch/<name>` for every test (the global scratch rule: never `/tmp`). Deploying to the real home is the operator's job.
- `sh-agent install` rebuilds the binary from its own git checkout with `go build` and lets that build install (BUG-006); the running binary is copied only when `go` is missing, when the tree is a plain copy with no `.git` (whose binary was built from that same tree and which may be read-only to the installing user, BUG-015), or when the source is the harness embedded in a release binary, where there is nothing to build from (ADR-033, D13). After a change under `cmd/` or `internal/`, `sh-agent status` from the target home must show a `binary` line without a warning before the install is called done.
- Adding a tool to an agent's `tools:` list means adding it to `toolMap` in `internal/harness/derive.go`. An unmapped tool fails the install on purpose: Antigravity hangs a subagent that names a tool it does not recognise.
- Agent `description:` values are one line, double-quoted if they contain a colon, and free of the words "dispatch prompt" or "harness binary" — there is no binary behind them anymore.
- Every skill directory needs a `SKILL.md` with a `description`; the installer refuses one without it. It also refuses any frontmatter value that a strict YAML parser reads as a comment — an unquoted `#` at the start or after a space or tab — because Antigravity silently drops such a skill (BUG-004); quote the value.
- Adding a skill, agent, command, or rule is a file drop — never a code change. If a new kind of customization needs code, that is a design question, not a quick patch.
- Third-party skills are vendored verbatim and listed in `skills/VENDORED.md`; never edit them in place — fork into a differently named skill.
- Stdlib only in the Go code. The binary is a file copier with a manifest; it does not need a dependency.
- Commit and push after each task, per the global git workflow. This repo is on `main`.

## Testing

There is no test framework beyond `go test`. A change to `context/`, `agents/`, or `skills/` is verified by `go run ./cmd/sh-agent install --home .agents/scratch/x` succeeding and `go run ./cmd/sh-agent status --home .agents/scratch/x` reporting no drift, then by one start of each CLI in the sandbox (`.agents/docs/testing/2026-09-03-verify-sandbox.md`) with its log grepped for parse errors. A change to `agents/*.md` also gets a behavioural run: `claude -p --agents <json built from the working-tree file>` in a scratch clone of a small repository, several runs, with the agent's artifact carrying something only the new definition writes, so the run proves which definition ran — the dispatch text does NOT repeat the changed rule (the definition alone is under test), the prompt goes on stdin because `--allowedTools` is variadic and swallows a trailing argument, `--allowedTools` names the agent's tools or every call is refused at the permission gate and the run proves nothing, the tool-call order is read from the run's subagent transcript under `~/.claude/projects/`, and the result is reported as n of N runs, never "works" (BUG-083); a change to `context/` that should change what the main agent does is run the same way with the new text given through `claude -p --append-system-prompt`. The real home is never touched. A CLI started from a scratch `HOME` cannot authenticate (credentials live in the real home), so the tier-4 CLI start is the sandbox procedure or the real-home smoke checks, never `HOME=<scratch> claude -p`. A verification log states what it did not exercise; "complete" or "all skills" is said only after the full matrix in `.agents/docs/testing/skill-matrix.md` has run. Tested means tested from where it runs: after any install into a real home, the three smoke checks in `.agents/docs/testing/2026-09-03-verify-real-home.md` (context questions — including the canary word GANNET that ends the rendered context, which a truncated file fails, BUG-085 — one skill load, one agent dispatch) run from that home on each CLI before the install is called done. A change to `context/`, `agents/` or `skills/` also gets `/review --docs` over its range; here the living docs also include `context/`, `agents/`, and every skill directory under `skills/` not listed in `skills/VENDORED.md`, searched at `skills/` because `.agents/skills` is a symlink.

## Release

Releases are SSH-signed annotated `vX.Y.Z` tags on `main`, cut only by `/preflight vX.Y.Z` (the gate itself — review, security, verification, docs freshness — is `skills/preflight/SKILL.md` Steps 1–5) from the operator's own turn, from this checkout, on a machine holding the operator's ed25519 key. The tag is signed per command (`-c gpg.format=ssh -c user.signingkey=...`, never a `git config` change, and `git` never under `sudo`) and checked before the push — bound first (the tag takes no target and sits on HEAD, the release commit, whose parent is the gated head and whose source tree is the gated one, on `main`), then annotated, last armor header `-----BEGIN SSH SIGNATURE-----`, `verify-tag` against the allowed-signers line — with the exact commands in `.agents/docs/deploy/remote.md` § Releases. The tag is the only version tracking (no version file; preflight Step 6.2); a pushed tag is never moved; a bad release gets the next PATCH. What a tag delivers, and to which place, is `.agents/docs/specs/sh-agent/delivery.md` § The model (rule 6 above).

## Key files

| Path | Purpose |
|---|---|
| `context/AGENTS.md` | The global context every session reads |
| `internal/harness/plan.go` | What gets installed where |
| `internal/harness/derive.go` | Claude → Antigravity agent translation |
| `REQUIREMENTS.md` | Non-negotiables for the harness |
| `.agents/docs/plans/2026-09-03-rebuild.md` | The rebuild checklist |
| `.agents/docs/deploy/remote.md` | The four places, the resync loop, and the signed-tag release procedure (ADR-030, ADR-033) |
| `.agents/docs/specs/sh-agent/delivery.md` | The delivery model and the feedback loop, end to end |
| `internal/harness/report.go` | `sh-agent report`: collect, redact, send to the feedback repository |
