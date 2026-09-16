---
owner: reviewer
date: 2026-09-16
---
# Review — bug-035-docs (round 2)   (2026-09-16)

Range reviewed: working tree diff (`git diff HEAD --`) over `README.md`, `cmd/sh-agent/main.go`, `cmd/sh-agent/main_test.go`, `internal/harness/apply.go`, `internal/harness/harness_test.go`

## Verdict: APPROVE

Round 1's medium finding is fixed correctly and completely, and the fix introduces no new defect. All three message sites — `README.md:53`, `cmd/sh-agent/main.go:309`, `internal/harness/apply.go:37` — now state the same rule with the same qualifier ("move only rules meant for every machine … never credentials or personal or host-specific notes"), the two test assertions (`cmd/sh-agent/main_test.go:347`, `internal/harness/harness_test.go:257-258`) require that exact phrase and both pass, `go vet`/`gofmt -l`/`go test ./...` are clean, and no living doc retains the old wording. The one open item — BUG-035's record still says `status: open` / "Not fixed" while the code now matches its own remediation — is the expected next step (bugs agent closes it in this commit, per round 1's ledger and this repo's established fix→review→bugs→commit sequence), not a defect in this diff.

## Findings
- **[low / confidence: high]** `.agents/docs/bugs/BUG-035-claude-md-guidance-ships-personal-content-everywhere.md` and `INDEX.md` — record still reads `status: open` / "Not fixed" although the code fix matching its own remediation text is now in the working tree. Not blocking: this is the documented next step in the workflow (bugs agent closes it as part of this commit), carried over from round 1's ledger, not a new gap. Flagging only so it is not forgotten before commit.

## Verification performed
- `git diff HEAD -- <the five files>` read in full; all three message strings compared word-for-word.
  - README: "Move only rules meant for every machine into a rule file under `rules/` in this repository — `rules/` is deployed to every install and every downstream mirror, so credentials and personal or host-specific notes never go there —"
  - `main.go` note: "move only rules meant for every machine into a rule file under rules/ in the harness repo — rules/ ships to every install and mirror, never credentials or personal or host-specific notes —"
  - `apply.go` keep message: same core clause as `main.go`, "into the harness repo — rules/ ships to every install and mirror, never credentials or personal or host-specific notes —"
  - The three convey an identical rule (destination = rules meant for every machine; exclusion = credentials/personal/host-specific notes) with matching qualifiers; README's "deployed to every install and every downstream mirror" vs. Go's "ships to every install and mirror" is a wording difference, not a narrowing — both name the same scope (every install + every mirror).
- `gofmt -l .` clean, `go vet ./...` clean, `go test ./...` — both packages `ok`.
- `git grep` for the old wording ("move its content into a rule file") across the repo, excluding `changes/`, `reviews/`, `retros/`, `handoffs/`, `bugs/`: the only two remaining hits are in `.agents/docs/testing/2026-09-15-verify-harness-practices-done.md`, a dated testing log excluded by the docs-lens rule — no living doc retains the stale wording.
- Checked `.agents/docs/testing/sh-agent/guide.md`, `skill-matrix.md`, `.agents/docs/deploy/*.md` for any other copy of this guidance — none found; nothing else to update.
- Confirmed the newest `changes/` entry (`2026-09-16-bug-038-agents-git.md`) is for the prior, unrelated fix — no changes entry yet for BUG-035, consistent with review-before-record-close ordering.

## Criteria (docs lens — n/a for most, per the `--docs` template)
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens (wording-consistency check above stands in for it)
4. Directory layout — n/a — docs lens
5. Simplicity — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
Did not re-review the underlying code fix's Go correctness beyond compiling/testing it (that is the code-lens review's job, already delivered as `2026-09-16-review-bug-035-code.md`, not re-litigated here). Did not check commit history beyond the newest `changes/` entry. Did not examine files outside the five named in scope.

## Outcome

Two rounds; the report above is round 2's.

- **Round 1 [medium], the Go messages narrowed against `README.md:53`** — `fixed` at `cmd/sh-agent/main.go:309` and `internal/harness/apply.go:37`: both now say "never credentials or personal or host-specific notes", the README's qualifier; the assertions in `cmd/sh-agent/main_test.go` and `internal/harness/harness_test.go` require that phrase, and both tests pass. Round 2 confirms the three sites state the same rule with the same qualifiers.
- **Round 1 [low], BUG-035's record still open** — `fixed`: closed in this commit, the Fix section naming the final wording.
