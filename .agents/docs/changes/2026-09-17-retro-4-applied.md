---
owner: main
date: 2026-09-17
---
# Third retro applied: six prevention rules graduated, two release rules learned     (2026-09-17)

## Summary

The session that shipped v1.1.0 logged no user correction — the operator's turns were approvals and one question — so this retro's findings came from the trail it wrote: seven bug records and thirteen gate findings. Nine proposals, all approved and applied here in the harness checkout, so every place picks them up with `sh-agent update`.

Six graduate a bug's prevention rule into the global context, where they now reach every session on every machine: a guard in front of a destructive action fails closed (BUG-069 deleted the only copy of unpushed work when it could not count it); a side effect goes below a read-only mode's early returns, with a test that the read-only mode does NOT perform it (BUG-070 made `--check` delete files); a generated or embedded artifact is pinned to its source in BOTH directions (BUG-066 embedded a gitignored `.env` while every test stayed green); a must-fail mutation asserts its target is unique, restores from a saved copy and ends with `git diff` read (BUG-068); a test's unusual setup is measured before the test is trusted (BUG-071); and the tree stays frozen while a verifier reads it (BUG-072 cost a whole verification run). Two came from cutting a real release for the first time: release artifacts are uploaded one call at a time, and a non-blocking gate finding is still a decision the gate report records. One is local to this repository: § Releases now says what publishing does not change.

## Type

task

## Files changed

- `context/AGENTS.md` — three new Hard rules bullets (destructive guard, read-only mode, generated artifact), two clauses added to the existing must-pass/must-fail sentence (mutation hygiene, measured premise), and the Definition of Done's verifier row gains the frozen-tree sentence
- `skills/preflight/SKILL.md` — Step 6's release-artifacts item now says to create the release and upload each artifact as its own call, with the measured reason; a new paragraph says a non-blocking finding is still a recorded decision, and that fixing after the gate invalidates the verification row
- `.agents/docs/deploy/remote.md` — § Releases states that publishing changes no repository visibility and gives the check, and documents the one-call-per-asset upload
- `.agents/docs/retros/2026-09-17-retro-4.md` (new) — the record; `.agents/docs/bugs/INDEX.md` — the graduated rows

## Agents dispatched

`bugs`, to move the graduated rows to `INDEX.md`'s Graduated table, each verified against `context/AGENTS.md` before moving. BUG-067 deliberately stays in the main table: its rule was already BUG-036's graduated rule and was simply not followed.

A second `bugs` dispatch recorded BUG-073: writing this retro, the main agent overwrote `2026-09-17-retro-3.md`, which already held an earlier session's third retro of the day, because it chose the suffix from what it remembered writing rather than from `ls`. The clobbered record was restored from 8af1875 (verified byte-identical), this one became `-retro-4`, and the changes entry, the harness queue and the index line were repointed. Nothing was lost only because git held it.

## Verification

`gofmt`, `go vet` and `go test -race -count=1 ./...` clean. `go run ./cmd/sh-agent install --home .agents/scratch/retro3-home` wrote 260 files with `status` reporting drift none, and each of the six new context passages was found in BOTH deployed copies (`~/.claude/AGENTS.md` and `~/.gemini/config/rules/AGENTS.md`); `sh-agent list` parses every skill, which covers the `skills/preflight/SKILL.md` edit. Behavioural check for the context change, as this repo's Testing section prescribes: `claude -p --append-system-prompt "$(cat …/.claude/AGENTS.md)"` asked a fresh session to quote the destructive-guard and read-only-mode rules, and it returned both verbatim — the text is loadable and read, not merely written. NOT exercised: the two-CLI sandbox start, because `/home/kesokaj/sandbox/home` no longer exists on this box and re-creating it needs the operator's own CLI logins; the frontmatter lint in the installer plus the `list` parse are what stand in for it here, and no frontmatter changed.

## Docs

`.agents/docs/deploy/remote.md` is the docs half of this change and is included above. The retro record carries every proposal's diff and rationale; the harness queue file `~/.config/sh-agent/proposals/2026-09-17-retro.md` carries the eight harness-scope ones with their statuses.

## Requirements

None touched. R4 holds — these are file edits, no code. R8 and R9 are unaffected: nothing about how the harness reaches a machine changed, only what the context tells a session once it is there.
