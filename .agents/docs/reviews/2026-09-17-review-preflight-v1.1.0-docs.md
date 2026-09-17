---
owner: reviewer
date: 2026-09-17
---
# Review — preflight v1.1.0 docs-lens   (2026-09-17)

Range reviewed: v1.0.0..086deab

## Verdict: APPROVE

No blocking or advisory-severity contradictions found between the living docs and the code at HEAD. The release-binary rework (two ways in, `sh-agent update` on a release binary, `sh-agent report` keeping no clone, the `feedback` → `feedback-log` rename) is documented consistently across `README.md`, `REQUIREMENTS.md` R8/R9, `AGENTS.md` § How the harness works, `.agents/docs/deploy/{local,remote}.md`, `.agents/docs/specs/sh-agent/delivery.md`, `.agents/docs/testing/guide.md`, and the living skills. This is a wide range (38 commits, 148 files) but nearly all of the docs churn in it was already the subject of two prior, narrower `/review --docs` passes in this same range (`.agents/docs/reviews/2026-09-17-review-edbddca-1a5578b-docs.md` and `2026-09-17-review-4c63232-feedback-log-rename.md`), whose findings are fixed and whose fixes I re-verified independently against HEAD rather than trusting the prior reports.

## Findings

None. Checked specifically and found no contradiction:

- **README map vs R8/R9 vs `AGENTS.md` rules 1–2** (finding class 3): all four texts describe the same two-ways-in-one-way-back model; the plain-copy/"unpacked bundle" fallback is consistently and explicitly demoted to "no longer one of the ways in" in `README.md:135`, `.agents/docs/deploy/remote.md:28`, `.agents/docs/deploy/local.md:53`, and `.agents/docs/specs/sh-agent/delivery.md` — none of them present it as a way IN (class 2 check: negative, correctly).
- **"A tag delivers nothing"** (class 2): not found anywhere in the living docs; `AGENTS.md` rule 2 and `.agents/docs/deploy/remote.md` § Releases both now say a tag is a release marker that also carries binaries a binary-only place updates to, while a checkout place still follows `main` — consistent with each other and with `internal/harness/release.go`.
- **Kept feedback clone** (class 2): not found; `local.md` § Report back explicitly documents the old `~/.config/sh-agent/feedback` clone as something `update` now *removes* (when clean and fully pushed), matching `README.md:116`'s daily-commands row for `update` and the `report` code path (no clone kept; throwaway clone removed after push).
- **`feedback` skill rename** (class 2): grepped every living doc (`README.md`, `AGENTS.md`, `REQUIREMENTS.md`, `ROADMAP.md`, `context/*.md`, every `skills/*/SKILL.md`, `agents/*.md`) for the bare skill name `feedback` outside of `feedback-log`, `feedback.jsonl`, and `styrman-harness-feedback` (the repo name, correctly unchanged) — zero hits. `skills/feedback-log/`, `skills/distill/SKILL.md`, `skills/retro/SKILL.md`, `skills/handoff/SKILL.md` all reference the skill correctly by its new name or by "feedback log" (the file, not the skill).
- **Commands/flags a reader would run** (class 4): every command in `.agents/docs/deploy/remote.md` § Releases checked against `internal/harness/release.go` and `embed_test.go`/`release_test.go` — `ReleaseNamespace = "sh-agent-release"`, `releaseSigner` string, `sumsVersionMarker` regex (`^#\s*styrman-harness\s+vX.Y.Z\s*$`, and `sha256sum -c`'s documented `#`-comment tolerance), `TestEmbedHoldsOnlyTrackedFiles` (exists at repo root, package `styrman`), and `TestReleaseSignerMatchesTheDeployGuide` (pins the signer line in the deploy guide against the compiled-in constant, so drift here is caught by `go test`) all match what the doc tells the operator to run. `README.md`'s Daily commands table and `cmd/sh-agent/main.go`'s usage/flag strings agree verb-for-verb and flag-for-flag (`install`/`sync`, `update --check`, `report --dry-run --feedback-repo`, etc.), and skill/agent counts (27 skills, 7 agents) match `ls skills | wc -l` / `ls agents | wc -l`.
- **Skills/agents describing changed harness behaviour** (class 5): `skills/preflight/SKILL.md` is generic (applies to any project) and does not describe this repo's specific delivery model, so it needed no change here and has none contradicting the code. No agent (`agents/*.md`) describes delivery/update/report mechanics.

## Criteria

Not the six-criteria code lens — this dispatch is the docs-vs-code lens (five finding classes above), all checked negative.

## Not examined

- The full 148-file / 38-commit diff was not read line-by-line; I read the living docs in full and cross-checked their claims against the relevant code (`internal/harness/release.go`, `cmd/sh-agent/main.go`, `embed_test.go`, `release_test.go`) rather than diffing every commit.
- `context/dev-setup.md` and the rest of `context/AGENTS.md` beyond the learning-loop/git sections were skimmed, not read end to end for unrelated drift (they are not central to this release's delivery-model change).
- `.agents/docs/architecture/DECISIONS.md` (ADR-030/031/033) and the `changes/`, `bugs/`, `reviews/`, `plans/`, `handoffs/` trail directories were sampled via grep, not read as living docs — they are append-only trail per this repo's own rules, out of the dispatch's "living docs" list, and not expected to be edited for currency.
- Did not run any command (per the hard limits); all checks were read-only comparison of doc text against source/test text.
- Did not verify the Antigravity-derived `agent.md` output or run `sh-agent install`/`status` against a scratch home; relied on static reading of `internal/harness/derive.go`'s existence being unchanged in a way that would contradict docs (not itself inspected in depth this pass).
