#!/usr/bin/env bash
# run.sh — seeded-weakness check of the `security` agent (plan 2026-10-07-lean-process-3, ADR-054a's open
# question: does effort high still find what max found?). usage: bash .agents/docs/testing/seed-security/run.sh <label> <effort>
# One commit plants six known flaws in a small Go service; Claude Code (headless, the same isolation as
# testing/ab-kit/run.sh: kit-preamble.sh, env -i allowlist, broker off, credential linked and removed) is asked
# only to dispatch `security` over that commit; the installed security.md gets `effort: <effort>`. The score is
# how many flaws the agent's artifact names, by keyword per flaw (printed per flaw, so a miss can be read; the
# keywords are loose, so the artifact is read by hand too). The fixture is a remote-code-execution service by
# design: the kit never starts it, and it binds 127.0.0.1:0 (loopback, a random port) if anything does. A server
# the agent itself starts survives a timeout: check `ss -ltnp` after a run that timed out. Homes live under $SEED_OUT/homes, outside every checkout.
set -euo pipefail
# Inherited GIT_* variables (GIT_DIR, GIT_WORK_TREE, …) would point the fixture's git at another repository.
for v in $(compgen -e | grep '^GIT_' || true); do unset "$v"; done
label=${1:?usage: run.sh <label> <effort>}; effort=${2:?effort: low|medium|high|xhigh|max}
case "$effort" in low|medium|high|xhigh|max) ;; *) echo "run.sh: effort must be low|medium|high|xhigh|max" >&2; exit 2 ;; esac
KIT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd); REPO=$(cd "$KIT/../../../.." && pwd)
REAL_HOME=$HOME; OUT=${SEED_OUT:-$REAL_HOME/sandbox/seed-security}; RUN=$OUT/runs/$label
export KIT_CLONE=$RUN/app KIT_HOME=$OUT/homes/$label
for p in "$RUN" "$KIT_HOME"; do [ ! -e "$p" ] || { echo "run.sh: $p exists; move it away to re-run" >&2; exit 2; }; done
mkdir -p "$KIT_CLONE/cmd/svc" "$KIT_HOME/.claude"
cd "$KIT_CLONE"; git init -q
[ "$(git rev-parse --show-toplevel)" = "$KIT_CLONE" ] || { echo "run.sh: $KIT_CLONE is not its own repository" >&2; exit 2; }
gi() { git -c user.name=fixture -c user.email=f@localhost "$@"; }
printf 'module example.com/svc\n\ngo 1.24\n' > go.mod
cat > cmd/svc/main.go <<'GO'
// Command svc is a small internal file service.
package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	routes(mux)
	if err := http.ListenAndServe("127.0.0.1:0", mux); err != nil {
		slog.Error("stopped", "err", err)
		os.Exit(1)
	}
}

func routes(mux *http.ServeMux) {}
GO
gi add -A; gi commit -qm "svc skeleton"
cat > cmd/svc/main.go <<'GO'
// Command svc is a small internal file service.
package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const dataDir = "/srv/files"

var apiToken = os.Getenv("SVC_TOKEN")

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	routes(mux)
	if err := http.ListenAndServe("127.0.0.1:0", mux); err != nil {
		slog.Error("stopped", "err", err)
		os.Exit(1)
	}
}

func routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /files", authed(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(dataDir, r.URL.Query().Get("name")))
	}))
	mux.HandleFunc("POST /thumbnail", authed(func(w http.ResponseWriter, r *http.Request) {
		out, err := exec.Command("sh", "-c", "convert "+r.FormValue("file")+" -resize 64x64 /tmp/thumb.png").CombinedOutput()
		if err != nil {
			http.Error(w, string(out), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("done"))
	}))
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("user") == "admin" && r.FormValue("password") == "Winter2026!" {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: fmt.Sprint(rand.Int63())})
			w.Write([]byte("welcome"))
			return
		}
		http.Error(w, "no", http.StatusUnauthorized)
	})
	mux.HandleFunc("POST /admin/purge", func(w http.ResponseWriter, r *http.Request) {
		os.RemoveAll(dataDir)
		w.Write([]byte("purged"))
	})
}

func authed(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token") != apiToken {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		h(w, r)
	}
}
GO
gi commit -qam "svc: files, thumbnails, login and admin purge"
. "$REPO/.agents/docs/testing/kit-preamble.sh"
go run -C "$REPO" ./cmd/sh-agent install --home "$KIT_HOME" --vendor claude > "$RUN/install.log" 2>&1
sed -i "s/^effort: .*/effort: $effort/" "$KIT_HOME/.claude/agents/security.md"
grep -q "^effort: $effort\$" "$KIT_HOME/.claude/agents/security.md"
printf '{"attribution":{"commit":"","pr":""}}\n' > "$KIT_HOME/.claude/settings.json"
real_cred=$REAL_HOME/.claude/.credentials.json
[ -f "$real_cred" ] && [ ! -L "$real_cred" ] || { echo "run.sh: $real_cred is not a regular file" >&2; exit 2; }
ln -s "$real_cred" "$KIT_HOME/.claude/.credentials.json"
t0=$(date +%s); rc=0
( cd "$KIT_CLONE" && env -i PATH="$PATH" TERM="${TERM:-dumb}" LANG="${LANG:-C.UTF-8}" USER="${USER:-$(id -un)}" LOGNAME="${LOGNAME:-$(id -un)}" SHELL=/bin/bash \
    STYRMAN_SKIP_CRED_PULL=1 HOME="$KIT_HOME" CLAUDE_CONFIG_DIR="$KIT_HOME/.claude" \
    timeout --foreground 40m claude -p "Dispatch the security agent over the last commit (the diff is git diff HEAD~1..HEAD), date $(date +%F), slug seed, a budget of 20 tool calls, artifact first. Do not review it yourself and do not ask me anything. Finish with the artifact path." \
    --model claude-opus-5-5 --strict-mcp-config --output-format json --permission-mode auto --disallowedTools WebFetch WebSearch \
    < /dev/null > "$RUN/result.json" 2> "$RUN/stderr.log" ) || rc=$?
[ -L "$KIT_HOME/.claude/.credentials.json" ] || echo "run.sh: WARNING the credential link was replaced" >&2
rm -f "$KIT_HOME/.claude/.credentials.json"
secs=$(( $(date +%s) - t0 ))
art=$(ls "$KIT_CLONE"/.agents/docs/reviews/*security* "$KIT_CLONE"/docs/reviews/*security* 2>/dev/null | head -1 || true)
[ -n "$art" ] || { echo "run.sh: $label: no security artifact (claude exit $rc, ${secs}s)"; exit 1; }
found=0
for f in "injection:command injection|shell injection|sh -c|os/exec|exec.Command" \
         "traversal:path traversal|directory traversal|\.\./|filepath.Join" \
         "password:hard-?coded (admin )?password|Winter2026|hardcoded credential|hard-coded credential" \
         "purge-authz:/admin/purge|purge.{0,80}(auth|unauthenticated)|unauthenticated.{0,80}purge" \
         "timing:constant[- ]time|timing|subtle.ConstantTimeCompare" \
         "rand:math/rand|predictable session|insecure random|crypto/rand"; do
  name=${f%%:*}; pat=${f#*:}
  if grep -qiE -- "$pat" "$art"; then found=$((found+1)); echo "  found $name"; else echo "  MISSED $name"; fi
done
echo "run.sh: $label (effort $effort): found $found/6 in ${secs}s, cost \$$(python3 -I -c 'import json,sys;print(round(json.load(open(sys.argv[1])).get("total_cost_usd",0),2))' "$RUN/result.json"); artifact $art"
