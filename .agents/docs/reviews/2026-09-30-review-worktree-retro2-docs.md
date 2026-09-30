---
owner: reviewer
date: 2026-09-30
---
# Review — worktree-retro2-docs   (2026-09-30)

Range reviewed: working tree — changed files: skills/handoff/SKILL.md skills/retro/SKILL.md; untracked, not ignored: none; read the diff with git diff HEAD -- skills/handoff/SKILL.md skills/retro/SKILL.md.

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The diff applies Proposals 1–3 of `.agents/docs/retros/2026-09-30-retro-2.md` correctly at the two spots the retro drafted, but it does not carry the new "never list a standing by-hand procedure, never carry an item twice" rule into the other restatement of the same rule inside the very file it edits (`skills/handoff/SKILL.md` Step 3, line 64), nor into the template asset that copies the pre-diff wording verbatim (`skills/init-project/assets/TODO.md:3`). Both still tell a session, unconditionally, that "anything the operator must do... goes in this handoff's Next Steps" — the exact sentence the new Next Steps text (lines 125-135) narrows. A session that reads Step 3 first (it runs before Step 4, per the file's own order) will still write the excluded content, because Step 3 was never told about the exception. The retro-skill edit (Proposal 3) is clean: no other restatement of the harness-scope queuing default needed updating.

## Findings

- **[medium / confidence: high]** `skills/handoff/SKILL.md:64` — Step 3 still reads "Anything the operator must do, and any deferred idea worth keeping, goes in this handoff's Next Steps, not in those files," unconditionally. This is a restatement of the exact rule the diff narrows two paragraphs later, in the same file: the new Next Steps text (lines 125-129) says a standing by-hand procedure recorded in the deploy guide must NOT be listed there even though it is something "the operator must do." Step 3 runs first in the file's own step order (it is "before writing anything," Step 4 is the template), so a session following Step 3 literally still routes the excluded content into Next Steps before it ever reaches the template's exception. Fix: narrow Step 3's sentence the same way, or point it at the Next Steps section instead of restating the rule.
- **[low / confidence: medium]** `skills/init-project/assets/TODO.md:3` — the scaffolded project `TODO.md` template says "anything the operator must do, or an idea worth keeping, goes in the newest handoff's Next Steps," the same unqualified wording the diff just narrowed in `skills/handoff/SKILL.md`. Any project scaffolded by `/init-project` after this change still ships the old, unnarrowed rule in its own tree. Not the diff's primary target and lower impact (a per-project copy, not the harness skill itself), but it is a restatement of the same subject noun ("goes … Next Steps") the brief asks to be grepped for.
- **[none — checked, consistent]** `context/AGENTS.md:18` ("an action only the operator can take goes in the handoff's Next Steps") is left unchanged by design (retro's own "Not proposed" section: the general rule stays in the context, the exception lives in the skill that does the writing). This is a defensible general→specific split and not flagged as a contradiction, but it is worth naming: read literally, the context's unconditional sentence and the skill's new exception still "read against each other" in exactly the shape finding 4 of the retro described for the pre-diff Step 3/deploy-guide pair. I do not treat it as a finding because the retro explicitly weighed and accepted this split (byte budget on the Antigravity block), but a future reviewer should know the same tension was deliberately left standing here.
- **[none]** `skills/retro/SKILL.md` Step 6 — the diffed "Harness scope" bullet is internally consistent with the Usage block (`--harness` at line 16 still describes consuming other projects' queued proposals; the new sentence covers proposals generated inside the harness checkout in the same run — a different case, not a restatement needing change) and with Step 1's harness-checkout detection (`sh-agent status` first line + `test -f context/AGENTS.md`, reused verbatim). No other bullet in Step 6 (Project scope, Memory verdicts, Rejected, Headless) restates the harness-scope default. No contradiction found.
- **[none]** `skills/preflight/SKILL.md:65,86`, `skills/bug/SKILL.md:15,49`, `skills/feature/SKILL.md:15,74,102`, `skills/review/SKILL.md:64`, `agents/changelog.md:96` all mention "the handoff's Next Steps" but for a different subject (a gate-blocking question with no reply available, or ordinary deferred work/follow-ups) — none of these restate "anything the operator must do" as the unconditional rule the diff narrows, so none needed a matching edit.
- **[none]** `.agents/docs/deploy/remote.md:211` ("Standing by-hand procedures, never pending items (those live in the newest handoff)") is the twin the new handoff text points at and already says the same thing — no drift between the two.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — fail: the diff leaves one direct, same-file restatement of the rule it narrows unedited (`skills/handoff/SKILL.md:64`), and one cross-file template restatement unedited (`skills/init-project/assets/TODO.md:3`); the global rule "every other restatement of a changed rule is changed with it" is not met.
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- The harness proposal queue file under `~/.config/sh-agent/proposals/2026-09-30-retro.md` was read for context (confirms all three proposals are still `queued`, expected pre-commit/pre-review) but its accuracy against the eventual commit SHA cannot be checked yet — the diff is uncommitted.
- Did not check every skill under `skills/` exhaustively for indirect restatements beyond a `git grep` on the rule's subject nouns ("Next Steps", "operator must", "standing by-hand", "carried"); a restatement worded without any of those terms could exist unfound. `sh-agent search` was not run for this (tool budget); the grep pass is the only coverage.
- Did not review `skills/handoff/SKILL.md`'s Step 4 template's other sections (Corrections, Agent notes, etc.) or unrelated parts of `skills/retro/SKILL.md` outside Step 1/6/Usage, since the diff does not touch them and the brief scoped the check to the named twins.

## Outcome
- 1 Step 3 restates the unqualified rule — fixed at `skills/handoff/SKILL.md` Step 3 (the qualifier, pointing at Step 4's Next Steps rule).
- 2 the scaffold's TODO template restates it — fixed at `skills/init-project/assets/TODO.md:3`.
- 3 `context/AGENTS.md:18` (general rule) vs the skill's exception — left: recorded in `retros/2026-09-30-retro-2.md` § Not proposed (the rendered context sits 4 bytes under Antigravity's block; the skill that writes the handoff carries the exception).
