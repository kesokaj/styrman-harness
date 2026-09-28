# Requirements — styrman-harness

> Immutable without the operator's explicit approval. Agents never edit this file.

## R1 — One source, two vendors

`context/AGENTS.md` (plus `context/*.md`) is the single global context, deployed to Claude Code and the Antigravity CLI. Subagents are authored once in Claude Code format under `agents/`; the Antigravity flavour is derived at install time. Nothing generated is committed.

## R2 — Standards, not vendor formats

The context follows the agents.md standard and is named `AGENTS.md` at every level. A project has no `CLAUDE.md`: Claude Code 2.1.277+ reads `AGENTS.md` natively, and a one-line `@AGENTS.md` stub is written by hand only where it cannot (before 2.1.277, Bedrock, Vertex, Foundry, or telemetry disabled). `~/.claude/CLAUDE.md` stays a one-line import: Claude reads a user-level `AGENTS.md` only as an ancestor file, and a project with its own `CLAUDE.md`, or one outside the home directory, would then lose the global context. Native reading also loads an `AGENTS.md` in any directory above the project and, on read, in any subdirectory; a template `AGENTS.md` never sits beside source. Skills follow the Agent Skills `SKILL.md` standard so that anything from skills.sh or any other source drops into `skills/` unchanged.

## R3 — No hooks, no gates, no CI

The harness is context, skills, and subagents. Behaviour comes from what is written in them, not from lifecycle hooks, commit gates, or pipelines. Git is backup and versioning only.

## R4 — One binary, stdlib only, easy on and off

`sh-agent` is a single Go binary with no dependencies. `install` is idempotent and overwrites only files it manages; `uninstall` removes exactly what its manifest lists; `update` pulls the repo from git and reinstalls. `update` refuses a checkout with uncommitted changes on every path, because it deploys the working tree rather than the commit it fetched; `install` stays the deliberate way to deploy a dirty tree. A hand-written `~/.claude/CLAUDE.md` is never overwritten. `install` also sets one key in the operator's `~/.claude/settings.json`, `attribution` (commit and PR trailers off), only when the key is absent; it never changes or removes another key, and `uninstall` leaves it.

## R5 — Hand-added customizations are first class

Adding a skill, subagent, command, or rule is: put the file in `skills/`, `agents/`, `commands/`, or `rules/` (by hand, by `npx skills add` from the repo root, or by `sh-agent adopt`), then `sh-agent install`. No registration, no manifest edits, no code. The installer maps each directory to its place under both vendors.

## R6 — Development never touches the real home

Every test and every development run of the installer uses `--home` with a scratch directory. Deploying to the operator's real `$HOME` is the operator's action.

## R7 — Projects remember and stay resumable

Every project keeps its scratch in `.agents/scratch/` (gitignored), its failures as numbered records in `.agents/docs/bugs/`, its progress in plan files with task checkboxes, and its retro and distill records in `.agents/docs/retros/`, so a zero-context session can resume and does not repeat a mistake. The harness's own state (manifest, feedback log, the harness proposal queue, the cross-project index) lives in `~/.config/sh-agent/` and nowhere else; a project's retro and distill records are the project's own state, in its trail (ADR-015, clarified 2026-09-13 with the operator's approval).

## R8 — One delivery channel, and it is `sh-agent`

The harness reaches a machine one of three ways: a checkout the operator owns; a release's `sh-agent` binary with the harness built in; or a zip of the repository that the operator downloads with their own credential and passes to `sh-agent update --from-zip <file>`. Each is kept current with `sh-agent update`: the checkout and the binary as today, and a zip place by the operator repeating the download and `--from-zip`. There is no other channel — no vendored mirror inside a consumer project, no per-project copy, no service. A project that wants the harness uses `sh-agent`; this repository names no consumer project and depends on none. Experience flows back the one way too: `sh-agent report` sends a COPY of this machine's corrections, proposals and a metrics record to the private feedback repository. The originals stay in `~/.config/sh-agent/`, which R7 keeps as the one place that state lives; a copy leaving for the operator's own private index is not a second home for it. (Approved by the operator 2026-09-17, on the exact text above; ADR-030. First sentence replaced 2026-09-17 on the exact text above, approved by the operator in their own turn — "Approve as written" — ADR-033. First sentence replaced 2026-09-28 on the exact text above, approved by the operator in their own turn — "approve as written" — ADR-038.)

## R9 — Four procedures, set in stone

The harness is operated by four procedures and no others, drawn in `README.md` § How it works:

- **Install:** a place either clones `styrman-harness` with the operator's own credential and runs `sh-agent install`, or downloads a release's `sh-agent` binary with that credential — `gh release download`, `go install` over SSH, or the Releases page — and runs `sh-agent install`, which installs the harness built into the binary. Either way installed files go into the vendor directories and no repository is created.
- **Update:** the operator edits, verifies, commits and pushes `main` in their checkout. A checkout place runs `sh-agent update` when the operator decides — fetch, refuse a dirty checkout, fast-forward, rebuild, reinstall. A binary-only place runs `sh-agent update`, which downloads the newest release, verifies its signed checksums, replaces the binary and reinstalls; it changes only when a release is cut. Nothing is pushed to a machine.
  - **Zip place:** the operator downloads the repository zip by hand and runs `sh-agent update --from-zip <file>`. It unpacks the zip into a temporary directory, refuses any entry that would leave it, installs that tree with the running binary, and deletes the directory. Nothing is fetched, no repository is created, and nothing in the zip replaces the binary. When the release's `SHA256SUMS` and `SHA256SUMS.sig` sit beside the zip, the zip is verified against them and refused if either check fails; a zip with neither beside it is refused unless `--unverified` is passed, and is then installed marked `UNVERIFIED`; a zip with only one of them beside it is refused.
- **Report back:** each place runs `sh-agent report`, which sends a redacted copy of `~/.config/sh-agent/` and a metrics record to `styrman-harness-feedback` through a throwaway clone that holds only that place's folder and is deleted after the push; the harness is improved from the reports in the operator's checkout, and that change returns through Update.
- **Release:** `/preflight vX.Y.Z`, from the operator's own turn, gates the tree, pushes an SSH-signed tag on `main`, and publishes that tag's binaries and a zip of its tree (`git archive`), with checksums signed by the same key. Checkout places follow `main`; binary-only places update to the newest release.

The four procedures and the map change together, only with the operator's explicit approval, an ADR, and one commit carrying the code, the docs and the map. (Approved by the operator 2026-09-17, on the exact text above; ADR-031. The four bullets replaced 2026-09-17 on the exact text above, approved by the operator in their own turn — "Plan approved" on the release-binary plan, which quoted them verbatim — ADR-033. The zip-place sub-bullet under Update added 2026-09-28 on the exact text above, approved by the operator in their own turn — "approve as written" — nested inside Update rather than a fifth procedure, on the operator's "Go with your recommendations" — ADR-038. Its last sentence reworded 2026-09-28 from "the binary itself is not replaced", which the running binary's install made imprecise, approved by the operator in their own turn — "All approved". The Release bullet's zip and the zip place's verification sentence added 2026-09-28, approved by the operator in their own turn — "All approved" — ADR-039. Its unsigned-zip clause replaced 2026-09-28 on the exact text above (refused unless `--unverified`; one file alone refused), proposed in the v1.9.0 handoff and approved by the operator in their own turn — "Go approved all" — ADR-040.)

## Known limits

- **Antigravity global rules are one 24,000-byte block.** CLI 1.2.7 builds one block from `~/.gemini/config/rules/*.md` and cuts it at 24,000 bytes with a `<truncated N bytes>` marker; a second file costs the first one room and is dropped without a marker (BUG-085, measured 2026-09-21; the vendor page says 12,000 characters per file). The deployed context is pinned under 22,000 bytes by a test, and `rules/` files share the block.
- **Antigravity global rules path.** The CLI's customization root is `~/.gemini/config/`; rules there are documented for the workspace form (`.agents/rules/*.md`) and by analogy for the global one. Verified on the Linux box 2026-09-03 (`.agents/docs/testing/2026-09-03-verify-real-home.md`): the CLI reads rules, skills, and agents from there. Confirmed on the Mac for Antigravity on 2026-09-04 (`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`); Claude Code on the Mac is unverified.
- **Antigravity tool names cannot be enumerated.** An unrecognised name in an agent's `tools:` hangs the subagent, so `sh-agent install` fails on any tool it cannot map rather than guessing.
- **Antigravity IDE is out of scope.** Only the `agy` CLI is targeted; the IDE's global paths are unverified.

## Non-goals

- Hooks of any kind, on either vendor. Researched and rejected 2026-09-03 (ADR-007, `.agents/docs/plans/2026-09-03-research-hooks.md`); the harness is enforced by reading, not by gates.
- Plugin or marketplace packaging on either vendor: both namespace skills and agents, which would rename all of them.
- Editing `settings.json` or any vendor setting: the harness is files in the customization directories, nothing else — with one exception, R4's `attribution` key, which `install` writes only when absent because that setting exists solely to switch off a vendor behaviour the context's Git rule forbids; never a permission, never another key.
- Symlinking the vendor directories to one shared store: neither CLI documents reading the other's directory, and a copy has no dangling-symlink failure mode.
