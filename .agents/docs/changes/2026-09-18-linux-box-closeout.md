---
owner: main
date: 2026-09-18
---
# Everything the Linux box can settle about v1.1.0, settled     (2026-09-18)

## Summary

The operator asked for everything that can be done here before they update their Mac. Four things closed, one defect found and fixed.

`go install github.com/kesokaj/styrman-harness/cmd/sh-agent@v1.1.0` works, which closes the last assumption in the plan's Risks — and it exercised the one version path nothing else could reach: that binary carries no `-ldflags`, so its `v1.1.0` comes from `debug.ReadBuildInfo().Main.Version` (D6, D21). Installing from it wrote 260 files with `source embedded v1.1.0` and drift none. All three published binaries were downloaded and verified against the signed `SHA256SUMS` — `OK` for each, good signature over the file as published, and `file` confirms the cross-compiled ones are Mach-O arm64 and ELF aarch64, though neither can run here. The README's `gh release download` path runs verbatim.

The defect: the README's `go install` line, as shipped in v1.1.0, does not work. The comment said "over SSH" but the command gave git no way to use SSH, so on a machine without an insteadOf rewrite it fails with `could not read Username for 'https://github.com'`. It is BUG-075, fixed the same day: README now gives a form that supplies the rewrite for that one command through `GIT_CONFIG_*` and changes no git config, and mentions the permanent alternative as the reader's own choice.

Also done: `sh-agent update` on the operator's real home at their request (up to date, 260 files unchanged, the old kept feedback clone correctly kept), and that clone measured — it holds nothing unique, so it is theirs to delete. The misleading message that prompted the measurement is BUG-074, fixed yesterday.

## Type

bugfix

## Files changed

- `README.md` — § Install without cloning's `go install` block replaced with the form that runs
- `.agents/docs/bugs/BUG-075-readme-go-install-ssh.md` (new) and `INDEX.md` — the record and its row
- `.agents/docs/testing/2026-09-17-release-binary-e2e.md` — a § Added 2026-09-18 section with every run above, and what only the Mac can settle
- `.agents/docs/plans/2026-09-17-release-binary.md` — the `go install` Risks row marked settled; the three Phase 5 boxes that were `[~]` for owed work now `[x]` (the rehearsal ran before the release, the notes were built and read, the planted-file refusal was exercised)

## Agents dispatched

`bugs`, for BUG-075; its record was read against the README diff and my own measurements before committing (context § Definition of Done step 4). No reviewer: the code diff is empty — this change is one README block, one bug record, and trail updates.

## Verification

`gofmt`, `go vet`, `go test -race -count=1 ./...` clean. Every claim above is a command that ran, with its output in the e2e log: the `go install` (exit 0, `sh-agent v1.1.0`, `go version -m` showing no `-ldflags`), the install from that binary (260 files, drift none), `sha256sum -c` over all three published binaries, `ssh-keygen -Y verify` over the published `SHA256SUMS`, the README `gh` line verbatim, and `sh-agent update` on the real home. NOT verified here, and only the Mac can: running any `darwin/arm64` binary, the macOS quarantine note, `report` from a release-binary install, and the swapped-signature edge case against the real release.

## Docs

`README.md` is the docs half and is included. The e2e log now separates what this box settled from what the Mac still owes, so the next session does not re-run the former or claim the latter.

## Requirements

None touched.
