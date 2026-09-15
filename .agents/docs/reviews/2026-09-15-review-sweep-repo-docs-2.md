---
owner: reviewer
date: 2026-09-15
---
# Review — sweep-repo-docs-2   (2026-09-15)

Range reviewed: no diff. Named files read whole from the working tree on 2026-09-15:
- `.agents/docs/changes/2026-09-15-research-beads-elephant-goldfish.md` (styrman-harness, HEAD d1ff612 plus uncommitted changes)
- `styrman-crew/image/harness/VENDORED.md` and `image/harness/skills/VENDORED.md` (crew HEAD 52538f7, mirror COMMIT cceee54)
- `argolis-devbox/deploy/harness/VENDORED.md` and `deploy/harness/skills/VENDORED.md` (argolis HEAD 9120516, mirror COMMIT 8046e4d)

Lens: `--docs --project`. Report only; no source, config or test was touched.

## Verdict: APPROVE WITH NITS

Nothing in scope is broken. Both mirrors are verbatim copies of their `COMMIT`. The changes entry matches its commit and the research record. Three low findings, worst first:
1. The harness's summary of the mirrors (`deploy/remote.md:33`) says every refresh bumps the project's version. argolis-devbox has no version: its refresh is a commit and a push.
2. `remote.md:33` and `AGENTS.md:35` both say the binary is built "inside their own images". argolis-devbox builds it on the VM, as root.
3. `skills/VENDORED.md:36`, mirrored into both projects, lists an anti-ui-slop MANIFEST file that is not in any tree.

## Findings

- **[low / confidence: high] divergent** `.agents/docs/deploy/remote.md:33` vs `argolis-devbox/deploy/harness/VENDORED.md:39-40` (and `argolis-devbox/AGENTS.md:160`).
  - **Says:** "Each mirror's `VENDORED.md` holds the refresh recipe (rsync from a clean checkout, record `COMMIT` and checksums, bump the project's version, push — the VMs rebuild on their next update)."
  - **Actual:** argolis-devbox has no project version ("Git commit hash IS the version", `AGENTS.md:160`). Its recipe is "commit and push: the VM applies the new `COMMIT` within one auto-update (five minutes)" (`VENDORED.md:39-40`, `deploy/systemd/devbox-autoupdate.timer:6`, `032_setup-harness.sh:215,255`). Only styrman-crew cuts a release (`image/harness/VENDORED.md:47-48`), and there each project re-seeds on its next start, keyed on `/styrman/VERSION` (`image/entrypoint.sh:131-133`). That is not a VM rebuild.
  - **Should say:** "styrman-crew: cut a release; each project re-seeds on its next start. argolis-devbox: commit and push; the VM applies the new `COMMIT` within one auto-update."
- **[low / confidence: high] divergent** `.agents/docs/deploy/remote.md:33` and `AGENTS.md:35` vs `argolis-devbox/deploy/harness/VENDORED.md:11-13` and `argolis-devbox/deploy/scripts.d/032_setup-harness.sh:215-231`.
  - **Says:** remote.md, "build `sh-agent` from it inside their own images"; AGENTS.md, "a downstream mirror, whose image built that binary from the same tree".
  - **Actual:** that is true of styrman-crew (`image/Dockerfile:8-13,163-164`). argolis-devbox has no image: root builds `/opt/devbox/bin/sh-agent` on the VM from `/opt/devbox/deploy/harness` whenever `COMMIT` changes. The rule itself (a plain copy never rebuilds) still holds for both mirrors. Only the stated reason is wrong for one of the two.
  - **Should say:** "build `sh-agent` from it before installing: styrman-crew in its runtime image, argolis-devbox on the VM as root". AGENTS.md:35 should say "whose image or VM build step built that binary from the same tree".
- **[low / confidence: high] stale** `skills/VENDORED.md:36`, mirrored verbatim at `styrman-crew/image/harness/skills/VENDORED.md:36` and `argolis-devbox/deploy/harness/skills/VENDORED.md:36`.
  - **Says:** "`anti-ui-slop` includes its own LICENSE/NOTICE/MANIFEST from the upstream archive".
  - **Actual:** `skills/anti-ui-slop/` holds `LICENSE`, `NOTICE`, `references/` and `SKILL.md`. A search for any file named like manifest (`find -iname '*manifest*'`) finds nothing in the harness or in either mirror. `git log --all --diff-filter=AD` under `skills/anti-ui-slop/` shows no such file was ever added or deleted there. It may have been dropped before the skill reached this path; that part is unverified.
  - **Should say:** "includes its own LICENSE and NOTICE from the upstream archive". Fix it in the harness, then re-mirror. The mirrors must not be edited in place.

No `dropped:` or `narrowed:` findings. The changes entry is a research record that changes and removes no rule, command, flag or caveat, so the brief's old-wording grep has nothing to search for. The two mirror `VENDORED.md` files are procedures owned by their projects; nothing in them was removed relative to the harness docs that point at them.

## Classification (`--project`)

| File | Class | Proof |
|---|---|---|
| `styrman-harness/.agents/docs/changes/2026-09-15-research-beads-elephant-goldfish.md` | accurate | See "Changes entry" below. |
| `styrman-crew/image/harness/VENDORED.md` | accurate | See "styrman-crew image/harness/VENDORED.md" below. |
| `argolis-devbox/deploy/harness/VENDORED.md` | accurate, at its pin | See "argolis-devbox deploy/harness/VENDORED.md" below. |
| `styrman-crew/image/harness/skills/VENDORED.md` | outdated | Line 36 (MANIFEST, finding 3). Otherwise byte-identical to harness `skills/VENDORED.md`. The six listed skills are present; `gcloud-auth-verification`, `grill-me` and `skill-repair` are absent as it says; `find-skills/SKILL.md` still names `skill-repair`; `gcloud-setup/SKILL.md:106` has the merged heading; the licence lines in `gke-basics/SKILL.md:3` and `terraform/SKILL.md:4` are present. |
| `argolis-devbox/deploy/harness/skills/VENDORED.md` | outdated | Same file, same line 36. |
| `styrman-harness/.agents/docs/deploy/remote.md` (points at the mirrors) | outdated | Line 33 (findings 1 and 2). The other lines were not checked against the mirrors. |

Proof for the three accurate rows:

**Changes entry.**
- The counts "11 add now, 15 add next, 6 later, 8 skip" match the `### ` sections of `.agents/docs/plans/2026-09-15-research-beads-elephant-goldfish-goldfish.md` (40 `| G` rows).
- "Files changed" matches `git show --stat 3803b6c`: four text files, the entry included.
- The "removed gate" claim holds: `git show 87a862d -- context/AGENTS.md` removes the "Four gates … never skipped on the strength of an 'autonomous' instruction" line.
- "Seven open plan boxes" holds: 6 in `plans/2026-09-03-rebuild.md` and 1 in `plans/2026-09-13-learning-loop-project-level.md`, as the research record's `:115` says.
- The research file is text, not binary (the numstat shows 158 insertions).

**styrman-crew `image/harness/VENDORED.md`.**
- `diff -rq` between `git archive cceee54` and `image/harness` prints exactly the three lines the recipe expects (`:29`).
- `sha256sum -c` exits 0.
- Build and seeding match: `image/Dockerfile:8-13,163-164` and `image/entrypoint.sh:131-133,157-159`.
- `deploy/verify-commit.sh:115` runs only `sha256sum -c` and `go vet`, never the mirror's own tests (`:32-33`). `:216-219` compares the skill and agent counts.
- "BUG-015 fixed in 2c8990e" holds: `git show 2c8990e`.
- The status strings quoted at `:58` and `:63-64` exist in `cmd/sh-agent/main.go:145,234` at cceee54.
- The claims at `:34-37` hold at the pin: BUG-058 exists, and BUG-016 is fixed in cceee54.
- ADR-046, ADR-047, ADR-062 D1-D4 and the versioning spec's `:89-92` all exist and say what is cited.

**argolis-devbox `deploy/harness/VENDORED.md`, at its pin 8046e4d.**
- `diff -rq` against `git archive 8046e4d` prints the three expected lines, and `sha256sum -c` exits 0.
- The update path matches: `032_setup-harness.sh:215` builds when the binary is missing or the build stamp differs from `COMMIT`; `:255` installs when the manifest is missing or the install stamp differs; `devbox-autoupdate.service:9` runs `devbox-cli update`, and `update.sh:452` runs `scripts.d`; the timer fires every 5 minutes.
- `deploy/harness/go.mod` exists, so the root module's `./...` does not reach the mirror.
- The checks at `.agents/docs/deploy/local.md:55` match `:19-22`.
- The status strings at `:51-54` are in `main.go:145,234` at 8046e4d.
- "(D4)" at `:55` points at `.agents/docs/plans/2026-09-04-vm-harness-replacement.md:37`, not at ADR-012, which has no numbered decisions. It resolves, just not from the ADR cited at `:9`.

## Notes (not findings)

- **Stale on the next refresh:** `argolis-devbox/deploy/harness/VENDORED.md:20-22`.
  - It says `TestDryRunNeverBuilds` "presumes the tree is its own git checkout and fails in a plain copy, by construction".
  - That is true at the pinned 8046e4d, where `cmd/sh-agent/main_test.go:201` uses `../..`.
  - It is false from cceee54 onward, where BUG-016 gave the test its own checkout-shaped fixture.
  - The argolis refresh recipe does not prompt anyone to update this paragraph. styrman-crew already did (`image/harness/VENDORED.md:34-37`).
- A refresh from harness HEAD would not trip crew's symlink gate (`verify-commit.sh:117`). The only tracked symlink is `.agents/skills`, and the recipe excludes `.agents`. No gitignored file outside `.agents/` exists in the checkout, so rsync would not carry extra files.

## Unverified

- **Changes entry, agent counts:** "6 readers and 3 proposers", "37 … failed", and "49 agents, 0 errors" with its lens breakdown. No run logs are in the repo. The breakdown does sum to 49.
- **argolis, VM behaviour:** "the VM applies the new `COMMIT` within one auto-update … without a reboot" was read from the script and timer, not observed on a VM.
- **styrman-crew, inside a project:** the status output in a running project and the counts `verify-commit.sh` prints were read from code, not run.
- **`skills/VENDORED.md` pins and digests:** `6654f6b`, `435076e` and `sha256:c3dc1a03…09bc08` were not re-fetched.
- **Mirror tests:** `go test` was not run in either mirror.

## Not examined

- Harness `README.md`, `.agents/docs/testing/sh-agent/guide.md` and `.agents/docs/deploy/local.md` were searched for mirror and vendoring references only, not read whole. They are comparison surfaces, not scope.
- `.agents/docs/specs/` is empty.
- In harness `AGENTS.md`, only line 35 was checked against the mirrors.
- The uncommitted working-tree edits listed by `git status` (`skills/{bug,retro,review}`, `AGENTS.md`, `README.md`, the plan, BUG-018, `INDEX.md`, the testing guide) are not in this scope.

## Criteria

1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Re-check greps

```bash
H=/home/kesokaj/workspace/styrman-harness; C=/home/kesokaj/workspace/styrman-crew; A=/home/kesokaj/workspace/argolis-devbox
# finding 1: version bump claimed for both mirrors
git -C "$H" grep -n -E "bump the project.s version|VMs rebuild on their next update" -- AGENTS.md README.md context agents skills .agents/docs/deploy .agents/docs/testing/sh-agent
# finding 2: "image" as the only build site
git -C "$H" grep -n -E "inside their own images|whose image built that binary" -- AGENTS.md README.md context agents skills .agents/docs/deploy .agents/docs/testing/sh-agent
# finding 3: MANIFEST that is not in the tree (harness, then both mirrors)
git -C "$H" grep -n "LICENSE/NOTICE/MANIFEST" -- skills/VENDORED.md
git -C "$C" grep -n "LICENSE/NOTICE/MANIFEST" -- image/harness/skills/VENDORED.md
git -C "$A" grep -n "LICENSE/NOTICE/MANIFEST" -- deploy/harness/skills/VENDORED.md
find "$H/skills/anti-ui-slop" -iname '*manifest*'
# note: argolis test sentence, stale once the pin passes cceee54
git -C "$A" grep -n "fails in a plain copy, by construction" -- deploy/harness/VENDORED.md
# changes entry: group count (expect 40)
git -C "$H" grep -c "^| G[0-9]" -- .agents/docs/plans/2026-09-15-research-beads-elephant-goldfish-goldfish.md
```
