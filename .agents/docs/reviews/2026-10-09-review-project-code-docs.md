---
owner: reviewer
date: 2026-10-09
---
# Review — project-code-docs   (2026-10-09)

Range reviewed: the whole repository at HEAD (c267990), no diff.

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The worst problem is that the dream job's 10-minute CLI timeout does not actually bound the run. `internal/dream/run.go:197-203` kills the direct child but sets no `WaitDelay`. If a grandchild of `claude`/`agy` still holds the output pipes, `cmd.Run` blocks with no limit while `dream.lock` is held, and every later nightly run exits 0 with "already running". The second problem is in the delivery claim's one home (AGENTS.md rule 6): `delivery.md:77` says the newest tag is v2.4.1 and leaves `metrics` out of the shipped verbs, while the repository is at v2.6.2 and `metrics` shipped in v2.6.0. The rest is low: usage errors exit 1 in some verbs although the docs say "usage errors exit 2", a few stale statements in README, AGENTS.md and a skill reference, and one stale comment. This was a whole-repo pass on a 35-call budget. Large parts were not read (see Not examined), so "no finding" outside the ranges below means "not checked", not "clean".

## Findings

- **[medium / confidence: medium]** `internal/dream/run.go:197-203`. `call` bounds the CLI with `exec.CommandContext(ctx…)` and `cliTimeout` (10 min, `run.go:25`). It sets no `cmd.WaitDelay`, and the child runs in sh-agent's own process group. Stdout and stderr are `bytes.Buffer`s, so Go copies them through pipes, and `Run` waits for those copies to reach EOF after the kill. If `claude` or `agy` started a grandchild that inherited the pipes (a language server, an MCP server, a helper), killing the direct child leaves them open and `Run` never returns. `realRun` then holds `dream.lock` (`run.go:54-62`) for as long as that grandchild lives, and each nightly run after it prints "already running" and exits 0 (`run.go:58-60`). Nothing fails, so the operator sees no error. This breaks the spec's row `.agents/docs/specs/sh-agent/dream.md:82` ("times out → the project fails … the others run"). Fix: set `cmd.WaitDelay` (for example 30 s; go.mod is go 1.24.4, and the field exists since Go 1.20), and preferably `SysProcAttr{Setpgid: true}` with `cmd.Cancel` killing the group. Test: a fake CLI on PATH that backgrounds `sleep 600` holding stdout and then exits or sleeps, with a shortened timeout; assert that `call` returns. (Confidence is medium because I did not run a reproduction. The mechanism is the documented `os/exec` behaviour.)
- **[medium / confidence: high]** docs, `.agents/docs/specs/sh-agent/delivery.md:77`. It says "Built and shipped (newest tag v2.4.1, 2026-10-05)" and lists the shipped verbs without `metrics`. `git tag` lists v2.5.0, v2.6.0, v2.6.1 and v2.6.2, and the `metrics` verb (commit 5ac69e0) is first contained in v2.6.0. Under AGENTS.md rule 6 this file is the one place that says what a tag delivers, and the README, deploy guides and ADRs point here, so every reader who follows a pointer gets the wrong current release and verb set. Fix: bring the paragraph to v2.6.2, add `metrics` (ingest/show/report/export) and the v2.5.0–v2.6.2 entries. Re-check: `git grep -n 'newest tag' -- .agents/docs/specs README.md AGENTS.md .agents/docs/deploy`.
- **[low / confidence: high]** usage-error exit codes, shared vocabulary (§ 7). `delivery.md:77` says "the verb table with one FlagSet per verb (usage errors exit 2)". README.md:189 and `dream.md:79` say every dream usage error exits 2 with its reason on stdout. These verbs instead return a plain error, which exits 1 with the message on stderr:
  - `dream --vendor bogus` (`cmd/sh-agent/main.go:360-363`; every other dream usage error is `ExitError{2,…}`, `internal/dream/dream.go:68-81`).
  - `search` and `index` validation: a bad `--limit`, a bad `--kind`, no mode, or mixed modes (`internal/search/search.go:81-118`, returned through `main.go:603-613`).
  - `adopt` with the wrong number of arguments (`main.go:756-757`).

  Only a flag the FlagSet does not know (`main.go:229-236`), stray arguments (`:237-242`), `memory` (`:853-855`) and `metrics` exit 2. Fix: either map these to `exitCode(2)` with the reason on `c.out`, or narrow the docs to say which verbs exit 2.
- **[low / confidence: high]** docs, `README.md:177`. The `sh-agent status` row says status shows "the running binary's build". The code reads the installed file `~/.local/bin/sh-agent` with `buildinfo.ReadFile` (`cmd/sh-agent/main.go:930-933`), and its own comment says "The line describes the installed file, not this process".
- **[low / confidence: high]** docs, layout. `AGENTS.md` § Layout (the `internal/...` lines) does not list `internal/trail`. That package is imported by `internal/dream/collect.go`, `internal/memory/memory.go`, `internal/metrics/events.go` and `internal/search/{reconcile,search}.go`, and specs `search.md:42` and `metrics.md:275` and testing `guide.md:145-148` name it. README § Layout (`README.md`, the `cmd/sh-agent` line in the Layout block) leaves out both `internal/metrics` and `internal/trail`, and the top-level `embed.go`.
- **[low / confidence: high]** stale instruction, `skills/harness-audit/references/areas.md:36`. The baseline step says to check "the newest tag's signature". Since ADR-053 no key signs anything (AGENTS.md § Release, `delivery.md`), and from v2.4.1 on there is no signature to check. Fix: say "the newest release's SHA256SUMS and assets".
- **[low / confidence: high]** stale comment, `cmd/sh-agent/main.go:414-415`. It says "runMetrics serves `metrics ingest` and `metrics export`", but the function also serves `show` and `report` (`:453-505`).
- **[low / confidence: low]** possibly unreachable branches, `internal/search/search.go:97-104`. `validate`'s `index` branch (words, `--refs`, `--open`, `--kind` refused) and the `o.Rebuild` refusal for `search` cannot be reached from the CLI. The `index` FlagSet does not register `--kind`, `--refs` or `--open` (`main.go:182-185`), `index` is not positional (`:237`), and `search` does not register `--rebuild` (`:175-181`). They are defensive checks in case a library caller sets these options; I did not check whether a test or another caller sets them. If none does, delete them, or keep them and say in a comment that they guard library callers only.

## Criteria

1. **Dead code: pass, with one low/low item.** I saw no commented-out code or debug output in the ranges read. `validate`'s index and rebuild branches look unreachable from the CLI (finding above). Not run: `go vet`/`staticcheck` for unused exports across the 95 files.
2. **Flow tracing: fail (one medium).** I traced every verb in `main.go`'s `verbs` table (`:145-200`) to its leaf entry: install/sync → `Reinstall`/`installFiles`; update → `UpdateFromRelease` / `update` / `updateFromZip`; status; diff; uninstall → `dream.RemoveTimer` + `harness.Uninstall`; report → `SendReport`; search/index → `search.Search`/`Index`; memory; dream → timer or `dream.Run`; metrics ingest/show/report/export; list; render; adopt; version.
   - Release update: download into a `MkdirTemp` directory under the config dir, removed on every path (`release.go:120-124`); the tag must match vX.Y.Z and be newer (`:88-104`); the replacement is atomic (`apply.go:356-390`).
   - Report: the throwaway clone is removed on every path (`report.go:888-895`); `RemoveKeptClone` keeps the clone whenever a check errors (`:683-744`).
   - The dream child process is the one gap: `WaitDelay` is missing (finding).
   - Context: there is no `context.Context` on the git or `gh` calls in report, update or release. They are interactive and can be interrupted, so I did not file it as a finding. The dream CLI call does have a timeout.
3. **Standards conformance: pass, with the low exit-code mismatch.** § 7 vocabulary, checked identical:
   - Verbs and flags: every verb that README, AGENTS.md, context, agents, skills, specs, the testing guide and deploy invoke exists in `verbs`, and every flag that appears in a `sh-agent …` invocation is registered (one enumeration grep over the living docs).
   - Search kinds: `main.go:62` and `internal/search/parse.go:250`.
   - Verifier statuses PASS/FAIL/BLOCKED/IN PROGRESS: `internal/metrics/events.go:22`, `internal/metrics/numbers.go:85`, `agents/verifier.md:12,49` and `metrics.md:211`.
   - Finding severity format: `events.go:23`, `agents/reviewer.md:76`, `agents/security.md:63` and `metrics.md:212`.
   - Memory caps 200 lines and 25 KB: `memory.go:28-29` and README.
   - Mismatch: the usage-error exit status (finding).

   § 8:
   - The embedded artifact is pinned both ways to git (`embed_test.go:54-93`, `TestFilesMatchTree`).
   - Guards keep data when their check errors: `RemoveKeptClone`; `update`'s dirty check, which also catches StatusUnknown and sits above the `--check`/`--dry-run` returns (`main.go:1051-1063`).
   - Dry runs: `update --dry-run` is `--check` (`main.go:691-692`); uninstall, the timer and install dry runs print and do not act. `update --from-zip --dry-run` does unpack beside the zip, which is documented (IN10, `main.go:1312-1319`) and removed afterwards.
4. **Directory layout: pass for the code, fail (low) for its description.** `cmd/sh-agent` plus `internal/*`, with tests beside the code. AGENTS.md and README Layout leave out `internal/trail` (and README `internal/metrics`).
5. **Simplicity (DIY-except-security): pass.** Stdlib plus `modernc.org/sqlite` only (`go.mod`). `crypto/sha256` is used as designed for release and zip checksums (`release.go:205`); there is no hand-rolled primitive or protocol. Since ADR-053 (no signing) the checksums are integrity-only, as the docs say.
6. **Decoupling: pass on what was read.** `internal/trail` is the shared seam for git-derived dates, so neither dream nor metrics imports the other (`trail.go` doc comment). The verb table gives each verb its own FlagSet. Not verified: that `internal/harness` imports no SQLite (AGENTS.md says a test pins this; I did not open the test).
7. **Frontend / TypeScript: not applicable.** There is no `web/` or `ui/`, and I made no frontend findings.

## Docs lens

Restatement grep run: `restatement-grep.sh -i -- 'newest tag' 'running binary.s build' 'usage error is exit 2'`. Living hits: `delivery.md:77`, `README.md:189`, `skills/harness-audit/references/areas.md:36`; the other hit was in a plans record and is excluded.

Classification of living docs (`--project`):

| Doc | Class | Proof |
|---|---|---|
| `README.md` | outdated | `:177` (status reads the installed file), Layout block (no `internal/metrics`, `internal/trail`), `:189` exit-2 claim vs `main.go:360-363` |
| `AGENTS.md` | outdated | § Layout lacks `internal/trail` |
| `.agents/docs/specs/sh-agent/delivery.md` | outdated | `:77` newest tag v2.4.1 vs v2.6.2; `metrics` missing; "usage errors exit 2" |
| `.agents/docs/specs/sh-agent/dream.md` | outdated (minor) | `:79` usage error exits 2 vs `main.go:360-363` |
| `.agents/docs/specs/sh-agent/metrics.md` | accurate on the rows checked | `:61-68`, `:283-294` match `main.go:416-599` |
| `.agents/docs/specs/sh-agent/memory-loop.md` | accurate on the row checked | `:38` exit 2 matches `main.go:853-855`; rest unverified |
| `.agents/docs/specs/sh-agent/search.md` | unverified | only `:42` checked (`internal/trail.Command` exists) |
| `.agents/docs/testing/guide.md` | accurate on the lines checked | `:145-148` names `internal/trail`; test commands carry `-timeout 30m` |
| `.agents/docs/deploy/remote.md` | accurate on the lines checked | `:104,125` embed pin command exists; the v2.4.0 bridge block is labelled historical |
| `skills/harness-audit/references/areas.md` | outdated | `:36` "newest tag's signature" |
| `agents/verifier.md`, `agents/reviewer.md`, `agents/security.md` | accurate on the vocabulary lines | `verifier.md:12,49`; `reviewer.md:76`; `security.md:63` |
| `REQUIREMENTS.md`, `ROADMAP.md`, `TODO.md`, `PRIORITY_TODO.md`, `context/*.md`, other `agents/*.md`, `deploy/local.md`, `deploy/gcp.md` | unverified | not read |

`README.md` § How it works (the map): the verbs it names (`install`, `update`, `dream --all`, `report`) all exist, and `TestHarnessMap` pins this (AGENTS.md rule 5; I did not run it). All verbs and flags the binary accepts are documented in the `usage` const (`main.go:46-78`) and README § Daily commands. `--feedback-repo`, `--force`, `--all-labels` and `--out` appear only in those two places, which is enough for discoverability.

Re-check commands:
```
git grep -n 'newest tag' -- .agents/docs/specs README.md AGENTS.md .agents/docs/deploy
git grep -n "running binary's build" -- README.md .agents/docs/specs
git grep -n -i 'usage error' -- README.md .agents/docs/specs AGENTS.md
git grep -n 'internal/search\b' -- AGENTS.md README.md   # Layout blocks: add internal/trail (and internal/metrics in README)
git grep -n -i "tag's signature\|signed SHA256SUMS" -- skills ':!skills/google-catalog/catalog'
git grep -n 'serves `metrics ingest`' -- cmd
```

## Not examined

- **Go code not read:** `internal/harness` beyond `release.go:40-375`, `report.go:560-959` and `apply.go:300-441`. That leaves out plan, derive, frontmatter/lint, inventory, zip, repo, source, the redaction needles in `report.go:1-559` and `960-1061`, and pathguard.
- **More Go code not read:** `internal/search` beyond `search.go:70-120`; `internal/memory` beyond the caps; `internal/metrics` beyond the regex lines; `internal/dream` beyond `run.go:40-409` (dream.go, collect.go and timer.go not read); `internal/trail` beyond its doc comment.
- **Tests:** none of the `*_test.go` files except `embed_test.go`. I ran no test, vet, gofmt or govulncheck.
- **Docs:** the bodies of specs `search.md`, `memory-loop.md` and `dream.md` (only the exit rows); the testing guide beyond the grep hits; `deploy/local.md` and `deploy/gcp.md`; `context/*.md`; `REQUIREMENTS.md`, `ROADMAP.md`, `TODO.md`; the newest `changes/` entry (2026-10-09-1651-release-v2.6.2.md); `sh-agent search --refs <path> --kind bug` for each file (not run: whole-repo scope); `.agents/docs/bugs/INDEX.md`, so recorded-bug reintroduction was not checked.
- **README:** the README grep output was saved to a file and only partly read (its first 2 KB). I then read the README headings and the Daily commands table separately; the rest of README § How it works prose was not read.
