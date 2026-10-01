---
owner: main
date: 2026-10-01
---
# Research — Google skills import: measurements and decisions (2026-10-01)

Request (operator, 2026-10-01): add every skill of github.com/google/agents-cli and of google/skills `skills/cloud` and `skills/developers` to the harness.

## What was measured

- Candidates: 141 skills — agents-cli `skills/` 7 (pin 2c39459), its LangChain extension template 1, google/skills `cloud` 131 and `developers` 2 (pin d5d9052). Apache-2.0 both. 5.7 MB, 778 files (614 md, 74 py, 44 yaml, 20 sh). Descriptions all within 1,024 chars; 67,698 chars together (average 480).
- Harness lint on a scratch tree with all of them: 2 refused for `allowed-tools` (`agent-platform-alert-configuration`, `cloud-run-alert-configuration`); none for an unreadable frontmatter line; 26 files carry an exec bit (the embed test refuses those under `skills/`); one name clash, `gke-basics` (the harness's pinned fork). Rendered context with 173 skills: 19,548 bytes of the 22,000 pin.
- Installed into a sandbox home and asked each CLI what it sees (`~/sandbox/audit-fix/imp-*.out`):
  - **Antigravity 1.2.14 lists 109 of 173 skills.** Two runs agree. It drops the longest descriptions first (kept average 382 chars, dropped 690; about 44,000 characters of name plus description survive). 16 of the harness's own 34 are dropped: `feature`, `review`, `golang`, `spec`, `terraform`, `cloud-run-deploy`, `data-gcp`, `requirements`, `testing-strategy`, `project-memory`, `dream`, `explain-project`, `operate-gcp`, `api-auth-gcp`, `json-api-style`, `mermaid-syntax`.
  - **Claude Code lists all 186 names but descriptions for only the first 41 alphabetically** (through `data-gcp`), about 19,800 characters; every later skill — `feature`, `golang`, `handoff`, `init-project`, `preflight`, `review`, `terraform` — is a bare name with no description, so nothing tells the model when to load it.
  - Today's 34 harness skills use 19,310 of those ~19,800 characters on Claude Code: the harness is already at the budget. The agents-cli seven plus the two developers skills alone add 5,749.
- Conclusion: installing the candidates as global skills breaks skill selection on both CLIs, for the harness's own workflow skills first. Even nine more global skills cost existing descriptions on Claude Code.

## The decisions

1. **Where the 141 live.** (a) Global skill directories, as asked — measured broken on both CLIs. (b) **A vendored catalog plus one router skill (recommended):** the 141 are vendored verbatim in the harness under a new directory that is not a skills directory (`catalog/google/<name>/`), installed to a non-skill location on each machine (`~/.config/sh-agent/catalog/google/`), and one small harness skill routes to them: its description triggers on Google Cloud, Gemini, Agent Platform, ADK and agents-cli tasks; its body tells the agent to find the matching entry in the catalog's index and read that SKILL.md and its references from the local copy. One description in the listing instead of 141; identical on both CLIs (R10); offline and pinned. (c) A curated handful installed globally and the rest fetched live by Google's own `finding-google-skills` from GitHub at use time — unpinned instruction text from the network at runtime, and it still costs budget. (d) Per-project only: the catalog is vendored and a project copies the few skills it needs into its own `.agents/skills/` — no global cost, but nothing is available until someone copies.
2. **What the router reads.** The local pinned copy (recommended: reviewed once, same on every machine, works offline) or Google's live `index.json` (always current, but unreviewed text and a network dependency inside a session).
3. **Budget room for the router.** Claude Code has about 500 characters left. A router description of at most 300 characters fits; recommended to also shorten `cloud-run-deploy`'s description (the longest) by the same amount so the listing keeps headroom, and to add a test that pins the sum of all skill descriptions under the measured budget so the next skill drop cannot silently push one out.
4. **The agents-cli seven and the LangChain template skill.** Into the catalog with the rest (recommended), or installed globally (costs about 3,900 characters and pushes six harness descriptions out of Claude Code's listing). The LangChain skill is a project template's own skill ("use when working in this project"); recommended: catalog it, never global.
5. **`gke-basics`.** Keep the harness fork installed as it is; upstream's current `gke-basics` and its 34 `gke-*` siblings go into the catalog, where names cannot clash. Retiring the fork in favour of the catalog's GKE set is a later, separate choice.
6. **The two `allowed-tools` skills.** In a catalog they are plain files the router reads, not skills a CLI loads, so the key grants nothing; kept verbatim, and the installer's permission-key lint does not apply to the catalog directory. (As global skills they would have to be forked without the key or refused.)
7. **Executable bits.** Stored 0644 (the embed test requires it); the mode is the only local change, recorded once in the catalog manifest; the router tells the agent to run scripts through their interpreter.
8. **Project opt-in.** A later phase: a way to copy a catalog skill into a project's `.agents/skills/` when the project lives on that product (so it loads natively there). Recommended: defer; the router covers the need first.
9. **REQUIREMENTS.** R5 lists the directories the installer maps; a catalog directory is a new kind, so R5 gains one sentence — the exact text needs the operator's yes and an ADR. No delivery path changes; the README map is untouched. Release binaries and the zip grow by about 5.7 MB.
10. **Provenance and refresh.** One manifest for the catalog (repository, path, pin, licence, date, the mode-only local change) instead of 141 rows in `VENDORED.md`, with a test that every catalog entry's files match `git archive <pin>` content (the BUG-111 rule, automated), and a documented refresh: fetch at a new pin, read the diff, re-run the security slice on it.
11. **Release rule.** "One Antigravity run of each changed skill" covers harness skills; catalog entries are data. Recommended: the router runs on both CLIs in the matrix, plus a sample of six catalog entries per release (one per family), not 141.
12. **Security.** A security review slice over the imported content before it lands (instruction-shaped text, `curl | sh`, credentials in samples, scripts that write outside the project), and over the diff at each refresh.
13. **Version.** MINOR after this lands (v2.1.0), together with the eleven gate lows and the timer fix already on `main`.

## Work shape if 1(b) is chosen

Feature skill: plan file with phases — (1) REQUIREMENTS R5 text, ADR; (2) installer: the catalog directory kind, its install destination, manifest test, status and uninstall; (3) vendor the 141 at their pins, index file, manifest; (4) the router skill and budget test, description trims; (5) security slice and reviews; (6) sandbox runs on both CLIs (router picks the right entry for six tasks), docs, matrix; (7) release.

## Second measurement — the catalog nested inside one router skill (2026-10-01)

A scratch tree with one new skill, `skills/google-catalog/`, holding the 141 entries under `catalog/<name>/` (exec bits cleared) and a generated `catalog/INDEX.md` (36,435 bytes, one line per entry): the unchanged installer lists 35 skills, renders 15,155 bytes, installs without a refusal (no secret-shaped or symlinked path among the 778 files). Installed into a sandbox home (`~/sandbox/audit-fix/sbImp2`) and asked on both CLIs (`imp2-*.out`):

- Claude Code lists 48 skills (35 + its 13 built-ins), Antigravity 37 (35 + its 2): every harness skill present on both, `typescript`'s description quoted in full on both, `google-catalog` listed, and neither `gke-storage` nor `bigquery-basics` in either list — a `SKILL.md` two levels below a skill directory is not discovered as a skill by either CLI.
- So the catalog needs no new directory kind, no installer change and no REQUIREMENTS edit: the existing skill deployment carries it, and a skill may read its own files on both CLIs.
- First routing probe ("alerting on Cloud Run error rate, as Terraform"): Claude Code loaded the harness's own `operate-gcp` and its alerts reference — correct precedence, wrong probe; Antigravity was denied a shell command under a per-name grant list (BUG-141's shape). The second probe, three tasks the harness does not cover, is recorded in the plan's Phase 4 log.
