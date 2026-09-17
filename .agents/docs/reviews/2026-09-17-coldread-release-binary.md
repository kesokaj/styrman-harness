---
owner: reviewer
date: 2026-09-17
---
# Cold-reader check — release-binary plan (2026-09-17)

Range reviewed: `.agents/docs/plans/2026-09-17-release-binary.md` (plan document, not a diff), plus the code surfaces it names.

## Round 1

### (1) Gaps

- **G1 (BLOCKING-shaped, high) — `go:embed` is compiled into every build, not only release/`go install` builds, but the Failure-modes table and the "cannot happen" claims assume otherwise.** Interfaces § Failure modes row 1 says a missing-source case "cannot happen for a release or `go install` binary; a dev build prints today's `harness repo not found …`". But `embed.go`'s `//go:embed all:context …` at the module root embeds the five directories into *any* binary built from this module, dev included — there is no build tag gating it. Given the plan's own Source-selection order (`--repo` → env → cwd/parent → recorded repo → embedded, install-release-binary.md:50), a `go run ./cmd/sh-agent status` invoked from outside any checkout with no recorded state would silently fall through to the embedded copy in a dev build too — not error as the table claims — and that embedded copy is whatever was on disk when the binary was last compiled, which can be stale relative to the live checkout a developer is iterating on. The plan does not decide whether the embedded fallback should be disabled for dev/untagged builds (e.g. gated on `buildVersion() != ""`) or accepted as dev behaviour. Two implementers will build different things here, and the difference changes what `sh-agent status`/`install` do during ordinary local testing (`.agents/scratch/x`, per this repo's own testing convention). Cite: plan lines 48–52, 75; `cmd/sh-agent/main.go:100-116` (today's unconditional `FindRepo` hard-error, which the embedded fallback removes for every build, not just release ones).

- **G2 (BLOCKING, high) — the `install`/`sync` rebuild-vs-copy branch (`cmd/sh-agent/main.go:117-158`) is not re-specified for `Dir == ""`.** Today, after `FindRepo` returns a directory, `install` calls `knownCheckout(*repoFlag, home, repo)` (main.go:123, :420) and computes `own := harness.IsCheckout(repo)` (main.go:137) before deciding to rebuild via `go build` or install the running binary as-is. `IsCheckout(repo string)` (`internal/harness/repo.go:79`) does `os.Lstat(filepath.Join(repo, ".git"))` — called with `repo == ""` (the plan's sentinel for the embedded source, "Source" ¶, line 48) this resolves to `Lstat(".git")` relative to the process's *current working directory*, not "no checkout." Likewise `knownCheckout` reads `st.Repo` and calls `samePath(st.Repo, repo)` with the same empty-string ambiguity. The plan's "Source selection" and "Source" paragraphs describe which *readers* take a `Source`, but never say how this specific control-flow block — which decides "rebuild from source" vs "install running binary as-is" — treats an embedded source. The natural answer (treat it like `own=false`, "plain copy, install as-is," since embedded is exactly a running binary that already contains its own vendor data) is plausible but is not decided anywhere in the plan, and getting it wrong via the existing `IsCheckout("")` cwd-relative bug is a concrete latent defect, not a style question. Cite: plan §Interfaces "Source"/"Source selection" (lines 48–52), §Surfaces touched (line 33, 36); `cmd/sh-agent/main.go:117-158`, `internal/harness/repo.go:79-82,420-432`.

- **G3 (gap, high) — R8 restates the exact two-channel model the plan changes, but R8 is absent from Scope, Surfaces touched, and the R9 replacement text, and the plan does not decide what happens to it.** R8 ("One delivery channel, and it is `sh-agent`") reads: "a checkout the operator owns … or an unpacked bundle installed with `sh-agent install`. There is no other channel" (`REQUIREMENTS.md:35`). The plan's approved-pending R9 replacement Install bullet (plan line 43) drops the bundle wording entirely and describes only "clones … or downloads a release's `sh-agent` binary" — a genuinely different two-way split (checkout vs. release binary) than R8's (checkout vs. bundle), with the bundle's code path kept but undocumented (scope item 10, plan line 22: "The unpacked bundle leaves the docs and the map; its code path stays"). After this change, R8 and R9 describe two different, mutually inconsistent worlds unless R8 is also edited in the same commit — and REQUIREMENTS.md is immutable without explicit operator approval per project rule. Review-standards criterion 3 calls this out by name: "every other restatement of [a changed rule] is changed with it." Neither Scope, Surfaces touched, nor Phase 1/5 name an R8 edit or ask the operator to approve one. Cite: `REQUIREMENTS.md:33-35`; plan line 22, 26, 37, 41-46.

- **G4 (gap, medium) — the Source-selection fallback order does not say whether an explicitly-named but invalid `--repo`/`$SH_AGENT_REPO` still hard-errors, or falls through to the next candidate.** Today's `FindRepo` comment is explicit: "A path the operator named is used or refused, never silently replaced by one the walk finds" (`internal/harness/repo.go:14-18`), and the code returns an error immediately when `explicit`/env is set but `!IsRepo(named)` (repo.go:24-26). The plan's flat list — "`--repo` → `$SH_AGENT_REPO` → the working directory or a parent → the recorded repo when it is a directory → the embedded files" (line 50) — reads as one uniform fallthrough and does not say the first two candidates keep their hard-error-on-mismatch behaviour while the rest degrade gracefully. A typo'd `--repo` silently landing on embedded (possibly stale, per G1) instead of erroring is exactly the failure mode the original code was written to prevent.

- **G5 (gap, low, documentation clause missing on a command/flag change) — Phase 5's preflight task says "`skills/preflight/SKILL.md` Step 6 … cross-compile the three binaries, `SHA256SUMS`, `SHA256SUMS.sig`, `gh release create vX.Y.Z`" but does not name where in Step 6's existing six numbered sub-steps this lands, nor how `ssh-keygen -Y sign` for `SHA256SUMS.sig` is meant to run non-interactively inside an agent-run `/preflight` (today's Step 6.3 already signs the git tag with the operator's key on a machine holding it; the plan is silent on whether checksum-signing reuses that same already-unlocked-key moment or needs its own). This is exactly a "task changing a … deploy step with no docs clause spelling out the mechanics" case; the plan does name the doc file to edit (satisfying the letter of the "docs clause" checklist) but not the ordering/interaction, which a first-pass implementer restructuring Step 6 would have to guess at.

- Already self-flagged by the plan (noted here only for completeness, not double-counted as a gap): the "Open" line on `go install …@latest` before `v1.1.0` exists states its own default; D5 and the R9/D1 text are explicitly gated on operator approval in Phase 1 task 1.

### (2) Scope

Reviewed against `## Why` (single-binary download to avoid a `.agents/docs` trail beside projects; `report`'s full clone leaking every place's history to every other place). Nothing found that isn't traceable to a D-row or an operator-confirmed scope item:
- Three-platform cross-compile (D10) — operator scope item 5.
- Source selection extending to `status`/`diff`/`list`/`render`, not just `install` — necessary consequence (a binary-only place has no checkout, so these commands would otherwise always error), not an added shape.
- `version` command's version-string logic (D6) — implementer decision needed to make `status`/`update` work at all for a release binary, not a new feature.
- No extra endpoint, no second shape for the same operation, no configuration invented for a value that doesn't vary was found.

### (3) Readiness

Blocking questions (a first-pass implementer must ask, or two implementers diverge):
1. Does the embedded-source fallback apply to dev/untagged builds too, and if so is that acceptable, or must it be disabled for builds without an injected version (G1)?
2. How does the `install`/`sync` rebuild-vs-plain-copy branch (`own`, `knownCheckout`, `IsCheckout`) treat `Source{Dir:""}` — is it forced down the "plain copy, install as-is" path, and is `IsCheckout`/`knownCheckout`/`samePath`'s `repo string` signature changed to avoid the empty-string/cwd ambiguity (G2)?
3. Does R8 get edited in the same commit as R9, and has the operator approved that text too (G3)? Today's R8 and the plan's proposed R9 describe different channel models.
4. Does an explicit but wrong `--repo`/`$SH_AGENT_REPO` still hard-error, or now silently fall through toward embedded (G4)?

Non-blocking (reasonable default available, or ordinary implementer discretion):
- Exact placement of the binary-attach step inside preflight Step 6, and whether checksum-signing shares the tag-signing key moment (G5) — a competent implementer picks the obvious slot (after tag push, before/alongside "Deploy via the project's deploy guide") and documents it.
- Whether the compiled-in allowed-signer line is a `const` or a package `var` overridable by white-box tests for the fake-`gh` integration tests — ordinary Go test-design discretion.
- `go install …@latest` before `v1.1.0` — plan already states its own default.

### Verified against source (not just recon)

- `TestHarnessMap` **does exist** (`cmd/sh-agent/main_test.go:987`), contrary to the recon brief's claim that it doesn't (`2026-09-17-recon-release-binary.md` line 98, "second pass" findings) — recon is stale/wrong on this point but the plan's own citation (line 36, dispatch's line-987 pointer) is correct, so this did not propagate into a plan defect. Flagged only so an implementer reading the recon brief instead of the plan doesn't chase a nonexistent problem.
- 15 `.gitkeep` files confirmed on disk (`find … -name .gitkeep | wc -l` → 15), matching the plan's Verification bullet (line 89) exactly — recon's "13" was wrong, plan's own number is right.
- `report` command already bypasses `FindRepo` today with a comment anticipating exactly this problem (`cmd/sh-agent/main.go:106-109`) — consistent with the plan's design, no issue.
- `REQUIREMENTS.md` R4/R8/R9 read in full; R8 and R9 both carry "(Approved by the operator 2026-09-17 …)" — same day as this plan, both existing texts, R9's *replacement* still pending approval per Phase 1.

## Not examined

- `internal/harness/report.go`/`report_test.go` bodies (`ResolveFeedback`, `verifyFeedback`, `commitAndPush`, the four pinned tests) — read only via the recon brief's summary, not the source directly; the throwaway-clone rewrite's interaction with existing test assertions is taken on the recon's word.
- `internal/harness/state.go`, `apply.go`, `frontmatter.go`, `derive.go`, `adopt.go` — not read; `adopt`'s "requires a checkout and says so" (plan line 50) not verified against `adopt.go`'s actual error paths.
- `.agents/docs/deploy/remote.md` § Releases full text (only grepped, not read in full) and § The four places — not verified against the plan's D-rows in detail.
- `cmd/sh-agent/main_test.go` and `main_test.go`'s fixtures beyond `TestHarnessMap` — not read.
- ADR-030/ADR-031/ADR-029 full texts — not read; only their conclusions as summarized in REQUIREMENTS.md/AGENTS.md.
- `internal/harness/harness_test.go` — not opened.

design needs revision (G1, G2, G3, G4 are gaps the plan does not decide)

implementation not ready — 4 BLOCKING questions listed above (G1–G4)

## Round 2

Fresh read of the plan (2026-09-17-release-binary.md, current version with D1–D16) before reading Round 1 again.

### My own read (before re-reading Round 1)

#### (1) Gaps

- **G6 (BLOCKING, high) — `status` and `update` are not re-specified for `Source{Dir:""}` the way `install` was (D13 covers install only).** `cmd/sh-agent/main.go` calls `harness.FindRepo` once, shared by `install`, `status`, `diff`, `update`, `list`, `render` (main.go: the single `repo, err := harness.FindRepo(...)` before the `switch cmd`). D14/D15 make `FindRepo` (or its `Source`-returning successor) hand back the embedded source instead of erroring for a release binary with no tree — but two existing functions built on the old two-state world ("real checkout" or "plain copy, git refused") are never revisited:
  - `status(repo, home, vendorFlag, out)` (main.go:262) calls `harness.Inspect(repo)` unconditionally and, when `info.Commit == ""` (true for both a plain copy *and* an embedded source), prints `"repo %s (plain copy, not its own git checkout; update is refused)\n"` (main.go:272). For an embedded install this is doubly wrong: the repo field is empty/`""` so the line reads `repo        (plain copy...)`, and it says *"update is refused"* — directly contradicting the plan's own Update interface bullet (line 44) and its E2E criterion (line 98, `sh-agent update` on the embedded place prints `up to date (v1.1.0)`, i.e. update plainly works there). The plan's only status-side change (line 56, `"source embedded vX.Y.Z"`) adds a line but never says the old plain-copy line is suppressed or reworded for this case.
  - `update(repo, home, checkOnly, ...)` (main.go:364) is pure git: `harness.Behind(repo)`, `harness.Inspect(repo)`, `harness.Pull(repo, out)`, `harness.Reinstall(repo, args, out)`. None of these are meant to run for a binary-only place — the plan's new release-download flow (lines 60–64) is a wholly separate procedure — but nothing in `cmd/sh-agent/main.go`'s dispatch (`case "update":`, main.go ~"if err := knownCheckout(...); return update(repo, home, ...)") is named as the branch point that sends an embedded/binary-only place to the new flow instead of into this function. Falling through unchanged, `harness.Behind("")` → `IsCheckout("")` → `Lstat(".git")` relative to the process's cwd — the exact ambiguity round 1's G2 flagged for `install`, reproduced here for `update` because D13's fix text says "Binary handling **on install**" only.
  - The "Source selection" paragraph itself (line 52) lists the commands it governs as "`install`, `status`, `diff`, `list`, `render`" — **`update` is not in that list**, even though `update` is exactly the command that most needs to know "checkout vs. embedded" to decide which of two completely different procedures to run. Either the list is an oversight or `update`'s branch is decided by some other unstated mechanism.
  Cite: plan lines 44, 52, 56, 60–64, 98; `cmd/sh-agent/main.go:262-275` (`status`), `:364-408` (`update`), `:117-119` (shared `FindRepo` call).

- **G7 (gap, medium) — the plan never states `FindRepo`'s new return type**, only the ordered list of candidates ("Source selection", line 52) and the `Source{FS, Dir}` shape ("Source" ¶, line 50). `internal/harness/repo.go:20` today is `func FindRepo(explicit, home string) (string, error)`. Surfaces touched (line 34) lists `repo.go: FindRepo, IsCheckout, Inspect, Reinstall` as touched but not what `FindRepo` returns afterward (a `Source`? a `(string, bool, error)` triple where the bool means "embedded"?) or whether `IsCheckout`/`Inspect`/`Reinstall`/`Behind`/`Pull` keep taking `repo string` (untouched for the checkout path) while a new, separate embedded branch in `main.go` short-circuits around all of them. This is the same class of gap as round 1's G2, generalized: G2 asked "how does the install rebuild-vs-copy branch treat `Dir==""`", D13 answered it for install specifically, but the underlying question — "what type does the plumbing between `FindRepo` and every caller carry now, and which existing `repo string` functions must never see `Dir==""`" — is still open for the other five(six) call sites.

#### (2) Scope

Nothing found beyond what a D-row or the operator's own scope list traces to; same conclusion as Round 1 (three-platform cross-compile → D10/scope 5; source selection spanning multiple commands is a necessary consequence, not an added shape; `version` logic → D6, needed for `status`/`update` to function at all).

#### (3) Readiness

Blocking (in addition to whatever Round 1 still has open):
1. G6 — what does `status` print for an embedded install in place of "(plain copy...; update is refused)", and what does `update`'s dispatch in `main.go` look like for `Source.Dir==""` — a new `if` before calling the existing git-based `update()`, or does `update()` itself grow a branch? Two implementers asked to "make update work for a binary-only place" independently will not produce the same `status`/`update` code without this decided.
2. G7 — `FindRepo`'s post-change signature/return shape, and confirmation that `IsCheckout`, `Inspect`, `Behind`, `Pull` are never called with an empty/embedded `Dir` from any of the five-or-six command paths, not just `install`.

Non-blocking: none new beyond Round 1's own non-blocking list (G5-equivalent, defaults stated).

design needs revision — G6, G7 are undecided
implementation not ready — 2 new BLOCKING questions (G6, G7), on top of Round 1's set until verified below

---

### Comparing against Round 1's four blocking items

Read Round 1's report after forming the above. Round 1's four blocking items map to the plan's D13–D16:

1. **G2 (install rebuild-vs-copy branch, `IsCheckout("")` ambiguity) — RESOLVED.** D13 and the new "Binary handling on install" paragraph (plan line 54) directly answer it: embedded → copy, never rebuild, `IsCheckout` never called with `""`.
2. **G1 (embedded fallback scope — dev builds vs release only) — RESOLVED.** D14 and "Source selection" (line 52) now gate the embedded fallback to a release binary carrying a `vX.Y.Z` version; a dev build with no tree found keeps today's hard error.
3. **G3 (R8 restates the two-channel model R9 changes, left unedited) — RESOLVED at the design level, pending execution.** The plan now carries a proposed exact R8 replacement (line 48) and a dedicated Phase 1 task ("Operator approves the R8 sentence," line 106, currently unchecked). Not yet approved, but no longer an undecided gap — it is a tracked, gated step.
4. **G4 (explicit but wrong `--repo`/`$SH_AGENT_REPO` — hard error or fall-through?) — RESOLVED.** D15 and "Source selection" (line 52, "never a fall-through (D15)") restore today's hard-error behaviour explicitly.

All four of Round 1's blocking items are resolved in this revision. However, the fix for items 1–2 (the `install`-specific "Binary handling" and "Source selection" paragraphs) was scoped narrowly to `install`, and in closing that gap the plan left the structurally identical problem open for `status` and `update` (G6/G7 above) — the same class of defect, one command over.

### Not examined (Round 2)

- `internal/harness/plan.go`, `inventory.go`, `render.go`, `frontmatter.go`, `derive.go`, `apply.go`, `state.go` bodies — not read; only signatures grepped. The claim that these five/six readers "move from `repo string` to a `Source`" cleanly (plan line 33) is taken on the plan's word, not verified line-by-line against e.g. `checkManagedTree(repo string)` (inventory.go:106) or `frontmatter.go`'s file-reading calls.
- `internal/harness/report.go`, `report_test.go` — not read in this round; Round 1 already flagged this as unexamined and it remains so here.
- `.agents/docs/deploy/remote.md` § Releases, § The four places — not read in full.
- `skills/preflight/SKILL.md` Step 6 — not read; G5 from Round 1 (preflight step ordering/signing-key-moment) not independently re-checked.
- `cmd/sh-agent/main_test.go`, `main_test.go` fixtures — not read beyond what Round 1 already checked (`TestHarnessMap` existence).
- ADR-029/030/031, `.agents/docs/specs/sh-agent/delivery.md` — not read.

## Round 3

Fresh read of the plan (2026-09-17-release-binary.md, current version with D1–D17, Phase 1 tasks 1–2 checked) plus `REQUIREMENTS.md`, `AGENTS.md`, `internal/harness/{repo,plan,inventory,render,state,report}.go`, `cmd/sh-agent/main.go`, `.agents/docs/deploy/remote.md` § Releases, `DECISIONS.md` ADR-030/ADR-031 — before re-reading Rounds 1–2.

### My own read (before re-reading Rounds 1–2)

#### (1) Gaps

- **G8 (BLOCKING, high) — the plan's own R9 replacement makes a released tag deliver binaries to binary-only places, but two other restatements of "a tag delivers nothing / every place follows `main`" are left unedited and unnamed.** `AGENTS.md` § How the harness works rule 2 (`AGENTS.md:13`): *"A tag delivers nothing. `/preflight vX.Y.Z` makes a signed release marker; every place follows `main`."* `.agents/docs/deploy/remote.md` § Releases, first sentence: *"It is a MARKER naming a version, not a way the harness reaches a machine: every place follows `main` with `sh-agent update`, and a tag changes nothing for any of them (ADR-030 D2, withdrawing ADR-026's tag channel)."* `DECISIONS.md` ADR-031 D3 lists "a tag delivers nothing" as one of AGENTS.md's five strict rules. All three are contradicted by the plan's own proposed R9 Release bullet (plan Interfaces: *"publishes that tag's binaries with checksums signed by the same key. Checkout places follow `main`; binary-only places update to the newest release."*) — for a binary-only place the tag is now exactly how the harness reaches the machine; it is no longer true that "a tag changes nothing for any of them." Surfaces touched names only `AGENTS.md` § How the harness works **rule 1** and `remote.md` § Releases (for the cross-compile/signing mechanics), never rule 2's sentence, never ADR-031 (amend or supersede?), never `remote.md`'s opening sentence. Review-standards criterion 3 ("every other restatement of a changed rule is changed with it — grep its subject nouns") and the global context's describing-text rule both name this exact failure mode. Cite: `AGENTS.md:13`; `.agents/docs/deploy/remote.md` § Releases opening line; `DECISIONS.md` ADR-031 (D3); plan § Interfaces "Release" bullet, § Surfaces touched.
- This is a design-level gap, not an implementation one: it needs a decision (does rule 2 change to something like "a tag delivers a release; a checkout follows `main`, a binary-only place downloads its assets", is ADR-031 amended or superseded by ADR-033, does `remote.md`'s opening sentence get reworded) before Phase 5's docs task can be executed without a second round of doc drift.

- **G9 (gap, low, mostly resolved but one prose omission survives from Round 2's G7) — the "Source selection" paragraph's command list still omits `update`.** Interfaces § "Source selection" reads: *"Source selection (`install`, `status`, `diff`, `list`, `render`): …"* — `update` is not named, exactly as Round 2 flagged for the old `FindRepo`. Functionally this is now covered: the "Per command with an embedded source" table explicitly gives `update`'s behaviour for `src.Dir == ""` (the release flow), so an implementer reading the whole Interfaces section will not miss it. Flagged only because the prose list itself still reads as if `update` takes a different path than `ResolveSource`, which is not true per the table — a wording gap, not a design gap. Non-blocking: the table resolves it.

#### (2) Scope

Checked against `## Why` (a single downloadable binary so a session doesn't inherit `.agents/docs`; `report` no longer mirrors every place's history to every other place). Everything in Scope and the Phase list traces to a D-row or a quoted operator scope item (three platforms → D10/scope 5; throwaway clone → D8/"throwaway clone for 12"; old-clone cleanup → D9/scope 12; the embedded-source reach into `status`/`diff`/`list`/`render`/`adopt` → a necessary consequence of a binary-only place having no checkout, not an added shape). No extra endpoint, no second shape for the same operation, no configuration invented for a value that never varies. Same conclusion as Rounds 1–2: nothing found.

#### (3) Readiness

Blocking:
1. G8 — is `AGENTS.md` rule 2, `remote.md` § Releases' opening sentence, and ADR-031 D3 edited/amended in the same R9 commit, and with what exact replacement text? Two implementers will otherwise ship code where a binary-only place's `update` downloads a release's binary while the project's own strict rules still say, in three places, that a tag delivers nothing.

Non-blocking:
- G9 — `update`'s omission from the "Source selection" prose list; the per-command table already answers it, so an implementer who reads the whole section is not stuck.
- Carried from Round 1 (still non-blocking, unchanged): exact insertion point for checksum-signing inside preflight Step 6 (G5) — a competent implementer places it right after tag verification, reusing the already-unlocked signing key, and documents it.

### Comparing against Rounds 1–2's blocking items

1. **G1 (embedded fallback scope, dev vs release builds) — RESOLVED.** D14 and "Source selection" gate the embedded fallback to a release binary whose version matches `^v\d+\.\d+\.\d+$`; a dev build with no tree keeps today's hard error.
2. **G2 (install rebuild-vs-copy branch, `IsCheckout("")` ambiguity) — RESOLVED.** D13 and "Binary handling on install" directly answer it: embedded → copy, never rebuild, `IsCheckout` never called with an empty path.
3. **G3 (R8 restates the channel model R9 changes, left unedited) — RESOLVED at the design level.** The plan carries a proposed exact R8 replacement and a Phase 1 approval task, now checked `[x]` with the operator's quoted approval.
4. **G4 (explicit but wrong `--repo`/`$SH_AGENT_REPO` — hard error or fall-through) — RESOLVED.** D15 and "Source selection" ("never a fall-through (D15)") keep today's hard-error behaviour.
5. **G5 (preflight Step 6 insertion point / signing-key moment) — still open, still non-blocking**, as in Round 1: a reasonable default is available and stated by the reviewer, not the plan, but it was never blocking.
6. **G6 (`status`/`update` not re-specified for `Source{Dir:""}`) — RESOLVED.** The "Per command with an embedded source" table now gives `status`'s replacement line and `update`'s branch to the release flow explicitly.
7. **G7 (`FindRepo`'s post-change return type/signature undefined) — RESOLVED.** `ResolveSource(explicit, home, version string) (Source, error)` is now the stated signature, with git helpers "called only with `src.Dir != \"\"`" and "with an embedded source none is called (D17)".

All seven of Rounds 1–2's blocking items are resolved in this revision. Round 3 finds one new blocking item (G8), one class apart from G1–G7: those were all about the *code* not deciding what an embedded source does; G8 is about the *project's own rule text* (`AGENTS.md`, `remote.md`, ADR-031) asserting something the plan's own design now makes false, in places the plan's Surfaces-touched list does not name.

### Not examined (Round 3)

- `internal/harness/frontmatter.go`, `derive.go`, `apply.go` bodies — not read; only referenced from `plan.go`/`inventory.go`. The claim that all readers "move from `repo string` to a `Source`" cleanly is still not verified line-by-line against these three files (carried forward from Round 2).
- `internal/harness/report.go` was read this round (lines 1–400: redaction, `Collect`, `Record`, `BuildRecord`, `ResolveFeedback` through its start); `commitAndPush`, `verifyFeedback`, and the four clone-pinning tests in `report_test.go` were not read in full — only their names were grepped.
- `internal/harness/main_test.go`/`cmd/sh-agent/main_test.go` bodies beyond earlier rounds' `TestHarnessMap` check — not read.
- `skills/preflight/SKILL.md` Step 6 full text — not read (G5 carried forward unverified independently, as in prior rounds).
- `.agents/docs/deploy/local.md`, `.agents/docs/specs/sh-agent/delivery.md`, ADR-029/ADR-033(draft doesn't exist yet) — not read.
- `internal/harness/adopt.go` read only its first ~40 lines; confirms it takes a disk `repo string` throughout what was read, consistent with "adopt requires a checkout," but the full file was not read.

design needs revision (G8 is a gap the plan does not decide; G1–G7 from Rounds 1–2 are resolved)

implementation not ready — 1 BLOCKING question (G8)
