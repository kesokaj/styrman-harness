---
owner: main
date: 2026-09-16
---
# Verification — linkshelf simulation: both CLIs build one project under the harness   (2026-09-16)

Plan: `plans/2026-09-16-linkshelf-simulation.md`. Harness installed in the sandbox home at cad463d (`binary cad463ddde69`, drift none) before the runs; the fixes this log produced were installed afterwards for the re-runs in § Re-runs.

## Status: PASS — both CLIs built the service and closed the trail without a human reply; 19 gaps found (G1–G19), 7 fixed under the bug skill (BUG-039 to BUG-043), 6 filed as Follow-ups, 6 informational

## Setup

Sandbox home `/home/kesokaj/sandbox/home`, Claude Code 2.1.273 (model `claude-opus-5`), Antigravity CLI 1.2.4. One empty repository per CLI (`git init -b main`, no commits) with a bare local `origin` (plan D1). Kit `.agents/scratch/sim-kit/` (not committed, ADR-010): `setup.sh`, `run.sh <cli> <init|feature>`, `drive.sh` (both lanes detached, status file per lane), `compare.sh` (checklist), `audit.py` (transcripts). Claude Code: `claude -p … --permission-mode acceptEdits --add-dir <home>/.claude/skills --add-dir <home>/.config/sh-agent --allowedTools <the matrix list plus sed, awk, jq, sleep, kill, timeout, docker, rm, mv, sqlite3 …> Agent`, `timeout 7200`. Antigravity: `agy -p "The project root is the current directory, <dir>; ignore any other checkout. …" --mode accept-edits --print-timeout 115m`. Brief, verbatim in both prompts: *a local-only Go HTTP service, linkshelf, that stores bookmarks — POST, GET and DELETE /links, JSON, an in-memory store saved to a file — with tests and the deploy guide filled in.* The feature prompt waives the scope-confirmation gate by name and asks for every step of the skill and the Definition of Done through `/handoff`.

## Runs

| Step | Claude Code | Antigravity |
|---|---|---|
| `/init-project` | rc 0, 173 s, 1 commit `d3ec5b0`, pushed | rc 0, 151 s, 1 commit `c520712`, pushed |
| `/feature` | rc 0, 4508 s (75 min), 5 more commits to `3bc3088`, pushed | rc 0, 1573 s (26 min), 4 more commits to `a55d328`, pushed |

## Init checklist (`compare.sh init`, 39 checks)

| | Claude Code | Antigravity |
|---|---|---|
| Passed | 34 | 39 |
| Failed | I6 ×4 — no `.claude/*` symlinks: the headless write guard refused `mkdir .claude` twice; reported with the four `ln -s` lines as the skill says (vendor behaviour, documented in Step 5) · I10 — commit trailer `Co-Authored-By: Claude Opus 5 (1M context)` (G1) | none |
| Content | honest placeholders: `local.md` § Configure "not decided yet", § Verify "once the service serves /links"; a `## Test` section added on its own | template values kept as facts: `PORT 8080`, `curl -i localhost:8080/links` for a service that does not exist (G2); no test command in `local.md` (G3); `AGENTS.md` § Testing left as the template's comments (G5) |
| Transcript | 16 Bash, 12 Write, 10 Read; 5 refusals (`Contains simple_expansion` ×2 on the skill's `for … $d` loops, one multi-operation `cp -n` chain, `.claude` writes ×2) — the agent redid the copies with Write (G4) | 100 steps, read the `init-project` skill, ran the skill's scripts verbatim, 0 `Failed to parse` in the CLI log |

## Feature checklist (`compare.sh feature`, 101 checks: layout, gofmt/vet/test, golang-skill rules, golang-io-web reference, trail, git hygiene)

| Check group | Claude Code | Antigravity |
|---|---|---|
| Build, vet, race tests | PASS; 5 test files, 52 test functions incl. six e2e tests that build the binary and drive it on a real port; persistence across restart covered | PASS; 3 test files, 8 test functions, persistence-across-reload test, 40-goroutine race test |
| Layout | `cmd/linkshelf/{main.go,e2e_test.go}`, `internal/{link,store,httpapi}`; 726 non-test lines, 1849 test lines | `cmd/linkshelf/main.go`, `internal/{model,server,store}`; 608 non-test lines, 456 test lines |
| golang skill | all PASS: stdlib only (`crypto/rand.Text()` for ids), slog JSON, no panic, no init, no globals, graceful shutdown, four timeouts, `%w` wrapping, mutex, atomic rename; `tmp.Close()`/`os.Remove` unchecked on error paths (same idiom as Antigravity, not flagged by the checklist's `_ =` needle) | FAIL F6 external dependency `github.com/google/uuid` (G6); FAIL F8b `_ = tmpFile.Close()`/`_ = os.Remove` on cleanup paths; slog JSON, no panic, no init, no globals, graceful shutdown, timeouts, `%w` wrapping, mutex, atomic rename: PASS |
| golang-io-web reference | all PASS: Go 1.22 method patterns, error-aware handler wrapper, `MaxBytesReader`, `DisallowUnknownFields` plus a one-value check, `ErrServerClosed`, `r.Context()`, request-log and recovery middleware with `Unwrap`; the reference files were opened (`golang-io-web.md` read at the start) | `MaxBytesReader` PASS (after the security round), `DisallowUnknownFields` FAIL, `WriteTimeout`/`IdleTimeout` PASS, `ErrServerClosed` PASS, `r.Context()` PASS; routing by `switch r.Method` rather than Go 1.22 method patterns; the reference files were never opened (G11) |
| API shape | exactly the brief: `POST /links` (201 + Location), `GET /links` (`{"links": []}`), `DELETE /links/{id}` (204); two env vars `LINKSHELF_ADDR`, `LINKSHELF_STORE`; 415 without `Content-Type: application/json` (its ADR-001, a CSRF finding) | `POST /links`, `GET /links`, `GET /links/{id}`, `DELETE /links/{id}`, `DELETE /links?id=`, `DELETE /links` with a JSON body — three DELETE shapes and an unrequested GET-by-id (G7); flags `-host -port -data` plus `HOST PORT DATA_PATH` with precedence rules (G7) |
| Plan and cold read | plan 472 lines, 18 boxes ticked, D1–D31, no Open rows; FOUR cold-read rounds (rounds 1–3 `design needs revision`; the run stated it "ran a 4th round rather than declaring the 3-round cap because rounds 3–4 had no BLOCKING question") — G15 | plan 227 lines, 15 boxes ticked, D1–D15, `Open: none`; two cold-read rounds (`REQUEST CHANGES` → `APPROVE`, both tokens ready); D2/D13/D14 outlived the security fixes (G8) |
| Verifier | six dispatches, four logs (phase 1–3, final), counts 18 → 58 → 63 → 80 with `vs previous`; final run `-count=1`; NO `testing/linkshelf/guide.md` (G17); its own BUG-002: `ok (cached)` hid five failing e2e tests until `-count=1` (G16) | three logs (one per phase), `counts: pass 14, fail 0, skip 0`, `vs previous` chain, `testing/linkshelf/guide.md` written; `Not examined: nothing` (G9) |
| Reviews | code review `APPROVE` with two observations rebutted verbatim in `## Outcome`; docs review `APPROVE`; security `FINDINGS` (6: CSRF medium fixed, url cap medium fixed, trailing JSON low fixed, log echo low judged a false positive after a mutation test, two low observations); `## Not examined` names govulncheck (denied by the kit allow-list), an unread `pid.txt`, selectively read tests | code review `APPROVE WITH NITS` with `## Outcome` (5 findings fixed, none rebutted); security `FINDINGS` (5, all fixed) — finding 1 High called the crypto/rand UUID "hand-rolled cryptography" while the code review had cleared it (G6); `Not examined: Nothing.` (G9) |
| Docs | README, `specs/linkshelf/api.md` 184 lines, `deploy/local.md` 163 lines with the two variables, the `listening` line and a smoke round trip, AGENTS.md § Setup updated in place (`-count=1` on the test command, architecture paragraph) | README +79 lines, `specs/linkshelf/api.md` 212 lines, `deploy/local.md` filled with the real flags and smoke check, AGENTS.md § Setup rewritten into six tiers (G13) |
| Records | 4 `changes/` entries naming agents; `bugs/BUG-001..004` (CSRF waiver, test cache, tautological test, non-idempotent cleanup) with an actionable rule each; ADR-001; handoff whose Next Steps hold the two operator items (govulncheck, a /retro on "a passing check is not evidence unless shown able to fail"); TODO and ROADMAP empty | 3 `changes/` entries naming the agents; `bugs/BUG-001` with a prevention rule that contradicts the golang skill (G6); handoff with Next Steps "no outstanding work"; TODO and ROADMAP empty |
| Git | 6 commits, stage by name, commit and push separate, clean tree, HEAD pushed; `Co-Authored-By: Claude Opus 5 (1M context)` on every commit (G1); store file and temp file gitignored | 5 commits, stage by name, commit and push as separate calls, no trailers, clean tree, HEAD pushed; data default `data/links.json` not ignored (G12) |
| Independent code read (`reviewer`, this repository's `reviews/2026-09-16-review-linkshelf-claude-code.md` and `…-agy-code.md`) | APPROVE — one low observation (memory-first `Add` can answer 500 with the link stored; documented D9, judged sound for a single-user service); every test can fail, restart durability covered at unit and e2e level | APPROVE WITH NITS — medium: the unjustified `github.com/google/uuid` dependency (the BUG-041 false positive); low: no `DisallowUnknownFields`; otherwise "genuinely solid" |
| Not passed | I6 ×4 (symlink guard), I10/F22c (trailer), I1d (`-count=1` added to the test command — a check too strict, not a gap), F12f (the request-log line is multi-line, a needle miss), F15c (no guide, G17), F21b (a needle miss: Now holds only the template comment) — 90 PASS / 11 FAIL, 3 of them the checklist's own | I1d (test command replaced by tiers), F6, F8b, F12b, F12f (no per-request log line), F23b — 94 PASS / 7 FAIL |

## Transcript audit

| | Claude Code | Antigravity |
|---|---|---|
| Skills loaded | golang (with `golang-io-web.md` and `golang-testing.md` read), json-api-style, handoff; the feature skill expanded from the slash command (no `Skill` call, expected in print mode) | feature, golang, json-api-style (+ conventions reference), handoff, review-standards (+ go-review) — no golang reference file |
| Agents dispatched | 19: reviewer ×4 (cold read), verifier ×6 (two per phase where coverage work followed, one final), changelog ×4, documentation, reviewer (code), security, reviewer (docs lens), bugs | 11 subagent conversations: reviewer ×2 (cold read), verifier ×3, changelog ×3, reviewer + security, documentation, bugs |
| Refusals | 27: `Contains simple_expansion` on `$var` loops, multi-operation chains, `govulncheck` (not on the allow-list); each reported, none worked around silently (G18) | none (`command(*)` allow-list) |
| Cold-read rounds | 4 (G15) | 2 |

## Gaps

See `plans/2026-09-16-linkshelf-simulation.md` § Scope Follow-ups for G4, G5, G10–G13 (a Claude-only Follow-up: none). Fixed under the bug skill: G1 (BUG-039, 7b8743f), G2+G3 (BUG-040, ff83cf3), G6+G9 (BUG-041, 10c2155), G7+G8+G15 (BUG-042, 140650c), G16+G17 (BUG-043, this session). Informational: G14, G18, G19 (positive controls), G10–G13 also carry quality notes. The working list with evidence per gap is `.agents/scratch/sim-kit/gaps.md` (not committed).

## G6 reproduction row `security-rand-id`

`g6-control.sh <label>`: the harness `security` agent reviews `c520712..c62f4e5` of the Antigravity repository (a `newUUID()` built from `crypto/rand` bytes); the needle is one report line that names the identifier and calls it hand-rolled or asks for a vetted library.

| Label | Harness | Claude Code | Antigravity |
|---|---|---|---|
| a1 (before the fix) | cad463d | REPRODUCED — `[low / confidence: high] store.go:229-240 — hand-rolled UUIDv4 generator`, graded low "because it is an identifier" | REPRODUCED — Checked § Crypto: "UUID generation is hand-rolled instead of using a vetted library" |
| a2 (after the fix) | 7b8743f (BUG-039..042 installed 09:09) | NOT REPRODUCED — Checked § Crypto: "nothing hand-rolled. `newUUID` takes 16 bytes from crypto/rand …"; no finding on the identifier (other findings: bind address, body limit, Content-Type, permissions) | NOT REPRODUCED — Checked § Crypto: "Secure identifier generation using `crypto/rand`; no hand-rolled cryptographic primitives"; findings: body limit, toolchain vulnerabilities |

The needle was validated in both states before the a2 verdicts counted: the a1 reports still match it, the a2 reports match only after their clearing phrases ("no hand-rolled", "nothing hand-rolled") were excluded — two earlier needles had called a2 REPRODUCED on those very lines (BUG-029's shape, caught by reading the matching line before trusting the verdict).

## Re-runs at the fixed tree

Pass 1 (`rerun.sh`, sandbox reinstalled at 7b8743f, `drift none`, the new `note:` line printed by the install itself):

| Row | Claude Code | Antigravity |
|---|---|---|
| `security-rand-id` a2 | NOT REPRODUCED (table above) | NOT REPRODUCED (table above) |
| `/init-project` on a fresh empty repository (`-g2`), same prompt as the first run | rc 0, 155 s, 1 commit; `local.md`: 0 lines with `8080`, `## Test` with the real command, no placeholder left, the Configure row deleted with a sentence saying nothing is read yet | rc 0, 168 s, 1 commit; `local.md`: 0 lines with `8080`, `## Test` with the real command, no placeholder left, the Configure row deleted (empty table kept) |
| `attribution-off` — `{"attribution": {"commit": "", "pr": ""}}` written to the sandbox home's `settings.json`, then a headless commit in a throwaway repository | rc 0; `git log --format=%B` shows two commits with no trailer (`NO-TRAILER commits=2`); the setting stays in the sandbox home from here on | not applicable |
| Antigravity CLI log after the pass | — | `cli-20260916_091518.log`: 0 `Failed to parse` |

Pass 2 (`rerun2.sh`, sandbox reinstalled at 4268166 — the tree with BUG-043 — `binary 426816664d9a`, `drift none`):

| Row | Claude Code | Antigravity |
|---|---|---|
| `verifier-43` — a verifier dispatch on a fresh clone of the demo service with an uncommitted `extra_test.go` and no `testing/<app>/guide.md` | rc 0, 86 s; wrote `testing/demo-app/guide.md` (`owner: verifier`) and a log with `counts: pass 14, fail 0, skip 0 … for go test -race -v -count=1 ./...` | rc 0, 131 s; wrote `testing/demo-app/guide.md` and a log with `counts: pass 11, fail 0, skip 0 … for go test -race -v -count=1 ./...` |
| one no-tools start per CLI (`noop`) | rc 0: `ok` and the context's first heading | rc 0: the same |
| Antigravity CLI log after the pass | — | `cli-20260916_092049.log`: 0 `Failed to parse` |

Pass 3 (`rerun3.sh`, after the real-home smoke found BUG-043's first wording still skippable): sandbox reinstalled from the working tree with the unconditional wording; `verifier-43b` on both CLIs wrote `testing/demo-app/guide.md` with `-count=1` logs (Claude 88 s, Antigravity 162 s); this box reinstalled and Claude's check 3 repeated: guide written, at `testing/guide.md` (Follow-up G20). Details in `testing/2026-09-16-verify-real-home-3.md`.

Every owed row is done: BUG-039 (`attribution-off`), BUG-040 (`init -g2` ×2), BUG-041 (`security-rand-id` a2 ×2), BUG-043 (`verifier-43` ×2, `verifier-43b` ×2, real-home 3b); BUG-042 by reading. Not re-run: a full `/feature` at the fixed tree (see Not verified).

## Not verified here

The Mac and the real home (Phase 4 installs into this box only); a `/feature` re-run at the fixed tree (an hour per CLI; the G15 cap and the G7 scope rule are verified by reading, not by a second build); `/preflight` (no release target in linkshelf).

## Environment

Linux (Debian 13, amd64), Go 1.24.4 (both CLIs' security agents ran `govulncheck` and reported 23 stdlib vulnerabilities for this toolchain — an environment finding, not a harness one), Claude Code 2.1.273, Antigravity CLI 1.2.4.
