# Accessibility

Reference file for the frontend skill. Read this when implementing accessible interfaces. Targets WCAG 2.2 Level AA, with 2.4.13 (AAA) as the measurable target for focus, because it is cheaper to meet than to argue about.

## Quick Reference: Key Success Criteria

| SC | Name | Level | Requirement |
|----|------|-------|-------------|
| 1.4.3 | Contrast (Minimum) | AA | 4.5:1 normal text, 3:1 large text (≥18pt or ≥14pt bold) |
| 1.4.11 | Non-text Contrast | AA | 3:1 for UI component boundaries, states and graphics against adjacent colors — this is the rule the focus ring must pass |
| 2.4.7 | Focus Visible | AA | A visible focus indicator on every interactive element; `:focus-visible` (technique C45) is sufficient |
| 2.4.11 | Focus Not Obscured | AA | Sticky/fixed elements must not hide the focused element |
| 2.4.13 | Focus Appearance | AAA | Indicator area ≥ a 2 CSS px perimeter, 3:1 between focused and unfocused pixels — a 2px outline with 2px offset in a token color meets it without per-component measuring |
| 2.5.8 | Target Size | AA | Interactive targets ≥ 24×24 CSS px (aim for 44×44) |
| 1.3.1 | Info and Relationships | A | Semantic HTML; heading hierarchy; landmarks |
| 4.1.2 | Name, Role, Value | A | Interactive elements need an accessible name and role |

---

## Color Contrast

The theme's tokens are the only colors (`theme-system.md`; `--color-*: initial` removes Tailwind's palette). The pairs below were measured with the WCAG relative-luminance formula on 2026-10-01; **measure a pair before adding it**, and re-measure when a token changes.

### Safe pairings, both themes

| Foreground | Background | Light | Dark | Use |
|---|---|---|---|---|
| `text-fg` | `bg-bg` | 11.2:1 | 8.3:1 | Body text |
| `text-fg` | `bg-bg-tertiary` | 10.1:1 | 7.0:1 | Body text on wells |
| `text-fg-heading` | `bg-bg` | 16.5:1 | 12.6:1 | Headings |
| `text-fg-secondary` | `bg-bg` / `bg-bg-tertiary` | 6.2:1 / 5.6:1 | 7.1:1 / 6.1:1 | Muted text, labels, placeholders |
| `text-fg-inverse` | `bg-accent` | 6.3:1 | 7.5:1 | Button text |
| `text-fg-inverse` | `bg-accent-hover` | 7.1:1 | 5.5:1 | Button text, hovered |
| `text-accent` | `bg-bg` | 6.3:1 | 7.5:1 | Links |
| `outline-ring` | any `bg-*` surface | ≥ 5.7:1 | ≥ 6.3:1 | Focus indicator (1.4.11, 2.4.13) |

### Pairs that do NOT pass as text

| Foreground | Light | Dark | Rule |
|---|---|---|---|
| `text-fg-tertiary` on `bg-bg` | 3.1:1 | 2.9:1 | Disabled text only (1.4.3 exempts disabled controls). Placeholders use `text-fg-secondary`. |
| `text-success` / `text-error` / `text-warning` on `bg-bg` | 3.4:1 / 3.4:1 / 3.4:1 | 8.8:1 / 5.5:1 / 7.2:1 | Light fails 4.5:1: semantic colors go on icons, borders and fills (`bg-error`, `border-success`, an SVG `fill-warning`), and the word is `text-fg`. Dark would pass, but one rule for both themes is the point. |
| `border-border-input` on `bg-bg` | 1.6:1 | 1.8:1 | Below 3:1 as a component boundary. An input is identified by its `<label>` and its `bg-bg-input` fill; the focused state is carried by `outline-ring`, never by a border color change. |

### Rules

- Do not rely on color alone to convey state. Pair with icons, labels, or shape changes (a `✓`, the word "Failed", a stripe).
- An error message is `text-fg` with `role="alert"` next to a `border-error` input, not red text.
- When placing text over images, add a gradient overlay (`bg-overlay`) or a solid tonal area. Test against both the lightest and the darkest image regions.
- Large text (≥ 24px, or ≥ 18.66px bold) needs 3:1, so `text-success` and friends are acceptable as **large** headline accents in light mode too; body copy never.

---

## Focus States

### The one rule

Never remove `outline` without a replacement; show the ring for keyboard users with `:focus-visible`. `styles.css` (theme-system.md) sets it globally:

```css
:focus-visible { outline: 2px solid var(--ring); outline-offset: 2px; }
:focus:not(:focus-visible) { outline: none; }
```

Why this shape: 2.4.13 asks for a 2 CSS px perimeter and 3:1 between focused and unfocused states; a 2px outline at 2px offset is that perimeter, and `--ring` (= `--accent`, 5.7–7.5:1 on every token surface) is the contrast. A component that overrides the global rule uses the same utilities:

```html
<button class="focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring">
  Click me
</button>
```

Never `focus:ring-*` with an alpha (`ring-accent/20`) as the only indicator: a 20 % ring is below 3:1 on every surface. Never a border-color change alone (`border-border-input` → `border-accent` is 1.6:1 → 6.3:1 for a 1px line that 2.4.13 does not count as a perimeter).

### Focus Not Obscured

Sticky headers, floating CTAs, and cookie banners can hide the focused element. Ensure:

- `scroll-padding-top` equals the height of any fixed header
- Floating elements have `pointer-events: none` except on their interactive children
- Test by tabbing through the entire page with the keyboard

```css
html {
  scroll-padding-top: 5rem; /* match your fixed header height */
}
```

---

## Keyboard Navigation

### Principles

1. **Logical tab order.** Let DOM order drive focus order. Never use `tabindex > 0`.
2. **No keyboard traps.** Every modal, dropdown, or drawer must be closeable with `Escape`. Native `<dialog>` + `showModal()` gives that, plus the focus trap and the top layer, for free (component-patterns.md).
3. **Skip link.** Add a visually hidden "Skip to main content" link as the first focusable element.
4. **Use native elements.** `<button>`, `<a>`, `<input>`, `<select>`, `<dialog>` are keyboard-accessible by default. A custom `<div onClick>` needs `role`, `tabIndex={0}` and `onKeyDown` — avoid it.

### Skip Link

```tsx
function SkipLink() {
  return (
    <a
      href="#main-content"
      className="sr-only focus:not-sr-only focus:fixed focus:top-4 focus:left-4 focus:z-[999]
        focus:rounded-lg focus:bg-accent focus:px-4 focus:py-2 focus:text-sm focus:text-fg-inverse"
    >
      Skip to main content
    </a>
  )
}

// In the root layout:
<SkipLink />
<Header />
<main id="main-content">{/* ... */}</main>
```

### Roving Tabindex for Tab Groups

When building tab bars or toolbars, only one item is tabbable. Arrow keys move focus.

```tsx
function TabBar({ tabs, active, onChange }: {
  tabs: string[]; active: number; onChange: (i: number) => void
}) {
  const handleKeyDown = (e: React.KeyboardEvent, i: number) => {
    if (e.key === "ArrowRight") onChange((i + 1) % tabs.length)
    if (e.key === "ArrowLeft") onChange((i - 1 + tabs.length) % tabs.length)
  }

  return (
    <div role="tablist" className="flex gap-1">
      {tabs.map((tab, i) => (
        <button key={tab} type="button" role="tab" aria-selected={i === active}
          tabIndex={i === active ? 0 : -1}
          onKeyDown={(e) => handleKeyDown(e, i)}
          onClick={() => onChange(i)}
          className={`px-4 py-2 text-sm rounded-lg transition-colors
            ${i === active ? "bg-bg-active text-fg-heading" : "text-fg-secondary hover:text-fg"}`}>
          {tab}
        </button>
      ))}
    </div>
  )
}
```

---

## Semantic HTML Checklist

- One `<h1>` per page. Headings are hierarchical (`h1` > `h2` > `h3`). No skipping levels.
- Use `<main>`, `<nav>`, `<header>`, `<footer>`, `<aside>`, `<section>` for landmarks.
- Use `<button>` for actions, `<a>` for navigation. Never `<div onClick>`.
- Lists of items use `<ul>`/`<ol>` + `<li>`.
- Every `<img>` needs `alt`. Decorative images use `alt=""`.
- Form inputs need `<label htmlFor="id">`.
- Tables need `<thead>`, `<th scope="col">`.

---

## ARIA Patterns

### When to Use ARIA

The first rule of ARIA: **don't use ARIA if native HTML can do the job.** Only reach for ARIA when building custom widgets with no native equivalent.

### Common Patterns

| Widget | Key Attributes |
|--------|---------------|
| Modal | native `<dialog>` (role, modality and `aria-modal` are implicit); a custom one needs `role="dialog"`, `aria-modal="true"`, `aria-labelledby` |
| Toast/Alert | `role="alert"` or `aria-live="polite"` |
| Tab | `role="tablist"` + `role="tab"` + `role="tabpanel"`, `aria-selected` |
| Accordion | `aria-expanded`, `aria-controls` |
| Dropdown | `aria-haspopup`, `aria-expanded`, `role="menu"` + `role="menuitem"` |
| Toggle | `aria-pressed` (button) or `role="switch"` + `aria-checked` |
| Loading | `aria-busy="true"` on the container |

### Common Mistakes

- Adding `role="button"` to a `<button>` (redundant).
- Using `aria-label` when visible text exists (use `aria-labelledby` instead).
- Forgetting `aria-live="polite"` on dynamically inserted content.
- Using `aria-hidden="true"` on focusable elements (creates ghost focus).

---

## Reduced Motion

Respect `prefers-reduced-motion` to avoid triggering motion sickness or seizures.

### CSS Fallback

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```

### Motion Integration

Use the `useReducedMotion` hook to swap animated variants for instant ones.

```tsx
import { useReducedMotion, motion } from "motion/react"
import type { ReactNode } from "react"

function FadeIn({ children }: { children: ReactNode }) {
  const reduced = useReducedMotion()

  return (
    <motion.div
      initial={reduced ? false : { opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={reduced ? { duration: 0 } : { duration: 0.5 }}
    >
      {children}
    </motion.div>
  )
}
```

What to disable vs. simplify:
- **Disable:** parallax, scroll-hijacking, auto-playing carousels, large translate/scale animations
- **Simplify:** fade-in (keep but make instant), hover scale (reduce to opacity change)

---

## Testing

### Automated — a Playwright spec against the Go-served build

The SPA has no `main.tsx` to hook `@axe-core/react` into, and the harness has no CI: the accessibility gate is one e2e spec that runs against the Go binary serving `dist/client`, with `@axe-core/playwright`. `typescript/references/testing.md` has the `playwright.config.ts` that starts the Go binary as `webServer`.

```ts
// web/e2e/a11y.spec.ts
import { test, expect } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'

for (const theme of ['light', 'dark'] as const) {
  test(`home has no axe violations (${theme})`, async ({ page }) => {
    await page.addInitScript((t) => localStorage.setItem('theme', t), theme)
    await page.goto('/')
    const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze()
    expect(results.violations).toEqual([])
  })
}

test('every focused control shows the ring', async ({ page }) => {
  await page.goto('/')
  for (let i = 0; i < 10; i++) {
    await page.keyboard.press('Tab')
    const outline = await page.evaluate(() => {
      const el = document.activeElement
      return el && el !== document.body ? getComputedStyle(el).outlineWidth : null
    })
    if (outline === null) break            // tabbed past the last control
    expect(outline).toBe('2px')
  }
})
```

Also: **Lighthouse** accessibility score 90+ as a baseline; **eslint-plugin-jsx-a11y** for common JSX mistakes at lint time.

### Manual Checklist

1. Tab through entire page — can you reach and operate every interactive element?
2. Is focus always visible, on both themes?
3. Can you close every modal/drawer with Escape?
4. Does content make sense with a screen reader? (VoiceOver on Mac: Cmd+F5)
5. Zoom to 200% — does layout remain usable?
6. Enable "Reduce motion" in OS settings — do animations respect it?
