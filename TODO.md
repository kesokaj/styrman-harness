# TODO — styrman-harness

## Open

- [ ] **The changelog row needs a coarser scope than the accountability rows.**
      `changeset:<digest>` moves on every edit, so on a live tree a changelog entry is
      stale the moment it is written, and the changelog agent's own anti-duplication
      rule then declines to write another — observed 2026-09-03: the row failed while
      three valid entries already existed. The accountability rows do not have this
      problem; reviewer and security *should* re-run when code moves. A resumability
      row records a task, which is a human-meaningful unit, not a byte-level one.
      Design question, not a bug: decide the unit before changing the code.
- [ ] **Reviewer row owed for the delta after review round 5.** The fixes made in
      response to that round — collapsing the gate onto `Owed()`, deleting `CoversHead`
      and `HeadScope` — are themselves unreviewed. Five rounds each found real defects,
      so this is a genuine gap rather than a formality.


Everything raised during the 2026-09-03 rebuild is either done — recorded in
`.agents/docs/changes/` — or deliberately settled and written down where it belongs:

- Design decisions and accepted limits live in `REQUIREMENTS.md` under **Known limits**
  and **Non-goals**, not here. A limit that has been investigated and accepted is not a
  task.
- Architectural choices live in `.agents/docs/architecture/DECISIONS.md` as ADRs.
- Planned work lives in `ROADMAP.md`.

Add an item here only when it is genuinely deferred work: something that should be done,
that nobody is doing yet, with the reason it was put off. Keep it short — an item that
has sat here through three sessions is a decision nobody has made, and belongs in
`REQUIREMENTS.md` or the bin.
