---
owner: bugs
date: 2026-09-30
id: BUG-118
status: fixed
area: internal/memory, internal/dream (locks)
---
# BUG-118: pid-file locks raced and could lose an index line or stop the nightly job

## Symptom
Code review round 1, findings 4 and 5: two concurrent `memory add` could lose an index line; after a crash a reused pid made the dream lock look held and silently stopped the nightly job.

## Repro
Read the lock code: an empty lock file was read as stale; two waiters could both remove a stale lock; unlock removed whichever holder's lock was present; the dream lock had no age bound.

## Root cause
Staleness rules layered on a pid file are check-then-act sequences with no atomicity, and a pid is reused by the OS. The mechanism was wrong, not one of its edge cases.

## Fix
`syscall.Flock` on both lock files (released by the kernel on exit). Tests: TestHeldLockRefusesAfterRetry, TestLeftoverLockFileIsNoLock (memory), TestRunLock (dream).

## Prevention
A lock between processes on one machine is a flock, never a pid file with staleness rules.

## Related
BUG-117 (same review round), BUG-113 (skills/dream).
