---
owner: reviewer
date: 2026-09-15
---
# Review — sweep-skills-1   (2026-09-15)

Range reviewed: working tree at HEAD d1ff612 plus uncommitted edits (only `skills/bug/SKILL.md:29` changed in scope). This was not a diff review. Every file in these nine directories was read in full: skills/bash/, skills/bug/, skills/cloud-run-deploy/, skills/distill/, skills/docker/, skills/feature/, skills/feedback/, skills/frontend/, skills/gcloud-setup/. Lens: `--docs --project`. They were checked against README.md, AGENTS.md, context/, agents/, skills/retro, skills/handoff, skills/preflight, skills/init-project assets, .agents/docs/testing/sh-agent/guide.md, .agents/docs/testing/skill-matrix.md, .agents/docs/deploy/*.md, the newest changes/ entry (2026-09-15-research-beads-elephant-goldfish.md), internal/harness/adopt.go, internal/harness/repo.go and cmd/sh-agent/main.go. `.agents/docs/specs/` is empty and there is no `.agents/rules/`.

## Verdict: REQUEST CHANGES

The worst problem is `skills/cloud-run-deploy/SKILL.md:3-9`. It calls itself the "End-to-end release workflow" and triggers on "release" and "deploy this service", but it never mentions `/preflight`. The global context (`context/AGENTS.md:45`) says those words mean `/preflight`, and `skills/feature/SKILL.md:78` calls preflight "the only path that gates, tags, and deploys". So a "release this" can load the deploy recipe directly and skip the Definition of Done, the reviewer and security passes, the SHIP report and the tag. Next worst: `/distill` writes feedback entries with no `project` field, which trips the retro's promotion grep in the same project (a false harness promotion). `/distill` also tells the user to run `sh-agent adopt <name>` from `.agents/skills/`, a path adopt never looks in. The Makefile reference pushes a native-arch image with an invented region. The frontend references still describe a Tailwind v3 `tailwind.config.ts` / `index.html` setup that the rewritten frontend skill (2026-09-13) has dropped. The bug, feedback, gcloud-setup auth rules, docker health-check and design references are otherwise sound. The BUG-018 fix at `skills/bug/SKILL.md:29` is in place, and no living doc still carries the old wording.

## Findings

- **[medium / confidence: medium]** divergent: `skills/cloud-run-deploy/SKILL.md:3-9` vs `context/AGENTS.md:45`, `skills/preflight/SKILL.md:3` and `skills/feature/SKILL.md:78`.
  - Text: "End-to-end release workflow … Use when the user says 'deploy this service', 'push to Cloud Run', 'release' … even a 'just redeploy it' request."
  - Why it fails: the context says "Ship", "deploy", "release" from the user means `/preflight`, and preflight triggers on the same words. Nothing in cloud-run-deploy sends the agent back to preflight, so the gates can be bypassed. Confidence is medium because it depends on which skill the model picks.
  - Should say: requests to ship, release or deploy go through `/preflight`. This skill is the procedure that preflight step 5 follows (`skills/preflight/SKILL.md:87`), loaded directly only for a local build/push rehearsal or when `/preflight` calls it. Drop "release" from the triggers.
- **[medium / confidence: high]** divergent: `skills/distill/SKILL.md:158-168` vs `skills/feedback/SKILL.md:13,43` and `skills/retro/SKILL.md:65`.
  - Text: the feedback entry distill appends has no `"project"` field.
  - Why it fails: feedback says every entry carries `project`, which is what makes cross-project matching work. The retro promotion check is `grep rule_ref … | grep -vE '"project": ?"<this project>"'`, and an entry with no project is never filtered out. So a later retro in the same project that cites `skill:<name>` gets one "cross-project" hit and adds a `Scope: harness` promotion proposal.
  - Should say: add `"project": "<absolute project path>"` to the example, as in `skills/feedback/SKILL.md:43`.
- **[medium / confidence: high]** stale: `skills/distill/SKILL.md:137`.
  - Text: "`sh-agent adopt <name>` copies it from `.agents/skills/<name>` into the harness repo."
  - Why it fails: `internal/harness/adopt.go:16-21` resolves a bare name as `<name>` relative to the working directory, then `$HOME/.agents/skills/<name>`, `~/.claude/skills/<name>` and `~/.gemini/config/skills/<name>`. It never looks in the project's `.agents/skills/`, so from the project root it fails with "no skill directory with a SKILL.md found".
  - Should say: "`sh-agent adopt .agents/skills/<name>` (a path) from the project root copies it into the harness repo."
- **[medium / confidence: high]** divergent: `skills/bash/references/makefile-patterns.md:128-139,193,196` vs `context/dev-setup.md:4,10,11`.
  - Text: `docker-build` runs a native `docker build` and `docker-push` pushes that image. `REGISTRY := us-docker.pkg.dev/...`, the tag comes from `git describe` and is also tagged `latest`, and the deploy targets hard-code `--region=us-central1`.
  - Why it fails: on Apple Silicon this pushes an arm64 image, and Cloud Run fails with `exec format error`. The dev setup says to always build amd64 for Google, use `${REGION}-docker.pkg.dev/${PROJECT}/${REPO}/${SERVICE}:${GIT_SHA}`, and never assume or invent a region.
  - Should say: build with `docker buildx build --platform linux/amd64 --load`, use `$(REGION)-docker.pkg.dev`, tag with `git rev-parse --short HEAD`, and take the region from `docs/deploy/gcp.md`, or point at `cloud-run-deploy`.
- **[medium / confidence: medium-high]** outdated: `skills/frontend/references/theme-system.md:341-398` and `:217-239` vs `skills/frontend/SKILL.md:34,40,61`.
  - Text: the tokens are wired through `tailwind.config.ts` (`theme.extend.colors`), and the FOWT script goes "to your `index.html`" with `<div id="root">`.
  - Why it fails: the skill says `vite.config.ts` is "The ONLY config file", pins `tailwindcss` 4.3, and makes `__root.tsx` own `<head>`. Tailwind v4 does not read a JS config unless the CSS declares `@config`, so `bg-primary`/`text-heading` are never generated and the "one class, both themes" claim fails silently. TanStack Start has no `index.html`. The references are unchanged since e521826 (2026-09-03); SKILL.md was rewritten in e5833a3 (2026-09-13).
  - Should say: tokens go in `src/index.css` through the v4 `@theme` block, and the FOWT script goes in the `__root.tsx` head.
- **[low-medium / confidence: high]** narrowed: `skills/frontend/SKILL.md:77-82,42` → `skills/frontend/references/component-patterns.md:199-200,219`.
  - Text: "For forms that submit to server functions, pair with `createServerFn`", and the example calls `fetch("/api/signup")` inline.
  - Why it fails: the reference drops the SKILL's qualifier that `createServerFn` exists only for SSR projects, and bypasses `src/lib/api.ts`, "the only fetch layer".
  - Should say: SPA (the default) submits through `src/lib/api.ts`; `createServerFn` only in a project that enables SSR.
- **[low-medium / confidence: medium]** divergent: `skills/frontend/references/component-patterns.md` (e.g. `:22-32,122,224`), `skills/frontend/references/animation-recipes.md` (e.g. `:37-40,148,166`) and `skills/frontend/references/accessibility.md:21-37` vs `skills/frontend/references/theme-system.md:3,7,411` and `skills/frontend/SKILL.md:25`.
  - Text: every pattern hard-codes a dark zinc/amber palette, and the contrast table covers "Dark Mode Safe Pairings" only.
  - Why it fails: the skill says light mode is the default and components reference tokens. Copied patterns render dark surfaces in light mode and ignore `data-theme`.
  - Should say: use token classes or `var(--…)` in the patterns, and give contrast pairs for both token sets.
- **[low / confidence: high]** stale: `skills/frontend/references/accessibility.md:258`.
  - Text: "In main.tsx (development only)".
  - Why it fails: the project structure in `skills/frontend/SKILL.md:36-40` has `router.tsx` and `__root.tsx`, and no `main.tsx`.
  - Should say: load axe from `__root.tsx` behind `import.meta.env.DEV`.
- **[low / confidence: high]** narrowed: the gate line in `87a862d^:context/AGENTS.md:59` → `context/AGENTS.md:16` and `skills/bug/SKILL.md:49`.
  - Old text: "a review or design loop that hits its round cap" is "never skipped on the strength of an 'autonomous' instruction".
  - Why it fails: the qualifier survives only for features (`skills/feature/SKILL.md:15`). For a bug, the review-cap gate lives at `context/AGENTS.md:16` ("three rounds without convergence means stop and ask the user") with no autonomy clause, and the bug skill has none. This is the same pattern as BUG-018, for the gate the fix did not cover.
  - Should say: `skills/bug/SKILL.md` says a review loop that hits its cap waits for the user whatever autonomous instruction is in effect.
- **[low / confidence: high]** narrowed: `context/AGENTS.md:16,19` → `skills/feature/SKILL.md:82` (and `skills/bug/SKILL.md:49` for the first half).
  - Text: feature Step 6 says "Steps 1–5 of the global Definition of Done: … commit and push. Then a report in chat". The bug skill has the same review row.
  - Why it fails: both restate the review row without the global rule that a change touching an interface, command or deploy step also gets `/review --docs`. The feature skill also drops step 5's `/handoff` and "a handoff on disk that says nothing is queued".
  - Should say: add "`/review --docs` over the range when an interface, command, or deploy step changed" to both, and "then `/handoff`" to feature Step 6.
- **[low / confidence: high]** stale: `skills/feature/SKILL.md:19`.
  - Text: "Check … `docs/bugs/INDEX.md`".
  - Why it fails: the same file defines `<docs-root>` as `.agents/docs` first (`:11`), and `skills/bug/SKILL.md:29` uses `<docs-root>/bugs/INDEX.md`. In an `.agents/docs` project that path does not exist, so the check silently finds nothing.
  - Should say: `<docs-root>/bugs/INDEX.md`.
- **[low / confidence: medium]** divergent: `skills/docker/SKILL.md:105` vs `context/dev-setup.md:8` (conditioned in cceee54).
  - Text: "Linux: `sudo -n systemctl start docker`".
  - Why it fails: the dev setup allows a root step only when the operator says so AND the project's `.claude/settings.json` records the standing permission; otherwise the agent hands back the exact command.
  - Should say: "Linux: the root step `systemctl start docker`, run only under the dev-setup root-step rule, otherwise handed to the operator."
- **[low / confidence: medium]** outdated: `skills/docker/SKILL.md:108`.
  - Text: "`CGO_ENABLED=0` keeps it a pure Go cross-compile in practice".
  - Why it fails: the Dockerfile at `:57` is `FROM golang:1.24 AS build`, with no `--platform=$BUILDPLATFORM` and no `GOARCH=$TARGETARCH`. Under `buildx --platform linux/amd64` the Go compile runs under QEMU, which `skills/cloud-run-deploy/SKILL.md:54-55,116` also says.
  - Should say: either add `FROM --platform=$BUILDPLATFORM golang:1.24 AS build`, `ARG TARGETOS TARGETARCH` and `GOOS=$TARGETOS GOARCH=$TARGETARCH go build`, or say plainly that the build is emulated.
- **[low / confidence: medium]** divergent: `skills/docker/SKILL.md:50-51` vs `context/AGENTS.md:57`.
  - Text: `docker run --rm -p 8080:8080 <app>:dev &` followed immediately by `curl -fsS http://localhost:8080/healthz`.
  - Why it fails: the curl races the container start. A process already on 8080 can answer instead, which the hard rule forbids ("probes the process it started rather than 'something answers'"), and the backgrounded container is never stopped.
  - Should say: `docker run -d --name <app>-smoke …` (fail if the port is bound), retry the curl against it, then `docker rm -f <app>-smoke`.
- **[low / confidence: high]** stale: `skills/docker/references/setup.md:12`.
  - Text: `web/.output` in `.dockerignore`.
  - Why it fails: that is the Nitro/Vinxi output directory, and `skills/frontend/SKILL.md:43,59` says the build writes `dist/client/` and there is no Nitro.
  - Should say: drop `web/.output`. Also do not ignore `web/dist/client` when `embed.go` embeds it into the image.
- **[low / confidence: medium]** divergent: `skills/docker/SKILL.md:84-86` vs `skills/frontend/SKILL.md:43,96-98`.
  - Text: the docker skill adds a `node:<pinned>` build stage that produces the bundle. The frontend skill says `dist/client/` is "committed when Go embeds it" and "`go build` needs no Node".
  - Why it fails: the two skills give opposite answers on whether the image build needs Node.
  - Should say: pick one. Either a committed bundle and no node stage, or a node stage and an uncommitted `dist/`.
- **[low / confidence: high]** divergent: `skills/bash/references/testing-linting.md:71-80,282-308`, `skills/gcloud-setup/SKILL.md:125-134` and `skills/frontend/references/accessibility.md:239` vs `context/AGENTS.md:45`.
  - Text: GitHub Actions workflows (including the unpinned `ludeeus/action-shellcheck@master`), "Service Account Keys (CI/CD only)" with `--key-file`, and "in CI via `@axe-core/react`".
  - Why it fails: the context says "no CI, no hooks, ever". The key-file path also contradicts `skills/gcloud-setup/SKILL.md:109` ("Never fix credentials by code or alternative tools").
  - Should say: remove the CI sections, or point to the local `make lint-sh`/`bats` targets.
- **[low / confidence: high]** incorrect statement: `skills/bash/references/bash-fundamentals.md:84-88,357`.
  - Text: `${NAME:-default_value}` is labelled "Use default if unset" (it applies when unset OR empty). `${NAME:-}` is labelled "Use default if unset or empty" (it supplies no default; it only suppresses `set -u`). `"${array[@]}"` is described as "word-split" (it keeps each element as one word).
  - Should say: `${NAME-default}` applies when unset; `${NAME:-default}` when unset or empty; `"${array[@]}"` gives one word per element.
- **[low / confidence: high]** stale: `.agents/docs/testing/skill-matrix.md:45`.
  - Text: the check still accepts `gcloud-auth-verification`.
  - Why it fails: that skill was merged into gcloud-setup and removed on 2026-09-13 (`skills/VENDORED.md:16`). The check still passes, on the other alternative.
  - Should say: `grep -qi "gcloud-setup" "$OUT"`.
- **[low / confidence: high]** stale: `skills/distill/SKILL.md:4`.
  - Text: `argument-hint: "[--from-log]"`.
  - Why it fails: the Usage block at `:15` documents `/distill --global` as well.
  - Should say: `"[--global] [--from-log]"`.

### Unverified
- `skills/feature/SKILL.md:72`: "BUG-033/041 in styrman-crew": that project was not checked.
- `skills/feedback/SKILL.md:81`: "the 2026-09-04 and 2026-09-08 retros": not in this repo's `retros/`, which holds only 2026-09-14. They may be in other projects or the proposals queue.
- `skills/frontend/SKILL.md:60-62`: the version pins "by the reference project on 2026-09-13".
- `skills/docker/references/health-checks.md:41`: the `--startup-probe httpGet.path=/healthz` flag syntax.
- `skills/cloud-run-deploy/SKILL.md:91`: whether `gcloud run services logs read` is GA.
- `skills/distill/SKILL.md:173`: a Claude Code built-in skill-creator.
- The Tailwind v4 JS-config behaviour in the theme-system finding comes from framework knowledge, not a build run here.

## Classification (named living docs)
| File | Class | Proof |
|---|---|---|
| skills/bash/SKILL.md | accurate | `:39-46` core rules agree with the references and the context; `:53-59` reference paths exist |
| skills/bash/references/bash-fundamentals.md | outdated | `:84-88` wrong parameter-expansion semantics |
| skills/bash/references/testing-linting.md | outdated | `:71-80,282-308` CI workflows vs `context/AGENTS.md:45` |
| skills/bash/references/makefile-patterns.md | outdated | `:133-139,193` native push, invented region vs `context/dev-setup.md:4,10` |
| skills/bug/SKILL.md | outdated | `:49` review row lacks `/review --docs` (`context/AGENTS.md:16`); `:29` BUG-018 fix present |
| skills/cloud-run-deploy/SKILL.md | outdated | `:3-9` release triggers vs `context/AGENTS.md:45` |
| skills/distill/SKILL.md | outdated | `:137` vs `internal/harness/adopt.go:16-21`; `:158-168` no `project` vs `skills/feedback/SKILL.md:43` |
| skills/docker/SKILL.md | outdated | `:105` vs `context/dev-setup.md:8`; `:108` vs `:57` |
| skills/docker/references/health-checks.md | accurate | `:40-42` agrees with cloud-run-deploy; `:8` `/healthz` matches `skills/cloud-run-deploy/SKILL.md:90` |
| skills/docker/references/setup.md | outdated | `:12` `web/.output` vs `skills/frontend/SKILL.md:43,59` |
| skills/feature/SKILL.md | outdated | `:19` `docs/bugs/INDEX.md` vs `:11`; `:82` omits `/handoff` vs `context/AGENTS.md:19` |
| skills/feedback/SKILL.md | accurate | `:54` vocabulary matches `skills/retro/SKILL.md:49`; `:81` matches `skills/handoff/SKILL.md:34,103` |
| skills/frontend/SKILL.md | accurate | `:77-99` SPA/SSR model consistent with `context/dev-setup.md:6`; current as of e5833a3 |
| skills/frontend/references/design.md | accurate | `:6` points at "Apps & Dashboards" (`skills/frontend/SKILL.md:101`); `:112` link exists |
| skills/frontend/references/theme-system.md | outdated | `:346` `tailwind.config.ts`, `:217` `index.html` vs `skills/frontend/SKILL.md:34,40,61` |
| skills/frontend/references/component-patterns.md | outdated | `:200,219` vs `skills/frontend/SKILL.md:42,80` |
| skills/frontend/references/accessibility.md | outdated | `:258` `main.tsx`; `:21-37` dark-only pairs vs `skills/frontend/references/theme-system.md:3` |
| skills/frontend/references/animation-recipes.md | outdated | `:148,166` hard-coded palette vs `skills/frontend/references/theme-system.md:7,411` |
| skills/gcloud-setup/SKILL.md | outdated | `:125-134` CI key files vs `:109` and `context/AGENTS.md:45`; `:106-112` auth rules accurate |
| .agents/docs/testing/skill-matrix.md (pointed at by the scope's matrix rows) | outdated | `:45` `gcloud-auth-verification` vs `skills/VENDORED.md:16` |

Also compared, and nothing in them contradicts the scope: `README.md:88` (the feature and bug skills load themselves; true per their descriptions), `.agents/docs/deploy/gcp.md:10` (template plus cloud-run-deploy; the template exists), `.agents/docs/deploy/local.md`, `.agents/docs/deploy/remote.md`, `.agents/docs/testing/sh-agent/guide.md`, and the newest `changes/` entry.

## Re-check commands
```bash
git grep -n "'release'" -- skills/cloud-run-deploy/SKILL.md; git grep -n 'preflight' -- skills/cloud-run-deploy/SKILL.md
git grep -n -A10 '"ts": "2026-05-25T14:00' -- skills/distill/SKILL.md | grep -c '"project"'
git grep -n 'adopt <name>` copies it from' -- skills/distill/SKILL.md
git grep -n 'us-docker.pkg.dev\|us-central1\|\$(DOCKER) build' -- skills/bash/references/makefile-patterns.md
git grep -n 'tailwind.config.ts\|index.html' -- skills/frontend/references/
git grep -n 'createServerFn\|fetch("/api' -- skills/frontend/references/component-patterns.md
git grep -n 'zinc-950\|amber-400' -- skills/frontend/references/ | wc -l
git grep -n 'main.tsx' -- skills/frontend/references/accessibility.md
git grep -n -i 'autonomous\|round cap\|hits its cap' -- skills/bug/SKILL.md
git grep -n 'review --docs\|/handoff' -- skills/feature/SKILL.md skills/bug/SKILL.md
git grep -n '`docs/bugs/INDEX.md`' -- skills/feature/SKILL.md
git grep -n 'sudo -n' -- skills/
git grep -n 'pure Go cross-compile\|BUILDPLATFORM' -- skills/docker/
git grep -n 'docker run --rm -p 8080:8080 <app>:dev &' -- skills/docker/SKILL.md
git grep -n 'web/.output' -- skills/
git grep -n 'node:<pinned>' -- skills/docker/SKILL.md; git grep -n 'committed when Go embeds' -- skills/frontend/SKILL.md
git grep -n -i 'github actions\|key-file\|in CI' -- skills/bash skills/gcloud-setup skills/frontend
git grep -n 'Use default if unset\|word-split' -- skills/bash/references/bash-fundamentals.md
git grep -n 'gcloud-auth-verification' -- .agents/docs/testing/skill-matrix.md
git grep -n 'argument-hint' -- skills/distill/SKILL.md
```

## Criteria
1. Dead code: n/a (docs lens)
2. Flow tracing: n/a (docs lens)
3. Standards conformance: n/a (docs lens)
4. Directory layout: n/a (docs lens)
5. Simplicity (DIY-except-security): n/a (docs lens)
6. Decoupling: n/a (docs lens)

## Not examined
- Every named file was read in full; nothing in the scope was skipped.
- Vendored skills (`skills/VENDORED.md`) were not searched, per the brief.
- Other projects' trails (styrman-crew, argolis-devbox mirrors) were not opened; claims citing them are listed as unverified.
- No command in the skills was executed (no gcloud, docker, buildx, Tailwind build). Correctness of the external CLI and framework behaviour rests on reading, not running.
