---
owner: main
date: 2026-09-29
---
# Adoption and efficiency — clean A/B of the whole harness, with and without `sh-agent search`   (2026-09-29)

Plan: `plans/2026-09-29-search-release-readiness.md`, Phases 1 and 2. It replaces the setup of `testing/2026-09-29-search-efficiency.md`, whose two arms shared the installed v1.12.1 skills and carried the context twice.

## Status: RECORDED — same quality in every arm; no cost saving overall; with the resume step the main agent searches in 6 of 6 large-trail sessions

## Setup

Claude Code 2.1.284, headless, `--setting-sources project,local`. A probe showed what that loads: a project skill is listed, no user-level skill is, and the context's canary word is absent. Each project carries its arm's harness in `.claude/{skills,agents}`, copied from `sh-agent install --vendor claude` of that arm's tree and excluded from git. The arm's rendered context is appended once. `sh-agent` is a wrapper pinning `--home` to a sandbox home (BUG-107). Projects live in `~/sandbox/search/eff/cwork/`, outside every checkout (BUG-106). Kit: `.agents/scratch/eff-kit/` (`clean-session.sh`, `lane-clean.sh`, `probe.sh`, `metrics-clean.py`, `score-trail.py`).

| Arm | Harness |
|---|---|
| base | v1.12.1 |
| new | this branch before the resume step: the "Search first" entry, the agents' and skills' lines |
| new2 | the same plus `sh-agent search --open` in the resume reading order (plan D9) |

**Large-trail lane:** a clone of this repository at v1.12.1 (about 1,000 files, a trail of 850) with one seeded drift. It runs three fresh sessions:

1. **s1:** resume, then fix the drift as a Task under the Definition of Done.
2. **s2:** resume, then answer five trail questions with `file:line`.
3. **s3:** resume, then add the Go test `TestUsageMatchesDocComment`.

**Sample-app lanes:** as in the first comparison (build, add a feature, trail check). Two runs per arm per lane.

## Large-trail lane — quality

Every session delivered: 18 of 18. In all six lanes the drift is fixed, the answers score 5 of 5, and the new test runs and passes.

## Large-trail lane — cost and search use (median, range; n = 2 per row)

| Session | Arm | Input tokens | Cost USD | Minutes | Read chars | Main-agent searches | Subagent searches |
|---|---|---|---|---|---|---|---|
| s1 fix drift | base | 2.5M (2.4M–2.6M) | 0.85 (0.85–0.86) | 1.6 | 40k (40k–41k) | 0 | 0 |
| s1 fix drift | new | 3.4M (3.2M–3.7M) | 1.09 (1.08–1.09) | 2.3 | 58k (54k–63k) | 0 | 0 |
| s1 fix drift | new2 | 4.0M (3.5M–4.6M) | 1.16 (1.07–1.26) | 2.3 | 49k (49k–50k) | 2 | 0 |
| s2 trail questions | base | 1.8M (1.6M–1.9M) | 0.66 (0.63–0.70) | 1.5 | 27k (25k–29k) | 0 | 0 |
| s2 trail questions | new | 2.0M (1.8M–2.1M) | 0.72 (0.69–0.76) | 1.4 | 37k (28k–47k) | 1 | 0 |
| s2 trail questions | new2 | 1.5M (1.2M–1.7M) | 0.58 (0.55–0.60) | 1.3 | 12k (0k–23k) | 3 | 0 |
| s3 add a test | base | 2.0M (1.9M–2.1M) | 0.70 (0.69–0.71) | 1.9 | 21k (19k–24k) | 0 | 0 |
| s3 add a test | new | 2.5M (1.9M–3.0M) | 0.83 (0.69–0.96) | 2.4 | 46k (32k–61k) | 2 | 0 |
| s3 add a test | new2 | 2.4M (1.6M–3.2M) | 0.80 (0.61–0.99) | 2.3 | 47k (35k–59k) | 3 | 0 |

**Adoption (plan D3, bar 4 of 6):** sessions in which the main agent ran `sh-agent search` were 3 of 6 for `new` and **6 of 6 for `new2`**. With v1.12.1 skills in both arms, the first comparison measured 0 of 28.

## Reading it

- **The resume step works on a real-sized trail.** Every `new2` session ran `sh-agent search --open` while resuming, and the trail-question sessions then searched by words.
- **Where search fits, it saves.** In the trail-question session `new2` read 12k characters against the baseline's 27k and cost 0.58 against 0.66 USD, with the same 5 of 5 answers.
- **Over the whole lane it does not save.** The search arms cost 15 to 19% more than the baseline (2.54 and 2.64 against 2.21 USD per lane, from the medians), most of it in s1. There the search arms' agents checked the new bug number against the repository's other branch and renumbered it (BUG-105 to BUG-108); the baseline wrote BUG-105 without checking. That is more careful work, and it is not attributable to search: the drift itself was found with grep in every arm, which is the right tool for an exact string.
- **n is 2 per cell.** The s1 and s2 ranges do not overlap between base and `new2`; the s3 ranges do.

## The probe on a small fixture

The `search-first` question ("resume, then: where did this project decide how sessions are stored…") ran on the 65-file `notebox` fixture, in clean mode. `new2` searched in 0 of 3 sessions and base in 0 of 3; all six answered correctly in 15 to 21 seconds with `cat` and one grep. On a trail of a few files the resume step is skipped, and nothing is lost by it.

## Sample-app lanes

Every session delivered: 24 of 24. All eight lanes pass `go vet` and `go test`, have the session-2 feature in their code, hold three handoffs and `ANSWERS.md`, and pushed to their origin. Arms: base and `new` (the app lanes ran before the resume step existed).

| Session | Arm | Input tokens | Cost USD | Minutes | Read chars | Main-agent searches | Subagent searches |
|---|---|---|---|---|---|---|---|
| s1 build | base | 39.0M (28.9M–50.5M) | 9.43 (8.27–9.88) | 29.3 (26.4–31.8) | 391k (308k–483k) | 0 | 0 |
| s1 build | new | 41.1M (28.0M–47.6M) | 9.11 (7.94–11.16) | 28.5 (25.9–33.7) | 387k (313k–581k) | 13 | 19 |
| s2 feature | base | 22.3M (19.5M–27.6M) | 6.94 (5.09–7.45) | 20.9 (17.0–25.1) | 388k (330k–410k) | 0 | 0 |
| s2 feature | new | 25.5M (18.0M–41.7M) | 6.67 (5.85–8.79) | 20.6 (18.7–27.2) | 374k (346k–461k) | 14 | 17 |
| s3 trail check | base | 5.1M (3.4M–9.2M) | 1.40 (1.18–2.46) | 3.9 (2.9–6.8) | 79k (46k–191k) | 0 | 0 |
| s3 trail check | new | 7.1M (4.1M–13.3M) | 2.07 (1.22–3.18) | 5.3 (3.3–8.6) | 107k (48k–169k) | 1 | 3 |

n = 4 per row (two apps, two runs). Per lane of three sessions:

| Arm | Lane costs USD (linkshelf 1, 2; tasklog 1, 2) | Median | Minutes (median) |
|---|---|---|---|
| base | 19.30, 17.78; 14.54, 18.25 | 18.02 | 56 |
| new | 18.71, 20.52; 16.97, 17.63 | 18.17 | 56 |

**No measurable difference in cost or time**, and every range overlaps. With the new skills in place the main agent ran `sh-agent search` in 9 of 12 sessions (28 calls) and the subagents made 39 calls; the first comparison, with v1.12.1 skills in both arms, had measured 0 main-agent calls in 28 sessions.

## Fixture notes

- The large-trail clone keeps a local `main` branch at this repository's HEAD beside the `v1.12.1` work branch, so later bug numbers are visible through git. That is what the search arms' renumbering reacted to. It is the same in every arm.
- `metrics-clean.py` first crashed on its own output file while its stderr was discarded, and an early adoption line (1 of 4) was stale. It was re-run with errors visible; the numbers above are from that run.

## What was not measured

- Antigravity: its app lanes did not compare in the first round (print mode ends before subagents return), and it records no token counts. The gate runs of `testing/2026-09-29-search-agy-gate.md` cover whether each changed skill and agent completes there.
- Days of work in one project; interactive sessions.
- Significance: no test would hold at n = 2.
