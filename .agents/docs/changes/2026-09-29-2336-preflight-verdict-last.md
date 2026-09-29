---
owner: main
date: 2026-09-29
---
# Preflight writes its verdict last; BUG-105 and BUG-083 closed   (2026-09-29-2336)

Asked for by the operator: nothing waiting before the release. `sh-agent search --open` listed two open bug records.

- **`skills/preflight/SKILL.md` Step 5** gained BUG-105's prevention rule: the verdict line is written last, after every item under "What was not checked" is read against every blocking condition, the project's own included.
- **Behavioural runs** (`testing/2026-09-29-preflight-verdict-last-behavioural.md`): with the rule the verdict is written last in 6 of 6 reports, without it in 0 of 6; NO-SHIP in all twelve; one Antigravity run completed. Four earlier runs were void: the kit's project-level agent copies tripped preflight's shadow-agent check.
- **BUG-105** set to fixed and its index row moved to Graduated. **BUG-083** set to fixed: 31 of 31 verifier dispatches measured today wrote their log within their first four commands.

Agents dispatched: `bugs`. The diff is two lines of skill text; read by the main agent, no reviewer (under the threshold). The release gate's docs lens covers the range.

Verification: `go run ./cmd/sh-agent install --home .agents/scratch/home-pf` exit 0, `status` drift none; both vendors' installed skill carries the sentence.
