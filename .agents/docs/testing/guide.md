---
owner: verifier
date: 2026-09-17
---
# Verification guide — sh-agent

The standing procedure for this repo. Written under the rule in `agents/verifier.md`
("What to run"): the first run under a new project writes this guide; later runs cite
it and record outcomes only. Update this file in the same run that changes the steps.

## 1. Go checks (every change, including one that touches only `context/`, `skills/`, `agents/`, `commands/` or `rules/`: the root package's embed tests compare those trees with the binary's embedded copy, so a new file or mode there fails here, BUG-137)

```bash
gofmt -l cmd internal     # must print nothing; NOT `gofmt -l .` — it walks .agents/scratch,
                          # where a stray Go file fails a gate that has nothing to do with it (BUG-048)
go vet ./...
go test -race ./...
```

## 2. Harness self-check (every change to `context/`, `agents/`, or `skills/`)

```bash
go run ./cmd/sh-agent list      # sanity: every skill and agent parses; its section headers print the counts
go run ./cmd/sh-agent render | wc -l   # rendered context/AGENTS.md must stay under 100 lines
go test ./internal/harness/... -run 'TestRenderedContextFitsAntigravityBlock|TestRenderedContextUnderTarget' -v   # a change to context/*.md
  # or rules/*.md must pass this: it renders the real checkout plus every rules/*.md and pins the
  # result under 22,000 bytes, inside Antigravity's 24,000-byte global rules block, and under 15,000 bytes,
  # the harness's own ceiling (goal 14,500; plan audit-fix D2b) (REQUIREMENTS §
  # Known limits, ADR-035); already covered by `go test ./...` in § 1, named here explicitly
```

## 3. Scratch-home install and drift check

Never install into the operator's real `$HOME` from a dev session.

```bash
rm -rf .agents/scratch/<slug>-home
mkdir -p .agents/scratch/<slug>-home
go run ./cmd/sh-agent install --home .agents/scratch/<slug>-home
go run ./cmd/sh-agent status --home .agents/scratch/<slug>-home   # must report "drift none"
```

Confirm every `agents/*.md` and its derived Antigravity `agent.md` under
`.gemini/config/agents/<name>/agent.md` in the scratch home carry the same
`## Records` section (exactly one each) when an agent definition changed.

## 3a. Plain-copy install (change touching checkout-vs-mirror rebuild logic)

When a change affects `IsCheckout` / rebuild-decision logic in `cmd/sh-agent` or
`internal/harness/repo.go`, the scratch-home tier above only exercises the checkout
path. Also build the real binary and run it against a read-only, `.git`-less copy:

```bash
go build -o .agents/scratch/<slug>/sh-agent ./cmd/sh-agent
mkdir -p .agents/scratch/<slug>/mirror
for d in context skills agents commands rules; do
  mkdir -p ".agents/scratch/<slug>/mirror/$d"
  rsync -a --exclude .git "$d/" ".agents/scratch/<slug>/mirror/$d/"
done
chmod -R a-w .agents/scratch/<slug>/mirror
.agents/scratch/<slug>/sh-agent install --repo .agents/scratch/<slug>/mirror \
  --home .agents/scratch/<slug>/home2   # must say "plain copy" and not rebuild
ls -a .agents/scratch/<slug>/mirror     # must show no .agents/ created
.agents/scratch/<slug>/sh-agent status --repo .agents/scratch/<slug>/mirror \
  --home .agents/scratch/<slug>/home2   # first line "plain copy"; drift none
chmod -R u+w .agents/scratch/<slug>/mirror   # restore before cleanup
```

Then confirm the checkout path still rebuilds with
`<built binary> install --repo <this checkout> --home .agents/scratch/<slug>/home3 --dry-run`
→ `[dry-run] the binary would be rebuilt from <checkout>`.

## 3b. `report`, the feedback loop back (change to `internal/harness/report.go`, `cmd/sh-agent/main.go`'s report wiring, or `reorder`'s `takesValue` map)

Never run `report` without `--dry-run` from a dev or verification session — it pushes
to the real feedback repository (`git@github.com:kesokaj/styrman-harness-feedback.git`)
and never touch the real `$HOME`.

```bash
go run ./cmd/sh-agent report --dry-run --home .agents/scratch/<slug>-home
```

Must print the destination, the collected/absent files, and the redaction list (what was
caught — it cannot show a miss), and write nothing (`git status` in a real feedback checkout stays byte-identical before and
after). `go test ./cmd/... -run TestReorder` covers a new value-flag being added to
`takesValue`; `-run TestReport` alone does NOT run every report test (`TestSendReport`,
`TestResolveFeedbackFailures`, `TestSendReportRefusesRedirectsAndLinks`,
`TestSendReportPushesOnlyToFeedbackURL` don't match that prefix) — use
`go test ./internal/... -run 'TestReport|TestSendReport|TestResolveFeedback'` for
collection, redaction, and the git side, against local bare repositories under
`.agents/scratch/` or Go's own `t.TempDir()`, never the real remote, and confirm with
`go test ./internal/... -list 'TestReport|TestSendReport|TestResolveFeedback'` that it
lists every one of them.
A real end-to-end run against the feedback repository is an e2e-tier check the main
agent runs and logs, not a bare `verifier` dispatch.

## 3c. Verb × place-kind matrix (change to `cmd/sh-agent` or `internal/harness` that touches a verb's behavior or a refusal message)

`.agents/docs/testing/sh-agent-matrix.md` is the reproducible spec: every `sh-agent`
verb run against a checkout clone, a `git archive` plain copy, and a release-stamped
local build, each from a fresh scratch home, driven by `.agents/scratch/matrix/run.sh`
(gitignored; `git clone` of the checkout for the checkout kind, `go build ./cmd/sh-agent`
for the development binary both the checkout and plain-copy kinds run,
`go build -trimpath -ldflags "-X main.version=vX.Y.Y"` for the release-shaped binary,
`git archive HEAD | tar -x` for the plain copy). `report` runs only `--dry-run`; `update`
runs with a PATH holding no `gh`, so nothing reaches the network.

Run it, then update the matrix file's table and its "How the cells were run" line (tree
SHA and date) in the same change. Every cell is `exercised` (ran, output logged),
`refused` (a message naming the cause and the alternative), `n/a` (the verb has no
meaning there), or `remainder` (no machine here can run it) — none blank. Pass means: no
silent cell; every `refused` cell's message names both the cause and the alternative;
and `status` reports drift none from a fresh home on each of the three place kinds.

## 3d. Search index (change to `internal/search`, or to the `search`/`index` wiring in `cmd/sh-agent`)

```bash
go test -race ./internal/search ./cmd/sh-agent                 # unit tests: reconcile, parsing, queries, refusals, the global search index
SH_AGENT_RECALL=1 go test ./internal/search -run 'TestRecall$' -v   # recall measurement over testdata/recall.tsv: prints per-question recall and the means; not a gate
go build -o .agents/scratch/sh-agent ./cmd/sh-agent
.agents/scratch/sh-agent index --rebuild --home .agents/scratch/s   # this checkout: under 10 s
.agents/scratch/sh-agent search --home .agents/scratch/s <words>    # warm, nothing changed: under 100 ms end to end
```

A new or changed verb also gets its rows in `sh-agent-matrix.md` (section 3c). A test written after its code is checked by the mutation protocol (`review-standards` § 2): the mutated code compiles and a named test fails.

## 3e. Memory, dream and the timer (change to `internal/memory`, `internal/dream`, or the `memory`/`dream` wiring in `cmd/sh-agent`)

```bash
go test -race ./internal/memory ./internal/dream ./internal/search ./cmd/sh-agent
```

- The unit tests put fake `claude`, `agy`, `systemctl` and `launchctl` on PATH; nothing calls a real CLI or manager.
- A real `sh-agent dream` run in development goes against a scratch `--home` (under `.agents/scratch/`) whose search.db indexes only fixture projects — never `--all` against the real home: a real run sends redacted project content and transcripts to a model provider.
- The timer is verified by `systemd-analyze verify` on units rendered through a fake manager. The real `--install-timer` is the operator's step (R6).
- Corpus measurement of the Antigravity reader (a change to `internal/search/dream.go`): `SH_AGENT_CORPUS=<dir of copied Antigravity conversation databases> go test ./internal/search -run TestAgyCorpusLeaves -count=1 -v` prints `corpus: N databases, L prose leaves kept, K with a tool-argument key`. K must be 0. The 2026-09-30 measurement over 95 databases was 971 and 0, before and after the reader change. The copies live under `.agents/scratch/` and are deleted after.
- The kit is in `testing/2026-09-30-verify-parity-memory-dream-phase-4.md`.

## 4. Sandbox start of each CLI (change to `context/`, `agents/`, or `skills/`)

One start of Claude Code and one of Antigravity in the sandbox (paths and flags in
`.agents/docs/testing/2026-09-03-verify-sandbox.md`), each asked a context question
that quotes the new wording, log grepped for parse errors. The full skill-matrix
scenarios are `.agents/docs/testing/skill-matrix.md` — run the matrix in full only
before saying "complete" or "all skills"; a targeted change runs only the matrix row(s)
it affects.
A sandbox home that must be logged in is built under `.agents/scratch/` (gitignored, never
embedded: `go:embed` covers only the five source directories); any credential file copied into
it from the real home is `chmod 600`, lives only for the run, and is deleted by name right
after (security review 2026-09-21). Prefer the real-home smoke checks of § 5 when the change
needs no separate home.

A change that must hold on each kind of machine (R10: Claude Code only, Antigravity only, both)
is checked in three sandbox homes, one per `sh-agent install --home <sb> --vendor claude|agy|all`,
with `HOME=<sb>` for the whole run so every `sh-agent` call the agent makes lands there (BUG-107).
Claude Code runs with `CLAUDE_CONFIG_DIR=<sb>/.claude` and a copy of `~/.claude/.credentials.json`
in it; Antigravity with its token and a settings file written for the run (§ Behavioural runs below).
A vendor setting the session cannot write, such as a project's `.claude/settings.local.json`,
reaches Claude Code through `claude -p --settings '<json>'` instead. Kit and first use:
`testing/2026-09-30-verify-parity-memory-dream-phase-2.md`.

### Behavioural runs (a change to `agents/*.md` or to a `context/` rule that changes what the main agent does)

Moved here from `AGENTS.md` § Testing on 2026-10-01 (plan audit-fix, C17); the project file points at this section.

- A change to `agents/*.md` gets a behavioural run: `claude -p --agents <json built from the working-tree file>` in a scratch clone of a small repository, several runs, with the agent's artifact carrying something only the new definition writes, so the run proves which definition ran. The dispatch text does NOT repeat the changed rule (the definition alone is under test); the prompt goes on stdin because `--allowedTools` is variadic and swallows a trailing argument; `--allowedTools` names the agent's tools or every call is refused at the permission gate and the run proves nothing; the tool-call order is read from the run's subagent transcript under `~/.claude/projects/`; the result is reported as n of N runs, never "works" (BUG-083). A change to `context/` that should change what the main agent does runs the same way with the new text given through `claude -p --append-system-prompt`.
- The real home is never touched. A CLI started from a scratch `HOME` cannot authenticate (credentials live in the real home), so the CLI start of § 4 is the sandbox procedure or the real-home smoke checks of § 5, never `HOME=<scratch> claude -p`. A verification log states what it did not exercise; "complete" or "all skills" is said only after the full matrix in `skill-matrix.md` has run.
- **Every kit sources `kit-preamble.sh` first** (`. .agents/docs/testing/kit-preamble.sh` at the top of the kit script built in scratch, with `KIT_CLONE`, `KIT_HOME` and, when the dispatch states fixture facts, `KIT_ASSERT_GREP` (newline-separated `file:pattern`) and `KIT_ASSERT_COMMIT` (one short SHA)). The preamble is the setup proof BUG-087, 090, 091, 106, 107, 115, 116 and 123 each asked for: it asserts every fixture fact before the CLI starts, refuses a clone whose parent is a git checkout (an ancestor `AGENTS.md` would load and an agent resolved the docs root upward into this repository's trail, BUG-106), lists the real `~/.config/sh-agent` before the run and again on exit, and FAILS the kit (exit 1) when it changed (a kit's own `--home` does not reach the agent's calls: `search` registered test clones in the real global index, BUG-107, and a `/distill` probe appended to the real feedback and retro logs, BUG-123), and puts a wrapper first on `PATH` that pins `--home "$KIT_HOME"` on every `sh-agent` call the agent makes.
- A sandbox home for headless agy gets its own `~/.gemini/antigravity-cli/settings.json`, written for the run with the fewest grants the run needs (a `command(<name>)` grant allows every argument of that command) and the scratch clone as its only trusted workspace, and deleted by name with the copied token — never a copy of the real home's file, which approves every tool and trusts real checkouts; a run counts only when its stderr has no "auto-denied" line and its `cli.log` no "invalid allow entry" line.
- The must-fail control for a change to `context/` runs in a project that does not contain the harness source, never in a clone of this repository, whose own `context/` carries the new text.
- A headless probe of a working-tree skill passes `--setting-sources project,local`, or first shows the real home has no installed skill of that name: a user-level skill shadows the project link (plan `2026-09-29-dreamer-memory-loop.md` D15).

## 5. Real-home smoke checks (after any install into a real home)

The three checks in `.agents/docs/testing/2026-09-03-verify-real-home.md` (context
questions — including the canary word GANNET that ends the rendered context, which a
truncated file fails (BUG-085) — one skill load, one agent dispatch) from that home, on each
CLI, before the install is called done. The Mac procedure is
`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`.
Run them from a fresh clone of the demo service, never from inside a harness checkout: the
canary word is written in this repo's README and AGENTS.md, so a session started here could
answer it from the project context without the global context having loaded. On Antigravity
the answer counts only with the store check beside it (`strings` the newest
`~/.gemini/antigravity-cli/conversations/*.db`: one `<RULE[...AGENTS.md]>` block, the canary
line inside it, no `<truncated` line). A headless Antigravity check that needs tools uses the
real settings' `permissions.allow` entries; `--dangerously-skip-permissions` only inside a
throwaway clone, and never as the standing way.

After any `sh-agent update`, also run the deploy guides' post-update check from that
home: `hash -r` (bash) or `rehash` (zsh), then `sh-agent version` unqualified and
confirm it matches the version just installed, then `sh-agent status` and confirm it
prints no `shadowed` line (BUG-079). When the version is wrong, `type sh-agent` then
`which -a sh-agent` name which file the shell is still reaching.

The by-hand install chains (README § Install without cloning, `remote.md` § First
install) are exercised only in `.agents/docs/testing/2026-09-17-release-binary-e2e.md`,
not by any tier in this guide — an accepted e2e-only check.

## Full matrix

`.agents/docs/testing/skill-matrix.md` is the reproducible spec of every skill/agent
sandbox scenario. Cite it rather than duplicating rows here; add a row there when a
skill or agent gains new behavior worth an end-to-end check.
