---
owner: reviewer
date: 2026-09-15
---
# Review — 641eccd-worktree-code   (2026-09-15)

Range reviewed: 641eccd..fe21671 (`git diff 641eccd..fe21671 -- context agents skills README.md AGENTS.md .agents/docs/deploy .agents/docs/testing/skill-matrix.md`) plus the uncommitted working tree (`git diff fe21671 -- skills .agents/docs/testing/skill-matrix.md`: `skills/init-project/SKILL.md`, `.agents/docs/testing/skill-matrix.md`)

## Verdict: REQUEST CHANGES

The worst problem is in `/preflight` Step 4, the mechanism ADR-023 just settled. A rebutted docs finding is carried forward only from the newest prior gate report's `## Outcome`, and no step writes the carried rebuttal into the next report's `## Outcome`. So one gate report later the finding is raised again and blocks. That is the failure ADR-023 was written to remove ("instead of blocking every later release"). A second defect sits in the same mechanism: the docs reviewer is never told about the rebuttal, so its REQUEST CHANGES verdict can still block. A third: Step 6 still says SHIP "has no findings left", which contradicts Step 4's new `left` gaps. These three findings are the third round on this one mechanism (BUG-026, ADR-023). Under Definition of Done step 2 and BUG-026's prevention rule, take them to the operator before a third patch. Two findings sit outside preflight. README's new `--vendor all` "widens" claim is false for `sh-agent update`. And the binary's two CLAUDE.md messages still give the procedure README just replaced. The working-tree init-project update is sound. Its shell snippets were run in scratch clones and behaved as documented, and the new `init-update` matrix check passes a correct update and fails two wrong ones.

## Findings

- **[medium / confidence: high]** `skills/preflight/SKILL.md:65` — A rebuttal is carried forward for one gate report only.
  - Step 4 reads "the `## Outcome` of the newest prior `reviews/*-preflight-*` report" and cites a rebutted finding "in the gate report". It does not say the citation goes into that report's `## Outcome`.
  - Step 6 (line 81) fills `## Outcome` only with "one line per finding" of the gate it resolves. A cited finding is by definition not raised, so it gets no line. On SHIP the report "needs none".
  - Failure path: gate N-1 has the rebuttal in its Outcome. Gate N reads it and cites the finding in its body. Gate N+1 reads gate N's Outcome, finds no rebuttal, and raises and blocks the finding again. That happens whenever two NO-SHIPs follow each other, or a SHIP is followed by the next release.
  - Same day, same version: Step 5 writes the same path, `$(date +%F)-preflight-<version>.md`, so the re-run overwrites the report whose Outcome held the rebuttal. Git still has it; the next run does not read git.
  - The "same statement at the same file:line" key also breaks when a fix for another finding shifts lines in the same doc. That key is the operator's ADR-023 wording, so this is noted for the operator, not re-argued.
  - Fix (operator's call, per BUG-026): Step 4 reads the Outcome of every prior preflight report, or each gate report copies carried rebuttals into its own `## Outcome`. The shorter option: a rebutted docs finding gets a doc change so it stops being raised, and the carry-forward mechanism is deleted.
- **[medium / confidence: medium]** `skills/preflight/SKILL.md:65,77` — A rebutted finding can still block through the reviewer's verdict.
  - Step 4 dispatches the docs reviewer "with the docs brief exactly as `/review` Step 3 gives it". That brief (`skills/review/SKILL.md:48`) carries no prior rebuttals, so the reviewer raises the rebutted finding again and sets its verdict from it.
  - Step 5 blocks on "a Step 4 REQUEST CHANGES verdict or high-severity docs finding, not counting a finding a prior gate report rebutted". The exclusion names findings, not the verdict.
  - Result: a rebutted finding that is the only reason for REQUEST CHANGES still blocks every run.
  - Fix: the verdict blocks only when a non-rebutted finding supports it. Or drop the verdict clause and block on high-severity non-rebutted findings alone.
- **[medium / confidence: high]** `skills/preflight/SKILL.md:81` — Step 6 contradicts the new `left` gaps.
  - Stale text: "the SHIP path has no findings left and needs none" (no `## Outcome`).
  - Line 65, changed in this range, now says "every other gap is recorded as `left` with an owner", and ADR-023 says "A low-severity docs gap can ship, recorded with its owner". A SHIP report therefore does have `left` findings.
  - No step says where they go: Step 5's section list has no place for them, and Step 6 says SHIP needs none.
  - An agent following Step 6 on SHIP drops the owners ADR-023 requires. This is the describing-text rule (review-standards §3).
  - Fix: Step 5 names where `left` gaps go on SHIP (for example a `## Outcome` with `left` lines), and line 81 drops "has no findings left".
- **[medium / confidence: high]** `README.md:32` (also `.agents/docs/deploy/local.md:27`, which is scoped to install) — `--vendor all` does not widen through `update`.
  - README says every command accepts `--vendor` and "`--vendor all` widens a single-vendor install to both". That holds for `install`: `cmd/sh-agent/main.go:137` passes an explicit flag through the rebuild.
  - It fails for `update`: `cmd/sh-agent/main.go:348` appends `--vendor` only when `vendor != harness.VendorAll`. The child `install --repo …` then resolves the default from the manifest (`resolveVendor`, main.go:216-224) and reinstalls the old single vendor.
  - Silent failure: `sh-agent update --vendor all` on a claude-only install never installs Antigravity.
  - Checked by reading the code. Not run, because `update` fetches and pulls this repo.
  - Fix: pass `--vendor` whenever `vendorFlag != ""`, as install does, with a test. Or narrow the README sentence to `install`.
- **[medium / confidence: high]** `internal/harness/apply.go:37`, `cmd/sh-agent/main.go:309` — The binary contradicts README's new CLAUDE.md procedure.
  - README.md:53 now tells the operator to move a hand-written `~/.claude/CLAUDE.md` into `rules/`, remove it, and reinstall, so the file becomes the one-line import R2 requires.
  - The binary still prints the replaced procedure: `install` says "add the line `@AGENTS.md` to it yourself", and `status` says "make sure it contains the line `@AGENTS.md`".
  - Following the binary leaves a CLAUDE.md with content plus an import, which is not the one-line import of REQUIREMENTS.md R2. The operator also gets two conflicting instructions.
  - This is a describing-text miss (review-standards §3; BUG-027's "grep the subject, not the old wording"). The binary is outside the diff's paths but is what README describes.
  - Fix: change both messages to point at the README procedure. No test asserts on either string (grep of `cmd/`, `internal/`).

## Criteria
1. Dead code — pass. The change is docs and skill text only, with no orphaned file or dangling reference. The `anti-ui-slop` MANIFEST mention was correctly removed: the directory holds only LICENSE and NOTICE.
2. Flow tracing — fail. `--vendor all` traced through `update` → `Reinstall` → child `install` loses the flag (finding 4). Preflight Step 4 → Step 5 → Step 6 → the next Step 4 loses the rebuttal (findings 1, 2).
   - What held: README `--dry-run`, `--force` and repo discovery via parents match `main.go`/`repo.go`. The remote.md mirror paragraph matches both mirrors' `VENDORED.md`, crew's `image/Dockerfile` builder stage, argolis's `032_setup-harness.sh` (run as root) and its 5-minute `devbox-autoupdate.timer`. The Mac log cited in remote.md exists.
   - init-project Steps 3–6, run in `.agents/scratch/review-641eccd/upd` (a clone of the older demo-app scaffold): update added only `.agents/docs/retros/.gitkeep`, and a second run added nothing.
   - The `init-update` matrix check exits 0 after a correct update and 1 for a modified `TODO.md` and for a scaffold-message commit (BUG-024's must-pass rule met).
   - In a fresh repo (`.agents/scratch/review-641eccd-new`) the loops produce exactly the documented tree and four resolving links. A wrong-target link and a real directory are reported as "kept".
3. Standards conformance — fail on describing text (findings 3, 5). The init-project staging and commit wording matches the global hard rules (stage by name, read `--cached --stat`, commit in its own call). `go run ./cmd/sh-agent list` and `install --dry-run --home .agents/scratch/review-641eccd-home` exit 0: the new quoted frontmatter passes the installer's lint, and the dry run wrote no file.
4. Directory layout — pass. No new files. The new matrix row sits in the existing table.
5. Simplicity (DIY-except-security) — concern, not a separate finding. The preflight rebuttal carry-forward is the third round of patches on one mechanism (BUG-026). Deleting it (fix the doc so the finding stops recurring) is shorter than defending it. No security-sensitive code.
6. Decoupling — pass. Skill edits stay inside their own files. Preflight reuses the `/review` docs brief by reference, which is also why finding 2 needs a change on the preflight side, not in the brief.

## Not examined
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` (modified in the working tree; outside the dispatched scope).
- The `init-update` matrix row was not run through either CLI. Only its check was exercised, against states built by hand.
- The init-project shell snippets were run on Linux only; macOS `readlink`/`find`/`sort` behaviour was not run.
- The commit message's verification claims (`go test -race`, `gofmt`, `go vet`, render) were not re-run: no Go file is in scope.
- `sh-agent update` was not executed (it fetches and pulls this repo), so finding 4 is from reading the code.
- Of argolis `032_setup-harness.sh`, only the header was read; crew's `entrypoint.sh` was not read.
- `agents/documentation.md` was not traced for the Step 4 fix routing.
- Skill-matrix rows other than the new `init-update` row.
