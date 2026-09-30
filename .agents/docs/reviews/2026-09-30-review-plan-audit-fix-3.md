---
owner: reviewer
date: 2026-09-30
---
# Review — plan audit-fix, cold read round 3   (2026-09-30)

Range reviewed: `.agents/docs/plans/2026-09-30-audit-fix.md` (cold read, round 3, no diff), checked against the code it names. Round 2's 10 questions were checked against the plan as it stands.

## Status: COMPLETE

## Verdict: REQUEST CHANGES

Round 2's four blockers are closed (REQUIREMENTS approval path D29a; shared early-write core plus per-agent trigger lines; SD7 fixture now a zeroed SQLite header, which does yield NOTADB; IN3 fixture now two links, one that the chmod would follow and one the rename would replace; C1 is a hard 15,000 pin with the arithmetic summing to about 14,650). One new BLOCKING question remains, in IN8: the lock is to be taken "around `Install` and `UpdateFromRelease`", but both re-exec a child `sh-agent install` (`internal/harness/release.go:177`, `internal/harness/repo.go:241`), so a parent that holds the lock across the child makes the child refuse with the plan's own message. The rest are NON-BLOCKING gaps: IN3 leaves `Status`/`Diff` and the resync loop reading through or removing the link; SD7's helper meets a lazy `sql.Open`; two tasks have no test or docs clause.

## (1) Gaps

- **[high / high] BLOCKING** IN8 (Phase 8, plan line 54) — the flock self-conflicts. `UpdateFromRelease` ends by exec'ing `<dest> install ...` (`release.go:177`); a checkout `Install` re-executes through a fresh build (`repo.go:225-241`, `main.go:283-286`); the checkout `update` path reinstalls the same way. flock(2) on a second open file description conflicts even with the parent, so a lock held by the parent across the child makes the child's `Install` refuse with "another sh-agent install holds ...". Decide where the lock is taken: only around the inner in-process write (`harness.Install` as the child runs it, not around the parent's download/verify/exec), or an env handoff (the child already gets `harness.ChildEnv`, `main.go:283`). Also state what the lock covers for a parent that downloads and replaces the binary before the child runs (unguarded window, or guarded by a second lock?). No test is named for contention or for the release-update path; one would catch this at once. No docs clause either (`specs/sh-agent/delivery.md` § install/update, README install row: a new refusal message on a command the map draws).
- **[medium / high]** IN3 (line 87) — the plan says "`status` counts it as kept, not drift", but `Status` (`apply.go:170-204`) reads with `os.ReadFile` (follows the link), has no `kept` state (only ok/missing/differs/foreign/stale), and `Diff` (`apply.go:207-233`) likewise follows. A link whose target has other bytes will report `differs` (drift, exit 1) forever after IN3 makes install skip it. Also the resync loop (`apply.go:97-113`): a path that an earlier install wrote (in `old.Files`) and that the operator later turned into a link is not in the new manifest (by the plan), so the loop's `Lstat` succeeds and `os.Remove` deletes the operator's link as "no longer in the source". That contradicts "sh-agent never writes through a link". Name Status/Diff/resync as in scope, pick the state name, and add a test row (status on a link; link whose path is in the old manifest).
- **[medium / medium]** SD7 (line 92) — `sql.Open` is lazy (`dream.go:62,120,206`): with a zeroed header the NOTADB error surfaces on the first `db.Query`, not on open. The plan says "when the open fails with `isUnreadable`". `openCopy` must force a read (`Ping`, or the caller's first query moves into the retry scope); say which. Also `copyFile` opens the destination `O_EXCL` (`dream.go:254`) for the db and for `-wal`/`-shm`, so "removes the copy" must remove all three copies or the second `copyDB` fails at once with EEXIST. `removeDB` (`store.go:103`) is the existing three-name helper but it is the index's own; say whether it is reused.
- **[low / high]** Struck skills' matrix rows: `testing/skill-matrix.md:56` (`python`) and `:67` (`anti-ui-slop`) stay after Phase 6 strikes the skills; the plan's strike grep (line 162) does not cover the matrix and AGENTS.md § Testing says "all skills" is claimed only after the full matrix runs. Add the deletion (or a void note: the matrix is not in the append-only set).
- **[low / high]** Early-write core (line 45, AG10): `agents/explorer.md:13` heading is "Write the brief EARLY", not "Write the artifact EARLY", and the verify extracts "between `## Write the artifact EARLY` and the trigger line". Say explorer's heading is renamed. The reviewer block the core is cut from carries reviewer-only words (the "range line"); the plan's list of what the core contains omits it, which is right, but say it is dropped so verifier/explorer/dreamer do not inherit it.
- **[low / medium]** IN4 docs clause cites README:162, which is the release-binary paragraph (no dirty tree exists there); the checkout `update` behaviour is the README:183 row and `.agents/docs/specs/sh-agent/delivery.md` § update. Also `update --dry-run` already maps to `--check` (`main.go:330-333`), so one mechanism covers both; the plan's "both" is fine but the test should assert `--dry-run` separately only if it is a separate code path (it is not).
- **[low / medium]** `--resume` lens (Phase 5): the lens list at `skills/review/SKILL.md:~22` defines `--all` as "every lens"; say whether `--resume` joins it (an explorer recon is not a review of a range).
- **[low / low]** Phase 8 S8 ("`Lint` as a `Parse` option") changes the `Lint(front string)` signature (`frontmatter.go:102`) that Phase 4 AG12b extends with the 1,024 limit; say the Phase 8 refactor keeps `TestLintDescriptionLimit` green unchanged.

## (2) Scope

Nothing built beyond the request without a D-row or an audit ID: README template is G13 now; `TestRenderedContextUnderTarget` is D2b; the flock is IN8; `explain-project` D31; `--resume` D14; `restatement-grep.sh` and `kit-preamble.sh` are TR1 (D11). The `lock_other.go` stub is covered by the stated release platforms. No finding.

## (3) Readiness — questions I would have to ask

1. BLOCKING. IN8: where exactly is the lock taken, given `UpdateFromRelease` and checkout `Install` both re-exec a child `install`? (Lock the child's in-process `Install` only, or hand the lock across with `ChildEnv`?) And what is the contention test?
2. NON-BLOCKING. IN3: does `Status`/`Diff` get a `kept` state for a link, and does resync skip a link whose path is in the old manifest? (Default: yes to both, following the stated intent.)
3. NON-BLOCKING. SD7: where does the retry force the lazy open (Ping), and does the retry remove all three copies? (Default: Ping in `openCopy`; remove all three.)
4. NON-BLOCKING. Delete the skill-matrix rows of the two struck skills? (Default: delete.)
5. NON-BLOCKING. Explorer's heading rename for the AG10 extraction. (Default: rename to the shared heading.)
6. NON-BLOCKING. IN4 docs location. (Default: README:183 row and delivery.md § update.)
7. NON-BLOCKING. Does `--all` include `--resume`? (Default: no.)

## Round 2 questions: status

1 REQUIREMENTS approval: answered (D29a, Phase 3 lines 116-117, Failure modes, Risks agree). 2 early-write block: answered (core + trigger lines; `grep -c 'after EVERY command'`), residual heading gap above. 3 SD7: answered in design (header zeroing verified to yield NOTADB via `isUnreadable`, `store.go:94`; seam `copyFile` var; injectable wait), residual lazy-open gap above. 4 IN3: answered for install; Status/resync residual above. 5 C1 hard/goal: answered (15,000 fails, 14,500 reported); arithmetic checked line by line: first list 4,600 + second 2,950 + 200 - 500 add-back leaves 14,656. 6 stale restatements: answered (D2b, D5a, D8a, D28 supersede; superseded rows marked). 7 IN4 `--check`: answered (D32, exit 1 on dirty, clean keeps exit 2). 8 behavioural runs: answered (Phase 4 covers all seven agents; context rules mapped to skill rows in Phase 3). 9 matrix rows: answered (Phase 5 and 7 write them on the demo service). 10 exec bit, preamble format, onboarding ownership, 46 rows: answered (`bash <path>`, `KIT_ASSERT_*` formats, documentation lane, 44 data rows now, 46 after two moves, checked: the Graduated table has 46 `|` lines today).

## Criteria (the plan as a design document)

1. Dead code — pass: strikes python and anti-ui-slop and dead pointers; matrix rows for the struck skills are the one leftover (gap above).
2. Flow tracing — fail: IN8 lock vs the re-exec chain; IN3 vs `Status`/`Diff`/resync; SD7 lazy open.
3. Standards conformance — fail: IN8 has no test and no docs clause; IN4 docs pointer at the wrong README row.
4. Directory layout — pass.
5. Simplicity (DIY-except-security) — pass: stdlib flock, no dependency; `api-auth-gcp` names vetted libraries.
6. Decoupling — pass: phase 3/4 edit the same agent files in sequence, stated by order.

## Not examined

The technical claims in D20-D26 (Terraform, Cloud Run, Tailwind, Vitest, Go version facts) and the five research briefs beyond the lines the plan cites; the first-list byte estimates per line were summed from line lengths of `context/AGENTS.md` (total 7,555 bytes across the named lines) but not re-measured after each proposed edit; Phase 6 and 7 skill content; the `--resume` lens parser at `skills/review/SKILL.md` was read only at the lens table; whether `sh-agent update` on a checkout reinstalls by re-exec (inferred from `main.go:372` and `repo.go:241`, not run).

design needs revision
implementation not ready
