---
owner: bugs
date: 2026-09-30
id: BUG-117
status: fixed
area: internal/dream
---
# BUG-117: the agy prompt went on argv, over the 131072-byte execve argument limit

## Symptom
Code review round 1, high, finding 1 (`.agents/docs/reviews/2026-09-30-review-cb799ff-f1ca525-code.md`): a large project's dream failed every night for good, its watermark never advancing.

## Repro
A project whose collected text exceeds 131072 bytes (the dream cap is 256 KiB) run through the agy driver: execve fails with E2BIG.

## Root cause
Linux caps one execve argument at 131072 bytes (MAX_ARG_STRLEN); the driver passed the whole prompt as the `-p` argument while the cap allowed twice that. Argv is also readable by any user in /proc.

## Fix
`internal/dream/run.go` `call` sends the prompt on stdin as one stream-json message (`agy --input-format stream-json --output-format stream-json --mode plan -p ""`), the answer is read from the `result` event. Probed on this box. Test: TestRunAgyPrivateHomeAndStdin.

## Prevention
Large or private data to a child process goes on stdin or in a file, never as one argument; check the payload cap against the OS per-argument limit.

## Related
BUG-116, BUG-119 (same run.go invocation), BUG-113 (skills/dream).
