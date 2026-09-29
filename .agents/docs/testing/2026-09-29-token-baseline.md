---
owner: main
date: 2026-09-29
---
# Token baseline — before `sh-agent search`   (2026-09-29)

Plan: `plans/2026-09-29-search-index.md`, Phase 0. The same probe runs again after Phase 3 and over the Phase 5 sandbox transcripts, so the numbers compare like for like.

## Status: BASELINE RECORDED

## What was measured

Every Claude Code transcript of this project on the Linux box (`~/.claude/projects/-home-kesokaj-workspace-styrman-harness/`), main sessions and subagents separately, as of 2026-09-29 around 18:00. Each tool result is classified by the tool call that produced it. "Share of tool output" is characters of tool results, not tokens. Input-side tokens are the sum of the `usage` fields on every assistant message.

## Result

## main: 23 transcripts

input-side tokens 3,252,227,065 (cache read 99%), output tokens 10,279,439; Read calls whole-file 18, ranged 5

| kind | share of tool output | calls | mean chars |
|---|---|---|---|
| read | 33.8% | 398 | 3112 |
| bash-other | 28.9% | 1202 | 882 |
| search | 13.3% | 258 | 1889 |
| agent-result | 10.7% | 365 | 1075 |
| git | 6.7% | 392 | 625 |
| web | 2.6% | 23 | 4111 |
| other | 1.7% | 116 | 529 |
| edit | 1.2% | 255 | 179 |
| build-test | 1.0% | 76 | 491 |

## subagents: 365 transcripts

input-side tokens 599,901,213 (cache read 90%), output tokens 3,863,021; Read calls whole-file 589, ranged 248

| kind | share of tool output | calls | mean chars |
|---|---|---|---|
| read | 59.2% | 1590 | 9031 |
| git | 14.5% | 695 | 5059 |
| bash-other | 12.9% | 1239 | 2523 |
| search | 9.5% | 1277 | 1800 |
| build-test | 2.0% | 476 | 1021 |
| web | 1.2% | 114 | 2512 |
| edit | 0.6% | 705 | 200 |
| other | 0.1% | 367 | 60 |

## Reading it

- Most input-side tokens are the conversation re-read from cache on every turn. An index cannot shrink that.
- File reads are the largest share of tool output: 34% in main sessions and 59% in subagents. Subagents load a whole file in 589 of 837 Read calls. This is where snippet search can save tokens.
- Search output is 13% of main-session and 10% of subagent tool output.

## What was not measured

- Antigravity sessions. Their logs were not read for this baseline.
- Tokens per tool result. Characters are the proxy.
- Whether a read was needed. The probe counts volume, not waste.

## The probe (run with `python3 <file> <transcript dir>`; ADR-010: a procedure is documented, not committed as a script)

```python
#!/usr/bin/env python3
"""Classify tool-result volume and token usage in Claude Code transcripts.
Usage: token-probe.py <transcript dir>   (e.g. ~/.claude/projects/<project-slug>)"""
import json, glob, os, re, sys, collections

def kind(name, inp):
    if name in ('Grep', 'Glob'): return 'search'
    if name == 'Read': return 'read'
    if name in ('Edit', 'Write', 'MultiEdit', 'NotebookEdit'): return 'edit'
    if name in ('Agent', 'Task'): return 'agent-result'
    if name in ('WebFetch', 'WebSearch'): return 'web'
    if name == 'Bash':
        c = re.sub(r'^cd [^&;]+(&&|;)\s*', '', inp.get('command', '').strip())
        if re.match(r'sh-agent search|go run \./cmd/sh-agent search', c): return 'sh-agent-search'
        if re.match(r'(git\s+grep|grep|rg|find|ls|fd)\b', c): return 'search'
        if re.match(r'(cat|sed -n|head|tail|nl|wc)\b', c): return 'read'
        if re.match(r'git\b', c): return 'git'
        if re.match(r'(go (test|vet|build|run)|gofmt|govulncheck)', c): return 'build-test'
        return 'bash-other'
    return 'other'

def size(c):
    if isinstance(c, str): return len(c)
    if isinstance(c, list): return sum(len(x.get('text', '')) for x in c if isinstance(x, dict))
    return 0

def scan(files):
    chars, calls, usage = collections.Counter(), collections.Counter(), collections.Counter()
    whole = ranged = 0
    for f in files:
        pending = {}
        for line in open(f, errors='ignore'):
            try: m = json.loads(line)
            except ValueError: continue
            msg = m.get('message') or {}
            if m.get('type') == 'assistant':
                for k, v in (msg.get('usage') or {}).items():
                    if isinstance(v, int): usage[k] += v
            content = msg.get('content')
            if not isinstance(content, list): continue
            for p in content:
                if p.get('type') == 'tool_use':
                    inp = p.get('input') or {}
                    pending[p['id']] = kind(p.get('name'), inp)
                    if p.get('name') == 'Read':
                        if inp.get('offset') or inp.get('limit'): ranged += 1
                        else: whole += 1
                elif p.get('type') == 'tool_result':
                    k = pending.get(p.get('tool_use_id'), 'unknown')
                    chars[k] += size(p.get('content')); calls[k] += 1
    return chars, calls, usage, whole, ranged

root = os.path.expanduser(sys.argv[1])
for label, files in (('main', glob.glob(root + '/*.jsonl')),
                     ('subagents', glob.glob(root + '/**/subagents/*.jsonl', recursive=True))):
    chars, calls, usage, whole, ranged = scan(files)
    total = sum(chars.values()) or 1
    inp = usage['input_tokens'] + usage['cache_read_input_tokens'] + usage['cache_creation_input_tokens']
    print(f'## {label}: {len(files)} transcripts\n')
    print(f'input-side tokens {inp:,} (cache read {100*usage["cache_read_input_tokens"]/max(inp,1):.0f}%), output tokens {usage["output_tokens"]:,}; Read calls whole-file {whole}, ranged {ranged}\n')
    print('| kind | share of tool output | calls | mean chars |\n|---|---|---|---|')
    for k, v in chars.most_common():
        print(f'| {k} | {100*v/total:.1f}% | {calls[k]} | {v//max(calls[k],1)} |')
    print()
```
