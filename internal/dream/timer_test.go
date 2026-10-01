package dream

import (
	"errors"
	"os"
	"os/user"
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
	// The test's home counts as the real one, so the bare unit names apply.
	prev := userHomeDir
	userHomeDir = func() (string, error) { return home, nil }
	t.Cleanup(func() { userHomeDir = prev })
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
	service, timer := renderSystemd("/h/.local/bin/sh-agent", "/h", "/h/.local/bin:/usr/bin")
	for _, want := range []string{"Type=oneshot", `Environment="PATH=/h/.local/bin:/usr/bin"`, `ExecStart="/h/.local/bin/sh-agent" dream --all --home "/h"`} {
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
	plist := renderLaunchd("com.styrman.sh-agent-dream", "/Users/o & p/.local/bin/sh-agent", "/Users/o & p", "/usr/bin")
	for _, want := range []string{"<string>com.styrman.sh-agent-dream</string>", "<string>/Users/o &amp; p/.local/bin/sh-agent</string>", "<string>dream</string>", "<string>--all</string>", "<string>--home</string><string>/Users/o &amp; p</string>", "<key>Hour</key><integer>3</integer>", "<key>PATH</key><string>/usr/bin</string>"} {
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

func TestInstallTimerRefusesUnsafePaths(t *testing.T) {
	fakeManager(t, "systemctl", 0)
	var ee *ExitError
	for _, path := range []string{"/p\"\nExecStartPre=/bin/evil", "/p\\x", "/p\r"} {
		home := timerHome(t)
		if err := InstallTimer(home, path, "linux"); !errors.As(err, &ee) || ee.Code != 2 || TimerInstalled(home) {
			t.Errorf("PATH %q: %v", path, err)
		}
	}
	service, _ := renderSystemd("/h/100%/sh-agent", "/h/100%", "/p%n")
	if !strings.Contains(service, "/h/100%%/sh-agent") || !strings.Contains(service, "PATH=/p%%n") || !strings.Contains(service, `--home "/h/100%%"`) {
		t.Errorf("a %% was not escaped:\n%s", service)
	}
}

// TestInstallTimerKeepsOnlyAbsolutePath: a relative PATH entry (".", "bin",
// an empty one) would resolve against the unit's working directory at 03:00,
// not where the operator's shell was; only absolute entries are frozen in.
func TestInstallTimerKeepsOnlyAbsolutePath(t *testing.T) {
	home := timerHome(t)
	fakeManager(t, "systemctl", 0)
	if err := InstallTimer(home, "/a/bin::.:node_modules/.bin:/b", "linux"); err != nil {
		t.Fatal(err)
	}
	service := read(t, filepath.Join(home, ".config", "systemd", "user", "sh-agent-dream.service"))
	if !strings.Contains(service, `Environment="PATH=/a/bin:/b"`) {
		t.Errorf("service:\n%s", service)
	}
}

// A home other than the real one gets its own unit name and label, so a scratch
// install or removal never reaches the operator's timer through the shared user
// manager (v2.0.0 gate, security code slice; R6).
func TestUnitNamesPerHome(t *testing.T) {
	real := timerHome(t)
	if unit, plist := UnitFor(real); unit != "sh-agent-dream" || plist != "com.styrman.sh-agent-dream" {
		t.Fatalf("real home: %s %s", unit, plist)
	}
	other := t.TempDir()
	unit, plist := UnitFor(other)
	if !strings.HasPrefix(unit, "sh-agent-dream-") || len(unit) != len("sh-agent-dream-")+8 || plist != "com.styrman."+unit {
		t.Fatalf("other home: %s %s", unit, plist)
	}
	if again, _ := UnitFor(other); again != unit {
		t.Fatalf("not deterministic: %s then %s", unit, again)
	}
	if u2, _ := UnitFor(t.TempDir()); u2 == unit {
		t.Fatal("two homes share a unit name")
	}
	write(t, filepath.Join(other, ".local", "bin", "sh-agent"), "#!/bin/sh\n", fresh)
	calls := fakeManager(t, "systemctl", 0)
	if err := InstallTimer(other, "/p", "linux"); err != nil {
		t.Fatal(err)
	}
	if got := callsOf(t, calls); got != "--user daemon-reload\n--user enable --now "+unit+".timer\n" {
		t.Errorf("systemctl calls:\n%s", got)
	}
	service := read(t, filepath.Join(other, ".config", "systemd", "user", unit+".service"))
	if !strings.Contains(service, `dream --all --home "`+other+`"`) {
		t.Errorf("a scratch timer must dream over its own home:\n%s", service)
	}
	os.Remove(calls)
	if err := RemoveTimer(other, "linux"); err != nil {
		t.Fatal(err)
	}
	if got := callsOf(t, calls); got != "--user disable --now "+unit+".timer\n--user daemon-reload\n" {
		t.Errorf("remove side keeps the bare name:\n%s", got)
	}
	// macOS: the label and the plist path carry the same suffix, and bootout names that path.
	mcalls := fakeManager(t, "launchctl", 0)
	if err := InstallTimer(other, "/p", "darwin"); err != nil {
		t.Fatal(err)
	}
	pl := read(t, macPlistPath(other))
	if !strings.Contains(pl, "<string>"+plist+"</string>") || !strings.Contains(pl, "<string>--home</string><string>"+other+"</string>") || !strings.HasSuffix(macPlistPath(other), plist+".plist") {
		t.Errorf("darwin plist for a scratch home:\n%s", pl)
	}
	os.Remove(mcalls)
	if err := RemoveTimer(other, "darwin"); err != nil {
		t.Fatal(err)
	}
	if got := callsOf(t, mcalls); !strings.Contains(got, "bootout gui/") || !strings.Contains(got, plist+".plist") {
		t.Errorf("launchctl bootout:\n%s", got)
	}
	// A symlinked spelling of the real home is the real home.
	link := filepath.Join(t.TempDir(), "home-link")
	if err := os.Symlink(real, link); err != nil {
		t.Skip("no symlinks here")
	}
	if u, _ := UnitFor(link); u != "sh-agent-dream" {
		t.Errorf("symlinked real home got %s", u)
	}
}

// TestUnitForIgnoresHomeOverride: a sandbox run sets HOME to its own
// directory and passes that directory as --home. The service manager is the
// user's, not the sandbox's, so that home is not the real one and must not get
// the bare unit names (security review of the v2.1.0 gate: with HOME as the
// yardstick, a sandbox --remove-timer reached the operator's real timer).
func TestUnitForIgnoresHomeOverride(t *testing.T) {
	sandbox := t.TempDir()
	t.Setenv("HOME", sandbox)
	// The yardstick is the platform's own account lookup, not the function
	// under test: where that lookup falls back to HOME (a pure-Go build on
	// macOS), there is nothing to tell apart.
	u, err := user.Current()
	if err != nil || u.HomeDir == "" || filepath.Clean(u.HomeDir) == filepath.Clean(sandbox) {
		t.Skipf("this platform resolves the account's home from HOME (%v): nothing to tell apart", err)
	}
	unit, plist := UnitFor(sandbox)
	if unit == unitName || plist == plistName {
		t.Errorf("a HOME override got the real home's unit names: %s, %s", unit, plist)
	}
}

// TestRenderSystemdEscapesDollar: systemd expands $VAR and ${VAR} in
// ExecStart, inside quotes too; $$ is a literal dollar sign.
func TestRenderSystemdEscapesDollar(t *testing.T) {
	service, _ := renderSystemd("/h/a$b/sh-agent", "/h/${X}", "/p")
	if !strings.Contains(service, `ExecStart="/h/a$$b/sh-agent" dream --all --home "/h/$${X}"`) {
		t.Errorf("a $ in ExecStart was not escaped:\n%s", service)
	}
}
