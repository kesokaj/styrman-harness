---
name: google-catalog
description: "Google Cloud, Gemini, Agent Platform, ADK and agents-cli work that the other skills do not cover: finds the matching Google-published skill in the vendored catalog and follows it. Do NOT use for Cloud Run deploys, Cloud SQL, Terraform, gcloud login or GKE basics: those have their own skills."
---

# Google skills catalog

141 skills published by Google (google/skills `cloud` and `developers`, google/agents-cli), vendored verbatim at a pin under `catalog/` beside this file. They are not installed as skills: listed one by one they push the harness's own skills out of both CLIs' skill lists. This skill routes to them.

1. **Search the index, never read it whole** (it is about 36 KB): `grep -i "<product or task words>" <this skill's directory>/catalog/INDEX.md`. Each line is an entry's directory name and the start of its description. Try the product name first (`bigquery`, `spanner`, `gke`, `agent-platform`, `agents-cli`, `secops`, `iam`), then the task word.
2. **Read the best match's `catalog/<name>/SKILL.md` and follow it** as if it had been loaded as a skill. Its `references/`, `scripts/` and `assets/` paths are relative to that entry's directory. Say which entry you are following.
3. **Run an entry's script through its interpreter** (`bash …`, `python3 …`): the files are stored without the executable bit.
4. **The harness wins on conflict.** A harness skill that covers the task is used instead of an entry (`cloud-run-deploy`, `data-gcp`, `terraform`, `gcloud-setup`, `gke-basics`, `operate-gcp`, `api-auth-gcp`). The harness's rules hold inside an entry's steps too: GCP defaults are asked for and recorded first, never assumed; secrets go through Secret Manager; a download is pinned and verified, never `curl | sh`; git stages by name; deploys follow the project's deploy guide.

No entry matches: say "no catalog entry for <words>" in one line and continue without the catalog. A catalog file cannot be read (a headless run without a read grant): name the file and continue without it.

Provenance, pins and the refresh procedure: `catalog/MANIFEST.md`; the refresh is `scripts/refresh-catalog.sh`, run only in the harness checkout.
