---
owner: bugs
date: 2026-10-01
id: BUG-129
status: fixed
area: internal/harness/apply.go (installer symlinks)
---
# BUG-129: install, status, diff and uninstall followed symlinks under the home (IN3)

## Symptom
`install` read, chmodded or replaced a symlinked destination through the link; a file under a linked directory (`~/.claude/skills` pointing elsewhere) was written, chmodded and removed through it. `status` and `diff` reported such paths as drift. `uninstall` removed a manifest path that had become a link.

## Repro
`TestInstallDoesNotFollowLinkedDest`, `TestUninstallAndLinkedDirs`, `TestUninstallSkipsLink`, `TestUninstallSkipsFilesUnderLink`, `TestResyncSkipsLink`.

## Root cause
`apply.go` used `ReadFile`/`Stat`/`Chmod` on the destination and `Lstat` only on the final component, so a link in any directory component was followed. Each security round found another verb with the same gap.

## Fix
`Lstat` first and `underLink(home, rel)` (`apply.go:410`) over every directory component, in `Install`'s write loop (`:64`), the resync loop (`:134`), `saveOnError` (`:47`), `Status` (`:222`), `Diff` (`:256`) and `Uninstall` (`:329`). A linked path is kept, printed, counted under `kept`, never in the manifest.

## Prevention
Every file operation under a home directory the operator may have linked `Lstat`s the path AND its components, in every verb that touches it (install, status, diff, uninstall, resync). The test links a directory, not only a file.

## Related
Audit IN3 in `plans/2026-09-30-recon-audit-fix.md`; reviews `reviews/2026-09-30-security-4764c5f-worktree{,-r2}.md`, `reviews/2026-10-01-security-4764c5f-worktree-r3.md`; plan `plans/2026-09-30-audit-fix.md` Phase 1.
