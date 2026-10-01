---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-p8-docs   (2026-10-01)

Range reviewed: working tree per `git status --short --untracked-files=all` (Phase 8 of the audit-fix plan: installer refactor in `cmd/sh-agent/` and `internal/harness/`, README rows, `sh-agent-matrix.md` row, search.md/delivery.md/guide.md edits); `git diff HEAD -- <files>`. Docs lens only. Round 2 adds the fixes made after round 1.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing broken. No living doc prints a command line that IN5 would now refuse, no living doc sends the dry-run unpack or signers file to the system temp dir, and `TestHarnessMap` passes. Three stale or missing description lines remain, all low.

## Findings
- **[low / confidence: high]** `.agents/docs/testing/guide.md:73,86` — § 3b still says a change to "`reorder`'s `takesValue` map" and that `TestReorder` "covers a new value-flag being added to `takesValue`". `takesValue` is now built from the verb's FlagSet inside `reorder` (`cmd/sh-agent/main.go:928-934`), so no map entry is hand-added any more. The heading and sentence describe a step that no longer exists. Fix: heading says "or `reorder`", and the sentence says a new value flag is picked up from its FlagSet and `TestReorder` covers it. Same stale wording in code comments `cmd/sh-agent/search_test.go:16`, `cmd/sh-agent/memory_test.go:65`, `cmd/sh-agent/main_test.go:54` (code, not a doc, so noted only).
- **[low / confidence: high]** `.agents/docs/specs/sh-agent/search.md` (end of § Commands, "Exit status: 0 with hits or without; 1 on a refusal or an error") — IN5 adds a third case: a flag the verb does not declare (for example `index --kind bug`) is a usage error, exit 2 with the verb's usage (`main.go` `run`, `exitCode(2)`; matrix row added). The spec's exit-status line omits it. Fix: add "2 on a usage error, including another verb's flag".
- **[low / confidence: medium]** `.agents/docs/specs/sh-agent/delivery.md:32` and `README.md` § Install from a zip — `update --from-zip` now writes its signers scratch beside the zip (`internal/harness/release.go`, `MkdirTemp(dir, ".sh-agent-verify-*")`), so verifying a release zip needs a writable zip directory, and with a missing `~/.config/sh-agent` the DRY run unpacks beside the zip (`main.go` `updateFromZip`). delivery.md says only "a temporary directory under `~/.config/sh-agent/`", true for the real run. Neither doc states the zip-directory requirement or the dry-run exception. This is the command's documentation, not the delivery claim, so it belongs in the README zip section or `deploy/local.md`, not in rule 6's home. Fix: one sentence.

## Evidence for the checks you asked for
- Task 2 (IN5): ran `grep -rnE 'sh-agent <verb> … --<other-verb flag>'` over README, AGENTS.md, skills, context, agents, specs, `testing/guide.md`, matrix and `deploy/*.md`, then subtracted each verb's own flags: zero lines left. The `search.md` synopsis (lines 16-19) gives `index` only `--global`/`--rebuild`; `guide.md` § 3b-3e lines (39-40, 61-66, 80, 121-122) and `deploy/remote.md:162` (`update --from-zip … --home … --dry-run`) are all legal. Multi-line commands were read by hand (`guide.md:61-66`). Global flags `--repo --home --vendor` are accepted by every verb.
- Task 3 (IN10): the only `temp dir` hits are README:183 (`diff`: "never the system temp dir", correct), `remote.md:167,181` (GitHub's upload error, unrelated), and `/tmp` bans in README:292 and AGENTS.md. No doc claims a system-temp unpack or signers file. `delivery.md` § The model matches the real run.
- Task 4: the old names appear only in dated plans, handoffs and `guide.md:73,86` (the finding above). `attributionHint`, `readSettings` and `pathHint` appear in no living doc.
- Task 5: AGENTS.md § Layout (`cmd/sh-agent, internal/harness   the binary and its tests`) and README:279 stay correct with `settings.go` added; § Key files lists three files by name and needs no row. The `toolMap` rule is unaffected (`derive.go` not in the diff). `go test ./cmd/sh-agent -run TestHarnessMap` passed (ok, 0.003s). The map text itself was not in the diff.
- Task 6: `restatement-grep.sh` output was too large to read in full (46 KB, only the head of the first key word). I used the targeted greps above instead; the script's results beyond that head are not reached.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- Full output of `skills/review/scripts/restatement-grep.sh` for `takesValue`, `system temp`, `--check --force`, `settings.json`, `signers` (output too large; replaced by targeted greps; `settings.json` and `--check --force` restatements not separately traced).
- The full test suite (`go test -race ./...`, `gofmt`, `go vet`): only `TestHarnessMap` was run.
- The Go code of the refactor beyond `main.go` `run`/`updateFromZip`, `zip.go`, `pathguard.go`, the `release.go` signers hunk; the other changed `internal/harness` files, `settings.go` and the new tests were not read.
- README map labels and `delivery.md` § The model beyond line 32; `deploy/gcp.md`.
- Whether the `status` `binary` line rule in AGENTS.md still holds after the refactor (not run).

## Outcome
- [low] guide § 3b still names `takesValue` — fixed (`.agents/docs/testing/guide.md:73,86`; the three test comments too)
- [low] `search.md` exit statuses lack 2 — fixed (`.agents/docs/specs/sh-agent/search.md:29`)
- [low] `delivery.md` and README: the zip's directory must be writable; the dry run unpacks beside the zip — fixed (`delivery.md` § The model, `README.md` § Install from a zip)

## Round 2

Range reviewed: same working tree, after the round 1 fixes and the code/security fixes (every verb without positional arguments refuses a stray argument, exit 2).

### Verdict: APPROVE WITH NITS

All three round 1 fixes hold. One low gap remains: the matrix row does not cover the stray-argument refusal. No living doc prints a command a verb now refuses.

### Findings
- **[low / confidence: high]** `.agents/docs/testing/sh-agent-matrix.md:42` — the row names only flags (`list --check`, `index --kind bug`, `version --dry-run`) and the message `flag provided but not defined: -check`. The new stray-argument refusal (`list foo`, `status foo`, `version x`) is not in the matrix; its message is different (`list takes no arguments: ["foo"]`, `main.go:220`) and it has no usage block. The README (`:199`) and `search.md:29` document it, and `main_test.go:1812,1860` test it, but the per-place matrix the guide uses for manual runs lacks it. Fix: add `list foo`, `version x` to the row, with `<verb> takes no arguments: ["foo"]`, exit 2.

### Evidence
- Stray-argument grep (your command, plus a wider one over `sh-agent|cmd/sh-agent <verb>` followed by a non-flag word, over README, AGENTS.md, specs, guide, matrix, `deploy/`, skills, context, agents): zero hits where a verb gets a stray word. Remaining hits are the verb alone, a trailing comment, or prose (`README.md:78,153-155,173,205-207,252`, `deploy/local.md:19-85`, `remote.md:15-74`, `skills/distill/SKILL.md:135`, `AGENTS.md:28`, `guide.md:23`); `README.md:195` is `memory add --type … --name … --description …`, legal; `search.md:19` is `index [--global] [--rebuild]`, legal.
- Matrix row vs measured behaviour: I built the tree to `.agents/scratch/sh-agent-r2` and ran it with `--home .agents/scratch/home-r2`. `list --check` printed `flag provided but not defined: -check`, then `Usage of list:`, exit 2; `index --kind bug` and `version --dry-run` the same with their verbs' usage, exit 2. The row matches. `list foo`, `status foo`, `version x` printed `<verb> takes no arguments: ["foo"]`, exit 2, with no usage block (the row's gap). `list -- foo --home …` kept what followed `--` positional and was refused (the message listed `--home` as an argument only because my own command put it after the `--`).
- Rule 6: `delivery.md:32` (§ The model) states the claim once and carries the new parenthetical (dry-run unpack beside the zip, writable zip folder). `README.md:176`'s new text (where the unpack goes, the zip's folder must be writable, why) describes what the COMMAND does and does not say which place gets what from a tag: allowed. The sentences at the end of the paragraph ("Moving a zip place forward means repeating the download…") are not in this diff (`git show HEAD:README.md` has them), so they are an older restatement, not a Phase 8 finding. `REQUIREMENTS.md:43` says "a temporary directory" and stays true; it is not in the diff.
- `restatement-grep.sh -i --` (output filtered to living files): `takesValue` only in `main.go:936-954` (a local variable, correct); `beside the zip` in `delivery.md:32`, `README.md:176`, `REQUIREMENTS.md:43` (the signature pair, a different meaning), and code comments that agree; `takes no arguments` only `main.go:220`, no doc quotes it; `exit 2`: `deploy/local.md:29,39`, `remote.md:75`, `README.md:185,194,195,199`, `search.md:29`, `dream.md:26,69,70`, `memory-loop.md:38` all agree with the code (`update --check` exit 2 for "behind" is stated beside the new usage-error sentence at `README.md:199`).
- `--` handling (what follows stays positional) is not described in README, search.md or guide: not a finding, since it is an edge of the stray-argument rule that README:199 already states.

### Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

### Not examined
- The full text of every `restatement-grep.sh` hit for `exit 2` in dated plans, handoffs, changes and reviews (filtered out: append-only trail, not living docs).
- The full test suite; only the built binary's behaviour on the cases above was run.
- Skills and agents beyond the grep hits; the Go code.

## Outcome (round 2)
- [low] matrix row covers only flag refusals — fixed (`.agents/docs/testing/sh-agent-matrix.md`: a second row for the stray-argument refusals, with the measured message and exit 2)
