---
owner: main
date: 2026-09-30
---
# Preflight — v1.13.1   (2026-09-30)

Range reviewed: v1.13.0..8529f0c

Release request: the operator, 2026-09-30, in their own turn: "Go with your recommendations. Then cut a new version". A PATCH: fixes to released `update` and `report` behaviour, one shared atomic-write helper, provenance for two forked skills, an accepted-risk ADR. No skill, agent, command, rule or context file changed (`git diff --stat v1.13.0..f00d059 -- skills agents context commands rules` shows `skills/VENDORED.md` alone, a manifest the installer does not deploy), so the project's own blocking condition (an Antigravity run of each changed skill and agent) has nothing to require.

## Findings, worst first

None blocking. Every finding was fixed before the tag or accepted with its Outcome in the report it came from. The gate agents ran at f00d059; the pre-tag fixes and their records are commit 8529f0c on top of it (the gated head), whose source diff over f00d059 (65 lines, three of six files tests) the main agent read and round 2 of the gate log verified, so the release commit carries records only.

| Severity / confidence | Finding | Where | Outcome |
|---|---|---|---|
| low / high | `sk_test_`/`rk_test_` prefixes match inside `task_test_`, wiping a feedback record | gate code review, 1 (`reviews/2026-09-30-review-v1.13.0-f00d059-code.md`) | fixed before the tag: a length-pinned `stripe key` needle; two benign rows fail without it (gate log § Round 2, mutation a) |
| low / high | the deny-list misses `ghr_`, GitLab `gldt-`/`glrt-`/`glptt-`, a Google refresh token and families the guide never names, while the guide reads as whole-vendor coverage | gate security, 1 (`reviews/2026-09-30-security-v1.13.0-f00d059.md`) | fixed in part before the tag: the four prefixes added (mutation b), the guide's wording narrowed to what is matched and naming what is not; the other families left: operator — deny-list residue accepted (the guide says unmatched shapes are sent; the operator's read before each report is the guarantee) |
| low / medium | ADR-046 rests on a false premise (a signature check would add an independent root) and names `sh-agent install --from-zip`, which does not exist | gate security, 2 | fixed by ADR-046a, appended (`architecture/DECISIONS.md`) |
| low | two tests check source text; the settings-write error branch has no test | gate code review, 2 | fixed in part: the ordering pin reads code lines only, `TestAttributionSetReportsAnUnwritableSettingsFile` drives the branch (mutation c); rebutted for the two pins themselves, which cover what no behavioural test can observe and say so |
| low | `delivery.md:56` leaves out the Slack webhook rule | gate code review, 3 | fixed before the tag |
| low | the two forks keep their upstream names against the vendoring rule | gate code review, 5 (older than the range) | left: operator — the recommendation the operator took; `AGENTS.md` and `VENDORED.md` state the exception |

## Rows

- **Reviewer** (`model: opus`): `reviews/2026-09-30-review-v1.13.0-f00d059-code.md`, APPROVE WITH NITS over the whole range, every regexp valid RE2, the removed names gone; Outcome appended. Earlier over parts of the range: `reviews/2026-09-30-review-worktree-code.md` (three rounds over e59638c's source, final APPROVE), the seven project-review artifacts of 12ef66b. The pre-tag fix diff (65 lines, six files, three of them tests) was read by the main agent, not re-dispatched.
- **Security** (`model: opus`): `reviews/2026-09-30-security-v1.13.0-f00d059.md` read all 42 release-diff paths end to end, trail included, nothing from the release diff under Not examined; two lows, Outcome appended. Earlier: `reviews/2026-09-30-security-project.md`, `-security-worktree.md` (three rounds), `-security-worktree-redaction.md`, each with Outcome. `govulncheck`: the known GO-2026-5024 only (ADR-045a).
- **Verification**: `testing/2026-09-30-verify-v1.13.1-gate.md`, PASS at f00d059 (372 pass, 0 fail, 2 skips with reasons; embed, context-size and map tests; scratch install with no drift; three cross-builds, the linux/amd64 binary printing v1.13.1; `report --dry-run`) and PASS again in § Round 2 over the pre-tag fixes (373 pass; mutations a, b, c each compiled and failed the named test; cross-builds repeated). Earlier logs in the range: `testing/2026-09-30-verify-project.md`, `-verify-review-fixes.md` (three rounds), `-verify-redaction-rules.md`, all PASS.
- **Docs freshness** (Step 4): `reviews/2026-09-30-review-v1.13.0-f00d059-docs.md`, APPROVE, no findings — every old name and wording the range removed has no living-doc hit; the delivery claim stays in its one home.
- **Trail freshness** (Step 3): pass — the newest `changes/` entry before the release entry (`2026-09-30-0910-redaction-per-kind-rules.md`) describes the range's last work; no future `date:`; `changes/` gained entries only and `DECISIONS.md` gained ADR-046 and ADR-046a with no past line touched (`git log --diff-filter=M` over the range shows e59638c, whose diff to those paths is 69 insertions and 0 deletions). `ROADMAP.md` and `TODO.md` were not touched since the base: both are empty by design, no plan is active; lag reported, not a violation.
- **Requirements**: R4 (one dependency, unchanged), R8 and R9 (no path changed; ADR-046 accepts a property of an existing path) hold; no `REQUIREMENTS.md` edit.
- **Shadowing agents**: none under `.claude/agents/` or `.agents/agents/`.

## Process notes

- The tree was not frozen for every dispatch this session: `documentation` ran beside the verifier of the redaction task and its deploy-guide edit landed mid-run; the gate's reviewer and security artifacts landed under `.agents/docs` while the gate verifier ran. In both cases only trail files moved, no source, and each verifier named it. Recorded in `changes/2026-09-30-0910-redaction-per-kind-rules.md`; the next session keeps the tree frozen for every agent.
- The round-3 verifier of the review fixes and the round-2 gate verifier each stopped with background work after their log was complete; both were stopped once the artifact was on disk.

## What was not checked

- The Mac and the devbox: installed after the release by the operator's `sh-agent update` there; `darwin/arm64` and `linux/arm64` binaries were built, not run.
- The real-home smoke checks (`testing/2026-09-03-verify-real-home.md`) on either CLI: no installed content changed in this range (no context, skill, agent, command or rule), so no install into a real home preceded the tag; the operator's `sh-agent update` at each place is where they run next.
- The sandbox CLI start tier: same reason.
- A real `sh-agent report` push and a real `update` from this release: the dry run and the fake-release update tests ran; the live paths run at the places after the upload.
- `WriteAtomic` on a network filesystem (reasoned from POSIX; the ordering is pinned).
- Deny-list residue named under the security row: accepted, not checked further.

Every item above was read against the blocking conditions: a verifier FAIL, BLOCKED or IN PROGRESS log (none: both rounds PASS); a high security finding (none); a missing security report, or one naming a release-diff path under Not examined (none: all 42 read); a reviewer REQUEST CHANGES or high finding (none); a `REQUIREMENTS.md` violation (none); a project agent shadowing a review agent (none); the project's own condition, an Antigravity run per changed skill or agent (none changed). No item is one a condition names.

VERDICT: SHIP
