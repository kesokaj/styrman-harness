---
owner: research
date: 2026-09-30
---

# Research brief: fixing the `frontend` references and the `typescript` skill

Scope: the `frontend` skill's four references and `skills/typescript/SKILL.md`. Every version was read on 2026-09-30 from the npm registry (`registry.npmjs.org/<pkg>`, dist-tags + publish time) unless another source is named. Contrast ratios were computed with the WCAG luminance formula, not estimated.

## 1. Current versions and canonical setup

| Package | `latest` | Published | Source |
|---|---|---|---|
| `@tanstack/react-start` | 1.168.60 | 2026-09-30 | registry; [TanStack releases](https://github.com/TanStack/router/releases) "Release 2026-09-30 17:47" |
| `@tanstack/react-router` | 1.170.41 | 2026-09-30 | same |
| `vite` | 8.3.1 | 2026-09-24 | registry; [Vite 8 announcement](https://vite.dev/blog/announcing-vite8) 2026-03-12: Rolldown replaces Rollup+esbuild, Node 20.19+/22.12+ |
| `@vitejs/plugin-react` | 6.1.1 | 2026-08-28 | registry |
| `tailwindcss`, `@tailwindcss/vite` | 4.3.3 | 2026-07-16 | registry; [endoflife.date/tailwind-css](https://endoflife.date/tailwind-css): 4.3 supported, 3.4 LTS to Feb 2027 |
| `react`, `react-dom` | 19.3.0 | 2026-09-09 | registry; [react.dev/versions](https://react.dev/versions) |
| `typescript` | 7.0.2 | 2026-07-08 | registry; [TS 7.0 announcement](https://devblogs.microsoft.com/typescript/announcing-typescript-7-0/): the Go compiler ships as `typescript@7`; `@typescript/typescript6` keeps the JS API |
| `zod` | 4.6.5 | 2026-09-13 | registry |
| `@tanstack/react-form` | 1.33.5 | 2026-08-11 | registry |
| `vitest` | 5.0.3 | 2026-09-30 | registry; [Vitest 5.0 post](https://vitest.dev/blog/vitest-5.html) 2026-09-03: needs Node >= 22.12, Vite >= 6.4 |
| `@testing-library/react` | 16.3.3 | 2026-08-27 | registry |
| `@playwright/test` | 1.63.0 | 2026-09-04 | registry |
| `motion` | 13.4.6 | 2026-09-29 | registry |

SKILL.md's pins `vite 8.2` and `React 19.2` lag one minor each; the other pins hold.

**TanStack Start canonical shape** ([build-from-scratch](https://tanstack.com/start/latest/docs/framework/react/build-from-scratch), [SPA mode](https://tanstack.com/start/latest/docs/framework/react/guide/spa-mode)):
- Plugin import `@tanstack/react-start/plugin/vite`; `viteReact()` after `tanstackStart()` (SKILL.md already has this).
- SPA key `tanstackStart({ spa: { enabled: true } })`, with `spa.maskPath` (default `/`) and `spa.prerender: { outputPath, crawlLinks, retryCount }`. The shell is "prerendering your application's root route only", written to `/_shell.html` by default; "a CDN that can serve static assets is all you need". The [hosting guide](https://tanstack.com/start/latest/docs/framework/react/guide/hosting) names `dist/client` as the publish directory. The Go catch-all must serve that shell: set `spa.prerender.outputPath: '/index.html'` or serve `_shell.html` by name (verify with one build).
- `<head>` is owned by `src/routes/__root.tsx`: `createRootRoute({ head: () => ({ meta, links, scripts }) })`, rendered by `<HeadContent />` in `<head>` and `<Scripts />` at the end of `<body>` ([head management](https://tanstack.com/router/latest/docs/framework/react/guide/document-head-management)). No `index.html`, no `main.tsx`.
- CSS: [Tailwind's Start guide](https://tailwindcss.com/docs/installation/framework-guides/tanstack-start) and [Start's CSS guide](https://tanstack.com/start/latest/docs/framework/react/guide/css-styling) show `import appCss from '../styles.css?url'` + `links: [{ rel: 'stylesheet', href: appCss }]`; a side-effect `import '../styles.css'` is the alternative "for app-wide CSS resets, design tokens" and keeps Start's early hints and inlining. Pick one per project.

**Tailwind 4** ([upgrade guide](https://tailwindcss.com/docs/upgrade-guide), [theme](https://tailwindcss.com/docs/theme), [dark mode](https://tailwindcss.com/docs/dark-mode)):
- `@import "tailwindcss"` replaces the three `@tailwind` directives. "JavaScript config files are still supported for backward compatibility, but they are no longer detected automatically in v4" – a new project has **no** `tailwind.config.*`; `content:` is gone (automatic source detection; `@import 'tailwindcss' source('../')` if the root must be pinned, per [Start's Tailwind guide](https://tanstack.com/start/latest/docs/framework/react/guide/tailwind-integration)).
- Tokens live in `@theme { --color-mint-500: … }` and generate `bg-/text-/border-/ring-/outline-mint-500`. Values that reference other variables need `@theme inline { --color-x: var(--y) }` ("without `inline`, your utility classes might resolve to unexpected values"). `--color-*: initial` clears the default palette.
- Dark variant: `@custom-variant dark (&:where([data-theme=dark], [data-theme=dark] *));` (docs give `.dark` and `data-theme` forms). Unneeded when every colour is a token flipping at `:root`; harmless for the odd `dark:` one-off.
- Defaults changed: `ring` is now 1px (was 3px), `border-*` defaults to `currentColor` (was gray-200). Browser floor: Safari 16.4, Chrome 111, Firefox 128.

## 2. Theme: light/dark, system preference, toggle, no flash

Minimal correct shape for Start + Tailwind 4 (tokens in CSS, utilities generated from them, script inline in the root head):

```css
/* src/styles.css */
@import "tailwindcss";
@custom-variant dark (&:where([data-theme=dark], [data-theme=dark] *));

:root { color-scheme: light; --bg: #ffffff; --fg: #3b3b3b; --accent: #005fb8; --ring: #005fb8; /* … */ }
:root[data-theme="dark"] { color-scheme: dark; --bg: #16191d; --fg: #abb2bf; --accent: #61afef; --ring: #61afef; }

@theme inline {
  --color-*: initial;                 /* only tokens; no zinc-*, amber-* */
  --color-bg: var(--bg); --color-fg: var(--fg);
  --color-accent: var(--accent); --color-ring: var(--ring);
}
```

```tsx
// src/routes/__root.tsx
import appCss from '../styles.css?url'
const themeScript = `(function(){try{var t=localStorage.getItem('theme');
if(t!=='light'&&t!=='dark'){t=matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light'}
document.documentElement.dataset.theme=t}catch(e){}})()`

export const Route = createRootRoute({
  head: () => ({
    meta: [{ charSet: 'utf-8' }, { name: 'viewport', content: 'width=device-width, initial-scale=1' },
           { name: 'color-scheme', content: 'light dark' }],
    links: [{ rel: 'stylesheet', href: appCss }],
    scripts: [{ children: themeScript }],       // rendered by <HeadContent />, before the stylesheet paints
  }),
  component: () => (<html><head><HeadContent /></head><body><Outlet /><Scripts /></body></html>),
})
```

Sources: `scripts[].children` is the documented inline slot and the head guide's own example is this theme script ([document head management](https://tanstack.com/router/latest/docs/framework/react/guide/document-head-management)); `color-scheme` makes form controls, scrollbars and the canvas follow the theme, and the `<meta name="color-scheme">` "before any CSS" stops the white flash on a dark OS ([MDN](https://developer.mozilla.org/en-US/docs/Web/CSS/color-scheme)); the localStorage/`matchMedia` cascade is Tailwind's recipe ([dark mode](https://tailwindcss.com/docs/dark-mode)), `removeItem('theme')` meaning "follow the OS". The toggle gains that third state and swaps its inline `style` for token classes.

## 3. Forms

**Default: no library.** `FormData` is Baseline since 2015 and a valid `fetch` body ([MDN](https://developer.mozilla.org/en-US/docs/Web/API/FormData)); React 19's `<form action={fn}>` calls `fn(formData)` in a transition, resets fields on success, needs no `preventDefault`, and pairs with `useActionState`/`useFormStatus` ([react.dev](https://react.dev/reference/react-dom/components/form)). Validation is the Go handler's job; the client shows what the API returns:

```tsx
const [state, submit, pending] = useActionState(async (_prev: FormState, fd: FormData) => {
  const res = await api.signup(Object.fromEntries(fd))   // api.ts: fetch('/api/signup', { method: 'POST', body: JSON.stringify(…) })
  return res.ok ? { ok: true } : { ok: false, errors: await res.json() }   // Go returns { field: message }
}, { ok: true })
<form action={submit} noValidate>
  <label htmlFor="email">Email</label>
  <input id="email" name="email" type="email" required aria-invalid={!!state.errors?.email} aria-describedby="email-err" />
  {state.errors?.email && <p id="email-err" role="alert">{state.errors.email}</p>}
  <button disabled={pending}>Sign up</button>
</form>
```

**When a library is justified** (multi-step, arrays of fields, cross-field rules, optimistic per-field errors): `@tanstack/react-form` 1.33.5 — `useForm({ defaultValues, onSubmit: async ({ value }) => … })`, `<form.Field name="…">{(field) => <input value={field.state.value} onChange={(e) => field.handleChange(e.target.value)} />}</form.Field>`, `validators: { onChange: zodSchema }` accepting any Standard Schema library ([quick start](https://tanstack.com/form/latest/docs/framework/react/quick-start)). Same vendor as the router, headless, no resolver package — so `react-hook-form` 7.89.0 + `@hookform/resolvers` 5.9.1 are two dependencies the stack does not need.

**Zod 4**: "The method forms (`z.string().email()`) still exist and work as before, but are now deprecated" — use top-level `z.email()` ([zod.dev/v4/changelog](https://zod.dev/v4/changelog)); import from `"zod"`. Zod itself is only justified when a library needs a schema or the client validates something the server cannot.

## 4. tsconfig, unit tests, e2e

**TypeScript 7.0.2 is `latest`** ([announcement](https://devblogs.microsoft.com/typescript/announcing-typescript-7-0/), 2026-07-08); its defaults come from 6.0 ([6.0 announcement](https://devblogs.microsoft.com/typescript/announcing-typescript-6-0/), 2026-03-23): `strict: true`, `module: esnext`, floating `target` (es2025), `types: []`, `noUncheckedSideEffectImports: true`; removed: `moduleResolution: node/classic`, `baseUrl`, `esModuleInterop: false`, `target: es5`. Vite "only performs transpilation … and does NOT perform type checking" with Oxc, so `isolatedModules` (or `verbatimModuleSyntax`, which the [tsconfig reference](https://www.typescriptlang.org/tsconfig/) says supersedes `isolatedModules`, `preserveValueImports`, `importsNotUsedAsValues`) is required, and Vite ignores tsconfig `target` ([Vite features](https://vite.dev/guide/features.html)). The [create-vite react-ts template](https://raw.githubusercontent.com/vitejs/vite/main/packages/create-vite/template-react-ts/tsconfig.app.json) is the reference shape; the recommended `web/tsconfig.json`:

```jsonc
{ "compilerOptions": {
    "target": "es2023", "lib": ["ES2023", "DOM", "DOM.Iterable"], "module": "esnext",
    "moduleResolution": "bundler", "jsx": "react-jsx", "types": ["vite/client"],
    "verbatimModuleSyntax": true, "moduleDetection": "force", "erasableSyntaxOnly": true,
    "allowImportingTsExtensions": true, "noEmit": true, "skipLibCheck": true,
    "strict": true, "noUncheckedIndexedAccess": true, "noImplicitOverride": true,
    "noUnusedLocals": true, "noUnusedParameters": true, "noFallthroughCasesInSwitch": true },
  "include": ["src"] }
```

`erasableSyntaxOnly` forbids `enum`, runtime `namespace`, parameter properties, `import =`/`export =`, `<T>x` assertions ([docs](https://www.typescriptlang.org/tsconfig/erasableSyntaxOnly.html)), matching the skill's "unions, not enums" rule. The template's build is `tsc -b && vite build` (types first); SKILL.md has it reversed.

**Unit tests: Vitest 5.0.3** (needs Node >= 22.12, Vite >= 6.4; `clearMocks` on by default; browser providers are separate packages such as `@vitest/browser-playwright` — [migration guide](https://vitest.dev/guide/migration.html)). Two tiers: (a) `environment: 'jsdom'` (jsdom 30.1.1) + `@testing-library/react` 16.3.3 + `@testing-library/jest-dom` 7.0.1 for component logic; (b) [browser mode](https://vitest.dev/guide/browser/) with `provider: playwright()` and `vitest-browser-react` 2.3.0 when real layout or `:focus-visible` matters (Vitest "explicitly discourages" `@testing-library/user-event` there and ships its own `userEvent`). **E2E: `@playwright/test` 1.63.0** against the Go binary serving `dist/client`, with `@axe-core/playwright` 4.13.0 for the accessibility gate.

## 5. Focus ring contrast

WCAG 2.2 ([1.4.11](https://www.w3.org/WAI/WCAG22/Understanding/non-text-contrast.html), AA): UI components and their states need 3:1 "against adjacent color(s)"; an external indicator is measured against the background the component sits on, an internal one against the component. [2.4.7](https://www.w3.org/WAI/WCAG22/Understanding/focus-visible.html) (AA) lists C45 `:focus-visible` as a sufficient technique. [2.4.13 Focus Appearance](https://www.w3.org/WAI/WCAG22/Understanding/focus-appearance.html) (AAA): area >= a 2 CSS px perimeter and 3:1 between focused and unfocused pixels — a 2px `outline` with 2px `outline-offset` on a token colour satisfies both without measuring per component. `:focus-visible` is Baseline since March 2022 ([MDN](https://developer.mozilla.org/en-US/docs/Web/CSS/:focus-visible)).

Measured (WCAG formula):

| Pair | Ratio | Verdict |
|---|---|---|
| dark `--focus-ring` #3e4452 on #16191d / #1e2227 (current theme-system.md) | 1.81 / 1.64 | **fails** 1.4.11 |
| dark `--accent` #61afef on #16191d / #1e2227 / #23272e | 7.46 / 6.77 / 6.34 | passes |
| light `--focus-ring` #005fb8 on #fff / #f8f8f8 | 6.31 / 5.94 | passes |
| accessibility.md `amber-400` #fbbf24 on white | 1.67:1 | **fails**; passes only on zinc-950 (11.9:1) |
| light `--error` / `--success` / `--warning` as **text** on white | 3.35 / 3.37 / 3.38 | icon/border only; fail 4.5:1 as text |
| `--fg-tertiary` as text, light on white / dark on #16191d | 3.08 / 2.92 | disabled only (exempt); not placeholders |

Token pattern that passes on both themes (one rule, no per-theme override):

```css
:root { --ring: var(--accent); }            /* light #005fb8, dark #61afef — both >= 5.9:1 on every surface */
:focus-visible { outline: 2px solid var(--ring); outline-offset: 2px; }
:focus:not(:focus-visible) { outline: none; }
```
Tailwind form: `focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring` (with `--color-ring` in `@theme inline`). Never `focus:ring-*` with a `/20` alpha as the only indicator (component-patterns.md's input does this).

## 6. Node/npm hygiene for `web/` inside a Go repo

- `npm ci` "must have an existing package-lock.json", errors on a package.json/lock mismatch, removes `node_modules` first and "will never write to package.json or any of the package-locks" ([npm ci](https://docs.npmjs.com/cli/v11/commands/npm-ci)): commit `web/package-lock.json`, `npm ci` in every scripted build, `npm install` only to change dependencies.
- `engines` "is advisory only" unless `engine-strict` is set; `devEngines` validates the contributor's environment before `install`/`ci`/`run` ([package.json](https://docs.npmjs.com/cli/v11/configuring-npm/package-json)). Set `"engines": { "node": ">=22.12" }` (Vite 8 and Vitest 5 floors) and an `.nvmrc`/`.node-version` with the pinned major.
- Node on 2026-09-30 ([nodejs.org](https://nodejs.org/en/about/previous-releases)): v26 Current, **v24 Krypton LTS** (recommended), v22 maintenance LTS, v20 and below EOL. Docker: a `node:24-alpine` stage (exact tag, never `latest`) runs `npm ci && npm run build`; the Go stage copies `web/dist/client` before `go build`. [docker-node best practices](https://github.com/nodejs/docker-node/blob/main/docs/BestPractices.md) (`NODE_ENV=production`, the non-root `node` user) apply only to a stage that *runs* Node — the SPA image never does.
- "Plain V8" concretely: nothing under `web/src/` imports `node:*`, `fs`, `path`, `process`, `Buffer` or `__dirname`; `types: ["vite/client"]` (not `node`) in the app tsconfig makes that a compile error, `@types/node` living only in `tsconfig.node.json` for `vite.config.ts`. Only `VITE_*` keys of `import.meta.env` reach the bundle.

## Proposed edits

**theme-system.md**
- Replace the "Tailwind CSS Integration" section (`tailwind.config.ts`, `content:`) with the `@theme inline` block from §2 and one line: "no `tailwind.config.*`, no `darkMode`, no `dark:` — the tokens flip at `:root`" (keep the `@custom-variant dark` line for one-offs).
- Replace the "For TanStack Start, add the FOWT script to `index.html`" block with the `__root.tsx` `head()` shape from §2 (`scripts: [{ children }]`, `meta color-scheme`, `links` with `?url`); drop `<div id="root">`.
- Rewrite `ThemeToggle` to Tailwind classes on tokens, three states (light / dark / system = `localStorage.removeItem`), and a `matchMedia` change listener while in system mode; keep the SVGs.
- Change dark `--focus-ring` from `#3e4452` to `#61afef` (or define `--focus-ring: var(--accent)` once); add a "measured 2026-09-30" note with the two ratios so the next edit re-measures.
- Note on `--fg-tertiary`: "disabled only; placeholders use `--fg-secondary`".
- Usage checklist: step 2 becomes "`@import` tokens after `@import "tailwindcss"` in `styles.css`", step 3 "add the script to `__root.tsx` `head()`", and add "run the 5 tests as one Playwright spec".

**accessibility.md**
- Replace the zinc/amber "Dark Mode Safe Pairings" table with the *theme's* token pairs and measured ratios (§5), plus the rule "measure a pair before adding it".
- Replace both focus-ring snippets (`var(--color-accent, #fbbf24)`, `focus-visible:outline-amber-400`) with the token pattern in §5 and add SC 2.4.13's 2 px perimeter as the reason for `outline-2 outline-offset-2`.
- Skip-link and TabBar: `bg-amber-400 text-zinc-950 bg-zinc-800 text-zinc-400` -> `bg-accent text-fg-inverse bg-bg-active text-fg-secondary`.
- Quick-reference table: add 2.4.13 (AAA, the measurable target) and note 1.4.11 covers the ring.
- Testing: replace the `main.tsx` + `@axe-core/react` block with an `@axe-core/playwright` spec against the Go-served build.

**component-patterns.md**
- Forms: replace the RHF+resolvers+`z.string().email()` example with the native `<form action>` + `useActionState` shape (§3); add a three-line "when a library" paragraph pointing at `@tanstack/react-form` with its `useForm`/`form.Field` shape and `z.email()` if Zod is used. Delete "pair with `createServerFn`" (SPA mode has no server functions; SKILL.md already says so).
- Every example: replace hard-coded `zinc-*`/`amber-*`/`emerald-600`/`red-600`/`rgba(9,9,11,…)` with tokens (`bg-bg-secondary`, `text-fg`, `bg-accent`, `bg-success`, `bg-error`); the header's `useMotionTemplate` background becomes `color-mix(in oklab, var(--bg) X%, transparent)` or a token with alpha.
- Input focus style: `focus:ring-amber-400/20` -> `focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring`.
- Modal: prefer native `<dialog>` + `showModal()` (focus trap, Escape, top layer for free); keep the portal version as the animated fallback.

**animation-recipes.md**
- Same token substitution across all eight recipes (`amber-400`, `zinc-*`, `bg-white`, `bg-black/70` -> `bg-accent`, `bg-bg`, `text-fg`, a `--overlay` token).
- Recipe 8 `SafeMotion`: `animate={reduced ? animate : animate}` is a no-op branch; delete it. Pin `motion` 13.4 (2026-09-29) once at the top.

**typescript/SKILL.md**
- Compiler Baseline: state the 2026 floor (TS 7.0 native compiler; 6.0 defaults: `strict`, `module esnext`, `types []`), then the app-config essentials: `moduleResolution bundler`, `verbatimModuleSyntax`, `moduleDetection force`, `erasableSyntaxOnly`, `noUncheckedIndexedAccess`, `noImplicitOverride`, `types: ["vite/client"]`, `noEmit` (Vite transpiles, `tsc -b` checks; build is `tsc -b && vite build`).
- Add a "Tests" section: Vitest 5 (jsdom + RTL 16 by default, browser mode when layout matters), Playwright 1.63 for e2e, one line each with the version and date.
- Add to Modules: "`erasableSyntaxOnly` means no `enum`, no runtime `namespace`, no parameter properties — `as const` objects and unions instead" (replaces the softer enum line).
- Boundaries: "Zod 4 (`z.email()`, top-level formats) only where a schema is reused; otherwise a hand-written type guard".

**`skills/typescript/references/` — warranted.** SKILL.md is 339 words; two copy-paste artifacts would double it:
1. `references/tsconfig.md` — the app tsconfig above, a `tsconfig.node.json` for `vite.config.ts` (`types: ["node"]`), the project-references root, and a one-line rationale per flag with the TS-docs link.
2. `references/testing.md` — `vitest.config.ts` for both tiers (jsdom+RTL, browser+Playwright), a `playwright.config.ts` that starts the Go binary as `webServer`, the axe spec, and the npm scripts (`test`, `test:browser`, `e2e`) so `verifier` has named commands.

**SKILL.md (frontend)**, minor: bump the pin line to `vite 8.3`, `React 19.3`, add `typescript 7.0`, `vitest 5.0`, `@playwright/test 1.63` (registry, 2026-09-30); fix the build script order to `tsc -b && vite build`; add one line on the SPA shell filename the Go catch-all serves (`spa.prerender.outputPath`).

## What was not verified
- No build was run: the `_shell.html`/`index.html` contract for the Go embed and `@theme inline` + `--color-*: initial` together are documented, not executed.
- `motion` 13's API was not re-read against the recipes.
- Vitest 5's environment changes come from the migration summary only.
