---
owner: main
date: 2026-09-18
---
# Preflight — v1.4.0     (2026-09-18)

Range reviewed: v1.3.0..42fad8f

Base: the last tag, `v1.3.0` (release commit fc90cbb). Head: 42fad8f, `main`, clean tree. The range is the second retro's application (42fad8f: three harness proposals, two project proposals, the delivery claim consolidated into `delivery.md` § The model; BUG-083/084 graduated) plus four record commits (the v1.3.0 handoff, three matrix rows). No Go source changed: the binary is v1.2.3's code with the new harness embedded. Version: MINOR — an installed place gains a skill step, a context rule and three agent behaviours; the operator's "new versions minor".

## VERDICT: SHIP

Nothing blocks. The operator's instruction for this release: apply all, cut a minor, no further retro, then stop; the review cap decisions below were the main agent's under that instruction and are recorded.

## Findings, worst first

None high. None medium open.

- [medium → fixed before the commit] The application's own review, round 1 (code): the three agents' templates had no Status field while the new rule made `IN PROGRESS` load-bearing — the BUG-084 shape, made in the diff that wrote BUG-084's rule; fixed (`## Status: IN PROGRESS | COMPLETE` with the never-at-hand-back comment in `reviewer`, `security`, `explorer`). Round 1 (docs): README:59/:61 still restated the delivery claim because the sweep grep matched only the singular; fixed, grep widened. `reviews/2026-09-18-review-worktree-retro-2-{code,docs}.md` § Outcome.
- [medium → rebutted] Round 2 (code) open items on rule 6's scope (README:154/:163/:182, `deploy/local.md:40-43`, `deploy/remote.md:59` describe what `update` does at a place): rule 6 now states that a guide's description of a COMMAND is that command's documentation, not the claim; `delivery.md:37` trimmed and § The model completed (the dirty-checkout refusal, the canonical-and-newer tag, `SHA256SUMS.sig`, the install path). Decided by the main agent, no round 3, per the operator's instruction to finish; recorded in the code review's § Outcome.
- [low / on record] Three prose edits after review round 2 (rule 6's scope sentence, `delivery.md` § The model and `:37`) were read by the main agent, not by a reviewer; the verifier's round 3 covers the committed head. Stated here, not hidden.
- [low / carried] `skills/preflight/SKILL.md:77`'s IN PROGRESS clause is still unqualified (accepted at v1.3.0; the handoff's first Task); BUG-083's tally stands at 1 of 3 and 1 of 1.

## Rows

| Row | Result | Artifact |
|---|---|---|
| Verification | PASS ×3 — round 1 on the working tree, round 2 on the frozen tree after the review fixes, round 3 at 42fad8f: `-race` suite ok, built-binary install (fresh home 260 written) / status drift none / list 27 skills 7 agents / dry-run 0 written; `Range reviewed: v1.3.0..42fad8f`; `govulncheck` none on go1.27.1 (main thread, bare) | `testing/2026-09-18-verify-retro-2-applied.md` |
| Review (opus, code lens) | round 1 REQUEST CHANGES (3 + 3 lows) → round 2 APPROVE WITH NITS (rule 6 scope, rebutted above); `## Outcome` per finding. Covers 42fad8f less the three post-round prose edits (read by the main agent) | `reviews/2026-09-18-review-worktree-retro-2-code.md` |
| Docs lens | round 1 REQUEST CHANGES (README:59/:61) → round 2 APPROVE, no narrowing of any removed qualifier | `reviews/2026-09-18-review-worktree-retro-2-docs.md` |
| Security | no row: no changed path touches auth, crypto, secrets, IAM, input parsing or dependencies; the staged-diff secret scan's two hits are the `-----BEGIN SSH SIGNATURE-----` armor header the project `AGENTS.md` names as a tag check | — |
| Trail freshness | pass: newest `changes/` entry `2026-09-18-retro-2-applied.md` is this range's work; no future `date:`; no modified past `changes/`/`DECISIONS.md` entry since v1.3.0; `ROADMAP.md`/`TODO.md` untouched (both empty — lag, not a violation) | — |
| Shadowing project agent | none | — |
| `REQUIREMENTS.md` | untouched since v1.3.0; § The model states R8/R9 | — |
| Signing key | present, fingerprint `SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU` from the private key | — |
| Embed guard | nothing untracked or ignored under the five embedded directories; `TestEmbedHoldsOnlyTrackedFiles` ok | — |
| Repository visibility | `PRIVATE` | — |

## What was not checked

- No CLI start with 42fad8f's files: the real home on this box is at b06268d (the timer will move it); no sandbox start today.
- No behavioural run of the three agents' new early-write sentence (BUG-083: a definition alone binds in some runs; the dispatch sentence carries it).
- No Mac or arm64 machine ran anything in this range.
- The three post-round-2 prose edits were not reviewed by an agent (above).
- No Dockerfile; no container tier.
