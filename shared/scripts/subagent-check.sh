#!/bin/bash
# SubagentStart / SubagentStop hook (Claude Code) — the deterministic half of
# "nothing of record lives only in chat": when a DoD owner (reviewer, security,
# verifier, changelog, explorer) stops without having written its artifact,
# block the stop once so the owner writes it. Never breaks the session: always
# exits 0; blocking happens only via the JSON decision.
#
#   subagent-check.sh start   # SubagentStart: stamp a marker for this agent_id
#   subagent-check.sh stop    # SubagentStop:  require an artifact at least as
#                             # new as the marker (or written in the last hour
#                             # when no marker exists — a resumed agent has no
#                             # SubagentStart). When the agent's final message
#                             # names its file, that file is what must exist,
#                             # so a sibling's artifact cannot stand in for it.
#
# Activates only in projects whose docs root has a changes/ directory (the
# harness opt-in marker). Antigravity has no subagent events; Claude-only.

mode=${1:-stop}
input=""
# Hook stdin is a closed pipe in practice; a manual run with an open, silent
# stdin must not hang forever (2 s cap where `timeout` exists).
read_stdin() {
  if command -v timeout >/dev/null 2>&1; then timeout 2 cat 2>/dev/null || true
  else cat 2>/dev/null || true; fi
}
[ ! -t 0 ] && input=$(read_stdin)
command -v jq >/dev/null 2>&1 || exit 0

j() { printf '%s' "$input" | jq -r "$1 // empty" 2>/dev/null; }

# Portable mtime (epoch seconds): GNU stat first, then BSD stat.
mtime_of() {
  stat -c %Y "$1" 2>/dev/null || stat -f %m "$1" 2>/dev/null || echo 0
}

# agent_type is matched against a fixed list below; agent_id names a marker
# file, so keep it to a safe character set (jq output is untrusted input).
agent=$(j .agent_type)
id=$(j .agent_id | tr -c 'A-Za-z0-9._-' '_' | cut -c1-128)
dir=$(j .cwd)
[ -n "$dir" ] && [ -d "$dir" ] || dir="$PWD"
[ -n "$agent" ] || exit 0

# Each owner's dated artifact: <docs>/<class>/<today><infix>*.md.
case "$agent" in
  reviewer)  class=reviews; infix="-review-" ;;
  security)  class=reviews; infix="-security-" ;;
  verifier)  class=testing; infix="-verify-" ;;
  changelog) class=changes; infix="-" ;;
  explorer)  class=plans;   infix="-recon-" ;;
  *) exit 0 ;;  # documentation writes undated files; other agents own no dated artifact
esac

root=$(git -C "$dir" rev-parse --show-toplevel 2>/dev/null)
[ -n "$root" ] || exit 0
if [ -d "$root/.agents/docs" ]; then docs="$root/.agents/docs"
elif [ -d "$root/docs" ]; then docs="$root/docs"
else docs="$root/.agents/docs"; fi
[ -d "$docs/changes" ] || exit 0

marker_dir="${TMPDIR:-/tmp}/subagent-check-$(id -u)"
marker="$marker_dir/${id:-$agent}"
# The marker dir must be ours and not a symlink — a dir pre-created by another
# local user could forge or drop markers. Without a usable dir, run markerless.
marker_ok=1
if [ -L "$marker_dir" ]; then
  marker_ok=0
elif [ ! -d "$marker_dir" ]; then
  mkdir -m 700 "$marker_dir" 2>/dev/null || marker_ok=0
fi
[ "$marker_ok" = 1 ] && [ ! -O "$marker_dir" ] && marker_ok=0

if [ "$mode" = "start" ]; then
  [ "$marker_ok" = 1 ] && : >"$marker" 2>/dev/null
  exit 0
fi

# stop: never block twice for the same agent.
if [ "$(j .stop_hook_active)" = "true" ]; then
  rm -f "$marker" 2>/dev/null
  exit 0
fi

today=$(date +%F)
# Since when must the artifact have been written? The start marker when we
# have one; otherwise (resumed agent, no SubagentStart) the last hour.
since=$(( $(date +%s) - 3600 ))
[ "$marker_ok" = 1 ] && [ -e "$marker" ] && since=$(mtime_of "$marker")

# Attribution: when the agent's final message names a file of its class, that
# file (not any sibling's) must exist and be fresh. Otherwise scan the class.
named=$(j .last_assistant_message | grep -oE "${today}${infix}[A-Za-z0-9._-]*\.md" | sort -u)
found=0
if [ -n "$named" ]; then
  while IFS= read -r b; do
    [ -n "$b" ] || continue
    f="$docs/$class/$b"
    # >= not >: an artifact written in the marker's own second still counts.
    if [ -e "$f" ] && [ "$(mtime_of "$f")" -ge "$since" ]; then found=1; break; fi
  done <<<"$named"
else
  for f in "$docs/$class/$today$infix"*.md; do
    [ -e "$f" ] || continue
    if [ "$(mtime_of "$f")" -ge "$since" ]; then found=1; break; fi
  done
fi
if [ "$found" = 1 ]; then
  rm -f "$marker" 2>/dev/null
  exit 0
fi
# Blocking: keep the marker — the next stop carries stop_hook_active and clears it.

reason="Definition of Done: the $agent agent stopped without writing its artifact — expected a file matching $class/$today$infix*.md under ${docs#"$root"/}/ written during this run (when your final message names the file, that exact file must exist). Write it now (a refusal, a no-op, or an empty result is still a report with a verdict), name its path in your final message, then stop."
jq -cn --arg r "$reason" '{decision:"block", reason:$r}'
exit 0
