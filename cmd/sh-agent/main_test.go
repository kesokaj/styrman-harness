package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime/debug"
	"strings"
	"syscall"
	"testing"

	"github.com/kesokaj/styrman-harness/internal/harness"
)

func TestReorder(t *testing.T) {
	cases := map[string][]string{
		// BUG-002: the flag package stops at the first positional argument.
		"adopt my-skill --home /x":          {"--home", "/x", "adopt", "my-skill"},
		"--vendor agy install":              {"--vendor", "agy", "install"},
		"adopt a --force":                   {"--force", "adopt", "a"},
		"adopt -- --looks-like-a-flag":      {"adopt", "--looks-like-a-flag"},
		"--repo r --home h adopt x --force": {"--repo", "r", "--home", "h", "--force", "adopt", "x"},
		// A value flag missing from takesValue swallows the next argument:
		// --feedback-repo would take "--dry-run" as its value and the run
		// would push when a dry run was asked (BUG-002's shape).
		"report --feedback-repo /x --dry-run": {"--feedback-repo", "/x", "--dry-run", "report"},
		"report --dry-run":                    {"--dry-run", "report"},
	}
	for in, want := range cases {
		if got := reorder(strings.Fields(in)); !reflect.DeepEqual(got, want) {
			t.Errorf("reorder(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestResolveVendor(t *testing.T) {
	home := t.TempDir()
	if v, err := resolveVendor("", home); err != nil || v != harness.VendorAll {
		t.Errorf("no flag, no manifest: %q, %v; want all", v, err)
	}
	if v, err := resolveVendor("agy", home); err != nil || v != harness.VendorAntigravity {
		t.Errorf("flag agy: %q, %v", v, err)
	}
	if _, err := resolveVendor("vscode", home); err == nil {
		t.Error("unknown vendor should fail")
	}
	if err := harness.SaveState(home, harness.State{Vendor: harness.VendorClaude}); err != nil {
		t.Fatal(err)
	}
	if v, err := resolveVendor("", home); err != nil || v != harness.VendorClaude {
		t.Errorf("manifest vendor: %q, %v; want claude", v, err)
	}
	if v, err := resolveVendor("all", home); err != nil || v != harness.VendorAll {
		t.Errorf("flag beats manifest: %q, %v", v, err)
	}
}

func TestShouldRebuild(t *testing.T) {
	cases := []struct{ goOnPath, child, own, want bool }{
		{true, false, true, true},
		{true, true, true, false},
		{false, false, true, false},
		{false, true, true, false},
		{true, false, false, false}, // a plain copy never rebuilds (BUG-015)
		{false, false, false, false},
	}
	for _, c := range cases {
		if got := shouldRebuild(c.goOnPath, c.child, c.own); got != c.want {
			t.Errorf("shouldRebuild(%v, %v, %v) = %v", c.goOnPath, c.child, c.own, got)
		}
	}
}

func TestBinaryOlder(t *testing.T) {
	cases := []struct {
		rev, commit string
		want        bool
	}{
		{"b64ffe886729", "b64ffe8", false},
		{"b64ffe8", "b64ffe886729", false},
		{"b64ffe886729+dirty", "b64ffe8", false},
		{"50ac771abcde", "b64ffe8", true},
		{"(devel)", "b64ffe8", true},
		{"unknown", "b64ffe8", true},
		{"", "b64ffe8", false},
		{"b64ffe886729", "", false},
	}
	for _, c := range cases {
		if got := binaryOlder(c.rev, c.commit); got != c.want {
			t.Errorf("binaryOlder(%q, %q) = %v", c.rev, c.commit, got)
		}
	}
}

func TestVersionFromInfo(t *testing.T) {
	stamp := func(rev, modified string) *debug.BuildInfo {
		return &debug.BuildInfo{
			Main:     debug.Module{Version: "(devel)"},
			Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: rev}, {Key: "vcs.modified", Value: modified}},
		}
	}
	cases := []struct {
		name string
		info *debug.BuildInfo
		want string
	}{
		{"clean", stamp("b64ffe886729abcdef0123456789abcdef012345", "false"), "b64ffe886729"},
		{"dirty", stamp("b64ffe886729abcdef0123456789abcdef012345", "true"), "b64ffe886729+dirty"},
		{"short revision kept whole", stamp("b64ffe8", "false"), "b64ffe8"},
		{"no revision", &debug.BuildInfo{Main: debug.Module{Version: "v1.2.3"}}, "v1.2.3"},
	}
	for _, c := range cases {
		if got := versionFromInfo(c.info); got != c.want {
			t.Errorf("%s: versionFromInfo = %q, want %q", c.name, got, c.want)
		}
	}
}

// TestRunInstallStatusUninstall drives the CLI end to end against this
// checkout and a scratch home. ChildEnv is set so install takes the plain
// path: under `go test` the running binary is the test binary, and a rebuild
// here would be a real `go build` of the repo for nothing.
func TestRunInstallStatusUninstall(t *testing.T) {
	t.Setenv(harness.ChildEnv, "1")
	home := t.TempDir()
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run([]string{"install", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("install: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "files managed") {
		t.Errorf("install output:\n%s", out.String())
	}
	if _, err := os.Stat(filepath.Join(home, harness.BinPath)); err != nil {
		t.Errorf("binary not installed: %v", err)
	}
	out.Reset()
	if err := run([]string{"status", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("status: %v\n%s", err, out.String())
	}
	for _, want := range []string{"binary     ", "drift      none"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("status missing %q:\n%s", want, out.String())
		}
	}
	// The binary line reads the installed file's build info: the test binary
	// was copied there, so the line must carry its version, never "unknown".
	if strings.Contains(out.String(), "binary     unknown") {
		t.Errorf("status could not read the installed binary's build info:\n%s", out.String())
	}
	out.Reset()
	if err := run([]string{"uninstall", "--home", home}, &out); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	// The vendor roots themselves are kept on purpose (pruneEmptyDirs stops
	// there); everything the manifest listed under them is gone.
	for _, rel := range []string{filepath.Join(harness.ClaudeRoot, "AGENTS.md"), filepath.Join(harness.ClaudeRoot, "skills"), filepath.Join(harness.AgyRoot, "agents"), harness.BinPath} {
		if _, err := os.Stat(filepath.Join(home, rel)); !os.IsNotExist(err) {
			t.Errorf("~/%s should be gone after uninstall, stat err = %v", rel, err)
		}
	}
	out.Reset()
	if err := run([]string{"status", "--repo", repo, "--home", home}, &out); err == nil {
		t.Error("status after uninstall should report exit 1")
	}
}

func TestKnownCheckout(t *testing.T) {
	home := t.TempDir()
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("SH_AGENT_REPO", "")
	if err := knownCheckout("", home, repo); err != nil {
		t.Errorf("first install (no manifest) must pass: %v", err)
	}
	if err := harness.SaveState(home, harness.State{Repo: repo}); err != nil {
		t.Fatal(err)
	}
	if err := knownCheckout("", home, repo); err != nil {
		t.Errorf("same checkout must pass: %v", err)
	}
	other := t.TempDir()
	if err := knownCheckout("", home, other); err == nil {
		t.Error("a checkout found only by the directory walk must be refused when another is installed")
	}
	if err := knownCheckout(other, home, other); err != nil {
		t.Errorf("explicit --repo must pass: %v", err)
	}
	t.Setenv("SH_AGENT_REPO", other)
	if err := knownCheckout("", home, other); err != nil {
		t.Errorf("$SH_AGENT_REPO must pass: %v", err)
	}
}

// TestDryRunNeverBuilds: with a toolchain present, no child guard and a
// checkout-shaped repo, a dry run must report the rebuild instead of
// performing it, and write nothing.
func TestDryRunNeverBuilds(t *testing.T) {
	t.Setenv(harness.ChildEnv, "")
	// A checkout-SHAPED repo of its own, not this checkout (`../..`): when
	// the test runs inside a plain copy with no .git — an unpacked bundle, or
	// once a copy vendored inside a consumer project — `../..` is not a
	// checkout, the dry run rightly prints the plain-copy line, and the
	// assertion below failed for a shape the test never meant to cover
	// (styrman-crew BUG-058). IsCheckout asks only for a `.git` entry and
	// Inspect tolerates one git cannot read, so an empty directory is enough.
	repo := checkoutShaped(t)
	home := t.TempDir()
	var out bytes.Buffer
	if err := run([]string{"install", "--dry-run", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("dry run: %v\n%s", err, out.String())
	}
	for _, want := range []string{"[dry-run] the binary would be rebuilt from", "[dry-run] "} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("dry-run output missing %q:\n%s", want, out.String())
		}
	}
	// Reinstall removes its build directory itself, so its absence proves
	// nothing; what a rebuild would leave is an installed binary and a
	// child's output, and a dry run must leave neither.
	for _, rel := range []string{harness.BinPath, ".claude", ".config"} {
		if _, err := os.Stat(filepath.Join(home, rel)); !os.IsNotExist(err) {
			t.Errorf("dry run wrote ~/%s: stat err = %v", rel, err)
		}
	}
	if strings.Contains(out.String(), "building sh-agent from") {
		t.Errorf("dry run ran a build:\n%s", out.String())
	}
}

// checkoutShaped is plainCopy plus the one thing IsCheckout looks for: a
// `.git` entry. It is what the dry-run test needs — "this is a checkout"
// without being THIS checkout — and it is not a git repository, which
// Inspect tolerates (empty fields) and nothing on the install path needs.
func checkoutShaped(t *testing.T) string {
	t.Helper()
	dir := plainCopy(t)
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// plainCopy copies this checkout's deployable trees into a directory with
// no .git: what an unpacked bundle is (R8).
func plainCopy(t *testing.T) string {
	t.Helper()
	src, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	dst := t.TempDir()
	for _, top := range []string{"context", "skills", "agents", "commands", "rules"} {
		root := filepath.Join(src, top)
		err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(src, p)
			if d.IsDir() {
				return os.MkdirAll(filepath.Join(dst, rel), 0o755)
			}
			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(dst, rel), b, 0o644)
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if harness.IsCheckout(dst) {
		t.Fatal("fixture: the plain copy must not carry .git")
	}
	return dst
}

// TestPlainCopyNeverRebuilds: an unpacked bundle's binary was built from
// the same tree, the tree may be read-only to the installing user, and it
// cannot be updated with git; install from it must deploy the running
// binary without a rebuild, without the no-toolchain warning, and
// without creating a scratch directory inside the copy (BUG-015).
func TestPlainCopyNeverRebuilds(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go not on PATH: the rebuild path is unreachable, nothing to prove")
	}
	t.Setenv(harness.ChildEnv, "")
	os.Unsetenv(harness.ChildEnv)
	copyDir := plainCopy(t)
	home := t.TempDir()

	var out bytes.Buffer
	if err := run([]string{"install", "--repo", copyDir, "--home", home}, &out); err != nil {
		t.Fatalf("install from a plain copy: %v\n%s", err, out.String())
	}
	for _, never := range []string{"building sh-agent from", "would be rebuilt", "warning: go is not on PATH"} {
		if strings.Contains(out.String(), never) {
			t.Errorf("install from a plain copy printed %q:\n%s", never, out.String())
		}
	}
	if _, err := os.Stat(filepath.Join(copyDir, ".agents")); !os.IsNotExist(err) {
		t.Errorf("install created %s inside the plain copy (err=%v)", filepath.Join(copyDir, ".agents"), err)
	}
	if _, err := os.Stat(filepath.Join(home, harness.BinPath)); err != nil {
		t.Errorf("binary not installed from the plain copy: %v", err)
	}

	out.Reset()
	if err := run([]string{"install", "--dry-run", "--repo", copyDir, "--home", home}, &out); err != nil {
		t.Fatalf("dry run from a plain copy: %v\n%s", err, out.String())
	}
	if strings.Contains(out.String(), "would be rebuilt") {
		t.Errorf("dry run from a plain copy announced a rebuild:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "[dry-run] plain copy") {
		t.Errorf("dry run from a plain copy must say so as a dry run:\n%s", out.String())
	}
}

// A hand-written ~/.claude/CLAUDE.md survives install, and status names it
// with where its content belongs, so the one-line import can replace it (R2).
func TestStatusNamesForeignClaudeMD(t *testing.T) {
	t.Setenv(harness.ChildEnv, "1")
	home := t.TempDir()
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	claudeMD := filepath.Join(home, harness.ClaudeRoot, "CLAUDE.md")
	if err := os.MkdirAll(filepath.Dir(claudeMD), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(claudeMD, []byte("my own notes\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run([]string{"install", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("install: %v\n%s", err, out.String())
	}
	out.Reset()
	if err := run([]string{"status", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("status: %v\n%s", err, out.String())
	}
	want := "note       ~/.claude/CLAUDE.md is yours; move only rules meant for every machine into a rule file under rules/ in the harness repo — rules/ ships to every install on every machine, never credentials or personal or host-specific notes"
	if !strings.Contains(out.String(), want) {
		t.Errorf("status missing %q:\n%s", want, out.String())
	}
}

// Claude Code adds a Co-Authored-By trailer to every commit unless its own
// attribution setting is off; the context's Git rule alone did not stop it in
// the sandbox (three of three commits). install and status say so until the
// setting is there (settings.json is the operator's file: never written).
func TestAttributionNote(t *testing.T) {
	t.Setenv(harness.ChildEnv, "1")
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	// %s is the settings.json path under the scratch home: the note names the file it read.
	const body = `%s does not set "attribution": {"commit": ""}; Claude Code adds a Co-Authored-By trailer to every commit until it does (context § Git)`
	prefix := map[string]string{"install": "note: ", "status": "note       "}
	cases := []struct {
		name     string
		settings string // "" means no file; "DIR" makes settings.json a directory, unreadable as a file
		vendor   string // extra --vendor flag, "" for the default (all)
		want     string // the note body expected after the command's prefix; "" means no note line at all
	}{
		// install writes the key when it is absent (attributionSet), so the note is never printed for these two
		{"no settings file", "", "", ""},
		{"trailer switched off", `{"attribution": {"commit": "", "pr": ""}}` + "\n", "", ""},
		{"trailer switched off with false", `{"attribution": {"commit": false}}` + "\n", "", ""},
		{"custom trailer", `{"attribution": {"commit": "Co-Authored-By: me"}}` + "\n", "", body},
		{"other settings only", `{"permissions": {"allow": ["Bash(git:*)"]}}` + "\n", "", ""},
		{"antigravity only never notes", "", "agy", ""},
		{"malformed json is advisory, no echo", `{"attribution": ` + "\n", "", "%s is not valid JSON, so the attribution setting could not be checked; Claude Code adds"},
		{"a directory is advisory", "DIR", "", "%s is not a regular file, so the attribution setting could not be checked"},
		{"a fifo is advisory and never opened", "FIFO", "", "%s is not a regular file, so the attribution setting could not be checked"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			home := t.TempDir()
			settingsPath := filepath.Join(home, harness.ClaudeRoot, "settings.json")
			if c.settings != "" {
				dir := filepath.Join(home, harness.ClaudeRoot)
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
				if c.settings == "DIR" {
					if err := os.Mkdir(filepath.Join(dir, "settings.json"), 0o755); err != nil {
						t.Fatal(err)
					}
				} else if c.settings == "FIFO" {
					// opening a FIFO with no writer blocks forever: the Stat guard must keep the command from opening it
					if err := syscall.Mkfifo(filepath.Join(dir, "settings.json"), 0o600); err != nil {
						t.Fatal(err)
					}
				} else if err := os.WriteFile(filepath.Join(dir, "settings.json"), []byte(c.settings), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			for _, cmd := range []string{"install", "status"} {
				args := []string{cmd, "--repo", repo, "--home", home}
				if c.vendor != "" {
					args = append(args, "--vendor", c.vendor)
				}
				var out bytes.Buffer
				if err := run(args, &out); err != nil {
					t.Fatalf("%s: %v\n%s", cmd, err, out.String())
				}
				noted := strings.Contains(out.String(), prefix[cmd]+settingsPath)
				if c.want == "" && noted {
					t.Errorf("%s: unexpected settings.json note:\n%s", cmd, out.String())
				}
				want := prefix[cmd] + fmt.Sprintf(c.want, settingsPath)
				if c.want != "" && !strings.Contains(out.String(), want) {
					t.Errorf("%s: missing note %q:\n%s", cmd, want, out.String())
				}
				if strings.Contains(out.String(), "invalid character") || strings.Contains(out.String(), "unexpected end") {
					t.Errorf("%s: the note echoes the JSON error:\n%s", cmd, out.String())
				}
			}
		})
	}
}

// install writes Claude Code's attribution setting when the key is absent —
// the one vendor setting the harness owns (R4, ADR-025): the file is created
// when missing, every other byte of an existing file stays, a present key of
// any value is left alone, and nothing about the file fails the install.
func TestAttributionSet(t *testing.T) {
	t.Setenv(harness.ChildEnv, "1")
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	const fragment = `"attribution": {"commit": "", "pr": ""}`
	pretty := "{\n  \"model\": \"claude-fable-5-1[1m]\",\n  \"theme\": \"dark\",\n  \"autoMode\": {\n    \"environment\": [\"**Repository visibility**: Not queryable here — assume private\"]\n  }\n}\n"
	run1 := func(t *testing.T, home string, args ...string) string {
		t.Helper()
		var out bytes.Buffer
		if err := run(append([]string{"install", "--repo", repo, "--home", home}, args...), &out); err != nil {
			t.Fatalf("install: %v\n%s", err, out.String())
		}
		return out.String()
	}
	write := func(t *testing.T, home, content string, mode os.FileMode) string {
		t.Helper()
		dir := filepath.Join(home, harness.ClaudeRoot)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		p := filepath.Join(dir, "settings.json")
		if err := os.WriteFile(p, []byte(content), mode); err != nil {
			t.Fatal(err)
		}
		return p
	}
	keys := func(t *testing.T, data []byte) map[string]json.RawMessage {
		t.Helper()
		var m map[string]json.RawMessage
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatalf("not a JSON object: %v\n%s", err, data)
		}
		return m
	}

	t.Run("missing file is created with only the key", func(t *testing.T) {
		home := t.TempDir()
		out := run1(t, home)
		p := filepath.Join(home, harness.ClaudeRoot, "settings.json")
		if !strings.Contains(out, "settings: "+p+" created with "+fragment) {
			t.Errorf("missing settings line:\n%s", out)
		}
		if strings.Contains(out, "note: "+p) {
			t.Errorf("note printed although the key was just written:\n%s", out)
		}
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		m := keys(t, data)
		if len(m) != 1 || string(m["attribution"]) != `{"commit": "", "pr": ""}` {
			t.Errorf("file holds %s", data)
		}
		if again := run1(t, home); strings.Contains(again, "settings:") {
			t.Errorf("second install wrote again:\n%s", again)
		}
	})

	t.Run("other keys, formatting and mode stay", func(t *testing.T) {
		home := t.TempDir()
		p := write(t, home, pretty, 0o600)
		out := run1(t, home)
		if !strings.Contains(out, "settings: "+p+" — "+fragment+" added; other keys untouched") {
			t.Errorf("missing settings line:\n%s", out)
		}
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		// everything before the closing brace of the original object is untouched, byte for byte
		prefix := strings.TrimRight(strings.TrimSuffix(strings.TrimRight(pretty, "\n"), "}"), " \n")
		if !strings.HasPrefix(string(data), prefix) {
			t.Errorf("original bytes changed:\n%s", data)
		}
		before, after := keys(t, []byte(pretty)), keys(t, data)
		for k, v := range before {
			if string(after[k]) != string(v) {
				t.Errorf("key %q changed: %s -> %s", k, v, after[k])
			}
		}
		if len(after) != len(before)+1 || string(after["attribution"]) != `{"commit": "", "pr": ""}` {
			t.Errorf("attribution not added exactly once: %s", data)
		}
		if !strings.HasSuffix(string(data), "}\n") {
			t.Errorf("trailing newline lost: %q", data[len(data)-3:])
		}
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0o600 {
			t.Errorf("mode changed to %v", info.Mode().Perm())
		}
	})

	t.Run("crlf file keeps its line endings", func(t *testing.T) {
		home := t.TempDir()
		p := write(t, home, "{\r\n  \"model\": \"x\"\r\n}\r\n", 0o644)
		run1(t, home)
		data, _ := os.ReadFile(p)
		if strings.Count(string(data), "\n") != strings.Count(string(data), "\r\n") {
			t.Errorf("mixed line endings: %q", data)
		}
		if m := keys(t, data); string(m["model"]) != `"x"` || m["attribution"] == nil {
			t.Errorf("file holds %s", data)
		}
	})

	t.Run("empty object", func(t *testing.T) {
		home := t.TempDir()
		p := write(t, home, "{}\n", 0o644)
		run1(t, home)
		data, _ := os.ReadFile(p)
		m := keys(t, data)
		if len(m) != 1 || m["attribution"] == nil {
			t.Errorf("file holds %s", data)
		}
	})

	for _, c := range []struct{ name, content, line string }{
		{"custom value is left alone", `{"attribution": {"commit": "Co-Authored-By: me"}}` + "\n", "note: "},
		{"off value is left alone", `{"attribution": {"commit": ""}}` + "\n", ""},
		{"malformed is left alone", `{"attribution": ` + "\n", "note: "},
		{"not an object is left alone", `["x"]` + "\n", "note: "},
		{"null is left alone", "null\n", "note: "},
	} {
		t.Run(c.name, func(t *testing.T) {
			home := t.TempDir()
			p := write(t, home, c.content, 0o644)
			out := run1(t, home)
			data, _ := os.ReadFile(p)
			if string(data) != c.content {
				t.Errorf("file changed:\n%s", data)
			}
			if strings.Contains(out, "settings: "+p+" —") || strings.Contains(out, "created with") {
				t.Errorf("a write was reported:\n%s", out)
			}
			if c.line != "" && !strings.Contains(out, c.line+p) {
				t.Errorf("expected a %q line for %s:\n%s", c.line, p, out)
			}
			if c.line == "" && strings.Contains(out, p) {
				t.Errorf("unexpected line about the file:\n%s", out)
			}
		})
	}

	t.Run("a file over the read bound is left alone", func(t *testing.T) {
		home := t.TempDir()
		big := `{"pad": "` + strings.Repeat("x", 1<<20) + `"}` + "\n"
		p := write(t, home, big, 0o644)
		out := run1(t, home)
		if !strings.Contains(out, "settings: "+p+" is larger than 1 MiB") {
			t.Errorf("missing size line:\n%s", out)
		}
		data, _ := os.ReadFile(p)
		if string(data) != big {
			t.Errorf("file changed (len %d)", len(data))
		}
	})

	t.Run("a symlinked settings.json is never written through", func(t *testing.T) {
		home := t.TempDir()
		target := filepath.Join(t.TempDir(), "dotfiles-settings.json")
		if err := os.WriteFile(target, []byte("{}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		dir := filepath.Join(home, harness.ClaudeRoot)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		p := filepath.Join(dir, "settings.json")
		if err := os.Symlink(target, p); err != nil {
			t.Fatal(err)
		}
		out := run1(t, home)
		if !strings.Contains(out, "settings: "+p+" is a symlink; add the key to its target by hand") {
			t.Errorf("missing symlink line:\n%s", out)
		}
		if info, err := os.Lstat(p); err != nil || info.Mode()&os.ModeSymlink == 0 {
			t.Errorf("the symlink was replaced: %v %v", info, err)
		}
		if data, _ := os.ReadFile(target); string(data) != "{}\n" {
			t.Errorf("the target was written: %s", data)
		}
	})

	t.Run("dry run writes nothing", func(t *testing.T) {
		home := t.TempDir()
		out := run1(t, home, "--dry-run")
		p := filepath.Join(home, harness.ClaudeRoot, "settings.json")
		if _, err := os.Stat(p); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("dry run created %s", p)
		}
		if !strings.Contains(out, "[dry-run] would create "+p) {
			t.Errorf("missing dry-run line:\n%s", out)
		}
	})

	t.Run("antigravity only writes nothing", func(t *testing.T) {
		home := t.TempDir()
		out := run1(t, home, "--vendor", "agy")
		if _, err := os.Stat(filepath.Join(home, harness.ClaudeRoot, "settings.json")); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("settings.json created for the agy vendor")
		}
		if strings.Contains(out, "settings:") {
			t.Errorf("settings line for the agy vendor:\n%s", out)
		}
	})

	t.Run("status alone keeps the note", func(t *testing.T) {
		home := t.TempDir()
		run1(t, home) // a full install writes the key …
		p := filepath.Join(home, harness.ClaudeRoot, "settings.json")
		if err := os.Remove(p); err != nil { // … and an operator can delete the file again
			t.Fatal(err)
		}
		var out bytes.Buffer
		if err := run([]string{"status", "--repo", repo, "--home", home}, &out); err != nil {
			t.Fatalf("status: %v\n%s", err, out.String())
		}
		if !strings.Contains(out.String(), "note       "+filepath.Join(home, harness.ClaudeRoot, "settings.json")+" does not set") {
			t.Errorf("status lost the note:\n%s", out.String())
		}
	})
}

// gitCleanEnv strips the caller's git environment. GIT_DIR or GIT_WORK_TREE
// set in the parent would point every fixture command at the operator's own
// checkout — the shape of BUG-045, where a probe's commands fell through to
// the real repository — and GIT_CONFIG_* would let their config decide what
// the fixture does.
func gitCleanEnv() []string {
	var env []string
	for _, kv := range os.Environ() {
		switch {
		case strings.HasPrefix(kv, "GIT_DIR="),
			strings.HasPrefix(kv, "GIT_WORK_TREE="),
			strings.HasPrefix(kv, "GIT_INDEX_FILE="),
			strings.HasPrefix(kv, "GIT_OBJECT_DIRECTORY="),
			strings.HasPrefix(kv, "GIT_COMMON_DIR="),
			strings.HasPrefix(kv, "GIT_CONFIG"):
			continue
		}
		env = append(env, kv)
	}
	return env
}

// gitFixture builds a real checkout with an upstream: a bare origin, a clone
// with a trivial cmd/sh-agent that records that it ran, and a first commit
// pushed. Behind() needs a real remote, and Reinstall() needs something to
// build, so the guard under test is exercised on the path it actually
// guards. Returns the work tree and the marker path the built binary writes.
func gitFixture(t *testing.T) (repo, marker string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go is not on PATH")
	}
	base := t.TempDir()
	origin := filepath.Join(base, "origin.git")
	repo = filepath.Join(base, "work")
	marker = filepath.Join(base, "installed.txt")
	run := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		cmd.Env = append(gitCleanEnv(), "GIT_CONFIG_GLOBAL="+filepath.Join(base, "gitconfig"), "GIT_CONFIG_NOSYSTEM=1")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v: %v\n%s", args, err, out)
		}
	}
	run(base, "git", "init", "--quiet", "--bare", "-b", "main", origin)
	run(base, "git", "clone", "--quiet", origin, repo)
	run(repo, "git", "config", "user.email", "fixture@example.invalid")
	run(repo, "git", "config", "user.name", "fixture")
	write := func(rel, content string) {
		t.Helper()
		p := filepath.Join(repo, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module fixture\n\ngo 1.24\n")
	// Stands in for the real installer: records that an install ran at all.
	write("cmd/sh-agent/main.go", "package main\n\nimport \"os\"\n\nfunc main() {\n\tos.WriteFile(os.Getenv(\"FIXTURE_MARKER\"), []byte(\"installed\"), 0o644)\n}\n")
	write("context/AGENTS.md", "# AGENTS.md\n")
	run(repo, "git", "add", "-A")
	run(repo, "git", "commit", "--quiet", "-m", "fixture")
	run(repo, "git", "push", "--quiet", "-u", "origin", "HEAD:refs/heads/main")
	run(repo, "git", "branch", "--quiet", "--set-upstream-to=origin/main")
	t.Setenv("FIXTURE_MARKER", marker)
	return repo, marker
}

// TestUpdateKeptCloneCleanupOnlyWhenUpdating: `update` removes the feedback
// clone older versions kept in this home (D9), but `--check` reports and
// changes NOTHING — the cleanup sat above its returns, so the read-only mode
// deleted files (security review 2026-09-17, finding 2).
func TestUpdateKeptCloneCleanupOnlyWhenUpdating(t *testing.T) {
	repo, marker := gitFixture(t)
	home := t.TempDir()

	// A clone that passes RemoveKeptClone's checks without a network: a local
	// repository whose origin URL is set to the feedback URL it verifies.
	plant := func(t *testing.T) string {
		t.Helper()
		base := t.TempDir()
		origin := filepath.Join(base, "feedback.git")
		kept := filepath.Join(home, ".config", "sh-agent", "feedback")
		_ = os.RemoveAll(kept)
		if err := os.MkdirAll(filepath.Dir(kept), 0o755); err != nil {
			t.Fatal(err)
		}
		for _, args := range [][]string{
			{"init", "--quiet", "--bare", "-b", "main", origin},
			{"clone", "--quiet", origin, kept},
			{"-C", kept, "remote", "set-url", "origin", harness.FeedbackURL},
		} {
			cmd := exec.Command("git", args...)
			cmd.Dir = base
			cmd.Env = append(gitCleanEnv(), "GIT_CONFIG_GLOBAL="+filepath.Join(base, "gitconfig"), "GIT_CONFIG_NOSYSTEM=1")
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("git %v: %v\n%s", args, err, out)
			}
		}
		return kept
	}

	kept := plant(t)
	var out bytes.Buffer
	if err := update(repo, home, true, "", "", &out); err != nil {
		t.Fatalf("update --check: %v\n%s", err, out.String())
	}
	if _, err := os.Stat(kept); err != nil {
		t.Errorf("--check removed the kept clone: %v\n%s", err, out.String())
	}
	if strings.Contains(out.String(), "removed the kept feedback clone") {
		t.Errorf("--check reported a removal:\n%s", out.String())
	}

	// The same update without --check does remove it.
	out.Reset()
	if err := update(repo, home, false, "", "", &out); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}
	if _, err := os.Stat(kept); !os.IsNotExist(err) {
		t.Errorf("update did not remove the kept clone: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "removed the kept feedback clone") {
		t.Errorf("the removal was silent:\n%s", out.String())
	}
	if _, err := os.Stat(marker); err != nil {
		t.Errorf("the install did not run: %v", err)
	}
}

// TestUpdateNeverInstallsDirtyTree is BUG-046: `update` deploys the working
// tree, not the commit it fetched, so a dirty checkout must not be installed
// on EITHER path. Pull already refused it when behind; the up-to-date path
// had no guard, and an hourly timer put a session's uncommitted edits into
// the operator's real home. --check only reports, so it still works dirty.
func TestUpdateNeverInstallsDirtyTree(t *testing.T) {
	installed := func(marker string) bool {
		_, err := os.Stat(marker)
		return err == nil
	}
	dirty := func(t *testing.T, repo string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(repo, "context", "AGENTS.md"), []byte("# edited, not committed\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	behind := func(t *testing.T, repo string) {
		t.Helper()
		other := t.TempDir()
		base := filepath.Dir(repo)
		run := func(dir string, args ...string) {
			t.Helper()
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Dir = dir
			cmd.Env = append(gitCleanEnv(), "GIT_CONFIG_GLOBAL="+filepath.Join(base, "gitconfig"), "GIT_CONFIG_NOSYSTEM=1")
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("%v: %v\n%s", args, err, out)
			}
		}
		run(other, "git", "clone", "--quiet", filepath.Join(base, "origin.git"), filepath.Join(other, "c"))
		c := filepath.Join(other, "c")
		run(c, "git", "config", "user.email", "fixture@example.invalid")
		run(c, "git", "config", "user.name", "fixture")
		if err := os.WriteFile(filepath.Join(c, "context", "AGENTS.md"), []byte("# upstream moved\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		run(c, "git", "commit", "--quiet", "-am", "upstream")
		run(c, "git", "push", "--quiet", "origin", "HEAD:refs/heads/main")
	}

	t.Run("clean and up to date installs", func(t *testing.T) {
		repo, marker := gitFixture(t)
		var out bytes.Buffer
		if err := update(repo, t.TempDir(), false, "", "", &out); err != nil {
			t.Fatalf("clean update: %v\n%s", err, out.String())
		}
		if !installed(marker) {
			t.Error("a clean checkout should have been installed")
		}
	})

	t.Run("dirty and up to date refuses", func(t *testing.T) {
		repo, marker := gitFixture(t)
		dirty(t, repo)
		var out bytes.Buffer
		err := update(repo, t.TempDir(), false, "", "", &out)
		if err == nil {
			t.Fatal("a dirty checkout must not be installed")
		}
		if !strings.Contains(err.Error(), "uncommitted changes") {
			t.Errorf("error should name the uncommitted changes, got %v", err)
		}
		if installed(marker) {
			t.Error("nothing should have been installed from a dirty tree")
		}
	})

	t.Run("dirty and behind refuses without pulling", func(t *testing.T) {
		repo, marker := gitFixture(t)
		before := harness.Inspect(repo).Commit
		behind(t, repo)
		dirty(t, repo)
		var out bytes.Buffer
		err := update(repo, t.TempDir(), false, "", "", &out)
		if err == nil || !strings.Contains(err.Error(), "uncommitted changes") {
			t.Fatalf("dirty and behind must refuse, got %v", err)
		}
		if installed(marker) {
			t.Error("nothing should have been installed")
		}
		if after := harness.Inspect(repo).Commit; after != before {
			t.Errorf("the tree moved from %s to %s; a refused update pulls nothing", before, after)
		}
	})

	t.Run("clean and behind pulls then installs", func(t *testing.T) {
		repo, marker := gitFixture(t)
		before := harness.Inspect(repo).Commit
		behind(t, repo)
		var out bytes.Buffer
		if err := update(repo, t.TempDir(), false, "", "", &out); err != nil {
			t.Fatalf("clean and behind: %v\n%s", err, out.String())
		}
		if after := harness.Inspect(repo).Commit; after == before {
			t.Error("a clean checkout behind upstream should have been fast-forwarded")
		}
		if !installed(marker) {
			t.Error("it should have been installed after the pull")
		}
	})

	t.Run("--check reports on a dirty tree and installs nothing", func(t *testing.T) {
		repo, marker := gitFixture(t)
		dirty(t, repo)
		var out bytes.Buffer
		if err := update(repo, t.TempDir(), true, "", "", &out); err != nil {
			t.Fatalf("--check only reports, it must not fail on a dirty tree: %v", err)
		}
		if !strings.Contains(out.String(), "up to date") {
			t.Errorf("--check should report the state, got %q", out.String())
		}
		if installed(marker) {
			t.Error("--check installs nothing")
		}
	})
}

// TestAbandonTempReportsRemoveFailure forces the error path the verifier
// reported as uncovered: the settings write failed AND its temp file could
// not be removed. Both lines must reach the operator; neither may panic.
func TestAbandonTempReportsRemoveFailure(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores the directory permissions this test relies on")
	}
	dir := t.TempDir()
	tmp := filepath.Join(dir, ".settings-tmp.json")
	if err := os.WriteFile(tmp, []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o700) })

	var out bytes.Buffer
	abandonTemp(tmp, filepath.Join(dir, "settings.json"), errors.New("disk full"), &out)
	got := out.String()
	if !strings.Contains(got, "disk full") {
		t.Errorf("the cause should be reported, got %q", got)
	}
	if !strings.Contains(got, "left behind") {
		t.Errorf("a temp file that could not be removed should be named, got %q", got)
	}
	if !strings.Contains(got, tmp) {
		t.Errorf("the temp path should be named so the operator can remove it, got %q", got)
	}
}

// TestReadSettingsProblems covers the reader's refusals, including the open
// failure the verifier reported as untested. A problem is named without a
// byte of the file's contents.
func TestReadSettingsProblems(t *testing.T) {
	dir := t.TempDir()

	if _, problem := readSettings(filepath.Join(dir, "absent.json")); problem != "could not be read" {
		t.Errorf("missing file: %q", problem)
	}

	sub := filepath.Join(dir, "adirectory.json")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, problem := readSettings(sub); problem != "is not a regular file" {
		t.Errorf("directory: %q", problem)
	}

	if os.Geteuid() == 0 {
		t.Skip("root can open a mode-000 file, so the open-failure leg cannot run")
	}
	unreadable := filepath.Join(dir, "unreadable.json")
	secret := `{"attribution":{"commit":"SECRET"}}`
	if err := os.WriteFile(unreadable, []byte(secret), 0o000); err != nil {
		t.Fatal(err)
	}
	data, problem := readSettings(unreadable)
	if problem != "could not be read" {
		t.Errorf("unreadable regular file: %q", problem)
	}
	if len(data) != 0 {
		t.Errorf("no bytes may be returned when the open failed, got %q", data)
	}
	// The problem strings are a closed set of fixed phrases, so asserting one
	// does not contain the secret proves nothing; what is worth asserting is
	// that the set IS closed — a future reader returning the OS error verbatim
	// would leak a path, and returning file bytes would leak the file.
	for _, allowed := range []string{"could not be read", "is not a regular file", ""} {
		if problem == allowed {
			return
		}
	}
	t.Errorf("readSettings returned an unrecognised problem string %q; it must stay a closed set so nothing from the file or the OS error escapes", problem)
}

// TestUpdateRefusesWhenStatusIsUnknown closes the coverage gap the round-2
// verifier named: `update`'s StatusUnknown branch had no cmd-level test.
// "git could not answer" must refuse, not deploy — a checkout whose state
// cannot be read is not a clean one.
func TestUpdateRefusesWhenStatusIsUnknown(t *testing.T) {
	repo, marker := gitFixture(t)
	// Break git's view of the tree while leaving the checkout recognisable:
	// the objects are gone, so `status` fails where `rev-parse` may not.
	if err := os.RemoveAll(filepath.Join(repo, ".git", "objects")); err != nil {
		t.Fatal(err)
	}
	if info := harness.Inspect(repo); !info.StatusUnknown && !info.Dirty {
		t.Skip("this git reports a readable status without its object store; the branch is covered in internal/harness")
	}
	var out bytes.Buffer
	err := update(repo, t.TempDir(), false, "", "", &out)
	if err == nil {
		t.Fatal("a checkout whose status git cannot report must not be installed")
	}
	if _, statErr := os.Stat(marker); statErr == nil {
		t.Error("nothing should have been installed")
	}
}

// TestReportWithoutCheckout: report runs before ResolveSource, so a machine that
// installed from an unpacked bundle — no checkout anywhere — still reports.
// A dry run with the default feedback path absent clones nothing.
func TestReportWithoutCheckout(t *testing.T) {
	t.Setenv("SH_AGENT_REPO", "")
	home, notRepo := t.TempDir(), t.TempDir()
	var out bytes.Buffer
	if err := run([]string{"report", "--repo", notRepo, "--home", home, "--dry-run"}, &out); err != nil {
		t.Fatalf("report --dry-run with no checkout: %v\n%s", err, out.String())
	}
	for _, want := range []string{"not a harness checkout", "harness_commit, drift, skills and agents are null", "[dry-run] would clone", "report.json"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, out.String())
		}
	}
	if entries, _ := os.ReadDir(home); len(entries) != 0 {
		t.Errorf("dry run wrote into home: %v", entries)
	}
}

// An embedded install's manifest records a release, not a path: knownCheckout
// says so and never resolves it as one (security review 2026-09-17).
func TestKnownCheckoutRefusesEmbeddedManifest(t *testing.T) {
	t.Setenv("SH_AGENT_REPO", "")
	home := t.TempDir()
	if err := harness.SaveState(home, harness.State{Repo: harness.EmbeddedPrefix + "v1.2.3"}); err != nil {
		t.Fatal(err)
	}
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	err = knownCheckout("", home, repo)
	if err == nil || !strings.Contains(err.Error(), "installed from the release binary v1.2.3") {
		t.Errorf("knownCheckout = %v, want a refusal naming the release", err)
	}
	// Named explicitly, it is the operator's decision, as for any checkout.
	if err := knownCheckout(repo, home, repo); err != nil {
		t.Errorf("knownCheckout with --repo: %v", err)
	}
}

// TestReleaseBinaryUpdateShedsKeptClone: a place with no checkout sheds the
// full feedback clone earlier versions kept, exactly as a checkout place does
// (plan § Old kept clone cleanup: "update, both kinds of place"). Only the
// checkout path did it, so the machine that prompted this whole feature — a
// Mac installed from a release binary — would have kept holding every other
// place's reports (found 2026-09-18). Not under --check.
func TestReleaseBinaryUpdateShedsKeptClone(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("SH_AGENT_REPO", "")
	t.Chdir(t.TempDir())
	defer func(v string) { version = v }(version)
	version = "v1.0.0"

	// A gh that answers the one call --check and the up-to-date path make.
	bin := t.TempDir()
	script := "#!/bin/sh" + "\n" + "printf " + "'v1.0.0" + "\\n" + "'" + "\n"
	if err := os.WriteFile(filepath.Join(bin, "gh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	// PATH holds this directory alone, so the real gh is unreachable — but
	// RemoveKeptClone runs git, so git is linked in beside the fake gh rather
	// than adding a system directory that would bring the real gh with it
	// (the mistake caught on 2026-09-17).
	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git not installed")
	}
	if err := os.Symlink(gitBin, filepath.Join(bin, "git")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)

	plant := func(t *testing.T, home string) string {
		t.Helper()
		base := t.TempDir()
		origin := filepath.Join(base, "feedback.git")
		kept := filepath.Join(home, ".config", "sh-agent", "feedback")
		if err := os.MkdirAll(filepath.Dir(kept), 0o755); err != nil {
			t.Fatal(err)
		}
		for _, args := range [][]string{
			{"init", "--quiet", "--bare", "-b", "main", origin},
			{"clone", "--quiet", origin, kept},
			{"-C", kept, "remote", "set-url", "origin", harness.FeedbackURL},
		} {
			cmd := exec.Command(gitBin, args...)
			cmd.Dir = base
			cmd.Env = append(gitCleanEnv(), "GIT_CONFIG_GLOBAL="+filepath.Join(base, "gitconfig"), "GIT_CONFIG_NOSYSTEM=1")
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("git %v: %v\n%s", args, err, out)
			}
		}
		return kept
	}

	// --check changes nothing, here as on the checkout path.
	home := t.TempDir()
	kept := plant(t, home)
	var out bytes.Buffer
	if err := run([]string{"update", "--check", "--home", home}, &out); err != nil {
		t.Fatalf("update --check: %v\n%s", err, out.String())
	}
	if _, err := os.Stat(kept); err != nil {
		t.Errorf("--check removed the kept clone: %v\n%s", err, out.String())
	}

	// The real update sheds it, even when the release is already current.
	out.Reset()
	if err := run([]string{"update", "--home", home}, &out); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "up to date (v1.0.0)") {
		t.Errorf("expected the up-to-date path:\n%s", out.String())
	}
	if _, err := os.Stat(kept); !os.IsNotExist(err) {
		t.Errorf("a binary-only place did not shed the kept clone: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "removed the kept feedback clone") {
		t.Errorf("the removal was silent:\n%s", out.String())
	}
}

// TestReportOnReleaseBinaryIgnoresRepoFlags: `report` reads a source, never
// installs one, so a release binary reports the harness it carries even when
// $SH_AGENT_REPO is exported — otherwise the refusal ResolveSource owes
// `install` would send a record with harness_version, drift, skills and agents
// all null, which is the useless-record shape BUG-077 was about (security
// review 2026-09-18).
func TestReportOnReleaseBinaryIgnoresRepoFlags(t *testing.T) {
	t.Chdir(t.TempDir())
	defer func(v string) { version = v }(version)
	version = "v9.9.9"
	t.Setenv("SH_AGENT_REPO", t.TempDir()) // exported in a shell profile, say

	home := t.TempDir()
	var out bytes.Buffer
	// A feedback path that is not a checkout stops it before any network use;
	// the record's notes are printed before that.
	_ = run([]string{"report", "--home", home, "--dry-run", "--feedback-repo", t.TempDir()}, &out)
	if strings.Contains(out.String(), "is refused") {
		t.Errorf("report refused over an exported $SH_AGENT_REPO:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "release binary v9.9.9") {
		t.Errorf("the record does not name the release binary as its harness:\n%s", out.String())
	}
}

// TestReleaseFromInfo: a release build uses -trimpath, which keeps -ldflags
// out of the build info, so status reads an installed release binary's
// version from the module version Go stamps at a clean tag; anything else
// falls back to the revision.
func TestReleaseFromInfo(t *testing.T) {
	rev := debug.BuildSetting{Key: "vcs.revision", Value: "b64ffe886729abcdef"}
	cases := []struct {
		name string
		info *debug.BuildInfo
		want string
	}{
		{"built at a clean tag", &debug.BuildInfo{Main: debug.Module{Version: "v1.2.3"}, Settings: []debug.BuildSetting{rev, {Key: "vcs.modified", Value: "false"}}}, "v1.2.3"},
		{"go install at a tag", &debug.BuildInfo{Main: debug.Module{Version: "v1.2.3"}}, "v1.2.3"},
		{"pseudo-version", &debug.BuildInfo{Main: debug.Module{Version: "v1.2.4-0.20260917090000-b64ffe886729+dirty"}, Settings: []debug.BuildSetting{rev, {Key: "vcs.modified", Value: "true"}}}, "b64ffe886729+dirty"},
	}
	for _, c := range cases {
		if got := releaseFromInfo(c.info); got != c.want {
			t.Errorf("%s: releaseFromInfo = %q, want %q", c.name, got, c.want)
		}
	}
}

// TestReleaseBinaryWithoutCheckout: a release binary run where no harness tree
// is found installs the files built into it, and every command follows the
// per-command table in the release-binary plan (D13, D14, D15, D17). A
// development build in the same place keeps today's error.
func TestReleaseBinaryWithoutCheckout(t *testing.T) {
	t.Setenv("SH_AGENT_REPO", "")
	t.Setenv(harness.ChildEnv, "")
	t.Chdir(t.TempDir())
	home := t.TempDir()
	var out bytes.Buffer

	if err := run([]string{"install", "--home", home}, &out); err == nil || !strings.Contains(err.Error(), "harness repo not found") {
		t.Fatalf("development build, no tree: err = %v, want harness repo not found", err)
	}
	if _, err := os.Stat(filepath.Join(home, harness.ClaudeRoot)); !os.IsNotExist(err) {
		t.Fatalf("a refused install wrote into home: %v", err)
	}

	defer func(v string) { version = v }(version)
	version = "v9.9.9"

	// A release binary takes no tree at all, so --repo is refused whatever it
	// points at — and refused, not ignored (ADR-034).
	out.Reset()
	if err := run([]string{"install", "--repo", t.TempDir(), "--home", home}, &out); err == nil || !strings.Contains(err.Error(), "never a tree on disk") {
		t.Errorf("release binary, --repo at a non-tree: err = %v, want the refusal", err)
	}
	out.Reset()
	realTree := filepath.Join(t.TempDir(), "harness")
	for rel, body := range map[string]string{
		"context/AGENTS.md": "# ctx\n",
		"skills/s/SKILL.md": "---\nname: s\ndescription: \"x\"\n---\nbody\n",
	} {
		full := filepath.Join(realTree, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := run([]string{"install", "--repo", realTree, "--home", home}, &out); err == nil || !strings.Contains(err.Error(), "never a tree on disk") {
		t.Errorf("release binary, --repo at a real tree: err = %v, want the same refusal", err)
	}

	out.Reset()
	if err := run([]string{"install", "--home", home}, &out); err != nil {
		t.Fatalf("install: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "release binary v9.9.9") || strings.Contains(out.String(), "building") {
		t.Errorf("install output does not name the release binary, or built:\n%s", out.String())
	}
	if _, err := os.Stat(filepath.Join(home, harness.BinPath)); err != nil {
		t.Errorf("the running binary was not installed: %v", err)
	}
	st, found, err := harness.LoadState(home)
	if err != nil || !found || st.Repo != "embedded:v9.9.9" || st.Commit != "" {
		t.Errorf("state = %+v (found %v, err %v), want repo embedded:v9.9.9 and no commit", st, found, err)
	}

	out.Reset()
	if err := run([]string{"status", "--home", home}, &out); err != nil {
		t.Fatalf("status: %v\n%s", err, out.String())
	}
	for _, want := range []string{
		"source     embedded v9.9.9 (release binary; update downloads the newest release)",
		"from release v9.9.9",
		"binary     ",
		"drift      none",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("status lacks %q:\n%s", want, out.String())
		}
	}
	if strings.Contains(out.String(), "repo       ") {
		t.Errorf("status prints a repo line for an embedded install:\n%s", out.String())
	}

	out.Reset()
	if err := run([]string{"diff", "--home", home}, &out); err != nil {
		t.Errorf("diff: %v\n%s", err, out.String())
	}
	out.Reset()
	if err := run([]string{"list", "--home", home}, &out); err != nil || !strings.Contains(out.String(), "### Skills") {
		t.Errorf("list: %v\n%s", err, out.String())
	}
	out.Reset()
	if err := run([]string{"render", "--home", home}, &out); err != nil || !strings.HasPrefix(out.String(), harness.Marker) {
		t.Errorf("render: %v\n%.200s", err, out.String())
	}
	out.Reset()
	if err := run([]string{"adopt", "anything", "--home", home}, &out); err == nil || err.Error() != "adopt needs a harness checkout; this install is a release binary" {
		t.Errorf("adopt: err = %v", err)
	}
	// `update` on a release binary goes to the release flow (Phase 4). PATH is
	// emptied for this call so the real gh is never reached: an earlier version
	// of this test called the real repository over the network (2026-09-17).
	out.Reset()
	func() {
		t.Setenv("PATH", t.TempDir())
		if err := run([]string{"update", "--home", home}, &out); err == nil || !strings.Contains(err.Error(), "install gh and run gh auth login") {
			t.Errorf("update with no gh: err = %v, want the gh refusal", err)
		}
	}()
	out.Reset()
	// The record's notes print before the feedback checkout is resolved, so
	// they show whether or not this non-checkout feedback path is refused.
	_ = run([]string{"report", "--home", home, "--dry-run", "--feedback-repo", t.TempDir()}, &out)
	if !strings.Contains(out.String(), "release binary v9.9.9") {
		t.Errorf("report does not say the harness is a release binary:\n%s", out.String())
	}
}

// TestReleaseBinaryUpdateCheckExitsTwo: on a release binary, `update --check`
// reports a newer release and exits 2, the same code a checkout that is behind
// gives (plan § Per command with an embedded source).
func TestReleaseBinaryUpdateCheckExitsTwo(t *testing.T) {
	t.Setenv("SH_AGENT_REPO", "")
	t.Chdir(t.TempDir())
	defer func(v string) { version = v }(version)
	version = "v1.0.0"

	// A gh that answers the one call --check makes, and nothing else on PATH.
	bin := t.TempDir()
	gh := filepath.Join(bin, "gh")
	script := "#!/bin/sh" + "\n" + "printf " + "'v9.9.9" + "\\n" + "'" + "\n"
	if err := os.WriteFile(gh, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)

	var out bytes.Buffer
	err := run([]string{"update", "--check", "--home", t.TempDir()}, &out)
	var code exitCode
	if !errors.As(err, &code) || int(code) != 2 {
		t.Fatalf("err = %v, want exit code 2", err)
	}
	if !strings.Contains(out.String(), "newer release v9.9.9 available") {
		t.Errorf("output: %q", out.String())
	}
}

// TestHarnessMap pins the map in README.md § How it works to the binary: the
// map is the contract for install, update, report back and release (R9), so it
// must exist, name only verbs the binary has, and carry the three verbs those
// procedures run through. A verb renamed or removed fails here, not silently in
// a picture nobody re-reads.
func TestHarnessMap(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	_, section, ok := strings.Cut(string(data), "\n## How it works\n")
	if !ok {
		t.Fatal("README.md has no § How it works; the map is required by R9")
	}
	section, _, _ = strings.Cut(section, "\n## ")
	_, block, ok := strings.Cut(section, "```mermaid\n")
	if !ok {
		t.Fatal("README.md § How it works has no mermaid map")
	}
	block, _, _ = strings.Cut(block, "```")
	seen := map[string]bool{}
	for _, m := range regexp.MustCompile(`sh-agent ([a-z]+)`).FindAllStringSubmatch(block, -1) {
		seen[m[1]] = true
		if !strings.Contains(usage, "\n  "+m[1]+" ") {
			t.Errorf("the map names `sh-agent %s`, which the binary does not have", m[1])
		}
	}
	for _, verb := range []string{"install", "update", "report"} {
		if !seen[verb] {
			t.Errorf("the map does not show `sh-agent %s`", verb)
		}
	}
}

// A verified release binary was installed to /usr/local/bin and then ignored,
// because two earlier PATH entries held other sh-agents; install said it
// succeeded and status said drift none, both about a binary that was not the
// one the shell ran (BUG-078). install warns and status carries a `shadowed`
// line, and both go quiet when the earlier entry is a link to the same file.
func TestInstallAndStatusNameAPathShadow(t *testing.T) {
	t.Setenv(harness.ChildEnv, "1")
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	installedBin := filepath.Join(home, harness.BinPath)

	// An earlier PATH entry holding a different sh-agent.
	earlier := t.TempDir()
	other := filepath.Join(earlier, "sh-agent")
	if err := os.WriteFile(other, []byte("#!/bin/sh\necho not-the-installed-one\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", earlier+string(os.PathListSeparator)+filepath.Dir(installedBin))

	var out bytes.Buffer
	if err := run([]string{"install", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("install: %v\n%s", err, out.String())
	}
	wantInstall := "warning: " + other + " comes before " + installedBin + " on PATH, so `sh-agent` runs that one"
	if !strings.Contains(out.String(), wantInstall) {
		t.Errorf("install missing %q:\n%s", wantInstall, out.String())
	}

	out.Reset()
	if err := run([]string{"status", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("status: %v\n%s", err, out.String())
	}
	wantStatus := "shadowed   " + other + " comes first on PATH, so `sh-agent` runs that, not the binary above"
	if !strings.Contains(out.String(), wantStatus) {
		t.Errorf("status missing %q:\n%s", wantStatus, out.String())
	}

	// The fix the operator applied: point the earlier entry at the installed
	// binary. It runs the same file, so neither command may still warn.
	if err := os.Remove(other); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(installedBin, other); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err := run([]string{"status", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("status after symlink: %v\n%s", err, out.String())
	}
	if strings.Contains(out.String(), "shadowed") {
		t.Errorf("status still reports a shadow for a link to the installed binary:\n%s", out.String())
	}
}

// One situation, one message. When this home's bin dir is not on PATH at all,
// pathHint owns it and names what answers instead; the shadow warning stays
// quiet, so a scratch-home install — which this repo does constantly — does
// not print a warning the operator learns to skip past.
func TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice(t *testing.T) {
	t.Setenv(harness.ChildEnv, "1")
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	binDir := filepath.Join(home, filepath.Dir(harness.BinPath))

	elsewhere := t.TempDir()
	other := filepath.Join(elsewhere, "sh-agent")
	if err := os.WriteFile(other, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", elsewhere) // this home's bin dir is NOT on PATH

	var out bytes.Buffer
	if err := run([]string{"install", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("install: %v\n%s", err, out.String())
	}
	want := "note: " + binDir + " is not on PATH, so `sh-agent` runs " + other + " instead; put " + binDir + " ahead of it"
	if !strings.Contains(out.String(), want) {
		t.Errorf("install missing %q:\n%s", want, out.String())
	}
	// Not the generic "warning:" — setting PATH to one directory also drops
	// `go`, and install rightly warns about that. Only the shadow warning
	// must be absent.
	if strings.Contains(out.String(), "comes before") {
		t.Errorf("install warned twice about one situation:\n%s", out.String())
	}

	out.Reset()
	if err := run([]string{"status", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("status: %v\n%s", err, out.String())
	}
	if strings.Contains(out.String(), "shadowed") {
		t.Errorf("status reported a shadow for a home whose bin dir is off PATH:\n%s", out.String())
	}

	// A dry run installs nothing, so it says nothing about PATH: on a fresh
	// home it would otherwise call the bin dir "not on PATH" whether it is or
	// not, a false diagnosis from a read-only mode.
	fresh := t.TempDir()
	out.Reset()
	if err := run([]string{"install", "--dry-run", "--repo", repo, "--home", fresh}, &out); err != nil {
		t.Fatalf("install --dry-run: %v\n%s", err, out.String())
	}
	// The bin-dir phrases only: PATH here is one directory, so install rightly
	// warns that `go` is not on it, and that warning is not what this pins.
	for _, claim := range []string{"add it so `sh-agent` is available everywhere", "so `sh-agent` runs", "comes before"} {
		if strings.Contains(out.String(), claim) {
			t.Errorf("install --dry-run made a PATH claim about a binary it did not write (%q):\n%s", claim, out.String())
		}
	}
}

// `update` replaces ~/.local/bin/sh-agent and nothing else. On devbox it was
// run by /usr/local/bin/sh-agent (v1.2.1), which stayed v1.2.1 while
// ~/.local/bin became v1.2.2, and bash's command hash kept running the old
// file — `which -a` showed the right order and `sh-agent version` printed the
// wrong version. update now says when the binary it ran from is not the one
// it installed and a shell can reach it by name (BUG-079).
func TestUpdateNamesTheBinaryItRanFrom(t *testing.T) {
	t.Setenv("SH_AGENT_REPO", "")
	t.Chdir(t.TempDir())
	defer func(v string) { version = v }(version)
	version = "v1.0.0"

	// A gh that answers the one call the up-to-date path makes.
	ghDir := t.TempDir()
	script := "#!/bin/sh" + "\n" + "printf " + "'v1.0.0" + "\\n" + "'" + "\n"
	if err := os.WriteFile(filepath.Join(ghDir, "gh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	// A PATH directory whose sh-agent IS this process's executable — the
	// shape of /usr/local/bin/sh-agent on devbox.
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	elsewhere := t.TempDir()
	byName := filepath.Join(elsewhere, "sh-agent")
	if err := os.Symlink(self, byName); err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	installed := filepath.Join(home, harness.BinPath)

	want := "note: this update ran from " + byName + ", not from the installed binary " + installed + "; update replaces only the installed one, so the file it ran from stays as it is. A shell that has already resolved `sh-agent` to that path keeps running it until `hash -r` (bash) or `rehash` (zsh); point it at the installed binary (`ln -sfn`) or remove it"

	t.Run("no installed binary: silent", func(t *testing.T) {
		// An up-to-date update into a home that has no binary installed
		// nothing; naming a file that does not exist and calling the running
		// one stale would send the reader to remove their only binary.
		t.Setenv("PATH", ghDir+string(os.PathListSeparator)+elsewhere)
		var out bytes.Buffer
		if err := run([]string{"update", "--home", home}, &out); err != nil {
			t.Fatalf("update: %v\n%s", err, out.String())
		}
		if strings.Contains(out.String(), "not from the installed binary") {
			t.Errorf("update named an installed binary that does not exist:\n%s", out.String())
		}
	})

	// From here on the home holds an installed binary — a different file
	// from this process — as after any install.
	if err := os.MkdirAll(filepath.Dir(installed), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(installed, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	t.Run("reachable by name: named", func(t *testing.T) {
		t.Setenv("PATH", ghDir+string(os.PathListSeparator)+elsewhere)
		var out bytes.Buffer
		if err := run([]string{"update", "--home", home}, &out); err != nil {
			t.Fatalf("update: %v\n%s", err, out.String())
		}
		if !strings.Contains(out.String(), want) {
			t.Errorf("update did not name the binary it ran from:\nwant %q\n%s", want, out.String())
		}
	})

	t.Run("not on PATH at all: silent", func(t *testing.T) {
		// A test binary or a `go run` build is never reachable by name.
		t.Setenv("PATH", ghDir)
		var out bytes.Buffer
		if err := run([]string{"update", "--home", home}, &out); err != nil {
			t.Fatalf("update: %v\n%s", err, out.String())
		}
		if strings.Contains(out.String(), "not from the installed binary") {
			t.Errorf("update named a binary no shell can reach by name:\n%s", out.String())
		}
	})

	t.Run("--check says nothing", func(t *testing.T) {
		t.Setenv("PATH", ghDir+string(os.PathListSeparator)+elsewhere)
		var out bytes.Buffer
		if err := run([]string{"update", "--check", "--home", home}, &out); err != nil {
			t.Fatalf("update --check: %v\n%s", err, out.String())
		}
		if strings.Contains(out.String(), "not from the installed binary") {
			t.Errorf("--check, which changes nothing, described a replacement:\n%s", out.String())
		}
	})
}

// The checkout path of update has the same shape: Reinstall writes
// ~/.local/bin/sh-agent, and the process that ran update — a copy at a
// system path, or a `go run` build — is not that file (BUG-079).
func TestCheckoutUpdateNamesTheBinaryItRanFrom(t *testing.T) {
	repo, _ := gitFixture(t)
	home := t.TempDir()
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	elsewhere := t.TempDir()
	byName := filepath.Join(elsewhere, "sh-agent")
	if err := os.Symlink(self, byName); err != nil {
		t.Fatal(err)
	}
	// Appended, so git and go stay reachable; no earlier entry holds THIS
	// executable under the name, so the report is the appended one.
	t.Setenv("PATH", os.Getenv("PATH")+string(os.PathListSeparator)+elsewhere)
	installed := filepath.Join(home, harness.BinPath)
	// The fixture's built binary writes a marker, not ~/.local/bin/sh-agent;
	// the note requires an installed file to be behind, so plant one — a
	// different file from this process, as after any real install.
	if err := os.MkdirAll(filepath.Dir(installed), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(installed, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	want := "note: this update ran from " + byName + ", not from the installed binary " + installed

	var out bytes.Buffer
	if err := update(repo, home, true, "", "", &out); err != nil {
		t.Fatalf("update --check: %v\n%s", err, out.String())
	}
	if strings.Contains(out.String(), "not from the installed binary") {
		t.Errorf("--check, which changes nothing, described a replacement:\n%s", out.String())
	}
	out.Reset()
	if err := update(repo, home, false, "", "", &out); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), want) {
		t.Errorf("checkout update did not name the binary it ran from:\nwant %q\n%s", want, out.String())
	}
}
