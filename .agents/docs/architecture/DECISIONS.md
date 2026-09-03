# Architecture Decisions

Append-only. To reverse a decision, append a new ADR and mark the old one
`Superseded by ADR-NNN` — never rewrite it.

## ADR-001: One Go binary replaces the shell and Python harness tooling     (2026-09-03) — Accepted

**Context:** Twelve separately deployed scripts across bash and Python drifted
independently and silently — the deployed harness was found running a constitution and
a hook set that had never been shipped, with two hook scripts referenced but absent.
The hooks also run on every tool call, where interpreter startup is a per-call tax.

**Decision:** Consolidate the three setup scripts, two Python helpers, and seven
hook/gate scripts into one Go binary, `harness`, deployed to `~/.agents/bin/harness`.
The only remaining shell is a bootstrap that runs `go run ./cmd/harness install`.

**Consequences:** One artifact to deploy and one to check for drift. Measured hook
startup drops from ~22.8 ms (Python) to ~2.6 ms. Installing now requires a Go
toolchain, which the operator already has. Cross-compilation for arm64 and amd64 is
enforced in CI.

## ADR-002: Gates inform; they block only at a commit boundary     (2026-09-03) — Accepted

**Context:** The trail gate blocked the Stop hook at the end of every turn in which
source had changed and the day's trail was incomplete. Measured cost of one wrong
block: eight forced continuations, $0.236 and 37.8 s, for a subagent whose entire
output was one word — and the held agent degraded into refusing outright.

**Decision:** The trail gate reports what is owed as injected context, and blocks at
exactly one boundary: a commit the changelog does not cover. The subagent gate holds an
owner that stopped without its artifact at most once, guarded both by
`stop_hook_active` and by an independent hold counter.

**Consequences:** Work in progress is no longer interrupted every turn. The contract is
still enforced, but at the boundary it was always defined on. Both guards fail open, so
an unfamiliar payload costs one wasted turn rather than eight.

## ADR-003: Definition-of-Done rows are proportional     (2026-09-03) — Accepted

**Context:** Every row fired on any source change. For a solo operator, running
reviewer, security and verifier subagents on a three-line edit spends model time for no
benefit.

**Decision:** Split the rows. Resumability rows (`changes/`, verification, handoffs)
always fire — they are what lets a zero-context agent resume, and they are cheap.
Accountability rows fire in proportion to the change, thresholded per project in
`.agents/harness.json`. Verification runs the project's own commands as a script; a
model is spent only to explain a failure.

**Consequences:** Small changes stop paying for large ceremony. The judgement moves
from the model into a computation that can be inspected and tuned. Every dispatch
records its cost, so the thresholds can be argued with using data rather than feel.

## ADR-004: Artifacts are identified by frontmatter, not filename     (2026-09-03) — Accepted

**Context:** Gates matched artifacts with a filename glob, `<date>-<infix>*.md`. A
dispatch prompt that named an output path silently broke the match — observed live,
where a recon brief written as `2026-09-03-harness-recon.md` did not match the required
`-recon-` infix and would have held the agent in a loop.

**Decision:** Every Definition-of-Done artifact declares `owner:` and `date:` in YAML
frontmatter. Gates locate artifacts by that frontmatter. Filenames are free.

**Consequences:** A dispatch prompt can name any path it likes. Files without the
frontmatter are ignored rather than miscounted, so a stray markdown file in `reviews/`
can no longer satisfy a row it does not own.

## ADR-005: The credential gate asks git, it does not parse the command     (2026-09-03) — Accepted

**Context:** The first version of the credential gate decided by parsing the shell
command's arguments. An independent security review defeated it fifteen ways out of
twenty-one probed: `git -C . add -A` (git's own global options consumed the
subcommand), `git add .e*` (a glob the shell expands), `git add secrets/` (a
directory), `git add -u`, `git add .en'v'` (shell quoting), and — worst — a credential
already in the index, which any later bare `git commit` would ship without ever
looking like a staging command.

**Decision:** Stop parsing. For any command mentioning a git staging subcommand, ask
git what it actually holds: scan the index, and scan the credential-shaped files in
the work tree that a stage could pick up. The command's spelling stops mattering. To
avoid failing closed, a command that unambiguously names existing ordinary files is
still allowed, and a plain `git commit` with a clean index is not tree-scanned.

**Consequences:** The gate is robust to spellings nobody enumerated, at the cost of
two `git` invocations per staging command. `.env.example` and its siblings are
explicitly exempted — a gate that refuses to commit a template is a gate that gets
turned off. All fifteen bypasses are now regression tests.

## ADR-006: Formatters run only from outside the tree being edited     (2026-09-03) — Accepted

**Context:** The PostToolUse formatter ran `npx --no-install prettier`, which
deliberately resolves the repository's own `node_modules/.bin/prettier`. The hook
fires automatically on any Edit or Write to markdown, JSON, YAML or TypeScript, on
both vendors, with no operator action — so opening a hostile repository and editing a
single file executed that repository's code. It also cost a measured 419 ms per edit.

**Decision:** Resolve every formatter through PATH, follow symlinks, and refuse any
binary that resolves inside the working repository or a `node_modules` directory. No
`npx`.

**Consequences:** Formatting now happens only where the operator has installed the
tool themselves. A project relying solely on a local prettier gets no auto-formatting
from the harness — its own tooling still applies. The 419 ms probe is gone.

## ADR-007: Rows are satisfied by scope, not by date     (2026-09-03) — Accepted

**Context:** `HasArtifact(row, day)` matched any artifact owned by the row and dated
today, so a second substantive change on the same day owed nothing. The gap hid owed
rows five separate times on the day it was found — including for a 251-line round and
for changes to the gate logic itself, both of which had to be dispatched by hand. A
parallel convention had also grown up alongside it: `/preflight` matched release ranges
against a `Scope:` line in a report's body, while the gates matched frontmatter. Two
mechanisms for one question.

**Decision:** Every artifact declares `scope:` in its frontmatter, and a row is
satisfied only by an artifact whose scope matches the current one exactly. A scope is
`changeset:<digest>`: a digest over the paths and contents of the changed source files,
one form spanning both sides of a commit. An artifact with no scope covers nothing.

Two earlier designs were written and rejected in review before shipping. A two-form
scope — `worktree:<digest>` before a commit, `commit:<sha>` after — re-owed every row the
moment work was committed, a gate rejecting what it had just accepted; content rather
than diff is what makes one form possible, since committing alters no byte on disk. A
`<base>..<head>` range form was unreachable, since nothing produced a scope it could
match; range coverage is a judgement the release gate makes from frontmatter, not
something the binary computes.

**Consequences:** The worktree digest changes the moment the tree does, so an artifact
stops covering work that has moved on — the property dates never had. Artifacts written
before this existed satisfy nothing and their rows come due again, which is correct
rather than unfortunate. `/preflight` reads the same frontmatter as the gates, so there
is one mechanism. Cost: two extra git invocations and a hash per gate evaluation.

## ADR-008: The dispatcher's CLI contract is tested against a stub     (2026-09-03) — Accepted

**Context:** The dispatcher shells out to `claude -p` and `agy -p`. Both need
credentials and cost money per run, so the contract — argument shape, result envelope,
cost field, timeout behaviour, artifact detection — was only ever exercised by hand. In
this project, every seam checked only by hand turned out to contain a bug.

**Decision:** Test against a stub script that mimics the CLIs: it parses the prompt,
writes a conforming artifact, and prints the result envelope. Behaviours are selected
per test — writes an artifact, writes nothing, hangs.

**Consequences:** CI covers the contract with no credentials and no spend. The stub
immediately found that `exec.CommandContext` cancellation does not bound wall time —
a one-second timeout took thirty seconds, because `Output()` waits for inherited pipes
to close; fixed with `WaitDelay`. The stub tests the harness's side of the contract,
not the vendors': a vendor changing its envelope would still only surface in use.

## ADR-009: The changelog row is scoped to a task, not to a code state     (2026-09-03) — Accepted

**Operator approval:** `REQUIREMENTS.md` is immutable without the operator's explicit
approval. R7 was rewritten twice on 2026-09-03 — once for scope-based coverage, once for
this narrative distinction. The operator was asked directly and chose "Approve both
changes"; this note exists because a review correctly observed that the approval left no
trace on disk, and an approval nobody can find is indistinguishable from one never given.

**Context:** ADR-007 made every row satisfied by a byte-exact `changeset:<digest>`. That is
right for rows recording a property of the code — a review of the old code is not a review
of the new code — but wrong for the changelog. A digest moves on every edit, so on a
working tree the entry was stale the moment it was written, and the changelog agent's own
anti-duplication rule then declined to write another. Observed with three valid entries on
disk and the row still failing.

**Decision:** Distinguish narrative rows from code-state rows. The changelog is the only
narrative row: it records what was done and why, a human-meaningful unit that does not stop
being true because a line moved. While the working tree is dirty, the changelog row is
satisfied by any entry from today that declares some scope. At a commit the tree is clean,
the exception lapses, and the entry must cover HEAD's scope exactly.

**Consequences:** The changelog stops being re-dispatched on every keystroke — the "agents
working for no benefit" cost the whole redesign exists to remove. The original bug stays
fixed, because it was about commits: a second commit on the same day still owes an entry,
and the blocking gate fires only at that boundary, where the exception does not apply. The
exception still requires a properly-formed artifact, so "an artifact without a scope covers
nothing" is untouched. Accountability rows are deliberately excluded: reviewer, security and
verifier re-run whenever the code moves, which is the point of them.
