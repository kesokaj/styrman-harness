---
owner: main
date: 2026-09-17
---
# E2E — v1.1.0 published and consumed from GitHub (2026-09-17)

Range reviewed: v1.1.0 (tag cb26d2b, release commit 42698d2)

## What ran

The release was cut by `/preflight v1.1.0` from the operator's own turn, on the Linux box (`dave`), following `.agents/docs/deploy/remote.md` § Releases. This log covers the part no rehearsal could reach: the real GitHub release, downloaded and consumed by the shipping code.

1. **Publish.** `gh release create v1.1.0 --verify-tag` with the three binaries, `SHA256SUMS` and `SHA256SUMS.sig` failed TWICE with `HTTP 500: Error creating asset temp dir` from `uploads.github.com`, a server-side error; both attempts rolled back, leaving no release (confirmed with `gh api repos/.../releases`, which listed none) and the tag intact. The release was then created without assets, which succeeded, and the five assets uploaded one at a time with `gh release upload --clobber` and up to three attempts each — every one succeeded on its first attempt. Final state: five assets, `SHA256SUMS` 287 bytes, `SHA256SUMS.sig` 310, the binaries 6.7–7.0 MB.
2. **Download and verify as a consumer.** `gh release download v1.1.0 --pattern …` into a scratch directory, then the checks a place runs: `ssh-keygen -Y verify -n sh-agent-release` over the DOWNLOADED `SHA256SUMS` → `Good "sh-agent-release" signature for ostling.simon@gmail.com with ED25519 key SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU`; `sha256sum -c --ignore-missing` → `OK`; first line `# styrman-harness v1.1.0`, the marker `update` requires; `./sh-agent-linux-amd64 version` → `sh-agent v1.1.0`.
3. **Golden path, install.** Run from `/` with `SH_AGENT_REPO` unset, into a scratch home: `release binary v1.1.0: installing the harness built into it and the running binary`, 260 files written. `status` → `source     embedded v1.1.0 (release binary; update downloads the newest release)`, `installed … from release v1.1.0 (vendor: all, 260 files …)`, `binary     v1.1.0`, `drift      none — 259 files match the repo`, exit 0.
4. **Golden path, update.** `update --check` against the real GitHub → `up to date (v1.1.0)`, exit 0.

## What this settles

- **The published release is installable by the shipping code**, which the rehearsal could only show for locally built assets.
- **D21 holds in the wild.** A release binary built with `-trimpath` at a clean tag reports `v1.1.0` for its own version — `status`'s `binary v1.1.0` line proves the module-version read works for a real release, not just in the unit test.
- **D26's marker survives the round trip** through GitHub's asset storage and is read back as the code requires.
- **The repository stayed private** throughout: `gh repo view --json visibility` reports `PRIVATE` after publishing, and the operator asked about exactly this. A release in a private repository is visible only to accounts that can already read it, and asset download needs authentication — which is what README § Install without cloning relies on.

## Not exercised

- **Every work Mac leg of the plan's e2e criteria.** This ran on the Linux box only. Owed on the work Mac: the README `gh` one-liner, `sh-agent install`, `status` showing `source embedded v1.1.0`, `sh-agent report` leaving no `feedback` or `report-*` directory, and the edge case (`update` with no newer release, and a swapped signature file refused).
- **`darwin/arm64` and `linux/arm64` were built and published but never run** — no machine of either kind was available here.
- **`go install github.com/kesokaj/styrman-harness/cmd/sh-agent@v1.1.0`** — untested; the tag now exists, so it can be.
- **A real `sh-agent report` from this published install.** The operator authorised exactly one real report today (D25) and it ran earlier from the checkout; a second was not authorised, so the golden path's report leg is owed on the work Mac.
- **`update` actually replacing a binary from a real release.** `--check` reached GitHub and compared versions, but no newer release exists to update TO. The replacement path is covered against local assets by the rehearsal (`.agents/scratch/v110-rehearsal/`) and by the unit tests.
