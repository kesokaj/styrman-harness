---
owner: main
date: 2026-09-18
---
# install names every file it changes     (2026-09-18)

## Summary

`sh-agent install` (and the install `update` runs) prints one line per file it adds, updates, corrects the mode of, or removes, on every run — not only under `--dry-run`, which used to be the only run that listed anything. Unchanged files stay a count. The operator asked for it ("more info … It should be clear what is installed new or things like that. Go with 1", 2026-09-18). Phase 1 of `.agents/docs/plans/2026-09-18-sh-agent-works-everywhere.md`. Unreleased.

## Type

feature (output), phase checkpoint

## What changed

- **`internal/harness/apply.go` `Install`**: `add     ~/<rel>` when no previous file existed; `update  ~/<rel>` when the bytes differed; `update  ~/<rel> (previous file unreadable: <err>)` when the previous file could not be read — written over as before, the line says why (D6); `mode    ~/<rel> (0666 -> 0644)` when only the permission differed — the probe runs on a dry run too, only the `chmod` is gated (D7, D8), and the file stays in the `unchanged` count; `remove  ~/<rel> (no longer in the source)` on the real run as on the dry run (it was dry-run only). The line for a written file prints after the write succeeds, so a write that fails — a directory at the path, which the rename cannot replace (`EISDIR`) — leaves no line claiming a change (Phase 1 review). The summary line's shape is unchanged; `written` = add + update.
- **README § Daily commands** (`install` row) and **`specs/sh-agent/delivery.md`** say so.

## Tests

`TestInstallNamesEveryChange`: first install lists every managed file as `add` and nothing else; an identical install lists nothing; a changed source file lists exactly its deployed copies as `update`; a removed source file lists `remove … (no longer in the source)`; a mode-only drift lists one `mode` line and is counted unchanged; a `0000` previous file (D9) lists `update … (previous file unreadable: …)` and is written; `--dry-run` lists the same `update` and `mode` lines and leaves bytes and modes as they were; a directory at a destination fails the install with no line for it. Eight must-fail mutations in `.agents/scratch/phase1-mustfail/mutate.py`, each failing as a test failure. Live: 260 `add` lines then none; a touched skill in a scratch clone → its two vendor copies as `update`; a chmodded file under `--dry-run` → one `mode` line, mode untouched.

## Definition of Done

`verifier` PASS (`testing/2026-09-18-verify-sh-agent-works-everywhere-phase-1.md`, 227/0/0, all four checkpoint legs live — the agent hit its turn limit once with an empty log and finished on a continuation; and `…-phase-1-2.md` over the final range); `reviewer` APPROVE WITH NITS (`reviews/2026-09-18-review-b393525-028a43b.md`, one low fixed — the line now prints after the write); `/review --docs` APPROVE (`…-docs.md`); the docs clause done in the same change. Found on the way, filed as Phase 2 boxes: a relative `--home` with `--repo` resolves differently on a real install and a dry run (the verifier's observation); the recon's four coverage gaps.
