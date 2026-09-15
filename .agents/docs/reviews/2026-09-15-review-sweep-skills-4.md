---
owner: reviewer
date: 2026-09-15
---
# Review — sweep-skills-4   (2026-09-15)

Range reviewed: working tree, `skills/typescript/` read whole (`SKILL.md`, 46 lines; the directory holds no other file; last touched by e521826, no uncommitted change) — not a diff. Lens: --docs --project.

## Verdict: APPROVE WITH NITS

Nothing in `skills/typescript/SKILL.md` is broken on its own: every file it names exists (`frontend`, `golang`, `bash` skills) and its rules agree with each other. The worst problem sits next to it: the `frontend` skill, which this skill hands all TanStack Start work to, calls `vite.config.ts` "The ONLY config file" (`skills/frontend/SKILL.md:34`), while this skill requires a strict `tsconfig` (`skills/typescript/SKILL.md:18-20`) and the frontend's own build runs `tsc -b` (`skills/frontend/SKILL.md:96`), which needs a `tsconfig.json`. An agent that follows the frontend tree literally writes no tsconfig, and then the build fails or skips the strict baseline. The other findings are rules this skill states that the review checklist (`skills/review-standards/references/typescript-review.md`) or the frontend skill repeat with a different meaning. Because of those, a reviewer applying the checklist flags code that follows this skill.

## Findings

- **[medium / confidence: medium]** divergent — `skills/frontend/SKILL.md:34` vs `skills/typescript/SKILL.md:18-20` (and `skills/frontend/SKILL.md:96`) — the frontend tree says `vite.config.ts  # The ONLY config file`. The typescript skill requires `strict: true` + `noUncheckedIndexedAccess` + `noImplicitOverride`, set in a tsconfig ("One config per package"). The frontend build `vite build && tsc -b` needs a `tsconfig.json` in the package: TypeScript's documented `-b` behaviour, not run here. The line came in with e5833a3 (2026-09-13) and most likely means "the only Vite/Start config, no `app.config.ts`". Read literally, it forbids the tsconfig that the other two lines need. Should say: `vite.config.ts  # the only build config: tanstackStart(), react(), tailwindcss()` and add `tsconfig.json  # strict baseline per the typescript skill` to the tree.
- **[low / confidence: high]** divergent — `skills/typescript/SKILL.md:24` vs `skills/frontend/SKILL.md:95` — the typescript skill says "**No `any`.**" with no exception, and so does the checklist (`typescript-review.md:7`). The frontend skill says "No `any` where a real type is feasible". The typescript description (`:8`) says not to load it for "UI, styling, components", so for component code the looser frontend wording is often the only rule in context. Should say (frontend:95): "Type every prop and API payload. No `any`: use `unknown` and narrow (typescript skill)."
- **[low / confidence: high]** divergent — `skills/typescript/SKILL.md:46` vs `skills/review-standards/references/typescript-review.md:9` — the skill accepts `@ts-expect-error` "(with reason)" and ranks it above `@ts-ignore`. The checklist lists `// @ts-ignore` among "disabling flags" that are findings outright, and flags `@ts-expect-error` "without a linked issue". So a reviewer flags a reasoned `@ts-expect-error` that follows the skill. The checklist names the skill as the language authority (`typescript-review.md:3`, `review-standards/SKILL.md:35-36`), so the checklist should change. Should say: "`// @ts-ignore` anywhere, or `// @ts-expect-error` without a reason comment".
- **[low / confidence: high]** divergent — `skills/typescript/SKILL.md:40` vs `skills/review-standards/references/typescript-review.md:29` — the skill bans default exports only "for shared modules". The checklist flags `export default` "for anything but framework-mandated files (routes)". The checklist's example is also wrong for this stack: TanStack routes use the named `export const Route` (`skills/frontend/SKILL.md:71`), and the framework-mandated default export is `vite.config.ts` (`skills/frontend/SKILL.md:53`). Should say (checklist:29): "flag `export default` in shared modules; framework-mandated files (`vite.config.ts`) are exempt".
- **[low / confidence: medium]** divergent — `skills/typescript/SKILL.md:45` vs `skills/review-standards/references/typescript-review.md:8` — the skill requires a comment on every `as` cast and `!` assertion ("each one needs a comment explaining why the compiler is wrong"). The checklist accepts "a runtime check **or** an adjacent comment". So a cast with a runtime check and no comment passes review but breaks the skill. Pick one wording. The skill is the authority, so the checklist should say "an adjacent comment proving the invariant (a runtime check does not replace it)". The other option is to widen the skill.
- **[low / confidence: medium]** stale verification — `.agents/docs/testing/skill-matrix.md:41` — the `typescript` row's pass check is `grep -qi "typescript" "$OUT"`. The prompt asks for "a strictly typed TypeScript function", so any answer with a ```` ```typescript ```` fence or the word TypeScript passes, whether or not the skill loaded. The check cannot fail. Procedure step 3 (`skill-matrix.md` §Procedure 3) reads the transcript for `Skill` tool uses, and that transcript read is the only real signal. Should say: a check that also requires a rule found only in this skill, e.g. `grep -qiE "unknown|noUncheckedIndexedAccess|satisfies" "$OUT"`, plus the transcript `Skill` use recorded as the pass condition.

## Classification (--project)

| File | Class | Proof |
|---|---|---|
| `skills/typescript/SKILL.md` | accurate | `:8-9`, `:14` name `frontend`, `golang`, `bash` and all three exist (`skills/frontend/SKILL.md:1-2`, `skills/golang/`, `skills/bash/SKILL.md:1-3`). The routing matches `context/dev-setup.md:6`, `skills/review-standards/SKILL.md:35` and `agents/reviewer.md:41`. Its rules agree with each other and with `review-standards/SKILL.md:28`, which requires caught errors typed `unknown`, as `:41` does. The divergences above belong to the neighbouring files, which name this skill as the authority. Not redundant: it is the one language-level source, and the checklist defers to it (`typescript-review.md:3`). |

## Unverified

- That `tsc -b` fails without a `tsconfig.json`: this comes from TypeScript's documented build-mode behaviour. No `tsc` was run here.
- Whether `README.md`, `<docs-root>/specs/` (empty), `testing/sh-agent/guide.md` and `deploy/*.md` should mention this skill at all: none of them name it, and none claims to list skills one by one, so the absence is not a finding.
- The newest `changes/` entry (`2026-09-15-research-beads-elephant-goldfish.md`) does not touch TypeScript. There is nothing to compare.

## Criteria

1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined

- `skills/frontend/references/*` and `skills/review-standards/references/go-review.md` were only grepped (for `tsconfig`, `any`, `export default`, `ts-ignore`), not read whole.
- Vendored skills (`skills/VENDORED.md`) were excluded from every search, as the brief requires.
- `.agents/rules`, `.agents/commands` and `.agents/agents` do not exist in this repo, so they were not searched.

## Re-check commands

```bash
git grep -n 'ONLY config file\|tsc -b' -- skills/frontend/SKILL.md
git grep -n -i 'tsconfig' -- skills/frontend skills/typescript
git grep -n -E 'No `any`' -- skills/frontend/SKILL.md skills/typescript/SKILL.md skills/review-standards/references/typescript-review.md
git grep -n -i 'ts-ignore\|ts-expect-error' -- skills/typescript skills/review-standards
git grep -n -i 'default export\|export default' -- skills/typescript/SKILL.md skills/review-standards/references/typescript-review.md skills/frontend/SKILL.md
git grep -n 'runtime check or an adjacent comment\|explaining why the compiler is wrong' -- skills/review-standards skills/typescript
git grep -n '^| `typescript`' -- .agents/docs/testing/skill-matrix.md
```
