---
owner: verifier
date: 2026-09-18
---
# Verification — sh-agent-works-everywhere-phase-2   (2026-09-18)

Range reviewed: 160aa16..3209cf3

## Status: PASS

Scope: complete — every command in the dispatch ran, plus the three live checks (a/b/c), against the guide (`.agents/docs/testing/guide.md`). Go sources were confirmed frozen (`git status --short` clean before and only my own docs artifact plus another agent's docs review appeared as untracked after; no source diff).

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -count=1 -v ./...` → exit 0; full output in `.agents/scratch/verify-phase-2-logs/go-test-all.txt`; `ok` for all three packages (root, `cmd/sh-agent`, `internal/harness`).
  - counts: pass 244, fail 0, skip 0 at 3209cf3 for `go test -race -count=1 -v ./...` (counted via `grep -cE -- '--- (PASS|FAIL|SKIP):'`, which matches indented subtest lines too — an initial `^--- PASS` anchor undercounted at 129 because subtest lines are indented; corrected).
  - vs previous: `.agents/docs/testing/2026-09-18-verify-sh-agent-works-everywhere-phase-1-2.md`, pass 227→244, skip 0→0. Matches the dispatch's expectation (227 + 9 new in `internal/harness/audit_test.go` + 8 new in `cmd/sh-agent/main_test.go` = 244) exactly; no unexplained difference.
- Targeted `internal/harness` tests (9 named): `TestIsReleaseVersionIsCanonical`, `TestNonCanonicalTagNeverReachesTheDownload`, `TestNoteRunningElsewhereOnlyForTheRealHome`, `TestOwnedByCaller`, `TestAdoptStripsGroupAndOtherWrite`, `TestThrowawayParentIsPrivate`, `TestDiffUsesTheHomeNotTmp`, `TestCellIsRuneSafe`, `TestYamlQuoteEscapesControlCharacters` → all PASS, exit 0.
- Targeted `cmd/sh-agent` tests (8 named): `TestUpdateOnPlainCopyNamesTheAlternative`, `TestAdoptOnPlainCopySaysWhereToCommit`, `TestAdoptOnReleaseBinaryNamesTheAlternative`, `TestStatusOnPlainCopyNamesTheRefusal`, `TestSyncIsInstall`, `TestDiffVerbShowsDrift`, `TestRelativeHomeResolvesAgainstTheInvokingDirectory`, `TestStatusWarnsOnOversizedSettings` → all PASS, exit 0.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase-2` → succeeded, "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed", settings.json created with attribution off.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-phase-2` → "drift none — 259 files match the repo" (the +1 vs 260 written is the settings.json which status doesn't count as a managed/tracked file in that line — consistent with install's own report, not a discrepancy investigated further as it is outside the named scope).
- `govulncheck ./...` → "No vulnerabilities found." (go1.27.1 toolchain in use, per AGENTS.md's expectation of clean on go1.25.8+).

## Live checks
(a) Relative `--home` with `--repo`, BUG-080: from a scratch cwd (`.agents/scratch/verify-phase-2-cwd`), ran `go run "$REPO/cmd/sh-agent" install --repo "$REPO/.agents/scratch/verify-phase-2-clone" --home rel-home`. Result: `.agents/scratch/verify-phase-2-cwd/rel-home/.config/sh-agent/state.json` exists; `.agents/scratch/verify-phase-2-clone/rel-home` does not exist. Confirms `--home` is resolved against the invoking cwd, not the `--repo` dir, on a real (non-dry-run) install.

(b) Plain copy refusal: copied `context skills agents commands rules` into `.agents/scratch/verify-phase-2-plain` (no `.git`), then `go run ./cmd/sh-agent update --repo .agents/scratch/verify-phase-2-plain --home .agents/scratch/verify-phase-2-plain-home` → exit 1:
```
sh-agent: not its own git checkout; a plain copy cannot be updated with git — `install` still deploys it as it is; to move forward, use a checkout (`git clone git@github.com:kesokaj/styrman-harness.git`) or a release binary (README § Install without cloning)
```
Names `install`, a checkout, and a release binary as required.

(c) `diff` after drift: `echo drift >> .agents/scratch/verify-phase-2/.claude/AGENTS.md` then `go run ./cmd/sh-agent diff --home .agents/scratch/verify-phase-2` → produced a unified-diff hunk ending `-drift` (the repo's rendered version lacks the injected line). `ls .agents/scratch/verify-phase-2/.config/sh-agent/` after the diff shows only `state.json` — no `sh-agent-diff-*` temp file left behind.

## Root cause  (FAIL only)
N/A — Status is PASS.

## Coverage of changed code
Cross-checked against `.agents/docs/reviews/2026-09-18-review-harness-merged.md` § Outcome (19 items). Every "fixed" item that names a function is named by a test I ran:
- item 1 (adopt write bits) → `TestAdoptStripsGroupAndOtherWrite`
- item 2 (update plain-copy message) → `TestUpdateOnPlainCopyNamesTheAlternative`, plus live check (b)
- item 6 (canonical release tags / compareRelease refusal) → `TestIsReleaseVersionIsCanonical`, `TestNonCanonicalTagNeverReachesTheDownload`
- item 7 (BUG-079 note real-home-only, ownership) → `TestNoteRunningElsewhereOnlyForTheRealHome`, `TestOwnedByCaller`
- item 8 (adopt on plain copy message) → `TestAdoptOnPlainCopySaysWhereToCommit`
- item 9 (four coverage gaps) → `TestAdoptOnReleaseBinaryNamesTheAlternative`, `TestStatusOnPlainCopyNamesTheRefusal`, `TestSyncIsInstall`, `TestDiffVerbShowsDrift`
- item 10 (relative --home, BUG-080) → `TestRelativeHomeResolvesAgainstTheInvokingDirectory`, plus live check (a)
- item 11 (one settings bound) → `TestStatusWarnsOnOversizedSettings`
- item 13 (throwaway parent 0700) → `TestThrowawayParentIsPrivate`
- item 14 (Diff temp under home, not /tmp) → `TestDiffUsesTheHomeNotTmp`, plus live check (c) (no stray temp file)
- item 15 (cell rune-safe) → `TestCellIsRuneSafe`
- item 16 (yamlQuote control escapes) → `TestYamlQuoteEscapesControlCharacters`
- items 3, 4, 19 (docs-only) → out of scope for this verifier dispatch (documentation agent's domain); not re-checked here.
- items 5, 12, 17, 18 (left/structural, no behaviour change, or folded into item 11/12's shared helper) → no dedicated test named, consistent with "left" disposition in the Outcome record; no gap found.

No changed non-test function was found without a naming test within the scope of this dispatch.

## Not examined
- Items 3, 4, 19 of the review's Outcome (documentation-only fixes to README/delivery.md/local.md) — not independently verified against file content in this run; only their existence as "fixed by documentation" was cross-checked against the record, not the rendered docs themselves.
- No behavioural CLI run of `claude -p --agents ...` or `--append-system-prompt` sandbox runs (guide.md's agent/context behavioural-run requirement) — this range's diff is confined to `cmd/` and `internal/`, not `context/` or `agents/*.md`, so that tier does not apply; confirmed no such paths in the range via the review record and commit list.
- The `.agents/docs/reviews/2026-09-18-review-160aa16-3209cf3-docs.md` file that appeared as untracked during this run was not written by me and was not read/verified — it is another agent's concurrent artifact, outside this dispatch's scope.

## Environment
- OS: Linux 6.12.107+deb13-amd64
- Go: go1.27.1 linux/amd64 (user-level SDK ahead of /usr/bin/go, per this repo's AGENTS.md)
- govulncheck: present at ~/.go/bin/govulncheck, ran clean
- No Docker/Colima involved (no Dockerfile in this range)
