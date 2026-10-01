---
owner: bugs
date: 2026-10-01
id: BUG-151
status: fixed
area: internal/harness/frontmatter.go (Split), installer lint
---
# BUG-151: the frontmatter reader closed a block at a fence followed by two carriage returns

## Symptom
`Split` took any line that was `---` after trimming every trailing `\r` and `\n` for a fence. A line `---\r\r` closed the block, and whatever followed it up to a real `---` (a permission key such as `allowed-tools`) was body text the lint never checked. A reader that does not take `---\r\r` for a fence would read that key as frontmatter. Present in v2.0.0 to v2.2.0.

## Repro
`Split("---\nname: x\n---\r\r\nallowed-tools: Bash\n---\nbody\n")` returned front `name: x\n` with no error; a skill written that way loaded.

## Root cause
`strings.TrimRight(l, "\r\n")` trims a run of characters, where a line ending is one `\n` after at most one `\r`. The lint and a stricter reader therefore disagreed on where the block ends, and the lint checked only its own reading. Found by the security reads of the v2.2.0 gate (`reviews/2026-10-01-security-a62eb9d-768efce.md` finding 5), left as row L5 of `reviews/2026-10-01-preflight-v2.2.0.md`.

## Fix
A fence is `---` after removing one `\n` and then at most one `\r`, for the opening and the closing line (`internal/harness/frontmatter.go`). Test `TestSplitFenceIsOneLine` in `internal/harness/harness_test.go`, red first on three assertions: the block closing early, an opening fence with two carriage returns accepted, a skill with the key behind such a fence loading. Unknown and not checked: whether Claude Code or Antigravity reads such a file's key.

## Prevention
A check that decides where untrusted text ends compares the exact shape, never a trimmed one: trim one known terminator, not a set of characters, and test the neighbour shapes (two of the terminator, another space character, a byte-order mark) as must-refuse cases.

## Related
BUG-143 (same lint family) and the `opensFrontmatter` spellings of the v2.2.0 gate; BUG-152 (same gate, same lesson: a class or exact shape, not a list); `reviews/2026-10-01-security-a62eb9d-768efce.md`.
