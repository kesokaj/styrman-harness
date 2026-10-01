---
owner: main
date: 2026-10-01
---
# Verification — the changed and new skills on Antigravity (2026-10-01)

Range reviewed: 2ba4a50 installed into the sandbox home `~/sandbox/audit-fix/sbAgy` (`sh-agent install --home … --vendor agy`, drift none, 133 files)

## Status: PASS — 24 of 24 rows

Procedure: `testing/skill-matrix.md` § Procedure (Antigravity line), run by `.agents/scratch/audit-fix/kit-agy-matrix.sh <row>` (not committed, ADR-010): `testing/kit-preamble.sh` sourced; a fresh clone of the demo service per row outside every checkout; `settings.json` written per run with the grants and the clone as the only trusted workspace; `agy -p "The project root is the current directory, <clone>; ignore any other checkout. <prompt>" --mode accept-edits --print-timeout 24m`; the row's check cell; the token copied in for each batch and deleted by name after it; `xargs -P 2`. A run counts only when stderr has no `auto-denied` line and `cli.log` no `invalid allow entry`.

| row | result | round | note |
|---|---|---|---|
| `typescript`, `frontend`, `golang`, `json-api`, `testing-strategy`, `operate`, `data`, `api-auth`, `gcloud-auth`, `gke`, `terraform`, `review-docs`, `explain-project` | PASS (13) | 1 | chat-only rows and rows whose artifact the agent wrote with the shell; each output names its skill; `review-docs` dispatched the `reviewer` and its report named `README.md:7` — the first Antigravity run of a derived agent in this plan besides the real-home smoke's verifier |
| `distill-from-log` | PASS | 2 | the record quotes the seeded `gap` timestamps |
| `requirements`, `spec`, `review-resume` | PASS | 3 | `review-resume`'s round 1 ended with the main agent announcing the `explorer` dispatch and stopping before the brief landed (the headless rule says the turn ends only after every artifact is on disk); round 3 wrote the brief |
| `bash`, `mermaid`, `docker`, `init-deploy`, `init-update`, `handoff`, `cloud-run` | PASS (6) | 5 | rounds 1–3: no output, `auto-denied`, on the per-name grant list; round 4 added `write_to_file(*)` on a wrong diagnosis (the grant store knows `command`, `read_file`, `write_file`, `mcp` only; `cli.log` skipped it as `ignoring invalid allow entry`), and only `handoff` passed, because that run happened to choose granted commands; the conversation databases then showed the denied item was a shell command every time — `node` (mermaid, trying to render), `docker`, `gcloud`, the `for d in …; do` loop of init-project Step 5b (both init rows), a multi-line `grep -rn` (bash) — so round 5 ran these six with `command(*)` in their clone (BUG-141). The `invalid_allow=6` count on their result lines is `cli.log`'s cumulative count from round 4; round 5 added none |

Result: 24 of 24 PASS. 17 rows on the documented per-name grant list (`handoff` on its second attempt), 6 with `command(*)`. Each PASS is the row's check cell of `testing/skill-matrix.md`: the artifact exists and has the required shape, or the chat output names the skill and the terms the cell lists. The kit's `skill_dir_in_db` probe stayed 0 on every row — Antigravity's conversation database records the skill's content, not its path — so skill LOADING is evidenced by the output naming the skill and by the artifact's shape, not by the database. Round 4's token was copied into the sandbox at 11:09 and deleted after the batch; round 5's likewise (`ls` of the token path after the batch: no such file).

Found on the way: BUG-141 (fixed in this change: the grant list in `testing/skill-matrix.md` step 2 and `testing/guide.md` § 4 now names the four actions and the database diagnosis, and the kit's round-5 line uses `command(*)` for the rows whose skill text holds a shell construct). The `review-resume` round-1 stop (main agent ending its turn while the explorer ran) is the behaviour `context/AGENTS.md`'s headless rule names; it passed on the third attempt without a change to the harness and is recorded here, not as a bug, because the rule already exists and the run count is 1 of 3.

## What was not exercised

The `init-memory`, `dream`, `dreamer`, `preflight` and `review` (code lens) rows on Antigravity; every row on the Mac; the real-home matrix (the real-home procedure is the three smoke checks of `testing/2026-10-01-verify-real-home.md`, 6 of 6 PASS on both CLIs).
