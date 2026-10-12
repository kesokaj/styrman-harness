---
owner: main
date: 2026-10-12
---
# Preflight — v2.7.0   (2026-10-12)

Range reviewed: v2.6.2..f0bc3e6

A MINOR release, the operator's choice of version (2026-10-12, "v2.7.0 (Recommended)").
- **New behaviour:** the agent frontmatter allow-list (ADR-056, ADR-056a), the stub list's home in the init-project skill (ADR-057, with R2 edited on approval), and wider credential redaction.
- **Fixes:**
  - BUG-167, BUG-168, BUG-169 (the project review);
  - BUG-170, the dream job under nohup;
  - BUG-171, `--vendor` checked first;
  - the dream signal relay and the timer's home check.
- **Removed:** nothing a harness place uses.
- **The one change a place may notice:** a hand-added agent with a key outside the allow-list fails `install`.

## Failures and their clearing (FAIL first)

| # | What failed | Severity / confidence | Status |
|---|---|---|---|
| 1 | Antigravity `harness-audit` row, first run: print timeout after 25 m with the turn in progress, no report (`.agents/scratch/gate-v270/agy-rel-harness-audit.log`); its credential line flagged nine files, every one a test fixture of this repo or the kit's own regex text (checked by pattern and against `git grep`), deleted by name and re-checked clean as the kit says | — | cleared: re-run `harness-audit-v270b` (45 m print timeout, `GOCACHE` outside the sandbox home) met the row's criterion, applied to the files it committed in its throwaway clone (the kit's `git status` check misses committed files; the widened check fails on the first run and passes here); `testing/2026-10-12-verify-real-home-v2.7.0.md` cites both runs |
| 2 | `testing/2026-10-12-verify-ee20d00.md` was BLOCKED on tier 9 (the context question on both CLIs) | — | cleared in place: Status PASS, `Previous run:` line, on the real-home `pointer` rows (`testing/2026-10-12-verify-real-home-v2.7.0.md`) |
| 3 | The abandoned gate of 2026-10-09: reviewer stub IN PROGRESS, smoke runs failed on expired logins | — | superseded: this gate's reviewer, security and smoke ran over the whole range; the stub was deleted, its security review committed (`reviews/2026-10-09-security-v2.6.2-66f3412.md`) |

## Step 1 — pre-gate work

- **Commits:**
  - `ee20d00`, `0d8bb88`, `13c5bde`: the resume items (`changes/2026-10-12-0848-resume-items-allowlist-bug170.md`).
  - `f0bc3e6`: the gate fixes (`changes/2026-10-12-1143-gate-v2.7.0-fixes.md`, BUG-171, ADR-056a).
- **Every DoD row ran:** reviewer and security over both, verifier after, docs by the main agent (under 30 lines each), and the bug records.
- **Secret scan of the staged diff** (no gitleaks on this box): no hit.

## Step 2 — review and verification

- **Reviewer** (opus, code+docs), `reviews/2026-10-12-review-v2.6.2-13c5bde-code-docs.md`: Status COMPLETE, APPROVE WITH NITS in both rounds. `Range reviewed: v2.6.2..13c5bde plus the working-tree fixes of round 2`.
- **Security**, `reviews/2026-10-12-security-gate-v2.6.2-13c5bde.md`: Status COMPLETE, FINDINGS, all low or info.
  - Every release-diff path was read; `## Not examined` names none.
  - The range touches input parsing (frontmatter lint, redaction, the timer's path check) and secrets (redaction).
- **The earlier rounds on the same work:** `reviews/2026-10-12-review-worktree-merged.md` (the resume items).
- **Verification item 1:** `testing/2026-10-12-verify-v2.7.0-gate.md`, owner verifier, `Range reviewed: v2.6.2..f0bc3e6`, Status PASS.
  - Tests: gofmt and vet clean; `go test -race` gave 697 pass, 0 fail, 5 skip.
  - govulncheck: nothing new.
  - Install and report: drift none; `report --dry-run` sends nothing.
  - Refusals: three refusal messages checked against § 3c.
  - Timer: `systemd-analyze verify` clean.
  - 13 new tests pass.
- **Verification item 2 (commands):** the range changed one verification command, govulncheck by its full path (`AGENTS.md`), which the gate verifier ran. `testing/guide.md` is unchanged in the range, and no plan was added or modified that names tiers.
- **Verification item 2 (main-agent tiers):** the Release rule's real-home smoke on both CLIs, plus one Antigravity run of each changed skill (`init-project`, `review`, `harness-audit`). These are in `testing/2026-10-12-verify-real-home-v2.7.0.md` (owner main, `Range reviewed: v2.6.2..f0bc3e6`):
  - smoke 8 of 8, with the Antigravity store check;
  - scaffold PASS, review PASS;
  - harness-audit 1 of 1 against its criterion (row 1 above).
- **Verification item 3:** no FAIL or BLOCKED log is left standing (row 2 above).

## Step 3 — trail freshness

- The newest `changes/` entry describes this work: pass (the release entry follows in the release commit).
- No `date:` in the future: pass.
- No verification log deleted in the range: pass (`git log --diff-filter=D` prints nothing).
- `changes/` append-only: pass (no modified past entry). `DECISIONS.md` append-only: pass (0 deleted lines in the range).
- `ROADMAP.md` touched (Last reviewed); `TODO.md` empty: report only.

## Step 4 — docs freshness

- The reviewer's docs lens found no stale wording in either round. Its round-2 docs finding (`update --vendor all`) was fixed in `f0bc3e6`.
- The release commit updates `specs/sh-agent/delivery.md` § Where things stand for v2.7.0: the newest tag, the refused-key clause and a v2.7.0 sentence. This was named by both lens agents; `TestDeliverySpecCurrent` runs on it before the tag.

## Decisions on non-blocking findings (Step 6)

- **Fixed before the gate, in `f0bc3e6`:** every round-1 finding, and round 2's quoted OpenSSL `pass:`, the five pinned over-redactions, the `update --vendor all` docs and `disallowedTools` (ADR-056a).
  - The operator's cap decision (2026-10-12): "Fix small ones, no round 3 (Recommended)".
  - Those round-2 fixes were not re-reviewed by the lens agents. Each has a test that failed without it, and the gate verifier ran after them.
- **Accepted, named in the release entry:**
  - After the self-kill, `dream --all` could start the next project's CLI in the microseconds before the process dies. The reviewer rated it low, with low confidence. A fix cannot be tested without killing the test process.
  - Redaction's known limits: `net use … /user:`, not chosen by the operator, and any shape the command-line class misses.

## What was not checked

- **The Mac real home** (`2026-09-04-verify-mac-real-home.md`): no Mac session here.
- **The § 3c verb × place-kind matrix in full:** the kit is absent. The changed refusal messages were checked instead.
- **The § 3g real-data ingest check:** the range's `internal/metrics` change is a refactor (`trail.DocsRoot`, `readSmall` through a `LimitReader`) covered by unit tests, including `TestReadSmallStopsAtOneMiB`. No plan names § 3g as a tier.
- **Mutation runs by the lens agents** (outside their lane). The main agent ran every new test red before its fix, and ran one must-fail mutation (`readSmall`'s cap).
- **The round-2 fixes of `f0bc3e6`** were not re-reviewed by the lens agents (operator's cap decision, above).
- **A live Claude Code install that refuses an agent with a now-refused key:** the lint is unit-tested; the real-home install ran with the harness's own eight agents.
- **The `index --rebuild` timing:** the gate verifier measured 46.1 s under the parallel load of three Antigravity runs and the race suite. The main agent re-timed it idle at 7.6 s, against the guide's 10 s.

Every blocking condition of preflight Step 5 and AGENTS.md § Release was read against the sections above and the list of what was not checked:
- every Antigravity run of a changed skill exists (`init-project`, `review`, `harness-audit`), and no agent changed;
- the previous release, v2.6.2, was not a bridge;
- nothing here names a missing tier, an unexamined release-diff path or a standing failure.

VERDICT: SHIP
