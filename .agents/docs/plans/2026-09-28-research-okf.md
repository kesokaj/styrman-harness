---
owner: main
date: 2026-09-28
---
# Research — OKF (Open Knowledge Format) against the harness   (2026-09-28)

## Question

The operator asked: "Have we implemented OKF in styrman-harness, should we do it?" This record answers both questions and ends with what approving it would do.

## Answer in brief

1. **Not implemented.** There is no trace of OKF anywhere in the repo, its history, the machine state or the six projects on this box, and no directory is a conformant bundle.
2. **Recommendation: A, not now and not never.**
   - Nothing reads OKF: no coding-agent CLI does, and no project here produces or consumes it.
   - The spec has one committer and has changed without a version bump.
   - Every way of adopting it costs more than it returns.
   - Only one realistic trigger would reopen the question: a project that actually has to import, publish or maintain an OKF bundle.
3. **Approved by the operator on 2026-09-28.** They answered §6's question: "no BigQuery work on other machines". So A stands.

## Method

- **Eight agents** in one workflow run on 2026-09-28, all on Opus 5.5:
  - four independent read-only readers: spec, ecosystem, harness map, consumer projects;
  - one drafter;
  - three adversarial verification lenses: external facts, harness claims and compliance, and a devil's advocate arguing for adoption.
- **Verification:** all three lenses returned "holds with corrections", and none changed the recommendation. Their corrections are applied throughout and listed in § Corrections from verification.
- **Spec source:** a clone of `GoogleCloudPlatform/open-knowledge-format` at `ad30107` (2026-08-21). `SPEC.md:n` below means that clone.
- **Other sources:**
  - the v0.1 text from `knowledge-catalog` at `d44368c15e`;
  - the three Google Cloud blog posts (2026-06-12, 2026-07-25, 2026-08-27);
  - live vendor changelogs;
  - eight community skill repositories;
  - this repo at `914dc90`.
- **Notes, clones and saved changelogs** are in `.agents/scratch/research/okf/` (gitignored): `spec.md`, `ecosystem.md`, `harness.md`, `consumers.md`, `draft.md`, `factcheck-external.md`, `check-harness.md`, `counter.md`.
- **Re-checked by the main agent** against the clone, the saved changelogs and the repo:
  - the conformance and trust wording (`SPEC.md:3,57-63,403-410,736-764`);
  - zero OKF hits in the Claude Code (2.1.283), agy (1.2.12), Gemini CLI and Codex changelogs;
  - the contents of the fabricioctelles skill;
  - ADR-005, ADR-006, ADR-007;
  - `embed.go:9`;
  - `deploy/remote.md:56`;
  - the frontmatter keys our own skills use.

## 1. Have we implemented it?

**No.** OKF appears nowhere in:
- this repo;
- its full history (`git log --all -S OKF`, `--grep okf`);
- `~/.config/sh-agent`;
- the six projects under `~/workspace`.

The only hit for `git grep -i okf` is `@hookform` in `skills/frontend/references/component-patterns.md:204`.

**Nothing is conformant either.** OKF requires every non-reserved `.md` file to carry YAML frontmatter with a non-empty `type` (`SPEC.md:736-744`). Of the 743 tracked `.md` files, 661 have frontmatter and none has `type:`. In `.agents/docs` (628 files):
- 525 carry only `owner:` and `date:`;
- 87 bug records add `id`, `status` and `area`;
- 9 retros add `session`;
- 3 deploy guides add `target`;
- 4 have no frontmatter.

## 2. What OKF is: the facts that decide this

- **The format.** A bundle is a directory of markdown files.
  - **What conformance requires:** every non-reserved `.md` file has parseable YAML frontmatter with a non-empty `type`, from an open, unregistered vocabulary. `index.md` and `log.md` are optional; when present, they follow §8 and §9 (`SPEC.md:736-744`).
  - **What consumers MUST NOT reject:** unknown keys, unknown types, missing optional fields, broken links, or a missing `index.md` (`SPEC.md:755-762`).
- **Versions.**
  - **v0.1** shipped on 2026-06-12 inside Google's `knowledge-catalog` samples.
  - **v0.2** shipped on 2026-07-24/25, with two breaking renames under a minor bump: `timestamp` became `generated.at`, and the body `# Citations` list became `sources`. It also added optional fields: `sources`, `usage_window`, `generated`, `verified`, `status: draft|stable|deprecated` and `stale_after`. Finally, it added an `Attested Computation` type, whose attester code runs on the consumer's side (`SPEC.md:287-433,584-608,803-811`).
  - The spec moved to its own repo on 2026-08-11.
  - **Changed in place:** on 2026-08-21, PR #6 turned every timestamp-valued key (`stale_after`, `sources[].last_modified`, `usage_window`) into a datetime with a UTC offset, and did not bump the version. "0.2" now names two texts. Issue #24 about this is open.
- **What it is for.** Knowledge about data catalogs.
  - All four sample bundles are analytics datasets.
  - The reference agent is a *producer*, built on Python, ADK, BigQuery and Gemini.
  - Google ships no consumer agent, no MCP server, no validator and no JSON Schema. Issue #8, which asks for a schema, is open with no replies.
  - On 2026-08-27 Google's third blog post showed ingestion into Knowledge Catalog. It came after the last spec commit, so product work continues. Whether that ingestion is preview or GA is not stated.
- **Governance.**
  - One Google engineer made all 6 commits on the current main.
  - There are no tags and no releases.
  - No maintainer has replied on any OKF thread in either repo: about 361 human comments across 119 threads. Maintainers are active elsewhere in `knowledge-catalog`.
- **Security.**
  - There is no security section.
  - The trust fields are unsigned and self-asserted, and the trust tiers derived from them are "advisory signals, not access control" (`SPEC.md:410`).
  - Attester sandboxing is deferred (`SPEC.md:788`).
- **Coding agents.** None of these reads OKF natively as of 2026-09-28, each checked through its changelog: Claude Code (2.1.283), agy (1.2.12), Gemini CLI (all 641 releases, plus a code search), Codex (500 releases) and Cursor (2026-08-19 to 2026-09-23).
  - This is unlikely to change, because OKF deliberately prescribes no runtime (`SPEC.md:57-63`). An agent reads a bundle with Read and grep, which both CLIs already do.
- **Relation to AGENTS.md.** The spec never mentions AGENTS.md, Agent Skills or MCP.
  - Google's v0.1 blog calls "the AGENTS.md / CLAUDE.md family of convention files" one bespoke example of the pattern that OKF formalises.
  - The largest producer by stars among those found, LangChain OpenWiki (16.8k), still writes `AGENTS.md` and `CLAUDE.md` so that agents find its bundle (its README lines 310 and 361-374). Writing `CLAUDE.md` into a project conflicts with R2.
  - *Inference:* OKF works in practice as a knowledge layer beside AGENTS.md, not a replacement for it.

## 3. The approaches

- **A. Watch.** Commit this record as the answer and change nothing else. §5 names the triggers that would reopen the question.
- **B. Vendor a community skill.**
  - **Best candidate:** `fabricioctelles/skills` → `okf-open-knowledge-format` (Apache-2.0).
    - 694-line SKILL.md with a description of about 750 characters.
    - 104 KB in all, including copies of both spec versions.
    - Ships `scripts/validate.sh` (8.5 KB of bash that greps frontmatter rather than parsing it) and offers to install the third-party `okflint`.
    - Repo HEAD is `4fcc2bf` (2026-09-27); the skill directory last changed at `875f499` (2026-09-09).
  - **Rejected candidates:**
    - `scaccogatto`: a Stop hook, a `uv run` MCP server, and a backfill skill that reads `~/.claude` transcripts.
    - `inkeep`: needs its own MCP server, even though it is the most-installed OKF skill (about 6.9k installs).
    - `rakibtg`: targets v0.1 and says to use it proactively in all coding work.
    - `siculo` and `catancs`: break the Agent Skills name rules.
    - `hamakyo`: maintains an `okf/` bundle beside the trail after every durable change, which would be a second trail.
    - `xSAVIKx`: `go build`s 13 binaries into `$HOME`.
- **C. A harness-owned authoring skill.**
  - Prose only, about 120 lines, pinned to `SPEC.md@ad30107`, with no script.
  - Carries the one safety rule the spec lacks: never run attester code that arrives inside a bundle without the operator's approval.
  - Must take a side on the spec's open contradictions: the spec recommends `/` links while Google's own reference agent forbids them, the relative-path base is unclear, and so is frontmatter on `log.md`.
- **D. Make the harness's own trail conformant, and have `/init-project` scaffold it.**
  - New trail files get `type:`.
  - The bug `status: fixed|open|wontfix` key is renamed, because it collides with OKF's lifecycle `status`.
  - `bugs/INDEX.md` gets `type:`, or a new name: on macOS's case-insensitive disk it collides with the reserved `index.md`.
  - One context sentence changes.
  - This touches at least 21 files: 1 context file, 7 agents, 7 skills, `init-project/SKILL.md`, 3 deploy templates, and the `DECISIONS.md` and `INDEX.md` templates. It amends rebuild decision #11 (owner and date), since both keys stay.
- **E. `sh-agent` validates OKF.** Either:
  - a YAML dependency, which breaks R4;
  - a hand-written parser for the YAML subset OKF uses (several hundred lines, unmeasured) in "a file copier with a manifest" (`AGENTS.md:57`);
  - or a check that only `type` is present, about 30 lines reusing `internal/harness/frontmatter.go`, which does not prove conformance.
- **F. A project pilot.**
  - The closest shape is `unifi-home-network`, "a written record of what the network actually looks like" (its `AGENTS.md:5-7`). A device and WLAN inventory with `verified` and `stale_after` would express in fields what its `AGENTS.md:49` says in prose.
  - Its real gap is different: its config snapshot is stranded in a handoff, and plain trail files fix that without OKF.
- **G. Reject it as a non-goal.**
  - Add a REQUIREMENTS line and an ADR. This would be the first non-goal added that way: ADR-007 only annotated the hooks line, which already existed.

## 4. Comparison

| Approach | Pros | Cons | Complexity | Cost | Maturity |
|---|---|---|---|---|---|
| A. Watch | Costs nothing. Nothing to re-pin while the spec moves. A grep for "OKF" finds this record. | No hands-on experience. The harness ships no watcher (ADR-001: the installer is the only executable), so a trigger is noticed only when someone looks. | None | 0 context bytes; 3-4 trail files | n/a |
| B. Vendor skill | Quickest route. No Go. Drops in under R2. | **Breaks the zero-scripts rule** (ADR-006, ADR-007): `validate.sh` would ship, and `embed.go:9` would build it into every release binary. Offers a third-party linter. The broad trigger fits our own trail and would fire on routine trail work; it cannot be trimmed in place (a trimmed fork is C). One maintainer. Nobody here would use it. | Low | +29 B of context. A ~750-character description in every session on both CLIs. 104 KB to vet for payloads (arXiv:2606.19191, PhantomSkill). A re-vendor on every spec change. | Low |
| C. Harness skill | Our text. Short. Carries the missing safety rule. Precedent exists for generic skills no project uses yet (`gke-basics`, `terraform`, `json-api-style`). | No known user (YAGNI). OKF is not in the declared stack. Must take sides on contradictions no maintainer has resolved. Needs a re-pin whenever "0.2" changes in place. | Medium | +14 B of context. Description tokens in every session. Reviewer and a sandbox check on both CLIs. | Low |
| D. Conformant trail | Typed documents. Approvals could be recorded as `verified: {by: human:<id>}`. | **`type:` adds nothing**: every trail file is already typed by its directory, and every reader selects by path. Append-only records never gain `type`, so every trail stays mixed and non-conformant. Needs the `status` rename and has the `INDEX.md` collision. `stale_after` would not have caught this harness's actual stale-text bugs (BUG-027, BUG-042, BUG-084), which are restatements, not age. | High | At least 21 files, an ADR, and a MINOR release that changes every project's trail | Low |
| E. Validator | A mechanical check. | R4 (a dependency) or several hundred lines of parser. Nothing to check without D. No schema or corpus to test against. May only run by hand (R3). | High (Medium for type only) | A new command, tests, a security review, a pinned spec SHA | Low |
| F. Project pilot | The right home under global-vs-project. No harness bytes. | Nothing reads it there either. The actual gap closes without OKF. | Low | Work in that project only | Low |
| G. Reject | Closes the question. | Claims more than the evidence supports: "nothing reads it" and "immature" can both change. A REQUIREMENTS edit to add, and another to undo. | Low | 1 REQUIREMENTS line and an ADR | n/a |

Fit with this project's own rules ("ok" fits, "strains" fits awkwardly, "conflict" breaks it, "–" does not apply):

| Rule | A | B | C | D | E | F | G |
|---|---|---|---|---|---|---|---|
| R2: standards, skills drop in unchanged | ok | ok | ok | ok for the trail (see T5 for skills) | ok | – | ok |
| R3: no hooks, gates or CI | ok | ok | ok | ok | strains: by hand only | – | ok |
| R4: stdlib only | ok | ok for the binary | ok | ok | conflict or strains | – | ok |
| R5 / file drop, no code | ok | ok | ok | ok | new code: a design question | – | ok |
| Zero scripts (ADR-006, ADR-007) | ok | **conflict** | ok | ok | ok | – | ok |
| Vendored verbatim (`AGENTS.md:56`) | – | binding: no trimming in place | – | – | – | – | – |
| Append-only records | – | – | – | conflict: old records stay untyped | – | – | – |
| ADR-035 budget (20,329 of 22,000 B, re-measured) | 0 | +29 B | +14 B | tens of bytes | 0 | 0 | 0 |
| Global vs project (`context/AGENTS.md:37,41`) | ok | ok | strains: no known user | ok in form, no project asks | strains | ok | ok |

R1, R6 to R9 and the other non-goals are unaffected by every approach.

On precedent: the operator asked for agents.md, and the reason recorded in ADR-005 is just that ("The operator asked for the agents.md standard"). R2 adopts SKILL.md for ecosystem drop-in. No ADR states "adopt a standard when something outside the harness reads it", so that criterion is this report's inference, not a rule. That is why it is not a column above.

## 5. Recommendation

**A: not now, but not never.** Commit this record and change nothing else.

- **No reader and no producer here.** No CLI reads OKF, and nothing the operator runs emits a bundle that should drop in. A conformant trail would be read exactly as today's markdown is, because readers select by directory, not by `type`.
- **No need on this box.** None of the six projects uses BigQuery, Knowledge Catalog, Looker or dbt. The two needs shaped like OKF are both met without it:
  - `unifi-home-network`'s stranded snapshot;
  - `ranta-rankare`'s BUG-005, about where constants come from and when they were checked.
- **The spec cannot be pinned with confidence.** It has one committer, no maintainer triage, no tags, "0.2" naming two texts, and Google's own tooling contradicting the spec on links.
- **Every other option costs more than it returns.**
  - B breaks the zero-scripts rule.
  - C is honest work for no user.
  - D's `type:` adds nothing the directories do not already say, and leaves every existing trail non-conformant forever.
  - E needs D and a parser.
  - F closes a gap that plain files already close.
- **Not G:** both objections are observable and can change. Waiting costs nothing, so no REQUIREMENTS line is needed.

**Triggers that reopen it.** T2 is the realistic one. T1 is unlikely to ever fire, because OKF asks no CLI to do anything beyond what Read and grep already do. So A works like G unless a real need turns up in a session. That is acceptable because waiting costs nothing, and grepping for "OKF" finds this record.

| # | Trigger | How it is noticed | Smallest first step |
|---|---|---|---|
| T2 | A project must import, publish or maintain an OKF bundle: BigQuery or Knowledge Catalog work, a bundle handed over, an Outline export | Scoping that work with `/feature` | **F in that project:** a prose project skill under a non-harness name, pinned to a spec SHA. No harness change. If OpenWiki is the producer, it writes `CLAUDE.md`, which R2 forbids, so bring that to the operator as a harness question. |
| T3 | A second project needs the same | A second project skill like T2's exists | **C** through `/retro --harness`. The project copies are removed or renamed in the same change, because a project skill never reuses a harness name. |
| T1 | Claude Code or agy documents native OKF reading | Grep their changelogs for `open knowledge format\|\bokf\b` (0 hits on 2026-09-28) | A probe project holding a three-concept bundle, asked with `claude -p` / `agy -p` in the real home or from a sandbox home per `testing/guide.md` §4, as `testing/2026-09-21-verify-agents-md-native.md` did. Then reconsider D, limited to what the CLI actually reads. |
| T4 | The spec stabilises: a tag or release, an answer on #24, or a Google validator or JSON Schema (#8) | Reading the spec repo | Not a trigger on its own; a precondition for D and E. Re-read SPEC.md at the tag, and only then cost out E. |
| T5 | `type: Skill` lands in the canonical SPEC.md (proposed only in `knowledge-catalog` #312, an issue against the frozen copy), or agentskills.io adopts `type` | Reading the spec repo or agentskills.io | Check that both CLIs load a SKILL.md with a top-level `type:` key. The agentskills `skills-ref` validator rejects unknown keys, but 8 of our skills already carry `argument-hint`, so CLI loading is the real bar, not `skills-ref`. |

## 6. The one question that could change this

The consumer survey covered the six projects on this Linux box. The harness is also installed on the work Mac (Antigravity only) and on two other Linux hosts where the operator works on other projects (`deploy/remote.md:56`), and none of those was surveyed. If work there touches BigQuery, Knowledge Catalog or OKF bundles today, T2 has already fired. The answer then becomes F on that machine's project, and C once a second project needs it. If it does not, A stands as written.

**Answered 2026-09-28.** The operator, in their own turn: "Approve A, no BigQuery work on other machines". A stands. This is the operator's statement; those machines were not surveyed.

## 7. What you would be approving

1. **This record:** `.agents/docs/plans/2026-09-28-research-okf.md`.
2. **The change entry:** `.agents/docs/changes/2026-09-28-research-okf.md`, naming all eight agents dispatched (four readers, a drafter, three verification lenses) and the main agent's spot checks.
3. **The Definition of Done for a docs-only change over ten lines:**
   - `verifier` runs the project's checks on the unchanged source tree (step 1) and writes its log.
   - No `reviewer`, because there is no source diff (`agents/reviewer.md:3`).
   - No `documentation` dispatch, because no behaviour changed.
4. **Closing the trail:** a handoff that says nothing is queued, then commit and push.

**Cost:** four or five markdown files, docs only. Nothing changes in `context/`, the skills, the agents, the templates, the Go code, `REQUIREMENTS.md` or `DECISIONS.md`. It uses 0 bytes of the context budget and needs no release.

**Found along the way, not part of this approval** (details in `.agents/scratch/research/okf/check-harness.md` and `consumers.md`):

| # | Problem | Recommendation |
|---|---|---|
| 1 | `owner:` means different things: our deploy guides say `main`, the `/init-project` deploy templates say `documentation`, and `testing/guide.md` says `verifier` | defer: decide whether `owner:` names the writer or the lane |
| 2 | `plans/2026-09-21-context-audit-prompt.md` has no frontmatter (`context/AGENTS.md:65`). `handoffs/2026-09-04-argolis-devbox-instructions.md` also lacks it, but handoffs are immutable. | **fixed** in this record's commit, as the operator asked on 2026-09-28; the handoff is left alone |
| 3 | unifi-home-network: the config snapshot and MCP quirks live only in its 2026-09-14 handoff (Next Step 5, marked optional) | fix, in that project |
| 4 | ranta-rankare: BUG-005's prevention rule is half applied; `app.py:43-45` and `:251` name no source or as-of date | fix, in that project |
| 5 | valkompass: the README names a Cloud Run revision that was deleted on 2026-09-22 | defer: the project is torn down |

## Corrections from verification

The three lenses corrected the draft before this record was written.

- **Harness lens:**
  - B conflicts with the zero-scripts rule, which the draft had missed.
  - ADR-005's reason is "the operator asked". An "outside consumer" rule was never stated.
  - The two-project promotion is a retro trigger, not a gate against C.
  - T1's probe method must use the real home or a sandbox home, never `--home`, which is an `sh-agent` flag.
  - Handoffs are immutable, so found-along-the-way #2 shrank to one file, and #1 became a question.
  - The change entry must name all eight agents, and the verifier's step applies.
  - D touches at least 21 files, not 18.
  - D amends rebuild decision #11 rather than reversing it.
- **Facts lens:**
  - A third Google blog post exists (2026-08-27).
  - PR #6 changed the timestamp-valued keys only.
  - "Advisory" describes the trust tiers, not the fields.
  - The maintainer silence covers OKF threads only (about 361 human comments, not 405).
  - The OpenWiki citation was on the wrong lines.
  - `usage_window` was missing from the v0.2 list.
  - #312 is an issue against the frozen copy.
  - The Gemini CLI check was widened to all 641 releases.
- **Devil's-advocate lens:**
  - The survey was limited to this box (§6).
  - T1 is unlikely to fire by design, so T2 is the realistic trigger.
  - T2 as drafted could never fire for unifi-home-network, so F there is declined on its merits.
  - OpenWiki conflicts with R2.
  - D's `type:` duplicates the directory.
  - `stale_after` would not catch this harness's stale-text bugs.
  - B's trigger would fire on our own trail work.
- **Main agent:** the draft had called skills-ref's rejection of a top-level `type` a conflict. That is weak, because 8 of our skills already carry `argument-hint`, which skills-ref also rejects.

## Not verified

- **Absence findings:** no maintainer replies and no spec commits since 2026-08-21. These come from API reads on 2026-09-28 and could change any day.
- **Cursor:** checked only through its changelog for 2026-08-19 to 2026-09-23.
- **Knowledge Catalog:** whether its OKF ingestion is preview or GA is not stated anywhere read.
- **Third-party validators** (`okflint`, `okf-graph`, `W4G1/okf`, `jfomhover/okf-validator`): known only by name.
- **Size of a hand-written YAML-subset validator:** an estimate, not a measurement.
- **Case of the reserved filenames:** the spec does not say whether matching is case-sensitive. The `INDEX.md` collision comes from APFS being case-insensitive by default.
- **Other install places:** the projects on the work Mac and the other two Linux hosts were not surveyed. The operator answered for them instead (§6).
- **Measured benefit:** no eval compared an agent with and without a bundle. "No measurable benefit" rests on there being no reader and no use case.
- **Popularity counts:** star and install figures are as the GitHub and skills.sh APIs reported them, and some look inflated.
