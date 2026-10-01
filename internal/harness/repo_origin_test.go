package harness

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseOrigin: every shape the installer writes into the manifest's repo
// value is read by one parser (IN6), so no reader cuts prefixes by hand.
func TestParseOrigin(t *testing.T) {
	cases := []struct {
		repo string
		want Origin
	}{
		{"/home/x/styrman-harness", Origin{Kind: OriginPath, Path: "/home/x/styrman-harness"}},
		{EmbeddedPrefix + "v1.17.0", Origin{Kind: OriginRelease, Release: "v1.17.0"}},
		{ZipOrigin("abcdef0123456789abcdef", ""), Origin{Kind: OriginZip, Sum: "abcdef0123456789abcdef", Trust: "unverified"}},
		{ZipOrigin("abcdef0123456789abcdef", "v1.16.2"), Origin{Kind: OriginZip, Sum: "abcdef0123456789abcdef", Trust: "verified v1.16.2"}},
		// a manifest from before the verdict was recorded has no ";…" part
		{ZipPrefix + "sha256:abc", Origin{Kind: OriginZip, Sum: "abc"}},
		{"", Origin{Kind: OriginPath}},
	}
	for _, c := range cases {
		if got := ParseOrigin(c.repo); got != c.want {
			t.Errorf("ParseOrigin(%q) = %+v, want %+v", c.repo, got, c.want)
		}
		if got := (State{Repo: c.repo}).Origin(); got != c.want {
			t.Errorf("State.Origin(%q) = %+v, want %+v", c.repo, got, c.want)
		}
	}
	if s := ParseOrigin(ZipOrigin("abcdef0123456789abcdef", "")).ShortSum(); s != "abcdef012345" {
		t.Errorf("ShortSum = %q, want the first 12 hex digits", s)
	}
	if s := (Origin{Sum: "abc"}).ShortSum(); s != "abc" {
		t.Errorf("ShortSum of a short sum = %q", s)
	}
}

// TestSourceOriginString: what Install records for each kind of source.
func TestSourceOriginString(t *testing.T) {
	if got := DirSource("/tree").OriginString(); got != "/tree" {
		t.Errorf("tree: %q", got)
	}
	if got := EmbeddedSource("v1.2.3").OriginString(); got != EmbeddedPrefix+"v1.2.3" {
		t.Errorf("embedded: %q", got)
	}
	zip := Source{Dir: "/unpacked", Origin: ZipOrigin("ff00", "v1.2.3")}
	if got := zip.OriginString(); got != ZipPrefix+"sha256:ff00;verified v1.2.3" {
		t.Errorf("zip: %q", got)
	}
	// and the round trip: Install's record is read back as the same kind
	for _, src := range []Source{DirSource("/tree"), EmbeddedSource("v1.2.3"), zip} {
		o := ParseOrigin(src.OriginString())
		switch {
		case src.Origin != "" && o.Kind != OriginZip,
			src.Origin == "" && src.Dir == "" && o.Kind != OriginRelease,
			src.Origin == "" && src.Dir != "" && o.Kind != OriginPath:
			t.Errorf("round trip of %+v gave %+v", src, o)
		}
	}
}

// TestResolveSourceLenient: report's source resolution (IN7). A release
// binary reports its embedded files and notes a --repo it ignores; a
// development build with no checkout anywhere gets the zero Source and the
// reason as a note, never an error.
func TestResolveSourceLenient(t *testing.T) {
	t.Setenv("SH_AGENT_REPO", "")
	home := t.TempDir()
	src, note := ResolveSourceLenient("/elsewhere", home, "v9.9.9")
	if src.Version != "v9.9.9" || src.FS == nil {
		t.Errorf("release: want the embedded source, got %+v", src)
	}
	if !strings.Contains(note, "--repo is not used here") {
		t.Errorf("release with --repo: note %q", note)
	}
	if _, note := ResolveSourceLenient("", home, "v9.9.9"); note != "" {
		t.Errorf("release without --repo: note %q, want none", note)
	}
	// a development build, run from a directory that is no checkout, with no manifest
	wd, _ := os.Getwd()
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(wd) })
	src, note = ResolveSourceLenient("", home, "(devel)")
	if src.FS != nil || src.Dir != "" {
		t.Errorf("no checkout: want the zero Source, got %+v", src)
	}
	if !strings.HasPrefix(note, "report: ") || !strings.Contains(note, "harness repo not found") {
		t.Errorf("no checkout: note %q", note)
	}
	// a named checkout is found the same way ResolveSource finds it
	repo := filepath.Join(t.TempDir(), "h")
	for _, d := range []string{"context", "skills"} {
		if err := os.MkdirAll(filepath.Join(repo, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(repo, "context", "AGENTS.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if src, note = ResolveSourceLenient(repo, home, "(devel)"); src.Dir != repo || note != "" {
		t.Errorf("named checkout: %+v, note %q", src, note)
	}
}
