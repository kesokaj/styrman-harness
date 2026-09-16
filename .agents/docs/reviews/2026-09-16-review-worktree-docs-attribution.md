---
owner: reviewer
date: 2026-09-16
---
# Review — worktree-docs-attribution   (2026-09-16)

Range reviewed: working tree

Docs lens on the uncommitted change (`git diff` / `git status --short`, excluding `.agents/docs/reviews/` and `.agents/docs/testing/`): `.agents/docs/architecture/DECISIONS.md`, `.agents/docs/deploy/local.md`, `.agents/docs/deploy/remote.md`, `README.md`, `REQUIREMENTS.md`, `ROADMAP.md`, `cmd/sh-agent/main.go`, `cmd/sh-agent/main_test.go`, `context/AGENTS.md`. Verified `gofmt -l`, `go vet ./...`, and `go test ./cmd/sh-agent/... -run TestAttribution` — all clean/pass. Grepped `attribution` across README.md, context/, agents/, skills/ (excluding VENDORED.md entries), `.agents/docs/deploy/`, REQUIREMENTS.md for stale wording — no leftover "never touches vendor settings" / "add it yourself" / "say so while it is missing" phrasing found outside this diff's own new text.

## Verdict: REQUEST CHANGES

Three of the four docs that restate `status`'s note-firing condition (README.md, `.agents/docs/deploy/local.md`, `context/AGENTS.md`) narrow it to "the key is present with a different value," but the code — and the diff's own new test, `TestAttributionSet/status_alone_keeps_the_note` — show the note also fires when `settings.json` is simply absent (e.g., before first install, or after an operator deletes the file). That is a statement that does not match the code it describes. A second, smaller issue: README.md line 13 lost its antecedent noun and now reads as a dangling "it."

## Findings

- **[medium / confidence: high]** `README.md:22`, `.agents/docs/deploy/local.md:35`, `context/AGENTS.md:45` — each states the `status` advisory note fires "when `~/.claude/settings.json` sets `attribution` to a different value" / "only when the key is present with a different value" / "when the file sets `commit` to anything else." Tracing `attributionOff` → `attributionHint` in `cmd/sh-agent/main.go`: when the file does not exist, `attributionOff` returns `(off=false, problem="")`, and `attributionHint` prints the note because `!off` — i.e. it fires on an absent file too, not only on a present-but-wrong value. `TestAttributionSet/status_alone_keeps_the_note` in the same diff exercises exactly this (deletes the file after install, then asserts `status` still prints "does not set"), confirming the behaviour. All three doc statements should add the absent-file case, or say "does not set the key to an off value" rather than implying the key must be present.
- **[medium / confidence: high]** `README.md:13` — "Make sure `~/.local/bin` is on your `PATH`. Requires Go 1.24+ and git; nothing else. `install` writes it when absent; status says when a different value is set (see context § Git)." The old version's lead-in ("If `~/.claude/settings.json` does not set `\"attribution\": {\"commit\": \"\"}`…") was dropped, so "writes it" and "a different value is set" now have no antecedent in this paragraph — the previous sentence is about the Go/git requirement, not attribution. Reads as a dangling pronoun to anyone hitting this paragraph cold; needs the noun (`attribution`/the key) restored.
- **[low / confidence: low, unverifiable]** `.agents/docs/architecture/DECISIONS.md:25` (ADR-025 Consequences) — "The Mac receives the key on its next `sh-agent update`." Traced `update` → `harness.Reinstall` (re-execs the rebuilt binary's own `install` code path), which is consistent with this claim, but I did not run `update` end-to-end on a second machine/checkout to confirm; noting as unverified rather than a defect.

## Criteria

1. Dead code — n/a — docs lens.
2. Flow tracing — n/a — docs lens (code was read only to check doc claims against it).
3. Standards conformance — fail on the "describing-text rule": the `status` note-condition restatement changed in three places but landed narrower than the code in all three (see findings above); no third restatement was left un-updated, but all three were updated to the same inaccurate wording.
4. Directory layout — n/a — docs lens.
5. Simplicity — n/a — docs lens.
6. Decoupling — n/a — docs lens.

## Not examined

- `.agents/docs/plans/2026-09-16-attribution-write.md` and the three round-N cold-read reviews under `.agents/docs/reviews/` — out of scope per the dispatch (reviews/testing excluded) and not diffed against code.
- `agents/*.md` and `skills/*/SKILL.md` were grepped for the word "attribution" only (no hits besides vendored-license noise); not read end-to-end for unrelated drift.
- Did not run `sh-agent update` on a second checkout to verify the ADR-025 "Mac receives the key" claim (see low-confidence finding above).
- `.agents/docs/deploy/gcp.md` and other deploy files not touched by this diff were not re-checked (out of scope).

## Outcome

- Finding 1 (medium) — fixed: `README.md` § Install and the `status` row, `.agents/docs/deploy/local.md` § Verify and `context/AGENTS.md` § Git now say the note fires when the file is missing or sets `commit` to anything else.
- Finding 2 (medium) — fixed: the README sentence names the setting again ("`install` writes Claude Code's attribution setting, `"attribution": {"commit": "", "pr": ""}` in `~/.claude/settings.json`, when the key is absent …").
- Finding 3 (unverified) — left as stated in ADR-025; the Mac's `sh-agent update` is the operator's next step and its outcome will be recorded in that session's log.
- Round 1 closed 2026-09-16.
