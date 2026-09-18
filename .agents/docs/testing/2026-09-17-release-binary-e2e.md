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

## Added 2026-09-18 — everything the Linux box can still settle

Run after the release, on `main` at c7a6116 (the tag's content plus the kept-clone message fix), all in `.agents/scratch/`:

- **`go install` from the tag works, and its version comes from build info.** `GOPATH=<scratch> GOPRIVATE='github.com/kesokaj/*' go install github.com/kesokaj/styrman-harness/cmd/sh-agent@v1.1.0` → exit 0 once git was told to fetch over SSH; the binary reports `sh-agent v1.1.0`, and `go version -m` shows `mod … v1.1.0` with NO `-ldflags` setting — so this is the `debug.ReadBuildInfo().Main.Version` branch of D6/D21, which nothing before could exercise. Installing from it, run from `/` with `SH_AGENT_REPO` unset: 260 files, `source embedded v1.1.0`, drift none. **This closes the last assumption in the plan's Risks.**
- **The README's first install path runs verbatim.** `gh release download --repo kesokaj/styrman-harness --pattern 'sh-agent-linux-amd64' --dir …` then `mv`/`chmod` → `sh-agent v1.1.0`, and it installs 260 files. The only deviation from the README was `--dir` pointing at scratch instead of `~/.local/bin`, so the operator's installed binary was left alone.
- **The README's `go install` path did NOT run verbatim** — recorded as BUG-075 and fixed in the same task. As shipped in v1.1.0 it failed with `could not read Username for 'https://github.com': terminal prompts disabled`: the comment said "over SSH" but the command gave git no way to use SSH. README now carries the working form.
- **All three published binaries verify.** `gh release download v1.1.0` with no pattern fetched all five assets; `sha256sum -c SHA256SUMS` → `OK` for `sh-agent-linux-amd64`, `sh-agent-darwin-arm64` and `sh-agent-linux-arm64`; `ssh-keygen -Y verify` over the file as published → good signature. `file` reports the darwin asset as `Mach-O 64-bit arm64` and the linux/arm64 asset as `ELF 64-bit LSB … ARM aarch64`, so the cross-compilation produced the right architectures even though neither can be executed here.
- **`sh-agent update` on the real home** (the operator asked): `up to date with upstream`, the old kept feedback clone correctly KEPT, 260 files unchanged. The clone was then measured — local `HEAD` `e973bac`, remote `main` `09586ff`, tracking ref stale at `71f749e`, and after a fetch `git log --oneline HEAD --not FETCH_HEAD` empty — so it holds nothing unique and is the operator's to delete. The misleading wording that prompted the check is BUG-074, fixed.

### Still only the Mac can settle these

- Running a `darwin/arm64` binary at all: the README one-liner, `sh-agent install`, `status` showing `source embedded v1.1.0`, and the macOS `xattr` quarantine note, which no Linux box can exercise.
- `sh-agent report` from a release-binary install, leaving no `feedback` or `report-*` directory behind.
- The edge case: `sh-agent update` with a swapped signature file refusing and leaving the binary untouched, run against the real release rather than a fixture.
- Whether both CLIs start cleanly from a home installed by a release binary, which needs the sandbox home that no longer exists here.

## Added 2026-09-18 — v1.2.0, and the first real cross-release update

v1.2.0 was cut with the CORRECTED procedure (create the release with notes alone, then one upload per asset). It published first time, with no rollback — the shape v1.1.0's two `HTTP 500`s taught.

- **Published and verified as a consumer:** `gh release download v1.2.0` (all five assets), `sha256sum -c SHA256SUMS` → OK for all three binaries, `ssh-keygen -Y verify` over the file as published → good signature, `./sh-agent-linux-amd64 version` → `sh-agent v1.2.0`.
- **The cross-release update path, which could not exist until now.** A scratch place was installed from the v1.1.0 binary (`source embedded v1.1.0`), then `sh-agent update` ran against the real GitHub: `downloading sh-agent-linux-amd64 v1.2.0`, `verified and installed`, the new binary's install, and `4 written, 256 unchanged`. The binary's hash changed, and `status` now reports `source embedded v1.2.0`, `binary v1.2.0`, drift none.
- **Those four files are the point of the release:** `context/AGENTS.md` and `skills/preflight/SKILL.md`, one copy per vendor. A binary-only place picking up new rules is exactly what the release binary was built for, and this is the first time it has happened end to end.
- **The embed-contamination guard ran twice**, as the corrected guide says: once among the pre-tag checks and once immediately before the builds.

### Still only the Mac can settle these

Unchanged from the section above, now owed against v1.2.0 rather than v1.1.0: running any `darwin/arm64` binary at all, the macOS `xattr` quarantine note, `report` from a release-binary install leaving nothing behind, the swapped-signature edge case, and a clean start of both CLIs from a home installed by a release binary.
