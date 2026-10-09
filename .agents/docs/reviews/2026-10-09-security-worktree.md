---
owner: security
date: 2026-10-09
---
# Security review — worktree   (2026-10-09)

Range reviewed: HEAD (f7c607a)..working tree — `.agents/docs/specs/sh-agent/search.md`, `internal/search/reconcile.go`, `internal/search/search.go`, `internal/trail/trail.go`, `internal/trail/trail_test.go`, untracked `internal/search/hardening_test.go` (read whole). Blast radius followed: every `exec.Command` site in `cmd/` and `internal/`, `gitView` end to end, `readAll`/`readIndexable`, `internal/memory` git calls (named by the new `trail.Command` doc claim).

## Status: COMPLETE (round 1 and round 2)

## Verdict: FINDINGS

The diff does what BUG-165 asks for search: `git status` is gone, and every git call in `internal/search` now goes through `trail.Command`. A scratch probe confirmed that the hardened `ls-files --others` runs no fsmonitor, while a plain `ls-files` in the same fixture does. The worst problem is a claim the diff adds. The new `trail.Command` doc comment (`internal/trail/trail.go:112-113`) says "Every package that runs git in a repository builds the call here (BUG-165)". That is false: `internal/memory/memory.go:128` still runs a plain `git check-ignore` in the project root. A scratch probe confirmed that this command starts the program set in the checkout's `core.fsmonitor`. This is the same class of bug as BUG-164 and BUG-165, and it is now hidden by a doc claim saying the class is closed. Second, the new `--debug` parser in `indexStats` can be made to attribute one entry's stat data to another path. A tracked filename can contain a newline followed by a forged header. A probe reproduced the mis-attribution. Its effect is limited to which hash is skipped (search integrity, no execution).

## Findings
- **[medium / confidence: high]** `internal/memory/memory.go:128` (claim at `internal/trail/trail.go:112-113`): `exec.Command("git", "-C", root, "check-ignore", "-q", rel)` runs without `-c core.fsmonitor=false`. Reading the index makes git run the hook named by `core.fsmonitor`. Probe in `.agents/scratch/secprobe-1791536335/`: with `core.fsmonitor` set to a script, plain `git check-ignore -q .agents/memory/x.md` created the marker file. Reached by `sh-agent memory add` in the operator's project. Threat model is as for BUG-164: the checkout's config is untrusted. Remediation: build this call (and `memory.Root`'s `rev-parse` at `:63`, which starts no program, for uniformity) with `trail.Command`. Add a marker test like `TestGitViewRunsNoCleanFilter` that uses a `core.fsmonitor` script. Alternatively, narrow the doc claim to `internal/search` and `internal/trail` and file the rest under BUG-165. The diff is not where the claim should be fixed alone: the claim and the code must agree.
- **[low / confidence: high]** `internal/search/reconcile.go:110-145` (`indexStats`): the next entry's header is taken as the text after the *last* newline before the NUL. A path name is printed verbatim under `-z` and may contain `\n`. A tracked file named `"zz\n100644 <oid> 0\ta.md"` therefore yields a header for `a.md` carrying the crafted entry's ctime, mtime and size. Because it sorts later, it overwrites `a.md`'s real stats in the map; reproduced with git 2.47.3 (`ls-files -s --debug -z` output in the probe). Effect: `a.md` is compared against the wrong size and mtime. A false "unchanged" keeps the index's blob ID while the bytes differ (stale or mis-keyed search content). It needs equal size and nanosecond mtime, so it is hard to trigger, and nothing is executed. Remediation: parse by position. Git's `print_debug` always prints exactly five newline-terminated lines (ctime, mtime, dev/ino, uid/gid, size/flags) after each NUL, so consume five lines and take the rest up to the next NUL as the header. Also drop a stat whose header oid differs from `view[p].oid`. Add a test with a newline-bearing path.
- **[low / confidence: low]** `internal/search/reconcile.go:41-48`: search does not refuse a partial clone the way `trail.Log`/`trail.changed` do (`partialClone`). `trail.Command`'s own comment says `GIT_NO_LAZY_FETCH` is ignored before git 2.44. If a sparse checkout marks a `.gitignore` skip-worktree, `ls-files --exclude-standard` and `check-ignore` may read it from the index blob. In a partial clone that blob may be missing, which would trigger a lazy fetch through the repository's `core.sshCommand`. This was not reproduced. It is the same exposure the removed `git status` had, so it is not a regression. Remediation: call a shared `trail.PartialClone(root)` in `gitView`/`toplevel` and refuse, as trail does.
- **[info]** `internal/search/reconcile.go:213` stats `filepath.Join(root, p)` with `os.Lstat`, which follows a symlinked *parent* directory out of the checkout. Only size and mtime are compared. Every read still goes through `fsRoot` (`readIndexable`, `:532`), which refuses the escape. No finding; noted so that the fast path never becomes a read path.

## Checked
- Secrets: none in the diff or the new test (test identities use `example.invalid`). The secret-shaped-name filter is unchanged and still applies to untracked files (`reconcile.go`, `harness.SecretShaped` loop).
- Injection: every git argument is a constant or `root`. Paths from `ls-files -z` are read only through `os.Root`. The test fixture's filter program is written into `t.TempDir()`. The `--debug` parse ambiguity is above.
- Program execution from checkout config (the subject of BUG-165): probe confirmed that hardened `ls-files --others --exclude-standard` runs no fsmonitor. Clean filters are no longer reachable: `TestGitViewRunsNoCleanFilter` checks its own must-fail state first. Targeted tests pass: `go test -race -run 'TestGitView|TestCommands|TestDates' ./internal/search ./internal/trail` gives ok/ok. No search test pins `core.fsmonitor=false` (it relies on `trail.Command`).
- Authz/authn: not applicable (local CLI; `guardProject`'s symlink and check-ignore guard is kept, now hardened).
- Crypto: `blobID` is SHA-1 of the git blob framing, used as a git-compatible content ID, not a security primitive. It is unchanged and is not hand-rolled crypto.
- Dependencies: no `go.mod` change in the diff. `govulncheck` is not installed on this box (`command not found`, absent from `~/go/bin`), so it was not run. npm: no ecosystem in scope.
- GCP surface: none touched.

## Not examined
- `govulncheck` (tool unavailable).
- The full `go test -race ./...` suite; only the targeted tests listed above were run.
- `internal/harness/repo.go` git calls (operate on the harness's own trusted checkout; outside this diff's blast radius) and `internal/dream` beyond its use of `memory.Store`.
- Behaviour on git older than 2.44 (partial-clone finding not reproduced).
- Probe directory `.agents/scratch/secprobe-1791536335/` is left in gitignored scratch.

## Round 2

Status: COMPLETE. Scope: the fixes in the working tree after round 1. These are `internal/memory/memory.go` and its test, `internal/trail/trail.go:106-120`, `internal/search/reconcile.go` (`indexStats`, the marker existence loop) and `internal/search/hardening_test.go` (all five tests), plus the README row. I re-read every declaration the fixes touched and grepped every reference to it.

### Verdict round 2: FINDINGS (low only)

Both round-1 items the fixes target (medium 1, low 2) are fixed and pass a mutation check. The fixes add no security defect. The only new item is a low one in a new test: it discards Add's error, so it could later pass vacuously. Low 3 (no partial-clone refusal in search) is left as a stated limit. That is acceptable at low severity, but it should be written down where the limit lives (see the note below).

### Per finding
- **Medium 1 (memory check-ignore ran core.fsmonitor): FIXED, confidence high.** `Root` (`:64`), `trusted` via `storeTracked` (`:308`) and `Add`'s check-ignore (`:131`) all build through `trail.Command`. A grep of `cmd/` and `internal/` finds `exec.Command("git"` only in `internal/trail/trail.go:118` and `internal/harness/repo.go` (the harness's own checkout, as the narrowed comment says). The comment at `trail.go:111-115` now matches the code. Its list (dream, metrics, search, memory) agrees with the callers: `metrics/events.go`, `search/{reconcile,search}.go`, `memory/memory.go`, and `dream/collect.go:75` via `trail.Dates`. Mutation check on a scratch copy of the tracked tree (`.agents/scratch/secmut2-*`): the baseline passed. With `Add`'s call reverted to plain `exec.Command` and `core.fsmonitor=false` removed from `trail.Command`, `TestAddRunsNoProjectHook`, `TestStoreRunsNoProjectHook` and `TestGitViewRunsNoFsmonitor` all FAIL ("... ran the ... fsmonitor hook").
- **Low 2 (indexStats forged header): FIXED, confidence high.** The call is now `ls-files -s --debug` without `-z`. A scratch probe ran with `core.quotePath=false` set by the checkout, which is the attacker's lever. `\n`, `\t` and `\r` names are still C-quoted (`"zz\n100644 … \ta.md"`, `"q\r"`), and only bytes of 0x80 and above print raw. Raw bytes match the `-z` path, so the stat is attributed correctly, and a quoted name is skipped and hashed. A debug line always starts with two spaces and a header never does. A real name beginning with `"` is itself quoted, so the `HasPrefix(p, "\"")` skip cannot drop a real unquoted name. `TestIndexStatsIgnoresForgedNames` names its forged file `evil\n…`, which sorts after `a.md`, so it would overwrite under the old parser. It asserts `a.md`'s size is exactly 5, so a missing entry also fails it.
- **Low 3 (partial clone in search): NOT CHANGED, by choice.** Still low / low confidence. `trail.Command`'s comment says outright that an older git ignores `GIT_NO_LAZY_FETCH`, yet `gitView` and `guardProject` never call `partialClone`. The search spec's sentence "None of them starts a program …" (`.agents/docs/specs/sh-agent/search.md`, Freshness) does not name this limit. A stated limit should be stated there or in BUG-165.

### New in round 2
- **[low / confidence: high]** `internal/memory/memory_test.go:422`: `TestAddRunsNoProjectHook` discards `Add`'s error. The test is red today (mutation above). But if a later change makes `Add` refuse before the check-ignore line, the hook never runs and the test passes vacuously. Remediation: assert `err == nil`, and that the memory file exists, before checking the marker.
- **[info, records]** The BUG-164 record (`.agents/docs/bugs/BUG-164-git-log-ran-a-repos-gpg-program.md:33,37`) and its `INDEX.md:116` row still name `internal/trail.command` (now `Command`). They also list the permitted project-root commands as `log`, `ls-tree`, `hash-object --no-filters`, `rev-parse`. Project roots now also run `ls-files` (`-s`, `--debug`, `--others --exclude-standard`), `check-ignore` and `config --get`. Each has a marker test now, but the prevention rule as written no longer describes the code. BUG-165 still reads `status: open` with `Fix: Open`. These are the main agent's records; there is no security impact.
- The marker loop (`reconcile.go:211-215`) now runs over every marker regardless of prefix, before the tracked loop. A marker that is a symlink still counts (fsRoot.Lstat succeeds), which is the conservative direction (more is left out). No issue.

### Checked round 2
- `go test -race -count=1` of the new and changed tests in `internal/search`, `internal/memory` and `internal/trail`: all PASS (TestGitViewRunsNoCleanFilter, RacilyCleanEditIsHashed, DeletedMarkerOutsidePrefix, IndexStatsIgnoresForgedNames, GitViewRunsNoFsmonitor, StoreTrackedArgs, StoreRunsNoProjectHook, AddRunsNoProjectHook, CommandsNeverLazyFetch).
- `gofmt -l` clean and `go vet` ok on the three packages.
- Scan for functions in the three touched files with no caller: none. `storeTrackedArgs` is gone from code, and only historical review and testing records name it. `trail.Command` has callers in four packages, so the export is used, not orphaned API.
- Secrets, crypto, dependencies, GCP: unchanged from round 1 (no `go.mod` change; `govulncheck` still unavailable).

### Not examined round 2
- The full `go test -race ./...` suite and `govulncheck` (tool absent).
- `TestGitViewDeletedMarkerOutsidePrefix` was run but not mutation-checked. It is a correctness fix from the code reviewer, not a security one.
- Git older than 2.44 (low 3).
- Scratch dirs left (gitignored): `.agents/scratch/secprobe-*`, `.agents/scratch/secmut-*`, `.agents/scratch/secmut2-*`.

## Outcome

Written by the main agent, 2026-10-09.

Round 1:
1. **Medium.** Memory's `check-ignore` starts fsmonitor: **fixed**. All memory git calls go through `trail.Command`, under `TestAddRunsNoProjectHook` (red before the fix). Recorded under BUG-165.
2. **Low.** A name holding a newline can forge `--debug` lines: **fixed**. `indexStats` reads line by line without `-z` and skips quoted names, under `TestIndexStatsIgnoresForgedNames` (red before the fix).
3. **Low.** Search does not refuse a partial clone: **left: operator, as a stated limit**. It is in `specs/sh-agent/search.md` § Freshness and in BUG-165's record. Not a regression; refusing would stop search in partial clones on every git.

Round 2:
4. **Low.** `TestAddRunsNoProjectHook` ignored `Add`'s error: **fixed**. The test now asserts that `Add` succeeded.

No round 3: round 2 reported one test low, and it is fixed.

Update, 2026-10-09 (the operator: "try to fix all in the next release"): item 3 is now **fixed**. Search refuses a partial clone only when git is older than 2.44 (`internal/trail` `LazyFetchSafe`, `internal/search/reconcile.go` `gitView`), under `TestGitViewPartialCloneNeedsNewGit` (red before the fix) and `TestVersionAtLeast`. The release gate's security review covers it.
