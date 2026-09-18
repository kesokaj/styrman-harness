---
owner: main
date: 2026-09-18
---
# The human install route names its version     (2026-09-18)

## Summary

After v1.2.2 was tagged, the three by-hand install chains — README § Install without cloning (`gh` route and browser route) and `remote.md` § First install — were bound to a version the reader states, the way `sh-agent update` is bound to the tag it asks for. Until now a genuinely signed OLDER release, served whole, passed every by-hand check: the signature proves the operator signed it once, not that it is the release you meant. The browser chain also moved into a scratch directory, lifts the macOS quarantine before it runs the binary, and checks all three files are present before moving any. Docs only; nothing in a release artifact.

## Type

fix (docs) — the gate report for v1.2.2 accepted this for the tag and promised it in the commit after the push; this is that commit.

## What changed

- Each chain starts with `V=vX.Y.Z`, downloads `"$V"`, and after the signature check requires the SIGNED checksum file to name that version — `grep -qxF "# styrman-harness $V" SHA256SUMS` — and the binary to be it — `test "$(./<asset> version)" = "sh-agent $V"`. The guarantee, stated in the comment: the version you named, or nothing — never "the newest". `V` comes from the operator who cut it or a checkout's changelog, and the comment says why not the Releases page (same server as the assets).
- README's browser chain: `ls` of the three names before the `mv` into `mktemp -d`; `{ xattr -d com.apple.quarantine <asset> 2>/dev/null || true; }` after `chmod` and before `version`, so a browser download on a Mac can run at all; the standalone `xattr` paragraph folded into it.
- `remote.md`'s comment states only what the chain does, including that `update` also refuses a not-newer tag, which no first install can mirror.
- `.agents/docs/testing/2026-09-17-release-binary-e2e.md` § Added 2026-09-18 — v1.2.2: the published release consumed from a fresh download.

## Verification

Against the real v1.2.2 bundle: `V=v1.2.2` → both lines exit 0; `V=v1.2.3` (the downgrade shape: a genuine bundle that is not the version named) → both exit 1; a marker rewritten to `v1.2.1` → exit 1; `V=v1X2X2` (what BRE dots would have matched) → exit 1; the literal placeholder and an empty `V` → exit 1. `bash -n` clean on all three chains. Verifier logs `testing/2026-09-18-verify-post-v1.2.2-docs.md` (the first, bundle-only binding — superseded) and `…-docs-2.md` (this binding, must-pass and both must-fail directions); the `-F` change after that log was re-run by hand as above.

## Definition of Done

Three rounds on the human route, the middle one put to the operator (`.agents/scratch/human-route-downgrade-options.md`, "Go with 1"):

- Round 1 — the gate's security report asked for a version binding; my first form bound the binary to the bundle beside it. `security` (`reviews/2026-09-18-security-f3ec5dc-worktree.md`): 1 medium — that binding does not stop a signed older bundle; 4 low. `reviewer --docs` APPROVE. `verifier` PASS on that form.
- Round 2 — option 1. `verifier` PASS in both directions; `security` (`…-worktree-2.md`): mechanism right, 1 medium (the Releases page named as a source for `V`) + 2 low (a hardcoded `V`; BRE dots), all fixed, `## Outcome`; `reviewer --docs` (`…-review-worktree-docs-2.md`) APPROVE WITH NITS: the same BRE point (fixed) and a `guide.md` tier question (left to `documentation`).
- The code-review row: under 50 lines, read by the main agent — the diff is three shell chains and their comments.
- Every finding has an `## Outcome` line in its report.

## Left, with owners

- `testing/guide.md`: no tier for the manual install chain, and § 5 does not yet ask for the unqualified `sh-agent version` or "no `shadowed` line" — `documentation`, in the handoff's Next Steps.
- The browser chain has never run on a Mac (`shasum --ignore-missing`, Gatekeeper, `xattr`) — the operator, one run on the work Mac.
