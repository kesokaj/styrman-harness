---
owner: changelog
date: 2026-09-03
scope: changeset:713c8dc5e016
---

# Stop gate scope unification     (2026-09-03)

## Summary

The Stop gate no longer computes its own scope — it now uses the same `Owed()` computation that `harness dispatch` satisfies. This eliminates the misalignment where coverage was decided with one function and the block message printed another, leaving the gate blocking on dirty trees even when following its instructions exactly.

## Type

bugfix

## Files changed

- `internal/gate/trailgate.go` — added `commitBoundary()` helper that asks `Owed()`; deleted `CoversHead()`, `HeadScope()` as dead code
- `internal/gate/gate_test.go` — added `TestBlockMessageNamesTheScopeThatClearsIt` regression test
- `claude/agents/reviewer.md` — removed body-`Scope:` convention, replaced with `Range reviewed:` in body and `scope:` in frontmatter
- `claude/agents/verifier.md` — removed body-`Scope:` convention, replaced with `Range reviewed:` in body and `scope:` in frontmatter
- `shared/skills/preflight/SKILL.md` — removed specification of unsupported `commit:<sha>` and `<base>..<head>` scope forms

## Verification

- `TestBlockMessageNamesTheScopeThatClearsIt` verifies the property: extract the scope from the block message, write exactly that entry, and the block clears.
- `CheckFreshness` also moved onto `Owed()` and documents why it deliberately answers a wider question than the Stop gate (reports on the working tree; gate fires only at commit boundary).

## Deferred work — Open problem

A content-digest scope is too fine-grained for the changelog row: a digest moves on every edit, so on a live tree the entry is stale as soon as it is written. The changelog records a human-meaningful unit (a task), not file content state. This means the changelog agent's own anti-duplication rule declines to write another entry when code changes after the first write, even though the row is no longer satisfied. Accountability rows (reviewer, security) do not have this problem because they genuinely should re-run when code moves.

**Action:** Add item to `TODO.md` to re-evaluate scope scheme. This affects `changeset:<digest>` only, not the `<base>..<head>` range form used in preflight reports.

## Requirements

n/a
