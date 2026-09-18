---
owner: main
date: 2026-09-18
---
# The sh-agent audit: every verb on every kind of place     (2026-09-18)

## Summary

Phase 2 of `plans/2026-09-18-sh-agent-works-everywhere.md` — the operator's Definition of Done for this session: "work until the delivery mechanism sh-agent is flawless". Five review lenses read `cmd/sh-agent` and `internal/harness` whole, an explorer mapped every verb × place-kind cell, and everything they found was fixed with a test that failed first or rebutted in writing; three review rounds of the fixes followed, two of them stopped at the operator's door because a fix had introduced the next round's defect. Phase 3, the matrix on this box, ran on the result: 53 cells, none silent. Unreleased.

## Type

fix (audit); phase checkpoint

## What changed

Behaviour a user of `sh-agent` sees:

- **Release tags are canonical.** `v1.02.2`, `v01.2.2`, `v1.2.02` and a component that does not fit an int are not release tags: `IsReleaseVersion` refuses them, so `update` never downloads such a tag, and a tag that is not newer than the running binary is refused with "the newest release X is not newer than this binary Y; nothing installed".
- **Refusals name the alternative.** `update` on a plain copy: "not its own git checkout; a plain copy cannot be updated with git — `install` still deploys it as it is; to move forward, use a checkout (`git clone …`) or a release binary (README § Install without cloning)". `adopt` on a release binary names the checkout to clone; on a plain copy its "next:" line says to copy the skill into a checkout and commit there, not "commit and push" a tree with no repository.
- **The BUG-079 note** prints only for the user's real home (never under `--home` elsewhere), identity by `os.SameFile` so a symlinked `$HOME` counts, and offers `ln -sfn` only when the file's directory is the caller's — otherwise "it is not yours to relink — remove it, with the rights that own its directory, or leave it and run the installed one by its path".
- **`~/.config/sh-agent` is private.** Created 0700; an existing looser one is tightened by `install` before it writes anything, with `mode    ~/.config/sh-agent (0755 -> 0700)`; a dry run names it and touches nothing; a symlinked directory and one the caller does not own are left alone; a failed `chmod` is a note. `diff` and `report` create it private and never change an existing one.
- **A relative `--home` beside `--repo`** lands under the invoking directory on a real install as on a dry run (BUG-080).
- **`adopt`** strips group/other write bits like `install` does. **`diff`** renders under `~/.config/sh-agent/`, never `/tmp`, and removes its temp file on every path. **`install`/`status`** report an oversized `settings.json` with one note. A description with a control character renders as valid YAML for Antigravity; a table cell never cuts a multi-byte character.
- **Docs**: the `keep` line, the `knownCheckout` refusal (README's "works from anywhere" qualified), `report: --repo is not used here`, and every changed message quoted where it is described.

Not changed, by decision: the ~220-line `settings.json` editor stays in `main.go` (a Follow-up); `LoadState` does not sanitise `Repo` (reason narrowed in the record); loose modes inside the runtime dir; ownership as the proxy for relinkability.

## Tests

`internal/harness/audit_test.go` (11 tests, one per finding) and `owned_unix_test.go`; eight tests in `cmd/sh-agent/main_test.go` (plain-copy `update`/`adopt`/`status` through the CLI, `adopt` on a release binary, `sync`, `diff` with drift, the relative `--home`, the oversized settings). Nineteen must-fail mutations in `.agents/scratch/phase2-mustfail/mutate.py`, each failing as a test failure; the ones that could not compile or could not be reached were redone or dropped and say so in the reports. Suite 250, `govulncheck` clean, `GOOS=windows` and `darwin/arm64` builds clean.

## Definition of Done

- **Audit** (`reviews/2026-09-18-review-harness-merged.md`, 19 findings, `## Outcome`): code lens APPROVE WITH NITS; security (opus) 1 medium + 7 low; simplify 1 medium + 3 low; docs 2 medium + 1 low; tests PASS with 2 gaps. Explorer recon `plans/2026-09-18-recon-sh-agent-audit.md`: no silent cell, four coverage gaps.
- **Round 1** (`3209cf3`): verifier PASS (244); reviewer (opus) APPROVE WITH NITS, 3 low; security 1 medium (the 0700 was a no-op after `SaveState`'s 0755) + 2 low; docs lens APPROVE.
- **Round 2** (`e7fec88`): verifier PASS (248, cross-builds); reviewer APPROVE WITH NITS — two mediums on the new `ConfigDir` (symlink followed; `EPERM` fails the run); security 1 medium — the same mechanism, worse (an install aborted after every write, no manifest). **The operator's call** (option 1): tighten as a `mode` line at the top of `install`, never fail.
- **Round 3** (`42840fb`): verifier PASS (250; four live legs quoted, the `mode` line before the first `add`); reviewer APPROVE, one low/low (README's overpromise); security four lows, no high, the medium closed structurally. After the cap, in the checkpoint commit: the two skips print a `note:` (the symlink one pinned by a test), README's clause corrected; Phase 4's gate reviews those lines.
- `documentation` ×3 (audit docs findings; the Phase 2 messages; the runtime-dir mode line); `bugs` for BUG-080.

## Where it stops

Unreleased; Phase 4 of the plan is the release. The Follow-ups are in the plan.
