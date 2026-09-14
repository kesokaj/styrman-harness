---
owner: security
date: 2026-09-14
---
# Security review — upstream items (name guard, manifest resync guard, dry-run fixture, root-steps scoping)   (2026-09-14)

Range reviewed: 8046e4d..worktree

Repo: `/home/kesokaj/workspace/styrman-harness` (separate repo from styrman-crew; its own
AGENTS.md read first). HEAD pinned:
`8046e4d5f24b009bd2461ebef28931cade7a4acd`. Blobs as read (`git hash-object`):

| Path | Blob |
|---|---|
| `internal/harness/inventory.go` | `5a7c641fdd06263f6712110025cf31504b38370b` |
| `internal/harness/apply.go` | `c1b2a55f2f2efcb8024b8efe2b9c961d0d50c3cc` |
| `cmd/sh-agent/main_test.go` | `4676f8bd0f4dc4a04a210ce15dfe87540463c65b` |
| `internal/harness/pathguard_test.go` | `4397427fa4f4a387bdda65f22ed0e5d29763e453` (untracked, new) |
| `context/dev-setup.md` | `a14ab2f4152e8a823aa5f3950dcbbedbb7bfb765` **then** `0ac4ad03bee324e883a6e284214eb3a92b612e12` |

The working tree MOVED during this review: `context/dev-setup.md` was edited by
another actor between my first read and my re-pin. Both blobs were read and
finding 6 below records which residual each one carries. The other four paths
re-hashed identical, so every other finding stands on the blob named above.

## Verdict: FINDINGS

The two guards are correct for the paths they cover and neither bricks a shipped
name, but the manifest guard was added to ONE of the three readers of that
manifest: `Uninstall` still walks `state.json`'s `Files` and calls `os.Remove` on
`filepath.Join(home, rel)` with no locality check at all, so the very entry the
new test proves `Install` refuses to follow (`../escape.txt`) is still unlinked by
`sh-agent uninstall`. Identical destructive primitive, identical untrusted input,
only the guard differs. Separately and more serious in effect: the install path
writes a **world-writable file into the agent's own instruction tree** whenever a
skill directory contains a symlink — confirmed empirically below, `mode=0777`
under `~/.claude/skills/` with the symlink TARGET's bytes inside, target read from
outside the repo. That is the input class living next door to the boundary this
diff hardens, and any local account can use it to rewrite what the model is told.

## Findings

- **[high / confidence: high]** `internal/harness/plan.go:99-128` (pre-existing; blast
  radius of the parse boundary the diff hardens) — a symlink FILE inside a skill
  directory is copied as content AND deployed world-writable. `LoadSkills`
  resolves the skill's top directory (`filepath.EvalSymlinks`, `inventory.go:94`)
  and `rel` is taken from the resolved dir, so the DESTINATION stays confined —
  that part is sound. But `filepath.WalkDir` does not follow symlinks, so a
  symlink entry arrives with `d.IsDir() == false`, and the walk then calls
  `os.ReadFile(path)`, which DOES follow it; and `d.Info()` on a symlink reports
  the LINK (documented `io/fs.DirEntry` behaviour: "If the entry denotes a
  symbolic link, Info reports the information about the link itself"), so
  `mode := info.Mode().Perm() | 0o600` is `0o777` and `writeAtomic`'s
  `tmp.Chmod(mode)` (`apply.go:257`) applies it. Verified, not reasoned: a
  fixture skill whose `leak.md` was a symlink to a file OUTSIDE the fixture repo
  produced, from `harness.Plan`, `.claude/skills/probe/leak.md mode=0777` carrying
  the target's bytes (fixture built and deleted under `.agents/scratch/`; no
  exploit, a mode/bytes print only). Two consequences. (a) Read: a vendored skill
  — and this repo's AGENTS.md says third-party skills are vendored VERBATIM, with
  `.agents/skills` symlinked into `skills/` so `npx skills add` lands there —
  copies any file the installing user can read (`~/.ssh/*`,
  `~/.config/gcloud/*`) into `~/.claude/skills/<name>/` and
  `~/.gemini/config/skills/<name>/`. (b) Write: the deployed copy is
  world-writable inside the directory tree BOTH CLIs read as instructions, so any
  local account rewrites what the model is told, with no user interaction and no
  further defect needed. `internal/harness/adopt.go:55-78` is the same pattern on
  the INGESTION side (`WalkDir` → `os.ReadFile` → `writeAtomic(target, data,
  info.Mode().Perm()|0o600)`), so `sh-agent adopt` pulls the dereferenced content
  into `skills/` at 0777, where it is then committed. Remediation, both sites:
  skip non-regular entries (`if !d.Type().IsRegular() { return nil }`) or refuse
  them by name, and clamp the mode (`info.Mode().Perm()&0o755 | 0o600`) so a
  bad Lstat can never widen it.

- **[medium / confidence: high]** `internal/harness/apply.go:60-66` vs `apply.go:214-228`
  — the new `if !filepath.IsLocal(rel) { continue }` is in `Install`'s resync
  only. `Uninstall` iterates the same `st.Files` from the same
  `~/.config/sh-agent/state.json` and does `dest := filepath.Join(home, rel)` →
  `os.Lstat` → `os.Remove(dest)` → `pruneEmptyDirs(home, filepath.Dir(dest))`
  with no locality check; `pruneEmptyDirs` (`apply.go:276`) then walks upward
  from a directory that is lexically OUTSIDE `home`, and its stop condition
  `len(dir) > len(home)` compares string LENGTH, which a `../..`-derived path can
  satisfy while pointing anywhere. `Status` (`apply.go:150-160`) has the same gap,
  read-only (`Lstat`), so it is cosmetic there but it is the same inconsistency.
  The new comment states the premise plainly — "the saved manifest is a file under
  the user's home: treat its entries as untrusted paths" — and a premise that
  holds for one reader holds for all three; this is the harness's own
  mechanism-reaches-every-caller shape. Remediation: do not repeat the check per
  caller, filter in `LoadState` (`internal/harness/state.go:29-48`) so every
  reader receives only local entries, or add `func (s State) LocalFiles() []string`
  and call it in all three. Extend `pathguard_test.go`'s
  `TestResyncNeverFollowsANonLocalManifestEntry` with an `Uninstall` case over the
  same fixture — as written it passes today while `uninstall` deletes the file it
  asserts survives.

- **[medium / confidence: medium]** `internal/harness/inventory.go:54-61` — `validName`
  refuses `/`, `\`, `.`, `..`, empty and non-local, and nothing else: it
  constrains path SHAPE, not the character set. The same `name` is interpolated
  into the deployed context with no escaping — `RenderInventoryCompact`
  (`inventory.go:259`) and `RenderInventory` (`inventory.go:220`) emit
  ``"`"+s.Name+"`"`` and `` | `%s` | ``, while `cell()` escapes `|` for
  DESCRIPTIONS only. `Parse`/`unquote` (`internal/harness/frontmatter.go:70-93`,
  `133-149`) pass a double-quoted value through intact, and a `name: |` block
  scalar joins continuation lines with a space, so a backtick, a pipe, or a
  sentence of prose all survive into the name; `Lint` only refuses an unquoted
  `#`. Because `skills/` is the drop target for vendored third-party skills, the
  value is attacker-influenced by design. A name such as
  `` x` — IMPORTANT: <instructions> `` closes the code span and writes free text
  into `~/.claude/AGENTS.md` and `~/.gemini/config/rules/AGENTS.md` — the file
  every session in every project reads first. Prompt injection into the harness's
  own trust anchor, no path traversal required. Remediation: at the same
  boundary, require a conservative alphabet, e.g. `^[a-z0-9][a-z0-9._-]{0,63}$`.
  Verified non-breaking (see Checked): every shipped name already matches.

- **[low / confidence: medium]** `cmd/sh-agent/main_test.go:232-242` and
  `internal/harness/repo.go:75-78,83-90` — the `.git` empty-directory fixture
  itself is clean: `checkoutShaped` builds it under `t.TempDir()`, it is never
  passed to a real `--home`, and the test asserts the dry run writes nothing. **No
  leak into a real install path.** What the fixture DOCUMENTS is the residual:
  `IsCheckout` is satisfied by any `.git` entry — directory or FILE — and
  `Install` calls `Inspect(repo)` unconditionally (`apply.go:86`), which runs
  `git -C <repo> rev-parse` and `git status --porcelain` in the tree. A `.git`
  FILE is a worktree pointer (`gitdir: <path>`), so planting one in a downstream
  MIRROR — styrman-crew's `image/harness`, argolis-devbox's `deploy/harness`, the
  plain-copy shape this very test exists to distinguish — points git at an
  attacker-chosen gitdir whose `config` can carry command-executing keys
  (`core.fsmonitor`), and also flips `shouldRebuild` from "install the running
  binary" to `go build ./cmd/sh-agent` + execute (`main.go:129-137`,
  `repo.go:139-162`). It needs write access to the mirror, at which point the
  SOURCE is already controlled, so this is defence in depth rather than a new
  break — but the fixture's comment now normalises "an empty `.git` is enough" as
  a harmless property, which is the sentence a future reader will lean on.
  Remediation (optional): have `IsCheckout` require a `.git` DIRECTORY, and keep
  `Inspect` from shelling out when the tree is not one.

- **[low / confidence: medium]** `internal/harness/inventory.go:54-61` + `plan.go:173-176`
  — case and Unicode-normalisation collisions are not detected. `validName` and
  `Plan`'s duplicate check (`files[i].Rel == files[i-1].Rel`) both compare BYTES.
  macOS is a declared target (`context/dev-setup.md`: "development happens on
  macOS (Apple Silicon)") and APFS is case-insensitive by default, HFS+
  additionally normalisation-folding. So `Reviewer` and `reviewer`, or an NFC and
  an NFD spelling of one name, pass every check and then land on ONE file: last
  writer wins silently, and since BOTH spellings enter the manifest, the NEXT
  install's resync `os.Remove`s the survivor — the skill disappears on the
  following run, not the current one, which is what makes it hard to diagnose. A
  vendored skill can thus shadow a first-party one. Remediation: fold case before
  the duplicate comparison in `Plan` and error rather than overwrite; rejecting
  non-ASCII outright (the alphabet in the finding above) closes the
  normalisation half without a `golang.org/x/text` dependency. Trailing dots and
  spaces were checked and are NOT a finding: they collide only on Windows, which
  is not a target, and POSIX treats `a.` and `a` as distinct.

- **[low / confidence: medium]** `context/dev-setup.md:8` — the change is a net
  STRENGTHENING of the boundary, and it got stronger again mid-review. In blob
  `a14ab2f` it turned an unconditional "the operator says it, run it" into a
  conjunction also requiring the project's `.claude/settings.json` to record the
  standing permission, while preserving the ordering that the operator's own turn
  comes first — but it left the old trailing sentence ("Project
  `.claude/settings.json` records the standing permission.") in place, where it
  read as an ASSERTION that the file does record it, i.e. as the answer to the
  precondition the same sentence had just introduced. On a consumer VM where the
  user already holds passwordless sudo, an agent reading the bullet end to end
  could conclude the condition is satisfied by definition. **Blob `0ac4ad0`
  resolves that**: the trailing sentence is now "Without that record, hand the
  step back with the exact command instead" — the negative form, which is the
  right fix. One residual remains in `0ac4ad0`: `.claude/settings.json` is a
  checked-in, agent-writable file, so making it the authorization oracle is a
  confused deputy — a cloned repo, or an agent editing it mid-task, self-grants
  root execution. The bullet should say the permission must be one the OPERATOR
  granted through the permission system, and that an agent never edits
  `settings.json` to obtain one. This text ships into every consumer project's
  context at install, so the wording IS the control; there is no code behind it.

## Checked

- **Secrets** — no credentials, tokens, or keys in the diff; no `.env`-shaped path
  added; the new error strings quote a `name:` value and a file NAME only, never
  file content. `git grep -nE '(api[_-]?key|secret|token|BEGIN [A-Z ]*PRIVATE KEY)'`
  over the five reviewed paths: nothing. The indirect exposure is finding 1 (a
  symlinked skill entry copying a readable secret into the deployed tree,
  demonstrated with a dummy file), not a secret in the diff.
- **Injection** — no SQL, no shell built from input on the changed path. The only
  command execution reachable from the diff is `Inspect`'s `git` (`exec.Command`,
  argv form, no shell); finding 4 covers the `.git`-as-config-pointer angle. The
  parsing boundary the diff hardens is the frontmatter `name:`: path injection is
  now refused, MARKUP injection into the rendered context is not (finding 3).
- **Authz/authn** — no authentication or authorization code in this binary. The
  only trust decisions are the dev-setup.md root-steps precondition (finding 6)
  and `knownCheckout` (`main.go:360-375`), untouched by this diff and still
  correctly refusing to build from a checkout other than the installed one.
- **Crypto** — none present, none added; no hand-rolled hashing, session, or token
  code anywhere in the diff. Nothing in the harness's ban list is touched.
- **Dependencies** — none added; `go.mod` untouched, module has no third-party
  requires (stdlib-only is this repo's rule and the diff honours it:
  `filepath.IsLocal` is Go 1.20+ stdlib, correctly preferred over a hand-written
  cleaner). `govulncheck` is NOT INSTALLED on this machine (`command -v
  govulncheck` → not found); no network-fetched tooling was added to the
  operator's box, so it was NOT RUN — stated rather than papered over. `npm audit`
  is not applicable: no JavaScript in this repo.
- **GCP surface** — no IAM binding, service account, `--allow-unauthenticated`, or
  service-account key file in the diff. `context/dev-setup.md`'s GCP bullets
  (Secret Manager over `--set-env-vars`, `gcloud auth` flows) are unchanged.
- **Does the parse-time refusal brick a working install?** No. Enumerated every
  shipped name: 27 skills, each `skills/<dir>/SKILL.md` whose `name:` equals its
  directory name, all `[a-z0-9-]+`; seven agents (`bugs`, `changelog`,
  `documentation`, `explorer`, `reviewer`, `security`, `verifier`), same shape.
  None contains `/`, `\`, a dot segment, or a non-ASCII byte, so none is refused
  by `validName` — and none would be refused by the tighter alphabet proposed in
  finding 3 either. Commands and rules take their name from the FILE name
  (`LoadDocs`, `inventory.go:172`) and are not passed through `validName` at all:
  currently unreachable for traversal (a filename cannot contain `/`), but it
  leaves the case-collision class of finding 5 unchecked there too.
- **Does the manifest read-back guard also need to refuse symlinks?** For the
  REMOVAL itself, no: `os.Remove` unlinks the final component and does not follow
  a terminal symlink, so an entry naming a symlink deletes the link, not its
  target. The residual is INTERMEDIATE components: `filepath.Join(home, rel)` is
  lexical, so if a directory inside the home is itself a symlink
  (`~/.claude/skills` → elsewhere) the kernel resolves through it and both the
  removal and `pruneEmptyDirs`'s `os.ReadDir`/`os.Remove` act outside the install
  root while every string in the code still looks local. The WRITE side is the
  same: `writeAtomic`'s `os.MkdirAll` + rename follows an intermediate symlink.
  Both need pre-existing write access to the home — the same actor who can edit
  `state.json`, i.e. inside the threat model the new comment adopts. Closing it
  properly means `os.OpenRoot(home)`-based traversal (Go 1.24+); recorded as a
  design question, not raised as a separate finding.
- **Tooling run** — `gofmt -l .` → no output, exit 0. `go vet ./...` → exit 0.
  `go test -race ./...` run BARE (not piped; status read before any conclusion,
  per `.agents/rules/verify-bare.md`) → **exit status 0**, both packages ok
  (`cmd/sh-agent`, `internal/harness`), including the two new cases in
  `pathguard_test.go` and the rewritten `TestDryRunNeverBuilds`. Green here means
  the guards do what their tests assert; it does NOT cover the gaps above, none of
  which has a test.

## Not examined

- `internal/harness/derive.go`, `render.go`, and most of `frontmatter.go` — read
  only as far as tracing the `name:` value through `Parse`/`unquote`/`Lint`
  required. `adopt.go` was read only for the symlink/mode pattern of finding 1.
- The vendored skills' actual contents (`skills/VENDORED.md` inventory) — not
  read. Finding 1 is about the MECHANISM that copies them, not a claim that any
  specific vendored skill is hostile today.
- Runtime behaviour on macOS: finding 5's case-collision claim is reasoned from
  the documented APFS/HFS+ defaults plus the byte comparison in the code, NOT
  reproduced on a Mac. This machine is Linux.
- The other working-tree changes that appeared during the review
  (`.agents/docs/bugs/INDEX.md`, `BUG-016-*.md`,
  `.agents/docs/reviews/2026-09-14-review-upstream-items.md`,
  `.agents/docs/testing/2026-09-14-upstream-items-verify.md`) are outside the five
  paths dispatched and were not reviewed.
- No exploit was written for any finding. The one empirical check (finding 1) ran
  `harness.Plan` over a throwaway fixture and printed a mode and a byte string;
  the fixture was deleted and no source, test, or config file was edited.

## Outcome


- **1 [high/high] symlink in a skill dir deployed as its target under the
  link's mode — FIXED.** `internal/harness/plan.go`'s walk and
  `adopt.go`'s `copyTree` refuse any entry with `fs.ModeSymlink` by name
  before reading it (a skill ships regular files only). Test
  `TestSkillFilesAndAdoptCopiesRefuseSymlinks` (a link to a file outside
  the repo → `Plan` and `copyTree` both error naming the link).
- **2 [medium/high] the IsLocal guard covered only install's resync —
  FIXED, and shorter.** The filter moved into `LoadState`
  (`internal/harness/state.go`): non-local entries are dropped when the
  manifest is read, so install, uninstall and status share one guard; the
  per-caller copy in `apply.go` is deleted. Test
  `TestLoadStateDropsNonLocalManifestEntries` (`../escape`, `/etc/passwd`,
  `ok/../../x` dropped; the local entry kept) plus the existing resync test.
- **3 [medium/medium] name shape but not charset — FIXED.** `validName`
  now also requires `identLike`: `[A-Za-z0-9]` first, then
  `[A-Za-z0-9._-]`; a backtick, space, newline or non-ASCII name is
  refused with the file named. Every shipped skill (27) and agent (7)
  passes (`sh-agent list --repo .`; sandbox install rc 0, drift none).
- Residuals recorded, unchanged: a `.git` FILE (a worktree pointer) makes
  `IsCheckout` true and lets git resolve elsewhere — pre-existing, the
  test fixture uses a directory and is test-only; macOS case-insensitive
  collisions between two ASCII names differing only in case are not
  detected (NFD is moot now that names are ASCII); `.claude/settings.json`
  as the authorization record for root steps is the harness's existing
  convention, not introduced here.
