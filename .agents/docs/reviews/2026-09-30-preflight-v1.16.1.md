---
owner: main
date: 2026-09-30
---
# Preflight — v1.16.1   (2026-09-30)

Range reviewed: v1.16.0..7bf0978

Release request: the operator, 2026-09-30, in their own turn: "Before you cut a new release. Do an audit if this happens on other places first. Then preflight a new release, update all docs so this is known." Then, on the list of what was left: "Ok, do all of these." This is a PATCH. It fixes released behaviour: `/init-project` wrote `.claude/` on a machine installed for Antigravity alone (BUG-113), and `/dream` resolves its store on such a machine. Nothing is added, and nothing an existing place must react to changes; a place installed `all` behaves as before.

## Findings, worst first

None blocking.

- **low / high**, docs lens 1 and 2: on an Antigravity-only machine that still holds Claude Code memory at the default path, `/dream`'s hint "run /init-project to move them" points at a migration init skips there. Advisory (Step 4). Owner: the main agent, handoff Next Steps (`reviews/2026-09-30-review-v1.16.0-7bf0978-docs.md` § Outcome).
- **process, recorded**: the first agy-only sandbox runs used a copy of the real home's Antigravity settings (BUG-115). They were replaced by compliant runs in this gate (`testing/2026-09-30-verify-real-home-v1.16.1.md`).
- **process, noted**: the docs-lens reviewer detached HEAD in the operator's checkout and reattached it. The state was checked clean afterwards (same Outcome).

## Rows

- **Reviewer**: `reviews/2026-09-30-1721-review-agy-only-init.md` covers the source of the range (both skills, the context line, the test, the spec, the matrix). Its REQUEST CHANGES rested on one high finding: the context line took the Antigravity block to 22050 bytes. That was fixed in 7bf0978 (21998 bytes, `TestRenderedContextFitsAntigravityBlock` fails on the reviewed text and passes on the fix; BUG-114), and the fix is recorded in its Outcome. The rest of the range is prose docs (README, `deploy/local.md`), covered by the docs lens, plus records.
- **Security**: no row. No changed path touches auth, crypto, secrets, IAM, input parsing or dependencies. The skill's new grep only reads the operator's own install manifest.
- **Verification**: `testing/2026-09-30-verify-v1.16.1-gate.md` PASS at 7bf0978: 375 pass, 0 fail, 2 skips with reasons; `TestInitProjectVendorCheckMatchesManifest` and the context-size test pass; scratch installs `--vendor agy` and default both show no drift; three cross-builds; govulncheck clean; installed copies byte-identical. The fix's own log is `testing/2026-09-30-1721-verify-agy-only-init.md`, PASS.
- **The project's own blocking condition** (`AGENTS.md` § Release): `testing/2026-09-30-verify-real-home-v1.16.1.md`, PASS. It covers the real home at 7bf0978 (the context check on both CLIs, with the Antigravity store check) and one run of each changed skill (`init-project`, `dream`) on each CLI along the `all` path. It also covers both skills on Antigravity in a home installed `--vendor agy`, along the new path. No auto-denied line, no invalid allow entry.
- **Docs freshness** (Step 4): `reviews/2026-09-30-review-v1.16.0-7bf0978-docs.md`, APPROVE WITH NITS. No living text still claims the old behaviour; the two lows are above. The docs pass the operator asked for ("update all docs so this is known") changed README, `deploy/local.md`, `specs/sh-agent/memory-loop.md`, `context/AGENTS.md` and the skill matrix (`changes/2026-09-30-1725-init-project-agy-only.md`).
- **Trail freshness** (Step 3): pass. The newest `changes/` entry before the release entry describes the range's work. No future `date:` (the one hit is the template's `YYYY-MM-DD`). No modification to `changes/` or `DECISIONS.md` in the range. `ROADMAP.md` and `TODO.md` are untouched and empty.
- **Requirements**: untouched.
- **Shadowing agents**: none (no project agent under `.agents/agents/` or `.claude/agents/`).

## What was not checked

- The Mac and the devbox: they pick the release up with `sh-agent update`; the binaries are built, not run there.
- A machine installed `--vendor agy` from a release binary, as spgo will be: this gate installed the agy-only home from the checkout. The skill text is the same bytes either way (the embed test pins it).
- Interactive sessions on either CLI.
- The `/dream` hint path of the docs-lens lows.

Every item above was read against the blocking conditions, and none triggers one: a verifier log that is FAIL, BLOCKED or IN PROGRESS (none); a high security finding (no row owed); a missing security report that is owed (none owed); a standing reviewer REQUEST CHANGES or high finding (the one high is fixed and recorded in the Outcome); a `REQUIREMENTS.md` violation (none); a shadowing agent (none); the project's own condition, an Antigravity run of a changed skill named here as not done (both changed skills ran on Antigravity, on both paths).

VERDICT: SHIP
