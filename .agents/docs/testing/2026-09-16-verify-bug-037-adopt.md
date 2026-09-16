---
owner: verifier
date: 2026-09-16
---
# Verification — bug-037-adopt (re-run, two added tests)   (2026-09-16)

Range reviewed: HEAD (d49320a) with uncommitted changes to
`internal/harness/adopt.go` and `internal/harness/harness_test.go` (BUG-037 fix).
Guide cited: `.agents/docs/testing/sh-agent/guide.md`.

## Status: PASS

Scope: guide §1-§3 run in full against the working tree; §3a judged not
applicable; the e2e/real-home matrices not re-run (cited, not exercised).

## Commands

- `gofmt -l .` → prints nothing on tracked source (one hit on gitignored
  scratch `.agents/scratch/mustfail3/internal/harness/harness_test.go`,
  confirmed via `git check-ignore -v`, not part of this change)
- `go vet ./...` → exit 0
- `go test -race -v ./...` → exit 0; output saved to
  `.agents/scratch/verify-bug-037-2/go-test-1.txt`
  counts: pass 38, fail 0, skip 1 at d49320a+dirty for `go test -race -v ./...`
  vs previous: `.agents/docs/testing/2026-09-16-verify-bug-037-adopt.md`
  (this file's own prior run), pass 36→38, skip 1→1 — rise of exactly 2,
  matching the two named added tests (`TestAdoptRefusesDestructiveSources`'s
  fourth subtest and `TestWithin`), both confirmed present and PASS in the
  new log (`TestAdoptRefusesDestructiveSources/a_relative_path_into_the_repo's_own_skills/`,
  `TestWithin`).
  vs `.agents/docs/testing/2026-09-15-verify-harness-practices-done.md`,
  pass 32→38 at 8202d6a→d49320a+dirty (six more, all from this bug's fix
  and its two rounds of added tests plus the prior round's), skip 1→1.
- `go run ./cmd/sh-agent list` → exit 0, all skills/agents parse, section
  headers print counts (27 skills, 7 agents) — no drift from expected shape.
- `go run ./cmd/sh-agent render | wc -l` → 93 (matches expectation).
- Scratch home install: `go run ./cmd/sh-agent install --home
  .agents/scratch/verify-bug-037-home-2` → exit 0, "260 written, 0 unchanged,
  0 removed, 0 kept".
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-bug-037-home-2`
  → exit 0, `drift none — 259 files match the repo`, `binary
  d49320ab3b7e+dirty` line present with no warning.
  No agent definition changed in this diff, so the `## Records`
  comparison across `agents/*.md` and derived Antigravity `agent.md` is not
  required this run.

## §3a — plain-copy install tier

Not applicable: this change touches only `internal/harness/adopt.go`
(the `Adopt`/`within` logic) and its test file; it does not touch
`IsCheckout` or any rebuild-decision logic in `cmd/sh-agent` or
`internal/harness/repo.go`, so the read-only mirror tier does not apply.

## Coverage of changed code

- `Adopt` (internal/harness/adopt.go) — covered by
  `TestAdoptRefusesDestructiveSources` (4 subtests: `.`, `..`, a skill
  adopted onto itself, and the newly added "a relative path into the
  repo's own skills/" subtest exercising `Adopt(repo, home, "skills/alpha",
  true)` from the repo root).
- `within` (internal/harness/adopt.go, new helper) — covered by
  `TestWithin`, a table test over: exact equality, parent-child both
  directions, prefix-sharing siblings (`/a/skills` vs `/a/skills-old` and
  `/a/skillsdelta`), a trailing separator, and an unrelated path.

No added `//nolint`, `t.Skip`, `.only(`, or similar suppression appears in
this diff.

### Independent reproduction (this run, not the main agent's claim)

Reproduced in `.agents/scratch/` (removed after use, not committed):
took `git archive HEAD`, applied this diff with `patch -p1`, confirmed
`TestWithin` passes as shipped, then replaced `within`'s body with
`strings.HasPrefix(child, parent)` in that scratch copy only. Result:
`TestWithin` FAILs exactly as claimed —
`within("/a/skills", "/a/skills-old") = true, want false` and
`within("/a/skills", "/a/skillsdelta") = true, want false` — confirming
`TestWithin` catches the prefix-string regression on siblings that share a
prefix.

Did not independently reproduce the second claim (the new
`TestAdoptRefusesDestructiveSources` subtest failing against
`git show HEAD:internal/harness/adopt.go`): that pre-fix file has no
`within` function, so applying only the test-file diff on top of the old
`adopt.go` fails at compile time (`undefined: within`) rather than at
`go test` runtime, and reproducing the pre-fix runtime behavior would
require hand-reconstructing the old dest-computation path outside the
diff, which risks testing a fabrication rather than the actual regression.
This is recorded as not reproduced, per the instruction to reproduce
"either", not both.

## Root cause

N/A — Status is PASS, no FAIL to explain.

## Not examined

- `.agents/docs/testing/2026-09-15-verify-harness-practices-e2e.md` and
  `.agents/docs/testing/2026-09-15-verify-real-home.md` — cited, not
  re-run; neither CLI was started and the skill/agent matrix was not
  driven this run, per dispatch.
- The second reproduction claim (subtest failing against the pre-fix
  `adopt.go`) was not independently reproduced; see Coverage section above
  for why.
- No files beyond the named diff (`internal/harness/adopt.go`,
  `internal/harness/harness_test.go`) were reviewed for this bug; the
  three untracked docs files present in `git status` (plan, two review
  reports) were not read or evaluated — out of scope for this verify pass.

## Environment

Linux 6.12.107+deb13-amd64, `go` from PATH (module
`github.com/kesokaj/styrman-harness`), git repo at
`/home/kesokaj/workspace/styrman-harness`, branch `main`, HEAD `d49320a`
with the named uncommitted diff. No Docker/Colima involved (no Dockerfile
in this change).
