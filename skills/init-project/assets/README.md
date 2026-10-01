# <project>

<one line: what it does and for whom — the project overview the user gave /init-project>

## Quick start

```bash
<build command>      # e.g. go build ./...
<run command>        # e.g. go run ./cmd/<app>
<test command>       # e.g. go vet ./... && go test -race ./...
```

The same commands are the contract in `AGENTS.md` § Setup commands; `.agents/docs/testing/guide.md` has the tiers behind them.

## Architecture

<two or three sentences: the entry point, what it talks to, where state lives>

```
cmd/<app>/        # the binary: flags, wiring, listener
internal/         # the packages, one responsibility each
web/              # the TanStack Start frontend, when there is one (built into the binary)
.agents/docs/     # the trail: plans, changes, decisions, specs, testing, deploy guides, bugs
```

## Configuration

Flags and environment variables, with their defaults; a runtime data path (a database file, an upload directory) defaults to a location OUTSIDE the repository, never beside the source. A secret-bearing setting's row says `Secret Manager: <name>` and has no default: a credential is never a code default.

| Setting | Flag / env | Default |
|---|---|---|
| `<PORT>` | `PORT` | `8080` |

## Deploy

How it runs here, on a remote box and on GCP is `.agents/docs/deploy/{local,remote,gcp}.md`; a release goes through `/preflight`.

## Where things are

Decisions: `.agents/docs/architecture/DECISIONS.md`. What changed: `.agents/docs/changes/`. What went wrong and the rule that came out of it: `.agents/docs/bugs/INDEX.md`. The next session starts from the newest `.agents/docs/handoffs/` entry.
