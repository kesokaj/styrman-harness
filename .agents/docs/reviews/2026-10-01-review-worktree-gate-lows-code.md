---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-gate-lows-code   (2026-10-01)

Range reviewed: 8b6c6de..working tree

## Status: COMPLETE

## Verdict: REQUEST CHANGES

One finding blocks: the Makefile deploy example in `skills/bash/references/makefile-patterns.md` now passes `--project=$(PROJECT)`, but the file's own variable is `GCP_PROJECT`; `PROJECT` is never defined, so the added flag expands to an empty `--project=` and the example's `gcloud run deploy` fails (or, depending on gcloud, silently targets nothing). The Go changes are sound: `go test -race ./internal/memory ./internal/dream ./internal/harness` passes, and each of the four new tests fails when its fix is reverted (mutations below).

## Findings

- **[medium / confidence: high]** `skills/bash/references/makefile-patterns.md:203,206` — `--project=$(PROJECT)`. `grep -n PROJECT` shows the variable is `GCP_PROJECT ?=` (line 134), used at 136, 146-147 and 238-239; `PROJECT` appears only in the two new lines. Make expands an undefined variable to the empty string, so the command becomes `--project=`. Fix: `--project=$(GCP_PROJECT)` on both lines (the `ifndef GCP_PROJECT` guard at 238 is the pattern's own check; the `deploy` recipe itself has no guard like `docker-push`'s `test -n`, so also consider the same `test -n "$(GCP_PROJECT)"` line there).
- **[low / confidence: medium]** `skills/docker/SKILL.md:88` (rule 1, last sentence, unchanged line) — "When the whole build must be reproducible byte for byte, pin the Go and Node images by digest too." The diff makes the Go stage digest-pinned in the template and in the same rule, so "the Go ... too" is now a stale restatement; only the Node stage remains for that sentence. Reword to "pin the Node image by digest too" (describing-text rule).

## Answers to the brief's points

1. **`RedactJSONL` marker.** On a hit the whole line becomes `{"redacted":"[redacted…]"}` plus the line's own terminator (`marker(…, jsonl=true)`, report.go:231-237, 268). In `feedback()` the entry's `ts` was already parsed from the original line before redaction, so ordering and the since-filter are unaffected; the prompt entry is a small JSON stub, which reads fine as data in the dream prompt. The marker text says "by sh-agent report", cosmetic only. Key blocks (`blockOpen`) also emit the JSON marker, consistent.
2. **Round trip.** `SplitAfter` on `line+"\n"` gives `[line+"\n", ""]`; the empty tail is skipped and a clean line is written back as `raw` unchanged, so `TrimRight(…,"\n")` returns the original bytes. `bufio.Scanner` lines contain no `\n`, so nothing of the line's own is stripped. Correct.
3. **`readLines` Lstat refusal.** `Add`: `trusted()` → `check-ignore` → pre-check `readLines` (error ignored, falls through) → `MkdirAll`, lock → second `readLines(index)` returns the refusal, all before either `WriteAtomic`. Only the store dir and `.lock` can be created; neither the entry file nor the index is written, and `WriteAtomic` renames over the path rather than following it. `Show` returns the `*Refusal`; `cmd/sh-agent/main.go:614` maps `IsRefusal` to `exitCode(2)`. Confirmed. The Lstat→Open gap is a negligible race for a store the operator owns.
4. **`allowed-tools` lint.** The spelling `allowed-tools` is the Agent Skills spec's (repo's `reviews/2026-09-03-vendor-format-check.md:17` quotes the spec; Claude Code uses the same key). `LoadSkills` has one non-test caller, `inventory.go:132`, which loads the harness source (checkout, embedded, zip); project skills under `.agents/skills` never pass through it, so a project skill with a legitimate `allowed-tools` is not refused anywhere else. Confirmed by `grep -rn "LoadSkills("`. The key is case-sensitive (`Allowed-Tools` would pass); acceptable, the spec key is lowercase.
5. **fsmonitor hook test.** Git here is 2.47.3. The test's control uses `git status`, which does call the hook; but `Store` runs `ls-files`, a different command, so I checked the test is not vacuous: on a scratch copy with `-c core.fsmonitor=false` deleted from `storeTrackedArgs`, `TestStoreRunsNoProjectHook` fails ("Store ran the project's fsmonitor hook") alongside `TestStoreTrackedArgs`, both after a clean build. So `ls-files` does consult the hook on this git and the test detects the missing guard. Deterministic across versions: it skips when the control does not fire, so on a git that never calls the hook it reports SKIP (not PASS) and cannot give false assurance; the only risk is a git where `status` fires the hook but `ls-files` does not, in which case the test passes vacuously. Low, not raised as a finding: the argument-list test `TestStoreTrackedArgs` still pins the flag there.
6. **Docs.** Docker digest placeholder `@sha256:<digest>` with `docker buildx imagetools inspect golang:<major.minor>` is consistent with the runtime image's digest rule (finding 2 aside). data-gcp: Cloud SQL's `cloudsqlsuperuser` does carry CREATEROLE and CREATEDB; the claim holds, but the em-dash parenthetical splits the sentence ("a `cloudsqlsuperuser` — … — member with DDL at runtime"), so the noun "member" is far from its article; readable, not a finding. README: `git add context/dev-setup.md && git commit -m` and `npx skills@<version> add … ` (no `-y`) are consistent with the global git rule; no other copy of `commit -am` was found in `README.md`, `skills/`, `context/`, `agents/` (only the vendored `find-skills` keeps `npx skills add … -y`, which is verbatim and untouchable). `skills/feature/SKILL.md` issue-body-is-data sentence: no conflict found.

## Evidence

- `go test -race ./internal/memory ./internal/dream ./internal/harness`: all three `ok`.
- Mutations on a scratch copy under `.agents/scratch/mut` (removed afterwards; the working tree untouched), each built and a named test failed: guard `-c core.fsmonitor=false` removed → `TestStoreRunsNoProjectHook`, `TestStoreTrackedArgs` fail; `feedback()` back to `sc.Text()` → `TestFeedbackRedactsEscapedSecret` fails; Lstat pointed at a nonexistent path → `TestSymlinkedIndexRefused` fails (Show and Add both followed); lint key renamed `allowed-toolz` → `TestLintRefusesAllowedTools` fails. The original file was never restored by reverse edit: the whole scratch copy was deleted.

## Criteria
1. Dead code — pass: `RedactJSONL`'s one caller is `dream.feedback`; nothing orphaned or commented out.
2. Flow tracing — pass for Go (callers and error paths traced above; tests are mutation-verified); fail for the Makefile example (finding 1).
3. Standards conformance — fail, low: one stale restatement in `skills/docker/SKILL.md` (finding 2); the changed rule's other copies (`golang:` base-image lines) were grepped, `gke-basics/references/gke-app-onboarding.md:28` (`golang:1.22`) predates the diff and is already flagged in the gke-basics SKILL.
4. Directory layout — pass: tests beside their code, no new files in the scope.
5. Simplicity (DIY-except-security) — pass: stdlib only; no new dependency; no hand-rolled crypto or auth.
6. Decoupling — pass: `dream` reaches redaction through the exported `harness.RedactJSONL`, a one-function seam; no reach-ins.

## Not examined
Docs-root changes in the working tree beyond the scope (BUG-143 and `INDEX.md` edits, BUG-145, the other reviews and verifier logs, `testing/2026-10-01-verify-agy-plan-mode-probe.md`) were not read. `go vet`, `gofmt` and the size pins in `internal/harness/size_test.go` ran only as part of the three-package test run (the rest of `go test ./...` and `govulncheck` not run). Antigravity's handling of `allowed-tools` was not tested. The `TestStoreRunsNoProjectHook` result on a git other than 2.47.3 was not tried.

## Round 2

Re-read the two spots in the working tree.

- Outcome, finding 1 (Makefile `PROJECT`): fixed. `makefile-patterns.md:203,206` now use `--project=$(GCP_PROJECT)`; `grep -n PROJECT` shows no remaining `$(PROJECT)`, and the variable is defined at line 134. The optional `test -n` guard on the deploy recipe was a suggestion, not a finding, and is not required.
- Outcome, finding 2 (docker rule 1): fixed. `skills/docker/SKILL.md:88` now reads "pin the Node image by digest too".

Round verdict: APPROVE. No open findings. Not re-run: the Go tests (no Go file changed in this round).

## Round 3

Scope: new hunks in `internal/harness/{frontmatter,inventory,harness_test}.go` and `internal/memory`. `go test -race -count=1 ./internal/harness ./internal/memory`: both `ok`.

Docker note: `skills/docker/SKILL.md` is back at HEAD (`git diff HEAD --stat -- skills/docker` empty), so round-1 finding 2 is moot, and the Makefile fix (round 2) still stands (`GCP_PROJECT` used at 203/206).

### Findings
- **[low / confidence: medium]** `internal/harness/inventory.go:298-306` (`LoadDocs`) — the permission lint and `ParseStrict` run only when `Split` succeeds. `Split` (`frontmatter.go:22-34`) is exact: the first line must be `---` (a leading BOM fails it) and the closer must be a bare `---` (`--- ` with a trailing space fails it). A command/workflow file in that shape is treated as plain body and deployed verbatim, with no lint, while a reader that strips a BOM or tolerates trailing spaces would still read `allowed-tools`. Skills and agents refuse on a `Split` error, so only commands, workflows and rules have the gap. `commands/` and `rules/` are empty today and the author is the operator, so this is a hardening gap, not a live hole. Fix: when `Split` fails but the text starts with a BOM or `---`, refuse instead of falling through. I did not verify how Claude Code or Antigravity parse a BOM-prefixed command.
- **[low / confidence: high]** `internal/harness/harness_test.go` (`TestLintPermissionKeysByKind`) — the LoadDocs refusal is tested only through `commands/` and `allowed-tools`; `hooks`, `workflows/` and `rules/` are the same function and loop but have no case. Not a defect in the code (one shared loop); the test would still pass if `hooks` were dropped from the LoadDocs key list. Add one `hooks` row in `rules/`.

### Answers
- **`unreadLine` strictness.** `keyLine` is `^[A-Za-z0-9_-]+:(?:[ \t]+(.*))?$`, so a top-level line that is a `---`/`...` document marker (a multi-document file), a key with a dot or slash, a `? ` complex key, a one-line flow map `{a: b}`, `key:value` without a space, or a column-0 wrapped plain scalar is now refused. Real keys in the checkout are all plain (`name description argument-hint tools model maxTurns effort skills metadata license`), and `LoadInventory(../..)` passes, which loads every skill, agent and command. For an installer that reads only name and description, refusing is acceptable and fail-closed: these lines are exactly the shapes `entries` would skip silently, so a name-keyed lint would never see them. Cost to note: one such line in a vendored skill (which cannot be edited in place, `skills/VENDORED.md`) fails the whole install, not just that skill; the error message names the line, and the remedy is the fork rule. Not a finding.
- **`LoadDocs` and `rules/`.** The refusal sits inside the per-file loop, and `LoadInventory` calls `LoadDocs` for `commands`+`workflows` and for `rules`, so rules are covered. `hooks` has no meaning in a rule (Claude Code rule files take `paths` only) and `allowed-tools` has none either; refusing is harmless, and the message says "command or rule". `allowed-tools` is meaningful in commands.
- **Benign list.** The test's list (metadata map, license, comment, column-0 list) plus the checkout run covers every key the checkout's own files use; `compatibility` (column-0 list) is the one shape the checkout does not use, and it is exercised synthetically.
- **Memory.** Third assertion of `TestSymlinkedIndexRefused` (link still a link after the refused Add) is right: without the refusal, `WriteAtomic` replaces the link with a regular file and the Lstat assertion fails. The `readLines` change is otherwise as in round 1.
- **Mutations** (scratch copy, deleted afterwards; each compiled and a named test failed): `unreadLine` neutralised → `TestLintPermissionKeysByKind` fails on the quoted, single-quoted and spaced-colon keys; `hooks` dropped from the agent key list → fails on the agent `hooks` case.

Criteria for the round: dead code pass (`Keys` has one caller, `permissionKey`; exported but used only in-package, not a finding); flow tracing pass apart from the `Split` gap above; standards pass; layout pass; simplicity pass (stdlib, no hand-rolled auth); decoupling pass.

Round verdict: APPROVE WITH NITS (two low findings, neither blocking).

Not examined in round 3: Claude Code's and Antigravity's actual parsing of BOM or trailing-space frontmatter; the full `go test ./...`.

## Status: COMPLETE
