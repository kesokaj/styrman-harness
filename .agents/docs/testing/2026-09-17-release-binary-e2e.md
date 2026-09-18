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

## Added 2026-09-18 — the work Mac, on v1.2.0 (the operator's own runs)

Host `spgo-mac.roam.internal`, darwin/arm64, Antigravity only. Run by the operator from `/tmp`, i.e. with no harness checkout anywhere near — the shape the whole feature exists for. **This is the first time a `darwin/arm64` binary has been executed at all.**

- **Download and run.** `gh release download v1.2.0 --pattern 'sh-agent-darwin-arm64' --dir ~/.local/bin`, `mv`, `chmod +x` → `sh-agent v1.2.0`, 6709378 bytes, byte-identical in size to the published asset. No `xattr` step needed: a `gh` download is not quarantined, as README says.
- **Install.** `sh-agent install --vendor antigravity` → `release binary v1.2.0: installing the harness built into it and the running binary`, `2 written, 128 unchanged, 0 removed, 0 kept — 130 files managed`. 130 is one vendor's half of 260, and the 2 written are `context/AGENTS.md` and `skills/preflight/SKILL.md` — exactly v1.2.0's delta over what that home already had.
- **Status.** `source     embedded v1.2.0 (release binary; update downloads the newest release)`, `installed … from release v1.2.0 (vendor: antigravity, 130 files …)`, `binary     v1.2.0`, `drift      none — 129 files match the repo`.
- **Report, the first ever from a release-binary install.** `--dry-run` named the throwaway clone (`report-7033`, blob-less, limited to `reports/spgo-mac.roam.internal/`) and the note `harness is the release binary v1.2.0, not a checkout: harness_commit is null, harness_version is set`. The real run sent `reports/spgo-mac.roam.internal/2026-09-18-103648 (c4b13a8)`. Verified independently from the Linux box by cloning the feedback repository: that commit exists and its `report.json` reads `{"host":"spgo-mac.roam.internal","os":"darwin","arch":"arm64","vendor":"antigravity","files":130,"harness_commit":null,"harness_version":"v1.2.0","drift":0,"skills":27,"agents":7,"go":"go1.27.1"}` — the `harness_version` shape existed only in fixtures until now. The push took about 1m 9s.
- **Nothing kept.** `ls -la ~/.config/sh-agent/` afterwards: `state.json` and no `report-*` directory. The throwaway was created and removed on a second machine.
- **`update --check`** → `up to date (v1.2.0)`, exit 0.
- **`adopt`** → `adopt needs a harness checkout; this install is a release binary`. (A typo, `adpot`, first printed the usage text, which incidentally confirmed v1.2.0's updated `update` and `report` descriptions shipped.)
- **CLI start.** `agy -p 'List three skills you have.'` from a scratch directory answered cleanly, naming `golang`, `typescript` and `cloud-run-deploy` with paths under `~/.gemini/config/skills/` — the files the binary carried. This is the tier the Linux box could not run, because its sandbox home no longer exists.

### What the Mac exposed

`ls ~/.config/sh-agent/` also showed `feedback/`, dated 2026-09-17 — the full clone earlier versions kept. `report` ignores it by design, but `update` is supposed to remove it, and the plan says so for BOTH kinds of place (§ Old kept clone cleanup). Only the checkout path called `RemoveKeptClone`; the release-binary branch returned before reaching it. So a binary-only place would have kept holding every other machine's reports — the exact complaint this feature began with, surviving in the one place it was meant to fix. Fixed the same day with a test that fails without the call; it reaches that Mac in the next release.

### Still not exercised anywhere

- The signature edge case on darwin (a swapped `SHA256SUMS.sig` refused). Proven on Linux against the real releases; not run on the Mac.
- The Claude Code side on that Mac: it was installed `--vendor antigravity` only, so `~/.claude` there holds nothing from v1.2.0 and `claude` was not started from it.
- `linux/arm64`: built and published for both releases, never executed on any machine.

## Added 2026-09-18 — v1.2.1 on two machines the operator ran himself

Both runs are from the operator's own shell, on machines with no harness checkout in play. Between them they close every leg that had been owed since v1.1.0 except three, named at the end.

### devbox (Debian, linux/amd64) — first release-binary install on this machine

- **Download verified before anything was moved.** `gh release download v1.2.1` (binary, `SHA256SUMS`, `SHA256SUMS.sig`) into a fresh `mktemp -d`; `ssh-keygen -Y verify -n sh-agent-release` → `Good "sh-agent-release" signature for ostling.simon@gmail.com with ED25519 key SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU`; `sha256sum -c` → `sh-agent-linux-amd64: OK`; `./sh-agent-linux-amd64 version` → `sh-agent v1.2.1`, run from the download directory before the install.
- **The install was then lost to PATH shadowing** — recorded as BUG-078, and the reason this section exists. `/usr/local/bin/sh-agent` held the verified v1.2.1 binary (7052333 bytes, and `/usr/local/bin/sh-agent version` printed `v1.2.1` when called by absolute path), but `which -a sh-agent` found `/home/admin/.local/bin/sh-agent` (PATH 1) and `/opt/devbox/bin/sh-agent` (PATH 5) ahead of it at PATH 11. Unqualified `sh-agent` therefore ran a `(devel)` build, which resolved the stale plain copy at `/opt/devbox/deploy/harness` and installed THAT — reporting `drift none` the whole time, because drift is measured against whatever source the running binary resolved.
- **Proof the content was stale, and proof it stopped being so.** With the shadowed binary, `claude -p 'Quote the rule about a guard in front of a destructive action.'` quoted the `rm -r` bullet — the fail-closed guard rule at `context/AGENTS.md:58` was not in that tree. After `/opt/devbox/bin/sh-agent` was pointed at `/usr/local/bin/sh-agent`, `grep -c 'A guard in front of a destructive' ~/.claude/AGENTS.md` → `1`.
- **Install, once the right binary ran.** `sh-agent install --vendor=all` → `release binary v1.2.1: installing the harness built into it and the running binary`, `29 written, 231 unchanged, 2 removed, 0 kept — 260 files managed`. The 2 removed are files the stale tree had and v1.2.1 does not.
- **Status.** `source embedded v1.2.1 (release binary; update downloads the newest release)`, `installed 2026-09-18T10:26:00Z from release v1.2.1 (vendor: all, 260 files …)`, `binary v1.2.1`, `drift none — 259 files match the repo`.
- **Report, from a release binary, `--vendor all`.** `report: harness is the release binary v1.2.1, not a checkout: harness_commit is null, harness_version is set`; `skipped feedback.jsonl (absent)`; sent `reports/devbox/2026-09-18-102654 (1aefa0e)` — 9 files, 1627 insertions, six queued retro proposals among them.
- **Nothing kept.** `ls -la ~/.config/sh-agent/` afterwards: `proposals/`, `retros.jsonl`, `state.json`. No `feedback/`, no `report-*`. First time this has been shown on a `--vendor all` release-binary place.

### The work Mac (darwin/arm64) — the kept clone, shed for real

- **`sh-agent update` on v1.2.1** → `removed the kept feedback clone /Users/spgo/.config/sh-agent/feedback (report no longer keeps one)`, then `up to date (v1.2.1)`. This is the fix for what § What the Mac exposed found on 2026-09-18: the release-binary branch of `update` returned before `RemoveKeptClone`. The clone that machine had held since 2026-09-17 is gone.
- The preceding `update` on that Mac, which fetched v1.2.1 itself, could not have shed it: `git show v1.2.0:cmd/sh-agent/main.go | grep -c shedKeptClone` → `0`, against `4` for v1.2.1. `update` replaces the binary and then runs the NEW binary's `install`, never its `update`, so the shedding waited for a second `update` — worth knowing for any future cleanup that ships this way.
- **`ls -la ~/.config/sh-agent/`** → `state.json` alone.
- **`report`** → sent `reports/spgo-mac.roam.internal/2026-09-18-122933 (789e2d9)`, two files, with the same release-binary metrics note. `feedback.jsonl`, `retros.jsonl` and `proposals/` were all absent and each was named as skipped.

### What these two runs settle

- A `darwin/arm64` binary and a `linux/amd64` release binary both run, install, report and update on machines with no checkout.
- `update` sheds the old kept clone on a release-binary place — the last piece of "report keeps nothing on the machine".
- A release delivers context rules to a binary-only place: the fail-closed guard rule, graduated on 2026-09-17, is now on both machines and was quoted back from neither before this.

### Still not exercised anywhere

- The swapped-signature edge case on darwin (a tampered `SHA256SUMS.sig` refused, binary untouched). Proven on Linux against the real releases; never run on the Mac.
- `linux/arm64`: built and published for all three releases, never executed on any machine.
- The Claude Code side of the work Mac: still installed `--vendor antigravity` only, so `~/.claude` there holds nothing from any release and `claude` has not been started from it. devbox covers the Claude side on Linux.
- `sh-agent update` has never run on devbox: that machine reached v1.2.1 by download and `install`, not by updating from an earlier release.

## Added 2026-09-18 — v1.2.2 published and consumed

v1.2.2 (release commit f3ec5dc on gated head f91257c, tag signed and verified before the push) was cut by `/preflight v1.2.2` from the operator's own turn and published with the create-then-upload-per-asset shape: five assets, each first attempt. Gate: `.agents/docs/reviews/2026-09-18-preflight-v1.2.2.md`.

- **Published and verified as a consumer**, from a fresh scratch directory: `gh release download v1.2.2` (all five assets) → `ssh-keygen -Y verify` over `SHA256SUMS` as published → `Good "sh-agent-release" signature`; `sha256sum -c SHA256SUMS` → OK for all three binaries; `./sh-agent-linux-amd64 version` → `sh-agent v1.2.2`; and the new binding the docs now ask for, `test "$(./sh-agent-linux-amd64 version)" = "sh-agent $(sed -n 's/^# styrman-harness //p' SHA256SUMS)"` → exit 0. The downloaded `linux/amd64` binary is byte-identical to the one built here (`cmp`).
- **Repository visibility** after the upload: `PRIVATE`.

### Not exercised

- `sh-agent update` from v1.2.1 to v1.2.2 on any place: the operator's Mac and devbox are the binary-only places, and both updates are theirs to run (the handoff's Next Steps). The cross-release update path itself is proven for v1.1.0 → v1.2.0 (Linux, scratch) and v1.2.0 → v1.2.1 (the Mac).
- `darwin/arm64` and `linux/arm64` v1.2.2 binaries: built, checksummed, signed and published; not executed.
- The README browser chain on a Mac, including its `shasum -a 256 -c - --ignore-missing` form.

### Added 2026-09-18 — the v1.2.1 → v1.2.2 update, on all three machines (the operator's runs)

- **Work Mac (darwin/arm64):** `sh-agent update` → `downloading sh-agent-darwin-arm64 v1.2.2`, `verified and installed sh-agent-darwin-arm64 v1.2.2`, the new binary's install `0 written, 130 unchanged` — correct, v1.2.2 changed no embedded file; `sh-agent version` → `sh-agent v1.2.2`; `install --vendor antigravity` → `0 written, 130 unchanged`. Clean.
- **devbox (linux/amd64):** `sh-agent update` → `downloading sh-agent-linux-amd64 v1.2.2`, `verified and installed …`, `0 written, 260 unchanged`; then **`sh-agent version` → `sh-agent v1.2.1`**. `which -a` showed `~/.local/bin` first and `/usr/local/bin` last, so PATH order was right and the BUG-078 guard correctly printed nothing. The `update` had been run BY `/usr/local/bin/sh-agent` (v1.2.1), which `UpdateFromRelease` never replaces — it writes `~/.local/bin/sh-agent` — and bash's command hash kept pointing at `/usr/local/bin`; after `sudo rm` of that file the shell said `-bash: /usr/local/bin/sh-agent: No such file or directory`, which is the hash table, not PATH. `source .bashrc` reset it: `sh-agent version` → `v1.2.2`, `install --vendor all` → `0 written, 260 unchanged`, `status` → `source embedded v1.2.2`, `binary v1.2.2`, drift none, no `shadowed` line. **The documented unqualified `sh-agent version` check is what caught it.** Recorded as BUG-079.
- **This Linux box (checkout place):** `sh-agent update` → `up to date with upstream (790550c)`, rebuilt from the checkout, `0 written, 260 unchanged`; `sh-agent version` → `sh-agent 790550c1f7a6`; `status` → `binary 790550c1f7a6`, drift none, no `shadowed` line. The kept feedback clone here is still declined by `update` (two commits it cannot prove pushed), as intended.

Still not exercised: the README browser chain on the Mac; `linux/arm64`; the swapped-signature edge on darwin.

## Added 2026-09-18 — v1.2.3 published and consumed

v1.2.3 (release commit 91fcf13 on gated head 30f7488; gate `.agents/docs/reviews/2026-09-18-preflight-v1.2.3.md`, SHIP) was cut by `/preflight v1.2.3` on the operator's "approved go go go" and published create-then-upload-per-asset: five assets, each first attempt; repository `PRIVATE` after.

- **Published and verified as a consumer**, from a fresh directory, with the by-hand chain's own checks as README now writes them: `ssh-keygen -Y verify` → good signature; `grep -qxF "# styrman-harness v1.2.3" SHA256SUMS` → 0; `sha256sum -c` → OK ×3; `./sh-agent-linux-amd64 version` = `sh-agent v1.2.3`; byte-identical to the build.
- **The published binary installed into a scratch home**: `release binary v1.2.3: installing the harness built into it and the running binary`, `260 written`; `status` → `source embedded v1.2.3`, `binary v1.2.3`, drift none; `~/.config/sh-agent` created `700`.

### The v1.2.2 → v1.2.3 update, on all three machines (the operator's runs, 2026-09-18)

Logged cell by cell in `sh-agent-matrix.md` § Real machines: this box (checkout place), the work Mac (darwin/arm64 release binary), devbox (linux/amd64 release binary). Each printed the `mode ~/.config/sh-agent (0755 -> 0700)` line on its first run, `version` unqualified reported the new version, `type` named the installed path, `status` drift none with no `shadowed` line, and the runtime dir was `700` after. The `darwin/arm64` binary has now run on the Mac; `linux/arm64` still has not run anywhere.

### Not exercised

- The README browser chain on the Mac (`shasum -a 256 -c - --ignore-missing`, Gatekeeper, `xattr`) — the operator's, when convenient.
- `linux/arm64`: built, signed, published for every release since v1.1.0; never executed.
- The darwin swapped-signature edge case of `update`: proven on Linux against the real releases; not run on the Mac.
