---
owner: bugs
date: 2026-10-05
id: BUG-160
status: fixed
area: testing (behavioural run kits, sandbox homes)
---
# BUG-160: live credentials left in scratch homes while the records said they were deleted

## Symptom
Security review `reviews/2026-10-05-security-v2.4.0-signer-rotation.md` (medium, confidence high): `.agents/scratch/{ct-ctl,ct-old,ct-new,dr-home}/` held `.ssh/id_ed25519` (unencrypted copy of the box's GitHub login key, account-wide push rights), `.config/gcloud/application_default_credentials.json` and `credentials.db`, and in dr-home and ct-new `.gemini/antigravity-cli/antigravity-oauth-token`. Records said otherwise: `testing/2026-10-05-verify-critical-thinking.md:19`, `-r2.md:14`, `-verify-dreamer-global-rules.md:14`, `-verify-v2.3.3-gate-rows.md:13`; handoffs 2026-10-05-0935 ("credentials deleted"), -1039, -1129 ("no credentials left"). Runs in those homes used `--dangerously-skip-permissions`; agents there tried `git push` on their own (remote was local). Files are gitignored, never committed.

## Repro
Not reproduced as a command; found by listing the homes. Before the fix, `find .agents/scratch -name id_ed25519 -o -name credentials.db` listed them.

## Root cause
The deletion is a per-run step in each throwaway kit script, "deleted by name" as `guide.md:158-162` says, and the list of names was the two kinds the kit author remembered (Claude `.credentials.json`, the Antigravity token). Credential kinds the sandbox gained otherwise (SSH key, gcloud files) were never on the list. Nothing proved the result: `kit-preamble.sh` has no credential handling at all (it fingerprints the real `~/.config/sh-agent` only), so the claim "deleted" was written from the kit's intent, not from a check. A cleanup that deletes by an enumerated list cannot be right about kinds it does not know it copied.

## Fix
Main agent deleted those files by name on 2026-10-05; a `find` for the credential names under `.agents/scratch` then returned 0. Mechanism landed in commit ee28fe3: `.agents/docs/testing/kit-preamble.sh` gained `kit_credential_leftovers`, run from `kit_on_exit`. It fails the run (exit 1, "kit: FAIL credentials left in KIT_HOME (BUG-160)", paths only) when KIT_HOME holds a credential by file name (id_rsa, id_ecdsa, id_ed25519, id_dsa, .credentials.json, application_default_credentials.json, credentials.db, access_tokens.db, `*oauth-token*`, hosts.yml) or a token or private key inside any file (`gh[opsur]_` tokens, `github_pat_`, a PEM private key header, `sk-ant-`, "refresh_token"), skipping the kit's own bin/. A run counts only with "kit: no credentials left in KIT_HOME". `testing/guide.md` § 4 names the check. Environment dumps a CLI writes (the addendum's `term1.env`) are covered by the content scan.
Also done outside the code: GitHub Actions was turned off, and the operator was told the gh token can only be revoked by them, in GitHub's Authorized OAuth Apps.

## Prevention
Before a kit copies any credential into a sandbox home, it records every path it copies and deletes exactly those, then proves none is left with a `find` over the home for the credential names (`id_ed25519*`, `.credentials.json`, `application_default_credentials.json`, `credentials.db`, `*oauth-token*`); that `find` is run once against a home that still holds a planted credential (must-fail seen) before its empty result is cited. A record or handoff never says "credentials deleted" without that find's output.
Mechanism now enforces this: the preamble's exit check (`kit_credential_leftovers`) plus `TestKitPreambleRefusesLeftoverCredentials` in `internal/harness/catalog_test.go` (five leftover kinds fail; a clean home and a lone `.pub` pass; with the check disabled all five cases failed, the mutation run). A kit sourcing the preamble needs no list of its own; cite the line "kit: no credentials left in KIT_HOME" as the proof.

## Related
Existing rules not followed: `guide.md:158-162` (deleted by name, security review 2026-09-21); global "silence from your own instrumentation is not evidence" (context Hard rules). Recommended home when the kit fix lands: `kit-preamble.sh` exit check plus a behaviour test beside `TestKitPreambleAssertsCommit`; actionable (a specific find catches it) and generalizable (any sandbox copy of a secret). Same finding: the retired signing key left in a transcript (finding 2, not this record). Status moved to fixed 2026-10-05 with commit ee28fe3; INDEX row graduated to the preamble check and its test.

## Addendum (main, 2026-10-05)

Security review round 2 (`reviews/2026-10-05-security-v2.4.0-signer-rotation-r2.md`, finding 3) found a further kind: the live `gh` token, exported as `GH_TOKEN`, written by Antigravity into `.agents/scratch/gate-agy/.gemini/antigravity-cli/brain/<id>/.system_generated/terminals/term1.env`. Deleted by name; a scan of `.agents/scratch` for `gh[pousr]_`, `github_pat_` and `BEGIN OPENSSH PRIVATE KEY` then found only repository copies naming the patterns. The kit fix must cover environment dumps a CLI writes, not only files the kit copies; the operator was advised to rotate the token.
