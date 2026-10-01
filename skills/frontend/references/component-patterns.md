# Component Patterns

Reference file for the frontend skill. Read this when building interactive UI components in TanStack Start (SPA mode: the Go API is the only backend; there are no server functions).
Every class is a theme token utility (`theme-system.md`): `bg-bg`, `text-fg`, `bg-accent`, `outline-ring`. No palette colors, no raw hex, no inline `style` objects for colors.

## Navigation

### Sticky Header with Scroll-Aware Transparency

The header overlays the hero and gains a blurred background on scroll. Avoids consuming viewport budget. The background is the `--bg` token with a scroll-linked alpha, so it follows the theme.

```tsx
import { motion, useScroll, useTransform, useMotionTemplate } from "motion/react"

function Header() {
  const { scrollY } = useScroll()
  const alpha = useTransform(scrollY, [0, 100], [0, 85])      // percent
  const blur = useTransform(scrollY, [0, 100], [0, 12])

  return (
    <motion.header
      className="fixed top-0 inset-x-0 z-50 border-b border-border/50"
      style={{
        backgroundColor: useMotionTemplate`color-mix(in oklab, var(--bg) ${alpha}%, transparent)`,
        backdropFilter: useMotionTemplate`blur(${blur}px)`,
      }}
    >
      <nav className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
        <a href="/" className="text-lg font-semibold text-fg-heading">Brand</a>
        <div className="hidden md:flex items-center gap-8">
          <a href="#features" className="text-sm text-fg-secondary hover:text-fg transition-colors">Features</a>
          <button type="button" className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-fg-inverse hover:bg-accent-hover">
            Get Started
          </button>
        </div>
      </nav>
    </motion.header>
  )
}
```

### Mobile Off-Canvas Navigation

Slide-in panel: focus moves to the close button, ESC closes, backdrop click closes.

```tsx
import { AnimatePresence, motion } from "motion/react"
import { useEffect, useRef, useState } from "react"

function MobileNav() {
  const [open, setOpen] = useState(false)
  const closeRef = useRef<HTMLButtonElement>(null)

  useEffect(() => { if (open) closeRef.current?.focus() }, [open])
  useEffect(() => {
    const h = (e: KeyboardEvent) => { if (e.key === "Escape") setOpen(false) }
    document.addEventListener("keydown", h)
    return () => document.removeEventListener("keydown", h)
  }, [])

  return (
    <>
      <button type="button" onClick={() => setOpen(true)} aria-label="Open menu" className="md:hidden p-2 text-fg">☰</button>
      <AnimatePresence>
        {open && (
          <>
            <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}
              onClick={() => setOpen(false)} className="fixed inset-0 z-40 bg-overlay" aria-hidden="true" />
            <motion.nav initial={{ x: "100%" }} animate={{ x: 0 }} exit={{ x: "100%" }}
              transition={{ type: "spring", damping: 25, stiffness: 200 }}
              className="fixed right-0 top-0 bottom-0 z-50 w-72 bg-bg-secondary p-6 shadow-lg"
              role="dialog" aria-modal="true" aria-label="Navigation">
              <button ref={closeRef} type="button" onClick={() => setOpen(false)} aria-label="Close menu"
                className="mb-8 p-2 text-fg-secondary hover:text-fg">✕</button>
              <ul className="space-y-4">
                <li><a href="#features" className="text-lg text-fg">Features</a></li>
                <li><a href="#pricing" className="text-lg text-fg">Pricing</a></li>
              </ul>
            </motion.nav>
          </>
        )}
      </AnimatePresence>
    </>
  )
}
```

---

## Modals & Dialogs

**Default: native `<dialog>` with `showModal()`.** The browser gives the focus trap, Escape, `aria-modal`, inert background and the top layer (no `z-index` fights, no body-scroll hack beyond one CSS line). Compound components for the header, body and footer.

```tsx
import { createContext, useContext, useEffect, useRef, type ReactNode } from "react"

const DialogCtx = createContext<{ close: () => void }>({ close: () => {} })

function Dialog({ open, onClose, title, children }: {
  open: boolean; onClose: () => void; title: string; children: ReactNode
}) {
  const ref = useRef<HTMLDialogElement>(null)

  useEffect(() => {
    const el = ref.current
    if (!el) return
    if (open && !el.open) el.showModal()
    if (!open && el.open) el.close()
  }, [open])

  return (
    <dialog
      ref={ref}
      onClose={onClose}                                     // Escape and el.close() both land here
      onClick={(e) => { if (e.target === e.currentTarget) onClose() }}   // backdrop click
      aria-labelledby="dialog-title"
      className="m-auto w-full max-w-lg rounded-xl bg-bg-secondary p-6 text-fg shadow-lg
        backdrop:bg-overlay open:flex open:flex-col"
    >
      <DialogCtx.Provider value={{ close: onClose }}>
        <div className="mb-4 flex items-center justify-between">
          <h2 id="dialog-title" className="text-lg font-semibold text-fg-heading">{title}</h2>
          <button type="button" onClick={onClose} aria-label="Close" className="text-fg-secondary hover:text-fg">✕</button>
        </div>
        {children}
      </DialogCtx.Provider>
    </dialog>
  )
}

Dialog.Body = ({ children }: { children: ReactNode }) => <div className="text-fg">{children}</div>
Dialog.Footer = ({ children }: { children: ReactNode }) => <div className="mt-6 flex justify-end gap-3">{children}</div>
export const useDialog = () => useContext(DialogCtx)
```

```css
/* styles.css — one line keeps the page from scrolling behind a modal dialog */
html:has(dialog[open]:modal) { overflow: hidden; }
```

**Animated fallback.** When the open/close animation matters more than the free semantics, the portal + `AnimatePresence` version in `animation-recipes.md` § 6 applies; it then owes the focus trap and Escape handling itself. Key rules either way: `role="dialog"` + `aria-modal="true"` (implicit on `<dialog>`), Escape closes, backdrop click closes, focus returns to the opener on close (`<dialog>` does this).

---

## Toast Notifications

Context/Provider pattern for global access. Non-intrusive — use a dialog for critical errors. Color is a stripe plus an icon, the text stays `text-fg` (the semantic tokens fail 4.5:1 as light-mode text — `accessibility.md`).

```tsx
import { AnimatePresence, motion } from "motion/react"
import { createContext, useCallback, useContext, useState, type ReactNode } from "react"

type Toast = { id: string; message: string; type: "success" | "error" | "info" }
const ToastCtx = createContext<{ toast: (msg: string, type?: Toast["type"]) => void }>({ toast: () => {} })
export const useToast = () => useContext(ToastCtx)

const stripe: Record<Toast["type"], string> = {
  success: "border-l-success", error: "border-l-error", info: "border-l-accent",
}

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([])
  const toast = useCallback((message: string, type: Toast["type"] = "info") => {
    const id = crypto.randomUUID()
    setToasts((p) => [...p, { id, message, type }])
    setTimeout(() => setToasts((p) => p.filter((t) => t.id !== id)), 4000)
  }, [])

  return (
    <ToastCtx.Provider value={{ toast }}>
      {children}
      <div className="fixed bottom-6 right-6 z-[100] flex flex-col gap-2" role="status" aria-live="polite">
        <AnimatePresence>
          {toasts.map((t) => (
            <motion.div key={t.id} initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              className={`rounded-lg border border-border border-l-4 bg-bg-secondary px-4 py-3 text-sm text-fg shadow-lg ${stripe[t.type]}`}>
              <span className="sr-only">{t.type}: </span>{t.message}
            </motion.div>
          ))}
        </AnimatePresence>
      </div>
    </ToastCtx.Provider>
  )
}
```

---

## Forms

**Default: no form library.** React 19's `<form action={fn}>` calls `fn(formData)` in a transition, resets uncontrolled fields after EVERY action (an error too: re-fill them from the returned state with `defaultValue` when the input should survive) and needs no `preventDefault`; `useActionState` holds the result and the pending flag. Validation is the Go handler's job; the client shows what the API returns: an AIP-193 error body (the `json-api-style` skill) whose `error.details[]` holds a `google.rpc.BadRequest` with `fieldViolations[{field, description}]`. `FormData` is a valid `fetch` body or, as here, becomes JSON with `Object.fromEntries`.

```tsx
import { useActionState } from "react"

type FormState = { ok: true } | { ok: false; errors: Record<string, string> }

type ApiError = {
  error: { code: number; status: string; message: string; details?: Array<{ "@type": string; fieldViolations?: Array<{ field: string; description: string }> }> }
}

function fieldErrors(body: ApiError): Record<string, string> {
  const out: Record<string, string> = {}
  for (const d of body.error.details ?? []) for (const v of d.fieldViolations ?? []) out[v.field] = v.description
  return Object.keys(out).length ? out : { _: body.error.message }   // no field violation: show the message once
}

async function signup(_prev: FormState, fd: FormData): Promise<FormState> {
  const res = await fetch("/api/signup", {                   // src/lib/api.ts in a real app
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(Object.fromEntries(fd)),
  })
  if (res.ok) return { ok: true }
  const body = (await res.json().catch(() => null)) as ApiError | null        // a 502 page or a mux 404 is not JSON
  return { ok: false, errors: body ? fieldErrors(body) : { _: `${res.status} ${res.statusText}` } }
}

export function SignupForm() {
  const [state, submit, pending] = useActionState(signup, { ok: true })
  const err = (f: string) => (state.ok ? undefined : state.errors[f])

  return (
    <form action={submit} noValidate className="space-y-4">
      <div>
        <label htmlFor="email" className="mb-1 block text-sm text-fg-secondary">Email</label>
        <input id="email" name="email" type="email" required
          aria-invalid={err("email") ? true : undefined}
          aria-describedby={err("email") ? "email-err" : undefined}
          className={`w-full rounded-lg border bg-bg-input px-4 py-2.5 text-fg
            focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring
            ${err("email") ? "border-error" : "border-border-input"}`} />
        {err("email") && <p id="email-err" role="alert" className="mt-1 text-sm text-fg">{err("email")}</p>}
      </div>
      {err("_") && <p role="alert" className="text-sm text-fg">{err("_")}</p>}   {/* no field named (401, 429, 500, a bare 400, a non-JSON reply): still shown */}
      <button type="submit" disabled={pending}
        className="w-full rounded-lg bg-accent py-2.5 text-sm font-medium text-fg-inverse
          hover:bg-accent-hover disabled:opacity-50 transition-colors">
        {pending ? "Submitting…" : "Sign up"}
      </button>
    </form>
  )
}
```

**When a library is justified** — multi-step forms, arrays of fields, cross-field rules, optimistic per-field errors — use `@tanstack/react-form` (same vendor as the router, headless, no resolver package): `useForm({ defaultValues, onSubmit: async ({ value }) => … })`, `<form.Field name="email">{(field) => <input value={field.state.value} onChange={(e) => field.handleChange(e.target.value)} />}</form.Field>`, `validators: { onChange: schema }` accepting any Standard Schema library. `react-hook-form` + `@hookform/resolvers` are two dependencies this stack does not need.

**Zod 4**, only when a library needs a schema or the client validates something the server cannot: top-level formats — `z.email()`, `z.url()`, `z.uuid()` — not the deprecated `z.string().email()`; import from `"zod"`.

Form rules:
- Every input needs a visible `<label>` with `htmlFor`.
- Error association with `aria-invalid` and `aria-describedby`; the message is `text-fg` + `role="alert"`, the input gets `border-error`.
- Show errors inline, not in toasts.
- Disable the submit button while `pending`.
- No `createServerFn`: SPA mode has no server functions; the Go API is the only backend (SKILL.md).

---

## Loading States

### Skeleton Screens

Prefer skeletons over spinners. They communicate structure and prevent layout shift.

```tsx
function Skeleton({ className }: { className?: string }) {
  return <div className={`animate-pulse rounded-lg bg-bg-tertiary ${className ?? ""}`} aria-hidden="true" />
}

function CardSkeleton() {
  return (
    <div className="space-y-3" aria-busy="true" aria-label="Loading">
      <Skeleton className="h-48 w-full" />
      <Skeleton className="h-4 w-3/4" />
      <Skeleton className="h-4 w-1/2" />
    </div>
  )
}
```

### Route-Based Code Splitting

TanStack Router automatically code-splits per route. For component-level splitting, use `lazy()`:

```tsx
import { lazy, Suspense } from "react"
const PricingTable = lazy(() => import("./PricingTable"))

function PricingPage() {
  return (
    <Suspense fallback={<PricingTableSkeleton />}>
      <PricingTable />
    </Suspense>
  )
}
```

For route-level lazy loading, define a `.lazy.tsx` file alongside the route:

```tsx
// routes/pricing.lazy.tsx — loaded on demand when the route is visited
import { createLazyFileRoute } from '@tanstack/react-router'

export const Route = createLazyFileRoute('/pricing')({
  component: PricingPage,
})
```

---

## Section Layout Patterns

### Full-Bleed Section

```tsx
import type { ReactNode } from "react"

function FullBleedSection({ children, className }: { children: ReactNode; className?: string }) {
  return (
    <section className={`w-full ${className ?? ""}`}>
      <div className="mx-auto max-w-7xl px-6">{children}</div>
    </section>
  )
}
```

### Alternating Two-Column (Cardless Features)

```tsx
function Feature({ title, description, imageSrc, imageAlt, reverse }: {
  title: string; description: string; imageSrc: string; imageAlt: string; reverse?: boolean
}) {
  return (
    <div className={`flex flex-col gap-8 py-16 md:flex-row md:items-center md:gap-16
      ${reverse ? "md:flex-row-reverse" : ""}`}>
      <div className="flex-1">
        <h3 className="text-2xl font-semibold text-fg-heading">{title}</h3>
        <p className="mt-3 text-fg-secondary leading-relaxed">{description}</p>
      </div>
      <div className="flex-1">
        <img src={imageSrc} alt={imageAlt} className="rounded-xl" width={600} height={400} loading="lazy" />
      </div>
    </div>
  )
}
```

### Responsive Data Table

On mobile, stack rows vertically. On desktop, standard table.

```tsx
function ResponsiveTable({ headers, rows }: { headers: string[]; rows: string[][] }) {
  return (
    <>
      <table className="hidden w-full md:table">
        <thead>
          <tr className="border-b border-border">
            {headers.map((h) => <th key={h} scope="col" className="py-3 text-left text-sm font-medium text-fg-secondary">{h}</th>)}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={i} className="border-b border-border/50">
              {row.map((cell, j) => <td key={j} className="py-3 text-sm text-fg">{cell}</td>)}
            </tr>
          ))}
        </tbody>
      </table>
      <div className="space-y-4 md:hidden">
        {rows.map((row, i) => (
          <div key={i} className="rounded-lg border border-border p-4 space-y-2">
            {row.map((cell, j) => (
              <div key={j} className="flex justify-between">
                <span className="text-sm text-fg-secondary">{headers[j]}</span>
                <span className="text-sm text-fg">{cell}</span>
              </div>
            ))}
          </div>
        ))}
      </div>
    </>
  )
}
```
