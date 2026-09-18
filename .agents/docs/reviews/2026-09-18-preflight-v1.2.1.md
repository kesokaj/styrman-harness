---
owner: main
date: 2026-09-18
---
# Release gate — v1.2.1     (2026-09-18)

Range reviewed: v1.2.0..86d57d8

## VERDICT: SHIP

No blocking finding stands. The code review returned APPROVE WITH NITS and security FINDINGS (three low); six findings in all, every one fixed before the tag, each with a test that fails without it, and the verification row re-taken at the fix commit. Both reviews were given the same range and both independently corrected an error in my dispatch text (see "What was not checked"). No project agent shadows a harness review agent.

## What this release is

A PATCH, and the reason it exists is a silent failure the operator found by installing v1.2.0 on a second machine. `argolis-devbox` carried an old unpacked bundle at `/opt/devbox/deploy/harness` recorded in its manifest; the release binary installed THAT tree's months-old content and reported `drift none` against it. It surfaced only because that machine's Claude Code quoted the `rm -r` rule instead of the fail-closed rule shipped in v1.2.0.

A release binary now installs the harness built into it and never a tree on the machine — not the working directory, not the recorded path, not one named with `--repo`/`$SH_AGENT_REPO`, which is refused rather than ignored (ADR-034, the operator: "the binary should only use the published git repo in github always and only"). A binary-only place also sheds the old kept feedback clone, which only the checkout path did before. Development builds are untouched: installing the working tree is how this harness is developed.

## Rows

| Row | Verdict | Artifact |
|---|---|---|
| Verification, v1.2.0..288e52d | PASS — 103 + 86, 0 fail, 0 skip; all four release-binary legs | `testing/2026-09-18-verify-preflight-v1.2.1.md` |
| Verification, after the gate's fixes | PASS — 104 + 87, 0 fail, 0 skip at 86d57d8; the four legs again, including the report regression | `testing/2026-09-18-verify-preflight-v1.2.1-fixes.md` |
| Code review | APPROVE WITH NITS — 1 medium, 2 low, all fixed; `## Outcome` written | `reviews/2026-09-18-review-preflight-v1.2.1-code.md` |
| Security | FINDINGS — 3 low, all fixed; trust boundary confirmed CLOSED | `reviews/2026-09-18-security-preflight-v1.2.1.md` |
| Docs lens | Not dispatched; the code review read the docs in range and produced the medium finding below | — |

## Findings, worst first

All six fixed in 86d57d8.

1. **[medium/high] `README.md:132` still described the old source order** and offered `--repo` unqualified — false for a release binary as of this release. My own miss: the changes entry claimed README needed nothing because I checked § Install without cloning and not the § Daily commands paragraph that documents the flags. The sentence now scopes the order to a build made from a checkout and states the release-binary rule.
2. **[low/high] an exported `$SH_AGENT_REPO` made `report` send a null record** on a release binary — a regression introduced by ADR-034 itself, an hour before the gate. `report` reads a source and never installs one, so it now takes the embedded files directly; a `--repo` passed to it gets a note, not a refusal. This is the shape BUG-077 was about, reappearing in the fix for BUG-077.
3. **[low/med] `RemoveKeptClone` asked git without refusing `GIT_DIR` and friends**, on a path that DELETES and that this release newly extends to binary-only places, while `SendReport` refuses them a few lines away for exactly that reason. Same `gitRedirects` loop now guards it.
4. **[low/high] the refusal named the path but not its origin.** With `$SH_AGENT_REPO` exported in a shell profile every command on the machine fails, and "drop it" does not say what to drop. It now reads `--repo=/path` or `$SH_AGENT_REPO=/path`.
5. **[low/high] the working-directory walk was pinned by no test** — pre-existing; deleting it kept the suite green. Now pinned from inside the tree and from a nested directory under it.
6. **[low/med] the trail commits internal identifiers** (hostnames, machine paths, a report body). No secret among them, which is checked every release, but publication is a separate decision from releasing — now written into `remote.md` § Releases where that decision would be made.

Two lower notes were answered rather than fixed, in the review's `## Outcome`: `adopt --repo` was already unreachable from a release binary before this change, and `shedKeptClone` running before the download is deliberate — what it sheds is dead weight, and deferring it would leave a machine holding other places' reports for another cycle.

## What was not checked

- **No docs-lens dispatch.** The code review read the docs in range and found the medium above; a separate lens would have added little. Recorded rather than skipped silently.
- **Both reviewers were handed a wrong premise by me**: my dispatch asked them to review two `report.go` message changes that actually shipped in v1.2.0 (`git diff v1.2.0..HEAD -- internal/harness/report.go` was empty of them). Each checked and said so independently. Nothing was reviewed on a false basis; the error was in my brief.
- **One of my own tests was too weak and the mutation run caught it**: the refusal test asserted the message named the release but not which flag, so the mutation that removed the origin passed. Strengthened, then it failed correctly. Four of five mutations failed first time.
- **Nothing was run on darwin.** `darwin/arm64` is built and published; the Mac legs for this release are owed after it lands, and the devbox reset is owed regardless.
- **Security ran no tests and no mutations** (the verifier's tier) and reasoned rather than measured finding 3; the main agent ran the five mutations.
- **`update` from v1.2.0 to v1.2.1 against the real release** cannot be tested before the tag exists. The cross-release path was exercised for real at v1.2.0 and is covered by tests here.

## Version

v1.2.1 — PATCH. Everything in it fixes released behaviour; no installed place must react, and the source-order change only makes a release binary do what its name says. The operator was offered MINOR and kept PATCH.
