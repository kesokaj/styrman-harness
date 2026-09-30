---
owner: security
date: 2026-09-30
---
# Security review — f1ca525-e2cffc1   (2026-09-30)

Range reviewed: f1ca525..e2cffc1 (1 commit, 30 files). This is round 2, over the fixes for `reviews/2026-09-30-security-cb799ff-f1ca525.md`. Line numbers are pinned to e2cffc1 (`git show e2cffc1:<path>`).

Blast radius followed into:
- `docsRootOf` (`internal/search/reconcile.go:68`), `skipReason`, `mdFiles`/`readFile`;
- `memory.trusted`, `harness.ConfigDir` callers, `keepFailure`;
- `claude --help` and `agy --help` (1.2.14);
- the Antigravity step payloads on this box, through a read-only probe that printed counts only.

## Status: COMPLETE

## Verdict: FINDINGS

**Worst problem: the Antigravity reader still sends what the agent did with its tools.** The round 1 fix keeps step types 14 and 15, and 15 is not prose only. The model's reply step also holds its tool calls with their JSON arguments: the command lines it ran and the whole files it wrote.

What the probe found, counts only, over 88 conversation databases:
- The reader's own extraction keeps 1960 text runs from type-15 steps.
- 737 of those runs carry a `"CommandLine"`.
- 59 carry a `"CodeContent"`, the longest 15,691 bytes.
- 10 carry a `"ReplacementContent"`.

Command lines are exactly the shapes round 1 showed `harness.Redact` misses (`mysql -pSECRET`, `curl -u user:pass`). The Claude Code reader drops tool inputs, so the inconsistency round 1 named still holds, now for tool inputs rather than tool outputs.

Second, the new symlink guard checks the docs root and the record directory, but not `.agents`. A repository that commits `.agents` as a link still has the job read wherever it points.

The other fixes hold against their code and tests:
- the private agy HOME with its cleanup and sweep;
- the Claude Code flags;
- stdin stream-json;
- the flock locks;
- folding of one-line fields;
- the `~` paths;
- `%` escaping and the refusal of unsafe unit text.

Two Outcome lines claim more than was done:
- the disclosure half of round 1 finding 3;
- the PATH half of finding 10.

No hand-rolled crypto, no new dependency, and govulncheck is clean.

## Findings
- **[medium / confidence: high]** `internal/search/dream.go:223` (`WHERE step_type IN (14, 15)`), with `textRun` at `:92`: type-15 steps carry tool-call arguments.
  - **Measured** (read-only `immutable=1` opens, no content printed). Of 1536 type-15 steps:
    - 737 hold `"CommandLine"`;
    - 476 hold `"AbsolutePath"`;
    - 69 hold `"TargetFile"`;
    - 59 hold `"CodeContent"`;
    - 10 hold `"ReplacementContent"`.

    The reader's regex and three-space rule keep them. The numbers are in the Verdict.
  - **Why it matters:**
    - Inline credentials in commands, and whole config files the agent wrote (`.pgpass`, `.netrc`, docker `config.json`, `.env`), reach the prompt, guarded only by the deny-list.
    - Under vendor `claude` or `all` they go to Anthropic. Harness proposals can carry them on to the feedback repository.
    - The spec's "never tool output" (`specs/sh-agent/dream.md` step 1) is literally true but misleading.
    - The Claude Code reader sends text parts only.
    - `TestAgyConversationText` has no type-15 step holding a tool call, so it cannot see this.
  - **Remediation (stdlib):**
    - Read only the reply-text field of the type-15 message. Find its field number with one probe, then walk the protobuf wire format with `encoding/binary.Uvarint` (about 30 lines). `google.golang.org/protobuf` is not needed to read one string field.
    - Failing that, drop every run that holds a JSON argument object (`{"` followed by a key and `":`).
    - Add a fixture type-15 step with prose plus `{"CommandLine":"… -pS3cret …"}`: the prose must stay and the command must go.
    - Until then, the spec says the agent's commands and written files are sent.
- **[low / confidence: high]** `internal/dream/collect.go:116-124`: `linked` does not check `.agents`.
  - **What it checks:** `linked` Lstats `root/<docs>` and `root/<docs>/<dir>`. For the default `.agents/docs`, Lstat follows every path component but the last.
  - **The path through a committed link:** say a repository commits `.agents` as a symlink to another checkout's `.agents`, or to any directory with `docs/{bugs,reviews,…}/*.md`.
    - `docsRootOf` (`reconcile.go:68`, `os.Stat`) accepts it, and so does `skipReason` (`dream.go:106`, `os.Stat`).
    - `linked` sees two real directories.
    - `mdFiles` reads the target's regular `.md` files into the prompt.
  - **The mismatch:** `memory.trusted` (`memory.go:77-95`), the model the round 1 remediation named, checks both `.agents` and `.agents/memory`. So for such a checkout the memory is refused while the trail is read. `TestCollectDoesNotFollowLinkedTrail` links `reviews` only.
  - **Remediation:** Lstat every component from the root down, or require `filepath.EvalSymlinks(dir)` to lie under `EvalSymlinks(root)`. Add the test with `.agents` linked.
- **[low / confidence: high]** Outcome 3 marks round 1 finding 3 "fixed", but its disclosure half was not done.
  - **What was asked:** state the cross-vendor flow in `specs/sh-agent/dream.md` and README.
  - **What e2cffc1 says:** `git grep -i "provider\|anthropic\|google"` over README, AGENTS.md and both specs finds only the memory sentences (`README.md:223`, `memory-loop.md:11`).
  - **What is not stated:**
    - under `claude` or `all`, the Antigravity conversations go to Anthropic;
    - under `antigravity`, the Claude Code transcripts go to Google.
  - AGENTS.md rule 1's "only to the model of the CLI the machine already uses" hides this second provider.
  - **Remediation:** one sentence in `dream.md` step 1 or 4 and in README's dream bullet.
- **[low / confidence: medium]** `internal/dream/timer.go:92-101`: Outcome 10 fixes the escaping only.
  - **What is fixed:** the refusal of `"`, `\`, CR and LF, and the `%%` escaping, are in place and tested.
  - **What is not:** round 1 called the practical risk "not injection". The unit freezes the installing shell's PATH, including absolute project-local directories that direnv, nvm or a checkout's `bin/` prepend. The 03:00 run can then resolve `claude` or `agy` from a repository. That half is neither fixed nor rebutted in the Outcome; README now only describes it ("the PATH of the shell that installed it").
  - **Partial mitigation:** Go 1.24's `exec.LookPath` already refuses relative entries (ErrDot).
  - **Remediation:** resolve the CLI with `exec.LookPath` at install time and write its absolute path. Or keep only absolute PATH entries outside indexed project roots. Or rebut it in the Outcome as accepted.
- **[low / confidence: low]** `internal/dream/run.go:198-222` and `:249`: the agy call fails only on a denied tool.
  - **The settings file:** `{"permissions":{"allow":[]},"trustedWorkspaces":[]}` has no `toolPermission` key, so the refusal rests on agy's default for a missing key. Earlier logs show that default loaded as `request-review`; no log in the range shows it for this file.
  - **The answer check:** `answer` rejects a result only when it names `denied_actions`. A tool agy runs without asking (whatever plan mode allows unprompted, such as a workspace read or a fetch) leaves no trace that `answer` checks.
  - **An unused option:** `agy --help` offers `--sandbox` (terminal restrictions), which the call does not use.
  - **Remediation:**
    - Write `"toolPermission":"request-review"` explicitly and assert it in `TestRunAgyPrivateHomeAndStdin`.
    - Treat any stream event that records a tool call as no answer, not only a denied one.
    - Consider `--sandbox`.
    - One probe of the stream's event types settles the confidence.

## Outcome check (round 1 → e2cffc1)
| # | Claim | Code / test | Holds |
|---|---|---|---|
| 1 | agy private HOME, grants nothing | `run.go:228-256`. The cleanup is deferred after `unlock`, so it runs first, under the lock. `TestRunAgyPrivateHomeAndStdin` covers the private HOME, no `always-proceed`, `"allow":[]`, the login copied and the HOME removed; there is also verify mutation (b) | yes, for the operator's grants; see finding 5 |
| 2 | Claude `--strict-mcp-config --setting-sources local` | `dream.go:177-188`, `TestCommand`, the args assertion in `TestRunWritesProposalsDraftsAndState`; `claude --help` confirms both flags | yes; managed settings still apply (admin-owned) |
| 3 | agy reader drops tool output | `search/dream.go:223`, a type-132 fixture, mutation (a) | partly: findings 1 and 3 |
| 4 | prompt off argv | `call` sends the prompt on stdin via `json.Marshal`; the argv is `-p ""`; the test asserts no prompt in argv | yes: the size limit and `/proc` visibility are both closed |
| 5 | harness proposals | `retro --harness` has the data-not-instructions rule and no reuse as drafted. `o.line` folds CR and LF, and `strings.Fields` also folds U+0085, U+2028 and U+2029. `TestRenderFoldsLineBreaks` covers it, and the rebuttal is accepted | yes as asked; the multi-line body can still hold `### N.` lines, and since the model can add real proposals through the JSON anyway, `/retro`'s rule is the guard |
| 6 | privacy statements | project-memory, AGENTS rule 1, README and the memory-loop spec are restated | yes, for memory |
| 7 | `~` in `project:` and `left out` | `o.line(c.Root)`, and `Dropped` is scrubbed (`collect.go` loop); `/retro` source 9 uses the `~` form | code yes; no assertion of the `~` form in the range's test diff |
| 8 | symlinked trail | `linked`, `TestCollectDoesNotFollowLinkedTrail` | partly: finding 2 |
| 9 | flock locks | both locks use flock; `TestRunLock` (a real flock), `TestHeldLockRefusesAfterRetry`, `TestLeftoverLockFileIsNoLock` | yes |
| 10 | unit escaping | `timer.go:34,98-100`, `TestInstallTimerRefusesUnsafePaths` | partly: finding 4 |
| 11 | leftover copies swept | `removeCopies` under the lock, `TestRunRemovesLeftoverCopies` | yes, with the notes under Checked |

## Checked
- **Secrets.** Needle counts over every added line:
  - `ya29.`, `sk-ant-`, `ghp_`, `github_pat_`, `-----BEGIN`, `glpat-`, `GOCSPX-` and `refresh_token` hit once each. All eight are on the round 1 report's Checked line that lists the needles by name: names, not values.
  - Full AKIA and AIza shapes, `xox[abps]-`, assignment shapes and the home path: 0.
  - The test login is the literal `token`.
  - The login copy is created with `O_EXCL` 0600 in a 0700 `MkdirTemp` under the config dir. Only regular files are copied (Lstat). It is removed by defer and swept by the next real run. Errors name the file's path, never its contents.
  - The trail markdown (3 review reports, the change entry, the plan hunk, specs, deploy guides, skill-matrix rows, skills, README, AGENTS.md) holds no credential or private data beyond posture notes. Its instruction-shaped text is the harness's own procedure.
- **Injection.**
  - The stdin message is built with `json.Marshal`, so prompt text cannot add an NDJSON line. The argv is constant.
  - `answer` scans line by line with a 16 MiB cap and fails closed to `keepFailure`, which scrubs (`run.go:456`).
  - The SQL is a constant with a literal `IN` list.
  - Unit text is refused or escaped as described in Outcome 10.
- **Authz and authn.**
  - `memory.trusted` is unchanged and still runs before the lock.
  - Neither lock opens with `O_NOFOLLOW`, and the dream lock truncates its file. Both files sit in operator-owned places: the config dir is 0700, and the store is refused when it is symlinked or tracked. There is no attacker path.
  - Go opens files `O_CLOEXEC`, so the CLI child does not hold the dream lock.
  - `withHome` replaces HOME only. This shell sets no GEMINI, ANTIGRAVITY or XDG_CONFIG_HOME variable.
  - Hardening note, not a finding: `agyHome` uses `$HOME` unchecked. With HOME empty and an explicit `--home`, the login would be read relative to the cwd. `filepath.IsAbs(real)` closes this, and without `--home`, `main.go:131` fails first.
- **Outcome 11 notes (not security defects).**
  - `removeCopies` runs after `readState`, so a corrupt state keeps a killed run's login copy until the state is fixed.
  - A dry run takes no lock, so a concurrent real run can delete its copies. The comment "no live run owns them" does not hold for dry runs.
  - With no signal handling, SIGINT and SIGTERM skip the defers.
- **Crypto:** nothing added, nothing hand-rolled.
- **Dependencies:**
  - `go.mod` and `go.sum` have 0 changed lines.
  - `govulncheck ./...` exits 0: "Your code is affected by 0 vulnerabilities". It also reports 1 vulnerability in a required module that the code does not call; it is not identified, as in round 1.
  - There is no `package.json`, so `npm audit` does not apply.
- **GCP surface:** none in the range.

## Not examined
- **Release-diff paths.** Every hunk was read in full, and the new files (3 reviews and the change entry) were read whole. Code outside the hunks that was not re-read: `run.go` `write()` and `parse()`, unchanged since round 1; `harness.Redact`; README outside its two hunks.
- **Uncommitted changes that appeared during this review.** They are outside the range and were not reviewed:
  - modified: `internal/dream/{collect,dream,run}.go` and their tests, `specs/sh-agent/dream.md`, the plan and `bugs/INDEX.md` (+138/−23 against e2cffc1);
  - untracked: BUG-116 to BUG-122, `reviews/…-review-f1ca525-e2cffc1-code.md`, and `testing/2026-09-30-verify-parity-review-round-1.md`, which was seen only through grep hits.
- **agy runtime.** Neither CLI was run, because a probe sends data to a provider. Unverified:
  - which tools run unasked with no `toolPermission` key and an empty allow list;
  - the stream-json event types;
  - whether agy expands `@path` mentions or slash commands inside a stream-json user message (the prompt carries transcripts that contain them).
- **Claude Code under the new flags.** Unverified: whether `--setting-sources local` still loads the user CLAUDE.md and skills (text only under `--tools ""`), and how claude.ai connectors behave under `--strict-mcp-config`. The Outcome's probe log is not in the range.
- macOS launchd on a real Mac.
- The one uncalled module vulnerability, which was not identified.

## Outcome
1. medium, reply steps carry tool-call arguments: fixed. A run that is a JSON argument object (`{"…`) is dropped (internal/search/dream.go). `TestAgyConversationText` has a reply step with a `CommandLine` and a `CodeContent` object that must stay out, and the prose beside them must stay in. On a copy of a real conversation (the gate18 handoff run) the rule keeps 21 prose runs and drops 32 argument runs, the largest 1692 bytes.
2. low, `.agents` itself as a link: fixed in the round-2 code fixes. Every component from the root down is checked (`TestCollectDoesNotFollowALinkedAgentsDir`).
3. low, remaining gaps:
   - The cross-vendor flow is now stated in `specs/sh-agent/dream.md` § Where the data goes.
   - The timer keeps only absolute PATH entries (`TestInstallTimerKeepsOnlyAbsolutePath`).
   - `answer()` needs `status: SUCCESS` (`TestAnswerNeedsASuccessfulResult`).
   - The private settings set `toolPermission: request-review` explicitly. A real sandbox run with the operator's `always-proceed` beside it answered with 2 proposals.
