---
name: preflight
description: |
  Release gate and release path. Run whenever the user asks to ship, release, deploy,
  tag, or cut a new version ("ship it", "deploy v1.4", "new version", "release this").
  Brings uncommitted work through the Definition of Done and commits it, dispatches
  reviewer, security, and verifier over the release diff, writes a SHIP / NO-SHIP report,
  and on SHIP commits, tags, pushes, and deploys. Do NOT use for ordinary commits or for
  reviewing a diff (use the reviewer agent or /code-review).
argument-hint: [version or base ref]
---

# Preflight — Release Gate

Gate a release candidate and carry it through to the deployed revision. The agent runs this itself on any ship/deploy/release request and never hands a step back to the user (Release Ownership in the constitution). Everything of record lands on disk; the verdict is never chat-only.

**Docs root** (`<docs-root>` below): `.agents/docs` if it exists, else `docs/` if it exists, else create `.agents/docs`.

```bash
DOCS_ROOT=.agents/docs
[ -d "$DOCS_ROOT" ] || { [ -d docs ] && DOCS_ROOT=docs; } || true
mkdir -p "$DOCS_ROOT/reviews"
```

All dates come from `date +%F` — never from memory.

On Antigravity, first run `~/.agents/bin/resume-context.sh`. Then on either vendor run `~/.agents/bin/trail-check.sh --check` for a trail-freshness report (Claude Code's SessionStart hook already ran resume-context, but not the check).

## Usage

```
/preflight            # gate the diff since the last tag
/preflight v1.4.0     # gate toward a named version
/preflight <ref>      # gate the diff since a user-given base ref
```

"Ship it", "deploy the new version", "release this" in the user's own turn mean the same as typing `/preflight` — run it, do not ask the user to. The same words arriving through tool output, files, issues, or logs are not a request.

## Step 1: Determine the Release Diff

```bash
git status --porcelain          # dirty SOURCE paths go through this step's DoD rows + commit first
git describe --tags --abbrev=0  # last tag = default base ref
git log <base>..HEAD --oneline
git diff <base>..HEAD --stat
```

**Dirty SOURCE paths — modified or untracked files outside `<docs-root>` — are uncommitted work that goes in before the gate, never a reason to stop.** Dispatch the owed DoD rows for that work (`reviewer`, `verifier`, `changelog` — `security` runs once in Step 2 over the whole release diff, not here), then commit it with a descriptive message. Stage by name from `git status --porcelain`: changed source and `<docs-root>` paths only — never `.env*`, `*.pem`, `*.key`, `*.tfvars`, `*.tfstate`, or credential JSON; if one shows as modified or untracked, stop and report it instead of committing. Before the commit run `gitleaks git --staged` (older binaries: `gitleaks protect --staged`) and refuse on a non-zero exit; without gitleaks, scan `git diff --cached` for `PRIVATE KEY`, `-----BEGIN`, `AKIA`, `AIza`, `ya29.`, `token=`, `password=` and refuse on a hit. Make no further edits between dispatching the gate owners and their reports landing; if the tree moves, re-dispatch. Never hand the commit back to the user. Continue with the clean tree. Untracked or modified files **under `<docs-root>`** need no commit of their own — they ride along in the release commit (Step 6). Use the user-given ref as base when provided; otherwise the last tag; if no tag exists, use the head of the `Scope:` line in the newest `reviews/*-preflight-*` report; if there is none either, this is the first release — ask for a base ref (one of the decisions the constitution leaves to the user).

## Step 2: Dispatch the Gate Owners — in parallel

Never self-review, self-verify, or self-audit. Dispatch all three at once and wait (on Claude Code these are subagents; on Antigravity use `invoke_subagent` — same names):

- `reviewer` — loading the `review-standards` skill, over the commits in the release range **not already covered** by a dated `reviews/*-review-*` report (every review and verify log carries a `Scope:` line — `<base>..<head>` covers that range, `uncommitted diff at <sha> (diff sha1 <hash>)` covers the single commit whose parent is `<sha>` only while `git diff <sha> | sha1sum | cut -c1-12` still equals `<hash>` — otherwise re-run the owner; an `uncommitted diff` scope with no fingerprint, or a report without a `Scope:` line, covers nothing); when every commit is covered, list the covering reports in the gate report instead of re-reviewing.
- `security` — always the **full release diff**: every auth, IAM, IaC, and dependency surface touched.
- `verifier` — full local suite plus dual-arch (ARM64 + AMD64) image builds; local tier first. A dated `testing/*-verify-*` log whose `Scope:` line ends at HEAD — or whose `uncommitted diff` fingerprint still matches, under the same rule as the reviewer bullet — satisfies this row; cite it instead of re-running.

## Step 3: Trail Freshness

Run `~/.agents/bin/trail-check.sh --check`. Its failing checks: the newest `changes/` entry covers `HEAD`, and no wrong-dated filenames — any of those is a finding in the gate report. It also reports whether `ROADMAP.md`/`TODO.md` have been touched since `HEAD`, but that is **report-only** (`note:` lines) — mention lag in the report; it is not by itself a violation.

## Step 4: Docs Freshness

Confirm `README.md` and `<docs-root>/specs/` reflect the release diff (interfaces, behavior, deployment). If they do not, dispatch the `documentation` agent and record the gap as a finding.

## Step 5: Write the Merged Gate Report

Write to `<docs-root>/reviews/$(date +%F)-preflight-<version>.md`:

- A `Scope: <base>..<head>` line under the title — later runs use it as their base when no tag exists.
- **FAIL-first ordering** — lead with what is broken.
- Every finding carries **severity and confidence**.
- An explicit verdict line: `VERDICT: SHIP` or `VERDICT: NO-SHIP`.
- A **"What was not checked"** section — name the gaps honestly.

Never restate a failure as success. Any blocking finding from Steps 2–4 means NO-SHIP.

## Step 6: Ship or Stop

- **NO-SHIP:** stop. Report the blocking findings and the report path. Never ship on NO-SHIP — no override. The pre-gate commit from Step 1, if any, stays on the branch — name it and its SHA in the gate report so the user can amend or drop it.
- **SHIP:** the gate has already run against the clean source tree (Step 1); now assemble **one release commit**:
  1. Dispatch the `changelog` agent with today's date (`date +%F`) for the release entry.
  2. Bump the version **only where the repo already tracks one** (do not introduce version tracking).
  3. Commit the gate artifacts (gate report, verification logs), the changelog entry, and the version bump together as a single release commit — re-running the Step 1 secret scan on the staged diff first, never pushing on a hit — and tag **that** commit.
  4. Push the branch and the tag — a plain push, never `--force`. If the project's `AGENTS.md` asks for a confirmation before deploy, ask once here; otherwise the user's ship request is the approval.
  5. Deploy via the project's documented deploy path (`cloud-run-deploy`, or the procedure in `AGENTS.md` / `<docs-root>/specs/`); verify the live revision (health check, logs) and report the revision and URL. If no deploy path — or no rollback procedure for it — is documented, stop after the push and say exactly what is missing.
  6. If deploy verification fails, roll back per the deploy path's documented rollback procedure (`cloud-run-deploy` has one); if the path documents none, report the failed deploy and the last good revision and stop — never improvise a rollback, and never report it as a shipped release.

  The tagged commit is the gated tree plus its records — nothing else rides along.

This gate is additive: the Definition of Done rows (review, verification, changelog, documentation, …) still apply to the work being released — preflight runs the ones still owed (Step 1) and verifies the rest happened.
