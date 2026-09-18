---
owner: reviewer
date: 2026-09-18
---
# Review — plan-sh-agent-works-everywhere-cold-read-3   (2026-09-18)

Range reviewed: plan file, round 3 — `.agents/docs/plans/2026-09-18-sh-agent-works-everywhere.md` (with D8) and `internal/harness/apply.go`

## (1) Gaps

- **BLOCKING (new)** Phase 1 task 1's chosen fixture for "previous file unreadable" (`apply.go:33-62`, D6) is "a directory in its place, which `ReadFile` refuses". I verified with a standalone `os.Rename` probe (regular file → existing directory) that this fails with `EEXIST`/`file exists`, not a silent overwrite. In the default branch, `writeAtomic` ends with `os.Rename(name, dest)` (`apply.go:298`); when `dest` is a directory, that rename errors, and `Install` returns `fmt.Errorf("writing %s: %w", dest, err)` (`apply.go:61`) instead of printing `update ~/<rel> (previous file unreadable: …)` and succeeding. The plan's own claim — "lists `update … (previous file unreadable: …)` and is written" (Phase 1, line 47) — cannot be true for the fixture it names; the install would instead return an error (or, if `writeAtomic` is extended to `RemoveAll` a directory first, that is new behavior the plan never decides). Either the fixture must change (e.g. a regular file that returns a non-`ErrNotExist` read error some other way — permission bits are unreliable as root/CI often bypass them) or `writeAtomic` must gain directory-removal handling the plan doesn't currently authorize. This blocks Phase 1 task 1 as written: two implementers would diverge on whether the test asserts success or an error.
- **Resolved from round 2**: D8 plus the added Phase 1 clause ("`--dry-run` over a deployed file `chmod`ed to 0644 lists the `mode` line and leaves the mode as it was") closes round 2's blocking gap about the `os.Stat` probe being trapped inside `if !dryRun`. Confirmed against `apply.go:49-55`, which today wraps both `Stat` and `Chmod`, matching D8's cited line range.
- **Carried over, still low severity**: Phase 1's checkpoint (line 49) still states "260 at the time of writing" for the add-line count with no re-derivation caveat beyond the parenthetical already there — acceptable as written, not blocking.

## (2) Scope vs Why

No new endpoints, verbs, flags, or config keys found beyond rounds 1-2. D1-D8 each carry an attributed operator quote or an implementer note with citation; nothing in the plan builds beyond the `## Why`/`## Scope` request.

## (3) Readiness

1. **BLOCKING** — Does the "previous file unreadable" test fixture (a directory standing in for the deployed file) actually exercise the `update … (unreadable)` line, or does `os.Rename` over a directory error out first? Verified: it errors first. Not stated; the plan asserts an outcome ("is written") that contradicts what the named fixture produces with today's `writeAtomic`.
2. **NON-BLOCKING** — Exact wording/formatting of the `mode` line's `(0644 -> 0755)` old→new values (octal vs `-rwxr--r--` style) — Verification line 40 makes wording non-binding; a reasonable implementer follows the interface line's literal example.
3. **NON-BLOCKING** — The "260 add lines" checkpoint figure will drift with unrelated content changes; a reasonable implementer re-derives it at run time.

## Verdict

`design needs revision` / `implementation not ready`: one new BLOCKING gap. D8 correctly closed round 2's finding, but Phase 1 task 1's own added test fixture for the D6 "unreadable file" case is internally inconsistent with `apply.go`'s rename-based `writeAtomic` — verified by direct `os.Rename` probe (file over existing directory → `EEXIST`). The plan needs either a working fixture for "read fails but write can proceed" or an explicit decision that this case is instead a refusal/error path.

## Not examined

Did not re-read `main.go`, `plan.go`, `repo.go`, `source.go`, `release.go`, `report.go`, `adopt.go`, `state.go`, `render.go`, `inventory.go`, `derive.go`, `pathguard.go`, `apply_test.go`, or any docs surface (`README.md`, `.agents/docs/deploy/{local,remote}.md`, `.agents/docs/specs/sh-agent/delivery.md`, `.agents/docs/testing/guide.md`); Phase 2-5 content beyond re-confirming D8 against `apply.go`. This round re-verified round 2's fix (D8) and probed the one new fixture Phase 1 added since round 2, given the 20-tool-call budget.
