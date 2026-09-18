---
owner: main
date: 2026-09-18
---
# Preflight — v1.3.0     (2026-09-18)

Range reviewed: v1.2.3..10a9fe2

Base: the last tag, `v1.2.3` (release commit 91fcf13). Head: 10a9fe2, `main`, clean tree. The range is the retro application (1358bcc: four context sentences, two `feedback-log` triggers, the verifier's early-log rule, the feature skill's cold-read clause; BUG-078–083 records) plus three pre-gate commits made inside this gate (94bec82, 8e1b044, 10a9fe2) that carried its review findings. No Go source changed since v1.2.3: the binary is the v1.2.3 code with the new harness embedded. Version: MINOR — an installed place gains rules, triggers and an agent behaviour (`remote.md` § Releases, version rule); the operator's "Cut v1.3.0".

## VERDICT: SHIP

Nothing blocks. The gate took four review rounds; the operator was asked at the second-round cascade ("fix it") and at the three-round cap ("1"), and each non-blocking finding below carries its decision.

## Findings, worst first

None high. None medium open.

- [low / medium] **`skills/preflight/SKILL.md:77` — the new "verifier log still IN PROGRESS" blocking clause is unqualified**, while `:52` lets a cut-off log be replaced by a fresh dispatch; a gate that correctly re-ran the verifier and got a PASS still finds a blocking condition on disk (a false NO-SHIP, never a false ship). (reviewer round 4.) **Decision: accepted for this tag, handoff Next Steps** — qualify it to "the log cited for the Verification row, or any FAIL or BLOCKED log in the range"; a one-line edit that would re-open a skill four rounds deep, in the safe direction.
- [medium → fixed] rounds 1–3: the `feedback-log` schema line without `redirection`; the verifier's template Status without `IN PROGRESS`; "after the FIRST command" unsatisfiable before Pre-flight; preflight's blocking list without `IN PROGRESS`; project `AGENTS.md:63`, `remote.md:8/91/195`, `delivery.md:7/56` still stating the ADR-030 delivery model; preflight re-dispatching over a FAIL log. **All fixed in the three pre-gate commits**, verified by the reviewer's next round each time; recorded as BUG-084 (vocabulary added in prose, not in the file's own enumerations) — `reviews/2026-09-18-review-v1.2.3-3b5b71a-code.md` § Outcome.
- [low / high] **BUG-083, open** — the verifier's rewritten early-log rule held in 1 of 3 scratch runs with the ORIGINAL wording and 1 of 1 with the fixed wording ("the FIRST verification command"); the dispatch sentence (`context/AGENTS.md` § Work) remains the enforcement. **Decision: on record**, not a release defect; more runs are the next session's.
- [note] the verifier's round-1 observation that `go run … install --dry-run` reports the binary as "1 written" (an ephemeral build ID) is a testing-recipe gap, not a code defect: the guide's built-binary form gives 0 written, as rounds 2–4 show.

## Rows

| Row | Result | Artifact |
|---|---|---|
| Verification | PASS — 250/0/0 with `-race`, `govulncheck` none on go1.27.1, three cross-builds, install 260 / drift none / dry-run 0 written from a built binary, list 27 skills 7 agents; rounds 2–4 re-run at each pre-gate head, `Range reviewed: v1.2.3..10a9fe2` | `testing/2026-09-18-verify-preflight-v1.3.0.md` |
| Review (opus), four rounds | round 1 REQUEST CHANGES (3) → round 2 REQUEST CHANGES (2, cascade) → round 3 REQUEST CHANGES (3 + low, the sweep) → round 4 APPROVE WITH NITS (1 low, above); `## Outcome` per finding | `reviews/2026-09-18-review-v1.2.3-3b5b71a-code.md` |
| Security | no row: no changed path in the range touches auth, crypto, secrets, IAM, input parsing or dependencies (no Go, no `go.mod`); the trail markdown added was read by the reviewer and the docs lens; the staged-diff secret scan ran before each commit — its two hits are the `-----BEGIN SSH SIGNATURE-----` armor header the project `AGENTS.md` names as a tag check | — |
| Docs lens | APPROVE WITH NITS, no findings in range; the worktree docs review's two fixes hold; the README tag sentence it flagged was fixed in 94bec82 | `reviews/2026-09-18-review-v1.2.3-3b5b71a-docs.md` |
| Trail freshness | pass: newest `changes/` entry `2026-09-18-retro-applied.md` is this range's work; no future `date:`; no modified past `changes/` or `DECISIONS.md` entry since v1.2.3; `ROADMAP.md` touched at the plan close, `TODO.md` empty | — |
| Shadowing project agent | none (`.claude/agents`, `.agents/agents` absent) | — |
| `REQUIREMENTS.md` | untouched since v1.2.3; the delivery sentences fixed in this gate now match R8/R9 | — |
| Signing key | present, fingerprint `SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU` from the private key | — |
| Embed guard | nothing untracked or ignored under the five embedded directories; `TestEmbedHoldsOnlyTrackedFiles` ok, re-run at 10a9fe2 | — |
| Repository visibility | `PRIVATE` | — |
| Behavioural | `claude -p --append-system-prompt` quoted the new Hard-rules sentence with BUG-078/079; installed-context smoke checks from this box's real home at d278f5a on both CLIs (`sh-agent-matrix.md` § Real machines); the fixed verifier definition via `claude -p --agents` in a scratch clone, 1 of 1 writing IN PROGRESS after the first verification command | `changes/2026-09-18-retro-applied.md`, `testing/sh-agent-matrix.md` |

## What was not checked

- No CLI was started from a scratch `HOME` holding 10a9fe2's files: neither CLI authenticates outside the real home; the real home is at d278f5a (four text lines behind), and updating it is the operator's call.
- Antigravity's skill-load and agent-dispatch legs were not re-run on this box for this range (its context question passed at d278f5a).
- No Mac or arm64 machine ran anything in this range; the README browser chain, Gatekeeper and `shasum --ignore-missing` remain the standing remainder.
- The BUG-081–084 records and the three verifier scratch logs were read by their authors and the main agent, not by a reviewer.
- No Dockerfile in the range; no container tier.
