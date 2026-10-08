#!/usr/bin/env bash
# run.sh — the A/B kit for the harness's process cost (plans 2026-10-07-lean-process, -2).
# usage: bash .agents/docs/testing/ab-kit/run.sh control|harness|<name>-full
#   control      no harness at all
#   harness      the harness installed from this checkout; the prompt as-is (unscaffolded: the Task path)
#   <name>-full  the harness, plus prefix-full.md: scaffold and follow the Feature workflow and Definition of Done
# Then: python3 -I accept.py <clone> <work dir>; python3 -I metrics.py <result.json> <home>/.claude <clone>.
#
# Built on kit-preamble.sh (BUG-162: the first version skipped it, and harness runs registered their scratch
# apps in the real global search index). The clone lives under $AB_OUT, default ~/sandbox/ab, outside every
# checkout (BUG-106); the sandbox home is $AB_OUT/homes/<arm>, outside every checkout (a git run from ~ reaches no repository); the CLI runs with HOME and
# CLAUDE_CONFIG_DIR there, so every sh-agent call the agent makes is pinned to it (BUG-107).
# The credential is a SYMLINK to the real ~/.claude/.credentials.json, not the guide's copy: a token refresh
# then lands in the real file instead of leaving the real login with a stale refresh token. It is deleted by
# name before the preamble's exit check, which must print its three lines (guide § Behavioural runs).
# On a styrman box every login shell's credential broker (/etc/profile.d/zz-styrman-credentials.sh) writes the
# SSH key, gcloud ADC and the Antigravity token into $HOME: with HOME=KIT_HOME they landed in the sandbox (the
# preamble's exit check caught it, 2026-10-07; the source of BUG-160's leftovers). The CLI therefore runs with
# STYRMAN_CRED_BROKER unset and STYRMAN_SKIP_CRED_PULL=1, the broker's own off switches. The CLI starts from
# `env -i` with an allowlist (PATH, TERM, LANG, USER, LOGNAME, SHELL, HOME, CLAUDE_CONFIG_DIR, git identity): the
# operator's GH_TOKEN, the launching session's tokens and the real Go caches stay out of the agent's
# environment (security review 2026-10-07); the broker still answers any process that asks, and ssh reads the
# real ~/.ssh. It is NOT a confinement: same uid, sudo, the real home reachable by path (guide § 4's
# accepted design, operator 2026-10-01). `timeout --foreground` lets Ctrl-C reach the agent; a
# server the agent leaves running is killed by the preamble's exit check, which fails the run (BUG-163); a kit killed by SIGKILL never runs that trap, so check `ss -ltnp` after one. No web tools: the harness arms
# install vendored skills, and the operator's acceptance (guide § Behavioural runs) covers text the harness wrote.
# A killed run leaves the link: the preamble's exit check then FAILs on it by name; delete it by name.
# Costs real tokens: about $0.5 (control) to $8 (a -full arm) on Opus.
set -euo pipefail
# Inherited GIT_* variables (GIT_DIR, GIT_WORK_TREE, …) would point the fixture's git at another repository.
for v in $(compgen -e | grep '^GIT_' || true); do unset "$v"; done
arm=${1:?usage: run.sh control|harness|<name>-full}
case "$arm" in control|harness|*-full) ;; *) echo "run.sh: unknown arm '$arm'" >&2; exit 2;; esac
KIT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
REPO=$(cd "$KIT/../../../.." && pwd)
REAL_HOME=$HOME
OUT=${AB_OUT:-$REAL_HOME/sandbox/ab}
RUN=$OUT/runs/$arm
export KIT_CLONE=$RUN/app KIT_HOME=$OUT/homes/$arm
for p in "$RUN" "$KIT_HOME"; do [ ! -e "$p" ] || { echo "run.sh: $p exists; move it away to re-run" >&2; exit 2; }; done
mkdir -p "$KIT_CLONE" "$KIT_HOME/.claude"
git -C "$KIT_CLONE" init -q
[ "$(git -C "$KIT_CLONE" rev-parse --show-toplevel)" = "$KIT_CLONE" ] || { echo "run.sh: $KIT_CLONE is not its own repository" >&2; exit 2; }
. "$REPO/.agents/docs/testing/kit-preamble.sh"
if [ "$arm" != control ]; then
  go run -C "$REPO" ./cmd/sh-agent install --home "$KIT_HOME" --vendor claude > "$RUN/install.log" 2>&1
fi
printf '{"attribution":{"commit":"","pr":""}}\n' > "$KIT_HOME/.claude/settings.json"
real_cred=$REAL_HOME/.claude/.credentials.json
[ -f "$real_cred" ] && [ ! -L "$real_cred" ] || { echo "run.sh: $real_cred is not a regular file; refusing" >&2; exit 2; }
ln -s "$real_cred" "$KIT_HOME/.claude/.credentials.json"
prompt=$(cat "$KIT/prompt.md")
case "$arm" in *-full) prompt="$(cat "$KIT/prefix-full.md")"$'\n\n'"$prompt" ;; esac
t0=$(date +%s)
rc=0
( cd "$KIT_CLONE" && env -i PATH="$PATH" TERM="${TERM:-dumb}" LANG="${LANG:-C.UTF-8}" USER="${USER:-$(id -un)}" LOGNAME="${LOGNAME:-$(id -un)}" SHELL=/bin/bash \
    STYRMAN_SKIP_CRED_PULL=1 HOME="$KIT_HOME" CLAUDE_CONFIG_DIR="$KIT_HOME/.claude" \
    GIT_AUTHOR_NAME=ab-kit GIT_AUTHOR_EMAIL=ab-kit@localhost GIT_COMMITTER_NAME=ab-kit GIT_COMMITTER_EMAIL=ab-kit@localhost \
    timeout --foreground 90m claude -p "$prompt" --model claude-opus-5-5 --strict-mcp-config --output-format json --permission-mode auto --disallowedTools WebFetch WebSearch \
    < /dev/null > "$RUN/result.json" 2> "$RUN/stderr.log" ) || rc=$?
echo $(( $(date +%s) - t0 )) > "$RUN/wall_seconds"
# A refresh must have gone through the link into the real file (plan lean-process-2 D5).
[ -L "$KIT_HOME/.claude/.credentials.json" ] || echo "run.sh: WARNING the credential link was replaced by a file; the real login may be stale" >&2
rm -f "$KIT_HOME/.claude/.credentials.json"
[ -f "$real_cred" ] && [ ! -L "$real_cred" ] || echo "run.sh: WARNING $real_cred is no longer a regular file" >&2
[ "$rc" != 124 ] || echo "run.sh: claude timed out after 90m; the run is recorded as timed out, not re-run" >&2
echo "run.sh: $arm finished (claude exit $rc) in $(cat "$RUN/wall_seconds")s; clone $KIT_CLONE; transcripts $KIT_HOME/.claude/projects"
exit "$rc"
