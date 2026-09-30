package dream

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeManager puts a recording systemctl or launchctl first on PATH; code is
// its exit status. The real user manager is never reached (R6).
func fakeManager(t *testing.T, name string, code int) (calls string) {
	t.Helper()
	bin := t.TempDir()
	calls = filepath.Join(t.TempDir(), "calls")
	script := "#!/bin/sh\necho \"$*\" >> \"" + calls + "\"\nexit " + string(rune('0'+code)) + "\n"
	if err := os.WriteFile(filepath.Join(bin, name), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+"/usr/bin:/bin")
	return calls
}

func timerHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	write(t, filepath.Join(home, ".local", "bin", "sh-agent"), "#!/bin/sh\n", fresh)
	return home
}

func callsOf(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return ""
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestRenderSystemdUnits(t *testing.T) {
	service, timer := renderSystemd("/h/.local/bin/sh-agent", "/h/.local/bin:/usr/bin")
	for _, want := range []string{"Type=oneshot", `Environment="PATH=/h/.local/bin:/usr/bin"`, `ExecStart="/h/.local/bin/sh-agent" dream --all`} {
		if !strings.Contains(service, want) {
			t.Errorf("service lacks %q:\n%s", want, service)
		}
	}
	for _, want := range []string{"OnCalendar=*-*-* 03:00", "Persistent=true", "WantedBy=timers.target"} {
		if !strings.Contains(timer, want) {
			t.Errorf("timer lacks %q:\n%s", want, timer)
		}
	}
}

func TestRenderLaunchdPlist(t *testing.T) {
	plist := renderLaunchd("/Users/o & p/.local/bin/sh-agent", "/usr/bin")
	for _, want := range []string{"<string>com.styrman.sh-agent-dream</string>", "<string>/Users/o &amp; p/.local/bin/sh-agent</string>", "<string>dream</string>", "<string>--all</string>", "<key>Hour</key><integer>3</integer>", "<key>PATH</key><string>/usr/bin</string>"} {
		if !strings.Contains(plist, want) {
			t.Errorf("plist lacks %q:\n%s", want, plist)
		}
	}
}

func TestInstallTimerLinux(t *testing.T) {
	home := timerHome(t)
	calls := fakeManager(t, "systemctl", 0)
	if err := InstallTimer(home, "/p:/q", "linux"); err != nil {
		t.Fatal(err)
	}
	unitDir := filepath.Join(home, ".config", "systemd", "user")
	service := read(t, filepath.Join(unitDir, "sh-agent-dream.service"))
	if !strings.Contains(service, `ExecStart="`+filepath.Join(home, ".local", "bin", "sh-agent")+`" dream --all`) || !strings.Contains(service, `PATH=/p:/q`) {
		t.Errorf("service:\n%s", service)
	}
	if got := callsOf(t, calls); got != "--user daemon-reload\n--user enable --now sh-agent-dream.timer\n" {
		t.Errorf("systemctl calls:\n%s", got)
	}
	if !TimerInstalled(home) {
		t.Error("TimerInstalled = false after install")
	}
	// Identical files: a rerun is fine.
	if err := InstallTimer(home, "/p:/q", "linux"); err != nil {
		t.Errorf("identical rerun: %v", err)
	}
	// Another PATH renders other text: the file on disk is kept and named.
	os.Remove(calls)
	err := InstallTimer(home, "/other", "linux")
	var ee *ExitError
	if !errors.As(err, &ee) || ee.Code != 2 || !strings.Contains(ee.Msg, "sh-agent-dream.service") || !strings.Contains(ee.Msg, "--remove-timer") {
		t.Errorf("a differing unit: %v", err)
	}
	if callsOf(t, calls) != "" {
		t.Error("systemctl ran although a unit was kept")
	}
	if !strings.Contains(read(t, filepath.Join(unitDir, "sh-agent-dream.service")), "PATH=/p:/q") {
		t.Error("the differing unit was overwritten")
	}
}

func TestInstallTimerRefusals(t *testing.T) {
	var ee *ExitError
	// No binary at BinPath.
	fakeManager(t, "systemctl", 0)
	if err := InstallTimer(t.TempDir(), "/p", "linux"); !errors.As(err, &ee) || ee.Code != 2 || !strings.Contains(ee.Msg, "sh-agent") {
		t.Errorf("no binary: %v", err)
	}
	// No systemctl on PATH.
	home := timerHome(t)
	t.Setenv("PATH", t.TempDir())
	if err := InstallTimer(home, "/p", "linux"); !errors.As(err, &ee) || ee.Code != 2 || !strings.Contains(ee.Msg, "systemctl") {
		t.Errorf("no systemctl: %v", err)
	}
	if TimerInstalled(home) {
		t.Error("files written although systemctl is missing")
	}
	// systemctl fails: the files stay and are named.
	home = timerHome(t)
	fakeManager(t, "systemctl", 1)
	if err := InstallTimer(home, "/p", "linux"); !errors.As(err, &ee) || ee.Code != 2 || !strings.Contains(ee.Msg, "left in place") {
		t.Errorf("failing systemctl: %v", err)
	}
	// An OS with no user timer.
	if err := InstallTimer(timerHome(t), "/p", "windows"); err == nil {
		t.Error("installed on an unsupported OS")
	}
}

func TestRemoveTimerLinux(t *testing.T) {
	home := timerHome(t)
	calls := fakeManager(t, "systemctl", 0)
	if err := RemoveTimer(home, "linux"); err != nil {
		t.Fatalf("absent is fine: %v", err)
	}
	if callsOf(t, calls) != "" {
		t.Error("systemctl ran with no unit installed")
	}
	if err := InstallTimer(home, "/p", "linux"); err != nil {
		t.Fatal(err)
	}
	os.Remove(calls)
	if err := RemoveTimer(home, "linux"); err != nil {
		t.Fatal(err)
	}
	if got := callsOf(t, calls); got != "--user disable --now sh-agent-dream.timer\n--user daemon-reload\n" {
		t.Errorf("systemctl calls:\n%s", got)
	}
	if TimerInstalled(home) {
		t.Error("unit files left after remove")
	}
}

func TestTimerMacOS(t *testing.T) {
	home := timerHome(t)
	calls := fakeManager(t, "launchctl", 0)
	if err := InstallTimer(home, "/p", "darwin"); err != nil {
		t.Fatal(err)
	}
	plist := macPlistPath(home)
	if _, err := os.Stat(plist); err != nil {
		t.Fatal(err)
	}
	if got := callsOf(t, calls); !strings.HasPrefix(got, "bootstrap gui/") || !strings.Contains(got, plist) {
		t.Errorf("launchctl calls:\n%s", got)
	}
	os.Remove(calls)
	if err := RemoveTimer(home, "darwin"); err != nil {
		t.Fatal(err)
	}
	if got := callsOf(t, calls); !strings.HasPrefix(got, "bootout gui/") {
		t.Errorf("launchctl calls:\n%s", got)
	}
	if TimerInstalled(home) {
		t.Error("plist left after remove")
	}
}
