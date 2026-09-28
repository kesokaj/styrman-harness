---
owner: reviewer
date: 2026-09-28
---
# Cold-read review, round 2 — `.agents/docs/plans/2026-09-28-update-from-zip.md`   (2026-09-28)

Range reviewed: plan file `.agents/docs/plans/2026-09-28-update-from-zip.md` and the code surfaces it names
(`internal/harness/source.go`, `repo.go`, `apply.go`, `state.go`, `cmd/sh-agent/main.go`, `REQUIREMENTS.md` R8/R9,
`.agents/docs/architecture/DECISIONS.md` ADR list, `.agents/docs/specs/sh-agent/delivery.md`, `README.md`,
`.agents/docs/deploy/{local,remote}.md`). Round-1 report read only after forming my own view, per the dispatch.

## Status: COMPLETE

## Verdict: design needs revision (round-1's blocker is fixed; one new blocking gap found, one round-1 non-blocker still open)

Round 1's high-severity finding (the symlink-refusal scope contradiction) is resolved: D6 and the current
"Checks run at two scopes" text in `## Interfaces` § `UnpackZip` explicitly scope the symlink/mode/size checks to
the five extracted trees and keep name-safety/entry-count archive-wide, with `.agents/skills` named as the case it
exists to admit. Round 1's findings 2 and 3 (Origin-vs-Dir precedence in `apply.go`, and `status`'s `LoadState`
reordering) are also now spelled out explicitly in the current `## Interfaces` text. Good revision.

But the plan still has a real hole: nothing in it gives a zip-installed home, running a non-release (development)
build with no local checkout, any way to run `status`, `list`, `render`, `diff`, `install`, or a plain `update`
afterward — every one of those calls `harness.ResolveSource` before reaching the command's own logic, and
`ResolveSource`'s only change is to *skip* a `zip:` manifest value as a directory candidate (repo.go:75), never to
succeed because of it. On a machine with no reachable checkout and a dev-versioned binary, `ResolveSource` falls
through to `"harness repo not found and ..."` and the command errors out before printing anything — including
`status`, whose whole job for this feature is to report the zip install. The plan's own Scope line promises
`--from-zip` "from any binary, release or development build," and its Golden-path E2E criterion sidesteps the gap
by using "a binary built with a canonical release version" (release binaries short-circuit `ResolveSource` via the
embedded-files branch regardless of the zip marker) — so the untested path is exactly the dev-build one the scope
claims to support.

## Findings

- **[high / confidence: high]** `internal/harness/repo.go:22-88` (`ResolveSource`) vs. plan `## Interfaces` §
  `ResolveSource` ("skips it as a candidate") — for a zip-installed home whose binary is NOT a release version (no
  local checkout reachable, no `--repo`/`$SH_AGENT_REPO`), every command that calls `ResolveSource` before its own
  logic (`status` at `main.go:184-185`, `diff`, `install`, `list`, `render`, plain `update`) fails with
  `"harness repo not found and ..."` before doing anything. The plan only removes the `zip:` value from the
  candidate walk (repo.go:75); it adds no success path for this case. This directly contradicts Scope's claim that
  `--from-zip` works "from any binary, release or development build" — a dev build that installs from a zip cannot
  run `status` on itself afterward. The plan's own Golden-path E2E (Phase 3) uses a release-versioned build, which
  masks the gap because the release branch of `ResolveSource` never consults the manifest at all. **Must decide**:
  either `ResolveSource` gains a zip-aware success branch (e.g., a `zip:`-prefixed `st.Repo` on a non-release binary
  returns some usable `Source`, or a dedicated error naming `--from-zip` as the way to reinstall), or the plan
  narrows Scope to "release binaries only" for the parts of the feature that run after the initial `--from-zip`
  call, and says so.
- **[medium / confidence: medium]** Plan `## Phases`, Phase 3, task 2's verify clause: `grep -iE 'two ways|no third'
  over the living docs prints nothing that describes the delivery model` — checked against the current text of the
  four docs the plan itself lists as surfaces to touch, and the pattern misses restatements that don't use the
  literal words "two ways" or "no third": `README.md:182` ("`sh-agent` is the only way the harness reaches a
  machine (R8): a checkout ... or a release binary ...") and `.agents/docs/deploy/remote.md:61` ("it is no longer
  one of the ways in") both restate the two-ways delivery model without the grepped phrase, and would read as clean
  even if left unedited. This is the exact failure mode `AGENTS.md`'s own describing-text rule warns about (grep
  subject nouns, not old wording) — the fix is a broader pattern (e.g. `-iE 'two ways|one way|only way|no third|
  two alternatives'`) or a manual read of the four listed files rather than trusting this one grep as sufficient
  before the Phase-3 checkpoint.
- **[low / confidence: medium, carried from round 1, still open]** `.agents/docs/specs/sh-agent/delivery.md:39`,
  `README.md:182`, `.agents/docs/deploy/local.md:45-55`, `.agents/docs/deploy/remote.md:61` all describe a
  pre-existing, already-shipped "plain copy (no `.git`) — the unpacked zip" install path: download the repo zip by
  hand, unpack it, run `sh-agent install --repo <dir>` from a dev build. That is materially the same operator
  problem this plan solves (no SSH/`gh`), through a different, already-documented mechanism that leaves the tree on
  disk permanently rather than deleting it. The plan lists these four files as surfaces to touch but never states
  whether the old path is deprecated, kept as a parallel fallback, or how the two should be distinguished in the
  docs. Not fixed since round 1.

## Scope

Everything in `## Scope: In` still traces to a D-row or the operator's "approve as written, go all but number 6."
No unrequested endpoint, flag, or configuration surface found. One scope-adjacent observation, not a violation:
Scope's claim "from any binary, release or development build" is broader than what the plan's own design (see
finding 1) currently supports end-to-end — that's a design gap, not unrequested scope, since the operator's request
implies "works after install too," which the plan doesn't yet deliver for a dev build.

## Readiness

1. **BLOCKING** — How does `status` (and every other `ResolveSource`-calling command) work on a zip-installed home
   running a development build with no local checkout? (Finding 1. Two implementers diverge: one leaves it broken
   as an accepted limitation, the other adds a new `ResolveSource` branch — different code, different docs claim.)
2. **NON-BLOCKING** — Is the Phase-3 `grep -iE 'two ways|no third'` check trusted as sufficient, or does whoever
   runs Phase 3 also manually re-read the four listed doc files? (Finding 2. A careful implementer reads the files
   anyway and catches the missed restatements; the plan's stated verify clause alone would not.)
3. **NON-BLOCKING** — Does the docs update (Phase 3) mention the pre-existing plain-copy/unpacked-zip path at all,
   and how it relates to `--from-zip`? (Round-1 finding 4, still open. A reasonable `documentation` pass can resolve
   this inline without asking first, but it's undecided.)

## Not examined

- `.agents/docs/plans/2026-09-28-recon-from-zip.md` (the recon this plan cites) — not read; the plan is expected to
  stand on its own for this check.
- `.agents/docs/testing/skill-matrix.md` and `reviews/2026-09-28-review-v1.7.0-381e86d-code.md` findings 1-3 (Phase
  4) — not read in detail; Phase 4 is downstream of the Phase 1-3 design questions this review focused on.
- `cmd/sh-agent/main_test.go` existing test patterns, to check `TestUpdateFromZip`/`TestStatusAfterZip`/
  `TestKnownCheckoutZip` naming conforms to existing convention — not read.
- `README.md` § How it works (the map) itself, and its dashed/solid line convention (AGENTS.md rule 4) — not opened.
- ADR-034's full text (only the portion quoted in `repo.go`'s comments was read) — the plan's claim that ADR-038
  "supersedes ADR-034 for this one path only" was not verified against ADR-034's full wording.
- `internal/harness/embed.go`, `release.go`, `report.go` — not read in this round; not central to the new findings.

design needs revision
