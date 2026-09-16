---
owner: reviewer
date: 2026-09-16
---
# Review — retro-2-release-procedure-docs   (2026-09-16)

## Round 1

Range reviewed: `8b95f22..worktree` (docs lens). Changed/tracked: `context/AGENTS.md`, `agents/verifier.md`, `skills/review/SKILL.md`, `skills/preflight/SKILL.md`, `skills/retro/SKILL.md`, `AGENTS.md`, `README.md`, `.agents/docs/deploy/remote.md`, `.agents/docs/architecture/DECISIONS.md`. Untracked: `.agents/docs/retros/2026-09-16-retro-2.md` (read whole).

### Verdict: REQUEST CHANGES

The worst problem: the new release procedure is documented in three places (ADR-026, `remote.md` § Releases, `AGENTS.md` § Release) as something `/preflight` executes at "Step 6" — SSH-signed tagging, a pre-push `verify-tag` check against the allowed-signers line — but `skills/preflight/SKILL.md` Step 6, the actual text an agent runs, was not touched by this change set and still says only "tag **that** commit" / "Push the branch and the tag", with no signing, no `ssh-keygen` capability check, and no verify-tag gate. An agent following the preflight skill literally, on this project, would cut a plain unsigned tag — exactly what ADR-026 says a release must never be — and the styrman-crew channel (signature-gated) would refuse it. This is a real gap between the doc set and the mechanism it describes, not a style nit: the six restated-rule checks the dispatch named are otherwise consistent (see Findings), so this is the one place the review-standards docs lens is built to catch.

### Findings

- **[high / confidence: high]** `skills/preflight/SKILL.md:79-90` (Step 6) — Step 6 has no reference to a project's own signed-tag procedure, unlike its Step 6 item 5 ("Deploy via the project's deploy guide"), which already has that hook for deploy. `ADR-026` (`.agents/docs/architecture/DECISIONS.md:214`) and `.agents/docs/deploy/remote.md:46` ("Procedure (what `/preflight` runs at Step 6 here...)") both assert that `/preflight` performs the `ssh-keygen` capability check, the per-command SSH-signed `git tag -s`, and the pre-push `verify-tag` check — none of which exist in `skills/preflight/SKILL.md`. `AGENTS.md:83-85` § Release repeats the same claim ("cut only by `/preflight vX.Y.Z`... checked before the push exactly as styrman-crew's channel checks it... with the exact commands in `remote.md` § Releases") without preflight itself pointing there. Fix: add a sentence to Step 6 item 3 (or a new item) instructing the agent to check the project's `AGENTS.md`/deploy docs for a release/signing procedure before tagging, the same way item 5 does for deploy — otherwise this is `dropped: skills/preflight/SKILL.md:85` (the signing step has no home in the skill that is supposed to run it).

- **[medium / confidence: high]** `.agents/docs/retros/2026-09-16-retro-2.md:24` — Proposal 1 states "Reused as drafted; the anchor clause is unchanged at line 50," but the applied text differs from the queued crew draft by more than the anchor: the crew draft (`~/.config/sh-agent/proposals/2026-09-16-retro.md` § From styrman-crew, Proposal "The docs lens greps the test suites too") reads `…the project's test suites as the same \`AGENTS.md\` Testing section names them…`; the text actually applied in `skills/review/SKILL.md:167` reads `…the project's test suites as its \`AGENTS.md\` Testing section names them…` ("the same" → "its"). The dispatch that produced this review was explicitly told this word change exists and must be recorded; the retro record does not say so. Fix: amend Proposal 1's line to note the wording change from the queued draft (this is a record-accuracy defect, not a defect in the applied text itself, which reads fine).

### Criteria

1. Dead code — n/a — docs lens.
2. Flow tracing — n/a — docs lens.
3. Standards conformance (restatements) — the five explicitly-named rule edits check out with one exception: the crypto trigger-list (preflight Security bullet, `context/AGENTS.md:16`, `agents/security.md:3,33`, `agents/reviewer.md:43`, `skills/review-standards/SKILL.md:56`) is now uniformly "auth, crypto, secrets, IAM..." everywhere — consistent, BUG-023's gap is closed in text (see Not examined re: index status). The verifier's "Long suites, once" bullet has no other restatement in the repo to conflict with — pass. The review skill's `--docs` grep list now names test suites and matches how this very dispatch's Lens paragraph is worded — pass. The context's post-dispatch untracked-files sentence (`context/AGENTS.md:51`) is new with no older restatement elsewhere to leave stale — pass. The retro skill's status-timing bullets (Step 6, both Project-scope and Harness-scope) were both updated in the same edit and agree with the pre-existing format at `skills/retro/SKILL.md:118` — pass. The one failure is not a restatement gap but a **procedure gap**: see the high finding above (ADR-026/remote.md/AGENTS.md describe a Step 6 that `skills/preflight/SKILL.md` Step 6 doesn't contain).
4. Directory layout — n/a — docs lens.
5. Simplicity (DIY-except-security) — n/a — docs lens. Noting only: the release procedure's crypto content (SSH-signed tags via `git`'s built-in `-c gpg.format=ssh`) is git's own signing mechanism, not hand-rolled — no DIY concern.
6. Decoupling — n/a — docs lens.

### Not examined

- `REQUIREMENTS.md` — read; ADR-026's claim that R3/R4 are unaffected checks out, no contradiction found, but I did not re-verify every other requirement line against the new release text beyond R3/R4.
- `.agents/docs/bugs/INDEX.md` — changed on disk mid-review (a concurrent process); not re-diffed against the retro record's Proposal 3/6 closure claims (both proposals are still `Status: queued` in the retro record, so no index mismatch is expected yet, but I did not re-read the file's final state after the notified change).
- Live verification of the signed-tag rehearsal itself — not repeated; relied on `.agents/docs/testing/2026-09-16-verify-retro-2-and-release-procedure.md`, which independently rehearsed the must-pass/must-fail cases and appears sound.
- Sandbox CLI parse-and-load smoke checks for Claude Code/Antigravity — out of docs-lens scope and explicitly deferred to the main agent per that verification log.
- Cosmetic: `.agents/docs/architecture/DECISIONS.md:207` ADR-026 heading has irregular multiple-space run before the date — not a finding, noted only in case it's a paste artifact.

## Round 2

Range reviewed: `8b95f22..worktree` (docs lens, second pass). Changed since round 1: `.agents/docs/architecture/DECISIONS.md`, `.agents/docs/bugs/BUG-023-preflight-security-trigger-omits-crypto.md`, `.agents/docs/bugs/BUG-044-security-probe-debris-repo-root.md`, `.agents/docs/bugs/INDEX.md`, `.agents/docs/deploy/remote.md`, `AGENTS.md`, `README.md`, `agents/verifier.md`, `context/AGENTS.md`, `skills/preflight/SKILL.md`, `skills/retro/SKILL.md`, `skills/review/SKILL.md`. Untracked, read whole: `.agents/docs/retros/2026-09-16-retro-2.md`, `.agents/docs/reviews/2026-09-16-security-release-procedure-adr-026.md`, `.agents/docs/testing/2026-09-16-verify-retro-2-and-release-procedure.md`.

### Verdict: APPROVE

Both round-1 findings are fixed correctly and completely, and the four artifacts named in the round-1 fix summary — `skills/preflight/SKILL.md` Step 6.3, `remote.md` § Releases, `AGENTS.md` § Release, ADR-026 — now describe one procedure with no contradiction found between them; no other harness restatement of Step 6's tagging still describes a plain tag (`git grep` of "tag **that** commit" / "Push the branch and the tag" turns up only the one, now-fixed, `skills/preflight/SKILL.md:85-86` pair). The security review's six findings against the procedure (binding, armor-check exit status, sudo/cmp, the wrong-workaround caveat, select-then-verify fallback, and the public-key-only fingerprint check) are each visible as fixed text in `remote.md`/`DECISIONS.md`/`AGENTS.md`. One gap surfaced during this round — the verification log's original rehearsal of the fingerprint check used the superseded public-key-only form, so its "run verbatim" coverage claim was briefly false — but the testing file changed on disk mid-review: the verifier appended its own "Round 2" section that re-rehearses the corrected `ssh-keygen -y -f ~/.ssh/id_ed25519 | ssh-keygen -lf -` command (`.agents/docs/testing/2026-09-16-verify-retro-2-and-release-procedure.md:90`) plus the full rewritten block's `test` assertions, target/ancestor binding and all four must-fail cases, and closes the gap in the same file. Recorded below as history; no open finding remains.

### Findings

None open. One transient gap, resolved during this review round — recorded as history, not a live finding:

- **[medium / confidence: high, RESOLVED]** `.agents/docs/testing/2026-09-16-verify-retro-2-and-release-procedure.md:27` (round 1 text) — the rehearsal's original "Must-pass" fingerprint step ran `ssh-keygen -lf ~/.ssh/id_ed25519.pub` (public key file only), while the procedure it claimed to verify — `.agents/docs/deploy/remote.md:49` (and the by-hand echo at `remote.md:79`) — had been rewritten during round 1's fix to `ssh-keygen -y -f ~/.ssh/id_ed25519 | ssh-keygen -lf -`, deriving the fingerprint from the PRIVATE key specifically because "a copied `.pub` alone proves nothing." The log's Coverage section then over-claimed "every command in the procedure block run verbatim." Resolution observed in this same review pass: the testing file gained a "Round 2" section (`:62-127`) that re-rehearses the corrected fingerprint command at line 90, plus the full rewritten block's `test` assertions and all four must-fail cases, closing the gap in place. No edit of mine was needed or made.

### Criteria

1. Dead code — n/a — docs lens.
2. Flow tracing — n/a — docs lens.
3. Standards conformance (restatements) — pass, with the one exception above. The BUG-023 trigger-list fix (`skills/preflight/SKILL.md` Step 2 Security bullet now includes "crypto") is consistent with `context/AGENTS.md` and `agents/security.md`, its prevention rule is graduated into `skills/review-standards/SKILL.md` §3, and `.agents/docs/bugs/INDEX.md` moves BUG-023 to Graduated and removes it from the open table — consistent with the BUG-029/BUG-040/BUG-043 precedent already in that file (graduated bugs are not duplicated in the main table). BUG-044's row is graduated to `context/AGENTS.md` § Hard rules the same way, and that bullet (`context/AGENTS.md:51`) matches the text the bug record cites. The retro-2 record's Proposal 1 fix (round-1 medium finding) is applied correctly: line 24 now reads "Reused as drafted with one word changed: the crew draft's 'as the same `AGENTS.md` Testing section' became 'as its `AGENTS.md` Testing section' ('the same' had no antecedent at that point in the sentence)" — matches the applied text in `skills/review/SKILL.md`. No other restatement of "status `applied`" (the retro skill's Proposal 5 fix) exists elsewhere in the harness to leave stale (`git grep` over README/skills/agents/context/AGENTS.md finds none). `REQUIREMENTS.md` R3/R4 checked again directly: no contradiction with ADR-026's claim.
4. Directory layout — n/a — docs lens.
5. Simplicity (DIY-except-security) — n/a — docs lens. The security review already confirmed the SSH-signing mechanism is git/OpenSSH's own, not hand-rolled; nothing in round 2's diff changes that.
6. Decoupling — n/a — docs lens.

### Not examined

- Sandbox CLI parse-and-load smoke checks for Claude Code/Antigravity — out of docs-lens scope, deferred to the main agent per the testing log.
- `.agents/scratch/security-release-procedure/probe/` and `.agents/scratch/verify-retro-2/` — not inspected; both agents report them gitignored and self-contained, and cleanup is out of a docs-lens reviewer's remit.
- Live re-verification of `/etc/styrman/allowed_signers` on the crew host, or of the crew-side updater `deploy/update/styrman-update` — outside this repository, as both the security review and round 1 already noted.
- The corrected fingerprint command and the rest of the testing log's Round 2 rehearsal were not independently re-run by me; I read its output and cross-checked the commands quoted against `remote.md`'s current text (they match), but did not execute `ssh-keygen`/`git` myself a third time.

## Outcome

- Round 1 [high] `skills/preflight/SKILL.md` Step 6 has no hook to a project's release procedure — fixed at `skills/preflight/SKILL.md` Step 6 item 3: follows the project's `AGENTS.md` Release section or deploy guide for tag form, signing command and pre-push checks; tags HEAD with no target argument; confirms `git rev-parse <tag>^{commit}` equals `git rev-parse HEAD` before pushing.
- Round 1 [medium] retro record says "reused as drafted" for Proposal 1 — fixed at `.agents/docs/retros/2026-09-16-retro-2.md` Proposal 1 (the one-word change is recorded).
- Round 2: APPROVE, no finding. After round 2 the main agent changed the same lines once more for the security review's round 2 (the tag takes no target argument; `git log --oneline @{u}..HEAD` and `git tag` added to the post-dispatch check in `context/AGENTS.md`; a probe-repository sentence in `agents/security.md`); those edits were not re-reviewed by this lens — the operator was waiting on the session — and the main agent's own restatement grep for the superseded forms (`<release commit>`, "prints its SHA", `ssh-keygen -lf ~/.ssh/id_ed25519.pub`) over README, AGENTS.md, context/, skills/, agents/, deploy/ and ADR-026 found none.
