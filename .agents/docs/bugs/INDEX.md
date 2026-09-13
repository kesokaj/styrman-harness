# Bugs

<!-- The project's memory of what went wrong. One row per BUG-NNN record in this directory; newest last. Maintained by the `bugs` agent. Read at resume; consult before touching an area listed here. -->

| ID | Date | Area | Status | Symptom | Prevention |
|---|---|---|---|---|---|
| BUG-001 | 2026-09-03 | installer/status | fixed | `status` reported the installed binary as stale right after `install` | status and diff plan without the binary; exclude `BinPath` from stale checks and test it |
| BUG-002 | 2026-09-03 | cli/flags | fixed | `sh-agent adopt <path> --home DIR` failed with a usage error | Go's flag package stops at the first positional; reorder flags ahead of positionals before parsing |
| BUG-003 | 2026-09-03 | context | fixed | `dev-setup.md` still sent scratch to `/tmp` after the hard rule moved it to `.agents/scratch/`; Claude flagged the contradiction | when a rule moves, grep every context file for the old wording before deploying |
| BUG-004 | 2026-09-03 | skills/frontmatter | fixed | Antigravity dropped `bug` and `feature`: unquoted `argument-hint: [… #issue]` failed its YAML parser; Claude Code was lenient | installer lints frontmatter for unquoted `#`; after a frontmatter change run each CLI once in the sandbox and grep its log |
| BUG-005 | 2026-09-03 | agents/verifier | fixed | `isolation: worktree` put the verifier in a fresh worktree at HEAD, blind to the uncommitted change it was asked to verify | a Claude-only agent key is not free; after any agent frontmatter change, dispatch it once in the sandbox on uncommitted work |
| BUG-006 | 2026-09-13 | installer/binary | fixed | `install` kept copying its own running (stale) binary; real home stuck on the 2026-09-03 build for days | a self-updater must build from source, never copy its own running binary; `status` prints the binary's build stamp vs the repo |
| BUG-009 | 2026-09-13 | cmd/tests | fixed | new uninstall test asserted `~/.claude` itself was removed and failed on first run | read the documented stop/keep list (`pruneEmptyDirs`) before asserting what an operation removes |
