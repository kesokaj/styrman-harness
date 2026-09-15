---
owner: main
date: 2026-09-15
---
# Review — 641eccd..fe21671 and the working tree, code and docs lenses, merged   (2026-09-15)

Range reviewed: `641eccd..fe21671` (`git diff 641eccd..fe21671 -- context agents skills README.md AGENTS.md .agents/docs/deploy .agents/docs/testing/skill-matrix.md`) plus the uncommitted `/init-project` update path and `init-update` matrix row (`git diff fe21671 -- skills .agents/docs/testing/skill-matrix.md`). Phase 7 of `plans/2026-09-15-harness-practices-g01-g26.md`.

## Verdict per lens

- code — REQUEST CHANGES — `reviews/2026-09-15-review-641eccd-worktree-code.md`; the `/init-project` update path and its matrix check were run in scratch clones and passed and failed as they should.
- docs — REQUEST CHANGES — `reviews/2026-09-15-review-641eccd-worktree-docs.md`; no stale text around the update path.

## Findings

- **P1** [medium] code — `skills/preflight/SKILL.md:65`: a rebutted docs finding is carried forward only from the newest prior gate report's `## Outcome`, and nothing writes it into the next report's Outcome; a same-day re-run also overwrites the report that held it.
- **P2** [medium] code — `skills/preflight/SKILL.md:77`: the exception covers a rebutted finding, not a REQUEST CHANGES verdict that rests on it, and the docs reviewer never sees prior rebuttals.
- **P3** [medium] code, docs — `skills/preflight/SKILL.md:81`: Step 6 still says the SHIP path "has no findings left and needs none", against Step 4's `left` gaps.
- **P4** [medium] code, docs — `internal/harness/apply.go:37` and `cmd/sh-agent/main.go:309`: the installer and `status` still say to add `@AGENTS.md` to a hand-written `CLAUDE.md`, against `README.md:53` and R2.
- **P5** [medium] code, docs — `README.md:32`, `.agents/docs/deploy/local.md:27`: `--vendor all` widens a single-vendor install through `install` only; `update` keeps the recorded vendor (`cmd/sh-agent/main.go:344-350`).
- **P6** [low] docs — `.agents/docs/deploy/local.md:12` still calls both vendors the default.
- **P7** [low] docs — `cmd/sh-agent/main.go:122-126`: the comment still says a plain copy's "image built" its binary.
- **U1** unverified, docs — `skills/init-project/SKILL.md:46`: "Antigravity never reads the content that file holds".
- **U2** unverified, docs — `.agents/docs/deploy/remote.md:33`: argolis-devbox has "no version to bump".

## Round 1

- P1, P2, P3 — `left: operator` — the third round of findings on the release gate's docs routing (Done review rounds 1 and 2, then these), so under Definition of Done step 2 no patch was attempted: BUG-026 is reopened, and the operator is asked what the mechanism is for, with a recommendation to reuse Step 2's coverage rule for Step 4 in place of the rebuttal carry-forward.
- P4 — `fixed at internal/harness/apply.go:37` and `cmd/sh-agent/main.go:309` — both say to move the content into a rule file under `rules/` in the harness repo and let install write the one-line import; `internal/harness/harness_test.go:257` asserts the new guidance.
- P5 — `fixed at README.md:32` and `.agents/docs/deploy/local.md:27` — `--vendor all` widens through `install`; `update` keeps the recorded vendor (`documentation`, checked against `cmd/sh-agent/main.go:137-138` and `:344-350`).
- P6 — `fixed at .agents/docs/deploy/local.md:12` — the vendor recorded by the existing install, else both.
- P7 — `fixed at cmd/sh-agent/main.go:124` — "its image or VM build step built the running binary".
- U1 — `rebutted:` `.agents/docs/architecture/DECISIONS.md:39` records that Claude Code reads `CLAUDE.md` and Antigravity reads `AGENTS.md` or `GEMINI.md`, which is the claim.
- U2 — `rebutted:` `/home/kesokaj/workspace/argolis-devbox/AGENTS.md:160` says the git commit hash is that project's version, and its `deploy/harness/VENDORED.md` refresh ends in a commit and a push with no bump.

After the fixes: `go test -race ./...` exit 0, `go vet ./...` exit 0, `gofmt -l` clean, a scratch install with `status` drift none and a `binary` line with no warning, `render` 93 lines.

## Outcome

- P1, P2, P3 — `fixed` by replacement in e31c3a2: the operator chose option C, so `/preflight` Step 4 no longer carries rebuttals forward or blocks on docs findings; its docs gaps are advisory (`skills/preflight/SKILL.md:65`, `:77`, `:81`; ADR-024 supersedes ADR-023; BUG-026 fixed). The replacement's own docs review is `reviews/2026-09-15-review-c514936-preflight-advisory-docs.md`, APPROVE WITH NITS, with its Outcome.
- P4 — `fixed` in c514936 at `internal/harness/apply.go:37` and `cmd/sh-agent/main.go:309`, asserted at `internal/harness/harness_test.go:257`.
- P5 — `fixed` in c514936 at `README.md:32` and `.agents/docs/deploy/local.md:27`.
- P6 — `fixed` in c514936 at `.agents/docs/deploy/local.md:12`.
- P7 — `fixed` in c514936 at `cmd/sh-agent/main.go:124`.
- U1 — `rebutted`: `.agents/docs/architecture/DECISIONS.md:39` records which file each CLI reads.
- U2 — `rebutted`: `/home/kesokaj/workspace/argolis-devbox/AGENTS.md:160` makes the git commit hash that project's version.

## Round 2

Both lenses re-ran over `641eccd..fcf3b79`, which includes round 1's fixes and the option C replacement (e31c3a2), each given round 1's fixed lines and "Verify each fix is correct and complete; look for defects the fixes introduced."

- code — APPROVE, no findings — `reviews/2026-09-15-review-641eccd-fcf3b79-code.md`; the five fixes re-derived from the diff, the assertion at `internal/harness/harness_test.go:257` fails when the string at `internal/harness/apply.go:37` is reverted.
- docs — APPROVE, no findings — `reviews/2026-09-15-review-641eccd-fcf3b79-docs.md`; `git grep` at fcf3b79 finds no restatement of superseded wording in the living docs.

The loop stops here: round 2 reports no finding beyond the Outcome's `fixed` and `rebutted` lines.
