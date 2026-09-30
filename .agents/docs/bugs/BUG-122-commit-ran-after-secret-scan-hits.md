---
owner: bugs
date: 2026-09-30
id: BUG-122
status: fixed
area: process (commit procedure)
---
# BUG-122: a commit went through after a secret scan that reported hits

## Symptom
For e2cffc1 the staged-diff secret scan printed 3 hits and the commit ran anyway.

## Repro
Scan and commit were chained with `;`, so the scan's output and status did not gate the commit.

## Root cause
Chaining the verification and the commit in one call: the commit ran regardless of the scan. No secret was committed (all three hits were prose naming the scan's own patterns, in a deploy-guide line and review text, read after the fact), so this was luck, not control.

## Fix
No code change; the procedure was corrected in this session (scan in its own call, read, then commit in a separate call).

## Prevention
The global rule already covers this: a verification step and the commit are separate calls, and the verification's result is read first. This was a failure to follow it.

## Related
Recurrence of an existing rule: a graduation candidate for the retro (actionable: a scan that prints hits blocks the commit; generalizable: applies to every commit). Recommend the retro decide its permanent home; no record count of three is reached by this alone.
