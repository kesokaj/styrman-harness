---
owner: main
date: 2026-09-18
---
# BUG-078 — a verified release binary shadowed by an earlier PATH entry     (2026-09-18)

## Summary

On devbox the operator downloaded v1.2.1, verified its signature and checksum, saw `./sh-agent-linux-amd64 version` print `v1.2.1`, moved it to `/usr/local/bin` — and `sh-agent` kept running a stale `(devel)` build from `/opt/devbox/bin`, which installed an old tree and reported `drift none`. Two PATH entries sat ahead of the new binary and nothing said so. This change makes `install` and `status` say so, and makes both install guides verify the *file* before it moves and the *shell's resolution* after. Unreleased: it reaches binary-only places at the next release.

## Type

bugfix

## What changed

- **`internal/harness/pathguard.go`** (new): `PathShadows(pathEnv, binPath) (shadows []string, onPath bool)` walks PATH as a shell does. `onPath` is true when the walk reaches an executable that `os.SameFile`s the installed file — through its own directory, a symlinked directory, a symlink or a hard link; there is no lexical directory comparison anywhere. `shadows` are the executables of that name found before it, in order, without repeats; an empty PATH entry prints as `./sh-agent`. `executableAt` is what it calls a binary (`os.Stat`, not a directory, execute bit set).
- **`cmd/sh-agent/main.go`**: `install` prints `warning: <path> comes before <installed> on PATH, so `sh-agent` runs that one` when PATH reaches the installed file but only after an earlier entry; when PATH does not reach it and something else answers, the old "not on PATH" note names it instead (`… so `sh-agent` runs <path> instead; put <bin dir> ahead of it`). `status` gains a `shadowed   <path> comes first on PATH, so `sh-agent` runs that, not the binary above` line under the same condition. `install --dry-run` prints no PATH claim: it wrote no binary. One `PathShadows` call feeds each; there is no second definition of "on PATH".
- **README § Install without cloning**: the `gh` route is one `&&` chain — download into `mktemp -d` with `--clobber`, the real signer line written to `allowed_signers`, `ssh-keygen -Y verify`, `sha256sum -c`, `./<asset> version` on the file, then the one `mv` and the scratch dir removed — so a failed check stops the move and no checksum or signer file lands in a PATH directory. The `go install` route says what it trusts (SSH to GitHub and git's hashes; no checksum database, no release signature) and the browser route is the same chain over the files in `~/Downloads`; each route is its own fence, so a paste runs one route. Then: `sh-agent version` *unqualified* first, with the sentence that it tells stale from fresh and never genuine from forged.
- **`.agents/docs/deploy/remote.md` § First install**: the same procedure as one chained block for a release-binary machine; § Verify asks for the unqualified `version` and no `shadowed` line. **`local.md` § Verify** says the same for a checkout place.
- **`.agents/docs/bugs/BUG-078-…`** and `INDEX.md`: `fixed`, with what the code half does *not* catch stated plainly (below).

## What it does not do

On devbox as it happened, the wrong binary ran the install and wrote itself to `~/.local/bin`, first on PATH — from its own point of view nothing was shadowed, so no in-binary check could have fired. The doc half catches that ordering: `sh-agent version` unqualified and `which -a`. The code half catches the other ordering, where `/usr/local/bin`, `/opt/…` or a Homebrew dir precedes `~/.local/bin`, which is common on a Mac.

## Tests

`TestPathShadows` (14 cases over real temp files: installed dir first, earlier entry, order and later entries, symlink, hard link, symlinked directory before and after a shadow, repeated entry, empty directory, non-executable, missing installed dir, `//`, `/.`, empty entry via `t.Chdir`, missing installed file), `TestExecutableAt`, `TestInstallAndStatusNameAPathShadow` (warning, `shadowed`, both quiet after a symlink to the installed binary), `TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice` (the off-PATH note names what answers; no double warning; `--dry-run` makes no PATH claim). Ten must-fail mutations in `.agents/scratch/bug078-mustfail/mutate.py`, each failing as a test failure — three first attempts did not compile (an unused variable each) and were redone, since a build failure proves nothing.

## Definition of Done

Four review rounds, because two landed on one mechanism and the context rule then sends the mechanism to the operator:

- Round 1 (`52574a0`): `verifier` PASS (`testing/2026-09-18-verify-bug078-path-shadow.md`); `reviewer` REQUEST CHANGES — high: `binOnPath` (string compare) and `PathShadows` (`Clean`) disagreed on `.local/bin//`, hiding a real shadow; medium: a hard link reported as a shadow; low: an unbacked "for an hour" in README. `security` 1 medium (README's signer line a placeholder, the file left in a PATH dir) + 6 low. `/review --docs` APPROVE, 2 low (`local.md` narrowed; no changes entry). All fixed except two lows rebutted (`reviews/…-4cdc0b5-52574a0*.md` § Outcome).
- Round 2 (`8ce58a0`): `verifier` PASS; `reviewer` REQUEST CHANGES — the same mechanism again, lexical `onPath` vs `SameFile` on a symlinked directory. `security`: two mediums, both docs (`remote.md`'s block a sketch with `$D` unassigned; README unchained at the verify), plus the rebuttal reasons corrected. **Stopped here and asked the operator what the gate was for**; their choice: "on PATH" is the identity walk reaching the file, nothing lexical.
- Round 3 (`b029ea1`): `verifier` PASS, the alias repro run by hand; `reviewer` APPROVE; `security` all low. The three-round cap; the operator granted one bounded extra round for the cheap fixes.
- Round 4 (`28e4747`): `verifier` PASS (215 tests/subtests, unchanged; a fresh-home `--dry-run` printed no PATH line); `security` four lows, no high — it validated the dry-run gate in both states itself. One low was a docs-only edit made in the closing commit (the browser fence became the same `&&` chain as the `gh` route); the others corrected two `left` reasons and confirmed the hung-mount reason (`reviews/2026-09-18-security-4cdc0b5-28e4747.md` § Outcome).
- `documentation` for `local.md` § Verify; `bugs` for the record and its `fixed` status, read against the diff and corrected (a mistyped SHA and `os.SameFile` attributed to the wrong commit).
- Every round's findings have a `## Outcome` line: fixed with file:line, rebutted with a checkable reason, or `left` with an owner.

## Where it stops

Unreleased. `context/AGENTS.md` does not yet carry the rule ("a verified install is what the shell resolves, not what was moved into place"); the bug record marks it a graduation candidate.
