---
owner: main
date: 2026-09-17
---
# Verification — one real `sh-agent report` from this box (2026-09-17)

## Why this run happened

Phase 3 of `.agents/docs/plans/2026-09-17-release-binary.md` rests on an assumption no local test can settle: that GitHub serves a `--filter=blob:none` partial clone over SSH and accepts a push from a `--depth 1` clone. The operator authorised exactly one real run for it in their own turn ("Yes, when Phase 3's code is ready"), recorded as plan decision D25. It ran after the local-bare-repository tests passed.

## What ran

Working tree: branch `release-binary` at 0eec14a, clean. Binary: `go run ./cmd/sh-agent` from that tree. Real home, real feedback repository (`git@github.com:kesokaj/styrman-harness-feedback.git`). Host `dave`.

1. `go run ./cmd/sh-agent report --dry-run` → exit 0. Named the throwaway clone (`/home/kesokaj/.config/sh-agent/report-102687`, "blob-less and limited to reports/dave/, and remove it after the push") and the five files it would write: `feedback.jsonl`, `retros.jsonl`, `state.json`, `proposals/2026-09-17-retro.md`, `report.json`. No redaction fired, so nothing was reported as caught — and by the dry run's own documented limit, that is not evidence that nothing needed catching.
2. `go run ./cmd/sh-agent report` → exit 0. Staged stat: five files, 430 insertions. Output: `report: sent reports/dave/2026-09-17-200615 (09586ff)`.

## Result

- **The partial clone works over SSH.** The clone step (`--filter=blob:none --no-checkout --depth 1`) succeeded against GitHub, and `git sparse-checkout set --no-cone reports/dave/` then `checkout` produced a working tree holding this host's directory only.
- **A push from a shallow, blob-less clone is accepted.** `push -u <url> HEAD:main` landed as `09586ff`.
- **Nothing was kept.** `ls -d ~/.config/sh-agent/report-*` → none. `~/.config/sh-agent/` holds `feedback.jsonl`, `retros.jsonl`, `state.json`, `proposals/` and the old `feedback` clone, which this run correctly ignored.
- **Timing:** the push stalled about four minutes with the SSH session established and bytes queued in its send buffer (`ss -tnp` showed `ESTAB` to GitHub, the process polling), then completed with exit 0. Slow, not refused; the first attempt exceeded a 120-second foreground timeout and was watched in the background.

## Not exercised by this run

- That the real clone never held another machine's report CONTENT. The tests check this through a hook before removal (`TestSendReportThrowawayClone`, with `GIT_NO_LAZY_FETCH=1`); the real clone was removed as designed, so it cannot be inspected after the fact. What this run shows is that the server accepted the filter, not what the clone contained.
- The rejected-push retry against the real remote: no other place pushed during this run. Covered locally by `TestSendReportRetriesRejectedPush`.
- `RemoveKeptClone` against this box's real `~/.config/sh-agent/feedback`. It is still there on purpose: removing it is `sh-agent update`'s job in the operator's own home, not this session's.
- Everything Phase 4 and Phase 5 own.
