---
owner: security
date: 2026-09-13
---
# Security review — improvement-items-round-2   (2026-09-13)

Range reviewed: `git diff 5cdd4d6 -- internal/harness/repo.go internal/harness/harness_test.go .gitignore`
(working tree, uncommitted; tree held still). One mechanism only: round-1 finding 1,
"git inherits the environment, so a nested copy can pass `ownCheckout`".

Blobs pinned (`git hash-object` at review time):

| blob | path |
|---|---|
| `14b5a42fbdac8f87cd10c4fe60b3106b53e1f5b3` | `internal/harness/repo.go` |
| `5792e0977158b3b9d7de59fcf127fa869c7be88d` | `internal/harness/harness_test.go` |
| `14ad0bfbb2495395650c3c6b1026cf1045bb45e4` | `.gitignore` |

## Verdict: FINDINGS — the reported hole is closed, the mechanism has NOT converged

`gitEnv()` does close round-1 finding 1 exactly as written, and I verified it two
ways. At the git level: in an isolated fixture, `git -C nested rev-parse
--show-toplevel` returns the enclosing repo normally but returns `nested` when
either `GIT_DIR` **or** `GIT_WORK_TREE` is exported — both are now scrubbed. At the
Go level: I copied `internal/harness` to scratch, reverted the two
`cmd.Env = gitEnv()` lines, and `TestNestedCopyIsNotACheckout` failed on precisely
the new assertions (`Inspect = {Branch:master Commit:278f0f9 Dirty:true}, want
empty`). The test has teeth and the fix is the thing that gives it teeth. I also
checked every remaining repository-location variable in `git help git` against the
real fixture: `GIT_NAMESPACE`, `GIT_DISCOVERY_ACROSS_FILESYSTEM`, and
`core.worktree` injected through config all leave `--show-toplevel` pointing at the
enclosing repo, and `GIT_CEILING_DIRECTORIES` only makes discovery fail, which
fails closed. So the *discovery* half of git's environment surface is covered, and
the seven-name list has no gap in it.

What is not covered is the *configuration* half of the same mechanism.
`GIT_CONFIG_COUNT`/`GIT_CONFIG_KEY_<n>`/`GIT_CONFIG_VALUE_<n>` and
`GIT_CONFIG_GLOBAL`/`GIT_CONFIG_SYSTEM` are documented in git-config(1) rather than
git(1), survive the scrub, and reach `Pull`'s `git pull --ff-only` — which `update`
follows immediately with `Reinstall`, a `go build` of that tree and an exec of the
result. Verified read-only in the fixture: with `GIT_CONFIG_COUNT` exported,
`git ls-remote --get-url origin` resolves to a host of the variable's choosing
instead of the configured remote, and `git rev-parse --git-path hooks` resolves to a
directory of the variable's choosing. That is the same sentence as round 1 — "an
inherited environment variable makes git answer for something other than what the
guard checked" — landing a second time on the same mechanism, which is the point at
which the global rule says stop extending the list and ask what the mechanism is
for. **My recommendation is not an eighth, ninth and tenth prefix.** Enumerating
git's environment is unbounded (`GIT_SSH_COMMAND`, `GIT_ASKPASS`,
`GIT_EXTERNAL_DIFF`, `GIT_PAGER`, `GIT_PROXY_COMMAND`, and `PATH` itself are all
still there and all still run programs), and a deny-list that must grow with git's
manual page is a maintenance liability in a repo whose stated virtue is that adding
things is a file drop. Two shorter mechanisms are available; see finding 1.

Calibration, stated plainly so it can be rebutted rather than argued around: an
attacker who can set `GIT_CONFIG_COUNT` in the operator's environment can usually
set `PATH` too, and then nothing in this binary matters. The reason I still rate it
medium is that the round-1 finding was accepted on exactly that premise ("no
attacker-controlled input beyond the environment is needed"), that `GIT_CONFIG_GLOBAL`
*is* set ambiently by some sandboxes and test harnesses, and that the consequence
here is worse than round 1's — `core.hooksPath` is code execution during `git pull`
itself, before the build. "The environment is trusted, so `update` is only safe in an
environment you control" is a legitimate rebuttal; if it is taken, write it down next
to `Pull` so round 3 does not re-derive this.

## Findings

- **[medium / confidence: high]** `internal/harness/repo.go:169-182` — `gitEnv()`
  scrubs the seven variables that redirect git's *repository discovery* but not the
  ones that redirect its *configuration*. `GIT_CONFIG_COUNT` with
  `GIT_CONFIG_KEY_<n>`/`GIT_CONFIG_VALUE_<n>` (equivalent to `git -c`), and
  `GIT_CONFIG_GLOBAL`/`GIT_CONFIG_SYSTEM` (which replace `~/.gitconfig` and
  `/etc/gitconfig`), pass through `gitEnv()` untouched into `Pull`'s
  `git pull --ff-only` (`repo.go:120-121`) and into every `gitOut` call. Verified
  read-only in a scratch fixture on git 2.47.3: with those variables set,
  `url.<host>.insteadOf` rewrote `ls-remote --get-url origin` from
  `https://github.com/kesokaj/styrman-harness.git` to an unrelated host, and
  `core.hooksPath` moved `rev-parse --git-path hooks` to an arbitrary directory. The
  first means `update` fast-forwards the harness from content the operator did not
  choose and `Reinstall` (`repo.go:138-163`) then builds and executes it; the second
  means `git pull`'s own hook resolution points at an attacker-chosen directory, so
  execution does not even need the build. Remediation, in order of preference:
  **(a)** stop enumerating and bound the environment instead — build git's env from
  an allow-list (`PATH`, `HOME`, and the `GIT_*_NAME`/`GIT_*_EMAIL` pair if any
  command ever needs it) rather than from `os.Environ()` minus a deny-list; this is
  the same number of lines and does not grow when git adds a variable. **(b)** If the
  environment is declared trusted, delete the deny-list too and say so in a comment
  on `Pull` — a half-scrub that stops `GIT_DIR` but not `core.hooksPath` invites the
  reader to believe git is sandboxed here, which is the more dangerous state to be
  in. Do not simply append three more prefixes.

- **[low / confidence: high]** `internal/harness/repo.go:170` — `var env []string`
  means `gitEnv()` returns **nil** when every variable in the process environment
  matches a scrubbed prefix, and `os/exec` documents that a nil `Cmd.Env` makes the
  child *inherit the parent's environment*. Verified with a five-line scratch
  program: `c.Env = env` (nil) printed `GIT_DIR=[/poisoned]`, while
  `c.Env = []string{}` printed `GIT_DIR=[]`. So an invocation such as
  `env -i GIT_DIR=... sh-agent update` re-opens exactly the hole the function was
  written to close, silently and with the test still green. Reachability is narrow —
  it needs an environment containing nothing but git redirect variables — hence low,
  but the defect is certain and the fix is one word: `env := []string{}`. Worth a test
  case that calls `gitEnv()` with a stubbed environ and asserts the result is non-nil.

- **[low / confidence: high]** `internal/harness/apply.go:79-81` — the guard is
  applied inconsistently: `Install` reads the commit with a raw
  `gitOut(repo, "rev-parse", "--short", "HEAD")` that never passes through
  `ownCheckout`, then guards only the `+dirty` suffix with `Inspect(repo).Dirty`.
  For a harness-shaped plain copy nested in another repository — the exact tree this
  whole mechanism exists to reject — `Inspect` now correctly returns `{}` so `Dirty`
  is false, while the raw call still returns the *enclosing* repository's commit. The
  manifest therefore records "installed from `<outer commit>`, clean" for files that
  came from an unmanaged working tree. No authorization is bypassed (`knownCheckout`
  at `main.go:350-365` keys on `st.Repo`, a path, not on the commit), so this is
  provenance rather than access: `sh-agent status` and every future drift comparison
  attribute the deployed context to a commit that never contained it. Remediation:
  have `Install` take the commit from `Inspect(repo).Commit` so one function owns the
  question, or return early from the stamping block when `!ownCheckout(repo)`.

- **[low / confidence: medium]** `internal/harness/harness_test.go:625-641` — the new
  polluted-environment variant proves the fix for `GIT_DIR` and `GIT_WORK_TREE` only,
  and only against `Inspect` and `Behind`. The other five scrubbed names, and `Pull`
  under a polluted environment, are untested, so five sevenths of the deny-list is
  asserted by reading rather than by the suite; a later edit that drops a prefix stays
  green. The teardown also mixes `t.Setenv("GIT_DIR", "")` with a direct
  `os.Unsetenv("GIT_DIR")` — harmless today because `t.Setenv`'s cleanup restores the
  original value at test end, but the two lines together read as if one of them is
  load-bearing when neither is. Remediation: table-drive the variant over all seven
  names and include `Pull`; drop the manual unset.

- **[low / confidence: high]** `.gitignore:7-13` — round-1 finding 4's secret-shaped
  patterns landed (`.env`, `.env.*`, `*.pem`, `*.key`, `*.tfvars`, `*.tfstate`), which
  is the improvement this file needed; recording two gaps rather than re-opening the
  finding. Credential JSON is still not covered — the global rule names
  "credential JSON" and the GCP posture makes a stray `*-key.json` or
  `application_default_credentials.json` the most likely secret to appear in this
  repo's tree, yet no pattern matches it. And the unanchored `.claude/agent-memory/`
  line was removed rather than kept: `/.claude/` on line 11 is root-anchored, so a
  nested `*/.claude/agent-memory/` is untracked-visible again. Remediation: add a
  credential-JSON pattern and restore the unanchored agent-memory line; a stale ignore
  costs nothing.

- Round-1 finding 2 (`govulncheck` on go1.24.4) — **rebutted, accepted.** `TODO.md`
  carries the toolchain upgrade. Not re-run this round; it is not in this delta.

## Checked

- **Secrets** — the `.gitignore` delta adds secret-shaped patterns and removes none
  that protected a secret; see the finding for the two gaps. No credentials, tokens,
  or keys in the changed Go code. `gitEnv()` does not log or echo the environment it
  filters, and no error path in `repo.go` prints an environment value — `gitOut`'s
  error wrapping includes only argv and git's stderr (`repo.go:191`), and argv here is
  constant strings plus `repo`.
- **Injection** — the environment *is* the injection channel under review; see
  finding 1. Argument construction is unchanged and safe: every `exec.Command` in the
  delta takes an argv slice, and `repo` is `filepath.Abs`'d by `FindRepo` before it
  reaches one. No new shell, template, SQL, or path interpolation.
- **Authz/authn** — no authentication surface. The authorization analogue, "which
  tree may be built and executed", is now guarded twice: `knownCheckout`
  (`main.go:350-365`) before `update()`, and `ownCheckout` inside
  `Inspect`/`Behind`/`Pull`. The ordering is still correct and the guards are still
  independent. `ownCheckout` is now environment-hardened for discovery; it is not
  environment-hardened for configuration, and `apply.go:79` skips it entirely
  (findings 1 and 3). No client-supplied ID, no reordered check, no skipped middleware.
- **Crypto** — none added, none present, nothing hand-rolled anywhere in the module.
  Unchanged by this delta.
- **Dependencies** — `go.mod` still requires nothing beyond `go 1.24.4`; the delta
  adds one stdlib import (`errors`, in the test file) and no third-party package, so
  the stdlib-first rule holds with nothing to justify. `govulncheck` was **not** re-run
  this round: its round-1 result is a rebutted finding against the toolchain, not
  against this delta, and nothing in this delta changes the reachable call set.
  `npm audit` is not applicable — no `package.json`, no JavaScript in the repo.
  Verification actually run on the delta: `gofmt -l .` clean, `go vet ./...` clean,
  `go test -race -run 'TestNestedCopyIsNotACheckout|TestBehindAndPull|TestFindRepo'
  ./internal/harness/` passed, and the reverted-fix scratch copy failed as described.
- **GCP surface** — nothing in this delta. No IAM, no service accounts, no Cloud Run
  flags, no Terraform, no secrets-as-env-vars, no key files.

## Convergence

**Not converged.** Round 1 and round 2 are the same sentence about the same
mechanism, and the round-2 answer is a longer list rather than a smaller surface.
Under the global rule this is the stop-and-ask point. The question to put to the
operator is not "which variables next" but: does `ownCheckout` need to run git at
all? `os.Stat(filepath.Join(repo, ".git"))` answers "is this tree the top of its own
checkout" for both a `.git` directory and a `git worktree` gitfile, cannot be
redirected by any environment variable, and removes one subprocess from every
`Inspect` call — that is a fix that makes the code shorter. It leaves only `Pull`'s
own execution needing a bounded environment, which is finding 1(a) and is one
allow-list, fixed in size. If instead the answer is that the operator's environment
is trusted, then the deny-list is theatre and should go, with a comment in its place.
Either way the next change should shrink this mechanism. A third round of prefixes
should not happen.

## Not examined

- Everything outside the three named files. The rest of the `5cdd4d6..working tree`
  diff (the ~3,000 lines of skill prose, `context/`, `agents/`, `cmd/sh-agent/`,
  `derive.go`) was reviewed in round 1 and is unchanged by this delta; I did not
  re-read it and this round makes no claim about it. Round 1's own "Not examined"
  section still stands in full.
- `git pull`'s command-executing environment beyond the config channel —
  `GIT_SSH_COMMAND`, `GIT_SSH`, `GIT_ASKPASS`, `GIT_EXTERNAL_DIFF`, `GIT_PAGER`,
  `GIT_PROXY_COMMAND`, `GIT_EDITOR` — was enumerated from `git help git` but not
  individually exercised. They are named as evidence that the deny-list cannot be
  completed, not as separate findings.
- No end-to-end run of `sh-agent update` against a poisoned environment. Every
  observation above was made at the `git` level in an isolated scratch fixture, or by
  running the existing unit test against a scratch copy of the package. A working
  end-to-end demonstration would be an exploit, not a finding, and was deliberately
  not built. No file in the repository was modified by this review.
- `.agents/docs/bugs/BUG-010-batch-edit-partial-apply.md` and the new testing record
  were not read; they are trail records with no execution surface.
