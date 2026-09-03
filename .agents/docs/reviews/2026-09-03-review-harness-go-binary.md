---
owner: reviewer
date: 2026-09-03
---

# Review — twelve scripts to one Go binary   (2026-09-03)

Scope: uncommitted working-tree diff (`git status --porcelain -uall`, 79 entries) at the
09:09 snapshot — 11 staged script deletions, 19 modified tracked files, 45 untracked
files including `cmd/harness/` and `internal/{trail,gate,hookio,install,dispatch}/`.

**Caveat on the target:** the tree moved during this review. `internal/trail/check.go`,
`internal/trail/check_test.go`, `internal/dispatch/dispatch_test.go` appeared and
`cmd/harness/main.go`, `internal/install/install.go`, `internal/install/doctor.go` were
rewritten between 09:02 and 09:08 while I read them. Findings below are re-verified
against the tree as of 09:09; anything written after that is unreviewed. An earlier
finding ("three skills call `harness check`, which does not exist") was resolved by
those concurrent edits and is withdrawn.

## Verdict: REQUEST CHANGES

The worst problem is that the one security control this change claims (`R8`: "staging or
committing a credential-shaped path is refused") does not survive contact with the git
invocations agents actually use. `git -C <dir> add -A` and `git -c k=v add -A` are
allowed through with no scan at all — verified against a built binary on a tree
containing a `.env`. Second: the proportionality rule that decides whether a reviewer
ever runs measures zero lines for untracked files, so a session that creates all-new
code — the exact shape of this change set, 3965 lines of new Go — scores 0 and skips the
reviewer row unless the file count happens to trip. Everything else is smaller: dead
exports, a documentation row that is declared but can never fire, and a 419 ms `npx`
spawn on every markdown edit inside a change whose entire stated rationale is a 2.6 ms
hook budget. The Go itself is clean — gofmt/vet/`go test -race` all pass, zero
third-party dependencies, nothing hand-rolled in a security-sensitive area.

## Findings

- **[high / confidence: high]** `internal/gate/secret.go:78` `mentionsGitStaging` —
  git's global options are not skipped, so the first non-flag token after `git` is taken
  as the subcommand. For `git -C /path add -A` that token is `/path`, the loop `break`s,
  and the gate returns `Allow()`. Verified end-to-end against a tree containing `.env`:

  ```
  git add -A                  -> {"decision":"deny", ...}
  git -C <dir> add -A         -> {}
  git -c core.x=1 add -A      -> {}
  ```

  `git -C` is the standard form for any command run outside the repo root, which is most
  agent-issued git. Fix: consume the value of the value-taking global options (`-C`,
  `-c`, `--git-dir`, `--work-tree`, `--namespace`, `--exec-path=`) before reading the
  subcommand, and test that form explicitly.

- **[high / confidence: high]** `internal/gate/secret.go:101` `stagesEverything` misses
  the two other blanket forms. Verified allowed: `git add -u` (stages every tracked
  modification) and `git commit -am wip` (same, plus commits). A tracked `.tfvars` or
  `.env` that was edited this session travels through both. Fix: match `-u`/`--update`,
  and treat a combined short-flag cluster containing `a` (`-am`, `-av`) as `-a`.

- **[medium / confidence: high]** `internal/trail/git.go:71` `ChangedLines` counts only
  `git diff --numstat` and `--cached`, both of which are blind to untracked files, while
  `DirtySource` (`:25`) deliberately counts untracked files via `-uall`. The two halves
  of the proportionality rule disagree about what "changed" means. Verified in a scratch
  repo: a single new 1-line file reports `0 lines across 1 file(s)`; in this repo the
  3965 lines of new untracked Go count as 0, and the reviewer row fired only because
  1738 lines of *staged deletions* happened to be there. A scaffolding session that
  creates four large new files gets no reviewer at all. Fix: add untracked files to the
  count, e.g. `git diff --numstat --no-index /dev/null <path>` per untracked source
  path, or `git add -N` semantics.

- **[medium / confidence: high]** `internal/trail/source.go:45` — `IsSourcePath` returns
  false for every `.md` path. In this repo the product *is* markdown: `shared/AGENTS.md`
  is the constitution, `claude/agents/*.md` are the six subagents, and 28 skills are
  markdown. Verified: a repo whose `AGENTS.md`, `agents/reviewer.md` and
  `skills/x/SKILL.md` were all rewritten reports `nothing owed: the trail covers the
  current working tree`. Today's diff changes 12 agent files and 5 skills and none of
  them counts. `.agents/harness.json` offers no per-project source-path override either.
  Fix: either count markdown outside the docs root as source, or add a `sourcePaths`
  override to `Config` and set it here.

- **[medium / confidence: high]** `internal/gate/secret.go:53` — the named-secret scan
  runs over every whitespace token of the command, not over paths, and denies. Verified:
  `git commit -m "handle .env.example properly"` is denied because the *message* contains
  `.env.example`; `git add .env.example` is denied because `.env.*` matches a file
  repos are supposed to commit. `REQUIREMENTS.md` R4 says every guard fails open and
  prices a wrong block at a full round-trip; this one fails closed on ordinary work.
  Fix: exclude `.env.example`/`.env.sample`/`.env.template` from the glob set, and skip
  tokens that follow `-m`/`--message` (or only scan tokens that exist on disk).

- **[medium / confidence: high]** `internal/gate/misc.go:51` — `Format` shells out to
  `npx --no-install prettier` for `.ts .tsx .js .jsx .mjs .cjs .css .scss .json .md
  .yaml .yml .html`. Measured on this machine: 419 ms for a no-op PostToolUse on a
  one-line `.md`, against the ~2.6 ms process budget `REQUIREMENTS.md` R2 and the
  `main.go` package comment use to justify the whole rewrite. Agents write markdown
  constantly, so this is the most frequently hit hook path. Fix: look for
  `node_modules/.bin/prettier` (walking up from the file) and skip silently when it is
  absent, instead of paying `npx` to tell you the same thing.

- **[medium / confidence: high]** `internal/trail/trail.go:27` + `trail.go:42` +
  `internal/dispatch/dispatch.go:233` — `RowDocumentation` is declared, has a dispatch
  prompt, and appears in the constitution's DoD table, but `Owed` never adds it, and
  `Class()` returns `""` for it, so `HasArtifact` can never be true and `Artifacts("")`
  returns nil. The documentation row cannot fire and could not be satisfied if it did.
  `--only documentation` errors out. Either wire it (a `docs`/`specs` class plus a
  trigger) or delete the constant and the dead `case`; shipping a contract row that is
  silently unreachable is worse than not having it.

- **[medium / confidence: medium]** `AGENTS.md:28` states `go run ./cmd/harness install
  --dry-run` "regenerate[s] in-repo without deploying". It does not: `refreshGenerated`
  takes `dryRun` and only prints (`install.go:205`), and `.github/workflows/ci.yml`
  asserts the dry run writes nothing. So the documented remedy for the CI drift check
  (`generated` job) does not exist — the only real path is a full `harness install`,
  which `AGENTS.md:40` forbids an agent from running. Add a `harness generate` (or
  `install --generate-only`) or correct the documentation.

- **[low / confidence: high]** Dead code, all with zero callers — delete:
  `internal/install/derive.go:30` `AgyToolNames` (its stated consumer, `doctor.go:210`,
  explicitly says agy exposes no way to enumerate tool names, so it will never have
  one); `internal/trail/trail.go:50` `Row.Resumability`; `internal/dispatch/cli.go:23`
  `runOutput.IsError` (set at `:97`, never read — note the headless `is_error` envelope
  is therefore ignored, which may itself be a bug worth fixing rather than deleting);
  `internal/hookio/hookio.go:28,29,32` `LastAssistantMsg`, `HookEventName`, `ToolName`.
  Also `internal/trail/source.go:104`: `sourceNames[path.Base(rel)] &&
  isDependencyManifest(...)` — every name in `isDependencyManifest` is already in
  `sourceNames`, so the first conjunct is unreachable-by-construction noise.

- **[low / confidence: high]** `internal/trail/source.go:93` — `securityPatterns` are
  substring-matched against the whole lowercased path, so `AuthorCard.tsx` (`auth`),
  `internal/parser/` (`parse`), any `validate*`, and any `policy*` file fire the paid
  security row; `deployDirs` includes `base`, so `web/src/base/config.json` counts as a
  deployment manifest. Over-firing is the safe direction but R6 says spend a model only
  where judgement is needed. Anchor on path segments rather than raw substrings.

- **[low / confidence: high]** `internal/dispatch/cli.go:70` — `cmd.Output()` discards
  stderr, so a failed dispatch reports only `claude: exit status 1` with no diagnosis in
  the table or the ledger. Capture `ExitError.Stderr` (tail it) into `Result.Detail`.

- **[low / confidence: high]** `internal/dispatch/verify.go:55` — on failure the `log`
  builder holding the actual command output is dropped on the floor; only the escalated
  agent's artifact records anything, and its prompt (`rowPrompt` default branch) does not
  contain the failure, so it must re-run the suite to see it. Write the FAIL log to
  `testing/` with `verdict: FAIL` before escalating, and name it in the prompt.

- **[low / confidence: medium]** `internal/install/doctor.go:250` — `dirDiff` walks only
  `src`, so files present in a mirrored destination but absent from the repo are never
  reported as drift, even though `install` prunes them. Doctor will say "in sync" for a
  `~/.claude/skills` containing a stale skill until the next install silently deletes
  it. Walk `dst` too.

- **[low / confidence: high]** `internal/install/install.go:49` — `fs :=
  flag.NewFlagSet(...)` shadows the imported `io/fs` package inside `Run` (the same file
  uses `fs.FileMode` at `:44`). It compiles, but shadowing a package identifier is the
  kind of thing the `golang` skill exists to prevent; `Doctor` already uses `fset`. Rename.

- **[low / confidence: medium]** `internal/gate/misc.go` is a grab-bag: the formatter
  hook, the feedback capture, the status line, and shared map helpers, three unrelated
  features in one 181-line file. It is the one place in this design where two engineers
  working adjacent features would collide. Split into `format.go`, `feedback.go`,
  `statusline.go`.

## Criteria

1. **Dead code — fail.** `AgyToolNames`, `Row.Resumability`, `runOutput.IsError`, three
   unread `hookio.Payload` fields, an unreachable conjunct in `OwesSecurity`, and the
   unreachable `RowDocumentation` branch in `rowPrompt`. No commented-out blocks and no
   orphaned files — the twelve scripts are genuinely gone from the index.
2. **Flow tracing — fail.** Two breaks: `ChangedLines` and `DirtySource` disagree about
   untracked files, and `RowDocumentation` reaches a `Class()` that returns `""`. The
   positive half traces clean: every hook in `claude/settings.json` and
   `antigravity/plugins/harness/hooks.json` now names `harness hook <event>`, every event
   name has a `case` in `runHook`, the five skills that referenced deleted scripts now
   call `harness resume`/`harness check` and both subcommands exist, and the only
   remaining `*.sh` is covered by the CI shellcheck job. Errors are wrapped with `%w` or
   deliberately swallowed at hook boundaries with a stated reason; contexts carry a
   timeout into both `cli.run` and `runShell`; goroutines in `runRows` all join on `wg`.
3. **Standards conformance — pass with nits.** `gofmt -l` empty, `go vet` clean,
   `go test -race ./...` green. Package comments on all five packages, comments explain
   *why*. Nits: the `fs` shadow in `install.Run`, and `internal/hookio` has no tests
   despite being the payload-parsing seam both vendors feed.
4. **Directory layout — pass.** `cmd/harness/` for the entrypoint, `internal/` for
   everything else, no `pkg/`, tests beside their code, one bootstrap `setup.sh` at the
   root exactly as `REQUIREMENTS.md` R2 allows. Nothing landed anywhere surprising.
5. **Simplicity (DIY-except-security) — pass.** `go.mod` declares zero external
   dependencies; every parse is stdlib. No hand-rolled auth, crypto, password hashing,
   session handling, or JWT/OIDC anywhere in the diff — the secret gate is filename-glob
   matching, not cryptography, which is the right amount of DIY. No speculative
   interfaces: `cliTool` is a struct with two instantiations, not an interface with one.
6. **Decoupling — pass with one nit.** Dependency direction is acyclic and shallow:
   `install` depends on nothing internal, `trail` on nothing internal, `gate` on
   `trail`+`hookio`, `dispatch` on `trail`, `main` on all four. No cross-package
   reach-ins and no test importing another package's unexported helpers. The one seam
   that would serialize edits is `internal/gate/misc.go`.
