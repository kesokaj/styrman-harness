---
owner: main
date: 2026-09-29
---
# Recall baseline — grep, before `sh-agent search`   (2026-09-29)

Plan: `plans/2026-09-29-search-index.md`, Phase 0. Phase 4 adds the index arm to the same test and the same questions.

## Status: BASELINE RECORDED

## Command

```bash
SH_AGENT_RECALL=1 go test ./internal/search -run 'TestRecall$' -v
```

## What was measured

41 questions in `internal/search/testdata/recall.tsv`, drawn from this repository's history and worded in other terms than their answers use. 26 are ordinary questions, 10 are drift questions (every living file that restates a rule, answer sets hand-checked beyond the literal wording), 3 ask for references to an ID or a path, and 2 ask for open records.

The grep arm is one mechanical pass: each question word (stop words and words under three characters dropped) through `git grep -c -i -F`, files ranked by distinct words matched, then total hits. Ordinary and drift questions score recall over the first ten files; reference and open questions score over the complete list.

## Result, grep arm

| Class | Mean recall | Questions |
|---|---|---|
| Ordinary (recall@10) | 0.269 | 26 |
| Drift (recall@10) | 0.353 | 10 |
| References (whole list) | 1.000 | 3 |
| Open records (whole list) | 1.000 | 2 |
| All | 0.379 | 41 |

The scorer emitted both extremes on this run: 0.00 on 20 questions and 1.00 on 13, so it is not stuck at either value.

## What this does not measure

- An agent that greps, reads and greps again. One pass is a floor for grep, not what an agent reaches after several turns; the Phase 5 sandbox runs measure agents.
- Which words an agent would choose. The arm uses the question's own words.
- The open-record answers depend on the day: the plan-with-unticked-boxes answer is this plan, which closes when the work is done.
