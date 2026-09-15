---
owner: reviewer
date: 2026-09-15
---
# Review — 641eccd-fcf3b79-docs   (2026-09-15)

Range reviewed: 641eccd..fcf3b79 (`git diff 641eccd...fcf3b79`)

## Verdict: APPROVE

All five round-1 fixes are correct, complete, and match the code they describe; `git grep` over the living docs at fcf3b79 finds no restatement of any superseded wording. No new defect was introduced by any of the five fixes. This range also carries a sixth, unlisted fix (preflight's docs routing rewritten a second time to the "advisory" rule, ADR-024) whose own docs-lens review already ran, found four low findings, and recorded an Outcome fixing three and rebutting one — that loop is closed and consistent with the current `skills/preflight/SKILL.md`.

## Findings

None. Each of the five fixes was checked against the code it describes:

- `skills/preflight/SKILL.md:65,77,81` — Step 4's docs review reads as advisory (owner, severity, confidence, filed under Follow-ups/Next Steps); Step 5's blocking list has no docs clause ("Step 4's docs gaps never block"); Step 6's "needs none" sentence is gone. Matches ADR-024 and the D49/D46 plan rows.
- `internal/harness/apply.go:37`, `cmd/sh-agent/main.go:309` — both now say to move a hand-written `CLAUDE.md`'s content into a rule file under `rules/` and let install write the one-line import; `internal/harness/harness_test.go:257` asserts the new wording and the full test suite passes (`go test -race ./...` exit 0).
- `README.md:32`, `.agents/docs/deploy/local.md:27` — `--vendor all` widens a single-vendor install through `install`; `update` keeps the recorded vendor. Matches `cmd/sh-agent/main.go`'s usage line (`default: from manifest if installed, else all`) and the `Install`/`Update` vendor-handling code.
- `.agents/docs/deploy/local.md:12` — "the vendor recorded by the existing install, else both" — consistent with the same default logic.
- `cmd/sh-agent/main.go:124` — comment now says "its image or VM build step built the running binary", matching `AGENTS.md`'s own updated line and `deploy/remote.md`'s new description of the two mirrors' different build steps (styrman-crew images it, argolis-devbox's VM builds it locally via a setup script).

Also spot-checked and consistent: `.agents/docs/deploy/remote.md`'s "the working directory or any of its parents" claim matches `internal/harness/repo.go`'s `FindRepo` walk; `skills/VENDORED.md`'s "LICENSE and NOTICE" (MANIFEST dropped) matches `skills/anti-ui-slop/`'s actual file list; the new `init-update` matrix row and `skills/init-project/SKILL.md`'s update-mode rewrite agree with each other and with plan D47/BUG-025's resolution; BUG-031 (open, unfixed) makes no false claim anywhere in the living docs — no doc asserts headless runs wait for dispatched agents.

## Not examined

- No end-to-end run of `/preflight` or `/init-project` in this session — behaviour was judged from text and existing test/tool output only (same limit the round-1 reviewer already recorded as Unverified).
- `.agents/docs/testing/skill-matrix.md` was read only for the `init-update` row's addition, not diffed line-by-line against every other row.
- `.agents/docs/reviews/2026-09-15-review-c514936-preflight-docs.md` (the round-4 REQUEST CHANGES report) was not opened directly; its content was cross-checked via BUG-026's summary and the later `-advisory-docs.md` report's Outcome, both of which are internally consistent with the current `skills/preflight/SKILL.md`.
- Downstream-mirror refresh work (task 7.6, styrman-crew/argolis-devbox) is outside this repo and was not checked.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens
