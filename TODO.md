# TODO — styrman-harness

Open items. Closed work lives in `.agents/docs/changes/`.

## Needs the operator

- [x] ~~Deploy the Go harness and confirm `harness doctor` passes.~~ Done
      2026-09-03: `./setup.sh` run on this machine, doctor 10/10. The install
      itself found three bugs (see `changes/`), all fixed in `7f57359`.
- [x] ~~Decide whether the agy plugin registered during research should stay.~~
      Moot: `harness install` now manages the registration itself.
- [ ] Record provenance for `gcloud-auth-verification`, `gke-basics`, `skill-repair`
      and `terraform` in `shared/skills/VENDORED.md` — each carries third-party
      `license`/`author` frontmatter but is unlisted. Licensing gap; only the
      operator knows how they were obtained.
- [ ] Decide the coverage model (see "Contract vs implementation" below). It is
      the largest open design question and it hid owed rows five times on
      2026-09-03 alone.
- [ ] Push. Three commits are local-only; nothing has been published.

## Deferred, with reasons

- [ ] Antigravity IDE plugin discovery is unverified — only `agy` is installed on the
      current machine. Confirm whether the IDE needs `agy plugin install` too.
- [ ] `disableBundledSkills`: several harness skills duplicate Claude Code built-ins.
      Needs a per-skill judgement, not a blanket setting.
- [ ] No *automated* integration test drives a real `claude -p` / `agy -p` dispatch.
      Verified by hand on 2026-09-03 — a scratch repo dispatched `changelog` and
      `security` as real subagents and wrote correct frontmatter — but nothing
      in CI covers the CLI contract, and CI has no credentials to run it.
- [ ] `harness doctor` cannot validate agy tool names — agy exposes no way to
      enumerate them. `agy plugin validate` is the closest available check.

## Contract vs implementation — the granularity gap

- [ ] **Rows are satisfied per DAY, but the contract says per COMMIT.**
      `Project.HasArtifact(row, day)` matches any artifact owned by the row and dated
      today, so a second substantive commit on the same day owes nothing — observed
      immediately: a 251-line follow-up round showed only the verifier as due, because
      the morning's review artifact already carried today's date.
      Fixing this needs artifacts to record *what they cover*, not just when. That is
      the same root as the `Scope:` question below, and both should be solved together:
      fold `scope:` into the artifact frontmatter (`<base>..<head>`, or a working-tree
      diff fingerprint) and have the gate match coverage rather than dates.
- [ ] **Two conventions for coverage now sit side by side.** `/preflight` matches
      release ranges against a `Scope:` line in the report *body*; the gates match
      `owner:`/`date:` in the frontmatter. One mechanism, not two.

## From the 2026-09-03 review round

- [ ] Antigravity `PreToolUse` is now declared in the plugin's `hooks.json`, but open
      bug reports say plugin hooks may not fire in the Antigravity IDE. Confirm the
      credential gate actually runs there before relying on it.
- [ ] `harness check` and the resume block have no test for the Antigravity
      `--agy` JSON wrapping; only the plain path is covered.

## Known gaps

- [ ] `documentation` is the one Definition-of-Done row the dispatcher cannot compute:
      "behaviour or interfaces changed" needs judgement. It stays model-triggered.
- [ ] No README. Removed deliberately in `a63a208`; `AGENTS.md` and `REQUIREMENTS.md`
      carry the orientation for now.
