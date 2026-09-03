---
owner: main
date: 2026-09-03
---
# Research — harness engineering (Böckeler) and elephant/goldfish (Rensin, Shvedov)   (2026-09-03)

## Question

The operator asked for two sources to be read in full and weighed against this harness: Birgitta Böckeler's "Harness engineering for coding agent users" (martinfowler.com, 2 April 2026) and the `vshvedov/elephant-goldfish` repository (Claude Code, Codex, Gemini CLI adapters of Dave Rensin's elephant/goldfish pattern). Which of their strategies should this harness adopt, given its rules: no hooks, no scripts, no CI, one Go installer, everything else prose?

## Summary

Both sources fit the harness better than the previous hook-driven design did, because both are mostly about *what the agent is told and what it is made to check*, not about lifecycle machinery. Böckeler gives a vocabulary and a priority order: **guides** (feedforward: context, skills, rules) and **sensors** (feedback: tests, linters, type checks, reviewers), each **computational** (cheap, deterministic, run on every change) or **inferential** (a model judging, expensive, run proportionally), improved through a **steering loop** whenever an issue repeats. The elephant/goldfish repo gives concrete workflows: design before code, a fresh subagent that sees only the artifact (never the intent) to stress-test a design, a diagnosis, or a diff; a failing test before any bug fix; findings fixed or rebutted verbatim with a hard cap; mandatory user gates at scope confirmation and at the cap.

Adopted into `context/AGENTS.md` today: the work protocol (scope → plan file with phases and task boxes → cold-reader check of the plan → phases as checkpoints → e2e before done), the bug path (problem note → goldfish diagnosis on symptom and repro only → failing test → smallest fix → bug record), the reviewer-sees-only-the-diff rule, fix-or-rebut with a three-round cap, and the steering loop as the learning-loop cadence plus the bug index. Not adopted: their "docs live in chat" default (this harness deliberately keeps the trail on disk so a zero-context agent can resume), their per-project bootstrap that rewrites command templates (the harness deploys one global set and lets projects add their own), and anything that needs a hook.

## Böckeler — what the article says

- **Harness = everything but the model.** For a coding-agent *user*, the outer harness is the context, skills, rules, tools, and checks wrapped around the vendor's built-in harness. Its two goals: raise the odds of a right first attempt, and self-correct before a human has to look.
- **Guides and sensors.** Guides anticipate (AGENTS.md, skills, bootstrap instructions, code mods). Sensors observe after the fact (tests, linters, structural tests, review instructions). Feedback-only harnesses repeat mistakes; feedforward-only ones never learn whether the rules worked. Sensor output should be written for the model — a linter message that says how to fix it is "a positive kind of prompt injection".
- **Computational vs inferential.** Deterministic checks are cheap enough to run on every change; model judgements are expensive and noisy and belong where semantic judgement is needed, with a model strong enough for the task.
- **Steering loop.** The human's job is to iterate on the harness: whenever an issue happens more than once, improve the guide or the sensor. Agents can draft the improvement.
- **Keep quality left.** Fast checks before a commit exists; expensive ones after integration; drift sensors (dead code, coverage quality, dependency scanners) continuously, outside the change lifecycle.
- **Three regulation categories.** Maintainability (tooling exists, easiest), architecture fitness (fitness functions: performance tests, logging conventions), behaviour (the hard one; AI-written tests are not enough, approved fixtures help selectively).
- **Harnessability and Ashby's law.** Typed languages, clear module boundaries, and opinionated frameworks are "ambient affordances". Committing to a small number of service topologies narrows the space the harness must govern; harness templates per topology are the likely next step, with the same drift problems service templates have.
- **The human.** A harness externalises what an experienced developer brings implicitly; it should direct human attention to where it matters, not remove it.

## Elephant/goldfish — what the repository does

- **Elephant** = the working session with full context; **goldfish** = a fresh subagent that receives only an artifact. The asymmetry is the test: if a goldfish given only the design doc cannot implement it, the doc is wrong. "Design is the new code."
- Five commands, adapted per stack at bootstrap: `eg-brainstorm` (parallel goldfish with different lenses), `eg-prd`, `eg-new-feature`, `eg-fix-bug`, `eg-precommit-review`.
- **New feature:** confirm scope in one or two sentences → design doc with required sections (why, scope, surfaces, interfaces, UX flow, failure modes, verification criteria, follow-ups) → **no-code gate** until two goldfish sign off: a critic (gaps) and a readiness check ("could a first-pass implementer ship this without a question?"), plus an informational comprehension pass → three-revision cap, then ask the user → layer-ordered implementation with per-layer verification → independent review → test gate incl. a real browser or device walkthrough → report, no commit.
- **Bug fix:** triviality gate → problem doc (symptom, repro, suspected area, hidden hypothesis, blast radius, "fixed means") → goldfish diagnosis on symptom and repro only, hypothesis withheld → converge or re-investigate → failing test that fails for the right reason → smallest fix → review → test gate → report.
- **Pre-commit review:** pre-flight checks run separately, new errors only are blockers → a goldfish gets the exact review template and the diff, never the intent → each finding fixed or rebutted verbatim, a visible ledger between rounds → hard cap of five rounds or a repeat finding at the same place, then a structured question to the user.
- **Interactivity is mandatory at decision points** and explicitly overrides "work autonomously" directives; only the user in the same turn can waive a gate.
- Bootstrap inspects the stack and rewrites `[BOOTSTRAP: …]` markers; updates are merged by prompt. Codex skills stay generic and inspect the repo at runtime instead.

## Mapping onto this harness

| Their idea | Here | Decision |
|---|---|---|
| Guides = context, skills, rules | `context/`, `skills/`, `rules/` | already the model; vocabulary adopted in the research, not the context |
| Computational sensors on every change | `gofmt`, `go vet`, `go test -race`, frontend runner, `docker build` | Definition of Done step 1; `verifier` runs them |
| Inferential sensors proportional to change | `reviewer` at 50+ lines / 5+ files, `security` on sensitive paths, `documentation` on behavior change | kept; made an ordered sequence |
| Steering loop | `/retro`, `/distill`, bug prevention rules that graduate into rules or tests | adopted: learning-loop cadence, bug index |
| Sensor messages for the model | verifier and reviewer logs written for a zero-context reader | already the rule |
| Keep quality left; drift sensors | fast checks before commit; `/audit` for drift | adopted in prose; no continuous job (no CI by rule) |
| Harness template per topology | this harness *is* the template for Go + TanStack Start + Cloud Run; `/init-project` instantiates it | adopted as framing; ADR-008 |
| Design before code, goldfish readiness check | plan file + `reviewer` with only the plan, three-round cap | adopted (Work protocol) |
| Goldfish bug diagnosis, failing test first | `explorer` with symptom and repro only; failing test before the fix | adopted (Work protocol, Bugs) |
| Reviewer sees only the diff; fix or rebut verbatim; hard cap | `reviewer` gets no intent; three rounds then ask | adopted (Definition of Done step 2) |
| Mandatory user gates | scope confirmation, no-code gate, review cap, missing repro | adopted as the four gates named in the Work protocol |
| Browser/device walkthrough before done | e2e tier: container run + curl, browser golden path + edge case | adopted (Work protocol step 6, tiers) |
| Docs live in chat unless there is a future reason | trail on disk always | not adopted — resume needs the disk |
| Per-project rewritten command templates | global skills + project `.agents/` layer | not adopted — one global set, projects add their own |
| Hooks (Stripe pre-push, OpenAI linters) | none | not adopted — ADR-007 |

## Open questions

- Behaviour harness: neither source has a good answer beyond tests plus manual walkthroughs. The e2e tier is the harness's answer for now; approved fixtures may be worth a skill later.
- Whether a `brainstorm`/`prd` pair (parallel goldfish with different lenses) earns its place: `/research` covers comparison; brainstorming stays a `grilling` conversation until a real need shows up.
