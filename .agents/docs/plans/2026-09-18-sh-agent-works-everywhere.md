---
owner: main
date: 2026-09-18
---
# Plan — sh-agent works everywhere   (2026-09-18)

## Why

The operator, 2026-09-18, own words: "I want us to work until the delivery mechanism sh-agent is flawless. This is dod for this long session." and, on scope: "All features of sh-agent need to work and make sure the harness is installed when you say install. And all other functions of sh-agent need to work. This is the core of the harness. In other sessions we will work on improving the harness but now focus on sh-agent." Three releases this week each found a gap on a real machine that no test here had seen (BUG-076, 077, 078, 079): a stale plain copy installed as if current, a verified binary shadowed on PATH, a shell's hash running the old file. The trail proves paths one at a time; nothing yet says every path is proven.

## Scope

In: every `sh-agent` verb — `install`, `update`, `status`, `report`, `adopt`, `list`, `render`, `diff`, `uninstall`, `version` — on every kind of place — a checkout, a release binary, a plain copy — and the two orderings that bit (first install into an empty home; update over an existing install, including from a binary outside `~/.local/bin`). For each cell: either exercised with its output logged, or refused with a message that names the cause. `install` puts the harness where each CLI reads it, and `status` afterwards reports drift none — proven from a fresh home, not assumed. A full-repository code and security audit of `cmd/sh-agent` and `internal/harness`, every finding fixed or rebutted in writing. The install output names each file it adds, updates or removes (the operator's "more info", option 1, 2026-09-18). A release at the end so binary-only places get all of it.

Out: the harness content itself (`context/`, `skills/`, `agents/`) — other sessions; new verbs; Antigravity-specific behaviour beyond what `install` writes for it; the feedback repository's reading side.

Follow-ups noted, not built: `linux/arm64` execution (no machine; the plan's Phase 4 names it as the remainder); the darwin swapped-signature edge (the Mac, the operator); graduating BUG-078/079's rules into `context/AGENTS.md` (a `/retro` decision); a `status --files` listing; moving the ~220-line `settings.json` editor from `cmd/sh-agent/main.go` into `internal/harness` (audit simplify lens, medium — a structural move with no behaviour change, found in: Phase 2 audit).

## Surfaces touched

`cmd/sh-agent/main.go` (every verb's wiring, `update` both paths, `pathHint`), `internal/harness/{apply,plan,repo,source,release,report,adopt,state,render,inventory,derive,pathguard}.go` and their tests; `README.md` § Install, § Install without cloning, § Daily commands; `.agents/docs/deploy/{local,remote}.md`; `.agents/docs/specs/sh-agent/delivery.md`; `.agents/docs/testing/guide.md`. The audit's recon brief: `.agents/docs/plans/2026-09-18-recon-sh-agent-audit.md` (Phase 2's first task writes it).

## Interfaces

- `install` output, one line per file it changes, on every run (`internal/harness/apply.go` `Install`, the per-file `switch`): `add     ~/<rel>` when the previous file does not exist (`errors.Is(readErr, fs.ErrNotExist)`); `update  ~/<rel>` when it exists with different bytes; `update  ~/<rel> (previous file unreadable: <err>)` when `ReadFile` fails for any other reason — today that branch writes over the file anyway (`writeAtomic` renames over it), and this plan keeps that and only names it (D6); `mode    ~/<rel> (0644 -> 0755)` when the bytes are equal and only the permission differs — the `os.Stat` probe runs on every run, `--dry-run` included, and only the `os.Chmod` is gated on `!dryRun` (D8); counted under `unchanged` because no bytes moved (D7); `remove  ~/<rel> (no longer in the source)` when the previous install's manifest names a file the source no longer has, on the real run as on `--dry-run` (today only `--dry-run` prints it); `keep    …` unchanged. `unchanged` files print nothing. The summary line keeps its shape, `N written, N unchanged, N removed, N kept — N files managed`, with `written` = add + update. `--dry-run` prints exactly the same lines and writes nothing (`write`/`remove` today become `add`/`update`/`remove`, so one vocabulary on both runs).
- `update` after a replacement (release path) or a rebuild (checkout path): the BUG-079 note when the process is not the installed file and a PATH entry reaches it by name; silent when nothing was written (D2).
- The verb × place matrix, `.agents/docs/testing/sh-agent-matrix.md`: rows are verbs, columns are checkout / release binary / plain copy, each cell `exercised: <log>` or `refused: "<message>"` or `n/a: <why>`; every cell filled before Done.
- No new flags, no new verbs, no config keys.

## Failure modes

- A cell that cannot be exercised on any machine here (arm64) is written `remainder: <machine needed>` and named in the handoff, never left blank or marked exercised.
- A verb that fails on a place kind refuses with a message naming the cause and the alternative (the BUG-076 rule); a silent no-op is a defect.
- An audit finding the implementer disagrees with is rebutted in the report's `## Outcome` with a checkable reason, never dropped.

## Verification

Tiers: unit (`go test -race`, each verb's package tests) / integration (scratch homes under `.agents/scratch/`, every verb against each place kind built from this checkout: a clone for "checkout", the release binary built with `-ldflags -X main.version=vX.Y.Z` for "release binary", a `git archive` export for "plain copy") / e2e (the operator's Mac and devbox, and this box's real home, for the cells that need a real PATH, a real shell and a published release) / release (`/preflight`).

E2E criteria: (golden) on a machine with no harness, `sh-agent install` from a release binary produces a home where `claude -p` quotes a rule from `context/AGENTS.md` and `agy -p` lists a skill, and `sh-agent status` says drift none; (edge) `sh-agent update` run from a copy outside `~/.local/bin` names that copy, and after `hash -r` the shell runs the installed one. Both hold whatever the output lines are called.

## Phases

- [x] Cold-reader check — three rounds (`reviews/2026-09-18-review-plan-sh-agent-works-everywhere-cold-read-{1,2,3}.md`): round 1 two blocking (the `--project` scope; D6/D7 undecided), round 2 one (D8), round 3 `design ready` with one blocking item that was a test fixture, not a design gap — at the cap the operator chose to proceed with it decided as D9: "Go with 1", 2026-09-18

### Phase 1 — install says what it changed
- [x] `add`/`update`/`mode`/`remove` per-file lines on every run (`internal/harness/apply.go`), `written` = add + update, `--dry-run` uses the same words — verify: unit (`TestInstallNamesEveryChange`: a first install lists every file as `add`; a second with one changed source file lists one `update` and nothing else; a removed source file lists one `remove`; a deployed file `chmod`ed to 0644 lists one `mode` line and stays counted `unchanged`; a deployed file made unreadable (a regular file with mode `0000`, which `ReadFile` refuses for a non-root process and which `writeAtomic`'s rename replaces; the case is skipped when the test runs as root, D9) lists `update … (previous file unreadable: …)` and is written; `--dry-run` lists the same lines and writes nothing; `--dry-run` over a deployed file `chmod`ed to 0644 lists the `mode` line and leaves the mode as it was), who: AI — expect: the seven cases pass and a mutation that drops the lines fails them — docs: `README.md` § Daily commands, `.agents/docs/specs/sh-agent/delivery.md`. Evidence: `TestInstallNamesEveryChange` failed on all seven cases before the change and passes after; seven must-fail mutations (`.agents/scratch/phase1-mustfail/mutate.py`) each fail it; README `install` row and `delivery.md` say so, 2026-09-18
- [x] Existing tests that match the old `write`/`remove` dry-run lines updated in the same change — verify: `go test -race ./...`, who: AI — expect: green. Evidence: no test matched the old lines (grep), `go test -race -count=1 ./...` green on the first run after the change, 2026-09-18
- [x] Checkpoint — expect: `go run ./cmd/sh-agent install --home .agents/scratch/x` twice prints one `add` line per managed file (260 at the time of writing; the count is whatever `list` manages) then none, and a touched source file prints one `update` — evidence: `testing/2026-09-18-verify-sh-agent-works-everywhere-phase-1.md` (PASS, 227/0/0; first install 260 `add` = 260 written; second install 0 change lines; a touched skill in a scratch clone → its two vendor copies as `update` plus the rebuilt binary; `--dry-run` after `chmod 666` → one `mode … (0666 -> 0644)`, file left 0666) and `…-phase-1-2.md` over the final range d737282 (PASS, 227; a directory at a destination fails the install by name with no line). Reviews `reviews/2026-09-18-review-b393525-028a43b.md` (APPROVE WITH NITS, the low fixed) and `…-review-028a43b-d737282.md` (APPROVE); docs lens `…-b393525-028a43b-docs.md` APPROVE. Entry `changes/2026-09-18-install-names-every-change.md`. Commits 028a43b, d737282, 2026-09-18

### Phase 2 — audit
- [x] `explorer` recon: every verb × place-kind cell in the code, with file:line where each is handled or refused, and the cells no test covers — verify: brief written, who: AI — expect: `.agents/docs/plans/2026-09-18-recon-sh-agent-audit.md` with the matrix skeleton. Evidence: brief written 2026-09-18; no silent cell across ten verbs × three place kinds; four coverage gaps, below
- [x] From the recon: a CLI-level test drives `update` on a plain copy through `main.go`'s wrapping (the `errPlainCopy` refusal is unit-tested only at `Behind`/`Pull`); `status` on a plain copy pins its "update is refused" line; `adopt` on a release binary pins its refusal AND names the alternative (a checkout) as the refusal rule wants; `sync` as a literal verb — verify: unit, who: AI — expect: four tests, each failing without its code path
- [x] From the Phase 1 verifier: a relative `--home` given with `--repo` resolves against the repo directory on a real install (`Reinstall` sets `cmd.Dir = repo`) but against the invoking cwd on `--dry-run`; decide one (absolutize `--home` before anything runs) and pin it — `bug` skill, verify: unit (a relative `--home` lands in the same place on both runs), who: AI — expect: green; record BUG-080
- [x] `/review internal/harness cmd/sh-agent --all` — the paths scope, both directories read whole (`--project` is the whole repository and takes no paths: one scope per run), every lens: code, security, simplify, docs, and `--tests` mapping every source file in the two directories to its covering test — verify: the merged record `reviews/<date>-review-harness-merged.md`, who: AI — expect: every finding with severity and confidence; the verifier's map names every uncovered file
- [x] From the BUG-079 rounds (`reviews/2026-09-18-security-c868ed4-a078848.md` § Outcome): `IsReleaseVersion` rejects a non-canonical spelling (`v1.02.2`, `v01.2.2`) and `UpdateFromRelease` refuses `compareRelease(tag, version) <= 0` — verify: unit (must-fail for both spellings; must-pass for `v1.2.3` over `v1.2.2`), who: AI — expect: green, and the tag reaches no download
- [x] From the same: the BUG-079 note says "remove it" alone when `from` or its directory is not owned by the caller, and prints only for the default home (or names `--home` in the line) — verify: unit, who: AI — expect: a root-owned fixture (skip when not reproducible) gets no `ln -sfn`; `--home` scratch gets no note or a note naming the home
- [ ] Findings fixed or rebutted in writing, each with a test that fails without the fix — verify: `## Outcome` on the record, who: AI — expect: no `left` without an owner; up to three review rounds, then the operator
- [ ] Checkpoint — expect: the merged record's `## Outcome` complete; suite green — evidence: verifier log, date

### Phase 3 — the matrix, on this box
- [ ] `.agents/docs/testing/sh-agent-matrix.md` — every verb against a checkout clone, a release-stamped binary and a `git archive` plain copy, from fresh scratch homes; each cell `exercised`/`refused`/`n/a` with the command and the line that proves it — verify: integration, who: AI — expect: no blank cell; `install` then `status` reports drift none from a fresh home on each place kind — docs: `.agents/docs/testing/guide.md` (the matrix as a tier)
- [ ] Every `refused` cell's message names the cause and the alternative; a cell that is neither exercised nor refused is a bug — verify: read, who: AI — expect: zero silent no-ops, or `bug` skill runs
- [ ] Checkpoint — expect: the matrix complete for this box — evidence: the matrix file and verifier log, date

### Phase 4 — the matrix, on real machines
- [ ] Mac and devbox: `update` to the release that carries Phases 1–3, then the cells only a real home can settle (unqualified `sh-agent version`, `type sh-agent`, `status` with no `shadowed`, `report` leaving nothing, `claude -p` and `agy -p` reading the installed files) — verify: e2e, who: user (runs) + AI (logs) — expect: pasted output logged in `sh-agent-matrix.md` § Real machines
- [ ] Mac: `sh-agent install --vendor all` and one `claude -p`; the README browser chain once (`V=` the current release) — verify: e2e, who: user — expect: logged; closes the Claude-side and `shasum`/`xattr` legs
- [ ] Remainder named: `linux/arm64` (no machine), the darwin swapped-signature edge (the Mac, when the operator has a minute) — verify: read, who: AI — expect: both in the handoff's Next Steps, neither marked exercised
- [ ] Checkpoint — expect: every cell exercised, refused, n/a or remainder — evidence: the matrix file, date

### Phase 5 — release and close
- [ ] `/preflight vX.Y.Z` from the operator's own turn (PATCH if only fixes; MINOR if Phase 1's lines count as new behaviour — D3) — verify: the gate, who: user + AI — expect: SHIP, published, consumed from a fresh download — docs: release entry
- [ ] Mac and devbox `update` to it; `sh-agent version` unqualified — verify: e2e, who: user — expect: the new version, no note or a note that is true
- [ ] Checkpoint and Done — verifier, reviewer over the whole range, security, documentation, bugs, handoff — expect: `ROADMAP.md` "Now" empty — evidence: handoff path, date

## Risks

Assumptions: a `git archive` export behaves as the plain-copy places did (no `.git`; `Reinstall` copies the running binary) — checked by `go run ./cmd/sh-agent status --home <scratch>` on it saying "plain copy", first task of Phase 3; a release-stamped local build behaves as a downloaded release binary for every verb except `update` (which needs `gh` and the real release) — checked by `status` saying `source embedded vX.Y.Z` on it, first task of Phase 3; the fake-`gh` fixture in `internal/harness/release_test.go` can stand in for `update`'s network step in the matrix — checked by that test's existing PASS.

## Decisions

D1 — 2026-09-18 — Scope is every verb on every place kind, `install` proven to install, audit, release; harness content out — by: the operator, "All features of sh-agent need to work … now focus on sh-agent"
D2 — 2026-09-18 — The BUG-079 note fires only after a binary was actually written; an up-to-date or older-release run says nothing — by: the operator, "1" on the note-trigger table (the round-2 security finding: on the older-release branch the running file is the newer one)
D3 — 2026-09-18 — The release version at Phase 5 is decided at the gate from what shipped: per-file `add`/`update`/`remove` lines are output, not a new surface, so PATCH unless the audit adds behaviour — by: implementer; the operator can override at `/preflight`
D4 — 2026-09-18 — Per-file lines print on every run, no `--verbose` flag; `unchanged` never listed — by: the operator, "So, more info is the best option. It should be clear what is installed new or things like that. Go with 1"
D5 — 2026-09-18 — The audit (Phase 2) is in scope — by: the operator, "I want us to work until the delivery mechanism sh-agent is flawless" and "All other functions of sh-agent need to work"; a mechanism no independent reader has audited whole is not known to be flawless, only not known to be broken (cold read round 1, scope note)
D6 — 2026-09-18 — A previous file that exists but cannot be read is written over, as today, and the line says so: `update  ~/<rel> (previous file unreadable: <err>)`; the alternative, refusing the install, would turn a stray unreadable file into a broken place with no path forward — by: implementer, cold read round 1 (`apply.go:33`, the `default` branch takes every non-nil `readErr`)
D7 — 2026-09-18 — A mode-only correction prints `mode    ~/<rel> (<old> -> <new>)` and stays in the `unchanged` count: the contract is "names each file it changed", a permission is a change, and no bytes were written — by: implementer, cold read round 1 (`apply.go:49-55`)
D8 — 2026-09-18 — The mode probe (`os.Stat`, compare `Perm()`) runs on every run so `--dry-run` prints the same `mode` line the real run would; only the `os.Chmod` is inside `!dryRun` — by: implementer, cold read round 2 (`apply.go:49-55` wraps both today)
D9 — 2026-09-18 — The unreadable-previous-file fixture is a regular file with mode `0000`, skipped as root; a directory in its place cannot serve, because `writeAtomic`'s final `os.Rename` over a directory fails (the Phase 1 review measured `EISDIR`; the verifier's live run printed `file exists`, `EEXIST` — the kernel's choice depends on the directory; either way the rename refuses) and the install refuses rather than writes — and since the Phase 1 review, the per-file line prints AFTER the write, so that refusal leaves no line claiming a change — by: the operator at the cold-read cap, "Go with 1"; residue judged a fixture detail, not a design ambiguity
Open: none
Rejected: `--verbose` for the per-file lines — a flag nobody remembers on the day it matters; `status --files` — answers "what is installed", not "what just changed"; a Follow-up
