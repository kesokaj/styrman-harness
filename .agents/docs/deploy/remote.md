---
owner: main
date: 2026-09-17
target: remote
---
# Deploy — other machines

The harness reaches a machine one way: `sh-agent` (R8, ADR-030). Git is the transport, on the operator's own credentials; there is no push-based deploy, no CI, no hooks, no service.

## First install on a machine

```bash
git clone git@github.com:kesokaj/styrman-harness.git ~/workspace/styrman-harness
cd ~/workspace/styrman-harness
go run ./cmd/sh-agent install
sh-agent status
```

`install` writes the `attribution` key on this machine; check `sh-agent status` prints no note.

On a machine that installs from a RELEASE BINARY instead, the download is verified before it is moved and the install is verified after it — two different checks, and skipping the second is what BUG-078 cost:

```bash
# One chained command: a failed check stops the move. V is the version you expect, from OUTSIDE
# the download (you cut it, or read it on a machine you already trust); ASSET is sh-agent-linux-amd64,
# sh-agent-darwin-arm64 or sh-agent-linux-arm64; the signer line is § Releases' allowed_signers line,
# taken from this file in a checkout, never from the release page. macOS: shasum -a 256 -c -
V=vX.Y.Z && ASSET=sh-agent-linux-amd64 \
&& D=$(mktemp -d) && cd "$D" \
&& gh release download "$V" --repo kesokaj/styrman-harness --clobber \
     --pattern "$ASSET" --pattern 'SHA256SUMS' --pattern 'SHA256SUMS.sig' --dir "$D" \
&& printf '%s\n' 'ostling.simon@gmail.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBkx0BICH1xG18f0a0F/02Vw0OvBTfSOVNXkPkO62NCR' > allowed_signers \
&& ssh-keygen -Y verify -f allowed_signers -I ostling.simon@gmail.com -n sh-agent-release -s SHA256SUMS.sig < SHA256SUMS \
&& grep -qxF "# styrman-harness $V" SHA256SUMS \
&& grep -v '^#' SHA256SUMS | sha256sum -c - --ignore-missing \
&& chmod +x "$ASSET" \
&& test "$("./$ASSET" version)" = "sh-agent $V" \
&& mv "$ASSET" ~/.local/bin/sh-agent && cd ~ && rm -r "$D"
# The `grep` and the `test` tie the SIGNED checksum file and the binary to the version YOU named —
# what `sh-agent update` does with the tag it asked GitHub for (update also refuses a tag that is not
# newer than what it has, which no first install can mirror). Without them a genuinely signed
# OLDER release, served whole, passes the signature and checksum steps: the signature only proves
# the operator signed it once. The checksum and signer files leave with "$D". A system path
# instead of ~/.local/bin is a `sudo mv` on that one line and changes nothing else.

sh-agent version                      # unqualified: what the SHELL resolves, after the move
which -a sh-agent                     # every candidate, in PATH order, when the version is wrong
```

`sh-agent version` unqualified is the one that matters. A binary moved into a system path is not installed if an earlier `PATH` entry holds another `sh-agent`: that one runs, installs whatever source it resolves, and reports `drift none` about it. It tells stale from fresh, never genuine from forged — it is the suspect binary's own output; authenticity is the signature and checksum on the file, before the move. `install` and `status` name such a binary now (a `warning:` line and a `shadowed` line), and a symlink from the earlier entry to the installed binary settles it — it runs the same file, so it is not a shadow. `--clobber` matters for the same class of reason: without it a leftover file from an earlier attempt makes the download refuse, and the stale file is what gets moved.

## The four places

The operator installs the harness in four places: the work Mac (Antigravity only), the home Mac (both CLIs), and two Linux hosts where they work on other projects (both CLIs). Every one is reached on the operator's own SSH key and `gh` login — there are no deploy keys and no per-machine credentials, because the operator authenticates as themselves everywhere and a key would restrict nothing (ADR-030 D1). A place gets the harness one of two ways:

- **A checkout the operator owns** — the normal path. First install as above; after that, `sh-agent update` whenever the operator decides the harness has changed. A timer (§ The resync loop) is a convenience a machine may add, not the mechanism.
- **A release binary** — one downloaded `sh-agent` with `context/`, `skills/`, `agents/`, `commands/` and `rules/` built into it, for a place that only USES the harness and should not have a checkout of it sitting beside its projects. `sh-agent install` installs what the binary carries and copies the binary itself; `sh-agent status` shows `source embedded vX.Y.Z`; `sh-agent update` downloads the newest release, verifies its signed checksums and replaces the binary, so the place changes only when a release is cut. `sh-agent adopt` needs a checkout and says so. See `README.md` § Install without cloning and `local.md` § Update a release-binary install.

A plain copy with no `.git` — an unpacked bundle — still installs, and `update` still refuses it (`local.md` § From a plain copy); it is no longer one of the ways in, it is the fallback when neither of the two can reach GitHub.

Nothing else carries the harness: no copy vendored inside a project, no per-project copy, no service. A project that wants the harness installs it with `sh-agent` into a home. That home receives installed FILES — `sh-agent install` writes the context, skills, agents, commands and rules and creates no repository. That alone does not stop a change flowing back: changing the harness takes a credential that can push to this repository, not a checkout, and anything that reaches the operator's SSH key or `gh` login can clone and push. A home seeded for agents is therefore read-only only when it reaches no such credential — no forwarded ssh-agent, no `gh` token, no readable key. This repository names no consumer project and depends on none (R8).

Experience flows back the one way too: each place runs `sh-agent report`, which sends a redacted copy of its `~/.config/sh-agent/` and a metrics record to the private `styrman-harness-feedback` repository through a throwaway clone holding only that place's own folder, removed after the push — so no place holds any other place's reports (`local.md` § Report back).

## The resync loop

```bash
# machine A: change, verify, install, commit, push
sh-agent install && git commit -am "context: tweak" && git push
# machine B
sh-agent update            # fetch; refuse a dirty checkout; if behind: pull --ff-only; rebuild and reinstall
sh-agent update --check    # only reports; exit 2 when behind; works on a dirty checkout
```

`sh-agent update` finds the checkout from `$SH_AGENT_REPO`, then the working directory or any of its parents, then the path recorded at the last install, so it works from anywhere. It refuses a non-fast-forward pull; resolve that by hand in the checkout.

**`update` refuses a checkout with uncommitted changes, whether or not it is behind.** It deploys the WORKING TREE, not the commit it just fetched, so on a machine where the harness is also being edited it would put half-finished work into the live context — which is what an hourly timer did here on 2026-09-16 (BUG-046). Commit or stash, or run `sh-agent install` to deploy a dirty tree deliberately. That makes `update` safe to run from a timer:

```ini
# ~/.config/systemd/user/sh-agent-update.service  — with the matching hourly .timer
[Service]
Type=oneshot
ExecStart=%h/.local/bin/sh-agent update --repo %h/workspace/styrman-harness
Environment=PATH=%h/.local/bin:/usr/local/bin:/usr/bin:/bin
```

On a box where nobody edits the checkout the timer simply keeps it current; on a box where somebody does, it now declines instead of deploying their draft.

The timer follows `main` by fast-forward: its trust root is whoever can push to this repository, exactly as a hand-typed `sh-agent update` is. That is the operator, on their own machines, and the accepted position (ADR-030). A box where that is not acceptable runs `sh-agent update --check` on the timer and updates by hand. A tag is never what a machine installs from (§ Releases).

## Releases — signed tags, as release markers

A release is an SSH-signed annotated `vX.Y.Z` tag on `main`, and the binaries published with it. A tag marks a version; it also carries the `sh-agent` binaries that a BINARY-ONLY place updates to, while a checkout place still follows `main` with `sh-agent update` and a tag changes nothing for it (ADR-033, superseding ADR-030 D2's "a tag delivers nothing"; ADR-030 D2 withdrew ADR-026's tag channel, and this restores a narrower one: the assets, verified by signature, never a push to a machine). A release is cut by `/preflight vX.Y.Z` from the operator's own turn, from this checkout (never a plain copy), on a machine holding the signing key. Nothing else is tagged: no lightweight tag, no `-rc` suffix, no other prefix. The tag is the only version tracking; there is no version file.

The signer — the allowed_signers line `verify-tag` checks against:
```
ostling.simon@gmail.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBkx0BICH1xG18f0a0F/02Vw0OvBTfSOVNXkPkO62NCR
```
Fingerprint `SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU`; the key is `~/.ssh/id_ed25519` on the Linux box. The Mac cannot cut a release until it holds this key or the line above lists the Mac's own.

Procedure — what `/preflight` runs at its Step 6 here (the skill's Step 6.3 follows a project's Release section), each command its own call, its exit status read before the next; `<gated head>` is the head SHA of the gate report's `Range reviewed:` line, the one external input — the tag itself takes no target, so no placeholder can sit on both sides of a check:
```bash
test "$(ssh-keygen -y -f ~/.ssh/id_ed25519 | ssh-keygen -lf - | cut -d' ' -f2)" = SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU   # exit 0: the signing key is here, derived from the PRIVATE key — else stop: this machine cannot sign
# gate passed and the release commit made (preflight Step 6.3): HEAD is that commit; <gated head> is the head SHA of the gate report's `Range reviewed:` line
test -z "$(git status --porcelain --untracked-files=all)"                                              # exit 0: clean tree — the bare form inherits status.showUntrackedFiles (BUG-049)
# BEFORE the tag, never after the push: `go:embed` reads the WORKING TREE, so a stray file under the
# five embedded directories is compiled into every released binary. The clean-tree test above cannot
# see an IGNORED file, and a tag is unrecallable once pushed (D23, BUG-066; security review 2026-09-17)
test -z "$(git status --porcelain --ignored --untracked-files=all -- context skills agents commands rules)"   # exit 0: nothing untracked or ignored there
go test -count=1 -run TestEmbedHoldsOnlyTrackedFiles ./...                                             # exit 0: the embedded set is exactly what git tracks
test "$(git rev-parse HEAD^)" = "$(git rev-parse <gated head>)"                                        # exit 0: the release commit sits on the gated head
git diff --quiet <gated head> HEAD -- ':(top)' ':(top,exclude).agents/docs'                            # exit 0: records only — ':(top)' because a bare '.' is relative to the cwd and misses a change in a sibling directory (measured)
git -c gpg.format=ssh -c user.signingkey="$HOME/.ssh/id_ed25519.pub" tag -s vX.Y.Z -m "styrman-harness vX.Y.Z"   # no target argument: tags HEAD
# check the tag before any push — every line carries its verdict in its exit status:
test "$(git rev-parse vX.Y.Z^{commit})" = "$(git rev-parse HEAD)"                                      # exit 0: the tag is on the release commit
git merge-base --is-ancestor vX.Y.Z^{commit} main                                                       # exit 0: on main
test "$(git cat-file -t vX.Y.Z)" = tag                                                                  # exit 0: annotated ("commit" means lightweight)
test "$(git cat-file tag vX.Y.Z | grep -E '^-----BEGIN' | tail -1)" = '-----BEGIN SSH SIGNATURE-----'   # exit 0: SSH-signed (the bare pipeline exits 0 either way)
printf '%s\n' 'ostling.simon@gmail.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBkx0BICH1xG18f0a0F/02Vw0OvBTfSOVNXkPkO62NCR' > .agents/scratch/allowed_signers
git -c gpg.ssh.allowedSignersFile=.agents/scratch/allowed_signers verify-tag vX.Y.Z                     # exit 0: Good "git" signature for ostling.simon@gmail.com
rm -f .agents/scratch/allowed_signers                                                                   # written for the check above only
git push origin main
git push origin vX.Y.Z                                                                                  # plain pushes, never --force
```

Then the release's binaries, which is how a binary-only place updates (R9, ADR-033). Built here, from the tagged tree, never by a service:
```bash
# 0. the contamination guard AGAIN, here: it ran before the tag so a bad tree is
#    never tagged, and it runs again here because `go:embed` reads the working
#    tree at BUILD time — which is after the tag and the pushes, so a .DS_Store
#    or an editor swap file appearing under skills/ in between would otherwise be
#    embedded unguarded (D23, BUG-066; code review 2026-09-17 round 2). Both
#    commands are read-only, so running them twice costs nothing.
test -z "$(git status --porcelain --ignored --untracked-files=all -- context skills agents commands rules)"   # exit 0: nothing untracked or ignored there
go test -count=1 -run TestEmbedHoldsOnlyTrackedFiles ./...                 # exit 0: embedded set == git's
# 1. the three platforms, from the tag's own tree, reproducible flags
for target in linux/amd64 darwin/arm64 linux/arm64; do
  GOOS="${target%/*}" GOARCH="${target#*/}" CGO_ENABLED=0 GOTOOLCHAIN=local \
    go build -trimpath -ldflags "-X main.version=vX.Y.Z" \
    -o ".agents/scratch/release-vX.Y.Z/sh-agent-${target%/*}-${target#*/}" ./cmd/sh-agent
done
# 2. the checksums, with the VERSION as the first line — `sh-agent update` refuses a
#    release whose signed SHA256SUMS names another version or none (D26 of the
#    release-binary plan), and `sha256sum -c` ignores a '#' line
cd .agents/scratch/release-vX.Y.Z
printf '# styrman-harness %s\n' vX.Y.Z > SHA256SUMS
sha256sum sh-agent-linux-amd64 sh-agent-darwin-arm64 sh-agent-linux-arm64 >> SHA256SUMS
sha256sum -c SHA256SUMS                                                    # exit 0: the file describes these files
# 3. signed with the same key as the tag, in the namespace the binary verifies
ssh-keygen -Y sign -f ~/.ssh/id_ed25519 -n sh-agent-release SHA256SUMS
printf '%s\n' 'ostling.simon@gmail.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBkx0BICH1xG18f0a0F/02Vw0OvBTfSOVNXkPkO62NCR' > allowed_signers
ssh-keygen -Y verify -f allowed_signers -I ostling.simon@gmail.com -n sh-agent-release -s SHA256SUMS.sig < SHA256SUMS   # exit 0: what the places will check, checked here first
rm -f allowed_signers                                                      # it exists only for the check above: the line the PLACES verify against is compiled into the binary, and a signer file left on disk is the artifact release.go refuses to trust (security review 2026-09-17)
# 4. the release, with notes read by the operator before the upload — created with the
#    notes ALONE, then one upload per asset: a create carrying every asset rolls the whole
#    release back when a single upload fails, which is what GitHub did twice for v1.1.0
#    (`HTTP 500: Error creating asset temp dir`), while one call per file then worked first
#    time each. Retry an individual upload before calling it a failure.
gh release create vX.Y.Z --repo kesokaj/styrman-harness --verify-tag --title "styrman-harness vX.Y.Z" \
  --notes-file ../release-notes-vX.Y.Z.md
for a in SHA256SUMS SHA256SUMS.sig sh-agent-linux-amd64 sh-agent-darwin-arm64 sh-agent-linux-arm64; do
  gh release upload vX.Y.Z --repo kesokaj/styrman-harness "$a" --clobber || echo "retry $a by hand"
done
gh release view vX.Y.Z --repo kesokaj/styrman-harness --json assets --jq '.assets[].name'   # all five
cd -
```
`allowed_signers` is written only to run the check above and removed immediately; the line the PLACES verify against is compiled into the binary (`internal/harness/release.go`), so a file on any machine cannot decide what a genuine release is — which is exactly why this procedure does not leave one lying about. The release notes are built as § Release notes below says. Before this repository is ever made public, read the trail for what it names: `.agents/docs/` carries hostnames (`spgo-mac.roam.internal`, `devbox`), machine paths (`/opt/devbox/deploy/harness`), and a report record's contents, all of which are fine inside a private repository and are internal identifiers outside one. There is no secret in them — that is checked every release — but publication is a separate decision from releasing, and this is the paragraph that says so (security review 2026-09-18).

Publishing a release changes NOTHING about the repository's visibility: `styrman-harness` stays private, the release page is readable only by accounts that can already read the repository, and every asset download needs that credential — which is what README § Install without cloning depends on. `gh repo view kesokaj/styrman-harness --json visibility` after the upload is the check; only a deliberate change in the repository's settings makes anything public.

Upload the assets one call at a time, not all of them on `gh release create`: a single create-with-assets call ROLLS THE WHOLE RELEASE BACK when one upload fails, which is what GitHub's upload endpoint did twice for v1.1.0 (`HTTP 500: Error creating asset temp dir`, both attempts leaving no release and the tag intact). Create the release with `--notes-file` alone, then `gh release upload vX.Y.Z --repo kesokaj/styrman-harness <asset> --clobber` per file, retrying an individual upload before calling it a failure.

After the upload: `gh release view vX.Y.Z --repo kesokaj/styrman-harness --json assets --jq '.assets[].name'` lists the five assets, and nothing further is pushed anywhere — a binary-only place picks the release up the next time the operator runs `sh-agent update` there.

### Release notes

Every release carries notes, built from the `changes/` entries dated after the previous tag's commit and read by the operator before the upload (D19). Write `.agents/scratch/release-notes-vX.Y.Z.md` with one line per entry under the groups that have any:

```markdown
## New features
- <entry title> — <one line, from the entry's Summary>
## Improvements and optimisation
## Fixes
- <entry title> (BUG-NNN)
## Other updates
## Install
Download: see README § Install without cloning. Existing binary-only places: `sh-agent update`.
```

An entry's `## Type` decides its group: `feature` → New features; `task` or `refactor` that changes behaviour → Improvements and optimisation; `bugfix` → Fixes, with its BUG id; rules, skills, docs and context → Other updates. An empty group is left out. The same text goes into the release's own `changes/` entry.
A tag that fails any check: `git tag -d vX.Y.Z`, fix, re-gate. A pushed tag is never moved or deleted; a bad release gets the next PATCH. `git` is never run under `sudo`. Rehearsed 2026-09-16 in `.agents/scratch/tagcheck/`: the operator's key passes; a lightweight tag reads `commit`, an unsigned annotated tag has no armor line and fails `verify-tag` with "no signature found", a foreign key fails with "No principal matched". The stricter forms above (`test` assertions, tag-on-HEAD, parent-is-gated-head, records-only diff, ancestor, the private-key-derived fingerprint) were rehearsed the same day in the same scratch repo: passing on a records-only release commit tagged at HEAD, failing on a tag made on the earlier commit, on a release commit carrying a source change, on a dirty tree, and on the unsigned tag.

Version: MINOR for anything an installed place gains (a new or changed skill, agent, command, rule, context rule or installer behaviour), PATCH for a fix to released behaviour only, MAJOR for what an existing deployment must react to (a moved install path, a dropped vendor, a `state.json` or manifest shape change, a new required flag) — the context's Git rule. First tag `v1.0.0`, base ref for its gate 4aa19be. The gate's verification row: the `AGENTS.md` Setup commands, `go run ./cmd/sh-agent install --home .agents/scratch/home` with `status` reporting no drift, and `sh-agent list` — 27 skills and 7 agents today; a release that changes them says so in its changelog entry. There is no deploy step after the push: each place picks the commit up with `sh-agent update` when the operator runs it. The handoff after a release names the live tag.

## Verify

Same as `local.md`: `sh-agent version` unqualified matching what was installed, `sh-agent status` with no drift and no `shadowed` line, then one start of each CLI. On the Mac, additionally confirm the Antigravity paths — passed 2026-09-04 (`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`); re-run after a change that touches Antigravity-specific behavior and record the result there.

## Operator steps

Standing by-hand procedures, never pending items (those live in the newest handoff):

- On the Mac after a push from the Linux box: `sh-agent update`, `sh-agent status`, then the three smoke checks on each CLI (`.agents/docs/testing/2026-09-03-verify-real-home.md`; the Mac procedure is `2026-09-04-verify-mac-real-home.md`).
- After any change to `context/`, `skills/`, or `agents/` is pushed: `sh-agent update` at each of the four places, when the operator decides to. Nothing else needs refreshing.
- At a place that has been in use: read the originals under `~/.config/sh-agent/` (the dry run shows what is caught, never what is missed), `sh-agent report --dry-run`, then `sh-agent report` (`local.md` § Report back). Every place reports the same way, the work Mac included (plan D12).
- On a new machine, `ssh-keygen -y -f ~/.ssh/id_ed25519 | ssh-keygen -lf -` must print the fingerprint in § Releases before that machine cuts a release (derived from the private key: a copied `.pub` alone proves nothing).
- Go upgrades on the Mac: `brew upgrade go`. On the Linux box the distribution had nothing newer, so it runs a user-level SDK on PATH ahead of `/usr/bin/go` — the procedure is in `local.md` § Operator steps, and it needs no root.

## Roll back

`git checkout <last good commit> && go run ./cmd/sh-agent install` on the affected machine; `sh-agent update` will refuse to move until the branch is fast-forwardable again, so push a fix or reset the branch first. A plain copy rolls back by installing the earlier copy. A release-binary place rolls back by downloading the previous release's binary by hand and running `sh-agent install` — `update` never goes backwards, on purpose (README § Install without cloning).

## Access

Remote boxes are reached only on the operator's explicit request with details they supply. Nothing in this repo assumes a host.
