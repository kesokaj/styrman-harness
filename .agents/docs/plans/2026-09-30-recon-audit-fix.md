---
owner: main
date: 2026-09-30
---
# Recon — audit findings, consolidated (2026-09-30)

Source: two full audits of this repository on 2026-09-30 (first at f1ca525, second at 58bf00d), six read-only passes each (context and requirements; workflow skills; technical skills; agents; installer code; search, memory, dream and the trail), plus spot checks from the main session. Published as two private pages (`Styrman Harness Audit`, `Styrman Harness Audit II`). This file is the evidence base for plan `2026-09-30-audit-fix.md`; every task there cites an ID below. Confidence is high where a finding was reproduced (marked *repro*) or spot-checked (*checked*); otherwise it was read at the location named.

Two audit claims were refuted and are not listed: "feedback entries lack a project field" (all 15 carry one) and "the installed real-home context is behind the repo" (it carries the current memory line).

## C — global context and governing documents

| ID | Sev | Where | Finding |
|---|---|---|---|
| C1 | high | `context/AGENTS.md` | Rendered context 21,906 of 22,000 bytes (*checked*). No rule can be added; `rules/` cannot deploy. Hard rules ≈1,000 words. |
| C2 | high | `context/AGENTS.md:18` vs `:66` | "ROADMAP and TODO EMPTY at every checkpoint" contradicts "Now names the active plan while one is in progress". |
| C3 | high | `skills/init-project/assets/{TODO,ROADMAP}.md` | Templates carry High/Medium/Low, "Ideas / Nice-to-Have", "Later / Ideas" sections the context forbids as parking lots. This repo's own ROADMAP.md has "Later / Ideas" too. |
| C4 | med | `context/AGENTS.md:66` vs `REQUIREMENTS.md:11` | R2's CLAUDE.md stub list restated differently (Vertex, Foundry missing). |
| C5 | med | `context/AGENTS.md` | Docs-root shorthand `docs/` defined only on the last line; mixed notations `<docs-root>/testing/guide.md`, `docs/testing/guide.md`, bare `testing/`. |
| C6 | med | `context/AGENTS.md:9` | Resume order reads TODO.md second though line 18 keeps it empty. |
| C7 | med | `context/dev-setup.md:8` | Root steps gated on `.claude/settings.json`, a Claude-only file; on Antigravity the rule can never be satisfied (R10). |
| C8 | med | `assets/agents-tree/docs/bugs/INDEX.md` | No `## Graduated` table though the context tells `bugs` to move rows into one; this repo's INDEX.md:87 has the shape. |
| C9 | med | `context/AGENTS.md:52` | 169-word bullet chaining five unrelated checks (binary heredoc, escape, `rm -r`, ` D`, post-dispatch status/log/tag). |
| C10 | med | `context/AGENTS.md:58` | 120-word sentence restating four rules owned by review-standards § 2. |
| C11 | med | `context/AGENTS.md:16` | DoD step 2 is 246 words folding thresholds, rounds, "shorter beats defended", operator-choice and `/review --docs`. |
| C12 | med | `REQUIREMENTS.md` R3, R4, R8, R9, R10 | ≈400 words of approval provenance inside the immutable file; R9's parenthetical 188 words. Needs operator approval to change. |
| C13 | low | `context/AGENTS.md:46` | "bypass hooks" beside "no hooks, ever". |
| C14 | low | `context/AGENTS.md:21,25` | Docs-only ≤10 lines rule stated twice. |
| C15 | low | `context/AGENTS.md:30`; `REQUIREMENTS.md` R10 | `/goal` named but not a harness command; R10 promises a "known differences" list that no spec has. Candidates: `/goal`, the settings.json gate (C7), `.claude/*` links. |
| C16 | low | `context/AGENTS.md:9` | `sh-agent search --open` and `memory show` exit 2 in an unscaffolded checkout before `/init-project` is offered. |
| C17 | low | `AGENTS.md` (project) § Testing | One 397-word paragraph plus six BUG bullets; silent on `internal/memory`, `internal/dream`, guide § 3e and the D30 reader measurement; "protobuf" appears nowhere. |
| C18 | low | `assets/REQUIREMENTS.md:7-19` | Leads with Scale Targets and Latency SLA tables, a burden on a CLI or library. |
| C19 | low | `assets/AGENTS.md:42` | "supplement the global testing tiers"; the context defines no tiers. |
| C20 | low | `ROADMAP.md:3`, `assets/ROADMAP.md:3` | Header comment says "both sections" for a three-section file. |
| C21 | low | `README.md` map legend | Dreamer lines coloured amber (report back) though rule 1 says the timer is neither a way in nor a way back; what it sends to the model provider is not drawn. |
| C22 | low | `REQUIREMENTS.md:31` (R7) | Still "auto memory store"; replacement wording awaits operator approval (handoff 2206 Next Step 3). |
| C-scar | low | `context/AGENTS.md:15,34,46,52,57,58,60,61,62,66`; `dev-setup.md:8,10`; `REQUIREMENTS.md:55-56` | One-off incidents written as global law (tree-freeze, dispatch boilerplate, token file name, escape lesson, `hash -r`, mutation protocol, dry-run child cwd, embed pin, `&&` chain, date-only names, `du` before `rm -r`, `exec format error`, BUG-085 measurements). |

## WS — workflow skills

| ID | Sev | Where | Finding |
|---|---|---|---|
| WS1 | med | `skills/review/SKILL.md:51,52,64` | Single paragraphs of 1,819, 1,730 and 3,446 chars carry the docs brief, the prompts brief and the whole `--fix` loop. |
| WS2 | med | `skills/handoff/SKILL.md:125-135` | Two dated operator quotes inside the fenced handoff template. |
| WS3 | med | `skills/init-project/SKILL.md:112,115` | Step 5b paragraphs of 1,248 and 2,456 chars; `python3` dependency for the JSON write and parse. |
| WS4 | med | `skills/distill/SKILL.md:4,16,38,176`; `feedback-log:14` | `--from-log` advertised, no step reads the log; `--global` in usage not in argument-hint. |
| WS5 | med | `skills/dream/SKILL.md:70` | "Run on Antigravity — it reads Claude Code's store" contradicts `:22` (one store, R10). |
| WS6 | low | `skills/distill/SKILL.md:156-170` | Step 6.3 writes a synthetic feedback entry for the agent's own act. |
| WS7 | low | `skills/review-standards/SKILL.md:3,8` | "loaded by /review" is wrong (review never reviews itself); `/preflight` is right (reads small diffs itself). Drop `/review` only. |
| WS8 | low | `feature:25`, `bug:49`, `review-standards:37`, `preflight:77`, `retro:49,141`, `feedback-log:58,86`, `init-project:91`, `review:18` | Scar tissue: another project's name (styrman-harness), operator quotes and BUG ids inside global skills. |
| WS9 | low | `skills/preflight/SKILL.md:93` | GCP shape (amd64, Artifact Registry, deploy) prescribed for all three deploy guides. |
| WS10 | low | `skills/research/SKILL.md:47,22` | Always writes `specs/<app>/integration-spec.md`; self-lookup uses `--kind plan` while it writes into specs (`spec`). |
| WS11 | low | `skills/VENDORED.md` | Note that find-skills "still names skill-repair" is stale (0 hits). find-skills' `-g -y` step is verbatim vendored and overridden by the context: not a skill defect. |
| WS12 | low | `skills/feedback-log/SKILL.md:16,84` vs `:90`; `handoff:34,103` | Rule 1 "never mention logging" contradicted by rule 7 and by handoff quoting log timestamps. |
| WS13 | low | `feature:78,90`, `bug:53`, `preflight:40` | Four dispatches to reviewer/verifier/security lack the budget and "write the artifact early" clause the context requires. |
| WS14 | med | `skills/review/SKILL.md:60,64`; `preflight` Step 4 | `documentation` (no Bash) dispatched without the `sh-agent search --refs` output that feature and bug carry. |
| WS15 | med | `skills/preflight/SKILL.md:40` | Step 1 "brings dirty work through the DoD" but runs tests, changelog and a conditional reviewer only; documentation and bugs skipped. |
| WS16 | med | `skills/preflight/SKILL.md:50,79` | Reviewer row accepts a covering review by `Range reviewed:` alone; only the verifier row requires Status PASS. Five committed review artifacts are `Status: IN PROGRESS`. |
| WS17 | low | `skills/distill/SKILL.md:38`, `skills/retro/SKILL.md:45` | Both dispatch dreamer to the same artifact name; no reuse rule for a session running both. |
| WS18 | low | `skills/handoff/SKILL.md:166` | 1,204-char sentence with a 200-char awk program inline for the cold-reader check. |
| WS19 | low | `skills/review/SKILL.md:52` | `--prompts` brief hardcodes the two installed skill paths with no fallback. |
| WS20 | low | `skills/feature/SKILL.md:82` | Never says whether a re-dispatched cold reader overwrites its artifact; review says overwrite. |
| WS21 | low | `skills/grilling/SKILL.md:26` | "dispatch a sub-agent" with no name; Antigravity needs one. Vendored: rules line or fork. |
| WS22 | low | `skills/feedback-log/SKILL.md` | Writer should record the project root (it does today; keep it pinned). |

## TS — technical and vendored skills

| ID | Sev | Where | Finding |
|---|---|---|---|
| TS1 | high | `skills/terraform/` | 100% AWS (175 AWS identifiers, 0 Google, *checked*); CI, pre-commit, MCP demands; invalid HCL (`for_each` + `count`); contradictory `.terraform.lock.hcl` advice; attribution footer demand; SKILL.md:388 promises a LICENSE that is absent. |
| TS2 | high | `skills/golang/references/` | ≈30 dangling pointers (12 to `assets/examples/*.go`, *checked*; "See Makefile template", "Nil Safety Deep Dive", "HTTP Testing", "Transactions", `templates/*`, `samber/cc-skills@promql-cli`). |
| TS3 | high | `skills/golang/`, `skills/VENDORED.md` | Flattened samber/cc-skills golang skills (29 persona blocks, "5 parallel sub-agents" ×8, Context7 ×9, `golang-samber.md` 52 KB); no VENDORED row (*checked*), no LICENSE. |
| TS4 | high | `skills/anti-ui-slop/LICENSE` | Frontmatter and VENDORED say MIT; the file is Apache-2.0 text. |
| TS5 | high | `skills/frontend/references/theme-system.md:217-239,346-393,408`; `accessibility.md:258` | Tailwind v3 (`tailwind.config.ts` with `content:`) and Vite SPA (`index.html`, `main.tsx`) while SKILL.md pins Tailwind 4 `@theme` and TanStack Start `__root.tsx`. |
| TS6 | med | `golang-testing.md:252`, `golang-io-web.md:196,413,430-432,668,674,901-910`; `bash/references/makefile-patterns.md:216-218` | `is.NotNil`, testify, `uuid.NewString`, "sqlx or pgx" MUST, Cobra/Viper table, mockgen rule against SKILL.md:44. CLI preamble fixed. |
| TS7 | med | `golang-testing.md:141-186` | `synctest.Run` marked experimental; GA in 1.25 as `synctest.Test`; contradicts `golang-architecture.md:725`. |
| TS8 | med | `golang-ecosystem.md:240-256,335,378,379,413,417` | tools.go pattern; golangci-lint v1 paths and keys beside v2 commands. |
| TS9 | med | `golang-ecosystem.md:13-113,1081-1095` | Observability is Prometheus/Grafana/Datadog/Sentry; nothing on Cloud Logging JSON contract, Cloud Trace, Profiler, Error Reporting. |
| TS10 | med | `golang/references/*` | Off-harness workflow: personas ×29, `.modernize` in project root, `/golang-modernize`, CI badges, Dependabot, `perf(scope):`. |
| TS11 | med | `terraform/references/` | Covered by TS1. |
| TS12 | med | `gke-basics/references/cli-reference.md:5-11`, `gke-golden-path.md:11,68`, `gke-cluster-creation.md:132`, `mcp-usage.md:72`, `gke-security.md:117`, `golden-path-autopilot.yaml:3,127` | "MCP > gcloud"; one file says never ask for the project, another ALWAYS ask; two dangling asset paths. |
| TS13 | med | `anti-ui-slop/SKILL.md:3` vs `frontend/SKILL.md:4-10` | Both claim every web UI task; neither routes to the other; `operate.md` duplicates frontend § Apps & Dashboards. |
| TS14 | med | `frontend/references/component-patterns.md:199-238` | react-hook-form + resolvers + zod unjustified; `z.string().email()` deprecated; `createServerFn` against SPA restriction. |
| TS15 | med | `frontend/references/theme-system.md:59,163,243-337`; `accessibility.md:21-37`; `animation-recipes.md` | Dark focus ring ≈1.8:1 against the 3:1 rule; examples hard-code zinc/amber; ThemeToggle uses inline style objects against the Tailwind-and-tokens rule. |
| TS16 | med | `cloud-run-deploy/SKILL.md:67-95,120`; `docker/references/health-checks.md:40-41` | No `--no-traffic --tag` canary, no startup probe, ingress, VPC, Cloud SQL flags; `--platform managed` redundant; verify curls unauthenticated while IAM-gated is the default; docker's probe hint points at a flag the deploy skill never names. |
| TS17 | low | `docker/SKILL.md:60,75,114`; `gke-app-onboarding.md:28` | `golang:1.24` and `golang:1.22` base pins; box runs 1.27. |
| TS18 | low | `gcloud-setup/SKILL.md:14,24,26,33-36` | brew line is a cask; heading promises Linux, no Linux path; two "check in Cloud Console" answers in a CLI skill. |
| TS19 | low | `json-api-style/SKILL.md:13,32,55`; `references/json-api-conventions.md:183-287` | 2012 Google JSON Style Guide (`apiVersion`, `data`, `error.errors[].domain`) not AIP; camelCase vs the golang skill's snake_case struct tags. |
| TS20 | low | `bash/references/bash-fundamentals.md:43-44,419,503-517,634-635`; `makefile-patterns.md:128,131,132` | SC2155 pattern; SC2046 nested substitution shown as working; `/var/log` logging; `IMAGE_TAG` twice; `$(error)` at parse time. |
| TS21 | low | `codebase-design/references/DESIGN-IT-TWICE.md:30` | "CONTEXT.md" from upstream's sibling; vendored verbatim, note in VENDORED. |
| TS22 | low | `anti-ui-slop/NOTICE:7` | `reference/ios.md` vs `references/`; vendored verbatim. |
| TS23 | low | `python/SKILL.md:3-7` | No Python in the stack; broad trigger. |
| TS24 | med | `golang-core.md:61,230,628`; `golang-architecture.md:45,15,146`; `golang-samber.md:139,192,243,627-641,720-734` | `wg.Go` 1.24 vs 1.25; pre-1.24 map internals; discarded `sql.Open` error; "never nil slices" MUST vs zero-value rule; "4+ args one per line even when the prompt asks otherwise"; panic in a default branch; `lo` over the loop; two non-compiling samples. |
| TS25 | med | `frontend/SKILL.md:43` vs `docker/SKILL.md:87-90` | `dist/client` committed with no pin test versus a Node build stage; not reconciled. |
| TS26 | med | `docker/SKILL.md:80` vs `docker/references/setup.md:8-18` | `.dockerignore` list enumerated twice, differing. |
| TS27 | low | `gke-basics/references/client-library-usage.md:51-54`, `iac-usage.md:30-31` | `context.TODO` uncalled, errors discarded; `deletion_protection = false` in the only GCP Terraform sample. |
| TS28 | low | `mermaid-syntax/SKILL.md:167,373`; `other-diagrams.md:398-403` | Relies on `%%{init}` its own secure-settings rule blocks; pre-11 `packet-beta`; no provenance row. |
| TS29 | low | `typescript/SKILL.md` | 2.4 KB, no references, nothing on current tsconfig, test runner, or the plain-V8 restatement. |

Missing technical knowledge for this stack: Cloud SQL Postgres + migrations; API auth on GCP (IAP, ID tokens, Identity Platform); observability on GCP (Cloud Logging JSON contract, Trace, Profiler, Error Reporting); testing strategy across Go + web + container; Terraform for GCP; Cloud Run runtime and cost, Jobs, Scheduler, Tasks; release hygiene without CI; Node/npm hygiene for `web/`; i18n. Dead weight: python; anti-ui-slop (paid MCP unconfigured, duplicates frontend); `golang-samber.md`; terraform's AWS and CI references; gke-basics inference/batch/multitenancy references.

## AG — subagents and their dispatch

| ID | Sev | Where | Finding |
|---|---|---|---|
| AG1 | high | `agents/changelog.md:83-84` | "Read the actual diff" and "run `git show --stat`" with no Bash (tools: Read, Grep, Glob, Edit, Write); BUG-089's rule not applied one bullet down. |
| AG2 | med | `agents/documentation.md:99,131,146` | Claims ROADMAP Next/Later; context DoD step 4 keeps ROADMAP empty at checkpoints. |
| AG3 | med | `agents/documentation.md:133,147` vs `:157` | Claims godoc while "Never: source". |
| AG4 | med | `context/AGENTS.md:16` | DoD step 2 sends all `AGENTS.md` docs findings to the main agent; review and preflight already say "other `AGENTS.md` sections". |
| AG5 | med | `agents/reviewer.md:36`; `derive.go:30-33,54`; `harness_test.go:42,130-161` | review-standards named with no path; `skills:` dropped for Antigravity on an untested comment; the derive test rewrites a fixture path the real file lacks; no skill-matrix row. |
| AG6 | med | `agents/verifier.md:35,37,43,55,98` | 2,053 words, 27 commits; guide duty stated four times; dated incidents inline; "10 minutes on Claude Code" with Antigravity's limit unnamed. |
| AG7 | med | `agents/changelog.md:5-6,72-79,90`; `changes/` | haiku/low writes ADRs (BUG-061 graduated by pointing at the same text); 126 of 172 entries are `owner: main`. |
| AG8 | low | `agents/bugs.md:20,165,168,243` | "Read the diff yourself" with no Bash; two files at maxTurns 20 with no write order; "dispatch prompt" in the body (rule covers descriptions only). |
| AG9 | low | `review:44`, `preflight:50` | "on Antigravity subagents inherit the session model": true for reviewer, false for security (derive pins `pro`). |
| AG10 | low | reviewer/security/verifier/explorer/dreamer | "Write the artifact EARLY" block identical in two, variant in two, short in one, absent in bugs, changelog, documentation. |
| AG11 | low | `agents/explorer.md:40,67`; `context:9` | Recon briefs land in `plans/` (21 of 63 files), which resume reads as "the active plan"; only the name distinguishes. |
| AG12 | low | `internal/harness/inventory.go:245`, `frontmatter.go:102` | No test pins "one line, no 'dispatch prompt' / 'harness binary'" for descriptions. |
| AG13 | low | `agents/reviewer.md:25`, `security.md:131` | An IN PROGRESS artifact is "a failed dispatch", yet such files are committed with no void marker. |
| AG14 | low | `changes/2026-09-30-2205-*:29-30` | "docs lens" counted as a dispatched agent. |
| AG15 | low | `derive.go` | `effort` and `maxTurns` dropped for Antigravity: security's max effort and every turn ceiling lost there; nothing states it. |

Lane mismatches: ROADMAP (AG2), godoc (AG3), AGENTS.md sections (AG4). Merge candidate: changelog into DoD step 4. Missing: a cold reader for the trail (resume check) as a `/review` lens; the subagent definition does not name "a record written outside the main context" though three agents exist for it.

## IN — installer Go code (cmd/sh-agent, internal/harness)

| ID | Sev | Where | Finding |
|---|---|---|---|
| IN1 | med | `internal/harness/report.go:66-103`; `report_test.go:151-162` | Over-redaction (*repro*): `-b <word>` (so `git checkout -b feature/x`), `## Secrets:`, `?sig=<16>`, `token:` and `passwords:` in prose all vanish; none in the benign near-miss list. `Redact` feeds `sh-agent dream`. |
| IN2 | med | `internal/harness/apply.go` Install | Files written one by one, `SaveState` only at the end; a mid-run failure leaves files no manifest lists. |
| IN3 | low | `apply.go` Install unchanged branch | `os.ReadFile/Stat/Chmod(dest)` follow a symlinked destination while every other path uses `Lstat`. |
| IN4 | low | `cmd/sh-agent/main.go:690-719` | `update --dry-run` returns exit 2 before the dirty-tree refusal; promises a pull the real run refuses. |
| IN5 | low | `main.go:102-127,848` | One FlagSet for every verb (`list --check --force` accepted); hand-kept `takesValue` map. |
| IN6 | low | `repo.go:81-87`; `main.go:316,536,573,789,792,1228`; `apply.go:117` | Manifest origin string parsed by hand at six sites in two packages. |
| IN7 | low | `main.go:449-466` | `report` re-implements source resolution. |
| IN8 | low | `apply.go`, `release.go` | No lock across concurrent installs; `internal/dream` has a flock helper. |
| IN9 | low | `main.go:885-900,1057-1069` | `pathHint`, `attributionHint` return an always-nil error. |
| IN10 | low | `release.go:451`; `main.go:1173` + `zip.go:92` | System temp dir for the signers file and the no-config dry-run unpack. |
| IN11 | low | `main.go:17,64` | dream timer `--dry-run` missing from usage and synopsis. |
| IN12 | low | `main.go:200` | `--install-timer --remove-timer` together exits 1; every other dream usage error exits 2. |
| IN13 | low | `main.go:225,507` | dream and memory refusals written to `os.Stderr`, not the writer tests see. |
| IN14 | low | `main.go:~712` | `update` calls `Inspect` three times. |
| IN-S | — | | Simplifications: S1 move `attribution*` settings code to `internal/harness`; S2 split `run` (324 lines) table-driven; S3 collapse verbs (`sync` alias, `render`, `diff`, `index`); S4 delete `adopt` (177 lines); S5 cut incident-narrative comments (27/19/7 in report/main/repo); S6 merge PATH walkers; S7 one `readCapped`; S8 `Lint` as a flag on `Parse`. |

Coverage: cmd/sh-agent 84.6%, internal/harness 87.7%; build, vet, gofmt, govulncheck clean.

## SD — search, memory, dream

| ID | Sev | Where | Finding |
|---|---|---|---|
| SD1 | high | `internal/search/dream.go:330-339` | Reader tries the text leaf before the nested parse (*repro*, *checked*): a sub-message holding one 32–126-byte string passes `textField` with key and length bytes glued on; `prose()` fails the JSON guard on the glued byte; a tool argument reaches the dreamer. Fix: nested first, text only when nested fails; two shapes into `TestProseLeaves`; re-measure the 971-leaf corpus. |
| SD2 | med | `internal/search/search.go:533,577` | Duplicate hits in multi-word search (*repro*): `have` keyed by the offset-shifted line; OR round compares unshifted; same chunk printed and counted twice. Fix: key by chunk id. |
| SD3 | med | `internal/dream/run.go:427-428` | Proposal `text` written multi-line and unfenced (*repro*); a model can forge `### N.` / `**Status:** queued` sections `/retro` parses. |
| SD4 | low | `internal/dream/collect.go:80,199` | "Changed since" is file mtime; a clone marks the whole trail new; `fit()` drops transcripts first. |
| SD5 | low | `internal/dream/run.go:333` | `DisallowUnknownFields`: one extra key = nightly failure retried with the same prompt. |
| SD6 | low | `internal/search/search.go` reconcileGlobal | `--global` reconciles every registered checkout's full tree per query; spec says 15 ms. |
| SD7 | low | `internal/search/dream.go:242-267` | db/-wal/-shm copied file by file without a lock; a torn global-index copy fails the whole run. |
| SD8 | low | `internal/dream/prompt.go:434` | `<<<data … >>>` fence not escaped. |
| SD9 | low | `internal/search/dream.go:363-372` | `prose()` guards JSON only; a non-JSON string argument passes; the spec's assumption is unpinned. |
| SD10 | low | `internal/dream/collect.go:353-366` | Claude Code slash-command output (`<local-command-stdout>`) arrives as user string content and reaches the prompt. |
| SD11 | low | `internal/dream/run.go:202-233,85` | A stdout line over 16 MiB logged as "invalid output"; `removeCopies` uses wall clock, not the injected one. |
| SD-spec | med | `specs/sh-agent/{search,dream,memory-loop,delivery}.md`; the parity plan | search.md 15 ms ignores `--global`; dream.md "first JSON that decodes" vs unknown fields refused, "newer than watermark" vs mtime, "zero proposals write no file" vs dropped-draft case; memory-loop.md 25 KB for add (enforced only in show), frontmatter list omits `node_type`; **delivery.md:77 claims memory and dream "built and shipped"** (newest tag v1.16.2, *checked*); plan's claude command lacks `--strict-mcp-config --setting-sources local`. |

Coverage: search 83.9%, memory 89.1%, dream 89.1%.

## TR — trail and learning loop

| ID | Sev | Where | Finding |
|---|---|---|---|
| TR1 | high | `bugs/INDEX.md` | Prevention rules recur after graduation: restatement-grep 11 (+ delivery.md today), run-proves-its-setup 11 (BUG-123 added), duty-in-dispatch 5, two-rounds-stop 4, check-and-side-effect 4 (BUG-122 added). Each recurrence added a longer paragraph, never a mechanism. BUG-122/123 parked as "unverified" in the handoff instead of a retro. |
| TR2 | med | `handoffs/2026-09-30-2206-*` § Environment State | Says the real home was installed from a275b96 and lacks the decoder; `sh-agent status` says installed 22:03 from 032ece1. The real-home log cited as PASS covers v1.16.2..a275b96. |
| TR3 | med | `changes/2026-09-30-2145-*` | Time-named files not time-ordered (committed 20:33, before the 2120 and 2135 entries). BUG-073 class. |
| TR4 | med | `reviews/` | Five artifacts `Status: IN PROGRESS`, verdict pending, no void marker; preflight can cite one as coverage (WS16). |
| TR5 | low | `plans/` | 63 files, 37 without a checkbox; `2026-09-15-harness-practices-g01-g26.md` has 14 unticked boxes and is what `--open` would list as open while the handoff says nothing is queued. |
| TR6 | low | `~/.config/sh-agent` | The nightly dreamer has never run on this box: no timer, no state file, no `proposals/dream/`. |
| TR7 | info | `.agents/docs` | 7.72 MB; reviews 3.71 MB (48%); trail:source 23.9:1; 379 commits in 16 days. |

## G — workflow coverage gaps

| ID | Stage | Today | Gap |
|---|---|---|---|
| G1 | Requirements | template + grilling | nothing elicits or traces requirements; template opens with SLA tables |
| G2 | Spec | research (one fixed shape) | no spec skill or template; `assets/agents-tree/docs/specs/` is a .gitkeep |
| G3 | Design / ADR | changelog on haiku | no skill says what is structural or shapes the ADR |
| G4 | Plan / tasks | feature | nothing sequences several plans; Task class has no breakdown help |
| G5 | Test strategy | verifier writes the guide when absent | no skill on tiers, fixtures, e2e choice; no guide template |
| G6 | Operate | nothing | no runbook, logs, monitoring, rollback outside a release |
| G7 | Security / dependencies | security per diff | no govulncheck / npm audit cadence; no API auth guidance |
| G8 | Data | nothing | no SQL, migrations, Cloud SQL |
| G9 | Onboarding / explain | README | nothing explains a codebase to a person |
| G10 | Learn | 5 skills + 2 agents | over-weighted; the automated half never ran on a real home |
| G11 | Size tiers | one workflow | a dotfiles repo owes the same paperwork as a service; nothing for multi-service |
| G12 | Trail cold reader | none | nothing checks that resume works for a zero-context reader |
| G13 | Templates | assets | no README, spec, testing-guide, plan, handoff, change, bug, review, retro templates outside skills and agents |
