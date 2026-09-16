---
owner: reviewer
date: 2026-09-16
---
# Review — bug-036-docs (2026-09-16)

Range reviewed: `git diff HEAD -- internal/harness/inventory.go internal/harness/render.go internal/harness/pathguard_test.go` (working tree vs HEAD `d167cf8`). This overwrites my earlier same-slug report: the fix has widened since that pass from a leaf-only `readRegular` check to `checkManagedTree`, which `Lstat`s each of the six managed directories (`context`, `agents`, `commands`, `rules`, `workflows`, `skills`) themselves and every entry directly inside the five flat ones, called from `LoadInventory` and `RenderContext` before any reader runs. `skills/<name>` stays the one documented exception, resolved via `filepath.EvalSymlinks` in `LoadSkills`, with its own `SKILL.md` still guarded by `readRegular`.

## Verdict: APPROVE

No living doc says or implies that `context/`, `agents/`, `commands/`, `rules/`, `workflows/` or `skills/` itself — or a file directly inside the first five — may be a symlink, and none documents a workflow that would leave one there. The one path the code still permits a link at, `skills/<name>`, is exactly the path every doc that mentions symlinks describes (`.agents/skills` resolving into `skills/`, and `sh-agent adopt` copying — never symlinking — into `skills/<name>`). `AGENTS.md`'s own claim that `.agents/skills` is a symlink is untouched by this diff (`checkManagedTree` Lstats `repo/skills`, the real top-level directory, never `.agents/skills`) and remains accurate.

## Findings

None that meet the `--docs` lens bar. One item outside the lens's grep targets, carried forward from the round-1 report since it is still true:

- **[informational, not a docs-lens finding]** `.agents/docs/bugs/BUG-036-symlink-refusal-incomplete.md:3,20` — `status: open` and the `## Fix` line still read "Not fixed", but this diff now implements the remediation for both the leaf-file gap the record names and the directory-level gap the round-1 code review (`.agents/docs/reviews/2026-09-16-review-bug-036-code.md`) found on top of it. Bug records are excluded from the `--docs` grep targets, so this is not a finding, but the record will read stale once this fix lands — flagging for whoever closes BUG-036.

## Checked and clean

- `README.md` ("Adding things by hand" table, `npx skills add` snippet at :75, `sh-agent adopt`/`-g` note at :80), `skills/find-skills/SKILL.md:101,106`, `skills/VENDORED.md:38`, `AGENTS.md:23`, `context/AGENTS.md:37`, `REQUIREMENTS.md` R5 and the "Non-goals" symlink line — all describe `.agents/skills` resolving into `skills/`, or `adopt` copying a globally-installed skill into `skills/<name>`; none claims a link at any of the six managed directories or at a file directly inside the five flat ones, so the widened refusal contradicts nothing.
- `skills/init-project/SKILL.md` Step 5's symlink discussion is about `.claude/{skills,agents,commands,rules,workflows}` links in a *consuming project*, a different mechanism (that skill's own prompt logic) untouched by this diff's Go code.
- The newest `changes/` entry (`2026-09-16-bug-037-adopt.md`) is for a different bug and does not describe this fix's scope; no staleness introduced there.
- `.agents/docs/specs/` is empty; `.agents/docs/deploy/{local,remote,gcp}.md` and `.agents/docs/testing/sh-agent/guide.md` contain no symlink-scoped claims to check against the widened behaviour.
- `go build ./...`, `gofmt -l .`, `go vet ./...` all clean against the current working tree.

## Criteria

1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined

- Code-lens criteria (correctness of `checkManagedTree`, test soundness) — deliberately out of scope for a `--docs` pass; a separate code-lens review of this same widened diff was not re-run by me in this session (round-1 code review is `.agents/docs/reviews/2026-09-16-review-bug-036-code.md`, against the version this diff supersedes).
- `argolis-devbox`'s `verify_mirror`/`032_setup-harness.sh` gate, named as a remediation item in BUG-036 — a different repository, not reachable from here.
- Sandbox CLI starts and real-home smoke checks (guide §4/§5) — not re-run this session; not part of the docs lens.
