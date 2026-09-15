---
owner: reviewer
date: 2026-09-15
---
# Review — sweep-repo-docs   (2026-09-15)

Range reviewed: the named files as they are in the working tree, read whole (not a diff): `README.md`, `AGENTS.md`, `.agents/docs/testing/sh-agent/guide.md`, `.agents/docs/deploy/{gcp,local,remote}.md`. Lens: `--docs --project`. Uncommitted working-tree edits to README.md, AGENTS.md and guide.md (count removals, the new `/review --docs` sentence) were read as part of the tree and their old wording re-grepped.

## Verdict: REQUEST CHANGES

The worst problem: `README.md:32` and `.agents/docs/deploy/local.md:27` say `--vendor` defaults to `all`. The code (`cmd/sh-agent/main.go:216-224`, and the help text at `main.go:48`) defaults to the vendor recorded in the manifest when something is installed, and to `all` only otherwise. Take a home installed with `--vendor antigravity`, like the Mac (`.agents/docs/testing/2026-09-04-verify-mac-real-home.md:11`). There, a bare `sh-agent install` keeps writing Antigravity only, even though both guides say it writes both vendors. Other stale statements: `remote.md:37` sends the reader to an empty `TODO.md`; `README.md:57` lists a narrower '#' lint rule than the code enforces and than `AGENTS.md:38` states; the README command table leaves out `--dry-run` and `adopt --force`; README and remote.md say the checkout search looks only in the working directory, but it also searches the parent directories. `AGENTS.md`, the testing guide and `gcp.md` match the code.

## Classification (living docs in scope)

| File | Class | Proof |
|---|---|---|
| `README.md` | outdated | `:32` says `all` [default], but `main.go:216-224` uses the manifest's vendor first. `:57` leaves out "at the start", which `frontmatter.go:113` rejects too. The `:22-30` table has no `--dry-run` or `--force` (`main.go:35,38,40`). |
| `AGENTS.md` | accurate | `:35` matches `main.go:127-147` (IsCheckout gate, go-missing warning) and the `binary` line at `main.go:254-262`. `:38` matches `frontmatter.go:113`. Every path in `:46` and `:52-56` exists. |
| `.agents/docs/testing/sh-agent/guide.md` | accurate | `:56` and `:65` match `main.go:145` and `main.go:131`. `:59` matches `main.go:234`. `:23` has its budget from ADR-018 (`DECISIONS.md:147`), and render prints 93 lines. `:39`: every `agents/*.md` has exactly one `## Records`. `:7` cites `agents/verifier.md:50`. |
| `.agents/docs/deploy/local.md` | outdated | `:27` says `all` (default), which contradicts `main.go:216-224`. Everything else matches: `:8` R6 (`REQUIREMENTS.md:25`), the `:27` go-missing warning (`main.go:147`), the `:33` log grep (BUG-004), and `:37` ADR-010. |
| `.agents/docs/deploy/remote.md` | outdated | `:37` points to `TODO.md`, which is headers only (`TODO.md:1-5`); the Mac Antigravity check passed on 2026-09-04 (`2026-09-04-verify-mac-real-home.md:9`). `:29` leaves out "and its parents" (`repo.go:15,31-34`). `:33` mirror paths exist on this box. |
| `.agents/docs/deploy/gcp.md` | accurate | The `:10` template exists at `skills/init-project/assets/agents-tree/docs/deploy/gcp.md`. |

None is redundant as a whole. `README.md:94-99` repeats `AGENTS.md:11-14` with the same meaning.

## Findings

- **[medium / confidence: high] outdated** `README.md:32` — "`--vendor VENDOR` (target `all` [default], …)". Since `main.go:216-224`, when `--vendor` is omitted the code uses the vendor recorded in `~/.config/sh-agent/state.json`, and `all` only when nothing is installed (the help text at `main.go:48` says so). Should say: "`--vendor VENDOR` (`all`, `antigravity` / `agy`, or `claude`; default: the vendor of the existing install, else `all`)". The same applies to the next finding.
- **[medium / confidence: high] outdated** `.agents/docs/deploy/local.md:27` — "`all` (default)". `:19` says `install` "writes ~/.claude and ~/.gemini/config", but a home first installed with `--vendor antigravity` stays Antigravity-only. Should say what `main.go:48` says, and tell the reader to pass `--vendor all` to widen an existing single-vendor install.
- **[low / confidence: high] outdated** `.agents/docs/deploy/remote.md:37` — "On the Mac, additionally confirm the Antigravity paths (`TODO.md`, high priority)". `TODO.md` is headers only, and the Mac Antigravity paths passed on 2026-09-04 (`.agents/docs/testing/2026-09-04-verify-mac-real-home.md:9,17-20`). The Mac items still open are the operator's Claude Code checks (`.agents/docs/plans/2026-09-03-rebuild.md:99,113`). Should say: "Same as `local.md`." Any Mac-only step that is still pending belongs in the newest handoff's Next Steps, not in `TODO.md`.
- **[low / confidence: high] divergent / narrowed** `README.md:57` vs `AGENTS.md:38` — README says the installer refuses "an unquoted `#` after a space or tab". AGENTS.md (correctly) says "at the start or after a space or tab", and `internal/harness/frontmatter.go:113` rejects a leading `#` too. A reader going by README would expect `description: #foo` to install. Should say: "an unquoted `#` at the start of the value or after a space or tab".
- **[low / confidence: high] outdated (interface not documented)** `README.md:24,26,29` — the command table says nothing about `install --dry-run`, `uninstall --dry-run` or `adopt --force` (`main.go:35,38,40`). The testing guide's plain-copy tier depends on `--dry-run` (`guide.md:64-65`). Should add them to the matching rows.
- **[low / confidence: high] outdated** `README.md:32` and `.agents/docs/deploy/remote.md:29` — "finds the checkout from `$SH_AGENT_REPO`, the working directory, or the path it recorded". `internal/harness/repo.go:15,31-34` walks the working directory *and its parents*. Should say "the working directory or any parent".

### Unverified

- `README.md:77` — `git add skills/find-skills skills-lock.json`. `skills-lock.json` does not exist in the tree and never did (`git log --all -- skills-lock.json` is empty), and it is not gitignored. If `npx skills add` does not create it at the repo root, `git add` fails on the pathspec and the `&&` chain never commits. Whether the skills CLI writes that file could not be checked offline.
- `.agents/docs/deploy/local.md:37` — "Antigravity prompts the same way for paths outside the workspace". No log in scope records it.
- `.agents/docs/deploy/remote.md:33` — both mirror directories exist locally (`~/workspace/styrman-crew/image/harness`, `~/workspace/argolis-devbox/deploy/harness`). Their `VENDORED.md` refresh recipes were not read.

### Old wording re-grepped (working-tree edits)

- The removed counts "27 skills", "seven subagents" and "27 skills, 7 agents parse" have no remaining hits in `README.md`, `AGENTS.md`, `guide.md`, `deploy/`, `context/`, `agents/`, or the `skills/` directories that are not vendored. No stale finding.
- The old guide caveat "counts change only when a skill/agent is added or removed" was removed along with the fixed counts it qualified. The new `guide.md:22` still has the list check, and `list` prints the counts in its section headers. Nothing was dropped by meaning.
- The new `AGENTS.md:46` sentence (`/review --docs` for `context/`, `agents/`, `skills/`) adds to the global DoD rule (interface, command or deploy-step changes) without contradicting it. It matches the review skill's "plus any paths the project's AGENTS.md adds". Not divergent.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Re-check commands

```bash
git grep -nE '`all` \[default\]|`all` \(default\)' -- README.md .agents/docs/deploy
git grep -nE 'the working directory, or the path' -- README.md .agents/docs/deploy
git grep -n 'TODO.md`, high priority' -- .agents/docs/deploy
git grep -n 'unquoted `#` after a space or tab' -- README.md
git grep -nE 'dry-run|--force' -- README.md
git grep -n 'skills-lock.json' -- README.md
git grep -niE '27 skills|seven subagents|27 skills, 7 agents' -- README.md AGENTS.md .agents/docs/testing/sh-agent/guide.md .agents/docs/deploy context agents skills
```

## Not examined

- Every named file was read whole. `.agents/docs/specs/` is empty. The newest `changes/` entry (`2026-09-15-research-beads-elephant-goldfish.md`) was not read. The mirrors' own `VENDORED.md` files and the behaviour of the `npx skills` CLI were not checked.
