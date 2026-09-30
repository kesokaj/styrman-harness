---
owner: main
date: 2026-09-30
---
# `sh-agent dream` — the background dreamer

Status: built 2026-09-30 (plan `plans/2026-09-30-parity-memory-dream.md`). R3's named exception, R4's dream clause, R10 (ADR-047, ADR-047a). The memory half of the learning loop is `specs/sh-agent/memory-loop.md`.

## What it does

`sh-agent dream --all`, run nightly by a user timer (`--install-timer`), or by hand with `--all` or `--project <dir>`. It looks at every project in the global search index, `~/.config/sh-agent/search.db` `projects` (read from a copy, never opened live). A project is skipped with a log line when its root is gone, it is not a git checkout, or it lacks its docs root. Otherwise, when anything changed since the project's watermark:

1. **Collect** (`internal/dream/collect.go`), in prompt order:
   - the memory store `<root>/.agents/memory/`, whole, when the trust checks pass (otherwise it is named as left out);
   - `feedback.jsonl` lines for this root;
   - `bugs/`, `reviews/`, `retros/`, `changes/`, `handoffs/` records;
   - the installed skills' names and descriptions (global: never a change on their own);
   - Claude Code transcript lines whose `cwd` is the root or under it, text parts only (`CLAUDE_CONFIG_DIR` honoured);
   - Antigravity conversations whose `workspace_uris` names the root, by `last_modified_time` (whole conversations, from copies of the databases; `internal/search/dream.go`).

   Everything but memory and skills is taken only when newer than the watermark.
2. **Scrub**: every input passes `harness.Redact` (report's rules) and the home path becomes `~`.
3. **Cap**: at 256 KB. Transcripts are dropped oldest first, then trail, feedback, skills. A project whose memory alone passes the cap fails. The prompt names what was left out.
4. **Ask**: one CLI call with NO tools, from the empty `~/.config/sh-agent/dream-work/`, 10-minute timeout, the child keeping sh-agent's own environment.
   - Claude Code (manifest vendor `claude` or `all`): `claude -p --tools "" --no-session-persistence --output-format text`, the prompt on stdin.
   - Antigravity (`antigravity`): `agy -p <prompt> --mode plan --output-format text` (agy has no no-tools flag: plan mode and the prompt's opening warning, that every tool call is refused and ends the run with no answer, keep it from calling one; D7b).

   `--vendor` overrides the manifest.
5. **Write**: the first `{…}` of the answer that decodes as `{"proposals":[…],"skills":[…]}`, scrubbed again.
   - `proposals/dream/<date>-<project>.md`: project targets, each with `**Status:** queued`, plus the draft list.
   - `proposals/dream/harness/<date>-<project>.md`: harness targets.
   - `proposals/dream/drafts/<date>-<project>-<name>/SKILL.md`: each draft skill with a kebab-case name; an invalid name is listed as dropped.

   Same-day reruns get `-2`, `-3`. Zero proposals write no file.
6. **Advance** the watermark to the run's start, only after a successful parse and write.

## Consumers

- Project-level `/retro` (source 9) reads its project's `proposals/dream/*.md`, re-derives each proposal against the files, and sets its Status to applied or rejected.
- `/retro --harness` reads `proposals/dream/harness/*.md` beside the harness queue.
- `sh-agent report` sends `proposals/dream/harness/` (R8), never the project files, drafts, state, log or failures.
- `sh-agent status` prints `dream      last run <date>: <n> projects, <m> proposals; timer installed|not installed` (or `failed: <reason>`, or `never run`).

## State and files (all under `~/.config/sh-agent/`, 0600 in 0700 directories)

| Path | What |
|---|---|
| `dream-state.json` | `projects` (root → watermark) and `last_run` (at, projects, proposals, error) |
| `dream.log` | one line per project (`ok N proposals`, `unchanged`, `skip <why>`, `fail <why>`) and one per run |
| `dream.lock` | the running pid; a lock whose pid is gone is replaced, and a live one makes a second run exit 0 "already running" |
| `dream-work/` | the CLI's working directory, kept empty; not a project, so the run's own conversation never maps to one |
| `dream-failed/` | a failed project's stdout and stderr, scrubbed; files older than 30 days are pruned at the start of a run |
| `proposals/dream/` | the output above |

## Failure modes

| Situation | Behaviour |
|---|---|
| No CLI on PATH | exit 2 before any project; `last_run.error` and a run log line; `status` says `failed` |
| The CLI exits non-zero or times out | the project fails, its output goes to `dream-failed/`, its watermark stays; the others run; exit 1 |
| The answer has no decodable JSON | the same, logged `fail invalid output` |
| The store is a symlink, or git tracks it | memory is left out and named in the prompt and the proposal frontmatter; the rest is sent |
| Memory alone over 256 KB | the project fails with `memory over cap` |
| No indexed project | nothing to do; exit 0 |
| A second run while one runs | exit 0, "already running" |

## Verification

Unit tests use a fake CLI on PATH that records its arguments, working directory and stdin (`internal/dream/run_test.go`), so the tool-free flags, the cwd, the absence of a planted secret in the prompt, the file shapes and the watermark rule are pinned. The D7 probe (plan) ran both real CLIs tool-free from an untrusted, empty directory. The end-to-end rows are in `testing/skill-matrix.md` and `testing/sh-agent-matrix.md`.
