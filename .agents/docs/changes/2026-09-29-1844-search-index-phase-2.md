---
owner: main
date: 2026-09-29
---
# Search index, Phase 2: the global search index   (2026-09-29-1844)

Plan: `plans/2026-09-29-search-index.md`. Spec: `specs/sh-agent/search.md` § Where the indexes live.

- **`--global`** searches `~/.config/sh-agent/search.db`: `feedback.jsonl`, `retros.jsonl` and `proposals/*.md` (kinds `feedback`, `retro`, `proposal`) and the docs root of every checkout that has run `search` or `index`, printed as `<name>:<path>:<line>`. The code shipped in Phase 1's commit; this phase added its tests and the spec's registration paragraph.
- **Registration:** any run inside a checkout registers it, `--global` included, at most once a day. A registered root that is no longer a git checkout is dropped with a note. A second checkout with the same base name prints as `name~2`. `index --global --rebuild` drops the registry with a note.
- **Tests:** `internal/search/global_test.go`, six tests. Written after the code: mutations G1–G5 each compiled and failed a named test.

Agents dispatched: `verifier` (`testing/2026-09-29-verify-search-index-phase-2.md`, PASS over bfcedcc..7fd56d5).

Verification: `go test -race ./...` exit 0; a hand integration across two scratch checkouts found the other checkout's open bug with `--global`, `--global --open` and `--global --refs`.
