---
owner: main
date: 2026-09-30
---
# Docs: every description of the harness names the search index   (post-release docs fix, no release change)

The operator asked, after v1.13.0, whether every supporting document describes the shipped feature. A grep over the living docs for verb enumerations and descriptions of what the binary is found four that did not:

- `README.md`: the opening description, the Layout line, and the map's `~/.config/sh-agent/` label (a label, not a line: no path changes).
- `AGENTS.md` (project): the Layout block and the Key files table gain `internal/search`; the Testing section points a search change at `testing/guide.md` § 3d and the matrix.
- `.agents/docs/specs/sh-agent/delivery.md` § Where things stand: the verb list gains `search` and `index`.
- `.agents/docs/deploy/gcp.md`: "a file copier" now "plus a local search index".

Fourteen lines of docs, read by the main agent; `TestHarnessMap` and `TestRenderedContextFitsAntigravityBlock` pass. No agent dispatched.
