---
owner: verifier
date: 2026-09-18
---
# Verification — post-v1.2.2 docs (round 2)   (2026-09-18)

Range reviewed: f3ec5dc..worktree (uncommitted edits to README.md, .agents/docs/deploy/remote.md, .agents/docs/testing/2026-09-17-release-binary-e2e.md)

## Status: PASS

Scope: Go checks (per guide.md § 1, no .go file in diff but run anyway per dispatch); the three install-chain V=/grep-marker/version checks exercised against the real published v1.2.2 release in `.agents/scratch/verify-post-v122-2/`; must-pass and must-fail (downgrade shape, tampered marker) cases; `bash -n` syntax check of each extracted chain. Real-home smoke checks and the browser chain's `xattr` step (cannot run on Linux) were NOT run — see Not examined. `testing/guide.md` already exists (dated 2026-09-17); not modified.

## Commands

- `gofmt -l cmd internal` → prints nothing, exit 0.
- `go vet ./...` → exit 0.
- `go test -race -v -count=1 ./...` → exit 0. Output at `.agents/scratch/verify-post-v122-2/go-test.txt`.
  counts: pass 215, fail 0, skip 0 at f3ec5dc+dirty for go test -race -v -count=1 ./...
  vs previous: .agents/docs/testing/2026-09-18-verify-post-v1.2.2-docs.md, pass 215→215, skip 0→0 (unchanged; no Go source in this diff).

### Real-release download (the one permitted network step)

- `gh release download v1.2.2 --repo kesokaj/styrman-harness --clobber --pattern sh-agent-linux-amd64 --pattern SHA256SUMS --pattern SHA256SUMS.sig --dir .agents/scratch/verify-post-v122-2/dl` → exit 0. Fetched `sh-agent-linux-amd64` (7052750 bytes), `SHA256SUMS`, `SHA256SUMS.sig`.
- `SHA256SUMS` first line: `# styrman-harness v1.2.2`.
- `ssh-keygen -Y verify -f allowed_signers -I ostling.simon@gmail.com -n sh-agent-release -s SHA256SUMS.sig < SHA256SUMS` → "Good \"sh-agent-release\" signature ... key SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU", exit 0.
- `grep -v '^#' SHA256SUMS | sha256sum -c - --ignore-missing` → "sh-agent-linux-amd64: OK", exit 0.
- `chmod +x sh-agent-linux-amd64` → exit 0.

### 1. Must-pass (V=v1.2.2, genuine bundle)

- `V=v1.2.2; grep -qx "# styrman-harness $V" SHA256SUMS` → exit 0.
- `V=v1.2.2; test "$(./sh-agent-linux-amd64 version)" = "sh-agent $V"` → exit 0 (`./sh-agent-linux-amd64 version` printed `sh-agent v1.2.2`).

### 2. Must-fail, downgrade shape (V=v1.2.3, a version this bundle is not)

- `V=v1.2.3; grep -qx "# styrman-harness $V" SHA256SUMS` → exit 1.
- `V=v1.2.3; test "$(./sh-agent-linux-amd64 version)" = "sh-agent $V"` → exit 1.
- Confirms: the genuine, internally consistent v1.2.2 bundle is refused when asked for a version it is not.

### 3. Must-fail, tampered marker (copy of SHA256SUMS with first line `# styrman-harness v1.2.1`, V=v1.2.2)

- Made `SHA256SUMS.tampered` via `cp` + `sed -i '1s/.*/# styrman-harness v1.2.1/'`; verified first line reads `# styrman-harness v1.2.1`.
- `V=v1.2.2; grep -qx "# styrman-harness $V" SHA256SUMS.tampered` → exit 1.

### 4. `bash -n` over each chain

Chains extracted verbatim (from the `V=` line through the `rm -r "$D"` line) into `.agents/scratch/verify-post-v122-2/chains/`:
- `readme-gh.sh` (README.md:91-101, comment line above `V=` stripped, not part of the chain) → `bash -n` exit 0.
- `readme-browser.sh` (README.md:126-135) → `bash -n` exit 0.
- `remote-first-install.sh` (.agents/docs/deploy/remote.md:28-38) → `bash -n` exit 0.

## Root cause (FAIL only)

N/A — Status PASS.

## Coverage of changed code

This is a docs-only diff (README.md, .agents/docs/deploy/remote.md, .agents/docs/testing/2026-09-17-release-binary-e2e.md); no Go source changed, no test-coverage gap applicable. No added `//nolint`, `t.Skip`, `.skip(`/`.only(`, or similar suppressions found in the diff (none expected in prose/markdown files).

## Not examined

- Real-home smoke checks (`.agents/docs/testing/2026-09-03-verify-real-home.md`) — out of scope per dispatch (real $HOME is off-limits).
- The `go install` route (README.md:112-115) — unchanged by this diff, not re-verified.
- The browser chain's `ls` pre-check and `xattr -d com.apple.quarantine ... || true` group — this is a Linux box; `xattr` is a macOS-only step and cannot be exercised for real here. The chain's syntax was checked with `bash -n` only.
- `.agents/docs/testing/2026-09-17-release-binary-e2e.md` itself — diffed but not independently re-run beyond confirming it is consistent with the chains tested above (it documents the same V=/grep/test pattern); not read line-by-line against every claim in this pass.
- `darwin-arm64` and `linux-arm64` assets — only `sh-agent-linux-amd64` was downloaded and exercised, per the dispatch's explicit instruction.
- `sh-agent update`'s own internal downgrade/marker logic (`internal/harness` Go code implementing the same checks programmatically) — not re-tested here since no Go source changed; covered by the go test -race run's existing 215 passing tests.

## Environment

Linux 6.12.107+deb13-amd64 (Debian trixie). `go version go1.27.1 linux/amd64`. `gh` present and authenticated (`gh auth status`: logged in as kesokaj, github.com). `ssh-keygen`, `sha256sum` present. No Colima/Docker involved (no Dockerfile change).
