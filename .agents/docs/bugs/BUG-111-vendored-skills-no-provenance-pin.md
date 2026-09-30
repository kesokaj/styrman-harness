---
owner: bugs
date: 2026-09-30
id: BUG-111
status: fixed
area: skills/VENDORED.md
---
# BUG-111: two vendored skills carried no source, pin or digest, and rewritten local text was indistinguishable from upstream

## Symptom
`skills/VENDORED.md` listed `gke-basics` and `terraform` with "not recorded
when vendored" as their source and no pin, alongside the manifest's own
claim that vendored skills are copied verbatim; both had in fact been
rewritten locally. Found by the project security review, finding 5.

## Repro
Read `skills/VENDORED.md` before the fix: no repository, no commit, no
digest for either row, while their `SKILL.md` text differs from any upstream
copy a reader could compare against — there was no way to tell a local edit
from upstream drift.

## Root cause
The two skills were vendored before the manifest's provenance convention
existed and were never backfilled; because the manifest claims "verbatim",
nobody treated the local rewrites as forks needing their own pin and
diff record. Even the fix's first attempt (commit eff85cf) wrote a pin for
`gke-basics` that was itself wrong, discovered only when the round-2
security check diffed every candidate upstream commit with `git archive |
diff -r` instead of trusting the commit date.

## Fix
`skills/VENDORED.md`: both rows moved to a new "forks" table recording the
upstream repository, the verified base commit (`google/skills` 6f0b877;
`antonbabenko/terraform-skill` v1.6.0 2271bc4), and the local changes made
since. `AGENTS.md` names them as forks. `TestVendoredRowsCarryAPin` refuses
any row lacking a repository and a pin.

## Fix (round 2 correction)
The initial `gke-basics` pin was wrong; corrected to 6f0b877 (2026-04-23)
after the security round diffed the upstream history with `git archive
<commit> <path> | diff -r`.

## Prevention
A provenance pin is verified by `git archive <pin> <path> | tar -x` and
`diff -r` against the vendored directory before it is written, not inferred
from a commit date or message; the diff itself is recorded as the row's
"Local changes". "Vendored" means verbatim; anything with local edits is a
fork, and it says so in the manifest, with its own pin and diff record.

## Related
BUG-096 (cross-file citation not verified against the target); the
provenance-verification rule (`git archive | diff -r` before writing a pin)
is a good candidate for `skills/VENDORED.md`'s own header notes, which
already state the update convention.
