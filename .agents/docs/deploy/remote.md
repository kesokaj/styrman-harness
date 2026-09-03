---
owner: main
date: 2026-09-03
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

## The resync loop

```bash
# machine A: change, verify, install, commit, push
sh-agent install && git commit -am "context: tweak" && git push
# machine B
sh-agent update            # fetch; if behind: pull --ff-only, rebuild, reinstall
sh-agent update --check    # only reports; exit 2 when behind
```

`sh-agent update` finds the checkout from `$SH_AGENT_REPO`, the working directory, or the path recorded at the last install, so it works from anywhere. It refuses a non-fast-forward pull; resolve that by hand in the checkout.

## Verify

Same as `local.md`: `sh-agent status` with no drift, then one start of each CLI. On the Mac, additionally confirm the Antigravity paths (`TODO.md`, high priority).

## Roll back

`git checkout <last good commit> && go run ./cmd/sh-agent install` on the affected machine; `sh-agent update` will refuse to move until the branch is fast-forwardable again, so push a fix or reset the branch first.

## Access

Remote boxes are reached only on the operator's explicit request with details they supply. Nothing in this repo assumes a host.
