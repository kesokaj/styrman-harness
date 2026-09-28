package harness

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// ReleaseRepo is the repository a release binary updates from.
const ReleaseRepo = "kesokaj/styrman-harness"

// ReleaseNamespace is the ssh-keygen signature namespace of SHA256SUMS.sig.
// It is part of what is verified: a signature made for anything else does not
// verify here, even with the right key (measured 2026-09-17).
const ReleaseNamespace = "sh-agent-release"

// releaseSigner is the allowed-signers line the release's checksums are
// verified against — the same line `.agents/docs/deploy/remote.md` § Releases
// gives `git verify-tag`, pinned by TestReleaseSignerMatchesTheDeployGuide.
//
// It is COMPILED IN, never read from a file on the machine: a signer file
// anywhere on disk would otherwise decide what counts as a genuine release,
// which is the one thing the operator's key is supposed to decide. It is a
// variable only so a test can substitute a test key's line.
//
// What is NOT pinned this way is ssh-keygen itself, the program that checks
// the signature; sshKeygen prefers a system path for exactly that reason, but
// a machine whose system directories or PATH an attacker owns can defeat any
// check this binary makes (security review 2026-09-17).
var releaseSigner = "ostling.simon@gmail.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBkx0BICH1xG18f0a0F/02Vw0OvBTfSOVNXkPkO62NCR"

// sumsVersionMarker is the line `/preflight` writes at the top of
// SHA256SUMS, naming the release those hashes belong to. It makes the VERSION
// part of the signed bytes: without it a signed SHA256SUMS from any release
// verifies under any tag, so serving an older release's assets under a newer
// tag would pass every check (security review 2026-09-17). `sha256sum -c`
// ignores a line starting with '#', so the file stays readable by it.
var sumsVersionMarker = regexp.MustCompile(`(?m)^#\s*styrman-harness\s+(v[0-9]+\.[0-9]+\.[0-9]+)\s*$`)

// maxAssetBytes and maxSumsBytes bound what a release can make this process
// read into memory. Variables so a test can lower them.
var (
	maxAssetBytes int64 = 256 << 20
	maxSumsBytes  int64 = 1 << 20
)

// ErrNewerRelease is returned by `update --check` on a release binary when a
// newer release exists, so the command can exit 2 as it does for a checkout
// that is behind.
var ErrNewerRelease = errors.New("a newer release is available")

// AssetName is the binary asset's name for one platform.
func AssetName(goos, goarch string) string {
	return fmt.Sprintf("sh-agent-%s-%s", goos, goarch)
}

// LatestRelease asks gh for the newest release's tag. gh carries the
// operator's own GitHub credential, which is the point: a release binary has
// no checkout and no other way in (R8).
func LatestRelease() (string, error) {
	out, err := exec.Command("gh", "release", "view", "--repo", ReleaseRepo, "--json", "tagName", "--jq", ".tagName").Output()
	if err != nil {
		return "", fmt.Errorf("gh release view failed: %s", firstLineOfError(err))
	}
	tag := strings.TrimSpace(string(out))
	if tag == "" {
		return "", errors.New("gh release view returned no tag")
	}
	return tag, nil
}

// UpdateFromRelease is `sh-agent update` on a place that has no checkout: the
// newest release is downloaded, its checksums' signature and the binary's own
// checksum are verified, and only then does the binary get replaced and its
// install run. Nothing is replaced when any check fails, and the download
// directory is removed on every path out.
//
// version is the running binary's; installArgs are passed to the new
// binary's install.
func UpdateFromRelease(home, version string, checkOnly bool, installArgs []string, out io.Writer) error {
	if _, err := exec.LookPath("gh"); err != nil {
		return errors.New("this install has no checkout; install gh and run gh auth login, or download a newer release by hand (README § Install without cloning)")
	}
	tag, err := LatestRelease()
	if err != nil {
		return err
	}
	// The tag comes from the GitHub API, which is not signed: it reaches argv
	// as a positional argument and decides which release is installed, so it
	// must look like a release before either happens (a `-`-leading tag would
	// be read by gh as a flag).
	if !IsReleaseVersion(tag) {
		return fmt.Errorf("the newest release is tagged %q, which is not a vX.Y.Z release; nothing installed", tag)
	}
	switch {
	case tag == version:
		fmt.Fprintf(out, "up to date (%s)\n", version)
		return nil
	// Never sideways or backwards. An older release is signed just as
	// honestly as a newer one, so a deleted or re-pointed release would
	// otherwise install a genuine binary the operator has moved past
	// (security review 2026-09-17).
	case !IsReleaseVersion(version):
		fmt.Fprintf(out, "this binary reports version %q, which is not a release; leaving it alone — download %s by hand (README § Install without cloning)\n", version, tag)
		return nil
	case compareRelease(tag, version) <= 0:
		fmt.Fprintf(out, "the newest release %s is not newer than this binary %s; nothing installed\n", tag, version)
		return nil
	}
	if checkOnly {
		fmt.Fprintf(out, "newer release %s available\n", tag)
		return ErrNewerRelease
	}

	// ssh-keygen verifies the checksums. Without it nothing can be trusted,
	// so nothing is downloaded either.
	keygen, err := sshKeygen()
	if err != nil {
		return err
	}

	// A random directory, and one whose siblings from an interrupted earlier
	// update are cleared first: a predictable name is one another local
	// account can pre-create as a symlink for these writes to follow.
	base, err := ConfigDir(home)
	if err != nil {
		return err
	}
	for _, stale := range globOrNil(filepath.Join(base, "release-*")) {
		if err := os.RemoveAll(stale); err == nil {
			fmt.Fprintf(out, "removed %s, left by an interrupted update\n", stale)
		}
	}
	dir, err := os.MkdirTemp(base, "release-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	asset := AssetName(runtime.GOOS, runtime.GOARCH)
	data, err := downloadRelease(dir, tag, asset, out)
	if err != nil {
		return err
	}
	if err := verifyRelease(keygen, dir, tag, asset, data); err != nil {
		return err
	}

	// Replaced through a temporary file in the destination's own directory
	// and a rename, so an interrupted update leaves the old binary whole
	// rather than a half-written one.
	dest := filepath.Join(home, BinPath)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	if err := writeAtomic(dest, data, 0o755); err != nil {
		return fmt.Errorf("replacing %s: %w", dest, err)
	}
	fmt.Fprintf(out, "verified and installed %s %s\n", asset, tag)

	// The new binary installs the harness it carries; this process is the old
	// one and must not do it (BUG-006's rule, from the other direction).
	cmd := exec.Command(dest, append([]string{"install"}, installArgs...)...)
	cmd.Stdout, cmd.Stderr = out, out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s install: %w", dest, err)
	}
	NoteRunningElsewhere(home, out)
	return nil
}

// downloadRelease fetches the platform's binary and the checksum files into
// dir and returns the binary's bytes.
func downloadRelease(dir, tag, asset string, out io.Writer) ([]byte, error) {
	fmt.Fprintf(out, "downloading %s %s\n", asset, tag)
	cmd := exec.Command("gh", "release", "download", tag,
		"--repo", ReleaseRepo,
		"--pattern", asset,
		"--pattern", "SHA256SUMS",
		"--pattern", "SHA256SUMS.sig",
		"--dir", dir)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	downloadErr := cmd.Run()
	// A release without this platform's asset is the likelier cause of a
	// failure here than anything else, and gh's own message for it is not
	// one the operator can act on, so the missing file is named first.
	path := filepath.Join(dir, asset)
	if _, err := os.Stat(path); err != nil {
		if line := firstStderrLine(stderr.String()); line != "" {
			return nil, fmt.Errorf("release %s has no %s (gh: %s)", tag, asset, line)
		}
		return nil, fmt.Errorf("release %s has no %s", tag, asset)
	}
	data, err := readCapped(path, maxAssetBytes)
	if err != nil {
		return nil, fmt.Errorf("release %s: %s is not usable: %w", tag, asset, err)
	}
	if downloadErr != nil {
		return nil, fmt.Errorf("gh release download %s failed: %s", tag, firstStderrLine(stderr.String()))
	}
	return data, nil
}

// verifyRelease checks the signature over SHA256SUMS and then the binary's
// own line in it. Both must pass before anything is replaced: the signature
// says the checksum file is the operator's, the checksum says the binary is
// the one that file names.
func verifyRelease(keygen, dir, tag, asset string, data []byte) error {
	sums, named, err := verifySignedSums(keygen, dir, dir, "release "+tag)
	if err != nil {
		return err
	}
	if named != tag {
		return fmt.Errorf("the signed SHA256SUMS says %s but this is release %s", named, tag)
	}
	want, err := checksumFor(string(sums), asset)
	if err != nil {
		return fmt.Errorf("SHA256SUMS of %s: %w", tag, err)
	}
	sum := sha256.Sum256(data)
	if got := hex.EncodeToString(sum[:]); got != want {
		return fmt.Errorf("checksum mismatch for %s", asset)
	}
	return nil
}

// verifySignedSums checks sumsDir's SHA256SUMS against SHA256SUMS.sig with the
// compiled-in signer (written to signersDir for ssh-keygen) and returns the
// file and the release it names. what names the subject in messages.
func verifySignedSums(keygen, sumsDir, signersDir, what string) (sums []byte, tag string, err error) {
	sumsPath := filepath.Join(sumsDir, "SHA256SUMS")
	if _, err := os.Stat(sumsPath); err != nil {
		return nil, "", fmt.Errorf("%s has no SHA256SUMS", what)
	}
	sums, err = readCapped(sumsPath, maxSumsBytes)
	if err != nil {
		return nil, "", fmt.Errorf("%s: SHA256SUMS is not usable: %w", what, err)
	}
	if _, err := os.Stat(filepath.Join(sumsDir, "SHA256SUMS.sig")); err != nil {
		return nil, "", fmt.Errorf("%s has no SHA256SUMS.sig", what)
	}
	identity, _, ok := strings.Cut(releaseSigner, " ")
	if !ok || identity == "" {
		return nil, "", errors.New("the compiled-in release signer line has no identity")
	}
	signers := filepath.Join(signersDir, "allowed_signers")
	if err := os.WriteFile(signers, []byte(releaseSigner+"\n"), 0o600); err != nil {
		return nil, "", err
	}
	verify := exec.Command(keygen, "-Y", "verify",
		"-f", signers,
		"-I", identity,
		"-n", ReleaseNamespace,
		"-s", filepath.Join(sumsDir, "SHA256SUMS.sig"))
	verify.Stdin = strings.NewReader(string(sums))
	var verifyErr bytes.Buffer
	verify.Stderr = &verifyErr
	if err := verify.Run(); err != nil {
		// A verifier that could not RUN is not a release that failed to
		// verify; saying the latter would read as a tampered release when the
		// machine's OpenSSH is broken. Either way nothing is installed.
		// What ssh-keygen SAID rides along: an OpenSSH too old for `-Y verify`
		// exits 255 like a bad signature does.
		said := firstStderrLine(verifyErr.String())
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			if said != "" {
				return nil, "", fmt.Errorf("could not run %s to verify the release: %w (%s)", keygen, err, said)
			}
			return nil, "", fmt.Errorf("could not run %s to verify the release: %w", keygen, err)
		}
		if said != "" {
			return nil, "", fmt.Errorf("SHA256SUMS signature does not verify (%s: %s)", filepath.Base(keygen), said)
		}
		return nil, "", errors.New("SHA256SUMS signature does not verify")
	}
	// Only now that the signature holds is anything in the file believed —
	// including the version it names, which is why the marker is inside the
	// signed bytes rather than beside them.
	marker := sumsVersionMarker.FindSubmatch(sums)
	if marker == nil {
		return nil, "", fmt.Errorf("the signed SHA256SUMS of %s names no release version; it cannot be told apart from another release's", what)
	}
	return sums, string(marker[1]), nil
}

// checksumFor reads one file's hash out of a sha256sum-format listing, in
// which a binary-mode line marks the name with a leading '*'. The hash must
// look like one, and two lines naming the same file are refused rather than
// resolved by order: the file is signed, so an ambiguity in it is a reason to
// stop, not to guess.
func checksumFor(sums, name string) (string, error) {
	var found string
	for _, line := range strings.Split(sums, "\n") {
		hash, file, ok := strings.Cut(strings.TrimSpace(line), " ")
		if !ok {
			continue
		}
		if strings.TrimPrefix(strings.TrimSpace(file), "*") != name {
			continue
		}
		hash = strings.ToLower(hash)
		if len(hash) != 64 || strings.Trim(hash, "0123456789abcdef") != "" {
			return "", fmt.Errorf("the line naming %s does not carry a SHA-256", name)
		}
		if found != "" && found != hash {
			return "", fmt.Errorf("two different hashes are given for %s", name)
		}
		found = hash
	}
	if found == "" {
		return "", fmt.Errorf("no line names %s", name)
	}
	return found, nil
}

// sshKeygenSystemPaths are looked at before PATH. A variable so a test can
// point it at nothing and exercise the no-verifier refusal.
var sshKeygenSystemPaths = []string{"/usr/bin/ssh-keygen", "/bin/ssh-keygen"}

// sshKeygen locates the signature verifier, preferring a system path: the
// verifier's own verdict is the whole trust chain, so a copy earlier in PATH
// is the cheapest way to defeat it. PATH is still the fallback, because a
// machine may keep OpenSSH elsewhere.
func sshKeygen() (string, error) {
	for _, p := range sshKeygenSystemPaths {
		// Regular AND executable: a non-executable file of that name is not a
		// verifier, and silently preferring it over a working PATH copy would
		// stop every update (code review 2026-09-17).
		if info, err := os.Stat(p); err == nil && info.Mode().IsRegular() && info.Mode().Perm()&0o111 != 0 {
			return p, nil
		}
	}
	p, err := exec.LookPath("ssh-keygen")
	if err != nil {
		return "", errors.New("ssh-keygen not found; cannot verify the release")
	}
	return p, nil
}

// readCapped reads a file a release supplied, refusing one that would not fit
// in memory comfortably.
func readCapped(path string, max int64) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if info, err := f.Stat(); err == nil && info.Size() > max {
		return nil, fmt.Errorf("%d bytes is larger than the %d this reads", info.Size(), max)
	}
	return io.ReadAll(io.LimitReader(f, max))
}

func globOrNil(pattern string) []string {
	matches, _ := filepath.Glob(pattern)
	return matches
}

// compareRelease orders two vX.Y.Z tags. Both are known to match
// IsReleaseVersion, so every part parses.
func compareRelease(a, b string) int {
	pa, pb := releaseParts(a), releaseParts(b)
	for i := range pa {
		if pa[i] != pb[i] {
			if pa[i] < pb[i] {
				return -1
			}
			return 1
		}
	}
	return 0
}

func releaseParts(v string) [3]int {
	var out [3]int
	for i, part := range strings.SplitN(strings.TrimPrefix(v, "v"), ".", 3) {
		n, err := strconv.Atoi(part)
		if err != nil {
			return [3]int{}
		}
		out[i] = n
	}
	return out
}

// firstLineOfError renders a failed command's first line of stderr, which is
// the part that says what went wrong; the rest is usually a usage dump.
func firstLineOfError(err error) string {
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		if line := firstStderrLine(string(exit.Stderr)); line != "" {
			return line
		}
	}
	return err.Error()
}

func firstStderrLine(s string) string {
	return strings.TrimSpace(strings.SplitN(s, "\n", 2)[0])
}

// SwapReleaseSigner replaces the compiled-in release signer and returns the
// function that restores it. It exists for tests of commands that verify a
// release; nothing else calls it.
func SwapReleaseSigner(line string) (restore func()) {
	old := releaseSigner
	releaseSigner = line
	return func() { releaseSigner = old }
}

// ErrZipUnsigned: no SHA256SUMS and no SHA256SUMS.sig beside a zip — a zip of
// a branch, which nothing signs. The caller installs it marked UNVERIFIED.
var ErrZipUnsigned = errors.New("no signed SHA256SUMS beside the zip")

// VerifyZipSums verifies a release zip against the SHA256SUMS and
// SHA256SUMS.sig beside it: the signature with the compiled-in signer, then
// the zip's own line, by its file name, against sum. It returns the release
// the signed file names. Neither file present is ErrZipUnsigned; anything
// else that fails is a refusal (plan D7).
func VerifyZipSums(zipPath, sum string) (string, error) {
	dir := filepath.Dir(zipPath)
	_, sumsErr := os.Stat(filepath.Join(dir, "SHA256SUMS"))
	_, sigErr := os.Stat(filepath.Join(dir, "SHA256SUMS.sig"))
	// Only "does not exist" means unsigned. Any other error (a folder that
	// cannot be read, say) is a refusal: an unreadable signature is not the
	// absence of one.
	for _, e := range []error{sumsErr, sigErr} {
		if e != nil && !errors.Is(e, fs.ErrNotExist) {
			return "", fmt.Errorf("cannot check for signed checksums beside %s: %w", filepath.Base(zipPath), e)
		}
	}
	if sumsErr != nil && sigErr != nil {
		return "", ErrZipUnsigned
	}
	keygen, err := sshKeygen()
	if err != nil {
		return "", fmt.Errorf("%w; the zip has signed checksums beside it and cannot be verified without it", err)
	}
	signersDir, err := os.MkdirTemp("", "sh-agent-verify-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(signersDir)
	what := filepath.Base(zipPath)
	sums, tag, err := verifySignedSums(keygen, dir, signersDir, what)
	if err != nil {
		return "", err
	}
	want, err := checksumFor(string(sums), what)
	if err != nil {
		return "", fmt.Errorf("the signed SHA256SUMS beside it are release %s's and do not list %s: %w — put the zip in a folder with its own release's SHA256SUMS and SHA256SUMS.sig, under the release's own file name", tag, what, err)
	}
	if want != sum {
		return "", fmt.Errorf("checksum mismatch for %s against the signed SHA256SUMS of %s: the zip is not the one that release published — download it again", what, tag)
	}
	return tag, nil
}
