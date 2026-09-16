---
owner: reviewer
date: 2026-09-16
---
# Review — bug-g6-security-exemption   (2026-09-16)

Range reviewed: working tree, five named paths

Blobs read:
- agents/security.md 94eef34ea211d1649222a62c2a4aebca405b91d8
- agents/reviewer.md 22206f3e2efa50d629067338205ac696da29d0aa
- agents/verifier.md deee22d0f6ebb92e22ea9448cb0874cb6869314f
- context/AGENTS.md 9c2778e13fbb6fb6843816cda492732160643f79
- skills/review-standards/SKILL.md 7a35efe98321f6747e201bfe14bbfda39a2de30f

## Verdict: APPROVE WITH NITS

The four "hand-rolled" definitions (reviewer, security, context, SKILL) and the two "stdlib alternative" dependency clauses (security, SKILL) restate the same rule consistently and add real, checkable content (a named list of primitives, a concrete `crypto/rand.Text()` example that checks out against the installed go1.24.4 toolchain). The one real defect is a scope narrowing in security.md's `Not examined` line that silently loosens the rule the other two files state unconditionally, and reviewer.md's list of stdlib-safe crypto uses has a dangling modifier that a careful reader can parse two ways.

## Findings

- **[medium / confidence: high]** `agents/security.md:71` — "`nothing` is not an answer on a review of a whole service: name the files outside the range you did not read and the tools that did not run." This is scoped only to "a review of a whole service," unlike the parallel lines in reviewer.md:89 ("`nothing` only when every file in the range was read and every criterion applied; otherwise name the gap") and verifier.md:87 ("`nothing` only when every tier the plan names ran and every package has a test; otherwise name the gap"), which apply the rule to every run regardless of scope. A careful reader of security.md alone can conclude that on an ordinary diff-scoped dispatch (not "a whole service"), writing `nothing` without having actually read every touched surface is fine — the opposite of what the other two agents now require and of the "silence is not [evidence]" framing in security.md's own `Write the artifact EARLY` section (line 21-22). Fix: drop the "on a review of a whole service" qualifier so the sentence reads as an unconditional rule, matching reviewer.md/verifier.md's "only when X; otherwise name the gap" pattern.

- **[low / confidence: medium]** `agents/reviewer.md:43` — "`crypto/rand` bytes or an identifier, a checksum, `crypto/tls` used as designed are stdlib, not DIY." The trailing clause "used as designed" grammatically attaches most naturally to the last item in the list (`crypto/tls`) only, leaving it ambiguous whether "a checksum" and "`crypto/rand` bytes or an identifier" are unconditionally exempt or must also be "used as designed." security.md:33 and SKILL.md's parallel sentences avoid this by putting "consuming... as designed" before the list, governing every item. Not blocking — the surrounding files make the intended reading clear — but two careful readers of reviewer.md in isolation could disagree on whether a misused (non-standard) checksum construction is covered by the exemption.

## Criteria

1. Dead code — pass; no removed or orphaned text left behind.
2. Flow tracing — the four hand-rolled definitions and two "stdlib alternative" clauses are consistent in substance across security.md, SKILL.md, reviewer.md and context/AGENTS.md; the `crypto/rand.Text()` (Go 1.24) claim was checked against the installed toolchain (`go doc crypto/rand.Text`) and is accurate. The one inconsistency is the security.md `Not examined` scope narrowing above.
3. Standards conformance — n/a (prose rules, not code); the changes read as a legitimate narrowing of an overbroad "hand-rolled" definition, consistent with the project's git history (BUG-033..038 line of fixes).
4. Directory layout — n/a, no new files.
5. Simplicity/DIY-except-security — the changes *are* the security-exemption clarification; no new hand-rolled code introduced, no unjustified dependency added.
6. Decoupling — n/a, prose only; the four restatements are independently editable but must stay in lockstep by hand (no single source of truth), which is exactly how the one inconsistency above arose.

## Not examined

Did not grep the full repo for every prior occurrence of "hand-roll" / "vetted librar" (checked agents/*.md, context/AGENTS.md, skills/review-standards/SKILL.md, skills/golang/references/*, REQUIREMENTS.md, context/dev-setup.md — all consistent or unrelated); did not check commands/preflight.md's gate logic beyond the one line that reads `## Not examined` sections (skills/preflight/SKILL.md:77), which is unaffected by this diff. Did not review .agents/docs/plans or other in-flight sandbox-simulation work referenced in recent commits — out of scope for the five named files.

## Outcome

- Finding 1 (medium) — fixed: `agents/security.md` § Not examined now reads "`nothing` only when every file in the blast radius was read and every tool ran; otherwise name the gap", unconditional like the reviewer's and verifier's lines.
- Finding 2 (low) — fixed: `agents/reviewer.md` criterion 5 now reads "consuming `crypto/rand` bytes or an identifier, a checksum, or `crypto/tls` as designed is stdlib, not DIY", so "as designed" governs the whole list.
- Round 1 closed 2026-09-16; no second round owed (two wording nits, no mechanism finding).
