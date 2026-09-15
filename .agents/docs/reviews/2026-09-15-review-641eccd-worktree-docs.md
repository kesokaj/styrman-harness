---
owner: reviewer
date: 2026-09-15
---
# Review — 641eccd-worktree-docs   (2026-09-15)

Range reviewed: 641eccd..fe21671 (paths: context agents skills README.md AGENTS.md .agents/docs/deploy .agents/docs/testing/skill-matrix.md) plus the uncommitted working tree vs fe21671 (skills, .agents/docs/testing/skill-matrix.md). Lens: --docs. Head searched: the working tree, with `git grep` and no revision.

## Verdict: REQUEST CHANGES

The worst problem is inside the release-gate rule this range introduces (ADR-023). `skills/preflight/SKILL.md` Step 4 now lets a SHIP report carry `left` docs gaps and cited rebuttals. Step 6 still says a SHIP report has no findings and no `## Outcome`, and Step 4 reads rebuttals only from the newest prior report's `## Outcome`. So one SHIP between two runs can lose a rebuttal, and the finding it covered blocks again — the failure ADR-023 exists to remove.

Two of the flag descriptions fixed in the README do not match the code:
- The installer still prints the `~/.claude/CLAUDE.md` advice that README:53 replaced, and that advice contradicts R2.
- "`--vendor all` widens a single-vendor install" holds for `install` but not for `update`.

The rest of the scope checks out against the code, the sibling mirror repos and the sandbox: the uncommitted `/init-project` update path and its `init-update` matrix row, the mirror descriptions, the other flag descriptions, and the VENDORED and Mac-verification edits.

## Findings

- **[medium / confidence: high]** `skills/preflight/SKILL.md:81` vs `skills/preflight/SKILL.md:65` and `:77`
  - What changed: Step 4 (:65) records every non-blocking docs gap "as `left` with an owner". A finding the newest prior gate report's `## Outcome` rebutted is cited and not raised again, and :77 does not count it as blocking.
  - What was not updated: Step 6 (:81) still says `## Outcome` is appended only by the task that resolves a NO-SHIP, and that "the SHIP path has no findings left and needs none".
  - The contradiction: under the new rule a SHIP report can hold `left` gaps and rebuttals, and :81 gives them no Outcome section.
  - Why it fails: the rebuttal lookup reads only the newest prior report's `## Outcome`. After a SHIP that cited a rebuttal, the next run finds no Outcome there. It raises the finding again, and if the finding is high severity it blocks.
  - Fix, either one:
    - Make :81 match :65: a SHIP report with `left` or rebutted docs gaps writes its own `## Outcome`.
    - Have :65 read rebuttals from every prior `*-preflight-*` report.
- **[medium / confidence: high]** `internal/harness/apply.go:37` vs `README.md:53`
  - README now says to move a foreign `~/.claude/CLAUDE.md`'s content into `rules/`, delete the file and re-run install.
  - The installer's own `keep` line still prints the old advice this scope removed from README: "not written by sh-agent; add the line `@AGENTS.md` to it yourself".
  - An operator who follows the installer's output keeps a multi-line `~/.claude/CLAUDE.md`. R2 says that file is a one-line import.
  - Fix: word the message as README:53 does.
- **[medium / confidence: high]** `README.md:32` and `.agents/docs/deploy/local.md:27`
  - Both say "`--vendor all` widens a single-vendor install to both", and the README says every command accepts `--vendor`.
  - True for `install`: an explicit flag is passed to the rebuilt child (`cmd/sh-agent/main.go:137-138`).
  - False for `update`:
    - `cmd/sh-agent/main.go:348` drops `--vendor` when it resolves to `all`.
    - The child `install` (`internal/harness/repo.go:156`) then takes the vendor from the manifest (`resolveVendor`, main.go:220-223).
    - Result: `sh-agent update --vendor all` on a claude-only install reinstalls claude only.
  - Fix, either one:
    - Limit the sentence to `install`.
    - Have `update` pass `--vendor` whenever the flag was given, as `install` does.
- **[low / confidence: high]** `.agents/docs/deploy/local.md:12` vs `.agents/docs/deploy/local.md:27` and `README.md:32`
  - Old statement still present: Prerequisites says "the installer writes for both by default or selectively with `--vendor`".
  - This scope changed the default to the vendor recorded by the existing install, else `all` (`resolveVendor`). A plain `install` after a single-vendor install writes one vendor.
  - Fix: same wording as :27.
- **[low / confidence: high]** `cmd/sh-agent/main.go:122-126` (comment) vs `AGENTS.md:35` and `.agents/docs/deploy/remote.md:33`
  - The comment still says a plain copy's "image built the running binary from that same tree".
  - The scope corrected this wording to "image or VM build step": argolis-devbox builds on the VM via `deploy/scripts.d/032_setup-harness.sh`, not in an image.
  - The comment describes behaviour, so the global hard rule on describing text applies. It is outside the brief's living-docs list, so it is recorded as low.

### Unverified
- `skills/init-project/SKILL.md:46` says "Antigravity never reads the content that file holds" (a project `CLAUDE.md`). No verification log in `.agents/docs/testing/` exercises Antigravity against a project `CLAUDE.md`, so the claim could not be checked.
- `.agents/docs/deploy/remote.md:33` says argolis-devbox has "no version to bump". The repo root has no VERSION file, and its `deploy/harness/VENDORED.md` names no bump, only commit and push. A version stamp elsewhere in that repo was not searched for.

### Checked and matching (no finding)
- **Flags and repo lookup:**
  - `install --dry-run`, `uninstall --dry-run` and `adopt --force` exist (main.go:79-81, :100, :206).
  - The checkout is looked up in the working directory and its parents (repo.go:32).
- **Frontmatter lint:** it rejects a `#` at the start of a value (frontmatter.go:113).
- **VENDORED edit:** `anti-ui-slop` ships LICENSE and NOTICE and no MANIFEST.
- **Mac verification:** the log `2026-09-04-verify-mac-real-home.md` reads "PASS for Antigravity".
- **Mirror claims in `remote.md`**, checked against the two projects' own files:
  - `styrman-crew/image/harness/VENDORED.md` and `image/Dockerfile`: builder stage, release re-seed.
  - `argolis-devbox/deploy/harness/VENDORED.md`: runs as root via 032_setup-harness.sh, within five minutes.
  - Both recipes: `COMMIT`/`SHA256SUMS` plus a diff expecting three lines.
- **`init-update` matrix row base:** both sandbox demo repos have `changes/` entries and no `.agents/docs/retros/`, so the row runs the update path. The asset `agents-tree/docs/retros/.gitkeep` exists, so the row's check can pass.
- **Old-wording searches with no stale hit** in the living docs (vendored skills excluded):
  - init-project: "already-initialized", "Resume Protocol" as a stop, "four lines", `cp -R`, `ln -sfn`, "if absent".
  - preflight: "blocks the release", "never edits the tree", "a Step 4 docs finding".
  - Mirrors: "bump the project's version", "VMs rebuild".
  - Other: "MANIFEST", "`TODO.md`, high priority", "the working directory, or the path", "`all` [default]", "add the line `@AGENTS.md` to it yourself" (the only hit is in code, finding 2 above).
- **Removed "the gate never edits the tree" (old preflight :65):** it still has a home by meaning at `skills/preflight/SKILL.md:90`, "The tagged commit is the gated tree plus its records — nothing else rides along". Not dropped.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` (modified in the tree): outside the dispatched paths.
- A code-lens pass on the new shell loops in `skills/init-project/SKILL.md` Steps 4 and 5 (the docs lens only compared statements with behaviour). The `init-update` matrix row and the installer were not run.
- README statements this scope did not touch, such as the `rules/` row's `~/.claude/rules/` claim.
- The newest `changes/` entry (`2026-09-15-harness-practices-phase7-decisions.md`) covers fe21671 but not the uncommitted `/init-project` update path or the `init-update` row. That is expected before their commit, so it is not recorded as a finding.
- Vendored skills (grilling, codebase-design, find-skills, anti-ui-slop, gke-basics, terraform), excluded from the searches per the brief.
