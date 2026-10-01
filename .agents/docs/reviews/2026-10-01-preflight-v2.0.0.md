---
owner: main
date: 2026-10-01
---
# Preflight — v2.0.0   (2026-10-01)

Range reviewed: v1.16.2..9c6c3ed

Requested in the operator's own turn: "Update the architecture page, readme and all supporting files that this is the new stable state and make sure the REQUIREMENTS is fully updated, then /preflight v2.0.0". Version: MAJOR at the operator's word; the rule's test holds too — a v1.16.x place reacts to this release: two skills it had (`python`, `anti-ui-slop`) are gone, and a flag passed to the wrong `sh-agent` verb now exits 2 instead of being accepted.

## Findings, worst first

No blocking finding. Every medium found by the gate was fixed before the tag; every low is recorded with an owner.

| # | severity / confidence | finding | where | outcome |
|---|---|---|---|---|
| 1 | medium / high | `internal/memory`'s `git ls-files` ran a project's `core.fsmonitor` hook, so the nightly dream executed project code at 03:00 (security code slice) | `internal/memory/memory.go:85` | fixed d26de97 (`storeTrackedArgs`, `TestStoreTrackedArgs`; the agent's hook probe: old binary 1 run, new binary 0) — reviewer APPROVE WITH NITS: the test pins the arguments, a behavioural hook test is `left: next task — main agent` |
| 2 | medium / medium | `allowed-tools: Bash` in the bash and docker skills let every shell command run unprompted on Claude Code while the skill was active (security harness-text slice) | `skills/bash/SKILL.md:4`, `skills/docker/SKILL.md:9` | fixed 89f3308 (lines deleted; no other skill carries the key) — BUG-143 |
| 3 | medium / medium | golangci-lint installed by `curl … \| sh`; gcloud tarball fetched with no checksum step (security harness-text slice) | `skills/golang/references/service.md:404`, `skills/gcloud-setup/SKILL.md:49` | fixed 89f3308 → 19d3467 → 9c6c3ed over three reviewer rounds (download chain wrong on macOS and landing in the project root; then my rewrite deleted the usage sentence): now `go install …@v2.14.0` verified by the checksum database, the release-binary path as one sentence, the gcloud SHA256 through `sha256sum -c -`; reviewer round 3 APPROVE; security Outcome: closed — BUG-143 |
| 4 | medium / high | `deploy/remote.md` said `sh-agent list` shows 27 skills and 7 agents (docs lens) | `.agents/docs/deploy/remote.md:203` | fixed in this release commit (34 and 8) |
| 5 | low / medium | the frozen-tree rule lost its reason and the fix-that-cannot-wait clause; `review` Step 5 omitted `docs/onboarding.md` (docs lens) | `context/AGENTS.md:15`, `skills/review/SKILL.md:105` | fixed 89f3308 (render 15,137 / hand-written 14,252 bytes) |
| 6 | low / high | the dream timer's `systemctl`/`launchctl` calls are not scoped by `--home`: a scratch `uninstall` or `--remove-timer` disables the real timer (security code slice) | `internal/dream/timer.go:112,116,171` | `left: operator` — one-line choice in the handoff (a per-home unit name, or refusing the service-manager calls for a non-real home) |
| 7 | low / high | the dream redacts raw `feedback.jsonl` lines without `report`'s JSON-decoding pass; `memory show` follows a symlinked `MEMORY.md` (security code slice) | `internal/dream/collect.go:65,148`, `internal/memory/memory.go:243` | `left: next task — main agent` |
| 8 | low / low | R3's "no tools" for the Antigravity dream child rests on plan-mode behaviour no recorded probe verified (security code slice) | `internal/dream/dream.go:166` | `left: next task — main agent` (one recorded agy probe settles it) |
| 9 | low / medium | data-gcp's default role sits in `cloudsqlsuperuser` (CREATEROLE/CREATEDB unnamed); docker builder images on mutable minor tags; `feature` reads issue bodies with no data-only clause; an older `npx -y`, a Makefile deploy without `--project`, `commit -am` in a README sample (security harness-text slice) | `skills/data-gcp`, `skills/docker`, `skills/feature`, `skills/find-skills`, README | `left: next task — main agent` (the mutable tag is Phase 6 D26, by research brief, not an operator rebuttal) |
| 10 | low / medium | the kit preamble's real-home guard is a string compare and fingerprints one directory; `command(*)` rows run in a HOME-only sandbox with the operator's uid (security trail slice) | `testing/kit-preamble.sh:24,47`, `testing/skill-matrix.md:14` | `left: operator` — accept as the sandbox's design (recommended) or ask for real isolation; recorded in the handoff |

## Rows of the Definition of Done

- **Dirty source before the gate (Step 1):** none at the start besides the README repair, committed as 0b1b5b8 with its entry. The gate then fixed its own findings in four commits — 89f3308, d26de97, 19d3467, 9c6c3ed — each with a `changes/` entry, a reviewer round and a verifier re-run; the verifier that was in flight when d26de97 and 9c6c3ed landed was stopped and re-dispatched (BUG-142's rule), never cited.
- **Verification:** `testing/2026-10-01-verify-v2.0.0-gate.md` PASS at 0b1b5b8; `-r2` PASS at 89f3308; `-r3` stopped (tree moved, not cited); `-r4` PASS at 19d3467; `-r5` PASS at 9c6c3ed — gofmt clean, vet clean, six packages ok with `-race`, 342 top-level tests, both size pins (14,252 / 15,137), `TestEmbedHoldsOnlyTrackedFiles`, scratch install 34 skills / 8 agents / drift none, govulncheck 0 affecting, the coverage map of every changed Go file (r1). Status at every cited log: PASS.
- **Reviewer:** the nine phase commits are covered by their phase reports, each `## Status: COMPLETE` with a `Range reviewed:` line over the working tree that became the commit (`reviews/2026-10-01-review-worktree-p5..p9-{code,docs,prompts}.md`, `review-4764c5f-worktree-code*.md`, `review-c62e064-*`, `review-329a99b-*`, `review-1430d53-*`); the parity commits of 2026-09-30 by `review-cb799ff-f1ca525-*`, `review-f1ca525-e2cffc1-code.md`, `review-e2cffc1-a275b96-code.md` (the `review-a275b96-head-code.md` left IN PROGRESS is not cited: that commit is covered by the v1.17.0 gate re-verify and the Phase 1 code review); ddac2d8 (3 files) and 0b1b5b8 (docs) read by the main agent, below the threshold; the gate-fix commits by `reviews/2026-10-01-review-0b1b5b8-89f3308-code.md` (three rounds: REQUEST CHANGES → REQUEST CHANGES → APPROVE, Outcomes per round) and `reviews/2026-10-01-review-89f3308-d26de97-code.md` (APPROVE WITH NITS). The Phase 6 code review hit its three-round cap and was decided by the operator (D40).
- **Security:** three slices over the whole range, each `## Status: COMPLETE`, each `## Not examined` naming no release-diff path after the follow-up reads: `reviews/2026-10-01-security-v1.16.2-0b1b5b8-{code,harness-text,trail}.md`, with `## Outcome` sections for the fix commits (the code slice re-probed the fsmonitor hook against the new binary; the harness-text slice ran the new install path and the gcloud check). No high finding.
- **Docs lens:** `reviews/2026-10-01-review-v1.16.2-0b1b5b8-docs.md`, APPROVE WITH NITS; its three findings are rows 4 and 5 above (fixed) — the "GitHub Container Registry only on explicit request" clause it reported dropped stands in `skills/cloud-run-deploy/SKILL.md:19`.
- **Documentation agent:** not dispatched — the docs changes of this gate are the docs-lens findings and the two release-state lines (`specs/sh-agent/delivery.md` § Where things stand → v2.0.0; `deploy/remote.md` count), applied by the main agent and read against the diff.
- **Bugs:** BUG-143 (the two security mediums in skills; fixed). BUG-142 (open, this plan) applied twice in this gate as the rule it names.
- **Changelog:** six entries by the main agent (`changes/2026-10-01-1229-…`, `-1244-…`, `-1247-…`, `-1249-…`, `-1251-…`, and the release entry).

## Trail freshness (Step 3)

- Newest `changes/` entry describes the release: the release entry written in Step 6; before it, the gate-fix entries describe the gate's own commits — pass.
- No `date:` in the future under `.agents/docs` — pass.
- `changes/` append-only: no past entry modified in the range — pass. `architecture/DECISIONS.md` shows as modified in four commits (d6e396a, cb799ff, 1430d53, ddac2d8): each appends an ADR (047, 047a, 049, 048) — pass.
- `ROADMAP.md` and `TODO.md` touched in five commits of the range — report-only, no lag.

## The project's own release rule (`AGENTS.md` § Release)

- Real-home smoke on both CLIs: `testing/2026-10-01-verify-real-home-v2.0.0.md` 6 of 6 at 0b1b5b8; and 6 of 6 again at 9c6c3ed after the gate fixes (same log, § The final head).
- One Antigravity run of each changed skill and agent: Phase 9's 24 rows (`testing/2026-10-01-verify-skill-matrix-antigravity.md`) plus the 14 rows this gate still owed — `bug`, `review`, `feature`, `research`, `retro`, `feedback-log`, `preflight`, `dream`, `project-memory` and the agents `bugs`, `changelog`, `documentation`, `dreamer`, `security` — 14 of 14 PASS (`testing/2026-10-01-verify-skill-matrix-antigravity-v2.0.0.md`); the four skills the gate fixes changed re-run from 9c6c3ed: 4 of 4 PASS (same log, § Re-run).
- Tag: signed per `deploy/remote.md` § Releases from this checkout; the signing key's fingerprint check is the first command of Step 6.

## What was not checked

- The Mac (no access this session): neither the smoke nor a build there.
- Antigravity rows unchanged since their 2026-09-29/30 runs (`init-memory*`, `review-prompts`, `review-project`, the shadow rows): not re-run.
- `go test -cover` percentages; govulncheck's advisory ID in summary mode (0 affecting either way).
- Offline: the golangci-lint v2.14.0 tag's existence and the vendor's exact go-install caveat wording; the gcloud install page's SHA256 for the unversioned tarball.
- A behavioural test for the fsmonitor fix (row 1's nit): the fix is evidenced by the security agent's probe, not by `go test`.
- Row 8: the Antigravity plan-mode assumption behind R3's "no tools" — unverified in this range, not blocking on the slice's evidence.

None of these is a condition the blocking list names: no verifier log cited is FAIL, BLOCKED or IN PROGRESS; no high security finding; every security report is COMPLETE with a `## Not examined` section naming no release-diff path; no REQUEST CHANGES stands unresolved (the recipe review's round 3 is APPROVE); no `REQUIREMENTS.md` violation (`allowed-tools` was a tension with § Non-goals, not an edit to a vendor setting; the security agent agreed); no project agent shadows a harness one (`.agents/agents` and `.claude/agents` absent).

VERDICT: SHIP
