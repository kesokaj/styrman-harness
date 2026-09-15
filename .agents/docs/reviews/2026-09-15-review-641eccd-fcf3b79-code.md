---
owner: reviewer
date: 2026-09-15
---
# Review — 641eccd-fcf3b79-code   (2026-09-15)

Range reviewed: 641eccd..fcf3b79 (`git diff 641eccd...fcf3b79`); code lens (a parallel docs-lens report, `2026-09-15-review-641eccd-fcf3b79-docs.md`, already covers the Markdown living-docs surface).

## Verdict: APPROVE

Nothing is broken. This is the round-1 fix commit for the two-lens Phase 7 review (`2026-09-15-review-641eccd-worktree-merged.md`) plus the Phase 7 gate-advisory and init-update work. All five round-1 code/docs fixes (P4–P7, P6) are present, correct, and match the code they describe; the new test can fail; `go test -race ./...`, `go vet ./...` and `gofmt -l` are clean. The three Go files touched (`apply.go`, `main.go`, `harness_test.go`) carry only string/comment changes — no logic changed, so there is no new flow-tracing surface to break. I re-derived the round-1 fixes independently from the diff rather than trusting the prior report, and confirmed them:

- `internal/harness/apply.go:37` and `cmd/sh-agent/main.go:309` — both now tell the operator to move a hand-written `CLAUDE.md`'s content into a rule file under `rules/` and let `install` write the one-line import; matches `README.md:53` and R2 (`REQUIREMENTS.md`).
- `internal/harness/harness_test.go:257` — new assertion `strings.Contains(out.String(), "move its content into a rule file under rules/")`. This is a real, failable assertion: reverting the `apply.go:37` string change makes it fail (the old string was "add the line `@AGENTS.md` to it yourself"), and it does not duplicate the pre-existing `keep    ~/.claude/CLAUDE.md` check above it.
- `skills/preflight/SKILL.md:65,77,81` — Step 4 is advisory (owner/severity/confidence, filed under Follow-ups/Next Steps), Step 5's blocking list has no docs clause ("Step 4's docs gaps never block"), Step 6's "needs none" sentence is gone. Cross-checked against ADR-024 in `DECISIONS.md` and plan rows D46/D49 — consistent, and BUG-026 is recorded `fixed` in `INDEX.md` and its own record.
- `README.md:32`, `.agents/docs/deploy/local.md:12,27` — the `--vendor` default/widening text matches `resolveVendor` (`cmd/sh-agent/main.go:216`) and the `update` command's vendor handling (`:344-350`): `update` passes `--vendor` only when `vendor != VendorAll`, so an install recorded as a single vendor stays that vendor across `update`, and `--vendor all` on `install` only is the widening path the docs describe.
- `cmd/sh-agent/main.go:124` / `AGENTS.md:35` — comment says "image or VM build step"; matches `deploy/remote.md`'s new description of the two mirrors' different build steps.

I also re-verified several evidence claims recorded in this range's own review/changes artifacts rather than taking them on faith (the `grep -cF` counts in plan task 7.7 and the changes entry for the gate-docs-advisory work); all reproduce exactly as claimed.

No REQUIREMENTS.md violation, no recorded bug (BUG-INDEX) reintroduced, no hand-rolled auth/crypto, no new dependency.

## Findings

None that meet the correctness/security/rule bar. Two low-confidence observations, not findings (no concrete failure mode, noted for completeness):

- `skills/init-project/SKILL.md` Step 4's new per-file copy loop (`find ... | while read -r f; ... [ -e "$d" ] && continue ...`) replaces the old `cp -R`; I traced its `.gitkeep`-only-into-empty-directory branch and the `[ -e "$d" ]` skip-if-exists branch by hand and both are correct for the cases the accompanying changes entry claims were tested (fresh clone: 18 files + 4 links added; second run: nothing added). This is prose in a skill file, not executed by CI, so I did not run it myself; the changes entry (`…phase7-init-update.md`) says the code reviewer ran it in scratch clones with matching results, which I did not re-verify independently. Low confidence, not a finding.
- The `.gitignore`-in-gofmt-list noise (`gofmt -l .` also lists files under `.agents/scratch/adv-init/...` and `.agents/scratch/lintcheck/main.go`) is pre-existing scratch content, gitignored (`.gitignore:21`), untouched by this diff, and not part of the repo's tracked Go tree. Not a finding.

## Criteria

1. Dead code — pass. No unused exports, no commented-out code, no leftover debug output introduced by this range (checked with a diff grep for `console.log`/`fmt.Println(`/`debug`/`TODO`/`FIXME` additions outside test files — none found). The old "add the line `@AGENTS.md`" strings are fully replaced, not left as dead alternates.
2. Flow tracing — pass. The three Go edits are string/comment-only; no caller/callee shape changed. `resolveVendor`/`update`'s vendor-handling logic (unchanged in this range) matches what the new docs claim. The new test asserts a real, failable behaviour (see above).
3. Standards conformance — pass. `gofmt -l` clean on tracked files, `go vet ./...` clean, `go test -race ./...` exit 0. No golang-skill rule engaged beyond formatting (no error paths, contexts, or resource lifetimes touched).
4. Directory layout — pass. No new files land outside the established tree (`internal/harness/`, `cmd/sh-agent/`, `skills/`, `.agents/docs/{bugs,changes,plans,reviews,architecture}/`); reviews and bug records are named and placed per convention.
5. Simplicity (DIY-except-security) — pass. No new dependency, no speculative abstraction; the changes are string edits and prose. N/A for the security exemption (no auth/crypto/session code touched).
6. Decoupling — pass. `apply.go` and `main.go`'s edits are independent one-line string changes in code that already existed at those seams; no cross-package reach-in introduced.

## Not examined

- No `/preflight` or `/init-project` run in the sandbox or on either CLI this session; all judgments on prose behaviour (init-project's update-mode shell snippets, the matrix rows) are from reading, matching the "unverified"/"not examined" limits already recorded in the round-1 and docs-lens reports for this same range.
- Task 6.4 (the Phase 6 matrix comparison / BUG-029/030/031 investigation) and task 7.6 (downstream-mirror refresh in styrman-crew and argolis-devbox) are outside this repo's tree and were not checked.
- I did not independently re-run the `init-update` shell snippets from `skills/init-project/SKILL.md` Step 4/5 in a scratch clone; I traced the logic by hand only (see Findings).
- `.agents/docs/reviews/2026-09-15-review-c514936-preflight-docs.md` (the round-4 REQUEST CHANGES report referenced by BUG-026) was read only via BUG-026's summary and the later advisory-docs report's Outcome, not opened directly.
- Full line-by-line diff of `.agents/docs/testing/skill-matrix.md`'s other rows against the new `init-update` row was not performed beyond the row itself.
