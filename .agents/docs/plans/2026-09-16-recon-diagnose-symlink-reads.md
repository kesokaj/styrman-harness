---
owner: explorer
date: 2026-09-16
---
# Recon — diagnose-symlink-reads (2026-09-16)

## Question
Why does a symlink committed under `context/`, `rules/`, `commands/`, `workflows/` or `agents/` deploy as its target's bytes, while one under `skills/` or inside a directory `sh-agent adopt` copies is refused? Rank root causes, file:line, evidence, falsification; state per-directory whether affected and the code path.

## Summary
The symlink refusal (`d.Type()&fs.ModeSymlink != 0` → error) was added in commit `cceee54` ("skill walk and adopt refuse symlinks") in exactly two places: `plan.go`'s skill `WalkDir` and `adopt.go`'s `copyTree`. It was never extended to the three other readers that also use `os.ReadFile`/`filepath.Glob` on repo files: `LoadAgents`, `LoadDocs` (backs both `commands`/`workflows` and `rules`), and `RenderContext` (backs `context`). Those three call `os.ReadFile(p)` directly on a glob-matched or walked path with no `os.Lstat`/`fs.ModeSymlink` check first, so Go's `os.ReadFile` transparently dereferences the symlink and returns the target's bytes, which then get written out at fixed mode `0o644` (or, for context, concatenated into the deployed `AGENTS.md`). This is one repeated omission, not three independent bugs — the guard exists and works, it just was applied to 2 of 5 directories.

## Top candidate root causes (ranked)

1. **Missing symlink guard in `LoadDocs`** — `/home/kesokaj/workspace/styrman-harness/internal/harness/inventory.go:186-210` (loop `for _, p := range paths { raw, err := os.ReadFile(p) ... }`, no line before it checks `d.Type()`/`Lstat`). Backs both `commands`+`workflows` (aliased to the same loader, see `LoadInventory` at inventory.go:88 calling `LoadDocs(repo, "commands", "workflows")`) and `rules` (inventory.go:90, `LoadDocs(repo, "rules")`). Evidence: `filepath.Glob` returns the matched path unchanged; `os.ReadFile` follows a symlink per Go's `os` package docs; no `fs.ModeSymlink` check anywhere in this function (confirmed via `grep -n "symlink\|ModeSymlink\|Lstat" internal/harness/*.go` — no hits in inventory.go around this function). Falsify: add a symlink under `rules/x.md` pointing at `/etc/hostname`, run `LoadDocs`/`sh-agent install`; if the resulting `Doc.Raw` is the link's own path text (or install errors) rather than `/etc/hostname`'s content, this is not the cause.

2. **Missing symlink guard in `LoadAgents`** — `/home/kesokaj/workspace/styrman-harness/internal/harness/inventory.go:155-176` (`for _, p := range paths { raw, err := os.ReadFile(p) ... }`, `paths` from `filepath.Glob(agents/*.md)`). Same pattern and same evidence/falsification approach as (1), scoped to `agents/`. Note `Plan` then writes `a.Raw` verbatim to `~/.claude/agents/<name>.md` at `plan.go:150` with `Mode: 0o644` and, via `DeriveAgy`, into the Antigravity mirror — so the same bytes reach both vendors' agent definitions, not just Claude's.

3. **Missing symlink guard in `RenderContext`** — `/home/kesokaj/workspace/styrman-harness/internal/harness/render.go:20-40`. Two spots: `os.ReadFile(filepath.Join(dir, "AGENTS.md"))` (line ~23, the `context/AGENTS.md` case) and the loop over `extras` from `filepath.Glob(context/*.md)` (`os.ReadFile(p)` around line ~35). Both feed directly into the `[]byte` returned as `context`, which `Plan` (plan.go:88-97) writes to `~/.claude/AGENTS.md`, `~/.claude/CLAUDE.md`'s target, and `~/.gemini/config/rules/AGENTS.md` — this is the path that explains the symptom's extra detail that a `context/` link's bytes are "also concatenated into the AGENTS.md that both CLIs load and send to the model provider." Falsify: replace `context/x.md` with a symlink to `/etc/hostname`, run `sh-agent install --home <scratch>`, `grep hostname <scratch>/.claude/AGENTS.md` (or the file's contents) — if hostname's content is absent, look elsewhere.

These three are one root cause pattern (the guard added in `cceee54` was scoped to skills+adopt only), so treat them as a single fix surface rather than three unrelated defects.

## Per-directory affected/unaffected table

| Directory | Affected? | Code path that reads it |
|---|---|---|
| `context/` | **Yes** | `RenderContext` (render.go:20-40) — `os.ReadFile` on `context/AGENTS.md` and each `context/*.md` glob match, no symlink check. Feeds `Plan` (plan.go:88-97) → `~/.claude/AGENTS.md`, `~/.claude/CLAUDE.md` stub's target, `~/.gemini/config/rules/AGENTS.md`. |
| `rules/` | **Yes** | `LoadDocs(repo, "rules")` called from `LoadInventory` (inventory.go:90) → `LoadDocs` (inventory.go:186-210), `os.ReadFile` on each glob match, no symlink check. Feeds `Plan` (plan.go:167-174) → `~/.claude/rules/<name>.md`, `~/.gemini/config/rules/<name>.md`, both `Mode: 0o644`. |
| `commands/` | **Yes** | `LoadDocs(repo, "commands", "workflows")` (inventory.go:88) — same `LoadDocs` function as rules, same lack of guard. Feeds `Plan` (plan.go:157-164) → `~/.claude/commands/<name>.md`, `~/.gemini/config/workflows/<name>.md`. |
| `workflows/` | **Yes** | Same call and function as `commands/` (they are the two `dirs` args to one `LoadDocs`); identical code path, identical gap. |
| `agents/` | **Yes** | `LoadAgents` (inventory.go:155-176), `os.ReadFile` on each `filepath.Glob(agents/*.md)` match, no symlink check. Feeds `Plan` (plan.go:147-153) → `~/.claude/agents/<name>.md` directly, and via `DeriveAgy` (derive.go, not read in this pass) → `~/.gemini/config/agents/<name>/agent.md`. |
| `skills/` | **No — refused** | `Plan`'s skill walk (plan.go:99-135): `filepath.WalkDir` over each skill dir; line 108 `if d.Type()&fs.ModeSymlink != 0 { return fmt.Errorf(...) }` before any `os.ReadFile`. Note `LoadSkills` (inventory.go:113-146) does `filepath.EvalSymlinks` on the **skill directory itself** (to allow a symlinked skill dir at the top level, e.g. from `npx skills add -g`) — that is a deliberate, different mechanism from a symlinked *file inside* the directory, which the WalkDir check still catches. |
| (adopt target) | **No — refused** | `copyTree` (adopt.go:88-115), same guard at line 93-94, same error text pattern ("is a symlink — refused"), added in the same commit `cceee54` as the skills-walk guard. |

## Flows
- `Plan(repo, selfExe, vendor)` (plan.go:73) → `LoadInventory(repo)` (inventory.go:84) → `LoadSkills`, `LoadAgents`, `LoadDocs`×2 → `RenderContext(repo, inv)` (render.go:20) → assembles `[]File{...}` → `apply.go` (not read this pass) presumably writes them to `$HOME`.
- The guarded path: `Plan`'s own `filepath.WalkDir` closure for skills (plan.go:99) checks `d.Type()&fs.ModeSymlink` (line 108) before `os.ReadFile` (line 118) — this is the pattern that should exist, but does not, in `LoadAgents`/`LoadDocs`/`RenderContext`.
- `adopt.go`'s `copyTree` (adopt.go:88) has its own, separately-added identical guard at line 93 — confirms the guard is a per-call-site convention here, not centralized, which is itself part of why it did not propagate.

## Conventions
- The codebase's convention for "refuse rather than follow" is: check `d.Type()&fs.ModeSymlink != 0` (from a `DirEntry`, i.e. inside a `WalkDir` callback) and return a descriptive `fmt.Errorf` naming the path. This exists only where a `WalkDir`/`copyTree`-style walk was written (skills, adopt); the three flat-glob loaders (`LoadAgents`, `LoadDocs`, `RenderContext`) never call `WalkDir` and never call `os.Lstat` either, so they have no equivalent check.
- Comments in plan.go:104-107 and adopt.go:92-93 both cite "security review 2026-09-14" as the origin of the guard — worth checking that review's scope if further history is wanted (not read in this pass).

## Risks
- Because the guard is duplicated per call site rather than centralized (e.g. in a shared `readRegular`/`walkRegular` helper), any future new loader is likely to reintroduce the same gap — this is presumably how the gap arose in the first place (skills/adopt got walk-style code with the check; agents/docs/context got flat `os.ReadFile`/`Glob` code without it).
- `RenderContext`'s output also reaches the model directly (concatenated into `AGENTS.md`), so this gap is not just an arbitrary-file-read into the consumer's home but a prompt-injection / exfiltration vector if the symlink target is attacker-influenced content, which the security-review comments elsewhere in this codebase treat as the main threat (mirrors the "backtick or newline... inject prose into the global context" reasoning already in inventory.go:56-64 for `validName`).

## Not examined
- `derive.go` (Antigravity agent derivation) — not read; assumed to just transform `a.Raw` bytes already loaded by `LoadAgents`, so it inherits the same gap rather than being an independent cause, but not verified by reading the file.
- `apply.go` — not read in this pass; assumed to write `File.Data` verbatim to `$HOME` per `File.Mode`, consistent with the symptom description, but not confirmed.
- The actual "security review 2026-09-14" review document/commit body was not located or read; only inferred from comments citing it in plan.go and adopt.go.
- Did not build a live repro fixture or run `sh-agent install` against a scratch home to empirically confirm; conclusions are from static code reading of `os.ReadFile`/`filepath.Glob` semantics plus the absence of any guard, cross-checked against the one commit (`cceee54`) that is known to have added the guard elsewhere.
