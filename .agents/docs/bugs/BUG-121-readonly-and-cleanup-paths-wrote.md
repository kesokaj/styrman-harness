---
owner: bugs
date: 2026-09-30
id: BUG-121
status: fixed
area: internal/search (ListProjects), internal/dream (AgyConversationText)
---
# BUG-121: a read-only path and a cleanup path wrote to the home

## Symptom
Found by my own scratch runs, before commit. (a) `dream --dry-run` changed the home. (b) `AgyConversationText` left an empty `~/.config/sh-agent` behind.

## Repro
(a) Run `sh-agent dream --dry-run` and list the home before and after: `-wal`/`-shm` files appear beside the global index. (b) Run the reader in a sandbox home and list it afterwards.

## Root cause
(a) `search.ListProjects` opened the WAL-mode global index read-only, which still creates `-wal`/`-shm`. (b) Two defers ran in reverse order, so the directory removal ran before the file removal it depended on.

## Fix
Each fixed test-first. (a) read a copy of the index; TestListProjects compares the directory listing. (b) one deferred function; the test asserts the directory is gone.

## Prevention
A read-only or dry-run path gets a test that lists the directory before and after; cleanup order lives in one deferred function.

## Related
BUG-113 (skills/dream); the global rule that read-only modes assert they do NOT write.
