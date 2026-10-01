// Command sh-agent installs, resyncs, inspects, and removes the sh-agent
// harness — global context, skills, and subagents — for Claude Code and the
// Antigravity CLI.
//
//	sh-agent install [--dry-run]     deploy the repo into ~/.claude and ~/.gemini/config
//	sh-agent status                  show drift between the repo and what is installed
//	sh-agent diff                    unified diff of every drifted file
//	sh-agent update [--check|--dry-run]  checkout: git pull and reinstall; release binary: download and verify the newest release
//	sh-agent update --from-zip FILE  zip place: install the harness in a repository zip with the running binary
//	sh-agent uninstall [--dry-run]   remove exactly what install wrote, and the dream timer
//	sh-agent list                    print the skill and agent inventory
//	sh-agent render                  print the assembled context file
//	sh-agent adopt <name|path>       copy a skill from elsewhere into the repo
//	sh-agent report [--dry-run]      send this machine's learning-loop files and metrics to the feedback repo, keeping nothing
//	sh-agent search <words>          ranked search of this checkout's files and trail, reconciled against git first
//	sh-agent index [--rebuild]       reconcile (or rebuild) the search index and print its counts
//	sh-agent dream --all|--project DIR [--dry-run]  propose rules and skills from every project's memory, trail, feedback and transcripts
//	sh-agent dream --install-timer|--remove-timer [--dry-run]  install or remove the nightly timer
//	sh-agent memory add|show         write one memory (body on stdin) or print the index of this project's .agents/memory/
//	sh-agent version
package main

import (
	"bytes"
	"debug/buildinfo"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/kesokaj/styrman-harness/internal/dream"
	"github.com/kesokaj/styrman-harness/internal/harness"
	"github.com/kesokaj/styrman-harness/internal/memory"
	"github.com/kesokaj/styrman-harness/internal/search"
)

const usage = `sh-agent — lifecycle tool for the sh-agent harness

Commands:
  install [--vendor V] [--dry-run]  deploy context, skills, and agents into ~/.claude and/or ~/.gemini/config (alias: sync)
  status                  show what is installed and whether it drifted from the repo (exit 1 on drift)
  diff                    unified diff of every drifted file
  update [--check|--dry-run]  checkout: git fetch, pull --ff-only and reinstall when behind; release binary: download the newest release, verify its signed checksums, replace the binary and reinstall (--check only reports, exit 2 when behind; --dry-run reports what update would do and changes nothing)
  update --from-zip FILE [--dry-run]  zip place: unpack a repository zip you downloaded, install its harness and the running binary (nothing in the zip replaces the binary), delete the unpacked copy; with the release's SHA256SUMS and SHA256SUMS.sig beside it the zip is verified, otherwise it is installed marked UNVERIFIED
  uninstall [--dry-run]   remove exactly the files install wrote, and the dream timer when one is installed; nothing else
  list                    print the skill and agent inventory
  render                  print the assembled context file
  adopt <name|path>       copy a skill installed elsewhere into the repo's skills/ (--force to replace)
  report [--feedback-repo DIR] [--dry-run]  send a redacted copy of ~/.config/sh-agent and a metrics record to the feedback repo (through a throwaway clone holding only this machine's reports, removed after the push; --feedback-repo names a checkout to write into instead)
  search [--limit N] [--kind K] [--global] <words...>  ranked search of every file and trail record in this git checkout: path:line, kind, record status, snippet; the index (.agents/scratch/search.db) is reconciled against git first, so it never answers from another branch
  search --refs <BUG-NNN|ADR-NNN|path> [--kind K] [--global]  every record naming an ID, a file or a directory (trailing /), with the ID's own status
  search --open [--kind K] [--global]  open bugs and plans with unticked boxes
      --global searches the global search index (~/.config/sh-agent/search.db): ~/.config/sh-agent's logs and proposals and the .agents/docs trail of every harness-scaffolded checkout that has run search; --kind is one of bug adr plan change handoff review testing retro spec deploy trail code doc feedback proposal; --limit 1-50, default 10
  index [--global] [--rebuild]  reconcile the search index, or delete and rebuild it, and print its counts
  memory add --type T --name SLUG --description TEXT  save one memory, its body on stdin, into this project's .agents/memory/ and its MEMORY.md line (T: user, feedback, project or reference); every CLI writes memory this way (R10)
  dream --install-timer | --remove-timer [--dry-run]  install (systemd user timer on Linux, launchd agent on macOS) or remove the nightly "dream --all", run with the PATH of the installing shell; uninstall removes it too
  dream --all | --project DIR [--dry-run]  the background dreamer: for every indexed project (or one) that changed since its last dream, send its memory, new trail, feedback and both CLIs' transcripts, redacted, to a CLI with no tools, and write the proposals and draft skills it returns under ~/.config/sh-agent/proposals/dream/ (R3's exception; --vendor picks the CLI, default from the manifest)
  memory show             print this project's memory index, capped at what a session loads, plus a hint when Claude Code's native memory is still on here
  version                 print the build version

Flags accepted by every command:
  --repo DIR        harness checkout, for a build made from one (default: --repo, $SH_AGENT_REPO, the working directory or its parents, then the last install). A release binary installs the harness built into it and refuses --repo (ADR-034)
  --home DIR        install under DIR instead of $HOME (for testing)
  --vendor VENDOR   target vendor: all, antigravity (or agy), claude (default: from manifest if installed, else all)
`

// version is set by a release build, -ldflags "-X main.version=vX.Y.Z"; empty
// otherwise, and buildVersion reads the build info instead.
var version string

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
	feedbackFlag := fs.String("feedback-repo", "", "report: feedback checkout (never cloned into)")
	fromZip := fs.String("from-zip", "", "update: install the harness in this repository zip with the running binary")
	limit := fs.Int("limit", 0, "search: hits to print, 1-50 (default 10)")
	kind := fs.String("kind", "", "search: only hits of this kind")
	refs := fs.String("refs", "", "search: records naming this ID or path")
	openRecords := fs.Bool("open", false, "search: open bugs and plans with unticked boxes")
	global := fs.Bool("global", false, "search, index: the global search index")
	rebuild := fs.Bool("rebuild", false, "index: delete and rebuild the index")
	memType := fs.String("type", "", "memory add: user, feedback, project or reference")
	memName := fs.String("name", "", "memory add: the memory's slug")
	memDesc := fs.String("description", "", "memory add: one line")
	dreamAll := fs.Bool("all", false, "dream: every project in the global search index")
	dreamProject := fs.String("project", "", "dream: one checkout")
	installTimer := fs.Bool("install-timer", false, "dream: install the nightly user timer")
	removeTimer := fs.Bool("remove-timer", false, "dream: remove the nightly user timer")
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
	// Absolute before anything runs: the checkout path rebuilds and runs
	// the build's install with cmd.Dir = repo, so a relative --home landed
	// inside the repository on a real install and under the invoking
	// directory on --dry-run (BUG-080).
	if abs, err := filepath.Abs(home); err == nil {
		home = abs
	}

	// A zip place (R8, ADR-038): the operator named the archive, so no tree
	// is resolved — not the working directory, not the manifest, and not a
	// release binary's embedded files. Every refusal comes before any command
	// runs: `uninstall --from-zip` must not uninstall, `report --from-zip`
	// must not report, and an empty value is refused rather than read as no
	// flag — that would fall through to a normal update, which on a release
	// binary downloads and replaces the binary.
	zipSet := false
	fs.Visit(func(f *flag.Flag) { zipSet = zipSet || f.Name == "from-zip" })
	if zipSet {
		switch {
		case cmd != "update":
			return fmt.Errorf("--from-zip belongs to update: sh-agent update --from-zip <file>")
		case *fromZip == "" || strings.HasPrefix(*fromZip, "-"):
			// A value that looks like a flag was a flag: `--from-zip --dry-run`
			// must not read --dry-run as the file and run a real update.
			return errors.New("--from-zip needs a file: sh-agent update --from-zip <file>")
		case *check || *repoFlag != "":
			return errors.New("--from-zip installs the named zip; --check/--repo do not apply to it")
		}
		return updateFromZip(*fromZip, home, *vendorFlag, *dryRun, out)
	}

	switch cmd {
	case "version":
		fmt.Fprintln(out, "sh-agent", buildVersion())
		return nil
	case "uninstall":
		// The dream timer is not in the manifest (the next install would delete
		// it as stale); uninstall removes it by R4's dream clause, also when no
		// manifest is left, and a dry run only says so.
		if dream.TimerInstalled(home) {
			if *dryRun {
				fmt.Fprintln(out, "[dry-run] would remove the dream timer")
			} else if err := dream.RemoveTimer(home, runtime.GOOS); err != nil {
				return err
			} else {
				fmt.Fprintln(out, "removed the dream timer")
			}
		}
		return harness.Uninstall(home, *dryRun, out)
	case "report":
		// Before ResolveSource, which errors when no harness is found: a
		// machine with no checkout must still be able to report.
		return report(*repoFlag, *feedbackFlag, home, *dryRun, out)
	case "memory":
		// Before ResolveSource: memory needs no harness tree (R10).
		return memoryCmd(fs.Args(), memory.Entry{Type: *memType, Name: *memName, Description: *memDesc}, home, *vendorFlag, out)
	case "dream":
		// Before ResolveSource: the dream job needs no harness tree (R3's exception).
		vendor, err := resolveVendor(*vendorFlag, home)
		if err != nil {
			return err
		}
		switch {
		case *installTimer && *removeTimer:
			// A usage error exits 2 like every other dream usage error, and
			// says why on the caller's writer (audit 2026-09-30, IN12, IN13).
			fmt.Fprintln(out, "sh-agent dream: --install-timer or --remove-timer, not both")
			return exitCode(2)
		case (*installTimer || *removeTimer) && *dryRun:
			// A new side effect goes below the read-only mode (context § Hard
			// rules): the dry run names what would change and calls nothing.
			what := "install the nightly dream timer (sh-agent dream --all at 03:00)"
			if *removeTimer {
				what = "remove the nightly dream timer"
			}
			fmt.Fprintln(out, "[dry-run] would "+what+"; installed now: "+strconv.FormatBool(dream.TimerInstalled(home)))
		case *installTimer:
			err = dream.InstallTimer(home, os.Getenv("PATH"), runtime.GOOS)
			if err == nil {
				fmt.Fprintln(out, "dream timer installed: sh-agent dream --all runs nightly at 03:00")
			}
		case *removeTimer:
			err = dream.RemoveTimer(home, runtime.GOOS)
			if err == nil {
				fmt.Fprintln(out, "dream timer removed")
			}
		default:
			err = dream.Run(dream.Options{Home: home, Vendor: vendor, All: *dreamAll, Project: *dreamProject, DryRun: *dryRun, Out: out})
		}
		var ee *dream.ExitError
		if errors.As(err, &ee) {
			fmt.Fprintln(out, ee.Msg)
			return exitCode(ee.Code)
		}
		return err
	case "search", "index":
		// Before ResolveSource too: searching needs no harness tree, and a
		// release-binary or zip place has none (plan D10).
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		o := search.Options{Dir: wd, Home: home, Words: fs.Args(), Refs: *refs, Open: *openRecords,
			Kind: *kind, Limit: *limit, Global: *global, Rebuild: *rebuild}
		if cmd == "index" {
			return search.Index(o, out)
		}
		return search.Search(o, out)
	}

	src, err := harness.ResolveSource(*repoFlag, home, buildVersion())
	if err != nil {
		return err
	}
	// repo is the tree on disk, "" for the files embedded in a release binary;
	// no git helper is ever called with "" (D17).
	repo := src.Dir

	switch cmd {
	case "install", "sync":
		vendor, err := resolveVendor(*vendorFlag, home)
		if err != nil {
			return err
		}
		if repo == "" {
			// A release binary installs the harness built into it and itself:
			// there is no source to build from (D13).
			if *dryRun {
				fmt.Fprintf(out, "[dry-run] release binary %s: the harness built into it and the running binary would be installed\n", src.Version)
			} else {
				fmt.Fprintf(out, "release binary %s: installing the harness built into it and the running binary\n", src.Version)
			}
			return installFiles(src, home, vendor, *dryRun, out)
		}
		if err := knownCheckout(*repoFlag, home, repo); err != nil {
			return err
		}
		// The deployed binary must be the repo's own build: rebuild and let
		// that build install, unless this process is that build (ChildEnv),
		// there is no toolchain to build with, the tree is a plain copy, or
		// this is a dry run — a dry run never builds or runs a fresh binary;
		// like every install it still reads the checkout's git state, which
		// is why the checkout must be one the operator named or installed
		// before. A plain copy (no .git: an unpacked bundle)
		// never rebuilds — whoever unpacked it built the running binary from
		// that same tree, the tree is often read-only to the installing user,
		// and git cannot update it anyway (BUG-015).
		_, goErr := exec.LookPath("go")
		own := harness.IsCheckout(repo)
		rebuild := shouldRebuild(goErr == nil, os.Getenv(harness.ChildEnv) == "1", own)
		if *dryRun && rebuild {
			fmt.Fprintf(out, "[dry-run] the binary would be rebuilt from %s\n", repo)
		} else if rebuild {
			var pass []string
			if *homeFlag != "" {
				pass = append(pass, "--home", home)
			}
			if *vendorFlag != "" {
				pass = append(pass, "--vendor", vendor)
			}
			return harness.Reinstall(repo, pass, out)
		}
		if !own && *dryRun {
			fmt.Fprintln(out, "[dry-run] plain copy, not its own git checkout: the running binary would be installed as is")
		} else if !own {
			fmt.Fprintln(out, "plain copy, not its own git checkout: the running binary is installed as is")
		} else if goErr != nil {
			fmt.Fprintln(out, "warning: go is not on PATH; the running binary was copied, not rebuilt — install Go and run `sh-agent update`")
		}
		return installFiles(src, home, vendor, *dryRun, out)

	case "status":
		return status(src, home, *vendorFlag, out)

	case "diff":
		vendor, err := resolveVendor(*vendorFlag, home)
		if err != nil {
			return err
		}
		// A zip install is not compared: the zip is not kept, and the source
		// this binary resolves (its embedded files, a checkout) is not what was
		// installed (ADR-038).
		if st, found, err := harness.LoadState(home); err == nil && found && strings.HasPrefix(st.Repo, harness.ZipPrefix) {
			fmt.Fprintln(out, "diff       not compared — this home was installed from a zip that is not kept; run update --from-zip <file> again to change it")
			return nil
		}
		files, err := harness.Plan(src, nil, vendor)
		if err != nil {
			return err
		}
		return harness.Diff(home, files, out)

	case "update":
		// A dry run of update is its --check: it reports, and changes nothing —
		// no download, no binary replaced, no kept clone shed, no pull. It
		// says what a real update would do instead.
		dry := *dryRun
		if dry {
			*check = true
		}
		if repo == "" {
			// A release binary has no source to build from: the newest
			// release is downloaded and verified instead (Phase 4).
			args, err := installPassArgs(*homeFlag, home, *vendorFlag)
			if err != nil {
				return err
			}
			// Both kinds of place shed the clone earlier versions kept
			// (plan § Old kept clone cleanup, D9). Only the checkout path
			// did, so a binary-only place went on holding every other
			// place's reports — which is the complaint this whole feature
			// began with. Never under --check, which changes nothing
			// (BUG-070). Found on the work Mac, 2026-09-18.
			if !*check {
				shedKeptClone(home, out)
			}
			// The BUG-079 note, when this run wrote the installed binary, is
			// printed by UpdateFromRelease itself: only there is "replaced"
			// known, and only there can its positive case be tested.
			err = harness.UpdateFromRelease(home, src.Version, *check, args, out)
			if errors.Is(err, harness.ErrNewerRelease) {
				if dry {
					fmt.Fprintln(out, "[dry-run] update would download that release, verify its signed checksums, replace the binary and reinstall")
					return nil
				}
				return exitCode(2)
			}
			if err != nil {
				return fmt.Errorf("update: %w", err)
			}
			return nil
		}
		if err := knownCheckout(*repoFlag, home, repo); err != nil {
			return err
		}
		err := update(repo, home, *check, *homeFlag, *vendorFlag, out)
		var code exitCode
		if dry && errors.As(err, &code) && code == 2 {
			fmt.Fprintln(out, "[dry-run] update would pull those commits, rebuild and reinstall")
			return nil
		}
		return err

	case "list":
		inv, err := harness.LoadInventory(src)
		if err != nil {
			return err
		}
		fmt.Fprint(out, harness.RenderInventory(inv))
		return nil

	case "render":
		inv, err := harness.LoadInventory(src)
		if err != nil {
			return err
		}
		data, err := harness.RenderContext(src, inv)
		if err != nil {
			return err
		}
		_, err = out.Write(data)
		return err

	case "adopt":
		if fs.NArg() != 1 {
			return errors.New("usage: sh-agent adopt <skill-name-or-path> [--force]")
		}
		if repo == "" {
			return errors.New("adopt needs a harness checkout (`git clone git@github.com:kesokaj/styrman-harness.git`, then run adopt from a build of it); this install is a release binary and carries no tree to adopt into")
		}
		dest, err := harness.Adopt(repo, home, fs.Arg(0), *force)
		if err != nil {
			return err
		}
		next := "commit and push the repo"
		if !harness.IsCheckout(repo) {
			// A plain copy has no repository: the adopted skill lives only
			// here until it is copied into a checkout (audit 2026-09-18).
			next = "copy it into your harness checkout and commit there — this tree is a plain copy with no repository"
		}
		fmt.Fprintf(out, "adopted into %s\nnext: review it, run `sh-agent install`, then %s\n", dest, next)
		return nil
	}
	return fmt.Errorf("unknown command %q\n\n%s", cmd, usage)
}

// installFiles plans the source with the running binary and installs it,
// then the PATH and attribution notes every install prints.
func installFiles(src harness.Source, home, vendor string, dryRun bool, out io.Writer) error {
	self, err := selfExecutable()
	if err != nil {
		return err
	}
	files, err := harness.Plan(src, self, vendor)
	if err != nil {
		return err
	}
	if err := harness.Install(src, home, files, vendor, dryRun, out); err != nil {
		return err
	}
	// A dry run wrote no binary, so no statement about how PATH reaches it is
	// true; the hint would call the bin dir "not on PATH" whether it is or not
	// (security review round 3).
	if !dryRun {
		if err := pathHint(home, out); err != nil {
			return err
		}
	}
	if vendor == harness.VendorAll || vendor == harness.VendorClaude {
		attributionSet(home, dryRun, out)
		return attributionHint(home, "note: ", out)
	}
	return nil
}

func report(repoFlag, feedbackFlag, home string, dryRun bool, out io.Writer) error {
	// `report` READS a source, it never installs one, so a release binary
	// takes its own files directly rather than going through the refusal
	// ResolveSource owes `install`: an exported $SH_AGENT_REPO would otherwise
	// send a record with harness_version, drift, skills and agents all null —
	// a machine reporting nothing useful, which is the shape BUG-077 was about
	// (security review 2026-09-18).
	var src harness.Source
	if v := buildVersion(); harness.IsReleaseVersion(v) {
		src = harness.EmbeddedSource(v)
		if repoFlag != "" {
			fmt.Fprintf(out, "report: --repo is not used here; this release binary reports the harness built into it\n")
		}
	} else if resolved, err := harness.ResolveSource(repoFlag, home, buildVersion()); err != nil {
		fmt.Fprintf(out, "report: %v\n", err)
	} else {
		src = resolved
	}
	host, err := os.Hostname()
	if err != nil {
		host = ""
	}
	if err := harness.SendReport(harness.ReportRun{
		Home: home, Source: src, Feedback: feedbackFlag, URL: harness.FeedbackURL,
		Host: host, Now: time.Now(), DryRun: dryRun,
	}, out); err != nil {
		return fmt.Errorf("report: %w", err)
	}
	return nil
}

// stdin is where `memory add` reads the body; a test replaces it.
var stdin io.Reader = os.Stdin

// memoryCmd runs `memory add` and `memory show` in the checkout holding the
// working directory. A refusal prints its reason and exits 2.
func memoryCmd(args []string, e memory.Entry, home, vendorFlag string, out io.Writer) error {
	if len(args) != 1 || (args[0] != "add" && args[0] != "show") {
		return errors.New("usage: sh-agent memory add --type T --name SLUG --description TEXT (body on stdin) | sh-agent memory show")
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	var text string
	if args[0] == "add" {
		text, err = memory.Add(wd, e, stdin, time.Now())
		if err == nil {
			text = "saved " + text + "\n"
		}
	} else {
		vendor, verr := resolveVendor(vendorFlag, home)
		if verr != nil {
			return verr
		}
		text, err = memory.Show(wd, vendor != harness.VendorAntigravity)
	}
	if memory.IsRefusal(err) {
		fmt.Fprintln(out, "sh-agent memory:", err) // on the caller's writer (audit 2026-09-30, IN13)
		return exitCode(2)
	}
	if err != nil {
		return err
	}
	fmt.Fprint(out, text)
	return nil
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

func status(src harness.Source, home, vendorFlag string, out io.Writer) error {
	repo := src.Dir
	// Read first: the source line says when this home came from a zip, which
	// the embedded files or a checkout beside it do not describe (ADR-038).
	st, found, err := harness.LoadState(home)
	if err != nil {
		return err
	}
	zipSum, fromZip := strings.CutPrefix(st.Repo, harness.ZipPrefix+"sha256:")
	fromZip = found && fromZip
	zipNote := ""
	if fromZip {
		zipNote = "; this home was installed from a zip — run update --from-zip <file> again to change it"
	}
	var info harness.RepoInfo
	switch {
	case repo == "":
		// The embedded files: no git state to show, and none is asked (D17).
		fmt.Fprintf(out, "source     embedded %s (release binary; update downloads the newest release%s)\n", src.Version, zipNote)
	default:
		info = harness.Inspect(repo)
		dirty := "clean"
		if info.StatusUnknown {
			// Never render an unreadable tree as clean: git failing to report is
			// not evidence of anything (round-2 release review).
			dirty = "state unknown — git could not report it"
		} else if info.Dirty {
			dirty = "uncommitted changes"
		}
		if info.Commit == "" {
			fmt.Fprintf(out, "repo       %s (plain copy, not its own git checkout; update is refused%s)\n", repo, zipNote)
		} else {
			fmt.Fprintf(out, "repo       %s (%s @ %s, %s%s)\n", repo, info.Branch, info.Commit, dirty, zipNote)
		}
	}

	if !found {
		fmt.Fprintln(out, "installed  nothing — run `sh-agent install`")
		return exitCode(1)
	}
	vDisplay := st.Vendor
	if vDisplay == "" {
		vDisplay = harness.VendorAll
	}
	from := "commit " + st.Commit
	if release, ok := strings.CutPrefix(st.Repo, harness.EmbeddedPrefix); ok {
		from = "release " + release
	} else if fromZip {
		// A manifest from before the verdict was recorded has no ";…" part.
		sum, trust, _ := strings.Cut(zipSum, ";")
		from = "zip sha256:" + short(sum)
		if trust != "" {
			from += " (" + trust + ")"
		}
	}
	fmt.Fprintf(out, "installed  %s from %s (vendor: %s, %d files, manifest ~/.config/sh-agent/state.json)\n", st.InstalledAt, from, vDisplay, len(st.Files))
	// The line describes the installed file, not this process: status may
	// run from `go run` or from a build outside ~/.local/bin.
	binaryLine := "binary     unknown (installed file has no build info)"
	if bi, err := buildinfo.ReadFile(filepath.Join(home, harness.BinPath)); err == nil {
		rev := versionFromInfo(bi)
		if repo == "" {
			rev = releaseFromInfo(bi)
		}
		binaryLine = "binary     " + rev
		switch {
		case repo == "":
			// No checkout commit to compare against.
		case !binaryOlder(rev, info.Commit):
		case rev == "(devel)" || rev == "unknown":
			binaryLine += " — no build stamp; run `sh-agent update`"
		default:
			binaryLine += fmt.Sprintf(" — older than the repo (%s); run `sh-agent update`", info.Commit)
		}
	}
	fmt.Fprintln(out, binaryLine)
	// The line above describes the installed file; these say when that file is
	// not the one an unqualified `sh-agent` runs (BUG-078). Only when PATH
	// reaches the installed file at all: for a scratch home it never does, and
	// the operator's own home is the one where the distinction bites.
	if shadows, onPath := harness.PathShadows(os.Getenv("PATH"), filepath.Join(home, harness.BinPath)); onPath {
		for _, s := range shadows {
			fmt.Fprintf(out, "shadowed   %s comes first on PATH, so `sh-agent` runs that, not the binary above\n", s)
		}
	}
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
	fmt.Fprintln(out, "dream      "+dream.StatusLine(home))

	vendor, err := resolveVendor(vendorFlag, home)
	if err != nil {
		return err
	}

	inv, err := harness.LoadInventory(src)
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

	if fromZip {
		// The zip is not kept, so its content cannot be compared; say so and
		// check what can be checked — every file the install wrote is there.
		var missing []string
		for _, rel := range st.Files {
			if _, err := os.Lstat(filepath.Join(home, rel)); err != nil {
				missing = append(missing, rel)
			}
		}
		fmt.Fprintf(out, "drift      not compared — installed from a zip that is not kept; %d of %d manifest files present\n", len(st.Files)-len(missing), len(st.Files))
		for _, rel := range missing {
			fmt.Fprintf(out, "  missing  ~/%s\n", rel)
		}
		if len(missing) > 0 {
			return exitCode(1)
		}
		return nil
	}
	files, err := harness.Plan(src, nil, vendor)
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
		fmt.Fprintln(out, "note       ~/.claude/CLAUDE.md is yours; move only rules meant for every machine into a rule file under rules/ in the harness repo — rules/ ships to every install on every machine, never credentials or personal or host-specific notes — then remove it and run install to write the one-line import")
	}
	if vendor == harness.VendorAll || vendor == harness.VendorClaude {
		if err := attributionHint(home, "note       ", out); err != nil {
			return err
		}
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
	info := harness.Inspect(repo)
	if behind == 0 {
		fmt.Fprintf(out, "up to date with upstream (%s)\n", info.Commit)
	} else {
		fmt.Fprintf(out, "behind upstream by %d commit(s)\n", behind)
	}
	// A dirty checkout is never installed, on either path. `update` deploys
	// the WORKING TREE, not the commit it just fetched, so uncommitted work
	// reaches every vendor directory — which is R6 defeated (BUG-046).
	// `install` stays the deliberate way to deploy a dirty tree. The check
	// sits above the --check and --dry-run returns: a read-only mode that
	// promised a pull the real run refuses was no check of it (audit
	// 2026-09-30, IN4).
	if info.Dirty || info.StatusUnknown {
		what := "has uncommitted changes"
		if info.StatusUnknown {
			what = "could not be checked for uncommitted changes (git did not report its status)"
		}
		return fmt.Errorf("%s %s and `update` installs the working tree, not the fetched commit; commit or stash them, or run `sh-agent install` to deploy them deliberately", repo, what)
	}
	if checkOnly {
		if behind == 0 {
			return nil
		}
		return exitCode(2)
	}
	// Below every `--check` return: that mode reports and changes nothing
	// (security review 2026-09-17, finding 2).
	shedKeptClone(home, out)
	if behind > 0 {
		if err := harness.Pull(repo, out); err != nil {
			return err
		}
		fmt.Fprintf(out, "pulled to %s\n", harness.Inspect(repo).Commit) // the tree moved: inspected again
	}
	args, err := installPassArgs(homeFlag, home, vendorFlag)
	if err != nil {
		return err
	}
	if err := harness.Reinstall(repo, args, out); err != nil {
		return err
	}
	// Reinstall wrote ~/.local/bin/sh-agent; the process running this update
	// may be another file (BUG-079).
	harness.NoteRunningElsewhere(home, out)
	return nil
}

// shedKeptClone removes the full clone of the feedback repository that earlier
// versions kept in this home, so no machine holds another machine's reports
// (D9). `report` keeps nothing now; update is where the old one goes. Reported,
// never fatal: it is not what update came to do.
func shedKeptClone(home string, out io.Writer) {
	if err := harness.RemoveKeptClone(home, harness.FeedbackURL, out); err != nil {
		fmt.Fprintf(out, "note: could not remove the kept feedback clone: %v\n", err)
	}
}

// installPassArgs are the flags an update hands to the install it runs: the
// home it was told to use, and a vendor narrower than "all".
func installPassArgs(homeFlag, home, vendorFlag string) ([]string, error) {
	var args []string
	if homeFlag != "" {
		args = append(args, "--home", home)
	}
	vendor, err := resolveVendor(vendorFlag, home)
	if err != nil {
		return nil, err
	}
	if vendor != harness.VendorAll {
		args = append(args, "--vendor", vendor)
	}
	return args, nil
}

// knownCheckout refuses to build and install from a checkout that was only
// discovered by walking up from the working directory when a different
// checkout is already installed. Building runs that tree's code and deploys
// its context everywhere; standing in a harness-shaped directory (an unpacked
// copy, a clone under review) must not be enough. An explicit --repo or
// $SH_AGENT_REPO, or a first install, is.
func knownCheckout(repoFlag, home, repo string) error {
	if repoFlag != "" || os.Getenv("SH_AGENT_REPO") != "" {
		return nil
	}
	st, found, err := harness.LoadState(home)
	if err != nil {
		return fmt.Errorf("reading the manifest before building: %w", err)
	}
	if !found || st.Repo == "" {
		return nil
	}
	// A manifest written by an embedded install records a release, never a
	// path: it is not handed to samePath, which would resolve it against the
	// working directory (security review 2026-09-17, finding 2).
	if release, ok := strings.CutPrefix(st.Repo, harness.EmbeddedPrefix); ok {
		return fmt.Errorf("this home was installed from the release binary %s, not from a checkout; pass --repo %s to build and install from it", release, repo)
	}
	if sum, ok := strings.CutPrefix(st.Repo, harness.ZipPrefix+"sha256:"); ok {
		sum, _, _ = strings.Cut(sum, ";")
		return fmt.Errorf("this home was installed from a zip (sha256 %s), not from a checkout; pass --repo %s to build and install from it", short(sum), repo)
	}
	if samePath(st.Repo, repo) {
		return nil
	}
	return fmt.Errorf("checkout %s is not the installed one (%s); pass --repo %s to build and install from it", repo, st.Repo, repo)
}

func samePath(a, b string) bool {
	ra, errA := filepath.EvalSymlinks(a)
	rb, errB := filepath.EvalSymlinks(b)
	if errA != nil || errB != nil {
		return a == b
	}
	return ra == rb
}

// shouldRebuild decides whether install re-executes through a fresh build
// of the repo: yes whenever a toolchain is present and the repo is its own
// git checkout, unless this process is already the fresh build that
// Reinstall spawned. A plain copy is installed from the running binary.
func shouldRebuild(goOnPath, child, ownCheckout bool) bool {
	return goOnPath && !child && ownCheckout
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
	// Every flag that TAKES A VALUE belongs here, or it swallows the next
	// argument: `report --feedback-repo /x --dry-run` parses --dry-run as the
	// path and pushes when a dry run was asked (BUG-002's shape, caught by
	// the cold read of the delivery plan before it shipped).
	takesValue := map[string]bool{
		"-repo": true, "--repo": true,
		"-home": true, "--home": true,
		"-vendor": true, "--vendor": true,
		"-feedback-repo": true, "--feedback-repo": true,
		"-from-zip": true, "--from-zip": true,
		"-limit": true, "--limit": true,
		"-kind": true, "--kind": true,
		"-refs": true, "--refs": true,
		"-type": true, "--type": true,
		"-name": true, "--name": true,
		"-description": true, "--description": true,
		"-project": true, "--project": true,
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

// pathHint says, after an install, whether the next unqualified `sh-agent`
// runs the binary just written — one message per situation. PATH reaching the
// file but only after an earlier entry is the BUG-078 case: install reports
// success and `status` reports drift none, both about a binary the next
// command will not run. PATH not reaching it is the old note, naming what
// answers instead when something does. Every fact comes from one PathShadows
// call, so the two halves cannot disagree on what "on PATH" means.
func pathHint(home string, out io.Writer) error {
	bin := filepath.Join(home, harness.BinPath)
	binDir := filepath.Dir(bin)
	shadows, onPath := harness.PathShadows(os.Getenv("PATH"), bin)
	switch {
	case onPath:
		for _, s := range shadows {
			fmt.Fprintf(out, "warning: %s comes before %s on PATH, so `sh-agent` runs that one\n", s, bin)
		}
	case len(shadows) > 0:
		fmt.Fprintf(out, "note: %s is not on PATH, so `sh-agent` runs %s instead; put %s ahead of it\n", binDir, shadows[0], binDir)
	default:
		fmt.Fprintf(out, "note: %s is not on PATH; add it so `sh-agent` is available everywhere\n", binDir)
	}
	return nil
}

// attributionSet writes the one vendor setting the harness owns (R4,
// ADR-025): Claude Code's attribution off, so no commit or PR carries the
// trailer the context's Git rule forbids. Only an absent key is written — a
// present value of any kind is the operator's and attributionHint says when
// it is on. The write is a textual insertion before the object's closing
// brace, re-parsed before it lands, because encoding/json would reorder the
// operator's keys; every byte up to the object's last value stays, and the
// file keeps its own line endings. No state of the file
// fails the install: the deployment already succeeded and this file is
// advisory. dryRun prints what would be written and touches nothing.
func attributionSet(home string, dryRun bool, out io.Writer) {
	const fragment = `"attribution": {"commit": "", "pr": ""}`
	path := filepath.Join(home, harness.ClaudeRoot, "settings.json")
	// Lstat, not Stat: a settings.json that is a symlink (a dotfiles checkout)
	// is never written through — a rename would replace the link with a file.
	info, err := os.Lstat(path)
	switch {
	case err == nil && info.Mode()&os.ModeSymlink != 0:
		fmt.Fprintf(out, "settings: %s is a symlink; add the key to its target by hand\n", path)
		return
	case errors.Is(err, os.ErrNotExist):
		if dryRun {
			fmt.Fprintf(out, "[dry-run] would create %s with %s\n", path, fragment)
			return
		}
		content := []byte("{\n  " + fragment + "\n}\n")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			settingsUnwritable(out, path, err)
			return
		}
		// O_EXCL: the Lstat above is advisory; whatever appears at the path
		// between the two calls (a planted link, a FIFO) makes this fail
		// instead of being written through or blocked on.
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			settingsUnwritable(out, path, err)
			return
		}
		if _, err := f.Write(content); err != nil {
			f.Close()
			settingsUnwritable(out, path, err)
			return
		}
		if err := f.Sync(); err != nil {
			f.Close()
			settingsUnwritable(out, path, err)
			return
		}
		if err := f.Close(); err != nil {
			settingsUnwritable(out, path, err)
			return
		}
		fmt.Fprintf(out, "settings: %s created with %s — Claude Code's commit and PR trailers are off\n", path, fragment)
		return
	case err != nil, !info.Mode().IsRegular():
		return // attributionHint names it, and readSettings bounds the size
	}
	data, problem := readSettings(path)
	if problem != "" {
		return // attributionHint names it
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil || top == nil {
		// malformed, not an object, or the bare value null (a nil map with no
		// error): attributionHint names each with one line
		return
	}
	if _, present := top["attribution"]; present {
		return
	}
	// Insert after the last byte of the object's body, keeping the operator's
	// indentation style (two spaces, the file's own newline before the brace)
	// and whatever follows the closing brace.
	trimmed := bytes.TrimRight(data, " \t\r\n")
	closing := bytes.LastIndexByte(trimmed, '}')
	opening := bytes.IndexByte(trimmed, '{')
	if closing <= opening {
		fmt.Fprintf(out, "settings: %s could not be edited safely (the result would not parse); add the key by hand\n", path)
		return
	}
	nl := "\n"
	if bytes.Contains(data, []byte("\r\n")) {
		nl = "\r\n" // keep the file's own line endings
	}
	body := bytes.TrimRight(trimmed[:closing], " \t\r\n")
	insert := "," + nl + "  " + fragment + nl
	if len(bytes.TrimSpace(trimmed[opening+1:closing])) == 0 {
		body = trimmed[:opening+1]
		insert = nl + "  " + fragment + nl
	}
	result := make([]byte, 0, len(data)+len(insert))
	result = append(result, body...)
	result = append(result, insert...)
	result = append(result, trimmed[closing:]...)
	result = append(result, data[len(trimmed):]...)
	if !json.Valid(result) {
		fmt.Fprintf(out, "settings: %s could not be edited safely (the result would not parse); add the key by hand\n", path)
		return
	}
	if dryRun {
		fmt.Fprintf(out, "[dry-run] would add %s to %s\n", fragment, path)
		return
	}
	if err := harness.WriteAtomic(path, result, info.Mode().Perm()); err != nil {
		settingsUnwritable(out, path, err)
		return
	}
	fmt.Fprintf(out, "settings: %s — %s added; other keys untouched\n", path, fragment)
}

// settingsUnwritable is the one sentence for every way the settings file
// could not be written; five call sites said it five times.
func settingsUnwritable(out io.Writer, path string, err error) {
	fmt.Fprintf(out, "settings: %s could not be written (%v); add the key by hand\n", path, err)
}

// readSettings reads the operator's settings file, bounded to a regular file
// of at most one megabyte; problem names why it could not be read, without a
// byte of its contents.
func readSettings(path string) (data []byte, problem string) {
	info, err := os.Stat(path)
	switch {
	case err != nil:
		return nil, "could not be read"
	case !info.Mode().IsRegular():
		return nil, "is not a regular file"
	case info.Size() > 1<<20:
		// One bound, here, so install and status say the same thing about
		// a file this tool will not rewrite (audit 2026-09-18).
		return nil, "is larger than 1 MiB, more than this tool will read or rewrite"
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, "could not be read"
	}
	defer f.Close()
	data, err = io.ReadAll(io.LimitReader(f, 1<<20))
	if err != nil {
		return nil, "could not be read"
	}
	return data, ""
}

// attributionHint says when Claude Code will still add its Co-Authored-By
// trailer. The context's Git rule forbids the trailer, but in the sandbox the
// rule alone did not hold (three of three headless commits carried it); the
// CLI's own setting is what stops it, and settings.json is the operator's
// file, so install and status only say so. prefix is each command's note
// style ("note: " for install, "note       " for status). Advisory only: no
// state of that file may fail install or status, and nothing from it is
// echoed.
func attributionHint(home, prefix string, out io.Writer) error {
	path := filepath.Join(home, harness.ClaudeRoot, "settings.json")
	off, problem := attributionOff(path)
	switch {
	case problem != "":
		fmt.Fprintf(out, "%s%s %s, so the attribution setting could not be checked; Claude Code adds a Co-Authored-By trailer to every commit unless it sets \"attribution\": {\"commit\": \"\"} (context § Git)\n", prefix, path, problem)
	case !off:
		fmt.Fprintf(out, "%s%s does not set \"attribution\": {\"commit\": \"\"}; Claude Code adds a Co-Authored-By trailer to every commit until it does (context § Git)\n", prefix, path)
	}
	return nil
}

// attributionOff reads the operator's settings file, bounded to a regular
// file of at most one megabyte. off is true when attribution.commit is "" or
// false; problem names why the file could not be judged, without a byte of
// its contents (a JSON error string carries the input it choked on).
func attributionOff(path string) (off bool, problem string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false, ""
	}
	data, problem := readSettings(path)
	if problem != "" {
		return false, problem
	}
	var settings struct {
		Attribution struct {
			Commit any `json:"commit"`
		} `json:"attribution"`
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		if json.Valid(data) {
			return false, "is not a JSON object"
		}
		return false, "is not valid JSON"
	}
	switch v := settings.Attribution.Commit.(type) {
	case string:
		return v == "", ""
	case bool:
		return !v, ""
	}
	return false, ""
}

// buildVersion describes the running process, for `sh-agent version`.
// `sh-agent status` reads the installed file instead (see status).
func buildVersion() string {
	if version != "" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	return versionFromInfo(info)
}

// releaseFromInfo renders an installed release binary's version: the module
// version when it is a release tag, which Go stamps for a build at a clean
// tag and `go install …@vX.Y.Z` records — a release build's -ldflags are not
// in the build info, since -trimpath leaves them out — and the VCS stamp
// otherwise.
func releaseFromInfo(info *debug.BuildInfo) string {
	if harness.IsReleaseVersion(info.Main.Version) {
		return info.Main.Version
	}
	return versionFromInfo(info)
}

// versionFromInfo renders a binary's VCS stamp as the short revision with a
// "+dirty" suffix for a modified tree, or the module version when the build
// carried no stamp (`go run`, `go test`, a build outside git).
func versionFromInfo(info *debug.BuildInfo) string {
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

// updateFromZip installs the harness in the repository zip the operator
// downloaded, with the running binary: unpack, install, delete the unpacked
// copy (R8, R9's zip place; ADR-038). Nothing is fetched, and nothing in the
// zip replaces the binary: the running one is installed as in every install.
func updateFromZip(zipPath, home, vendorFlag string, dryRun bool, out io.Writer) (err error) {
	vendor, err := resolveVendor(vendorFlag, home)
	if err != nil {
		return err
	}
	// A dry run writes nothing, so it never creates ~/.config/sh-agent: when
	// that directory is absent it unpacks under the system temp directory
	// instead (plan D10).
	parent := filepath.Join(home, ".config", "sh-agent")
	if dryRun {
		if _, err := os.Stat(parent); err != nil {
			parent = ""
		}
	} else {
		// Created here, and removed again if anything below fails — the unpack
		// or the install that can still refuse the zip's content — so a
		// refused zip leaves the home as it was (BUG-093).
		// `=`, not `:=`: the deferred check must see the function's result, not
		// a block-scoped copy of it.
		var created []string
		if created, err = harness.MkdirAllTracked(parent); err != nil {
			return err
		}
		defer func() {
			if err != nil {
				harness.RemoveCreated(created)
			}
		}()
	}
	// The zip place is the fallback for machines that cannot take the normal
	// way in, so it accepts the zip it is given (ADR-041). A release zip with
	// its SHA256SUMS and SHA256SUMS.sig beside it is verified before a byte is
	// written, and refused if that check fails; without the pair it is
	// installed marked UNVERIFIED (R9).
	verifiedTag := ""
	check := func(sum string) error {
		tag, err := harness.VerifyZipSums(zipPath, sum)
		if errors.Is(err, harness.ErrZipUnsigned) {
			return nil
		}
		if err != nil {
			return err
		}
		verifiedTag = tag
		return nil
	}
	dir, sum, err := harness.UnpackZip(zipPath, parent, check)
	if err != nil {
		return err
	}
	defer func() {
		if rerr := os.RemoveAll(dir); rerr != nil {
			fmt.Fprintf(out, "warning: could not remove the unpacked copy %s: %v\n", dir, rerr)
		}
	}()
	prefix := ""
	if dryRun {
		prefix = "[dry-run] "
	}
	trust := "UNVERIFIED — no signed SHA256SUMS beside it; trust rests on your own download"
	if verifiedTag != "" {
		trust = "VERIFIED against the signed SHA256SUMS of release " + verifiedTag
	}
	fmt.Fprintf(out, "%szip %s (sha256 %s, %s): installing its harness and the running binary %s; nothing is downloaded and nothing in the zip replaces the binary\n", prefix, zipPath, sum, trust, buildVersion())
	// The manifest keeps the verdict too, so status can tell a verified install
	// from an unverified one after the output has scrolled away.
	origin := harness.ZipPrefix + "sha256:" + sum + ";unverified"
	if verifiedTag != "" {
		origin = harness.ZipPrefix + "sha256:" + sum + ";verified " + verifiedTag
	}
	src := harness.Source{FS: os.DirFS(dir), Dir: dir, Origin: origin}
	if err := installFiles(src, home, vendor, dryRun, out); err != nil {
		return err
	}
	// Only after a successful install: a refused or dry run changes nothing,
	// not even a stale leftover from an interrupted run.
	if !dryRun {
		harness.SweepStaleDirs(parent, "zip-", time.Hour, out)
	}
	return nil
}

// short is the first 12 hex digits of a sha256, enough to tell zips apart.
func short(sum string) string {
	if len(sum) > 12 {
		return sum[:12]
	}
	return sum
}
