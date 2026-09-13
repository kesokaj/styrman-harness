---
owner: bugs
date: 2026-09-13
id: BUG-010
status: fixed
area: tooling/batch-edits
---
# BUG-010: Python batch-edit script silently left later files unedited after an assert failed on an earlier one

## Symptom
Twice in one day, a Python script of the form `for file: for (old,new): assert old in text; write` was used to
apply a set of text replacements across many files. In both runs the `assert` failed partway through — once on a
stale no-op pair left in the script, once on a README table row whose wording differed from what was assumed.
The script raised and stopped there. Files scheduled before the failing one had already been written; files
scheduled after it had not. `go test` still passed, so nothing flagged the half-applied tree.

## Repro
Run a script shaped like:
```python
for file, pairs in edits:
    text = read(file)
    for old, new in pairs:
        assert old in text
        text = text.replace(old, new)
    write(file, text)
```
with one `old` string that does not match the file's actual current content. Files after the failure point are
never touched; the script's non-zero exit is the only signal.

## Root cause
One script editing many files, fail-fast in the middle, no per-file confirmation and no verification that every
edit landed. The agent treated "assert then write" as safe without noticing that a mid-run failure leaves earlier
files written and later files silently untouched, with no summary distinguishing the two groups.

## Fix
Not a code fix — the remaining edits were identified from the traceback and re-applied by hand both times; no
wrong content shipped, but it required noticing the failure and manually diffing.

## Prevention
Before writing any file, check every `old` string against every target file in a first pass (or process one file
per script invocation) so a stale assumption fails before anything is written. Have the script print a
"written: <file>" line per file so success/failure is legible from output, not inferred. After any batch edit,
run `git status --short` and `git diff --stat` against the list of files the edit was supposed to touch, before
trusting the change is complete.

## Related
BUG-009 (an assertion written without first reading the code's own documented rule). Both are "agent asserted a
fact about the tree instead of verifying it first" — three occurrences of this shape would make a case for a
project rule on verify-before-assert in scripted edits.
