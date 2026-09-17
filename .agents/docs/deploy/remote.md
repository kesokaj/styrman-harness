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

## The four places

The operator installs the harness in four places: the work Mac (Antigravity only), the home Mac (both CLIs), and two Linux hosts where they work on other projects (both CLIs). Every one is reached on the operator's own SSH key and `gh` login — there are no deploy keys and no per-machine credentials, because the operator authenticates as themselves everywhere and a key would restrict nothing (ADR-030 D1). A place gets the harness one of two ways:

- **A checkout the operator owns** — the normal path. First install as above; after that, `sh-agent update` whenever the operator decides the harness has changed. A timer (§ The resync loop) is a convenience a machine may add, not the mechanism.
- **An unpacked bundle** — a copy of the repository with no `.git`, for a major remodel or a machine that cannot authenticate. `go run ./cmd/sh-agent install` from the unpacked directory installs the running binary as is and never rebuilds (BUG-015); `sh-agent status` names it a plain copy with no commit; `update` refuses it, by design. To move it forward, unpack a newer bundle and install again. See `local.md` § From an unpacked bundle.

Nothing else carries the harness: no copy vendored inside a project, no per-project copy, no service. A project that wants the harness installs it with `sh-agent` into a home. That home receives installed FILES — `sh-agent install` writes the context, skills, agents, commands and rules and creates no repository. That alone does not stop a change flowing back: changing the harness takes a credential that can push to this repository, not a checkout, and anything that reaches the operator's SSH key or `gh` login can clone and push. A home seeded for agents is therefore read-only only when it reaches no such credential — no forwarded ssh-agent, no `gh` token, no readable key. This repository names no consumer project and depends on none (R8).

Experience flows back the one way too: each place runs `sh-agent report`, which sends a redacted copy of its `~/.config/sh-agent/` and a metrics record to the private `styrman-harness-feedback` repository (`local.md` § Report back).

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

A release is an SSH-signed annotated `vX.Y.Z` tag on `main`. It is a MARKER naming a version, not a way the harness reaches a machine: every place follows `main` with `sh-agent update`, and a tag changes nothing for any of them (ADR-030 D2, withdrawing ADR-026's tag channel). A release is cut by `/preflight vX.Y.Z` from the operator's own turn, from this checkout (never a plain copy), on a machine holding the signing key. Nothing else is tagged: no lightweight tag, no `-rc` suffix, no other prefix. The tag is the only version tracking; there is no version file.

The signer — the allowed-signers line `verify-tag` checks against:
```
ostling.simon@gmail.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBkx0BICH1xG18f0a0F/02Vw0OvBTfSOVNXkPkO62NCR
```
Fingerprint `SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU`; the key is `~/.ssh/id_ed25519` on the Linux box. The Mac cannot cut a release until it holds this key or the line above lists the Mac's own.

Procedure — what `/preflight` runs at its Step 6 here (the skill's Step 6.3 follows a project's Release section), each command its own call, its exit status read before the next; `<gated head>` is the head SHA of the gate report's `Range reviewed:` line, the one external input — the tag itself takes no target, so no placeholder can sit on both sides of a check:
```bash
test "$(ssh-keygen -y -f ~/.ssh/id_ed25519 | ssh-keygen -lf - | cut -d' ' -f2)" = SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU   # exit 0: the signing key is here, derived from the PRIVATE key — else stop: this machine cannot sign
# gate passed and the release commit made (preflight Step 6.3): HEAD is that commit; <gated head> is the head SHA of the gate report's `Range reviewed:` line
test -z "$(git status --porcelain --untracked-files=all)"                                              # exit 0: clean tree — the bare form inherits status.showUntrackedFiles (BUG-049)
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
git push origin main
git push origin vX.Y.Z                                                                                  # plain pushes, never --force
```
A tag that fails any check: `git tag -d vX.Y.Z`, fix, re-gate. A pushed tag is never moved or deleted; a bad release gets the next PATCH. `git` is never run under `sudo`. Rehearsed 2026-09-16 in `.agents/scratch/tagcheck/`: the operator's key passes; a lightweight tag reads `commit`, an unsigned annotated tag has no armor line and fails `verify-tag` with "no signature found", a foreign key fails with "No principal matched". The stricter forms above (`test` assertions, tag-on-HEAD, parent-is-gated-head, records-only diff, ancestor, the private-key-derived fingerprint) were rehearsed the same day in the same scratch repo: passing on a records-only release commit tagged at HEAD, failing on a tag made on the earlier commit, on a release commit carrying a source change, on a dirty tree, and on the unsigned tag.

Version: MINOR for anything an installed place gains (a new or changed skill, agent, command, rule, context rule or installer behaviour), PATCH for a fix to released behaviour only, MAJOR for what an existing deployment must react to (a moved install path, a dropped vendor, a `state.json` or manifest shape change, a new required flag) — the context's Git rule. First tag `v1.0.0`, base ref for its gate 4aa19be. The gate's verification row: the `AGENTS.md` Setup commands, `go run ./cmd/sh-agent install --home .agents/scratch/home` with `status` reporting no drift, and `sh-agent list` — 27 skills and 7 agents today; a release that changes them says so in its changelog entry. There is no deploy step after the push: each place picks the commit up with `sh-agent update` when the operator runs it. The handoff after a release names the live tag.

## Verify

Same as `local.md`: `sh-agent status` with no drift, then one start of each CLI. On the Mac, additionally confirm the Antigravity paths — passed 2026-09-04 (`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`); re-run after a change that touches Antigravity-specific behavior and record the result there.

## Operator steps

Standing by-hand procedures, never pending items (those live in the newest handoff):

- On the Mac after a push from the Linux box: `sh-agent update`, `sh-agent status`, then the three smoke checks on each CLI (`.agents/docs/testing/2026-09-03-verify-real-home.md`; the Mac procedure is `2026-09-04-verify-mac-real-home.md`).
- After any change to `context/`, `skills/`, or `agents/` is pushed: `sh-agent update` at each of the four places, when the operator decides to. Nothing else needs refreshing.
- At a place that has been in use: read the originals under `~/.config/sh-agent/` (the dry run shows what is caught, never what is missed), `sh-agent report --dry-run`, then `sh-agent report` (`local.md` § Report back). Whether the work Mac reports at all is the operator's call — its corrections come from work sessions.
- On a new machine, `ssh-keygen -y -f ~/.ssh/id_ed25519 | ssh-keygen -lf -` must print the fingerprint in § Releases before that machine cuts a release (derived from the private key: a copied `.pub` alone proves nothing).
- Go upgrades on the Mac: `brew upgrade go`. On the Linux box the distribution had nothing newer, so it runs a user-level SDK on PATH ahead of `/usr/bin/go` — the procedure is in `local.md` § Operator steps, and it needs no root.

## Roll back

`git checkout <last good commit> && go run ./cmd/sh-agent install` on the affected machine; `sh-agent update` will refuse to move until the branch is fast-forwardable again, so push a fix or reset the branch first. An unpacked bundle rolls back by installing the earlier bundle.

## Access

Remote boxes are reached only on the operator's explicit request with details they supply. Nothing in this repo assumes a host.
