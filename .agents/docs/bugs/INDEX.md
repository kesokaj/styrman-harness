# Bugs

<!-- The project's memory of what went wrong. One row per BUG-NNN record in this directory; newest last. Maintained by the `bugs` agent. Read at resume; consult before touching an area listed here. -->

| ID | Date | Area | Status | Symptom | Prevention |
|---|---|---|---|---|---|
| BUG-001 | 2026-09-03 | installer/status | fixed | `status` reported the installed binary as stale right after `install` | status and diff plan without the binary; exclude `BinPath` from stale checks and test it |
| BUG-002 | 2026-09-03 | cli/flags | fixed | `sh-agent adopt <path> --home DIR` failed with a usage error | Go's flag package stops at the first positional; reorder flags ahead of positionals before parsing |
