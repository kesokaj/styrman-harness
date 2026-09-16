---
owner: bugs
date: 2026-09-16
id: BUG-050
status: fixed
area: internal/harness/plan.go (deploy walk)
---
# BUG-050: the install walk copied secret-shaped files into every vendor directory

## Symptom
`install` copies every file in a skill directory into both vendor
directories on every machine with no filter on what the file is. The
security review planned an install against a scratch repo holding
`skills/foo/creds.key` and got `write ~/.claude/skills/foo/creds.key` plus
the Antigravity mirror write. `.gitignore` lists exactly those shapes
(`*.key`, `*.pem`, `.env*`), so git never shows such a file and no review
would ever catch it from the diff — the one step that touches the file on
disk is the one that spreads it. Not triggered in a real deploy; found by
review.

## Repro
Place a file matching one of the secret shapes inside any
`skills/<name>/` directory and run `install`; pre-fix, the file is copied to
`~/.claude/skills/<name>/` and the Antigravity mirror.

## Root cause
The deploy walk in `internal/harness/plan.go` refused symlinks and
`.DS_Store` but never asked what a file IS — no check against the
secret-shaped-path list the context's Git rule already applies to commits.

## Fix
`secretShaped` helper (stdlib `filepath.Match`, `internal/harness/plan.go:207`)
refuses `.env`, `.env.*`, `*.pem`, `*.key`, `*.tfvars`, `*.tfstate`,
`*credentials*.json`, `*-key.json` inside a skill directory, naming the
path; called at `plan.go:120`. `TestPlanRefusesSecretShapedFiles` covers
seven shapes plus one control file that must still deploy.

## Prevention
The same secret shapes that stop a commit (`context/AGENTS.md` § Hard
rules, "Ask git what actually changed" bullet) stop a deploy too: any step
that copies project files outward applies that list, not just the commit
step.

## Related
BUG-044/045 are the same family (a secret- or debris-shaped path reaching
somewhere it shouldn't); with this record the "apply the commit's
secret-shape list to every outward-copying step" rule has recurred enough
to be worth a named check in `internal/harness` rather than three separate
fixes — it already graduated to `secretShaped` and its test here. Source:
`.agents/docs/reviews/2026-09-16-security-clear-the-board.md`.
