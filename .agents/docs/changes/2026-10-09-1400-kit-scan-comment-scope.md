---
owner: main
date: 2026-10-09
---
# v2.6.1 gate: the kit scan's comment states its real scope

From the security review's round 3 (`reviews/2026-10-09-security-v2.6.0-f6d9854.md`, one low and one nit; both round-2 fixes confirmed).
- `testing/kit-preamble.sh`: the comment above `unset -f` said the scan holds against a function "a kit or an agent" defined. That claimed more than the code does. It now says the scan trusts the kit's own shell: a kit that redefines a builtin (`printf`, `read`, `[`, `unset`) can still blind it, and the agent under test defines nothing there. The unused function-level `local rc` is removed.
- Comment and declaration only. `TestKitPreambleRefusesLeftoverCredentials` passes, and the gate's verifier runs after this commit.
