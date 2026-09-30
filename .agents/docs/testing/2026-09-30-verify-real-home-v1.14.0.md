---
owner: main
date: 2026-09-30
---
# Verification — real home on the Linux box before v1.14.0, both CLIs   (2026-09-30)

Range reviewed: v1.13.1..7d6dca0, installed into the real home as 7d6dca0 (`sh-agent update` from this checkout at 10:41: 11 written, 255 unchanged, 266 managed, binary 7d6dca016a1d, `status` with no warning line); every changed file byte-identical between the tree and both vendors' installed copies (`diff` on `agents/verifier.md` and the four `SKILL.md` files, no output). The checks run from the real `$HOME` in fresh clones of a fixture outside every checkout.

## Status: PASS on both CLIs — the three smoke checks, plus one run of each changed skill and of the changed agent on each CLI

The rule that makes this a gate condition: project `AGENTS.md` § Release (no tag until the real-home smoke on both CLIs and one Antigravity run of each changed skill and agent exist).

## Fixture

`~/sandbox/gate14/fixture`: `greetbox`, a Go module with `cmd/app` (a `-home` flag), `internal/greet` (tested), `internal/sweep` (a `filepath.Glob` pattern built from the flag's path — the shape the new `golang` rule forbids), an `AGENTS.md` with Setup commands, `.agents/docs/testing/guide.md`, an empty bug index, `REQUIREMENTS.md`, `ROADMAP.md`, `TODO.md`. `~/sandbox/gate14/run.sh <cc|agy> <scenario> <n>` clones it to `~/sandbox/gate14/runs/<cli>-<scenario>-<n>` (the runs directory is not inside a checkout), proves each scenario's fact before the CLI starts (the glob line present; the seeded regression failing `go test`), and runs from the clone with the REAL home. A wrapper first on `PATH` pins every `sh-agent` call to `~/sandbox/gate14/home` (BUG-107): `~/.config/sh-agent/` lists `feedback.jsonl proposals retros.jsonl search.db state.json` before and after, and the pinned home stayed empty. For the verifier scenario a background writer drops `.agents/docs/reviews/<date>-security-probe.md` two seconds after the verifier's log appears.

## Claude Code 2.1.284 (`claude -p`, real home, installed skills and agents)

| Check | Result |
|---|---|
| Context questions, no tools | PASS — "AGENTS.md — global agent context", 28 skills and 8 subagents, GANNET, `bug` for a failing Go test |
| `golang` skill (review `internal/sweep/sweep.go` and `cmd/app/main.go`, chat only) | PASS — loaded `golang` and `review-standards`; finding 1, high: the home path inside a `filepath.Glob` pattern, resting on the new rule |
| `review-standards` skill (the criteria on `sweep.go`, chat only) | PASS — loaded `review-standards` with its Go checklist; the glob finding again |
| `verifier` dispatch with the foreign artifact | PASS — log written early, Status FAIL for the fixture's untested packages (intended); the foreign file "another agent's artifact (a security review) that appeared while I ran; seen, not mine, not read. No source, test or config path changed" |
| `/bug` on the seeded regression (`Hello` drops its space) | PASS — failing test kept, fix at `internal/greet/greet.go:5`, `BUG-001` written with "Class check: `git grep '\"hello'` found the prefix at this one site only", committed, no push; 105 s |
| `/feature` (a `greetbox search` verb, stop after the plan and its cold read) | PASS — `plans/2026-09-30-greetbox-search.md` with unticked boxes and a "Consumers (who would use the capability, writers included)" list of four, one marked "not a consumer" with the reason; the cold read hit its three-round cap and said so; no code; 732 s |

## Antigravity CLI 1.2.13 (`agy -p`, real home, `--mode accept-edits`, no skip flag)

No run's stderr has an `auto-denied` line; `cli.log` has no `invalid allow entry`.

| Check | Result |
|---|---|
| Context questions, no tools | PASS — the same four answers, GANNET |
| `golang` skill | PASS — loaded `golang`; finding 1 "Dynamic Path Passed to `filepath.Glob`" quoting the new rule verbatim |
| `review-standards` skill | PASS — loaded `review-standards` and `golang`; findings on `sweep.go` |
| `verifier` dispatch with the foreign artifact (the derived `agent.md`) | PASS — log written; "Tree is unchanged across source, tests, and configuration. New untracked files are artifacts under `<docs-root>` and `.agents/scratch/`"; 133 s |
| `/bug` on the seeded regression | PASS — fix committed (`bb81efe`), `go test` green, `BUG-001-hello-spacing.md` written; its record carries no class-check line (the regression is one site, and the rule fires on a class); 514 s |
| `/feature` | PASS — `plans/2026-09-30-greetbox-search.md` with a "Consumers:" list naming a writer consumer and a reader consumer; the plan path returned; 458 s |

## What was not exercised

- The Mac and the devbox: they pick the release up with `sh-agent update`; `darwin/arm64` and `linux/arm64` binaries were built by the gate verifier, not run.
- Interactive sessions on either CLI; every run was headless.
- The `bug` skill's class check on a real class on Antigravity (the fixture's regression has one site); the rule's wording was run on Claude Code with the same one-site fixture and produced the line.
- The `review-standards` § 2 additions as a reviewer would apply them to a consolidation (no consolidation in the fixture); the skill loaded and applied its criteria on both CLIs.
- A second attempt of any scenario: every scenario passed on its first run.
