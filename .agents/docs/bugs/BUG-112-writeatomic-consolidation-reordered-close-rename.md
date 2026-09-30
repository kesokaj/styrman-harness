---
owner: bugs
date: 2026-09-30
id: BUG-112
status: fixed
area: internal/harness/apply.go (WriteAtomic)
---
# BUG-112: consolidating the settings writer onto WriteAtomic dropped a diagnostic, then reordered close and rename

## Symptom
Round 1 of the project review consolidated the duplicated atomic-write
sequence (temp file, write, sync, chmod, rename), used by both the settings
writer and the binary installer, onto one `internal/harness.WriteAtomic`
helper. The consolidation dropped the settings writer's
temp-file-left-behind diagnostic (round-1 review finding). The round-2 fix
restored the diagnostic (`abandon`) but moved the temp file's `Close()` to
after `os.Rename`, so a close failure was reported against a file that was
already in its final place, and on Linux a binary could be briefly renamed
into place while still open (round-2 review and security findings).

## Repro
Round 1 diff: `WriteAtomic` had no equivalent of the settings writer's
"temp file left behind" error wrapping. Round 2 diff: `tmp.Close()` called
after `os.Rename(tmp.Name(), dest)` instead of before — no test caught the
ordering because no behavioural test can observe close-vs-rename order from
outside.

## Root cause
The refactor treated "same steps, same result" as sufficient; it did not
compare the two original implementations' step ORDER line by line before
merging them, and the property that makes the order matter (a close failure
must be attributable to a file not yet renamed; ETXTBSY on Linux for an open
binary) is not observable through any behavioural test, so nothing forced
the ordering question to be asked before merging.

## Fix
`internal/harness/apply.go` `WriteAtomic`: write, chmod, sync, close, then
rename (close restored before rename; `abandon` restored for the
diagnostic). `TestWriteAtomicClosesBeforeRename` pins the order at source
level, since no behavioural test can distinguish the two orderings.

## Prevention
When a refactor merges or consolidates a sequence whose ORDER is the
contract (write, sync, close, rename; or any lock/unlock, open/close
pairing), diff the merged sequence against each original line by line before
dispatching it for review, not just its net effect. Any such ordering that
no behavioural test can observe gets its own source-level pin test in the
same change, not added later after a regression.

## Related
Review `.agents/docs/reviews/2026-09-30-review-project-merged.md` Round 1-3;
`security-worktree.md` round-2 finding; this is the same class of "a fix
took more than one attempt on one mechanism" that BUG-026 and BUG-060
recorded — if a fourth record on this pattern appears, the "two rounds on
one mechanism stop the patching" rule (already in `context/AGENTS.md` §
Definition of Done step 2) should be cited explicitly by the reviewer
dispatch rather than rediscovered.
