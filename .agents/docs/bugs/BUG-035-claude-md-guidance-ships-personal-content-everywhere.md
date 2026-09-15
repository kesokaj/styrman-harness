---
owner: bugs
date: 2026-09-16
id: BUG-035
status: open
area: installer/status (guidance text)
---
# BUG-035: installer and status advise moving a hand-written CLAUDE.md into rules/, which ships everywhere

## Symptom
Security review of the a81ef4b mirror refresh (`styrman-crew/.agents/docs/reviews/2026-09-15-security-worktree.md`, finding 5) found the guidance at `internal/harness/apply.go:37` and `cmd/sh-agent/main.go:309` (and `README.md:53`) tells the owner of a hand-written `~/.claude/CLAUDE.md` to "move its content into a rule file under rules/ in the harness repo, remove the file, and run install again" with no warning that `rules/` is deployed to both vendors on every machine, pushed, and mirrored into every downstream consumer (styrman-crew's runtime image, reaching every tenant VM).

## Repro
Not reproduced by exploit; read from the message text at `internal/harness/apply.go:37`. A user follows the printed instruction literally, puts personal content (hostnames, tokens, private context) into `rules/`, and it ships to every machine that installs the harness.

## Root cause
The guidance in BUG-032's own change (c514936, fe21671) tells the reader where a foreign file's content can go so `sh-agent status` stops flagging it, but never distinguishes "rules meant for every machine" from "personal content that happens to be in a CLAUDE.md" — the two look identical to the code that detects the foreign file.

## Fix
Not fixed. Remediation named by the finding: the guidance in both messages and `README.md:53` separates personal content (keep it outside the harness — a project's `AGENTS.md`, or an unmanaged file) from rules meant for every machine, and the covering tests (`TestStatusNamesForeignClaudeMD`, `internal/harness/harness_test.go:257`) change with the text. Not reachable today: the mirror's `rules/` holds only `.gitkeep`.

## Prevention
Guidance text that tells a user to move content into a harness-managed, deployed-everywhere path (`rules/`, `context/`, `skills/`) states what must never go there (credentials, personal/host-specific notes) in the same sentence, not as a separate caveat a reader can miss.

## Related
BUG-032 (the change that introduced this guidance, whose review found the coverage gap but not this content-scope gap); styrman-crew `2026-09-15-security-worktree.md` finding 5, merged in `2026-09-15-review-worktree-merged.md`.
