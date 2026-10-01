#!/usr/bin/env bash
# refresh-catalog.sh — vendor the Google skills catalog at two pins.
#
#   bash skills/google-catalog/scripts/refresh-catalog.sh [--dest <dir>] [--work <dir>]
#        [--exclude <entry>]... <agents-cli pin> <google-skills pin>
#
# Run from the harness checkout, never from a session in another project. It
# clones both repositories in full (so any pin resolves) into the work
# directory, builds the catalog there, and only then replaces the entries under
# the destination. Vendored files are copied byte for byte; the one local
# change is the file mode (no executable bit: the embed has none).
#
# It refuses, naming the entry and leaving the destination untouched: a pin
# that does not resolve; an entry without SKILL.md frontmatter (name and
# description); an entry name in two source sets; a symlink; a go.mod, go.sum
# or *.go file (the Go tool and the embed would both mishandle the directory).
# The operator then passes --exclude <entry> and records it under the
# manifest's "## Excluded". What it cannot see cheaply — a secret-shaped file
# name, a path git ignores — is caught by the suite it tells you to run next
# (TestCatalogInstallsClean, TestEmbedHoldsOnlyTrackedFiles).
#
# INDEX.md and CHECKSUMS are written by the Go tests in write mode, so the
# index is derived by the installer's own frontmatter parser.
#
# Needs git, go and (for the default sources) the network. Two environment
# variables replace the source repositories, for the test's local fixtures:
# CATALOG_AGENTS_CLI_REPO and CATALOG_GOOGLE_SKILLS_REPO. Relative --dest and
# --work are resolved against the caller's directory.
set -euo pipefail

usage() { sed -n '2,6p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//' >&2; }
fail() { echo "refresh-catalog: $*" >&2; exit 1; }

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
repo=$(cd "$script_dir/../../.." && pwd)
dest="$repo/skills/google-catalog/catalog"
work="$repo/.agents/scratch/catalog-refresh"
excludes=""
while [ $# -gt 0 ]; do
  case "$1" in
    --dest) dest=$2; shift 2 ;;
    --work) work=$2; shift 2 ;;
    --exclude) excludes="$excludes $2"; shift 2 ;;
    --) shift; break ;;
    -*) usage; exit 2 ;;
    *) break ;;
  esac
done
[ $# -eq 2 ] || { usage; exit 2; }
# Absolute before anything runs: the Go writer below runs with another working
# directory, where a relative path would name nothing (BUG-146).
abs() { case "$1" in /*) printf '%s\n' "$1" ;; *) printf '%s/%s\n' "$PWD" "$1" ;; esac; }
dest=$(abs "$dest")
work=$(abs "$work")
pin_a=$1
pin_g=$2
url_a=${CATALOG_AGENTS_CLI_REPO:-https://github.com/google/agents-cli}
url_g=${CATALOG_GOOGLE_SKILLS_REPO:-https://github.com/google/skills}

# The four source sets: label, clone, path below it.
sets="agents-cli|agents-cli|skills
agents-cli-langchain|agents-cli|extensions/langchain/template/.agents/skills
google-skills-cloud|google-skills|skills/cloud
google-skills-developers|google-skills|skills/developers"

mkdir -p "$work"
fetch() { # <clone name> <url> <pin>: a full clone, reused; checked out at the pin
  local dir="$work/$1"
  if [ -d "$dir/.git" ]; then
    git -C "$dir" remote set-url origin "$2"
    git -C "$dir" fetch -q origin '+refs/heads/*:refs/remotes/origin/*' || fail "cannot fetch $2"
  else
    git clone -q "$2" "$dir" || fail "cannot clone $2"
  fi
  git -C "$dir" -c advice.detachedHead=false checkout -q --detach "$3" 2>/dev/null || fail "pin $3 does not resolve in $2"
}
fetch agents-cli "$url_a" "$pin_a"
fetch google-skills "$url_g" "$pin_g"
short_a=$(git -C "$work/agents-cli" rev-parse --short=7 HEAD)
short_g=$(git -C "$work/google-skills" rev-parse --short=7 HEAD)

build="$work/build"
rm -rf "$build"
mkdir -p "$build"
origins="$work/origins.txt"
: > "$origins"
excluded() { case " $excludes " in *" $1 "*) return 0 ;; esac; return 1; }

while IFS='|' read -r label clone sub; do
  root="$work/$clone/$sub"
  [ -d "$root" ] || fail "source set $label: no $sub at the pin"
  for d in "$root"/*/; do
    [ -d "$d" ] || continue
    name=$(basename "$d")
    excluded "$name" && continue
    [ -f "$d/SKILL.md" ] || fail "entry $name ($label) has no SKILL.md (--exclude $name)"
    prev=$(awk -v n="$name" '$1 == n {print $2}' "$origins")
    [ -z "$prev" ] || fail "entry $name is in two source sets, $prev and $label (--exclude $name drops both)"
    bad=$(find "${d%/}" -type l -print | head -1)
    [ -z "$bad" ] || fail "entry $name ($label) holds a symlink: ${bad#"$root"/} (--exclude $name)"
    bad=$(find "${d%/}" -type f \( -name go.mod -o -name go.sum -o -name '*.go' \) -print | head -1)
    [ -z "$bad" ] || fail "entry $name ($label) holds a Go file the module would pick up: ${bad#"$root"/} (--exclude $name)"
    mkdir -p "$build/$name"
    cp -R "${d%/}/." "$build/$name/"
    printf '%s %s\n' "$name" "$label" >> "$origins"
  done
done <<SETS
$sets
SETS
[ -s "$origins" ] || fail "no entry found in the four source sets"

find "$build" -type f -exec chmod 644 {} +
find "$build" -type d -exec chmod 755 {} +
[ -f "$work/agents-cli/LICENSE" ] && [ -f "$work/google-skills/LICENSE" ] || fail "a repository has no LICENSE file at the pin"
cp "$work/agents-cli/LICENSE" "$build/LICENSE-agents-cli"
cp "$work/google-skills/LICENSE" "$build/LICENSE-google-skills"
chmod 644 "$build/LICENSE-agents-cli" "$build/LICENSE-google-skills"

# INDEX.md and CHECKSUMS, by the Go tests in write mode over the build; an
# entry without frontmatter name and description fails here, before the swap.
if ! (cd "$repo" && CATALOG_DIR="$build" CATALOG_WRITE=1 go test ./internal/harness \
    -run '^TestCatalog(MatchesChecksums|IndexMatchesEntries|FilesAreCarryable)$' -count=1 > "$work/writer.log" 2>&1); then
  grep -E 'catalog_test|catalog entry|FAIL' "$work/writer.log" | head -20 >&2 || true
  fail "the catalog built from these pins is refused (lines above; --exclude the entry named)"
fi

# The manifest: rows rewritten, the fetched date kept for a row whose pin did
# not change, the hand-written sections kept byte for byte.
old="$dest/MANIFEST.md"
today=$(date +%F)
count() { awk -v l="$1" '$2 == l' "$origins" | wc -l | tr -d ' '; }
fetched() { # <path cell> <pin>: the old row's date when its pin is the same
  if [ -f "$old" ]; then
    d=$(awk -F'|' -v p=" $1 " -v pin=" $2 " '$3 == p && $4 == pin {gsub(/ /, "", $6); print $6}' "$old" | head -1)
    [ -n "$d" ] && { echo "$d"; return; }
  fi
  echo "$today"
}
row() { printf '| %s | %s | %s | Apache-2.0 | %s | %s |\n' "$1" "$2" "$3" "$(fetched "$2" "$3")" "$(count "$4")"; }
table="| Repository | Path | Pin | License | Fetched | Entries |
|---|---|---|---|---|---|
$(row "$url_a" "skills" "$short_a" agents-cli)
$(row "$url_a" "extensions/langchain/template/.agents/skills" "$short_a" agents-cli-langchain)
$(row "$url_g" "skills/cloud" "$short_g" google-skills-cloud)
$(row "$url_g" "skills/developers" "$short_g" google-skills-developers)"
tail_sections="## Excluded

Upstream entries the refresh refused, with the reason; passed as \`--exclude\` on every refresh.

- none

## Notes

Per-entry notes from the security review; the vendored text itself is never edited.

- none"
manifest_date=$today
if [ -f "$old" ]; then
  if grep -q '^## Excluded' "$old"; then tail_sections=$(sed -n '/^## Excluded/,$p' "$old"); fi
  old_table=$(sed -n '/^| Repository | Path |/,/^$/p' "$old" | sed '/^$/d')
  if [ "$old_table" = "$table" ]; then manifest_date=$(sed -n 's/^date: //p' "$old" | head -1); fi
fi
cat > "$build/MANIFEST.md" <<MANIFEST
---
owner: main
date: $manifest_date
---
# Google skills catalog — provenance

Skills published by Google, vendored verbatim at the pins below by \`scripts/refresh-catalog.sh\` and reached through the router skill one directory up. \`CHECKSUMS\` pins every vendored file; \`INDEX.md\` is derived from the entries' frontmatter.

$table

Local changes: file modes only (no executable bit).

$tail_sections
MANIFEST

# Excluded entries the operator did not pass again.
sed -n '/^## Excluded/,/^## Notes/p' "$build/MANIFEST.md" | sed -n 's/^- `\([^`]*\)`.*/\1/p' | while read -r name; do
  excluded "$name" || echo "note: $name is listed under ## Excluded but was not passed as --exclude"
done

# The swap: entries replaced, entries upstream dropped removed.
mkdir -p "$dest"
before=$(cd "$dest" && find . -mindepth 1 -maxdepth 1 -type d | sed 's|^\./||' | sort)
after=$(awk '{print $1}' "$origins" | sort)
find "$dest" -mindepth 1 -maxdepth 1 -type d -exec rm -rf {} +
cp -R "$build/." "$dest/"

while IFS='|' read -r label clone sub; do
  echo "$label: $(count "$label") entries"
done <<SETS
$sets
SETS
added=$(comm -13 <(printf '%s\n' "$before") <(printf '%s\n' "$after") | tr '\n' ' ')
removed=$(comm -23 <(printf '%s\n' "$before") <(printf '%s\n' "$after") | tr '\n' ' ')
[ "$before" = "" ] && added="(first import: $(printf '%s\n' "$after" | wc -l | tr -d ' ') entries)"
added=${added% }; echo "added: ${added:-none}"
removed=${removed% }; echo "removed: ${removed:-none}"
echo "next: git add skills/google-catalog && go test ./..."
