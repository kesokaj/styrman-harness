---
owner: security
date: 2026-09-13
---
# Security review — improvement-items   (2026-09-13)

Range reviewed: 5cdd4d6..working tree (uncommitted; tree held still)

Blobs pinned (`git hash-object` at review time):

| blob | path |
|---|---|
| `48f389773760920358944ac3be366923984864d1` | `internal/harness/repo.go` |
| `21b08a070691744a62965c5b064e6b31a0b058af` | `cmd/sh-agent/main.go` |
| `be2175f0fcc185f26b9d29fc0e39aae6bae85e95` | `internal/harness/derive.go` |
| `be1a984998ce6ce8aee4e5e4ee85ae1e31a76f9d` | `internal/harness/plan.go` (unchanged, read as blast radius) |
| `727bba7bbc8ab36ae958cadaf3e6836909ee4e80` | `cmd/sh-agent/main_test.go` |
| `0af142235be56f1a7076ad3bb0b4f38967cb2144` | `internal/harness/harness_test.go` |
| `d7276f962cd7ce1585d6651c9e89dba14a827c86` | `skills/gcloud-setup/SKILL.md` |
| `8fd8f3afc8f4eaa9273e17a56967a6443433d761` | `skills/gcloud-auth-verification/SKILL.md` (deleted; read at 5cdd4d6) |
| `654b6d0f88e0e18ea1ec72cbd19ef3f22549560b` | `context/dev-setup.md` |
| `137c4a6d9d37303586d194c94538e4b0fd5b39a8` | `context/AGENTS.md` |
| `2212cd466a0d89da3085b8d2dbe8e9e2675797a4` | `REQUIREMENTS.md` |
| `5b0ce78bbcbe19b4cbbdb22fcded18b07e1fe20e` | `.gitignore` |

## Verdict: FINDINGS

The new `ownCheckout` guard does close the round-3 nested-copy path it was written
for — `Inspect`/`Behind`/`Pull` all refuse a harness-shaped tree nested inside
another repository, and the new test proves it — but the guard asks git the
question in a way git will answer wrongly when `GIT_DIR` or `GIT_WORK_TREE` is in
the environment. I verified this empirically: with `GIT_DIR` set,
`git -C <nested-plain-copy> rev-parse --show-toplevel` returns the nested path
itself, so `ownCheckout` returns true for exactly the tree it is meant to reject,
and the subsequent `git -C repo pull --ff-only` then runs against the foreign git
directory with the harness tree as its work tree. Since `update` proceeds straight
into `Reinstall`, which runs `go build ./cmd/sh-agent` in that tree and executes
the result, the failure mode of the bypassed guard is local code execution from a
repository the operator did not choose. Same environment blind spot exists on the
pre-existing `gitOut` calls; the fix is one helper, not a redesign. Everything else
is low: `govulncheck` flags two stdlib symbols the code calls, but on reading both
advisories neither is reachable given the call shapes here.

## Findings

- **[medium / confidence: high]** `internal/harness/repo.go:73-81` — `ownCheckout`
  calls `gitOut(repo, "rev-parse", "--show-toplevel")`, and `gitOut`
  (`repo.go:164-175`) builds `exec.Command` without setting `cmd.Env`, so the
  child git inherits `GIT_DIR`, `GIT_WORK_TREE`, `GIT_COMMON_DIR`,
  `GIT_INDEX_FILE`, and `GIT_OBJECT_DIRECTORY` from the process. Verified in a
  scratch fixture: with an outer repo at `outer/` and a plain harness-shaped copy
  at `outer/nested/`, `git -C nested rev-parse --show-toplevel` prints
  `.../outer` normally but prints `.../outer/nested` when `GIT_DIR=.../outer/.git`
  is exported — so `ownCheckout` returns true and the guard is silently off. `Pull`
  (`repo.go:120`) then runs `git -C repo pull --ff-only` under the same inherited
  `GIT_DIR`, checking another repository's upstream content into the harness tree,
  which `update` immediately hands to `Reinstall` (`repo.go:137-162`) to build and
  execute. `direnv`, a git-hook context, or any wrapper that exports `GIT_DIR` is
  enough to reach this; no attacker-controlled input beyond the environment is
  needed. Remediation: scrub the git environment in one place — give `gitOut` and
  `Pull`'s `exec.Command` a `cmd.Env = append(os.Environ(), "GIT_DIR=",
  "GIT_WORK_TREE=", "GIT_COMMON_DIR=", "GIT_INDEX_FILE=", "GIT_OBJECT_DIRECTORY=")`
  (or `os.Unsetenv`-equivalent filtering), and add a test that exports `GIT_DIR`
  around the existing `TestNestedCopyIsNotACheckout` body.

- **[low / confidence: high]** toolchain — `govulncheck ./...` on go1.24.4 reports
  two stdlib vulnerabilities at **called-symbol** level, but I read both advisories
  and neither is reachable here; recording them so the next reviewer does not
  re-derive this, and because the bump is nearly free.
  `GO-2025-3956` / CVE-2025-47906 ("Unexpected paths returned from `LookPath`",
  fixed in go1.24.6) is traced to `internal/harness/repo.go:138`. The advisory's
  condition is that `LookPath` is passed `""`, `"."`, or `".."` while `PATH`
  contains executables rather than directories; `Reinstall` passes the constant
  `"go"`, so the condition cannot be met. (I initially wrote that a `.` element in
  `PATH` could make `exec.Command("go", …)` run a `go` from the working directory —
  that is **wrong** for this toolchain: since Go 1.19 such a resolution yields
  `exec.ErrDot` at `Run` instead of executing. Correction recorded rather than
  quietly dropped.)
  `GO-2026-4602` / CVE-2026-27139 ("FileInfo can escape from a Root in `os`", fixed
  in go1.25.8) is traced to `internal/harness/inventory.go:71`. The escape is
  specific to a `File` opened under an `os.Root`; this module never constructs an
  `os.Root`, so the plain `os.ReadDir` call is outside the advisory's shape.
  A third, `GO-2026-4970` ("Root escape via symlink plus trailing slash in `os`",
  fixed in go1.25.12), appears only at package level — not called.
  Remediation: bump the toolchain to go1.25.12 or later and re-run `govulncheck` so
  the scan is clean and future real hits are not lost in known noise. Both fixes
  are stdlib-only; the module has zero third-party requires, so nothing else moves.

- **[low / confidence: medium]** `cmd/sh-agent/main.go:245` —
  `buildinfo.ReadFile(filepath.Join(home, harness.BinPath))` newly feeds an
  operator-supplied path (`--home`) to a binary-format parser (ELF/Mach-O/PE).
  It fails closed (`err != nil` leaves the `unknown` line) and prints only the VCS
  revision, so nothing is leaked and nothing is executed; the residual risk is
  parser robustness against a malformed file placed at `$home/.local/bin/sh-agent`,
  and the fact that a symlink there is followed without a check. In this repo's
  single-operator, local-only model the file is one the operator runs anyway, so
  this is informational rather than exploitable. Remediation, if wanted: `os.Lstat`
  the path and skip a symlink, and cap by confirming the mode is a regular file.

- **[low / confidence: medium]** `.gitignore:11` — the `# Per-project subagent
  memory.` / `.claude/agent-memory/` pair was removed (improvement item 6 asked for
  it). The root-anchored `/.claude/` on line 8 still hides a top-level
  `.claude/agent-memory/`, so the common case stays covered, but a nested
  `*/.claude/agent-memory/` in any subdirectory is now untracked-visible. Agent
  memory files accumulate project context and occasionally quoted output, which is
  exactly the class of file that should not reach a commit by accident. This
  matters more than it first looks, because the file now contains **no
  secret-shaped patterns at all**: after the removal `.gitignore` is `.DS_Store`,
  `/sh-agent`, `/.claude/`, `/agent/`, `.agents/scratch/*`. There is no `.env`,
  `*.pem`, `*.key`, `*.tfvars`, `*.tfstate`, or credential-JSON line, so the global
  behavioural rule ("a secret-shaped path stops the commit"; "stage by name, never
  `git add -A` with untracked files") is the only thing standing between an
  accidentally created secret file and a push — and a behavioural rule is not a
  control. Remediation: keep an unanchored `.claude/agent-memory/` (or
  `**/agent-memory/`) line even though the feature is gone — a stale ignore costs
  nothing — and add the secret-shaped patterns the global rules already name. The
  improvement review itself records that sibling projects (`crew`,
  `homey-troubleshooting`) have committed agent-memory files already.

- **[low / confidence: medium]** `skills/gcloud-setup/SKILL.md:125-136` — the
  merged skill is the *only* surviving gcloud auth guidance, and its "Service
  Account Keys (CI/CD only)" section still documents
  `gcloud auth activate-service-account --key-file=key.json` and
  `export GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json`. That is a
  service-account key file on disk, which the global GCP posture otherwise avoids
  (Secret Manager, Workload Identity). The section is unchanged by this diff and
  does carry a "Never commit service account keys to git" warning, but with the
  companion skill deleted it now reads as the harness's single recommendation.
  Remediation: narrow it to "Workload Identity Federation for CI; a key file only
  when a platform offers nothing else, and then from a secret store, never from the
  repo or a developer machine."

## Checked

- **Secrets** — no credentials, tokens, keys, or `.env`-shaped files in the diff.
  `git diff 5cdd4d6` grepped for key/token/secret/password patterns: only prose
  about Secret Manager and the `key.json` example above. `REQUIREMENTS.md` R7 and
  `context/*` add no secret paths. Nothing is echoed into logs: the only new output
  is `versionFromInfo`'s VCS revision string and the `plain copy` status line.
  Correction to an earlier draft of this file: this repo's `.gitignore` does **not**
  and never did block `.env`, `*.pem`, or `*.key` — it carries only `.DS_Store`,
  `/sh-agent`, `/.claude/`, `/agent/`, and `.agents/scratch/*` (verified against
  `git show 5cdd4d6:.gitignore`). See the `agent-memory` finding.
- **Injection** — no SQL, no templates, no shell strings. Every external command
  is `exec.Command` with an argv slice (`repo.go:120,148,154,165`); `repo` is
  always `filepath.Abs`'d by `FindRepo` before it reaches them, and no user string
  is interpolated into a command line. The one genuine injection-adjacent channel
  is the *environment*, not the arguments — see the `GIT_DIR` finding.
- **Authz/authn** — no authentication surface in the binary. The authorization
  analogue here is "which tree may be built and executed": `knownCheckout`
  (`main.go:350-365`) still gates `install` and `update` on an explicit `--repo`,
  `$SH_AGENT_REPO`, or a match with the installed manifest, and the new
  `ownCheckout` adds a second, independent gate on the git-driven half. The
  ordering is correct — `knownCheckout` runs before `update()` and therefore before
  `Behind`/`Pull`/`Reinstall`. The new `TestNestedCopyIsNotACheckout`
  (`harness_test.go:603-642`) covers the nested case in a clean environment and
  asserts the enclosing repo is still accepted, so the guard does not over-refuse.
  No test covers the polluted-environment case.
- **Crypto** — none added, none present. Nothing hand-rolled: no auth, hashing,
  session, or JWT code anywhere in the module. The agent-memory removal touches no
  crypto. `git hash-object` is used here only for review provenance, not by the
  code.
- **Dependencies** — `go.mod` requires nothing beyond `go 1.24.4`; the stdlib-only
  rule holds and the one new import, `debug/buildinfo`, is stdlib, correctly chosen
  over shelling out to `go version -m`. No new dependency is introduced by this
  diff, so there is nothing to justify against stdlib-first. `govulncheck ./...`
  ran (go1.24.4, linux/amd64) and its two called-symbol hits are analysed above —
  both unreachable given the call shapes, toolchain bump recommended anyway.
  `npm audit` was not run and is not applicable: `npm` is on PATH but the repo has
  no `package.json` and no JavaScript.
- **GCP surface** — no IAM bindings, no Terraform, no `gcloud run deploy` in the
  changed code. The prose changes are a net improvement: the deleted
  `gcloud-auth-verification` skill's security-critical rule ("**DO NOT** attempt to
  fix the credentials via code or alternative tools") survives verbatim in intent
  as `gcloud-setup` step 2, and the old "Project Discovery" ladder that told the
  agent to read `GCP_PROJECT`/`PROJECT_ID` out of a `.env` file was deleted — the
  agent is now told to ask once and record in `docs/deploy/gcp.md`, which keeps it
  away from secret-shaped files. No wording anywhere in the merged skill tells the
  agent to fetch, mint, or repair credentials itself. `context/dev-setup.md` keeps
  "Secrets reach Cloud Run through Secret Manager (`--set-secrets`), never
  `--set-env-vars`". The one residual is the key-file section noted above.
- **Agent memory removal** — clean and a small net positive: no `memory:` key
  remains in any `agents/*.md`, the `derive.go` comment matches, and
  `internal/harness` never parsed the key (nothing in `derive.go`'s `toolMap` or
  the frontmatter reader keys on it), so no install-time failure is introduced. The
  replacement sentences point at `INDEX.md` and the newest `changes/` entries,
  which are committed artifacts rather than an uncommitted per-machine cache.

## Not examined

- The large prose rewrites in `skills/docker/**`, `skills/golang/**` (including the
  two new untracked files `skills/golang/references/golang-samber.md` and
  `golang-structure.md`), `skills/frontend/SKILL.md`, `skills/handoff/SKILL.md`,
  `skills/{bug,feature,research,feedback,retro,init-project,cloud-run-deploy}/` and
  `skills/VENDORED.md` were not read line by line — roughly 3,000 lines of the
  diff. What I did run over the whole diff: a grep of added lines for
  credential-shaped strings (`password|secret|api_key|token|bearer|private key|
  BEGIN .* PRIVATE|AKIA|ghp_|sk-…`) and for dangerous shell shapes
  (`curl … | sh`, `eval $`, `chmod 777`, `--no-verify`, `--allow-unauthenticated`,
  `--set-env-vars`, `insecure`, `skip-verify`). The only hits were prose *advising
  against* those things (Docker layer-secrets guidance, the `crypto/rand` and
  Argon2id tables in `golang-ecosystem.md`) plus one `POSTGRES_PASSWORD: app`
  placeholder in a compose example, which is fine. A full read for
  prompt-injection-shaped guidance did not happen; I did note that several
  `golang/references/*.md` "Audit mode" lines instruct launching up to five
  parallel sub-agents via the Agent tool — a cost and turn-budget question for the
  reviewer, not a security one.
- `.agents/docs/deploy/remote.md`, `README.md`, `ROADMAP.md`, `TODO.md`, and the
  `.agents/docs/{changes,handoffs,reviews}` additions were not reviewed; they are
  trail records with no execution surface.
- No dynamic testing of `sh-agent install`/`update` against a scratch home was
  performed; the `GIT_DIR` bypass was demonstrated at the `git` level in an
  isolated fixture, not by driving `sh-agent` itself (deliberately — a working
  end-to-end demonstration is an exploit, not a finding).
- `internal/harness/apply.go`, `state.go`, `adopt.go`, and `inventory.go` were read
  only where the diff reaches them (`BinPath`, `LoadState`, `LoadSkills`); their
  own path-handling and symlink behaviour was not audited in this pass.
