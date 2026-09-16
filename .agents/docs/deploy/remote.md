---
owner: main
date: 2026-09-15
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
sh-agent update            # fetch; if behind: pull --ff-only, rebuild, reinstall
sh-agent update --check    # only reports; exit 2 when behind
```

`sh-agent update` finds the checkout from `$SH_AGENT_REPO`, then the working directory or any of its parents, then the path recorded at the last install, so it works from anywhere. It refuses a non-fast-forward pull; resolve that by hand in the checkout.

## Downstream mirrors

Two projects vendor this repo as a .git-less copy: styrman-crew `image/harness/` and argolis-devbox `deploy/harness/`. Each mirror's `VENDORED.md` holds the refresh recipe: rsync from a clean checkout, record `COMMIT` and `SHA256SUMS`, then diff to confirm only those two files plus `VENDORED.md` differ. What happens after the refresh differs by project — styrman-crew builds `sh-agent` inside its own image (`image/Dockerfile`, builder stage) and ships a refresh by cutting a release, since the image version stamp is what makes a running project re-seed; argolis-devbox commits and pushes with no version to bump, and its VM builds `sh-agent` on the box itself, as root, via `deploy/scripts.d/032_setup-harness.sh` on its next auto-update (within five minutes). A mirror is a plain copy either way: `install` from it never rebuilds (the tree may be read-only to the installing user — BUG-015), `sh-agent status` there shows no commit, and `update` refuses, by design. After any change to `context/`, `skills/`, or `agents/` here, refreshing both mirrors is a task in each of those projects; a retro that applies a harness-scope proposal says so in its report.

## Verify

Same as `local.md`: `sh-agent status` with no drift, then one start of each CLI. On the Mac, additionally confirm the Antigravity paths — passed 2026-09-04 (`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`); re-run after a change that touches Antigravity-specific behavior and record the result there.

## Operator steps

Standing by-hand procedures, never pending items (those live in the newest handoff):

- On the Mac after a push from the Linux box: `sh-agent update`, `sh-agent status`, then the three smoke checks on each CLI (`.agents/docs/testing/2026-09-03-verify-real-home.md`; the Mac procedure is `2026-09-04-verify-mac-real-home.md`).
- After any change to `context/`, `skills/`, or `agents/`: refresh the downstream mirrors from their own projects (section above).
- Go upgrades on the Mac: `brew upgrade go`.

## Roll back

`git checkout <last good commit> && go run ./cmd/sh-agent install` on the affected machine; `sh-agent update` will refuse to move until the branch is fast-forwardable again, so push a fix or reset the branch first.

## Access

Remote boxes are reached only on the operator's explicit request with details they supply. Nothing in this repo assumes a host.
