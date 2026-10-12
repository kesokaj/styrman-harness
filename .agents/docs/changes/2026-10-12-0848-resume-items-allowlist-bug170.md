---
owner: main
date: 2026-10-12
---
# Resume items after the abandoned v2.6.2..66f3412 gate: BUG-170, the -u class, the agent allow-list, rows 21, 23, 28 L2     (2026-10-12)

## Summary
Handoff 2026-10-09-1921 said nothing was queued, but a release gate had started after it and stopped. Its reviewer artifact was a stub, its security review complete, and its Claude and Antigravity smoke runs failed on expired logins, not on the harness. The operator approved every resume item ("Go all, fit all into next version preflight") and chose item 7's scope: the agent-key allow-list, rows 21 and 23, and row 28 L2 in; the dream login stop deferred. All in `ee20d00`:
- **BUG-170:** dream's `callWithin` watches SIGINT and SIGHUP only when they were not ignored at start, so `nohup sh-agent dream --all &` no longer loses its project at logout and runs on.
- **The `-u`/`--user` needle:** the gate's finding 2 put a quote, a backtick and a key's colon in front of the flag. This change's security review found more shapes, and the operator chose the class over a fourth batch: any non-word boundary, and whitespace, `=`, an argv-array separator or an escaped quote before the value.
- **ADR-056:** an agent's frontmatter is an allow-list (`agentKeys`), so `memory`, `initialPrompt`, `isolation` and a vendor's next key are refused at install.
- **Row 21:** the dream timer refuses every control character in the binary path, the home and PATH.
- **Row 23:** search's `index` and `--rebuild` refusals, which the CLI could not reach, are removed and the CLI refusals pinned.
- **ADR-057 (row 28 L2):** the list of setups needing the `@AGENTS.md` stub moves to the init-project skill, R2 edited with the operator's approval of the exact text, the context's pointer resolving at every place.
- **Smaller fixes:** the `readSmall` cap has a test; the BUG-167 record no longer says its fix put `mcpServers` in the paragraph; the dream spec and a derive comment are current.

## Type
task (resume items; one bug record, two ADRs, one requirements edit)

## Files changed
- `internal/dream/run.go`, `run_test.go` — BUG-170: `signal.Ignored` before `Notify`; `TestCallLeavesTheGroupOnAnIgnoredSignal`
- `internal/dream/timer.go`, `timer_test.go` — row 21: `unicode.IsControl` refusal; `TestInstallTimerRefusesEveryControlCharacter`
- `internal/harness/report.go`, `report_test.go` — the command-line credential rule as a class; shapes, near-misses, a pinned over-redaction, `TestReportRedactFlagAtValueStart`
- `internal/harness/inventory.go`, `frontmatter.go`, `harness_test.go`, `derive.go` — ADR-056: `agentKeys`, `unlistedKey`; lint tests per kind
- `internal/search/search.go`, `cmd/sh-agent/main_test.go` — row 23: two unreachable refusals removed; the CLI refusals pinned
- `internal/metrics/events_test.go` — `TestReadSmallStopsAtOneMiB`
- `REQUIREMENTS.md` (R2), `context/AGENTS.md`, `skills/init-project/SKILL.md`, `ROADMAP.md` — ADR-057
- `AGENTS.md`, `README.md`, `.agents/docs/specs/sh-agent/dream.md` — the allow-list and the signal exception documented
- `.agents/docs/architecture/DECISIONS.md` — ADR-056, ADR-057
- `.agents/docs/bugs/BUG-170-*.md`, `INDEX.md`, `BUG-167-*.md` — the new record, its row, the corrected line
- `.agents/docs/reviews/` — the abandoned gate's security review (2026-10-09), this change's three records; the gate's reviewer artifact, a 23-line IN PROGRESS stub with no findings, deleted untracked

## Agents dispatched
reviewer (code+docs, once), security (once), verifier (once, last at `ee20d00`). Artifacts: `reviews/2026-10-12-review-worktree-code-docs.md`, `reviews/2026-10-12-security-worktree.md`, `reviews/2026-10-12-review-worktree-merged.md` (record and Outcome), `testing/2026-10-12-verify-ee20d00.md`.

## Verification
`testing/2026-10-12-verify-ee20d00.md`: BLOCKED at `ee20d00` on tier 9 only. Tier 9 is the sandbox start of each CLI: no committed kit has a context-question row, and the logins are the operator's. Tiers 1–8 are green:
- **Tests:** 688 pass, 0 fail, 5 skip; govulncheck reports no new advisory.
- **Install and report:** a scratch install shows drift none with a clean `binary` line, and `report --dry-run` sends nothing.
- **Search:** `index --rebuild` 6.0 s, a warm `search` 0.21 s.
- **Timer:** the units rendered through a fake manager pass `systemd-analyze verify`.

The first verifier dispatch hung for over 30 minutes in a Bash call that never started (most likely an approval prompt a background agent cannot show) and was stopped. The second finished the tiers from its outputs. Tier 9 moves to the release gate's real-home smoke on both CLIs (guide § 5), together with the Antigravity run of the changed `init-project` skill.
Every code fix had its test red first: the two signal subtests, the five control characters, the four unknown agent keys, and the JSON-lines and plain shapes of both needle changes (6 failures on the JSON-lines test alone before the class). The `readSmall` cap was checked by mutation: raising the cap failed the test, and the restore was diffed clean. Over-redaction of the class rule on the trail: 15 lines at v2.6.2, 36 now, all in review records.

## Notes
- After the verifier, the timer refusal names its alternative (the verifier's § 3c finding): `internal/dream/timer.go`, asserted in `TestInstallTimerRefusesUnsafePaths`, red first. It is committed after `ee20d00`, so the release gate's verifier covers it.
- The release gate is next, by the operator's word: `/preflight` for the next version. The release commit updates `delivery.md` § Where things stand, including the refused-key clause (the merged record's finding 3).
- A `-u` shape the class misses is a known limit from now on, not a patch (the rule's comment, the operator's choice).
