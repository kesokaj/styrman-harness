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
# `--home "$KIT_HOME"` on every sh-agent call the agent makes (BUG-107). At exit
# it also fails the run when any credential is still in KIT_HOME (BUG-160): a
# kit deletes what it copied before it ends, and this check proves it did.
# The exit check is that EXIT trap and nothing else: a kit never sets an EXIT
# trap of its own after sourcing this file and never ends in `exec <cli>`,
# and a run counts only when its output holds "kit: real config unchanged" and
# "kit: no credentials left in KIT_HOME" (or a FAIL line).
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
  # Resolved first, and the full id must start with the value: git takes a
  # tag or branch whose name is hex digits (a date, `20261001`) before an
  # abbreviated id, and such a ref would pass as "the commit".
  kit_full=$(git -C "$KIT_CLONE" rev-parse --verify --quiet --end-of-options "${KIT_ASSERT_COMMIT}^{commit}" 2>/dev/null) || kit_full=""
  case "$kit_full" in
    "") echo "kit: commit $KIT_ASSERT_COMMIT not in the clone's log" >&2; exit 2 ;;
    "$KIT_ASSERT_COMMIT"*) ;;
    *) echo "kit: KIT_ASSERT_COMMIT must be a commit SHA: '$KIT_ASSERT_COMMIT' names a ref in the clone" >&2; exit 2 ;;
  esac
  git -C "$KIT_CLONE" merge-base --is-ancestor "$kit_full" HEAD 2>/dev/null || { echo "kit: commit $KIT_ASSERT_COMMIT not in the clone's log" >&2; exit 2; }
fi
# The real ~/.config/sh-agent before and after: a run that changes a name or
# a file's contents there fails the kit. The fingerprint covers names and
# contents, so an appended line to an existing file changes it; a names-only
# listing did not (docs-lens review 2026-10-01). It does not see: a rewrite
# with the same bytes, a `touch` or a `chmod`; a file made and removed again
# before the run ends; anything inside a directory it cannot list, or inside
# a subdirectory that is itself a symlink.
# The real home is remembered here: a kit overrides HOME for the CLI it runs,
# and the exit check must still look at the operator's own directory.
kit_real_config="$HOME/.config/sh-agent"
# Both tolerate a home that has no such directory yet: under pipefail a
# failing `ls` ended the preamble with status 2 and no message (BUG-158).
# The fingerprint is every path plus a checksum of every file, by tools macOS
# has too: `ls --time-style` and `md5sum` are GNU's, and without them the
# fingerprint was the same empty value before and after, proving nothing.
kit_list_real_config() { { ls -la "$kit_real_config/" 2>/dev/null || true; } | sed 's/^/  /'; }
kit_fingerprint() {
  # -H: the directory itself may be a symlink (a dotfiles setup), and find
  # would otherwise stop at it. Links inside are listed by name, never
  # walked (-L would follow one to `/` and read the whole machine). A link
  # that points at a regular file is read through, so that file's content
  # counts; any other link gives the path it points at, so re-pointing it
  # shows, and a link to a device or a fifo, which never ends, is not read.
  { find -H "$kit_real_config" -print 2>/dev/null || true
    find -H "$kit_real_config" -type f -exec cksum {} + 2>/dev/null || true
    find -H "$kit_real_config" -type l -exec sh -c 'for p; do if [ -f "$p" ]; then cksum "$p"; else printf "%s -> %s\n" "$p" "$(readlink "$p")"; fi; done' sh {} + 2>/dev/null || true; } | LC_ALL=C sort | cksum | cut -d' ' -f1
}
echo "kit: real $kit_real_config before the run:"; kit_list_real_config
kit_before=$(kit_fingerprint)
# Credentials left in the sandbox home (BUG-160): by NAME, the kinds a kit
# copies or a CLI logs in to (an SSH private key, Claude's credentials, gcloud's,
# Antigravity's and gh's tokens); and by CONTENT, a token or private key a CLI
# wrote somewhere else — Antigravity dumped GH_TOKEN into a terminal log. Paths
# only, never the matching text. The kit's own wrapper directory is skipped.
kit_credential_leftovers() {
  [ -d "$KIT_HOME" ] || return 0
  { find "$KIT_HOME" -path "$KIT_HOME/bin" -prune -o -type f \( \
      -name 'id_rsa' -o -name 'id_ecdsa' -o -name 'id_ed25519' -o -name 'id_dsa' -o \
      -name '.credentials.json' -o -name 'application_default_credentials.json' -o \
      -name 'credentials.db' -o -name 'access_tokens.db' -o -name '*oauth-token*' -o \
      -name 'hosts.yml' \) -print 2>/dev/null || true
    grep -rlIE --exclude-dir=bin -e 'gh[opsur]_[A-Za-z0-9]{30,}' -e 'github_pat_[A-Za-z0-9_]{30,}' \
      -e '-----BEGIN [A-Z ]*PRIVATE KEY-----' -e 'sk-ant-[A-Za-z0-9_-]{20,}' -e '"refresh_token"' \
      "$KIT_HOME" 2>/dev/null || true; } | LC_ALL=C sort -u
}
kit_on_exit() {
  local rc=$? failed=0 left
  echo "kit: real $kit_real_config after the run:"; kit_list_real_config
  if [ "$kit_before" != "$(kit_fingerprint)" ]; then
    echo "kit: FAIL the real ~/.config/sh-agent changed during the run (BUG-107, BUG-123); the run does not count" >&2
    failed=1
  else
    # Said in so many words: a run whose output holds neither this line nor
    # the FAIL line above did not reach the check.
    echo "kit: real config unchanged"
  fi
  left=$(kit_credential_leftovers)
  if [ -n "$left" ]; then
    echo "kit: FAIL credentials left in KIT_HOME (BUG-160); delete them and re-check before the run counts:" >&2
    printf '%s\n' "$left" | sed 's/^/  /' >&2
    failed=1
  else
    echo "kit: no credentials left in KIT_HOME"
  fi
  [ "$failed" = 0 ] || exit 1
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
