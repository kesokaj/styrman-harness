---
owner: security
date: 2026-09-28
---
# Security review, round 3 — update-from-zip   (2026-09-28)

Range reviewed: 1cb2608..worktree
- `git diff 1cb2608`: 18 modified files, 273+/34-. They did not change during this run.
- The untracked files: the 18 present at the start of this run, plus two that appeared during it and were read (the phase-3 verify log and the code-r3 review). This report is excluded.

## Status: COMPLETE

## Verdict: FINDINGS

**The worst problem is still the round-1 medium, open by decision.** `update --from-zip` installs an unsigned archive as the global instructions of two agents, and those agents run commands with the operator's credentials. Round 2 set two conditions for the `UNVERIFIED` line to count as an interim:
- Condition 2 is now met. The signed-zip option is recorded as an open operator decision, in the plan's Follow-ups and in ADR-038's Consequences.
- Condition 1 is still pending. No trail file records the operator accepting the risk. A gate must not read ADR-038 as that acceptance.

**Everything else is low.** Of the six round-2 items the dispatch named:
- four are closed;
- one is closed with an accepted residual (the 64 MiB cap);
- one is partly closed (the `RemoveAll` warnings).

**New this round, reproduced:** the stale-directory sweep runs before the zip is even opened. As a result:
- a refused run deletes stale `zip-*` directories;
- so does `update --from-zip --dry-run`, where the operator left out the file name.

A stale comment still says a refused zip writes nothing. The reach is only sh-agent's own leftovers, so this is not a security exposure. It is a refusal and read-only contract that the code does not keep.

## Findings

- **[medium / confidence: high]** Round-1 medium, **open by decision**: an unverified archive becomes agent instructions, from any binary, including a signed release binary.
  - Unchanged since round 2:
    - The `UNVERIFIED` line is adequate disclosure, but it is not a control.
    - `status` never says "unverified".
    - The printed sha256 has no reference value to compare against.
  - Condition 2 is met: plan § Scope "Follow-ups noted" ("A signed zip …, needs the operator's approval") and ADR-038 Consequences ("awaits the operator's decision").
  - Condition 1 is open: the operator's acceptance is not recorded anywhere.
  - Inconsistency: ADR-038's "Rejected" bullet still gives "a zip is unsigned" as the reason. That is a choice about the release process, not a constraint, and it reads as contradicting the Consequences, which call the same control pending. Reword the bullet to "deferred to the operator (see Consequences)".
  - Remediation, unchanged and needing no new dependency:
    - `/preflight` publishes a `git archive` zip of the tag and lists it in the already-signed `SHA256SUMS`.
    - `--from-zip` checks the zip's `crypto/sha256` against that list through `release.go`'s existing `ssh-keygen -Y verify` path.
    - Optionally, `status` says `unverified` for a zip install.
  - Until that is built, the operator's explicit acceptance belongs in ADR-038 or the release report before v1.8.0 ships.

- **[low / confidence: high]** `cmd/sh-agent/main.go:1052-1053, :1064, :1076` (`updateFromZip`): **a refused or mistyped run still deletes, and two comments overstate the contract.**
  - The order of work on a real run:
    1. `MkdirAllTracked(parent)` creates the chain.
    2. `sweepStaleZipDirs` deletes `zip-*` directories older than an hour.
    3. Only then does `UnpackZip` open the file.
  - On a refusal, the chain is undone. The sweep is not.
  - **Reproduced** in `.agents/scratch/secr3-5NjxVn/`. Each home held a 3-hour-old `zip-old` with a file in it. The runs used `env -i` with a scratch `HOME` and an empty `PATH`.

    | Probe | Command | Exit | `zip-old` afterwards |
    |---|---|---|---|
    | C | `update --from-zip <missing file> --home H` | 1, `open … no such file` | **deleted** |
    | D | `update --from-zip --dry-run --home H` | 1, `open --dry-run: no such file` | **deleted** |
    | E (control) | `update --from-zip <missing> --dry-run --home H` | 1 | kept |
    | F | a refused zip into a home that did not exist | 1 | not applicable: the whole created chain was removed again (`homeF` absent) |

    - In D, the flag package takes `--dry-run` as the file name, so the run is real and it deletes, although the operator asked for a dry run.
    - F shows the undo itself works.
  - Three texts are false or overstated:
    - `:1052-1053`: "UnpackZip creates the parent only after the zip's names pass, so a refused zip writes nothing". This is false since D13. Code r2 named this comment and code r3 names it again.
    - `:1064`: "so a refused zip leaves the home as it was (BUG-093)". This holds for the chain, not for the sweep.
    - BUG-093's prevention rule ("A command that can refuse its input creates nothing before the refusal point"). D13's create-then-undo meets the rule's test clause (an empty home stays empty) but not its first clause.
  - Reach: only this command's own `zip-*` leftovers, within the same user's permissions. That is why this is low.
  - Remediation (shortest):
    - Call `sweepStaleZipDirs` only after `UnpackZip` succeeds. A refused archive and probe D then delete nothing.
    - Replace the `:1052-1053` sentence.
    - Refuse a `--from-zip` value that starts with `-` as "needs a file" (`./-x.zip` still works).
    - Record in D13 that the chain is created and then undone, a deliberate reading of BUG-093's rule.
    - Test: extend `TestUpdateFromZipSweepsStaleDirs` with a refused real run that must keep `zip-old`.

- **[low / confidence: high]** `internal/harness/zip.go:88` and `cmd/sh-agent/main.go:1082`: **`RemoveAll` errors are still discarded at two of the three sites** (round-2 item, partly closed).
  - Only the sweep warns now. `UnpackZip`'s deferred `os.RemoveAll(tmp)` and `updateFromZip`'s `defer os.RemoveAll(dir)` stay silent.
  - What a failed removal leaves:
    - An unpacked tree left under `~/.config/sh-agent` is reported only when a later real run sweeps it, after an hour.
    - One left under the system temp directory (a dry run on a fresh home) is never swept or reported.
    - On an error path it also stops `RemoveCreated` from emptying the chain.
  - The content is the zip's five trees, inside a 0700 `MkdirTemp` directory, so there is no exposure. Hence low.
  - Remediation:
    - Make the `updateFromZip` defer a closure that prints `warning: could not remove …`. The function already has `out`.
    - `UnpackZip` joins its cleanup error into the error it returns.

- **[low / confidence: high]** Pre-existing, outside the diff: `internal/harness/report.go:411` `sweepInterruptedClones` still builds `filepath.Glob(filepath.Join(home, ".config", "sh-agent", "report-*"))` from an unescaped home.
  - Round 2's "give `report.go` the same fix" was not taken.
  - It is gated by the `report-<pid>` name and a dead pid, so its reach is another home's interrupted-report leftovers.
  - It does not block this gate. The same `os.ReadDir` + prefix fix applies.

- **[low / confidence: medium]** Carried over from round 2, unchanged since (mtime 17:26:34, before round 2 closed): `.agents/docs/testing/skill-matrix.md:13`.
  - The unrestricted `Bash(curl:*)` / `command(curl)` grant is still there.
  - So is the sentence "The protection is that the copied token lives only for the batch and is deleted by name after it". Deleting the local copy does not revoke a copy that was exfiltrated during the batch.
  - The operator approved option 1 (D5). Round 2's remediation still applies: `curl` to localhost only, state the residual risk, and sign out or revoke the sandbox session after each batch.

- **[low / confidence: high]** Carried over from round 2: **texts that state a removed property** (partly closed).
  - Closed: the usage line and the `updateFromZip` doc comment are now accurate ("install its harness and the running binary (nothing in the zip replaces the binary)").
  - Open:
    - `internal/harness/repo.go:23-33`: "A RELEASE binary installs the harness built into it and nothing else … not one named on the command line".
    - `DECISIONS.md:287`: ADR-034's header is still "— Accepted", with no "Superseded in part by ADR-038" mark, which the file's own header requires.
    - ADR-038 "What is unpacked" omits the 64 MiB archive cap and the 1,024-byte / 32-level name bounds.
    - `release.go:350`: `readCapped`'s comment ("reads a file a release supplied") now also covers the operator's zip.
    - R9's "the binary itself is not replaced" is routed to the release report by D12, which is acceptable if that happens.

## Round-2 status — the six items the dispatch named

1. **Glob sweep: CLOSED.**
   - The new code: `os.ReadDir(parent)`, a `zip-` name prefix, and the `DirEntry` type check (a symlink or a non-directory is skipped), with `Lstat` for the age. No pattern is matched against a path.
   - `TestUpdateFromZipSweepWithGlobCharsInHome` rebuilds round 2's probe-B arrangement: `h[1]` beside `h1`, both with a `zip-old` aged 2 h. It asserts both directions: the home's own directory is swept and the other home's is kept.
   - The `feedback-old` case pins the prefix.
   - Green here. The must-fail mutation is the main agent's, reported by code r3; I did not repeat it.
   - The residual is in `report.go` (finding 4).
2. **Plan-refusal leftover: CLOSED.**
   - `updateFromZip` now owns the chain: `MkdirAllTracked` runs before `UnpackZip`, and `RemoveCreated` runs in a defer keyed on the named result.
   - The `:=` at the function's outer scope reuses the named `err`; it does not shadow it.
   - The defers run last-in, first-out: the `zip-*` tree first, then the chain.
   - `os.Remove` never deletes a non-empty directory, so a partial install keeps what it wrote.
   - `TestUpdateFromZipRefusedContentWritesNothing` uses round 2's probe-A arrangement (a skill with no description, on a fresh home). Code r3's shadowed-`err` mutation compiled and failed it. Probe F here confirms that a missing home is created and then removed in full.
   - Nit: the test asserts only `err != nil`. Asserting that the error names "no description" would prove the refusal came from `Plan` and not from an earlier stage.
   - The ordering has a cost: finding 2.
3. **64 MiB cap: CLOSED, residual accepted.**
   - `maxZipArchiveBytes = 64 << 20` is about 23 times the real archive (2,853,062 bytes) and a quarter of round 2's cap.
   - Round 2's worst-case peak scales down to roughly 400–500 MB. That is reasoned, not measured.
   - Not taken: reading into a buffer sized from `Stat` with `io.ReadFull`. `readCapped` still uses `io.ReadAll`, whose growth can transiently double the 64 MiB. This is acceptable for a file the operator chose.
4. **`MkdirAllTracked` treating `EEXIST` as ours: CLOSED in code.**
   - `errors.Is(err, fs.ErrExist)` → `continue`. The path is never appended, so it is never removed.
   - Untested: only a race reaches it, and a deterministic test needs a seam (a package-level `mkdir` variable). Given how little it exposes, I do not ask for one.
5. **`RemoveAll` warnings: PARTLY CLOSED.** The sweep warns; two sites still do not (finding 3).
6. **TMPDIR dry-run test: CLOSED.**
   - `TestUpdateFromZipDryRunCleansTempDir` points `TMPDIR` at a fresh directory, after `isolate` has created the test's temp root elsewhere. It dry-runs on a fresh home and compares the directory's entries before and after.
   - A missing `defer os.RemoveAll(dir)` would leave a `zip-*` directory there.
   - Limit: the test proves nothing is left in `TMPDIR`, not that the unpack went there. `TestUpdateFromZipDryRunWritesNothing` covers the home side.

Test isolation (the BUG-094 class, D13), checked by reading:
- 12 of the zip CLI tests call `isolate` (empty `PATH`, scratch working directory, `SH_AGENT_REPO` empty).
- `TestFromZipRefusedOnOtherCommands` empties `PATH` but keeps its working directory. Code r3 flags that D13 overstates this.
- The two dry-run tests use a fake `gh` and a local bare origin respectively.
- No red path I traced reaches the network or the real checkout's remote.
- Optional hardening: `isolate` does not set `HOME`. Every call passes `--home`, but `t.Setenv("HOME", t.TempDir())` would cover a regression in home resolution.

## Checked

- **Secrets.**
  - I regex-scanned the full `git diff 1cb2608` and every untracked file, including the two that appeared mid-run. The patterns:
    - GitHub, AWS, Slack, Google and Anthropic token shapes;
    - private-key armour and `ssh-ed25519`/`ssh-rsa` keys;
    - `refresh_token`, `client_secret` and `access_token`;
    - `password=` and `api_key=`;
    - email addresses;
    - "ignore previous instructions" and "disregard" phrasing.
  - Positive control: a synthetic `ghp_` token and a synthetic address both fired (2 of 2).
  - Hits:
    - in the diff: none;
    - in the untracked files: only rounds 1 and 2's own descriptions of their scans, and `f@example.invalid`, a test git identity.
- **Trail markdown**, read end to end for secrets, credentials, private data and instruction-shaped text:
  - Files changed since round 2: the update-from-zip plan, the e2e log, the ADR-038 hunk of `DECISIONS.md`, and round 2's report.
  - The code-r2 review: round 2's Not-examined gap, now closed.
  - Files that appeared during this run: the phase-3 verify log and the code-r3 review.
  - Result: no credentials. The only private data is the operator's GitHub handle and scratch paths, both already in this private repository. Nothing is instruction-shaped: no text aimed at an agent redirects it.
- **Covered by round 2 and unchanged since** (mtime earlier than round 2's report at 17:38:31; relied on per the dispatch):
  - BUG-093, BUG-094 and the recon plan;
  - the three cold reads, the code-r1 and docs reviews, and the security round-1 report;
  - the phase-1 and phase-2 verify logs;
  - `internal/harness/zip_test.go`;
  - the hunks in `bugs/INDEX.md`, `deploy/local.md`, `deploy/remote.md`, `delivery.md`, `skill-matrix.md`, `AGENTS.md`, `README.md`, `REQUIREMENTS.md`, `ROADMAP.md` and `VENDORED.md`.
  - I re-read the code hunks of `apply.go`, `plan.go`, `repo.go`, `source.go` and `.gitignore` anyway.
- **Code read this round:**
  - `internal/harness/zip.go`, end to end;
  - `cmd/sh-agent/zip_test.go`, end to end;
  - the whole `main.go` and `main_test.go` diff;
  - `installFiles`, `readCapped`, `sweepInterruptedClones`, and `repo.go:18-40`.
- **Injection:** there is no shell, SQL or template use. The sweep joins `ReadDir` names to its own parent and nothing else. Entry names are validated before any write, unchanged since round 2.
- **Authz/authn:** unchanged since round 2:
  - The refusal block runs before any command.
  - `--check` and `--repo` are refused with `--from-zip`.
  - A `zip:` value is never a candidate path.
  - The trust model is finding 1.
- **Crypto:** only `crypto/sha256` from the standard library, used as a checksum. Nothing is hand-rolled.
- **Dependencies:**
  - No new module.
  - `govulncheck ./...` under go1.27.1 printed "No vulnerabilities found." and exited 0.
  - `npm audit` does not apply: there is no `package.json`.
- **GCP surface:** none touched.
- **Tools:**
  - `gofmt -l cmd internal`: clean.
  - `go vet ./...`: exit 0.
  - `go test -race -count=1 -run 'Zip|FromZip|DryRun|SecretShaped' ./cmd/sh-agent ./internal/harness`: ok for both packages.
- **Probes:** `.agents/scratch/secr3-5NjxVn/` (gitignored, about 7.4 MB, left in place).
  - The binary was built from the worktree.
  - Every run used `env -i` with a scratch `HOME` and an empty `PATH`.
  - The scratch `HOME` and `PATH` directories were still empty afterwards, and the real home was never named.
  - HEAD is still 1cb2608, with no tag on it.

## Not examined

- Nothing in the release diff is unread: every path is either read this round or covered by round 2 and unchanged since (listed above). A file that changes after this report is not covered.
- `main.go` outside its diff hunks and `installFiles`; `release.go` beyond `readCapped`; `report.go` beyond its sweep; `Install`, `inventory.go` and `state.go`. None of these was re-read end to end.
- No mutation runs of my own. The mutations cited are code r3's and the main agent's, as those reports record them.
- The full `go test -race ./...`: not rerun by me. The phase-3 verify log records 160 passing and 0 failing.
- The memory figure under item 3 is reasoned, not measured. No hostile archive was built.
- macOS.

## Outcome

See `changes/2026-09-28-1752-update-from-zip.md` and plan D11/D13. After this round:
- **Fixed:** the stale sweep now runs only after a successful install; `--from-zip` refuses a value that starts with `-`; a failed `RemoveAll` of the unpacked copy prints a warning; the stale comment, usage and package doc are updated; `TestFromZipRefusedOnOtherCommands` uses `isolate`; ADR-038's Rejected entry now says deferred; ADR-034 is marked superseded in part.
- **Left:** the unsigned-zip medium, open by decision and put to the operator; the pre-existing `report.go:411` glob, outside this diff and a follow-up; the matrix curl grant, accepted under option 1 (the allowlist is stated not to be a boundary).
