---
owner: main
date: 2026-09-17
---
# Release binary — Phase 2 of the release-binary plan: the embedded source     (2026-09-17)

## Summary

The harness's own files now live inside the binary, and every reader reaches them through one abstraction. `embed.go` (package `styrman`, module root) embeds `context/`, `skills/`, `agents/`, `commands/` and `rules/` with `//go:embed all:…`, so the 15 `.gitkeep` files come along. `Source{FS, Dir, Version}` replaces the `repo string` every reader took: `Dir` is a checkout or plain copy (`FS` is `os.DirFS(Dir)`), and an empty `Dir` with a `Version` is the files built into a release binary. The symlink refusals of BUG-036 run when `Dir` is set — `embed` cannot hold a link — and are unchanged in meaning for a tree on disk. `ResolveSource` replaces `FindRepo` and adds one last candidate after today's four: the embedded files, but only when the running binary's version is a `vX.Y.Z` tag (D14), so a development build outside a tree keeps today's error instead of installing a stale copy. An explicit `--repo` or `$SH_AGENT_REPO` naming something that is not a harness tree is still refused, never replaced by the embedded files (D15). With an embedded source, install copies the running binary and never builds (D13), `state.json` records `"repo": "embedded:vX.Y.Z"` with no commit, `status` prints a `source     embedded vX.Y.Z` line instead of the `repo` line, `diff`, `list` and `render` read the embedded files, `adopt` and `update` refuse with a message naming the release, the report record carries a new `harness_version` field with `harness_commit` null, and no git helper is called with an empty path (D17). Phase 4 replaces `update`'s refusal with the release flow (D22).

## Type

task

## Files changed

- `embed.go` (new, committed in 7f3c0f9) — package `styrman`, `//go:embed all:context all:skills all:agents all:commands all:rules`
- `embed_test.go` — `TestFilesMatchTree` (disk ⊆ embed, `.gitkeep` included); `TestEmbedHoldsOnlyTrackedFiles` (embed = `git ls-files`, both directions) from the security review's finding 1
- `plan_test.go` (new, package `styrman_test`) — `TestEmbeddedPlanEqualsCheckout`: the embedded plan equals the checkout plan in paths, bytes, modes and KeepForeign, for all three vendors
- `internal/harness/source.go` (new) — `Source`, `DirSource`, `EmbeddedSource`, `EmbeddedPrefix`, `IsReleaseVersion`, `Source.path`, `Source.readRegular`
- `internal/harness/source_test.go` (new) — MapFS plan, nine refusals from a MapFS, non-regular `SKILL.md`, DirSource = MapFS, a linked `skills/<name>` followed and a dangling one refused, embedded install state and report record
- `internal/harness/inventory.go`, `render.go`, `plan.go` — every reader takes a `Source`; `fs.ReadDir`/`fs.Glob`/`fs.ReadFile`/`fs.WalkDir` in place of the `os`/`filepath` calls; `checkManagedTree` runs only for `Dir != ""`
- `internal/harness/repo.go` — `ResolveSource` replaces `FindRepo`; `absSource`
- `internal/harness/apply.go` — `Install` takes a `Source` and records `embedded:<version>` with no commit for one
- `internal/harness/report.go` — `ReportRun.Source`, `BuildRecord(home, Source, …)`, `Record.HarnessVersion` (`harness_version`)
- `cmd/sh-agent/main.go` — `var version` (release `-ldflags`), `ResolveSource`, `installFiles`, the per-command embedded branches, `releaseFromInfo`, `knownCheckout` refuses an `embedded:` manifest, usage text names the new source
- `cmd/sh-agent/main_test.go` — `TestReleaseBinaryWithoutCheckout`, `TestReleaseFromInfo`, `TestKnownCheckoutRefusesEmbeddedManifest`
- `internal/harness/harness_test.go`, `pathguard_test.go`, `report_test.go` — calls moved to `DirSource(repo)`, assertions unchanged; `TestFindRepo`/`TestFindRepoRefusesNamedNonRepo` became `TestResolveSource`/`TestResolveSourceRefusesNamedNonRepo`
- `.agents/docs/plans/2026-09-17-release-binary.md` — Phase 2 boxes ticked with evidence; D19–D23; one new Phase 5 box (preflight refuses untracked or ignored files under the five directories) and one for the `AGENTS.md` `govulncheck` call-path note

## Agents dispatched

- `verifier` — PASS, `.agents/docs/testing/2026-09-17-release-binary-phase-2.md`
- `reviewer` — APPROVE, 0 findings, `.agents/docs/reviews/2026-09-17-review-f7ee817-21ba8dc-code.md`
- `security` — 1 medium, 6 low, `.agents/docs/reviews/2026-09-17-security-f7ee817-21ba8dc.md`; three fixed, five rebutted in its `## Outcome`, one raised to the operator

## Verification

`gofmt -l cmd internal *.go` silent, `go vet ./...` clean, `go test -race -count=1 ./...` green (verifier: 81 pass, 0 fail, 0 skip — the one former skip, `TestFindRepo`, became `TestResolveSource`, which moves to a temporary working directory and always runs). Each guard was seen failing before it was trusted, then restored: the link refusals disabled → four symlink tests fail; `skills` embedded without `all:` → 52 files planned from the tree only; a skill file made 0755 → a mode mismatch on both vendors; the release-version gate removed → `TestResolveSource` and `TestReleaseBinaryWithoutCheckout` fail; the `embedded:` skip removed → `TestResolveSource` fails; an explicit non-tree falling through → two tests fail; an ignored `rules/.env` → `TestEmbedHoldsOnlyTrackedFiles` fails while every other test stays green; `knownCheckout`'s new guard removed → its test fails. Integration, run from `/` with `SH_AGENT_REPO` unset, from a `CGO_ENABLED=0 go build -trimpath -ldflags "-X main.version=v0.0.0"` binary: install wrote 260 files, `status` showed `source     embedded v0.0.0`, `from release v0.0.0` and drift none, `state.json` recorded `embedded:v0.0.0` with no commit, `update` refused with exit 1; the same homes compared against a checkout install have the same file set, bytes and modes. Not exercised: a `go install …@vX.Y.Z` binary (no tag carrying the embed exists yet), the `binary` line of a clean-tag build (the test build was from a dirty tree, so it showed a revision), the vendor-CLI smoke runs (no `context/`, `agents/` or `skills/` content changed), and everything Phases 3–5 own.

## Docs

No living doc changed. `README.md` § Install without cloning, the map, `AGENTS.md` rules 1 and 2, `REQUIREMENTS.md` R8/R9, the deploy guides and the delivery spec land together in Phase 5's single commit to `main` (D11, R9).

## Requirements

R8 and R9 still describe two ways in, one of them the unpacked bundle; the approved replacement text sits in the plan's Interfaces and is applied in Phase 5, in one commit with the code, the docs and the map. Nothing in this phase is reachable on a machine until a release carrying the embed exists, so no R8/R9 claim is broken meanwhile. R4 holds: stdlib only, no dependency added.
