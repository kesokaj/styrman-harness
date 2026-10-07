# usage: python3 -I metrics.py <result.json> <claude config dir> <clone>
# Cost and time from the CLI's result; turns, tokens and subagents from the transcripts under <config>/projects
# (claude -p's own num_turns and duration cover only the last segment of a long run).
import collections, glob, json, os, subprocess, sys
def rows(p):
    for l in open(p):
        try: yield json.loads(l)
        except ValueError: continue
res, cfg, app = sys.argv[1], sys.argv[2], sys.argv[3]
r = json.load(open(res))
def scan(p):
    u = collections.Counter(); seen = set(); peak = 0
    for l in open(p):
        try: d = json.loads(l)
        except ValueError: continue
        if d.get('type') != 'assistant': continue
        m = d['message']
        if m.get('id') in seen: continue
        seen.add(m['id']); x = m.get('usage') or {}
        for k, v in x.items():
            if isinstance(v, int): u[k] += v
        peak = max(peak, x.get('input_tokens', 0) + x.get('cache_creation_input_tokens', 0) + x.get('cache_read_input_tokens', 0))
    return len(seen), u, peak
mt, mpeak, sends = 0, 0, 0
for p in glob.glob(f'{cfg}/projects/*/*.jsonl'):
    t, _, pk = scan(p); mt += t; mpeak = max(mpeak, pk)
    sends += sum(1 for d in rows(p) if d.get('type') == 'assistant' for c in (d['message'].get('content') or []) if c.get('type') == 'tool_use' and c.get('name') == 'SendMessage')
agents = []
for p in sorted(glob.glob(f'{cfg}/projects/*/*/subagents/*.jsonl')):
    meta = p[:-6] + '.meta.json'
    m = json.load(open(meta)) if os.path.exists(meta) else {}
    t, u, _ = scan(p)
    ts = [d.get('timestamp') for d in rows(p) if d.get('timestamp')]
    calls = sum(1 for d in rows(p) if d.get('type') == 'assistant' for c in (d['message'].get('content') or []) if c.get('type') == 'tool_use')
    agents.append(f"{m.get('agentType', '?')}: {m.get('description', '')} | {ts[0][11:19]}-{ts[-1][11:19]} | {t} turns, {calls} tool calls")
# git runs inside the agent-written repository: a minimal environment, never the operator's tokens.
# HOME is a path that does not exist, and core.fsmonitor/hooks are off: the repository's own .git/config is agent-written.
GITENV = {'PATH': os.environ['PATH'], 'HOME': '/nonexistent', 'GIT_CONFIG_NOSYSTEM': '1', 'GIT_CONFIG_GLOBAL': '/dev/null'}
def sh(cmd): return subprocess.run(cmd.replace('git ', 'git -c core.fsmonitor=false -c core.hooksPath=/dev/null '), shell=True, cwd=app, capture_output=True, text=True, env=GITENV).stdout.strip()
print(json.dumps(dict(
  cost_usd=round(r.get('total_cost_usd', 0), 2), cost_by_model={k: round(v.get('costUSD', 0), 2) for k, v in (r.get('modelUsage') or {}).items()},
  wall_min=round(int(open(os.path.join(os.path.dirname(res), 'wall_seconds')).read()) / 60, 1),
  main_turns=mt, main_peak_ctx=mpeak, sendmessage_resumes=sends, subagents=len(agents), agents=agents,
  commits=sh('git rev-list --count HEAD 2>/dev/null || echo 0'), files=sh('git ls-files -co --exclude-standard | wc -l'),
  bug_records=sh('ls .agents/docs/bugs/BUG-*.md 2>/dev/null | wc -l')), indent=1))
