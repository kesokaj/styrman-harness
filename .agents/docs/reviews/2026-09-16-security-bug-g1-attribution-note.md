---
owner: security
date: 2026-09-16
---
# Security review — bug-g1-attribution-note   (2026-09-16)

Range reviewed: working tree, cmd/sh-agent/main.go and main_test.go

Blobs read (`git hash-object`):
- `cmd/sh-agent/main.go` — `3a7663e894a36bcca0dd568be8df3786254d762a`
- `cmd/sh-agent/main_test.go` — `22e56970b0d2e2ff122c82fe3548e023b610d2c4`

## Verdict: FINDINGS

The new `attributionHint` makes a file the harness explicitly does not own — `<home>/.claude/settings.json` — a hard dependency of `sh-agent install` and `sh-agent status`. Every `os.ReadFile` error other than `ENOENT` is returned up and aborts the command: reproduced against a scratch home, a mode-000 file, a directory, and a symlink loop each made `status` exit 1 *before* printing its drift verdict, and made `install` exit 1 *after* writing all 260 files — a fully successful install reporting failure. A FIFO at that path hung the process indefinitely (killed at 10s), because `os.ReadFile` on a non-regular file has neither a size bound nor a timeout. Nothing here leaks a secret in any quantity that matters, but the parse-error branch does echo file bytes to stdout (one character for a syntax error, a whole numeric literal for a number overflow), and `os.ReadFile` follows symlinks, so the bytes need not come from settings.json at all. No credentials in the diff, no injection surface, no crypto, no new dependency (`encoding/json` and `io/fs` are stdlib; `go.mod` still has zero requires), nothing GCP.

## Findings

- **[medium / confidence: high]** `cmd/sh-agent/main.go:483-486` — `data, err := os.ReadFile(path)` followed by `if err != nil && !errors.Is(err, fs.ErrNotExist) { return err }` turns any unreadable `settings.json` into a fatal error for `install` and `status`. Reproduced against `--home .agents/scratch/secreview-home`:
  - directory at the path → `sh-agent: read …/settings.json: is a directory`, exit 1
  - mode 000 file → `open …: permission denied`, exit 1 — and for `install` this fires *after* `260 files managed`, so a complete install exits non-zero
  - symlink loop → `too many levels of symbolic links`, exit 1
  - FIFO → process hung; `timeout 10` killed it (exit 124), no read bound, no deadline

  This matters beyond ergonomics: the repo's own completion gate is "`sh-agent status` from the target home must show a `binary` line without a warning" and `/preflight` reads install exit codes; a root-owned or 0600 `settings.json` baked into a mirror's image (the BUG-015 plain-copy path) turns a good install into a reported failure, and an unattended run has no way to tell that apart from a real one. Remediation: treat *every* read error as "cannot check" — fall through to the same `…could not be checked` note the malformed-JSON branch already prints and `return nil`; additionally `os.Open` + `Stat` and skip anything that is not a regular file, then `io.LimitReader` the read (a settings file is kilobytes; 1 MiB is generous). Never let an advisory note about someone else's file decide this command's exit status.

- **[low / confidence: high]** `cmd/sh-agent/main.go:495` — `fmt.Fprintf(out, "%s~/.claude/settings.json is not valid JSON (%v); …", prefix, err)` prints `json.Unmarshal`'s error verbatim, and that text carries file content. Measured against Go 1.24.4:
  - `*json.SyntaxError` embeds one character of the input: `{"apiKey": "sk-live-…", "attribution": broken}` → `invalid character 'b' looking for beginning of value`
  - number overflow produces `*json.UnmarshalTypeError` with the **whole literal**: `{"attribution":{"commit":1e99999999}}` → `json: cannot unmarshal number 1e99999999 into Go struct field .attribution.commit of type float64`
  - wrong-type errors echo only Go type names, not values — those are safe

  `os.ReadFile` follows symlinks, so the byte need not come from a settings file: a symlink at `<home>/.claude/settings.json` pointing at an env file was read and the note printed `invalid character 'A' looking for beginning of value` — the first byte of `AWS_SECRET_ACCESS_KEY=…`. One byte per run is not an exfiltration channel, and an attacker who can plant that symlink already controls the directory holding the agent definitions the CLI executes, so the practical risk is low. It is still the wrong default for a file where Claude Code users keep `apiKeyHelper`, `env` blocks and tokens, and `sh-agent status` output is routinely pasted into handoffs and issues. Remediation: do not print `err`. Print the byte offset only (`var se *json.SyntaxError; errors.As(err, &se)` → `se.Offset`), or just "could not be parsed".

- **[low / confidence: medium]** `cmd/sh-agent/main.go:495,509` — both notes name the literal string `~/.claude/settings.json` while the function actually read `filepath.Join(home, harness.ClaudeRoot, "settings.json")`. Under `--home .agents/scratch/x` — which this repo's `AGENTS.md` *mandates* for every development install — the note tells the operator to fix a file that was not the file checked; the same holds in a container or mirror home. `pathHint`, twenty lines above, prints the real resolved path (`note: .agents/scratch/secreview-home/.local/bin is not on PATH`), so the inconsistency is within one screen of code. The security flavour is false assurance: the note is the *only* signal that the anti-trailer control is absent, and an operator who edits the wrong home believes a control is in place that is not. Remediation: print `path`, as `pathHint` does.

- **[low / confidence: medium]** `cmd/sh-agent/main.go:494-497` — on a parse failure the function prints "could not be checked" and returns, so the actual warning ("Claude Code adds a Co-Authored-By trailer to every commit until it does") is never printed. The commit message for this change says the context rule alone did not hold in the sandbox — three of three headless commits carried the trailer — which makes this note the enforcement, and a malformed or JSONC-commented settings file silently downgrades it to a shrug. Remediation: print the parse note *and* the warning; unparseable is not evidence the setting is present.

- **[low / confidence: high]** `go.mod:3` — pre-existing, not introduced by this diff, recorded because the review ran the tool: `govulncheck ./...` on go1.24.4 reports two *called* stdlib vulnerabilities — GO-2026-4602 (`os`, FileInfo can escape from a Root; trace `internal/harness/inventory.go:164` `os.ReadDir`, fixed in go1.25.8) and GO-2025-3956 (`os/exec`, unexpected paths from `LookPath`; trace `internal/harness/repo.go:140`, fixed in go1.24.6). The new `os.ReadFile` call is not on either trace. Remediation: bump the toolchain to ≥ go1.25.8 and re-run.

- **[low / confidence: low]** `cmd/sh-agent/main.go:21,76` — the new `io/fs` import is shadowed inside `run()` by `fs := flag.NewFlagSet(cmd, flag.ContinueOnError)`. It compiles today because `run` never references the package, but `errors.Is(err, fs.ErrNotExist)` written inside `run` by the next author will not mean what it reads as. Not exploitable; a trap worth closing with an import alias (`iofs "io/fs"`) or renaming the local to `flags`.

## Checked

- **Secrets** — no credential, token or key added in either file; test fixtures are `{"attribution": …}` and `{"permissions":{"allow":["Bash(git:*)"]}}`, both inert. Confirmed the change only ever *reads* `settings.json`: it appears in no `harness.Plan` entry (`internal/harness/plan.go`), and `pruneEmptyDirs` (`internal/harness/apply.go:269-289`) stops at `filepath.Join(home, ClaudeRoot)`, so install never writes it and uninstall never removes it — the doc comment's claim holds. Secrets *reaching stdout*: bounded to the byte-echo in finding 2.
- **Injection** — no shell, SQL, template or filesystem-write path takes file content. The only value derived from external input is `settings.Attribution.Commit`, consumed by a type switch that compares it and discards it; it is never printed, never joined into a path, never passed to `exec`. Path is `filepath.Join(home, ".claude", "settings.json")` with two compile-time constants — no traversal is reachable from file content; `home` comes only from `--home` or `os.UserHomeDir()`, both the operator's own.
- **Authz/authn** — no auth surface in this binary; the change is a read-only advisory about a local CLI setting and grants nothing. No IDs, no middleware, no confused deputy.
- **Crypto** — none added or touched; nothing hand-rolled.
- **Dependencies** — `encoding/json` and `io/fs`, both stdlib; `go.mod` declares zero `require` lines, so stdlib-first is satisfied without argument. `govulncheck ./...` ran (installed at `/home/kesokaj/.go/bin/govulncheck`) — output in finding 5. `npm audit` not applicable: no `package.json` in the repo (`npm` is on PATH, the ecosystem is not present).
- **GCP surface** — nothing: no IAM binding, service account, Cloud Run flag, Secret Manager reference or key file anywhere in the diff.
- **Supporting checks** — `gofmt -l .` clean, `go vet ./...` clean, `go test -race ./...` passes both packages. The added `TestAttributionNote` uses `t.TempDir()` and `t.Setenv(harness.ChildEnv, "1")`, so it neither rebuilds nor touches the real `$HOME`; it does run 8 full installs against the live checkout read-only.

## Not examined

- The **behaviour of Claude Code itself** given `"attribution": {"commit": ""}` — the note asserts this key suppresses the trailer; that claim was read, not verified against the CLI, and no schema for `settings.json` was consulted. If the key or its accepted values are wrong, the note is confidently wrong and every finding above is moot for its stated purpose.
- The `bool` arm of the type switch (`case bool: if !v`) — no evidence was gathered that Claude Code ever accepts a boolean there; it may be dead code or may be the real spelling. Untested either way.
- **`--vendor antigravity`** paths: the guard `vendor == harness.VendorAll || vendor == harness.VendorClaude` was read but only the default (`all`) vendor was exercised at runtime.
- The **`/dev/zero` (character device) case** was reasoned about — `os.ReadFile` sizes from `Stat`, gets 0, and grows unbounded — but deliberately **not run**, to avoid an OOM on the operator's box. The FIFO case was run and hung; the unbounded-memory case is inferred, not observed.
- The other eight modified files in the working tree (`agents/*.md`, `context/AGENTS.md`, `skills/*`) — out of scope by the dispatch, not read.
- No **install into the real `$HOME`** and no CLI smoke check; every run used `--home .agents/scratch/secreview-home`, which is also the home whose byte-echo and hang behaviour was measured.

## Outcome

- Finding 1 (medium) — fixed: `attributionOff` now `os.Stat`s first and judges only a regular file, reads it through `io.LimitReader(f, 1<<20)`, and every failure (missing → silent, unreadable, not regular, truncated, malformed) is an advisory note with `return nil`; a FIFO is never opened. Tests `a directory is advisory` and `a fifo is advisory and never opened` (the FIFO case would hang the old code under `go test`'s own timeout).
- Finding 2 (low) — fixed: the parse failure prints "is not valid JSON" and nothing from the error; the test `malformed json is advisory, no echo` also fails if "invalid character" or "unexpected end" ever reaches stdout.
- Finding 3 (low) — fixed: both notes print the path that was read (`filepath.Join(home, ".claude", "settings.json")`), as `pathHint` does; the tests assert on that path under the scratch home. For the operator's real home that path is `~/.claude/settings.json`, which is what the README and `context/AGENTS.md` § Git name.
- Also recorded, "parse failure suppresses the actual warning" — fixed: the advisory note carries the trailer warning too ("… could not be checked; Claude Code adds a Co-Authored-By trailer to every commit unless it sets …").
- `io/fs` shadowing (low) — moot: the import is gone; the code compares against `os.ErrNotExist`.
- Not examined item "whether the setting actually suppresses the trailer" — owed: a sandbox row `attribution-off` (a headless commit with the setting present in the sandbox home) runs after the Claude lane ends, recorded in `testing/2026-09-16-linkshelf-simulation.md`.
- Round 1 closed 2026-09-16; no second round owed (robustness on an advisory read; no mechanism finding).
