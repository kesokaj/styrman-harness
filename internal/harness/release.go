package harness

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
var releaseSigner = "ostling.simon@gmail.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBkx0BICH1xG18f0a0F/02Vw0OvBTfSOVNXkPkO62NCR"

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
	if tag == version {
		fmt.Fprintf(out, "up to date (%s)\n", version)
		return nil
	}
	if checkOnly {
		fmt.Fprintf(out, "newer release %s available\n", tag)
		return ErrNewerRelease
	}

	// ssh-keygen verifies the checksums. Without it nothing can be trusted,
	// so nothing is downloaded either.
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		return errors.New("ssh-keygen not found; cannot verify the release")
	}

	dir := filepath.Join(home, ".config", "sh-agent", fmt.Sprintf("release-%d", os.Getpid()))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	asset := AssetName(runtime.GOOS, runtime.GOARCH)
	data, err := downloadRelease(dir, tag, asset, out)
	if err != nil {
		return err
	}
	if err := verifyRelease(dir, tag, asset, data); err != nil {
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
	downloadErr := cmd.Run()
	// A release without this platform's asset is the likelier cause of a
	// failure here than anything else, and gh's own message for it is not
	// one the operator can act on, so the missing file is named first.
	data, readErr := os.ReadFile(filepath.Join(dir, asset))
	if readErr != nil {
		return nil, fmt.Errorf("release %s has no %s", tag, asset)
	}
	if downloadErr != nil {
		return nil, fmt.Errorf("gh release download %s failed: %s", tag, firstLineOfError(downloadErr))
	}
	return data, nil
}

// verifyRelease checks the signature over SHA256SUMS and then the binary's
// own line in it. Both must pass before anything is replaced: the signature
// says the checksum file is the operator's, the checksum says the binary is
// the one that file names.
func verifyRelease(dir, tag, asset string, data []byte) error {
	sums, err := os.ReadFile(filepath.Join(dir, "SHA256SUMS"))
	if err != nil {
		return fmt.Errorf("release %s has no SHA256SUMS", tag)
	}
	if _, err := os.Stat(filepath.Join(dir, "SHA256SUMS.sig")); err != nil {
		return fmt.Errorf("release %s has no SHA256SUMS.sig", tag)
	}
	identity, _, ok := strings.Cut(releaseSigner, " ")
	if !ok || identity == "" {
		return errors.New("the compiled-in release signer line has no identity")
	}
	signers := filepath.Join(dir, "allowed_signers")
	if err := os.WriteFile(signers, []byte(releaseSigner+"\n"), 0o600); err != nil {
		return err
	}
	verify := exec.Command("ssh-keygen", "-Y", "verify",
		"-f", signers,
		"-I", identity,
		"-n", ReleaseNamespace,
		"-s", filepath.Join(dir, "SHA256SUMS.sig"))
	verify.Stdin = strings.NewReader(string(sums))
	if err := verify.Run(); err != nil {
		return errors.New("SHA256SUMS signature does not verify")
	}
	want, ok := checksumFor(string(sums), asset)
	if !ok {
		return fmt.Errorf("SHA256SUMS does not name %s", asset)
	}
	sum := sha256.Sum256(data)
	if got := hex.EncodeToString(sum[:]); got != want {
		return fmt.Errorf("checksum mismatch for %s", asset)
	}
	return nil
}

// checksumFor reads one file's hash out of a sha256sum-format listing, in
// which a binary-mode line marks the name with a leading '*'.
func checksumFor(sums, name string) (string, bool) {
	for _, line := range strings.Split(sums, "\n") {
		hash, file, ok := strings.Cut(strings.TrimSpace(line), " ")
		if !ok {
			continue
		}
		if strings.TrimPrefix(strings.TrimSpace(file), "*") == name {
			return strings.ToLower(hash), true
		}
	}
	return "", false
}

// firstLineOfError renders a failed command's first line of stderr, which is
// the part that says what went wrong; the rest is usually a usage dump.
func firstLineOfError(err error) string {
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		if line := strings.TrimSpace(strings.SplitN(string(exit.Stderr), "\n", 2)[0]); line != "" {
			return line
		}
	}
	return err.Error()
}
