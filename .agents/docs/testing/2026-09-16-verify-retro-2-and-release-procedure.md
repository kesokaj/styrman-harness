---
owner: verifier
date: 2026-09-16
---
# Verification — retro-2-and-release-procedure   (2026-09-16)

Range reviewed: 8b95f22..worktree

## Status: PASS

Scope: setup commands, install/status/list drift check, and a full rehearsal of the § Releases — signed tags procedure in a throwaway repo (must-pass and all three must-fail cases). No Go source in the change set, so `-race` coverage is unchanged from the prior baseline; sandbox CLI starts are the main agent's separate task, not run here.

## Commands

- `gofmt -l .` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -v -count=1 ./...` → exit 0; output saved to `.agents/scratch/verify-retro-2/go-test-1.txt`; counts: pass 79, fail 0, skip 1 at 8b95f22+dirty for `go test -race -v -count=1 ./...`; vs previous: `.agents/docs/testing/2026-09-16-verify-retro-2026-09-16-applied.md`, pass 64→79, skip 1→1 (same `TestFindRepo` pre-existing skip). Rise in pass count is not a FAIL condition; no fall in pass, no rise in skip.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-retro-2/home` → exit 0; `260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed`; settings.json created with attribution off.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-retro-2/home` → exit 0; `drift none — 259 files match the repo`; `repo has 27 skills, 7 agents, 0 commands, 0 rules`.
- `go run ./cmd/sh-agent list` → exit 0; table confirms **27 skills, 7 agents** — matches the counts styrman-crew's channel smoke-checks and matches § Releases' "27 skills and 7 agents today" claim.

## Release procedure rehearsal (throwaway repo `.agents/scratch/verify-retro-2/rehearsal/`)

`git init`, one empty commit (`c8a203e`). No command touched the real repository's git state, no `git config` in the real repo, no `git tag` in the real repo.

Must-pass:
- `ssh-keygen -lf ~/.ssh/id_ed25519.pub` → `256 SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU ostling.simon@gmail.com (ED25519)` — matches the fingerprint quoted in `remote.md` § Releases.
- `git -c gpg.format=ssh -c user.signingkey="$HOME/.ssh/id_ed25519.pub" tag -s v0.0.1 -m "styrman-harness v0.0.1" HEAD` → exit 0.
- `git cat-file -t v0.0.1` → `tag`.
- `git cat-file tag v0.0.1 | grep -E '^-----BEGIN' | tail -1` → `-----BEGIN SSH SIGNATURE-----`.
- `printf '%s\n' 'ostling.simon@gmail.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBkx0BICH1xG18f0a0F/02Vw0OvBTfSOVNXkPkO62NCR' > .agents/scratch/verify-retro-2/allowed_signers` then `git -c gpg.ssh.allowedSignersFile=<that file> verify-tag v0.0.1` → exit 0, `Good "git" signature for ostling.simon@gmail.com with ED25519 key SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU`.

Must-fail (each its own tag in the same scratch repo):
- Lightweight tag `git tag v0.0.1-lightweight HEAD` → `git cat-file -t v0.0.1-lightweight` = `commit` (not `tag`), as the doc claims.
- Annotated unsigned tag `git tag -a v0.0.1-unsigned -m "unsigned annotated tag" HEAD` → `git cat-file tag v0.0.1-unsigned | grep -E '^-----BEGIN'` prints nothing (no armor line); `verify-tag v0.0.1-unsigned` → exit 1, `error: no signature found`.
- Foreign key: generated a throwaway ed25519 keypair under `.agents/scratch/verify-retro-2/foreignkey/` (`ssh-keygen -t ed25519 -N "" -f ... -C throwaway@example.com`), signed `v0.0.1-foreignkey` with it (tag command exit 0, signature itself valid but under a different key), then `verify-tag v0.0.1-foreignkey` against the scratch `allowed_signers` → exit 1, `No principal matched.` — exact wording the doc's rehearsal note claims.

Did not read `/etc/styrman/allowed_signers` (root-only; not attempted, per instruction).

## Root cause (FAIL only)

N/A — Status is PASS.

## Coverage of changed code

No Go file changed in this change set (`git diff --stat HEAD` over the tracked files touched: `agents/verifier.md`, `context/AGENTS.md`, `skills/preflight/SKILL.md`, `skills/retro/SKILL.md`, `skills/review/SKILL.md` — all `.md`; plus `AGENTS.md`, `README.md`, `.agents/docs/deploy/remote.md`, `.agents/docs/architecture/DECISIONS.md`, and the new `.agents/docs/retros/2026-09-16-retro-2.md`). There is no Go test to map these to. Coverage for this instruction-text change is:
- `context/AGENTS.md`, `agents/verifier.md`, `skills/review/SKILL.md`, `skills/preflight/SKILL.md`, `skills/retro/SKILL.md` — exercised indirectly by `sh-agent install`/`status`/`list` (steps 2 above: the files install cleanly, drift-free, and the skill/agent counts parse and match).
- `.agents/docs/deploy/remote.md` § Releases — exercised directly and exhaustively by the rehearsal above (every command in the procedure block run verbatim, all must-pass and must-fail cases confirmed).
- Sandbox CLI starts (context/skill/agent parse-and-load smoke test for each CLI) — **not run here**; that is the main agent's separate task per the dispatch.
- No suppressions (`t.Skip`, `//nolint`, `eslint-disable`, `@ts-ignore`, `.only(`) were added in this change set — diff is documentation/instruction text only, no code.

## Not examined

- Sandbox CLI parse-and-load smoke checks for Claude Code and Antigravity (explicitly deferred to the main agent per the dispatch).
- `/etc/styrman/allowed_signers` itself (root-only, not read per instruction).
- No live push, no real tag, no real `git config` change — intentionally out of scope per the dispatch (rehearsal only, throwaway repo).

## Environment

Linux 6.12.107+deb13-amd64. `go` and `git` from PATH (versions not separately queried — not required by the commands run). `ssh-keygen`/OpenSSH available and signing key present at `~/.ssh/id_ed25519` matching the fingerprint in the doc. No Docker/Colima involved (no Dockerfile changed).

## Round 2

Range reviewed: 8b95f22..worktree

### Status: PASS

Scope: re-rehearsal of the REWRITTEN § Releases — signed tags block in `.agents/docs/deploy/remote.md` (stricter `test` assertions, target/ancestor binding, private-key-derived fingerprint, `sudo -n cmp` line) plus a fresh `skills/preflight/SKILL.md` Step 6 read, in a new throwaway repo `.agents/scratch/verify-retro-2/rehearsal2/`; Go setup commands re-run (no Go source changed since round 1); install/status/list re-checked for drift.

### Correction to Round 1

Round 1's `vs previous:` line compared against `.agents/docs/testing/2026-09-16-verify-retro-2026-09-16-applied.md` (pass 64). The correct newest earlier counts line for this exact command is in `.agents/docs/testing/2026-09-16-verify-attribution-write.md` (pass 79, same sha lineage). Round 2 compares against that log.

### Commands

- `gofmt -l .` → exit 0, no output.
- `go vet ./...` → exit 0, no output.
- `go test -race -v -count=1 ./...` → exit 0; output saved to `.agents/scratch/verify-retro-2/go-test-round2.txt`; counts: pass 79, fail 0, skip 1 at 8b95f22+dirty for `go test -race -v -count=1 ./...`; vs previous: `.agents/docs/testing/2026-09-16-verify-attribution-write.md`, pass 79→79, skip 1→1 — no change (no Go source in this range).
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-retro-2/home2` → exit 0; `260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed`; settings.json created with attribution off.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-retro-2/home2` → exit 0; `drift none — 259 files match the repo`; `repo has 27 skills, 7 agents, 0 commands, 0 rules`.
- `go run ./cmd/sh-agent list` → exit 0; `### Skills (27)`, `### Subagents (7)` — matches the doc's claim, unchanged since round 1 despite the Step 6 wording change (Step 6 is instructions, not a skill/agent count driver).

### Release procedure rehearsal, rewritten block (throwaway repo `.agents/scratch/verify-retro-2/rehearsal2/`)

`git init -b main`, one empty commit `b00283b070c9b40951483b7ad43949aa125adcca`. No command touched the real repository's git state; no `git config`, `git tag`, or push in the real repo.

Substitutions: `vX.Y.Z` → `v0.0.2`, `<release commit>` → `b00283b070c9b40951483b7ad43949aa125adcca` (the SHA `git rev-parse HEAD` printed).

Must-pass, every command up to and excluding the two `git push` lines, each its own call, exit read before the next:
- `ssh-keygen -y -f ~/.ssh/id_ed25519 | ssh-keygen -lf -` → `256 SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU ostling.simon@gmail.com (ED25519)` — matches the doc's fingerprint, derived from the private key.
- `git status --porcelain` → exit 0, no output (clean).
- `git rev-parse HEAD` → `b00283b070c9b40951483b7ad43949aa125adcca`.
- `git -c gpg.format=ssh -c user.signingkey="$HOME/.ssh/id_ed25519.pub" tag -s v0.0.2 -m "styrman-harness v0.0.2" b00283b070c9b40951483b7ad43949aa125adcca` → exit 0.
- `test "$(git rev-parse v0.0.2^{commit})" = "b00283b070c9b40951483b7ad43949aa125adcca"` → exit 0.
- `git merge-base --is-ancestor v0.0.2^{commit} main` → exit 0.
- `test "$(git cat-file -t v0.0.2)" = tag` → exit 0.
- `test "$(git cat-file tag v0.0.2 | grep -E '^-----BEGIN' | tail -1)" = '-----BEGIN SSH SIGNATURE-----'` → exit 0.
- `printf '%s\n' 'ostling.simon@gmail.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBkx0BICH1xG18f0a0F/02Vw0OvBTfSOVNXkPkO62NCR' > allowed_signers` (rehearsal dir, not `.agents/scratch/allowed_signers` since this whole rehearsal is already under scratch) → exit 0.
- `sudo -n cmp /etc/styrman/allowed_signers allowed_signers` → **not run: root step, operator's** (per instruction — no sudo used).
- `git -c gpg.ssh.allowedSignersFile=allowed_signers verify-tag v0.0.2` → exit 0, `Good "git" signature for ostling.simon@gmail.com with ED25519 key SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU`.
- `git push origin main` / `git push origin v0.0.2` → **not run** (no remote configured in the throwaway repo; excluded by the dispatch's must-pass scope, which stops before the two push lines).

Must-fail (each in the same rehearsal repo):
- (a) Second commit `2ad0be9fada36ba2e6c8fd85f17797e1e042cd86`, `git tag v0.0.3 <second sha>`, then `test "$(git rev-parse v0.0.3^{commit})" = "b00283b070c9b40951483b7ad43949aa125adcca"` → exit 1, as required.
- (b) Lightweight `git tag v0.0.4 main`: `test "$(git cat-file -t v0.0.4)" = tag` → exit 1 (`git cat-file -t v0.0.4` prints `commit`).
- (c) Annotated unsigned `git -c gpg.format=ssh tag -a v0.0.5 -m "unsigned annotated" main`: the armor `test` → exit 1 (no `-----BEGIN` line, grep prints nothing); `git verify-tag v0.0.5` (with the scratch `allowedSignersFile`) → exit 1, `error: no signature found`.
- (d) Fresh throwaway key `foreignkey2/id_ed25519` (freshly generated under the rehearsal dir, `ssh-keygen -t ed25519 -N "" -C throwaway2@example.com`), tag `v0.0.6` signed with it (tag command itself exits 0 — signature is valid, just under an unlisted key) → `git verify-tag v0.0.6` against the scratch `allowed_signers` → exit 1, stderr `No principal matched.` — exact wording the doc claims.

All must-pass and must-fail outcomes for the rewritten block match the doc's claims exactly, including the new `test` forms, the target-SHA and ancestor-of-`main` binding, and the private-key-derived fingerprint check. The `sudo -n cmp` line is syntactically present in the block and unexecuted here by design (root step, belongs to the operator on the crew host, not reproducible in this throwaway repo which has no `/etc/styrman/allowed_signers`).

### `skills/preflight/SKILL.md` Step 6 re-read

Re-read in full (Step 6 "Ship or stop"). It defers to a project's `AGENTS.md`/deploy guide for tag form, signing command and pre-push checks ("tag exactly as it says and run those checks before the push") and requires confirming the tag's target before pushing (`git rev-parse <tag>^{commit}` prints its SHA) — consistent with, and satisfied by, the rewritten `remote.md` block rehearsed above. No command in Step 6 itself was rehearsed standalone; it is instructions that point at the deploy-guide procedure already exercised end to end.

### Coverage of changed code

No Go file changed in this range. Documentation/instruction-text coverage as in round 1, plus the newly rewritten `remote.md` § Releases block — exercised directly and exhaustively above, both must-pass and all four must-fail cases. `skills/preflight/SKILL.md` Step 6 — exercised indirectly via `sh-agent install`/`status`/`list` (parses and installs drift-free) and by inspection against the rehearsed procedure it defers to.

### Not examined

- Sandbox CLI parse-and-load smoke checks (main agent's separate task, as in round 1).
- `/etc/styrman/allowed_signers` and the `sudo -n cmp` line (root-only, operator's machine; not run here).
- No live push, no real tag, no real `git config` change in the actual repository.

### Environment

Same as round 1: Linux 6.12.107+deb13-amd64, `go`/`git`/`ssh-keygen` from PATH, signing key at `~/.ssh/id_ed25519` matching the doc's fingerprint. No Docker/Colima involved.
