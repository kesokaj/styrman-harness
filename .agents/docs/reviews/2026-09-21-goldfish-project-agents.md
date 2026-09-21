---
owner: reviewer
date: 2026-09-21
---
# Goldfish review — project `AGENTS.md`   (2026-09-21)

Range reviewed: not a diff — goldfish audit of `AGENTS.md` (repo root) and the files it references by path, per the audit prompt Step 5.

## Status: COMPLETE

Files read for Part A ("this file and what it references"): `AGENTS.md`, `README.md`, `.agents/docs/specs/sh-agent/delivery.md`, `.agents/docs/deploy/remote.md`, `REQUIREMENTS.md`, `.agents/docs/plans/2026-09-03-rebuild.md` (first 60 lines), `internal/harness/plan.go` (first 80 lines), `internal/harness/derive.go` (full, short file), `internal/harness/report.go` (first 60 lines), `context/AGENTS.md`, `context/dev-setup.md`, `skills/VENDORED.md`, `.agents/docs/testing/2026-09-03-verify-sandbox.md` (first ~39 lines), `.agents/docs/testing/2026-09-03-verify-real-home.md` (full), `.agents/docs/testing/skill-matrix.md` (first 30 lines), `.agents/docs/deploy/local.md` (full). Part B added: `.agents/docs/changes/2026-09-18-release-v1.4.0.md` (full), `skills/preflight/SKILL.md` (full — deliberately NOT read for Part A, since `AGENTS.md` names it only as the `/preflight` command, never as a path).

`AGENTS.md` itself: 11,853 bytes → **~2,963 tokens** (bytes ÷ 4).

## Part A

### G1 — Comprehension

styrman-harness is a personal agent harness: one global context file (`context/AGENTS.md` + `context/dev-setup.md`, agents.md standard), skills (`SKILL.md` standard), subagents (`agents/*.md`, Claude format, Antigravity derived at install), and `sh-agent`, a single stdlib-only Go binary that installs/updates/inspects/removes all of it for Claude Code (`~/.claude`) and the Antigravity CLI (`~/.gemini/config`). No hooks, no gates, no CI — behaviour comes only from what is written and read.

How the system works (README's map, pointed to by `AGENTS.md` § How the harness works, not reproduced there): edit `context/`, `skills/`, `agents/`, `commands/`, `rules/` in the checkout → `gofmt`/`vet`/`test -race` → commit + push `main` → `/preflight vX.Y.Z` gates the tree and cuts an SSH-signed annotated tag, publishing that tag's three platform binaries with a signed `SHA256SUMS`. Two kinds of "places" exist: (a) a checkout, kept current by `sh-agent update` (fetch, refuse dirty tree, `pull --ff-only`, rebuild, reinstall — a tag changes nothing for it); (b) a binary-only release install, which `update` moves forward only when a newer signed release exists (download, verify signature+checksum, replace the binary, reinstall). `sh-agent install` writes files into `~/.claude` and `~/.gemini/config`; no repository is ever created there. Feedback flows back one way only: `sh-agent report` sends a redacted copy of `~/.config/sh-agent/` to a private feedback repo through a throwaway clone, removed after the push.

Must (per `AGENTS.md` itself): run the Setup commands gate before any completion claim; use `--home .agents/scratch/<name>` for every installer test, never the real `$HOME`; after a `cmd/`/`internal/` change, confirm `sh-agent status`'s `binary` line is warning-free before calling install done; add every new agent tool to `toolMap` in `derive.go`; keep agent `description:` one line, quoted if it has a colon, free of "dispatch prompt"/"harness binary"; every skill needs a `SKILL.md` with `description`, frontmatter values quoted; add a skill/agent/command/rule as a file drop only; never edit vendored skills in place; Go code is stdlib-only; commit and push after each task, on `main`.

Must not: run the installer against the operator's real `$HOME` during development; bypass `toolMap` for a new tool; treat a path/behaviour change to install/update/report/release as done without updating the map (rule 3) in the same commit.

### G2 — Critic (kept only if behaviour-changing) — 2 findings

1. **Rule 6's claim-vs-documentation line is genuinely ambiguous for the exact restatement grep it mandates.** It says a guide/README row describing what a COMMAND does at a place "is that command's documentation, not the claim" — the boundary between "documents what update does" and "restates which place gets what from a tag" is a judgment call with no worked example. An agent grepping subject nouns after a delivery-model ADR (as this rule and the global "restatement" hard rule both require) could over-flag command-behaviour rows as stale claims, or under-flag an actual restatement disguised as command documentation.
2. **`govulncheck` guidance names the two known-stale advisories only informally** ("`os` and `os/exec`, reached from the inventory readers and from `Reinstall`'s `LookPath`"), not by advisory ID, and admits "the exact call chain a run names will differ" — plausible to misclassify a genuinely new advisory as "the known two" if wording drifts.

### G3 — Readiness: release v1.4.0, using only `AGENTS.md` + its by-path references — 4 questions, 2 blocking

1. **[BLOCKING]** What are `/preflight`'s own gate steps 1–5 (review/security/verification dispatch, trail freshness, docs freshness) — the steps that precede the "Step 6" `remote.md` documents ("what `/preflight` runs at its Step 6 here... the skill's Step 6.3 follows a project's Release section")? Neither `AGENTS.md` nor any file it names by path gives these; `skills/preflight/SKILL.md` is named only as the command `/preflight`, never as a path.
2. **[NON-BLOCKING]** Is `v1.4.0` the correct next version for the pending changes? The MINOR/PATCH/MAJOR rule itself IS answered (`context/AGENTS.md` § Git, a by-path reference); confirming it against the actual pending diff just needs an ordinary `git log`, not a missing document.
3. **[BLOCKING]** What is `<gated head>` for this release — which commit did the gate actually approve? `remote.md`'s tag procedure takes it as a given external input ("the gate report's `Range reviewed:` line") but nothing in the reference set explains how that gate report is produced — same root cause as Q1.
4. **[NON-BLOCKING]** Where do this release's `changes/` entries (needed for release notes) live — the format is documented (`remote.md` § Release notes) but enumerating "entries dated after the previous tag's commit" is an ordinary `git log`/`ls`, not a documentation gap.

## Part B — comparison against the repository

### G1 errors against README.md, remote.md, delivery.md, changes/2026-09-18-release-v1.4.0.md, skills/preflight/SKILL.md

**No factual errors found** — G1 was built from README.md, remote.md and delivery.md directly (all three are in the by-path reference set), and nothing in `changes/2026-09-18-release-v1.4.0.md` or `skills/preflight/SKILL.md` contradicts it.

**Missing points** (both traceable to the same reference-set gap as G3 Q1/Q3, not to a misreading):
- G1's "how the system works" narrative compresses `/preflight` to "gates the tree and cuts a tag." The actual gate (`skills/preflight/SKILL.md` Steps 1–5) dispatches `reviewer` (opus) and `security` conditionally, requires a `verifier` PASS log reaching HEAD, checks trail and docs freshness, and writes a SHIP/NO-SHIP report before anything is tagged — none of which G1 could state, since the skill file is never linked by path from `AGENTS.md`.
- G1 does not (could not) capture that showing the release notes to the operator is explicitly **not** a gate in this repo's actual practice — `changes/2026-09-18-release-v1.4.0.md` records "no pause after showing the notes," matching `skills/preflight/SKILL.md` Step 6.1's "showing is not a gate: print them and continue... in the same run." This is real, current behaviour that a goldfish reading only `AGENTS.md` + its references would not know to expect.

### Proposed fixes (fewest bytes, in the file a goldfish would look at)

- **Fix for G3 Q1/Q3 (blocking) and the missing points above** — one clause added to `AGENTS.md` § Release, immediately after "`/preflight vX.Y.Z`": `` (the gate itself — review, security, verification, trail/docs freshness — is `skills/preflight/SKILL.md` Steps 1–5) ``. This turns the only command-name reference into a path reference, closing the gap for every future goldfish read at the cost of one parenthetical.

## Verdicts

**SUFFICIENT — NO, for the release task; YES, for routine work.** For day-to-day work in this repo (setup gate, testing tiers, layout, hand-added customizations, git rules) `AGENTS.md` + its by-path references are sufficient to act correctly on the first pass. For "release v1.4.0 as this repo did it" specifically, the file is **not** sufficient on its own: the actual gate mechanics live in `skills/preflight/SKILL.md`, which `AGENTS.md` names only as a command, never as a path — an agent following only "this file and what it references" would tag and push without ever dispatching `reviewer`/`security`/`verifier` or writing the SHIP/NO-SHIP report Step 5 requires.

**MINIMAL — two spans safe to cut without causing a mistake** (both are "why" narrative, not "what to do"; the operative instruction survives their removal):
- `AGENTS.md` line 29, the middle of the toolchain paragraph — from "pinning a patched toolchain made the module unbuildable..." through "...which fails before any of this code executes." The operative rule ("`go.mod` carries NO `toolchain` directive... the fix is upgrading the Go on the machine, never pinning it here") does not depend on the ADR-028/ADR-029 history to be followed correctly; the history is what `.agents/docs/architecture/DECISIONS.md` is for.
- `AGENTS.md` line 18 (rule 6), the closing sentence — "The v1.3.0 gate spent three review rounds finding the same claim stale in six files after one ADR (ADR-030 → ADR-033); a change to the claim is one file, plus the map's lines when a path moves (rule 3)." The anecdote (first clause) is trail colour; only the second clause ("a change to the claim is one file...") is the operative rule an agent must follow.

## Not examined

`internal/harness/plan.go` past line 80, `internal/harness/report.go` past line 60, `.agents/docs/plans/2026-09-03-rebuild.md` past line 60, and `.agents/docs/testing/skill-matrix.md` past line 30 were only partially read — sufficient for gist (all four are long implementation/procedure files whose remainder is more of the same pattern already sampled), not exhaustively verified line by line. No other files in the by-path reference set, or in the Part B comparison set, were skipped.
