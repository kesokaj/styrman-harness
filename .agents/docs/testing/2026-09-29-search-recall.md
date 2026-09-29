---
owner: main
date: 2026-09-29
---
# Recall — the search index against grep, and the vector gate   (2026-09-29)

Plan: `plans/2026-09-29-search-index.md`, Phase 4. Baseline: `testing/2026-09-29-recall-baseline.md`.

## Status: RECORDED — the index beats one grep pass; the vector gate is not passed

## Commands

```bash
SH_AGENT_RECALL=1 go test ./internal/search -run 'TestRecall$' -v -count=1
.agents/scratch/venv/bin/python .agents/scratch/vecprobe/probe.py   # the vector probe (scratch only, never committed)
```

## The index arm against the grep arm (41 questions)

Both arms use the same question words. The index arm runs `sh-agent search` (words), `--refs` or `--open`; the grep arm is one `git grep -c -i -F` pass per word. A hit in `internal/search/testdata/` (the questions themselves) is dropped from both arms before scoring.

| Class | grep | index |
|---|---|---|
| Ordinary, recall@10 (26) | 0.269 | 0.692 |
| Drift, recall@10 (10) | 0.353 | 0.376 |
| References, whole list (3) | 1.000 | 1.000 |
| Open records, whole list (2) | 1.000 | 1.000 |
| All (41) | 0.379 | 0.653 |

The first index run scored 0.077 on ordinary questions: `search` ANDed every word, matched two chunks (one of them the question file itself) and stopped, since the OR retry ran only on zero hits. `search` now ranks chunks with every word first and fills the remaining slots with chunks holding any word (plan D27), and the recall test drops the question file from both arms.

Where the index still scores 0 on ordinary questions, the question's words and the answer's words do not overlap (for example "antigravity cut off the end of the global rules file" against BUG-085's "truncated at 24,000 bytes"). Drift questions stay the weak class for both arms: an answer set of up to ten files, each wording the rule differently.

## The vector gate (36 ordinary and drift questions)

Every chunk of this repository's index (10,403) embedded with all-MiniLM-L6-v2 (384 dimensions, `fastembed` 0.8.1 on ONNX Runtime, 236 s on this box's CPU), ranked by cosine similarity, fused by reciprocal rank fusion with k = 60.

| Arm | Ordinary (26) | Drift (10) | All (36) |
|---|---|---|---|
| grep | 0.269 | 0.353 | 0.293 |
| index | 0.692 | 0.386 | 0.607 |
| vec | 0.692 | 0.153 | 0.543 |
| index+grep (RRF) | 0.577 | 0.383 | 0.523 |
| index+grep+vec (RRF) | 0.769 | 0.416 | 0.671 |
| index+vec (RRF) | 0.769 | 0.415 | 0.671 |

**Gate: not passed.** Against the best non-vector arm, the index alone, fusing in vectors adds 6.4 points overall and 3.0 on drift, under the plan's 10-point bar. Read literally, the plan compares against "the index arm plus the grep arm"; their fusion (0.523) is weaker than the index alone, and against it vectors add 14.8 points. A gate that a weaker baseline can pass measures the baseline, so the index alone is the comparison used (plan D30).

**Delivery condition not met either.** The adoption rule asks for an in-process model with no cgo, no daemon and no network. `github.com/clems4ever/all-minilm-l6-v2-go`, named in the research as pure Go, loads ONNX Runtime from a shared library at run time (`WithRuntimePath("libonnxruntime.so")`, `github.com/yalue/onnxruntime_go`); no pure-Go route was measured.

## What was not measured

- An agent that searches, reads and searches again; the Phase 5 sandbox runs measure agents.
- Other embedding models (code-specific ones such as jina-embeddings-v2-base-code), chunk sizes, or hybrid weights other than equal RRF.
- Vector storage and query speed in `modernc.org/sqlite/vec`: the probe kept vectors in memory, since the gate is about retrieval quality.
