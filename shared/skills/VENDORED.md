# Vendored third-party skills — provenance manifest

Skills below are vendored verbatim from the open skills.sh ecosystem (SKILL.md standard).
Update by re-fetching from the pinned source and bumping the entry here; never edit
vendored content in place (fork it into a differently-named skill instead).

| Skill | Source | Pin | License | Vendored | Update via |
|-------|--------|-----|---------|----------|------------|
| grill-me | github.com/mattpocock/skills (`skills/productivity/grill-me`) | 6654f6b | MIT | 2026-08-31 | `npx skills add mattpocock/skills -s grill-me` |
| grilling | github.com/mattpocock/skills (`skills/productivity/grilling`) | 6654f6b | MIT | 2026-08-31 | `npx skills add mattpocock/skills -s grilling` |
| codebase-design | github.com/mattpocock/skills (`skills/engineering/codebase-design`) | 6654f6b | MIT | 2026-08-31 | `npx skills add mattpocock/skills -s codebase-design` |
| find-skills | github.com/vercel-labs/skills (`skills/find-skills`) | 435076e | MIT | 2026-08-31 | `npx skills add vercel-labs/skills -s find-skills` |
| anti-ui-slop | uizze.com `.well-known/skills` (agentskills discovery v0.2.0), v1.2.13 | sha256:c3dc1a03…09bc08 (digest verified) | MIT | 2026-08-31 | re-fetch archive from `https://uizze.com/.well-known/skills/index.json` |

## Apache-2.0 skills — attribution recorded from their own frontmatter

These four were vendored without a row here. Each declares its licence and publisher
in its own `SKILL.md` frontmatter, which is the authority for attribution; the fetch
origin was not recorded at the time and is marked unknown rather than guessed.

| Skill | Publisher / author | Version | License | Source |
|-------|--------------------|---------|---------|--------|
| gcloud-auth-verification | google | v1 | Apache-2.0 | not recorded when vendored |
| gke-basics | Google Cloud | 1.0.0 | Apache-2.0 | not recorded when vendored |
| skill-repair | google | v1 | Apache-2.0 | not recorded when vendored |
| terraform | Anton Babenko | 1.6.0 | Apache-2.0 | not recorded when vendored |

Apache-2.0 permits redistribution with attribution, which the table above provides.
It also asks that a copy of the licence travel with the work: this repository is not
published, and if that changes, add `LICENSE` files for these four before it does.
Re-vendoring any of them should fill in the Source column at the same time.

Notes:
- `grill-me` is a user-invocable launcher (`disable-model-invocation: true`) whose engine is
  the `grilling` skill — they ship as a pair. Its body invokes the Claude-Code-specific
  Skill tool, so `/grill-me` is inert on Antigravity — use the `grilling` skill there directly.
- `anti-ui-slop` includes its own LICENSE/NOTICE/MANIFEST from the upstream archive; the
  optional paid UIZZE MCP is NOT configured — the free playbooks work standalone.
- Installing skills with `npx skills add -g` writes directly into agent dirs and will be
  wiped by the harness deploy — `harness install` mirrors this directory and prunes
  anything the repo does not have. To keep a skill, vendor it here.
