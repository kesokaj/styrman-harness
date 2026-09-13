// Command sh-agent installs, resyncs, inspects, and removes the sh-agent
// harness — global context, skills, and subagents — for Claude Code and the
// Antigravity CLI.
//
//	sh-agent install [--dry-run]     deploy the repo into ~/.claude and ~/.gemini/config
//	sh-agent status                  show drift between the repo and what is installed
//	sh-agent diff                    unified diff of every drifted file
//	sh-agent update [--check]        git fetch; pull --ff-only and reinstall when behind
//	sh-agent uninstall [--dry-run]   remove exactly what install wrote
//	sh-agent list                    print the skill and agent inventory
//	sh-agent render                  print the assembled context file
//	sh-agent adopt <name|path>       copy a skill from elsewhere into the repo
//	sh-agent version
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/kesokaj/styrman-harness/internal/harness"
)

const usage = `sh-agent — lifecycle tool for the sh-agent harness

Commands:
  install [--vendor V] [--dry-run]  deploy context, skills, and agents into ~/.claude and/or ~/.gemini/config (alias: sync)
  status                  show what is installed and whether it drifted from the repo (exit 1 on drift)
  diff                    unified diff of every drifted file
  update [--check]        git fetch the harness repo; pull --ff-only and reinstall when behind (--check only reports, exit 2 when behind)
  uninstall [--dry-run]   remove exactly the files install wrote, nothing else
  list                    print the skill and agent inventory
  render                  print the assembled context file
  adopt <name|path>       copy a skill installed elsewhere into the repo's skills/ (--force to replace)
  version                 print the build version

Flags accepted by every command:
  --repo DIR        harness checkout (default: --repo, $SH_AGENT_REPO, the working directory or its parents, then the last install)
  --home DIR        install under DIR instead of $HOME (for testing)
  --vendor VENDOR   target vendor: all, antigravity (or agy), claude (default: from manifest if installed, else all)
`

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		var code exitCode
		if errors.As(err, &code) {
			os.Exit(int(code))
		}
		fmt.Fprintln(os.Stderr, "sh-agent:", err)
		os.Exit(1)
	}
}

// exitCode carries a non-error exit status such as "drift found".
type exitCode int

func (e exitCode) Error() string { return fmt.Sprintf("exit %d", int(e)) }

func run(args []string, out io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(out, usage)
		return nil
	}
	cmd, rest := args[0], args[1:]

	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	fs.SetOutput(out)
	repoFlag := fs.String("repo", "", "harness checkout")
	homeFlag := fs.String("home", "", "install under this directory instead of $HOME")
	vendorFlag := fs.String("vendor", "", "target vendor: all, antigravity (or agy), claude")
	dryRun := fs.Bool("dry-run", false, "print what would change, change nothing")
	check := fs.Bool("check", false, "update: only report whether the repo is behind")
	force := fs.Bool("force", false, "adopt: replace an existing skill")
	if err := fs.Parse(reorder(rest)); err != nil {
		return err
	}

	home := *homeFlag
	if home == "" {
		h, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		home = h
	}

	switch cmd {
	case "version":
		fmt.Fprintln(out, "sh-agent", buildVersion())
		return nil
	case "uninstall":
		return harness.Uninstall(home, *dryRun, out)
	}

	repo, err := harness.FindRepo(*repoFlag, home)
	if err != nil {
		return err
	}

	switch cmd {
	case "install", "sync":
		vendor, err := resolveVendor(*vendorFlag, home)
		if err != nil {
			return err
		}
		// The deployed binary must be the repo's own build: rebuild and let
		// that build install, unless this process is that build (ChildEnv)
		// or there is no toolchain to build with.
		_, goErr := exec.LookPath("go")
		if shouldRebuild(goErr == nil, os.Getenv(harness.ChildEnv) == "1") {
			var pass []string
			if *homeFlag != "" {
				pass = append(pass, "--home", home)
			}
			if *vendorFlag != "" {
				pass = append(pass, "--vendor", vendor)
			}
			if *dryRun {
				pass = append(pass, "--dry-run")
			}
			return harness.Reinstall(repo, pass, out)
		}
		if goErr != nil {
			fmt.Fprintln(out, "warning: go is not on PATH; the running binary was copied, not rebuilt — install Go and run `sh-agent update`")
		}
		self, err := selfExecutable()
		if err != nil {
			return err
		}
		files, err := harness.Plan(repo, self, vendor)
		if err != nil {
			return err
		}
		if err := harness.Install(repo, home, files, vendor, *dryRun, out); err != nil {
			return err
		}
		return pathHint(home, out)

	case "status":
		return status(repo, home, *vendorFlag, out)

	case "diff":
		vendor, err := resolveVendor(*vendorFlag, home)
		if err != nil {
			return err
		}
		files, err := harness.Plan(repo, nil, vendor)
		if err != nil {
			return err
		}
		return harness.Diff(home, files, out)

	case "update":
		return update(repo, home, *check, *homeFlag, *vendorFlag, out)

	case "list":
		inv, err := harness.LoadInventory(repo)
		if err != nil {
			return err
		}
		fmt.Fprint(out, harness.RenderInventory(inv))
		return nil

	case "render":
		inv, err := harness.LoadInventory(repo)
		if err != nil {
			return err
		}
		data, err := harness.RenderContext(repo, inv)
		if err != nil {
			return err
		}
		_, err = out.Write(data)
		return err

	case "adopt":
		if fs.NArg() != 1 {
			return errors.New("usage: sh-agent adopt <skill-name-or-path> [--force]")
		}
		dest, err := harness.Adopt(repo, home, fs.Arg(0), *force)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "adopted into %s\nnext: review it, run `sh-agent install`, then commit and push the repo\n", dest)
		return nil
	}
	return fmt.Errorf("unknown command %q\n\n%s", cmd, usage)
}

func resolveVendor(flagVal, home string) (string, error) {
	if flagVal != "" {
		return harness.NormalizeVendor(flagVal)
	}
	st, found, err := harness.LoadState(home)
	if err == nil && found && st.Vendor != "" {
		return harness.NormalizeVendor(st.Vendor)
	}
	return harness.VendorAll, nil
}

func status(repo, home, vendorFlag string, out io.Writer) error {
	info := harness.Inspect(repo)
	dirty := "clean"
	if info.Dirty {
		dirty = "uncommitted changes"
	}
	fmt.Fprintf(out, "repo       %s (%s @ %s, %s)\n", repo, info.Branch, info.Commit, dirty)

	st, found, err := harness.LoadState(home)
	if err != nil {
		return err
	}
	if !found {
		fmt.Fprintln(out, "installed  nothing — run `sh-agent install`")
		return exitCode(1)
	}
	vDisplay := st.Vendor
	if vDisplay == "" {
		vDisplay = harness.VendorAll
	}
	fmt.Fprintf(out, "installed  %s from commit %s (vendor: %s, %d files, manifest ~/.config/sh-agent/state.json)\n", st.InstalledAt, st.Commit, vDisplay, len(st.Files))
	rev := buildVersion()
	binaryLine := "binary     " + rev
	switch {
	case !binaryOlder(rev, info.Commit):
	case rev == "(devel)" || rev == "unknown":
		binaryLine += " — no build stamp; run `sh-agent update`"
	default:
		binaryLine += fmt.Sprintf(" — older than the repo (%s); run `sh-agent update`", info.Commit)
	}
	fmt.Fprintln(out, binaryLine)
	if last, ok := harness.LastRetro(home); ok {
		days := int(time.Since(last).Hours() / 24)
		hint := ""
		if days > 14 {
			hint = " — past 14 days, run /retro at the end of the next session"
		}
		fmt.Fprintf(out, "retro      last proposal %s (%d days ago)%s\n", last.Format("2006-01-02"), days, hint)
	} else {
		fmt.Fprintln(out, "retro      none yet — run /retro at the end of a session with corrections")
	}

	vendor, err := resolveVendor(vendorFlag, home)
	if err != nil {
		return err
	}

	inv, err := harness.LoadInventory(repo)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "repo has   %d skills, %d agents, %d commands, %d rules\n", len(inv.Skills), len(inv.Agents), len(inv.Commands), len(inv.Rules))
	if vendor == harness.VendorAll || vendor == harness.VendorClaude {
		fmt.Fprintf(out, "claude     ~/%s  (AGENTS.md, CLAUDE.md, skills/, agents/, commands/, rules/)\n", harness.ClaudeRoot)
	}
	if vendor == harness.VendorAll || vendor == harness.VendorAntigravity {
		fmt.Fprintf(out, "agy        ~/%s  (rules/AGENTS.md, skills/, agents/, workflows/, rules/)\n", harness.AgyRoot)
	}

	files, err := harness.Plan(repo, nil, vendor)
	if err != nil {
		return err
	}
	report, err := harness.Status(home, files)
	if err != nil {
		return err
	}
	for _, e := range report.Entries {
		if e.State != "ok" {
			fmt.Fprintf(out, "  %-8s ~/%s\n", e.State, e.Rel)
		}
	}
	if report.Counts["foreign"] > 0 {
		fmt.Fprintln(out, "note       ~/.claude/CLAUDE.md is yours; make sure it contains the line `@AGENTS.md`")
	}
	if !report.Drift() {
		fmt.Fprintf(out, "drift      none — %d files match the repo\n", report.Counts["ok"])
		return nil
	}
	fmt.Fprintf(out, "drift      %d missing, %d differ, %d stale — run `sh-agent diff`, then `sh-agent install`\n",
		report.Counts["missing"], report.Counts["differs"], report.Counts["stale"])
	return exitCode(1)
}

func update(repo, home string, checkOnly bool, homeFlag, vendorFlag string, out io.Writer) error {
	behind, err := harness.Behind(repo)
	if err != nil {
		return err
	}
	if behind == 0 {
		fmt.Fprintf(out, "up to date with upstream (%s)\n", harness.Inspect(repo).Commit)
		if checkOnly {
			return nil
		}
	} else {
		fmt.Fprintf(out, "behind upstream by %d commit(s)\n", behind)
		if checkOnly {
			return exitCode(2)
		}
		if err := harness.Pull(repo, out); err != nil {
			return err
		}
		fmt.Fprintf(out, "pulled to %s\n", harness.Inspect(repo).Commit)
	}
	var args []string
	if homeFlag != "" {
		args = append(args, "--home", home)
	}
	vendor, err := resolveVendor(vendorFlag, home)
	if err != nil {
		return err
	}
	if vendor != harness.VendorAll {
		args = append(args, "--vendor", vendor)
	}
	return harness.Reinstall(repo, args, out)
}

// shouldRebuild decides whether install re-executes through a fresh build
// of the repo: yes whenever a toolchain is present, unless this process is
// already the fresh build that Reinstall spawned.
func shouldRebuild(goOnPath, child bool) bool {
	return goOnPath && !child
}

// binaryOlder reports whether the running binary is not the repo's current
// build: unstamped (`go run`, or a build outside git) or stamped with a
// different revision. Prefixes are compared both ways because git's short
// hash length is not fixed.
func binaryOlder(rev, commit string) bool {
	rev = strings.TrimSuffix(rev, "+dirty")
	if rev == "" || commit == "" {
		return false
	}
	if rev == "(devel)" || rev == "unknown" {
		return true
	}
	return !strings.HasPrefix(rev, commit) && !strings.HasPrefix(commit, rev)
}

// reorder moves flags in front of positional arguments so that
// `sh-agent adopt my-skill --home /tmp/x` parses; the flag package stops at
// the first positional argument otherwise.
func reorder(args []string) []string {
	takesValue := map[string]bool{
		"-repo": true, "--repo": true,
		"-home": true, "--home": true,
		"-vendor": true, "--vendor": true,
	}
	var flags, positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--":
			positional = append(positional, args[i+1:]...)
			i = len(args)
		case strings.HasPrefix(a, "-"):
			flags = append(flags, a)
			if takesValue[a] && i+1 < len(args) {
				flags = append(flags, args[i+1])
				i++
			}
		default:
			positional = append(positional, a)
		}
	}
	return append(flags, positional...)
}

// selfExecutable returns the running binary so install can copy it to
// ~/.local/bin. On the normal path this is the build Reinstall just made;
// without a toolchain it is whatever binary is running, and status says so.
func selfExecutable() ([]byte, error) {
	p, err := os.Executable()
	if err != nil {
		return nil, err
	}
	if p, err = filepath.EvalSymlinks(p); err != nil {
		return nil, err
	}
	return os.ReadFile(p)
}

func pathHint(home string, out io.Writer) error {
	bin := filepath.Join(home, filepath.Dir(harness.BinPath))
	for _, p := range filepath.SplitList(os.Getenv("PATH")) {
		if p == bin || strings.TrimSuffix(p, "/") == bin {
			return nil
		}
	}
	fmt.Fprintf(out, "note: %s is not on PATH; add it so `sh-agent` is available everywhere\n", bin)
	return nil
}

func buildVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	rev, modified := "", ""
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = s.Value
		case "vcs.modified":
			if s.Value == "true" {
				modified = "+dirty"
			}
		}
	}
	if rev == "" {
		return info.Main.Version
	}
	if len(rev) > 12 {
		rev = rev[:12]
	}
	return rev + modified
}
