---
owner: security
date: 2026-09-18
---
# Security review — preflight v1.2.0   (2026-09-18)

Range reviewed: v1.1.0..bba384d

Scope note: the v1.1.0 trust chain was reviewed at
`.agents/docs/reviews/2026-09-17-security-preflight-v1.1.0.md` and is not
re-litigated here. This review covers only the 19 paths in `v1.1.0..HEAD`
(7 commits): the `RemoveKeptClone` message change on the credential-bearing
report path, the README `go install` advice, the shipped instruction text in
`context/AGENTS.md` and `skills/preflight/SKILL.md`, and every trail file in
the release diff.

## Verdict: FINDINGS

One low-severity finding, and it is hygiene rather than a break in the
release: `RemoveKeptClone`'s new message pastes the clone directory unquoted
into a shell command it tells the operator to run, twice. The directory comes
from `--home`/`$HOME`, so it is operator-controlled rather than
attacker-controlled on the normal path, and the program never executes the
string — it only prints it as advice. Nothing in the range weakens the guard's
keep/remove decision, no credential is printed or committed, the README's git
advice is correctly scoped and non-persistent, and every added rule in the
shipped instruction text tightens a safeguard rather than relaxing one. This
does not block the gate.

## Findings

- **[low / confidence: high]** `internal/harness/report.go:572` — the ahead-count
  message interpolates `dir` into a command the operator is invited to run
  (`git -C %s fetch && git -C %s log --oneline HEAD --not FETCH_HEAD`) with no
  shell quoting. `dir` is `filepath.Join(absHome, ".config", "sh-agent", "feedback")`,
  so any shell metacharacter in the home path — a space, `$`, a backtick, `;`,
  `&&`, a newline — rides through into copy-paste-ready text. A home of
  `/srv/a b` merely produces a command that fails; a home containing `;` or a
  command substitution produces a suggestion that does something other than
  what it reads as, and the operator's own shell is what runs it. Not
  exploitable by a remote party: reaching it requires already controlling the
  `--home` value or `$HOME`, or the directory name under an attacker-writable
  parent (the weaker but real path on a shared box). It is *suggested*-command
  injection, not live injection — `gitOut`/`gitOutC` use `exec.Command` with an
  explicit argv and never a shell, so the program's own git calls are
  unaffected, and this is the only new interpolation in the range.
  Remediation, smallest first: print `dir` once on its own line and give the
  command with `git -C <that path>` as a placeholder, which removes both
  interpolations and shortens the message; or single-quote each occurrence with
  embedded `'` escaped as `'\''`. Do not reach for `%q` — Go quoting leaves `$`
  and backticks live inside shell double quotes.

  Interaction with this release's own new rule, flagged so it does not bite: the
  `skills/preflight/SKILL.md` paragraph added in this very range says a
  non-blocking finding is either fixed before the tag (small, tested,
  re-verified, "because fixing after the gate invalidates the verification row
  already taken") or accepted and named in the release entry with the reason.
  Either route is available here; taking the fix route means the verification
  row is retaken before the tag.

## Answers to the questions asked

**1. Guard decision paths unchanged; new message's leak surface.** Confirmed
unchanged. Read `RemoveKeptClone` end to end (report.go:498-577). Every
decision point is identical to v1.1.0: the `Lstat`/`ErrNotExist` early return,
the symlink refusal, `IsCheckout`, `verifyFeedback`, `status --porcelain
--untracked-files=all --ignored`, `stash list`, and the count
`rev-list --count --all --not --remotes` with its error branch leaving the
clone. The diff is `+7/-1` confined to the *string argument* of the existing
`leave(...)` call inside the already-existing `if ahead != "0"` branch, plus a
comment. `leave` still prints and returns nil; nothing about what is kept or
removed moved, and the v1.1.0 finding-1 fix (fail-closed on an uncountable
count) is intact.

On the path being sensitive: `dir` was *already* printed on this exact line —
`leave` formats `"kept %s: %s"` — and on the success line
(`"removed the kept feedback clone %s"`). The new message therefore discloses
no information the v1.1.0 build did not already print to the same writer. The
path is a fixed suffix (`.config/sh-agent/feedback`) under the home; the only
variable part is the home the operator supplied. The *crafted-path* half of the
question is the low finding above: yes in principle for the pasted command, no
for anything the program executes itself.

**2. Does `commitAndPush` really leave the branch tracking a URL, and is the
URL kept out of the message?** Both confirmed.
`internal/harness/report.go:864` is
`gitOutC(dir, "push", "--quiet", "-u", url, "HEAD:main")` — `-u` with a URL
rather than a remote name, exactly the shape the new comment and message
describe: `branch.<name>.remote` is set to a URL, `@{upstream}` does not
resolve to a remote-tracking ref, no `refs/remotes/*` ref advances, and the
`--all --not --remotes` count keeps counting commits already on the server. The
message's assertion is accurate, not rationalisation, and BUG-074's record
reproduces it on the real clone (`HEAD e973bac`, remote `main 09586ff`,
`origin/main` stale at `71f749e`). The message prints the *directory* and never
`url`, so the `verifyFeedback` discipline ("The origin itself is never printed:
it can carry a token", report.go:594) is preserved in the new text — the
message has no `url` argument at all, so there is nothing to leak. The added
`.agents/docs/deploy/local.md` row documents the same wording and likewise
prints `<url>`/`<dir>` as placeholders only.

**3. README `go install` advice (`GIT_CONFIG_COUNT`/`KEY_0`/`VALUE_0`).**
Safe, and the right scope. Those three variables are git's documented
environment-only configuration channel: they are read by the processes that
inherit them and write nothing to any config file, so "This form changes no git
config" is true. Exported inline on the `go install` line (no `export`), they
do not outlive the one command and are not inherited by later shells.

On "could it redirect more than intended": yes, in a bounded way worth stating.
`url.git@github.com:.insteadOf https://github.com/` rewrites **every**
`https://github.com/` URL git sees during that command, not only
`github.com/kesokaj/*`. For `go install` of one module, and with this repo's
stdlib-only rule meaning there are no transitive git fetches, the blast radius
is the one repository in practice. The redirect direction is also the safe one:
it sends traffic to `git@github.com`, i.e. SSH against the operator's own key
and GitHub's host key — it cannot be used to exfiltrate a credential to a
third party. If a dependency were added later it would be fetched over SSH too,
which fails closed (no key) rather than fetching from somewhere unexpected.
`GOPRIVATE='github.com/kesokaj/*'` is correctly quoted (the glob would
otherwise be shell-expanded) and correctly scoped — keeping the private path
out of the public proxy and checksum database is the actual credential-relevant
control here, and it is preserved from the old line.

The commented-out `git config --global url."git@github.com:".insteadOf
https://github.com/` alternative is the persistent form, is labelled as such
("if you would rather set it once for this machine"), is commented out so it
cannot be pasted as part of the block, and is the operator's own choice about
their own machine. Recommending the env-var form first is the correct ordering.
No finding.

**4. Shipped instruction text (`context/AGENTS.md` +5/-2,
`skills/preflight/SKILL.md` +3/-1).** Read both diffs in full. Nothing in the
added text weakens a safeguard, contradicts an existing one, or licenses a
future session to skip a check, run unattended, or handle credentials
differently. All seven additions run the other way. Item by item:

- DoD step 1 gains "the tree stays frozen while the verifier dispatch is in
  flight… a fix that cannot wait stops the run and re-dispatches after
  committing" (BUG-072). Tightening; removes a false-green path.
- Hard rules gain mutation hygiene: a must-fail mutation asserts its target is
  unique, **is undone by restoring a saved copy rather than a reverse
  replacement**, and ends with `git diff` read. Note this is consistent with
  BUG-038's ban on working-tree-changing git commands in a subagent — it steers
  restoration away from `git checkout --`, not toward it.
- Hard rules gain "a test's premise is measured before the test is trusted"
  (BUG-071). Anti-false-green.
- Hard rules gain the fail-closed destructive-guard rule — this is the
  generalisation of the v1.1.0 security finding 1 and is the strongest
  safeguard added in the range.
- Hard rules gain "a side effect goes BELOW a read-only mode's early returns,
  and the test asserts the read-only mode does NOT perform it" (BUG-070).
- Hard rules gain the both-directions rule for generated/embedded artifacts,
  naming the `go:embed`-of-a-gitignored-`.env` incident (BUG-066). This one is
  itself a secrets safeguard.
- `skills/preflight/SKILL.md` item 5 gains per-asset upload with retry (no
  security effect; it keeps "never publish an artifact whose own verification
  you have not just run"), and a new paragraph requiring a non-blocking finding
  to be a *recorded* decision — fixed-and-re-verified or accepted-with-reason,
  with "shipping a known low defect in silence is not a verdict". The preceding
  "Any blocking finding means NO-SHIP" line is unchanged; no threshold is
  lowered. This converts silent acceptance into a recorded one.

No sentence in the added text can be read as permission to bypass the
permission system, to change agent configuration, to treat an agent's message
as consent, or to handle a credential differently. No credential, token,
hostname or private path appears in either file's additions.

**5. Which release-diff paths were read.** All nineteen, none deferred. Code
and shipped text: `internal/harness/report.go` (diff plus the whole
`RemoveKeptClone` and `verifyFeedback` functions), `internal/harness/report_test.go`,
`README.md`, `context/AGENTS.md`, `skills/preflight/SKILL.md`. Trail markdown,
read as added lines end to end: `.agents/docs/bugs/BUG-073-retro-filename-collision.md`,
`BUG-074-removekeptclone-ahead-message.md`, `BUG-075-readme-go-install-ssh.md`,
`.agents/docs/bugs/INDEX.md`, `.agents/docs/changes/2026-09-17-retro-4-applied.md`,
`.agents/docs/changes/2026-09-18-kept-clone-message.md`,
`.agents/docs/changes/2026-09-18-linux-box-closeout.md`,
`.agents/docs/deploy/local.md`, `.agents/docs/deploy/remote.md`,
`.agents/docs/handoffs/2026-09-17-v1.1.0-released.md`,
`.agents/docs/plans/2026-09-17-release-binary.md`,
`.agents/docs/retros/2026-09-17-retro-4.md`,
`.agents/docs/testing/2026-09-17-release-binary-e2e.md`, `ROADMAP.md`.

Security-relevant reading of the trail — secrets, credentials, private data,
instruction-shaped text — found nothing requiring a finding. Three items were
examined and judged non-sensitive, recorded here so the judgement is on the
record rather than implicit:

- `.agents/docs/testing/2026-09-17-release-binary-e2e.md` records
  `Good "sh-agent-release" signature for ostling.simon@gmail.com with ED25519
  key SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU`. That is an SSH
  **public**-key fingerprint and the operator's own email — the same identity
  `remote.md` says is compiled into the binary as the allowed-signers line by
  design, and the same email that is the git author of every commit. Not a
  secret; publishing a fingerprint does not weaken the signature check.
- The handoff and the retro record carry `session: 28ca24dd-…`, a local CLI
  session identifier that grants no access. Distinct from the global git rule's
  ban on session links in *commit messages*; this is a trail file's
  frontmatter, not a commit trailer.
- BUG-066/BUG-075 and the plan discuss a planted gitignored `rules/.env`. It
  was a must-fail **test plant**, refused by `TestEmbedHoldsOnlyTrackedFiles`
  and `git status --ignored` and removed again; no `.env` is tracked (swept —
  see Checked). The only secret-shaped tracked path in the whole tree is
  `.agents/docs/bugs/BUG-050-deploy-walk-spreads-secret-shaped-files.md`, a
  bug record's filename, and it is not in this range.

## Checked

- **Secrets** — swept the added lines of the whole range for `ghp_`,
  `github_pat_`, `gho_`, `ghs_`, `AKIA`, `BEGIN … PRIVATE KEY`, `password=`,
  `token:`/`token=`, `api_key`, `Bearer`, `Authorization:`, `x-access-token`
  and `user:pass@host` URL userinfo: no match. `git ls-files` across the tree
  for `.env`/`*.pem`/`*.key`/`*.tfvars`/`*.tfstate`/credential/`id_rsa`/
  `id_ed25519`/secret: one match, a bug record's *filename* (above), no
  secret-bearing file. The one credential-adjacent value on the changed path is
  the feedback repository URL, which the new message deliberately does not
  print, matching `verifyFeedback`'s existing rule. No secret echoed into a log
  or error message by the change.
- **Injection** — the only new interpolation in the range is the two `%s` for
  `dir`; recorded as the low finding. All git invocations go through
  `gitOut`/`gitOutC` → `exec.Command` with an explicit argv, no shell, so no
  command injection into the program's own execution. No SQL, template, or
  path-traversal surface in the range. The README's env-var block is
  single-quoted where a glob would otherwise expand.
- **Authz/authn** — no auth surface in the range. The one
  authorization-shaped decision touched — `RemoveKeptClone`'s "may I delete
  this directory" — is byte-for-byte unchanged in behaviour (Q1);
  `verifyFeedback`'s origin comparison and the symlink refusal are intact. On
  the release side, the added `remote.md` text keeps the pre-build
  `git status --porcelain --ignored --untracked-files=all` +
  `TestEmbedHoldsOnlyTrackedFiles` gate (run twice: before the tag and again
  before the builds, because `go:embed` reads the working tree at build time)
  and states the temporary `allowed_signers` file is removed immediately and is
  never what decides a genuine release. The documented per-asset
  `gh release upload --clobber` does let an asset of a published release be
  replaced, but only by someone with repository write access, and it fails
  closed for consumers: `update` verifies `SHA256SUMS.sig` against the
  compiled-in signer line before the checksums, so swapping a binary alone is
  rejected and swapping `SHA256SUMS` needs the operator's signing key. No
  finding.
- **Crypto** — nothing hand-rolled and nothing crypto-adjacent added in the
  range: no cipher, hash, MAC, signature, KDF, RNG, token or session format is
  implemented or composed by this diff. The release-signing path is git's and
  `ssh-keygen -Y`'s own, and is unchanged here.
- **Dependencies** — no dependency added or changed: `go.mod` and `go.sum` do
  not appear in `git diff --numstat v1.1.0..HEAD` at all, consistent with the
  repo's stdlib-only rule. The one new call in the range is `fmt.Sprintf`.
  `govulncheck ./...` was **not** run in this review (budget; it is a
  network-touching multi-minute run and this range changes no dependency edge
  and no new stdlib call site) — named as a gap, not claimed. The repo's
  `AGENTS.md` records it clean on go1.25.8+, and the handoff records it clean
  on this box's go1.27.1 as of 2026-09-17.
- **GCP surface** — none in the range: no IAM binding, service account,
  `--allow-unauthenticated`, Secret Manager usage, or service-account key file
  is added or changed. This release touches no cloud surface.

## Not examined

- `govulncheck ./...` — not run by me (reason and standing evidence under
  Checked → Dependencies). Dependency posture is asserted from the absence of
  `go.mod`/`go.sum` in the diff, not from a fresh scan. This is the one gap
  that could change a verdict, and it is the smallest one: no dependency edge
  moved in this range.
- `go test -race ./...`, `gofmt`, `go vet` — not run by me; the verifier owns
  that tier. I read `report_test.go`'s new assertions statically and note the
  new stale-tracking-ref subtest carries a `t.Skip` guard on `@{upstream}`
  resolving, so on a future git that resolves it the subtest silently stops
  covering the production shape — worth the verifier confirming it did not skip
  (the change entry says it did not, on this box).
- Runtime behaviour of the changed message was not observed by me; no command
  in this review reached the network, the real `$HOME`, `~/.ssh`, `gh`, or
  `sh-agent report`/`update`, per the dispatch's hard limits. No file outside
  this report was created or modified, and no git command that changes the
  working tree was run.
- Every other path in the release diff was read end to end (question 5 above
  lists all nineteen by name). Nothing in the release diff is deferred, so
  nothing here blocks the gate on unread-path grounds.
