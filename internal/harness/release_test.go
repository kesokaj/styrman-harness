package harness

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// The fake gh: enough of the two calls UpdateFromRelease makes, driven by the
// environment so each case can bend one thing. Absolute paths throughout,
// because the cases that test a missing tool run with a PATH that holds
// almost nothing.
const fakeGh = `#!/bin/sh
if [ "$1" = "release" ] && [ "$2" = "view" ]; then
  if [ -n "$GH_FAIL_VIEW" ]; then
    echo "gh: HTTP 401: Bad credentials (https://api.github.com/graphql)" >&2
    exit 1
  fi
  printf '%s\n' "$GH_TAG"
  exit 0
fi
if [ "$1" = "release" ] && [ "$2" = "download" ]; then
  dest=""
  while [ $# -gt 0 ]; do
    if [ "$1" = "--dir" ]; then shift; dest="$1"; fi
    shift
  done
  n=0
  for f in "$GH_ASSETS"/*; do
    [ -e "$f" ] || continue
    /bin/cp "$f" "$dest/" && n=$((n+1))
  done
  if [ "$n" -eq 0 ]; then
    echo "gh: no assets match the file pattern" >&2
    exit 1
  fi
  exit 0
fi
echo "fake gh: unexpected call: $*" >&2
exit 1
`

// The release binary this fixture ships: a script that records the arguments
// it was installed with, so the test can tell that the NEW binary ran.
const fakeReleaseBinary = `#!/bin/sh
printf 'install %s\n' "$*" >> "$SH_AGENT_FAKE_MARKER"
`

type release struct {
	assets string // what gh "downloads" from
	tag    string
	signer string // allowed-signers line for the key that signed SHA256SUMS
	key    string // that key, for a case that re-signs a tampered file
}

// newRelease builds a signed release: the platform's binary, SHA256SUMS over
// it, and a real ssh-keygen signature over that file.
func newRelease(t *testing.T, tag string) release {
	t.Helper()
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		t.Skip("ssh-keygen not installed")
	}
	dir := t.TempDir()
	keys := t.TempDir()
	key := filepath.Join(keys, "key")
	mustRun(t, "", "ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-C", "test@release", "-f", key)
	pub, err := os.ReadFile(key + ".pub")
	if err != nil {
		t.Fatal(err)
	}
	fields := strings.Fields(string(pub))
	if len(fields) < 2 {
		t.Fatalf("unexpected public key %q", pub)
	}

	asset := AssetName(runtime.GOOS, runtime.GOARCH)
	writeFile(t, filepath.Join(dir, asset), fakeReleaseBinary)
	signRelease(t, dir, asset, key, tag)
	return release{assets: dir, tag: tag, key: key, signer: "test@release " + fields[0] + " " + fields[1]}
}

// signRelease writes SHA256SUMS over the asset as it is on disk now, and
// signs it. Called again after changing the asset when a case needs the
// checksum to match, and NOT called again when it must not.
func signRelease(t *testing.T, dir, asset, key, tag string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, asset))
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	// The marker /preflight writes, which puts the VERSION inside the signed
	// bytes (plan D26).
	writeFile(t, filepath.Join(dir, "SHA256SUMS"),
		"# styrman-harness "+tag+"\n"+hex.EncodeToString(sum[:])+"  "+asset+"\n")
	_ = os.Remove(filepath.Join(dir, "SHA256SUMS.sig"))
	mustRun(t, dir, "ssh-keygen", "-Y", "sign", "-q", "-f", key, "-n", ReleaseNamespace, "SHA256SUMS")
}

func mustRun(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out)
	}
}

// releaseEnv puts the fake gh on PATH, points it at this release, and
// compiles in the test key's signer line. withSSHKeygen false leaves
// ssh-keygen off PATH; withGh false leaves gh off it.
func (r release) env(t *testing.T, withGh, withSSHKeygen bool) {
	t.Helper()
	bin := t.TempDir()
	if withGh {
		gh := filepath.Join(bin, "gh")
		writeFile(t, gh, fakeGh)
		if err := os.Chmod(gh, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	// PATH holds this directory ALONE, so "gh missing" means missing: adding
	// the real ssh-keygen's directory would have brought the real gh with it,
	// and the first run of this test called the real repository (caught
	// 2026-09-17). ssh-keygen is linked in instead when a case wants it.
	if withSSHKeygen {
		real, err := exec.LookPath("ssh-keygen")
		if err != nil {
			t.Skip("ssh-keygen not installed")
		}
		if err := os.Symlink(real, filepath.Join(bin, "ssh-keygen")); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", bin)
	t.Setenv("GH_TAG", r.tag)
	t.Setenv("GH_ASSETS", r.assets)
	t.Setenv("GH_FAIL_VIEW", "")
	os.Unsetenv("GH_FAIL_VIEW")

	old := releaseSigner
	releaseSigner = r.signer
	t.Cleanup(func() { releaseSigner = old })
}

// installedBinary is the sentinel this machine already has, so every refusal
// can be checked against "the binary was not replaced".
func installedBinary(t *testing.T, home string) string {
	t.Helper()
	dest := filepath.Join(home, BinPath)
	writeFile(t, dest, "#!/bin/sh\necho the old binary\n")
	if err := os.Chmod(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	return dest
}

func TestUpdateFromReleaseReplacesAndInstalls(t *testing.T) {
	r := newRelease(t, "v9.9.9")
	r.env(t, true, true)
	home := t.TempDir()
	dest := installedBinary(t, home)
	marker := filepath.Join(home, "installed.txt")
	t.Setenv("SH_AGENT_FAKE_MARKER", marker)

	var out bytes.Buffer
	// The process running this update is not the file it installs — the
	// devbox shape — and a PATH entry reaches it by name (BUG-079).
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	elsewhere := t.TempDir()
	byName := filepath.Join(elsewhere, "sh-agent")
	if err := os.Symlink(self, byName); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", os.Getenv("PATH")+string(os.PathListSeparator)+elsewhere)
	if err := UpdateFromRelease(home, "v1.0.0", false, []string{"--home", home}, &out); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}
	wantNote := "note: this update ran from " + byName + ", not from the installed binary " + dest + "; update replaces only the installed one"
	if !strings.Contains(out.String(), wantNote) {
		t.Errorf("a replacing update did not name the binary it ran from:\nwant %q\n%s", wantNote, out.String())
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != fakeReleaseBinary {
		t.Errorf("the installed binary is not the release's:\n%s", got)
	}
	if info, err := os.Stat(dest); err != nil || info.Mode().Perm()&0o100 == 0 {
		t.Errorf("the installed binary is not executable: %v %v", info.Mode(), err)
	}
	ran, err := os.ReadFile(marker)
	if err != nil || !strings.Contains(string(ran), "install --home "+home) {
		t.Errorf("the new binary's install did not run: %q %v", ran, err)
	}
	for _, want := range []string{"downloading sh-agent-", "verified and installed sh-agent-", "v9.9.9"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, out.String())
		}
	}
	// The download directory is gone, whatever happened.
	left, _ := filepath.Glob(filepath.Join(home, ".config", "sh-agent", "release-*"))
	if len(left) != 0 {
		t.Errorf("left behind: %v", left)
	}
}

func TestUpdateFromReleaseUpToDateAndCheck(t *testing.T) {
	r := newRelease(t, "v1.0.0")
	r.env(t, true, true)
	home := t.TempDir()
	dest := installedBinary(t, home)
	before, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}

	// The newest release IS this binary: nothing is downloaded or replaced,
	// on either mode.
	for _, checkOnly := range []bool{false, true} {
		var out bytes.Buffer
		if err := UpdateFromRelease(home, "v1.0.0", checkOnly, nil, &out); err != nil {
			t.Fatalf("checkOnly=%v: %v", checkOnly, err)
		}
		if !strings.Contains(out.String(), "up to date (v1.0.0)") {
			t.Errorf("checkOnly=%v output: %q", checkOnly, out.String())
		}
	}

	// A newer release, --check: named, exit-2 sentinel, nothing replaced.
	t.Setenv("GH_TAG", "v9.9.9")
	var out bytes.Buffer
	err = UpdateFromRelease(home, "v1.0.0", true, nil, &out)
	if !errors.Is(err, ErrNewerRelease) {
		t.Errorf("--check with a newer release: err = %v, want ErrNewerRelease", err)
	}
	if !strings.Contains(out.String(), "newer release v9.9.9 available") {
		t.Errorf("--check output: %q", out.String())
	}
	if after, _ := os.ReadFile(dest); string(after) != string(before) {
		t.Error("--check replaced the binary")
	}
}

// TestUpdateFromReleaseRefusals: every way the release can fail to prove it
// is the operator's leaves the installed binary exactly as it was, names the
// check that failed, and keeps nothing (plan § Failure modes).
func TestUpdateFromReleaseRefusals(t *testing.T) {
	for _, tc := range []struct {
		name    string
		want    string
		gh      bool
		keygen  bool
		prepare func(t *testing.T, r release)
	}{
		{
			name: "gh missing", want: "install gh and run gh auth login",
			gh: false, keygen: true,
		},
		{
			name: "gh cannot reach the repository", want: "gh release view failed: gh: HTTP 401",
			gh: true, keygen: true,
			prepare: func(t *testing.T, _ release) { t.Setenv("GH_FAIL_VIEW", "1") },
		},
		{
			// Neither a system copy nor one on PATH: nothing can verify the
			// release, so nothing is downloaded.
			name: "ssh-keygen missing", want: "ssh-keygen not found; cannot verify the release",
			gh: true, keygen: false,
			prepare: func(t *testing.T, _ release) {
				old := sshKeygenSystemPaths
				sshKeygenSystemPaths = []string{filepath.Join(t.TempDir(), "absent")}
				t.Cleanup(func() { sshKeygenSystemPaths = old })
			},
		},
		{
			name: "no asset for this platform", want: "has no sh-agent-",
			gh: true, keygen: true,
			prepare: func(t *testing.T, r release) {
				if err := os.Remove(filepath.Join(r.assets, AssetName(runtime.GOOS, runtime.GOARCH))); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "signature made by another key", want: "SHA256SUMS signature does not verify",
			gh: true, keygen: true,
			prepare: func(t *testing.T, r release) {
				other := filepath.Join(t.TempDir(), "other")
				mustRun(t, "", "ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-C", "other@release", "-f", other)
				signRelease(t, r.assets, AssetName(runtime.GOOS, runtime.GOARCH), other, r.tag)
			},
		},
		{
			name: "SHA256SUMS altered after signing", want: "SHA256SUMS signature does not verify",
			gh: true, keygen: true,
			prepare: func(t *testing.T, r release) {
				writeFile(t, filepath.Join(r.assets, "SHA256SUMS"),
					"# styrman-harness "+r.tag+"\n"+strings.Repeat("0", 64)+"  "+AssetName(runtime.GOOS, runtime.GOARCH)+"\n")
			},
		},
		{
			name: "binary does not match its signed checksum", want: "checksum mismatch for sh-agent-",
			gh: true, keygen: true,
			prepare: func(t *testing.T, r release) {
				// The asset changes AFTER the checksums were signed, so the
				// signature still verifies and only the hash catches it.
				writeFile(t, filepath.Join(r.assets, AssetName(runtime.GOOS, runtime.GOARCH)),
					"#!/bin/sh\necho not the binary that was signed\n")
			},
		},
		{
			// The trust model is that the signer is COMPILED IN. A release
			// that ships its own allowed-signers file must not be able to
			// nominate the key that vouches for it.
			name: "release ships its own signer file", want: "SHA256SUMS signature does not verify",
			gh: true, keygen: true,
			prepare: func(t *testing.T, r release) {
				other := filepath.Join(t.TempDir(), "other")
				mustRun(t, "", "ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-C", "other@release", "-f", other)
				signRelease(t, r.assets, AssetName(runtime.GOOS, runtime.GOARCH), other, r.tag)
				pub, err := os.ReadFile(other + ".pub")
				if err != nil {
					t.Fatal(err)
				}
				fields := strings.Fields(string(pub))
				// The SAME identity as the compiled-in line, so only the KEY
				// differs: nothing but where the signer line came from can
				// tell these two apart.
				line := "test@release " + fields[0] + " " + fields[1] + "\n"
				for _, name := range []string{"allowed_signers", "allowed_signers_from_release"} {
					writeFile(t, filepath.Join(r.assets, name), line)
				}
			},
		},
		{
			// The signature is honest and the hash matches; only the version
			// the signed file names says this is another release's SHA256SUMS
			// served under this tag.
			name: "signed checksums of another release", want: "the signed SHA256SUMS says v1.2.3 but this is release v9.9.9",
			gh: true, keygen: true,
			prepare: func(t *testing.T, r release) {
				asset := AssetName(runtime.GOOS, runtime.GOARCH)
				data, err := os.ReadFile(filepath.Join(r.assets, asset))
				if err != nil {
					t.Fatal(err)
				}
				sum := sha256.Sum256(data)
				writeFile(t, filepath.Join(r.assets, "SHA256SUMS"),
					"# styrman-harness v1.2.3\n"+hex.EncodeToString(sum[:])+"  "+asset+"\n")
				if err := os.Remove(filepath.Join(r.assets, "SHA256SUMS.sig")); err != nil {
					t.Fatal(err)
				}
				mustRun(t, r.assets, "ssh-keygen", "-Y", "sign", "-q", "-f", r.key, "-n", ReleaseNamespace, "SHA256SUMS")
			},
		},
		{
			// A signed SHA256SUMS with no version in it cannot be told from
			// any other release's, which is what the marker exists to stop.
			name: "signed checksums naming no version", want: "names no release version",
			gh: true, keygen: true,
			prepare: func(t *testing.T, r release) {
				asset := AssetName(runtime.GOOS, runtime.GOARCH)
				data, err := os.ReadFile(filepath.Join(r.assets, asset))
				if err != nil {
					t.Fatal(err)
				}
				sum := sha256.Sum256(data)
				writeFile(t, filepath.Join(r.assets, "SHA256SUMS"), hex.EncodeToString(sum[:])+"  "+asset+"\n")
				if err := os.Remove(filepath.Join(r.assets, "SHA256SUMS.sig")); err != nil {
					t.Fatal(err)
				}
				mustRun(t, r.assets, "ssh-keygen", "-Y", "sign", "-q", "-f", r.key, "-n", ReleaseNamespace, "SHA256SUMS")
			},
		},
		{
			name: "no signature in the release", want: "has no SHA256SUMS.sig",
			gh: true, keygen: true,
			prepare: func(t *testing.T, r release) {
				if err := os.Remove(filepath.Join(r.assets, "SHA256SUMS.sig")); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "no checksums in the release", want: "has no SHA256SUMS",
			gh: true, keygen: true,
			prepare: func(t *testing.T, r release) {
				if err := os.Remove(filepath.Join(r.assets, "SHA256SUMS")); err != nil {
					t.Fatal(err)
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newRelease(t, "v9.9.9")
			r.env(t, tc.gh, tc.keygen)
			if tc.prepare != nil {
				tc.prepare(t, r)
			}
			home := t.TempDir()
			dest := installedBinary(t, home)
			before, err := os.ReadFile(dest)
			if err != nil {
				t.Fatal(err)
			}
			marker := filepath.Join(home, "installed.txt")
			t.Setenv("SH_AGENT_FAKE_MARKER", marker)

			var out bytes.Buffer
			err = UpdateFromRelease(home, "v1.0.0", false, nil, &out)
			if err == nil {
				t.Fatalf("update succeeded, want a refusal naming %q\n%s", tc.want, out.String())
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("err = %v, want it to contain %q", err, tc.want)
			}
			if after, _ := os.ReadFile(dest); string(after) != string(before) {
				t.Error("the installed binary was replaced by a release that did not verify")
			}
			if _, err := os.Stat(marker); err == nil {
				t.Error("an install ran although the release did not verify")
			}
			left, _ := filepath.Glob(filepath.Join(home, ".config", "sh-agent", "release-*"))
			if len(left) != 0 {
				t.Errorf("left behind: %v", left)
			}
		})
	}
}

// TestReleaseSignerMatchesTheDeployGuide: the line compiled in here and the
// allowed-signers line `.agents/docs/deploy/remote.md` § Releases gives
// `git verify-tag` are the same operator key. Two copies of a key are two
// keys the day one of them is rotated, so this fails when they drift.
func TestReleaseSignerMatchesTheDeployGuide(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".agents", "docs", "deploy", "remote.md"))
	if err != nil {
		t.Fatal(err)
	}
	var documented []string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "ostling.simon@gmail.com ssh-") {
			documented = append(documented, strings.TrimSpace(line))
		}
	}
	if len(documented) == 0 {
		t.Fatal("remote.md § Releases names no allowed-signers line")
	}
	for _, line := range documented {
		if line != releaseSigner {
			t.Errorf("remote.md carries\n  %s\nbut the binary verifies releases against\n  %s", line, releaseSigner)
		}
	}
}

func TestChecksumFor(t *testing.T) {
	sums := "# styrman-harness v1.2.3\n" +
		strings.Repeat("a", 64) + "  sh-agent-linux-amd64\n" +
		strings.Repeat("B", 64) + " *sh-agent-darwin-arm64\n" +
		"not-a-line\n" +
		"deadbeef  sh-agent-linux-arm64\n" +
		strings.Repeat("c", 64) + "  twice\n" +
		strings.Repeat("d", 64) + "  twice\n" +
		strings.Repeat("e", 64) + "  same\n" +
		strings.Repeat("e", 64) + "  same\n"
	for _, tc := range []struct {
		name, want, wantErr string
	}{
		{name: "sh-agent-linux-amd64", want: strings.Repeat("a", 64)},
		// A binary-mode '*' and an upper-case hash are both sha256sum's own
		// spellings.
		{name: "sh-agent-darwin-arm64", want: strings.Repeat("b", 64)},
		{name: "sh-agent-linux-arm64", wantErr: "does not carry a SHA-256"},
		{name: "sh-agent-windows-amd64", wantErr: "no line names"},
		// A signed file that says two things about one name is a reason to
		// stop, not to take the first.
		{name: "twice", wantErr: "two different hashes"},
		{name: "same", want: strings.Repeat("e", 64)},
	} {
		got, err := checksumFor(sums, tc.name)
		switch {
		case tc.wantErr != "":
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("checksumFor(%q) err = %v, want %q", tc.name, err, tc.wantErr)
			}
		case err != nil:
			t.Errorf("checksumFor(%q): %v", tc.name, err)
		case got != tc.want:
			t.Errorf("checksumFor(%q) = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestCompareRelease(t *testing.T) {
	for _, tc := range []struct {
		a, b string
		want int
	}{
		{"v1.0.0", "v1.0.0", 0},
		{"v1.0.1", "v1.0.0", 1},
		{"v1.1.0", "v1.0.9", 1},
		{"v2.0.0", "v1.99.99", 1},
		{"v1.0.0", "v1.0.1", -1},
		{"v1.9.0", "v1.10.0", -1}, // not a string comparison
		{"v0.9.0", "v1.0.0", -1},
	} {
		if got := compareRelease(tc.a, tc.b); got != tc.want {
			t.Errorf("compareRelease(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

// TestSSHKeygenPrefersASystemPath: the verifier's verdict IS the trust chain,
// so a copy planted earlier in PATH must not be the one that runs (security
// review 2026-09-17).
func TestSSHKeygenPrefersASystemPath(t *testing.T) {
	if _, err := os.Stat("/usr/bin/ssh-keygen"); err != nil {
		t.Skip("no system ssh-keygen to prefer")
	}
	bin := t.TempDir()
	planted := filepath.Join(bin, "ssh-keygen")
	writeFile(t, planted, "#!/bin/sh\nexit 0\n")
	if err := os.Chmod(planted, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	got, err := sshKeygen()
	if err != nil {
		t.Fatal(err)
	}
	if got == planted {
		t.Errorf("sshKeygen() = %q, the planted copy", got)
	}
	if got != "/usr/bin/ssh-keygen" {
		t.Errorf("sshKeygen() = %q, want the system path", got)
	}
}

// TestUpdateFromReleaseNeverGoesBackwards: an older release is signed exactly
// as honestly as a newer one, so a release that was deleted or re-pointed
// would otherwise install a genuine binary the operator has moved past, and
// `--check` would keep announcing it (security review 2026-09-17, finding 1).
func TestUpdateFromReleaseNeverGoesBackwards(t *testing.T) {
	r := newRelease(t, "v0.9.0")
	r.env(t, true, true)
	home := t.TempDir()
	dest := installedBinary(t, home)
	before, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(home, "installed.txt")
	t.Setenv("SH_AGENT_FAKE_MARKER", marker)

	for _, checkOnly := range []bool{false, true} {
		var out bytes.Buffer
		if err := UpdateFromRelease(home, "v1.0.0", checkOnly, nil, &out); err != nil {
			t.Fatalf("checkOnly=%v: %v", checkOnly, err)
		}
		if !strings.Contains(out.String(), "the newest release v0.9.0 is older than this binary v1.0.0; nothing installed") {
			t.Errorf("checkOnly=%v output: %q", checkOnly, out.String())
		}
		if after, _ := os.ReadFile(dest); string(after) != string(before) {
			t.Errorf("checkOnly=%v replaced the binary with an older release", checkOnly)
		}
		if _, err := os.Stat(marker); err == nil {
			t.Errorf("checkOnly=%v ran an install", checkOnly)
		}
	}

	// A tag that is not a release never reaches argv, where a leading dash
	// would be read by gh as a flag.
	for _, tag := range []string{"-rf", "v1", "latest", "v1.0.0-rc.1", ""} {
		t.Setenv("GH_TAG", tag)
		var out bytes.Buffer
		err := UpdateFromRelease(home, "v1.0.0", false, nil, &out)
		if err == nil || !strings.Contains(err.Error(), "not a vX.Y.Z release") {
			if tag == "" && err != nil && strings.Contains(err.Error(), "no tag") {
				continue // gh answering nothing at all is its own error
			}
			t.Errorf("tag %q: err = %v, want a refusal", tag, err)
		}
	}

	// A binary with no release version of its own does not install over
	// itself from a release either: it says so and leaves the file alone.
	t.Setenv("GH_TAG", "v9.9.9")
	var out bytes.Buffer
	if err := UpdateFromRelease(home, "b64ffe886729+dirty", false, nil, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "which is not a release; leaving it alone") {
		t.Errorf("development build: %q", out.String())
	}
	if after, _ := os.ReadFile(dest); string(after) != string(before) {
		t.Error("a development build was overwritten from a release")
	}
}

// TestVerifierMessageCarriesWhatTheToolSaid: an OpenSSH too old for
// `-Y verify` exits 255, exactly as a bad signature does, so without the
// tool's own words a broken machine reads as a tampered release. The verdict
// is unchanged — nothing installs either way — but the cause is visible
// (code review 2026-09-17, round 2).
func TestVerifierMessageCarriesWhatTheToolSaid(t *testing.T) {
	bin := t.TempDir()
	fake := filepath.Join(bin, "ssh-keygen")
	writeFile(t, fake, "#!/bin/sh\necho \"unknown option -- Y\" >&2\nexit 255\n")
	if err := os.Chmod(fake, 0o755); err != nil {
		t.Fatal(err)
	}
	old := sshKeygenSystemPaths
	sshKeygenSystemPaths = []string{fake}
	t.Cleanup(func() { sshKeygenSystemPaths = old })

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "SHA256SUMS"), "# styrman-harness v9.9.9\n")
	writeFile(t, filepath.Join(dir, "SHA256SUMS.sig"), "not really a signature\n")
	keygen, err := sshKeygen()
	if err != nil {
		t.Fatal(err)
	}
	err = verifyRelease(keygen, dir, "v9.9.9", AssetName(runtime.GOOS, runtime.GOARCH), []byte("x"))
	if err == nil {
		t.Fatal("a verifier that cannot verify must be an error")
	}
	if !strings.Contains(err.Error(), "unknown option -- Y") {
		t.Errorf("err = %v, want it to carry ssh-keygen's own message", err)
	}
	if !strings.Contains(err.Error(), "does not verify") {
		t.Errorf("err = %v, want the verdict kept", err)
	}
}
