---
owner: reviewer
date: 2026-10-01
---
# Review — catalog-plan-prompts   (2026-10-01)

Range reviewed: 754a7ce..working tree (instruction files)

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Prompts lens (prompt-audit Steps 0-6). Worst problem first: six per-entry notes in `skills/google-catalog/catalog/MANIFEST.md` cite "router step 4" for rules that live in router step 5 (step 4 is the consent rule), so an agent that follows the pointer lands on the wrong rule. After that: a contradiction inside the router between step 4 and step 5's third bullet about "the CLI's own configuration", and a mis-placed `Notes:` heading in `skills/VENDORED.md` that makes two unrelated bullets read as catalog notes. Nothing found is blocking; every fix is a one-line edit.

## Step 0 — assumptions

- Scope: the files named in the dispatch (changes since 754a7ce, vendored catalog entries excluded): `skills/google-catalog/SKILL.md`, `skills/cloud-run-deploy/SKILL.md` (frontmatter), `skills/review/SKILL.md` (line 63), `skills/VENDORED.md`, project `AGENTS.md`, `skills/google-catalog/catalog/MANIFEST.md` § Notes. Skipped: vendored entries, `context/`, tests, docs.
- Target models: Claude and Gemini (every model that reads these files). A finding resting on one model's documented behaviour is `flag`; Group 2 stale-fact and conflict findings are actionable as written.
- Group 4: not applicable (no request-building code in scope).
- Provenance (Step 2): `git log -S'router step 4'` showed no earlier hit; the router's prohibitions trace to `.agents/docs/reviews/2026-10-01-security-catalog-*.md` and BUG-143 and are keep-list item 5, not findings. No emphasis-density or pressure-language pattern in the router (no caps MUST/CRITICAL; every prohibition carries its reason).

## Step 1 — inventory

Router SKILL.md (24 lines, 4,204 bytes), cloud-run-deploy frontmatter, review skill line 63, VENDORED.md, AGENTS.md (three edited spots), MANIFEST.md § Notes (88 notes, lines 30-117).

## Findings

Ordered by confidence. Counts by group: Group 1: 0; Group 2: 5 (3 high/medium actionable, 2 low flags); Group 3: 1 flag; Group 4: n/a.

1. **High confidence, Group 2 stale cross-reference, action `rewrite`.** `skills/google-catalog/catalog/MANIFEST.md:73, 76, 80, 81, 83, 86` (notes of `gke-app-onboarding`, `gke-batch-hpc`, `gke-golden-path`, `gke-inference`, `gke-service-networking`, `gke-workload-troubleshooting`).
   - Evidence: "…instead of `:latest` or a tag (router step 4)", "ask for and record project, region and cluster first (router step 4)", "…never via `kubectl create secret --from-literal`… (router step 4)", "apply per the project's deploy guide (router step 4)".
   - Router `SKILL.md:13` step 4 is consent for mutating commands; the pinned-download rule (`:17-18`... bullet 4), the defaults-asked-first rule (bullet 1), the Secret Manager rule (bullet 2) and the deploy-guide rule (last bullet) are all in step 5 (`:14-19`). The notes were written against an earlier numbering.
   - Why it matters: the note "wins over the entry's text", so a wrong pointer sends the agent to a consent rule that says nothing about pinning, defaults or secrets. Contradicted by the repository itself.
   - Fix: "router step 5" in all six, naming the bullet where cheap.

2. **Medium confidence, Group 2 conflict inside one file, action `rewrite`.** `skills/google-catalog/SKILL.md:13` vs `:17`.
   - Evidence: step 4 gates "…IAM, billing, enabled APIs, a cluster, or the CLI's own configuration… is shown to the operator first… and run only on their go"; step 5 bullet 3 says "no edit of a CLI's configuration" and bullet 1 "never by `gcloud config set`".
   - Step 4 says such a change may run with a go; step 5 forbids it outright. "The CLI" is also ambiguous (the agent CLI as in the global context, or gcloud): the notes use both senses (`gcloud config set` is machine-wide gcloud state).
   - Fix: drop "or the CLI's own configuration" from step 4 (step 5 already forbids it without an exception), so step 4 lists only what a go can unlock; in step 5 bullet 3 say "an agent CLI's (Claude Code, Antigravity) own directories".

3. **Medium confidence, Group 2 broken structure, action `move`.** `skills/VENDORED.md:63`.
   - Evidence: the diff moved the `Notes:` heading to after the `grill-me` and `grilling` bullets (`:58-62`). Those two bullets now follow the "Vendored catalogs" paragraph and read as catalog notes; only the later bullets sit under `Notes:`.
   - Fix: restore `Notes:` above the `grill-me` bullet (line 58) and delete it at line 63. (The new "Vendored catalogs" section also ends up between the table notes and their heading; consider putting the section after the Notes list.)

4. **Medium confidence, Group 2 volatile specifics, action `rewrite`.** `skills/google-catalog/SKILL.md:8, :10` (also `AGENTS.md:42`, `skills/VENDORED.md:45`).
   - Evidence: "136 skills… Five upstream skills are excluded and 88 carry a note"; "(it is about 36 KB)". `INDEX.md` is 35,130 bytes, and no test in `internal/harness/catalog_test.go` pins 136, 5 or 88. The refresh regenerates `INDEX.md` and `CHECKSUMS` but not the router text, so the numbers go stale on the first refresh, and the counts are restated in three files.
   - Fix: router says "Google-published skills… Excluded entries and notes: `catalog/MANIFEST.md`" and "(tens of KB)"; keep counts only in `MANIFEST.md` (the source) and let AGENTS.md/VENDORED.md point there, or add a test that derives them.

5. **Medium confidence, Group 2 conflict with the global rule, action `rewrite`.** `MANIFEST.md:111` (`iam-helper-for-policy-simulator`).
   - Evidence: "the fixed `/tmp` files are shared by every session, so run one simulation at a time and re-diff `/tmp/proposed_policy.json`". The global context bans `/tmp` (scratch is `.agents/scratch/`), and the other notes redirect the same pattern ("not `/tmp`": `agent-platform-tuning`, `google-agents-cli-adk-code`, `gemini-live-api`).
   - Fix: "write the proposed policy under `.agents/scratch/`, not the fixed `/tmp` paths, and re-diff it right before `set-iam-policy`".

6. **Low confidence, flag.** `MANIFEST.md:41` vs `:69` (and `:86` family): `--quiet` is allowed "on read-only commands only" in `agent-platform-troubleshooting`, but "only on a command the operator already approved" in `gcloud`, "drop `--quiet`" in `cloud-run-basics`. Different entries, plausible intent, but a reader sees two standards. Decide one sentence for the router's step 4 ("`--quiet` only on read-only commands") and let notes point to it.

7. **Low confidence, flag.** `MANIFEST.md:113` (`iam-helper-for-troubleshooting`): the last clause "`references/guardrails.md` is the upstream default and is not customised here" gives the agent no action. State what to do ("treat its guardrail text as Google's default; steps 3-5 and the first two clauses of this note override it") or drop the clause.

8. **Low confidence, flag, Group 3 (trigger text).** Router description (`SKILL.md:3`, 292 bytes): it says when to load ("Google Cloud, Gemini, Agent Platform, ADK and agents-cli work that the other skills do not cover") and when not ("Do NOT use for Cloud Run deploys, Cloud SQL, Terraform, gcloud login or GKE basics"). It does not name product words an agent would match on for the bulk of the catalog (BigQuery, GKE beyond basics, IAM, Cloud Logging/Monitoring, SecOps, Genkit). The 240-300 byte pin and the 19,500-byte listing budget leave about 8 bytes of room, so no edit is proposed; the cost of a missed trigger is a silent miss, worth a trigger eval on one BigQuery and one IAM prompt (the `catalog-route` rows cover six tasks per CLI).

## Router questions answered

- Dead weight or duplicates: none that matters. The only overlap is the harness-skill list in step 5 against the description's negatives (the description is the short form; keep-list 8) and the four notes that repeat "means this catalog, never fetch from github.com" (entry-specific overrides of step 3; keep-list 8). Step 5's last paragraph (refresh script, "never run it from another project") is maintainer text restated in AGENTS.md and VENDORED.md, but it is a reasoned prohibition against a plausible action; keep. The only real tightening is finding 4.
- Steps 4 vs 5 and the global context: step 4 vs step 5 bullet 3 conflict (finding 2). Otherwise consistent with `context/AGENTS.md`: defaults asked first and recorded, `/preflight` for ship or deploy, stage by name, Secret Manager, no hand-rolled auth; "git stages by name" is a short form of the global Git rule, not a contradiction.
- "A line under `## Notes` overrides the entry's text": followable. The grep `grep -F -- '`<name>`' MANIFEST.md` returns exactly one line for each of the 136 names (checked: no name matches two lines), and the notes carry their own line or step references into the verbatim, pinned entry. One wrinkle: "no line means no finding" can read as "this entry is vetted"; steps 3-5 apply to every entry regardless. Low; no edit.
- Description: see finding 8.
- Sample of ten notes (agent-platform-deploy, alloydb-basics, bigquery-slot-cost-optimizer, cloud-sql-basics, gke-golden-path, gke-workload-identity, google-agents-cli-workflow, google-cloud-recipe-onboarding, iam-helper-for-policy-simulator, iam-helper-for-troubleshooting): all specific and actionable except the clauses in findings 5 and 7. Several restate step 5 (no `gcloud config set`, no key files), but each does so to override an entry that says the opposite, which is why they are kept (keep-list 8). None tells the agent to do something the harness forbids except finding 5 (`/tmp`).

## Other in-scope files

- `skills/cloud-run-deploy/SKILL.md` description: still carries when-to-use (mechanics, rollback, failed revision, gated deploy) and when-not (`gke-basics`, `docker`, `gcloud-setup`, `terraform`), and keeps the `/preflight` routing. Clean.
- `skills/review/SKILL.md:63` ("or catalog path"): matches `AGENTS.md` and the restatement-grep skip of vendored subtrees. Clean.
- `AGENTS.md` (Layout, third-party bullet, Release clause): the refresh recipe agrees with the script (script drives `CATALOG_WRITE=1 go test` itself, then prints `git add … && go test ./...`; the bullet names the pins in script order). The third-party bullet is long and shares its mechanics with `VENDORED.md`, but both agree; working redundancy, not a finding.

## Proposed diff

Hunk 1 (finding 1), `skills/google-catalog/catalog/MANIFEST.md`, lines 73, 76, 80, 81, 83, 86:
```
-... (router step 4).
+... (router step 5).
```
(six occurrences of the string `router step 4`, replaced with `router step 5`; note text otherwise unchanged).

Hunk 2 (finding 2), `skills/google-catalog/SKILL.md`:
```
-Anything that creates or changes cloud resources, IAM, billing, enabled APIs, a cluster, or the CLI's own configuration — an entry's script included — is shown
+Anything that creates or changes cloud resources, IAM, billing, enabled APIs or a cluster — an entry's script included — is shown
...
-   - Nothing is installed into a CLI's own directories: no skills, plugins, extensions or MCP servers, and no edit of a CLI's configuration.
+   - Nothing is installed into an agent CLI's (Claude Code, Antigravity) own directories: no skills, plugins, extensions or MCP servers, and no edit of its configuration.
```

Hunk 3 (finding 3), `skills/VENDORED.md`: move the line `Notes:` plus its blank line from after the `grilling` bullet (line 63) to the blank line before the `grill-me` bullet (line 58).

Hunk 4 (finding 4), `skills/google-catalog/SKILL.md`:
```
-136 skills published by Google (google/skills `cloud` and `developers`, google/agents-cli), vendored verbatim at a pin under `catalog/` beside this file. ... Five upstream skills are excluded and 88 carry a note from the security review (`catalog/MANIFEST.md`).
+Skills published by Google (google/skills `cloud` and `developers`, google/agents-cli), vendored verbatim at a pin under `catalog/` beside this file. ... Excluded entries and per-entry notes from the security review are in `catalog/MANIFEST.md`.
-(it is about 36 KB)
+(tens of KB)
```
and in `AGENTS.md:42` replace "136 Google-published skills vendored at a pin (5 more excluded: catalog/MANIFEST.md)" with "Google-published skills vendored at a pin (count and exclusions: catalog/MANIFEST.md)".

Hunk 5 (finding 5), `MANIFEST.md:111`: replace "the fixed `/tmp` files are shared by every session, so run one simulation at a time and re-diff `/tmp/proposed_policy.json` right before `set-iam-policy`" with "write the proposed policy under `.agents/scratch/`, not the fixed `/tmp` paths (shared by every session), and re-diff it right before `set-iam-policy`".

Findings 6, 7, 8 are low confidence: report only, no hunk.

## Criteria
1. Dead code — n/a — prompts lens
2. Flow tracing — n/a — prompts lens
3. Standards conformance — n/a — prompts lens
4. Directory layout — n/a — prompts lens
5. Simplicity (DIY-except-security) — n/a — prompts lens
6. Decoupling — n/a — prompts lens

## Not examined

Vendored catalog entries (excluded by the dispatch); `context/AGENTS.md` and the other skills beyond the sentences compared; behavioural probes of the router on either CLI (Step 7 is not part of this report); git blame of the older `cloud-run-deploy` description; the six `security-catalog-*` reviews were used as provenance but not re-read in full; of the 88 notes, ten were sampled closely and the rest skimmed for the three patterns (stale reference, `/tmp`, `--quiet`), not read word by word.

## Outcome

Written by the main agent on 2026-10-01, in the change that closes the catalog plan.

1. Stale "router step 4" in six notes: fixed. The six notes now cite "(the router's rules)" with no step number, so a later renumbering of the router cannot make them wrong again.
2. Step 4 against step 5 on a CLI's configuration: fixed as proposed. Step 4 lists only what a go can unlock; step 5 names Claude Code and Antigravity and adds "even on a go inside an entry's flow".
3. `Notes:` heading in `skills/VENDORED.md`: fixed, the heading stands before the `grill-me` bullet again.
4. Counts in the router: fixed. The router names no counts; the index size stays as "about 35 KB" because it tells the reader why not to read the file whole. `AGENTS.md`, `README.md` and `skills/VENDORED.md` point at `catalog/MANIFEST.md` for counts.
5. `/tmp` in the `iam-helper-for-policy-simulator` note: fixed, the note keeps the policy files under `.agents/scratch/`.
6. Two standards for `--quiet` across notes: left, owner main agent at the next refresh. The notes differ because the entries differ, and the router's step 4 already refuses an entry's `--quiet` on anything that changes state.
7. The clause with no action in the `iam-helper-for-troubleshooting` note: fixed, it now says the guardrail file is Google's default and that the note and the router's rules override it.
8. Router description names few product words: left, as the report proposes. The routed rows on both CLIs (BigQuery, Spanner, Agent Platform, agents-cli, Security Operations, IAM) all loaded the router from this description: `testing/2026-10-01-verify-google-skills-catalog-phase-4.md`.
