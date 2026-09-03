// Command harness is the single deployable artifact of the agent harness: the
// installer, the drift check, the vendor-agnostic hooks, and the Definition-of-
// Done dispatcher.
//
// It replaces twelve shell and Python scripts. Two reasons for one binary
// rather than a directory of scripts. First, latency: the hooks run on every
// tool call, and a Go process starts in ~2.6 ms against ~22.8 ms for a Python
// one. Second, drift: twelve separately deployed files drift independently and
// silently, which is exactly how this harness came to be running a
// configuration nobody had shipped.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/kesokaj/styrman-harness/internal/dispatch"
	"github.com/kesokaj/styrman-harness/internal/gate"
	"github.com/kesokaj/styrman-harness/internal/hookio"
	"github.com/kesokaj/styrman-harness/internal/install"
	"github.com/kesokaj/styrman-harness/internal/trail"
)

const usage = `harness — agent harness tooling

Usage:
  harness install [--dry-run] [--target claude|antigravity]
        Deploy the harness. Overwrites what it manages; runs doctor at the end.
  harness doctor [--quiet]
        Report drift between this repo and the deployed harness. Exit 1 on drift.
  harness derive <claude-agent.md> <out.md>
        Derive an Antigravity agent definition from a Claude Code one.
  harness dispatch [--dry-run] [--json]
        Run the Definition-of-Done rows this working tree owes.
  harness resume [--agy]
        Print the Resume Protocol orientation block for the current directory.
  harness check
        Trail-freshness report. Exit 1 when the trail does not cover the code.
  harness hook <event>
        Run a hook. Events: sessionstart, userpromptsubmit, pretooluse,
        posttooluse, stop, subagentstart, subagentstop,
        agy-preinvocation, agy-stop, statusline.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	args := os.Args[2:]
	switch os.Args[1] {
	case "install":
		exit(install.Run(args))
	case "doctor":
		exit(install.Doctor(args))
	case "derive":
		exit(install.Derive(args))
	case "dispatch":
		exit(dispatch.Run(args))
	case "resume":
		runResume(args)
	case "check":
		runCheck()
	case "hook":
		if len(args) < 1 {
			fmt.Fprint(os.Stderr, usage)
			os.Exit(2)
		}
		runHook(args[0])
	case "-h", "--help", "help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "harness: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

func exit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "harness:", err)
		os.Exit(1)
	}
}

func runResume(args []string) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	text := gate.ResumeContext(dir)
	if len(args) > 0 && args[0] == "--agy" {
		hookio.EmitAgyContext(text)
		return
	}
	hookio.EmitContext(text)
}

// runCheck prints the trail-freshness report. It exits non-zero when the trail
// is stale, so /preflight can gate a release on it.
func runCheck() {
	wd, err := os.Getwd()
	if err != nil {
		exit(err)
	}
	proj, ok := trail.Resolve(wd)
	if !ok {
		fmt.Fprintf(os.Stderr, "harness: not a git repository: %s\n", wd)
		os.Exit(1)
	}
	report := proj.CheckFreshness()
	for _, line := range report.Lines {
		fmt.Println(line)
	}
	if report.Stale {
		os.Exit(1)
	}
}

// runHook dispatches one hook event. Every branch is silent on error and exits
// zero: a hook that breaks the session is worse than a hook that does nothing.
func runHook(event string) {
	switch event {
	case "statusline":
		raw, _ := io.ReadAll(os.Stdin)
		fmt.Print(gate.StatusLine(raw))
		return
	}

	p := hookio.Read()
	switch event {
	case "sessionstart":
		hookio.EmitContext(join(gate.ResumeContext(p.Dir()), gate.TrailNotice(p.Dir())))
	case "agy-preinvocation":
		if p.InvocationNum > 1 {
			hookio.EmitAgyContext("") // orient on the first invocation only
			return
		}
		hookio.EmitAgyContext(join(gate.ResumeContext(p.Dir()), gate.TrailNotice(p.Dir())))
	case "userpromptsubmit":
		gate.CaptureFeedback(p)
		hookio.EmitContext(gate.NoticeOnce(p))
	case "pretooluse":
		hookio.Emit(gate.SecretDecision(p))
	case "posttooluse":
		gate.Format(p)
		hookio.Emit(hookio.Allow())
	case "stop":
		hookio.Emit(gate.StopDecision(p))
	case "agy-stop":
		hookio.Emit(gate.AgyStopDecision(p))
	case "subagentstart":
		gate.SubagentStartMark(p)
		hookio.Emit(hookio.Allow())
	case "subagentstop":
		hookio.Emit(gate.SubagentStopDecision(p))
	default:
		fmt.Fprintf(os.Stderr, "harness: unknown hook event %q\n", event)
		os.Exit(2)
	}
}

func join(parts ...string) string {
	out := ""
	for _, p := range parts {
		if p == "" {
			continue
		}
		if out != "" {
			out += "\n"
		}
		out += p
	}
	return out
}
