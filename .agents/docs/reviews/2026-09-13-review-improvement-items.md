---
owner: reviewer
date: 2026-09-13
---
# Review — improvement-items application   (2026-09-13)

Range reviewed: 5cdd4d6..working tree, scoped paths only (per dispatch):
`internal/harness/{repo,apply}.go`, `cmd/sh-agent/main.go` + tests, `agents/*.md`,
`context/AGENTS.md`, `context/dev-setup.md`, `skills/{handoff,bug,feature,research,
feedback,retro,init-project,gcloud-setup}/SKILL.md`, `skills/init-project/assets/*`,
`skills/VENDORED.md`, deleted `skills/gcloud-auth-verification/`, `skills/docker/SKILL.md`,
`skills/frontend/SKILL.md`. `gofmt -l .`, `go vet ./...`, `go build ./...`, and
`go test -race ./...` all pass clean on the working tree.

## Verdict: APPROVE WITH NITS

Everything in the Go layer (item 5, 17) is correct and covered by new tests (nested-copy
refusal, GIT_DIR/GIT_WORK_TREE not fooling the check, buildinfo read of the installed
file) that pass under `-race`. Every cross-file contract I traced for items 1–4, 6–14,
16, and 19 (TODO/ROADMAP ownership, memory removal, bug-record-in-fix-commit, ADR
template matching what `changelog` actually writes, dedup, gcloud-setup merge, Colima
split) is consistent end to end. The one real problem is in `skills/frontend/SKILL.md`
(item 15): the new "SPA is the default, no Node at runtime" framing was added at the
top of the file but the `createServerFn` section below it — code example and the
"Secrets... only inside `createServerFn`" rule — was left exactly as it was when SSR
was the assumed default, so the file now tells a model two incompatible things about
where server-only code and secrets live.

## Findings

- **[medium / confidence: high]** `skills/frontend/SKILL.md:77-96` — Self-contradiction introduced by item 15's rewrite. Line 20-22 states SPA is the default and "no Node runs in production." Lines 77-79 then say `createServerFn` "need[s] the SSR build: in SPA mode there is no Node at runtime, so the Go API is the only backend and every call goes through `src/lib/api.ts`" — i.e., in the default mode `createServerFn` cannot be used for server-only logic. But the code block immediately below (lines 83-86) still shows `createServerFn({...}).handler(async () => { return db.query(...) })` as *the* pattern for server-only DB access, and the Rules bullet at line 91 still says "Secrets and privileged calls **only** inside `createServerFn` handlers — never in components..." with no SPA/SSR qualifier. A model following this skill for the default (SPA) case has no consistent instruction: put secrets in `createServerFn` (which the prose just said doesn't run), or route everything through `src/lib/api.ts` to the Go backend (which the Rules bullet doesn't mention as the SPA-mode alternative). Also line 72's loader comment "runs on server during SSR" is likewise stale for the SPA-default framing introduced above it. Fix: state plainly that in SPA mode (default) secrets/DB access live in the Go backend behind `src/lib/api.ts`, `createServerFn` is SSR-only and used when a project opts into SSR, and adjust the loader comment and code example (or clearly gate them under an "SSR only" heading) to match.

- **[low / confidence: medium]** `skills/frontend/SKILL.md:59-62` — Version pins ("`vite` 8.2", "`@tanstack/react-start` 1.168") are stated as fact ("Pinned by the reference project on 2026-09-13") without a way for me to verify them against styrman-crew's actual `package.json`; I did not have that file in scope to cross-check. Flagging as low/uncertain rather than asserting an error — worth a quick `cat web/package.json` cross-check next time this file is touched.

## Criteria

1. **Dead code** — Pass. The deletions (`gcloud-auth-verification/SKILL.md`, the `memory: project` frontmatter lines, the `.claude/agent-memory/` gitignore entries, the "Suggested Skills"/global-handoff sections) all have zero remaining references outside historical `.agents/docs/` records (out of scope, already excluded) and `VENDORED.md`'s deliberate note explaining the merge.
2. **Flow tracing** — Pass for the Go change: `ownCheckout` is checked at the top of `Inspect`, `Behind`, and `Pull` (every git-touching entry point), `Install` gets its commit through `Inspect` (so a plain copy records `commit: ""` rather than a foreign repo's HEAD), and `status`/`update` both handle the empty-commit and `errPlainCopy` cases without a crash or a misleading message. Fail (one item) for the prose: item 15's `createServerFn`/SPA contradiction in `frontend/SKILL.md` (finding above) — a model tracing "where do secrets live" through that file lands in a dead end.
3. **Standards conformance** — Pass. `gofmt`, `go vet`, `go build`, and `go test -race` all clean; the Go changes match the `golang` skill's error-wrapping and stdlib-only conventions (no new dependency; `debug/buildinfo` is stdlib).
4. **Directory layout** — Pass. No new files landed outside their expected trees; the deleted skill directory and its `SKILL.md` are fully gone; report written to `.agents/docs/reviews/`.
5. **Simplicity (DIY-except-security)** — Pass. `ownCheckout` via `os.Lstat(repo/.git)` is simpler and more robust than the review item's own suggested `git rev-parse --show-toplevel` comparison (avoids a git invocation and is immune to `GIT_DIR`/`GIT_WORK_TREE` redirection, which the new test explicitly exercises) — a case where the implementation improved on the proposal without adding complexity. No new dependencies anywhere in this diff.
6. **Decoupling** — Pass. The ownership split (main agent owns `TODO.md`/ROADMAP "Now"; `changelog` owns `changes/`+ADRs; `documentation` owns ROADMAP Next/Later) is now stated identically in `context/AGENTS.md`, `agents/changelog.md`, `agents/documentation.md`, and `skills/handoff/SKILL.md` — no agent's lane overlaps another's after this change.

## Not examined

Per the dispatch's explicit exclusion list: `skills/golang/references/*`, `skills/docker/references/*`, `README.md`, `ROADMAP.md`, `TODO.md`, `.agents/docs/*` (said to be already verified mechanically). Also not examined, outside the stated scope: `REQUIREMENTS.md` (spot-checked only for item 19's R7 wording, which is applied and consistent), `skills/cloud-run-deploy/SKILL.md` (noticed in passing while checking item 8's cross-reference — its Colima/gcloud-setup edits looked consistent, but it was not in the scope list so I did not review it line by line), and `web/package.json` for the frontend version-pin claim (finding 2 above). No security-sensitive surface (auth/crypto/secrets) was touched in this diff.
