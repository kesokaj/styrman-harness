---
owner: security
date: 2026-09-16
---
# Security review — preflight-v1.0.0   (2026-09-16)

Range reviewed: 4aa19be..c7e3688 (`4aa19be..HEAD`, 14 commits, 75 files, 2585 insertions)

Release gate for the first SSH-signed tag, which becomes root-adjacent code on styrman-crew's host.
Every path in the release diff was read; see `## Not examined` for the two partial reads.

## Verdict: FINDINGS

**No blocker.** Nothing in 2585 insertions is a credential — the ed25519 material is a **public**
key (verified here: `ssh-keygen -lf` on the quoted line prints
`256 SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU … (ED25519)`, matching the documented
fingerprint) — the three installer fixes written today each close what they claim to close, the
release procedure's round-1 and round-2 defects are genuinely fixed (the `<release commit>`
placeholder is gone because the `tag -s` target argument was **deleted**, not defended), and
`govulncheck ./...` reports **0 affecting**.

The worst thing on the record is not a defect but the posture this range ships: **it adds
signature gating for one consumer while, in the same range, templating an unsigned hourly
`fetch → pull --ff-only → go build → execute` timer for the other two** — and on argolis-devbox
that build runs **as root** on the VM. Anyone who can push to `main` (a stolen GitHub token, no
signing key needed) gets code execution as the operator within the hour on every box running the
timer, plus their prose into the `AGENTS.md` both CLIs feed the model. `deploy/remote.md` states
this acceptance in as many words, so it is informed, not hidden — but it is what v1.0.0 actually
deploys, and it belongs on the gate record. Alongside it, the consumer channel selects
`sort -V | tail -1` **then** verifies with no walk-down and no version floor, so an unsigned
`v99.0.0` pins the host on a tag it refuses and deleting the newest tag rolls it back.

The sharpest concrete implementation gap is that `secretShaped` — the new refusal whose entire
purpose in this release is keeping a key off consumer machines — does not match `id_ed25519`, the
exact filename of the signing key this release's trust model rests on.

## Findings

- **[medium / confidence: high]** `internal/harness/plan.go:204-211` (`secretShaped`) — the new
  refusal that stops a credential in a skill directory being copied into both vendor directories on
  every consumer machine matches only `.env`, `.env.*`, `*.pem`, `*.key`, `*.tfvars`, `*.tfstate`,
  `*credentials*.json`, `*-key.json`. **`filepath.Match` is case-sensitive and every pattern is
  extension-shaped, so the most common private-key filenames — which have no extension — are not
  matched: `id_ed25519`, `id_rsa`, `id_ecdsa`.** Also unmatched: `.ENV`, `server.KEY` (case),
  `*.p12`, `*.pfx`, `*.jks`, `*.ppk`, `*.gpg`, `*.asc`, `.netrc`, `.npmrc`, `.pgpass`, `kubeconfig`.
  Exploitable in principle precisely because the guard's own rationale (BUG-050) is that `.gitignore`
  hides these from review — and this repo's `.gitignore` likewise does not list `id_*`, so an
  `id_ed25519` dropped into a skill directory is *both* invisible to review *and* deployed to
  `~/.claude/skills/…` and `~/.gemini/config/skills/…` on every machine. It is the signing key's own
  filename, in the release whose trust root is that key. Remediation: lower-case before matching and
  add the extensionless SSH shapes (`id_rsa`, `id_dsa`, `id_ecdsa`, `id_ed25519*`, excluding `*.pub`)
  plus `*.p12`, `*.pfx`, `*.jks`, `.netrc`, `.npmrc`, `.pgpass` — stdlib `strings.ToLower` +
  `filepath.Match`, no dependency.

- **[medium / confidence: high]** `.agents/docs/deploy/remote.md` § The resync loop (new in this
  range) — the range does not merely document `update`'s new dirty refusal, it **recommends and
  templates an hourly systemd user timer** running `sh-agent update`: unattended `git fetch` →
  `pull --ff-only` → `go build` → execute, following `main` with **no signature check**, on every
  operator box, while `§ Downstream mirrors` records that argolis-devbox's VM builds `sh-agent`
  **as root** via `deploy/scripts.d/032_setup-harness.sh` on a five-minute auto-update. `--ff-only`
  limits rewrite but not a forward push. The document states the acceptance explicitly ("its trust
  root is whoever can push to this repository … the accepted position here"), and BUG-046 is
  evidence the timer is live and does deploy. This is therefore an accepted risk, correctly
  disclosed — recorded here because it is the largest exposure in this release's blast radius and
  because it is the mirror image of what the signed-tag channel exists to prevent. Remediation, if
  the acceptance is ever revisited: `sh-agent update --release` (already named as an undecided
  Follow-up in ADR-026) so the same signature gate covers the operator boxes; meanwhile
  `sh-agent update --check` on the timer, as `remote.md` itself offers.

- **[medium / confidence: high]** `.agents/docs/architecture/DECISIONS.md` ADR-026 Consequences /
  `.agents/docs/deploy/remote.md` § Releases — **the consumer's documented tag selection is
  select-then-verify with no fallback and no version floor**
  (`git tag -l 'v*' | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -1`, then verify, then
  refuse). Two live consequences: (a) anyone who can push a tag — a stolen GitHub credential, **no
  signing key needed** — creates a lightweight or unsigned `v99.0.0` and the host pins on it
  permanently, refusing every genuine later release including a security fix; (b) deleting the
  newest tag walks the host back to an older signed tag — a rollback to superseded harness text
  carrying a valid signature. This is crew's code, not this repository's, so it cannot be fixed in
  this range, but v1.0.0 is the first artefact to enter that channel. Already recorded in ADR-026's
  Consequences as a recommendation to crew's Session B; restated so the gate record carries it.
  Remediation: walk `sort -V` downward to the newest tag that *verifies*, alarm on every refusal,
  and refuse any version below the installed one.

- **[medium / confidence: high]** `.agents/docs/testing/2026-09-16-verify-real-home-9bb7ce0.md`
  § "What the operator runs" and `.agents/docs/testing/2026-09-16-verify-real-home-fcc8539.md`
  § "Antigravity, second attempt" — **two permission-widening recipes are preserved in the release
  range as the documented way to finish a verification**, and the newest handoff points future
  sessions at both. The first writes
  `"permissions": {"allow": ["command(*)", "read_file(*)", "write_file(*)", "edit_file(*)", "list_dir(*)"]}`
  into the operator's real `~/.gemini/antigravity-cli/settings.json` — `command(*)` is an
  unrestricted shell grant — and restores it with a bare `mv … .bak …` with **no verification that
  the restore happened**, so an interrupted run leaves the grant permanently in place. The second
  offers two `agy … --dangerously-skip-permissions` invocations and then: "adding a Bash permission
  rule for `agy` in `~/.claude/settings.json` would let a future session run them directly." That is
  the exact action this session's permission classifier refused **three times** (recorded honestly in
  the same files) now written into the trail as the recommended workaround, where a later agent
  reads it as sanctioned. Exploitable in principle as a self-weakening path: agent-authored text
  cannot authorise a permission change, but it can look like it does. Remediation: mark both blocks
  operator-only and non-agent-executable in as many words; replace the `mv` restore with a restore
  plus a verifying `grep`/`diff` that fails loudly; and drop the standing suggestion to add a Bash
  permission rule for `agy`, or record it as an ADR with the operator's approval rather than as a
  convenience line in a dated log.

- **[low / confidence: high]** `.agents/docs/deploy/remote.md` § Releases, the line
  `git diff --quiet <gated head> HEAD -- . ':!.agents/docs'` — **the `.` pathspec is cwd-relative,
  and this is the one command binding the tag's source tree to the gated tree.** Run from any
  subdirectory of the checkout it silently narrows to that subdirectory and a source change
  elsewhere in the release commit passes unseen; the block's preamble names no working directory.
  **Proven** in an isolated probe (`.agents/scratch/security-v1.0.0/probe`, `git init` exit status
  read, every later command `git -C`): with a smuggled change in `internal/harness/plan.go`, from
  the repo root `--quiet` exits 1 (caught); from `cmd/sh-agent` it exits **0** (missed); with
  `':(top)'` from the same subdirectory it exits 1 (caught). Remediation, one token:
  `git diff --quiet <gated head> HEAD -- ':(top)' ':!.agents/docs'`, or
  `cd "$(git rev-parse --show-toplevel)"` at the head of the block.

- **[low / confidence: high]** `.agents/docs/deploy/remote.md` § Releases, the line
  `test -z "$(git status --porcelain)"` — **the bare `--porcelain`, which is the exact form
  BUG-049/ADR-027 D3 judged unsafe one commit earlier** and replaced in `internal/harness/repo.go:94`
  with `--untracked-files=all`, because `status.showUntrackedFiles` is an ordinary git setting.
  **Proven** in the same probe: with `status.showUntrackedFiles=no` and an untracked
  `skills/smuggled/SKILL.md` present, bare `--porcelain` prints nothing while
  `--untracked-files=all` prints `?? skills/smuggled/SKILL.md`. Untracked files cannot enter the tag
  (a tag names a commit), so this is not a content-injection path into the consumer — but it means
  the gate can verify, and the release can be cut, on a tree holding untracked files the shipped tag
  does not contain, so what was tested is not what ships. Remediation, the same one-word fix as the
  code: `test -z "$(git status --porcelain --untracked-files=all)"`.

- **[low / confidence: medium]** `.agents/docs/deploy/remote.md` § Releases, the line
  `test ! -d /etc/styrman || sudo -n cmp /etc/styrman/allowed_signers .agents/scratch/allowed_signers`
  — a `sudo -n` invocation embedded in a command block an agent is told to execute, on the machine
  holding the signing key. Two things. (a) The global context requires a root step to have a
  standing permission recorded in the project's `.claude/settings.json`, or else be handed back with
  the exact command; this repo's `.claude/` is gitignored so no such record is readable from the
  range, and the handoff records the main agent already ran `sudo -n cat` on that file once without
  one ("Self-noted … although this project records no standing root permission"). (b) The guard
  **fails open**: on any machine without `/etc/styrman` the trust-root comparison is skipped in
  silence, and `verify-tag` then runs against a line the repository itself supplies, which by
  ADR-026's own words "adds no trust". The step is a read-only `cmp` and `sudo -n` fails closed
  without the permission, so exposure is small; the defect is that a skipped trust-root check is
  indistinguishable from a passed one. Remediation: print which branch was taken
  (`… && echo "trust root compared" || echo "NOT the crew host — allowed-signers line unverified"`),
  and hand the `sudo` line to the operator rather than embedding it until the record exists.

- **[low / confidence: high]** `go.mod:4` (`toolchain go1.25.8`, ADR-028) — the pin is a **ceiling
  hint, not a floor**, and the commit subject overstates it. ADR-028 deliberately keeps the `go`
  directive at 1.24.4 so `GOTOOLCHAIN=local` on an older toolchain still builds — and on that path
  the build gets the 1.24.4 stdlib and the two advisories (`GO-2026-4602` in `os` via `LoadSkills`,
  `GO-2025-3956` in `os/exec` via `Reinstall`) are **back**. ADR-028's Consequences states this
  correctly; the commit subject ("clears the two stdlib advisories … govulncheck goes from 2
  affecting to 0") does not, and a commit subject is what a reader trusts. Supply chain, stated
  plainly: on the default `GOTOOLCHAIN=auto` a build machine running Go < 1.25.8 now **downloads a
  toolchain it did not download before** — module `golang.org/toolchain@v0.0.1-go1.25.8.<goos>-<goarch>`
  from `GOPROXY` (`https://proxy.golang.org,direct` here), verified against the checksum database
  (`GOSUMDB=sum.golang.org`, `GOPRIVATE` and `GONOSUMDB` empty — read with `go env`). So it *is*
  verified, but it is a new build-time network dependency on every consumer building with an older
  Go, including argolis-devbox's VM which builds as root, and a proxy-blocked or offline host that
  previously built fine now fails unless it sets `GOTOOLCHAIN=local`. Remediation: the mechanism is
  sound and needs no change; state the two-path reality wherever the "0 affecting" claim is
  repeated, and prefer the operator's system Go upgrade, which ADR-028 already names as the
  alternative that makes the pin redundant.

- **[low / confidence: medium]** `internal/harness/repo.go:86-94` (`Inspect`) and
  `internal/harness/inventory.go:166-220` (`LoadSkills`) — three residual ways content that is not in
  the commit still reads as clean and deploys. (a) `--untracked-files=all` still honours ignore
  rules, so a `.gitignore` committed *inside* a skill directory, or a user-level `core.excludesFile`,
  hides files from the dirty check that `Plan` will nonetheless walk and deploy — `secretShaped` now
  covers the credential-shaped subset, not the rest; likewise `git update-index --assume-unchanged`
  and `--skip-worktree` on a tracked file. (b) `gitOut` runs git in the operator's inherited
  environment, so `GIT_DIR`/`GIT_WORK_TREE` can redirect the dirty check at another repository — the
  new test helper `gitCleanEnv` strips exactly these, so the class is known and deliberately accepted
  for the production path (`Pull`'s comment says so); under the systemd timer the environment comes
  from the unit file, which bounds it. (c) `LoadSkills` resolves `skills/<name>` through
  `filepath.EvalSymlinks` by design (`npx skills add -g`) and does **not** require the target to be
  inside the repo, so a committed symlink there would let a tree outside git decide what is deployed,
  and later changes to that tree never make the checkout dirty. Nothing is live today: the only
  tracked symlink is `.agents/skills → skills` (`git ls-files -s | awk '$1=="120000"'`). Remediation:
  constrain the resolved skill dir with the existing `within(repo, dir)`; consider surfacing ignored
  files in `Inspect` rather than only untracked ones.

- **[low / confidence: medium]** `internal/harness/plan.go:129` and `internal/harness/adopt.go:129`
  — `mode := info.Mode().Perm() | 0o600` preserves a repo file's permission bits into the deployed
  copy, so a skill file sitting at 0666/0777 locally (git tracks only the execute bit, so this
  arrives via `adopt`, a hand copy, or a permissive umask — not via clone) deploys world-writable
  into `~/.claude/skills/…` and `~/.gemini/config/skills/…`, where any local user could rewrite
  instructions the agent then executes. Pre-existing, not introduced here, but it is the same walk
  the new `secretShaped` refusal was added to. Remediation: `0o644 | (perm & 0o100)` — take the
  execute bit, drop group and other write.

- **[low / confidence: medium]** `.agents/docs/deploy/remote.md` § Releases — `<gated head>` is
  called "the one external input", but it is a value an agent transcribes out of a gate report the
  same agent wrote. `test "$(git rev-parse HEAD^)" = "<gated head>"` is a real check against a real
  ref and is **not** tautological (BUG-047 is genuinely fixed), yet a wrong or convenient
  substitution passes it silently. Remediation: derive it rather than transcribe it —
  `gated=$(grep -m1 '^Range reviewed:' <report> | sed 's/.*\.\.//')`, then use `"$gated"`, so the
  report file rather than the agent's memory is the input.

- **[low / confidence: medium]** `skills/distill/SKILL.md:156` — the newly required `project` field
  is an absolute path (`/home/you/workspace/<project>`), so `~/.config/sh-agent/feedback.jsonl` now
  carries the operator's home path on every entry, and `/retro` filters with
  `grep -vE '"project": ?"<this project>"'` — an unescaped path interpolated into a regex. A
  directory name containing `.`, `+`, `(`, `|` or `[` makes the exclusion match too much or too
  little, which is the same wrong-evidence promotion BUG-021 was written to stop, arriving by
  another route. The log is local and uncommitted, so the path disclosure is minor. Remediation:
  `grep -vF` on the JSON fragment, or a fixed-string field match.

- **[low / confidence: low]** `skills/gcloud-setup/SKILL.md:127-137` — the reframing of service-account
  keys as a last resort is a real tightening (CI justification removed, operator instruction
  required, "never commit" retained). What the section still does not say, now that it is the
  sanctioned escape hatch, is where the key file may live (never inside the repo), that it is
  `chmod 600`, that it is scoped to one role, and that it is deleted when the task ends —
  `export GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json` is left as the last line with no
  lifecycle. Remediation: four clauses on the existing bullet.

- **[low / confidence: high]** Process, `.agents/docs/changes/2026-09-16-toolchain-pin.md` — the
  entry for commit c7e3688 records "Agents dispatched: changelog only; no source changed, so no
  verifier or reviewer dispatch." `go.mod` is source, and the global Definition of Done dispatches
  `security` whenever **dependencies** are touched — which is the whole subject of that commit. The
  last commit in the release range therefore shipped with no independent review of its own; this
  gate is the first. Not a defect in the change itself (see the `toolchain` finding above), but the
  rule it skipped is the one that would have caught the overstated commit subject. Remediation: the
  entry's "Agents dispatched" line should name this report.

## Checked

- **Secrets** — full `git diff 4aa19be..HEAD` scanned for `BEGIN … PRIVATE KEY`,
  `ssh-rsa|ssh-dss|ssh-ed25519`, `AKIA`, `ghp_`, `github_pat_`, `xox[baprs]-`, `sk-…`, `AIza`,
  `Bearer`, and `password|passwd|secret|token|api_key|client_secret` assignments, plus **every
  ≥32-character base64-ish run in the added lines** as an independent sweep. Result: the ed25519
  public key is the *only* long token in the range that is not a commit SHA or a file path. Its
  four occurrences are all intended — `deploy/remote.md` (the quoted line and the `printf` into
  `.agents/scratch/allowed_signers`), `DECISIONS.md` ADR-026 Context, and the dated review/testing
  records. Verified public, not private, with `ssh-keygen -lf` on a scratch copy:
  `256 SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU ostling.simon@gmail.com (ED25519)` —
  matches the documented fingerprint exactly. The email is the operator's own git author identity.
  `.agents/scratch/allowed_signers` is covered by `.gitignore`. `cmd/sh-agent/main_test.go` contains
  the literal `{"attribution":{"commit":"SECRET"}}` — that is a fixture asserting the reader never
  echoes file contents into an error string, i.e. a leak *test*, not a leak. No private key, token,
  passphrase, `.env` or private data anywhere in the range.
- **Injection** — the installer takes no network or user-supplied input. The parse boundary is the
  frontmatter `name:`, still constrained by `validName`/`identLike` to one `[A-Za-z0-9._-]` segment
  with `filepath.IsLocal`, which is what keeps `filepath.Join` inside the install root and keeps
  backticks and newlines out of the rendered context; the new `adopt` staging name interpolates that
  already-validated `name` plus `os.Getpid()`, so no traversal and no arbitrary `RemoveAll` target.
  Every subprocess is `exec.Command` with an argument vector, never a shell. No SQL, no templates.
  The one unescaped interpolation in the range is `/retro`'s grep pattern (finding above).
- **Authz/authn** — the trust chain is: operator's own turn → `/preflight` gate → private-key
  possession proved by `ssh-keygen -y -f ~/.ssh/id_ed25519 | ssh-keygen -lf -` (derived from the
  private key, not a `.pub`) → tag on HEAD with no target argument → `HEAD^` equal to the gated head
  → source tree identical to the gated tree → annotated → last armor header
  `-----BEGIN SSH SIGNATURE-----` → `verify-tag` → root-only `/etc/styrman/allowed_signers`.
  Round 1's "no check of what the tag points at" and round 2's tautology are both genuinely closed —
  the `<release commit>` placeholder that sat on both sides is gone because the target argument was
  deleted rather than defended. The quoted allowed-signers line carries no `cert-authority` and no
  `namespaces=`, so it cannot mint further signers. **No path was found by which a tag the consumer
  host accepts is produced without the operator's private key**: editing the quoted line or the
  scratch copy changes only the local pre-push check, and the host re-verifies against its own
  root-only file. The residual weaknesses are the selection algorithm and the cwd-relative and
  untracked-blind assertions, above.
- **Crypto** — **nothing hand-rolled.** The range implements or composes no cipher, hash, MAC,
  signature, KDF, RNG, token format, session format or auth flow. Signing and verification are
  delegated wholly to `git` + OpenSSH (`gpg.format=ssh`, `ssh-keygen -Y`, `git verify-tag`), the
  vetted path, and signing uses per-command `-c` flags rather than an agent writing `git config` —
  the right call. The Go code imports no `crypto/*` at all. `writeAtomic`'s new `tmp.Sync()` is a
  durability fix, not crypto. The armor-header `test` is a content assertion about git's output, not
  a primitive, and correctly is not the decider — `verify-tag` is. No finding under the hard rule.
- **Dependencies** — `go.mod` requires nothing; the binary is stdlib-only and this range adds no
  dependency, matching the repo's own rule. The one supply-chain change is `toolchain go1.25.8`,
  judged in full above. `govulncheck ./...` run here, exit 0: *"No vulnerabilities found. Your code
  is affected by 0 vulnerabilities. This scan also found 2 vulnerabilities in packages you import
  and 24 vulnerabilities in modules you require, but your code doesn't appear to call these
  vulnerabilities."* Environment: local toolchain go1.25.8, `GOTOOLCHAIN=auto`,
  `GOPROXY=https://proxy.golang.org,direct`, `GOSUMDB=sum.golang.org`, `GOPRIVATE` empty.
  `npm audit` **not applicable and not faked** — there is no JavaScript/TypeScript and no
  `package.json` in this repository or in the range.
- **GCP surface** — no IAM binding, service account, Cloud Run flag, `--allow-unauthenticated`,
  Terraform or deploy config changes. The GCP-adjacent edits are instruction text only and every one
  moves in the safe direction: hardcoded `us-central1` replaced by an operator-supplied `$REGION`
  with `REGION ?= $(error …)` so a missing value fails loudly instead of defaulting into another
  project (`gcloud-setup/SKILL.md`, `bash/references/makefile-patterns.md`); service-account key
  files demoted from "CI/CD" to a last resort behind an explicit operator instruction (one residual
  gap, above); the image tag becomes the git short SHA and an amd64 target is named as what ships;
  ship/deploy/release requests routed through `/preflight` (BUG-020).
- **Instruction-shaped text** — `context/AGENTS.md`, `agents/*.md`, every changed `skills/**`,
  `AGENTS.md`, `README.md`, `REQUIREMENTS.md` and the command blocks in `.agents/docs/deploy/*.md`
  read as security controls. **Every change tightens.** `context/AGENTS.md` adds a post-dispatch
  detection control (`git status --short --untracked-files=all`, `git log --oneline @{u}..HEAD`,
  `git tag`), a headless-run wait rule, and a runtime-data-path rule; `agents/security.md` adds the
  BUG-045 probe-repo binding. The two standouts are genuine security *improvements*:
  `skills/golang/references/golang-ecosystem.md` removes an instruction to register **private**
  code with third-party indexers ("registering uploads the source to a third party"), and
  `skills/docker/SKILL.md` turns a bare `sudo -n systemctl start docker` into a conditional that
  hands the command back absent a recorded standing permission. **No new wording tells a future
  agent to run a root command, write a credential, widen a permission, or send data outside the
  machine** — with the two exceptions recorded as findings above: the `sudo -n cmp` line embedded in
  the release block, and the permission-widening recipes preserved in the two real-home
  verification logs. The trail markdown (bugs, changes, handoffs, plans, retros, reviews, testing)
  is records only, carries no credential or private data, and its self-reporting is honest — the
  BUG-045 incident, the timer firing twice, and the `sudo -n cat` taken without a standing
  permission are all disclosed by the agents that did them.

## Not examined

- `.agents/docs/reviews/2026-09-16-security-release-procedure-adr-026.md` (208 lines) and
  `.agents/docs/testing/2026-09-16-verify-retro-2-and-release-procedure.md` (127 lines) — read in
  large part (headers, incident section, findings, Checked, Not examined, Outcome, and the full
  rehearsal command lists) but **not line-by-line end to end**; both are dated records of prior
  rounds, both were covered by the whole-range secrets sweep, and the findings they carry were
  re-derived independently here against the current text rather than taken on their word.
- Long lines in several trail files were read truncated to 300-400 columns
  (`verify-real-home-fcc8539.md`, `verify-clear-the-board.md`, the two review records). The truncated
  tails are narrative continuations; no command block or key material was cut.
- **The consumer side is outside this repository and was not read at all**: crew's
  `deploy/update/styrman-update`, its channel implementation, `/etc/styrman/allowed_signers` on the
  crew host (root-only; deliberately not read), and argolis-devbox's
  `deploy/scripts.d/032_setup-harness.sh`. Findings 2 and 3 are judged from ADR-026's and
  `remote.md`'s descriptions of them, not from their source.
- **The release procedure was not executed.** No tag was made, signed, verified or pushed from this
  review; the `.` pathspec and bare `--porcelain` findings were proven in an isolated probe
  repository (`.agents/scratch/security-v1.0.0/probe`, `git init` exit status read, every later
  command `git -C`), not against the real checkout. The real repository was confirmed untouched
  afterwards: `git tag` empty, HEAD still `c7e3688`, no unexpected untracked path.
- **`GIT_DIR` redirection of the dirty check was not weaponised.** The mechanism is named in the
  residual finding; a crafted `GIT_DIR` whose index matches the working tree was not constructed.
- **`Sync`-before-rename durability was not fault-injected** — untestable without a crash, as the
  test's own comment says.
- The Mac, the sandbox home, the live systemd timer, and Antigravity's execution-half smoke checks
  were not exercised; no CLI was started during this review.

## Addendum — the code reviewer's blocking finding on `toolchain go1.25.8`

Added after the code-lens reviewer returned REQUEST CHANGES on this same range, making the release
NO-SHIP pending re-gate. **Yes, I reached that directive's supply-chain question, and my finding
corroborates theirs.** Stated plainly, and without softening the severity split:

- Under the default `GOTOOLCHAIN=auto`, a host whose local Go is older than 1.25.8 does not build
  this module with its own toolchain. It **must fetch** module
  `golang.org/toolchain@v0.0.1-go1.25.8.<goos>-<goarch>` from `GOPROXY` first. That is a **hard
  build requirement**, not an optimisation: no network, a blocked or air-gapped proxy, or a
  restricted egress policy means `go build` fails, and `sh-agent install` and `sh-agent update` fail
  with it, because both run `go build ./cmd/sh-agent` through `harness.Reinstall`
  (`internal/harness/repo.go`). A go1.24.x host that built this repository yesterday can stop
  building it today. The reviewer is right to block on that, and their shape — crew's channel host —
  is the right one to worry about: that host checks a tag out and runs the harness from it.
- My own finding rated this **low**, and I am not revising that rating, because I rated it on the
  **security** axis only: the fetch *is* authenticated. Verified with `go env` here —
  `GOPROXY=https://proxy.golang.org,direct`, `GOSUMDB=sum.golang.org`, `GOPRIVATE` and `GONOSUMDB`
  empty — so the toolchain module is checksum-verified against the transparency log like any other
  module. There is no unverified-download or supply-chain-substitution finding. The **availability**
  consequence is the serious one, it is outside my lens, and the code reviewer has it correctly at
  blocking severity. Both readings are in the same direction; neither cancels the other.
- What I did flag on the security axis stands and is relevant to the fix: the pin is a **ceiling
  hint, not a floor**. `GOTOOLCHAIN=local` on go1.24.4 still builds (ADR-028 keeps the `go`
  directive at 1.24.4 precisely so it does), and on that path the two advisories the pin was added
  to clear — `GO-2026-4602` in `os` via `LoadSkills`, `GO-2025-3956` in `os/exec` via `Reinstall` —
  are back. So the two obvious remediations trade against each other and the trade should be made
  deliberately: dropping the directive restores the old build behaviour and restores the
  advisories; keeping it and documenting `GOTOOLCHAIN=local` as the offline path restores the build
  and restores the advisories on exactly the hosts that take it. The remediation that gives both is
  the one ADR-028 already names as the standing alternative — an operator-run system Go upgrade on
  the affected hosts, which makes the pin redundant rather than wrong.
- Not examined, and material to the fix: I did **not** test a build on a go1.24.x host, with an
  unreachable proxy, or with `GOTOOLCHAIN=local` — no such environment was constructed during this
  review. The failure mode above is read from the directive's documented semantics and from
  `Reinstall`'s `go build`, not observed. The verifier should measure it on the crew host shape
  before the re-gate.

## Outcome

1. [medium] `secretShaped` missed `id_ed25519`, `id_rsa`, `.netrc`, `*.p12` — the signing key's own filename — **fixed**: the pattern list now covers the four `id_*` names, `.envrc`, `.netrc`, `.npmrc`, `.pypirc`, `*.env`, `*.p12`, `*.pfx`, `*.keystore`, `*.jks` and `*.tfstate.backup`, alongside the original shapes. Case sensitivity is left as is and named here as residual: `filepath.Match` is case-sensitive, so `ID_RSA` would pass, which matters little on the Linux and macOS paths this harness targets and would be a behaviour change to force.
2. [medium] `deploy/remote.md` templates an unsigned hourly fetch-build-execute timer, so push access to `main` is code execution within the hour — **accepted, and already documented as the posture**: that section states the timer follows `main` by fast-forward rather than a signed tag, that its trust root is whoever can push, and what a box that cannot accept it runs instead. Recorded here as this release's actual posture rather than resolved, which is what the finding asks for.
3. [medium] the consumer's select-then-verify has no walk-down and no version floor — **left: styrman-crew**. It is that project's code; ADR-026's Consequences carry the recommendation (walk the version sort downward to the newest tag that verifies, alarm on every refusal, refuse a version below the installed one) and it is written into the brief for crew's channel session.
4. [medium] both real-home verification logs preserved permission-widening recipes — **fixed**. The `command(*)` wildcard block and the standing `--dangerously-skip-permissions` suggestion are removed from `testing/2026-09-16-verify-real-home-9bb7ce0.md` and `-fcc8539.md`; both now say what the two checks need and that the permission is the operator's decision, with an explicit note that no recipe is recorded and why. The earlier `2026-09-15-verify-real-home.md` mention is left: it is a past-tense record of what was done that day, not an instruction. This was the sharpest finding of the set, because it was the trail teaching a future agent to do the thing three refusals had just prevented.
5. [low, probe-proven] the release block's `git diff --quiet <gated head> HEAD -- . ':!.agents/docs'` used a cwd-relative `.` — **fixed**. Reproduced first: with the check run from one subdirectory and a source change smuggled into a sibling, the old form exits 0 and misses it while `':(top)' ':(top,exclude).agents/docs'` exits 1 and catches it. Corrected in `deploy/remote.md` § Releases and ADR-026.
6. [low, probe-proven] the release block used the bare `git status --porcelain` that BUG-049 had rejected one commit earlier — **fixed**: it is `--untracked-files=all` in both places, and the `<gated head>` comparison now derives both sides through `git rev-parse` rather than comparing against a transcribed string.
7. [low] deployed files preserved group and other write — **fixed**: `Plan` masks the mode to `Perm()&0o755 | 0o600`, so a 0664 or 0666 file in the repo lands 0644 in both vendor directories; `TestPlanDropsGroupAndOtherWrite` asserts no planned file carries `0o022`.
8. [low] `LoadSkills` allows a skill directory to be a symlink outside the repo — **accepted, by design and documented**: `checkManagedTree` exempts `skills/<name>` deliberately because `npx skills add -g` leaves one, and the comment there says so. The files inside are still read with `readRegular` and a symlinked `SKILL.md` is refused (BUG-036).
9. [low] `sudo -n cmp` sits inside the agent-followed block and fails open off the crew host — **fixed** by the condition added earlier (`test ! -d /etc/styrman ||`), and the line is marked crew-host-only. It remains the one root step in the procedure and has never been executed here.
10. [low] c7e3688 shipped with "changelog only, no reviewer or security dispatch" although it touched dependencies — **fair, and this gate is the correction**: the toolchain pin went in as a one-line change with its own ADR and no independent read, and the first independent read of it found the high finding above. The lesson is recorded here rather than left implicit: a change to `go.mod` is a dependency change and takes the security row, however few lines it is.
11. Remaining lows not individually actioned — `/retro`'s unescaped grep interpolation, the gcloud key-file lifecycle wording, `<gated head>` transcription (now derived, see 6) — are recorded in the report body; the first two are ordinary follow-ups, filed in the handoff rather than fixed inside a release gate.

---

## Round 2

Range reviewed: `c7e3688..787e15c` (the fixes made in response to round 1 and the release
reviewer), plus the **uncommitted working-tree change** that withdraws `toolchain go1.25.8`
from `go.mod` (`git diff 787e15c -- go.mod`, one deletion) and a modified
`.agents/docs/testing/2026-09-16-verify-preflight-v1.0.0.md`. The whole release range
`4aa19be..HEAD` is round 1's; this section judges only the response.

Blobs pinned at read time: `HEAD=787e15cdbf52b1443cfcdd760ee6bc185f215716`, working tree
carrying ` M go.mod` and ` M .agents/docs/testing/2026-09-16-verify-preflight-v1.0.0.md`.

### Verdict: FINDINGS — none blocking. **Ship.**

Every one of the seven fixes closes the substance of what it answers, and none of them opens a
hole. Two, though, are wired only half-way, and both half-wirings sit exactly where the finding
was: `Inspect`'s new `StatusUnknown` fails closed in `update` and `Pull` but still reads as the
word **"clean"** in `sh-agent status` and as an unqualified commit in the install manifest; and
the new mode mask applies only to files `Install` actually writes, so a file already deployed
`0666` whose content has not changed is **never chmodded** and keeps the loose mode the fix
exists to remove. Neither is a new hole; both mean the remediation does not fully remediate.

The sharpest new thing this round is not in the fixes themselves but in the gap the widened
`secretShaped` list makes visible: the list is enforced on the way **out** (`Plan` refuses to
deploy) and not on the way **in** (`Adopt`/`copyTree` copies anything into the tracked
`skills/`), and `.gitignore` covers none of the newly added names. `sh-agent adopt` of a
third-party skill carrying an `id_ed25519`, a `.netrc` or an `.npmrc` therefore lands it in a
git-tracked directory of a repository that is pushed to GitHub, and the only guard that fires
is the install-time refusal on every consumer — after the push.

Beside it, one clause of round 1's finding 4 was never removed: the handoff — the file the next
session reads first — still recommends adding a Bash permission rule for `agy`, and now points at a
log whose commands were deleted. The two logs were cleaned; the file with the most authority over
the next session was not.

On the specific questions asked: the widened list **cannot** refuse any file in this repo's 27
skills today (checked by `find`); the mode mask **cannot** strip an execute bit (`0o755` keeps
all three); `adopt`'s backup leaves nothing new readable, but its error messages name a prepared
copy the deferred cleanup has already deleted; and `GOTOOLCHAIN=local` has no download-side
supply-chain consequence at all — its consequence is the opposite one, judged below.

### On shipping v1.0.0 with GO-2026-4602 and GO-2025-3956 present

**Ship. Do not hold the release for the toolchain upgrade** — and not because "it is only a
local developer tool", which is not a security argument. Because both advisories are matched by
symbol and **unreachable by precondition** in this module. I read the advisory records
themselves (`https://vuln.go.dev/ID/GO-2026-4602.json`, `.../GO-2025-3956.json`), not the
titles:

- **GO-2026-4602 / CVE-2026-27139** — *"when listing the contents of a directory using
  File.ReadDir or File.Readdir the returned FileInfo could reference a file outside of **the Root
  in which the File was opened**… limited to reading metadata provided by lstat."* The escape is
  an escape **from an `os.Root`**. This module constructs none: `grep -rn "os\.Root\|OpenRoot\|os\.DirFS"`
  over every `.go` file returns nothing. The three `os.ReadDir` call sites
  (`inventory.go:125`, `inventory.go:164`, `apply.go:290`) pass plain path strings with no
  containment boundary to breach, and `LoadSkills` already leaves the tree deliberately through
  `filepath.EvalSymlinks`. There is no confinement here for the bug to break.
- **GO-2025-3956 / CVE-2025-47906** — *"passing certain strings to LookPath (`""`, `"."`, and
  `".."`) can result in the binaries listed in the PATH being unexpectedly returned."* The
  trigger is the **argument value**. Every `LookPath` in non-test code is a literal:
  `exec.LookPath("go")` at `cmd/sh-agent/main.go:129` and `internal/harness/repo.go:154`. Every
  `exec.Command` name is a literal (`"diff"`, `"git"`, `"go"`) or `bin`, an absolute path. No
  call site can produce `""`, `"."` or `".."`, so the vulnerable branch is not reachable.

So the residual risk of shipping is: a consumer binary built on a go1.24.x host carries two
stdlib defects neither of which this code can trigger. Against that, the pin's measured cost was
that the harness could not be **built at all** on an offline host with an older Go — the outer
`go run` is itself subject to the directive, so the `Reinstall` fix never gets a chance. Withdrawing
the pin is the right trade and I endorse it: an unbuildable tool is a worse security property
than two unreachable symbol matches, because the tool is how every later fix is delivered.

Two conditions on that endorsement, neither blocking:

1. **`govulncheck ./...` now reports 2 affecting again** (run here, output quoted under *Checked*).
   Any future gate that keys on "govulncheck clean" will fail from now on. Record the exception
   in the trail with the reachability argument above, or a later session reads the red output,
   re-adds the pin, and re-breaks the offline build. Note also *why* the number moved, measured rather than
   assumed: with the pin withdrawn, `go version` in this checkout is **go1.24.4** — the system Go
   (`/usr/bin/go`, Debian's) — and `go env GOVERSION` agrees. go1.25.8 ran here only because the
   pin made the `go` command re-exec a toolchain it had downloaded. So govulncheck is reporting the
   toolchain that actually compiles, and the earlier "0 affecting" was a true measurement of a build
   that existed only because of the pin's download. Both numbers are honest; they describe different
   builds.
2. **The operator's system Go upgrade stays the only real fix** and belongs in the handoff's Next
   Steps, not in a plan file waiting on this release. Debian trixie shipping no patched package
   makes it a manual toolchain install; that is an operator action, correctly outside the gate.

### Findings

- **[medium / confidence: high]** `cmd/sh-agent/main.go:237-241` and `internal/harness/apply.go:79-82`
  — **the `StatusUnknown` fix is wired into 2 of its 4 call sites.** `update` and `Pull` now fail
  closed, which is the deploy decision and the important half. But `status` still computes
  `dirty := "clean"; if info.Dirty { dirty = "uncommitted changes" }`, so a tree whose state git
  could not report prints as **`(main @ abc1234, clean)`**; and `Install` still writes
  `if commit != "" && info.Dirty { commit += "+dirty" }`, so the manifest records the bare commit
  as the provenance of files that came from a working tree of unknown state. Reachable without any
  attacker, **measured in the probe repo**: with `.git/index` unreadable (`chmod 000`)
  `git rev-parse --short HEAD` exits 0 while `git status --porcelain --untracked-files=all` exits
  **128** (`fatal: .git/index: index file open failed`), and with a truncated/corrupt `.git/index`
  the same pair gives 0 and **128** (`index file smaller than expected`) — a killed git, a full
  disk, an unclean shutdown. So `Commit` is populated (the "plain copy" branch is not taken) and the
  tree's real state is unknown. A stale `.git/index.lock`, which I expected to do this, **does not**
  — git 2.x still answers `status` read-only with the lock present, exit 0; probed and negative. The harness's own `AGENTS.md` makes
  `sh-agent status` the check the operator reads before calling an install done, and BUG-046's
  whole lesson is that a wrong "clean" is what makes a bad deploy invisible. Remediation, two
  lines, no new mechanism: `if info.StatusUnknown { dirty = "status unknown — git could not report it" }`
  in `status`, and `commit += "+unknown"` in `Install`, mirroring the existing `+dirty`.

- **[medium / confidence: high]** `internal/harness/adopt.go:128-158` (`copyTree`) and `.gitignore:9-15`
  versus `internal/harness/plan.go:213-225` — **the widened credential list guards the exit and not
  the entrance, and git's own guard was not widened with it.** `Plan` refuses to deploy the new
  names; `Adopt` copies a foreign skill directory — including whatever `npx skills add -g` fetched
  from the internet — into the repo's **git-tracked** `skills/<name>/` with no such refusal. The
  repo's `.gitignore` lists only `.env`, `.env.*`, `*.pem`, `*.key`, `*.tfvars`, `*.tfstate`: it
  does **not** list `id_rsa`, `id_ed25519`, `id_ecdsa`, `id_dsa`, `.netrc`, `.npmrc`, `.pypirc`,
  `.envrc`, `*.env`, `*.p12`, `*.pfx`, `*.jks`, `*.keystore`, `*.tfstate.backup`,
  `*credentials*.json` or `*-key.json` — that is, exactly the shapes this change just taught the
  installer to fear, plus two that were there before. So the path is: adopt → file sits in
  `skills/<name>/` → `git add skills/<name>` (the documented follow-up to adopt) → pushed to a
  GitHub remote → and only *then* does the install-time refusal fire, on every consumer, after the
  secret is public. The comment above `secretShaped` says the names "match the context's Git rule";
  the file that mechanises that rule does not. Remediation, no code design question: call
  `secretShaped(d.Name())` in `copyTree`'s walk so adopt refuses at ingest with the same message,
  and mirror the full pattern list into `.gitignore`. Both are the stdlib/`filepath.Match` already
  present; no dependency.

- **[medium / confidence: high]** `internal/harness/apply.go:36-47` (`Install`) with
  `internal/harness/plan.go:141` — **the mode mask never reaches a file that is already deployed.**
  `Install` short-circuits on `readErr == nil && bytes.Equal(existing, f.Data)` and counts
  "unchanged"; it does not chmod. `writeAtomic`'s `tmp.Chmod(mode)` runs only on the write path. A
  machine that already has `~/.claude/skills/x/SKILL.md` at `0666` from an install before this fix
  therefore keeps `0666` through every future `sh-agent install` and `update` for as long as that
  file's bytes do not change — which, for a stable skill, is indefinitely. The new
  `TestPlanDropsGroupAndOtherWrite` asserts `Plan`'s output, so it is green while the deployed file
  on the affected machine is not. Exploitable in principle exactly as round 1 stated it: another
  local account rewrites instructions every session reads, without touching the repo. Remediation:
  on the unchanged branch, `os.Stat(dest)` and `os.Chmod(dest, f.Mode)` when `info.Mode().Perm() != f.Mode.Perm()`,
  counting it as written or as its own "tightened" line so the operator sees it happen.

- **[low / confidence: high]** `internal/harness/adopt.go:83` versus `:93-94` and the error strings
  at `:101`, `:109`, `:112` — **the failure messages name a prepared copy that `defer os.RemoveAll(staging)`
  deletes before the caller reads them.** The comment states the invariant in as many words —
  *"staging is never deleted on a failure: whatever else goes wrong, the prepared copy survives and
  is named"* — and the defer registered at line 83 runs on every return, including all three
  failure returns. The substance of the release reviewer's finding 3 **is** fixed (the original is
  moved aside and renamed back, so it is never lost); what is wrong is the promise. The new test
  asserts only `strings.Contains(err.Error(), ".agents/scratch")` — that the message names a path,
  never that the path exists — so the test enshrines the false claim. This is the "description that
  outlives its behaviour" case from the global rules, in a message the operator is meant to act on
  after a failed replace. Remediation: keep the copy (`keep := false` … `defer func(){ if !keep { os.RemoveAll(staging) } }()`)
  or delete the promise from the comment and the three messages.

- **[low / confidence: medium]** `internal/harness/plan.go:120-122, 213-225` — **the refusal is
  whole-plan fatal and its trigger is now much wider.** One matching filename anywhere in any skill
  makes `Plan` return an error, so `install`, `status` and `diff` all fail for **every** vendor and
  every file, not just the offending one. Fail-closed is right for a credential; the newly added
  `.npmrc`, `.envrc`, `.pypirc`, `.netrc` and `*.env` are *configuration* shapes that appear in
  legitimate skill scaffolding and templates, and skills are vendored verbatim from third parties
  (`skills/VENDORED.md`). A future `npx skills add` of an upstream skill that ships a template
  `.npmrc` would break `sh-agent update` on every box running the hourly timer at once, until a
  human intervenes. Checked and clean today: `find skills -type f` matching all fourteen new
  patterns returns **nothing** across the 27 skills, and no skill text even mentions those names,
  so there is **no false positive in this release**. Remediation if it ever bites: keep the hard
  refusal for the unambiguous credential shapes and skip-with-a-loud-line for the rc-file shapes,
  or name the offending file in a `--force`-able error.

- **[low / confidence: medium]** `internal/harness/repo.go:174` with `cmd/sh-agent/main.go:258-275`
  — **`GOTOOLCHAIN=local` moves the choice of stdlib to the consumer's host and nothing records
  which toolchain was used.** The fix is right and its download-side supply chain is *better*, not
  worse: with the pin withdrawn nothing is fetched from the toolchain endpoint at all, and with the
  pin still present `GOTOOLCHAIN=local` would have suppressed the only fetch. The residue is
  provenance: each consumer's binary is built with whatever `go` that box has (≥ the `go 1.24.4`
  directive, or the build fails — see the next finding), and `sh-agent status`'s `binary` line
  prints only the VCS revision from `versionFromInfo`. `debug.BuildInfo.GoVersion` is already in
  the `bi` the line reads. Remediation, one field, no dependency: print it, so
  "which Go built the thing that installs my agent context" is answerable on the box.

- **[low / confidence: medium]** `internal/harness/repo.go:174` — **`GOTOOLCHAIN=local` also removes
  the escape hatch for a host *below* the `go` directive.** With `auto`, a box on go1.23 built this
  module by downloading a toolchain; with `local` it cannot, and `go build` fails with
  `go.mod requires go >= 1.24.4`. That is a deliberate and correct trade (the alternative is a
  compiler download on every install), but it is a *new* failure for a host the old behaviour
  served, and neither `deploy/remote.md` nor the `Reinstall` comment states the resulting floor.
  Remediation: one sentence wherever the install prerequisites are listed — "Go ≥ 1.24.4 must be
  installed on the host; sh-agent never downloads a toolchain."

- **[low / confidence: low]** `internal/harness/plan.go:111` — `d.IsDir()` returns before
  `secretShaped(d.Name())`, so a **directory** named like a credential (`secrets.env/`,
  `credentials.json/`) is descended into and its contents deployed under a name the refusal was
  meant to stop. Contrived, and the files inside are each checked on their own names; recorded for
  completeness because the list-widening invites exactly this shape of miss.

- **[low / confidence: medium]** `internal/harness/adopt.go:157` (`copyTree`'s
  `info.Mode().Perm()|0o600`) — round 1's mode finding named `plan.go` **and** `adopt.go`; only
  `plan.go` was masked. A `0666` file adopted from `~/.agents/skills` lands `0666` inside the
  repo's `skills/`, where it stays until someone re-clones (git tracks only the execute bit), and
  where any local account can rewrite it between `adopt` and the next `install` — at which point
  `Plan` masks the *deployed* copy but faithfully ships whatever bytes that account wrote.
  Remediation: the same `&0o755 | 0o600` mask in `copyTree`.

### Fix-by-fix judgement

| # | Fix | Closes it? |
|---|---|---|
| 1 | `secretShaped` widened (14 new patterns) | **Yes** for the deploy boundary — `id_ed25519`, `id_rsa`, `id_ecdsa`, `id_dsa`, `.netrc`, `.npmrc`, `.pypirc`, `.envrc`, `*.env`, `*.p12`, `*.pfx`, `*.jks`, `*.keystore`, `*.tfstate.backup` all match via `filepath.Match`. No false positive in the 27 skills today. **Gaps:** the ingest boundary and `.gitignore` (finding 2), whole-plan fatality (finding 5), case sensitivity (already named residual — `ID_RSA`, `server.KEY` still pass; `.PEM` too). |
| 2 | `Perm()&0o755 \| 0o600` | **Partly.** New writes are correct — `0664→0644`, `0666→0644`, `0777→0755`, `0600→0600` — and **no execute bit is stripped**: `0o755` retains owner, group and other execute, so a script skill still ships executable. But already-deployed files are never chmodded (finding 3). |
| 3 | release block: `--untracked-files=all`, `':(top)'`, `rev-parse` both sides | **Yes.** Re-read in `deploy/remote.md` and ADR-026; both assertions are now cwd-independent and untracked-visible, and the gated head is derived rather than transcribed. |
| 4 | permission-widening recipes removed from the two testing logs | **Mostly.** Both `command(*)` wildcard blocks and the `--dangerously-skip-permissions` pair are gone from the two logs, replaced by what the check needs plus an explicit statement that the permission is the operator's decision. **But** the `agy` Bash-permission-rule suggestion survives verbatim in `handoffs/2026-09-16-board-cleared.md:72` (medium finding above), and the replacement paragraph adds a pointer to the 2026-09-03 log, which still lists the five permission strings (low finding above). |
| 5 | `Reinstall` builds with `GOTOOLCHAIN=local` | **Yes**, and the `append(os.Environ(), …)` override is correct — `exec.Cmd` de-duplicates keeping the last entry, so an inherited `GOTOOLCHAIN=auto` loses. Residues in findings 6 and 7. |
| 6 | `adopt` keeps a restorable backup | **Yes** for the data loss; the messages are wrong (finding 4). |
| 7 | `Inspect` → `StatusUnknown`, `update`/`Pull` fail closed | **Partly** (finding 1). |

### Two more findings on fix #4 (the permission-widening removal)

- **[medium / confidence: high]** `.agents/docs/handoffs/2026-09-16-board-cleared.md:72` — **the sentence
  round 1 asked to drop survives verbatim in the handoff, which is the file the next session reads
  first.** The testing logs were cleaned; this was not: *"The two `!` commands that close it are in
  `testing/2026-09-16-verify-real-home-fcc8539.md`; **a Bash permission rule for `agy` would let a
  future session run them directly**."* Round 1's Outcome records this fix as done. It is done in
  two of the three places. The global resume order puts the newest handoff ahead of the dated
  testing logs, so the one file that still carries the suggestion is the one with the most
  authority over the next session's first move — and it recommends exactly the action this
  session's permission classifier refused three times. Its pointer is also now stale: the two `!`
  commands it names were removed from that log, so a future agent follows it, finds nothing, and is
  left with only the permission-rule suggestion. Remediation: strike the clause after the semicolon
  and repoint the sentence at the log's new paragraph, in the same edit.

- **[low / confidence: high]** `.agents/docs/testing/2026-09-16-verify-real-home-9bb7ce0.md`, the
  replacement paragraph — the new text ends *"The 2026-09-03 procedure in this directory records how
  it was done once before, including restoring the file afterwards from a backup taken first."*
  `.agents/docs/testing/2026-09-03-verify-real-home.md:11` still spells out the five permission
  strings (`command(*)`, `read_file(*)`, `write_file(*)`, `edit_file(*)`, `list_dir(*)`) in prose.
  The block is gone; a signpost to where it still lives was added in its place. Weaker than the
  copy-paste block it replaced — past tense, no JSON, no commands — but it converts a dated record
  into a pointer, which is the direction the finding asked not to go. Remediation: drop the last
  sentence, or say what the earlier log shows *without* naming it as the place to look.
  Checked and **acceptable** by contrast: `.agents/docs/testing/skill-matrix.md:14` also carries the
  `command(*)` set, for a **sandbox** `HOME` under `.agents/scratch/`, never the operator's real
  home — a different action with a different blast radius, and correctly left alone.

### Live-fix status — the tree moved during this review

Two of the findings above were fixed in the working tree while I was writing, and the corrections
are correct. Recorded rather than deleted, because the record of what was found is the point:

- **`adopt`'s `defer os.RemoveAll(staging)`** (finding 4) — **fixed** in the working tree: the defer
  is gone, replaced by an explicit `os.RemoveAll(staging)` on the copy-failure path only (where
  nothing has been moved yet) and a `hadPrevious` flag that picks the right message. Re-read at
  `internal/harness/adopt.go:83-122`; every failure path now names a path that exists. One residue,
  **[low / confidence: high]**: a failed adopt now leaves a full copy of the source skill in
  `.agents/scratch/adopt-<name>-<pid>` indefinitely, and `copyTree` still applies no `secretShaped`
  refusal, so a credential inside an adopted third-party skill is now *persisted* there as well as
  landing in `skills/`. Same remediation as the medium finding above: refuse at ingest.
- **`sh-agent status` rendering `StatusUnknown`** (half of finding 1) — **fixed** in the working
  tree: `dirty = "state unknown — git could not report it"`. The other half stands:
  `internal/harness/apply.go:79-82` still writes the bare commit into the manifest, with no
  `+unknown` beside the existing `+dirty`. That manifest is the provenance record `sh-agent status`
  reads back.
- `go.mod`'s `toolchain` withdrawal and ADR-029, and the `AGENTS.md`/`README.md` lines describing
  it, are read and judged above; the ADR is honest about all three mismeasurements and names the
  third corner (`go run` under `auto` + `GOPROXY=off`) that only the last measurement covered.

### Proof of the mode finding (finding 3), end to end on the real install path

Run against a scratch home, never the operator's own (`--home .agents/scratch/security-v1.0.0-r2/home`):

```
install #1 exit 0 → ~/.claude/skills/bash/SKILL.md mode 644
chmod 666 that file
install #2 exit 0 → "0 written, 260 unchanged, 0 removed, 0 kept — 260 files managed"
mode after install #2: 666
find <home>/.claude <home>/.gemini -type f -perm /022 → .claude/skills/bash/SKILL.md
```

So the mask protects the first write and nothing else. A machine that already carries a loose mode
keeps it through every install and update, and the file stays group- and world-writable in the
directory whose contents every session reads as instructions.

## Checked (round 2)

- **Secrets** — `git diff c7e3688..787e15c` plus the uncommitted `go.mod`/docs changes scanned for
  private-key armor, `ssh-`/`AKIA`/`ghp_`/`github_pat_`/`xox`/`sk-`/`AIza`/`Bearer` and
  password/token assignments: none. The one long token in the range is the same ed25519 **public**
  key round 1 verified. `cmd/sh-agent/main_test.go`'s `"SECRET"` fixture changed shape — it now
  asserts the problem string belongs to a closed set rather than merely "does not contain SECRET",
  which is a **stronger** leak test and the right correction; the fixture is still a test of a leak,
  not a leak. The scratch home created for the mode probe holds a deployed copy of the harness and
  no credential.
- **Injection** — nothing new parses external input. `secretShaped` widened its `filepath.Match`
  pattern list; `filepath.Match` takes no regex metacharacters and every pattern is a literal, so no
  pattern injection. `Reinstall`'s new `build.Env` is an env assignment, not a shell; verified
  empirically that `exec.Cmd` de-duplicates keeping the **last** entry — a probe program that set
  `GOTOOLCHAIN=auto` in its own environment and appended `GOTOOLCHAIN=local` saw the child print
  `local`. All `exec.Command`/`LookPath` names remain literals.
- **Authz/authn** — the release procedure's three corrected assertions re-read in
  `deploy/remote.md` and ADR-026; both sides of the parent check now go through `git rev-parse`
  (a non-existent `<gated head>` makes `rev-parse` exit 128 and the `test` fail, so it fails
  closed), `--untracked-files=all` overrides `status.showUntrackedFiles`, and `':(top)'` makes the
  source-tree assertion cwd-independent as round 1's probe required. No change to the trust chain.
- **Crypto** — **nothing hand-rolled**, again. This round's diff implements or composes no cipher,
  hash, MAC, signature, KDF, RNG, token or session format and no auth flow; it adds a filename
  pattern list, a permission mask, an env var, a rename-with-backup and a boolean. Signing stays
  entirely inside `git` + OpenSSH. No `crypto/*` import in the module at all.
- **Dependencies** — `go.mod` still `require`s nothing; the withdrawal *removes* a directive rather
  than adding anything. `govulncheck ./...` run here, **2 affecting**: `GO-2026-4602` (`os`, via
  `LoadSkills` → `os.ReadDir`, `inventory.go:164`) and `GO-2025-3956` (`os/exec`, via `Reinstall` →
  `exec.LookPath`, `repo.go:154`), both matched against `os@go1.24.4`, which is this host's actual
  toolchain (`go version` → go1.24.4, `/usr/bin/go` → go1.24.4, `go env GOVERSION` → go1.24.4).
  Both judged unreachable by precondition above, from the advisory JSON rather than the titles.
  `npm audit` not applicable: no `package.json`, no JavaScript, in the repo or the range.
- **GCP surface** — untouched by this round's diff: no IAM, service account, Cloud Run flag,
  Terraform or deploy config. The only deploy-adjacent edits are the release-block assertions.
- **The two behaviour questions the dispatch asked** — (a) *can the widened list refuse a legitimate
  skill file?* Not in this release: `find skills -type f` against all fourteen new patterns plus the
  seven old ones returns nothing across 27 skills, and no skill body even mentions those filenames.
  The exposure is future and third-party, recorded as a low finding. (b) *can the mask strip a
  needed execute bit?* No: `0o755` keeps owner, group and other execute; only `0o022` is dropped.
  Verified arithmetically for `0600/0640/0644/0664/0666/0700/0755/0777` and observed on 260 deployed
  files, of which the only `perm /022` hit was the one I chmodded by hand.

## Not examined (round 2)

- **`internal/harness/harness_test.go`'s 113 added lines were read in the diff, not run.** I did not
  execute `go test -race ./...`, `go vet` or `gofmt` this round — that is the verifier's row, and
  the tree was being edited while I read, so a green run here would not have described the tree the
  gate will see. `TestPlanDropsGroupAndOtherWrite` I read and judged: it asserts `Plan`'s output and
  therefore cannot catch the `Install` gap proven above.
- **The uncommitted working tree is a moving target.** Between my first and last read, `adopt.go`,
  `repo.go`, `cmd/sh-agent/main.go`, `main_test.go`, `harness_test.go`, `go.mod`, `AGENTS.md`,
  `README.md` and three trail files all changed under me. Everything above was re-checked against
  the tree as of my last read; anything committed after that is unreviewed. **Nothing in this
  section has been re-verified against a commit**, because there is none yet — a re-gate should
  confirm the final commit matches what is described here.
- **The consumer side remains unread** — crew's channel, `/etc/styrman/allowed_signers`,
  argolis-devbox's setup script. Unchanged from round 1.
- **Not measured: a host with Go older than the `go 1.24.4` directive.** The floor finding was
  measured with a throwaway module declaring `go 1.99.0`: `GOTOOLCHAIN=local` exits 1 with
  `go.mod requires go >= 1.99.0 (running go 1.24.4; GOTOOLCHAIN=local)`, and `GOTOOLCHAIN=auto` with
  `GOPROXY=off` exits 1 with `toolchain not available`. That demonstrates the two failure shapes but
  is not a real go1.23 host, and I did not build the harness on one.
- **The release procedure was not executed** and no tag was made, signed or pushed. Probes ran only
  in `.agents/scratch/security-v1.0.0-r2/probe` (`git init` exit status read as 0; every later
  command `git -C <probe>`). For the record: one of those probe commands was a
  `git -C <probe> checkout-index -f -a` to undo my own probe edit **inside the probe repo only** —
  the real checkout received no tree-changing git command, and `git tag` in it is still empty.
- **`core.fsmonitor` as a way to make a dirty tree read clean: probed and negative.** A broken hook
  path and a hook that answers "nothing changed" both left
  `git status --porcelain --untracked-files=all` reporting ` M` and `??` lines with exit 0. The
  probe was validated in both states — the baseline showed the dirty lines, so the instrument could
  emit. No finding.
- Not re-read this round: the whole `4aa19be..c7e3688` range (round 1's scope), the 27 skill bodies,
  and the trail files outside the response diff except the four I grepped for permission recipes.

## Outcome — round 2

- **The advisory judgement is accepted as the basis for shipping, and recorded here so no later session silently re-adds the pin.** `GO-2026-4602` is an escape from an `os.Root`, and this module constructs none — no `os.Root`, `os.OpenRoot` or `os.DirFS` anywhere in it. `GO-2025-3956` triggers when `LookPath` is passed `""`, `"."` or `".."`, and the only call in this code passes the literal `"go"`. Both are symbol matches that `govulncheck` reports as affecting and neither is reachable in this program. Against that, the pin made the tool unbuildable offline, and an unbuildable tool is the worse property because it is how every later fix ships. v1.0.0 ships with both advisories listed and unreachable; the standing operator toolchain upgrade remains the way to clear the report itself.
- [medium] the mode mask never reached an already-deployed file — **fixed and reproduced end to end**: `Install`'s unchanged branch now corrects the mode in place, a failed chmod is a note rather than an install failure, and `--dry-run` changes nothing. The review's own scenario, run again: install → 644, `chmod 666`, install → "0 written, 260 unchanged", mode back to **644**. `TestInstallCorrectsTheModeOfAnUnchangedFile` covers it, including the dry-run leg.
- [medium] `secretShaped` guarded the exit but not the entrance — **fixed**: `copyTree` refuses a secret-shaped name as it copies, so `adopt` stops a credential before it reaches `skills/` rather than after it has been committed and pushed, and `TestAdoptRefusesSecretShapedFiles` asserts nothing lands. `.gitignore` now mirrors all twenty-one patterns with a comment naming `secretShaped` as the thing it is kept in step with.
- [medium] the handoff still carried the "a Bash permission rule for `agy`" clause after the testing logs were cleaned — **fixed**: that sentence is replaced by what the checks need, whose decision it is, and why no recipe is recorded. This was round 1's finding 4 surviving in the one file the next session reads first, which is exactly where it would have done the most harm.
- [medium] the manifest recorded a bare commit for a tree git could not describe — **fixed**: `+status-unknown`, alongside the `status` output fix.
- The review's four direct questions, answered by its own measurements and confirmed here: the widened list refuses no file in the 27 skills today; the mode mask strips no execute bit; `adopt`'s backup leaves nothing newly readable; `GOTOOLCHAIN=local` has no download-side supply-chain cost.

**Not re-reviewed:** no third security pass ran over these fixes. What replaced it is named above — the end-to-end reproduction of the mode scenario, and the new tests.
