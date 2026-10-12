---
owner: main
date: 2026-10-12
---
# Verification — real-home install at f0bc3e6, the guide § 5 smoke on both CLIs, and one Antigravity run of each changed skill (v2.7.0 gate)

Range reviewed: v2.6.2..f0bc3e6

**The install.** On 2026-10-12 at 07:04 the real home on the Linux box had become a release-binary place on v2.6.2. The operator chose to install the candidate there for the smoke (2026-10-12, "Install checkout there (Recommended)"), as the v2.6.0 gate did. The authority is the project's Release rule and the operator's release instruction of 2026-10-12 ("Go all, fit all into next version preflight").
- A plain `go run ./cmd/sh-agent install` was refused, by design: "this home was installed from the release binary v2.6.2, not from a checkout; pass --repo …".
- `go run ./cmd/sh-agent install --repo <this checkout>` exited 0: 11 written, 1779 unchanged, 0 removed.
- After `hash -r`: `type -a sh-agent` shows one `sh-agent`, `~/.local/bin/sh-agent`.
- `sh-agent version` printed `f0bc3e67af7c+dirty`.
- `sh-agent status`: installed from commit `f0bc3e6+dirty`, binary `f0bc3e67af7c+dirty`, drift none (1789 files match), no `shadowed` line.
- "+dirty" is the verifier's uncommitted log under `.agents/docs`, which installs nothing.

## Status: PASS

11 of 11 rows passed, each against its own criterion.

### Real-home smoke (guide § 5)

The three checks of `testing/2026-09-03-verify-real-home.md`, plus a fourth, `pointer`. The pointer check is the context question on the new wording that the verifier of ee20d00 left BLOCKED as tier 9 (guide § 4) and the v2.7.0 gate verifier named as the main agent's.
- The script is `.agents/scratch/realhome/smoke-v270.sh`, the v2.6.0 gate's `smoke.sh` with a `pointer` check added: `grep -qi 'init-project skill lists them'`.
- Each check ran in a fresh clone of the demo fixture under `~/sandbox/realhome`, outside every checkout.

| CLI | check | result | seconds |
|---|---|---|---|
| Claude Code | context: GANNET, 36 harness skills, the Antigravity poll rule "Yes", DoD step 2 quoted | 1 of 1, denied 0 | 18 |
| Claude Code | skill: Go review names `golang`, changes nothing | 1 of 1, denied 0 | 135 |
| Claude Code | agent: `verifier` over a comment edit, log written | 1 of 1, denied 0 | 94 |
| Claude Code | pointer: "(the init-project skill lists them)" quoted | 1 of 1, denied 0 | 10 |
| Antigravity | context | 1 of 1, denied 0 | 33 |
| Antigravity | skill: `golang` loaded, nothing changed | 1 of 1, denied 0 | 102 |
| Antigravity | agent: `verifier` log written within its 12-call budget | 1 of 1, denied 0 | 195 |
| Antigravity | pointer: "(the init-project skill lists them)" quoted | 1 of 1, denied 0 | 11 |

**The Antigravity store check** read a copy of the context run's conversation (`5046ad98…`, written at 11:44:46, the minute its output was written). The copy was deleted by name after the check.
- `<RULE[~/.gemini/config/rules/AGENTS.md]>`, `<RULE[…/critical-thinking.md]>` and the demo project's own `AGENTS.md` block each appear 2 times, once per stored copy of the prompt, so there is one block of each per prompt.
- GANNET is present 14 times.
- There is no `<truncated` line.
- The new wording "the init-project skill lists them" is present.

### One Antigravity run of each changed skill (AGENTS.md § Release)

The range changes `skills/init-project`, `skills/review` and `skills/harness-audit`, and no agent.

| Row | kit | result | seconds |
|---|---|---|---|
| scaffold (init-project) | `testing/agy-rel/run.sh scaffold`, `AGY_TAG=v270` | PASS; no processes left, real config unchanged, no credentials left | 272 |
| review | `testing/agy-rel/run.sh review`, `AGY_TAG=v270` | PASS; no processes left, real config unchanged, no credentials left | 198 |
| harness-audit | `.agents/scratch/realhome/harness-audit-v270b.sh` (`agy-rel/run.sh`'s shape; print timeout 45 m, `GOCACHE` outside `KIT_HOME`) | 1 of 1 against the row's criterion; no processes left, real config unchanged, no credentials left | 1308 |

None of the eleven rows showed an `auto-denied` line or an `invalid allow entry`.

**The harness-audit row took three tries.**
- **The first run (`harness-audit-v270`)** hit agy's 25-minute print timeout with its turn in progress: empty output, two recon briefs, no report. Its credential line flagged nine files:
  - eight Go build-cache objects, holding this repo's own test fixtures (the deliberately split fake AWS key at `catalog_test.go:1592`, the redaction tests' PEM headers and their Anthropic, GitHub and Google token shapes, packed together by the compiler);
  - the conversation database, whose two refresh-token hits are the kit's own regex text, which the agent read in `kit-preamble.sh`.

  Every hit was checked by pattern and against `git grep`, without printing a value. The files were deleted by name, as the kit says, and the scan re-run on that home: "no credentials left in KIT_HOME".
- **The re-run (`harness-audit-v270b`)** had a 45-minute print timeout and the Go build cache under the run directory, outside `KIT_HOME`. agy exited 0 after 1308 s with the recon brief, the per-area brief (AG ids, severities, `file:line`), the report and the fix plan with `### Phase` sections. The kit printed FAIL because the agent had COMMITTED its artifacts in the throwaway clone (`9133574`), and the kit's check finds the new files only through `git status`.
- **The criterion, applied to what the run made.** The row's criterion, unchanged, was applied to the committed and uncommitted new files together (`git diff --name-only f0bc3e6..HEAD` plus `git status`): PASS.
  - The widened check fails on the first run's output, which has no report, and passes here, so it is validated in both states.
  - The clone's commit was never pushed: this checkout's HEAD stayed `f0bc3e6`.
  - A committed kit row that counts committed files is a follow-up (the handoff's Next Steps).

## What was not checked

- The Mac (`2026-09-04-verify-mac-real-home.md`): no Mac session here. Checkout places pick the release up at their next `sh-agent update` (`deploy/remote.md` § Operator steps).
- The § 3c matrix in full (absent kit). The gate verifier checked the changed refusal messages instead.
