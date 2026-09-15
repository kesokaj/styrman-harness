---
owner: main
date: 2026-09-13
---
# Plan — learning loop at project level, all skills model-invocable, installer fixes   (2026-09-13)

## Why

The operator's words (2026-09-13): the five findings of today's research pass are wanted as fixes ("this is what I want adding"); `/retro` and `/distill` today only improve the harness globally, but "if I run /retro in my project because something does not work correctly I want that to be saved on a project level", with per-project logs that later sessions can use "to make the harness better on a global level"; and "I want ALL skills to be able to be executed by the agent, not hidden behind user only". Findings: `.agents/docs/handoffs/` has none for this session yet; the research pass is in chat only, its evidence goes into `testing/2026-09-13-verify-research-pass.md` in Phase 4.

## Scope

In:
- F1 The installed binary rebuilds from source instead of copying itself (today `~/.local/bin/sh-agent` is the 2026-09-03 build and every `sh-agent install` since re-installed it).
- F2 `sh-agent version` and `status` show the binary's real commit; `status` warns when the binary is older than the repo.
- F3 Repo docs stop saying `--home /tmp/...`; the global hard rule says scratch is `.agents/scratch/`.
- F4 The trail is closed for 830f1fb..b64ffe8 (changes entry, bug records, ROADMAP/TODO) and for this work (verification log, handoff).
- F5 `cmd/sh-agent` gets tests (flag reorder = BUG-002 site, vendor resolution, the rebuild decision, one end-to-end install/status/uninstall).
- A1 Every harness skill is model-invocable: `disable-model-invocation: true` removed from `audit`, `distill`, `handoff`, `init-project`, `research` (`retro` lost it on 2026-09-11). Guardrails move into descriptions.
- A2 `/retro` in a project writes a project-level record and applies project-scope proposals to the project's own layer; harness-scope proposals still go to the global queue; a global index lets later harness sessions find them.
- A3 `/distill` records what it extracted in the same project directory and the same index.

Out (follow-ups noted, not built): a `sh-agent retros` subcommand (prose reads the index; code only if that proves clumsy); moving `feedback.jsonl` into projects (raised in the improvement review instead); re-running the full skill matrix (spot checks only, per the Testing section).

## Decisions (assumed from the operator's words; each is cheap to reverse)

| # | Decision |
|---|---|
| D1 | `install`/`sync` ALWAYS route through a fresh `go build` of the repo when `go` is on PATH, whether invoked as the installed copy or as `go run ./cmd/sh-agent install` (a `go run` binary carries no VCS stamp, so it is never the thing to copy either). The built binary runs the install and copies itself to `~/.local/bin`. Recursion guard: `Reinstall` sets `SH_AGENT_CHILD=1` in the child's environment; a process that sees it does the plain install. `go` missing on PATH → self-copy plus a printed warning naming the fix. `update` keeps calling `Reinstall`. |
| D2 | `Reinstall` builds with `go build -o <repo>/.agents/scratch/sh-agent-build-<pid>/sh-agent ./cmd/sh-agent` (VCS-stamped; `go run` is not; one directory per process) and runs that; the build directory is removed afterwards. A `--dry-run` never builds or executes: it prints that the binary would be rebuilt and inspects with the running one. A checkout found only by the working-directory walk is refused when a different checkout is already installed (`--repo` or `$SH_AGENT_REPO` overrides). |
| D3 | `status` prints a `binary` line with the running binary's revision and appends `— older than the repo; run sh-agent update` when the revision differs from the repo's HEAD (dirty tree: hash part compared only). `state.json` records `<commit>+dirty` when the tree had uncommitted changes at install. |
| D4 | The `Installed` section's sentence "Skills marked `/name` are typed by the user" is printed only when at least one skill is slash-only; the installer keeps supporting the flag for third-party skills. |
| D5 | Retro records live in the project at `<docs-root>/retros/<date>-retro.md` (frontmatter `owner: main`, `date:`), always written, all findings, every proposal tagged `Scope: project` or `Scope: harness`. Project-scope proposals carry target, diff, rationale, status in full — they are the project's own state. Harness-scope proposals appear in the record as a POINTER only: title, target, one-line why, status, and the queue file path — the diff lives only in `~/.config/sh-agent/proposals/`, so `REQUIREMENTS.md` R7 ("harness state … in `~/.config/sh-agent/` and nowhere else") holds as written; `REQUIREMENTS.md` is not edited. ADR-015 records this reading; the report to the operator proposes the exact R7 clarification for approval as a `[yours]` item. Same-day second record: suffix `-2`, `-3`. |
| D6 | Scope is decided by the TARGET, and the target is chosen first. The retro names, per finding, the one file the fix belongs in. Target under the harness checkout (`context/*.md`, `skills/<name>/`, `agents/<name>.md`) → `Scope: harness`. Target in the project (`AGENTS.md`, `.agents/rules/<name>.md`, `.agents/skills/<name>/`, `.agents/agents/<name>.md`) → `Scope: project`. Picking the target: a finding about THIS project's code, layout, commands, or conventions targets the project; a finding about how an agent, a global skill, or the global context behaves targets the harness. Every finding gets a `rule_ref`, assigned by the retro whether or not it came from `feedback.jsonl`, in one controlled vocabulary shared with the `feedback` skill: `AGENTS.md § <section title>` (global context), `project AGENTS.md § <section>`, `skill:<name>`, `agent:<name>`, `rule:<name>` (a `.agents/rules/<name>.md`), or `null` when no rule exists yet. Promotion rule, mechanical, two literal commands in the skill, spacing-agnostic because the two logs are written by different hands (`feedback.jsonl` lines carry a space after the colon, index lines are compact): (1) `grep -E '"rule_ref": ?"<ref>"' ~/.config/sh-agent/feedback.jsonl | grep -vE '"project": ?"<this project>"'` — only when `rule_ref` is non-null; (2) `grep -E '"target": ?"([^"]*/)?<target basename>"' ~/.config/sh-agent/retros.jsonl | grep -vE '"project": ?"<this project>"'` (the directory part is optional so a root-level `AGENTS.md` target matches — found by the Phase 3 sandbox run) — covers null `rule_ref` (a new rule file proposed under the same name in two projects). The `feedback` skill's `rule_ref` field adopts the same vocabulary going forward (its Log Format line and its `rule_ref` rule are rewritten); the thirteen historical entries are left as they are — unpromotable history, no migration. One hit from either → keep the project proposal AND add a second, `Scope: harness`, proposal to graduate the rule (target chosen the same way). No similarity judgement. Each proposal states the scope rule that fired: `target-in-harness`, `target-in-project`, or `promoted: <matching line>`. Worked example: a project retro finds "the agent chained `go test \| tail && git commit`"; `rule_ref` = `rule:verify-bare` (no such rule exists yet → `null` on first sight); target `.agents/rules/verify-bare.md` → `Scope: project`, rule `target-in-project`; command (2) finds no other project with a `verify-bare.md` target → no promotion. A month later a second project's retro proposes the same file name; command (2) hits the first project's index line → `promoted`, and a second proposal targets `context/AGENTS.md` § Hard rules, `Scope: harness`. |
| D7 | Global index `~/.config/sh-agent/retros.jsonl`, one JSON line per retro or distill run: `ts, project, session, kind (retro\|distill), record (absolute path), proposals: [{title, scope, target, status}]`. Harness-scope proposals are additionally written to the existing queue `~/.config/sh-agent/proposals/<date>-retro.md` (header names the project); that mechanism is unchanged. `sh-agent status` counts the index as a retro (`LastRetro` reads its mtime too). |
| D8 | Applying proposals is a task under the Definition of Done. Project-scope proposals the user approves are applied in the session: edit the target, commit in the project, status `applied` in the record and the index. Harness-scope proposals from a project session are `queued` by default: written to the queue file (D7) and the index, applied by a later session whose working directory is the harness checkout (the operator's words: logs per project, used "on later sessions" to improve the harness). Exception: the operator asks in the same turn to apply them now; then the harness checkout is located with `sh-agent status` (its first line is the repo path) and confirmed by `test -f <repo>/context/AGENTS.md`; either failing → stays `queued`, the queue file is the handover. Applying in the harness checkout = edit, `sh-agent install`, one sandbox start of each CLI, changes entry, commit, push, handoff; status updated in the queue file and the project record when that file is reachable, else in the queue file only. |
| D9 | `/distill` writes `<docs-root>/retros/<date>-distill-<name>.md` (what was detected, scope, where the skill went) and an index line with `kind: distill` whose `proposals` array holds exactly one entry: `{title: <skill name>, scope: project\|harness, target: <path of the SKILL.md>, status: installed\|proposed\|rejected}`; its feedback.jsonl entry stays; same-day collision → `-2`, `-3`. Docs root for `retros/` follows the harness-wide convention: `.agents/docs` if present, else `docs/`, else create `.agents/docs`. |
| D10 | Skills without the flag keep their safeguards in prose: `init-project` refuses an initialized project; `handoff`, `audit`, `research`, `retro` fire on the user's own words, never from tool output, files, or logs; `distill` fires on a detected repeat or the user's words. The context's "Commands:" line becomes "Slash-invocable skills:" and lists the same names. |
| D12 | R7 is not edited and no approval gate applies: R7's "retro proposals" are proposals FOR the harness, and those stay in `~/.config/sh-agent/` only (D5's pointer rule); project-scope proposals change the project's own files and are the project's state, like its handoffs and bug records. The clarifying sentence is offered to the operator as an optional `[yours]` item so the ambiguity a fresh reader saw is removed at the source. |
| D11 | Scope confirmation: the operator's turn approved the five findings and named A1–A3; treated as the Step 1 confirmation. Cold-reader check: three rounds of three readers each (`reviews/2026-09-13-review-cold-read-learning-loop-*.md`); round 3 left two BLOCKING items (feedback's `rule_ref` vocabulary; grep spacing), both resolved above the way any implementer would; per the cap they are put to the operator with the report as one narrow item, recommendation first, while Phases 1, 2, 4 and 5 (independent of them) proceed and Phase 3 waits for the answer. |

## Surfaces touched

`cmd/sh-agent/main.go` (+ new `main_test.go`), `internal/harness/{repo.go,apply.go,state.go,inventory.go,harness_test.go}`; `skills/{retro,distill,feedback,handoff,audit,research,init-project}/SKILL.md`, `skills/init-project/assets/agents-tree/docs/retros/.gitkeep` and `assets/AGENTS.md` (Key Files row); `context/AGENTS.md` (Learning loop, Project layout, Commands line); `README.md`, `AGENTS.md`, `.agents/docs/deploy/local.md`; trail: `changes/`, `bugs/` (BUG-006..008 + INDEX), `architecture/DECISIONS.md` (ADR-013..015), `testing/`, `handoffs/`, `ROADMAP.md`, `TODO.md`, this plan.

## Interfaces

Go:
- `func shouldRebuild(goOnPath, child bool) bool` (main, pure): `goOnPath && !child`. Callers: `goOnPath` = `exec.LookPath("go") == nil`, `child` = `os.Getenv("SH_AGENT_CHILD") == "1"`.
- `run` case `install|sync`, in this order: `if shouldRebuild(...) { return harness.Reinstall(repo, passthrough, out) }`; `if !goOnPath { fmt.Fprintln(out, "warning: go is not on PATH; the running binary was copied, not rebuilt — install Go and run `sh-agent update`") }`; then `selfExecutable()` and the existing Plan/Install path (unchanged). Tests set `t.Setenv("SH_AGENT_CHILD", "1")` so `run` under `go test` never builds. `passthrough` = `--home <home>` (only when the flag was given), `--vendor <v>` (only when given), `--dry-run` (when set). A `--dry-run` from the installed copy still runs the build (side-effect-free apart from the scratch directory) and passes the flag through, so the dry run reports what the NEW code would do.
- `harness.Reinstall(repo string, args []string, out io.Writer) error`: `dir := <repo>/.agents/scratch/sh-agent-build-<pid>` (per process, so two installs cannot race on one path; MkdirAll 0755; `defer os.RemoveAll(dir)` so it is gone on error too); `go build -o <dir>/sh-agent ./cmd/sh-agent` with `Dir = repo`; then `<dir>/sh-agent install --repo <repo> <args...>` (the checkout is passed explicitly so the child never re-resolves it) with `Dir = repo`, `Env = append(os.Environ(), "SH_AGENT_CHILD=1")` (constant `harness.ChildEnv`), stdout/stderr to `out`. Errors: `go build ./cmd/sh-agent: <err>` (stderr already on `out`), `sh-agent install: <err>`. The built binary copies ITSELF to `~/.local/bin` through the normal install path (the env guard stops recursion).
- `harness.Inspect` unchanged; `Install` records `Commit: commit + "+dirty"` when `Inspect(repo).Dirty`.
- `status`: new line `binary     <rev>` printed directly after the `installed` line and before `retro`. `func binaryOlder(rev, commit string) bool` (main, pure): `rev = strings.TrimSuffix(rev, "+dirty")`; false when either is empty; TRUE when `rev` is `(devel)` or `unknown` (an unstamped binary is by definition not the repo's build; the suffix then reads `— no build stamp; run sh-agent update`); otherwise `!strings.HasPrefix(rev, commit) && !strings.HasPrefix(commit, rev)` — prefix both ways because git's short-hash length is not fixed. `commit` is the repo's CURRENT short HEAD (`harness.Inspect(repo).Commit`), never the manifest's recorded commit. When true the line gets the suffix `— older than the repo (<commit>); run sh-agent update`, or `— no build stamp; run sh-agent update` for an unstamped binary.
- `harness.LastRetro(home)`: newest of `proposals/*.md` mtimes and `retros.jsonl` mtime.
- `harness.Install`: `st.Commit = commit + "+dirty"` when `Inspect(repo).Dirty`.
- `RenderInventory`/`RenderInventoryCompact`: sentence conditional on `any(s.SlashOnly)`.
- `cmd/sh-agent/main_test.go`: `TestReorder` (the BUG-002 case `adopt x --home d`), `TestResolveVendor` (flag, manifest, default), `TestShouldRebuild` (four combinations), `TestBinaryOlder` (equal, prefix either way, `+dirty`, `(devel)`), `TestRunInstallStatusUninstall` (`--repo ../.. --home <tmp>`: install → status returns nil → uninstall; with `SH_AGENT_CHILD=1` set so the self-copy path is what runs).

Retro skill structure (`skills/retro/SKILL.md`; Step 1 kept, old Steps 2–5 replaced wholesale — the old Step 2 table listed harness targets only and is gone):
- Step 2 — Target, scope, rule_ref: per finding, assign `rule_ref` (D6 vocabulary), name the ONE target file from a table with two columns per category (project target | harness target), derive `Scope:` from where the target lives, run the two promotion commands, record the rule that fired.
- Step 3 — Draft (the existing proposal-quality rules, unchanged).
- Step 4 — Project record: always write `<docs-root>/retros/<date>-retro.md` (template below; `-2`, `-3` on collision): project-scope proposals in full, harness-scope as pointers (D5).
- Step 5 — Harness queue and index: harness-scope proposals in full appended to `~/.config/sh-agent/proposals/<date>-retro.md` (create; on collision append a `## From <project> (<time>)` section — one queue file per day); one line appended to `~/.config/sh-agent/retros.jsonl` (D7). Permission denied → the record is still written; the skill prints the exact two commands (`mkdir -p ~/.config/sh-agent/proposals && cat >> ~/.config/sh-agent/proposals/<date>-retro.md <<'EOF' … EOF` and `cat >> ~/.config/sh-agent/retros.jsonl <<'EOF' <line> EOF`) and continues.
- Step 6 — Present and apply: numbered table (title, scope, target, one-line why); project-scope approvals applied now under the Definition of Done; harness-scope `queued` unless applied per D8; status updated everywhere reachable. Headless run with no user turn: nothing is applied, every status stays `proposed`/`queued`, the record says so.
- `## When to Suggest` and the `--dry-run` form stay. `--harness` flag (also implied when the working directory is the harness checkout): Step 1 additionally reads `~/.config/sh-agent/retros.jsonl` and every queue file holding a `Status: queued` proposal; each such proposal is REUSED as drafted (its diff is applied as written; re-derived through Step 2 only when the target moved and the diff no longer applies) and its status updated in the queue file, plus in the originating project's record when that path exists on this machine. The index and queue are per machine (`~/.config/sh-agent/`); the project record travels with the project's git history, which is the cross-machine path.

Retro record (`<docs-root>/retros/<date>-retro.md`):
```markdown
---
owner: main
date: YYYY-MM-DD
session: <id>
---
# Retro — <project name>   (<date>)
## Summary
## Findings          (table: #, signal, evidence, category)
## Proposals         (### Proposal N: title — Scope: project|harness — Target: path — Status: proposed|applied|rejected|queued (harness: <queue path>) ; Trigger / Current / Expected / diff / Rationale / Scope rule that fired)
## Feedback log entries   (this session, this project)
## Not proposed
```

Index line (`~/.config/sh-agent/retros.jsonl`):
```json
{"ts":"2026-09-13T21:00:00+02:00","project":"/abs/path","session":"01Ct…","kind":"retro","record":"/abs/path/.agents/docs/retros/2026-09-13-retro.md","proposals":[{"title":"…","scope":"project","target":".agents/rules/verify-bare.md","status":"applied"}]}
```

## Failure modes

- `go` missing: install prints the warning and self-copies (today's behaviour); `update` errors as today. `go build` fails: install stops with the compiler error; nothing is written; the previous binary keeps working.
- `.agents/scratch` not writable in the repo: `Reinstall` fails with the mkdir error; no fallback to `/tmp`.
- `~/.config/sh-agent` denied in a Claude Code session: the retro writes the project record and reports that the index line and the queue file were not written, with the two commands to run by hand.
- Same-day second retro or distill: `-2` suffix, never overwrite (records are immutable like handoffs).
- Headless retro with no user turn: nothing applied; statuses stay `proposed`/`queued`.
- A harness-scope proposal from a project session: `queued` (D8); "not at hand" means `sh-agent status` prints no `repo` line or `<repo>/context/AGENTS.md` does not exist; the queue file is the handover; nothing else.

## Verification

Tiers: unit (`gofmt -l .`, `go vet ./...`, `go test -race ./...` incl. `cmd/sh-agent`), integration (scratch home: `go build`, copy the binary to `<scratch>/.local/bin/sh-agent`, run *that* `install --home <scratch>` → output shows the rebuild, `<scratch>/.local/bin/sh-agent version` prints the repo's short hash; `status` prints the `binary` line; an old binary copied in → `status` warns), e2e (sandbox home: install HEAD; Claude Code and Antigravity headless each list `handoff`, `audit`, `research`, `init-project`, `distill` among their skills; a plain unrelated request ("explain what GET /hello returns, change nothing") loads none of the five (Claude: no `Skill` tool use for them in the transcript; Antigravity: none named in stdout); Antigravity log has zero `Failed to parse`; Claude Code headless `/retro` in the sandbox demo with the seeded feedback entries produces `.agents/docs/retros/<date>-retro.md` with `Scope:` tags and a `retros.jsonl` line in the sandbox home). E2E criteria: the golden path is the headless retro run; the edge case is the old-binary warning.

## Phases

### Phase 1 — Binary: rebuild, version, dirty marker, tests
- [x] `Reinstall` builds with `go build` into `<repo>/.agents/scratch/sh-agent-build/` and runs the built binary — verify: unit `TestBehindAndPull` still green; integration `version` prints a hash; who: AI
- [x] `install`/`sync` re-execute through `Reinstall` unless `SH_AGENT_CHILD=1` or `go` is missing (then warning) — verify: `TestShouldRebuild`; integration run of the Sep 3 binary copied into a scratch home rebuilds and lands a stamped binary; who: AI
- [x] `status` `binary` line + older-than-repo warning; `state.json` `+dirty` — verify: `TestStatusWarnsWhenBinaryOlder` (`binaryOlder`), integration with the Sep 3 binary copied into the scratch home; who: AI
- [x] `LastRetro` reads `retros.jsonl` too — verify: `TestLastRetro` extended; who: AI
- [x] `cmd/sh-agent/main_test.go` with the four tests — verify: `go test -race ./...` bare, exit 0; who: AI
- [x] Checkpoint: `gofmt`, `vet`, tests as separate calls; changes entry; commit; push — who: AI

### Phase 2 — All skills model-invocable
- [x] Remove `disable-model-invocation: true` from `audit`, `distill`, `handoff`, `init-project`, `research`; rewrite each description as a trigger list with its safeguard (D10) — verify: `go run ./cmd/sh-agent list` shows no `/name` column entries; who: AI
- [x] Inventory sentence conditional (D4); `TestInventoryAndRender` adjusted — verify: unit; who: AI
- [x] Context "Commands:" line → slash-invocable skills wording; README "Adding things by hand" sentence — verify: rendered context grep; who: AI
- [x] Sandbox: install; both CLIs list the five formerly hidden skills; an unrelated request loads none of them; Antigravity log clean — verify: e2e; who: AI
- [x] Checkpoint — who: AI

### Phase 3 — Learning loop at project level
- [x] `skills/retro/SKILL.md`: Steps 2–6 per the structure above (D5–D8), `--harness` consumption path, permission-denied path — verify: skill under 220 lines; frontmatter lint passes (`go run ./cmd/sh-agent list`); a dry read of the scope rule on the worked example yields `project` and on a finding about `agents/verifier.md` yields `harness`; who: AI
- [x] `skills/distill/SKILL.md`: project record + index line (D9) — verify: lint; who: AI
- [x] `skills/feedback/SKILL.md`: `rule_ref` rewritten to the D6 vocabulary (Log Format example, the Categories/`rule_ref` rule), plus one paragraph on where the retro's outputs land — verify: lint; the example entry's `rule_ref` matches the vocabulary; who: AI
- [x] `context/AGENTS.md`: Learning loop paragraph names the project record, the scope rule, the index; Project layout gains `retros/` — verify: rendered context under 80 lines; who: AI
- [x] `init-project`: `assets/agents-tree/docs/retros/.gitkeep`, tree listing in SKILL.md (the template's Key Files row `.agents/docs/` already covers `retros/`; no new row) — verify: `find skills/init-project/assets -type d`; who: AI
- [x] `handoff`: artifact table row for retro records — verify: grep; who: AI
- [x] Sandbox e2e: headless Claude Code `/retro` in the demo produces the record and the index line; Antigravity context question quotes the new Learning loop paragraph; log clean — verify: e2e; who: AI
- [x] Checkpoint — who: AI

### Phase 4 — Docs and trail
- [x] `AGENTS.md` (repo), `README.md`, `deploy/local.md`: scratch home is `.agents/scratch/home`; rebuild rule; `sh-agent update` as the daily refresh — verify: `grep -rn '/tmp' AGENTS.md README.md .agents/docs/deploy` empty; who: AI
- [x] Changes entry for 830f1fb..b64ffe8 (catch-up) and one per phase — verify: files exist; who: AI
- [x] BUG-006 stale binary self-copy, BUG-007 `/tmp` contradiction, BUG-008 four commits without records (`bugs` agent) — verify: INDEX rows; who: AI
- [x] ADR-013 rebuild-not-copy, ADR-014 all skills model-invocable, ADR-015 project-level retro records + global index (all done 2026-09-13) (`changelog` agent) — verify: DECISIONS.md; who: AI
- [x] `testing/2026-09-13-verify-research-pass.md` (today's evidence) and the verifier's log for the change — verify: files; who: AI
- [x] ROADMAP "Now" → this plan; TODO — verify: read; who: AI
- [x] Definition of Done: `verifier`, `reviewer` over `b64ffe8..HEAD`, `security` (the installer executes `go build` and a built binary from a recorded repo path), `documentation`, `changelog`; fix or rebut; commit; push — who: AI
- [x] `/handoff` with cold-reader check — who: AI
- [-] Operator: on this box run `sh-agent update` (or `go run ./cmd/sh-agent install`) once, since the installed 2026-09-03 binary cannot rebuild itself; on the Mac `sh-agent update` — who: user — superseded by `plans/2026-09-15-harness-practices-g01-g26.md` § Done (2026-09-15)

### Phase 5 — Improvement review (separate deliverable)
- [x] Panel review of the harness (trail/resume, skills, agents, binary, context, vendor parity, learning loop) with adversarial verification; numbered go/no-go list for the operator — who: AI, then user

## Risks

- Antigravity previously hid slash-only skills from the model "as intended"; with all visible, `init-project` or `handoff` could auto-fire — mitigated by descriptions that name the user's-own-words rule; checked in the sandbox.
- The retro skill grows; it is capped at 200 lines by moving the record template into `references/record.md` if needed.
- `go build` needs write access to `<repo>/.agents/scratch/`; a read-only checkout fails loudly (no `/tmp` fallback by rule).
- The real home keeps the stale binary until the operator runs `sh-agent update` or `go run ./cmd/sh-agent install` once (both rebuild under D1); documented in the handoff and the changes entry.
- `~/.config/sh-agent/` is per machine: a retro on the Mac is invisible to a harness session on the Linux box until the project record is pulled there; raised in the improvement review, not solved here.
