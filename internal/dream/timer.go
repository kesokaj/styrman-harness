package dream

// The nightly timer: R3's exception runs `sh-agent dream --all` from an OS
// user timer that `--install-timer` installs and `--remove-timer` or
// `uninstall` removes (R4). The timer is not in the install manifest: the
// next `install` would delete a listed unit as stale (plan D12).

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kesokaj/styrman-harness/internal/harness"
)

const (
	unitName  = "sh-agent-dream"
	plistName = "com.styrman.sh-agent-dream"
)

func linuxServicePath(home string) string {
	return filepath.Join(home, ".config", "systemd", "user", unitName+".service")
}

// renderSystemd returns the service and timer unit texts. PATH is the one of
// the shell that ran --install-timer: the user manager's own PATH lacks
// ~/.local/bin, where claude and agy live (cold read r1).
func renderSystemd(bin, path string) (string, string) {
	// systemd expands %-specifiers inside quotes too; %% is a literal %.
	bin, path = strings.ReplaceAll(bin, "%", "%%"), strings.ReplaceAll(path, "%", "%%")
	service := fmt.Sprintf(`[Unit]
Description=sh-agent dream: the nightly background dreamer (R3)

[Service]
Type=oneshot
Environment="PATH=%s"
ExecStart="%s" dream --all
`, path, bin)
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
func renderLaunchd(bin, path string) string {
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
  <array><string>%s</string><string>dream</string><string>--all</string></array>
  <key>EnvironmentVariables</key>
  <dict><key>PATH</key><string>%s</string></dict>
  <key>StartCalendarInterval</key>
  <dict><key>Hour</key><integer>3</integer><key>Minute</key><integer>0</integer></dict>
</dict>
</plist>
`, plistName, x(bin), x(path))
}

// InstallTimer installs and starts the nightly dream timer for goos, with
// ExecStart at home's BinPath and PATH as given.
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
	var files map[string]string
	var manager string
	var start [][]string
	switch goos {
	case "linux":
		service, timer := renderSystemd(bin, path)
		files = map[string]string{linuxServicePath(home): service, linuxTimerPath(home): timer}
		manager = "systemctl"
		start = [][]string{{"--user", "daemon-reload"}, {"--user", "enable", "--now", unitName + ".timer"}}
	case "darwin":
		files = map[string]string{macPlistPath(home): renderLaunchd(bin, path)}
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

// RemoveTimer stops and removes the dream timer; absent is fine.
func RemoveTimer(home, goos string) error {
	var files []string
	var manager string
	var stop [][]string
	switch goos {
	case "linux":
		files = []string{linuxTimerPath(home), linuxServicePath(home)}
		manager = "systemctl"
		stop = [][]string{{"--user", "disable", "--now", unitName + ".timer"}, {"--user", "daemon-reload"}}
	case "darwin":
		files = []string{macPlistPath(home)}
		manager = "launchctl"
		stop = [][]string{{"bootout", "gui/" + strconv.Itoa(os.Getuid()), macPlistPath(home)}}
	default:
		return nil
	}
	present := false
	for _, p := range files {
		if _, err := os.Lstat(p); err == nil {
			present = true
		}
	}
	if !present {
		return nil
	}
	if _, err := exec.LookPath(manager); err == nil {
		// The first stop call may fail when the unit was never loaded; the
		// files are still removed, and the reload runs after.
		exec.Command(manager, stop[0]...).Run()
	}
	for _, p := range files {
		if err := os.Remove(p); err != nil && !errors.Is(err, os.ErrNotExist) {
			return &ExitError{2, "sh-agent dream --remove-timer: " + err.Error()}
		}
	}
	if len(stop) > 1 {
		if _, err := exec.LookPath(manager); err == nil {
			exec.Command(manager, stop[1]...).Run()
		}
	}
	return nil
}
