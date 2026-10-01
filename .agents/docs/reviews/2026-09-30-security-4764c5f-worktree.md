---
owner: security
date: 2026-09-30
---
# Security review — 4764c5f-worktree   (2026-09-30)

Range reviewed: 4764c5f..working tree (`git diff 4764c5f`, 16 files, +839/-94; the untracked `.agents/docs/testing/2026-09-30-verify-audit-fix-phase-1.md` is not in the range)

## Status: COMPLETE

## Verdict: FINDINGS

7 findings: 1 medium, 6 low.

The worst problem is a redaction regression in `internal/harness/report.go`. The new value-length check on the "secret assignment" needle lets through secrets that the needle at 4764c5f caught: a quoted value with spaces in it (a passphrase), and short passwords. `harness.Redact` is the one redaction shared by both paths out of the machine:
- `sh-agent report` pushes to the GitHub feedback repository;
- `sh-agent dream` sends project text to a model CLI, which `--vendor` may point at the other vendor.

A line that fails the check now leaves whole. I measured both states on the same inputs (finding 1).

The other six are low:
- the IN3 symlink promise does not hold on the failed-write path, on uninstall, or for a parent-directory link;
- the new transcript filter still sends `<task-notification>` output to the dreamer;
- the prompt fence escape leaves `>>>` behind for some runs of `>`;
- the nightly `git status` can take a project's index lock;
- `indent()` does not treat a bare CR as a line break.

No argument injection, live-database access or tool-argument leak was found on the paths the dispatch named.

## Findings
- **[medium / confidence: high]** `internal/harness/report.go:71-74` — the "secret assignment" needle and its `accept` func.
  - **Weakness.** The value is captured as `(\S*)`, so a quoted value stops at its first space. The length rule (8+ characters, or 6+ with a digit) then judges only the first word. When the rule rejects it, the whole line is kept.
  - **Measured.** A scratch probe (`.agents/scratch/sec-review/probe/main.go`, calling `harness.Redact` on the working tree) KEPT all four of these lines:
    - `export DB_PASSWORD="correct horse battery staple"`
    - `PASSWORD='my dog rex 2019'`
    - `password: letmein`
    - `api_key: "AbCdEf"`
  - **Before this range.** The needle exactly as it stood at 4764c5f, replayed in `.agents/scratch/sec-review/old/main.go`, matched all four lines (it printed `true` for each). So each of them was redacted before this change.
  - **Controls.** `password=hunter2` is still redacted. `## Secrets:` and `curl -b cookies.txt` are kept, as intended.
  - **Why no test fails.** Every must-redact value in `TestReportRedactShapes` is `FAKEVALUE`: 9 characters, no space. The new near-miss rows (`report_test.go:166-168`) pin only the relief side.
  - **Where it leaks.** Both egress paths above: report and dream.
  - **Remediation (no dependency).**
    1. When the value opens with a quote, capture up to the closing quote (`"([^"]*)"|'([^']*)'|(\S*)`) and judge the whole quoted value, spaces included.
    2. Add must-redact rows for a spaced passphrase and for a quoted value whose first word is short.
    3. Put one choice to the operator: whether a 6-7 letter unquoted value (`password: letmein`) is redacted, which costs the `token: short` near-miss. Record the answer in the needle comment and in the test.
  - **Related, older than this range.** A YAML block scalar (`secret: |`, value on the next line) leaks under both versions, because redaction is line-based. Name that limit in the same comment.
  - **The cookie change is fine.** It only stops matching `-b` followed by something with no `=` (a branch name, a cookie file, a number). A `name=value` pair after `-b`/`--cookie` is still redacted (probe control). The attached form (`-bsid=x`) was missed before this change too.
- **[low / confidence: high]** `internal/harness/apply.go:37-56` (`saveOnError`) and `apply.go:320-333` (`Uninstall`) — IN3's comment says a symlinked destination is "kept and left out of the manifest, so uninstall never removes it". That holds only after a successful install.
  - On a failed write, the union adds back every old manifest entry that `os.Lstat` finds on disk. `Lstat` succeeds on a link, so the operator's link is listed again.
  - `Uninstall` removes every listed path after only an existence `Lstat`.
  - Result: a failed install followed by `sh-agent uninstall` unlinks the operator's link. So does an uninstall run against any manifest written before this change. The link's target is untouched.
  - The union is no worse than the old manifest, which listed the same entries. It just does not deliver the new guarantee.
  - **Remediation:** skip symlinks in `Uninstall` with the same mode check the resync loop now uses (`apply.go:132-135`). That one line covers both paths. The test: a failed write, then uninstall, then the link is still there.
  - **Outside the range, same function:** `Uninstall` also removes a listed file the operator has since overwritten with their own content. It checks neither `IsOurs` nor the link.
- **[low / confidence: high]** `internal/harness/apply.go:62-66` and the printed line "sh-agent never writes through a link" — the check covers the final path component only.
  - When a parent directory is a link (`~/.claude`, `~/.claude/skills`, a skill directory), these all act inside the link's target, which can be outside the home: `os.ReadFile`, `WriteAtomic` (its temp file and rename), the `os.Chmod` in the unchanged branch, and the stale `os.Remove`.
  - `pruneEmptyDirs` then calls `os.Remove` on a directory link whose target has just emptied. Go's `os.Remove` tries unlink before rmdir, so the operator's link is deleted.
  - No privilege boundary is crossed: install runs as the home's owner, and neither README nor the deploy guides describe a root or other-user install (a grep for sudo with install/--home found none). This is not exploitable by another user.
  - The defect is the printed claim, which a reader takes as a guarantee.
  - **Remediation (operator choice):**
    - *Recommended, simplest:* reword the line to "a linked file is never written".
    - *Alternative:* confine every operation with `os.OpenRoot(home)` (stdlib since Go 1.24). This would refuse a dotfiles link that points outside the home, which is a behaviour change.
- **[low / confidence: medium]** `internal/dream/collect.go:425-431` (`contentText`) — the new filter drops only user string content that starts with `<local-command` or `<command-`.
  - I counted this machine's 413 Claude transcripts with a read-only Python pass over `~/.claude/projects/*/*.jsonl` that printed tag names and counts, never content.
  - 474 user string contents are `<task-notification>` records. 377 of them carry a `<result>`: a subagent's or background command's output (median 870 characters, largest 9,740).
  - A further 527 `isMeta` string records and 306 `isMeta` text parts are injected text, such as expanded skill bodies.
  - All of these still reach the dreamer as prose. That contradicts the rule this change carries out ("command output, not prose"; a tool step's output "the Claude Code reader never sends"). Line-based `Redact` is the only net left over them.
  - **Remediation:** drop records with `isMeta: true`, and user string content that opens with a harness tag (`^<[a-z-]+>`). A plain-prose allow-list is shorter than a deny-list that grows one prefix per round, which is the BUG-110 lesson.
- **[low / confidence: high]** `internal/dream/prompt.go:43-44` — the fence escape `strings.ReplaceAll(text, ">>>", "> > >")` runs once.
  - A run of 5 or 8 `>` characters still holds `>>>` afterwards: `">>>>>"` becomes `"> > >>>"`. Verified with the same call in `.agents/scratch/sec-review/fence/main.go`.
  - The opening `<<<data <kind>: <label>` line is not escaped at all.
  - So transcript text can still write what reads as its block's end, followed by a forged block labelled with another kind, for example memory.
  - It is low because the fence is a soft, model-level boundary, and `/retro` already treats dream output as data to be re-derived.
  - **Remediation:**
    - *Robust:* per-run fences carrying an unpredictable nonce from `crypto/rand.Text()` (stdlib, Go 1.24; consuming it as designed is not hand-rolling).
    - *Minimum:* repeat the replace until no `>>>` remains, and escape `<<<data` the same way.
    - Add a test row with a run of five `>`.
- **[low / confidence: medium]** `internal/dream/collect.go:209,226` (`trailDates`) — every indexed project now gets an unattended nightly `git log` and `git status`.
  - **Argument injection: none.** `root` is its own argv element after `-C`, and git takes it verbatim as the path, leading dashes included. `docs` follows `--`, so pathspec magic can at most widen the file set. No shell is involved.
  - **What remains:**
    - (a) `git status` without `--no-optional-locks` may refresh and rewrite `.git/index` under `index.lock`. A concurrent session's git command in that checkout can then fail with "index.lock exists". Git documents the flag for exactly this kind of background reader.
    - (b) It runs repository-local config such as `core.fsmonitor`, at night. That is the same exposure `sh-agent search` already has interactively (`internal/search/reconcile.go:120`), so it opens no new trust boundary.
  - **Remediation:** `git --no-optional-locks -c core.fsmonitor=false -C root …`. Also add `-z`, as `reconcile.go` does. Without it a quoted porcelain path (a space, non-ASCII) misses the `modified` set, and a modified record is dated by its old commit. That last point is correctness, not security.
- **[low / confidence: low]** `internal/dream/run.go:457-459` (`indent`, called at `run.go:444`) — `indent` breaks lines on `\n` only.
  - `line()` folds `\r` and every Unicode space. `o.clean()` (Redact, then the home path) keeps a bare CR.
  - So a proposal text that uses bare CRs can put `### N.` and `**Status:** queued` at column 0 for any reader that treats CR as a line end: CommonMark does, grep does not.
  - The forged-section test (`run_test.go:496`) covers `\n` only.
  - Confidence is low because I did not measure whether `/retro`'s reader, an LLM using its file tool, splits on a bare CR. The skill's own "text as data, re-derive each" rule mitigates it.
  - **Remediation:** normalise `\r\n` and `\r` to `\n` before indenting, using the replacer `line()` already has, and add a CR row to that test.

## Checked
- **Secrets:** the range adds no credential. Test values are `FAKEVALUE`, `hunter2` and `ghs_FAKETOKEN` shapes. The new agy scanner error wraps `bufio.ErrTooLong` with no stream content, and it goes only to the dream log. `ee.Msg` and the memory refusal now print to `out`; neither carries a secret.
- **Injection:**
  - *Git argv:* `trailDates` is safe (finding 6).
  - *Protobuf walker (`internal/search/dream.go:356-370`):* it now tries a nested message first. A JSON argument string opens with `{` or `[` (0x7b/0x5b, wire type 3), which never parses as a message. So it still falls through to `prose()` and is dropped whole.
  - *Crafting gains nothing:* a leaf that a nested parse yields is a substring of text the step already carried as prose.
  - *Residual, as D30 recorded:* a tool whose arguments are a non-JSON proto string would pass as prose. 032ece1 measured 0 such leaves among 971.
  - *Walk cost:* bounded at O(32·n).
  - *`openCopy` (`dream.go:252-279`):* it opens only `copyDB`'s copy in `tmp`, with `mode=ro`.
  - *Retry:* it removes only `tmp/<name>{,-wal,-shm}`. `name` is either fixed or an id already refused for `/`, `\`, `.` and `..`. The live files are only read (`os.Open`) and never opened by SQLite, so the retry cannot touch the live database.
  - *`parse()` ignoring unknown keys:* this removes no guard that mattered. An injected answer could already match the exact keys, and every field used is still folded (`line()`) or validated (`skillName`).
- **Authz/authn:** not touched by the range.
- **Crypto:** nothing added or changed; no hand-rolled primitive.
- **Dependencies:** `go.mod` and `go.sum` are not in the range. `govulncheck ./...` on go1.27.1 reported "No vulnerabilities found. Your code is affected by 0 vulnerabilities." (exit 0; output in `.agents/scratch/sec-review/govulncheck.txt`). There is no `package.json` (no `web/` or `ui/`), so `npm audit` does not apply.
- **GCP surface:** none in the range.
- **Also read:**
  - `cmd/sh-agent/main.go`: `update --check` now reports a dirty tree before its read-only return, which is the safe direction.
  - `internal/search/search.go`: hits are now deduplicated by chunk id. No security relevance.

## Not examined
- **Test bodies**, beyond the grep hits above: `harness_test.go` (+199), `dream_test.go` (+157), `run_test.go`, `collect_test.go`, `main_test.go`, `memory_test.go`, `search_test.go`. No test was run here; the verifier owns that.
- **`internal/dream/run.go`:** `write`, `unique` and `keepFailure`, and the `skillName` pattern that guards the draft path (unchanged code, not read).
- **`report.go`:** the `Collect` file list, including whether `dream-failed/` is collected. Also the open flags of `copyFileImpl`.
- **Unmeasured claims:**
  - The index rewrite by `git status` is taken from git's documentation, not reproduced.
  - How the `/retro` reader presents a bare CR (finding 7).
- **Out of range:** the untracked `.agents/docs/testing/2026-09-30-verify-audit-fix-phase-1.md`, and another agent's `.agents/docs/reviews/2026-09-30-review-4764c5f-worktree-code.md`.
- **Callers of `harness.Redact`** other than `scrub` were not traced.

## Outcome

Answered in the next round (see `2026-09-30-security-4764c5f-worktree-r2.md`); the fixes and their tests are in the phase commit. Every finding of this round was fixed, except the deploy-guide sentence, which the documentation pass of the phase fixes, and the nested-first prose-loss concern, rebutted by the corpus measurement (971 leaves and 0 argument keys before and after, `testing/2026-09-30-verify-audit-fix-phase-1.md`).
