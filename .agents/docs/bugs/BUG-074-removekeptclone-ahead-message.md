---
owner: bugs
date: 2026-09-18
id: BUG-074
status: fixed
area: internal/harness/report.go (RemoveKeptClone ahead-count message)
---
# BUG-074: RemoveKeptClone reported "unpushed commits" for commits already on the remote

## Symptom
`sh-agent update` on the operator's own box refused to remove the old kept feedback clone: `kept /home/kesokaj/.config/sh-agent/feedback: has 2 unpushed commit(s)`. The operator read this as work at risk and asked about it. Both commits were already on the remote; nothing was unpushed. The refusal itself was correct (BUG-069's fail-closed guard) — only the stated reason was wrong.

## Repro
On the real clone (2026-09-18): local `HEAD` `e973bac`, remote `main` `09586ff`, clone's `refs/remotes/origin/main` stale at `71f749e`. `git rev-parse @{upstream}` fails: "upstream branch 'refs/heads/main' not stored as a remote-tracking branch". After `git fetch <url> main`, `git merge-base --is-ancestor HEAD FETCH_HEAD` succeeds and `git log --oneline HEAD --not FETCH_HEAD` is empty — every local commit is on the remote.

## Root cause
The guard counts with `git rev-list --count --all --not --remotes`, measured against remote-tracking refs. `report`'s `commitAndPush` pushes with `-u <url>` (a URL, not a remote name), which sets `branch.main.remote` to that URL but never updates `refs/remotes/origin/*`. That tracking basis is frozen at clone time, so every later commit counts as "not held by any remote-tracking ref" forever. The count itself is a fine conservative trigger; the message converted it into a conclusion the measurement cannot support — "N commits no tracking ref holds" was reported as "N unpushed commits".

## Fix
Commit `c7a6116`: the message now names what was measured, notes the commits may already be pushed and why the refs go stale, and gives the settling commands (`git -C <dir> fetch && git -C <dir> log --oneline HEAD --not FETCH_HEAD`). The guard's keep/remove decision is unchanged. `TestRemoveKeptClone` gained a subtest reproducing the production shape (branch tracking a URL, a commit pushed, no tracking ref moved) asserting the clone is kept without the "unpushed" claim; three subtests fail against the old message. `.agents/docs/deploy/local.md` § Report back documents it.

## Prevention
A message that reports a count names what the count was measured against whenever that basis can be stale or indirect, and never restates it as a stronger conclusion than the measurement supports — "N commits no remote-tracking ref holds" is what `--not --remotes` establishes; "N unpushed commits" is not. When the difference affects a decision the reader will make (here: delete a directory or not), the message also gives the command that settles it.

## Related
Same family as BUG-060 (a claim about push redirection had to name the mechanisms actually checked) and BUG-032 (a coverage claim across several files had to name what was left uncovered). BUG-069 is NOT that family: its defect was a guard that failed open, and its fail-closed behaviour is unchanged here — this record corrects an earlier draft that grouped the three together.

`context/AGENTS.md` § Hard rules already carries the rule for claims the AGENT makes ("A claim that something works states HOW it was verified and on which path… Never generalise a single observation into a claim about 'both paths' or 'every case'"). What this bug adds is that the same discipline applies to text the SOFTWARE prints: a message is a claim the operator acts on, and this one sent them looking for work that was never lost. Extending that existing rule to cover emitted messages is the promotion candidate — for a retro to decide, not this record.
