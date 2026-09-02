#!/bin/bash
# Resume Protocol context injector — the deterministic half of the Resume Protocol.
# Shared by BOTH harnesses:
#   - Claude Code:      SessionStart hook (JSON on stdin; stdout is injected into
#                       model context — every line printed here reaches the agent)
#   - Antigravity CLI:  `resume-context.sh --agy` as a PreInvocation hook (JSON on
#                       stdin with invocationNum / workspacePaths): on the first
#                       invocation of a conversation, wraps the same lines as
#                       {"injectSteps":[{"ephemeralMessage":...}]}; otherwise {}.
#                       Also usable as a plain skill step (no stdin).
# Prints AT MOST 6 short lines orienting the agent in an existing harness project;
# prints nothing outside one. Never blocks the session — always exits 0. No network.

mode=${1:-}

# Repo-local config must not execute during a gate (core.fsmonitor runs a
# program on some git commands); the payload chooses the directory, so pin it off.
g() { git -c core.fsmonitor=false "$@"; }

# Hook stdin is a closed pipe in practice; a manual run with an open, silent
# stdin must not hang forever (2 s cap where `timeout` exists).
read_stdin() {
  if command -v timeout >/dev/null 2>&1; then timeout 2 cat 2>/dev/null || true
  else cat 2>/dev/null || true; fi
}

# Filenames become model-facing text: keep them to a safe character set.
safe() { tr -c 'A-Za-z0-9._/\n-' '_' | cut -c1-120; }

# --- resolve project dir: hook JSON .cwd (Claude) / .workspacePaths[0] (agy),
# else $PWD. Stdin may be absent (standalone/tty), empty, or garbage — tolerate all.
dir="$PWD"
invocation=1
if [ ! -t 0 ]; then
  input=$(read_stdin)
  if [ -n "$input" ] && command -v jq >/dev/null 2>&1; then
    hook_cwd=$(printf '%s' "$input" | jq -r '.cwd // .workspacePaths[0] // empty' 2>/dev/null)
    [ -n "$hook_cwd" ] && [ -d "$hook_cwd" ] && dir="$hook_cwd"
    invocation=$(printf '%s' "$input" | jq -r '.invocationNum // 1' 2>/dev/null)
  fi
fi

# Everything the agent gets to see for $dir. Prints nothing outside a git repo
# or a project that has not opted into the harness trail (no changes/ dir).
print_context() {
  local root docs f hf hb hf_epoch age adr n newest idx b newest3 change_date commit_date missing
  root=$(git -C "$dir" rev-parse --show-toplevel 2>/dev/null)
  [ -n "$root" ] || return 0
  if [ -d "$root/.agents/docs" ]; then docs="$root/.agents/docs"
  elif [ -d "$root/docs" ]; then docs="$root/docs"
  else docs="$root/.agents/docs"; fi
  [ -d "$docs/changes" ] || return 0

  # (1) The reading order — per AGENTS.md.
  echo "Resume Protocol: AGENTS.md -> REQUIREMENTS.md -> ROADMAP.md -> TODO.md (open items) -> newest handoff -> DECISIONS.md (index below; grep '^## ADR' for older) -> 3 newest changes/"

  # (2) Latest handoff — newest by FILENAME date. mtime lies after clone/
  # checkout/rsync; handoff filenames are date-prefixed by contract.
  hf=""
  for f in "$docs/handoffs"/[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-*.md; do
    [ -f "$f" ] || continue
    hf="$f"
  done
  if [ -n "$hf" ]; then
    hb=${hf##*/}
    hf_epoch=$(date -d "${hb:0:10}" +%s 2>/dev/null || date -j -f %Y-%m-%d "${hb:0:10}" +%s 2>/dev/null || echo 0)
    if [ "$hf_epoch" -gt 0 ]; then
      age=$(( ($(date +%s) - hf_epoch) / 86400 ))
      echo "Latest handoff: $(printf '%s' "${hf#"$root"/}" | safe) (${age}d old)"
    else
      echo "Latest handoff: $(printf '%s' "${hf#"$root"/}" | safe)"
    fi
  else
    echo "no handoffs recorded"
  fi

  # (2b) ADR index. DECISIONS.md is append-only and grows without bound, so the
  # Resume Protocol reads this index plus the newest three ADRs, never the
  # whole file.
  adr="$docs/architecture/DECISIONS.md"
  if [ -f "$adr" ]; then
    n=$(grep -c '^## ADR' "$adr" 2>/dev/null); n=${n:-0}
    if [ "$n" -gt 0 ]; then
      # ADR titles are repo content headed for model context: safe chars only.
      newest=$(grep '^## ADR' "$adr" | tail -n 3 | sed 's/^## //' | tr -c 'A-Za-z0-9._:()/\n -' '_' | cut -c1-80 | paste -sd '|' -)
      echo "ADRs: $n in ${adr#"$root"/} — newest: $newest"
    fi
  fi

  # (2c) Index of the three newest changes/ entries (date + slug from the
  # filename). The Resume Protocol opens a body only when this line does not
  # answer, so most sessions never read the entries at all.
  idx=""
  for f in "$docs/changes"/[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-*.md; do
    [ -e "$f" ] || continue
    b=${f##*/}
    idx="${b%.md}
$idx"
  done
  newest3=$(printf '%s' "$idx" | grep . | head -n 3 | safe | paste -sd '|' -)
  [ -n "$newest3" ] && echo "Changes (newest 3): $newest3"

  # (3) Newest changes/ entry date (from filename) vs last commit date.
  change_date=""
  for f in "$docs/changes"/[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-*.md; do
    [ -e "$f" ] || continue
    b=${f##*/}
    change_date=${b:0:10}
  done
  commit_date=$(g -C "$root" log -1 --format=%cd --date=short 2>/dev/null)
  if [ -n "$commit_date" ]; then
    if [ -z "$change_date" ] || [[ "$commit_date" > "$change_date" ]]; then
      echo "WARNING: trail behind code (last change entry ${change_date:-none}, last commit $commit_date)"
    fi
  fi

  # (4) Missing or empty root files the Resume Protocol expects.
  missing=""
  for f in REQUIREMENTS.md ROADMAP.md TODO.md; do
    [ -s "$root/$f" ] || missing="$missing $f"
  done
  [ -n "$missing" ] && echo "Missing or empty:$missing"
  return 0
}

# agy PreInvocation fires before every model call; orient on the first only.
if [ "$mode" = "--agy" ]; then
  command -v jq >/dev/null 2>&1 || { echo '{}'; exit 0; }
  [ "$invocation" = "1" ] || { echo '{}'; exit 0; }
  text=$(print_context)
  [ -n "$text" ] || { echo '{}'; exit 0; }
  jq -cn --arg t "$text" '{injectSteps:[{ephemeralMessage:$t}]}'
  exit 0
fi

print_context
exit 0
