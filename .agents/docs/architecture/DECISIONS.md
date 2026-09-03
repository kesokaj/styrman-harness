# Architecture Decisions

Append-only. To reverse a decision, append a new ADR and mark the old one `Superseded by ADR-NNN` — never rewrite it. The ADRs of the previous design (a hook-driven Go harness) were removed with it in commits 09fc5cc and 20a5abd; they remain readable in git history.

## ADR-001: The harness is context, skills, and subagents — no hooks, no gates     (2026-09-03) — Accepted

**Context:** The previous harness enforced its Definition of Done with lifecycle hooks, a commit gate, scope digests, and a dispatch ledger. It was reverted the same day it was finished: the machinery cost more attention than the rules it enforced were worth for a solo operator.

**Decision:** Behaviour lives in prose: the global context (`context/AGENTS.md`), the skills, and the subagent definitions. Nothing runs on tool calls or session events. The only executable is the installer.

**Consequences:** Rules are followed because they are read, not because they are enforced. A rule that is not followed is fixed by rewriting the rule (`/retro`), not by adding a gate. Both vendors are served by identical mechanisms, since neither needs hook support.

## ADR-002: One Go binary, `styrman`, with a git resync loop     (2026-09-03) — Accepted

**Context:** The operator works on several machines and wanted to push a harness change from one and pick it up on another without thinking. A bash script was considered first; it was replaced by Go so that `styrman update` (fetch, pull, rebuild, reinstall) and a manifest could be tested against a scratch home rather than by hand.

**Decision:** `cmd/styrman`, stdlib only. `install` writes an atomic copy of every managed file and records the manifest in `~/.config/styrman/state.json`; `update` fast-forwards the repo and re-runs `go run ./cmd/styrman install`, so a code change in the repo takes effect in the same update; the binary copies itself to `~/.local/bin/styrman`.

**Consequences:** Go is required on every machine (already true for the operator). Two copies of each skill exist on disk, one per vendor — a deliberate trade against symlinks, which neither CLI documents following across directories.

## ADR-003: Single source; Antigravity agents are derived at install     (2026-09-03) — Accepted

**Context:** The previous layout committed generated copies (`claude/CLAUDE.md`, `antigravity/GEMINI.md`, derived agents) that needed a drift check to stay honest.

**Decision:** `context/`, `agents/` (Claude Code format), and `skills/` are the only sources. `styrman install` assembles the context (AGENTS.md + dev-setup.md + generated inventory) and translates agent frontmatter for Antigravity in memory. Nothing generated is committed.

**Consequences:** No drift is possible between committed copies. The translation table (`toolMap`, `modelMap` in `derive.go`) is the one place vendor differences live; an unmapped tool fails the install instead of shipping a name Antigravity would hang on.

## ADR-004: Uninstall is manifest-only; hand-installed skills survive     (2026-09-03) — Accepted

**Context:** The previous installer mirrored whole directories and pruned anything the repo did not have, which wiped skills installed by hand with `npx skills add -g`. The operator wants adding skills by hand to stay easy.

**Decision:** `install` records every file it writes; `uninstall` and the resync step remove only paths from that manifest. A `~/.claude/CLAUDE.md` without the styrman marker is never overwritten. `styrman adopt` copies a hand-installed skill into the repo when it should become permanent.

**Consequences:** Directories under `~/.claude/skills` may hold both styrman's and the operator's skills; `status` reports only styrman's. A skill installed by hand into a vendor directory is not versioned and is not deployed to the other vendor until adopted.

## ADR-005: The context is named AGENTS.md everywhere     (2026-09-03) — Accepted

**Context:** The operator asked for the agents.md standard. Claude Code reads `CLAUDE.md`; Antigravity reads `AGENTS.md` or `GEMINI.md`.

**Decision:** The deployed file is `~/.claude/AGENTS.md` and `~/.gemini/config/rules/AGENTS.md`. `~/.claude/CLAUDE.md` and a project's `CLAUDE.md` are the single line `@AGENTS.md`. `/init-project` creates `.claude/{skills,agents,commands}` as symlinks into `.agents/`, matching what the skills.sh CLI does on its own.

**Consequences:** One canonical name at every level. Claude Code's import adds one hop; it is documented to load user-level imports without prompting.

## ADR-006: The binary is `sh-agent`; GCP defaults are asked for, never stored     (2026-09-03) — Accepted

**Context:** The operator asked for a less repo-specific binary name and, on reviewing `context/dev-setup.md`, rejected stored GCP defaults: a placeholder region or registry name is exactly the kind of value an agent will silently assume. They also asked for zero scripts in the repo.

**Decision:** The binary, its state directory (`~/.config/sh-agent`), its env var (`SH_AGENT_REPO`), and its file marker are named `sh-agent`; the repo keeps its name. The context says work is local by default and, when the user says GCP, the agent asks for project, region, Artifact Registry repository, service name, and visibility before doing anything, then records them in the project's `AGENTS.md`. The one shell script in the repo (`mermaid-syntax/scripts/validate-mermaid.sh`) was deleted along with the old hook scripts under `~/.agents/bin`.

**Consequences:** ADR-002's name is superseded by this one; its content stands. The first GCP action in a new project costs one question instead of one wrong assumption. Skills that need helper code ship it as documentation, not as scripts.
