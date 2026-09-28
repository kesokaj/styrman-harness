---
owner: main
date: 2026-09-28
---
# CLI token files stop a commit; skill-matrix sandbox permissions narrowed     (date: 2026-09-28)

## Summary
These are the two v1.7.0 gate deferrals, approved by the operator ("2.) approved. 3.) approved."):
- **Context rule.** The global context's secret-shaped-path rule now names "a CLI token file such as `antigravity-oauth-token`", matching preflight's list since v1.7.0.
- **Skill matrix.** The full-matrix procedure's sandbox permissions follow `AGENTS.md` § Testing:
  - no `bash`, `python3` or unrestricted `curl` for Claude;
  - named commands instead of `command(*)` for agy;
  - the invalid `edit_file(*)` and `list_dir(*)` entries removed;
  - clone-only trust and no `toolPermission` key;
  - the auto-denied and invalid-allow-entry checks.

Also recorded: the operator updated the work Mac to v1.7.0 ("verified and installed sh-agent-darwin-arm64 v1.7.0"; 6 Antigravity files written, 123 unchanged, 129 managed). Its canary check is still open.

## Type
task

## Files changed
- `context/AGENTS.md` — § Git, the secret-shaped-path list (+52 B; context now 20,550 of 22,000 B).
- `.agents/docs/testing/skill-matrix.md` — Procedure step 2, both permission lists.
- `.agents/docs/testing/2026-09-28-verify-token-rule-and-matrix.md` — the behavioural and settings evidence.
- `.agents/docs/testing/2026-09-28-verify-v1.8.0-candidate.md` — verifier, PASS.

## Agents dispatched
- `verifier`: PASS.
- The gate's reviewer, security and docs lens follow in `/preflight v1.8.0`.

## Verification
- **Context clause:** new 3 of 3 stopped the commit and asked; the control 3 of 3 committed the rest and left the token out. The token was committed 0 of 6 times.
- **agy settings:** the first check caught `edit_file(*)` as invalid, and `list_dir(*)`, the known-invalid control, fired as expected. The committed list passed with 0 invalid entries and no auto-deny.
- **Verifier:** 132 top-level and 251 total tests pass, govulncheck clean, drift none, the clause is installed on both vendors, and the context ends with GANNET.

## Docs
`testing/skill-matrix.md` is the procedure document.

## Requirements
n/a
