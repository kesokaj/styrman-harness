---
owner: verifier
date: 2026-09-30
---
# Verification — agy-only-init (2026-09-30)

Earlier run stopped: the first verifier run on this change was stopped mid-way because context/AGENTS.md had to change (it broke TestRenderedContextFitsAntigravityBlock); this log is the rerun over the final tree.

Range reviewed: 22ca863 (HEAD) working tree, uncommitted (5 modified files) + untracked init_project_test.go

## Status: PASS

Scope: Go tiers (§ 1, § 2 named test), § 3 scratch-home install for `--vendor agy` and default, § 4 by reading the main agent's runs plus the two matrix CHECK commands run here. Not complete: Claude Code path of the changed skills and any live CLI start were not exercised (see Not examined). guide.md exists (.agents/docs/testing/guide.md), followed.

## Commands
- `gofmt -l $(git ls-files '*.go') init_project_test.go` → exit 0, prints nothing
- `go vet ./...` → exit 0
- `go test -race -v -count=1 ./...` → exit 0; output .agents/scratch/verify-agy-only/go-test-1.txt; packages: ok root, cmd/sh-agent, internal/harness, internal/search; counts: pass 375, fail 0, skip 2 at 22ca863+dirty for `go test -race -v -count=1 ./...`; vs previous: .agents/docs/testing/2026-09-30-verify-v1.16.0-gate.md, pass 374→375, skip 2→2 (+1 pass = new TestInitProjectVendorCheckMatchesManifest). Skips: TestRecall (measurement, needs SH_AGENT_RECALL=1), TestHelperIndex (helper process).
  - `TestInitProjectVendorCheckMatchesManifest` → PASS
  - `TestRenderedContextFitsAntigravityBlock` → PASS
- `go build -o .agents/scratch/verify-agy-only/sh-agent ./cmd/sh-agent` → exit 0
- `install --home .agents/scratch/verify-agy-only/home-agy --vendor agy` → exit 0; `status --home …` → exit 0: "vendor: antigravity, 133 files"; "drift none — 132 files match the repo"
- `install --home .agents/scratch/verify-agy-only/home-def` (default) → exit 0; `status` → exit 0: "vendor: all, 266 files"; "drift none — 265 files match the repo"
- `cmp` installed vs tree skills/init-project/SKILL.md and skills/dream/SKILL.md: home-agy (.gemini/config/skills/) both identical; home-def both identical under .claude/skills/ and .gemini/config/skills/ (4 of 4).
- § 4 CHECK of matrix row `init-agy-only` in /home/kesokaj/sandbox/agyonly/proj (expect non-zero) → exit 1 (`.claude` exists) — as expected, pre-fix fixture
- same CHECK in /home/kesokaj/sandbox/agyonly/proj2 (expect zero) → exit 0
- CHECK of row `dream-agy-only` in proj2, OUT=/home/kesokaj/workspace/styrman-harness/.agents/scratch/agyonly-dream.out → exit 0

## § 4 evidence read (no CLI started)
- Pre-fix (must fail): .agents/scratch/agyonly-prefix.out and tasks/b0pejwi81.output: agent scaffolded proj with `.gitignore` lines `.claude/settings.local.json`, `.claude/worktrees/`; the listing shows `.claude/` with symlinks agents/commands/rules/skills and settings.local.json on an Antigravity-only home; fixture inspection: `.claude` exists in proj. The row's check exits 1 there: the defect reproduced.
- Post-fix: agyonly-postfix.out says "Manifest detected Antigravity-only (vendor: antigravity), skipping .claude/ symlinks and settings", store made mode 700; b7xc9vwlm.output: init exit 0, `ls .claude` "No such file", .gitignore holds `.agents/memory`, commit "init: scaffold agent tree and project docs"; proj2 has no `.claude`, has `.agents`.
- Dream: agyonly-dream.out first line `store: /home/kesokaj/sandbox/agyonly/proj2/.agents/memory (1 memory files; …)`, repeated as last line; b7xc9vwlm: dream exit 0.
- Caveat (per the dispatch): post-fix runs used the skill text before the context/AGENTS.md shortening; that edit touches neither skill (diff stat: skills/init-project 18 lines, skills/dream 3, context/AGENTS.md 2). Not re-run.

## Coverage of changed code
- skills/init-project/SKILL.md — covered by init_project_test.go (TestInitProjectVendorCheckMatchesManifest) plus matrix row init-agy-only
- skills/dream/SKILL.md — matrix row dream-agy-only (end-to-end only; no Go test asserts its text)
- context/AGENTS.md — TestRenderedContextFitsAntigravityBlock (size pin only, not wording)
- .agents/docs/specs/sh-agent/memory-loop.md, .agents/docs/testing/skill-matrix.md — docs, no code
- init_project_test.go — is the new test; no t.Skip, nolint, .only or other suppression added (grep exit 1 = none)

## Not examined
- Claude Code path of the changed skills (the release gate's real-home runs cover it).
- No Antigravity or Claude Code start in this run (credentials deleted); § 4 is the main agent's runs read, not re-executed. The post-fix runs predate the context/AGENTS.md edit.
- § 3a–3d, the sh-agent-matrix and full skill-matrix: not run (no cmd/internal Go change in the diff).
- Wording quality of the SKILL text; review lives in .agents/docs/reviews/2026-09-30-1721-review-agy-only-init.md (not read).

## Environment
Linux 6.12.107 debian amd64, go1.27.1; scratch homes under .agents/scratch/verify-agy-only/. Tree at end: `git status --porcelain` identical to start (same 5 M + 3 ??; the tree did not move).
