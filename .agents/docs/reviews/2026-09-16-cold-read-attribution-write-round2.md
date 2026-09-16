---
owner: reviewer
date: 2026-09-16
---
# Cold read round 2 — `.agents/docs/plans/2026-09-16-attribution-write.md`

Range reviewed: plan file

Round 1's five findings checked against the current plan: BLOCKING 1 (Non-goals contradiction) is resolved by D5 + Phase 1 task 4 (Non-goals exception sentence added verbatim, in the same commit as ADR-025). BLOCKING 2 (ADR-010 left standing) is resolved by the same task ("ADR-010's heading marked superseded in part") and by Surfaces naming `DECISIONS.md`. NON-BLOCKING 3 (`sync`) is resolved by D6 and by Scope naming the alias explicitly. NON-BLOCKING 4 (README line 16) is resolved by Surfaces and Phase 1 task 4 naming the reworded sentence verbatim. NON-BLOCKING 5 (malformed vs. non-object JSON) is **not** addressed — it carries forward below.

## (1) Gaps

- **BLOCKING — the plan reorders `attributionSet` before `attributionHint` inside the same `install` run, which invalidates two existing `TestAttributionNote` cases, and no task owns fixing them** — `cmd/sh-agent/main_test.go:374` and `:378`. `TestAttributionNote` drives `run(["install", ...])` then `run(["status", ...])` against the same scratch home and asserts a note is printed for *both* commands in the `"no settings file"` case (line 374, `want: body`) and the `"other settings only"` case (line 378, `want: body`). The Interfaces section has `attributionSet` run first and, in exactly these two cases, write the `attribution` key with an empty value before `attributionHint` runs its check — so under the new code, `attributionHint` finds `off == true` and prints *no* note, for `install` in the same invocation, and then `status` (second loop iteration, file already written by the prior `install`) also finds the key present and prints no note. Both existing cases go red as written. The plan's Surfaces list does name `cmd/sh-agent/main_test.go (TestAttributionNote)`, so the file is flagged as touched, but Phase 1's task list only has "`TestAttributionSet` written first and red" (task 1) and wiring `attributionSet` into the install case verified by `-run TestAttributionSet` (task 2) — neither task says these two `TestAttributionNote` cases must change, and the plan does not decide what they change *to* (delete the case, split expectations per command, or add an assertion for the new `settings:` line). A first-pass implementer hits this only at the Checkpoint's full `go test ./...`, as a surprise, with no recorded decision on the correct new expectation.
- **NON-BLOCKING, carried from round 1 (unresolved) — detection of "no attribution key" vs. "not a JSON object"** — Interfaces (plan line 21) still requires "valid JSON object, no `attribution` key" but never states the decoding target or how a syntactically valid but non-object top level (`[]`, `"x"`, `42`) is classified; the fallback branch's wording ("malformed... write nothing") doesn't explicitly cover it. A reasonable implementer folds it into the `json.Unmarshal`-into-`map[string]json.RawMessage` error path, matching D4's advisory-only spirit, but the plan still asserts rather than checks this.
- **NON-BLOCKING — `attributionSet`'s declared error contract doesn't match its sibling's actual behavior** — plan line 25 ("Never returns an error for any state of the file... the only error is one from `out`") vs. `cmd/sh-agent/main.go:488-492`, where `attributionHint` calls `fmt.Fprintf(out, ...)` without checking the returned error and unconditionally `return nil`s. If `attributionSet` is meant to actually surface `out` write failures (unlike its sibling), that's a behavioral divergence the plan doesn't flag as a decision; if it isn't, the `error` return type is decorative, same as today. Low impact — worth a one-line D-row, not a blocker.

## (2) Scope

Nothing built beyond `## Why` plus the two operator-approved additions (D2 sentence, D5 Non-goals exception) and the cold-read-sourced D6 (`sync`). No unnamed scope expansion found.

## (3) Readiness

1. **BLOCKING** — What do the `"no settings file"` and `"other settings only"` `TestAttributionNote` cases assert after this change, for both the `install` and `status` legs of the loop? The plan names the file as touched but never decides the new expected strings, and the natural fix (no note, since the key now gets written) changes correctness-relevant test coverage without a recorded decision.
2. NON-BLOCKING (carried) — malformed vs. non-object-JSON classification for the insert branch: reasonable default assumed, not stated.
3. NON-BLOCKING — does `attributionSet` need to check `out`'s write errors, unlike `attributionHint`, or does it inherit the same "never actually errors" behavior? Reasonable default (mirror the sibling) is available.
4. NON-BLOCKING (deferred by the plan itself) — new helper vs. three-state `attributionOff` — implementer's explicit choice, recorded as such.

## Verdict

**design needs revision** (round 1's two gaps are closed; one new gap found: the `TestAttributionNote` regression is a real, citable correctness gap, not yet decided).
**implementation not ready** (1 BLOCKING question).

## Not examined

`ADR-025`'s text (not yet written), `context/AGENTS.md`'s claimed 93-line render, `deploy/local.md`/`deploy/remote.md` full contents beyond the attribution string, and `internal/harness/repo.go` `Reinstall`'s full body (confirmed only that it re-execs a freshly built binary under `ChildEnv=1`, which is why the `rebuild` early-return at `main.go:141` does not itself skip `attributionSet`/`attributionHint` — the re-exec reaches that code on the child invocation, same as it does today for `attributionHint`).
