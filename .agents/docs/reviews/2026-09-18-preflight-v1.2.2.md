---
owner: main
date: 2026-09-18
---
# Preflight — v1.2.2     (2026-09-18)

Range reviewed: v1.2.1..f91257c

Base: the last tag, `v1.2.1` (release commit ea7a202). Head: f91257c, `main`, clean tree (`git status --porcelain --untracked-files=all` empty). Eight commits; source diff 6 files, 410 insertions, 15 deletions — the BUG-078 work: `internal/harness/pathguard.go` (new) and its test, `cmd/sh-agent/main.go` and its test, `README.md`, `ROADMAP.md`. Version: PATCH — it fixes released behaviour and adds no surface a user must react to.

## VERDICT: SHIP

Nothing blocks. Every non-blocking finding below carries the decision made about it.

## Findings, worst first

None high. None medium.

- [low / high] **Manual verify chains do not assert the version marker** — `README.md` § Install without cloning (`gh` route) and `remote.md` § First install strip `# styrman-harness vX.Y.Z` with `grep -v '^#'` and never check it, while `verifyRelease` in `release.go` does; an actor who can serve assets but not sign could hand the human route an older, genuinely signed release. (security, `reviews/2026-09-18-security-v1.2.1-f91257c.md`, finding 1; raised as a hardening note in the round-1 security report and not dispositioned then — an omission of mine.) **Decision: accepted for this tag, fixed in the commit immediately after the push.** README and the deploy guides are in no release artifact — the binary embeds `context/`, `skills/`, `agents/`, `commands/`, `rules/` — so the fix reaches every reader of `main` the same moment whether it lands before or after the tag, and a source change in the release commit would fail the records-only check in `remote.md` § Releases. The fix binds the binary to the signed marker without the reader knowing the version in advance: `test "$(./$ASSET version)" = "sh-agent $(sed -n 's/^# styrman-harness //p' SHA256SUMS)"`.
- [low / medium] **The browser fence works in `~/Downloads`** — `> allowed_signers` truncates and the final `rm` deletes three fixed names a user may hold there for something else; the `gh` route avoids this with `mktemp -d`. (reviewer, `reviews/2026-09-18-review-b029ea1-f91257c.md`, finding 1.) **Decision: accepted for this tag, fixed in the same post-push commit** — move the three files into a `mktemp -d` first, one-chain shape kept. Same reasoning: README ships in no artifact.
- [low / medium] **The browser fence's `shasum -a 256 -c - --ignore-missing` has never run on a Mac** — written from the manual; BUG-075's rule says a documented command is run before the release that documents it. Both GNU `sha256sum` and Perl `shasum` 6.x fail closed here (exit 1 on empty input and on no file verified — measured by the reviewer and by security independently); macOS's shipped version is not measured. (security, finding 2.) **Decision: accepted and named in the release entry; owner: the operator, one `shasum --help | grep ignore-missing` on the work Mac settles it** (handoff Next Steps). Not fixed by removal, because the route with no verification at all was the previous state and worse.
- [low / high] **`status` prints no `shadowed` line when the installed file is missing or not executable**, while both deploy guides make "no `shadowed` line" a Verify pass condition; the uncovered case is a present-but-non-executable file with a healthy `binary` line. (security, finding 3; `left` twice before.) **Decision: left, owner main agent** — a shadow of a file the shell cannot run is not a shadow; the honest gap is the Verify wording, filed with the docs follow-up below.
- [low / medium] **`testing/guide.md` § 5 does not ask for the unqualified `sh-agent version` or "no `shadowed` line"** that both deploy guides' Verify sections now do. (docs lens, `reviews/2026-09-18-review-v1.2.1-f91257c-docs.md`, finding 1.) **Advisory; owner `documentation`; filed in the handoff's Next Steps.** No plan is active.
- Informational: the trail carries internal hostnames and home paths (`spgo-mac.roam.internal`, `/home/admin`, `/Users/spgo`, `/opt/devbox/…`) — fine in a private repository; `remote.md` § Releases already says publication is a separate decision. Repository visibility checked before this gate: `PRIVATE`.

## Rows

| Row | Result | Artifact |
|---|---|---|
| Verification | PASS — gofmt, vet, `go test -race` 215/0/0 (unchanged vs the last BUG-078 log), `govulncheck` none, `TestEmbedHoldsOnlyTrackedFiles` pass, embed-dir contamination guard empty; live off-PATH install/status lines quoted | `testing/2026-09-18-verify-preflight-v1.2.2.md` |
| Review — code through b029ea1 | APPROVE (round 3 of the BUG-078 work; rounds 1–2 REQUEST CHANGES, each with `## Outcome`) | `reviews/2026-09-18-review-4cdc0b5-b029ea1.md` (covering; its predecessors `…-52574a0.md`, `…-8ce58a0.md`) |
| Review — commits after b029ea1 (opus) | APPROVE, 1 low (above); the `!dryRun` gate proved by mutation on a scratch copy; both `sha256sum`/`shasum` degenerate exits probed | `reviews/2026-09-18-review-b029ea1-f91257c.md` |
| Security (opus), whole range | FINDINGS, all low (above); 27 paths read end to end; `## Not examined` names no release-diff path; signer line byte-identical in all seven places; `govulncheck` re-run by the agent | `reviews/2026-09-18-security-v1.2.1-f91257c.md` |
| Docs lens | APPROVE WITH NITS, 1 low (above); the README map untouched across the range — `AGENTS.md` rule 3 holds, no line rerouted | `reviews/2026-09-18-review-v1.2.1-f91257c-docs.md` |
| Trail freshness | pass: newest `changes/` entry is `2026-09-18-bug-078-path-shadow.md`, the work in this diff; no future `date:`; no modified past entry in `changes/` or `DECISIONS.md` since v1.2.1; `ROADMAP.md` touched once (Now emptied), `TODO.md` untouched and empty | — |
| Shadowing project agent | none (`.claude/agents`, `.agents/agents` absent) | — |
| `REQUIREMENTS.md` | untouched since v1.2.1; no violation named by any row | — |
| Signing key | present, fingerprint `SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU` derived from the private key | — |

## What was not checked

- macOS: no command in this range was run on a Mac — the browser fence's `shasum` form, `mktemp -d` cleanup behaviour, and `install`/`status` output there. All Linux-verified.
- The `gh` and browser install chains were read and their exit-status semantics probed on this box; neither chain was executed end to end against the real release page (the release does not exist yet at gate time).
- The four dated review reports and the bug record were not re-read by this gate beyond their verdict lines and `## Outcome` sections.
- `linux/arm64` and the swapped-signature edge case on darwin remain unexercised, as the plan's Follow-ups record; neither is changed by this release.
- No Dockerfile in the range; no container tier.
