---
owner: main
date: 2026-09-28
---
# `sh-agent update --from-zip`: a third way in, the zip place     (date: 2026-09-28)

## Summary
A machine that can hold neither the operator's SSH key nor `gh` can now take the harness from a repository zip downloaded by hand. `sh-agent update --from-zip <file>`:
- reads the zip once, hashes it, and checks every entry before writing anything;
- unpacks only the five harness trees under `~/.config/sh-agent/zip-*`;
- installs them with the running binary, and deletes the unpacked copy.

Nothing is fetched, and the output flags the zip as `UNVERIFIED`. The manifest records `zip:sha256:<hex>`, and `status` and `diff` say that a zip install is not compared.

REQUIREMENTS R8/R9 carry the operator-approved text ("approve as written"), with the zip place nested under Update (D9: "Go with your recommendations"). ADR-038 supersedes ADR-034 in part. `AGENTS.md` rule 1, `delivery.md` § The model, the README map (drawn dashed, because a zip is fetched by hand) and the deploy guides all say three ways.

**Also in this change:**
- `update --dry-run` is now read-only on both paths; it used to act.
- `diff` has a zip branch.
- Leftover `zip-*` directories are swept after a successful run.
- `secretShaped` and `.gitignore` cover `*oauth-token*`.
- `testing/skill-matrix.md` takes matrix option 1: the lists the scenarios need, with the statement that the allowlist is not a security boundary.
- `skills/VENDORED.md` records that the `find-skills` Step 6 edit is not verbatim.

## Type
feature

## Files changed
- **New:**
  - `internal/harness/zip.go`: `UnpackZip`, `checkZipTrees`, `MkdirAllTracked`, `RemoveCreated`, the limits;
  - `internal/harness/zip_test.go`;
  - `cmd/sh-agent/zip_test.go`.
- **Changed code:**
  - `cmd/sh-agent/main.go`: the flag and its refusals, `updateFromZip`, the `status` and `diff` zip lines, `knownCheckout`, `update --dry-run`, the sweep, usage and the package doc;
  - `internal/harness/{source,apply,repo,plan}.go`: `Source.Origin`, the manifest record, the `ResolveSource` skip and its message, `*oauth-token*`;
  - `.gitignore`.
- **Requirements, decisions and project rules:** `REQUIREMENTS.md` R8/R9; `.agents/docs/architecture/DECISIONS.md` (ADR-038, and ADR-034 marked superseded in part); `AGENTS.md` rule 1.
- **Docs:** `README.md`; `.agents/docs/specs/sh-agent/delivery.md`; `.agents/docs/deploy/{local,remote}.md`; `.agents/docs/testing/skill-matrix.md`; `skills/VENDORED.md`.
- **Trail:**
  - the plan and recon;
  - cold reads r1–r3;
  - code reviews r1–r3, security r1–r3 and the docs lens;
  - verify logs for phases 1–3 and the E2E;
  - BUG-093 and BUG-094.

## Agents dispatched
- **Cold reads:** `explorer` (recon); `reviewer` ×3 cold reads (the last at the cap, with its residue settled by the operator as D9).
- **Code:** `reviewer` code r1 (REQUEST CHANGES), r2 (REQUEST CHANGES), r3 (APPROVE WITH NITS, nits fixed).
- **Docs:** `reviewer` docs lens (APPROVE WITH NITS, fixed).
- **Security:** `security` r1–r3, never a high. The medium, an unsigned zip, stays open by decision.
- **Build and records:** `documentation`; `verifier` for phases 1, 2 and 3 (PASS); `bugs` ×2.

## Verification
- **Suite:** `go test -race`: 160 top-level and 291 total pass; govulncheck clean; context 20,550 of 22,000 B; scratch install drift none, with GANNET on both vendors.
- **Mutations:** each fix was proven with a mutation that compiled and failed its named test: the dry run, `status`, the size limit, the duplicate name, `diff`, the sweep and its prefix, the empty value, the other commands, the checkout dry run, the temp-directory cleanup, and the created-chain removal.
- **E2E**, run four times against the real GitHub zipball with a release-versioned build (`testing/2026-09-28-verify-update-from-zip-e2e.md`):
  - the installed context is byte-identical to the zip's render;
  - a slip zip is refused with 0 entries written;
  - a dry run leaves the home unchanged.
- **Not exercised:** macOS, and a browser-downloaded zip (the same layout, covered by a unit test).

## What went wrong
- **BUG-093:** a refused zip left `~/.config/sh-agent` behind. It took three layers to fix: the directory was first created before the name check, then before the kind check, then before `Plan`.
- **BUG-094, open:** a red-phase test ran a real `report` and pushed commit `c5d47b2` to the feedback repository. It holds metrics only, from an empty scratch home. Whether to remove it is the operator's decision.

## Docs
As listed above. The README gained an "Install from a zip" section.

## Requirements
R8 and R9 changed, with the approval noted inline; ADR-038.
