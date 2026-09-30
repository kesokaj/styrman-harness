---
owner: main
date: 2026-09-30
---
# Preflight — v1.13.0   (2026-09-30)

Range reviewed: v1.12.1..c63f736

Release request: the operator, 2026-09-29, in their own turn: "Cut new version update docs. This is dod", and the version "v1.13.0 (Recommended)". A MINOR release: new verbs (`sh-agent search`, `sh-agent index`), one dependency, nothing an existing place must react to.

## Findings, worst first

None blocking. Every finding below was fixed before the tag or accepted by the operator, with its Outcome in the report it came from.

| Severity / confidence | Finding | Where | Outcome |
|---|---|---|---|
| low / high | a work-tree rename misparsed; a path outside the checkout indexed | gate security, finding 1 (`reviews/2026-09-29-security-v1.12.1-c872598.md`) | fixed at 505732d (`--no-renames`, BUG-108), test failed first |
| low / high | a skip-worktree file under a symlinked directory read outside the checkout | security check of the fixes, finding 2 (`reviews/2026-09-30-security-c872598-505732d.md`) | fixed at de1a88f (every read through an `os.Root`), test failed first; the check of that fix found no other way out |
| low / high | a `search.db` copied in with a project directory is trusted | gate security, finding 2 | left: operator — "Document the limit (Recommended)"; spec § Failure modes |
| low / high | GO-2026-4970 on `os.Root.Open` for a Go before 1.25.12 | security check of fix 2 (`reviews/2026-09-30-security-505732d-de1a88f.md`) | left: not exploitable (an `Lstat` of the same path precedes every open); the two known advisories named in `AGENTS.md` and README at c63f736 |
| low / high | the symlink test could not fail on its write-through check; an `O_EXCL` comment overpromised | gate code review (`reviews/2026-09-29-review-f41e5c0-c872598-code.md`) | fixed at 505732d |
| low / high | the read-through test covered one trigger shape | security check of fix 2, finding 2 | fixed at c63f736 (both shapes) |
| low / medium | the `gitView` root has no test of its own | security check of fix 2, finding 3 | left: git lists a modified or untracked file only where the path is a regular file, so no fixture reaches that read; recorded, not hidden |
| medium / high (docs, advisory) | two closed bug records' bodies still say "not fixed" | gate docs lens (`reviews/2026-09-29-review-v1.12.1-c872598-docs.md`) | rebutted: a record's body is history, its status the one later edit; the lens's brief excludes bug records |

## Rows

- **Reviewer** (`model: opus` on the gate dispatch): `reviews/2026-09-29-review-search-index-code.md` (c04f2d7..565bb1e, APPROVE), `-code-r2.md` (565bb1e..aae16af, APPROVE with Outcome), `reviews/2026-09-29-review-aae16af-f41e5c0-code.md` (APPROVE), `reviews/2026-09-29-review-f41e5c0-c872598-code.md` (APPROVE WITH NITS, fixed). v1.12.1..c04f2d7 is text only (R4, ADR-045, three doc lines), read by the docs lens `reviews/2026-09-29-review-8ee5d70-HEAD-docs.md`. c872598..c63f736 (the three gate-fix commits: 100 lines of Go and tests, two doc lines) is under the reviewer threshold and was read by the main agent and by the two security checks named below.
- **Security**: `reviews/2026-09-29-security-v1.12.1-c872598.md` read all 94 paths of the release diff, nothing under Not examined; `reviews/2026-09-30-security-c872598-505732d.md` and `reviews/2026-09-30-security-505732d-de1a88f.md` read the gate-fix diffs end to end. de1a88f..c63f736 (a test widened, `configView` back on the plain reads the first report read, two doc lines) had no security dispatch: it reverts to code that report covered.
- **Verification**: `testing/2026-09-30-verify-v1.13.0-gate-r4.md`, PASS at c63f736 (366 pass with subtests, 0 fail, 2 skips with reasons; three cross-builds and the Go 1.24.4 build; embed and context tests; govulncheck); the earlier gate logs at c872598, 505732d and de1a88f PASS.
- **Docs freshness** (Step 4): `reviews/2026-09-29-review-v1.12.1-c872598-docs.md` APPROVE WITH NITS, the nits rebutted above; the earlier docs lenses over the range fixed three nits (ADR-045a) and found nothing else.
- **Trail freshness**: pass — the newest `changes/` entry (`2026-09-30-0505-v1.13.0-gate-fixes.md`) describes the release's last work; no future `date:`; `changes/` gained eight additions and no modification; `DECISIONS.md` gained ADR-045 and ADR-045a and the status marks on ADR-002 and ADR-029, the one edit a past ADR allows; `ROADMAP.md` and `TODO.md` were touched four times since the base and are empty.
- **The project's own blocking condition** (`AGENTS.md` § Release: one Antigravity run of each changed skill and agent): `testing/2026-09-29-search-agy-gate.md` (10 of 10 scenarios), `testing/2026-09-29-preflight-verdict-last-behavioural.md` (the preflight skill, one Antigravity run), `testing/2026-09-29-search-index-behavioural.md` (`reviewer`, `explorer`, `documentation`, `bugs`, `changelog` on Antigravity). The gate-fix commits changed no skill, agent or context. Real-home smoke on both CLIs: `testing/2026-09-29-verify-real-home-search.md`, PASS.

## What was not checked

- The Mac and the devbox: installed after the release by the operator's `sh-agent update` there; `linux/arm64` and `darwin/arm64` binaries were built, not run.
- Interactive sessions on either CLI; every behavioural run was headless.
- Antigravity's `$(…)` auto-denial in headless mode with narrow grants: two skills ran there with the per-run skip flag inside the sandbox; from the real home the limit does not apply.
- A release range in which the old preflight text would have written SHIP against its own unchecked items; the new rule's ordering effect was measured, its verdict effect was not distinguishable on the fixtures.
- `gitView`'s root read (finding above), unreachable by a fixture.
- Efficiency over days of work in one project; the measured runs are single sessions and three-session lanes.

Every item above was read against the blocking conditions: a verifier FAIL, BLOCKED or IN PROGRESS log (none: r4 is PASS); a high security finding (none); a missing security report, or one naming a release-diff path under Not examined (none); a reviewer REQUEST CHANGES or high finding (none); a `REQUIREMENTS.md` violation (none: R4 as amended, R7, R8 hold); a project agent shadowing a review agent (none); the project's Antigravity-run rule (met for every changed skill and agent, the preflight skill included). None of the unchecked items is named by a blocking condition.

VERDICT: SHIP
