---
owner: main
date: 2026-10-09
---
# Before v2.6.0: toolchain, the retro matrix check, partial clones in search

Operator, 2026-10-09: "Go with your recommendations. ANd try to fix all in the next release".
- **Toolchain.** This box moved to go1.27.2: a user-level SDK at `~/sdk/go1.27.2`, symlinked at `~/.local/bin/go` ahead of `/usr/local/bin/go`, following `deploy/local.md`'s procedure. `govulncheck ./...` now reports 0 vulnerabilities the code reaches; GO-2026-6604, -6600 and -6599 are gone. The deploy guide's procedure and AGENTS.md's toolchain notes are updated, including where this box's `go` lives.
- **`testing/skill-matrix.md`, the `retro` row's check.** It now accepts a written retro record with a proposal of either scope; a harness-scope proposal must still be queued under `proposals/`. The old check failed every project-scope run, the 2c0f215 control included. The new check was run against the four kept sandbox runs (all pass), an empty clone (fails) and an unqueued harness record (fails).
- **Search and partial clones (BUG-165, security low 3, now closed).** Search refuses a partial clone only when git is older than 2.44, read once from `git version` outside any repository (`trail.LazyFetchSafe`, with `trail.PartialClone` exported). Pinned by `TestGitViewPartialCloneNeedsNewGit` (red before the fix) and `TestVersionAtLeast`. The search spec, BUG-165 and the security review's Outcome are updated.
- **Verification:** `go test -race -timeout 30m ./...` exit 0 on go1.27.2 (8 packages), gofmt and vet clean. The release gate's reviewer, security and verifier cover this commit within v2.5.0..HEAD.
