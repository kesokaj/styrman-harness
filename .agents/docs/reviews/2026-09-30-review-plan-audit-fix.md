---
owner: reviewer
date: 2026-09-30
---
# Review — plan audit-fix, cold read   (2026-09-30)

Range reviewed: `.agents/docs/plans/2026-09-30-audit-fix.md` (cold read, no diff), checked against the code and files it names.

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The plan is not implementation ready. The worst problem is the C1 headline: the listed context moves cannot reach "render <= 14,000 bytes". They save about 4,600 bytes from 21,906, which gives about 17,300, and no further cut is named. Nearby, the plan contradicts itself on `/goal` and on the warn threshold, and several verify clauses cannot produce the condition they claim (tools not installed, a flag the code rejects, a grep that matches by design). Most of the remaining gaps are NON-BLOCKING.

## (1) Gaps

- **[high / high]** C1 arithmetic. `context/AGENTS.md` render is 21,906 bytes (`sh-agent render | wc -c`). The plan needs -7,906. The lines it names are 16 (1,400 B), 34 (1,156), 52 (941), 56 (796), 57 (567), 58 (687), 59-62 (1,356), and each is cut to a stated size. That sums to about -4,600, and the G11 paragraph adds back about 250. The Phase 3 verify ("render <= 14,000") will fail, and the implementer must pick further rules to drop, which is an operator-level choice. Fix: list the extra cuts, or set the target at about 17,500, or name which whole rules leave.
- **[high / high]** `/goal` contradiction. Interfaces (line 36), D17 and the `C15` note say keep `/goal`. The Phase 3 task (line 110) says "`/goal` dropped".
- **[medium / high]** Warn threshold stated three ways. Interfaces and Phase 3 say 15,000. Failure modes says 18,500. D2a says the 18,000 target is superseded. `TestRenderedContextTarget` "warns, not fails", so it cannot fail and its verify ("prints the byte count") verifies nothing. The test is extra scope with no one asking (D2: "implementer"). The `render | wc -c` check already covers it. Fix: drop the test, or pin one number.
- **[high / high]** Phase 6 verify `grep -rn 'python\|anti-ui-slop' context skills/*/SKILL.md README.md` expects "only the frontend Do-NOT line". It already matches `skills/init-project/SKILL.md:115` (python3, which D13 keeps and WS3 documents more) and `skills/gke-basics/SKILL.md:49` ("kubernetes python"). It cannot reach the expected result. Fix: anchor the grep to the skill names (`/python`, `skills/python`, `anti-ui-slop`).
- **[high / high]** `--kind spec,plan` (research skill, Phase 5 WS10). `internal/search/search.go:84-92` accepts exactly one kind and errors on a comma list. The plan either needs a code change and a docs row in `specs/sh-agent/search.md`, or should read `--kind spec` and a second call.
- **[medium / high]** AG12b / "the one skill over 600 today (cloud-run-deploy, 1,007)". Not true. Single-line descriptions of `skills/dream` (691), `skills/project-memory` (629) and `skills/review` (775) also exceed 600, measured on the first `description:` line. The Phase 4 verify "`sh-agent list` expect no note" then fails unless those are trimmed. Nothing in the plan trims them. There is also a consistency trap: `Lint` failing at > 1,024 makes `install` refuse any vendored skill that is over the limit, and vendored skills may not be edited in place (AGENTS.md § Rules). The plan does not measure the vendored ones.
- **[medium / high]** IN2 (Phase 1). "manifest saved after each successful write", expect "manifest lists 1..N-1". `apply.go:95-110` removes stale files before writes and only saves state once at line 140. A manifest of just the written files drops the old install's not-yet-rewritten files from the manifest while they remain on disk: orphans, the defect the task is meant to fix. The manifest after a failure should be (old entries still on disk) ∪ (written). `SaveState` per file is also O(n) rewrites; saving once in an error path would do. The expected result and the approach both need stating.
- **[medium / high]** Tools not installed. `terraform`, `tflint` and `gcloud` are MISSING on this machine (`command -v`). Phase 6 verifies `terraform validate` and `gcloud run deploy --help`, and the Risks section says gcloud flags are "checked by `gcloud run deploy --help`". Those cannot run. Fix: a step to fetch them into scratch, or say those rows are reported as unverified.
- **[medium / medium]** Antigravity `skills:` key (D5, AG5/AG15). Phase 4 says "the agy behavioural run below proves the skill loads", but no agy behavioural run exists below in Phase 4; only `claude -p --agents` runs and a sandbox start. Phase 9's matrix covers skills, not the derived agent key. D5's "drawn solid only after one agy behavioural run" has no task. The reviewer would also carry two mechanisms for one operation (the `skills:` key and the explicit `~/.claude/skills/review-standards/SKILL.md` read). Antigravity hangs on unrecognised agent keys (AGENTS.md § Rules, toolMap), so a wrong `skills:` path form is a hang risk; the sandbox-start task sits last in the phase. Fix: move it first and add the agy run row.
- **[medium / medium]** C1 moves rules an author follows into `review-standards` (lines 56-62, "a review checks and an author follows"). `review-standards` is loaded by reviewers only. The four Go-flavoured rules also go to `skills/golang`, but TypeScript, bash and docs authors get none. State where the author-side copy lives for non-Go work, or say it is deliberately review-side only.
- **[medium / medium]** The derive/Lint changes alter `install` behaviour (refusal over 1,024; a `note:` line in `list`; `LoadAgents` refusing descriptions), and Phase 1 changes exit codes and refusal output (IN11-IN13). No docs clause names `specs/sh-agent/*.md`, README, or `.agents/docs/testing/guide.md` for any of these. Phase 5 adds `/review --resume` and `distill --from-log/--global` with no docs clause either (Phase 7 does have README rows for new skills).
- **[low / medium]** Phase 4 verifier behavioural check "log carries 'Step 4' marker". Nothing in the plan says the new definition instructs the agent to write a step marker. Without it the run proves nothing (AGENTS.md § Testing: the artifact must carry something only the new definition writes). The same applies to the changelog check (the agent would have to name the SHA it was given; spelled out there, but not for the verifier).
- **[low / medium]** SD7 `VACUUM INTO` "on a read-only connection". The earlier design copies the DB files so the operator's live Antigravity DB is never opened. A read-only WAL connection may still create or touch `-shm`. The plan needs a statement that the live DB is not modified, and a fixture that proves it (`TestCopyDBWhileWriting` only checks the copy opens).
- **[low / medium]** IN8. The flock helper moves from `internal/dream/run.go:514` into `internal/harness`. `internal/harness` has a `!unix` stub file (`owned_other.go`), so a non-unix stub is needed. Lock path, blocking or not, and behaviour on contention are undecided. Release builds are only linux/darwin (remote.md:139), so this is low.
- **[low / medium]** Interfaces text left as thinking-out-loud: "`sh-agent memory init-settings <root>`? — NO", "a `.agents/rules/` line is not global; ... ? No: ..." (grilling, WS21). The final decision is recoverable but the content to write for the VENDORED note and the keep-list is unstated; the claimed "(D-row.)" for grilling does not exist (D13 covers init-project only).
- **[low / medium]** `TestSkillLinksResolve`: "See <Title>" pointer rule is undefined, and it walks `skills/**/*.md`, including vendored skills that cannot be edited. Say which dirs it skips.
- **[low / medium]** TR1 `restatement-grep.sh`: inputs, output format and "living files" list are unspecified; the installed path on Antigravity (`~/.gemini/config/skills/review/scripts/`) differs from Claude's and skill bodies are not rewritten by derive (only agent bodies are, `derive.go:25,63`). "`kit-preamble.sh` sourced by every kit in guide § 4": kits there are prose procedures, not scripts.
- **[low / low]** E2E (a)/(b): AGENTS.md § Testing says a CLI started from a scratch HOME cannot authenticate; (a) and (b) say "both CLIs start and answer" / run `/init-project` then `/spec` then `/feature` without saying the home (the sandbox procedure or the real home after Phase 9). E2E (c) "sandbox `dream --project`" runs with fake CLIs on PATH per guide § 3e; say so.
- **[low / low]** Re-audit "published as Styrman Harness Audit III", verified by "page URL in the changes entry": the mechanism ("the same way as the two before") and a fallback if the external service is unavailable are not stated.
- **[low / low]** "Grep -c 'Write the artifact EARLY' ... expect 5 identical blocks": the "one shared wording (13 lines)" is not in the plan, and `explorer.md` has none today; the AG9 sentence "security is pinned to pro" likewise has no text to copy. NON-BLOCKING; the implementer writes it.
- **[low / low]** AGENTS.md says deploying to the real home is the operator's job; D15 cites "You can update the local harness everything you need it" as the authority. One line in the plan saying this is a one-time exception (and does not amend the rule) would close it.

## (2) Scope — built beyond the request in `## Why`

The Why names: fix all findings, add features from the audit, phases, verified, a PRIORITY_TODO, re-audit. Covered by "add new features from your audit report": the seven new skills (D7), `--resume`. Not named and with no D-row naming who asked:
- `TestRenderedContextTarget` (warn-only test; D2 is "implementer") — see above; it cannot fail.
- Two thresholds (`Lint` fails > 1,024 AND `list` notes > 600) for the same property: a second shape for one check (D19 cites the format brief, not a request).
- Two mechanisms for the agent-to-skill link on Antigravity (`skills:` key AND a body read-path AND the Budget prose line), when the Budget prose line duplicates `maxTurns` that the dispatcher already names per the context.
- `explain-project`, `testing-strategy` and `operate-gcp` write to `docs/onboarding.md` etc.; fine as audit features, but `explain-project` (a human walkthrough) has no audit ID cited in the plan (D7 lists G1, G2, G5-G9 only).

## (3) Readiness — questions I would have to ask

1. BLOCKING. C1: which additional rules leave the always-loaded context to reach <= 14,000, given the listed moves reach about 17,300? Or is the target 17,500?
2. BLOCKING. `/goal`: kept (Interfaces, D17) or dropped (Phase 3)?
3. BLOCKING. IN2: what must the manifest list after a failure on file N (old entries on disk plus written), and is it saved per file or once in an error path?
4. BLOCKING. `--kind spec,plan`: change `search.go` to accept a list (then a test, a spec row, and a docs clause), or use one kind?
5. BLOCKING. Phase 6 python/anti-ui-slop grep: expected result is unreachable as written; what should it match?
6. BLOCKING. AG12b: trim `dream`, `project-memory`, `review` (and any vendored over 1,024) or raise the note threshold? What happens when a vendored skill is over 1,024?
7. BLOCKING. Where does the agy behavioural run for the derived `skills:` key live, and what is its pass criterion? What if the vendor key form is wrong (hang risk)?
8. NON-BLOCKING. Terraform, tflint and gcloud are missing: fetch to scratch or report unverified? (Plan default: verify clauses assume they exist.)
9. NON-BLOCKING. Warn threshold number (15,000 or 18,500) and whether the test stays.
10. NON-BLOCKING. SD7: does the live Antigravity DB get opened even read-only? IN8: lock path and contention behaviour. SD9 "a non-JSON string argument corpus test": what corpus?
11. NON-BLOCKING. Content for grilling's VENDORED note, the AG9 "security is pinned to pro" sentence, the early-write block, the verifier step marker.
12. NON-BLOCKING. `TestSkillLinksResolve` "See <Title>" rule and skipped dirs; `restatement-grep.sh` interface and Antigravity path.
13. NON-BLOCKING. Which home E2E (a)/(b) run in; docs clauses for exit codes, `list` note, `install` refusal, `--resume`, `--from-log`.
14. NON-BLOCKING. Re-audit publish mechanism and fallback.

## Criteria (applied to the plan as a design document)

1. Dead code — no dead code yet; the plan strikes `python`/`anti-ui-slop` and deletes dead pointers; `TestRenderedContextTarget` is a test that cannot fail.
2. Flow tracing — fail: IN2 manifest semantics, `skills:` derive chain has no agy run, lint/install refusal of vendored skills not traced.
3. Standards conformance — fail: docs clauses missing for changed commands/exit codes; `/goal` restated both ways; AGENTS.md real-home rule needs an explicit exception line.
4. Directory layout — pass: `skills/review/scripts/`, `assets/agents-tree/docs/specs/`, `.agents/docs/testing/kit-preamble.sh` are where conventions put them (Antigravity path for scripts to confirm).
5. Simplicity (DIY-except-security) — mixed: stdlib flock reuse and no new dependency is right; duplicate mechanisms (two thresholds, `skills:` key plus read path plus Budget line) are the simplicity risk. `api-auth-gcp` names vetted libraries only (no hand-rolled crypto).
6. Decoupling — pass: Phases 1, 3-7 and 8 touch separate files; Phase 8 rewrites `main.go` after Phase 1 edits it (serialized, stated by order).

## Not examined

The recon brief and the five research briefs were read only where needed to check cited facts (not line by line); the research briefs' technical claims (Terraform, Cloud Run flags, Tailwind, Vitest versions) were not verified. The Graduated row count "46 rows" (TR1) was not reliably counted. `skills/**` vendored description lengths over 1,024 were not measured. The real content of G1-G13 and WS/AG IDs was not cross-checked against the recon brief beyond IN2, IN3, IN4, IN8, IN10, SD7 and C1.

design needs revision
implementation not ready
