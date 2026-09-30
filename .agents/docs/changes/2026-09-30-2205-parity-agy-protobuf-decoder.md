---
owner: main
date: 2026-09-30
---
# Antigravity reader decodes protobuf; gate-review fixes     (2026-09-30)

## Summary
The v1.17.0 gate's code and security reviews found the Antigravity reader's byte-run filter still leaking tool-call argument pieces. Some were glued to a printable length byte; others were split at a multi-byte character. On the operator's word (plan D30), the scan is replaced by a schemaless protobuf wire-format decoder. A tool argument is one whole JSON string field, dropped whole. On copies of the 95 local conversation databases: 971 prose leaves kept, 0 with a tool-argument key.

Also:
- `answer()` voids any stream in which a tool step ran;
- the child's `request-review` is asserted by a test;
- the docs pass covers the timer's absolute PATH, the current agy invocation, the auth and login limits, a guide § 3e for memory, dream and the timer, and "memory store" wording.

The gate's real-home runs found that a real-home `/distill` wrote a test line into the real `feedback.jsonl` and `retros.jsonl`. Both lines were removed, and the finding is recorded.

## Type
fix within the feature (release gate)

## Files changed
- `internal/search/dream.go` and test
- `internal/dream/run.go` and test
- README
- `.agents/docs/{deploy/local.md,deploy/remote.md,specs/sh-agent/dream.md,testing/guide.md}`
- the plan (D30)
- the gate reviews with Outcomes
- the gate verifier log, the real-home log and the docs-lens report

## Agents dispatched
reviewer, security, verifier, docs lens (gate), documentation.

## Verification
`go test -race -count=1 ./...` exit 0. A mutation back to the prefix rule failed `TestAgyConversationText`. A real agy sandbox run: 2 proposals.

## Requirements
R3, R4, R10.
