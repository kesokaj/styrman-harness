---
owner: reviewer
date: 2026-09-16
---
# Review — bug-038-code (2026-09-16)

Range reviewed: working tree — `git diff HEAD -- agents/` (`agents/explorer.md`, `agents/reviewer.md`, `agents/security.md`, `agents/verifier.md`)

## Verdict: APPROVE

The diff is a single identical bullet added to the `## Records` section of the four read-only agents that carry the Bash tool. It closes the gap BUG-038 describes, is scoped correctly (the three Edit/Write-only agents — `bugs`, `changelog`, `documentation` — have no Bash tool and are rightly left out), does not touch frontmatter, and is deployed and translated correctly to both CLIs. No defect found.

## Findings

None. This is a documentation/prevention-rule change to agent definitions with no code path to trace and no correctness risk introduced.

Two non-blocking observations (not findings — out of this diff's scope, since the dispatch scoped review to `agents/` only):
- `.agents/docs/bugs/BUG-038-verifier-git-checkout-discarded-fix.md` still has `status: open` and `Fix: Not fixed`, and `INDEX.md`'s BUG-038 row still reads `open`, even though this diff is exactly the fix the record's own Prevention section calls for. Someone (main agent or `bugs`) still needs to flip the record's status and Fix section once this diff lands — that's a records-hygiene step, not a defect in the four files reviewed here.
- The record's `Related` section recommends the rule eventually be promoted to `context/AGENTS.md`'s subagent section as a global hard rule "not per-agent text." This diff deliberately does the narrower per-agent fix first (matching the dispatch scope); the promotion is a listed follow-up, not something this diff was asked to do, so its absence here is not a gap.

## Checks requested by the dispatch

- **Closes the gap / other agents needing it:** Yes. The line is added verbatim to `verifier.md`, `reviewer.md`, `security.md`, `explorer.md` — exactly the four BUG-038 names. Confirmed by frontmatter inspection that `bugs.md`, `changelog.md`, `documentation.md` declare `tools: Read, Grep, Glob, Edit, Write` (no Bash), so they have no way to invoke `git checkout --` etc. in the first place; leaving them out is correct, not an oversight.
- **Could the wording forbid needed reads (`git diff`, `git show`, `git log`, `git status`, `git grep`)?** No realistic misreading. The rule is scoped by "a git command that **changes the working tree**" and enumerates only mutating commands (`checkout --`, `restore`, `stash`, `clean`, `reset`, `apply`); none of the five read commands appear in the forbidden list, and the sentence positively directs the agent to `git show <rev>:<path>` as the way to read old content. The bug record's own Prevention section additionally says "reading commands... stay allowed" — the shipped line doesn't repeat that allowlist explicitly, but since the forbidden list is closed and read commands are absent from it, no agent following this text would treat `git diff`/`git log`/`git status`/`git grep` as forbidden.
- **Frontmatter / description-line rules still hold:** Confirmed. `git diff HEAD -- agents/` shows every hunk lands inside the `## Records` section, below the frontmatter block; `description:` lines are unchanged (one line each, 404–564 chars, no unquoted `#`). Independently re-ran `go run ./cmd/sh-agent list` (27 skills, 7 agents, all parse) and `go run ./cmd/sh-agent install --home .agents/scratch/<tmp>` (260 files, drift none against the dirty tree) and grepped both derived copies: the Claude `.claude/agents/reviewer.md` and the Antigravity `.gemini/config/agents/*/agent.md` for all four agents each carry `## Records` once and the new BUG-038 line exactly once. This matches the main agent's claim. Scratch dir removed after the check.
- **Sentence placement — one rule set or contradiction?** Reads as one rule set. It's the fourth consecutive bullet under `## Records`, appended directly after "Never: source, tests, config, or another agent's records — ...". The two "Never" bullets are complementary, not overlapping: the first scopes which *files* an edit/write tool must never touch; the new one scopes which *git commands* must never run regardless of tool. Nothing in the new line contradicts or loosens the first — an agent that already writes nothing outside its own report gains an additional constraint on Bash-invoked git, not a conflicting one.

## Criteria

1. Dead code — pass; no code added, nothing to remove.
2. Flow tracing — pass; no function/module changed, so no caller/callee chain to trace. The rule text itself was traced against actual tool availability (Bash) for all seven agents, confirming correct scoping.
3. Standards conformance — pass; matches `AGENTS.md`'s per-project rule that agent definitions are file drops, body text only, quoted description untouched; wording is consistent with the bug record's own Prevention section.
4. Directory layout — pass; edits are in-place in the existing `agents/*.md` files, the correct location per this repo's layout.
5. Simplicity — pass; smallest possible fix (one identical line, four files), no abstraction, no new mechanism; matches the bug record's recommended text almost verbatim.
6. Decoupling — pass; each agent file remains independently editable, the added line is self-contained and doesn't reference or depend on the other three files.

## Not examined

- Whether `context/AGENTS.md` should also gain a global version of this rule (explicitly out of scope per the dispatch, and flagged by the bug record itself as a separate follow-up, not part of this fix).
- `internal/harness/derive.go` and `plan.go` source itself were not diffed/re-reviewed line-by-line (unchanged by this diff; only their *output* was spot-checked via a scratch install to confirm the new line survives translation).
- No other files in the working tree were reviewed; the dispatch scoped this review to `agents/` only, and `git status` confirms no other paths are modified.

## Outcome

- **No findings.** The two observations: BUG-038's record is closed in this commit; the promotion into `context/AGENTS.md` was done in the same commit after the docs lens raised it as a finding (`reviews/2026-09-16-review-bug-038-docs.md`), as one clause on the Git paragraph, with the render, scratch-install and sandbox evidence recorded there.
