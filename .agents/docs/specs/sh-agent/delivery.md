---
owner: main
date: 2026-09-17
---
# Delivery and the feedback loop — what this harness is building toward

The picture in one place, because the operator is discarding the agent's memory and the trail has to carry everything on its own. Written 2026-09-17. What is DECIDED lives in `REQUIREMENTS.md` R8 and R9, ADR-030 and ADR-033 (§ The model below is their one-paragraph statement); what is PLANNED lives in `.agents/docs/plans/2026-09-17-delivery-rework.md`; this document is the narrative that makes both readable to someone arriving cold.

## The problem, in the operator's words

> "I want the delivery to be around the sh-agent only. Nothing else. I'm the sole user of this harness and decide everything but I install it on multiple places and when I decide to update the repo I want a smart way to push any updates to places where I use it."

And the loop back:

> "Some sort of feedback function from each place installed reporting back to the main is good for a feedback loop so we can take experiences from each installed place to use to make the harness better."

## The four places

| Where | CLIs | How it is driven |
|---|---|---|
| Work Mac | Antigravity only | The operator runs `sh-agent update` |
| Home Mac | Claude Code and Antigravity | The operator runs `sh-agent update` |
| A Linux project host | both | The operator runs `sh-agent update` |
| A second Linux project host | both | The operator runs `sh-agent update` |

The normal flow is deliberately dull: the operator edits this repository, is happy with it, then goes to each place and runs one command. There is no push, no timer that has to work, no service that has to be up. A timer is a convenience some machines add; it is not the mechanism.

## The model

The one statement of what a tag delivers, to which kind of place, and how a place moves forward; every other living file points here (project `AGENTS.md` § How the harness works, rule 6). A **checkout place** follows `main`: `sh-agent update` fetches, refuses a dirty checkout, fast-forwards, rebuilds the binary from the checkout and reinstalls; a tag changes nothing for it. A **binary-only place** moves only when a new `vX.Y.Z` tag is published: `sh-agent update` downloads the newest release's binary (a canonical `vX.Y.Z` tag, newer than the running one), verifies `SHA256SUMS.sig` against the signer line compiled into the binary and the binary against its checksum, replaces `~/.local/bin/sh-agent` and installs the harness embedded in it. A tag is therefore a release marker that also carries the binaries a binary-only place updates to; nothing is ever pushed to any machine, and a tag without its assets is a marker no place can update to (R8, R9; ADR-033, superseding ADR-030 D2).

## How a machine gets the harness

Two ways, and no third:

1. **A checkout the operator owns**, kept current with `sh-agent update` (§ The model). This is the normal path.
2. **A release binary**, installed with `sh-agent install` and moved forward by `sh-agent update` (§ The model). One downloaded file carrying `context/`, `skills/`, `agents/`, `commands/` and `rules/` inside it, for a place that only USES the harness: the operator wanted no `.agents/docs` trail of this repository sitting beside their other projects, where a session opened there could mistake it for the project's own. What `update` does there is § The model. That place changes only when a release is cut (ADR-033).

A plain copy with no `.git` — the unpacked zip — still installs, and `update` refuses it with a message naming `install` still deploying it as it is and the two alternatives: a checkout, or a release binary; it is no longer one of the ways in: it is the fallback for a machine that can reach neither of them.

Authentication is the operator's own SSH key and `gh` login, present on every machine because they log in as themselves. That is why no deploy keys, tokens or per-machine credentials exist anywhere in this design — they would restrict nothing and add key management. The accepted cost is theirs, stated plainly: "if I don't login in with ssh or gh I can't get updates." A release download needs that same credential, because the repository is private; the plain copy is the fallback when that bites.

## What keeps consumers from changing the harness

Part shape, part credential. Only the four checkouts are repositories, and they belong to the operator. Everything downstream of them receives installed FILES: `sh-agent install` writes the context, skills, agents, commands and rules into a home and never creates a checkout there. But a checkout is not what changing the harness takes — a credential that can push to this repository is, and anything that reaches the operator's SSH key or `gh` login can clone and push, after which every place's `update` deploys it. So a home seeded for agents is read-only only when it reaches no such credential: no forwarded ssh-agent, no `gh` token, no readable key. The shape is provable here; the credential half is each host's to hold (security review 2026-09-17, finding 5).

State the shape as a property of `sh-agent install`, provable from this repository. Do not rest it on how any particular consumer seeds its homes — that was done once, in ADR-030 D4, and it is a dependency pointing the wrong way: if that consumer changes its entrypoint, a claim in this repository's decision record silently becomes false and nothing here would catch it. What `install` did is on its output, every run: one line per file it added, updated (including a previous file it could not read, written over as before), corrected the mode of, removed, or kept (a managed path already holds a file sh-agent did not write, left alone and counted separately), with the unchanged as a count; `--dry-run` prints the same lines and writes nothing (2026-09-18, the operator's "more info").

## The loop back

`sh-agent report` sends a COPY of one machine's experience to a private repository, `styrman-harness-feedback`, one directory per machine per run: the corrections log, the retro proposals, the cross-project index, and a metrics record naming the host, the installed commit or release, drift, the skill and agent counts and the toolchain that built the binary. It keeps nothing: the clone it writes through holds only that machine's own folder and is removed after the push, so no place can read another place's reports — which is what the operator found when they asked "But why can I see dave reports?".

It refuses any checkout whose origin is not that repository, and pushes to the feedback URL itself rather than to `origin` — so `origin`'s own `pushurl` or an extra `url` cannot redirect it, though a git URL rewrite (`insteadOf`, `pushInsteadOf`) or a remote named after that URL still applies, as it does to every push on the machine, trusted as the operator's own configuration rather than refused — and redacts before sending, structurally, with no JSON walk. Redaction is a safety net, not a guarantee: the guarantee is the operator reading what was added under `~/.config/sh-agent/` before each report. The credential needles — vendor token prefixes, and case-insensitive password, secret, token and api-key assignments, `Bearer` and authorization headers written `Authorization: …` (colon form, a quoted header name included; a header set by a call such as `Header.Set("Authorization", …)` is an unlisted shape), and URL credentials (including an empty user) — LOCATE a secret and never define how much is removed — a whole text line, a whole JSONL record (replaced by a JSON redaction marker), a whole `-----BEGIN`-to-`-----END` block, or the whole of `state.json` goes. This is the design the operator chose at the cold-read cap, and the reasoning is worth keeping: `AKIA` is a prefix, so replacing the matched token strips the marker and ships the key; a private key's header matches while its body carries no needle at all. Any span shorter than a structural unit leaks while reporting success.

Originals never leave `~/.config/sh-agent/`. R7 keeps that as the one place harness state lives; a copy going to the operator's own private index is not a second home for it, and R8 says so explicitly so the two rules do not contradict. `~/.config/sh-agent/` is created 0700 and every install tightens it to 0700 when it is looser and the directory is the caller's own and not a symlink (either case gets a `note:` instead); `diff` and `report` create it 0700 when absent but leave an existing one's mode alone.

## What is deliberately NOT being built

- **A webservice**, for delivery or for feedback. Git already carries both halves. A service adds hosting, per-machine authentication and an attack surface, for one user, against R3, R4 and local-by-default. A dashboard over collected data stays possible later without changing how any machine reports.
- **A signed-tag delivery channel.** ADR-026 designed one on 2026-09-16 to give a consumer host a trust root; ADR-030 withdrew it the next day because the operator authenticates as themselves everywhere, so it solved a problem that was not there. Tags survive as release markers — `v1.0.0` stands, and the versioning rule and signing recipe stay in `deploy/remote.md`; what a tag delivers today is § The model.
- **Vendored mirrors.** Two projects kept a `.git`-less copy refreshed by hand, which is why a mirror refresh kept appearing as a task in three repositories at once. They retire. Each consumer's own session does that work in its own repository.
- **Per-machine keys of any kind**, for the reason above.

## The rule that keeps this repository independent

> "styrman-harness and the feedback/trajectory should not be directly connected to any other project, if any other project wants to use the harness they need the sh-agent, that's it."

This repository names no consumer project in its living docs and depends on none. Knowing a consumer exists is fine; depending on its internals is not. The test is simple: **would the harness still be correct if every consumer vanished tomorrow?** If a document here needs another repository to be true, it is wrong, whatever it says.

This is easy to get wrong when several of those repositories sit side by side on one disk, because reading them costs nothing. It was got wrong once already, on 2026-09-17, in the same paragraph that stated the rule.

## Where things stand

**Built and shipped:** `sh-agent` with install, update, status, diff, uninstall, list, render, adopt, report and version; `update` refusing a dirty checkout on every path; the installer refusing secret-shaped files and duplicate skill names; `v1.0.0` tagged and signed. 27 skills, 7 agents.

**The loop back is built:** `sh-agent report` since 2026-09-17, run end to end from the Linux box into the private `styrman-harness-feedback` repository (`.agents/docs/testing/2026-09-17-delivery-rework-e2e.md`). How to run it and what fails: `.agents/docs/deploy/local.md` § Report back. Not yet exercised on real data: a redaction firing, and any place but that one box.

**Open, defaulting to no:** whether this repository becomes public. Only that would let a machine bootstrap with no credential at all, and it is a question about content rather than packaging.

**Built, 2026-09-17:** what was here as "a follow-up with its own plan — versioned release bundles" is now the release binary above (ADR-033, `.agents/docs/plans/2026-09-17-release-binary.md`). It needed no CI: the binary cross-compiles locally in `/preflight`, which signs the checksums with the same key as the tag. It does not remove the credential requirement while the repository is private.
