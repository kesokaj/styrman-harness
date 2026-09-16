---
owner: reviewer
date: 2026-09-16
---
# Review — attribution-write   (2026-09-16)

Range reviewed: working tree, `cmd/sh-agent/main.go` and `cmd/sh-agent/main_test.go`
Blobs read: `main.go` a5b537b78337b42df94f27df9091e8d465e8fa13, `main_test.go` d428d98860017d9b06e980400b003bdf26f3c311

## Verdict: APPROVE WITH NITS

The write path is sound for every realistic settings.json (small file, one flat top-level object): the `json.Unmarshal`-into-map gate filters comments, a second top-level value, and BOM-prefixed files into the "malformed, nothing written" branch before the textual-insertion code ever runs, `json.LastIndexByte`/`IndexByte` land on the true outer braces because that gate guarantees a single well-formed object, and the final `json.Valid(result)` is a real safety net, not decoration. `gofmt`, `go vet`, and `go test -race -count=1 ./cmd/sh-agent/` all pass clean. One concrete, verified bug: a top-level JSON `null` (a value D7 explicitly names) gets the wrong message and is untested — worth a one-line fix, not a blocker.

## Findings

- **[medium / confidence: high]** `cmd/sh-agent/main.go:509-535` (`attributionSet`) — a settings.json containing the single JSON value `null` prints `settings: <path> could not be edited safely (the result would not parse); add the key by hand` instead of the plan's D7 message `is not a JSON object; add the key by hand`. Verified live: `printf 'null\n' > .claude/settings.json && sh-agent install` → the "could not be edited safely" line, file untouched. Cause: `json.Unmarshal(data, &top)` does **not** error on `null` into a `map[string]json.RawMessage` (it leaves `top` nil with `err == nil`), so execution falls past the `json.Valid(data)` / "is not a JSON object" branch (which only fires on `err != nil`) into the textual-insertion path, where `trimmed` has no `{`/`}` and `closing <= opening` (both -1) triggers the generic "could not be edited safely" message. No data loss (nothing is written either way), but it is a plan (D7) conformance miss and the one non-object case with no test coverage — `TestAttributionSet`'s "not an object is left alone" case uses `["x"]`, not `null`. Fix: check `top == nil && err == nil` (or `bytes.Equal(bytes.TrimSpace(data), []byte("null"))`) alongside the `err != nil` branch and route it through the same "is not a JSON object" message.

- **[low / confidence: medium]** `cmd/sh-agent/main.go:592, 596, 601, 606` — `tmp.Close()` (cleanup call inside the write-error branch) and all four `os.Remove(tmpName)` calls on the error-cleanup paths discard their return values outright (not even `_ = `), which the go-review checklist's error-handling bullet flags. Low severity because these are best-effort cleanup after the primary error is already reported to the operator via `out` (a failed `Remove` just leaves a stray `.settings-*.json` temp file, not data corruption), but it is a real gap against the stated rule; worth `_ = os.Remove(tmpName)` for explicitness if the project wants the checklist literally satisfied.

- **[low / confidence: low]** `cmd/sh-agent/main.go:509` / `readSettings` — `readSettings`'s `io.LimitReader(f, 1<<20)` truncates silently (no error signals truncation). Before this diff that only fed the read-only "off" judgment; now the same truncated bytes can be re-serialized and written back via the temp-file/rename in `attributionSet`. For the narrow, unrealistic case of a settings.json >1 MiB whose top-level object happens to close exactly within the first MiB with further bytes beyond it, the write would silently drop everything past the truncation point. Real settings.json files are tiny, so this is not worth blocking on, but it's a behavior change in what the truncation feeds (read-only → potential overwrite source).

## Criteria

1. Dead code — pass. No orphaned functions; `readSettings` is a genuine shared helper used by both `attributionSet` and `attributionOff`, matching the plan's Open-item default.
2. Flow tracing — pass with the `null` gap above. `attributionSet` is wired into the install case exactly where the plan specifies (vendor all/claude, before `attributionHint`), `attributionOff`'s refactor onto `readSettings` preserves its old three-state contract (verified via `TestAttributionNote`'s directory/FIFO/malformed cases, still passing), and every new subtest can fail (checked by inspection: each asserts either a file-content invariant or a specific output line that the code must actively produce). The two must-fail demonstrations are plan-recorded (2026-09-16, custom_value/off_value and missing_file/other_keys) but not independently re-verified here (would require editing source, which I don't do).
3. Standards conformance — pass on gofmt/vet/test; one low-severity gap on unchecked `Close`/`Remove` returns on cleanup paths (see finding above). Doc-comments on `attributionSet`/`readSettings`/`attributionOff` describe the current behavior accurately.
4. Directory layout — pass. Both changes stay in `cmd/sh-agent/main.go` and its co-located `main_test.go`; no new files.
5. Simplicity / DIY-except-security — pass. Textual insertion instead of decode-reencode is justified in-comment (D3, avoids `encoding/json` key reordering); no new dependency; stdlib only (`bytes`, `encoding/json`, `os`, `io`).
6. Decoupling — pass. `readSettings` is a small, single-purpose helper at the seam both callers need; no cross-package reach-ins; `attributionSet` and `attributionOff` stay independently readable.

## Not examined

- `REQUIREMENTS.md`, `context/AGENTS.md`, `README.md`, `.agents/docs/architecture/DECISIONS.md`, and the deploy guides — out of scope per the dispatch (limited to `main.go`/`main_test.go`); the plan's Phase 1 doc tasks (R4/ADR-025, context/README/deploy wording) are still unticked in the plan file, consistent with not being in this diff.
- CRLF, BOM, and comment/second-top-level-value inputs were reasoned through against documented `encoding/json` semantics and one BOM/comment case is structurally identical to the already-tested "malformed" cases (both fail `Unmarshal` and `json.Valid` together), but none of the three was run live against the binary — only `null` and the dry-run/vendor-agy/no-write paths were empirically exercised.
- The plan's two must-fail demonstrations (Phase 1, row 4) were not independently reproduced by reverting code, per the reviewer's no-source-edits constraint; taken on the plan's own record.

## Outcome

- Finding 1 (medium) — fixed: a bare `null` (which `json.Unmarshal` leaves as a nil map without an error) now takes the "is not a JSON object; add the key by hand" branch; test `null is left alone`.
- Finding 2 (low) — fixed: the error paths call `abandonTemp`, which reports the write failure and names the temp file if even its removal fails; the one remaining bare `tmp.Close()` sits on a file discarded on the next line and says so in its comment.
- Finding 3 (low) — fixed: `attributionSet` refuses a file larger than 1 MiB before reading ("is larger than 1 MiB, more than this tool will rewrite; add the key by hand"), so the bounded reader can never feed a truncated write; test `a file over the read bound is left alone`.
- Added after the verifier's flag, before this review returned: a symlinked settings.json is never written through (`os.Lstat`, "is a symlink; add the key to its target by hand"); test shown to fail without the guard.
- Round 1 closed 2026-09-16; no second round owed (no mechanism finding).
