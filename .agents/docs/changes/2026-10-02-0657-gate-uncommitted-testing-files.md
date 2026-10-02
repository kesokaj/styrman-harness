---
owner: main
date: 2026-10-02
---
# Gate: an uncommitted testing file makes no log current (2026-10-02-0657)

`skills/preflight/SKILL.md` Step 2 item 1: an uncommitted change to the testing guide, a matrix, a kit or a script under the docs root now counts like an uncommitted source change: no log reaches HEAD until it is committed. Closes row 1 of the v2.3.1 gate report (docs lens finding A). One sentence, read by the main agent. Verification: `go test ./internal/harness`, `go run ./cmd/sh-agent list`.
