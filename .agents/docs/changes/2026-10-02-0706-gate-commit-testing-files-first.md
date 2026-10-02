---
owner: main
date: 2026-10-02
---
# Gate: commit testing files before the verifier runs (2026-10-02-0706)

The previous commit (30db86d) made an uncommitted testing file stop every log from counting, but said nothing about committing it, which left the gate with no way to pass (docs review of v2.3.1..30db86d, medium). `skills/preflight/SKILL.md`: Step 2 item 1 now says to commit such a file through Step 1's rows before the verifier runs, and to commit and re-dispatch when the verifier writes the guide; Step 1 names the same exception to "ride along"; the status check lists untracked files too. Two sentences, read by the main agent. Verification: `go test ./internal/harness`, `go run ./cmd/sh-agent list`.
