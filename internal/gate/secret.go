package gate

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/kesokaj/styrman-harness/internal/hookio"
	"github.com/kesokaj/styrman-harness/internal/trail"
)

// secretGlobs name paths that must never enter a commit. This guard replaces
// the old destructive-command guard deliberately: the agent runs inside an
// isolated microVM, so guarding the local filesystem buys nothing the VM does
// not already give. A secret does not stay in the VM — it is pushed to a
// remote, where no amount of local isolation can call it back.
var secretGlobs = []string{
	".env", ".env.*", "*.pem", "*.key", "*.p12", "*.pfx", "*.jks",
	"*.tfvars", "*.tfstate", "*.tfstate.backup",
	"id_rsa", "id_dsa", "id_ecdsa", "id_ed25519",
	"credentials.json", "service-account*.json", "*-sa.json", "gha-creds-*.json",
	".npmrc", ".pypirc", ".netrc",
}

// exampleSuffixes mark a file as a template rather than a credential. Denying
// `.env.example` would be a gate failing closed on a file whose whole purpose
// is to be committed.
var exampleSuffixes = []string{".example", ".sample", ".template", ".dist", ".tpl"}

// isSecretPath reports whether a path's basename looks like a credential.
func isSecretPath(p string) bool {
	base := path.Base(strings.Trim(strings.TrimSpace(p), `"'`))
	for _, suffix := range exampleSuffixes {
		if strings.HasSuffix(base, suffix) {
			return false
		}
	}
	for _, glob := range secretGlobs {
		if ok, _ := path.Match(glob, base); ok {
			return true
		}
	}
	return false
}

// stagingCommands are the git subcommands that can move a file toward a remote.
var stagingCommands = map[string]bool{"add": true, "commit": true, "stash": true}

// SecretDecision is the PreToolUse gate on shell commands.
//
// It answers one question: could this command move a credential toward a
// remote? An earlier version tried to decide that by parsing the command's
// arguments, and a security review defeated it fifteen ways — `git -C . add -A`,
// `git add .e*`, `git add secrets/`, `git add .en'v'`, and a secret already in
// the index that a bare `git commit` would ship. Shell and git argument syntax
// is not a thing a gate can reliably parse.
//
// So it asks git instead. For any staging command, it inspects what git
// actually holds: the index, and the credential-shaped files in the work tree
// that a stage could pick up. The command's spelling stops mattering.
func SecretDecision(p hookio.Payload) hookio.Decision {
	cmd := commandText(p)
	if cmd == "" || !mentionsStaging(cmd) {
		return hookio.Allow()
	}
	proj, ok := trail.Resolve(p.Dir())
	if !ok {
		return hookio.Allow()
	}

	// A secret already staged is the worst case: any later `git commit` ships
	// it, and that commit need not look like a staging command at all.
	if staged := secretsInIndex(proj.Root); len(staged) > 0 {
		return hookio.Deny(fmt.Sprintf(
			"Refusing: %s already staged in the index and matches a credential-shaped path. "+
				"Run `git restore --staged %s` and add it to .gitignore.",
			strings.Join(staged, ", "), staged[0]))
	}

	// A plain `git commit -m …` adds nothing to the index, so with the index
	// clean there is nothing more to check. Scanning the tree here would deny
	// an ordinary commit for an untracked .env it was never going to touch.
	if !couldAddFiles(cmd) {
		return hookio.Allow()
	}

	// Otherwise, deny only when the tree actually holds a credential this
	// command could pick up, and the command does not unambiguously name a set
	// of ordinary files. Denying a precise `git add main.go` because an
	// unrelated .env exists elsewhere would be the gate failing closed.
	addable := secretsInTree(proj.Root)
	if len(addable) == 0 || namesOnlyOrdinaryFiles(proj.Root, cmd) {
		return hookio.Allow()
	}
	return hookio.Deny(fmt.Sprintf(
		"Refusing this stage: the work tree contains credential-shaped file(s) %s that it could pick up. "+
			"Stage the paths you mean by name, or add these to .gitignore first.",
		strings.Join(addable, ", ")))
}

// commandText pulls the shell command from either vendor's payload shape.
func commandText(p hookio.Payload) string {
	if s, ok := p.ToolInput["command"].(string); ok && s != "" {
		return s
	}
	// Antigravity spells its run_command arguments differently.
	for _, key := range []string{"Command", "command", "CommandLine"} {
		if s, ok := p.ToolCall.Args[key].(string); ok && s != "" {
			return s
		}
	}
	return ""
}

// mentionsStaging reports whether a git staging subcommand appears anywhere in
// the command line. This is deliberately over-broad: a false positive costs one
// index scan, while a false negative is a leaked credential.
func mentionsStaging(cmd string) bool {
	for _, field := range strings.Fields(cmd) {
		if stagingCommands[strings.Trim(field, `"';&|`)] {
			return true
		}
	}
	return false
}

// couldAddFiles reports whether the command can put new content into the index:
// an `add` or `stash`, or a commit carrying a blanket flag. Coarse on purpose —
// it errs toward scanning, and the scan itself is what decides.
func couldAddFiles(cmd string) bool {
	for _, tok := range strings.Fields(cmd) {
		clean := strings.Trim(tok, `"';&|`)
		if clean == "add" || clean == "stash" {
			return true
		}
		if len(clean) > 1 && clean[0] == '-' && clean[1] != '-' {
			for _, r := range clean[1:] {
				if r == 'A' || r == 'a' || r == 'u' {
					return true
				}
			}
		}
		if clean == "--all" || clean == "--update" {
			return true
		}
	}
	return false
}

// namesOnlyOrdinaryFiles reports whether every path-shaped argument resolves to
// an existing regular file that is not credential-shaped. Anything else — a
// glob, a directory, a quoted fragment, a path that does not exist, a blanket
// flag — is treated as unresolvable, and the tree scan applies.
func namesOnlyOrdinaryFiles(root, cmd string) bool {
	named := 0
	for _, tok := range strings.Fields(cmd) {
		if tok == "" || strings.ContainsAny(tok, `'"$`+"`") {
			return false // quoting or expansion: cannot be read literally
		}
		if strings.HasPrefix(tok, "-") {
			// Blanket flags stage far more than they name.
			for _, r := range strings.TrimPrefix(tok, "-") {
				if r == 'A' || r == 'a' || r == 'u' {
					return false
				}
			}
			continue
		}
		if !looksLikePath(tok) {
			continue
		}
		if strings.ContainsAny(tok, "*?[") {
			return false // a glob: the shell decides what it means
		}
		fi, err := os.Stat(filepath.Join(root, tok))
		if err != nil || !fi.Mode().IsRegular() || isSecretPath(tok) {
			return false
		}
		named++
	}
	return named > 0
}

// looksLikePath distinguishes a file argument from a subcommand or a message.
func looksLikePath(tok string) bool {
	return strings.Contains(tok, "/") || strings.Contains(tok, ".")
}

// secretsInIndex lists credential-shaped paths git currently has staged.
func secretsInIndex(root string) []string {
	out, err := gitOutput(root, "diff", "--cached", "--name-only")
	if err != nil {
		return nil
	}
	return filterSecrets(strings.Split(out, "\n"))
}

// secretsInTree lists credential-shaped paths a stage could pick up: untracked
// or modified files that are not yet in the index.
func secretsInTree(root string) []string {
	out, err := gitOutput(root, "status", "--porcelain", "-uall")
	if err != nil {
		return nil
	}
	var paths []string
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 4 {
			continue
		}
		paths = append(paths, strings.Trim(line[3:], `"`))
	}
	return filterSecrets(paths)
}

func filterSecrets(paths []string) []string {
	var found []string
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p != "" && isSecretPath(p) {
			found = append(found, p)
			if len(found) >= 5 {
				break
			}
		}
	}
	return found
}

func gitOutput(root string, args ...string) (string, error) {
	full := append([]string{"-c", "core.fsmonitor=false", "-C", root}, args...)
	out, err := exec.Command("git", full...).Output()
	return strings.TrimSpace(string(out)), err
}
