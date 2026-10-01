---
owner: bugs
date: 2026-10-01
id: BUG-154
status: fixed
area: agents/verifier.md (log template), testing (real-home smoke kit, sandbox kit)
---
# BUG-154: the verifier followed the dispatch's description of its log over its own template (file name, status line)

## Symptom
At 4165f81 the real-home smoke row `realhome claude agent` printed `pass=0` and the sandbox row `ad claude verifier` printed `pass=0 artifact=none`. In both the verifier HAD run and written a PASS log. Real-home: `.agents/docs/testing/2026-10-01-2042-verify-greet-comment.md` held a plain line `Status: PASS`, not the template's `## Status: PASS`; the kit greps `^## Status`. Sandbox: the log was `2026-10-01-2026-greet-comment.md`, no `-verify-` in the name; the kit globbed `*verify*.md`.

## Repro
Fresh clone of the demo service, Claude Code print mode, prompt "Dispatch the verifier agent over the working tree (one changed comment in internal/greet/greet.go) with a budget of 15 tool calls; tell it to write the log after the first command. ...". 2 of 2 Claude Code runs at 4165f81 deviated (one in the name, one in the status line); 2 of 2 Antigravity runs wrote `2026-10-01-verify-greet.md` with `## Status: PASS`.

## Root cause
Read from the real-home run's session transcript: the dispatching session described the log itself in the dispatch prompt ("Take the date/time stamp for the file name ... from the shell (`date +%F-%H%M`...)", "ends in an explicit `Status: PASS` or `Status: FAIL` line"). The verifier followed the dispatch over its own template, because `agents/verifier.md` gave the template but did not say its name and status heading hold against a dispatch that describes them differently. The release gate (`skills/preflight/SKILL.md` Step 2) and the kits read exactly those two things. The sandbox run's dispatch was not read: same shape of deviation, cause assumed the same, not proven.
`Diagnosis missed:` my first reading was "the kit's glob is wrong" (true for the sandbox row only); I re-judged that row a pass from the saved clone. The second row showed the cause was upstream of the kit.

## Fix
`agents/verifier.md` § Verification log: "Two things in it are read by programs and by the release gate, so they follow this template even when the dispatch describes them differently: the file name holds `-verify-`, and the status is the heading line `## Status: <one word>`, never a plain `Status:` line". Sandbox kit (`.agents/scratch/audit-fix/kit-agent-data.sh`, scratch, untracked) now takes any new `testing/*.md` except `guide.md` as the artifact and judges the planted-instruction row on content; the real-home kit keeps its strict `^## Status` check, which found this. Row `verifier-log-shape` (a dispatch that describes the log differently on purpose) is added to `testing/skill-matrix.md` by the main agent; its results at the final head are in the gate's logs, not known when this was written.

## Prevention
An agent's artifact has parts a program reads (file name, status line): the definition names them and says they hold against the dispatch, and a dispatch never re-describes an agent's artifact format: it names the slug and the duties, the definition owns the shape. A kit's `pass=0` is read from the saved run before it is called a failure or a pass.

## Related
BUG-083 and BUG-064 (same agent, same family: a definition's rule loses to dispatch text on Claude Code; those say put the duty in the dispatch, this one is the converse: shape duties live in the definition and the dispatch must not contradict them). BUG-084 is a retro-vocabulary record, not this class. Earlier verifier records: BUG-095, BUG-062, BUG-063.
