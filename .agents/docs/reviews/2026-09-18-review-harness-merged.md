---
owner: main
date: 2026-09-18
---
# Review — harness, every lens     (2026-09-18)

Range reviewed: the named files as they are in the working tree, read whole — not a diff: internal/harness cmd/sh-agent. Tree at 160aa16 (Phase 1 of `plans/2026-09-18-sh-agent-works-everywhere.md` closed). The audit phase of that plan, run as `/review internal/harness cmd/sh-agent --all`: five agents, one per lens, each given the scope text and its brief, reading the same frozen tree. The explorer's matrix (`plans/2026-09-18-recon-sh-agent-audit.md`) is background: no silent cell across ten verbs × three place kinds; four coverage gaps.

## Verdicts

| Lens | Agent | Verdict | Artifact |
|---|---|---|---|
| code | `reviewer` | APPROVE WITH NITS — 1 medium, 1 low | `reviews/2026-09-18-review-harness-code.md` |
| security | `security` (opus) | FINDINGS — 1 medium, 7 low (one low/low) | `reviews/2026-09-18-security-harness.md` |
| simplify | `reviewer` | APPROVE WITH NITS — 1 medium, 3 low | `reviews/2026-09-18-review-harness-simplify.md` |
| docs | `reviewer` | APPROVE WITH NITS — 2 medium, 1 low | `reviews/2026-09-18-review-harness-docs.md` |
| tests | `verifier` | PASS — 2 coverage gaps; 77.1% / 86.9% | `testing/2026-09-18-verify-harness.md` |

No high anywhere. No FAIL. The release verification chain — signature, namespace, in-band version marker, hash over the exact bytes written — was read end to end and found sound (security).

## Findings, worst first

Each carries the lens, severity, confidence, and the disposition the fix round takes; `## Outcome` below records what happened.

1. **[security, medium/high]** `internal/harness/adopt.go:168` `copyTree` keeps group/other write bits (`Perm()|0o600`) where `plan.go:138` strips them (`&0o755|0o600`): an adopted third-party skill lands group-writable in the tree every install ships; git records only the x bit. → **fix**: plan.go's expression; test.
2. **[code, medium/high]** `cmd/sh-agent/main.go:461`, `repo.go:118-119,153` — `update` on a plain copy surfaces `errPlainCopy` with no alternative named, unlike every sibling refusal (the BUG-076 rule). No CLI-level test. → **fix**: name the alternative (a checkout or a release binary; `install` still deploys it as it is); CLI test (recon gap 1).
3. **[docs, medium/high]** `apply.go:35,143` — the `keep` line and the `kept` count are undocumented; README `install` row and `delivery.md` list add/update/mode/remove only. → **fix** (`documentation`).
4. **[docs, medium/medium]** `main.go:546` `knownCheckout` refuses an install/update from a checkout found by walking up from cwd when a different one is on record — undocumented, and README's "works from anywhere" reads as a contradiction. → **fix** (`documentation`): describe the refusal and qualify the claim.
5. **[simplify, medium/medium]** `main.go:676-895` — a ~220-line `settings.json` editor lives in the entrypoint, not `internal/harness`. → **left** (main agent): a move with no behaviour change, during an audit whose purpose is behaviour; a Follow-up of the plan.
6. **[security, low/high]** `release.go:106,116` — `v1.02.2`, `v01.2.2`, `v1.2.02` pass `IsReleaseVersion`, compare equal, reach the download (the signed marker stops substitution); `releaseParts` reads an overflowing component as `v0.0.0`. → **fix**: canonical spelling only, `compareRelease <= 0` refuses, overflow is an error; tests for all four (Phase 2 box).
7. **[security, low/high]** `pathguard.go:112-121` — the BUG-079 note is `--home`-blind: run with `--home <scratch>` it calls the scratch home's binary "the installed binary" and may send the reader to relink their real one at it; and `ln -sfn` is offered before ownership is known. → **fix**: the note prints only for the default home; "remove it" alone when `from` is not the caller's; tests (Phase 2 box).
8. **[code, low/high]** `main.go:245-257` — `adopt` on a plain copy prints "commit and push the repo" with no `.git` present. → **fix**: the message depends on `IsCheckout`; test.
9. **[recon/tests]** `adopt` on a release binary refuses without naming the alternative; `status` on a plain copy's "update is refused" line untested; `sync` never invoked by name; `Diff`'s body 0% (the one `diff` test has no drift). → **fix**: four tests, `adopt`'s message names a checkout.
10. **[verifier, Phase 1]** a relative `--home` with `--repo` resolves against the repo dir on a real install (`Reinstall` sets `cmd.Dir = repo`) and against cwd on `--dry-run`. → **fix** (`bug` skill, BUG-080): absolutize `--home` at entry; test.
11. **[simplify, low/medium]** `main.go:691,732-733` vs `readSettings` — the 1 MiB settings bound enforced two ways; `status` never warns on an oversized file. → **fix**: one bound in `readSettings`.
12. **[simplify, low/high]** the `settings: %s could not be written (%v); add the key by hand` format string ×7. → **fix**: one helper.
13. **[security, low/high]** `report.go:400-402,453-458` — the throwaway clone's parent created `0o755` where `release.go` uses `0o700`. → **fix**: `0o700`.
14. **[security, low/high]** `apply.go:220-237` `Diff` writes the rendered harness to `os.CreateTemp("")` — `/tmp`. → **fix**: the temp file under `~/.config/sh-agent/` (the home's runtime dir), removed after.
15. **[security, low/high]** `inventory.go:394-403` `cell` slices a description by bytes and can cut a rune. → **fix**: rune-aware.
16. **[security, low/low]** `frontmatter.go` `yamlQuote` escapes `\` and `"` only; `DeriveAgy` renders descriptions with it. → **assess**: the installer already refuses multi-line descriptions; if a control character can still reach it, escape it; else rebut.
17. **[security, low/high]** `repo.go:74-82` + `main.go:546-570` — `LoadState` sanitises `Files` but not `Repo`; `knownCheckout` compares `st.Repo` with a repo `st.Repo` itself chose. → **left** (main agent): the manifest is written by `sh-agent` as the user; a prior write as the user is the attacker already holding the user's rights; `knownCheckout`'s purpose (a DIFFERENT checkout found from cwd) still holds when the walk finds one.
18. **[simplify, low/low]** `main.go:557-562` a comment duplicated back to back. → **fix**.
19. **[docs, low/unverified]** `report: --repo is not used here` on a release binary undocumented. → **fix** (`documentation`), one clause.

## Not checked

- No lens ran on a machine other than this one; the Mac and devbox cells are Phase 4's.
- The test files were scanned for secrets, not read (security); the code lens did not read the test files in full.
- `govulncheck` was not run by the agents (no network); this session's last run (`testing/2026-09-18-verify-sh-agent-works-everywhere-phase-1-2.md`) was clean on go1.27.1 with zero third-party requires.
- `report.go`'s redaction list, `release.go`'s edge messages, `adopt.go`, `frontmatter.go`, `derive.go`, `state.go`, `source.go` were grep-checked by the docs lens, not read.
- The simplify lens's medium (5) and security's 17 are left by decision, not verified further.

## Outcome

Written by the main agent after the fix round (batch A: `internal/harness`; batch B: `cmd/sh-agent`), before the round's own review; each fix has a test that failed before it and a must-fail mutation (`.agents/scratch/phase2-mustfail/mutate.py`, 13 of 13 failing as test failures).

1. adopt keeps write bits — **fixed** (internal/harness/adopt.go:175); `TestAdoptStripsGroupAndOtherWrite`.
2. `update` on a plain copy names no alternative — **fixed** in the error itself (internal/harness/repo.go:119); `TestUpdateOnPlainCopyNamesTheAlternative` drives it through the CLI.
3. `keep`/`kept` undocumented — **fixed** by `documentation` (README `install` row, `delivery.md`).
4. `knownCheckout` refusal undocumented — **fixed** by `documentation` (README § Daily commands qualified; `local.md` § Configure).
5. the settings editor in `main.go` — **left**, main agent: a structural move with no behaviour change; the plan's Follow-ups.
6. non-canonical tags — **fixed**: canonical regex and an int-fit check (internal/harness/source.go:54); `compareRelease(tag, version) <= 0` refuses with "is not newer than" (internal/harness/release.go:117); `TestIsReleaseVersionIsCanonical`, `TestNonCanonicalTagNeverReachesTheDownload` (the tag never reaches the download).
7. the BUG-079 note under `--home` and before ownership is known — **fixed**: real home only (internal/harness/pathguard.go:118); "it is not yours to relink" when `from` is not the caller's (internal/harness/pathguard.go:134); `TestNoteRunningElsewhereOnlyForTheRealHome`, `TestOwnedByCaller`.
8. `adopt` on a plain copy says commit and push — **fixed** (cmd/sh-agent/main.go:267); `TestAdoptOnPlainCopySaysWhereToCommit`.
9. the four coverage gaps — **fixed**: `TestAdoptOnReleaseBinaryNamesTheAlternative` (and the message names `git clone`), `TestStatusOnPlainCopyNamesTheRefusal`, `TestSyncIsInstall`, `TestDiffVerbShowsDrift`.
10. relative `--home` — **fixed**, BUG-080 (cmd/sh-agent/main.go:107); `TestRelativeHomeResolvesAgainstTheInvokingDirectory` reproduced the stray directory inside a scratch clone before the fix.
11. two bounds on the settings file — **fixed**: one, in `readSettings` (cmd/sh-agent/main.go:848); `TestStatusWarnsOnOversizedSettings`; `TestAttributionSet`'s expectation follows the shared note.
12. the format string ×7 — **fixed**: `settingsUnwritable` (cmd/sh-agent/main.go:831).
13. throwaway parent `0o755` — **fixed** (internal/harness/report.go:457); `TestThrowawayParentIsPrivate`.
14. `Diff` in `/tmp` — **fixed**: under `~/.config/sh-agent/` (internal/harness/apply.go:227); `TestDiffUsesTheHomeNotTmp` also covers the 0% body.
15. `cell` by bytes — **fixed** (internal/harness/inventory.go:398); `TestCellIsRuneSafe`.
16. `yamlQuote` control characters — **fixed**, escaped as `\uXXXX` (internal/harness/frontmatter.go:180); `TestYamlQuoteEscapesControlCharacters` — whose expected string the file tool converted into the character itself on first write (BUG-057's class, caught by grepping for the escape).
17. `LoadState` and `s.Repo` — **left**, main agent, as reasoned above.
18. the duplicated comment — **fixed**.
19. `report: --repo is not used here` undocumented — **fixed** by `documentation`.
