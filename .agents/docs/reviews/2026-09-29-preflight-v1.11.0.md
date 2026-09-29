---
owner: main
date: 2026-09-29
---
# Preflight — v1.11.0   (2026-09-29)

Range reviewed: v1.10.0..814c257

VERDICT: SHIP

Nothing blocks the release:
- no verifier FAIL, BLOCKED or IN PROGRESS log at the head;
- no security row: the range changes rule text only (`context/`, `agents/`, `skills/`), no auth, crypto, secrets, IAM, input parsing or dependencies, and no Go file;
- the final code review is APPROVE WITH NITS, both nits closed in `814c257`;
- no REQUIREMENTS change;
- no project agent shadows a review agent.

**How the range got here.** Eight retro proposals were applied at once (`d3cc2b7`). Four review rounds followed:
- r1 (`reviews/2026-09-28-review-v1.10.0-d3cc2b7-code.md`): REQUEST CHANGES, three mediums in the new text; fixed in `cb6cf73`.
- r2 (`reviews/2026-09-28-review-v1.11.0-gate-r2.md`): REQUEST CHANGES, the same two mechanisms again; the operator struck both in `1d03e7a` (BUG-095).
- r3 (`reviews/2026-09-29-review-v1.10.0-1d03e7a-code.md`): REQUEST CHANGES, the long-suite rule a third time; the operator struck it in `740ca80`.
- r4 (`reviews/2026-09-29-review-v1.10.0-740ca80-code.md`): APPROVE WITH NITS; closed in `814c257`.

## What ships (the deployed diff, 6 files)

| File | Change |
|---|---|
| `context/AGENTS.md` | DoD step 2: a low finding whose fix adds a flag, a refusal or a step for the operator is put to the operator, recorded `left: operator`. § Git: "is it live?" is answered with the fact and an offer to run `/preflight`, never with the steps `/preflight` performs handed to the operator. § Work: a finished agent whose task still runs background work is stopped. |
| `agents/verifier.md` | `maxTurns: 30` → `45`. Nothing else. |
| `skills/feature/SKILL.md` | Inputs that always travel together are one input. |
| `skills/research/SKILL.md` | A finding about another project stays out of this one (BUG-088 graduated). |
| `skills/review/SKILL.md` | The `left: operator` outcome in the `--fix` round. |
| `skills/preflight/SKILL.md` | A `left: operator` finding is accepted into the release entry; the run does not stop. |

**Struck before release** (retro proposals 3, 5, 7): the verifier shared-call rule, the cold-read cap early exit, and the long-suite rewording. Proposal 7 is re-queued in `~/.config/sh-agent/proposals/2026-09-23-retro.md` with the note that it needs a behavioural run designed for it.

## Findings and decisions, worst first

| # | Source | Severity / confidence | Finding | Decision |
|---|---|---|---|---|
| 1 | code r4 | medium / high | Context sentences 1 and 8 had no behavioural run | Sentence 1: run in `814c257` (`testing/2026-09-29-verify-retro-probe2.md`, new 2 of 2 put the choice to the operator, control 0 of 2). Sentence 8: **accepted gap**, it needs a multi-agent session with a background child; named here and in the release entry. |
| 2 | code r4 | low / medium | The "is it live" sentence could misfire on steps only the operator can take | Fixed in `814c257`: scoped to "the steps `/preflight` itself performs". Probe: 2 of 2 still offer `/preflight`. |
| 3 | code r1–r3 | medium (×7) | Contradictions and an unverifiable check in the applied rule text; a verify log counting attempts as executions | All fixed or struck; each round's `## Outcome` says which. BUG-095 records the cause. |

## Reviews

- **Code:** the four rounds above; every round's report carries `## Outcome`. The `814c257` change (one clause) was read by the main agent, not re-reviewed, as the Definition of Done allows below 50 lines.
- **Security:** no row (see the verdict line).
- **Docs:** `reviews/2026-09-28-review-v1.10.0-d3cc2b7-docs.md` (APPROVE, no findings). The later commits only removed or narrowed text from that range, and the docs lens's subjects (README, specs, deploy guides) describe none of these rules.

## Verification

`testing/2026-09-29-verify-v1.11.0-gate-r5.md`, PASS at `814c257`:
- gofmt and vet are clean; `go test -race`: 171 pass, 0 fail, 0 skip; govulncheck is clean;
- context 21,444 of 22,000 bytes, GANNET at the end on both vendors;
- scratch install: drift none; 27 skills, 7 agents.

Behavioural evidence: `testing/2026-09-28-verify-retro-apply.md` (sentence 6 vs control, 2 of 2 vs 0 of 2; the verifier probe, whose shared-call subject was later struck) and `testing/2026-09-29-verify-retro-probe2.md` (sentence 1 vs control; sentence 6 as scoped).

## Trail freshness

- **Newest change entry:** `changes/2026-09-29-0706-islive-scope-probe2.md`; with `0659-retro-round3`, `0623-retro-strikes`, `2026-09-28-2228-retro-apply-fixes` and `2222-retro-apply` it covers the range.
- **Dates:** none is in the future.
- **Append-only:** `changes/` gained additions only; `DECISIONS.md` is untouched.
- **ROADMAP.md** "Now" and **TODO.md** are empty. The retro record and the three queue files carry final statuses (5 applied, 2 rejected, 1 queued).

## Version

MINOR, v1.11.0: installed places gain context rules and a verifier budget; nothing breaks, nothing becomes required.

## What was not checked

- Context sentence 8 behaviourally (above).
- The `left: operator` restatements in the review and preflight skills behaviourally (skill prose; install and parse check only).
- Antigravity for every probe; the skill matrix; the real-home smoke (operator, after `sh-agent update`).
