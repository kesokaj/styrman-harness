---
owner: reviewer
date: 2026-09-15
---
# Review — sweep-skills-3   (2026-09-15)

Range reviewed: Scope: the named files as they are in the working tree, read whole — not a diff: skills/handoff/, skills/init-project/, skills/json-api-style/, skills/mermaid-syntax/, skills/preflight/, skills/python/, skills/research/, skills/retro/, skills/review/, skills/review-standards/ (HEAD d1ff612 plus the uncommitted edits at skills/retro/SKILL.md:75 and skills/review/SKILL.md:48)

Lens: `--docs`, with the `--project` classification clause applied to the named files. Report only.

## Verdict: REQUEST CHANGES

The worst problem: `skills/retro/SKILL.md` says a retro is never started uninvited and is only *suggested*, and its trigger list leaves out the 14-day cadence. The global context (`context/AGENTS.md:41`) and the binary's own `sh-agent status` hint (`cmd/sh-agent/main.go:271`) both tell the agent to *run* `/retro`, including on that 14-day cadence. An agent following the skill never acts on the hint. Two rule restatements in `/preflight` are also narrower than their sources: the security trigger drops "crypto", and the docs-freshness step re-lists a shorter target set than the `--docs` brief now carries. Both init-project templates for ROADMAP and TODO still include "Ideas" sections, which the Definition of Done forbids. The rest are low-severity stale names and internal contradictions between skills and their reference files. Nothing is broken in the installer code this scope describes.

## Findings

- **[medium / confidence: high] divergent** `skills/retro/SKILL.md:3` ("Never start one uninvited") and `skills/retro/SKILL.md:25-32` ("When to Suggest … two or more … Suggest it naturally") versus `context/AGENTS.md:41` ("Run `/retro` at the end of a session with two or more corrections, and whenever `sh-agent status` shows more than 14 days since the last one") and `cmd/sh-agent/main.go:271` ("past 14 days, run /retro at the end of the next session"). The skill makes it a suggestion and leaves out the 14-day trigger (narrowed). The context and the binary tell the agent to run it. Should say: the skill's trigger section names both conditions from context:41 and uses the same verb as the context ("run", or change both to "suggest").
- **[medium / confidence: high] divergent** `skills/init-project/assets/ROADMAP.md:15-17` (`## Later / Ideas` — "Unscheduled directions worth keeping visible") and `skills/init-project/assets/TODO.md:13` (`## Ideas / Nice-to-Have`) versus `context/AGENTS.md:18` ("at Done both are EMPTY … neither file is a parking lot for ideas"), `skills/handoff/SKILL.md:64`, and each template's own line 3. Every new project is scaffolded with sections that invite the parking the Definition of Done forbids. Related: `agents/documentation.md:30,43` still assigns Next/Later to documentation ("move items between Next and Later"). Should say: no ideas sections, or a comment that they stay empty at Done and that ideas go to the handoff's Next Steps.
- **[medium / confidence: high] narrowed: `context/AGENTS.md:16` → `skills/preflight/SKILL.md:51`** — the security row triggers on "auth, IAM, secrets, input parsing, or dependencies" and leaves out **crypto**. Both the context and `agents/security.md:3` list "auth, crypto, secrets, IAM, input parsing, or dependencies". A release that touches only crypto code gets no security row. Should say: "auth, crypto, secrets, IAM, input parsing, or dependencies".
- **[low / confidence: high] narrowed: `skills/review/SKILL.md:48` → `skills/preflight/SKILL.md:65`** — preflight Step 4 dispatches "the docs brief" but re-lists its targets as README, specs, the testing guide and deploy/ only. It omits the newest `changes/` entry and the old-wording `git grep` over `AGENTS.md`, `.agents/{rules,skills,agents,commands}/` and the paths the project `AGENTS.md` adds (here `context/`, `agents/`, `skills/`), all of which the brief (now uncommitted at review:48) carries. Should say: "the `--docs` brief from `/review` Step 3, verbatim", with no second list.
- **[low / confidence: medium] narrowed: `context/AGENTS.md:16` → `skills/review/SKILL.md:60`** — the `--fix` stop rule keeps "a round's findings land on the mechanism the previous round fixed" but drops the context's second condition, "or a fix introduces the next round's defect". Should say: add that clause.
- **[low / confidence: medium] divergent** `skills/preflight/SKILL.md:87` — "`local.md` when the project is local-only: build for `linux/amd64`, push to Artifact Registry, deploy". The amd64/Artifact Registry steps belong only to `gcp.md` (`context/dev-setup.md:10-11`). The template `skills/init-project/assets/agents-tree/docs/deploy/remote.md:29-33` copies with `scp`, and this repo's `.agents/docs/deploy/local.md:14-23` installs into a home. Should say: follow the chosen guide's own procedure and verify step; amd64 build plus Artifact Registry push only for `gcp.md`.
- **[low / confidence: high] divergent** `skills/init-project/SKILL.md:45` — "note that agents never edit it" versus `context/AGENTS.md:52` (after approval and an ADR, the agent updates `REQUIREMENTS.md`) and the template it copies, `skills/init-project/assets/REQUIREMENTS.md:4` ("Never modify without explicit user approval"). Should say: "never edited without the user's explicit approval".
- **[low / confidence: high] stale** "Resume Protocol" at `skills/handoff/SKILL.md:34`, `:150`, `:154` and `skills/init-project/SKILL.md:18` (also `agents/changelog.md:84`). No living doc defines that name: the global context's section is `## Resume, then checkpoint` (`context/AGENTS.md:7`), and `git log -S 'Resume Protocol' -- context/` finds no commit. Should say: "the resume order in the global context (`Resume, then checkpoint`)".
- **[low / confidence: high] stale** `skills/init-project/assets/AGENTS.md:40` — "supplement the global testing tiers". Tiers left the global context in 87a862d ("move procedures into feature and bug skills") and are now named in `skills/feature/SKILL.md:48`. Should say: "supplement the tiers each plan names (feature skill)".
- **[low / confidence: medium] divergent** `skills/init-project/assets/agents-tree/docs/deploy/local.md:8` — "The `verifier` follows it when a plan names a local smoke test". `agents/verifier.md:37,42-49` reads only `AGENTS.md` and `testing/<app>/guide.md` and never names `deploy/local.md`. Should say: either verifier.md names the deploy guide as a source, or the template drops the claim.
- **[low / confidence: high] divergent** `skills/review-standards/SKILL.md:3` ("loaded by the reviewer agent, /review, and /preflight") and `:8` ("Whoever is reviewing (the `reviewer` agent, `/review`, `/preflight`) loads this file and judges the diff") versus `skills/review/SKILL.md:9` ("It never reviews by itself") and `skills/preflight/SKILL.md:48` ("Never self-review"). Only `agents/reviewer.md:6-7` loads the skill. Should say: "loaded by the `reviewer` agent, which `/review` and `/preflight` dispatch, and by the main agent when it reads a below-threshold diff itself".
- **[low / confidence: medium] divergent** `skills/handoff/SKILL.md:78` — "Session ID: <conversation-id> … matches the feedback log's session field" versus `skills/feedback/SKILL.md:44` (`"session": "first-8-chars-of-conversation-id"`). A full id never string-matches the log. Should say: "its first 8 characters are the feedback log's session field".
- **[low / confidence: medium] divergent** `skills/review-standards/references/typescript-review.md:29` — "flag `export default` for anything but framework-mandated files" versus `skills/typescript/SKILL.md:40` ("no default exports for shared modules"). The checklist is broader than the idiom authority, and `skills/review-standards/SKILL.md:36` says the language skill wins. Should say: scope the check to shared modules, as the typescript skill does.
- **[low / confidence: high] divergent** `skills/python/references/python-fundamentals.md:12,39` (marked ✅: `from typing import Optional`; `from collections import defaultdict` at :11) contradicts the same file's `:27` ("❌ Do not import individual classes/functions from modules") and `:158` ("Use X | None instead of Optional[X]"), and `skills/python/SKILL.md:35` rule 4. Should say: `import collections` / `import typing`, with no `Optional`.
- **[low / confidence: high] divergent** `skills/python/references/python-fundamentals.md:73` — `logger.error(f"Invalid integer input: {e}")` versus `skills/python/SKILL.md:48` rule 13 and the same file's `:255-259` (lazy `%` formatting in logging calls). Should say: `logger.error("Invalid integer input: %s", e)`.
- **[low / confidence: high] incorrect example** `skills/python/references/python-fundamentals.md:278` — `with threading.Lock():` creates a new lock on every entry, so nothing is ever excluded. Code copied from the "✅ Locks" example has a data race. Should say: `lock = threading.Lock()` at module or instance scope, then `with lock:`.
- **[low / confidence: medium] divergent** `skills/mermaid-syntax/SKILL.md:29,38,249` ("Always wrap text containing special characters in double quotes", with `? < >` in the list) versus the skill's own unquoted examples: `examples/flowchart-examples.md:46` (`Attempt{Attempts < 3?}`), `:93` (`Q1{Budget > 1000?}`), `:94`, and `references/flowchart-complete.md:324` (`Check{Condition?}`). Should say: quote those labels, or narrow the rule to the characters that actually break parsing.
- **[low / confidence: high] divergent** `skills/mermaid-syntax/references/other-diagrams.md:169` — heading "Bottom-to-Top (v11)", but the code at `:173` is `gitGraph TB:` (top to bottom). Should say: `gitGraph BT:`, or retitle the heading "Top-to-Bottom".
- **[low / confidence: medium] stale** `skills/review/SKILL.md:4` — `argument-hint` lists the lens flags plus `--fix` and `--high`, but not `--project`, which is parsed at `:25` and shown in usage at `:17`. Should say: add `[--project]`.

### Old-wording grep for what the scope changes or removes

The uncommitted edits in scope only add text: `skills/retro/SKILL.md:75` adds the removed-rule home/`dropped:` sentence, and `skills/review/SKILL.md:48` appends the grep and `--project` clauses to the `--docs` brief. Neither removes wording, so there is no old wording to grep for and no `dropped:` finding. The one restatement of the extended brief that did not follow the addition is the preflight:65 finding above.

### Unverified

- `skills/review/SKILL.md:21` — Antigravity dispatches through `invoke_subagent`. `internal/harness/derive.go:15-22` maps only agent `tools:` names, and no log under `.agents/docs/testing/` records the name.
- `skills/preflight/SKILL.md:44` — "a line ending in `..worktree` names its base". Nothing current writes that form: `skills/review/SKILL.md:37` writes the scope text, and `agents/reviewer.md:71` and `skills/preflight/SKILL.md:71` write `<base>..<head>`. It was added in cebd2c9, and I could not confirm whether a writer still exists.
- Mermaid inline `%%` comments after a statement (`skills/mermaid-syntax/SKILL.md:46-48,64-65,132-133`; `references/flowchart-complete.md:127-130,284`; `references/sequence-complete.md:300`). Mermaid documents comments on their own line, and in a sequence message the trailing text may render as message text. No diagram was rendered.
- `skills/mermaid-syntax/references/sequence-complete.md:14-24` — the heading says "JSON syntax" but the code shows `<<boundary>>` stereotypes, not `participant X@{ "type": ... }`. Not rendered.
- `skills/python/references/python-testing.md:190` — `AsyncClient(app=app)`. Recent httpx releases need `transport=ASGITransport(app=app)`; version not checked.
- External claims in `skills/mermaid-syntax/SKILL.md` (v11.0-11.4 arrowless bug, secure-setting limits, architecture label charset) and the Google JSON Style Guide URL at `skills/json-api-style/SKILL.md:55`.

## Classification of the named living docs

| File | Class | Proof |
|---|---|---|
| skills/handoff/SKILL.md | outdated | :34/:150/:154 "Resume Protocol" vs context/AGENTS.md:7; :78 vs feedback/SKILL.md:44 |
| skills/init-project/SKILL.md | outdated | :45 vs context/AGENTS.md:52; :18 "Resume Protocol" |
| skills/init-project/assets/AGENTS.md | outdated | :40 "global testing tiers" (removed from context in 87a862d) |
| skills/init-project/assets/CLAUDE.md | accurate | :1 `@AGENTS.md` = context/AGENTS.md:69 |
| skills/init-project/assets/REQUIREMENTS.md | accurate | :4 = context/AGENTS.md:52 |
| skills/init-project/assets/ROADMAP.md | outdated | :15-17 vs context/AGENTS.md:18 |
| skills/init-project/assets/TODO.md | outdated | :13 vs its own :3 and context/AGENTS.md:18 |
| skills/init-project/assets/agents-tree/docs/deploy/local.md | outdated | :8 vs agents/verifier.md:37 |
| skills/init-project/assets/agents-tree/docs/deploy/remote.md | accurate | :8 = context/dev-setup.md:7 |
| skills/init-project/assets/agents-tree/docs/deploy/gcp.md | accurate | :16-22 hold the five defaults of context/dev-setup.md:4; :27 = dev-setup.md:10-11 |
| skills/init-project/assets/agents-tree/docs/architecture/DECISIONS.md | accurate | :7 shape = agents/changelog.md:73 |
| skills/init-project/assets/agents-tree/docs/bugs/INDEX.md | accurate | :5 columns = agents/bugs.md:64 |
| skills/init-project/assets/agents-tree/**/.gitkeep | accurate | deployed: internal/harness/plan.go:111 skips only dirs and .DS_Store |
| skills/json-api-style/SKILL.md | accurate | :37 = references/json-api-conventions.md:107 |
| skills/json-api-style/references/json-api-conventions.md | accurate | :130 = SKILL.md:38 |
| skills/mermaid-syntax/SKILL.md | outdated | :38/:249 contradicted by examples/flowchart-examples.md:46 |
| skills/mermaid-syntax/references/flowchart-complete.md | outdated | :324 unquoted `?` vs SKILL.md:38 |
| skills/mermaid-syntax/references/sequence-complete.md | accurate (unverified :14-24) | :316 `#59;` = SKILL.md:72 |
| skills/mermaid-syntax/references/class-state-complete.md | accurate | :174-181 namespaces without relationships = SKILL.md:175 |
| skills/mermaid-syntax/references/other-diagrams.md | outdated | :169 vs :173 |
| skills/mermaid-syntax/examples/flowchart-examples.md | outdated | :46/:93 vs SKILL.md:29 |
| skills/mermaid-syntax/examples/sequence-examples.md | accurate | :236 `#59;` = SKILL.md:72 |
| skills/preflight/SKILL.md | outdated | :51 vs context/AGENTS.md:16; :65 vs review/SKILL.md:48; :87 vs dev-setup.md:10 |
| skills/python/SKILL.md | accurate | :48 = references/python-fundamentals.md:255-259 |
| skills/python/references/python-fundamentals.md | outdated | :12 vs :27/:158; :73 vs SKILL.md:48; :278 |
| skills/python/references/python-testing.md | accurate (unverified :190) | :207 `asyncio_mode = "auto"` consistent with :183-197 |
| skills/research/SKILL.md | accurate | :45 spec path = agents/documentation.md:26; "Step 1" = feature/SKILL.md:17 |
| skills/retro/SKILL.md | outdated | :3/:25-32 vs context/AGENTS.md:41 and cmd/sh-agent/main.go:271 |
| skills/review/SKILL.md | outdated | :60 vs context/AGENTS.md:16; :4 vs :25 |
| skills/review-standards/SKILL.md | outdated | :3/:8 vs review/SKILL.md:9 |
| skills/review-standards/references/go-review.md | accurate | :37 slog JSON = context/dev-setup.md:5 |
| skills/review-standards/references/typescript-review.md | outdated | :29 vs skills/typescript/SKILL.md:40 |

No named file is redundant. The overlaps (retro Step 2 and feedback's rule_ref vocabulary; init-project's tree and context/AGENTS.md:62-67) restate each other with the same meaning and point at each other.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Re-check commands

```bash
git grep -n -e 'Never start one uninvited' -e 'Suggest it naturally' -- skills/retro/SKILL.md
git grep -n -e 'Later / Ideas' -e 'Unscheduled directions' -e 'Ideas / Nice-to-Have' -- skills/init-project/assets
git grep -n 'touches auth, IAM, secrets' -- skills context agents
git grep -n "against \`README.md\`, \`<docs-root>/specs/\`" -- skills/preflight/SKILL.md
git grep -n 'land on the mechanism the previous round fixed' -- skills/review/SKILL.md
git grep -n 'local-only: build for `linux/amd64`' -- skills/preflight/SKILL.md
git grep -n 'agents never edit it' -- skills/init-project/SKILL.md
git grep -n 'Resume Protocol' -- skills agents context README.md AGENTS.md
git grep -n 'global testing tiers' -- skills
git grep -n 'verifier` follows it' -- skills/init-project/assets
git grep -n -e 'loaded by the reviewer agent, /review, and /preflight' -e 'Whoever is reviewing' -- skills/review-standards/SKILL.md
git grep -n "matches the feedback log's session field" -- skills/handoff/SKILL.md
git grep -n 'flag `export default`' -- skills/review-standards/references/typescript-review.md
git grep -n -e 'from typing import Optional' -e 'from collections import defaultdict' -e 'logger.error(f' -e 'with threading.Lock():' -- skills/python/references/python-fundamentals.md
git grep -n -e 'Attempts < 3?' -e 'Budget > 1000?' -e 'Need Speed?' -e 'Check{Condition?}' -- skills/mermaid-syntax
git grep -n -e 'Bottom-to-Top' -e 'gitGraph TB:' -- skills/mermaid-syntax/references/other-diagrams.md
git grep -n 'argument-hint' -- skills/review/SKILL.md
git grep -n -e 'invoke_subagent' -e '\.\.worktree' -- skills agents .agents/docs/testing
```

## Not examined

- No Mermaid diagram was rendered and no Python example was executed; the external-behaviour claims above are listed as unverified.
- `.agents/docs/specs/` is empty, so there was nothing to compare. The newest `changes/` entry (2026-09-15-research-beads-elephant-goldfish.md) does not touch the scope.
- `skills/golang`, `skills/frontend` and `skills/cloud-run-deploy` were grepped only for the points the scope cites (TS default exports, the rollback section at cloud-run-deploy:97), not read whole.
- The active plan `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` was not read, so some findings may already be queued there.
