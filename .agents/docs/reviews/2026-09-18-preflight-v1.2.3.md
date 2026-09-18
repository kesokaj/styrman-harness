---
owner: main
date: 2026-09-18
---
# Preflight — v1.2.3     (2026-09-18)

Range reviewed: v1.2.2..30f7488

Base: the last tag, `v1.2.2` (release commit f3ec5dc). Head: 30f7488, `main`, clean tree. Fourteen commits; source diff 21 files, 1303 insertions, 88 deletions: BUG-079 (four commits), the by-hand install route's version binding, Phase 1 (`install` names every change), Phase 2 (the audit and its three fix rounds), Phase 3's matrix, of `plans/2026-09-18-sh-agent-works-everywhere.md`. Version: PATCH — every change fixes released behaviour or adds output; no surface an existing place must react to (D3 of the plan).

## VERDICT: SHIP

Nothing blocks. Every non-blocking finding below carries the decision made about it.

## Findings, worst first

None high. None medium.

- [low / high] **An unwritable `~/.config/sh-agent` (e.g. 0500) passes the tightening and `install` deploys all 260 files before failing at `SaveState`** — a fully deployed home with no manifest, which `uninstall` then cannot clean. Pre-existing (measured at e7fec88 too); in range because `TightenConfigDir` sits at the top of `Install`. (security, `reviews/2026-09-18-security-v1.2.2-30f7488.md`, 1.) **Decision: accepted for this tag, filed as a Follow-up** (a writability probe that refuses BEFORE any write, with a test). A code fix now would re-open the verification row; the case needs a directory the user made unwritable on purpose.
- [low / medium] **`yamlQuote` misses U+0085, U+2028/2029 and the C1 range**, which YAML 1.1 parsers read as line breaks — the residual of the audit's finding 16. (security, 2.) **Decision: accepted, Follow-up** (one guard, one test); the reachable inputs are skill and agent descriptions the installer already requires to be one line.
- [low / high] **Two U+0646 characters survived in `reviews/2026-09-18-security-c868ed4-bc745dc.md`** (the report quoting the brief that was fixed). Trail text agents read back; not shipped (embed excludes `.agents/docs`). (security, 3.) **Decision: fixed in the release commit**, removed by character code; `grep` for the character reports 0.
- [low / medium] **`delivery.md:51` says every install tightens the runtime dir when it is the caller's own, omitting the symlink case** — README's twin sentence has both. (reviewer `reviews/2026-09-18-review-42840fb-30f7488.md`; docs lens `reviews/2026-09-18-review-v1.2.2-30f7488-docs.md`.) **Decision: fixed in the release commit** (a docs file under `.agents/docs`, which the records-only check allows).
- [low] a failed `filepath.Abs(home)` is swallowed (only when `Getwd` fails) — **accepted, Follow-up**; the tightening is install-only, so a place that only ever runs `report` keeps a loose dir — **accepted on record** (every `update` runs an install); the doc chains' `grep -qxF` marker and `release.go`'s marker regex agree but nothing pins them — **Follow-up** (a test).

## Rows

| Row | Result | Artifact |
|---|---|---|
| Verification | PASS — 250/0/0, three cross-builds, `govulncheck` none, `TestEmbedHoldsOnlyTrackedFiles`, the embed contamination guard empty, live install 260 `add` then none with `700` and drift none, § 3c matrix re-run: 53 cells matching the matrix verbatim | `testing/2026-09-18-verify-preflight-v1.2.3.md` |
| Review — every code commit but the last | APPROVE / APPROVE WITH NITS with `## Outcome`, per slice | `reviews/2026-09-18-review-c868ed4-{bc745dc,0f2da25,a078848}.md`, `…-review-b393525-028a43b.md`, `…-review-028a43b-d737282.md`, `…-review-160aa16-3209cf3.md` (opus), `…-review-3209cf3-e7fec88.md`, `…-review-e7fec88-42840fb.md` |
| Review — the post-cap commit (opus) | APPROVE; every guard order confirmed; one low (above) | `reviews/2026-09-18-review-42840fb-30f7488.md` |
| Security (opus), whole range | FINDINGS, six low (above); 80 paths read end to end; `## Not examined` names no release-diff path; signer line byte-identical in seven places; signature, marker and checksum gate sound end to end; `govulncheck` clean | `reviews/2026-09-18-security-v1.2.2-30f7488.md` |
| Docs lens | APPROVE; every quoted string matches the code; the README map byte-identical across the range — rule 3 holds | `reviews/2026-09-18-review-v1.2.2-30f7488-docs.md` |
| Trail freshness | pass: newest `changes/` entry is `2026-09-18-sh-agent-audit.md`, the work in this diff; no future `date:`; no modified past entry in `changes/` or `DECISIONS.md` since v1.2.2; `ROADMAP.md` touched once (Now names the plan), `TODO.md` empty | — |
| Shadowing project agent | none | — |
| `REQUIREMENTS.md` | untouched since v1.2.2; no violation named by any row | — |
| Signing key | present, fingerprint `SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU` from the private key | — |
| Repository visibility | `PRIVATE` | — |

## What was not checked

- No command in this range was run on a Mac or an arm64 machine; the by-hand browser chain, Gatekeeper and `shasum --ignore-missing` remain unrun there (the plan's Phase 5 and remainder).
- The not-owned branches of `NoteRunningElsewhere` and `TightenConfigDir`, and `TightenConfigDir`'s chmod-failure note, are unreachable by a non-root test and were reasoned from the source, not observed.
- `owned_other.go` (the `!unix` twin) is compiled for windows, never run; no windows asset ships.
- The dated reviews and logs of this range were read for their verdict lines and `## Outcome` sections, not re-read whole.
- No Dockerfile in the range; no container tier.
