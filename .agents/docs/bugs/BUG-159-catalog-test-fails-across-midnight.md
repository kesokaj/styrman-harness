---
owner: bugs
date: 2026-10-02
id: BUG-159
status: fixed
area: internal/harness/catalog_test.go (TestRefreshCatalogScript)
---
# BUG-159: TestRefreshCatalogScript failed once with no code change when a run crossed midnight

## Symptom
Two independent agents, both started just before 2026-10-02 00:00, each saw one red run with no code change: the gate's verifier got exit 1 on its first run in a copy (four reruns passed; `.agents/docs/testing/2026-10-01-verify-v2.3.0-gate-r9.md`), and the code review of 42670ff..0a60619 saw one failure at the check "a second run at the same pins changed the catalog", also at 42670ff.

## Repro
Run the test so that it straddles 00:00 local time. Not reproduced on demand; the reviewer's timing found the cause.

## Root cause
The refresh script writes today's date (`date +%F`) into the manifest when a row's pin is new. The test compares the catalog across two runs of the script; a test that crosses midnight compares two dates. Found by the code review (low, older than the release). The verifier had flagged the red run "unexplained" and asked that it be noted before tagging.

## Fix
The test waits out the last two minutes of a day before it starts (about twenty seconds). The script is unchanged: a manifest dated on the day of the refresh is the intended behaviour.

## Prevention
A test that compares two runs of something that stamps the date either fixes the date or does not straddle midnight. A one-off red run gets its cause found before it is called a flake.

## Related
BUG-155 (same test).
