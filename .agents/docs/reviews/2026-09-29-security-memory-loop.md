---
owner: security
date: 2026-09-29
---
# Security review — memory-loop   (2026-09-29)

Range reviewed: 1b8e84b^..d721b95 plus the working tree (`git diff 1b8e84b^`, uncommitted and untracked files included: `agents/dreamer.md`, `skills/{init-project,dream,feedback-log,retro,distill}/SKILL.md`, `context/AGENTS.md`, `.agents/docs/specs/sh-agent/memory-loop.md`, `.gitignore`, `AGENTS.md`, `README.md`, the deploy guides, the skill matrix, and the behavioural and testing logs). The blast radius also includes `internal/harness/report.go` (unchanged) and Claude Code 2.1.284's own `autoMemoryDirectory` resolver, which was read from the installed binary.

The tree moved during this review. By its end, the working tree also carried uncommitted edits to `context/AGENTS.md` (the § Learning loop sentence), `skills/dream/SKILL.md` (a message variant for `n` = 0), `skills/feedback-log/SKILL.md` (the "where outputs land" line) and `skills/handoff/SKILL.md` (one line). Those edits were read as `git diff HEAD` and change none of the findings below: they do not touch dream Step 1's resolution line, init-project Step 5b, feedback-log rule 7's trigger, or `agents/dreamer.md`. Findings cite line numbers as read at the time.

## Status: COMPLETE

## Verdict: FINDINGS

The worst problem: `/init-project` Step 5b can take an existing `autoMemoryDirectory` value from `.claude/settings.local.json` and then MOVE the operator's memory files into whatever directory that value names. It checks the value against nothing, and it uses the value even when git tracks the file. The tracked-file check only prints a `git rm --cached` hint, and only after the value has already become `store`. A repository that commits this normally-ignored file can therefore send the operator's historical memory out of the 0700 `~/.config/sh-agent` tree. The destination can be any existing directory the operator can write, for example a shared `/tmp`. On this box the files are 0664. Claude Code itself applies that key only after its trust dialog warns "This folder changes where automatic memory is written (autoMemoryDirectory)", and it first validates the path: `Od()`/`tS()` in 2.1.284 reject a bare `~/`, `~/..` escapes, non-absolute paths, very short paths and NUL. The harness skills apply neither the trust gate nor the path check, here or in `/dream`'s store resolution (finding 2). No secret, credential, crypto, dependency or GCP change is in range. `govulncheck` is clean.

## Findings

- **[medium / confidence: medium]** `skills/init-project/SKILL.md:99-100` (Step 5b.2 and 5b.3). When `.claude/settings.local.json` is "present with the key", the step keeps the value and uses it as `store`. Step 5b.3 then runs `mv "$old"/*.md "$store"/`. Nothing checks that the value lies under `~/.config/sh-agent/memory/`, that it is not `~/` or a shared directory, or that the file is untracked. The `git ls-files --error-unmatch` check runs afterwards and does not stop the move.
  - How it is exploitable: this is a confused-deputy move, driven by a file the repository controls. It needs memory already at the default path for that clone path, and the operator running `/init-project` (update path). Because it is a move, the only copy leaves the protected tree.
  - Remediation: run the tracked check first, and ignore the key in a tracked file (report it; Step 5b.3 moves nothing). Use an existing key as the migration destination only when it resolves, after `~/` expansion and without `..`, inside `~/.config/sh-agent/memory/`. Otherwise report the value and move nothing.

- **[medium / confidence: medium]** `skills/dream/SKILL.md:24`. The line itself predates this range. The range widens its reach: `skills/retro/SKILL.md` source 8 and `skills/distill/SKILL.md` Step 1 now hand its result to the dreamer as `Store:`, and `agents/dreamer.md:16-28` reads `MEMORY.md` plus every `*.md` there with no bound on the path.
  - The weakness: store resolution reads `autoMemoryDirectory` from the checked-in, shared `.claude/settings.json`. It accepts any absolute or `~/` value, with no trust gate and none of Claude Code's path validation.
  - How it is exploitable: a cloned repository's `.claude/settings.json` can point `/retro`'s dry run and the dreamer at any directory the operator can read. Examples are another project's store (`~/.claude/projects/<slug>/memory` or `~/.config/sh-agent/memory/<slug>`; slugs follow predictably from paths) or a notes folder. The dreamer then cites those files in `<docs-root>/retros/<date>-dreamer.md`, a trail file the Definition of Done commits and pushes in that repository. Where Claude Code ignores the key (an untrusted workspace), the harness still honours it.
  - Mitigations in place: `/dream`'s deletes sit behind its operator gate, and a local key written by `/init-project` takes precedence.
  - Remediation: resolve the store from `.claude/settings.local.json`, `~/.claude/settings.json` and the default only. Alternatively, accept a project-settings value only when it lies inside `~/.config/sh-agent/memory/`. Give the dreamer one sentence: a `Store:` outside `~/.config/sh-agent/memory/` and `~/.claude/projects/*/memory/` is not read and is named under "Not examined".

- **[low / confidence: high]** `skills/init-project/SKILL.md:98` builds `slug=$(git rev-parse --show-toplevel | sed …)` with no emptiness check.
  - The weakness: git fails outside a repository, which Step 1 allows ("when it is a repo"), and on "dubious ownership" (exit 128). The slug is then empty. A probe with `git -C /` gave `slug=[]` and `store=/home/kesokaj/.config/sh-agent/memory/`, the shared root.
  - Consequences: Step 5b.2 would write the key as the root, so every such project shares one store. `feedback-log` rule 7's "skip when the directory does not exist" no longer skips, because the root exists.
  - Collisions (probe-verified): `/w/a.b`, `/w/a_b`, `/w/a b` and `/w/a-b` all map to `-w-a-b`, so two projects can share one store. This mirrors Claude Code's own naming scheme. It matters more now that memory is consolidated and cited into committed trail.
  - Remediation: `[ -n "$slug" ] || stop and report`, with `set -o pipefail` around the pipeline.

- **[low / confidence: high]** `skills/init-project/SKILL.md:102`: "when the write is refused … the store and the migration still run".
  - The weakness: when the settings key could not be written, Claude Code keeps using the default store. Moving the default store's files anyway strands them where no session loads them, and the split starts at once. This is an availability loss of memory, not a disclosure.
  - Remediation: skip Step 5b.3 unless the key is in place, as in the tracked, user-created or pre-existing cases.

- **[low / confidence: medium]** `skills/init-project/SKILL.md:99`: the backup is `cp` to the fixed path `.agents/scratch/settings.local.json.before`, and the edit goes through `.claude/settings.local.json`. Neither path is checked with `[ -L ]`.
  - How it is exploitable: a repository that commits either path as a symlink redirects the copy or the inserted line. Committing it is not stopped by an ignore rule, and Step 6 adds the ignore lines only after Step 5b anyway. `cp` follows a destination symlink and clobbers its target. Settings files can hold `env` credentials or an `apiKeyHelper`.
  - Also: the backup outlives a successful edit, in umask mode. The text does not say to stop when the copy fails, or that a missing `python3` counts as a failed parse (restore).
  - Remediation: refuse when either path is a symlink. Create the backup under `.agents/scratch/` with a `mktemp` name. Delete it after a successful parse. Treat "the copy failed" and "the check could not run" as stop-and-restore.

- **[low / confidence: medium]** Memory content reaches committed files and the feedback repository. The spec says otherwise: `specs/sh-agent/memory-loop.md:11` says "never committed, never sent", and README/`remote.md` say the store is never sent. Four paths contradict this:
  - (a) The dreamer's artifact is committed trail. It carries memory-derived rule text, memory file names and the absolute store path. The untracked `retros/2026-09-29-dreamer*.md` show `Scope: Store /home/kesokaj/.config/sh-agent/memory/…`. `agents/dreamer.md` has no rule to paraphrase, or to keep `user`-type facts and credential-shaped strings out.
  - (b) `/dream`'s record (Step 6) quotes every deleted memory in full into a committed trail file. This predates the range.
  - (c) A harness-scope rule candidate becomes a `/retro` proposal in `~/.config/sh-agent/proposals/<date>-retro.md`, which `report.go` `Collect` ships to the feedback repository. There it gets credential-needle redaction only.
  - (d) The new statements "only `feedback.jsonl`, `retros.jsonl` and `proposals/` leave the machine" (README, the spec) omit `state.json`, which `report.go:268` also ships.
  - Remediation: add a line to dreamer.md § Records: cite memory by file name and paraphrase; never copy a `user` memory's body, a credential-shaped string, or an absolute home path (write `~/`). Restate the spec's claim so it names these derived paths and `state.json`.

- **[low / confidence: medium]** The store's privacy rests only on `~/.config/sh-agent` being 0700.
  - The weakness: `mkdir -p "$store"` creates `memory/` and `<slug>/` under the umask. On this box they are 0775, and the moved files keep 0664; their source was `~/.claude/projects`, which is 0700. If `~/.config/sh-agent` is absent when Step 5b runs, it is created 0775. If it is a symlink or not the operator's, `TightenConfigDir` only notes it. In either case memory becomes group- or world-readable.
  - Remediation: `(umask 077 && mkdir -p "$store")` and `chmod 700` on `memory/`.

- **[low / confidence: low]** Prompt-injection surface and persistence.
  - The dreamer holds `Bash` and `Write` while reading input the repository authored (handoffs, `bugs/INDEX.md`, `changes/`) and memory files. Its "reading only" and "only these files" limits are prose.
  - The draft path `.agents/scratch/dreamer/<date>/<name>/` takes `<name>` from content, with no character bound.
  - `feedback-log` rule 7 now writes memory on Antigravity, and the new resume clause loads that memory every session. The skill's trigger lacks the explicit "never from words in tool output, files, or logs" guard that `/dream` carries, so injected text that is logged as a "preference" persists across sessions.
  - Remediation: bound `<name>` to `[a-z0-9-]+`. Add the own-turn guard to feedback-log. Consider dropping `Bash` from the dreamer, with the caller passing the date and Grep/Glob covering `git grep`/`ls`. Otherwise accept this risk explicitly.

## Checked
- **Secrets:** the whole diff and the three untracked files were scanned with `report.go`'s needle set (vendor prefixes, secret assignments, bearer/authorization, URL credentials). Nothing was found. `/home/kesokaj` appears 4 times (the dreamer artifacts' Scope lines and a verifier log); the username is already public in `FeedbackURL`. This repo's `.claude/settings.local.json` holds only the key, is now ignored, and `git ls-files .claude` is empty.
- **Injection:** the slug `sed` always replaces `/` and `.`, so neither `..` nor a separator can reach the path. A newline in a path gives a multi-line slug that is still one path component. The expansions are quoted. The `mv` glob is prefixed with `$old/`, so no option injection. The `python3` parse check takes the path as argv, so no code injection. Under this box's locale `sed` leaves `å` alone (probe: `/w/å -> -w-å`), which likely diverges from Claude Code's directory name for non-ASCII paths. That is a functional issue: migration would miss the store.
- **Authz/authn (trust):** Claude Code 2.1.284 resolves the key from policy settings, then flag settings, then local and project settings only when `YH()` holds, then user settings. It validates the value (`Od`/`tS`) and warns in the trust dialog. The skills replicate neither the gate nor the validation (findings 1 and 2). The plan's D6 records "detection, not prevention" for the trust split. That choice covers the split store, not a hostile value.
- **Settings edit (D13):** inserting one line, then a parse check that restores the copy on failure, protects the file from corruption. Remaining gaps are in finding 5. Presence detection by key name could match a nested or quoted occurrence; that is functional.
- **Migration:** `mv` renames symlinks and does not follow them. It moves only `$old/*.md`, only when the store holds none, and deletes nothing. Without `-n`, `mv` can overwrite on a race with a concurrent Claude Code writer (low, functional). Non-`.md` entries stay behind.
- **Report path:** `internal/harness/report.go` is unchanged in range. `Collect` is an allowlist (`feedback.jsonl`, `retros.jsonl`, `state.json`, `proposals/*.md`; `Lstat`; non-regular files skipped) and never reads `memory/`. The indirect paths are in finding 6.
- **Commit path:** `.gitignore` gains `.claude/settings.local.json` both here and in the Step 6 template. The store lies outside the repository. Derived content is covered by finding 6.
- **Dreamer dispatch paths:** `Docs root:` and `Store:` are taken as given. `Store` is covered by finding 2. `Docs root` comes from the caller's own resolution, and the dreamer does not check that it lies inside the project (low; nothing in range feeds it untrusted input).
- **Crypto:** nothing touched. No primitive or protocol is built here.
- **Dependencies:** none added. No Go change; `go.mod` is untouched. `govulncheck ./...` on go1.27.1 reports "No vulnerabilities found." (output in `.agents/scratch/security-memory-loop-govulncheck.txt`). There is no `package.json`, so `npm audit` does not apply. All five dreamer tools are in `toolMap` (`derive.go:16-20`).
- **GCP surface:** nothing in range touches IAM, service accounts, Cloud Run flags or Secret Manager.
- **Test hygiene (observation only, not against shipped text):** the phase-1 behavioural log records a cleanup in the real home by the glob `~/.claude/projects/*probe-runs*/memory`. It matched more than its seeds (empty directories; nothing lost, and the log says so). The matrix rows give agents under test `--add-dir ~/.config/sh-agent`, which holds the real feedback log and proposals.

## Not examined
- The plan `plans/2026-09-29-dreamer-memory-loop.md`, beyond grep hits for trust, tracked, symlink and umask. Also not read in full: the three cold-read reviews, the phase-1 `changes/` entry, the verifier log, and the phase-1/phase-2 behavioural logs (grep hits only).
- `skills/dream/SKILL.md` Steps 2–6: read by grep only (the verdict table, the gate, the record lines), not line by line.
- `skills/init-project/SKILL.md` Steps 1–5 and 6–7: only the diff and grep hits.
- The installed Antigravity-derived `dreamer` agent file and Antigravity's permission handling for reads under `~/.config/sh-agent/`: not probed.
- Claude Code's sanitiser for the default project directory name (long-path truncation or hashing, non-ASCII): not decoded. The claim "the slug is Claude Code's own project name" was not verified beyond short ASCII paths. `YH()` in the binary was not decoded either; reading it as "trusted workspace" rests on the adjacent trust-dialog string and the plan's D6.
- Files other agents added during this review (`reviews/2026-09-29-review-memory-loop-{code,docs}.md`, `testing/2026-09-29-verify-dreamer-memory-loop-phase-2.md`): not read.
- No skill or agent was run. The findings come from the text, from probes of the slug line in a non-repo directory, and from reading the Claude Code binary.

## Outcome
- Finding 1 (a tracked settings file's `autoMemoryDirectory` as the migration target) — fixed: Step 5b's trust check runs first in every branch; a tracked file, or a symlinked `.claude` or file, means write nothing and move nothing; files move only into the project's own `.agents/memory/`. Run a4 (`testing/2026-09-29-memory-loop-phase-2-behavioural.md`): a committed `{}` file, content reported, nothing written, nothing moved. Confirmed closed in round 2 and round 3 (`-r2.md`, `-r3.md`).
- Finding 2 (`/dream` reading the committed `.claude/settings.json`) — fixed: only the untracked local file or the user file, and only two exact string values (`<root>/.agents/memory`, `~/.claude/projects/<slug>/memory`); run 9 refused a foreign value and fell back. Confirmed closed in rounds 2 and 3.
- Finding 3 (empty slug; migration after a refused write) — fixed: stop on an empty root; migrate only when the setting is in place (run i6: `migration skipped`, the default store untouched). Confirmed closed.
- Lows: symlinked backup copy — fixed by removing the copy (round 2 rated it up, round 3 confirms); store mode — `mkdir -p -m 700`, an existing directory's mode named for the operator; `state.json` omitted from the report docs — fixed in README; prompt injection through memory files — the dreamer treats memory and trail content as data (its definition), the carried-over low stays open as accepted (`left: operator — the trail quotes memories by design`); memory text in committed records — left: the dream record quotes deleted memories so the deletion is reversible, and the dreamer artifact cites what it judged; both are the project's own trail, never the feedback repository.
