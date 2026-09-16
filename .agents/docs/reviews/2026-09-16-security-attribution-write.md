---
owner: security
date: 2026-09-16
---
# Security review — attribution-write   (2026-09-16)

Range reviewed: working tree, cmd/sh-agent/main.go and main_test.go

## Blobs pinned — and a warning: the tree moved mid-review

| file | blob at HEAD `1254672` | blob I first read (11:10) | blob I finished against (11:16) |
|---|---|---|---|
| `cmd/sh-agent/main.go` | `27d03a0e3bc564ee1168ee272ba1fb24df61b339` | `672989cff1dfae6ae413bcc08fdd228b072ae7cd` | **`0657c04d311b738a90a0e47e48bc2275cff27cae`** |
| `cmd/sh-agent/main_test.go` | `fa2196a50afe393015ae6369a799b3556b542562` | `df476099a12753e037f70b4d65d932661b0288cc` | **`1b0a61bea89cb5a5462bf7683543f201d9bdd71b`** |

Both files were rewritten by another writer at `11:15:50` while this review was running, between my read and my first test run.
The `672989c`/`df47609` blobs were only ever hashed, never written to the object database (`git cat-file -p 672989c` -> "Not a valid
object name"), so they are unrecoverable and **the findings below are stated against `0657c04`/`1b0a61b`, the bytes on disk now**. The
bold column is authoritative. If the tree has moved again since, re-pin before acting on this. What changed in that window is recorded
in "Already fixed under me" so the caller can tell which of my observations were overtaken rather than wrong.

Command: `git diff -- cmd/sh-agent/main.go cmd/sh-agent/main_test.go`. The other modified files in the tree (`REQUIREMENTS.md`,
`ROADMAP.md`, `context/AGENTS.md`, the plan and the three cold-read reviews) were out of scope by instruction and were not read.

## Verdict: FINDINGS

Nothing high or medium. The write is careful in the ways that matter: it refuses a symlink outright via `os.Lstat`, refuses a
directory or FIFO via `!IsRegular()`, refuses anything over 1 MiB before reading it, parses the file into a map before touching it,
re-validates the spliced result with `json.Valid` before it lands, writes through a temp file in the same directory, restores the
original permission bits, and renames. I drove twenty crafted inputs through a binary built from the current blob and could not get it
to corrupt a file, write outside `--home`, hang, panic, or echo a byte of anyone's settings. What is left is a cluster of low-severity
items, the sharpest being that the create branch still opens with `os.WriteFile` — `O_WRONLY|O_CREATE|O_TRUNC`, no `O_EXCL`, no
`O_NOFOLLOW` — so the careful `Lstat` gate is advisory rather than enforced: anything that can write into `~/.claude` in the window
between the two syscalls can still redirect the write through a planted symlink, or wedge the install forever with a planted FIFO.
That precondition is an attacker who already owns the operator's config directory, which is why it is low and not medium, and one flag
closes it.

## Findings

- **[low / confidence: high]** `cmd/sh-agent/main.go:491-505` — TOCTOU: the `Lstat` gate is not enforced by the write, and a planted
  FIFO hangs the install with no timeout. `os.Lstat` decides the file is absent, then `os.WriteFile` opens
  `O_WRONLY|O_CREATE|O_TRUNC` with neither `O_EXCL` nor `O_NOFOLLOW`. A local process that can create entries in `<home>/.claude`
  during that window can (a) plant a symlink, which `open(2)` follows, putting the write on any path the operator can write — the one
  way this code still creates a file outside `--home`; or (b) plant a FIFO, where `open(2)` for writing blocks until a reader appears,
  and there is no timeout, deadline, or context anywhere on the path, so `sh-agent install` stops with no further output. Both are
  closed by one change with no behaviour cost: `os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)` instead of
  `os.WriteFile` — losing the race then yields `EEXIST`, which the existing "could not be written (%v)" line already handles, and a
  FIFO or symlink at the path becomes that same error instead of a hang or a redirect. Confirmed by test that a *pre-existing* FIFO or
  directory is refused cleanly (no hang, file untouched); it is only the race window that is open.
- **[low / confidence: high]** `cmd/sh-agent/main.go:562-583` — no `Sync` before the rename. The temp file is written and `Close`d but
  never fsynced, and `os.Rename` makes the directory entry atomic, not the data behind it. A crash or power loss in the window can
  leave `settings.json` zero-length or partial — and that file carries `permissions.allow`, the list deciding what Claude Code may run
  without asking. Losing it fails *open into prompting*, not into silent grants, so the security direction is safe; it is still the
  operator's security config being replaced without a durability barrier. Remediation: `tmp.Sync()` before `tmp.Close()`.
  `internal/harness/apply.go:240-263` has the identical pattern, so this is a codebase-wide nit rather than a regression — worth fixing
  in both or neither.
- **[low / confidence: high]** `cmd/sh-agent/main.go:521-524` + `:645` — two lines of the same output contradict each other for a file
  that is valid JSON but not an object. Input `["x"]` produces, in order: `settings: <path> is not a JSON object; add the key by hand`
  (from `attributionSet`, correct) and then `note: <path> is not valid JSON, so the attribution setting could not be checked` (from
  `attributionHint` -> `attributionOff`, incorrect — the file *is* valid JSON, it just does not unmarshal into that struct). Verified by
  running it. An operator reading the second line will go looking for a syntax error that is not there. The cause is `attributionOff`
  collapsing every `json.Unmarshal` failure into `"is not valid JSON"`; it predates this diff, but this diff is what puts the two lines
  side by side. Remediation: have `attributionOff` distinguish a type mismatch from a parse failure, or let `attributionSet` suppress
  the hint once it has already explained the file.
- **[low / confidence: high]** `cmd/sh-agent/main.go:481-482` (doc comment) — "every other byte of the file stays" is false for a CRLF
  file, and a description that overstates a guarantee is read as one. The splice does `bytes.TrimRight(trimmed[:closing], " \t\r\n")`
  on the body, so the `\r` ending the last member's line is consumed and replaced by `,\n`. Verified: input
  `{\r\n  "model": "x"\r\n}\r\n` comes back as `{\r\n  "model": "x",\n  "attribution": ...\n}\r\n` — valid JSON, other keys intact,
  but one byte of the operator's file changed outside the insertion and the result has mixed line endings. Not exploitable; filed
  because the comment is the contract a later reader will trust. Remediation: say "every other key's bytes stay" and note the line
  ending, or preserve the file's dominant terminator in `insert`.
- **[low / confidence: high]** `cmd/sh-agent/main.go:10,39` — `install` now writes a file `uninstall` does not remove, while the usage
  text still promises `uninstall [--dry-run]   remove exactly what install wrote, nothing else`. `grep -rn settings.json internal/`
  returns nothing, so `harness.Uninstall` has no knowledge of it and the manifest never lists it. Not a vulnerability; it is a stated
  guarantee that the diff made untrue, and the safest resolution is probably to amend the sentence (leaving an operator's settings
  alone on uninstall is the right behaviour) rather than to start deleting settings files.
- **[low / confidence: low]** `cmd/sh-agent/main.go:512-514` vs `:585-601` — the 1 MiB guard is checked against `info.Size()` from the
  `Lstat`, but the read that follows is a separate `open`+`io.LimitReader(f, 1<<20)` with no size re-check. A file that grows past
  1 MiB between the two is read truncated, and the truncated buffer is what `json.Valid` blesses and `rename` installs. Reaching it
  needs the file to be rewritten larger mid-install *and* the first megabyte to still parse as a complete JSON document, which for
  strict JSON means the overflow is whitespace. Very narrow — hence low confidence that it is reachable at all. Remediation if it is
  ever worth closing: read `1<<20 + 1` bytes and treat `len(data) > 1<<20` as the same refusal the size gate already prints.
- **[low / confidence: high]** dependencies — `govulncheck ./...` reports the repo's **toolchain** affected by two stdlib advisories.
  Pre-existing, not introduced here, and not reachable from the new code, but `GO-2025-3956` is squarely in this binary's habits:
  `harness.Reinstall` executes whatever `exec.LookPath("go")` returns. Details in Checked. A toolchain bump fixes both.

### Already fixed under me (present in `672989c`, gone in `0657c04`)

Recorded so the caller can tell these were real and are closed, not that I missed them. In the blob I first read, `attributionSet`
gated on `os.Stat`, which follows symlinks. That made two no-race defects: a **dangling** symlink at `<home>/.claude/settings.json`
reported `ErrNotExist`, so the create branch ran and `os.WriteFile` followed the link and created the file *at the link's target*,
outside `<home>`; and a **live** symlink (a settings.json linked into a dotfiles checkout — the pattern this project's own
`/init-project` uses for `.claude/*`) was read through and then `os.Rename` replaced the *link itself* with a regular file, orphaning
the version-controlled copy while printing "other keys untouched". The same blob also had no size gate, so a >1 MiB file was read
truncated and the truncated buffer could be written back, and `null` at top level relied on the `closing <= opening` arm to catch it.
`0657c04` fixes all four: `os.Lstat` with an explicit `ModeSymlink` refusal ordered first in the switch, an `info.Size() > 1<<20`
refusal, an explicit `top == nil` check for bare `null`, and an `abandonTemp` helper that reports a temp file it could not remove. I
re-ran the symlink cases against the current blob: dangling and live symlinks are both refused with
`settings: <path> is a symlink; add the key to its target by hand`, the link survives, and no file appears outside `<home>`.

### Crafted-input matrix — exercised, not reasoned

Black-box against a binary built from `0657c04` into `.agents/scratch/secrev/` (no source created), one scratch `--home` per case,
sha256 of the file compared before and after. Driver: `.agents/scratch/secrev/{drive,rest,links}.sh`.

| input | file after | message |
|---|---|---|
| UTF-8 BOM before `{` | unchanged | advisory: not valid JSON |
| CRLF | key added | correct; see the CRLF finding for the one trimmed byte |
| trailing `// comment` (JSONC) | unchanged | advisory: not valid JSON |
| second top-level value `{"a":1}\n{"b":2}` | unchanged | advisory: not valid JSON |
| bare `null` | unchanged | `is not a JSON object` |
| array `["x"]` | unchanged | `is not a JSON object` + a contradicting hint (finding above) |
| `{}` | key added, no leading comma | correct |
| `{"a":"}}}"}` (braces inside a string) | key added after the real closer | correct |
| existing `attribution` (any value) | unchanged | left to the operator, hint fires |
| directory at the path | n/a | `is not a regular file`, no hang |
| **FIFO** at the path | still a FIFO | `is not a regular file`, **no hang** (60 s timeout unused) |
| symlink, dangling | link intact, no file created anywhere | `is a symlink` |
| symlink, live | link intact, target unmodified | `is a symlink` |
| mode `0o600` | key added, **mode still 0600** | correct |
| mode `0o000` | unchanged | `could not be read` |
| 2 MB file whose first bytes parse | unchanged | `is larger than 1 MiB` |
| 20 000-deep nesting | unchanged, **no panic** | Go's scanner depth cap (10 000) rejects it first |
| `--dry-run`, file absent | not created | `[dry-run] would create …` |
| `--dry-run`, file present | unchanged | `[dry-run] would add …` |
| `--vendor agy`, file present | unchanged | no settings line at all |
| leftover `.settings-*` temp files after all runs | **none** | `find` across every scratch home |

Two structural properties make the splice sound and are worth keeping in the comment: once `json.Unmarshal(data, &top)` succeeds with
a non-nil map, `data` is a valid JSON document whose top-level value is an object, so the last non-whitespace byte *must* be that
object's `}` and the first `{` *must* be its opener — `LastIndexByte`/`IndexByte` cannot mis-target, which is what the `}}}` -in-a-string
case confirms. And `info.Mode().Perm()` masks setuid/setgid/sticky off, so a strangely-moded original cannot propagate those bits.

## Checked

- **Secrets** — none. The only literal written is `"attribution": {"commit": "", "pr": ""}` — two empty strings. No credential, token,
  key, or path to one appears in either file; grep for `token|secret|key|password|credential` over the diff returns only "keys"
  meaning JSON keys.
- **Injection** — no shell, SQL, or template in the new code. The only external input is the bytes of `settings.json`, flowing to
  `json.Unmarshal`, `json.Valid`, `bytes.TrimRight`/`IndexByte`/`LastIndexByte`, and a byte splice — never to `exec`, never as a format
  string. The path is `home` (operator flag) joined with the constants `harness.ClaudeRoot` and `"settings.json"`, so no path element
  is attacker-controlled; the residual path weakness is symlink resolution in the race window, filed above.
- **Authz/authn** — none in scope, but the diff does move a trust boundary and that is the headline for reviewers: `sh-agent install`
  now *mutates* a file it previously only read, and that same file holds Claude Code's `permissions.allow` — the list governing what
  the agent may execute without prompting. The mutation is correctly fenced (map-parse first, key-absent check, `json.Valid` on the
  result, other keys never reserialized, a present `attribution` of any value left alone), and I could not make it alter or drop
  another key in any of the twenty cases above. No finding, but it is the thing to re-review if the splice is ever loosened.
- **Crypto** — nothing hand-rolled and nothing crypto-adjacent: no primitive, token format, session, or protocol is implemented or
  composed here. `os.CreateTemp`'s randomized suffix is the stdlib's own and is not a security control in this design (the rename is).
- **Dependencies** — none added. `go.mod` is `module github.com/kesokaj/styrman-harness` / `go 1.24.4` with **zero** `require` lines,
  matching the repo's stdlib-only rule; the diff's one new import is `bytes`. `govulncheck ./...` ran
  (`/home/kesokaj/.go/bin/govulncheck`, Go 1.24.4) and reports two stdlib advisories, both reachable from pre-existing code and neither
  from the new code: **GO-2026-4602** (`FileInfo` can escape from a `Root` in `os`; found in `os@go1.24.4`, fixed in `os@go1.25.8`;
  trace `internal/harness/inventory.go:164` `os.ReadDir`) and **GO-2025-3956** (unexpected paths returned from `LookPath` in
  `os/exec`; fixed in `os/exec@go1.24.6`; trace `internal/harness/repo.go:140` `exec.LookPath`). It also counted 2 vulnerabilities in
  imported packages and 43 in required modules that the code does not call. `npm audit` not applicable — no JavaScript in this repo.
  `gofmt -l cmd/sh-agent/` clean; `go vet ./...` clean.
- **GCP surface** — untouched. No IAM binding, service account, `--allow-unauthenticated`, Secret Manager reference, or key file
  anywhere in the diff; the change is local-filesystem only.
- **Output hygiene (asked for explicitly)** — no byte of the file, and no byte of a JSON error, can reach stdout. The `Fprintf`
  arguments across the new code are exactly: `path` (operator-supplied), the `fragment` constant, `tmpName`, and `error` values from
  `os.MkdirAll`, `os.WriteFile`, `os.CreateTemp`, `(*File).Write`, `(*File).Close`, `os.Chmod`, `os.Rename` and `os.Remove` — all
  `*fs.PathError`, carrying a syscall name and a path, never content. The one value that *would* carry input bytes, the
  `json.Unmarshal` error, is deliberately discarded (`cmd/sh-agent/main.go:521`), and the `readSettings` refactor preserved that
  property for `attributionOff`, whose `problem` values are the three constants "could not be read", "is not a regular file", "is not
  valid JSON". The only disclosure is one bit — whether the file is valid JSON but not an object. Confirmed by diffing stdout against
  the file contents for every case in the matrix.
- **Temp-file lifetime (asked for explicitly)** — `os.CreateTemp(filepath.Dir(path), ".settings-*.json")`: same directory as the
  target (so the rename is atomic and never crosses a filesystem), mode `0600` at creation, randomized name, dot-prefixed so Claude
  Code's `settings.json` glob will not read it. Removed on all four post-create error returns, now through `abandonTemp`, which also
  prints the temp path if the removal itself fails. `find` over every scratch home after the full matrix: none left behind. The only
  residue path is a SIGINT/SIGTERM or panic between `CreateTemp` and `Rename`, which leaves a `0600` file containing the operator's own
  settings — no disclosure, and `uninstall` will not clean it (see the uninstall finding).

## Not examined

- **The current blob's own diff was not read as a diff.** `672989c` was never stored in git, so I could not `diff` what I first read
  against `0657c04`; I re-read the whole of `attributionSet`/`abandonTemp`/`readSettings` in the current file and re-ran the full
  matrix against a binary built from it, but the *rest* of `main.go` outside those functions was read only in the original diff. If
  the 11:15:50 rewrite touched anything beyond the attribution block, I did not see it.
- `main_test.go` at `1b0a61b` was read only by grep (confirming new `is larger than 1 MiB` and `is a symlink` subtests at lines 575 and
  587); the full current test file was not read line by line, and `go test -race ./...` was **not run** — the tests were inspected, not
  executed. Test gaps I can still name: no test covers a FIFO, a BOM, CRLF, mode `0000`, or a JSONC comment; every one of those is
  covered only by my scratch matrix, which is not committed and will not run again.
- Whether `{"commit": "", "pr": ""}` actually disables Claude Code's trailers, and whether `"pr"` is a key Claude Code reads at all.
  That is the change's entire purpose and it cannot be verified from this repo. A wrong key name would make the write inert while
  still carrying every risk above. Note that `attributionOff` only ever inspects `commit`, so `pr` is unverified on both sides — the
  tool cannot detect its own setting being wrong.
- Whether Claude Code accepts JSONC in `settings.json`. If it does, this code silently skips every commented settings file (verified
  above: a trailing `// comment` makes it refuse), which is a correctness gap, not a security one.
- Two `sh-agent install` processes racing on the same `--home`: last `rename` wins, and the loser's `json.Valid` was computed against a
  stale read. Not exercised; likely benign because both write the same fragment.
- Behaviour when `sh-agent install` runs as root via `sudo`: the file would be created root-owned under root's `$HOME`, and a later
  non-root install would hit `EACCES` and print the advisory line rather than fail. Reasoned from the code, not run.
- A symlinked `.claude` **directory** (as opposed to the settings file) still routes the write to the link's target, because `Lstat`
  only declines to follow the final component. This is pre-existing for the whole installer — `harness.Install` writes every harness
  file through that same directory — so it is out of this diff's blast radius, but it is the remaining answer to "can `--home` produce
  files outside itself": yes, if `<home>/.claude` is itself a link.
- The four out-of-scope modified files. If any of them *describes* what this write does, that prose was not checked against the code —
  and the CRLF finding shows the in-file comment already overstates the guarantee.

## Outcome

- The tree moved under the review (the verifier's symlink flag and the code review's null/size/cleanup findings landed while it ran); the report's own note says it finished against the later blobs, and its "four real defects closed" line is the record of that. This outcome is against the tree after its three findings were applied.
- Finding 1 (low) — fixed: the create path opens with `O_WRONLY|O_CREATE|O_EXCL`, so anything planted at the path after the `Lstat` fails the open instead of being written through or blocked on; the write is followed by `Sync` and a checked `Close`.
- Finding 2 (low) — fixed: `tmp.Sync()` before the rename on the insert path. The same pattern in `internal/harness/apply.go:240-263` is pre-existing and filed as a Follow-up in the plan, not changed here.
- Finding 3 (low) — fixed: a valid non-object file (`["x"]`, `null`) now gets one line, from `attributionHint` ("is not a JSON object, so the attribution setting could not be checked …"); `attributionSet` returns silently for it. A CRLF file keeps its line endings (the insert uses the file's own), test `crlf file keeps its line endings`; the doc comment says "every byte up to the object's last value stays".
- Its 20 crafted inputs left seven files named by their contents in the repository root (the probe wrote to the content string as a path); removed by name, recorded as BUG-044 — the agent's rule is scratch only.
- Round 1 closed 2026-09-16; no second round owed (three low findings, no mechanism finding).
