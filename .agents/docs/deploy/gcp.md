---
owner: main
date: 2026-09-03
target: gcp
---
# Deploy — GCP

**Not applicable.** The harness is a file copier into home directories; it has no image, service, or project. The guide exists so the convention every project follows (`docs/deploy/{local,remote,gcp}.md`) holds here too, and so an agent looking for a GCP path finds this answer instead of inventing one.

Projects that do deploy to GCP fill the template that `/init-project` ships (`skills/init-project/assets/agents-tree/docs/deploy/gcp.md`): recorded defaults, image, procedure via the `cloud-run-deploy` skill, secrets, verify, rollback.
