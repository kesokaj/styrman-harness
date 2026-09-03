---
owner: main
date: 2026-09-03
---
# Verification — every skill and agent, both CLIs, mock scenarios   (2026-09-03)

Range reviewed: 1aba06a..working tree (skill matrix after the sandbox fixes)

## Status: PASS — 31 scenarios × 2 CLIs; no harness defect found in this pass

## Method

One mock scenario per skill (30 skills; `feature`/`bug`/`init-project`/`preflight` also cover the six agents they dispatch, plus a direct `explorer` dispatch and a context question about the deploy guides). Each scenario runs headless in a fresh clone of the demo service (`/home/kesokaj/sandbox/runs/<cli>-<id>`), Claude Code with `--permission-mode acceptEdits` and a command allowlist, Antigravity with `--mode accept-edits` and `permissions.allow` in its settings. A shell check decides PASS/FAIL from files produced and the reply; the CLI's own transcript is then read for which skills loaded and which agents were dispatched. Cloud and container commands were never allowed for Claude and never run by Antigravity (audited: the only executions were read-only `gcloud config list`/`which gcloud`, and `git push` inside `/preflight`, which fails without a remote as intended). Kit: `.agents/scratch/skilltest/` (gitignored): `matrix.txt`, `prompts/`, `checks/`, `run-one.sh`, `summarize.py`, `results/`.

Seeds: three identical tooling corrections in `~/.config/sh-agent/feedback.jsonl` for `/retro` and `/distill`; a committed `TrimLeft` regression in `greet.go` for `/bug`; the bare three-file service at `8cbb128` for `/init-project`.

## Results

| Scenario | What it exercises | Claude Code | Antigravity |
|---|---|---|---|
| `init-deploy` | `/init-project` on the bare service; deploy guides filled / marked not applicable | PASS (155s) | PASS (130s) |
| `bug` | `/bug` trailing-space defect → failing test, fix, bug record, commit | PASS (170s) | PASS (340s) |
| `audit` | `/audit internal` → report under `reviews/` | PASS (manual) | PASS (213s) |
| `preflight` | `/preflight` with no tag, no remote → verdict and report | PASS (436s) | PASS (285s) |
| `feature` | `/feature` GET /metrics → plan file with phases and boxes, no code | PASS (manual) | PASS (manual) |
| `research` | `/research` light mode → stops at the continue question | PASS (107s) | PASS (50s) |
| `retro` | `/retro` over seeded feedback entries → proposal file in `~/.config/sh-agent/proposals/` | PASS (147s; the first run without `--add-dir ~/.config/sh-agent` failed on that access alone) | PASS (119s) |
| `distill` | `/distill --from-log` over seeded repeats → project-local skill or a reasoned no | PASS (137s, with the config directory allowed) | PASS (83s) |
| `handoff` | `/handoff` → file under `handoffs/` | PASS (manual) | PASS (manual) |
| `cloud-run` | "Deploy this service to Cloud Run." → asks the five defaults, points at `deploy/gcp.md`, runs nothing | PASS (48s) | PASS (manual) |
| `explorer` | `explorer` agent maps the repo, changes nothing | PASS (47s) | PASS (70s) |
| `ctx-deploy` | context: quote DoD step 3, the layout line, where GCP defaults live (no tools) | PASS (7s) | PASS (8s) |
| `feedback` | a correction in the prompt → entry appended to `feedback.jsonl` | PASS (44s) | PASS (116s) |
| `grill-me` | `/grill-me` → questions, no file changes | PASS (50s) | PASS (67s) |
| `grilling` | "grill me on this idea" → auto-loads `grilling` | PASS (59s) | PASS (29s) |
| `find-skills` | "is there a skill for Helm?" → `npx skills` path, no install | PASS (45s) | PASS (36s) |
| `skill-repair` | repair steps for a failed install, no changes | PASS (34s) | PASS (16s) |
| `golang` | review greet.go against the Go standards, no changes | PASS (25s) | PASS (39s) |
| `bash` | write and `bash -n` a health-check script | PASS (41s) | PASS (57s) |
| `python` | write and `py_compile` a test-summary script | PASS (25s) | PASS (64s) |
| `typescript` | typed client function, chat only | PASS (18s) | PASS (20s) |
| `frontend` | TanStack Start page plan, chat only | PASS (29s) | PASS (34s) |
| `docker` | Dockerfile for Colima, no build | PASS (59s) | PASS (59s) |
| `gke` | GKE Autopilot day-0 checklist, no gcloud | PASS (59s) | PASS (45s) |
| `gcloud-auth` | "ADC not found" → auth skills, runs nothing | PASS (18s) | PASS (14s) |
| `terraform` | Cloud Run module skeleton, no terraform run | PASS (28s) | PASS (29s) |
| `json-api` | error body design, chat only | PASS (41s) | PASS (22s) |
| `mermaid` | sequence diagram file | PASS (19s) | PASS (62s) |
| `codebase-design` | seam placement in deep-module vocabulary | PASS (42s) | PASS (36s) |
| `review-standards` | apply the review criteria to main.go | PASS (56s) | PASS (60s) |
| `anti-ui-slop` | critique a three-card status page (vendored skill) | PASS (38s) | PASS (32s) |

Times are wall-clock per run; "manual" marks the runs whose result line was lost when the driver script was edited while they were in flight (their outputs were intact and their checks were run afterwards).

## Evidence from the transcripts

- Every knowledge skill auto-loaded from its description on both CLIs (`golang`, `bash`, `python`, `typescript`, `frontend`, `docker`, `gke-basics`, `gcloud-setup` + `gcloud-auth-verification`, `terraform`, `json-api-style`, `mermaid-syntax`, `codebase-design`, `review-standards`, `anti-ui-slop`, `find-skills`, `skill-repair`, `feedback`, `grilling`). `/grill-me` delegates to `grilling` as designed.
- Workflow skills dispatched the agents they name: `/audit` → `explorer`, `reviewer`, `security`, `verifier`, `changelog`; `/preflight` → `reviewer`, `security`, `verifier` (+ `changelog` on SHIP); `/bug` → `verifier`, `bugs`, `changelog` (goldfish `explorer` skipped, allowed for a one-line fix in an obvious place); `/feature` → `reviewer` as the cold reader, `json-api-style` for the response shape.
- Deploy guides: both CLIs quote the new Definition of Done step and layout line; both `/init-project` runs filled `deploy/local.md` with the real commands and left `remote.md`/`gcp.md` marked *Not applicable yet* with a `## Deploy` section in `AGENTS.md`; both answered "Deploy this service to Cloud Run" by asking for the five defaults and naming `deploy/gcp.md` as where they will be recorded, running nothing.
- `/preflight` reached opposite verdicts on the same kind of tree: Claude Code took the empty tree as base for a first release, ran the reviewer on Opus, and returned NO-SHIP on three findings (listener timeouts among them); Antigravity took the base from the newest review's `Range reviewed:` line, reused that review, and returned SHIP. Both readings fit the skill's base-ref rule; the difference is model strength and the ambiguity of a `Range reviewed: <sha>..worktree` line. Recorded, not changed.

## Vendor behaviour (not harness defects)

- Claude Code headless cannot read or write `~/.config/sh-agent/` (outside the working directory) without `--add-dir`; `/retro` and `/distill` said so explicitly and degraded rather than inventing data, then passed with the directory added. Interactively this is one approval prompt.
- Antigravity ignores `--add-dir` semantics for the project root and may read sibling checkouts inside its trusted workspace.

## Not verified here

Interactive use on the Mac; real deploys (no target exists, by design); `skill-repair` and `anti-ui-slop` beyond loading and answering (vendored, no install failure or UIZZE access to exercise).

## Environment

Linux (Debian 13, amd64), Go 1.24.4, Claude Code 2.1.259, Antigravity CLI 1.1.25 (Gemini 3.8 Flash); harness at 1aba06a plus this working tree.
