---
owner: main
date: 2026-09-17
---
# Handoff — 2026-09-17 — delivery rework: Phase 1 done, Phase 2 one task in, session ended by a tool block

Session ID: 418afec9-9692-4e1b-a3b7-be09b67cf8aa
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Read this first

**This handoff was written without git.** Part-way through Phase 2 the session's permission classifier stopped evaluating shell commands — every one, including read-only — and said it would keep doing so for the rest of that conversation. The operator chose to start a fresh session rather than switch permission modes. So the state below is what was verified BEFORE the block, plus one file re-read afterwards with the file reader. **Run `git status --short --untracked-files=all` and `git log --oneline -3` first and trust those over this document.** This file itself is uncommitted for the same reason.

## Where the work is

The feature is `.agents/docs/plans/2026-09-17-delivery-rework.md`. It is committed, pushed, and it is the real source: interfaces, failure modes, every decision with its reasoning, and the phase checklist. Read it before anything else here.

**Committed and pushed, through `6af54f7`:**

- Phase 1 complete. ADR-030 (delivery is `sh-agent` and nothing else; both vendored mirrors retire; ADR-026's tag channel withdrawn as a delivery mechanism, tags stay as release markers) and `REQUIREMENTS.md` R8, approved by the operator on its exact text, committed together with the `ROADMAP.md` pointer as the trajectory rule requires.
- The plan, and the three-round cold read at `.agents/docs/reviews/2026-09-17-coldread-delivery-rework.md`.
- Everything earlier in the day: `v1.0.0` signed and pushed, the toolchain upgrade that cleared both advisories, an empty bug board.

**Uncommitted in the working tree** — two files, one finished task:

- `cmd/sh-agent/main_test.go` — `TestReorder` gained two cases, `report --feedback-repo /x --dry-run` and `report --dry-run`.
- `cmd/sh-agent/main.go` — `reorder`'s `takesValue` map gained `-feedback-repo` and `--feedback-repo`, with a comment naming why. Re-read after the block and confirmed present at roughly line 451.

That is Phase 2's first box. The test was run bare before the fix and failed with `["--feedback-repo" "--dry-run" "report" "/x"]` against the wanted order, then passed after. Both states were observed, so the box can be ticked on that evidence. It was not committed and the box is not ticked.

`internal/harness/report.go` was the next thing to be written and the write was the command that got blocked, so **it does not exist**. Nothing is half-written.

## Resume here

1. Confirm the tree with git. Commit the two files above as Phase 2's first task and tick its box in the plan with the must-fail/must-pass evidence.
2. Write `internal/harness/report.go`. The plan's `## Interfaces` specifies it completely, and the redaction rules are the part to follow exactly rather than improvise: the needles LOCATE and never define the span. A whole JSON string value, a whole text line, or a whole `-----BEGIN`-to-`-----END` block is removed. This matters because `AKIA` is a prefix, so replacing the matched token alone strips the marker and ships the key, and a PEM body carries no needle at all. `json.Valid` is asserted per line after re-encoding. The operator chose this design at the cold-read cap over refusing-with-override and over not scanning; D9 records it.
3. Then the remaining Phase 2 boxes, then Phases 3 and 4. Nineteen boxes were left when this session ended.

## The one thing that needs attention now

The working tree is dirty, and `sh-agent update` refuses a dirty checkout on every path (ADR-027, shipped earlier today). That is the guard working, but it means **the operator's machines will not pick up harness changes until this tree is clean**. Committing the two files above clears it. The hourly timer on this box will decline meanwhile and say so in `systemctl --user --failed`.

## Key decisions this session, all recorded in the plan

- Delivery is `sh-agent` only, four places, on the operator's own credentials; no deploy keys, because they authenticate as themselves everywhere. Read-only consumption is structural rather than a permission: everything downstream of the four checkouts receives installed FILES, not a repository.
- No webservice. Git carries both halves; a service adds hosting, per-machine auth and an attack surface for one user, against R3, R4 and local-by-default.
- Versioned release bundles are a follow-up with their own plan. They would formalise the zip the operator already uses, need no CI because the binary cross-compiles locally, and do NOT remove the credential requirement while the repository is private. `Open:` row in the plan, defaulting to private.
- The cold read never reached both readiness tokens and ended at its three-round cap, as the workflow prescribes. Twelve blocking findings, each verified against the code before being accepted; four changed the design. That is normal and recorded, not a failure.

## Still outstanding, not this plan

- **styrman-crew's Session B is building the signed-tag channel this supersedes.** It needs to be stopped and rescoped: a checkout on the host, `sh-agent update`, the entrypoint seeding project homes from it, the mirror retired. The operator has been told twice and it is still open.
- argolis-devbox Session C is aligned already and needs nothing new.
- Two advisory docs nits from the v1.0.0 gate: `README.md` and `AGENTS.md` do not mention the two new `install` refusals.
- The full skill matrix has never run from a real home on either CLI.

## Corrections logged this session

The operator set the form their questions take, on 2026-09-17: "When you want me to answer anything give me a short list and a short explanation and tell me what you recommend." Logged to `~/.config/sh-agent/feedback.jsonl` and saved to memory. Apply it: short list, one line of why per option, recommendation named.
