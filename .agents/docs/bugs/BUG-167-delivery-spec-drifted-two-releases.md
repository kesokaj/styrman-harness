---
owner: main
date: 2026-10-09
id: BUG-167
status: fixed
area: .agents/docs/specs/sh-agent/delivery.md (§ Where things stand, the delivery claim's one home)
---
# BUG-167: the delivery claim's one home named v2.4.1 and no `metrics` at v2.6.2

## Symptom
`.agents/docs/specs/sh-agent/delivery.md:77` said "Built and shipped (newest tag v2.4.1, 2026-10-05)" and listed the shipped verbs without `metrics`, while `git tag` ended at v2.6.2 and `metrics` had shipped in v2.6.0. Under AGENTS.md rule 6 this paragraph is the one place that says what a tag delivers; README, the deploy guides and the ADRs point at it. Found by the project review of 2026-10-09 (code+docs lens, medium); it had reached four releases.

## Repro
- `git tag | sort -V | tail -1` against `grep -n 'newest tag' .agents/docs/specs/sh-agent/delivery.md`.
- `TestDeliverySpecCurrent` (`cmd/sh-agent/delivery_test.go`) was red before the fix on both counts: the verb `metrics` missing, "newest tag v2.6.2" missing.

## Root cause
The release procedure (`deploy/remote.md` § Releases, run by `/preflight` Step 6) writes the release entry under `changes/` and the release notes, and never names this paragraph; nothing in the tree compared the paragraph with what the tree itself records (the newest release entry, the verb table). A claim with one home and no pin drifts at the first release that forgets it.

## Fix
- The paragraph: v2.6.2, `metrics (specs/sh-agent/metrics.md)` in the verb list, `mcpServers` in the refused keys, and one sentence each for v2.6.2, v2.6.1, v2.6.0 and v2.5.0 in the paragraph's own style.
- The pin: `TestDeliverySpecCurrent` reads the paragraph, requires every verb in `cmd/sh-agent`'s verb table to be named, and requires "newest tag vX.Y.Z" to match the highest `changes/*-release-vX.Y.Z.md` entry.
- The procedure: `deploy/remote.md` § Releases runs that test right after the release commit, before the tag, so a release commit that writes its entry without the paragraph cannot be tagged.

## Prevention
A doc that restates a fact the tree holds elsewhere (a version, a verb list, a count) is pinned to that fact by a test that reads both, or carries no number. The release commit updates `delivery.md` § Where things stand together with the release entry; `TestDeliverySpecCurrent` in the release procedure refuses the tag otherwise.
