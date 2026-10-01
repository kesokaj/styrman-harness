# Theme System

Reference file for the frontend skill. Every frontend ships with dark/light theme support. **"System" (follow the OS) is the default on first load**, which resolves to light when the OS states no preference; the toggle's three states are light, dark, system.

The shape, in one sentence: tokens are CSS custom properties on `:root` (light) and `:root[data-theme="dark"]` (dark), Tailwind 4 maps them to utilities with `@theme inline`, and a one-line script in `__root.tsx`'s `head()` sets `data-theme` before the stylesheet paints. There is no `tailwind.config.*`, no `darkMode` option, no `index.html`, no `main.tsx`.

## Design Tokens

All UI colors are defined as CSS custom properties. Components reference tokens — never raw hex values, never a Tailwind palette color (`--color-*: initial` removes the default palette so `bg-zinc-900` does not exist).

### 🌞 Light Mode — VS Code Light Modern

Source: [microsoft/vscode `light_modern.json`](https://github.com/microsoft/vscode/blob/main/extensions/theme-defaults/themes/light_modern.json)

| Token | Hex | Usage |
|-------|-----|-------|
| `--bg` | `#ffffff` | Page background |
| `--bg-secondary` | `#f8f8f8` | Cards, sidebar, panels |
| `--bg-tertiary` | `#f3f3f3` | Code blocks, wells |
| `--bg-hover` | `#f2f2f2` | List/row hover |
| `--bg-active` | `#e8e8e8` | Active selection |
| `--bg-input` | `#ffffff` | Input fields |
| `--fg` | `#3b3b3b` | Body text |
| `--fg-secondary` | `#616161` | Muted text, placeholders |
| `--fg-tertiary` | `#8b949e` | **Disabled text only** (3.1:1 on white — never placeholders or body text) |
| `--fg-heading` | `#1f1f1f` | Headings |
| `--fg-inverse` | `#ffffff` | Text on accent background |
| `--accent` | `#005fb8` | Buttons, links, focus ring |
| `--accent-hover` | `#0258a8` | Button hover |
| `--border` | `#e5e5e5` | Dividers, card borders |
| `--border-input` | `#cecece` | Input borders (decorative; the label and the fill identify the input, the ring carries focus) |
| `--success` | `#2ea043` | **Icons, borders, fills only** — 3.4:1 on white fails 4.5:1 as text; write the word in `--fg` |
| `--error` | `#f85149` | Icons, borders, fills only (3.4:1 as text); the message text is `--fg` with `role="alert"` |
| `--warning` | `#bb8009` | Icons, borders, fills only (3.4:1 as text) |
| `--ring` | `var(--accent)` | Focus outline — 6.3:1 on `--bg`, 5.9:1 on `--bg-secondary`, 5.7:1 on `--bg-tertiary` |

### 🌙 Dark Mode — One Dark Pro Night Flat

Source: [Binaryify/OneDark-Pro `OneDark-Pro-night-flat.json`](https://github.com/Binaryify/OneDark-Pro/blob/master/themes/OneDark-Pro-night-flat.json)

| Token | Hex | Usage |
|-------|-----|-------|
| `--bg` | `#16191d` | Page background |
| `--bg-secondary` | `#1e2227` | Cards, dropdowns |
| `--bg-tertiary` | `#23272e` | Elevated surfaces |
| `--bg-hover` | `#2c313a` | List/row hover |
| `--bg-active` | `#2c313a` | Active selection |
| `--bg-input` | `#1d1f23` | Input fields |
| `--fg` | `#abb2bf` | Body text |
| `--fg-secondary` | `#9da5b4` | Muted text, placeholders |
| `--fg-tertiary` | `#5c6370` | **Disabled text only** (2.9:1 on `--bg`) |
| `--fg-heading` | `#d7dae0` | Headings |
| `--fg-inverse` | `#16191d` | Text on accent background |
| `--accent` | `#61afef` | Buttons, links, focus ring |
| `--accent-hover` | `#528bff` | Button hover |
| `--border` | `#37393d` | Dividers, card borders |
| `--border-input` | `#3e4452` | Input/panel borders |
| `--success` | `#98c379` | Success states — passes as text in dark mode (8.8:1 on `--bg`) |
| `--error` | `#e06c75` | Error states — passes as text in dark mode (5.5:1 on `--bg`, 4.7:1 on `--bg-tertiary`) |
| `--warning` | `#d19a66` | Warning states — passes as text in dark mode (7.2:1 on `--bg`) |
| `--ring` | `var(--accent)` | Focus outline — 7.5:1 on `--bg`, 6.8:1 on `--bg-secondary`, 6.3:1 on `--bg-tertiary` |

Because the semantic colors pass as text only in dark mode, the rule that holds on both themes is the light one: semantic colors on icons, borders and fills, the words in `--fg`.

### Accent Palette (dark mode decorative colors)

Tags, badges, syntax highlights, decorative elements; dark mode only, never as the only carrier of meaning:

| Name | Hex | Role |
|------|-----|------|
| Coral | `#e06c75` | Destructive, alerts |
| Sand | `#d19a66` | Warnings, constants |
| Gold | `#e5c07b` | Types, highlights |
| Sage | `#98c379` | Success, confirmed |
| Sky | `#61afef` | Links, info |
| Lilac | `#c678dd` | Keywords, premium |
| Teal | `#56b6c2` | Info badges, metadata |

### Contrast, measured

Every ratio above was computed with the WCAG 2.x relative-luminance formula on 2026-10-01 (script: a dozen lines of Node; re-run it before changing a value). What passes and what does not:

| Pair | Light | Dark | Verdict |
|---|---|---|---|
| `--ring` (= `--accent`) on `--bg` / `--bg-secondary` / `--bg-tertiary` | 6.31 / 5.94 / 5.68 | 7.46 / 6.77 / 6.34 | passes 1.4.11 (3:1) and 2.4.13 everywhere |
| `--fg` on `--bg` | 11.2 | 8.3 | passes 1.4.3 |
| `--fg-secondary` on `--bg` / `--bg-tertiary` | 6.2 / 5.6 | 7.1 / 6.1 | passes 1.4.3 as text |
| `--fg-tertiary` on `--bg` | 3.08 | 2.92 | fails 4.5:1; disabled text only (exempt) |
| `--fg-inverse` on `--accent` / `--accent-hover` | 6.31 / 7.09 | 7.46 / 5.45 | passes as button text |
| `--success` / `--error` / `--warning` on `--bg` | 3.37 / 3.35 / 3.38 | 8.75 / 5.52 / 7.16 | light: icons and borders only; dark: text too |
| `--border-input` on `--bg` | 1.57 | 1.81 | below 3:1: an input is identified by its label and `--bg-input` fill, and its focused state by `--ring` |
| the previous dark `--focus-ring` `#3e4452` on `--bg` | — | 1.81 | **failed** 1.4.11; replaced by `--ring: var(--accent)` |

Rule: measure a pair before adding it to a table; "looks fine" on one monitor is not a measurement.

---

## Copy-Paste: `src/styles.css`

One file: Tailwind, the dark variant, the tokens, the `@theme inline` mapping, the focus ring, the body. Linked from `__root.tsx` (below).

```css
/* src/styles.css — Tailwind 4 + theme tokens
   Light: VS Code Light Modern (default)   Dark: One Dark Pro Night Flat */
@import "tailwindcss";

/* dark: variant for the odd one-off; the tokens flip at :root, so it is rarely needed */
@custom-variant dark (&:where([data-theme=dark], [data-theme=dark] *));

:root {
  color-scheme: light;                 /* form controls, scrollbars and the canvas follow */
  --bg: #ffffff;        --bg-secondary: #f8f8f8;  --bg-tertiary: #f3f3f3;
  --bg-hover: #f2f2f2;  --bg-active: #e8e8e8;     --bg-input: #ffffff;
  --fg: #3b3b3b;        --fg-secondary: #616161;  --fg-tertiary: #8b949e;
  --fg-heading: #1f1f1f; --fg-inverse: #ffffff;
  --accent: #005fb8;    --accent-hover: #0258a8;
  --border: #e5e5e5;    --border-input: #cecece;
  --success: #2ea043;   --error: #f85149;         --warning: #bb8009;
  --ring: var(--accent);               /* one rule, both themes: 5.7–7.5:1 on every surface */
  --overlay: rgb(0 0 0 / 0.5);
  --shadow-sm: 0 1px 2px rgb(0 0 0 / 0.05);
  --shadow-md: 0 4px 6px rgb(0 0 0 / 0.07);
  --shadow-lg: 0 10px 15px rgb(0 0 0 / 0.1);
}

:root[data-theme="dark"] {
  color-scheme: dark;
  --bg: #16191d;        --bg-secondary: #1e2227;  --bg-tertiary: #23272e;
  --bg-hover: #2c313a;  --bg-active: #2c313a;     --bg-input: #1d1f23;
  --fg: #abb2bf;        --fg-secondary: #9da5b4;  --fg-tertiary: #5c6370;
  --fg-heading: #d7dae0; --fg-inverse: #16191d;
  --accent: #61afef;    --accent-hover: #528bff;
  --border: #37393d;    --border-input: #3e4452;
  --success: #98c379;   --error: #e06c75;         --warning: #d19a66;
  --overlay: rgb(0 0 0 / 0.7);
  --shadow-sm: 0 1px 2px rgb(0 0 0 / 0.3);
  --shadow-md: 0 4px 6px rgb(0 0 0 / 0.4);
  --shadow-lg: 0 10px 15px rgb(0 0 0 / 0.5);
  /* decorative accent palette, dark only */
  --accent-coral: #e06c75; --accent-sand: #d19a66; --accent-gold: #e5c07b;
  --accent-sage: #98c379;  --accent-sky: #61afef;  --accent-lilac: #c678dd; --accent-teal: #56b6c2;
}

/* Tokens → utilities. `inline` because the values are var() references;
   `--color-*: initial` removes Tailwind's palette, so bg-zinc-900 is not a class. */
@theme inline {
  --color-*: initial;
  --color-bg: var(--bg);               --color-bg-secondary: var(--bg-secondary);
  --color-bg-tertiary: var(--bg-tertiary); --color-bg-hover: var(--bg-hover);
  --color-bg-active: var(--bg-active); --color-bg-input: var(--bg-input);
  --color-fg: var(--fg);               --color-fg-secondary: var(--fg-secondary);
  --color-fg-tertiary: var(--fg-tertiary); --color-fg-heading: var(--fg-heading);
  --color-fg-inverse: var(--fg-inverse);
  --color-accent: var(--accent);       --color-accent-hover: var(--accent-hover);
  --color-border: var(--border);       --color-border-input: var(--border-input);
  --color-success: var(--success);     --color-error: var(--error);  --color-warning: var(--warning);
  --color-ring: var(--ring);           --color-overlay: var(--overlay);
  --shadow-sm: var(--shadow-sm);       --shadow-md: var(--shadow-md); --shadow-lg: var(--shadow-lg);
}

/* Focus: 2px outline, 2px offset — WCAG 2.4.13's 2 CSS px perimeter, 3:1 by token choice */
:focus-visible { outline: 2px solid var(--ring); outline-offset: 2px; }
:focus:not(:focus-visible) { outline: none; }

body { background-color: var(--bg); color: var(--fg); }
```

Generated utilities: `bg-bg`, `bg-bg-secondary`, `text-fg`, `text-fg-heading`, `text-fg-inverse`, `bg-accent`, `hover:bg-accent-hover`, `border-border`, `border-border-input`, `text-error`, `outline-ring`, `shadow-md`, `bg-overlay`, plus every other color utility Tailwind derives from a `--color-*` token (`ring-accent`, `fill-success`, `decoration-accent`). Theme transitions: add `transition-colors` on the elements that should fade, not on `*`.

---

## FOWT Prevention (Flash of Wrong Theme) — `__root.tsx`

`src/routes/__root.tsx` owns `<head>`. The inline script goes in `head().scripts[].children`, which `<HeadContent />` renders before the stylesheet link, so `data-theme` is set before the first paint. `<meta name="color-scheme" content="light dark">` tells the browser, before any CSS, that both schemes are supported, which stops the white flash of the UA canvas on a dark OS.

```tsx
// src/routes/__root.tsx
import { HeadContent, Outlet, Scripts, createRootRoute } from '@tanstack/react-router'
import appCss from '../styles.css?url'

// No stored choice, or an unknown value, means "follow the OS".
const themeScript = `(function(){try{var t=localStorage.getItem('theme');if(t!=='light'&&t!=='dark'){t=matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light'}document.documentElement.dataset.theme=t}catch(e){}})()`

export const Route = createRootRoute({
  head: () => ({
    meta: [
      { charSet: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { name: 'color-scheme', content: 'light dark' },
      { title: 'App' },
    ],
    links: [{ rel: 'stylesheet', href: appCss }],
    scripts: [{ children: themeScript }],
  }),
  component: () => (
    <html lang="en">
      <head><HeadContent /></head>
      <body><Outlet /><Scripts /></body>
    </html>
  ),
})
```

**Priority cascade:**

| Priority | Source | Behavior |
|----------|--------|----------|
| 1 | `localStorage.theme` = `light` or `dark` | The user's explicit choice persists across sessions |
| 2 | `prefers-color-scheme` | No stored choice (the "system" state): follow the OS, live |
| 3 | Fallback | **Light** (no `matchMedia` match) |

`import appCss from '../styles.css?url'` + `links` is the Start-documented shape; a side-effect `import '../styles.css'` also works and keeps Start's CSS inlining. Pick one per project.

---

## Ready-Made Component: `ThemeToggle.tsx`

Three states — light, dark, system — cycled by one button; "system" removes the stored key and follows `prefers-color-scheme` changes while mounted. Tailwind classes on tokens, no inline `style` objects, inline SVG icons, no dependencies.

```tsx
// src/components/ThemeToggle.tsx
import { useEffect, useState } from 'react'

type Theme = 'light' | 'dark' | 'system'

function readTheme(): Theme {
  if (typeof window === 'undefined') return 'system'
  const stored = localStorage.getItem('theme')
  return stored === 'light' || stored === 'dark' ? stored : 'system'
}

function apply(theme: Theme) {
  const resolved =
    theme === 'system'
      ? matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
      : theme
  document.documentElement.dataset.theme = resolved
}

const next: Record<Theme, Theme> = { light: 'dark', dark: 'system', system: 'light' }
const label: Record<Theme, string> = { light: 'Light', dark: 'Dark', system: 'System' }

export function ThemeToggle() {
  const [theme, setTheme] = useState<Theme>(readTheme)

  useEffect(() => {
    if (theme === 'system') localStorage.removeItem('theme')
    else localStorage.setItem('theme', theme)
    apply(theme)
    if (theme !== 'system') return
    const mq = matchMedia('(prefers-color-scheme: dark)')
    const onChange = () => apply('system')
    mq.addEventListener('change', onChange)
    return () => mq.removeEventListener('change', onChange)
  }, [theme])

  return (
    <button
      type="button"
      onClick={() => setTheme((t) => next[t])}
      aria-label={`Theme: ${label[theme]}. Switch to ${label[next[theme]]}`}
      title={`Theme: ${label[theme]}`}
      className="inline-flex items-center justify-center rounded-lg border border-border p-2 text-fg hover:bg-bg-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
    >
      {theme === 'dark' ? <MoonIcon /> : theme === 'light' ? <SunIcon /> : <SystemIcon />}
    </button>
  )
}

function SunIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      <circle cx="12" cy="12" r="5" />
      <line x1="12" y1="1" x2="12" y2="3" /><line x1="12" y1="21" x2="12" y2="23" />
      <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" /><line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
      <line x1="1" y1="12" x2="3" y2="12" /><line x1="21" y1="12" x2="23" y2="12" />
      <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" /><line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
    </svg>
  )
}

function MoonIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
    </svg>
  )
}

function SystemIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      <rect x="2" y="3" width="20" height="14" rx="2" />
      <line x1="8" y1="21" x2="16" y2="21" /><line x1="12" y1="17" x2="12" y2="21" />
    </svg>
  )
}
```

---

## Tailwind 4 notes

- `@import "tailwindcss"` replaces the three `@tailwind` directives. A new project has **no** `tailwind.config.*`; `content:` is gone (automatic source detection; `@import "tailwindcss" source("../")` only if the scan root must be pinned).
- Tokens that reference other variables must be mapped with `@theme inline`; a plain `@theme` would bake the light value into the utility.
- No `darkMode: 'class'`, and `dark:` only for a one-off that is not a token flip; the `@custom-variant dark` line above makes it `data-theme`-aware for that case.
- Tailwind 4 defaults that changed from 3: `ring` is 1px (was 3px), `border-*` defaults to `currentColor` (was gray-200), so every border names a color (`border-border`).
- Browser floor: Safari 16.4, Chrome 111, Firefox 128.

---

## Usage Checklist

When starting a new frontend:

1. Copy `src/styles.css` above (Tailwind import, `@custom-variant dark`, tokens, `@theme inline`, focus ring, body).
2. Put the `head()` block above in `src/routes/__root.tsx`: the `color-scheme` meta, the stylesheet link, and the theme script.
3. Drop `<ThemeToggle />` into the header.
4. Use token utilities (`bg-bg-secondary`, `text-fg`, `outline-ring`) or `var(--token)` in raw CSS. **Never a raw hex value, never a palette color**: `--color-*: initial` makes `bg-zinc-900` a no-op class, which is the point.
5. Run the quick test below once as a Playwright spec (the `accessibility.md` testing section has the harness).

### Quick Test

- Page loads in light mode on a light OS with nothing stored (system mode resolving to light)
- Toggle: light → dark → system → light; refresh preserves each
- OS set to dark, `localStorage.theme` absent, refresh → dark; flip the OS while in system mode → page follows without a reload
- No flash of wrong theme on any page load (the script runs before the stylesheet paints)
- Tab through the page: every focused control shows the 2px `--ring` outline on both themes
