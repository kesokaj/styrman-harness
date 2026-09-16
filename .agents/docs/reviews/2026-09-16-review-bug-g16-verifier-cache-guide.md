---
owner: reviewer
date: 2026-09-16
---
# Review — bug-g16-verifier-cache-guide   (2026-09-16)

Range reviewed: working tree, agents/verifier.md
Blob sha (as read): c66fa86707fa17f8b03e82ab0f9fc8f17b4c9b72

## Verdict: APPROVE WITH NITS

The diff is two prose edits to the verifier agent definition. The `guide.md`
edit is a clear improvement — it replaces a subjective trigger ("no guide
describes" the procedure) with a decidable one (the file does not exist),
and stays consistent with `context/AGENTS.md`'s "created by the verifier on
first use" and `agents/documentation.md`'s identical phrasing; no stale
copy elsewhere needs updating. The `-count=1` clause is directionally sound
(Go's test cache genuinely can't see a subprocess-built binary or a file a
test reads at runtime) but its trigger condition is written with enough
grammatical and epistemic slack that two careful runs of this same agent,
on the same diff, could reasonably reach different verdicts — and the file
gives it no fallback the way its neighboring rules have one.

## Findings

- **[medium / confidence: medium]** `agents/verifier.md:43` — "with `-count=1`
  when any test builds or launches a binary or reads a file the change set
  touched" is not decidable the mechanical way the rest of this section is
  (line-counting `--- PASS:`, a `git diff <sha>` check, a `.only(` grep). It
  requires the agent to know, transitively, whether any test in the package
  under verification builds a binary, execs a subprocess, or reads a file at
  runtime — behavior that isn't always visible from reading the test source
  (a helper package, an indirect `os.ReadFile` two calls deep, a `testdata/`
  fixture reached through a third-party library) and that the agent has no
  instructed procedure (no named grep pattern, no "if unsure, do X") for
  establishing. The clause also parses two ways: "[builds or launches a
  binary] or [reads a file the change set touched]" reads as an unqualified
  binary/launch trigger plus a qualified file-read trigger — i.e., a test
  that runs `go build` or execs any binary trips the flag regardless of
  whether that binary is related to the diff, while a file-read only trips
  it when the file is one the change set touched. If that asymmetry is
  intended, spell it out; if not, the sentence needs a parenthetical to
  scope both halves to the change set the same way. As written, since a
  wrong "no" quietly leaves a cache-shadowed pass in the counts line (the
  failure mode the clause exists to prevent) while a wrong "yes" only costs
  a slower run, a verifier that can't decide should default to adding
  `-count=1` — that fallback isn't stated.
- **[low / confidence: high]** `agents/verifier.md:43` — checked for the
  contradiction the dispatch asked about and found none: "a counts line
  never comes from a cached result" governs the *current* run's Go test
  cache; "vs previous" governs comparing this run's freshly-produced counts
  against an *earlier log's* recorded counts. Different mechanisms, no
  conflict, and the counts-line format (`counts: pass N, fail N, skip N …
  for <the command as run>`) already accommodates a run that adds
  `-count=1` since the exact command is recorded verbatim. Not a finding —
  noted per the dispatch's question, verdict: no contradiction.

## Criteria

1. Dead code — n/a, no code removed or orphaned by this prose change; pass.
2. Flow tracing — the `guide.md` existence check is consistent with every
   caller of that concept (`context/AGENTS.md:17`, `agents/documentation.md:27,35`,
   `skills/init-project/SKILL.md:75`); the `-count=1` clause's own internal
   logic traces cleanly to its stated rationale (test cache blindness) but
   see the finding above on its trigger's decidability.
3. Standards conformance — no language skill applies (prose agent
   definition); frontmatter still parses (`go run ./cmd/sh-agent list` shows
   7 agents, `verifier` intact) and the file still contradicts nothing in
   `context/AGENTS.md` step 1 or step 3.
4. Directory layout — n/a, single existing file edited in place.
5. Simplicity — the added clause is one sentence, no new dependency, no new
   file; the ambiguity is in wording, not in unjustified complexity.
6. Decoupling — n/a, single-file prose change with no cross-file coupling
   introduced.

## Not examined

Did not exercise a live verifier dispatch against a fixture that builds a
binary or reads a file (would confirm in practice whether the trigger
condition resolves the way I read it); did not check the other five
uncommitted/untracked files the dispatch named as out of scope
(`.agents/docs/plans/2026-09-16-linkshelf-simulation.md`,
`.agents/docs/reviews/2026-09-16-review-linkshelf-agy-code.md`,
`.agents/docs/testing/2026-09-16-linkshelf-simulation.md`) — excluded per
the dispatch's scope, not reviewed here.

## Outcome

- Finding 1 (medium) — fixed by deletion of the condition: the counts run now always carries `-count=1` ("a verifier never takes a cached result"), so there is nothing to decide and no wrong "no"; the cost is a fresh run per dispatch, which is what a verifier is for.
- Finding 2 — agreed, not a finding.
- Finding 3 — no change needed.
- Round 1 closed 2026-09-16; a live verifier dispatch at the fixed definition is the sandbox row owed (Re-runs, `testing/2026-09-16-linkshelf-simulation.md`).
