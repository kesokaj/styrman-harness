---
owner: main
date: 2026-10-05
---
# Kit credential check fails closed (BUG-160 follow-up)

## Type

bugfix

## Summary

`testing/kit-preamble.sh` `kit_credential_leftovers`, after the code review (`reviews/2026-10-05-review-kit-credential-check.md`) and security review (`reviews/2026-10-05-security-v2.4.1-kit-credential-check.md`) reproduced seven passes with a credential present: a path `find` cannot enter or a file the shell cannot read is now reported; symlinks are followed, `KIT_HOME`'s own included (`find -H`); only the kit's top-level `bin/` is skipped; binary files are scanned (`grep -a`) except ELF and Mach-O executables, since the installed `sh-agent` carries the report redaction patterns; added names `.netrc`, `.git-credentials`, `*.p12`, `*.pfx`, `.aws/credentials`, `.docker/config.json`, and patterns `ya29.`, `AIza`, `AKIA`, `glpat-`, `sk-proj-`.

## Verification

`TestKitPreambleRefusesLeftoverCredentials`: 12 leftover cases fail the run, a lone `.pub` and an ELF binary naming the PEM header pass, a symlinked `KIT_HOME` holding a key fails. Against the previous check (ee28fe3) exactly the seven new gap cases fail (mutation run). No false alarm on today's real sandbox homes, `.agents/scratch/ct-new` (4,595 files) and `ct-old` (4,221), Antigravity conversation databases and the installed binary included. Test tooling only: no release.
