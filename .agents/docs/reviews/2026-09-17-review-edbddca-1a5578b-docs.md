---
owner: reviewer
date: 2026-09-17
---
# Review — docs-lens edbddca..1a5578b (2026-09-17)

Range reviewed: edbddca..1a5578b (docs-lens; code of record at 1a5578b in internal/harness/{release,report,source,repo,apply}.go and cmd/sh-agent/main.go)

## Verdict: REQUEST CHANGES

One stale paragraph in the map's own supporting prose still teaches the OLD two-ways-in story (checkout + unpacked bundle) and omits the release binary entirely, directly contradicting REQUIREMENTS.md R8/R9, AGENTS.md rule 1, and ADR-033 D1 ("the unpacked bundle leaves the docs and the map, its code path stays") — all of which this same commit correctly updated everywhere else. Everything else checked (delivery.md, DECISIONS.md, both deploy guides, the release procedure's exact commands against release.go's checks, report.go's throwaway-clone behaviour, exit codes) is internally consistent and matches the code at 1a5578b.

## Findings

- **[high / confidence: high]** `README.md:135` — "`sh-agent` is the only way the harness reaches a machine (R8): a checkout kept current with `update`, or an unpacked copy of the repository with no `.git`, installed with `go run ./cmd/sh-agent install` from that directory (it installs the running binary and `update` refuses it — unpack a newer copy to move forward). No project vendors a copy." This presents the unpacked bundle as one of "the" two ways in and never mentions the release binary at all. It contradicts: REQUIREMENTS.md R8 ("a checkout ... or a release's `sh-agent` binary with the harness built in"), R9's Install/Update bullets, AGENTS.md rule 1 ("a checkout ... or a release binary ... installed with `sh-agent install` and moved forward by `sh-agent update`"), and ADR-033 D1 ("The two ways in are a checkout and a release binary with the harness built in; the unpacked bundle leaves the docs and the map, its code path stays" — `DECISIONS.md:270`). Every other place in this same diff that used to make this claim was correctly rewritten (`README.md`'s own "How it works" § Install bullet, `delivery.md:33-35`, `remote.md:22-28`, `local.md:55`, `AGENTS.md` rule 1) — this one paragraph in "## The loop: edit on one machine, pull on another" was missed. Fix: rewrite to name the checkout and the release binary as the two ways in, and demote the plain-copy sentence to a fallback note the way `local.md:55` and `remote.md:28` already phrase it.

- **[low / confidence: low]** `README.md:115` (Daily commands, `sh-agent install`/`sync` row) — describes only the checkout-rebuild and plain-copy-as-is behaviours of `install`, omitting the release-binary case (`main.go:130-138`, its own distinct message "release binary %s: installing the harness built into it and the running binary"). Not a contradiction — the release-binary install path is documented correctly elsewhere (§ Install without cloning) — but a reader consulting only the command reference table sees two of the three code branches. Optional: add a clause for the embedded-source case.

## Criteria (docs-lens interpretation)

1. Dead code — n/a to a docs-lens review; no stale generated artifacts noticed.
2. Flow tracing (doc-to-code) — verified `UpdateFromRelease`, `verifyRelease`, `checksumFor`, `RemoveKeptClone`, and the `update`/`install` branch in `main.go` against `README.md`, both deploy guides, and `delivery.md`: messages, exit codes (`ErrNewerRelease` → exit 2, `Behind`/dirty-checkout paths), the `SHA256SUMS` version-marker regex (`^#\s*styrman-harness\s+(vX.Y.Z)\s*$`) against `remote.md`'s `printf '# styrman-harness %s\n'` line, the signer identity/namespace (`sh-agent-release`, `ostling.simon@gmail.com`) against the release procedure's `ssh-keygen -Y sign/verify` commands, and the three build targets against `AssetName` — all consistent. `report.go`'s throwaway-clone-then-remove behaviour and `RemoveKeptClone`'s exact conditions match `local.md` § Report back and the table of failure messages there.
3. Standards conformance / describing-text rule — one violation (finding 1 above): a restatement of the "ways in" rule was not updated everywhere its subject nouns appear, exactly the grep the global context's describing-text rule calls for.
4. Directory layout — n/a (no new files in this diff outside the review artifact).
5. Simplicity — n/a to a docs-lens review.
6. Map/REQUIREMENTS/AGENTS agreement (class 3) — `README.md` § How it works, `REQUIREMENTS.md` R8/R9, and `AGENTS.md` rules 1-2 agree with each other after this diff, EXCEPT `README.md:135`, which disagrees with all three.

## Not examined
- `apply.go` was only spot-checked (grepped for `EmbeddedPrefix`/manifest recording, not read end to end) — its interaction with D24 ("a release binary's install REPLACES a checkout install of the same home") was not traced in code, only noted as an ADR-033a claim; no doc currently spells out that replacement behaviour for an operator, and I did not confirm whether one should.
- D21's empirical claim ("`-trimpath` keeps `-ldflags` out of the build info, measured with go1.27.1", i.e. `status`'s `releaseFromInfo` reading `Main.Version` from a binary built with `-ldflags "-X main.version=..."` and `-trimpath`) was read in code (`main.go:849-859`, `status` at `:348-351`) but not empirically re-verified by building a binary — plausible given Go's VCS-stamping behaviour, not independently confirmed here.
- `internal/harness/repo.go`'s `Reinstall`, `Inspect`, `Pull` were read for the plain-copy/checkout distinction but their test files were not opened.
- Did not read `cmd/sh-agent/main.go` lines outside those touching `update`/`install`/`status`/`buildVersion` (e.g. `report()`, `installFiles`, `resolveVendor` bodies) in full; relied on grep + targeted reads.
- Did not check `.agents/docs/plans/2026-09-17-release-binary.md` (the plan named by ADR-033a) against the shipped docs, only the DECISIONS.md addendum summarizing it.
- `context/AGENTS.md` and `skills/` (other than `preflight/SKILL.md`, the only skill this diff touched) were grepped for the stale phrases but not read end to end; no hits found for "unpacked bundle" / "tag delivers nothing" / kept-clone language there.

## Outcome

Read against the diff by the main agent, 2026-09-17. Both findings were real and both are fixed; a repo-wide sweep found one more of the same class.

**Fixed**

1. **[high] `README.md` § The loop still framed R8 as "a checkout … or an unpacked copy of the repository with no `.git`".** Correct finding: that paragraph contradicted R8, R9, `AGENTS.md` rule 1 and ADR-033 D1, all rewritten in the same commit. It now names the two ways in — a checkout, or a release binary installed with `sh-agent install` and moved forward by `sh-agent update` — and states the plain copy as the fallback for a machine that can reach neither, which is what ADR-033 D1 says happens to it.
2. **[low] the `install`/`sync` row of § Daily commands omitted the release-binary branch.** Fixed: the row now names all three sources (checkout rebuilds, release binary deploys what it carries and copies the running binary, plain copy deploys the running binary as is).
3. **Found by the sweep the finding prompted, not by the review:** `.agents/docs/deploy/local.md` § Report back's "No harness checkout here (an unpacked bundle, or nothing)" row described only the two old cases, so a release-binary place's record was undocumented — and that place is exactly where the new `harness_version` field is set. The row now distinguishes a release binary (`harness_commit` null, `harness_version` set, counts from the embedded files), a plain copy, and no tree at all.

This is the fourth time a rename left a stale restatement behind in this repository (BUG-003, BUG-007, BUG-027). The rule those recorded — grep the rule's SUBJECT NOUNS, not its old wording — is what found item 3: the sweep searched for "unpacked bundle", "bundle place", "tag delivers nothing", "tag changes nothing" and the kept-clone phrasings across every living doc, `context/`, `skills/`, `agents/`, `commands/` and `rules/`, and read each hit in place rather than trusting the count.

**Not disputed**

The review's "not examined" list (apply.go's D24 interaction, D21's build-info claim, full reads of `report()` and `installFiles`) is accurate. D21's claim — `-trimpath` keeps `-ldflags` out of the build info — was measured this session with go1.27.1 (`go version -m` on two builds, one with `-trimpath` and one without) and is recorded in the plan's D21 row.
