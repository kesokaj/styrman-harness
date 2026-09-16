---
owner: bugs
date: 2026-09-16
id: BUG-041
status: fixed
area: agents/security, agents/reviewer, skills/review-standards, context (security exemption)
---
# BUG-041: security agent graded stdlib crypto/rand UUID generation as hand-rolled cryptography

## Symptom
Antigravity's `/feature` run had the harness `security` agent grade a UUID v4 built from `crypto/rand`
bytes as "hand-rolled cryptography, strictly prohibited" (High), while the harness `reviewer`, in the
same round, had judged the same lines as stdlib use that "does not violate the security exemption".
The main agent followed security, added `github.com/google/uuid` against the golang skill's
stdlib-first rule, and the `bugs` agent then wrote a project prevention rule ("always import vetted
UUID/crypto libraries") that itself contradicted the harness rule. Both agents' reports also closed
with `Not examined: nothing`.

## Repro
Sandbox `sim-agy-linkshelf` reviews `2026-09-16-linkshelf-security.md` and `…-review.md`,
`bugs/BUG-001-hand-rolled-uuid.md`. Reproduction row `security-rand-id` (kit `g6-control.sh`) at
harness `cad463d`: Claude Code's security agent filed the same code low ("because it is an
identifier"), Antigravity's filed it as "hand-rolled instead of using a vetted library" — the two
harness agents disagreed on one mechanism even in the control run.

## Root cause
The four restatements of the "never hand-roll crypto" rule (context § Hard rules, review-standards §
5, security § 4, reviewer criterion 5) never defined what "hand-rolled" means, so consuming
`crypto/rand` bytes to build an identifier read, to one agent, as rolling crypto rather than using it
as designed. The report templates' "Not examined: nothing" allowed a blank full-coverage claim with no
condition attached, so neither agent's false confidence was visible from its own artifact.

## Fix
All four restatements now define hand-rolled as a primitive or protocol implemented or composed in the
target code and named in the finding; stdlib `crypto/*` used as designed (bytes or an identifier from
`crypto/rand`, a checksum, TLS) is explicitly no finding. A finding that asks for a dependency must
name the stdlib alternative it rejected (e.g. Go 1.24 `crypto/rand.Text()`). The three `Not examined`
lines now allow `nothing` only under a stated full-coverage condition.

## Prevention
A ban must name the thing banned precisely enough that its stdlib counterpart is not caught by it.
When two harness agents contradict each other on the same mechanism in the same round, the rule is the
defect, not the code — fix the rule before trusting either report.

## Related
Review: `.agents/docs/reviews/2026-09-16-review-bug-g6-security-exemption.md`. Testing:
`.agents/docs/testing/2026-09-16-verify-linkshelf-sim-fixes.md`. Sandbox re-run owed (not yet run): G6
reproduction row `a2` (after the fix) on both CLIs, per
`.agents/docs/testing/2026-09-16-linkshelf-simulation.md` § G6 reproduction row / § Re-runs.
