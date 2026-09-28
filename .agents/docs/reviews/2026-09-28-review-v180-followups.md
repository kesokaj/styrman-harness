---
owner: reviewer
date: 2026-09-28
---
# Review — v1.8.0 gate follow-ups   (2026-09-28)

Range reviewed: e4ae34c..worktree (`git diff e4ae34c`, plus the untracked
`.agents/docs/testing/2026-09-28-verify-find-skills-routing.md`).

## Status: COMPLETE

## Verdict: APPROVE

All eight follow-up items from the v1.8.0 gate (preflight rows 3-5; handoff
Next Steps 1c, 1d, 2) are present, each does what it claims, and every claim
I could independently verify checked out. No defect found.

## Verification performed

- `skills/find-skills/SKILL.md` diffed byte-for-byte against
  `repos/vercel-labs/skills` at `435076e` via `gh api … --jq .content | base64
  -d | diff -` — **identical**, confirming the restore is truly verbatim.
- The new `context/AGENTS.md` § Global vs project sentence ("`npx skills add
  <repo> -s <skill>` from the checkout root, or `sh-agent adopt <skill>`
  after a global install"): `adopt` is a real subcommand
  (`cmd/sh-agent/main.go:303`), and `.agents/skills` is a symlink to
  `../skills` (confirmed with `ls -la`), so `npx skills add … -s` does land
  in `skills/` through it, as claimed.
- `.agents/docs/testing/2026-09-28-verify-find-skills-routing.md` (the
  untracked evidence file for that routing change) records a real
  behavioural run (2 of 2 with the new sentence naming the harness checkout
  and its two install commands; 0 of 2 in the control), a valid size check
  (20,822/22,000 bytes), and correctly discarded an invalid first control
  (the real home, which still carried the old local `find-skills` edit).
  Frontmatter and format are correct.
- `internal/harness/zip.go`'s `UnpackZip` cleanup-error join: confirmed the
  named return `err` is set by the `return` statement before the deferred
  closure runs even though callees use a shadowed `err :=` internally — the
  join is correct, and `RemoveAll` is only attempted on the failure path, so
  success is unaffected.
- `internal/harness/repo.go`'s reworded comment ("a zip the operator names
  with `update --from-zip`, which never passes through here"): confirmed by
  tracing `updateFromZip` (`cmd/sh-agent/main.go:1050`) — it builds
  `harness.Source` directly and calls `installFiles`, which never calls
  `ResolveSource`. Only `install` (line 150) and `report` (line 368) call
  `ResolveSource`; the zip path genuinely bypasses it.
- `internal/harness/report.go`'s `sweepInterruptedClones`: glob replaced with
  `os.ReadDir` + prefix/`IsDir`/non-symlink filter. Traced the new test
  `TestSweepInterruptedClonesGlobCharsInHome` (`internal/harness/zip_test.go`)
  by hand: with a home literally named `h[1]`, `filepath.Glob` would
  interpret `[1]` as a character class and match the sibling directory `h1`
  instead — i.e. the old code would sweep *another* home's clone and miss its
  own. The new test asserts exactly that inversion is fixed, and it would
  fail under the old glob-based code, so it is a real regression test, not a
  tautology.
- `cmd/sh-agent/zip_test.go`'s new `TestUpdateFromZipInstallRefusalSweepsNothing`:
  confirmed `sweepStaleZipDirs` in `main.go:1094-1100` sits strictly after a
  successful `installFiles` call and inside `if !dryRun`, matching the
  corrected doc-comment on the neighboring pre-existing test. `main.go` has
  no diff in this range, so this is new coverage of previously-untested
  existing behavior, not a behavior change.
- R9's reworded zip sentence in `REQUIREMENTS.md` ("nothing in the zip
  replaces the binary" vs. the old "the binary itself is not replaced"):
  confirmed `installFiles` → `harness.Plan` does write the running binary
  into the vendor bin dir on every install including the zip path, so the
  old wording was genuinely ambiguous and the new one is accurate. Grepped
  every other living restatement of this claim (`cmd/sh-agent/main.go`,
  `.agents/docs/deploy/remote.md`, `.agents/docs/specs/sh-agent/delivery.md`)
  — all already carry the corrected wording from earlier gate rounds; only
  `REQUIREMENTS.md` itself still had the old sentence, and this diff fixes
  it. The two hits still carrying the old text
  (`.agents/docs/testing/2026-09-28-verify-update-from-zip-e2e.md` and
  `…-phase-2.md`) are dated, append-only historical logs of what a specific
  run printed at the time — not living restatements of the current rule —
  so leaving them is correct.
- `agents/bugs.md`'s example row: now a template placeholder
  (`BUG-001 | YYYY-MM-DD | <area> | …`) instead of a live copy of this
  repo's own BUG-001, matching handoff item 1d.
- `.agents/docs/testing/skill-matrix.md`: `curl` dropped from both the
  Claude `--allowedTools` and the Antigravity `settings.json` allowlist.
  Confirmed the one scenario that mentions curl (`bash` row) only has the
  agent *write* a script containing a curl invocation, checked with
  `bash -n` (syntax only, never executed) — dropping curl from the run
  allowlist does not break it, matching the handoff instruction "drop curl
  if no row needs it." The wording change also now names both credential
  copies (`.credentials.json` and `antigravity-oauth-token`) as deleted by
  name after the batch, per handoff item 2.
- `gofmt -l cmd internal` and `go vet ./...`: both clean on the worktree.

## Findings

None. No correctness, security, or requirement-rule defect found in this
range.

## Criteria

1. Dead code — pass. Nothing added is unreferenced; the `agents/bugs.md`
   example row change replaces content rather than orphaning it.
2. Flow tracing — pass. `zip.go`'s error join, `repo.go`'s bypass claim, and
   `report.go`'s sweep rewrite were each traced caller-to-callee and hold;
   both new tests can fail against the prior code and were confirmed to
   exercise the intended regression.
3. Standards conformance — pass. `gofmt`/`vet` clean; every other living
   restatement of the R9 zip-binary claim was grepped and found already
   consistent, with `REQUIREMENTS.md` being the one this diff fixes.
4. Directory layout — pass. The one new file is
   `.agents/docs/testing/2026-09-28-verify-find-skills-routing.md`, correctly
   under `testing/`.
5. Simplicity — pass. `ReadDir` + prefix match is simpler and safer than the
   glob it replaces; `errors.Join` is stdlib; no new dependency.
6. Decoupling — pass. Changes are confined to the files their own concern
   touches (zip cleanup, report sweep, routing doc, vendored skill); no
   cross-package reach-ins introduced.

## Not examined

- Did not re-run `go test -race ./...` or `govulncheck ./...` myself (ran
  `gofmt -l` and `go vet` only); relying on the project's own verification
  gate for the full suite.
- Did not independently verify the Antigravity behavioural-run leg of
  `2026-09-28-verify-find-skills-routing.md` (only Claude Code's `sbnew`/
  `sbctl` pair is transcript-checkable from here); the log itself notes
  Antigravity was "not exercised" for this question, which is consistent
  with what I read.
- Did not check `.agents/docs/reviews/2026-09-28-preflight-v1.8.0.md` rows
  3-5 or the handoff file's other Next Steps items beyond 1c/1d/2, since
  those are out of this dispatch's scope.
