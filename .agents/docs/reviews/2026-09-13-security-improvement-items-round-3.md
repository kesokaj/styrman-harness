---
owner: security
date: 2026-09-13
---
# Security review — improvement-items-round-3   (2026-09-13)

Range reviewed: `git diff 5cdd4d6 -- internal/harness/repo.go internal/harness/apply.go internal/harness/harness_test.go`
(working tree, uncommitted; tree held still). One mechanism only: the plain-copy check
(`ownCheckout`), round 3 of round-1 finding 1 / round-2 finding 1.

Blobs pinned (`git hash-object` at review time):

| blob | path |
|---|---|
| `ecbbd8ebc36b06a20e9917492de58fdd8dea76bd` | `internal/harness/repo.go` |
| `500f8135b40d20f193009befebd236ded2203cb7` | `internal/harness/apply.go` |
| `491803af50fdb6292a9108ef0f306902309f958d` | `internal/harness/harness_test.go` |
| `21b08a070691744a62965c5b064e6b31a0b058af` | `cmd/sh-agent/main.go` (blast radius, unchanged by this delta) |

Round 2's blob for `repo.go` (`14b5a42…`) was never committed and is no longer in the
object store, so every size comparison below is against `5cdd4d6`, not against round 2.

## Verdict: FINDINGS (all low) — the mechanism is env-immune for the nested-copy case and has CONVERGED

The nested-copy case is closed, and closed by construction rather than by enumeration:
`ownCheckout` is `os.Lstat(filepath.Join(repo, ".git"))`, a syscall no environment
variable can redirect, and `gitEnv` is gone. Verified end to end against a real build,
not by reading: a harness-shaped copy at `outer/mirror` inside a real git repository
reports `(plain copy, not its own git checkout; update is refused)` from `sh-agent
status`, `sh-agent update` exits **1** with `not its own git checkout; a plain copy
cannot be updated with git` (status read bare, not through a pipe), no
`.agents/scratch/sh-agent-build-*` is created, so no build and no exec happen — and all
three observations are **identical** with `GIT_DIR`, `GIT_WORK_TREE`, `GIT_CONFIG_GLOBAL`,
`GIT_CONFIG_SYSTEM`, `GIT_OBJECT_DIRECTORY`, `GIT_INDEX_FILE`, `GIT_CEILING_DIRECTORIES`,
`GIT_DISCOVERY_ACROSS_FILESYSTEM`, `GIT_SSH_COMMAND` and `GIT_ASKPASS` all exported and
pointed at the enclosing repository. Installing from that copy under the same polluted
environment writes `"commit": ""` into `state.json` (round-2 finding 3 fixed; no foreign
commit is attributed). Nothing in this round is a new hole; the three findings below are
a missing written rationale and two residuals of the *accepted* trust premise, recorded
so a fourth round does not re-derive them. **Do not open round 4 on this mechanism.**

## Findings

- **[low / confidence: high]** `internal/harness/repo.go:108-116` (`Pull`), `:69-78`
  (`ownCheckout`) — the accepted rationale exists only in the dispatch message; it is
  written **nowhere on disk**. I grepped `.agents/docs/architecture/DECISIONS.md`
  (ADR-012…016), `REQUIREMENTS.md`, `TODO.md`, `AGENTS.md`, and the changes entries:
  no record says "git runs in the operator's inherited environment on purpose, because
  `sh-agent update` is `git pull` as the operator would run it by hand, and an
  allow-list would break custom SSH commands and proxies". Round 2 asked for exactly
  this sentence next to `Pull` if the trust premise was taken (its remediation (b)),
  and it did not land. Worse than absent, the comment that *did* land reads the other
  way: `ownCheckout`'s "git's answer can be redirected by its environment, so the check
  never asks git" tells a cold reader that git is sandboxed here, which is the
  misreading round 2 named as the more dangerous state. Three channels remain open by
  design and a reader cannot tell which are intentional: `$SH_AGENT_REPO` still selects
  which tree is built and executed (and `knownCheckout`, `main.go:350-353`, returns nil
  early when it is set), git's location variables still redirect the commands that run
  after the gate (finding 2), and git's config/auth variables still reach `git pull`.
  Remediation, one line, no code change: a comment on `Pull` (or an ADR paragraph)
  stating the premise and naming those three channels as accepted. Cheapest possible
  fix for the only defect that actually caused this review to be round 3.

- **[low / confidence: high]** `internal/harness/repo.go:83-91` (`Inspect`), reaching
  `internal/harness/apply.go:79-84` and `cmd/sh-agent/main.go:218-254` — the gate is
  env-immune; everything **after** the gate is not, and the leaked class is wider than
  the acceptance describes. Verified read-only against the built binary on two sibling
  scratch checkouts: `sh-agent status --repo realA` prints `(main @ e4d0dd0, clean)`
  normally and `(other @ 129af1e, uncommitted changes)` with `GIT_DIR` exported at an
  unrelated repository — the branch, commit and dirty flag of a repository that is not
  the one whose files were installed. That value is not cosmetic: `Install` writes it to
  `state.json` as provenance, `status` feeds it to `binaryOlder` (`main.go:249`,
  `:387-396`), and `Pull` would fast-forward against that other repository's upstream
  with `-C repo` as the work tree. The acceptance is stated for git's *configuration*
  and *authentication* variables (`GIT_CONFIG_*`, `GIT_SSH_COMMAND`, `GIT_ASKPASS`,
  `PATH`) — variables that say how git talks to a remote, where breaking them would
  break a legitimate custom setup. `GIT_DIR`, `GIT_WORK_TREE`, `GIT_INDEX_FILE` and
  `GIT_OBJECT_DIRECTORY` are a different class: they say *which repository is acted on*,
  no legitimate `sh-agent update` needs them, and they are set **ambiently** by git
  itself — any command run from a git hook, from `git rebase --exec`, or from a
  `git filter-branch`-style wrapper inherits `GIT_DIR`. So the reachable case needs no
  attacker at all, only an operator who runs `sh-agent update` from inside a hook or a
  rebase step, and the failure is silent (a wrong commit recorded, a pull against the
  wrong upstream). Remediation, in preference order: **(a)** accept it in writing —
  finding 1's one line covers this if it names these four variables, and no code moves;
  **(b)** only if it is ever revisited, `cmd.Env = append(os.Environ(), "GIT_DIR=",
  "GIT_WORK_TREE=", "GIT_INDEX_FILE=", "GIT_OBJECT_DIRECTORY=")` — a bounded set of four
  location variables that does not touch config, SSH, proxy or `PATH`, so the round-2
  objection to an unbounded deny-list does not apply. I am **not** asking for (b) in a
  round 4; (a) closes this finding.

- **[low / confidence: medium]** `internal/harness/repo.go:75-78` — `os.Lstat` accepts a
  `.git` **file**, so a tree can nominate any repository as its own and pass the gate.
  Verified: a plain copy carrying `.git` with `gitdir: <other>/.git` reports
  `(other @ 129af1e, uncommitted changes)` instead of `(plain copy…)`, i.e. the pointer
  inside the copied tree decides what `Inspect`, `Behind` and `Pull` act on. This is
  deliberate and the comment says so — `git worktree` and submodules need it, and a
  vendored harness submodule genuinely is its own updatable checkout. It is on the
  record only because the gate's question is "does this tree own its git state" and a
  gitfile answers "some repository does", with the pointer living in the tree being
  judged. Calibration, so it is not over-read: anyone who can write `.git` into a
  harness tree can also write `context/AGENTS.md` and `skills/*/SKILL.md`, which
  `install` deploys into the operator's home and which no gate inspects — strictly
  worse, and `knownCheckout` is the control for that, not this. **No change
  recommended;** recorded so round 4 does not present it as a discovery. If it is ever
  tightened, the cheap version is `st.IsDir() || <gitfile points inside repo>`, and
  that would break submodule checkouts — which is why I am not recommending it.

## Convergence

**Converged.** Round 1: "git inherits the environment, so a nested copy can pass
`ownCheckout`." Round 2: "the deny-list cannot be completed; stop asking git." Round 3
stops asking git — the question is now a filesystem stat, `gitEnv` and its seven-name
list are deleted, one subprocess is removed from every `Inspect` call on a plain copy,
and the refusal is a single sentinel (`errPlainCopy`) checked in three places. The
mechanism got smaller, which is the shape the global rule asks for; `repo.go` is 172
lines against 149 at `5cdd4d6`, and all 23 added lines are the gate, the sentinel and
their comments, with none of them a list that must grow when git adds a variable. The
same sentence cannot land a fourth time, because there is no longer an environment
channel into the *decision*: `os.Lstat` has no configuration. What remains (findings 2
and 3) is the trust premise the operator has accepted, and the only outstanding action
is to write that premise down once (finding 1). Round 4 on this mechanism would be a
process failure, not a security one.

## Checked

- **Secrets** — none in the delta. No credential, token or key in the three changed
  files; `errPlainCopy`'s message and `gitOut`'s error wrapping (`repo.go:161-172`)
  carry only constant argv and git's own stderr, never an environment value — this is
  strictly less exposure than round 2, where the scrubbed-variable list existed in
  memory. `status`'s new plain-copy line prints the repo path only. Carried over, still
  open, and outside this delta's named files: round-2 finding 5 on `.gitignore` —
  credential JSON (`*-key.json`, `application_default_credentials.json`) is still
  unmatched by any pattern, and the unanchored `.claude/agent-memory/` line is still
  absent while `/.claude/` is root-anchored.
- **Injection** — the environment was the injection channel under review and it no
  longer reaches the decision. Argument construction is unchanged and safe: every
  `exec.Command` takes an argv slice, `repo` is `filepath.Abs`'d by `FindRepo` before it
  reaches one, and `filepath.Join(repo, ".git")` is the only new path construction — it
  cannot escape `repo` because the second element is a constant. No shell, template,
  SQL or format-string interpolation added.
- **Authz/authn** — no authentication surface in the module. The authorization analogue,
  "which tree may be built and executed", is guarded by `knownCheckout` before both
  `install`/`sync` (`main.go:114`) and `update` (`main.go:167`) — I re-read the switch to
  confirm `update` is covered — and now by `ownCheckout` inside `Inspect`/`Behind`/`Pull`.
  Ordering is correct and the guards stay independent: the plain-copy gate fails closed
  *before* `git fetch`, so a refused tree never reaches `Reinstall`. `knownCheckout`
  keying on `st.Repo` (a path) rather than on the commit is why finding 2 is provenance
  and not privilege. `$SH_AGENT_REPO` short-circuiting `knownCheckout` is unchanged,
  pre-existing, documented, and consistent with the accepted premise — named in finding 1
  because the delta's new comments could be read as making tree selection env-proof.
- **Crypto** — none present, none added, nothing hand-rolled anywhere in the module.
  Unchanged by this delta.
- **Dependencies** — `go.mod` requires nothing beyond `go 1.24.4`; the delta adds one
  stdlib import (`errors`, in the test file) and no third-party package, so stdlib-first
  holds with nothing to justify. `govulncheck ./...` **was re-run this round** (repo.go
  is in the delta and one trace lands in it): two reachable stdlib vulnerabilities,
  unchanged from round 1 — **GO-2025-3956** (`exec.LookPath`, fixed go1.24.6) traced to
  `internal/harness/repo.go:135` in `Reinstall`, and **GO-2026-4602** (`os`, fixed
  go1.25.8) traced to `internal/harness/inventory.go:71`; plus 2 imported-but-uncalled
  and 43 required-but-uncalled. Pre-existing, rebutted in round 1, and **still tracked**
  — not in `TODO.md`, which is now empty, but in the newest handoff's Open list
  (`handoffs/2026-09-13-learning-loop-and-improvement-review.md:27`), which the file's
  own format note designates for operator items; I checked because round 2 cited
  `TODO.md:15`. `npm audit` not applicable and not run — no `package.json`, pure Go
  module. Also run on the delta: `gofmt -l .` clean, `go vet ./...` clean, `go test
  -race ./...` passing (both packages), `go version go1.24.4 linux/amd64`,
  `git version 2.47.3`.
- **GCP surface** — nothing in this delta. No IAM binding, service account, Cloud Run
  flag, Terraform, secret-as-env-var, or key file anywhere in the changed files.

## Test quality

`TestNestedCopyIsNotACheckout` has teeth for the case it names and is honest about it:
it asserts `Inspect`, `Behind` and `Pull` on a nested copy, then re-asserts `Inspect` and
`Behind` with `GIT_DIR`/`GIT_WORK_TREE` pointed at the enclosing repository, then
confirms the enclosing repository itself is *not* refused (the false-positive direction,
which is the one a careless "always refuse" fix would break). Two limits worth stating
rather than filing: the polluted-environment half can no longer fail for the right
reason — `ownCheckout` short-circuits before any git call, so the assertion now proves a
property of `os.Lstat`, and it would stay green even if every post-gate `gitOut` were
redirected (finding 2); and `Pull` is asserted only in the clean-environment half. The
`t.Setenv("GIT_DIR", "")` + `os.Unsetenv("GIT_DIR")` teardown pair flagged in round 2 is
still there and still harmless — `t.Setenv`'s cleanup restores the original — but both
lines together still read as if one is load-bearing when neither is. `TestBehindAndPull`
and `TestInstallMarksDirtyTree` continue to cover the own-checkout path, so the gate's
accept side is not asserted by the new test alone.

## Not examined

- Everything outside the three named files, plus `cmd/sh-agent/main.go`, which I read
  (status, update, `knownCheckout`, `binaryOlder`, the command switch) as blast radius
  but did not review as a diff — it is unchanged by this delta. The rest of the
  `5cdd4d6..working tree` diff (`context/`, `agents/`, `skills/**`, `README.md`,
  `ROADMAP.md`, `REQUIREMENTS.md`, `derive.go`, `.gitignore`) was reviewed in round 1 and
  is not re-read here; round 1's and round 2's "Not examined" sections still stand.
- The post-gate environment surface was demonstrated with `GIT_DIR` only (finding 2).
  `GIT_WORK_TREE`, `GIT_INDEX_FILE`, `GIT_OBJECT_DIRECTORY`, `GIT_CONFIG_COUNT`,
  `GIT_CONFIG_GLOBAL`, `core.hooksPath` and `url.*.insteadOf` were verified against the
  *gate* (all refused) but not individually re-verified against the commands that run
  after it; round 2 established that behaviour at the git level and nothing in this delta
  changes it.
- No end-to-end `git pull` against a poisoned upstream was performed, and no exploit was
  built. Every observation above is a read-only `status`/`update --check` run, a refused
  `update`, or one `install` into a scratch home under `--home`. Fixtures live in the
  session scratchpad, outside the repository; **no file in the repository was modified by
  this review** except this report.
- `.agents/docs/bugs/BUG-010-*`, the changes entries, and
  `testing/2026-09-13-verify-improvement-items.md` were opened only far enough to confirm
  where the plain-copy verification and the toolchain disposition are recorded; they are
  trail records with no execution surface.
