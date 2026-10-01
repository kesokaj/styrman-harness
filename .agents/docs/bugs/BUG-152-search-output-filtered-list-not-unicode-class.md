---
owner: bugs
date: 2026-10-01
id: BUG-152
status: fixed
area: internal/search/search.go (printable)
---
# BUG-152: search output filtered a list of characters, not Unicode's format class

## Symptom
`printable` replaced C0, DEL and C1 control characters only (in v2.2.0 also a hand-written list of line separators and direction controls). Zero-width characters, U+FEFF, the Arabic letter mark U+061C and the tag block U+E0000 to U+E007F passed, so text hidden in a file name or an indexed line reached the terminal and, through search output, a model; tag characters spell words no terminal shows. Present in v2.0.0 to v2.2.0.

## Repro
`printable("a\U000e0041\U000e0042b")` returned its input unchanged; the security read reproduced a sentence written in tag characters coming out whole after a search hit.

## Root cause
A deny-list of code points, extended by one more list when the first gap was reported, where the property that matters is the character's class. Found by the security read of the v2.2.0 gate (`reviews/2026-10-01-security-768efce-f3c3ee4.md` finding 1), left as row L2 of `reviews/2026-10-01-preflight-v2.2.0.md`.

## Fix
`printable` replaces every rune of class Cf (format), Zl and Zp, besides the control characters. Test `TestPrintableDropsLineAndDirectionControls`, five new cases, red first: tag characters, zero-width space and word joiner, U+061C, U+FEFF, soft hyphen. Cost: a zero-width joiner inside an emoji sequence or a zero-width non-joiner in Persian text shows as `?` in search output; the file itself is untouched.

## Prevention
Text shown to a terminal or handed to a model is filtered by Unicode class (`unicode.Cc`, `Cf`, `Zl`, `Zp`), not by a list of code points; a second report against one list is the signal to switch to the class.

## Related
BUG-110 (same lesson for redaction needles: a rule per kind, not another batch); BUG-151 (same gate, exact shape instead of a trimmed one); `reviews/2026-10-01-security-768efce-f3c3ee4.md`.
