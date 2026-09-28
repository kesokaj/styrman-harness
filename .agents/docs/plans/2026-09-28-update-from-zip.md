---
owner: main
date: 2026-09-28
---
# Plan — `sh-agent update --from-zip <file>`   (2026-09-28)

## Why

A machine that can neither hold the operator's SSH key nor run `gh`, such as a corporate workstation, cannot reach the harness today. A checkout needs git over SSH, and `sh-agent update` on a release binary downloads through `gh` (R8, R9). The operator can download the repository zip by hand in a browser.

The operator's words, 2026-09-28:
- "add the --from-zip option so i can run manually";
- then "approve as written, go all but number 6".

That approval covers the REQUIREMENTS text quoted under Interfaces, scope items 1–5 and 7, and, for item 7, the recommended matrix option 1. Item 6 is out: the pasted proposal's report stripping, `trigger: always_on` frontmatter, Cloudtop overlay and skill filtering.

Recon: `.agents/docs/plans/2026-09-28-recon-from-zip.md`.

## Scope

- **In:**
  - `sh-agent update --from-zip <file>` from any binary, release or development build (D7 says what other commands do afterwards). It unpacks the zip into a temporary directory under `~/.config/sh-agent/`, refuses any entry that would leave it, installs that tree with the running binary (no rebuild, no binary replacement), and deletes the directory.
  - The manifest marks the install as coming from a zip, and `status`, `knownCheckout` and `ResolveSource` handle that marker.
  - The REQUIREMENTS R8/R9 text as approved; ADR-038; the README map and the delivery docs; `AGENTS.md` rule 1.
  - Then item 7: the matrix lists, and the v1.8.0 release.
- **Out:**
  - fetching anything;
  - verifying the zip's signature, since a zip is unsigned and trust rests on the operator's own download;
  - replacing the binary from a zip;
  - item 6 of the pasted proposal;
  - a manifest schema change.
- **Follow-ups noted, not built:**
  - **A signed zip:** `/preflight` would publish a `git archive` zip under the release's signed `SHA256SUMS`, and `--from-zip` would verify it through the existing `ssh-keygen` path, refusing or clearly flagging an unverified zip. It changes R9's Release procedure, so it needs the operator's approval. Until then the command prints `UNVERIFIED` (security reviews r1 and r2, medium, open by decision).
  - `skills/find-skills/SKILL.md:97` names `styrman-harness` in a vendored skill: an old in-place edit, found during the item-4/5 sweep.

## Surfaces touched

- `internal/harness/zip.go` (new): `UnpackZip`. `internal/harness/zip_test.go` (new).
- `internal/harness/source.go`:
  - `ZipPrefix`;
  - a `Source.Origin` field the manifest records instead of `Dir`, when set.
- `internal/harness/apply.go:117-123`: the manifest's `Repo`/`Commit` for a zip source.
- `internal/harness/repo.go:75`: the `ResolveSource` candidate walk skips a `ZipPrefix` manifest value, as it skips `EmbeddedPrefix`.
- `cmd/sh-agent/main.go`:
  - the flag (`:84-90`);
  - the `reorder` value table (`:617`);
  - the usage string;
  - the `update` case (`:198-230`), with a `--from-zip` branch placed before `ResolveSource` (`:123`);
  - `status` (`:346-472`, the `from` line `:382-385` and the drift block `:444-471`);
  - `knownCheckout` (`:573`).
- `cmd/sh-agent/main_test.go`: the CLI tests.
- `REQUIREMENTS.md` R8 and R9 (approved text); `.agents/docs/architecture/DECISIONS.md` ADR-038.
- `README.md`: § How it works (the map) and the install and update sections. `.agents/docs/specs/sh-agent/delivery.md` § How a machine gets the harness and § The model. `.agents/docs/deploy/{local,remote}.md` where they state the number of ways. `AGENTS.md` rule 1 ("Two ways in").
- `.agents/docs/testing/skill-matrix.md` (item 7).

## Interfaces

- **The command.** `sh-agent update --from-zip <file> [--home <dir>] [--vendor <v>] [--dry-run]`. `--from-zip` takes a value and joins `reorder`'s value table.
  - `--from-zip` with `--check` or with `--repo` is refused before anything is read: "--from-zip installs the named zip; --check/--repo do not apply to it".
  - `$SH_AGENT_REPO` is ignored on this path. The operator named the file, so D15 holds.
- **Order of work.**
  1. `sha256` of the file.
  2. `cfg, err := harness.ConfigDir(home)` (it returns an error too), then `harness.UnpackZip(file, cfg)`, which returns the tree's directory.
  3. `IsRepo(dir)`; failing it is an error naming the zip.
  4. `src := harness.Source{FS: os.DirFS(dir), Dir: dir, Origin: harness.ZipPrefix + "sha256:" + hex}`.
  5. `installFiles(src, home, vendor, dryRun, out)`, the same helper `install` uses: the running binary via `selfExecutable`, never `Reinstall`.
  6. `os.RemoveAll` of the unpack directory in a `defer`, on every path, errors and dry run included.
- **Output.**
  - First line: `zip <file> (sha256 <hex>): installing its harness with the running binary <version>; the binary itself is not replaced`, or the same with a `[dry-run]` prefix.
- **`harness.UnpackZip(zipPath, parent string) (dir string, err error)`**
  - It creates `os.MkdirTemp(parent, "zip-*")`.
  - It reads the archive with `archive/zip`.
  - Checks run at two scopes (D6).
    - **Archive-wide, on every entry:** the name must be relative, contain no `..` element and no backslash, and be `filepath.IsLocal` after cleaning. The archive may hold at most 5,000 entries. A breach refuses the whole zip, with nothing installed and the temp directory removed.
    - **Only on entries under the five extracted trees:** a symlink or any non-regular, non-directory mode refuses the whole zip. So does more than 16 MiB uncompressed for one file, or more than 128 MiB in total across the extracted files, counted from bytes actually written through an `io.LimitReader`, never from the header.
    - An entry outside the five trees is never opened and never written, whatever its mode. That includes this repository's own tracked symlink `.agents/skills`, which GitHub's zip keeps as a symlink entry.
  - **Layout.**
    - When every entry sits under one top-level folder and that folder holds `context/AGENTS.md`, the folder is stripped. That is GitHub's "Download ZIP" (`styrman-harness-main/…`) and its API zipball (`<owner>-<repo>-<sha>/…`).
    - Otherwise, when `context/AGENTS.md` is at the root, the root is used.
    - Otherwise the error is "not a harness zip: no context/AGENTS.md at its root or under one top folder".
  - **Extracted:** only entries under `context/`, `skills/`, `agents/`, `commands/` and `rules/`, the five trees the binary embeds. Every other entry is still name-checked, but never written.
  - File mode `0o644`, directories `0o755`.
- **`State.Repo` for a zip install:** `zip:sha256:<hex>`. `State.Commit` is `""`. `ZipPrefix = "zip:"`.
  - In `Install` (`apply.go:117-123`, which today branches on `src.Dir != ""` only), `src.Origin != ""` is checked FIRST and wins: `repo, commit = src.Origin, ""`. Then comes the existing `Dir` branch (path plus commit), and then the embedded branch. A zip source sets both `Dir` (the temp directory) and `Origin`; the temp path must never reach the manifest.
  - `ResolveSource` (`repo.go:75`) skips it as a candidate.
  - When `ResolveSource` finds no tree and the manifest's `Repo` has `ZipPrefix`, the not-found error (`repo.go:87`) gains, first: "this home was installed from a zip (sha256 <first 12 hex>); a development build needs a checkout or --repo — or run `sh-agent update --from-zip <file>` again". The rest of the message is unchanged (D7).
  - `knownCheckout` (`main.go:573`) treats it like `EmbeddedPrefix`: "this home was installed from a zip (sha256 <hex>), not from a checkout; pass --repo <dir> to build and install from it".
- **`status` for a manifest whose `Repo` has `ZipPrefix`.** `status` moves `harness.LoadState(home)` (today at `main.go:370`) to the top of the function, before the `source`/`repo` line (`:349-368`), so that line can append the zip note. The `installed nothing` early return keeps its place and message after the source line.
  - `installed  <time> from zip sha256:<first 12 hex> (…)`.
  - The drift block is replaced by a presence check over the manifest's files: `drift      not compared — installed from a zip that is not kept; <n> of <m> manifest files present`, then `  missing  ~/<rel>` for each absent one. The exit status is 1 when any is missing, 0 otherwise.
  - The `source` line for a release binary adds: `; this home was installed from a zip — run update --from-zip <file> again to change it`.
- **REQUIREMENTS, exactly as approved.**
  - R8's first sentence is replaced by: "The harness reaches a machine one of three ways: a checkout the operator owns; a release's `sh-agent` binary with the harness built in; or a zip of the repository that the operator downloads with their own credential and passes to `sh-agent update --from-zip <file>`. Each is kept current with `sh-agent update`: the checkout and the binary as today, and a zip place by the operator repeating the download and `--from-zip`."
  - R9's Install/Update procedure gains, as an indented sub-bullet under **Update** (D9), not as a fifth top-level procedure: "**Zip place:** the operator downloads the repository zip by hand and runs `sh-agent update --from-zip <file>`. It unpacks the zip into a temporary directory, refuses any entry that would leave it, installs that tree with the running binary, and deletes the directory. Nothing is fetched, no repository is created, and the binary itself is not replaced."
  - Both keep their approval notes: "(Approved by the operator 2026-09-28, on the exact text above: "approve as written"; ADR-038.)".
- **ADR-038.** A third way in, the zip place.
  - It supersedes ADR-034 for this one path only: a release binary installs a tree on disk solely when the operator names a zip with `--from-zip`, never a recorded or discovered path.
  - It amends the "two ways, no third" claim in `delivery.md`.

## Failure modes

- **Not a zip, or corrupt.** Error: "<file>: not a readable zip: <err>". Nothing is installed, the temp directory is removed, and the exit status is non-zero.
- **An unsafe entry** (slip, symlink, over a limit). The error names the entry and the rule it broke. Nothing is installed and the temp directory is removed.
- **Not a harness zip.** The error says so, and names the expected layout.
- **The zip carries a harness that needs a newer binary.** The install proceeds with the running binary; the first output line says the binary is not replaced. When the zip's `skills/` or `agents/` use something this binary cannot map (for example an unmapped tool in `derive.go`), the existing plan error stops the install before any write. That is the existing guard, unchanged.
- **The temp directory cannot be created** (unwritable `~/.config/sh-agent`). Error, and nothing installed.
- **`status` after a zip install** reports presence, not content. Content drift against the zip cannot be computed; the output says so rather than pretending.
- **A zip-installed home, then a development build with no checkout nearby and no `--repo`.** `status`, `diff`, `list`, `render` and `install` fail as they do today for a development build with no tree, and the error names the zip install (D7). A release binary never hits this: its source is its embedded files.
- **Dry run.** It unpacks and plans, prints what would change, and removes the temp directory. It writes no installed file, no manifest and no binary.

## Verification

- **Tiers:**
  - unit (`zip_test.go`: zip fixtures built in the test with `archive/zip`'s writer);
  - CLI (`main_test.go`);
  - an e2e with a real GitHub zip, using a release-versioned build;
  - security review (input parsing, a new way in);
  - `/review --docs`;
  - `TestHarnessMap`, which pins verbs; `--from-zip` is a flag, but the map gains a line.
- **E2E criteria:**
  1. **Golden path.** A binary built with a canonical release version runs `update --from-zip <GitHub zipball of main> --home <scratch>`. It installs, `status --home <scratch>` names the zip and reports every manifest file present, and the installed `~/.claude/AGENTS.md` equals the zip's rendered context.
  2. **Edge case.** A zip with an entry `x/../../escape` is refused, nothing is written under `<scratch>`, and no `zip-*` directory is left under `<scratch>/.config/sh-agent/`.
  3. **Dry run.** It leaves `<scratch>` byte-identical to before (manifest absent or unchanged), with no `zip-*` directory left.

## Phases

- [x] Cold-reader check — rounds 1–3 (`reviews/2026-09-28-coldread-update-from-zip-r{1,2,3}.md`); round 3 ended at the cap with one narrow residue (R9's structure), put to the operator, whose answer was "Go with your recommendations" (2026-09-28): D9

### Phase 1 — Zip source (internal/harness)

- [x] `UnpackZip` with the layout rules, limits and refusals — verify: `go test -race -run 'TestUnpackZip' ./internal/harness`, AI — expect: tests for GitHub top-folder strip, root layout, `..` slip, absolute name, backslash name, symlink entry INSIDE `skills/` (refused), symlink entry OUTSIDE the five trees such as `.agents/skills` (accepted, not written), over-count, over-size (written bytes), not-a-harness, and "entries outside the five trees not written" all pass. Each is written first and seen failing against a zero-value `UnpackZip`.
- [x] `ZipPrefix`, `Source.Origin`, the manifest record (`apply.go`), and the `ResolveSource` skip — verify: `go test -race ./internal/harness`, AI — expect: a test that an `Origin` source records `zip:sha256:…` with an empty commit, and that `ResolveSource` never returns a `zip:` manifest value as a candidate.
- [x] Checkpoint — expect: `go test -race ./...` green, gofmt and vet clean — evidence: `verifier` log `testing/2026-09-28-verify-update-from-zip-phase-1.md` PASS (139/266 pass; the must-fail mutation compiled and failed `one_file_too_large`; coverage 80–100% on the new functions, with uncovered error branches named), 2026-09-28. **No commit at this checkpoint:** the path change lands in ONE commit with its docs and map (`AGENTS.md` rule 3; D2).

### Phase 2 — CLI, status, knownCheckout

- [x] `--from-zip` flag, the `reorder` table and usage; the `update` branch before `ResolveSource`; the `--check`/`--repo` refusals; dry run — verify: `go test -race -run 'TestUpdateFromZip' ./cmd/sh-agent`, AI — expect: golden path into a scratch home; refusals; dry run writes nothing and leaves no `zip-*` directory (asserted absent, not only "normal mode present"); a release-versioned run installs from the zip, not the embedded files. — docs: `README.md` § update
- [x] `status` zip lines and presence check; `knownCheckout` zip branch — verify: `go test -race -run 'TestStatusAfterZip|TestKnownCheckoutZip' ./cmd/sh-agent`, AI — expect: the `from zip` line, "not compared" with the present count, missing files listed with exit 1, and the `knownCheckout` refusal text.
- [x] Checkpoint — expect: the full suite green — evidence: `verifier` log `testing/2026-09-28-verify-update-from-zip-phase-2.md` PASS (mutations (a) dry-run and (b) status both compiled and failed their named tests; restored identical), 2026-09-28. No commit (D2).

### Phase 3 — Requirements, ADR, docs, map, e2e, one commit

- [x] `REQUIREMENTS.md` R8/R9 exactly as approved; ADR-038 — verify: `git diff REQUIREMENTS.md`, AI — expect: only the two approved passages and their approval notes.
- [x] `AGENTS.md` rule 1 ("Two ways in" → three); `documentation` updates `README.md` (map and install/update), `delivery.md` § The model and § How a machine gets the harness, and `deploy/{local,remote}.md` pointers — verify: `/review --docs` over the phase diff, and `go test -run TestHarnessMap ./...`, both — expect: no restatement of the old count left. `git grep -niP '\btwo ways\b|\bno third\b|only way the harness|\bways in\b|one of the ways'` (`\b` because a bare `ways in` matches "always include") over the living docs (`README.md`, `AGENTS.md`, `context/`, `skills/`, `agents/`, `.agents/docs/{specs,deploy}/`, `.agents/docs/testing/guide.md`) prints only lines that state the new three-way model, plus `skills/golang/references/golang-samber.md:431` ("two ways to compose operations", unrelated). Positive control: the same grep at the pre-feature HEAD prints `README.md:182` and `.agents/docs/deploy/remote.md:61` among its hits. The map shows the zip line. — docs: `README.md`, `delivery.md`
- [x] E2E 1–3 — verify: a release-versioned build (`-ldflags "-X main.version=v0.0.1"`) in `.agents/scratch/`; the zip from `gh api repos/kesokaj/styrman-harness/zipball/main`; hand-built slip and dry-run cases; AI — expect: the three criteria above. Proven before trusting: the zip has one top folder, and the slip fixture holds the `..` entry (`unzip -l`).
- [x] `security` over the feature diff; `reviewer` over the diff — expect: findings fixed or rebutted in `## Outcome`.
- [x] Checkpoint — expect: ONE commit carrying the code, tests, REQUIREMENTS, ADR-038, docs and map; pushed — evidence: `verifier` log `testing/2026-09-28-verify-update-from-zip-phase-3.md` PASS; the feature commit follows this tick; 2026-09-28.

### Phase 4 — Matrix (item 7) and release

- [x] `testing/skill-matrix.md`: restore the lists the scenarios need (`bash`, `python3`, `curl` for Claude; the matching named agy commands), keep `edit_file(*)` and `list_dir(*)` out, and state that the allowlist lets scenarios run and is not a security boundary, because the sandbox runs as the operator's OS user; the protection is the throwaway home, with the copied token deleted by name after each batch — verify: read against the v1.8.0 gate code review's findings 1–3 (`reviews/2026-09-28-review-v1.7.0-381e86d-code.md`), AI — expect: each finding closed or rebutted.
- [ ] `/preflight v1.8.0` over `v1.7.0..HEAD` — expect: VERDICT SHIP; tag, binaries, notes.
- [ ] Checkpoint — expect: the release published; handoff written — evidence: the gate report path and release URL, <date>.

## Risks

Assumptions:
- GitHub's API zipball and the web "Download ZIP" both have exactly one top folder. Checked by `unzip -l` on the fetched zipball, the first task of the Phase 3 e2e.
- `archive/zip` exposes a symlink entry through `FileInfo().Mode()&fs.ModeSymlink`. Checked by the symlink unit test's fixture, whose `Mode()` is asserted before the refusal is tested.

## Decisions

- D1 — 2026-09-28 — Unpack under `ConfigDir(home)` (`~/.config/sh-agent/zip-*`), never `/tmp`. This follows the code's convention (`apply.go` `Diff`, `release.go` downloads), and a scratch `--home` keeps tests inside the scratch tree — by: implementer.
- D2 — 2026-09-28 — Phases 1–3 commit once, at the end of Phase 3. `AGENTS.md` rule 3 puts a path change in ONE commit with the code, docs and map, so the feature skill's per-phase commit is waived for these phases. The plan file and verify logs hold the state between them — by: implementer, per `AGENTS.md` rule 3.
- D3 — 2026-09-28 — No manifest schema change. The marker is a `State.Repo` value (`zip:sha256:<hex>`, the `EmbeddedPrefix` precedent), and `status` reports presence instead of content for zip installs. A new manifest field is what the repo's versioning rule calls MAJOR — by: implementer.
- D4 — 2026-09-28 — Extract only the five embedded trees. Name-check every entry, write only those. The zip also carries `.agents/`, `cmd/` and `internal/`, which the install never reads — by: implementer.
- D6 — 2026-09-28 — The symlink, mode and size checks apply only to entries under the five extracted trees; name safety and the entry count apply archive-wide. Otherwise the repo's own GitHub zip, whose `.agents/skills` is a tracked symlink, would be refused — by: cold read round 1 (`reviews/2026-09-28-coldread-update-from-zip-r1.md`, finding 1).
- D7 — 2026-09-28 — `--from-zip` works from any binary. Afterwards, a development build with no checkout and no `--repo` has no source for the other commands; they fail with the existing not-found error, which then names the zip install first. It is not given a way to read the manifest as a source, because nothing of the zip is kept. The zip place's intended binary is a release binary (R8), whose source is its embedded files — by: cold read round 2 (`reviews/2026-09-28-coldread-update-from-zip-r2.md`, finding 1).
- D8 — 2026-09-28 — The existing plain-copy path (an unpacked tree installed by a binary built from it; `update` refuses it) stays exactly as it is, the fallback `delivery.md` already describes. `--from-zip` is the zip place's way in and forward. The docs name both and say which is which: a zip place keeps no tree, while a plain copy is a tree someone unpacked and keeps — by: cold read rounds 1–2, finding 3.
- D9 — 2026-09-28 — The zip text sits inside R9's existing Update procedure, as a nested sub-bullet with the exact wording approved. It is a new kind of place, not a new procedure, so R9's "Four procedures" title and ADR-031 stay true. ADR-038 says so — by: the operator, "Go with your recommendations", on cold read round 3's residue (`reviews/2026-09-28-coldread-update-from-zip-r3.md`).
- D10 — 2026-09-28 — A dry run unpacks under `~/.config/sh-agent/` only when that directory already exists, and otherwise under the system temp directory. `ConfigDir` creates the directory, which a dry run must not do — by: implementer, Phase 1.
- D11 — 2026-09-28 — Supersedes the Interfaces lines on the output, `UnpackZip`'s signature, and "`--from-zip` with `--check` or `--repo` is refused", as well as D10's place of directory creation. It follows the review round (`reviews/2026-09-28-review-update-from-zip-{code,docs}.md`, `reviews/2026-09-28-security-update-from-zip.md`):
  - `UnpackZip(zipPath, parent) (dir, sum string, err error)` reads the file once (at most 256 MiB), hashes and parses that one buffer, and validates everything before writing: names (at most 1,024 bytes, fewer than 32 levels), count, kinds, duplicates, and sizes counted by decompressing. Only then does it create the parent chain, and on any error it removes everything it created.
  - `--from-zip` is detected with `fs.Visit` and refused before ANY command runs, when used with a command other than `update`, with an empty value, or with `--check`/`--repo`.
  - The first output line says `UNVERIFIED`, and that the running binary is installed but nothing in the zip replaces it.
  - `diff` does not compare a zip install.
  - A real run sweeps `zip-*` directories older than an hour.
  - `update --dry-run`, on the checkout and release paths, is its `--check`: it reports and changes nothing.

  — by: code review REQUEST CHANGES findings 1–3 and lows; security lows 2 and 3; docs lens 1–3.
- D12 — 2026-09-28 — The approved R9 sentence "the binary itself is not replaced" is imprecise: the running binary is installed as in every install; what holds is that nothing in the zip replaces it. REQUIREMENTS is not edited without the operator's approval, so the docs and output say the precise thing, and the sentence is put to the operator in the release report — by: implementer, from security review low 3.
- D13 — 2026-09-28 — Supersedes D11's "creates the parent chain … removes everything it created" inside `UnpackZip`, following round 2 (`reviews/2026-09-28-review-update-from-zip-code-r2.md`, `reviews/2026-09-28-security-update-from-zip-r2.md`):
  - `UnpackZip` no longer creates any directory beyond its own `zip-*`. `updateFromZip` creates `~/.config/sh-agent` with `MkdirAllTracked` and removes what it created on ANY failure, the install's own refusals included (BUG-093, third layer).
  - The sweep matches `zip-` names from `os.ReadDir` and never globs; a glob built from `--home 'h[1]'` reached another home.
  - The archive cap is 64 MiB.
  - Every zip CLI test runs with an empty `PATH` and a scratch working directory (BUG-094).

  — by: code review r2 findings 1–3; security r2 lows 2–3.
- D5 — 2026-09-28 — Item 7 takes matrix option 1 (the recommendation), read from "go all but number 6" on a row whose recommendation named the matrix choice — by: implementer, stated to the operator in chat.
- Rejected:
  - reading the zip in place as an `fs.FS` with no unpack: the approved REQUIREMENTS text says it unpacks into a temporary directory;
  - routing through `update()`/`Reinstall`: that rebuilds, and the scope forbids it;
  - per-file hashes in the manifest: a schema change (D3).
