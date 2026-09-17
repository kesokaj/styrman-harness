---
owner: security
date: 2026-09-17
---
# Security review — preflight v1.1.0 release gate   (2026-09-17)

Range reviewed: v1.0.0..086deab

## Verdict: FINDINGS (2 low, both ordering/hygiene in the release PROCEDURE; nothing blocks the gate)

No unreviewed code ships in this release, and no fix made after a prior review
introduced a regression: every fix the cited reports claim in their `## Outcome`
sections is present and whole at HEAD, verified by locating it in the shipping
source rather than by trusting the report. The worst thing found is an ORDERING
defect in the release procedure, not in the code: `.agents/docs/deploy/remote.md`
§ Releases puts the `go:embed` contamination guard AFTER `git push origin
vX.Y.Z`, and `skills/preflight/SKILL.md` Step 6 encodes the same order (6.4
pushes the tag, 6.5 builds artifacts). Because "a pushed tag is never moved", a
stray ignored file under the five embedded directories is caught only once the
version is public and unrecallable. The guard still stops the contaminated
binary from being PUBLISHED, so the blast radius is a burnt version number, not
a shipped secret — hence low. Ship is not blocked by either finding.

## Tree identity (verified here, not taken from the dispatch)

`git rev-parse 24d08a8^{tree} cd0bbcb^{tree}` → both
`da8d57d59eb1ff49ea3b96d2cdbb1183a5fb8b93`, and `git diff 24d08a8 cd0bbcb
--stat` is empty. The trees the prior reports reviewed and the tree shipping as
v1.1.0 are byte-identical. This premise is what licenses a delta review instead
of a repeat, so it was checked before anything else rather than accepted.

## Prior coverage cited, not repeated

- `sh-agent report` (feedback loop, redaction, push target):
  `.agents/docs/reviews/2026-09-17-security-6af54f7-015cf2d.md` (8 findings),
  `.agents/docs/reviews/2026-09-17-security-015cf2d-012b4ca.md` (8 findings)
- Embedded source and `ResolveSource`: `.agents/docs/reviews/2026-09-17-security-f7ee817-21ba8dc.md`
- Throwaway clone and `RemoveKeptClone`: `.agents/docs/reviews/2026-09-17-security-73f2c36-83f3cf6.md`
- `update` from a signed release: `.agents/docs/reviews/2026-09-17-security-dab2565-4d45155.md`
- Also in range and counted as coverage, though the dispatch did not name them:
  `-security-012b4ca-worktree.md` and `-security-012b4ca-worktree-final.md`
  (the redaction round-3 probes).

## Part 1 — are the claimed fixes present and whole at HEAD?

**All present. None found partially applied.** Each was located in the shipping
source.

### Report round 1 (`-security-6af54f7-015cf2d.md`)
| Claim | At HEAD | Whole? |
|---|---|---|
| F1 wider needles, no bare `sk-` | `report.go:38-57` — `prefixNeedles` (19 vendor shapes, shortest `sk-ant-`/`glpat-`; no bare `sk-`), `keywordNeedles` (4 case-insensitive regexps) | yes |
| F2 no key path in `report.json` | `report.go:94-100` — `Redaction` has `File/Line/Pattern/Lines/Unterminated` and **no `Path` field at all**: the leak is structurally impossible, not merely unset | yes, strongest form |
| F3 push target + env redirects | `report.go:598` `gitRedirects` = the five vars, consumed `617-621` in `SendReport` **before** `Collect`; `verifyFeedback:579-594` requires `rev-parse --show-toplevel` to `EvalSymlinks`-equal `dir`, and `remote.origin.url` to equal the URL exactly | yes |
| F4 block rule before JSON; `blockEnd`; rewritten dashes | `report.go:159` the block rule is the FIRST branch of the loop; `blockEnd:104-133` searches from the header's offset and re-opens on a header following a footer; `blockOpen:63` matches `BEGIN[A-Z ]*PRIVATE KEY` | yes |
| F6 `runDir` refuses a linked `reports` | `report.go:775-779` — `os.Lstat` on `reports` and `reports/<host>`, refuses `!IsDir()`, before any `MkdirAll` | yes |

### Report round 2 (`-security-015cf2d-012b4ca.md`)
| Claim | At HEAD | Whole? |
|---|---|---|
| F1 push to the URL, not origin | `report.go:851` — `gitOutC(dir, "push", "--quiet", "-u", url, "HEAD:main")`; no remote name decides the destination | yes |
| F2/F3 per-value JSON walk deleted | `report.go:149-181` `redactLines` — the walk is gone; `decoded()` (`191+`) is used for MATCHING only and the whole line goes. The retirement is stated in the doc-comment at `145-148` | yes — mechanism deleted, not patched |
| F4 origin not echoed | `report.go:589-592` — names the expected URL only, with the reason in a comment | yes |

### Embed / source (`-security-f7ee817-21ba8dc.md`)
| Claim | At HEAD | Whole? |
|---|---|---|
| F1 `TestEmbedHoldsOnlyTrackedFiles` | `embed_test.go:54`; skips explicitly without git or outside a checkout (`55-59`) | yes |
| F2 `knownCheckout` refuses `EmbeddedPrefix` | `source.go:42` defines it; `main.go:341` and `main.go:515` both `CutPrefix` before comparing paths; `repo.go:42` also excludes an embedded manifest from `ResolveSource`'s candidates | yes, at all three sites |
| F1's deferred release-time half (D23) | present in `remote.md` § Releases step 0 (`git status --porcelain --ignored --untracked-files=all -- context skills agents commands rules` + the test) | present but **mis-ordered** — Finding 1 |

### Clone / update (`-security-73f2c36-83f3cf6.md`)
| Claim | At HEAD | Whole? |
|---|---|---|
| F1 `RemoveKeptClone` fails closed | `report.go:558` `rev-list --count --all --not --remotes`; `509` `filepath.Abs(home)`; `538` `status --porcelain --untracked-files=all --ignored` | yes |
| F2 `update --check` deletes nothing | read the whole function: the two `checkOnly` returns are `main.go:430` (`return nil`) and `435` (`exitCode(2)`); `RemoveKeptClone` is at `459`, below both, and below the dirty-tree guard as well (so a dirty checkout removes nothing either — fail-safe) | yes |
| F3 `sweepInterruptedClones` | defined `report.go:409`, called `report.go:701` at the start of the throwaway path | yes |
| F4 `scrubCredentials` | `report.go:494`, applied to both the URL and git's message at `472` | yes |
| F5 `LC_ALL=C` for the recognised message | `report.go:849-852` — `gitOutC` for that one push; the match covers `rejected` and `non-fast-forward` | yes |

### Release consumption (`-security-dab2565-4d45155.md`)
| Claim | At HEAD | Whole? |
|---|---|---|
| F1a tag validated before argv | `release.go:101-103` — `IsReleaseVersion(tag)` before the switch and before gh sees it | yes |
| F1b never backwards/sideways | `release.go:104-118` — `==` up to date; a non-release running binary left alone; `compareRelease < 0` refused | yes |
| F1c version inside the signed bytes | `release.go:47` `sumsVersionMarker`, checked `249-255` **after** `verify.Run()` at `243`. Nothing in the file is believed before the signature holds | yes, correctly ordered |
| F2 `sshKeygen` prefers a system path | `release.go:299-316` — `/usr/bin`, `/bin`, then PATH | yes |
| F3 random 0700 dir, stale siblings cleared | `release.go:138-147` — glob-and-remove `release-*`, then `os.MkdirTemp` | yes |
| F4 `checksumFor` strict | `release.go:272-295` — 64 chars, hex-only, refuses two DIFFERENT hashes, tolerates an exact duplicate | yes |
| F5 stderr captured | `release.go:190-209` — `strings.Builder` stderr; "no asset" distinguished by `os.Stat` | yes |
| F6 `readCapped` | `release.go:52-53` (256 MiB / 1 MiB), used at `203` and `222` | yes |

Two properties no prior report states, confirmed here: (a) `verifyRelease`'s
order is signature → version marker → checksum → byte compare, and the
downloaded bytes are never made executable or run until `writeAtomic`
(`release.go:165`), which runs only after `verifyRelease` returns nil — so the
`exec.Command(dest, "install")` at `172` cannot execute unverified bytes; (b) the
release branch of `update` passes only `installPassArgs(*homeFlag, home,
*vendorFlag)` — the operator's own flags — to the new binary, so no external
input reaches that argv.

## Part 2 — release-diff paths no security report names

`git diff v1.0.0..HEAD --name-only` (148 paths) minus every path and basename
appearing in the seven in-range `-security-*.md` reports leaves **97 paths, of
which none is Go source and none is a skill**: trail markdown
(`bugs/BUG-051..072`, `changes/`, `handoffs/`, `plans/`, `retros/`, non-security
`reviews/`, `testing/`) plus `ROADMAP.md` and three agent definitions
(`agents/changelog.md`, `documentation.md`, `verifier.md`). Every Go file in the
release diff is named by at least one prior security report whose tree is
byte-identical to HEAD. The uncovered paths were read for the surface that
applies to them — see Part 3.

## Part 3 — secrets, credentials, private data, instruction-shaped text

Added-line scan over the whole 13,690-line release diff for
`ghp_|gho_|ghs_|ghu_|github_pat_|glpat-|sk-ant-|sk-proj-|xox[baps]-|AKIA|ASIA|AIza|ya29.|GOCSPX-|-----BEGIN|PRIVATE KEY|BEGIN SSH SIGNATURE`
returns hits only inside redaction documentation and test-fixture discussion.
Every one is a deliberate placeholder or a vendor's published example:
`AKIAFAKEVALUE`, `AKIAIOSFODNN7EXAMPLE` (AWS's own doc example), `ghs_FAKETOKEN`,
`FAKEVALUE`, `MIIEvFAKEbody`, and `secret123` (pre-existing in the vendored
`skills/terraform/references/security-compliance.md`, not added here).
**No live credential is added by this release.** No `.env`, `*.pem`, `*.key`,
`*.tfvars`, `*.tfstate` or credential JSON appears in the diff.

Two things that look like secrets and correctly are not:
- `release.go:39` and `remote.md` § Releases carry
  `ostling.simon@gmail.com ssh-ed25519 AAAAC3Nza…62NCR`. A PUBLIC key, and its
  presence IS the control — the allowed-signers line must be compiled in rather
  than read from disk. `release.go:26-38` says so, including what it does not
  pin (ssh-keygen itself).
- `remote.md` names the private key's PATH (`~/.ssh/id_ed25519`) and its
  fingerprint. A path and a fingerprint are not key material.

The signing procedure does not expose the key: `ssh-keygen -y -f
~/.ssh/id_ed25519 | ssh-keygen -lf - | cut -d' ' -f2` derives the PUBLIC
fingerprint for comparison, so the private half is never printed, never placed
in a variable and never written to a file; `tag -s` and `-Y sign` both take the
key by path. Nothing in the procedure prints anything sensitive. Finding 2
covers the one residue it leaves.

**Instruction-shaped text.** This release changes `context/AGENTS.md`, nine
skills and three agent definitions, plus ~90 trail files, all of which ship to
every machine that clones this repository — a real delivery vector, so it was
read for it. Nothing found that redirects an agent to exfiltrate data, reach the
network, touch `~/.ssh` or a real `$HOME`, weaken a permission setting, or treat
an agent message as consent. The changes are the operator's own harness rules
and they *strengthen* controls rather than relax any:
- `context/AGENTS.md` — BUG-058's rule that a finding counts as fixed only when
  every trigger it names has a test failing without the fix; the "second round
  on the same mechanism → stop and ask" rule; the stop-sign checklist after any
  dispatch that can write (`git status`/`log`/`tag`, BUG-044/045).
- `skills/preflight/SKILL.md` — stage by name, never `.env*`/`*.pem`/`*.key`/
  `*.tfvars`/`*.tfstate`/credential JSON, and `gitleaks git --staged` (with an
  explicit needle-scan fallback) refusing on a hit before the release commit.
- `agents/verifier.md`, `documentation.md`, `changelog.md` — guide-path change
  (ADR-032), app resolution (BUG-062), guide-write duty (BUG-063), and BUG-061's
  every-claim-traces-to-the-diff rule.
- **Checked specifically:** `--dangerously-skip-permissions`,
  `bypassPermissions` and `acceptEdits` appear **nowhere** in `skills/`,
  `context/`, `commands/`, `rules/`, `AGENTS.md`, `REQUIREMENTS.md` or
  `testing/guide.md`. The flag occurs only inside dated testing logs under
  `.agents/docs/testing/`, each as a record that the operator approved it for
  specific Antigravity runs and that `settings.json` was not modified. Nothing
  prescribes it. No finding, but worth recording as an explicit negative
  because a prescribed bypass here would have shipped everywhere.

Related code check (no prior report covers it): `attributionSet`
(`main.go:615-701`) is the only vendor setting the installer writes. It is
unchanged in this release — only its call site moved into a helper — and it
cannot weaken an operator's `permissions` key: it writes only an ABSENT
`attribution` key, refuses a symlinked `settings.json` (`Lstat`), creates with
`O_EXCL` so a planted link between the two calls fails rather than being written
through, caps at 1 MiB, and inserts textually with a re-parse instead of
re-encoding the operator's file. `permissions` appears nowhere in the Go source.

## Part 4 — does the release procedure produce artifacts the shipping code accepts?

Traced end to end, `remote.md` § Releases against `release.go`. **It does.**

| Procedure produces | Code requires | Match |
|---|---|---|
| `printf '# styrman-harness %s\n' vX.Y.Z > SHA256SUMS` | `sumsVersionMarker` = `(?m)^#\s*styrman-harness\s+(v[0-9]+\.[0-9]+\.[0-9]+)\s*$`, then `named == tag` | yes — regex checked against the exact literal emitted |
| `sha256sum <three binaries> >> SHA256SUMS` (two-space separator) | `checksumFor` cuts at the first space, `TrimSpace`s the name, strips a `*` | yes; the `#` marker line cannot be misread as a hash line (its name field never equals an asset name) |
| `ssh-keygen -Y sign -n sh-agent-release SHA256SUMS` | `ReleaseNamespace = "sh-agent-release"`, `-Y verify -n` the same | yes |
| signer line written to a local `allowed_signers` | `releaseSigner` compiled in; identity cut at the first space | yes — and the procedure states in as many words that its own file is only a pre-upload check |
| assets `sh-agent-{linux-amd64,darwin-arm64,linux-arm64}` | `AssetName(runtime.GOOS, runtime.GOARCH)` | yes for those three; any other host (e.g. darwin/amd64) gets the clean refusal `release vX.Y.Z has no sh-agent-<goos>-<goarch>` — a coverage gap, not a security defect |
| `-trimpath`, `CGO_ENABLED=0`, `GOTOOLCHAIN=local`, `-X main.version=vX.Y.Z` | `version` is what `UpdateFromRelease` compares; `GOTOOLCHAIN=local` is consistent with ADR-029's no-`toolchain`-directive decision | yes |
| `gh release create … --verify-tag` | n/a (publish side) | consistent |

The version marker is the one that mattered: a v1.0.0-era `SHA256SUMS` carries
no marker, so `update` correctly refuses it (`release.go:251`) — the procedure
is the only thing that can make a release installable, and it writes the line.

## Findings

- **[low / confidence: high]** `.agents/docs/deploy/remote.md` § Releases (binaries block, step 0) and `skills/preflight/SKILL.md` Step 6.4-6.5 — **the `go:embed` contamination guard runs after the tag is pushed, which is the one irreversible step.** The documented order is: clean-tree test → `git tag -s` → verify tag → `git push origin main` → `git push origin vX.Y.Z` → *then* `test -z "$(git status --porcelain --ignored --untracked-files=all -- context skills agents commands rules)"` and `go test -run TestEmbedHoldsOnlyTrackedFiles`. The skill encodes the same sequence (6.4 "Push the branch and the tag", 6.5 "Release ARTIFACTS") and even anticipates the failure: "stop with the tag pushed and say so". Exploitable in principle: an ignored or untracked file under the five embedded directories — `.gitignore` covers `.env`, which `all:` embeds, which is BUG-066's own reproduction — is compiled into every released binary, and the guard that catches it fires only once `vX.Y.Z` is public and, by this repo's own rule ("a pushed tag is never moved"), unrecallable. The guard does prevent PUBLISHING the contaminated binary, so no secret ships; the cost is a burnt version number and a tag with no assets. The earlier clean-tree test does **not** close this: `git status --porcelain --untracked-files=all` without `--ignored` cannot see an ignored file. Remediation: move step 0's two commands above the `git tag` line in `remote.md` § Releases — they read only the working tree and nothing the tag creates — and add one sentence to `skills/preflight/SKILL.md` Step 6.5 saying a deploy guide's artifact PRECONDITIONS run before 6.3's tag, not with the build. No code change, no dependency.
- **[low / confidence: medium]** `.agents/docs/deploy/remote.md` § Releases — **the procedure writes an `allowed_signers` file in two places and removes neither**, leaving `.agents/scratch/allowed_signers` behind after every release plus one in the build directory. This is not key material (a public key) and `.agents/scratch/` is gitignored, so it cannot be committed. It is worth a line because a signer file resting on disk is precisely the artifact `release.go:26-38` argues must never decide what a genuine release is, and a later reader who finds one may wire something to it. Confidence medium: a maintenance trap rather than a reachable weakness today. Remediation: append `rm -f .agents/scratch/allowed_signers` after the `verify-tag` line and `rm -f allowed_signers` after the pre-upload `-Y verify`, each with its one-line reason.

No high or medium findings. Neither finding blocks the gate.

## Checked

- **Secrets** — added-line scan of the full release diff (13,690 lines) for 19 vendor token prefixes, PEM headers and SSH signature armor: all hits are documented fake fixtures or a vendor's published example. The compiled-in and documented `ssh-ed25519` line is a public key and is the control. The operator's email appears as the signer identity, which is its purpose. No secret-shaped path in the diff.
- **Injection** — the release path's only external inputs are gh's tag string and the release assets. The tag is regex-gated to `^v\d+\.\d+\.\d+$` before reaching argv (`release.go:101`), so no `-`-leading or shell-meaningful tag survives; asset names derive from `runtime.GOOS`/`GOARCH`, not input. `SHA256SUMS` parsing is length- and charset-strict and refuses ambiguity instead of resolving it. Every git and gh call is `exec.Command` with separate args — no shell anywhere in the codebase. `hostDir` reduces a hostname to one path element; `runDir` `Lstat`-refuses a non-directory `reports`/`reports/<host>`; `Source.readRegular` refuses a symlink or non-regular file, and embed cannot store a link.
- **Authz/authn** — no endpoints, no request handling. The two authorization decisions that exist are "is this release the operator's" (compiled-in allowed-signers line + namespace + in-signature version marker) and "is this checkout the feedback repository" (`verifyFeedback` + pushing to the URL literally). Both fail closed: a failed `gitOut` yields `""`, which never equals the expected URL.
- **Crypto** — nothing hand-rolled, and nothing in this release comes close. `crypto/sha256` used as a checksum (`release.go:260`); signature verification delegated to OpenSSH's `ssh-keygen -Y verify`; transport is git over SSH. No cipher, MAC, KDF, RNG, token format, session or auth flow is implemented or composed here. The one place a home-grown construction could have crept in — binding a version to signed bytes — is solved by requiring a comment line INSIDE the signed file rather than by inventing a signature format, which is the correct answer.
- **Dependencies** — `go.mod` requires nothing beyond the stdlib (the repo's stdlib-only rule) and this release adds none. `govulncheck ./...` on go1.27.1: **"No vulnerabilities found."** `gofmt -l cmd internal *.go` clean; `go vet ./...` clean. No npm ecosystem in the repo, so `npm audit` does not apply — not faked, simply absent.
- **GCP surface** — none. No IAM binding, service account, `--allow-unauthenticated`, Cloud Run flag, Secret Manager reference or service-account key file anywhere in the release diff. The only remote service is GitHub, reached via `gh` (carrying the operator's own credential, deliberate per R8/D7) and git over SSH.

## Not examined

Read end to end: `internal/harness/release.go` (all 378 lines),
`internal/harness/source.go` (all 86), `internal/harness/report.go` lines 1-200
(the redaction walk), `internal/harness/repo.go` lines 1-70 (`ResolveSource`,
`IsRepo`), `cmd/sh-agent/main.go` lines 180-215, 420-475 and 612-701,
`.agents/docs/deploy/remote.md` § Releases, `skills/preflight/SKILL.md` Step 6,
and the complete diffs of `agents/*.md`, `ROADMAP.md`, `context/AGENTS.md` and
all nine changed skills.

Verified by targeted location rather than full read — enough to answer "is the
fix present and whole", but NOT a full re-review:
`internal/harness/report.go` lines 200-860 (guard lines of `Collect`,
`BuildRecord`, `cloneThrowaway`, `commitAndPush`, `RemoveKeptClone` read; their
full control flow not re-traced), the remainder of `cmd/sh-agent/main.go` and
`internal/harness/repo.go`, and `apply.go`, `inventory.go`, `plan.go`,
`render.go` (only their `EmbeddedPrefix`/`secretShaped` sites). All are named by
a prior security report whose tree is byte-identical to HEAD, which is why the
delta check was scoped to the fixes.

Not read, and therefore named rather than cleared:
- **Test files in the range** — `internal/harness/{harness,report,release,source,pathguard}_test.go`, `cmd/sh-agent/main_test.go`, `plan_test.go`, `embed_test.go` (only the `TestEmbedHoldsOnlyTrackedFiles` header and skip guards read). Must-fail evidence for the claimed fixes was taken from the prior reports' `## Outcome` sections and **not re-derived here**; no mutation was run in this review. This is the largest single gap: the fixes are confirmed PRESENT, not confirmed TESTED, by this review.
- `AGENTS.md`, `README.md`, `REQUIREMENTS.md`, `.agents/docs/specs/sh-agent/delivery.md`, `.agents/docs/architecture/DECISIONS.md`, `.agents/docs/deploy/local.md`, `.agents/docs/testing/guide.md` — read only through the diff's scanned added lines, not end to end.
- The ~97 uncovered trail-markdown paths of Part 2: scanned for secrets (token regexes over added lines) and for hostile-instruction/exfiltration patterns (`ignore previous`, `curl`/`wget`/`nc`, `base64 -d`, pipe-to-shell, `dangerously`, `~/.ssh`, `settings.json`, `POST http`) across added lines, and read where those scans hit — not each file end to end.
- **Not run:** any mutation or must-fail probe; `sh-agent report --dry-run`; any `sh-agent update` in a scratch home; `go test ./...` (the verifier owns that tier — `.agents/docs/testing/2026-09-17-verify-preflight-v1.1.0.md`, which was UNTRACKED at review time). No network was reached and no `gh` command was run, per the dispatch's hard limits. No source file was edited; this report is the only file created.

Operational note for the gate, not a security finding: at review time
`git status --short --untracked-files=all` showed
`?? .agents/docs/testing/2026-09-17-verify-preflight-v1.1.0.md`. The procedure's
`test -z "$(git status --porcelain --untracked-files=all)"` will fail until that
file is in the release commit.
