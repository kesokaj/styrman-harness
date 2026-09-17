---
owner: main
date: 2026-09-17
---
# Release gate — v1.1.0     (2026-09-17)

Range reviewed: v1.0.0..f129aa7

## VERDICT: SHIP

No blocking finding. Nothing in this release is high severity, no reviewer returned REQUEST CHANGES, the verification row passes at the commit being released, no `REQUIREMENTS.md` violation stands, and no project agent shadows a harness review agent (`.agents/agents/` and `.claude/agents/` hold none — checked before any dispatch, BUG-033).

Thirteen findings were raised across five dispatches. All thirteen are FIXED — none was carried into the release as accepted risk — and each fix has a test that fails without it, or, for the three in documentation, a correction read against the code.

## What this release is

The harness can now be installed from one downloaded binary with `context/`, `skills/`, `agents/`, `commands/` and `rules/` built into it, and `sh-agent update` at such a place moves it forward from a signed GitHub release. `sh-agent report` no longer keeps a clone of the private feedback repository on any machine. `sh-agent report` itself is also new since v1.0.0. MINOR, not MAJOR: every v1.0.0 install keeps working, `state.json`'s shape is unchanged, and no flag was removed — an existing place gains behaviour and is required to react to nothing.

## Rows

| Row | Verdict | Artifact |
|---|---|---|
| Verification, release range | PASS — 99 top-level + 85 subtests, 0 fail, 0 skip; `govulncheck` clean; embedded install byte-identical to the checkout install | `testing/2026-09-17-verify-preflight-v1.1.0.md` (v1.0.0..086deab) |
| Verification, after the gate's own fixes | PASS — 101 + 85, 0 fail, 0 skip, reproducing the main agent's count exactly | `testing/2026-09-17-verify-preflight-v1.1.0-fixes.md` (v1.0.0..6ef7434) |
| Code review, release range | APPROVE WITH NITS — 0 high, 0 medium, 5 low | `reviews/2026-09-17-review-preflight-v1.1.0-code.md` |
| Code review, round 2 over the fixes | APPROVE WITH NITS — 0 high, 0 medium, 3 low | `reviews/2026-09-17-review-086deab-6ef7434-code.md` |
| Security, release range | FINDINGS — 2 low, both in the release procedure, none in the code | `reviews/2026-09-17-security-preflight-v1.1.0.md` |
| Docs lens (advisory) | APPROVE — 0 findings | `reviews/2026-09-17-review-preflight-v1.1.0-docs.md` |
| Release rehearsal | every check passed | § Rehearsal below |

Earlier phase coverage cited rather than repeated, each with its own `## Outcome`: `reviews/2026-09-17-review-f7ee817-21ba8dc-code.md`, `-security-f7ee817-21ba8dc.md`, `-security-73f2c36-83f3cf6.md`, `-security-dab2565-4d45155.md`, `-review-edbddca-1a5578b-docs.md`, `-review-6af54f7-015cf2d-code.md`, `-security-6af54f7-015cf2d.md`, `-security-015cf2d-012b4ca.md`. Two dispatches verified independently, rather than taking it from the dispatch, that `24d08a8^{tree}` and `cd0bbcb^{tree}` are the same tree (`da8d57d5`) — so those phase reviews describe exactly what ships — and that no `*.go` file changed between the merge and the gate.

## Findings, worst first

All fixed. Severity and confidence as the reporting dispatch gave them.

1. **[medium/high] no downgrade protection, and the version was not inside the signed bytes** (Phase 4 security). `tag == version` was an equality test and `SHA256SUMS` named only files and hashes, so a genuinely signed OLDER release installed and verified honestly, and one release's signed checksums verified under another release's tag. Fixed: the tag must parse as `vX.Y.Z` before it reaches argv, must be numerically newer, and `SHA256SUMS` must carry `# styrman-harness vX.Y.Z` matching the tag (D26, D27). This is why the release procedure below writes that line — a release cut the old way is correctly uninstallable.
2. **[medium/high] `RemoveKeptClone`'s guard failed open** (Phase 3 security). When it could not count unpushed work — a branch with no upstream — it deleted the clone holding the operator's only copy. Fixed to count over all refs and keep the clone whenever the count cannot be taken; a stash, ignored files and a symlink now keep it too (BUG-069).
3. **[medium/high] `update --check` deleted files** (Phase 3 security). The read-only mode removed the kept clone because the cleanup sat above its returns. Fixed and pinned in both directions (BUG-070).
4. **[medium/high] `go:embed` shipped untracked files** (Phase 2 security). An ignored `rules/.env` was compiled into a build while every other test stayed green, because nothing read the artifact in the other direction. Fixed by pinning the embedded set to `git ls-files`, and by the release guard below (BUG-066, D23).
5. **[low] `report --dry-run` looped forever** on an unreadable `reports/<host>`, treating "cannot read" as "name taken". Fixed; the test hangs without the fix.
6. **[low] a permanent push refusal was misdiagnosed as a race.** `! [remote rejected]` — a pre-receive hook, a protected branch, a read-only credential — was retried and then reported as something to re-run, forever. Narrowed to non-fast-forward shapes, with git's real wordings measured under `LC_ALL=C` by the round-2 review.
7. **[low] the embed guard ran after the tag push** (release security), so a contaminated tree was caught only once the tag was public and unrecallable. Moved above the tag — and then, on round 2's correct objection, restored before the builds as well, since `go:embed` reads the tree at build time. It runs in both places.
8. **[low] a dead `verifyErr` store**, which was exactly what separated a broken verifier from a bad release: an OpenSSH too old for `-Y verify` exits 255 like a bad signature. Now ssh-keygen's own first line rides with the verdict. My first attempt classified by wording and broke a refusal test; that attempt is gone.
9. **[low] `allowed_signers` was left on disk** by the release procedure — the very artifact the code refuses to trust. Now removed after each check.
10. **[low] `sshKeygen` accepted a non-executable file**, and **[low] the verifier came from `PATH`** before a system path (Phase 4 security). Both fixed (D28).
11. **[low] `checksumFor` accepted a malformed hash and the first of two conflicting lines**; **[low] unbounded reads**; **[low] a guessable download directory**; **[low] a clone error could repeat a URL carrying a token**; **[low] a dead `empty` return**. All fixed.
12. **[low] three stale restatements in documentation** — README § The loop still framed R8 as "a checkout or an unpacked copy"; an `install` table row omitted the release binary; the plan described the guard's old position. Fixed. This class has now appeared five times (BUG-003, 007, 027), and the rule those recorded — grep the subject nouns, not the old wording — is what found two of the three.
13. **[low] the plan promised a partial-clone fallback the code never had.** The plan is a requirement of record, so the row was corrected to what actually happens: git warns and sends every blob, the run is correct and slower, and the warning is discarded.

## Rehearsal of the release procedure

`.agents/scratch/v110-rehearsal/` (gitignored). The procedure in `.agents/docs/deploy/remote.md` § Releases was run end to end with a TEST key, then its output was fed to the shipping code:

- Three platform binaries built with the documented flags; `sha256sum -c` accepts `SHA256SUMS` WITH the `# styrman-harness v1.1.0` comment line, so the marker does not break the standard tool.
- Signed and verified with the consumer's own command before any upload.
- A binary built as "v1.0.0" with the test signer linked in via `-X` — no seam added to shipping source — updated ITSELF from those assets: verified, replaced itself with the linux/amd64 asset, installed 260 files, and reported `source embedded v1.1.0` with drift none.
- Three tampered variants of the same release were each refused for the right reason (flipped byte → checksum; edited `SHA256SUMS` → signature; missing `.sig` → refused), and the installed binary survived all three byte-identical.
- The contamination guard was checked in the must-fail direction too: a planted ignored `rules/.env` is seen by the `git status --ignored` guard and by `TestEmbedHoldsOnlyTrackedFiles`. That leg ran after every dispatch had finished, because it writes into the repository (BUG-072).

## What was not checked

- **No real release exists yet**, so every download path was exercised against a fake `gh`. The first genuine `gh release download` happens when a place updates to v1.1.0.
- **`darwin/arm64` and `linux/arm64` were built, never run.** This is a Linux/amd64 box; the Mac legs of the plan's e2e criteria are owed on the work Mac.
- **`go install …@v1.1.0` is untested** and cannot be tested before the tag exists.
- The round-2 review argued the must-fail mutations from the code rather than running them; the main agent ran them (twelve for the release path, five for the report path, each restored from a saved copy).
- The release-range security dispatch read no test file and ran no mutation: it confirmed the prior fixes are PRESENT, not re-confirmed they are TESTED. It also left 97 paths under `## Not examined` — all trail markdown, `ROADMAP.md` and three agent definitions, no Go source and no skill among them, which is why this does not block (BUG-034 covers a release-diff path named there; these are named and are not code).
- Neither CLI was started in the sandbox for this release: no `context/`, `skills/` or `agents/` CONTENT changed in the gate's own commits, though earlier commits in the range did and were exercised then.

## Version

v1.1.0 — MINOR. New behaviour for an installed place (a second way in, `update` from a release, `report` keeping nothing), nothing an existing deployment must react to. The tag is the only version tracking; no version file exists and none was introduced.
