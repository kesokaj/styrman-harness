---
owner: main
date: 2026-09-29
---
# Search index, Phases 3 and 4: instructions, recall, and the Done reviews   (2026-09-29-2009)

Plan: `plans/2026-09-29-search-index.md`. Spec: `specs/sh-agent/search.md`.

## Phase 3 — instructions

- **Context:** a "Search first" entry under § Work (D31), paid for by three cuts: the `find-skills`/`npx skills add` procedure shortened, the `/skills` and `/agents` sentence moved to `skills/init-project`, and "Read the index before touching an area it names" folded into the entry's `--refs` clause. The rendered context is 21,988 of 22,000 bytes.
- **Agents:** `explorer` (first command), `reviewer` (first pre-flight step), `security` (scope step), `dreamer` (beside its grep, and for unwritten prevention rules), and `documentation`, `bugs` and `changelog`. The last three use the search output the main agent puts in their dispatch; they have no shell (D26, the operator's choice).
- **Skills:** `review-standards`, `dream`, `handoff` (`--open --kind plan`, the awk scan as fallback), `bug` and `feature` (`--refs <path>`, and the dispatch rule for the three shell-less agents), `retro` (source 9: `--global`), `research` (earlier research first), and `init-project` (the moved sentence).
- **Behavioural runs** (`testing/2026-09-29-search-index-behavioural.md`): on Claude Code, round 2 had every changed agent search in 3 of 3 runs, and the context entry led to a first-call search in 2 of 3 runs (control 0 of 3). Round 1's weaker wording (0 or 1 of 3) led to the rewrite (D29). On Antigravity, the agents that completed searched; `dreamer` and one `security` run never completed, because print mode ended first.
- **Docs:** README rows for both verbs, the spec, `deploy/local.md`, `testing/guide.md` § 3d, and the `sh-agent-matrix.md` and `skill-matrix.md` rows.

## Phase 4 — recall and the vector gate

- Word search now puts chunks holding every word first and fills the remaining slots with any-word chunks (D27). The first recall run exposed AND-only matching (0.077).
- `--refs <bare file name>` prints a hint naming the tracked paths ending in that name (D28).
- Recall on 41 questions (`testing/2026-09-29-search-recall.md`): index 0.653 against grep 0.379, and 0.677 on the final build.
- Vectors: not adopted (D30). MiniLM fusion adds 6.4 points over the index alone, under the 10-point bar. The "pure Go" model library also needs ONNX Runtime.

## Done reviews

- `reviewer`, code over c04f2d7..565bb1e: APPROVE. Over the fix delta 565bb1e..aae16af: APPROVE with two lows, fixed (git runs with `LC_ALL=C`; a comment). Records: `reviews/2026-09-29-review-search-index-code.md`, `-code-r2.md` with Outcome.
- `reviewer`, docs lens over 8ee5d70..HEAD: APPROVE WITH NITS. Three nits fixed: ADR-045a, the `govulncheck` note, and `go mod tidy`.
- `security`, three rounds, with Outcome in `reviews/2026-09-29-security-search-index.md`. Round 1 fixed six findings, each with a test that failed first: DSN escaping, registration only for `.agents/docs` checkouts, own blob IDs, printable output, 0700 directories, secret-shaped names skipped. Round 2 fixed the global docs root, the hint and note lines, the case-insensitive secret list and the descriptions; registration stays automatic by the operator's choice (D34). Round 3, the cap, reverted the read-once change and documented the edit-during-lock race (the operator's choice, D35), and a check that errors now keeps the registration.
- `verifier`, final over c04f2d7..aae16af: PASS (`testing/2026-09-29-verify-search-index-final.md`). The `LC_ALL=C` fix landed after it: `go test -race ./...` exit 0 on HEAD, run bare by the main agent.
- `bugs`: BUG-106 (run clones inside the checkout) and BUG-107 (test runs registered checkouts in the real home's global search index). Both rules graduated into the project `AGENTS.md` § Testing.

Agents dispatched in these phases: `reviewer` ×3, `security` ×3, `verifier`, `bugs` ×3.
