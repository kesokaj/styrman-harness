# Requirements — styrman-harness

The agent harness: one contract, deployed to two vendors.

> Immutable without the operator's explicit approval. Agents never edit this file.

## R1 — One contract, two vendors

`shared/AGENTS.md` is the single source of the constitution. `claude/CLAUDE.md` and
`antigravity/GEMINI.md` are **generated copies** and are never hand-edited. The six
subagents are authored for Claude Code; the Antigravity flavours are derived from them.
CI fails if any generated file has drifted from its source.

## R2 — One deployable artifact

All harness tooling is one Go binary, `harness`. No second language, no directory of
scripts. The only shell in the repo is a bootstrap that runs `go run ./cmd/harness
install`, because the binary cannot install itself before it exists.

Rationale: hooks run on every tool call, where a Go process costs ~2.6 ms against
~22.8 ms for a Python one; and twelve separately deployed files drift independently
and silently, which is how this harness once ran a configuration nobody had shipped.

## R3 — Deployment overwrites, and drift is loud

`harness install` overwrites everything it manages. It never merges, never writes
backups. Machine-local settings live in `~/.claude/settings.local.json`, which the
installer never touches.

`harness doctor` runs at the end of every install and on every session start. It fails
when a managed file differs from the repo, when a hook names a command that does not
exist or is not executable, or when the Antigravity plugin is copied but not
registered.

## R4 — Gates inform; they block only at a boundary

A gate that blocks costs a full model round-trip whenever it is wrong — measured at
eight forced continuations, $0.236 and 37.8 s for a subagent whose entire output was
one word, and the held agent degraded into refusing outright.

So: the trail gate reports what is owed as context, and blocks at exactly one
boundary — a commit that the changelog does not cover. The subagent gate holds an
owner that stopped without writing its artifact at most once, and every guard in it
fails open.

## R5 — Definition-of-Done rows are proportional

Resumability rows (`changes/`, verification, `handoffs/`) always fire: they are what
lets a zero-context agent resume, and they are cheap. Accountability rows (reviewer,
security, documentation) fire in proportion to the change, because they are
re-derivable from the diff and running all of them on a three-line edit spends model
time for no benefit. Thresholds are per-project in `.agents/harness.json`.

## R6 — Spend a model only where judgement is needed

Verification runs the project's own commands as a script. A model is invoked only to
explain a failure. Every dispatched agent replies in at most ten lines; the detail
lives in its artifact.

Every dispatch records its cost to `.agents/docs/testing/dispatch-ledger.jsonl`, so
"this row is not worth its price" is an argument that can be settled with data.

## R7 — Artifacts are identified by content, and record what they cover

Every Definition-of-Done artifact declares `owner:`, `date:` and `scope:` in YAML
frontmatter. Gates locate artifacts by that frontmatter, never by filename glob — a
dispatch prompt that names an output path must not be able to silently break the
contract.

`scope` is what a row is satisfied *against*. Dates were the first answer and were far
too coarse: any artifact dated today satisfied its row, so a second change on the same
day owed nothing. That gap hid owed rows five separate times on the day it was found,
including for changes to the gate logic itself. A scope is `changeset:<digest>` — one
form, a digest over the paths and contents of the changed source files. Content rather
than diff is what lets it span the commit that lands the work: committing moves files
from untracked to tracked without altering a byte on disk, so the digest is identical
either side. An artifact with no scope covers nothing, so one written before the
mechanism existed cannot silently satisfy a row.

Rows scoped to a code state — reviewer, security, verifier — re-owe whenever the code
moves, which is the point of them. The changelog is scoped to a *task* instead: while work
is in progress, today's entry stands, because a narrative does not stop being true because
a line moved. At a commit the tree is clean and the entry must cover HEAD exactly. See
ADR-009.

Coverage decides the blocking gate too, not just the dispatcher. Deciding it by date
there would have left the harness's one blocking gate with the very bug this replaced:
two commits on one day, only the first with a changelog, and both dates reading "today".

This is also the single mechanism for coverage: `/preflight` reads the same frontmatter
rather than parsing a `Scope:` line out of a report body.

## R8 — Isolation covers the host; the gate covers the remote

The operator runs both CLIs inside an isolated microVM, so the harness does not gate
local destructive commands. It gates what leaves the VM: staging or committing a
credential-shaped path is refused.

That gate decides by asking git what it actually holds — the index, and the
credential-shaped files a stage could pick up — never by parsing the command's
arguments. A security review defeated the argument-parsing version fifteen ways
(`git -C . add -A`, `git add .e*`, `git add secrets/`, `git add .en'v'`, and a secret
already staged that a bare `git commit` would ship). Shell and git syntax is not
something a gate can reliably parse; git's own state is.

Two consequences the harness accepts rather than hides:

- **Verification executes the repository's code.** Running a project's tests runs its
  test files, and `verify` in `.agents/harness.json` is a repo-committed command list.
  This is inherent to verifying anything. The microVM is the boundary that makes it
  acceptable; `harness dispatch` prints every command before running it, and names
  whether it came from the repo or from auto-detection, so it is never silent.
- **Formatters run only from outside the tree.** The PostToolUse formatter resolves
  its binary through PATH and refuses anything inside the working repository or a
  `node_modules` directory. An earlier version ran `npx --no-install prettier`, which
  resolves the repo's own `node_modules/.bin/prettier` — editing a markdown file in a
  hostile repo executed that repo's code, with no operator action.

## Known limits

Accepted, not outstanding. Each was investigated and left as it is for the reason given.

- **The `documentation` row cannot be computed.** "Behaviour or interfaces changed" is
  a judgement, not a property of a diff. It stays model-triggered; the dispatcher
  deliberately declares no such row rather than declaring one that never fires.
- **agy tool names cannot be enumerated.** Antigravity offers no way to list the tool
  names valid in an agent's `tools:` list, and an unrecognised name hangs the subagent
  at dispatch. `harness derive` therefore fails the build on any tool it cannot map,
  and `harness doctor` runs `agy plugin validate` — the closest check that exists.
- **Antigravity IDE behaviour is unverified.** Only the `agy` CLI is installed on the
  operator's machine. Whether the IDE loads a plugin without `agy plugin install`, and
  whether it fires plugin `PreToolUse` hooks at all, are untested; open upstream bug
  reports suggest the latter may not work. The credential gate is therefore load-bearing
  on Claude Code and best-effort on the Antigravity IDE.
- **No README.** Removed deliberately in `a63a208`. `AGENTS.md` covers orientation and
  this file covers intent; a third document would only drift from both.
- **Bundled vendor skills stay enabled.** `disableBundledSkills` was considered and
  rejected: the harness skills are domain knowledge (Go, Terraform, GKE) rather than
  workflow duplicates, and `review-standards` is the reviewer's criteria library, not a
  competing review command. The constitution already leaves `/code-review` untouched.

## Non-goals

- Plugin or marketplace packaging on the Claude side. Verified: a plugin namespaces
  both skills and agents, which would rename all 28 skills and all 6 agents.
- Agent Plugins 1.0. Its v1 covers Skills and MCP only; hooks, agents and rules —
  most of this harness — are out of scope.
- Sandbox settings. The microVM already provides isolation.
- Workflow-tool orchestration for `/preflight`. It is Claude Code only, and the
  shell-driven dispatcher works on both vendors; the trade only shifts if Antigravity
  gains an equivalent.
- A dispatcher *agent*. Orchestration is a computation, not a judgement call.
