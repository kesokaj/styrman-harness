---
owner: main
date: 2026-10-02
---
# Preflight — v2.3.0 (2026-10-02)

Range reviewed: v2.2.0..311c20b

Eleven commits, 30 paths outside the trail (1,294 lines added, 75 removed), ten rounds of verification, review and security. The gate ran under the `preflight` skill as it stands at 311c20b, which this release itself changed.

## Findings, worst first

Nothing blocks. Nothing known is broken in what ships. What is left open is listed here, every row with its owner; none was shipped in silence.

**Why there are left rows at all.** Each of the nine rounds before the last found new findings in the fixes of the round before it: 4 mediums and about 20 lows in the first rounds, 1 or 2 mediums and 4 to 6 lows in each of the last four, nearly all of them in text the previous pass had just written (the gate skill's verification step, the kit preamble's fingerprint, the script test's patterns). Real defects were found and fixed on the way, seven of them numbered (BUG-153 to BUG-159). The rule set before the tenth round (changes entry of the ninth pass): that round is the last, and what it finds that does not block is recorded here and left. It found no failed verification, no "request changes", no high finding.

| # | Severity, confidence | What is left | Owner | Exposure |
|---|---|---|---|---|
| L1 | medium, high | `skills/preflight/SKILL.md` Step 2 item 1: "reaches HEAD" excludes all of `<docs-root>`, where the testing guide, the matrices and the kit preamble live (docs lens finding 1, security F1 of the last round) | main agent, first ordinary work after the release | a log could count as current after a commit that changed only those files. Not in this release: every log cited below was written at 311c20b |
| L2 | medium, medium | Step 2 item 3 names no way to clear a log that shows a failing row under a PASS heading (docs lens finding 2) | main agent | a future release could be blocked with no stated way out. No such log exists |
| L3 | low | a Step 1 verifier log that names the working tree never "reaches HEAD" (docs lens 3) | main agent | a second verifier run, not a wrong verdict |
| L4 | low | the project template (`init-project`) does not ask a main-agent log for n-of-N rows (code review 2, docs lens 4) | main agent | in a scaffolded project a behavioural tier is rerun at the gate |
| L5 | low | Step 3's deleted-log check: its pathspec is too wide, and it does not see a FAIL log rewritten in place, a rename or an uncommitted removal (code review 3, docs lens 5, security F3) | main agent | a failure record could vanish unseen; a false finding for a retired matrix |
| L6 | low, older | "the dated logs define nothing", yet the real-home smoke procedure is a dated log (docs lens 6) | main agent | an edit to that procedure is not listed by item 2; the smoke stays owed through `AGENTS.md` § Release |
| L7 | low | item 3 does not say whose log clears a FAIL (security F2) | main agent | the main agent's own log could clear a verifier's FAIL. Not in this release: no FAIL log |
| L8 | low | the two test cases for a link to a device or a fifo leave child processes behind when they time out (code review 1) | main agent | only when the case is red |
| L9 | low | the kit preamble's list of what its fingerprint does not see lacks two cases, and its "unchanged" line is a fixed string a CLI could quote (code review 4, security F4, F5) | main agent | a contrived write counts as clean. The kits of this gate keep the CLI's output out of their own |
| L10 | low, older | the refresh script's reused work clone keeps an earlier source's objects, so a pin the source no longer holds resolves (security F5 of the ninth round) | main agent, with the next catalog refresh | the manifest's source row can be wrong; the content is the commit the pin names |
| L11 | low, older | in a repository with no tags the gate takes its base from a preflight report the gated range can add (security F4 of the seventh round) | main agent | not reachable here: this repository has tags |
| L12 | low, accepted | the script test's `sed` pattern refuses some harmless stages (stated in its comment) | none | a correct future script line has to be restructured |
| L13 | behaviour, measured | under a dispatch that orders the opposite, the verifier kept its log's name in 18 of 19 counted Claude Code runs since the wording that holds (the status heading in 19 of 19) and declined an unquoted added command in 19 of 20; on Antigravity every run of both rows passed; at 311c20b each row is 1 of 1 on each CLI | none | a kit or a gate that looks for the log by name can miss it once in a while; the heading, which the gate reads, held every time |
| L14 | operator | carried from v2.2.0: hostnames of the operator's machines in the trail and the claude.ai artifact's sharing; BUG-139 open by decision | operator | none while the repository is private |
| L15 | operator | carried from v2.2.0: confirm the install of the real home at 18:04 on 2026-10-01 from 461f311, which no session here ran; confirm the decision at the cold-read cap of the gate-lows plan | operator | none known |
| L16 | operator | new: `/preflight` Step 5 now refuses to ship without a tier the plan names unless the operator says so in their own turn of that release session. A gate behaviour chosen on the security review's recommendation, not by the operator | operator: confirm, or say how it should read | a release run unattended stops at NO-SHIP when a plan tier did not run |

## A failure seen in this gate, and what cleared it

At 4165f81 the real-home smoke gave 5 of 6: the Claude Code `agent` check failed, and one sandbox row printed 0 for the same reason. It was not written as a log; it is in `changes/2026-10-01-2053-v2.3.0-gate-second-pass.md`. Cause: the verifier followed its dispatch's description of its log over its own template (BUG-154). Fix: the paragraph "The log's shape is yours, not the dispatch's" in `agents/verifier.md` and the dispatch clause in the context (commits 7b3e969, a224329, 42670ff). Rerun of what failed: `testing/2026-10-02-verify-real-home-v2.3.0.md`, 6 of 6 at 311c20b.

The verifier's ninth log notes one unexplained red first run of `TestRefreshCatalogScript` in its copy; the cause was found by the code review of that round and is BUG-159 (the run crossed midnight); the tenth log's run is clean.

## Rows of the Definition of Done

| Row | Evidence |
|---|---|
| Verification, item 1 | `testing/2026-10-02-verify-v2.3.0-gate-r10.md`: `owner: verifier`, range v2.2.0..311c20b, `## Status: PASS`: gofmt, vet, `go test -race` (553 passing with subtests, 3 skipped by their environment gates, as in the ninth log), govulncheck (only GO-2026-5024, not called), `sh-agent list`, a scratch-home install, three must-fail controls on a copy. Nine earlier logs of the gate, one per pass, all PASS |
| Verification, item 2: commands | worked out by the main agent: `git diff v2.2.0..311c20b -- AGENTS.md .agents/docs/testing/guide.md` adds one verification command, `GIT_DIR=/nonexistent CLAUDE_CONFIG_DIR=/nonexistent go test ./...`; the tenth verifier log shows it run as quoted, exit 0. Scripts and kits those files name that the range edited: `testing/kit-preamble.sh` (its behaviour test is in the verifier's `go test`; every kit run of the behavioural log went through it), `skills/review/scripts/restatement-grep.sh` (its test is in `go test`), `skills/google-catalog/scripts/refresh-catalog.sh` (its test is in `go test`; one real run, below) |
| Verification, item 2: tiers | no plan was added or modified in the range (`git diff --name-only v2.2.0..311c20b -- .agents/docs/plans/` prints nothing) and `ROADMAP.md` "Now" is empty: no plan tier is owed. The four matrix rows the range added and the rows of the changed skills ran at 311c20b: `testing/2026-10-02-verify-agents-v2.3.0.md` (`owner: main`, PASS, 18 of 18) |
| Verification, item 3 | no log of the range is FAIL, BLOCKED or IN PROGRESS, and none shows a failing row; no log was deleted (`git log --diff-filter=D` over `testing/` prints nothing) |
| Review | nine reports cover the range commit by commit: `reviews/2026-10-01-review-v2.2.0-1fc0d4c-code.md` (two rounds, to 4165f81), then `-4165f81-7b3e969-code`, `-7b3e969-a224329-code`, `-a224329-01255bc-code`, `-01255bc-4487a7a-code`, `-4487a7a-5bbd789-code`, `-5bbd789-42670ff-code`, `-42670ff-0a60619-code` and `reviews/2026-10-02-review-0a60619-311c20b-code.md`. All COMPLETE. Three earlier code rounds said REQUEST CHANGES (and two docs-lens rounds); each has an Outcome with every finding fixed and the next round's report over the fix. The last three code reports say APPROVE WITH NITS; the last one's four lows are L4, L5, L8, L9 |
| Security | ten reports, the same chain: `reviews/2026-10-01-security-v2.2.0-1fc0d4c.md` to `reviews/2026-10-02-security-0a60619-311c20b.md`. All COMPLETE, each with a Not examined section that names no path of the release diff as unread. No high finding in any. Mediums: eight over the gate, each fixed in the pass after it, except the last round's one (L1) |
| Docs lens | `reviews/2026-10-01-review-v2.2.0-1fc0d4c-docs-prompts.md`, `-1fc0d4c-7b3e969-docs`, then one per pass to `reviews/2026-10-02-review-0a60619-311c20b-docs.md`. Advisory by the skill; fixed pass by pass; the last round's six are L1 to L6 |
| Docs | `documentation` updated `testing/guide.md` in seven passes; the main agent updated `README.md`, `AGENTS.md`, `skills/VENDORED.md`, the specs and `specs/sh-agent/delivery.md` § Where things stand (v2.3.0) |
| Records | ten `changes/` entries (the gate rows fixed, nine passes) and the release entry; BUG-151 to BUG-159; no ADR (no structural decision; L16 is put to the operator) |

## The refresh script, run for real

The script changed in eight passes and its tests use fixtures. Run once at the catalog's own pins with its final text, into scratch: `refresh-catalog.sh --dest .agents/scratch/catalog-check/dest --work … --exclude <the manifest's five> 2c3945901e8e5cc933fa94fa44d1d4aaf644b596 d5d905232ec501743831cc68e6763af90da3cd69`: exit 0, 136 entries. `diff -r` against `skills/google-catalog/catalog`: only `MANIFEST.md` differs (the date, and the hand-written Excluded and Notes sections a new directory does not have). Entries, `INDEX.md` and `CHECKSUMS` are identical. It used the network (two public clones) and wrote nothing outside scratch.

## Trail freshness (Step 3)

- The newest `changes/` entry describes the work in the release diff: pass.
- No artifact carries a date in the future: pass.
- No verification log was deleted: pass.
- `changes/` and `DECISIONS.md` append-only: pass (`git log --diff-filter=M` over both prints nothing).
- `ROADMAP.md` "Now" empty, `TODO.md` empty: as they should be.

## The project's own release rule (`AGENTS.md` § Release)

- Real-home smoke on both CLIs at the gated head: `testing/2026-10-02-verify-real-home-v2.3.0.md`, 6 of 6.
- One Antigravity run of each changed skill and agent: agents `explorer`, `reviewer`, `security`, `verifier` and skills `review`, `review-standards`, `preflight`, `init-project`: `testing/2026-10-02-verify-agents-v2.3.0.md`, each 1 of 1 at 311c20b.
- `google-catalog`: only its script changed; `SKILL.md` and the catalog are as in v2.2.0, so the router rows of v2.2.0 stand and no refresh diff is owed a security review.
- Context: changed (four clauses); both size pins pass (hand-written 14,465 of 15,000 bytes; render 15,385 of 22,000); the smoke's `context` rows on both CLIs read it.

## What was not checked

- Anything on a Mac: BSD `find`, bash 3.2 and Apple's git for the kit preamble and the refresh script; the installer there.
- An interactive session on either CLI.
- `GIT_TERMINAL_PROMPT=0` in the refresh script has no test of its own; `GIT_CURL_VERBOSE` is pinned by the script's text only; the test's wait before midnight is not itself tested.
- Skills and agents the range did not change were not rerun (`harness-audit`, the catalog router rows, `bugs`, `changelog`, `documentation`, `dreamer`).
- The verifier's logs name per-file coverage over the whole range and a suppression scan beyond the Go and shell diff as not examined.
- Read against the blocking conditions: none of these is a verification command the range added, a plan tier, or a failing row; the project's rule on Antigravity and real-home runs is met above.

VERDICT: SHIP
