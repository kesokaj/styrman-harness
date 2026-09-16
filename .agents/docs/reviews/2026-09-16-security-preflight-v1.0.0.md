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
