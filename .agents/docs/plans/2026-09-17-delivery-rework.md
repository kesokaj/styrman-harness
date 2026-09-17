---
owner: main
date: 2026-09-17
---
# Plan — delivery through `sh-agent` only, with a feedback loop back   (2026-09-17)

## Why

The operator, in their own words on 2026-09-17: "I want the delivery to be around the sh-agent only. Nothing else. I'm the sole user of this harness and decide everything but I install it on multiple places and when I decide to update the repo I want a smart way to push any updates to places where I use it." Four places: the work Mac (Antigravity only), the home Mac (both CLIs), argolis-devbox (both), styrman-crew (both). The normal flow is the operator updating this repository, then running `sh-agent update` at each place.

Two things stand in the way today. The two downstream projects do not use `sh-agent update` at all — they vendor a `.git`-less copy of this repository that is refreshed by hand, which is why a mirror refresh keeps appearing as a task in three repositories at once. And a signed-tag channel was designed yesterday (ADR-026) for styrman-crew's host, which this supersedes: it existed to give that host a trust root, and the operator's own credentials are already on every machine, so it solves a problem that is not there.

The operator also asked for a loop back: "Some sort of feedback function from each place installed reporting back to the main is good for a feedback loop so we can take experiences from each installed place to use to make the harness better."

## Scope

In, this repository only: a decision record superseding ADR-026's delivery half; `REQUIREMENTS.md` R8; a new `sh-agent report` command with its tests; the deploy guides rewritten around the four places and around the release bundle; the mirror sections retired.

Out: the changes each consumer needs in its own repository, which are their own sessions under their own Definition of Done — styrman-crew retires `image/harness/` and keeps a checkout, argolis-devbox retires `deploy/harness/`. This plan names what they must end up with and nothing more, because the operator's rule is that this repository is not connected to any project.

Out: the webservice the operator floated. Assessed and struck — git already carries both halves, a service adds hosting, per-machine auth and an attack surface for one user, and it contradicts R3 (no pipelines), R4 (one stdlib binary) and the local-by-default rule. A dashboard over collected data stays possible later without changing how any machine reports.

Follow-ups noted, not built:
- **Versioned release bundles** (the operator: "Maybe we should release the agent as a package on GitHub that is versioned?"). Worth doing and formalises the zip path they already use for major remodelling, but it is a second feature with its own plan — found in: scope, 2026-09-17. Assessment in `## Decisions` D6.
- Chaining `report` onto a timer, once it has proved itself by hand — found in: scope.
- `sh-agent update --require-signed`, if a host ever needs a trust root that is not the operator's credential — found in: ADR-026 superseded.

## Surfaces touched

`cmd/sh-agent/main.go` (the `report` verb, its flags, its help), `cmd/sh-agent/main_test.go`; `internal/harness/report.go` (new) and `harness_test.go`; `internal/harness/state.go` (the recorded feedback-repo path); `REQUIREMENTS.md` (R8); `.agents/docs/architecture/DECISIONS.md` (ADR-030); `.agents/docs/deploy/local.md` and `remote.md`; `README.md`; `AGENTS.md`.

`skills/retro/SKILL.md` (its harness-scope bullet ends "name the downstream mirrors its `deploy/remote.md` lists as the next task in their projects" — with no mirrors, that instruction has no referent) and `AGENTS.md`'s Key-files row for `deploy/remote.md`.

The first draft of this plan claimed `skills/` was untouched and that this was "not a behaviour change for any agent". The cold read proved that false at `skills/retro/SKILL.md:132`. Not touched, and checked: `context/` and `agents/`.

## Interfaces

```
sh-agent report [--dry-run] [--feedback-repo DIR] [--home DIR]
```

- **Runs before `FindRepo`.** `main.go:105` calls `FindRepo` for every verb except `version` and `uninstall`, and it hard-errors when there is no checkout. `report` joins that early switch, because a machine that installed from an unpacked bundle has no checkout and must still be able to report. Every repo-derived field is then optional, not fatal.
- **Resolves the feedback checkout** two ways and no more: `--feedback-repo`, else `~/.config/sh-agent/feedback/`. When that default path does not exist, `report` clones `git@github.com:kesokaj/styrman-harness-feedback.git` into it — SSH, the same transport and the same credential the operator already uses for this repository, named as a constant in `report.go` and changed by editing that constant. A path given with `--feedback-repo` is never cloned into: if it is not a git repository, that is an error, because the flag exists for tests and for a checkout the operator placed deliberately. NOT `state.json` — `apply.go:106` rebuilds `State` wholesale on every install, so a path recorded there is erased by the next `sh-agent install` and the following `report` would silently clone a second checkout. NOT an environment variable, which would be a third mechanism for a value that never varies. `--feedback-repo` is added to `reorder`'s `takesValue` map (`main.go:447`) in the same task as the flag, or `report --feedback-repo X --dry-run` swallows `--dry-run` and pushes when a dry run was asked — BUG-002's exact shape, measured by the cold read.
- **Collects**, read-only, from `~/.config/sh-agent/`: `feedback.jsonl`, `retros.jsonl`, `proposals/*.md`, `state.json`. Each absent file is skipped and named in the output.
- **Redacts content before sending, structurally, and never refuses over it.** `secretShaped` (`plan.go:206`) matches a FILENAME and none of the collected names can match it, so reusing it here would be a check that cannot fire. Refusing on a match is worse than useless: this repository's own `skills/terraform/references/security-compliance.md:290-291` documents `password="secret123"` and `AKIAIOSFODNN7EXAMPLE`, so a proposal quoting that rule would jam the loop. And a needle-span replacement is worse still — `AKIA` is a PREFIX, so replacing the matched token strips the marker and ships the key, and `-----BEGIN` matches a PEM's header line while its body carries no needle at all, which would ship the key and report it redacted. The needles locate; they never define the span. The span is structural:

  **Needles** (the shapes the harness's Git rule already stops a commit for): `-----BEGIN`, `PRIVATE KEY`, `AKIA`, `AIza`, `ya29.`, `token=`, `password=`.

  **`feedback.jsonl` and `retros.jsonl`, line-delimited JSON.** Each line is unmarshalled into `any`. A line that does not parse is not dropped and not crashed on — it falls through to the text rule below. Every STRING value is walked recursively, through objects and arrays. A string containing any needle is replaced WHOLE by `[redacted by sh-agent report]` — the whole value, never a span inside it, which is exactly what makes a PEM pasted into a `user_quote` leave entirely rather than headless. The line is re-marshalled and `json.Valid` asserted on the result before it is written; a line that would not be valid JSON is a bug and fails the run loudly rather than shipping.

  **`state.json`, one JSON document.** The same walk over the whole document.

  **`proposals/*.md` and any other text.** A line containing `-----BEGIN` opens a block: that line through the matching `-----END` line inclusive becomes one `[redacted by sh-agent report: N lines]`. With no `-----END` before the end of file, redaction runs to the end and the record says `unterminated`. Any other line containing a needle is replaced whole.

  The original files are never modified; redaction happens in the copy being sent. Every redaction is recorded in `report.json` as file, line, the JSON key path when there is one, and which needle fired — never the text. `--dry-run` prints the same list. The record's location and needle name leak nothing: the needle is a source constant and the location is already marked in the shipped copy.

- **Generates** one metrics record, `report.json`. Fields from a checkout are `null` when there is none:

```json
{
  "ts": "2026-09-17T09:00:00+02:00",
  "host": "dave", "os": "linux", "arch": "amd64",
  "vendor": "all", "files": 260,
  "harness_commit": "c7af0de", "installed_at": "2026-09-16T21:59:15+02:00",
  "drift": 0, "skills": 27, "agents": 7,
  "go": "go1.27.1",
  "redactions": [{"file": "feedback.jsonl", "line": 12, "pattern": "token="}]
}
```

`go` is `runtime.Version()` of the RUNNING `sh-agent` binary, which is what built it — not `go version` from PATH, which can be a different toolchain entirely and is why `Reinstall` pins `GOTOOLCHAIN=local` (`repo.go:170`). `redactions` is `[]` when nothing fired.


- **Writes** to `<feedback>/reports/<host>/<date>-<HHMMSS>/`, one directory per run, never overwriting. Stages **by name** — the paths it just wrote, never `git add -A` — reads `git diff --cached --stat`, then commits and pushes as separate calls, exactly as the harness's own Git rule requires of everyone else. A pre-existing dirty file in the feedback checkout is left alone and named in the output; the first draft said it would be "committed along with the run", which would have routed around both the staging rule and the content scan.
- `--dry-run` prints the destination and every file it would write, runs the content scan, and writes nothing.
- `report` never writes to the harness repository and runs no git command there beyond `Inspect`.

No change to `install`, `update`, `status`, `uninstall`, `list`, `render`, `adopt` or `version`. `--message` was in the first draft and is dropped: it is a field the request never asked for, and it was the flag that triggered the `takesValue` defect.

## Failure modes

| What fails | What the operator sees |
|---|---|
| No feedback checkout and the clone fails | `report: could not reach <url>: <git's message>; clone it by hand or pass --feedback-repo`, exit 1. Nothing written. |
| `~/.config/sh-agent/` has nothing to report | `report: nothing to send from <path>` and the metrics record still goes, because "this machine is alive and at commit X" is itself the signal. Exit 0. |
| A collected file's CONTENT matches a credential shape | The token is replaced by `[redacted by sh-agent report]` in the sent copy, the original is untouched, and `report.json` lists file, line and pattern. The run continues, exit 0. Printed as `report: redacted 2 token(s) — see report.json`. |
| The push is rejected (someone else pushed first) | `report: push rejected; run `git -C <feedback> pull --rebase` and re-run`, exit 1. The commit stays local, so nothing is lost. |
| The feedback checkout is dirty | Untouched. `report` stages only the paths it wrote, so a stray file stays stray and is named in the output. |
| Run with no harness checkout at all | Works, because `report` runs before `FindRepo`. `harness_commit`, `drift`, `skills` and `agents` are `null`, and the output says which. |

## Verification

Tiers: unit for the collection, redaction and record shape; integration for the git side against a local bare repository; e2e for one real run from this box to the real feedback repository.

E2E criteria:
- Golden path: on this box, `sh-agent report` creates `reports/dave/<date>-<time>/` in the feedback repository with `report.json` plus the three collected files, commits, pushes, and a fresh clone elsewhere shows them.
- Edge case: a second run the same day creates a second directory and does not overwrite the first.
- Edge case: `--dry-run` against a real feedback checkout leaves `git status` clean.

Delivery is verified separately and needs no new code: `sh-agent update` on a checkout, and `sh-agent install` from an unzipped plain copy, both already measured on 2026-09-17 (plain copy installs with the running binary, `status` names it a plain copy, `update` refuses it).

## Phases

- [x] Cold-reader check — three rounds, `reviews/2026-09-17-coldread-delivery-rework.md`. Never reached both tokens; ended at the cap as the workflow prescribes. Twelve blocking findings over three rounds, each verified against the code before being accepted, and four changed the design materially: `report` runs before `FindRepo`, the feedback path is not kept in `state.json`, the redaction is structural rather than needle-span, and `skills/` is a touched surface. The one ambiguity that changes what gets built was put to the operator, who chose JSON-aware redaction (D9). The remainder is resolved in D9 and D10, 2026-09-17

### Phase 1 — The decision and the requirement (no code)

- [x] ADR-030 written, D1 to D6 in the file's own style; ADR-026's heading marked "Superseded in part by ADR-030". Verified append-only: one deleted line in `DECISIONS.md` and it is that heading, 2026-09-17
- [x] `REQUIREMENTS.md` R8 applied on the operator's approval of its exact text, given 2026-09-17 in reply to that text quoted in full. R7 unchanged; R8's closing sentence carries the reconciliation. `ROADMAP.md` Now points at this plan, 2026-09-17
- [x] Checkpoint — ADR-030, R8 and the ROADMAP pointer committed together as the trajectory rule requires; `git show --stat 9222762`: 3 files, 16 insertions, 2 deletions; pushed. The staged diff's one credential-shape hit was a context line of ADR-026's own prose about armor headers, zero added lines, 2026-09-17

### Phase 2 — `sh-agent report`

- [ ] `--feedback-repo` added to `reorder`'s `takesValue` map, and `TestReorder` gains `report --feedback-repo /x --dry-run` — verify: unit, `go test ./cmd/... -run TestReorder`, who: AI — expect: the case fails before the map entry exists and passes after; `--dry-run` survives as a flag rather than becoming the value of `--feedback-repo`
- [ ] `report` joins the pre-`FindRepo` switch beside `version` and `uninstall` — verify: unit, run `report --dry-run` with `--repo` pointing at a directory that is not a checkout, who: AI — expect: it reports rather than erroring; the repo-derived fields are `null`
- [ ] `internal/harness/report.go`: collect, redact, render the record — verify: unit, `go test ./internal/... -run TestReport`, who: AI — expect: absent files skipped and named; a JSONL fixture whose `user_quote` holds a full PEM has that whole value replaced and the line still satisfies `json.Valid`; a JSONL fixture with `AKIAIOSFODNN7EXAMPLE` mid-sentence has the whole value replaced, not the prefix; a malformed JSONL line falls through to the text rule instead of crashing; a markdown fixture with a BEGIN/END block collapses to one marker line; a fixture with none of them is sent byte-identical; every case continues to a written record and exit 0
- [ ] The git side: resolve, write the run directory, stage by name, read `--cached --stat`, commit and push as separate calls — verify: integration against a local bare repository under `.agents/scratch/`, who: AI — expect: two runs the same day make two directories; a pre-existing dirty file in the checkout is not staged; a rejected push leaves the commit local and says so
- [ ] The first push is into an EMPTY repository: a freshly created one has no commits, so the clone warns and leaves no upstream branch, and the first push must set it (`git push -u origin HEAD:main`). Every later run pushes plainly — verify: integration against a local bare repository with no commits, who: AI — expect: the first run creates `main` upstream and succeeds; the second run pushes without `-u`
- [ ] Clone-on-first-use, its own task because it ships to four machines and runs exactly once on each — verify: integration, point the constant at a local bare repository and run with the default path absent, who: AI — expect: the directory is created and populated; a second run reuses it and does not re-clone; a clone that fails exits non-zero, names the URL and git's message, and leaves no half-made directory
- [ ] `--dry-run`, and the invariant that `report` runs no git command in the harness repo beyond `Inspect` — verify: unit plus a grep of `report.go` for git invocations, who: AI — expect: `--dry-run` against a real feedback checkout leaves `git status` clean and still prints the redaction list; no `exec.Command("git"...)` in `report.go` targets the harness repo
- [ ] `report.json` shape: every field in Interfaces, `null` for the repo-derived ones with no checkout, `redactions: []` when nothing fired — verify: unit, who: AI — expect: the record parses and carries each key; the nothing-to-send case still writes the record and exits 0
- [ ] Checkpoint — expect: `gofmt -l cmd internal`, `go vet ./...`, `go test -race -count=1 ./...` all clean — evidence: verifier log path, 2026-09-17

### Phase 3 — The docs, rewritten around the four places

- [ ] `deploy/remote.md`: § Downstream mirrors and § Releases-as-delivery retire; a § The four places replaces them — verify: `git grep -nE 'image/harness|deploy/harness|Downstream mirrors|downstream mirrors' -- README.md AGENTS.md context skills agents .agents/docs/deploy ':!skills/mermaid-syntax'` is empty, who: AI — expect: empty. The strings are the ones that go stale, not the bare word `mirror`, which is ordinary prose at `context/AGENTS.md:57` and `skills/feature/SKILL.md:23`, nor `styrman-crew`, which is bug provenance at `agents/changelog.md:61` and `skills/feature/SKILL.md:88` — a grep that forced those to zero would delete true text and edit a vendored skill (cold read round 2) — docs: `.agents/docs/deploy/remote.md`
- [ ] `skills/retro/SKILL.md:132`: the harness-scope bullet stops telling an agent to "name the downstream mirrors its `deploy/remote.md` lists" and names the consumer step that replaces it — verify: the grep above, plus one headless start of each CLI, who: AI — expect: no skill refers to a section that no longer exists — docs: `skills/retro/SKILL.md`
- [ ] `AGENTS.md` Key-files row for `deploy/remote.md` drops "mirrors" and the ADR-026 reference; the Release section points at ADR-030 for delivery — verify: read, who: AI — expect: no living doc calls a tag a delivery channel — docs: `AGENTS.md`
- [ ] `deploy/remote.md:45` and `:93` and `README.md:53` ("every downstream mirror", singular) — the grep stays green while these go stale, because they say it in words the grep does not carry — verify: read all three, who: AI — expect: none of them describes a mirror or a per-project copy — docs: `.agents/docs/deploy/remote.md`
- [ ] `deploy/remote.md:82`'s versioning rule and its gate recipe SURVIVE — D2 keeps tags as release markers, so MINOR/PATCH/MAJOR and the signing procedure stay; only the prose calling a tag the way a consumer receives the harness goes — verify: read, who: AI — expect: the versioning rule present, no sentence describing a tag as a delivery channel — docs: `.agents/docs/deploy/remote.md`
- [ ] `deploy/local.md` and `README.md`: `update` for a checkout, `install` from an unpacked bundle, `report` for the loop back — verify: read them cold, who: AI — expect: a reader with only these docs can update any of the four places and send a report — docs: `README.md`
- [ ] Checkpoint — expect: `/review --docs` over the range finds no stale restatement — evidence: review report path, 2026-09-17

### Phase 4 — End to end and Done

- [ ] Create `styrman-harness-feedback`, private — verify: `gh repo view`, who: both — expect: it exists and is private
- [ ] The three E2E criteria above, from this box — verify: e2e, who: AI — expect: as written in Verification
- [ ] Definition of Done: verifier, reviewer over the diff, security (the command reads a credential-bearing directory and pushes to a remote), documentation, changelog, bugs, commit, push, `/handoff` — verify: each artifact on disk, who: AI — expect: every finding fixed, rebutted or owned

## Risks

- **Assumption:** every one of the four places has the operator's git and `gh` credentials, so no key is distributed. Stated by the operator on 2026-09-17. Checked by: the first `sh-agent update` at each place, which is the operator's own step. The accepted downside, in their words: "if I don't login in with ssh or gh I can't get updates." The unzipped bundle is the documented fallback.
- **Assumption, and it is about ANOTHER repository, so it is quoted rather than cited:** styrman-crew's ADR-062 D2 reads "`image/entrypoint.sh` seeds every non-MINIMAL project home with `sh-agent install --repo /styrman/harness --home $HOME` as the project user". An install writes files into a home; it does not create a checkout. So an agent working in such a home has the harness text and no repository to push to. Checked by reading that text on 2026-09-17. The plan leans on this for its read-only claim, and it cannot be re-verified from this repository, so the consumer session confirms it in its own repo before that claim is relied on there.
- `report` sends the operator's own words from `feedback.jsonl` to a private repository of their own. It carries no project source. It does not carry "no secrets" by construction — the redaction pass lowers the risk and reports what it caught, and anything it does not recognise goes. That is stated rather than claimed away, and it is why Phase 4 keeps the security row.
- The collected files name consumer projects: `feedback.jsonl` and `retros.jsonl` both carry a `project` field with an absolute path. That is not a breach of D5. D5 says this REPOSITORY does not depend on any project; the feedback repository is the operator's own cross-project index and naming projects is its whole job. Recorded here because the two rules look contradictory until read carefully.

## Decisions

D1 — 2026-09-17 — Delivery is `sh-agent update` from `main` on a checkout the operator owns, at four places; no other channel — by: the operator, "I want the delivery to be around the sh-agent only. Nothing else."
D2 — 2026-09-17 — The signed-tag channel is retired as a delivery mechanism; `v1.0.0` and any later tag stay as release markers — supersedes ADR-026's delivery half — by: the operator's turn of 2026-09-17, and the fact that their credentials are already on every machine
D3 — 2026-09-17 — No webservice — by: assessment, accepted by the operator's "Yes to all"; see `## Scope` Out
D4 — 2026-09-17 — The feedback repository is `styrman-harness-feedback`, private, one directory per machine per run — by: the operator asked for `styrman-harness-<something>`; "feedback" is the word the harness already uses for this loop (`feedback.jsonl`, the `feedback` skill, the learning loop), so it costs no new vocabulary
D5 — 2026-09-17 — This repository names no consumer project in its living docs; anything that wants the harness uses `sh-agent` — by: the operator, "styrman-harness and the feedback/trajectory should not be directly connected to any other project, if any other project wants to use the harness they need the sh-agent, that's it."
D6 — 2026-09-17 — Versioned release bundles are a follow-up, not this plan. Assessment: worth doing; it formalises the zip the operator already uses, gives a machine without a checkout a way in, and can be built locally by `/preflight` with `GOOS`/`GOARCH` cross-compiles and `gh release create`, so it needs no CI and does not touch R3 — measured 2026-09-17, both `darwin/arm64` and `linux/amd64` build here. It does NOT remove the credential requirement while the repository is private, because a release asset on a private repository still needs auth. Making it public is a separate question about content, not packaging — by: implementer, recorded for the operator
D7 — 2026-09-17 — Round 1 of the cold read returned `design needs revision` with seven blocking questions; six were verified against the code by the author before any was accepted. Four changed the design: `report` runs before `FindRepo` so a machine with no checkout can report; the feedback path is not kept in `state.json`, which `apply.go:106` rebuilds wholesale; the redaction scans content because `secretShaped` matches filenames none of the collected files can have; and `skills/retro/SKILL.md:132` makes `skills/` a touched surface. Two changed the writing: staging by name, and the mirror grep. `--message` was dropped as scope — by: cold read round 1, `reviews/2026-09-17-coldread-delivery-rework.md`
D8 — 2026-09-17 — Round 2 of the cold read returned `design needs revision` with five blocking questions; the two that changed the design were verified against the repository first. The content scan REFUSED on a match, and this repository's own `skills/terraform/references/security-compliance.md:290-291` documents `password="secret123"` and `AKIAIOSFODNN7EXAMPLE` — a proposal quoting that rule would have bricked the feedback loop permanently with no override. It now redacts the matched token, reports file, line and pattern, and continues. The Phase-3 grep was unsatisfiable on the bare word `mirror` and on `styrman-crew`, both of which are true prose in files the plan must not touch, one of them vendored; it now greps the strings that actually go stale — by: cold read round 2, `reviews/2026-09-17-coldread-delivery-rework.md` § Round 2
D9 — 2026-09-17 — Round 3, the cap, returned `design needs revision` with three ambiguities and three items the reader would resolve the same way. The operator was asked the one that changes what gets built and chose JSON-aware redaction over refusing with an override or not scanning at all. The needles are now stated to LOCATE and never to define the span: a whole JSON string value, or a whole text line, or a whole BEGIN/END block. That closes the two leaks round 3 found — a prefix needle shipping the rest of the key, and a PEM body with no needle in it — because no span is ever partial. The same-way items were applied without asking: the Phase-2 expect rewritten to redact-and-continue, `json.Valid` asserted per line, the three prose paths the grep cannot see, and the empty-repository first push — by: the operator's choice at the cap, 2026-09-17, and cold read round 3
D10 — 2026-09-17 — `deploy/remote.md`'s versioning rule and signing recipe survive the Phase-3 retirement; only the prose describing a tag as the way a consumer receives the harness goes. Round 3 called this ambiguous and it was — D2 keeps tags as markers, which settles it — by: implementer, resolving cold read round 3 item 3
Open: should the harness repository become public, which is what would let a machine bootstrap with no credential at all? — default: it stays private, and the credential requirement stands with the unzipped bundle as the fallback
Rejected: a webservice for delivery and feedback — hosting, per-machine auth and an attack surface for one user, against R3, R4 and local-by-default
Rejected: per-machine deploy keys — the operator authenticates as themselves everywhere, so they add key management and restrict nothing
Rejected: a branch in this repository for feedback — it would need write access to the harness, which is the thing being prevented
