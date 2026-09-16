---
owner: reviewer
date: 2026-09-16
---
# Review — worktree-docs   (2026-09-16)

Range reviewed: working tree

## Verdict: REQUEST CHANGES

The four docs subjects (crypto exemption, `Not examined`, cold-read cap/scope/`supersedes`, init-project template) are consistently restated everywhere `git grep` reaches, with no stale old wording left behind. The one real defect: four docs describe the new `sh-agent install`/`status` attribution note as requiring `"attribution": {"commit": "", "pr": ""}` in `~/.claude/settings.json`, but the code that produces that note (`cmd/sh-agent/main.go`) checks and prints only `"attribution": {"commit": ""}` — the `pr` key is never read, never required, and never appears in the printed message. An operator following the docs' quoted JSON literally would be following a shape the tool does not actually check.

## Findings
- **[medium / confidence: high]** `README.md:14`, `README.md:22`, `.agents/docs/deploy/local.md:31`, `.agents/docs/deploy/remote.md:19`, `context/AGENTS.md:45` — all state the setting the tool checks for is `"attribution": {"commit": "", "pr": ""}`. The actual implementation (`cmd/sh-agent/main.go:488,490`, `attributionOff` at `main.go:499-524`) only inspects `Attribution.Commit`; the printed note text is literally `does not set "attribution": {"commit": ""}` — no `pr` key anywhere in the code or in `main_test.go`'s expected `body` string (`cmd/sh-agent/main_test.go:373`). The `pr` key appears only inside one test's settings *input* fixture (`main_test.go:375`), never checked. Fix: either make the docs quote the actual checked/printed shape (`{"commit": ""}`, optionally noting `pr` as an extra recommendation for PR footers with a separate rationale), or make the code also check `pr` and update the printed message to match what the docs claim.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Subjects checked (git grep against living docs: README.md, AGENTS.md, context/, agents/, skills/ (excl. VENDORED.md entries), .agents/docs/specs/, .agents/docs/deploy/)
- Hand-rolled-crypto exemption (four restatements: `agents/reviewer.md:43`, `agents/security.md:33`, `context/AGENTS.md:49`, `skills/review-standards/SKILL.md:55`) — all four carry the same expanded "primitive or protocol implemented or composed here" / stdlib-`crypto/*`-is-not-DIY wording; no stale short-form left in any living doc (only unrelated non-crypto "hand-rolled" mentions remain, in `skills/golang/references/`, about worker pools and migrations).
- `Not examined … or nothing` — `agents/reviewer.md:89`, `agents/security.md:71`, `agents/verifier.md:87` all now require naming the gap and reserve `nothing` for full coverage; `agents/explorer.md` never had a `## Not examined` section body to update (only the pre-flight sentence, itself unchanged and consistent).
- Cold-reader cap / scope rule / `supersedes` D-row — all three present and mutually consistent in `skills/feature/SKILL.md` (Step 2 template, Step 3 brief and cap text, Step 6 Definition of Done); no other living doc restates the old two-question / no-scope-check / no-supersedes wording (the only hit for old "Answer two things" text is `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`, a dated plan, out of the living-docs set by the dispatch's own exclusion list).
- Cold-reader brief "Answer three things" — numbered (1) Gaps (2) Scope (3) Readiness, correctly renumbered and internally consistent.
- init-project `local.md` template placeholders and `## Test` section — `skills/init-project/SKILL.md:75` and `skills/init-project/assets/agents-tree/docs/deploy/local.md` agree: Build/Run/Test verbatim from setup answers, Configure row deleted when no settings, Verify left uncommitted until the endpoint exists; the SKILL.md question list (`skills/init-project/SKILL.md:47`) does ask for build/run/test commands, matching the new `## Test` section.
- `sh-agent install`/`status` settings.json note — present and cross-referenced in README.md, both deploy guides and `context/AGENTS.md` (see finding above for the one inaccuracy).

## Not examined
Did not run `go test`/`go vet` (docs-lens dispatch, code correctness out of scope beyond checking docs-vs-code text match). Did not check `.agents/docs/deploy/gcp.md` (unchanged, not in diff, and not named in the scope-narrowing subject list). Did not review `.agents/docs/reviews/`, `.agents/docs/testing/`, `.agents/docs/plans/` per the dispatch's exclusion. Did not check whether an ADR is owed for the attribution-note behavior (outside the named subjects).

## Outcome

- Finding 1 (medium) — fixed: `README.md` (§ Install and the `status` row), `.agents/docs/deploy/local.md`, `.agents/docs/deploy/remote.md` and `context/AGENTS.md` § Git now name the shape the tool checks and prints, `"attribution": {"commit": ""}`; the context sentence and the README § Install say once that `"pr": ""` beside it is the optional extra for PR descriptions and that the note checks `commit` only. The tool is unchanged: the Git rule is about commit authorship.
- Not examined "whether an ADR is owed" — no: the note is a status line, not a structural decision; the decision that settings.json stays the operator's file is R4 and is recorded in BUG-039.
- Round 1 closed 2026-09-16; no second round owed.
