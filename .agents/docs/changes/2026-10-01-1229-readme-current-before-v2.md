---
owner: main
date: 2026-10-01
---
# 2026-10-01-1229 — README brought current before the v2.0.0 gate

Trivial docs change (5 lines): the § Working on the harness comment block described the two size pins twice and garbled; it now states each once (hand-written context under 15,000, the full render under 22,000). The § Layout block names `internal/memory` and `internal/dream` beside `internal/search`. Checked against the tree at ddac2d8 for this release: REQUIREMENTS.md needs no edit (every claim re-read against the code: search.db and the memory store never sent by `report`, `uninstall` removes the dream timer, `update` refuses a dirty checkout, the four verbs lists match the verb table); `specs/sh-agent/delivery.md` § Where things stand is updated in the release commit, because it names the newest tag. Diff read by the main agent; no agent dispatched.
