# usage: python3 -I accept.py <app dir> <work dir>   — builds ./cmd/crm, runs it, checks the contract black-box.
import json, os, subprocess, sys, time, urllib.request, urllib.error, socket
app, work = os.path.abspath(sys.argv[1]), os.path.abspath(sys.argv[2]); os.makedirs(work, exist_ok=True)
# Agent-written code is built and run with a minimal environment: no operator tokens, a throwaway HOME and Go cache.
MINENV = {'PATH': os.environ['PATH'], 'HOME': work, 'GOCACHE': os.path.join(work, 'gocache'), 'GOFLAGS': '-mod=readonly', 'GOPROXY': 'off', 'GOTOOLCHAIN': 'local'}
binp = os.path.join(work, 'crm'); data = os.path.join(work, 'data.json')
if os.path.exists(data): os.remove(data)
b = subprocess.run(['go', 'build', '-o', binp, './cmd/crm'], cwd=app, capture_output=True, text=True, env=MINENV)
if b.returncode: print('BUILD FAIL', b.stderr); sys.exit(1)
s = socket.socket(); s.bind(('127.0.0.1', 0)); port = s.getsockname()[1]; s.close()
base = f'http://127.0.0.1:{port}'
def start():
    p = subprocess.Popen([binp], cwd=work, env={**MINENV, 'PORT': str(port), 'DATA_FILE': data}, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    for _ in range(100):
        try: urllib.request.urlopen(base + '/api/customers', timeout=1); return p
        except Exception: time.sleep(0.1)
    return p
def req(m, path, body=None):
    d = json.dumps(body).encode() if body is not None else None
    r = urllib.request.Request(base + path, data=d, method=m, headers={'Content-Type': 'application/json'})
    try:
        with urllib.request.urlopen(r, timeout=5) as f: raw = f.read(); code = f.status; ct = f.headers.get('Content-Type', '')
    except urllib.error.HTTPError as e: raw = e.read(); code = e.code; ct = e.headers.get('Content-Type', '')
    except OSError: return 0, '', ''
    try: return code, json.loads(raw), ct
    except Exception: return code, raw.decode(errors='replace'), ct
res = []
def raw(m, path, body):
    r = urllib.request.Request(base + path, data=body, method=m, headers={'Content-Type': 'application/json'})
    try:
        with urllib.request.urlopen(r, timeout=5) as f: return f.status
    except urllib.error.HTTPError as e: return e.code
    except OSError: return 0
def check(name, ok): res.append((name, bool(ok)))
p = start()
try:
    c, alice, _ = req('POST', '/api/customers', {'name': 'Alice', 'email': 'alice@example.com', 'phone': '123'})
    check('create customer 201 + string id', c == 201 and isinstance(alice, dict) and isinstance(alice.get('id'), str) and alice.get('name') == 'Alice')
    aid = alice.get('id') if isinstance(alice, dict) else ''
    check('empty name 400', req('POST', '/api/customers', {'name': '', 'email': 'x@example.com'})[0] == 400)
    check('bad email 400', req('POST', '/api/customers', {'name': 'X', 'email': 'not-an-email'})[0] == 400)
    check('duplicate email case-insensitive 409', req('POST', '/api/customers', {'name': 'A2', 'email': 'ALICE@Example.com'})[0] == 409)
    check('malformed JSON 400', raw('POST', '/api/customers', b'{bad') == 400)
    c, lst, _ = req('GET', '/api/customers'); check('list customers', c == 200 and isinstance(lst, list) and len(lst) == 1)
    check('get customer', req('GET', f'/api/customers/{aid}')[0] == 200)
    check('get unknown customer 404', req('GET', '/api/customers/nope')[0] == 404)
    c, bob, _ = req('POST', '/api/customers', {'name': 'Bob', 'email': 'bob@example.com', 'phone': ''}); bid = bob.get('id') if isinstance(bob, dict) else ''
    B = lambda cid, s, e, n='': req('POST', '/api/bookings', {'customerId': cid, 'start': s, 'end': e, 'note': n})
    c, b1, _ = B(aid, '2026-11-02T10:00:00Z', '2026-11-02T11:00:00Z', 'first')
    check('create booking 201 confirmed', c == 201 and isinstance(b1, dict) and b1.get('status') == 'confirmed' and b1.get('customerId') == aid and b1.get('note') == 'first')
    b1id = b1.get('id') if isinstance(b1, dict) else ''
    check('overlap 409', B(bid, '2026-11-02T10:30:00Z', '2026-11-02T11:30:00Z')[0] == 409)
    check('containing overlap 409', B(bid, '2026-11-02T09:00:00Z', '2026-11-02T12:00:00Z')[0] == 409)
    check('overlap across time zones 409', B(bid, '2026-11-02T11:30:00+01:00', '2026-11-02T12:30:00+01:00')[0] == 409)
    check('touching end==start 201', B(bid, '2026-11-02T11:00:00Z', '2026-11-02T12:00:00Z')[0] == 201)
    check('end<=start 400', B(aid, '2026-11-03T10:00:00Z', '2026-11-03T10:00:00Z')[0] == 400)
    check('bad time 400', B(aid, 'tomorrow', '2026-11-03T10:00:00Z')[0] == 400)
    check('unknown customer 400', B('nope', '2026-11-04T10:00:00Z', '2026-11-04T11:00:00Z')[0] == 400)
    B(aid, '2026-11-01T08:00:00Z', '2026-11-01T09:00:00Z', 'early')
    c, all_, _ = req('GET', '/api/bookings'); starts = [x.get('start') for x in all_] if isinstance(all_, list) else []
    check('list bookings sorted by start (3)', c == 200 and len(starts) == 3 and [x[:10] + x[11:13] for x in starts] == ['2026-11-0108', '2026-11-0210', '2026-11-0211'])
    c, fa, _ = req('GET', f'/api/bookings?customerId={aid}'); check('filter by customerId', c == 200 and isinstance(fa, list) and len(fa) == 2)
    c, fw, _ = req('GET', '/api/bookings?from=2026-11-02T10:30:00Z&to=2026-11-02T11:00:00Z'); check('from/to overlap window', c == 200 and isinstance(fw, list) and len(fw) == 1 and fw[0].get('id') == b1id)
    c, cb, _ = req('POST', f'/api/bookings/{b1id}/cancel'); check('cancel 200 cancelled', c == 200 and isinstance(cb, dict) and cb.get('status') == 'cancelled')
    check('cancel unknown 404', req('POST', '/api/bookings/nope/cancel')[0] == 404)
    check('cancelled slot rebookable', B(bid, '2026-11-02T10:00:00Z', '2026-11-02T11:00:00Z')[0] == 201)
    c, html, ct = req('GET', '/'); check('GET / html', c == 200 and 'html' in ct.lower())
    # concurrency: 20 parallel identical bookings, exactly one must win
    import concurrent.futures as cf
    with cf.ThreadPoolExecutor(20) as ex: codes = list(ex.map(lambda _: B(aid, '2026-12-01T10:00:00Z', '2026-12-01T11:00:00Z')[0], range(20)))
    check('concurrent double-booking: exactly one 201', codes.count(201) == 1 and codes.count(409) == 19)
finally:
    p.terminate(); p.wait()
p = start()
try:
    c, lst, _ = req('GET', '/api/customers'); c2, bl, _ = req('GET', '/api/bookings')
    check('persists across restart', isinstance(lst, list) and len(lst) == 2 and isinstance(bl, list) and len(bl) == 5)
finally:
    p.terminate(); p.wait()
for n, ok in res: print(('PASS ' if ok else 'FAIL ') + n)
print(f'SCORE {sum(ok for _, ok in res)}/{len(res)}')
