---
owner: reviewer
date: 2026-10-01
---
# Review — c62e064-worktree-docs   (2026-10-01)

Range reviewed: working tree against commit c62e064 (`git diff c62e064`, 12 tracked files) plus the untracked `.agents/docs/handoffs/2026-10-01-0601-audit-fix-phases-0-1.md`, `.agents/docs/testing/kit-preamble.sh`, `skills/review/scripts/restatement-grep.sh`. Docs lens (`/review --docs`).

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The two new mechanisms do not do what their documentation says. `kit-preamble.sh` claims to warn when the real `~/.config/sh-agent` changed (BUG-123's graduation rests on it) but compares `ls -R` name listings, so an append to an existing `feedback.jsonl` (the exact BUG-123 case) is not detected; checked by hand. `restatement-grep.sh` cannot search any word starting with `-` (the flags this lens exists for) and excludes all of `testing/`, so it misses guide.md, which hid two stale pointers to the moved AGENTS.md § Testing rule. The changed spec lines (delivery, search `--global`, dream zero-proposals, memory-loop 200/25 KB, parity plan flags, TR4/TR5 markers) match the code and I found no stale restatement of them.

## Findings
- **[medium / confidence: high]** `.agents/docs/testing/kit-preamble.sh:41,43` — change detection hashes `ls -R` (names only, no sizes or times). Verified in scratch: appending a line to an existing file in a copy of the directory leaves the hash equal, so the "WARNING the real ~/.config/sh-agent changed" never fires for BUG-123's scenario (a skill appended lines to the real feedback.jsonl and retros.jsonl). guide.md:167 and BUG-123's "graduated" status claim this protection. Fix: hash `find "$HOME/.config/sh-agent" -type f -printf '%p %s %T@\n' | sort` (or file contents) and test the append case; also state that the warning does not fail the kit (the trap changes no exit status), or make it.
- **[medium / confidence: high]** `skills/review/scripts/restatement-grep.sh:21-24` — every argument starting with `-` is "unknown option" (exit 2); `--` is also rejected. The lens's own key words (`--no-session-persistence`, `--append-system-prompt`, `--global`, `--strict-mcp-config`) cannot be searched. Fix: accept `--` as end of options and treat an unknown `-x` after it as a word.
- **[medium / confidence: high]** `restatement-grep.sh:31` — the living-files definition drops `.agents/docs/testing/` wholesale, but `testing/guide.md` (named living by this lens, by `skills/review/SKILL.md` --docs, and by AGENTS.md), `skill-matrix.md` and `sh-agent-matrix.md` are living docs; only the dated logs are append-only. Consequence found by hand, invisible to the script: `.agents/docs/testing/guide.md:156` "a settings file written for the run (AGENTS.md § Testing)" and `.agents/docs/testing/skill-matrix.md:14` "(`AGENTS.md` § Test…" still point at AGENTS.md for the agy settings.json rule that this range moved to guide § 4 § Behavioural runs (guide.md:171). Fix: exclude `testing/[0-9]{4}-` logs only, and change both sentences to "guide § 4 § Behavioural runs".
- **[low / confidence: high]** `.agents/docs/specs/sh-agent/memory-loop.md:12` — lists `node_type: memory` beside `name`, `description` as a top-level key; `internal/memory/memory.go:168` (and `memory_test.go:62`) write it under `metadata:`. Change to `metadata.node_type: memory`.
- **[low / confidence: medium]** `restatement-grep.sh` — does not exclude `.agents/docs/architecture/DECISIONS.md` (appeared as a hit for "reconcile") or vendored skills (`skills/review-standards/references/prompt-audit/` appeared), both of which the `--docs` lens in `skills/review/SKILL.md:51` and AGENTS.md exclude; and nothing references the script (`skills/review/SKILL.md`, `agents/reviewer.md`: no hit), so a reviewer is not told it exists. Header warning that untracked files are invisible is honest (the script cannot list itself or kit-preamble.sh until staged).
- **[low / confidence: medium]** `kit-preamble.sh:51-54` — an existing `$KIT_HOME/bin/sh-agent.real` is reused without rebuilding, so a re-run kit after a code change tests the old binary; a kit that sets its own EXIT trap afterwards silently drops the after-run check. Rebuild when `KIT_BIN` is unset; note the trap in the header.

Checked and fine: the refusal fires (clone parent inside this checkout: exit 2 with the BUG-106 message, run by hand); `--all` exit codes and logic (0 hits, 1 none, 2 usage); delivery.md:77 (internal/memory and internal/dream absent from the v1.16.2 tree, so "not yet released" is true); search.md:46; dream.md:34 against `internal/dream/run.go:395`; memory-loop.md 201st-line and 200/25 KB against `memory.go:27,134,161,214`; parity plan flags against `dream.go:177`; AGENTS.md § Testing pointers (§ 3e, § 4, § 5 exist); TR4 marker (one review) and TR5 (`grep -c '^- \[ \]'` = 4, the template lines); Graduated table 45 rows as the plan states.

## Criteria
1. Dead code — pass: nothing orphaned; the restatement-grep script is unreferenced (see low finding).
2. Flow tracing — fail: kit-preamble change detection and restatement-grep option/scope handling (findings 1-3).
3. Standards conformance — fail: two stale restatements (guide.md:156, skill-matrix.md:14), one wrong key placement (memory-loop.md:12).
4. Directory layout — pass: script under `skills/review/scripts/`, preamble beside guide in `testing/`.
5. Simplicity — pass: bash, stdlib tools only.
6. Decoupling — pass.

## Not examined
Go code behind the changed claims beyond the lines cited; the full-diff of the untracked handoff's content (read only for existence); plans older than 2026-09-30 for stale restatements (the script lists many historical plan hits, treated as records); README rows beyond lines 190-194; no Antigravity or CLI run of the preamble (only its refusal and its change-detection were exercised, in scratch).

## Outcome

Every finding fixed in the phase (see round 3, `2026-10-01-review-c62e064-worktree-docs-r3.md`, APPROVE with zero findings): the preamble fingerprints sizes and mtimes and fails the kit; the grep script takes dash-led words, keeps the guide, matrix and preamble as living files and excludes records and vendored skills; the two stale pointers and the six Graduated rows name guide § 4 § Behavioural runs and the preamble; `metadata.node_type`; the review skill names the script in the plain --docs brief; the preamble builds from its own checkout.
