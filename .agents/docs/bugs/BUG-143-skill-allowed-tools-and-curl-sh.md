---
owner: bugs
date: 2026-10-01
id: BUG-143
status: graduated (rule 1 to `internal/harness/inventory.go` `LoadSkills` + `TestLintRefusesAllowedTools`)
area: skills (bash, docker, golang, gcloud-setup)
---
# BUG-143: skills shipped `allowed-tools: Bash` and unverified install fetches (`curl | sh`), missed by four reviews

## Symptom
The v2.0.0 gate's security review of the harness text (`reviews/2026-10-01-security-v1.16.2-0b1b5b8-harness-text.md`) found two mediums that three phase reviews and the Audit III technical-skills pass had not:
1. `skills/bash/SKILL.md:4` and `skills/docker/SKILL.md:9` carried `allowed-tools: Bash`: on Claude Code every shell command runs without a prompt while the skill is active, and these are the skills most often loaded while reading untrusted scripts and Dockerfiles.
2. `skills/golang/references/service.md:404` (added in Phase 6) installed golangci-lint with `curl -sSfL https://golangci-lint.run/install.sh | sh`; `skills/gcloud-setup/SKILL.md:49-50` fetched the gcloud tarball with no checksum step.

## Repro
`grep -rn '^allowed-tools' skills/*/SKILL.md`; `grep -rn 'install.sh | sh' skills/`.

## Root cause
The reviews (Phase 6 code/docs/prompts lenses, the Phase 6 security review) read skill BODIES for the recipes they teach; none treated the frontmatter as a permission surface or a `curl | sh` line as a dependency fetch. `allowed-tools` was inherited from the first version of both skills (2026-08), survived the Phase 6 rewrite, and was never questioned. The install one-liner was copied from the tool's own install page, which recommends it.

## Fix
Commit 89f3308: both `allowed-tools` lines deleted; golangci-lint installed from the pinned release asset verified with `sha256sum -c --ignore-missing` against the release's checksums file; the gcloud tarball step checks the SHA256 from the install page before `install.sh`.

## Prevention
- A skill's frontmatter is a permission surface: `allowed-tools` is never set in a harness skill (the operator's permission mode decides). Graduated: `LoadSkills` in `internal/harness/inventory.go` refuses any skill whose frontmatter sets it ("allowed-tools is never set in a harness skill; the operator's permission mode decides (BUG-143)"), pinned by `TestLintRefusesAllowedTools` in `internal/harness/harness_test.go`.
- Every install instruction in a skill fetches a pinned artifact and verifies a checksum or signature before running it; `curl | sh`, `wget -O- | bash` and unpinned `npx -y` are findings. The security reviewer's brief for skill assets names them. Not yet mechanised: stays a review rule. Graduation target: `review-standards` § 8 or the security agent's definition.

## Related
BUG-111 (vendored pins), D40 (no `:latest`). No existing rule covered either item (not a repeat). Fix commit 89f3308; review `reviews/2026-10-01-security-v1.16.2-0b1b5b8-harness-text.md`.
