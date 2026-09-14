---
owner: reviewer
date: 2026-09-14
---
# Review — bug-015-docs (2026-09-14)

Range reviewed: uncommitted working tree against HEAD `e0f6830` (`/review --docs` pass; BUG-015, the `sh-agent install` plain-copy behavior change). Diffed files: `.agents/docs/deploy/remote.md`, `AGENTS.md`, `README.md`, `cmd/sh-agent/main.go`, `cmd/sh-agent/main_test.go`, `internal/harness/repo.go`, `.agents/docs/testing/sh-agent/guide.md`, `.agents/docs/bugs/INDEX.md` (append-only, out of scope), plus untracked `.agents/docs/bugs/BUG-015-install-rebuild-plain-copy.md` and `.agents/docs/testing/2026-09-14-verify-bug-015.md` (out of scope per dispatch — bugs/ and dated testing logs are historical records).

Note: `git diff HEAD --stat` truncated silently on first invocation (omitted `.agents/docs/bugs/INDEX.md` and `.agents/docs/testing/sh-agent/guide.md`, both present under `git status --short`); re-run confirmed the full 8-file diff. Worth knowing if this recurs elsewhere.

## Verdict: APPROVE WITH NITS

One file that describes this exact command still states the pre-change behavior with no mention of the plain-copy case: `.agents/docs/deploy/local.md:27`. Everything else that describes `sh-agent install` — README's command table, AGENTS.md's "Rules for working here", `.agents/docs/deploy/remote.md`, and the testing guide — was already updated correctly and consistently, including a new section (3a) documenting the exact new stdout line alongside the other install-time messages.

## Findings

- **[medium / confidence: high]** `.agents/docs/deploy/local.md:27` — "Every `install` rebuilds the binary with `go build` first; only without `go` on PATH is the running binary copied, with a warning." This is the pre-BUG-015 sentence: it omits the plain-copy case entirely (no `.git`, e.g. a downstream mirror), where install now also skips the rebuild and prints no toolchain warning. An operator following this file for a mirror install would expect a rebuild attempt or a "go is not on PATH" warning; neither happens. Replacement, mirroring the already-fixed `AGENTS.md:35` sentence:
  "From its own git checkout, `install` rebuilds the binary with `go build` first; the running binary is installed as is when `go` is missing (with a warning) or when the tree is a plain copy with no `.git` — a downstream mirror (BUG-015)."

## New stdout line documented?

Yes. `plain copy, not its own git checkout: the running binary is installed as is` (printed by `run()` in `cmd/sh-agent/main.go` when `!own`, replacing the toolchain warning for that branch) is documented in `.agents/docs/testing/sh-agent/guide.md`, new section "3a. Plain-copy install" (lines 41–65), in the same place the sibling `[dry-run] the binary would be rebuilt from …` message is documented (line 65) — consistent placement, both cited as literal expected substrings ("must say \"plain copy\"", "first line \"plain copy\""). README's install-table row (line 24) and AGENTS.md's rule (line 35) describe the behavior in prose without quoting the literal string, which is consistent with how the pre-existing go-missing warning is also described there (prose, not verbatim).

## Files checked and their status

- `README.md` (command table row 24, quick start, line 101) — accurate, already updated.
- `AGENTS.md` (Rules for working here, line 35; Testing section) — accurate, already updated.
- `.agents/docs/deploy/remote.md` (Downstream mirrors section) — accurate, already updated, cites BUG-015.
- `.agents/docs/deploy/local.md` — stale, see finding above.
- `.agents/docs/testing/sh-agent/guide.md` — accurate; new section 3a added and consistent with code.
- `context/AGENTS.md`, `context/dev-setup.md` — no mention of install/rebuild mechanics at all (out of scope for this behavior); no stale claims.
- `skills/*/SKILL.md` (all matches: cloud-run-deploy, docker, golang refs, terraform refs, init-project asset) — all matches are unrelated generic Go/Docker/Terraform "rebuild"/"copied" usage, none about `sh-agent install`. No findings.
- `cmd/sh-agent/main.go` usage/help constant and live `--help` output — generic ("deploy context, skills, and agents…"), never described rebuild mechanics before or after; not stale.
- Grep of the whole repo (outside changes/, architecture/DECISIONS.md, reviews/, handoffs/, retros/, bugs/) for rebuild/rebuilt/copied/running binary/plain copy/"go is not on PATH"/always rebuilds turned up only the files listed above plus historical `.agents/docs/plans/*.md` and `.agents/docs/testing/2026-09-03-verify-rebuild.md` and similar dated verification logs — treated as historical snapshots of an in-progress implementation/verification, not current-behavior documentation, and not corrected (see Not examined).

## Criteria

Not applicable in the usual code sense — this is a docs-only pass — but for completeness against the six:
1. Dead code — n/a (no source reviewed for removal).
2. Flow tracing — the one finding is exactly a flow-tracing gap: a description file that an installer's operator reads was not updated when the callee's behavior (the rebuild decision) changed.
3. Standards conformance — n/a.
4. Directory layout — n/a; report written to `.agents/docs/reviews/`.
5. Simplicity — n/a.
6. Decoupling — n/a.

## Not examined

- `.agents/docs/plans/2026-09-03-rebuild.md`, `.agents/docs/plans/2026-09-13-learning-loop-project-level.md`, `.agents/docs/plans/2026-09-14-cold-read-session-close.md` — flagged by grep but treated as historical planning records (like `changes/`) rather than current-behavior documentation; not corrected. If any of these is still treated as a live spec rather than a record, it needs its own pass.
- `.agents/docs/testing/2026-09-03-verify-rebuild.md` and other dated `.agents/docs/testing/*.md` verification logs (2026-09-13 items) — same treatment: dated logs of a specific run, not living docs; not corrected.
- `.agents/docs/bugs/BUG-015-install-rebuild-plain-copy.md`, `.agents/docs/bugs/INDEX.md`, `.agents/docs/testing/2026-09-14-verify-bug-015.md` — out of scope per dispatch (bugs/ excluded; the verify-bug-015 log is a dated verifier record, not a living description).
- No code review performed (not this pass's job): `cmd/sh-agent/main.go` and `internal/harness/repo.go` diffs were read only to establish ground truth for the docs check, not reviewed for correctness.
- Did not verify `go run ./cmd/sh-agent install` end-to-end against a real plain-copy mirror in this session (the testing guide's section 3a and the untracked verify-bug-015 log claim this was done separately by the verifier).

## Outcome

- Medium: `.agents/docs/deploy/local.md:27` omitted the plain-copy case — **fixed** with the replacement sentence above, in the same change.
- Nits: none required a change; the new stdout line stays documented in `.agents/docs/testing/sh-agent/guide.md` § 3a next to its dry-run sibling.
