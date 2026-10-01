---
owner: main
date: 2026-10-01
---
# 2026-10-01-1815 — The v2.1.0 gate's left items closed; the `harness-audit` skill

Plan `plans/2026-10-01-gate-lows-and-harness-audit.md`, Phases 1 and 2. The operator, own turn: "I want you to fix everything that is left to do. Fix it or skip it. I trust your recommendations. Plus i want to add one more skill. Harness audit …".

Fixed (gate report `reviews/2026-10-01-preflight-v2.1.0.md`, rows L1 to L10):

- **L1, timer**: `sh-agent dream --remove-timer` prints `dream timer removed` or `no dream timer found for <home> on this system`; `uninstall` says it removed a timer only when it did. A unit that v2.0.0 or earlier wrote under the bare name in a home whose own name is now suffixed is never stopped or deleted by the binary: `--remove-timer`, `uninstall` (dry runs included) and the `dream` line of `status` name it, with the check and the steps (`dream.RemoveTimer` returns what it removed; `LegacyTimerFiles`, `LegacyTimerNotice`). Tests `TestRemoveTimerSaysWhatItRemoved`, `TestLegacyTimerFilesAreNamedNotRemoved`, `TestDreamTimerSaysWhatItFound`.
- **L5, lint**: a command or rule whose first line opens a frontmatter block that cannot be read is refused (`TestLintRefusesUnreadableFrontmatterInDocs`).
- **L6, search**: the header of a words search and of `--open` names the directories a `.sh-agent-noindex` marker leaves out (`TestHeaderNamesLeftOutDirectories`). In the header, not in `status`, which knows no project (plan D1, D11).
- **L2, L7 and the override residue, refresh script**: a leading `//` in `--dest` is the same place; `--work` is never `/`, the home directory or the checkout, never nested with `--dest`, and a non-empty one needs the script's marker (written before the first clone) or an `origins.txt`; a source override is an `https://` URL or an absolute path.
- **L3, L4, notes**: two notes reworded; `TestCatalogNotesNeverAllowQuiet` judges the clause that names `--quiet` and refuses a carve-out word in it.
- **L10, tests**: every subprocess of the catalog tests gets an environment without `GIT_*` variables (`GIT_DIR=/nonexistent go test …` failed before, passes now).

Skipped, with the reason: L8 (hostnames in the trail, the artifact's sharing: the operator's before any publication); L9 accepted on the operator's words; other kinds of note exception than `--quiet` stay with the security review of a refresh (D2); BUG-139 stays open, a measured rate of 5 of 6 (D3); the interactive check per CLI and `printf %q` on macOS (only the operator, or a Mac, can do them).

Bug records: BUG-150 (new: a red-first run of the `--work` guard's test cloned the fixtures into the real checkout root; five untracked paths, checked and removed by name; the cases now run a copy of the script in a stand-in checkout with sources that cannot be cloned); BUG-142 and BUG-140 graduated.

New skill `skills/harness-audit/` (`SKILL.md`, `references/areas.md`): the audit of 2026-09-30 and 2026-10-01 as a procedure. A baseline the main agent measures, one read-only `explorer` per area (C, WS, TS, AG, IN, SD, TR, G), findings with permanent IDs, a report, and an action plan in the `feature` skill's shape that is written and not started; `--reaudit` re-checks every ID with one `reviewer` per area. Model-invocable (ADR-014); its 236-byte description brings the listing to 19,485 of 19,500 bytes, so the next skill needs a description shortened elsewhere.

Docs (documentation agent, read against the diff): README (the `status`, `uninstall`, timer and `search` rows, the refusal sentence, the skills table), `specs/sh-agent/{dream,search}.md`, `deploy/local.md`, `testing/guide.md` § 2, § 3f, § 4, `testing/skill-matrix.md` (three rows). Main agent: project `AGENTS.md` § Rules, `skills/VENDORED.md`, the script's header.

Verification so far (main agent): `gofmt -l cmd internal` empty, `go vet ./...`, `go test -race -count=1 ./...` six packages ok; every new test seen red first. Next: the skill's behavioural runs on both CLIs, then the release gate's verifier, reviews and security read. Cold read of the plan: three rounds, the cap passed on the operator's same-turn words with six Open rows (plan § Phases, first box).

Agents dispatched: reviewer ×3 (cold reads), documentation ×1, bugs ×1.
