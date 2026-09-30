# Vendored third-party skills — provenance manifest

Skills in the first table are vendored verbatim from the open skills.sh ecosystem
(SKILL.md standard). Update by re-fetching from the pinned source and bumping the entry
here; never edit vendored content in place (fork it into a differently-named skill
instead). The second table lists the two forks that predate that rule.

| Skill | Source | Pin | License | Vendored | Update via |
|-------|--------|-----|---------|----------|------------|
| grilling | github.com/mattpocock/skills (`skills/productivity/grilling`) | 6654f6b | MIT | 2026-08-31 | `npx skills add mattpocock/skills -s grilling` |
| codebase-design | github.com/mattpocock/skills (`skills/engineering/codebase-design`) | 6654f6b | MIT | 2026-08-31 | `npx skills add mattpocock/skills -s codebase-design` |
| find-skills | github.com/vercel-labs/skills (`skills/find-skills`) | 435076e | MIT | 2026-08-31 | `npx skills add vercel-labs/skills -s find-skills` |
| anti-ui-slop | uizze.com `.well-known/skills` (agentskills discovery v0.2.0), v1.2.13 | sha256:c3dc1a03…09bc08 (digest verified) | MIT | 2026-08-31 | re-fetch archive from `https://uizze.com/.well-known/skills/index.json` |
| prompt-audit reference (one file, not a skill: `review-standards/references/prompt-audit/prompt-audit.md`, with `LICENSE.txt` beside it) | github.com/anthropics/skills (`skills/claude-api/shared/prompt-audit.md`) | 8a1541c (sha256:c33200da…1878f2) | Apache-2.0 | 2026-09-29 | `curl -sL https://raw.githubusercontent.com/anthropics/skills/<sha>/skills/claude-api/shared/prompt-audit.md`, then bump the pin and digest here |

## Apache-2.0 skills — forks, kept under their upstream names

These two are NOT verbatim: each was rewritten here from an upstream base and keeps
the upstream name. Their provenance was "not recorded when vendored" until the project
review of 2026-09-30 (security finding 4) traced it by diffing each against upstream's history
(`git archive <commit> <path> | diff -r`); a third, `gcloud-auth-verification`
by google (v1, Apache-2.0), was merged into `gcloud-setup` on 2026-09-13 — its three
rules paraphrased under "Before asking the user to log in" — and removed. Each declares
its licence and publisher in its own `SKILL.md` frontmatter, which is the authority for
attribution. A re-fetch from the pin is a diff to read, never a copy to apply over the
local text.

| Skill | Publisher / author | Derived from | Pin | Local changes | License |
|-------|--------------------|--------------|-----|---------------|---------|
| gke-basics | Google Cloud | github.com/google/skills (`skills/cloud/gke-basics`, the one skill it was before upstream split it into `gke-*` sub-skills at 5d9d054 on 2026-06-23) | 6f0b877 (2026-04-23; the tree is unchanged upstream through 2026-05-10) | `SKILL.md` only: the description rewritten with this harness's cross-references (`cloud-run-deploy`, `docker`), a "When NOT to use" line added, the closing pointer to the Developer Knowledge MCP replaced by the GKE docs URL; `references/` and `assets/` verbatim | Apache-2.0 |
| terraform | Anton Babenko | github.com/antonbabenko/terraform-skill (`SKILL.md`, `references/`) | 2271bc4 (tag v1.6.0) | `SKILL.md` condensed (prose and examples to tables, ~140 lines shorter), renamed from `terraform-skill`, the description rewritten with this harness's cross-references (`gke-basics`, `docker`, `gcloud-setup`); `references/` verbatim but for "Claude" → "The AI" / "Gemini CLI" in two lines | Apache-2.0 |

Apache-2.0 permits redistribution with attribution, which the table above provides.
It also asks that a copy of the licence travel with the work: this repository is not
published, and if that changes, add `LICENSE` files for these two before it does.

Notes:
- `grill-me` (a launcher for `grilling`) and `skill-repair` (repairs a `skills`-CLI manifest the
  harness does not use) were pruned 2026-09-03 after the sandbox matrix: `/grilling` does the same
  job on both CLIs, and `sh-agent adopt` replaced the manifest repair. `find-skills` (vendored verbatim)
  still names `skill-repair` in its description; harmless — the model finds no such skill — and fixed
  by the next re-fetch, not by editing the vendored file.
- `find-skills` is verbatim again at `435076e` (restored 2026-09-28). Its old local Step 6 edit, which routed installs into this repository, now lives in the global context's § Global vs project.
- `anti-ui-slop` includes its own LICENSE and NOTICE from the upstream archive; the
  optional paid UIZZE MCP is NOT configured — the free playbooks work standalone.
- To add a skill: run `npx skills add <owner/repo> -s <skill> -y` from the repo root (it lands
  in `skills/` through the `.agents/skills` symlink), or copy any directory holding a `SKILL.md`
  into `skills/`, then `sh-agent install`. Add a row above. `sh-agent install` deploys exactly
  the repo's skills and removes only what it installed itself, so a skill installed by hand
  into `~/.claude/skills` survives but is neither versioned nor deployed to Antigravity —
  `sh-agent adopt <skill>` moves it in here.
