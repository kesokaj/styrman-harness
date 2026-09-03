package install

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// finding is one thing doctor noticed.
type finding struct {
	ok      bool
	subject string
	detail  string
}

// Doctor reports drift between this repo and the deployed harness.
//
// It exists because the harness spent an unknown stretch running a
// configuration that had never been shipped — a stale constitution, two hook
// scripts that did not exist on disk, and an Antigravity plugin that was
// copied but never registered. None of that announced itself. Every check here
// is one of those failures made loud.
func Doctor(args []string) error {
	fset := flag.NewFlagSet("doctor", flag.ContinueOnError)
	home := fset.String("home", "", "check under this directory instead of $HOME")
	quiet := fset.Bool("quiet", false, "print only failures")
	if err := fset.Parse(args); err != nil {
		return err
	}
	env, err := newEnv(*home)
	if err != nil {
		return err
	}

	var findings []finding
	findings = append(findings, checkDeployed(env)...)
	findings = append(findings, checkHookCommands(env))
	findings = append(findings, checkAgyRegistration(env))
	findings = append(findings, checkAgyPlugin(env))

	failed := 0
	for _, f := range findings {
		if !f.ok {
			failed++
		}
		if f.ok && *quiet {
			continue
		}
		mark := "ok  "
		if !f.ok {
			mark = "DRIFT"
		}
		fmt.Printf("%-5s %-34s %s\n", mark, f.subject, f.detail)
	}
	if failed > 0 {
		return fmt.Errorf("%d of %d checks failed — run `harness install`", failed, len(findings))
	}
	if !*quiet {
		fmt.Printf("all %d checks passed\n", len(findings))
	}
	return nil
}

// checkDeployed compares every managed file's content against the repo.
func checkDeployed(env Env) []finding {
	var out []finding
	for _, t := range []Target{TargetClaude, TargetAntigravity} {
		ops, err := Plan(env, t)
		if err != nil {
			out = append(out, finding{false, string(t), err.Error()})
			continue
		}
		for _, op := range ops {
			subject := trimHome(env.Home, op.Dst)
			if op.Mirror {
				diff, err := dirDiff(op.Src, op.Dst)
				switch {
				case err != nil:
					out = append(out, finding{false, subject, err.Error()})
				case len(diff) > 0:
					out = append(out, finding{false, subject,
						fmt.Sprintf("%d file(s) differ or missing, e.g. %s", len(diff), strings.Join(first(diff, 3), ", "))})
				default:
					out = append(out, finding{true, subject, "in sync"})
				}
				continue
			}
			same, err := sameContent(op.Src, op.Dst)
			switch {
			case err != nil:
				out = append(out, finding{false, subject, err.Error()})
			case !same:
				out = append(out, finding{false, subject, "differs from " + rel(env.Repo, op.Src)})
			default:
				out = append(out, finding{true, subject, "in sync"})
			}
		}
	}
	return out
}

// checkHookCommands verifies every hook command in the deployed settings can
// actually run. A hook naming a missing script fails on every tool call.
func checkHookCommands(env Env) finding {
	path := filepath.Join(env.Home, ".claude", "settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return finding{false, "hook commands", "cannot read " + trimHome(env.Home, path)}
	}
	var doc any
	if err := json.Unmarshal(data, &doc); err != nil {
		return finding{false, "hook commands", "settings.json is not valid JSON: " + err.Error()}
	}

	var missing []string
	for _, cmd := range collectCommands(doc) {
		bin := executablePath(env.Home, cmd)
		if bin == "" {
			continue // a shell builtin or a bare name resolved via PATH
		}
		fi, err := os.Stat(bin)
		if err != nil || fi.IsDir() || fi.Mode().Perm()&0o111 == 0 {
			missing = append(missing, trimHome(env.Home, bin))
		}
	}
	if len(missing) > 0 {
		return finding{false, "hook commands", "missing or not executable: " + strings.Join(missing, ", ")}
	}
	return finding{true, "hook commands", "all resolve and are executable"}
}

// collectCommands walks arbitrary settings JSON for hook "command" strings.
func collectCommands(node any) []string {
	var out []string
	switch v := node.(type) {
	case map[string]any:
		for key, val := range v {
			if key == "command" {
				if s, ok := val.(string); ok {
					out = append(out, s)
				}
				continue
			}
			out = append(out, collectCommands(val)...)
		}
	case []any:
		for _, item := range v {
			out = append(out, collectCommands(item)...)
		}
	}
	return out
}

// executablePath extracts the program from a hook command line and expands the
// home shorthands both vendors accept. It returns "" when the program is a
// bare name to be resolved through PATH.
func executablePath(home, cmd string) string {
	fields := strings.Fields(cmd)
	if len(fields) == 0 {
		return ""
	}
	prog := fields[0]
	switch {
	case strings.HasPrefix(prog, "~/"):
		prog = filepath.Join(home, prog[2:])
	case strings.HasPrefix(prog, "$HOME/"):
		prog = filepath.Join(home, strings.TrimPrefix(prog, "$HOME/"))
	case !strings.Contains(prog, "/"):
		return ""
	}
	return prog
}

// checkAgyRegistration verifies the plugin is imported, not merely copied.
// agy keeps its manifest in the real home regardless of where the harness was
// deployed, so this check is meaningless against a scratch home.
func checkAgyRegistration(env Env) finding {
	if !env.isRealHome() {
		return finding{true, "agy plugin registration", "skipped — scratch home"}
	}
	path := filepath.Join(env.Home, ".gemini", "config", "import_manifest.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return finding{false, "agy plugin registration",
			"no import manifest — copied but not registered; run `agy plugin install`"}
	}
	var manifest struct {
		Imports []struct {
			Name string `json:"name"`
		} `json:"imports"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return finding{false, "agy plugin registration", "import manifest is not valid JSON"}
	}
	for _, imp := range manifest.Imports {
		if imp.Name == "harness" {
			return finding{true, "agy plugin registration", "harness imported"}
		}
	}
	return finding{false, "agy plugin registration", "harness not in the import manifest"}
}

// checkAgyPlugin asks agy itself whether the plugin is well-formed. agy offers
// no way to enumerate valid tool names, so this is the closest available check
// on the derived agents' tool lists.
func checkAgyPlugin(env Env) finding {
	if _, err := exec.LookPath("agy"); err != nil {
		return finding{true, "agy plugin validation", "skipped — agy not on PATH"}
	}
	if !env.isRealHome() {
		return finding{true, "agy plugin validation", "skipped — scratch home"}
	}
	path := filepath.Join(env.Home, ".gemini", "config", "plugins", "harness")
	out, err := exec.Command("agy", "plugin", "validate", path).CombinedOutput()
	if err != nil {
		return finding{false, "agy plugin validation", strings.TrimSpace(string(out))}
	}
	return finding{true, "agy plugin validation", "plugin validates"}
}

func sameContent(a, b string) (bool, error) {
	ha, err := hashFile(a)
	if err != nil {
		return false, fmt.Errorf("repo source unreadable: %w", err)
	}
	hb, err := hashFile(b)
	if err != nil {
		return false, nil // not deployed yet counts as drift, not an error
	}
	return ha == hb, nil
}

func hashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// dirDiff lists relative paths that differ between src and dst: missing,
// changed, or present in dst but no longer in src.
//
// The second walk is not symmetry for its own sake. Where the installer hands a
// directory to an external tool rather than mirroring it, nothing runs
// pruneExtras — so a file deleted from the repo lingers in the deployment, and
// a one-directional walk would never mention it. "Drift is loud" has to include
// drift made of things that should no longer be there.
func dirDiff(src, dst string) ([]string, error) {
	var diff []string
	inSrc := map[string]bool{}

	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		inSrc[relPath] = true
		same, err := sameContent(path, filepath.Join(dst, relPath))
		if err != nil {
			return err
		}
		if !same {
			diff = append(diff, relPath)
		}
		return nil
	})
	if err != nil {
		return diff, err
	}

	// Anything deployed that the repo no longer has.
	err = filepath.WalkDir(dst, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // a destination that cannot be read is reported by the first walk
		}
		if d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(dst, path)
		if err != nil || inSrc[relPath] {
			return nil
		}
		diff = append(diff, relPath+" (deployed but not in the repo)")
		return nil
	})
	return diff, err
}

func first(list []string, n int) []string {
	if len(list) <= n {
		return list
	}
	return list[:n]
}
