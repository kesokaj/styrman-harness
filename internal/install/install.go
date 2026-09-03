// Package install deploys the harness and reports drift between the repo and
// what is actually deployed.
//
// Everything here overwrites. The previous installer deep-merged JSON and kept
// timestamped backups, which produced two failure modes: a merge that appended
// a changed hook beside its predecessor so both fired, and an unbounded pile of
// .bak files. Overwriting has one rule and no surprises; machine-local settings
// belong in settings.local.json, which this never touches.
package install

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Target is one vendor's deployment.
type Target string

const (
	TargetClaude      Target = "claude"
	TargetAntigravity Target = "antigravity"
)

// Env is where an install reads from and writes to. Home is separated from the
// process environment so an install can be exercised against a scratch
// directory without touching the operator's real configuration.
type Env struct {
	Repo string
	Home string
}

// Op is one deployment step: a file copy or a directory mirror.
type Op struct {
	Src    string // absolute, in the repo
	Dst    string // absolute, under Home
	Mirror bool   // mirror a directory, deleting destination extras
	Mode   fs.FileMode

	// AgyOwned marks a destination that `agy plugin install` writes itself.
	// The op still describes the desired state so doctor can check it, but the
	// installer must not copy into it: agy's canonical plugin location is this
	// same path, so handing agy this directory makes it copy a tree onto
	// itself — which truncates every file to zero bytes.
	AgyOwned bool
}

// Run deploys the harness.
func Run(args []string) error {
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	dryRun := fs.Bool("dry-run", false, "print what would happen, change nothing")
	home := fs.String("home", "", "deploy under this directory instead of $HOME")
	target := fs.String("target", "", "deploy only this target (claude|antigravity)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	env, err := newEnv(*home)
	if err != nil {
		return err
	}
	targets, err := selectTargets(*target)
	if err != nil {
		return err
	}

	// The repo's own vendor copies are generated, never hand-edited: one
	// source of truth means the deployed file and the committed file cannot
	// disagree about what the constitution says.
	if err := refreshGenerated(env, *dryRun); err != nil {
		return err
	}

	binary := filepath.Join(env.Home, ".agents", "bin", "harness")
	if err := buildBinary(env, binary, *dryRun); err != nil {
		return err
	}

	for _, t := range targets {
		ops, err := Plan(env, t)
		if err != nil {
			return err
		}
		// One predicate decides both halves. Guarding the skip on "agy exists"
		// while guarding the install on "this is the real home" left a gap
		// where neither ran: a scratch-home install — the procedure AGENTS.md
		// mandates for development — deployed the plugin by no mechanism at
		// all, then failed its own doctor, and re-running never fixed it.
		agyInstalls := agyWillInstall(env)
		for _, op := range ops {
			if op.AgyOwned && agyInstalls {
				continue // agy installs this itself, from the repo
			}
			if err := apply(op, *dryRun); err != nil {
				return fmt.Errorf("%s: %w", t, err)
			}
		}
		if t == TargetAntigravity {
			if err := registerAgyPlugin(env, *dryRun); err != nil {
				return err
			}
		}
	}

	if *dryRun {
		fmt.Println("[dry-run] no changes made")
		return nil
	}
	fmt.Println("deployed; running doctor")
	return Doctor([]string{"--home", env.Home})
}

func newEnv(home string) (Env, error) {
	if home == "" {
		h, err := os.UserHomeDir()
		if err != nil {
			return Env{}, err
		}
		home = h
	}
	repo, err := repoRoot()
	if err != nil {
		return Env{}, err
	}
	return Env{Repo: repo, Home: home}, nil
}

func selectTargets(name string) ([]Target, error) {
	switch name {
	case "":
		return []Target{TargetClaude, TargetAntigravity}, nil
	case string(TargetClaude):
		return []Target{TargetClaude}, nil
	case string(TargetAntigravity):
		return []Target{TargetAntigravity}, nil
	default:
		return nil, fmt.Errorf("unknown target %q (expected claude or antigravity)", name)
	}
}

// repoRoot finds the harness repo containing this executable's working
// directory.
func repoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "shared", "AGENTS.md")); err == nil {
			return dir, nil
		}
		if parent := filepath.Dir(dir); parent == dir {
			return "", errors.New("not inside the harness repo (no shared/AGENTS.md above the working directory)")
		}
	}
}

// Plan lists the deployment steps for one target.
func Plan(env Env, t Target) ([]Op, error) {
	switch t {
	case TargetClaude:
		home := filepath.Join(env.Home, ".claude")
		return []Op{
			{Src: env.repo("claude", "CLAUDE.md"), Dst: filepath.Join(home, "CLAUDE.md"), Mode: 0o644},
			{Src: env.repo("claude", "settings.json"), Dst: filepath.Join(home, "settings.json"), Mode: 0o644},
			{Src: env.repo("shared", "skills"), Dst: filepath.Join(home, "skills"), Mirror: true},
			{Src: env.repo("claude", "agents"), Dst: filepath.Join(home, "agents"), Mirror: true},
		}, nil
	case TargetAntigravity:
		home := filepath.Join(env.Home, ".gemini")
		return []Op{
			{Src: env.repo("antigravity", "GEMINI.md"), Dst: filepath.Join(home, "GEMINI.md"), Mode: 0o644},
			{Src: env.repo("shared", "skills"), Dst: filepath.Join(home, "config", "skills"), Mirror: true},
			{Src: env.repo("antigravity", "plugins", "harness"), Dst: filepath.Join(home, "config", "plugins", "harness"), Mirror: true, AgyOwned: true},
		}, nil
	}
	return nil, fmt.Errorf("unknown target %q", t)
}

func (e Env) repo(parts ...string) string {
	return filepath.Join(append([]string{e.Repo}, parts...)...)
}

// isRealHome reports whether this install targets the operator's actual home
// rather than a scratch directory. Steps that reach outside Home — anything
// shelling out to a vendor CLI — must run only when this is true.
func (e Env) isRealHome() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	return filepath.Clean(home) == filepath.Clean(e.Home)
}

// refreshGenerated regenerates the repo's vendor copies from their single
// source: the constitution from shared/AGENTS.md, the Antigravity agents from
// the Claude ones.
func refreshGenerated(env Env, dryRun bool) error {
	src := env.repo("shared", "AGENTS.md")
	for _, dst := range []string{
		env.repo("claude", "CLAUDE.md"),
		env.repo("antigravity", "GEMINI.md"),
	} {
		if err := apply(Op{Src: src, Dst: dst, Mode: 0o644}, dryRun); err != nil {
			return err
		}
	}

	agents, err := filepath.Glob(env.repo("claude", "agents", "*.md"))
	if err != nil {
		return err
	}
	outDir := env.repo("antigravity", "plugins", "harness", "agents")
	for _, in := range agents {
		out := filepath.Join(outDir, filepath.Base(in))
		if dryRun {
			fmt.Printf("[dry-run] derive %s -> %s\n", rel(env.Repo, in), rel(env.Repo, out))
			continue
		}
		if err := DeriveFile(in, out); err != nil {
			return fmt.Errorf("derive %s: %w", rel(env.Repo, in), err)
		}
	}
	return nil
}

// buildBinary compiles the harness and installs it where the hooks expect it.
// The hooks name an absolute path, so a missing binary breaks every tool call —
// which is why a failed build fails the install rather than warning.
func buildBinary(env Env, dst string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[dry-run] go build -o %s ./cmd/harness\n", dst)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	cmd := exec.Command("go", "build", "-o", dst, "./cmd/harness")
	cmd.Dir = env.Repo
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("building harness binary: %w", err)
	}
	fmt.Println("built", dst)
	return nil
}

// registerAgyPlugin imports the plugin into Antigravity. Copying the files
// into place is not enough — agy only loads plugins listed in its own import
// manifest, which is why the harness subagents were invisible to it.
//
// agy always writes to the real home: it has no equivalent of --home, and it
// copies the plugin into its own canonical location. So a scratch-home install
// must not call it, or a test run would reach out and modify the operator's
// actual configuration.
// The source is always the repo, never the deployed directory. agy's canonical
// plugin location is ~/.gemini/config/plugins/<name>, so passing it that path
// asks it to copy a tree onto itself — which truncated every file in the plugin
// to zero bytes the first time this ran against a real home.
func registerAgyPlugin(env Env, dryRun bool) error {
	src := env.repo("antigravity", "plugins", "harness")
	dst := filepath.Join(env.Home, ".gemini", "config", "plugins", "harness")
	if dryRun {
		fmt.Printf("[dry-run] agy plugin install %s\n", src)
		return nil
	}
	if !agyWillInstall(env) {
		if !env.isRealHome() {
			fmt.Println("skipping agy plugin install: --home is a scratch directory and agy would write to the real home;",
				"the plugin was mirrored to", dst, "instead")
			return nil
		}
		fmt.Println("warning: agy not on PATH — plugin mirrored to", dst,
			"but not registered; run `agy plugin install "+src+"` once agy is installed")
		return nil
	}
	cmd := exec.Command("agy", "plugin", "install", src)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("agy plugin install: %w", err)
	}
	return nil
}

// agyWillInstall reports whether `agy plugin install` is actually going to run
// for this install. Both the decision to skip mirroring an agy-owned directory
// and the decision to invoke agy must consult exactly this — if they can
// disagree, the plugin falls through the gap between them.
func agyWillInstall(env Env) bool {
	if _, err := exec.LookPath("agy"); err != nil {
		return false
	}
	// agy has no --home and always writes to the real one, so a scratch-home
	// install must not call it.
	return env.isRealHome()
}

func apply(op Op, dryRun bool) error {
	if op.Mirror {
		if dryRun {
			fmt.Printf("[dry-run] mirror %s -> %s\n", op.Src, op.Dst)
			return nil
		}
		return mirrorDir(op.Src, op.Dst)
	}
	if dryRun {
		fmt.Printf("[dry-run] copy %s -> %s\n", op.Src, op.Dst)
		return nil
	}
	return copyFile(op.Src, op.Dst, op.Mode)
}

func copyFile(src, dst string, mode fs.FileMode) error {
	if mode == 0 {
		mode = 0o644
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// mirrorDir makes dst an exact copy of src, removing anything in dst that src
// does not have.
func mirrorDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	want := map[string]bool{}

	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}
		want[relPath] = true
		target := filepath.Join(dst, relPath)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return copyFile(path, target, info.Mode().Perm())
	})
	if err != nil {
		return err
	}
	return pruneExtras(dst, want)
}

// pruneExtras removes destination entries the source no longer has, so a
// deleted skill actually disappears instead of lingering forever.
func pruneExtras(dst string, want map[string]bool) error {
	var extras []string
	err := filepath.WalkDir(dst, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(dst, path)
		if err != nil || relPath == "." {
			return err
		}
		if !want[relPath] {
			extras = append(extras, path)
			if d.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, path := range extras {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func rel(base, path string) string {
	if r, err := filepath.Rel(base, path); err == nil {
		return r
	}
	return path
}

// trimHome shortens a path for display.
func trimHome(home, path string) string {
	if strings.HasPrefix(path, home) {
		return "~" + strings.TrimPrefix(path, home)
	}
	return path
}
