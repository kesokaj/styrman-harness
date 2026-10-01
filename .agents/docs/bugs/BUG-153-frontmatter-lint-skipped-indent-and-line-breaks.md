---
owner: bugs
date: 2026-10-01
id: BUG-153
status: fixed
area: internal/harness/frontmatter.go (ParseStrict, unreadLine, entries), permission-key lint
---
# BUG-153: two frontmatter shapes hid a permission key from the installer's lint while a YAML reader still reads it

## Symptom
1. A block whose lines are all indented (`---`, `  description: d`, `  allowed-tools: Bash(*)`, `---`) loaded in a command or a rule with no keys seen: the key grammar took every indented line for a continuation, of nothing. YAML reads such a block as the top-level map, so a reader that parses it as YAML gets `allowed-tools`. Skills and agents were refused only by accident (no description found).
2. A line break the lint does not split on inside a key line (`description: d<CR>allowed-tools: Bash(*)`, or U+0085, U+2028, U+2029 for the CR) loaded in all four kinds with only the first key seen. YAML 1.2 treats a lone carriage return as a line break. Present in v2.0.0 to v2.2.0.

## Repro
The reviewer ran each shape through `LoadInventory` on a scratch copy, with must-fail controls (`reviews/2026-10-01-security-v2.2.0-1fc0d4c.md`, findings 1 and 2, both medium). Not verified on the reader side: Claude Code's YAML parser could not be run in isolation and Antigravity's reader was not examined, so whether either reads the key is unknown.

## Root cause
The lint's grammar is line-based and decides for itself what a line and a top-level line are, more narrowly than the YAML readers it guards against: it trusted indentation to mean "continuation" even with no key above, and `\n` to be the only line break. Each skipped shape was a place the other parser could find a key.

## Fix
`ParseStrict` refuses a block holding a line break that is not a newline (a carriage return not directly before a newline, U+0085, U+2028, U+2029) and a block whose first content line is indented (`indentedBeforeKey`), in `internal/harness/frontmatter.go`. Test `TestLintSeesEveryKeyAYAMLReaderSees` in `internal/harness/harness_test.go`: six shapes in each of the four kinds (24 cases, all red before the fix), each asserting the refusal is the lint's own and not a side effect, plus shapes that must keep loading (a nested map, an indented list, a CRLF file, a leading comment, a folded skill description).

## Prevention
A lint that guards against what ANOTHER parser will read must refuse everything it cannot itself read the same way: whatever it skips as harmless (indentation, a character it does not treat as a line break, a line it calls a continuation) is a place the other parser may find a key. For each skip rule, write the must-refuse case in which the skipped text holds the forbidden key, in every kind the lint covers.

## Related
BUG-143 (the lint this bypassed), BUG-151 (same family, fixed at the fence only, one commit earlier), BUG-152 (same gate, same lesson). Third record of the family (BUG-143, BUG-151, BUG-153): the rule has recurred three times. Recommended home: `skills/review-standards` § 8 (code-shape rules a review checks); actionable (a must-refuse case per skip rule catches it) and generalizable (any lint or filter guarding another parser). Not a `context/` addition, so no size-pin cost; merge into the existing § 8 lint line.

## Addendum 2026-10-01: a third shape (tab-led line)
The security read of 1fc0d4c..4165f81 (`reviews/2026-10-01-security-1fc0d4c-4165f81.md`, finding 1, low, confidence low) found it: a frontmatter line that STARTS WITH A TAB after a column-0 key (also after a nested map, a list, an empty value) still loaded, and the harness folded it into the key above as a continuation. js-yaml 4.3 and yaml 2.9 error on such a block; Claude Code's fallback nests the line; an older lenient reader would likely read a top-level key (not run).
Fix: `tabLed` in `internal/harness/frontmatter.go`, called first by `ParseStrict`: any frontmatter line starting with a tab and holding more than white space is refused ("starts with a tab: YAML has no tab indentation ..."). Test: five more shapes in each of the four kinds in `TestLintSeesEveryKeyAYAMLReaderSees` (tab-indented first line; tab-led key after a key, a nested map, a list, an empty key). 16 of the 20 cases were red before the fix (the tab-indented first line was already refused by `indentedBeforeKey`, with no test until now), plus two shapes that must keep loading (a tab inside a value, a tab-only line).
