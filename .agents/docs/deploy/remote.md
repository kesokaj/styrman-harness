---
owner: main
date: 2026-09-16
target: remote
---
# Deploy — other machines

Every machine the operator works on (the Mac, the Linux boxes) runs its own copy of the harness. Git is the transport; there is no push-based deploy, no CI, no hooks.

## First install on a machine

```bash
git clone git@github.com:kesokaj/styrman-harness.git ~/workspace/styrman-harness
cd ~/workspace/styrman-harness
go run ./cmd/sh-agent install
sh-agent status
```

`install` writes the `attribution` key on this machine; check `sh-agent status` prints no note.

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

The timer follows `main` by fast-forward, not a signed tag: its trust root is whoever can push to this repository, exactly as a hand-typed `sh-agent update` is. That is the operator's own machine and the accepted position here (§ Releases covers the one consumer that does verify signatures, styrman-crew). A box where that is not acceptable runs `sh-agent update --check` on the timer and updates by hand.

## Downstream mirrors

Two projects vendor this repo as a .git-less copy: styrman-crew `image/harness/` and argolis-devbox `deploy/harness/`. Each mirror's `VENDORED.md` holds the refresh recipe: rsync from a clean checkout, record `COMMIT` and `SHA256SUMS`, then diff to confirm only those two files plus `VENDORED.md` differ. What happens after the refresh differs by project — styrman-crew builds `sh-agent` inside its own image (`image/Dockerfile`, builder stage) and ships a refresh by cutting a release, since the image version stamp is what makes a running project re-seed; argolis-devbox commits and pushes with no version to bump, and its VM builds `sh-agent` on the box itself, as root, via `deploy/scripts.d/032_setup-harness.sh` on its next auto-update (within five minutes). A mirror is a plain copy either way: `install` from it never rebuilds (the tree may be read-only to the installing user — BUG-015), `sh-agent status` there shows no commit, and `update` refuses, by design. After any change to `context/`, `skills/`, or `agents/` here, refreshing both mirrors is a task in each of those projects; a retro that applies a harness-scope proposal says so in its report. styrman-crew is replacing its mirror with the signed-tag channel below and, by its own decision in that handoff, does not refresh `image/harness/` in the meantime (it stays at 4aa19be); argolis-devbox's mirror is refreshed as before until its `sh-agent update` timer lands.

## Releases — signed tags

styrman-crew's host will consume this repository by its newest SSH-signed `vX.Y.Z` tag through a host-side, signature-gated channel (crew's handoff `2026-09-16-session-close-harness-channel-next.md`; selection `git tag -l 'v*' | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -1`, lightweight and non-SSH-signed tags refused). `sh-agent update` on operator machines and on argolis-devbox keeps following `main`; a tag changes nothing for them. A release is cut by `/preflight vX.Y.Z` from the operator's own turn, from this checkout (never a mirror or plain copy), on a machine holding the signing key. Nothing else is tagged: no lightweight tag, no `-rc` suffix, no other prefix. The tag is the only version tracking; there is no version file. ADR-026.

The signer — the one line of crew's `/etc/styrman/allowed_signers` (root-only on that host; that file is the trust root, the copy here adds none):
```
ostling.simon@gmail.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBkx0BICH1xG18f0a0F/02Vw0OvBTfSOVNXkPkO62NCR
```
Fingerprint `SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU`; the key is `~/.ssh/id_ed25519` on the Linux box. The Mac cannot cut a release until it holds this key or crew lists the Mac's own.

Procedure — what `/preflight` runs at its Step 6 here (the skill's Step 6.3 follows a project's Release section), each command its own call, its exit status read before the next; `<gated head>` is the head SHA of the gate report's `Range reviewed:` line, the one external input — the tag itself takes no target, so no placeholder can sit on both sides of a check:
```bash
test "$(ssh-keygen -y -f ~/.ssh/id_ed25519 | ssh-keygen -lf - | cut -d' ' -f2)" = SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU   # exit 0: the signing key is here, derived from the PRIVATE key — else stop: this machine cannot sign
# gate passed and the release commit made (preflight Step 6.3): HEAD is that commit; <gated head> is the head SHA of the gate report's `Range reviewed:` line
test -z "$(git status --porcelain)"                                                                    # exit 0: clean tree
test "$(git rev-parse HEAD^)" = "<gated head>"                                                         # exit 0: the release commit sits on the gated head
git diff --quiet <gated head> HEAD -- . ':!.agents/docs'                                               # exit 0: it changed records only — the source tree is the gated one
git -c gpg.format=ssh -c user.signingkey="$HOME/.ssh/id_ed25519.pub" tag -s vX.Y.Z -m "styrman-harness vX.Y.Z"   # no target argument: tags HEAD
# check the tag the way the consumer does — every line carries its verdict in its exit status; before any push:
test "$(git rev-parse vX.Y.Z^{commit})" = "$(git rev-parse HEAD)"                                      # exit 0: the tag is on the release commit
git merge-base --is-ancestor vX.Y.Z^{commit} main                                                       # exit 0: on main
test "$(git cat-file -t vX.Y.Z)" = tag                                                                  # exit 0: annotated ("commit" means lightweight)
test "$(git cat-file tag vX.Y.Z | grep -E '^-----BEGIN' | tail -1)" = '-----BEGIN SSH SIGNATURE-----'   # exit 0: SSH-signed (the bare pipeline exits 0 either way)
printf '%s\n' 'ostling.simon@gmail.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBkx0BICH1xG18f0a0F/02Vw0OvBTfSOVNXkPkO62NCR' > .agents/scratch/allowed_signers
test ! -d /etc/styrman || sudo -n cmp /etc/styrman/allowed_signers .agents/scratch/allowed_signers     # crew host only, exit 0 elsewhere: the one root step here, a read — git is never run under sudo
git -c gpg.ssh.allowedSignersFile=.agents/scratch/allowed_signers verify-tag vX.Y.Z                     # exit 0: Good "git" signature for ostling.simon@gmail.com
git push origin main
git push origin vX.Y.Z                                                                                  # plain pushes, never --force
```
A tag that fails any check: `git tag -d vX.Y.Z`, fix, re-gate; a `sudo -n cmp` mismatch means crew's trust root moved — update the quoted line here from the host file, never the other way round. A pushed tag is never moved or deleted; a bad release gets the next PATCH. Rehearsed 2026-09-16 in `.agents/scratch/tagcheck/`: the operator's key passes; a lightweight tag reads `commit`, an unsigned annotated tag has no armor line and fails `verify-tag` with "no signature found", a foreign key fails with "No principal matched"; `verify-tag` against `/etc/styrman/allowed_signers` itself fails with "Permission denied" for the operator user — hence the scratch copy and the `cmp`; `git` is never run under `sudo`. The stricter forms above (`test` assertions, tag-on-HEAD, parent-is-gated-head, records-only diff, ancestor, the private-key-derived fingerprint) were rehearsed the same day in the same scratch repo: passing on a records-only release commit tagged at HEAD, failing on a tag made on the earlier commit, on a release commit carrying a source change, on a dirty tree, and on the unsigned tag. The `sudo -n cmp` line has not been run in this repository's trail yet: it is the operator's root step on the crew host.

Version: MINOR for anything a consumer gains (a new or changed skill, agent, command, rule, context rule or installer behaviour), PATCH for a fix to released behaviour only, MAJOR for what an existing deployment must react to (a moved install path, a dropped vendor, a `state.json` or manifest shape change, a new required flag) — the context's Git rule. First tag `v1.0.0`, base ref for its gate 4aa19be (the commit both mirrors carry). The gate's verification row: the `AGENTS.md` Setup commands, `go run ./cmd/sh-agent install --home .agents/scratch/home` with `status` reporting no drift, and `sh-agent list` — 27 skills and 7 agents today, the counts crew's channel smoke-checks; a release that changes them says so in its changelog entry. There is no deploy step after the push: the consumer's timer does the rest and its outcome is observed on the crew host, not from here. The handoff after a release names the live tag.

## Verify

Same as `local.md`: `sh-agent status` with no drift, then one start of each CLI. On the Mac, additionally confirm the Antigravity paths — passed 2026-09-04 (`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`); re-run after a change that touches Antigravity-specific behavior and record the result there.

## Operator steps

Standing by-hand procedures, never pending items (those live in the newest handoff):

- On the Mac after a push from the Linux box: `sh-agent update`, `sh-agent status`, then the three smoke checks on each CLI (`.agents/docs/testing/2026-09-03-verify-real-home.md`; the Mac procedure is `2026-09-04-verify-mac-real-home.md`).
- After any change to `context/`, `skills/`, or `agents/`: refresh argolis-devbox's mirror from its project; styrman-crew receives the next signed tag (§ Releases).
- On a new machine, `ssh-keygen -y -f ~/.ssh/id_ed25519 | ssh-keygen -lf -` must print the fingerprint in § Releases before that machine cuts a release (derived from the private key: a copied `.pub` alone proves nothing).
- Go upgrades on the Mac: `brew upgrade go`.

## Roll back

`git checkout <last good commit> && go run ./cmd/sh-agent install` on the affected machine; `sh-agent update` will refuse to move until the branch is fast-forwardable again, so push a fix or reset the branch first.

## Access

Remote boxes are reached only on the operator's explicit request with details they supply. Nothing in this repo assumes a host.
