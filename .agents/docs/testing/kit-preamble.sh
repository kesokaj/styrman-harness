#!/usr/bin/env bash
# kit-preamble.sh — sourced first by every behavioural or sandbox run kit built
# from guide.md § 4 (`. .agents/docs/testing/kit-preamble.sh` at the top of the
# kit script in scratch). A run proves its setup before its output counts
# (BUG-087, 090, 091, 106, 107, 115, 116, 123): this file is that proof, in
# one place, so a kit cannot skip a step by forgetting it. Sourced, never
# executed, so its mode is 0644 like every other file under the docs root.
#
# Inputs (environment):
#   KIT_CLONE          the fixture clone the CLI runs in (required)
#   KIT_HOME           the sandbox home the run uses (required; never the real home)
#   KIT_ASSERT_GREP    newline-separated `file:pattern` lines, each relative to
#                      KIT_CLONE; every pattern must match its file (optional)
#   KIT_ASSERT_COMMIT  one commit SHA (7 to 40 hex digits, never a ref name)
#                      that must be in the clone's history (optional)
#   KIT_BIN            the sh-agent build to wrap (default: a fresh `go build` of the
#                      harness checkout this file lives in, run from any cwd)
# Effects: asserts the fixture, refuses a clone that sits inside another git
# checkout (BUG-106), lists ~/.config/sh-agent before the run and traps EXIT to
# list it after (BUG-107, BUG-123), and puts a wrapper first on PATH that pins
# `--home "$KIT_HOME"` on every sh-agent call the agent makes (BUG-107).
set -euo pipefail
: "${KIT_CLONE:?KIT_CLONE is required}"
: "${KIT_HOME:?KIT_HOME is required}"
case "$KIT_HOME" in "$HOME"|"$HOME/") echo "kit: KIT_HOME is the real home; refusing" >&2; exit 2;; esac
[ -d "$KIT_CLONE/.git" ] || { echo "kit: $KIT_CLONE is not a git clone" >&2; exit 2; }
# The clone's parent is no checkout: an ancestor AGENTS.md would load (BUG-106).
if git -C "$(dirname "$KIT_CLONE")" rev-parse --git-dir >/dev/null 2>&1; then
  echo "kit: $(dirname "$KIT_CLONE") is inside a git checkout; the clone must live outside every checkout (BUG-106)" >&2; exit 2
fi
# Fixture facts the dispatch states, asserted before the CLI starts.
if [ -n "${KIT_ASSERT_GREP:-}" ]; then
  while IFS= read -r line; do
    [ -z "$line" ] && continue
    f="${line%%:*}"; p="${line#*:}"
    grep -q -- "$p" "$KIT_CLONE/$f" || { echo "kit: fixture assert failed: $f lacks '$p'" >&2; exit 2; }
  done <<< "$KIT_ASSERT_GREP"
fi
if [ -n "${KIT_ASSERT_COMMIT:-}" ]; then
  # Asked of git directly, never `git log | grep -q`: under pipefail grep's
  # early exit kills git log on a long history and the status is 141, which
  # refused a commit that was there (BUG-157).
  # A SHA, nothing else: `HEAD`, a branch or a tag name would be an ancestor
  # of HEAD whatever the fixture holds, and the assertion would prove nothing.
  [[ "$KIT_ASSERT_COMMIT" =~ ^[0-9a-f]{7,40}$ ]] || { echo "kit: KIT_ASSERT_COMMIT must be a commit SHA of 7 to 40 hex digits, got '$KIT_ASSERT_COMMIT'" >&2; exit 2; }
  git -C "$KIT_CLONE" merge-base --is-ancestor --end-of-options "${KIT_ASSERT_COMMIT}^{commit}" HEAD 2>/dev/null || { echo "kit: commit $KIT_ASSERT_COMMIT not in the clone's log" >&2; exit 2; }
fi
# The real ~/.config/sh-agent before and after: a run that writes there fails
# the kit. The fingerprint is `ls -lR` (names, sizes, mtimes), so an appended
# line to an existing file changes it; a names-only listing did not (docs-lens
# review 2026-10-01).
# The real home is remembered here: a kit overrides HOME for the CLI it runs,
# and the exit check must still look at the operator's own directory.
kit_real_config="$HOME/.config/sh-agent"
kit_list_real_config() { ls -la "$kit_real_config" 2>/dev/null | sed 's/^/  /'; }
kit_fingerprint() { ls -lR --time-style=full-iso "$kit_real_config" 2>/dev/null | md5sum | cut -c1-12; }
echo "kit: real $kit_real_config before the run:"; kit_list_real_config
kit_before=$(kit_fingerprint)
kit_on_exit() {
  local rc=$?
  echo "kit: real $kit_real_config after the run:"; kit_list_real_config
  if [ "$kit_before" != "$(kit_fingerprint)" ]; then
    echo "kit: FAIL the real ~/.config/sh-agent changed during the run (BUG-107, BUG-123); the run does not count" >&2
    exit 1
  fi
  exit "$rc"
}
trap kit_on_exit EXIT
# A wrapper first on PATH pins --home for every sh-agent call the agent makes.
mkdir -p "$KIT_HOME/bin"
# Built fresh every run (the build cache makes it cheap): a stale binary from
# an earlier kit measured the wrong code (docs-lens review 2026-10-01). A kit
# that must use a given binary sets KIT_BIN.
if [ -z "${KIT_BIN:-}" ]; then
  KIT_BIN="$KIT_HOME/bin/sh-agent.real"
  kit_repo=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)
  go build -C "$kit_repo" -o "$KIT_BIN" ./cmd/sh-agent || { echo "kit: go build of sh-agent in $kit_repo failed; set KIT_BIN to a build" >&2; exit 2; }
fi
printf '#!/usr/bin/env bash\nexec %q "$1" --home %q "${@:2}"\n' "$KIT_BIN" "$KIT_HOME" > "$KIT_HOME/bin/sh-agent"
chmod 0755 "$KIT_HOME/bin/sh-agent"
export PATH="$KIT_HOME/bin:$PATH"
hash -r
echo "kit: preamble ok — clone $KIT_CLONE, home $KIT_HOME, sh-agent wrapper first on PATH"
