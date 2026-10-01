package dream

// The nightly timer: R3's exception runs `sh-agent dream --all` from an OS
// user timer that `--install-timer` installs and `--remove-timer` or
// `uninstall` removes (R4). The timer is not in the install manifest: the
// next `install` would delete a listed unit as stale (plan D12).

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kesokaj/styrman-harness/internal/harness"
)

const (
	unitName  = "sh-agent-dream"
	plistName = "com.styrman.sh-agent-dream"
)

// userHomeDir is the real home; a test seam, since the tests' homes are
// temporary directories that must count as the real one.
var userHomeDir = accountHome

// accountHome is the home of the account the process runs as, from the
// account database and never from HOME: the user service manager belongs to
// the account, and a sandbox run overrides HOME (the testing guide's kits
// do), so HOME as the yardstick let a sandbox --remove-timer or uninstall
// reach the operator's real timer (BUG-149). An account the database does
// not list has no known home, and then no home gets the bare names:
// user.Current would answer with HOME there without saying so (os/user
// without cgo reads only /etc/passwd on Linux, so LDAP, SSSD and
// systemd-homed accounts and a container's bare uid are not listed; on
// macOS it asks the system's own resolver, with or without cgo).
func accountHome() (string, error) {
	u, err := lookupAccount()
	if err != nil {
		return "", err
	}
	if u.HomeDir == "" {
		return "", errors.New("the account has no home directory")
	}
	return u.HomeDir, nil
}

// lookupAccount is the account database lookup by uid; a test seam.
// LookupId, not Current: Current falls back to HOME when the lookup fails.
var lookupAccount = func() (*user.User, error) { return user.LookupId(strconv.Itoa(os.Getuid())) }

// UnitFor names the systemd unit and the launchd label for a home. The real
// home gets the bare names; any other home gets them suffixed with a short
// hash of its path, so a scratch `--install-timer`, `--remove-timer` or
// `uninstall` never touches the operator's own timer through the shared user
// service manager (R6; security review of the v2.0.0 gate).
func UnitFor(home string) (unit, plist string) {
	if real, err := userHomeDir(); err == nil {
		if filepath.Clean(real) == filepath.Clean(home) {
			return unitName, plistName
		}
		// A symlinked spelling of the real home is the real home (review 2026-10-01).
		if ri, err := os.Stat(real); err == nil {
			if hi, err := os.Stat(home); err == nil && os.SameFile(ri, hi) {
				return unitName, plistName
			}
		}
	}
	sum := sha256.Sum256([]byte(filepath.Clean(home)))
	suffix := "-" + hex.EncodeToString(sum[:4])
	return unitName + suffix, plistName + suffix
}

func linuxServicePath(home string) string {
	unit, _ := UnitFor(home)
	return filepath.Join(home, ".config", "systemd", "user", unit+".service")
}

// renderSystemd returns the service and timer unit texts. PATH is the one of
// the shell that ran --install-timer: the user manager's own PATH lacks
// ~/.local/bin, where claude and agy live (cold read r1).
func renderSystemd(bin, home, path string) (string, string) {
	// systemd expands %-specifiers inside quotes too; %% is a literal %.
	bin, home, path = strings.ReplaceAll(bin, "%", "%%"), strings.ReplaceAll(home, "%", "%%"), strings.ReplaceAll(path, "%", "%%")
	// ExecStart also expands $VAR and ${VAR}, inside quotes too; $$ is a
	// literal dollar sign. Environment= expands nothing, so path keeps its own.
	bin, home = strings.ReplaceAll(bin, "$", "$$"), strings.ReplaceAll(home, "$", "$$")
	service := fmt.Sprintf(`[Unit]
Description=sh-agent dream: the nightly background dreamer (R3)

[Service]
Type=oneshot
Environment="PATH=%s"
ExecStart="%s" dream --all --home "%s"
`, path, bin, home)
	timer := `[Unit]
Description=Run sh-agent dream nightly

[Timer]
OnCalendar=*-*-* 03:00
Persistent=true

[Install]
WantedBy=timers.target
`
	return service, timer
}

// renderLaunchd returns the launchd agent's plist.
func renderLaunchd(label, bin, home, path string) string {
	x := func(s string) string {
		var b strings.Builder
		for _, r := range s {
			switch r {
			case '&':
				b.WriteString("&amp;")
			case '<':
				b.WriteString("&lt;")
			case '>':
				b.WriteString("&gt;")
			default:
				b.WriteRune(r)
			}
		}
		return b.String()
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key><string>%s</string>
  <key>ProgramArguments</key>
  <array><string>%s</string><string>dream</string><string>--all</string><string>--home</string><string>%s</string></array>
  <key>EnvironmentVariables</key>
  <dict><key>PATH</key><string>%s</string></dict>
  <key>StartCalendarInterval</key>
  <dict><key>Hour</key><integer>3</integer><key>Minute</key><integer>0</integer></dict>
</dict>
</plist>
`, label, x(bin), x(home), x(path))
}

// InstallTimer installs and starts the nightly dream timer for goos, with
// ExecStart at home's BinPath, `--home home` so the nightly run is scoped to
// that home (a scratch timer never dreams over the real home, R6), and PATH
// as given.
func InstallTimer(home, path, goos string) error {
	bin := filepath.Join(home, harness.BinPath)
	if info, err := os.Stat(bin); err != nil || !info.Mode().IsRegular() {
		return &ExitError{2, fmt.Sprintf("sh-agent dream --install-timer: no binary at %s: run sh-agent install first", bin)}
	}
	path = absolutePath(path)
	// A quote, a backslash or a line break would end the unit's quoted value
	// and let the rest of the path become unit text (security review
	// 2026-09-30); such a home or PATH is refused, never escaped by guess.
	if strings.ContainsAny(bin+path, "\"\\\n\r") {
		return &ExitError{2, "sh-agent dream --install-timer: the binary path or PATH holds a quote, a backslash or a line break; the timer's unit cannot carry it safely"}
	}
	unit, label := UnitFor(home)
	var files map[string]string
	var manager string
	var start [][]string
	switch goos {
	case "linux":
		service, timer := renderSystemd(bin, home, path)
		files = map[string]string{linuxServicePath(home): service, linuxTimerPath(home): timer}
		manager = "systemctl"
		start = [][]string{{"--user", "daemon-reload"}, {"--user", "enable", "--now", unit + ".timer"}}
	case "darwin":
		files = map[string]string{macPlistPath(home): renderLaunchd(label, bin, home, path)}
		manager = "launchctl"
		start = [][]string{{"bootstrap", "gui/" + strconv.Itoa(os.Getuid()), macPlistPath(home)}}
	default:
		return &ExitError{2, "sh-agent dream --install-timer: no user timer on " + goos + " (systemd on Linux, launchd on macOS)"}
	}
	if _, err := exec.LookPath(manager); err != nil {
		return &ExitError{2, fmt.Sprintf("sh-agent dream --install-timer: no %s on PATH: the timer needs the user service manager", manager)}
	}
	// A unit on disk with other text is the operator's or an older install's:
	// kept and named, never overwritten.
	var paths []string
	for p, text := range files {
		paths = append(paths, p)
		if data, err := os.ReadFile(p); err == nil && string(data) != text {
			return &ExitError{2, fmt.Sprintf("sh-agent dream --install-timer: %s exists with other text (another PATH, or edited by hand); kept. Run sh-agent dream --remove-timer first", p)}
		}
	}
	for p, text := range files {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return &ExitError{2, "sh-agent dream --install-timer: " + err.Error()}
		}
		if err := harness.WriteAtomic(p, []byte(text), 0o644); err != nil {
			return &ExitError{2, "sh-agent dream --install-timer: " + err.Error()}
		}
	}
	for _, args := range start {
		if out, err := exec.Command(manager, args...).CombinedOutput(); err != nil {
			return &ExitError{2, fmt.Sprintf("sh-agent dream --install-timer: %s %s: %v %s; the unit files are left in place: %s",
				manager, strings.Join(args, " "), err, strings.TrimSpace(string(out)), strings.Join(paths, ", "))}
		}
	}
	return nil
}

// absolutePath keeps PATH's absolute entries: a relative one (".", "bin", an
// empty one) would resolve against the unit's working directory at 03:00,
// not where the installing shell was (security review round 2).
func absolutePath(path string) string {
	var keep []string
	for _, e := range filepath.SplitList(path) {
		if filepath.IsAbs(e) {
			keep = append(keep, e)
		}
	}
	return strings.Join(keep, string(filepath.ListSeparator))
}

// RemoveTimer stops and removes the dream timer of this home, for goos, and
// reports whether there was one: absent is fine, but the caller says which
// it was (a removal that found nothing once printed "removed").
func RemoveTimer(home, goos string) (bool, error) {
	unit, _ := UnitFor(home)
	var files []string
	var manager string
	var stop [][]string
	switch goos {
	case "linux":
		files = []string{linuxTimerPath(home), linuxServicePath(home)}
		manager = "systemctl"
		stop = [][]string{{"--user", "disable", "--now", unit + ".timer"}, {"--user", "daemon-reload"}}
	case "darwin":
		files = []string{macPlistPath(home)}
		manager = "launchctl"
		stop = [][]string{{"bootout", "gui/" + strconv.Itoa(os.Getuid()), macPlistPath(home)}}
	default:
		return false, nil
	}
	present := false
	for _, p := range files {
		if _, err := os.Lstat(p); err == nil {
			present = true
		}
	}
	if !present {
		return false, nil
	}
	if _, err := exec.LookPath(manager); err == nil {
		// The first stop call may fail when the unit was never loaded; the
		// files are still removed, and the reload runs after.
		exec.Command(manager, stop[0]...).Run()
	}
	for _, p := range files {
		if err := os.Remove(p); err != nil && !errors.Is(err, os.ErrNotExist) {
			return false, &ExitError{2, "sh-agent dream --remove-timer: " + err.Error()}
		}
	}
	if len(stop) > 1 {
		if _, err := exec.LookPath(manager); err == nil {
			exec.Command(manager, stop[1]...).Run()
		}
	}
	return true, nil
}

// LegacyTimerFiles lists the bare-named unit files under a home whose own
// unit name is suffixed: what v2.0.0 or earlier wrote there before names were
// per home (BUG-144), or what a home of an account the database does not
// list holds since BUG-149. They are not this home's unit any more, so
// install, remove and the installed check never see them; this is how the
// operator still learns of them. The real home's bare-named unit is its own.
func LegacyTimerFiles(home string) []string {
	if unit, _ := UnitFor(home); unit == unitName {
		return nil
	}
	var found []string
	for _, p := range []string{
		filepath.Join(home, ".config", "systemd", "user", unitName+".service"),
		filepath.Join(home, ".config", "systemd", "user", unitName+".timer"),
		filepath.Join(home, "Library", "LaunchAgents", plistName+".plist"),
	} {
		if _, err := os.Lstat(p); err == nil {
			found = append(found, p)
		}
	}
	return found
}

// LegacyTimerNotice is what remove, uninstall and status tell the operator
// about LegacyTimerFiles: the paths and the steps. Nothing here stops or
// deletes them, because the service manager is shared by every home of the
// account and a stop by name ends whichever unit it loaded, which may be
// the real home's.
func LegacyTimerNotice(home, goos string) string {
	files := LegacyTimerFiles(home)
	if len(files) == 0 {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "legacy dream timer files under %s, written by v2.0.0 or earlier under the name every home shared:\n", home)
	for _, f := range files {
		fmt.Fprintf(&b, "  %s\n", f)
	}
	b.WriteString("They are not removed or stopped here: a stop by name ends whichever unit the service manager loaded, which may be another home's.\n")
	if goos == "darwin" {
		uid := strconv.Itoa(os.Getuid())
		fmt.Fprintf(&b, "Check which file is loaded: launchctl print gui/%s/%s (its path line)\n", uid, plistName)
		fmt.Fprintf(&b, "If it is one of the files above: launchctl bootout gui/%s <that plist>\n", uid)
		b.WriteString("Then delete the files above.")
	} else {
		fmt.Fprintf(&b, "Check which file is loaded: systemctl --user show -p FragmentPath %s.timer\n", unitName)
		fmt.Fprintf(&b, "If it is one of the files above: systemctl --user disable --now %s.timer\n", unitName)
		b.WriteString("Then delete the files above and run: systemctl --user daemon-reload")
	}
	return b.String()
}
