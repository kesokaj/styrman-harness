---
owner: reviewer
date: 2026-09-18
---
# Review — plan-sh-agent-works-everywhere-cold-read-1   (2026-09-18)

Range reviewed: plan file, round 1 — `.agents/docs/plans/2026-09-18-sh-agent-works-everywhere.md` and `internal/harness/apply.go`, `cmd/sh-agent/main.go` (excerpts), `skills/review/SKILL.md`

## (1) Gaps

- **BLOCKING** Phase 2's mechanism cannot do what the plan says. Phase 2 task 2 reads: `/review --project --all` scoped to `cmd/sh-agent` and `internal/harness`. But `skills/review/SKILL.md:25` states "Exactly one scope class per run: paths (any number), or one range/ref, or one PR, or `--project`; a mix is the error 'one scope per run'" and `skills/review/SKILL.md:35` defines `--project` as "the whole repository at HEAD ... no diff." There is no syntax that scopes `--project` to two directories. Run as written it either errors ("one scope per run" if paths are also given) or reviews the *entire* repo — including `context/`, `skills/`, `agents/`, which the plan's own Scope line explicitly puts Out ("the harness content itself ... other sessions"). The implementer needs a decision: drop `--project` and pass the two paths instead (then `--all` lenses run on paths, which the skill supports), or accept a whole-repo review and filter findings by hand. This also breaks the `--tests` lens's "map every source file in the repository to its covering test" under `--project`, which the plan's task line paraphrases as scoped to the two packages.
- **Gap** Phase 1's own target code (`internal/harness/apply.go:35-63`) has two paths the plan's four verification cases don't cover: (a) `os.ReadFile` failing with something other than `os.ErrNotExist` (e.g. EACCES) falls into the same `default:` branch as "file absent" — under the new vocabulary this would be printed as `add ~/path` for a file that already exists but couldn't be read, which is misleading and untested; (b) the mode-correction branch (`apply.go:49-55`) writes to disk (`os.Chmod`) on an "unchanged" (byte-equal) file but the plan's Scope line claims the install output "names each file it adds, updates or removes" — is a mode-only fix an `update` line, or does it stay silent as today? Two implementers would diverge here and the Verification section's four test cases (first install/changed file/removed file/dry-run) don't settle it.
- **Assumption not checked, but checkable in one command**: the E2E golden path (`sh-agent install` from a release binary → `claude -p` quotes a rule) and the Phase 1 checkpoint's literal "260 add lines" were spot-checked (`go run ./cmd/sh-agent install --home <scratch> --dry-run` today prints `260 written`) — this matches, not a gap, but the plan should say the number is illustrative/re-derived at checkpoint time since it will drift the moment `context/`, `skills/`, or `agents/` change (Out-of-scope work in "other sessions" per the plan's own Scope line).
- **D3** (release version PATCH vs MINOR) is left to the implementer at the gate "from what shipped" — reasonable, but the plan never states whether the Phase 2 audit's fixes (which are unscoped until the audit runs) could themselves force MINOR/PATCH classification questions beyond D3's stated criterion (new output lines only). Low-severity ambiguity, not blocking since D3 explicitly allows the operator to override at `/preflight`.

## (2) Scope vs Why

- The `## Why` quote is about verbs working reliably across place kinds ("All features of sh-agent need to work ... make sure the harness is installed when you say install"). The Scope section additionally commits to "A full-repository code and security audit of `cmd/sh-agent` and `internal/harness`, every finding fixed or rebutted" — this is plausible under "this is the core of the harness" but has no explicit D-row quoting the operator authorizing an *audit* specifically (D1 only cites the verb/place-kind scope, not the audit mechanism). Not a hard violation — "flawless" plausibly implies it — but worth flagging since it is real, unbounded extra work (Phase 2) with no upper bound on rounds beyond the generic "three rounds, then the operator" rule.
- No other extra endpoints/shapes/config found: "No new flags, no new verbs, no config keys" is stated explicitly (Interfaces, line 28) and matches what Phase 1 actually changes (existing `write`/`remove` dry-run vocabulary renamed, not a new flag).

## (3) Readiness

Questions a first-pass implementer needs answered:

1. **BLOCKING** — How should Phase 2's audit actually be scoped, given `--project` cannot take paths per `skills/review/SKILL.md:25`? (see Gap above)
2. **BLOCKING** — Is a permission-denied (non-ENOENT) read of an existing destination file reported as `add`, `update`, or an error that stops the install? Today's code silently treats it as "no previous file." (apply.go:33-39)
3. **BLOCKING** — Does the mode-correction-only branch (apply.go:49-55) get an `update` line under the new "names each file it changed" contract, or stay silent as today?
4. **NON-BLOCKING** — Order of add/update/remove lines within one run: plan doesn't state it, but "same order the plan iterates files, resync loop after" is the obvious default and any reasonable implementer would keep the existing iteration order.
5. **NON-BLOCKING** — Whether the audit's extra scope needs an explicit D-row; a reasonable implementer proceeds under D1's broad "all features... flawless" and the operator can strike it at `/preflight` review time.

## Verdict

`design needs revision` — the Phase 2 audit mechanism as literally written (`/review --project --all` scoped to two directories) contradicts the review skill's own scope rules and needs a decision before Phase 2 starts. `implementation not ready` — three BLOCKING questions above (the audit scoping mechanism, and two untested code paths in the exact function Phase 1 rewrites) would produce different builds from two careful implementers.

## Not examined

Did not read: `internal/harness/{plan,repo,source,release,report,adopt,state,render,inventory,derive,pathguard}.go` in full (only apply.go and a slice of main.go were read; release.go, pathguard.go, report.go, inventory.go — all named in Surfaces touched — were not opened), `cmd/sh-agent/main.go` beyond lines ~220-360 (the full `update`, `status`, `pathHint`, `NoteRunningElsewhere` bodies were not read), `README.md` § Install/Daily commands, `.agents/docs/deploy/{local,remote}.md`, `.agents/docs/specs/sh-agent/delivery.md`, `.agents/docs/testing/guide.md`, `internal/harness/release_test.go` (the fake-`gh` fixture the Risks section leans on), and `.agents/docs/plans/2026-09-18-recon-sh-agent-audit.md` (does not yet exist — it's written by Phase 2's first task). Did not verify Phase 4/5 real-machine claims (out of scope for a static cold read). Given the 30-tool-call budget, this is a first-pass reading of the plan plus the one surface (`apply.go`) that Phase 1 directly rewrites, plus enough of `main.go`/`skills/review/SKILL.md` to check the two load-bearing claims (BUG-079 note already shipped; `/review --project` scoping).
