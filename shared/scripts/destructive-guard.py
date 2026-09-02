#!/usr/bin/env python3
"""PreToolUse guard for destructive shell commands (Claude Code, Bash tool).

Policy — the operator's rule: agents may do anything inside their scratch
roots; destructive commands anywhere else ask first.
  - Scratch roots: /tmp/claude-<uid>/ (Claude Code's own scratchpad root),
    /tmp/agy-<uid>/, and the $TMPDIR-based equivalents. Never the un-suffixed
    /tmp/claude or /tmp/agy: those are shared-writable names anyone can plant.
    A root counts only while it exists, is a real directory (not a symlink),
    and is owned by this uid.
  - Covered destructive commands (not every destructive command — git
    checkout/restore/stash and `find -delete` fall to the classifier as they
    always did): rm, rmdir, unlink, shred, trash (also by absolute path such
    as /bin/rm, and behind env/busybox); git clean; git reset --hard, judged
    by where git itself says the worktree is (`rev-parse --show-toplevel`),
    so a scratch repo whose config points its worktree elsewhere asks.
  - A single simple destructive command whose every path (after ~, relative
    paths and symlinks are resolved) lies strictly under a scratch root ->
    "allow" (no prompt, no classifier). Any ".." component, any "/" after a
    glob character, any unexpanded "$" in the command, any git global option
    other than one `-C <dir>`, and any environment assignment outside a small
    safe set -> "ask".
  - Otherwise -> "ask" (a prompt, deterministically — this replaces the
    static Bash(rm *) etc. ask rules, which cannot see paths).
  - Compound or unparseable commands that mention a destructive command -> ask.
  - Shell writes that name a deployed harness path (~/.agents/bin, ~/.claude,
    ~/.gemini) -> ask: the Edit rules cover the file tools, this covers Bash.
  - Commands that are none of the above -> no opinion (empty output), so the
    normal permission flow applies.
Never emits "deny"; any internal error degrades to "ask", never to "allow".
Reads the hook JSON on stdin; always exits 0.  `--selftest` runs the built-in
decision table (creates and removes /tmp/claude-<uid>/.guard-selftest).
"""
import json
import os
import re
import shlex
import stat
import subprocess
import sys
import tempfile

UID = os.getuid()
_TMP = tempfile.gettempdir()          # $TMPDIR on macOS (/var/folders/...), /tmp on Linux
ROOTS = [f"/tmp/claude-{UID}", f"/tmp/agy-{UID}",
         os.path.join(_TMP, f"claude-{UID}"), os.path.join(_TMP, f"agy-{UID}")]
# Compare like with like: /tmp is a symlink to /private/tmp on macOS, so a
# resolved path must be checked against resolved roots.
ROOTS_NORM = sorted({os.path.normpath(r) for r in ROOTS})
ROOTS_REAL = sorted({os.path.realpath(r) for r in ROOTS})
HOME = os.path.expanduser("~")
PROTECTED = [os.path.realpath(os.path.join(HOME, p)) for p in (".agents/bin", ".claude", ".gemini")]
PROTECTED_MENTIONS = [p for p in PROTECTED] + [f"~/{p}" for p in (".agents/bin", ".claude", ".gemini")] + \
                     [f"$HOME/{p}" for p in (".agents/bin", ".claude", ".gemini")] + \
                     [f"${{HOME}}/{p}" for p in (".agents/bin", ".claude", ".gemini")]
DESTRUCTIVE = {"rm", "rmdir", "unlink", "shred", "trash"}
WRAPPERS = {"timeout", "time", "nice", "nohup", "stdbuf", "command", "builtin", "noglob", "env", "busybox"}
SAFE_ENV = {"LANG", "LC_ALL", "LC_CTYPE", "NO_COLOR", "TZ", "TERM"}
WRITERS = {"tee", "cp", "mv", "install", "rsync", "ln", "chmod", "chown", "truncate", "dd", "sed", "perl",
           "python", "python3", "patch", "curl", "wget", "tar", "unzip", "git", "touch", "mkdir", "rm", "rmdir",
           "unlink", "shred", "trash", "bash", "sh", "zsh", "dash", "eval", "xargs", "env", "find"}
# "$" of any kind: the shell expands it, the guard cannot — ask.
CONTROL = re.compile(r"\|\||&&|;|\||&|\$|`|\n|>|<")
MENTION = re.compile(r"(?:^|[\s;&|(/])(?:rm|rmdir|unlink|shred|trash)(?:\s|$)|\bgit\b.*\b(?:clean|reset\s+--hard)\b", re.DOTALL)


def root_usable(root):
    """A scratch root counts only if it exists, is a real directory, and is ours."""
    try:
        st = os.lstat(root)
    except OSError:
        return False
    return stat.S_ISDIR(st.st_mode) and not stat.S_ISLNK(st.st_mode) and st.st_uid == UID


def under_root(path, cwd):
    p = os.path.expanduser(path)
    if not os.path.isabs(p):
        p = os.path.join(cwd, p)
    if any(seg == ".." for seg in p.split(os.sep)):
        return False                     # never reason about ".." — ask
    m = re.search(r"[*?\[]", p)
    if m:
        if os.sep in p[m.start():]:
            return False                 # a "/" after a glob can follow a symlink out — ask
        p = p[: m.start()]               # literal prefix before the first glob char
    norm = os.path.normpath(p)
    real = os.path.realpath(p)           # follows symlinks
    norm_ok = any(norm.startswith(r + os.sep) and root_usable(r) for r in ROOTS_NORM)
    real_ok = any(real.startswith(r + os.sep) and root_usable(r) for r in ROOTS_REAL)
    return norm_ok and real_ok


def _expand(tok, cwd):
    """Best-effort path normalisation of a shell token (no expansion of $VAR)."""
    t = tok.split("=", 1)[1] if tok.startswith("--") and "=" in tok else tok
    t = t.replace("$HOME", HOME).replace("${HOME}", HOME)
    t = os.path.expanduser(t)
    if not os.path.isabs(t):
        t = os.path.join(cwd, t)
    return os.path.realpath(t)           # a symlink into a protected dir counts as the dir


def mentions_protected(cmd, cwd):
    """True when any path-like token, or the cwd itself, sits under a deployed harness dir."""
    if any(p in cmd for p in PROTECTED_MENTIONS):
        return True
    try:
        toks = shlex.split(cmd)
    except ValueError:
        toks = cmd.split()
    cands = [t for t in toks if "/" in t or t.startswith(("~", ".")) or "HOME" in t]
    cands.append(cwd)
    for t in cands:
        n = _expand(t, cwd)
        if any(n == p or n.startswith(p + os.sep) for p in PROTECTED):
            return True
    return False


def git_toplevel(repo):
    """Where git itself says the working tree and common dir are, or None."""
    try:
        out = subprocess.run(["git", "-c", "core.fsmonitor=false", "-C", repo, "rev-parse",
                              "--show-toplevel", "--git-common-dir"],
                             capture_output=True, text=True, timeout=5)
    except (OSError, subprocess.SubprocessError):
        return None
    if out.returncode != 0:
        return None
    lines = [ln.strip() for ln in out.stdout.splitlines() if ln.strip()]
    if len(lines) < 2:
        return None
    top, common = lines[0], lines[1]
    if not os.path.isabs(common):
        common = os.path.join(repo, common)   # git prints it relative to its process cwd, the -C dir
    return top, os.path.realpath(common)


def decide(cmd, cwd):
    """Return (decision, reason) or (None, None) when the command is not ours."""
    if mentions_protected(cmd, cwd):
        # A shell write to a deployed harness file: the Edit rules cover the
        # file tools; this covers redirections and writer commands in Bash.
        toks_guess = cmd.split()
        if re.search(r">|>>", cmd) or any(os.path.basename(t) in WRITERS for t in toks_guess):
            return "ask", "shell write or removal touching a deployed harness path (~/.agents/bin, ~/.claude, ~/.gemini)"
    if CONTROL.search(cmd):
        if MENTION.search(cmd):
            return "ask", "compound command includes a destructive command; the scratch-root allow covers only a single simple command"
        return None, None
    try:
        toks = shlex.split(cmd)
    except ValueError:
        return ("ask", "unparseable command mentions a destructive command") if MENTION.search(cmd) else (None, None)
    unsafe_env = False
    while toks and ("=" in toks[0] and not toks[0].startswith("-")):
        if toks[0].split("=", 1)[0] not in SAFE_ENV:
            unsafe_env = True
        toks = toks[1:]                  # leading VAR=value assignments
    if toks and (os.path.basename(toks[0]) if "/" in toks[0] else toks[0]) in WRAPPERS:
        i = next((k for k, t in enumerate(toks) if os.path.basename(t) in DESTRUCTIVE or t == "git"), None)
        if i is None:
            return None, None
        # env/busybox may carry VAR=value between the wrapper and the command
        if any("=" in t and t.split("=", 1)[0] not in SAFE_ENV for t in toks[1:i]):
            unsafe_env = True
        toks = toks[i:]
    if not toks:
        return None, None
    head = os.path.basename(toks[0]) if "/" in toks[0] else toks[0]   # /bin/rm and ./bin/rm are rm
    if head in DESTRUCTIVE:
        if unsafe_env:
            return "ask", f"{head} run with an environment assignment"
        paths, after_dd = [], False
        for t in toks[1:]:
            if after_dd:
                paths.append(t)
            elif t == "--":
                after_dd = True
            elif t.startswith("-"):
                continue
            else:
                paths.append(t)
        if not paths:
            return "ask", f"{head} without a path argument"
        outside = [p for p in paths if not under_root(p, cwd)]
        if outside:
            return "ask", f"{head} touches a path outside the agent scratch roots (or uses .. / a glob before a slash): {' '.join(outside[:3])}"
        return "allow", f"{head} confined to an agent scratch root"
    if head == "git":
        if unsafe_env:
            return "ask", "git run with an environment assignment"
        # Global options: exactly one `-C <dir>` is understood; anything else
        # (--work-tree, --git-dir, -c, repeated -C, ...) is asked about.
        repo, rest, i, seen_c = cwd, toks[1:], 0, False
        while i < len(rest) and rest[i].startswith("-"):
            if rest[i] == "-C" and not seen_c and i + 1 < len(rest):
                repo, seen_c, i = rest[i + 1], True, i + 2
                continue
            if MENTION.search(cmd):
                return "ask", f"git global option {rest[i]} is not understood by the scratch-root guard"
            return None, None
        sub = rest[i:]
        if sub[:1] == ["clean"] or (sub[:1] == ["reset"] and "--hard" in sub):
            # Any pathspec or extra positional after the subcommand: ask.
            if any(not t.startswith("-") for t in sub[1:]):
                return "ask", f"git {sub[0]} with a pathspec is not judged by the scratch-root guard"
            if not under_root(os.path.join(repo, "x"), cwd):
                return "ask", f"git {sub[0]} outside the agent scratch roots"
            # The -C directory is under a root; now ask git where the worktree
            # really is (core.worktree, gitfiles, and worktrees can point out).
            top = git_toplevel(_expand(repo, cwd))
            if top is None:
                return "ask", f"git {sub[0]}: could not resolve the working tree"
            if under_root(os.path.join(top[0], "x"), cwd) and under_root(os.path.join(top[1], "x"), cwd):
                return "allow", "git working-tree reset confined to an agent scratch root"
            return "ask", f"git {sub[0]}: the repository's working tree or git dir resolves outside the scratch roots"
    return None, None


def emit(decision, reason):
    print(json.dumps({"hookSpecificOutput": {"hookEventName": "PreToolUse",
                                             "permissionDecision": decision,
                                             "permissionDecisionReason": reason}}))


def selftest():
    base = f"/tmp/claude-{UID}/.guard-selftest"
    os.makedirs(os.path.join(base, "sub"), exist_ok=True)
    link = os.path.join(base, "etclink")
    if not os.path.islink(link):
        os.symlink("/etc", link)
    # a scratch repo whose worktree points outside must ask
    esc = os.path.join(base, "escaperepo")
    outside = tempfile.mkdtemp(prefix="guard-outside-")
    subprocess.run(["git", "init", "-q", esc], check=True, capture_output=True)
    subprocess.run(["git", "-C", esc, "config", "core.worktree", outside], check=True, capture_output=True)
    okrepo = os.path.join(base, "okrepo")
    subprocess.run(["git", "init", "-q", okrepo], check=True, capture_output=True)
    os.makedirs(os.path.join(okrepo, "sub"), exist_ok=True)
    binlink = os.path.join(base, "binlink")
    if not os.path.islink(binlink):
        os.symlink(os.path.join(HOME, ".agents/bin"), binlink)
    table = [
        ("ask", f"git -C {esc} clean -fdx", "/"),
        ("ask", f"git -C {esc} reset --hard", "/"),
        ("ask", "git clean -fdx", esc),
        ("allow", f"git -C {okrepo} clean -fdx", "/"),
        ("allow", "git reset --hard", okrepo),
        ("ask", f"git -C {okrepo} clean -fdx ../../../../etc", "/"),
        ("ask", f"git -C {base}/nonrepo clean -fdx", "/"),
        ("ask", "cat x > ~/.agents//bin/trail-check.sh", "/"),
        ("ask", "cp x ~/.agents/./bin/y", "/"),
        ("ask", "curl -o ~/.claude/settings.json http://x", "/"),
        ("ask", "wget -O settings.json http://x", HOME + "/.claude"),
        ("ask", "tar -C ~/.agents/bin -xf x.tar", "/"),
        ("ask", "touch hooks/x", HOME + "/.claude"),
        (None, "ls ~/.agents/bin", "/"),
        ("ask", "./bin/rm -rf /etc/x", "/"),
        ("ask", "/usr/bin/env rm -rf /etc/x", "/"),
        ("ask", 'bash -c "cp /tmp/x ~/.agents/bin/trail-check.sh"', "/"),
        ("ask", "echo x | xargs -I{} cp {} ~/.claude/settings.json", "/"),
        ("ask", f"cp /tmp/x {base}/binlink/trail-check.sh", "/"),
        ("allow", f"git -C {okrepo}/sub clean -fdx", "/"),
        ("ask", "cat x > ~/.cla''ude/settings.json", "/"),
        ("allow", f"rm -rf {base}/sub", "/"),
        ("allow", "rm -f sub/a.txt", base),
        ("allow", f"rm -rf {base}/*.log", "/"),
        ("allow", f"timeout 5 rm -rf {base}/sub", "/"),
        ("allow", f"LANG=C rm -rf {base}/sub", "/"),
        ("ask", "git clean -fdx", base),          # not a repository: git cannot say where the worktree is
        ("ask", f"git -C {base} reset --hard", "/"),
        ("ask", f"rm -rf {base}/sub /etc/passwd", "/"),
        ("ask", f"rm -rf {base}/../../../etc", "/"),
        ("ask", f"rm -rf {base}/etclink/", "/"),
        ("ask", f"rm -rf {base}/*/../../../../../etc", "/"),
        ("ask", f"rm -rf {base}/*/", "/"),
        ("ask", "rm -rf */../../..", base),
        ("ask", f"rm -rf /tmp/claude-{UID}", "/"),
        ("ask", "rm -rf /tmp/claude/x", "/"),
        ("ask", "rm -rf ~/x", "/"),
        ("ask", "rm", "/"),
        ("ask", f"cd {base} && rm -rf sub", "/"),
        ("ask", f"rm -rf {base}/sub & true", "/"),
        ("ask", 'rm -rf "$(echo /etc)"', "/"),
        ("ask", f"rm -rf {base}/x > /dev/null", "/"),
        ("ask", f"FOO=1 rm -rf {base}/sub", "/"),
        ("ask", f"git -C {base} --work-tree=/ clean -fdx", "/"),
        ("ask", f"git --work-tree=/ -C {base} clean -fdx", "/"),
        ("ask", f"git -C {base} -C .. clean -fdx", "/"),
        ("ask", f"GIT_WORK_TREE=/ git -C {base} clean -fdx", "/"),
        ("ask", "git clean -fdx", HOME),
        ("ask", "git reset --hard HEAD~1", HOME),
        ("ask", f"git -C {HOME} reset --hard", HOME),
        ("ask", f"git -C {HOME} clean -fdx", HOME),
        ("ask", "git -C . clean -fdx", HOME),
        ("ask", "trash ~/file", "/"),
        ("ask", f"rm -rf {base}/$EVIL", "/"),
        ("ask", "rm -rf $HOME/.ssh", base),
        ("ask", "/bin/rm -rf /home/x", "/"),
        ("allow", f"/bin/rm -rf {base}/sub", "/"),
        ("allow", f"env rm -rf {base}/sub", "/"),
        ("ask", f"env FOO=1 rm -rf {base}/sub", "/"),
        ("ask", f"busybox rm -rf /etc/x", "/"),
        ("ask", f"git -c core.hooksPath=/x -C {base} clean -fdx", "/"),
        ("ask", "cat x > ~/.agents/bin/trail-check.sh", "/"),
        ("ask", f"cp x {HOME}/.claude/settings.json", "/"),
        ("ask", "sed -i s/a/b/ ~/.gemini/GEMINI.md", "/"),
        ("ask", f"rm {HOME}/.agents/bin/resume-context.sh", "/"),
        (None, "git status", "/"),
        (None, "ls -la /etc", "/"),
        (None, "echo rm", "/"),
        (None, "cat ~/.claude/CLAUDE.md", "/"),
        (None, "diff a ~/.agents/bin/trail-check.sh", "/"),
    ]
    failed = 0
    # root_usable: a symlink or a foreign-owned dir never counts as a root
    lnk = os.path.join(base, "rootlink")
    if not os.path.islink(lnk):
        os.symlink("/etc", lnk)
    for want, fn, arg in [(False, root_usable, lnk), (True, root_usable, base), (False, root_usable, "/etc"), (False, root_usable, base + "/missing")]:
        got = fn(arg)
        mark = "ok " if got == want else "FAIL"
        if got != want:
            failed += 1
        print(f"{mark} {str(want):5} {str(got):5} root_usable({arg})")
    os.unlink(lnk)
    import shutil
    for want, cmd, cwd in table:
        got = decide(cmd, cwd)[0]
        mark = "ok " if got == want else "FAIL"
        if got != want:
            failed += 1
        print(f"{mark} {str(want):5} {str(got):5} {cmd}")
    os.unlink(link)
    os.rmdir(os.path.join(base, "sub"))
    os.unlink(binlink)
    shutil.rmtree(esc)
    shutil.rmtree(okrepo)
    shutil.rmtree(outside)
    os.rmdir(base)
    print(f"selftest: {len(table) - failed}/{len(table)} passed")
    return failed


def main():
    if len(sys.argv) > 1 and sys.argv[1] == "--selftest":
        sys.exit(1 if selftest() else 0)
    raw = ""
    try:
        raw = sys.stdin.read()
        data = json.loads(raw) if raw.strip() else {}
        if data.get("tool_name", "Bash") != "Bash":
            return
        cmd = (data.get("tool_input") or {}).get("command") or ""
        cwd = data.get("cwd") or os.getcwd()
        decision, reason = decide(cmd, cwd)
    except Exception as exc:  # never fail open
        decision, reason = ("ask", f"destructive-guard error: {exc}") if MENTION.search(raw) else (None, None)
    if decision is None:
        return
    emit(decision, reason)


if __name__ == "__main__":
    main()
    sys.exit(0)
