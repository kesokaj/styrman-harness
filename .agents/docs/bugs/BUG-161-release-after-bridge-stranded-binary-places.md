---
owner: bugs
date: 2026-10-06
id: BUG-161
status: fixed
area: process / release (skills/preflight, deploy/remote.md)
---
# BUG-161: v2.4.1 was cut right after the v2.4.0 bridge, stranding a binary-only place

## Symptom
Operator, 2026-10-06, on the machine einherjar (a v2.3.2-or-older release binary): `sh-agent update` printed "downloading sh-agent-linux-amd64 v2.4.1", then "sh-agent: update: release v2.4.1 has no SHA256SUMS.sig". Every retry gave the same.

## Repro
From a binary up to v2.3.2 with v2.4.1 the newest release: `sh-agent update`. Those binaries require a `SHA256SUMS.sig` by the old key; v2.4.1 carries none (ADR-053: no later release carries one).

## Root cause
v2.4.0 (tagged 2026-10-05 16:35 UTC) was the bridge release: the only one carrying a `SHA256SUMS.sig` by the old key (ADR-052/053; delivery.md § The model). `sh-agent update` moves only to the NEWEST release. The main agent cut v2.4.1 (17:33 UTC, 58 minutes later, at the operator's request) without first saying it would close the bridge for every place not yet updated, and without asking which places had crossed. einherjar had not. The v2.4.0 handoff and release notes said "update while v2.4.0 is the newest"; the next release then made it not the newest. The gate (preflight) had no step that treats "the previous release was a bridge" as a condition of the next one.

## Fix
Operator replaced einherjar's binary by hand with the README § Install without cloning chain for v2.4.1 (checksum verified). `sh-agent version` v2.4.1, `status` drift none, `update` reports "up to date (v2.4.1)". The Mac had crossed to v2.4.0 and updates normally. Zip places are not stranded: a v2.3.2 binary installs the v2.4.1 zip (exit 0), marked UNVERIFIED, not refused (scratch run 2026-10-06).

## Prevention
Rule: after a release that bridges a trust or format change (one older binaries need in order to move forward), the next release is not cut until the operator confirms every binary-only place has crossed.
Graduated 2026-10-06 into a mechanism: `skills/preflight/SKILL.md` Step 6, the paragraph "Then the bridge check." (after "Before item 1"): /preflight names the bridge and asks, in the operator's own turn, whether every place has crossed it; without a yes it stops before item 1, and an unattended run never tags past a bridge. Project `AGENTS.md` § Release carries the one-clause pointer ("a release right after a bridge release waits for the operator's word that every place has crossed it (preflight Step 6, BUG-161)"). Both read and confirmed. The changed skill ships in the next release, which owes its Antigravity run (AGENTS.md § Release); until then the rule is prose, not yet exercised on either CLI.

## Related
ADR-052, ADR-052a, ADR-053 (`architecture/DECISIONS.md`); `specs/sh-agent/delivery.md` § The model; `reviews/2026-10-05-preflight-v2.4.0-r2.md`, `reviews/2026-10-05-preflight-v2.4.1.md`. Not a repeat of an earlier rule; first record of this class. Rule is actionable (a preflight step can grep the prior release's notes for "bridge") and generalizable (any trust or format change). It costs no context bytes if placed in the preflight skill.
