---
owner: main
date: 2026-09-30
---

# Merged review — `/review --project --all`

Range reviewed: the whole repository at HEAD (3478a1c), no diff.

No FAIL. Every lens passed or approved; the findings are one medium (simplify) and five low (security), nothing above.

## Verdict per lens

| Lens | Agent | Verdict | Artifact |
|---|---|---|---|
| code | reviewer | APPROVE — no correctness, security or rule finding | [review-project-code](2026-09-30-review-project-code.md) |
| security | security | FINDINGS, all low; no high or medium; govulncheck 0 reached | [security-project](2026-09-30-security-project.md) |
| simplify | reviewer | APPROVE WITH NITS — one medium duplication | [review-project-simplify](2026-09-30-review-project-simplify.md) |
| docs | reviewer | APPROVE — no statement out of step with the code | [review-project-docs](2026-09-30-review-project-docs.md) |
| prompts | reviewer | clean — no high/medium finding, empty proposed diff | [review-project-prompts](2026-09-30-review-project-prompts.md) |
| tests | verifier | PASS — 366 pass, 0 fail, 2 skip with reasons; every source file mapped to a covering test | [verify-project](../testing/2026-09-30-verify-project.md) |
| recon | explorer | brief written; no orphan skill, empty `commands/` and `rules/` noted | [recon-project](../plans/2026-09-30-recon-project.md) |

## Findings

| # | Severity / confidence | Lens | Location | Finding |
|---|---|---|---|---|
| 1 | medium / high | simplify | `internal/harness/apply.go:309-347`, `cmd/sh-agent/main.go:908-935` | The atomic write sequence (temp file, write, sync, chmod, rename) is implemented twice; apply.go's own comment records that the fsync fix landed in the settings writer first and was ported by hand. Proposed: one `internal/harness` helper returning `error`, with `attributionSet` wrapping it. |
| 2 | low / high | security | `internal/harness/release.go:144-148` | `update` sweeps leftover `release-*` download directories with `filepath.Glob` on a pattern built from the home path; a home containing `*`, `?` or `[` matches entries in other directories, and every match is removed with `os.RemoveAll` regardless of type or age. The two sibling sweeps (`report.go:419`, `main.go:1171`) already use `os.ReadDir` plus a prefix match. Proposed: the same, plus a test with a metacharacter in the home path. |
| 3 | low / high | security | `internal/harness/repo.go:197,225,241` | A checkout place builds and runs the tip of `main` after `git pull --ff-only` with no signature check; release places verify SSHSIG. By design (R9, AGENTS.md rule 3). Operator choice: record as accepted risk in an ADR line (recommended), or verify commit signatures after the pull. |
| 4 | low / medium | security | `internal/harness/report.go:39-58` | The report redaction deny-list misses JWTs, `sk_live_`/`rk_live_`/`whsec_`, `npm_`, `hf_`, `pypi-`, legacy `sk-` keys, Slack webhook URLs and `https://<token>@host` without a colon. Destination is the operator's private repository, hence low. Proposed: anchored `regexp` needles, each with a test that fails without it. |
| 5 | low / high | security | `skills/VENDORED.md` (Apache-2.0 table) | `gke-basics` and `terraform` have no recorded source, pin or digest; every other vendored row has one. Proposed: re-vendor from a pinned upstream commit and record the pin. |
| 6 | low / low | security | `internal/harness/report.go:901-914` | Push errors return git's stderr unscrubbed, while the clone failure path wraps the same kind of message in `scrubCredentials`. Proposed: wrap both push-error returns. |
| 7 | low / low | simplify, code | `cmd/sh-agent/main.go` (1199 lines) | The entrypoint holds the settings.json attribution editor (~230 lines) and the PATH-shadow reporter, neither CLI-dependent; the golang skill's "cmd/ stays thin" item. Both lenses call it non-blocking; the single-file grouping is also defensible. |

Observations that are not findings, recorded once:
- `deadcode ./...` flags `internal/harness/release.go:422 SwapReleaseSigner`; it is the documented cross-package test seam used by `cmd/sh-agent/zip_test.go` (code and simplify lenses agree: tool false positive).
- `commands/` and `rules/` hold only `.gitkeep` since commit fc89509 while the project AGENTS.md layout table describes them as holding deployed prompts and rules (explorer). Both are wired into `LoadDocs`, so the doc line is true of the mechanism and empty of content; a one-line operator check whether that is intended scaffolding.
- Two prompt-lens flags below the action threshold: `skills/mermaid-syntax/SKILL.md:159` (a Mermaid v11.5.0 upgrade note, uncontradicted) and the 2026-09-13 format-cutover note in `skills/feedback-log/SKILL.md:56` and `skills/distill/SKILL.md:156` (migration-relative but functionally necessary).
- Docs lens checked and did not flag: `testing/guide.md` omits the BUG-087/090/091/106/107 behavioural-run rules because they graduated into AGENTS.md § Testing; `deploy/remote.md:203`'s "27 skills and 7 agents" is framed as the v1.0.0 gate's historical count.

## Not checked

- Code lens: `apply.go`, `frontmatter.go`, `inventory.go`, `plan.go`, `release.go` body, `render.go`, `repo.go`, `source.go`, `state.go`, `internal/search/parse.go` and `store.go`, `embed.go` body, all `*_test.go`, and all markdown prose; `go test -race` was left to the verifier.
- Security lens: `pathguard.go`, `source.go` beyond `SH_AGENT_REPO`, `inventory.go`, `frontmatter.go`, `render.go`, `owned_*.go`, most of `state.go` and `apply.go`; `main.go` lines 1-94, 150-394, 416-635, 700-799, 940-1039; `internal/search/parse.go`, most of `search.go`, `reconcile.go`'s `gitView` and apply loop; test files; a full read of the skills and context text.
- Simplify lens: prose of `context/`, `agents/`, `commands/`, `rules/`, skill bodies; `store.go` schema versioning and `report.go` clone/push flow beyond signature level; test files.
- Docs lens: `specs/sh-agent/search.md` and `memory-loop.md`, `context/AGENTS.md` and `dev-setup.md`, individual agent and skill files, the testing matrices, the bulk of `changes/`, DECISIONS.md content beyond heading existence, `plans/`, `retros/`, `handoffs/`, `*_test.go`.
- Prompts lens: 19 of the 22 non-vendored SKILL.md files were covered by a signal grep, not a sentence-by-sentence read (full reads: `bug`, `feature`, `preflight`, the context, the root AGENTS.md, all eight agents).
- Tests lens: the sandbox and real-home CLI smoke tiers, `owned_other.go` (`!unix` build tag, not compilable on this box), the search recall and performance tier (no search change in scope), the report dry-run tier, the repo-wide suppression grep (no diff to compare against).
- Explorer: no staticcheck/deadcode-grade sweep of unexported functions; individual plan files not checked for orphan status against changes, handoffs and bug records; `.agents/docs/{changes,reviews,testing,bugs}` not inventoried.
