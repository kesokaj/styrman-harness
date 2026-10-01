---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-gate-lows-docs   (2026-10-01)

Range reviewed: 8b6c6de..working tree (`git diff HEAD`), docs lens

## Status: COMPLETE

## Verdict: REQUEST CHANGES

Docs lens only. Three living docs still say the old thing, and four behaviour changes have no home in the docs that list the same class of rule. The most concrete stale restatements are the unpinned `npx skills add ... -y` in `skills/VENDORED.md:57` (the README sample was fixed, this copy was not) and `git commit -am` in `.agents/docs/deploy/remote.md:72` (README copy fixed, this copy not). The three missing-home items below are lists of refusals that now omit one refusal each. No finding is high severity; none contradicts REQUIREMENTS.md.

## Findings

- **[medium / confidence: high]** `skills/VENDORED.md:57` — still says `npx skills add <owner/repo> -s <skill> -y`. The README diff (`README.md:251`) changed the same recipe to `npx skills@<version> add ... ` with no `-y`, and BUG-143's prevention rule names "unpinned `npx -y`" as a finding. Restatement left standing; make it match README § Adding things by hand. Same recipe as the README, so the two now disagree. (`skills/VENDORED.md:12-14` "Fetched with" cells record how three rows were fetched historically; left alone, not a finding.)
- **[low / confidence: high]** `.agents/docs/deploy/remote.md:72` — resync loop still shows `sh-agent install && git commit -am "context: tweak" && git push`; the diff changed the identical line in `README.md:205` to `git add context/dev-setup.md && git commit -m ...` (the global rule is stage by name, never `-a`). Two copies of one example now differ.
- **[medium / confidence: high]** dropped/unhomed: `allowed-tools` refusal (a) has no home in the lists of what the installer refuses. `AGENTS.md:58-59` (project rules: description, `#` values, 1,024-char limit) and `README.md:233` ("The installer refuses a `SKILL.md` without a `description`, and any frontmatter value ... read as a comment") enumerate refusals and omit `allowed-tools`; `.agents/docs/specs/sh-agent/delivery.md:77` lists "the installer refusing secret-shaped files, symlinked destinations, duplicate skill names and malformed agent descriptions", also without it. A skill author following README § Adding things by hand hits the refusal with no doc naming it. The only record is BUG-143 (INDEX Graduated row, line 156). Add one clause to `AGENTS.md` § Rules and `README.md:233`; delivery.md:77 is the "Built and shipped" list (rule 6's home for the claim, so a one-clause addition there is the owner's call).
- **[medium / confidence: high]** narrowed: memory-store refusals. `internal/memory/memory.go:245-248` now refuses a symlinked `MEMORY.md` (exit 2, via `refuse`) on `show`, on `add`, and in the 200-line precheck. The lists of refusals say only "`.agents` or the store is a symlink": `.agents/docs/specs/sh-agent/memory-loop.md:38`, `README.md:195` (`sh-agent memory add` row), `skills/project-memory/SKILL.md:52` ("a symlinked or git-tracked store"). Spec row at memory-loop.md:38 should name the index file; README:195 and project-memory:52 follow. Also no memory row in `testing/sh-agent-matrix.md` mentions a symlinked index (grep: none), so the sandbox matrix does not cover it; the Go test in `internal/memory/memory_test.go` does (read in diff, not run).
- **[low / confidence: medium]** dropped: dream feedback redaction (c). `.agents/docs/specs/sh-agent/dream.md:22` (Scrub) and `README.md:194` say inputs pass `harness.Redact` / "report's rules"; the diff makes feedback lines pass `harness.RedactJSONL` (the JSON-decoding pass `report` uses for `feedback.jsonl`), and `dream.md` § Hardening details (lines 41-46) has no bullet for it. The spec's wording "every input passes `harness.Redact`" is now inexact for feedback entries. Add one bullet to Hardening details and amend line 22. `RedactJSONL` is a new exported function in `internal/harness/report.go`; no doc lists it (internal API, not a finding by itself).
- **[low / confidence: medium]** `skills/data-gcp/SKILL.md:111` — the new text (a role `app` that is `cloudsqlsuperuser` "also carries CREATEROLE and CREATEDB") is not restated in `skills/golang/references/service.md:316`, which describes the same role as "`app` (a `cloudsqlsuperuser` member)" and points to data-gcp § Roles as "the one recipe", so it is not a contradiction. The sentence itself now reads awkwardly: "a `cloudsqlsuperuser` — ... — member with DDL" splits across lines 111-112. Unverified: the CREATEROLE/CREATEDB claim about Cloud SQL's `cloudsqlsuperuser` is a Google documentation fact I did not check (no network).

## Checked, no stale copy

- `golang:<major.minor>` without digest: `skills/docker/SKILL.md:130` (troubleshooting row, a tag in prose, fine), `skills/golang/SKILL.md:56` and `skills/operate-gcp/SKILL.md:93` say the tag tracks the patch and point at docker rule 1; neither shows a Dockerfile `FROM`. `skills/gke-basics/references/gke-app-onboarding.md:28` `FROM golang:1.22` is the noted gke-basics fork and `gke-basics/SKILL.md:23` tells the reader to use the docker skill's tag; it does not mention the digest now, acceptable (illustrative sample, flagged as such). `docker/SKILL.md` lines 67-69 and rule 1 agree with each other.
- Makefile `--project`: `git grep 'gcloud run deploy' skills` lines without `--project` are in `api-auth-gcp`, `cloud-run-deploy`, `data-gcp` etc.; these are multi-line commands, I read their continuation only in `cloud-run-deploy/SKILL.md:78,121`; not all lines checked, listed as unverified below.
- `commit -am`: no copy in `skills/init-project/assets/*` or any skill or context file (grep over `skills context agents README.md`).
- `npx -y`: no other copy in `skills/` outside `skills/find-skills` (vendored) and `skills/VENDORED.md:57` (finding above).
- feature issue text as data: `git grep 'gh issue view'` over skills/context/agents finds only `skills/feature/SKILL.md:19`; no second restatement.
- Graduation: `.agents/docs/bugs/INDEX.md:156` (Graduated row) and BUG-143 status agree with `inventory.go`/`TestLintRefusesAllowedTools`; BUG-145 row present.
- Stale tests: none found. `internal/harness/harness_test.go:2387-2396` asserts the new refusal; I did not search every fixture skill for `allowed-tools` beyond `git grep` (only that test hits).

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- Newest `changes/` entry: the newest is `2026-10-01-1340-timer-per-home-and-two-decisions.md` (before this scope); the working tree has no `changes/` entry for these lows (untracked: BUG-145, two testing logs). Expected at checkpoint; not a finding yet.
- Whether every `gcloud run deploy` multi-line snippet in skills carries `--project` (only the Makefile edit was in scope; the rest were not read).
- `.agents/docs/testing/guide.md` § 3e and `sh-agent-matrix.md` read by grep only (`symlink`, `feedback.jsonl`), not in full.
- `.agents/docs/deploy/local.md`, `gcp.md`; plans, retros, changes archives (historical, not living).
- Code correctness of the diff (not this lens); the Go tests were not run.
- The Cloud SQL role-attribute claim (external).
