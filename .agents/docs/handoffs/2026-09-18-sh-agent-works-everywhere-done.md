---
owner: main
date: 2026-09-18
---
# Handoff — 2026-09-18 — "sh-agent works everywhere" is done; v1.2.3 on all three machines

Session ID: 28ca24dd-a2d8-4278-85aa-70e24bbed445
Project: /home/kesokaj/workspace/styrman-harness
Branch: `main` at 6bebc57, clean and pushed. Tags: `v1.0.0` … `v1.2.3`, all pushed and published.

## Context

The operator's Definition of Done for this session: "work until the delivery mechanism sh-agent is flawless … All features of sh-agent need to work and make sure the harness is installed when you say install." The plan `.agents/docs/plans/2026-09-18-sh-agent-works-everywhere.md` carried it through five phases; every box is `[x]` or `[-]` with evidence; `ROADMAP.md` "Now" is empty. Nothing is queued.

## Completed

- **Every verb on every kind of place** — `testing/sh-agent-matrix.md`: 53 cells on this box (checkout, plain copy, release-stamped binary), none silent, every refusal naming its cause and the alternative; verified twice by a re-run at the final tree. Three real machines in § Real machines: this box (checkout place), the work Mac (darwin/arm64, Antigravity only) and devbox (linux/amd64, `--vendor all`) — each updated v1.2.2 → v1.2.3 with the new `mode ~/.config/sh-agent (0755 -> 0700)` line, `sh-agent version` unqualified reporting the new version, `type sh-agent` naming the installed path, `status` drift none with no `shadowed` line, the runtime dir `700`.
- **`install` installs the harness and says what it did** — `add`/`update`/`mode`/`remove`/`keep` per file, on every run; from a fresh home on each place kind and from the published v1.2.3 binary itself, `status` reports drift none; `agy -p` on the Mac and `claude -p` on this box quoted a context rule from the installed file.
- **Audited** — five lenses over `cmd/sh-agent` and `internal/harness` whole (`reviews/2026-09-18-review-harness-merged.md`): 19 findings, 15 fixed with tests that failed first, 4 left with reasons; three fix rounds; two stops at the operator's door where a fix had introduced the next defect (the note's trigger, D2; the runtime dir's tightening, option 1).
- **Released** — v1.2.3 (`reviews/2026-09-18-preflight-v1.2.3.md`, SHIP; `changes/2026-09-18-release-v1.2.3.md`), signed, published, consumed from a fresh download, installed from the published binary. Earlier today: v1.2.2 and the closeouts of BUG-078, BUG-079, BUG-080.

## Open / In Progress

- **Nothing in a plan.** The gate's accepted lows are the plan's Follow-ups: an unwritable `~/.config/sh-agent` (e.g. 0500) lets `install` deploy every file and then fail at `SaveState` with no manifest — a writability probe that refuses BEFORE any write, with a test, is the first Task next session; `yamlQuote` and the Unicode line-break characters (U+0085, U+2028/9, C1); the doc chains' `grep -qxF` marker unpinned to `release.go`'s regex; a swallowed `filepath.Abs(home)` error; the `settings.json` editor's move out of `main.go`.
- **Remainder no machine here can run:** the README browser chain on the Mac (`V=v1.2.3`; settles `shasum --ignore-missing`, Gatekeeper and the `xattr` step); `linux/arm64` (built and published for every release since v1.1.0, never executed); the darwin swapped-signature edge of `update`.
- The kept feedback clone on this box (`~/.config/sh-agent/feedback`, two commits `update` cannot prove pushed) — the operator's to delete by hand, as every `update` says.
- Requirements edits awaiting approval: none. `close-map-gaps`: ADR-035, still waiting on another place's proposal.

## Key Decisions

D1–D10 in the plan. Load-bearing: D2 — the BUG-079 note only after a binary was written (the operator's "1"); D4 — per-file lines always, no flag ("Go with 1"); D9 — the unreadable-file fixture at the cold-read cap ("Go with 1"); D10 — Phases 4 and 5 swapped (the release before the real-machine cells); option 1 at Phase 2's round-2 cap — tighten as a `mode` line at the top of `install`, never fail ("Go with 1"). The version: PATCH per D3, the operator's "approved go go go".

## Corrections

- none logged: `~/.config/sh-agent/feedback.jsonl` has no entry for this session.

## Agent notes (unverified)

- Three review caps were reached today, and each time the right move was to delete or simplify the mechanism, not repair it: the lexical "on PATH" beside a file identity; a note that fired when nothing was installed; an unconditional chmod after every write. A fix on a mechanism reviewers have already flagged should get its must-fail case in the ARRANGEMENT that failed before it is called fixed — a rule for `/retro`.
- The file tool converted an escape into its character twice (a test expectation, an explorer's brief), caught only by grepping for the character; the mutation harness reported a build failure as "fails as it must" twice before it printed the FAIL line. Both are `/retro` input; neither is a user correction.
- The `feedback-log` skill wrote nothing this session although the operator redirected the work at least twice ("The binary should only use the published git repo…" yesterday; "more info … Go with 1" today); why it does not fire on redirections is worth one line in the retro.
- `.agents/scratch/` is large (four releases' assets, three consumer downloads, matrix fixtures, three mutation harnesses, a dozen verify homes) — gitignored, deletable.

## Environment State

- `main` at 6bebc57, clean, pushed. Six tags — `v1.0.0`, `v1.1.0`, `v1.2.0`, `v1.2.1`, `v1.2.2`, `v1.2.3` — all published; repository `PRIVATE`.
- Suite: 250 tests/subtests, none failing or skipped; `gofmt`, `go vet`, `go test -race`, `govulncheck` clean on go1.27.1; `GOOS=windows`, `darwin/arm64`, `linux/arm64` builds clean.
- This box's real home: installed from `main` at 20ca4b4 on the operator's word, runtime dir `700`. The Mac and devbox: v1.2.3.
- `ROADMAP.md` "Now" empty; `TODO.md` empty.

## Artifacts to Read

- `.agents/docs/plans/2026-09-18-sh-agent-works-everywhere.md` — the whole plan, closed, D1–D10, Follow-ups.
- `.agents/docs/testing/sh-agent-matrix.md` — every cell, this box and three real machines, the remainder.
- `.agents/docs/reviews/2026-09-18-review-harness-merged.md` — the audit and what happened to each finding.
- `.agents/docs/reviews/2026-09-18-preflight-v1.2.3.md` — the gate and its accepted lows.
- `.agents/docs/changes/2026-09-18-release-v1.2.3.md` — what shipped, one page.

## Next Steps

1. **Operator, Mac, when convenient:** the README browser chain once with `V=v1.2.3` (the exact command is README § Install without cloning, third fence); paste the output — it closes the `shasum`/`xattr` legs.
2. **Next session:** the writability probe before any write (the gate's first accepted low) as a `bug`-skill Task (it gets the next BUG number when recorded; none exists yet); then the other Follow-ups as Tasks.
3. `/retro` — the three agent notes above; the graduation of BUG-078/079's rules into `context/AGENTS.md`.
4. `close-map-gaps` when the other place's proposal arrives (ADR-035).
