package harness

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadSettingsProblems(t *testing.T) {
	dir := t.TempDir()

	if _, problem := ReadSettings(filepath.Join(dir, "absent.json")); problem != "could not be read" {
		t.Errorf("missing file: %q", problem)
	}

	sub := filepath.Join(dir, "adirectory.json")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, problem := ReadSettings(sub); problem != "is not a regular file" {
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
	data, problem := ReadSettings(unreadable)
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

// TestSettingsWriteUsesHarnessWriteAtomic pins the settings editor to the one
// atomic-write helper: a second temp-write-sync-rename sequence is a fix to
// port by hand (project review 2026-09-30, simplify finding 1).
func TestSettingsWriteUsesHarnessWriteAtomic(t *testing.T) {
	src, err := os.ReadFile("settings.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(src), "os.CreateTemp(") {
		t.Error("settings.go writes through its own temp file; use WriteAtomic")
	}
	if !strings.Contains(string(src), "WriteAtomic(") {
		t.Error("settings.go does not call WriteAtomic")
	}
}

// TestAttributionSetReportsAnUnwritableSettingsFile drives the error branch
// behind harness.WriteAtomic: a settings directory that cannot take a temp
// file makes the edit fail, and the failure is reported with its cause and
// the by-hand remedy, the file left as it was (v1.13.1 gate review).
func TestAttributionSetReportsAnUnwritableSettingsFile(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores the directory permissions this test relies on")
	}
	home := t.TempDir()
	dir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "settings.json")
	if err := os.WriteFile(path, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o700) })
	var out bytes.Buffer
	AttributionSet(home, false, &out)
	got := out.String()
	if !strings.Contains(got, "could not be written") || !strings.Contains(got, "add the key by hand") {
		t.Errorf("the failure should be reported with the remedy, got %q", got)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "{}\n" {
		t.Errorf("the settings file should be untouched, got %q %v", data, err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Errorf("no temp file should be left behind, got %d entries", len(entries))
	}
}

// TestAttributionSetAndHint drives the settings editor in package harness
// (S1): the file is created when absent, the key inserted into an existing
// object with every other byte kept, a present key left alone, a dry run
// touches nothing, a symlink is never written through, and the hint names
// each state of the file in one line.
func TestAttributionSetAndHint(t *testing.T) {
	const fragment = `"attribution": {"commit": "", "pr": ""}`
	settings := func(home string) string { return filepath.Join(home, ClaudeRoot, "settings.json") }
	write := func(t *testing.T, home, content string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Join(home, ClaudeRoot), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(settings(home), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	read := func(t *testing.T, home string) string {
		t.Helper()
		data, err := os.ReadFile(settings(home))
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}

	t.Run("absent file is created", func(t *testing.T) {
		home := t.TempDir()
		var out bytes.Buffer
		AttributionSet(home, false, &out)
		if got := read(t, home); got != "{\n  "+fragment+"\n}\n" {
			t.Errorf("created file: %q", got)
		}
		if !strings.Contains(out.String(), "created with") {
			t.Errorf("no creation line: %q", out.String())
		}
		out.Reset()
		AttributionHint(home, "note: ", &out)
		if out.Len() != 0 {
			t.Errorf("hint after the key is set: %q", out.String())
		}
	})
	t.Run("dry run writes nothing", func(t *testing.T) {
		home := t.TempDir()
		var out bytes.Buffer
		AttributionSet(home, true, &out)
		if _, err := os.Stat(settings(home)); err == nil {
			t.Error("dry run created the file")
		}
		if !strings.Contains(out.String(), "[dry-run] would create") {
			t.Errorf("dry run line: %q", out.String())
		}
		write(t, home, "{\n  \"theme\": \"dark\"\n}\n")
		out.Reset()
		AttributionSet(home, true, &out)
		if !strings.Contains(out.String(), "[dry-run] would add") || read(t, home) != "{\n  \"theme\": \"dark\"\n}\n" {
			t.Errorf("dry run on an existing file: %q / %q", out.String(), read(t, home))
		}
	})
	t.Run("key inserted, other bytes kept, CRLF kept", func(t *testing.T) {
		home := t.TempDir()
		write(t, home, "{\r\n  \"theme\": \"dark\"\r\n}\r\n")
		var out bytes.Buffer
		AttributionSet(home, false, &out)
		got := read(t, home)
		if !strings.Contains(got, "\"theme\": \"dark\",\r\n  "+fragment+"\r\n}") {
			t.Errorf("insertion: %q", got)
		}
		if !strings.Contains(out.String(), "added; other keys untouched") {
			t.Errorf("line: %q", out.String())
		}
	})
	t.Run("empty object", func(t *testing.T) {
		home := t.TempDir()
		write(t, home, "{}\n")
		AttributionSet(home, false, &bytes.Buffer{})
		if got := read(t, home); got != "{\n  "+fragment+"\n}\n" {
			t.Errorf("empty object: %q", got)
		}
	})
	t.Run("present key of any value is the operator's", func(t *testing.T) {
		home := t.TempDir()
		const custom = "{\"attribution\": {\"commit\": \"Co-Authored-By: me\"}}\n"
		write(t, home, custom)
		var out bytes.Buffer
		AttributionSet(home, false, &out)
		if read(t, home) != custom || out.Len() != 0 {
			t.Errorf("present key: file %q, out %q", read(t, home), out.String())
		}
		AttributionHint(home, "note: ", &out)
		if !strings.Contains(out.String(), "does not set \"attribution\": {\"commit\": \"\"}") {
			t.Errorf("hint for a custom trailer: %q", out.String())
		}
	})
	t.Run("off with false, hint silent", func(t *testing.T) {
		home := t.TempDir()
		write(t, home, "{\"attribution\": {\"commit\": false}}\n")
		var out bytes.Buffer
		AttributionHint(home, "note: ", &out)
		if out.Len() != 0 {
			t.Errorf("hint: %q", out.String())
		}
	})
	t.Run("malformed and non-object files are advisory", func(t *testing.T) {
		// a bare null is a nil map for AttributionSet (left alone) and a zero
		// struct for the hint, which then reports the key as unset
		for content, want := range map[string]string{"{\"attribution\": \n": "is not valid JSON", "[1, 2]\n": "is not a JSON object", "null\n": "does not set \"attribution\""} {
			home := t.TempDir()
			write(t, home, content)
			var out bytes.Buffer
			AttributionSet(home, false, &out)
			if read(t, home) != content {
				t.Errorf("%q: file changed to %q", content, read(t, home))
			}
			AttributionHint(home, "note: ", &out)
			if !strings.Contains(out.String(), want) {
				t.Errorf("%q: hint %q, want %q", content, out.String(), want)
			}
			if strings.Contains(out.String(), "invalid character") {
				t.Errorf("%q: the JSON error is echoed: %q", content, out.String())
			}
		}
	})
	t.Run("symlink is never written through", func(t *testing.T) {
		home := t.TempDir()
		target := filepath.Join(t.TempDir(), "real.json")
		if err := os.WriteFile(target, []byte("{}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(home, ClaudeRoot), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, settings(home)); err != nil {
			t.Fatal(err)
		}
		var out bytes.Buffer
		AttributionSet(home, false, &out)
		if data, _ := os.ReadFile(target); string(data) != "{}\n" {
			t.Errorf("target written through the link: %q", data)
		}
		if !strings.Contains(out.String(), "is a symlink") {
			t.Errorf("line: %q", out.String())
		}
	})
	t.Run("oversized file", func(t *testing.T) {
		home := t.TempDir()
		big := append([]byte("{"), bytes.Repeat([]byte(" "), 1<<20)...)
		write(t, home, string(append(big, '}')))
		var out bytes.Buffer
		AttributionHint(home, "note: ", &out)
		if !strings.Contains(out.String(), "larger than 1 MiB") {
			t.Errorf("oversized: %q", out.String())
		}
	})
}

// A linked ~/.claude (a dotfiles checkout) is never written through: the
// installer refuses every other path under a linked directory (underLink,
// BUG-129's rule), and the settings file is no exception.
func TestAttributionSetRefusesALinkedClaudeDir(t *testing.T) {
	home := t.TempDir()
	target := t.TempDir()
	if err := os.Symlink(target, filepath.Join(home, ClaudeRoot)); err != nil {
		t.Skip("symlinks unavailable:", err)
	}
	var out bytes.Buffer
	AttributionSet(home, false, &out)
	if _, err := os.Lstat(filepath.Join(target, "settings.json")); err == nil {
		t.Fatalf("settings.json was written through the linked %s into %s", ClaudeRoot, target)
	}
	if !strings.Contains(out.String(), "symlink") {
		t.Fatalf("the refusal does not name the link: %q", out.String())
	}
	// The edit branch too: an existing file behind the link is left byte for byte.
	existing := []byte("{\n  \"theme\": \"dark\"\n}\n")
	if err := os.WriteFile(filepath.Join(target, "settings.json"), existing, 0o644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	AttributionSet(home, false, &out)
	got, _ := os.ReadFile(filepath.Join(target, "settings.json"))
	if string(got) != string(existing) {
		t.Fatalf("settings.json behind the link was edited: %q", got)
	}
	if !strings.Contains(out.String(), "symlink") {
		t.Fatalf("the refusal does not name the link on the edit path: %q", out.String())
	}
}
