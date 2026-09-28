---
owner: main
date: 2026-09-28
---
# Preflight — v1.8.0   (2026-09-28)

Range reviewed: v1.7.0..a962740

VERDICT: SHIP

No finding blocks. There is no verifier FAIL, no high security finding, no REQUEST CHANGES from a gate code review, no REQUIREMENTS violation (R8/R9 changed with the operator's approval quoted inline; ADR-038), and no project agent shadowing a review agent. An earlier gate attempt over `v1.7.0..381e86d` ended in REQUEST CHANGES on the matrix lists. The operator settled that as matrix option 1, which lands in `a962740` (`reviews/2026-09-28-review-v1.7.0-381e86d-code.md` § Outcome).

## Findings and decisions, worst first

| # | Source | Severity / confidence | Finding | Decision |
|---|---|---|---|---|
| 1 | security (r1–r3, gate) | medium / high | The zip is unsigned and its content becomes the agents' instructions; `UNVERIFIED` discloses the risk but does not control it | **Accepted by the operator.** The approved scope item 1 read: "A zip is unsigned, so trust rests on your own download; it gets a security review", answered "approve as written, go all but number 6" (2026-09-28). The control, a signed release zip verified by `--from-zip`, changes R9's Release procedure. It is recorded in ADR-038 § Consequences and in the plan's Follow-ups, and goes to the operator in the handoff. |
| 2 | security (gate) | medium / high (r2 and r3 rated it low) | `testing/skill-matrix.md` grants `bash`, `python3` and `curl`; the sandbox runs as the operator's user and reaches `~/.ssh` and tokens | **Accepted:** matrix option 1, approved with "go all" and "Go with your recommendations". The procedure states that the allowlist is not a security boundary. The matrix is not run in this release. Follow-up: per-row grants, deleting the Claude `.credentials.json` copy too, and dropping `curl` if no row needs it (gate code review findings 1–2). |
| 3 | code (gate) | low / high | `cmd/sh-agent/zip_test.go:514`'s test header describes the old sweep placement, and no test catches a sweep moved between `UnpackZip` and `installFiles` | Deferred to the next task (handoff). |
| 4 | security (gate) | low | `internal/harness/zip.go:86-90` discards the cleanup's `RemoveAll` error; the comment at `internal/harness/repo.go:23-28` still says a release binary installs only verified content | Deferred to the next task (handoff). A comment and a warning line; no behaviour at stake. |
| 5 | docs (gate) | low / high | `agents/bugs.md:67`'s example INDEX row is near-verbatim this repository's BUG-001, left behind by the universality cleanup | Deferred to the next task (handoff); docs gaps do not block. |

## Reviews

- **Code:** `reviews/2026-09-28-review-v1.7.0-a962740-code.md` (opus, APPROVE WITH NITS). Earlier coverage, per its `Range reviewed:`: `review-v1.7.0-381e86d-code.md` (with Outcome), `review-universal-cleanup.md`, and `review-update-from-zip-code{,-r2,-r3}.md`.
- **Security:** `reviews/2026-09-28-security-v1.7.0-a962740.md` (FINDINGS, no high; all 63 paths read or covered by a named earlier report). Earlier: `security-v1.7.0-381e86d.md` and `security-update-from-zip{,-r2,-r3}.md`.
- **Docs:** `reviews/2026-09-28-review-v1.7.0-a962740-docs.md` (APPROVE WITH NITS).

## Verification

`testing/2026-09-28-verify-v1.8.0-gate.md`: PASS over `v1.7.0..a962740`.
- gofmt and vet clean;
- `go test -race`: 162 top-level pass, 0 fail;
- govulncheck clean;
- context 20,550 of 22,000 B;
- scratch install drift none, with GANNET on both vendors;
- 27 skills and 7 agents;
- the embed guard is clean.

Feature evidence: the E2E runs on the real zipball (`testing/2026-09-28-verify-update-from-zip-e2e.md`) and the phase logs.

## Trail freshness

- **Newest change entry:** by the timed-name rule, `changes/2026-09-28-1752-update-from-zip.md`, which describes the work in `a962740`.
- **Dates:** none is in the future.
- **Append-only:** `changes/` holds added files only. `DECISIONS.md` gained ADR-038, and ADR-034's heading gained a "superseded in part by ADR-038" note, the supersession marking the changelog agent's lane allows.
- **ROADMAP.md** "Now" named the zip plan and is emptied in the release commit. **TODO.md** is empty.

## What was not checked

- macOS: `--from-zip` ran on this Linux box only; the README map draws the zip place dashed for that reason.
- The full skill matrix, and the real-home smoke after the tag (operator).
- BUG-094: the report `c5d47b2` that a red test pushed stays in `styrman-harness-feedback` pending the operator's decision.
