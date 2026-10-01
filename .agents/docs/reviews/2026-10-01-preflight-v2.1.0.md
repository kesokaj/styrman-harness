---
owner: main
date: 2026-10-01
---
# Preflight — v2.1.0   (2026-10-01)

Range reviewed: v2.0.0..159585a

Requested in the operator's own turn: "Continue through all phases, then /preflight v2.1.0". Version: MINOR by the rule — the release adds behaviour (one skill, `google-catalog`, with 136 Google-published skills behind it; a search marker; a listing-budget test) and removes nothing a v2.0.0 place used. One thing a place does after updating: where the nightly dream timer is installed, `sh-agent dream --remove-timer` then `--install-timer`, once (the unit's text changed).

## Findings, worst first

Nothing blocks. No verifier log at the gated head is FAIL, BLOCKED or IN PROGRESS; no security finding is high; every security report is COMPLETE with a `## Not examined` section that names no release-diff path; no reviewer verdict stands at REQUEST CHANGES; no `REQUIREMENTS.md` clause is violated; no project agent shadows a review agent.

What the gate found and fixed on the way (all in the range, each with tests seen failing first or failing with the fix removed):

| # | finding | severity / confidence | found by | state |
|---|---|---|---|---|
| 1 | The dream timer's "real home" was read from `HOME`; a sandbox run overrides `HOME`, so a sandbox `--remove-timer` or `uninstall` used the real timer's name (BUG-149). The first fix (`user.Current`) still fell back to `HOME` for an account outside `/etc/passwd` on a release binary | low, then medium / high | security slice at 7dd0645; code review at f62513b | fixed in f62513b and b4e3d79: lookup by uid, no fallback |
| 2 | The refresh script's swap deleted directories under any `--dest` (BUG-147), stopped silently half way (BUG-148), and its guard could be passed by a spelling: `..`, a symlink, a missing directory, an unreadable parent, an unset `HOME` | medium to low / high | security slices; code review rounds | fixed in 7dd0645, f62513b, b4e3d79, ff4c181 |
| 3 | Thirteen manifest notes read as an exception to a router rule (`--quiet` after approval, an MCP server "with approval", `gcloud config set project` once named) | low / medium | security reads at f62513b and b4e3d79 | fixed: the router and the manifest say a note never relaxes the router; notes reworded; `TestCatalogNotesNeverAllowQuiet` |
| 4 | The router's consent rule did not cover data leaving the machine | low / low | security slice at 7dd0645 | fixed in f62513b |
| 5 | Upstream names and paths printed raw in the script's refusals, one inside a hint the operator pastes | low / medium | code review round 3; security reads | fixed in f62513b, b4e3d79, ff4c181; one fixture per message in 159585a |
| 6 | A root `.sh-agent-noindex` marker was ignored, then left the global index serving the trail | low / high | code review rounds 1 and 2 | fixed in 7dd0645 |
| 7 | Verifier log `-r4` at ff4c181: FAIL on one coverage gap (two quoted messages with no test) | low | verifier | fixed in 159585a (test only); `-r5` PASS and says the finding is closed |
| 8 | Docs lens over the range: eight lows (a stale count, the timer row, two overstated sentences, the script's guards undocumented) | low | docs lens at 7dd0645 | fixed in f62513b; the release-state paragraph in this commit |

Left with an owner, each a low, recorded and not fixed before the tag (three fix rounds each moved the tree; the rule from the fourth pass on was: fix a medium or high, record a low). They are one task, the first after the release, and are in the handoff's Next Steps:

| # | left | owner | what the operator is exposed to until then |
|---|---|---|---|
| L1 | `--remove-timer`, `uninstall` and `status` do not say when they found nothing, and no longer find a timer that v2.0.0 or earlier wrote under the bare name from a home that is not the account's (a scratch home; every home of an account `/etc/passwd` does not list) | operator choice, recommendation: say what was found and print the cleanup command | such a timer keeps running nightly after a removal that reported success. Not this box: its real home has no timer installed. `deploy/local.md` § Operator steps has the by-hand procedure; `grep "^$(id -un):" /etc/passwd` tells which case applies |
| L2 | `--dest //<home>` passes the script's explicit guards (`pwd -P` keeps a leading `//`) | main agent | two mistakes needed: that spelling, and a `CHECKSUMS` file or an unlistable directory at the target |
| L3 | The `google-agents-cli-scaffold` note lets `--auto-approve`/`--force` through after a dry run; the `google-cloud-filestore-nfs-browser` note asks for a run without `--quiet` that the entry's script hard-codes | main agent | the router's steps 2 and 4 forbid both whatever the note says; a model that weighs the note above the router could overwrite scaffolding files, recoverable from git |
| L4 | `TestCatalogNotesNeverAllowQuiet` passes a note that puts an exception beside an allowed phrase, and pins only the literal `--quiet`, not other skip-prompt spellings or other kinds of exception | main agent, at the next refresh | a refresh's diff gets a security review (project `AGENTS.md` § Release), which is the control today |
| L5 | The permission-key lint is skipped for a command or rule whose frontmatter does not parse (`commands/` and `rules/` are empty) | operator choice, recommendation: refuse | none today |
| L6 | A `.sh-agent-noindex` marker hides a directory from search without a notice | operator choice, recommendation: one `noindex:` line in `status` | a stray marker under the docs root would empty `--refs` lookups; `git status` shows it |
| L7 | `--work` of the refresh script is unchecked (it removes only `<work>/build`) | operator | the operator's own scratch directory |
| L8 | Hostnames and paths of the operator's machines in `deploy/remote.md` and seven older trail files; a claude.ai artifact URL in the v2.0.0 handoff | operator | none while the repository is private; scrub with `git grep -l` before it is ever published, and confirm the artifact is not shared by link |
| L9 | Plan decision D24: the suite takes 2.1 times and the first install 5.8 times as long with the catalog (expectations were 2 and 3) | operator, recommendation: accept | seconds, once per place |
| L10 | The catalog tests pass the whole environment to git: an exported `GIT_DIR` or `GIT_INDEX_FILE` would aim the fixture's `add`, `commit` and `reset --hard` at another repository (pre-existing; `reviews/2026-10-01-security-ff4c181-159585a.md`) | main agent | only a shell that exports those variables, which git does inside hooks, and this project has none |

## Rows of the Definition of Done

- **Dirty source before the gate (Step 1):** the catalog plan's Phase 3 and review work, brought through the rows and committed as 7dd0645. The gate then fixed its own findings in four commits — f62513b, b4e3d79, ff4c181, 159585a — each with a `changes/` entry and a verifier run at its head.
- **Verification:** `testing/2026-10-01-verify-v2.1.0-gate.md` PASS at 7dd0645; `-r2` PASS at f62513b; `-r3` PASS at b4e3d79; `-r4` FAIL at ff4c181 on finding 7 above, everything else green; `-r5` PASS at 159585a, the gated head — `gofmt` clean, `go vet` clean, six packages ok with `-race`, 364 top-level tests (3 env-gated skips), the two timer tests running in a build without cgo, `govulncheck` clean, recall 0.633, both size pins (14,252 of 15,000; 15,155 of 22,000), the listing budget 19,236 of 19,500, `TestEmbedHoldsOnlyTrackedFiles`, a scratch-home install with no drift, three release targets compiling (19.9, 19.6, 19.2 MB), and two must-fail controls. The `-r4` FAIL is carried here as it is: its one finding is fixed by a test-only commit and `-r5` names it closed.
- **Reviewer:** every commit of the range is covered by a COMPLETE report whose last round is not REQUEST CHANGES — 8b6c6de: `reviews/2026-10-01-review-worktree-timer-per-home-code.md` (round 2 APPROVE); 425a444: `reviews/2026-10-01-review-worktree-gate-lows-code.md` (round 3 APPROVE WITH NITS); 754a7ce..7dd0645: `reviews/2026-10-01-review-catalog-plan-code.md` (round 3 APPROVE WITH NITS, its four lows fixed) and `reviews/2026-10-01-review-catalog-plan-prompts.md` (APPROVE WITH NITS); 7dd0645..ff4c181: `reviews/2026-10-01-review-7dd0645-f62513b-code.md` (round 1 REQUEST CHANGES on finding 1, round 3 APPROVE WITH NITS). ff4c181..159585a is one test file, +26/-7, below the threshold: read by the main agent against the diff. Earlier reports whose first-round verdict was REQUEST CHANGES are cited by their last round.
- **Security:** the whole range, in slices, each COMPLETE with an Outcome — the 136 vendored entries at cf9c5e2 in five reports (`reviews/2026-10-01-security-catalog-{ai-agents,gke,data-storage,ops-iam-secops,solutions-router}.md`: every high closed by excluding five upstream skills; the entries have not changed since, which `TestCatalogMatchesChecksums` pins); 8b6c6de..425a444 in `reviews/2026-10-01-security-worktree-gate-lows.md`; code and harness text at 7dd0645 (`…-security-v2.0.0-7dd0645-code-text.md`); the trail at 7dd0645 (`…-security-v2.0.0-7dd0645-trail.md`, 55 of 55 paths read); the catalog's manifest, index and checksums and the delta to f62513b (`…-security-v2.1.0-catalog-harness-files.md`); f62513b..b4e3d79 (`…-security-f62513b-b4e3d79.md`); b4e3d79..ff4c181 (`…-security-b4e3d79-ff4c181.md`); ff4c181..159585a (`…-security-ff4c181-159585a.md`). No high; one medium-equivalent (finding 1's second form, filed by the reviewer) fixed. The "between slices" gap one report named — the manifest after the entry reviews — is the subject of the fourth report.
- **Docs lens (Step 4):** `reviews/2026-10-01-review-v2.0.0-7dd0645-docs.md`, APPROVE WITH NITS, eight lows, all fixed (Outcome in the report). It confirms the delivery claim is unchanged: the catalog sits under `skills/`, already embedded in the binary and carried by the zip, so neither the README map nor `specs/sh-agent/delivery.md` § The model moves. The docs changes after 7dd0645 are those fixes and the gate's own sentences, read by the main agent against the diff.
- **Documentation agent:** dispatched once (README's timer row and frontmatter sentence, `memory-loop.md`, `testing/guide.md`); its edits read against the code. The release-state paragraph of `delivery.md` is rewritten in the release commit by the main agent.
- **Bugs:** BUG-146, BUG-147, BUG-148 (the refresh script) and BUG-149 (the timer's real home), all fixed in the range; BUG-145's record restored. BUG-139, BUG-140 and BUG-142 stay open from before the range; BUG-142's rule (a frozen tree under the verifier) held in every round of this gate.
- **Changelog:** entries by the main agent for each commit (`changes/2026-10-01-1555-…`, `-1615-…`, `-1636-…`, `-1657-…`, `-1717-…`) and the release entry.

## Trail freshness (Step 3)

- The newest `changes/` entry describes the release (the release entry of Step 6; before it, the gate's fix entries) — pass.
- No `date:` in the future under `.agents/docs` — pass (one template line reads `YYYY-MM-DD`).
- `changes/` and `architecture/DECISIONS.md` append-only: `git log --diff-filter=M v2.0.0..HEAD` over both shows no commit — pass. The 1615 entry's sentence about the timer's limit was wrong; it is corrected by the 1636 entry, not edited.
- `ROADMAP.md` touched in f82c0c1 and f62513b (Now is empty since the plan closed); `TODO.md` untouched and empty — report-only, no lag.

## The project's own release rule (`AGENTS.md` § Release)

- Real-home smoke on both CLIs: `testing/2026-10-01-verify-real-home-v2.1.0.md`, 6 of 6 at 159585a (context, one skill, one agent on Claude Code and on Antigravity).
- One Antigravity run of each skill the range changed: `testing/2026-10-01-verify-skill-matrix-antigravity-v2.1.0.md`, 6 of 6 (`bash`, `cloud-run-deploy`, `data-gcp`, `review`, `feature`, `project-memory`), on the skill text that ships. No agent changed.
- The catalog, as data: the router's rows `catalog-no-leak`, `catalog-route` (seven tasks) and `catalog-precedence` on both CLIs at ff4c181, whose router, notes and index are what ships (`testing/2026-10-01-verify-google-skills-catalog-phase-4.md` § Reruns during the v2.1.0 gate: 9 of 9 on each CLI, the manifest note read first in 14 of 14 routed runs); and the security review of the import, the five slice reports above.
- Tag: signed per `deploy/remote.md` § Releases from this checkout; the signing key's fingerprint check is the first command of Step 6.

## What was not checked

- The Mac: neither the smoke nor a build run there; `printf %q` and the refresh script on macOS's bash 3.2.
- One interactive session per CLI on the operator's machine (does it prompt before reading under the installed skill directory): the operator's, moved from the plan to the handoff. It is not an Antigravity run of a changed skill: the headless rows above are those runs.
- A real `systemctl` or `launchctl` call for BUG-149 (the tests stop at the unit name, by design: R6); an account outside `/etc/passwd` (simulated through a seam).
- The refresh script against the real GitHub remotes after its hardening (the fixtures are local repositories).
- Antigravity's no-leak from a record: its init event lists no skills, so that cell is two identical self-listings.
- The sentences that changed in the six harness skills, as behaviour: the rows show each skill loads and its procedure runs.
- `go test -cover` percentages; the Antigravity rows of skills the range did not change.

None of these is a condition the blocking list names, the project's rule included: every skill the range changed has its Antigravity run, and the real-home smoke exists on both CLIs.

VERDICT: SHIP
