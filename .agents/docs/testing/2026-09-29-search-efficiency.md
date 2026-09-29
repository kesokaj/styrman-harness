---
owner: main
date: 2026-09-29
---
# Efficiency — sample apps and a larger project, with and without `sh-agent search`   (2026-09-29)

Plan: `plans/2026-09-29-search-index.md`, Phase 5. Asked for by the operator: "run sandboxed projects that test the efficiency when building sample apps".

## Status: RECORDED — on Claude Code, no measurable efficiency difference, with the same quality in both arms; the main agent never used search, only subagents did. Antigravity's app lanes are not comparable (print mode ended before subagents returned).

## Setup

Kit: `.agents/scratch/eff-kit/` (gitignored, ADR-010). Projects: `~/sandbox/search/eff/work/`, outside every checkout (BUG-106).

| | Baseline arm | Search arm |
|---|---|---|
| Harness | v1.12.1 | this branch (as built at 5a92849, before the security fixes; those change registration, hashing and output hygiene, not ranking or cost) |
| Context | the real home's installed v1.12.1 context, plus the v1.12.1 render appended | the same installed context, plus the new render appended (both arms carry the same duplicate-context overhead) |
| Agents | all eight v1.12.1 definitions through `--agents` | all eight working-tree definitions through `--agents` |
| Skills | the real home's v1.12.1 skills | the same (the skills' new lines are not in this measurement) |
| `sh-agent` | the v1.12.1 build, through a wrapper pinning `--home` to a sandbox home (BUG-107) | the new build, through the same wrapper |
| Antigravity | a sandbox home installed from v1.12.1 | a sandbox home installed from the working tree |

Each app lane starts from the same scaffold (v1.12.1's `/init-project` assets, one commit) with a local bare `origin`. It runs three fresh sessions in order:

1. **s1:** build the app through `/feature` (gates waived by name; Definition of Done through `/handoff`, push).
2. **s2:** resume from the trail and add a feature (linkshelf: tags; tasklog: `edit`).
3. **s3:** resume and write `ANSWERS.md`, listing every decision and bug record touching storage, each prevention rule and its covering test, then add missing tests.

Apps: linkshelf (the 2026-09-16 brief) and tasklog (a Go CLI). The larger-project task runs in a clone of this repository at v1.12.1 with one seeded drift (`deploy/remote.md` says `--check` exits 3). It asks five trail questions in `ANSWERS.md` and for every restatement of the exit status to be made consistent. Claude Code: 2 runs per arm per app, plus 2 per arm on the larger task. Antigravity: 1 run of each.

## Claude Code (2.1.284) — quality

Every session succeeded: 28 of 28. Every app lane in both arms passes `go vet` and `go test`, has the s2 feature in its code, holds three handoffs and `ANSWERS.md`, and pushed to its origin. All four larger-task runs answered 5 of 5 and fixed the seeded drift in one commit.

## Claude Code — cost and time (median, range; n = 4 per session row, 2 per larger-task row)

| Session | Arm | Input tokens | Uncached input | Output | Cost USD | Wall minutes | Read chars | `sh-agent search` calls |
|---|---|---|---|---|---|---|---|---|
| s1 build | base | 47.8M (25.6M–98.4M) | 1531k (988k–2502k) | 440k (270k–992k) | 9.27 (6.03–14.38) | 25.7 (15.1–35.8) | 309k (167k–532k) | 0 |
| s1 build | search | 43.5M (28.7M–60.7M) | 1420k (1226k–1688k) | 364k (277k–433k) | 8.82 (6.60–12.56) | 22.5 (15.1–31.4) | 247k (197k–374k) | 22 |
| s2 feature | base | 32.9M (18.4M–55.8M) | 1348k (1106k–1749k) | 318k (173k–458k) | 6.67 (4.58–9.80) | 16.7 (11.1–23.9) | 273k (219k–359k) | 0 |
| s2 feature | search | 23.3M (18.0M–49.6M) | 1114k (939k–1706k) | 226k (169k–328k) | 5.49 (4.73–9.10) | 13.2 (9.7–19.0) | 227k (191k–347k) | 11 |
| s3 trail check | base | 12.8M (5.9M–18.5M) | 546k (250k–632k) | 94k (56k–142k) | 3.02 (1.75–3.62) | 6.3 (3.9–8.9) | 136k (74k–197k) | 0 |
| s3 trail check | search | 11.2M (5.7M–22.4M) | 440k (288k–828k) | 129k (63k–174k) | 2.78 (1.72–4.11) | 7.2 (4.0–8.8) | 83k (68k–174k) | 0 |
| larger project | base | 4.4M (3.2M–5.5M) | 302k (265k–340k) | 28k (28k–29k) | 1.25 (1.16–1.34) | 1.9 (1.6–2.1) | 21k (19k–22k) | 0 |
| larger project | search | 4.4M (4.1M–4.7M) | 160k (142k–178k) | 25k (24k–26k) | 1.28 (1.26–1.30) | 2.5 (2.0–2.9) | 33k (18k–48k) | 0 |

Input tokens include cache reads (most of the total); uncached input is fresh input plus cache writes; cost is Claude Code's own `total_cost_usd`. Main session and subagents are summed.

| Lane (s1+s2+s3) | Cost USD | Input tokens | Wall minutes |
|---|---|---|---|
| base linkshelf 1 / 2 | 23.32 / 26.77 | 121.9M / 165.0M | 65 / 65 |
| search linkshelf 1 / 2 | 24.13 / 22.04 | 127.8M / 101.4M | 53 / 55 |
| base tasklog 1 / 2 | 13.75 / 14.23 | 60.0M / 62.5M | 31 / 35 |
| search tasklog 1 / 2 | 13.51 / 13.32 | 58.3M / 53.7M | 31 / 35 |
| Median lane, base / search | 18.78 / 17.78 | | |

**Who searched:** all 33 `sh-agent search` calls came from subagents: `security` 22, `reviewer` 11. The main agent made 0 in 28 sessions, including the four trail-question sessions and the larger-project task, where search should help most.

## Antigravity (agy 1.2.13)

| Session | Arm | Minutes | Steps | Search output lines | Delivered |
|---|---|---|---|---|---|
| linkshelf s1–s3 | base | 13.6 / 18.6 / 5.1 | 411 / 769 / 526 | 0 | yes: tests pass, tags in code, handoff, `ANSWERS.md`, pushed |
| linkshelf s1–s3 | search | 9.0 / 1.6 / 5.8 | 221 / 198 / 206 | 18 / 18 / 14 | no: tags never built, no handoff (2 commits) |
| tasklog s1–s3 | base | 4.4 / 2.5 / 2.8 | 97 / 66 / 120 | 0 | void: s1 left no `go.mod`; s2 and s3 auto-denied again (third attempt); not re-run, since the app lanes do not compare anyway |
| tasklog s1–s3 | search | 6.2 / 3.7 / 6.7 | 307 / 216 / 266 | 17 / 3 / 9 | no: `edit` never built, no handoff (3 commits) |
| larger project | base | 5.7 | 256 | 0 | 5 of 5, drift fixed |
| larger project | search | 4.9 | 328 | 9 | 5 of 5, drift fixed |

The search arm's app sessions ended with "Verifier subagent dispatched … Standing by for its report": `agy -p` returned before the dispatched subagent did, so those sessions did less work. That explains their shorter times. The vendor limit is the one Phase 3 recorded, and it struck the two arms unevenly, so **the Antigravity app lanes do not compare**. The larger-project pair does: equal quality, 4.9 against 5.7 minutes, and more steps with search (328 against 256).

Grants: two agy batches were void, their sessions auto-denied a command the scenario needed (`which`, `rg`, `tee`, `true`, `chmod`, `ps`, `uname`, `false`, the app binaries). Each was added by name, per the project `AGENTS.md` § Testing rule, and the affected lanes re-ran. The rows above are from runs with no "auto-denied" line.

## Reading it

- **Efficiency on sample apps: no measurable difference.** The search arm's median session cost is 5 to 18% lower in s1 to s3, and its median lane cost 5% lower. Every range overlaps the baseline's, and n is 2 per app per arm. Quality was identical.
- **The larger project: no difference either.** The search arm's cost was the same (1.28 against 1.25 USD), its uncached input lower (160k against 302k), and its answers the same, 5 of 5. The main agent answered the trail questions with grep and reads in both arms.
- **Where the index is used, it is used by the agents told to use it as a fixed step** (`security`, `reviewer`). The context entry reached the main agent in the direct probe (2 of 3, `testing/2026-09-29-search-index-behavioural.md`) but not in real work (0 of 28).
- **Small trails cap the upside.** A new sample app has a trail of a few files, and grep over it is instant. This repository's own recall test (`testing/2026-09-29-search-recall.md`) is where the index is clearly ahead: 0.677 against 0.379 recall on paraphrased and structured questions.

## What was not measured

- More runs per cell; no significance test was run, and none would hold at n = 2.
- A project with a trail of hundreds of records being worked on for days, the case the operator named ("a powerhouse in larger projects"). The larger-project task was one short session.
- The skills' new lines (both arms used v1.12.1 skills), and the security fixes made after the kit was built.
- Antigravity tokens: its conversation store is undocumented protobuf, and its log holds no token counts.
