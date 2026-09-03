package gate

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kesokaj/styrman-harness/internal/hookio"
)

// Format runs the language's formatter over a just-edited file. The two
// vendors name the edited path differently, so both shapes are checked.
func Format(p hookio.Payload) {
	file := firstString(
		mapString(p.ToolInput, "file_path"),
		mapString(p.ToolInput, "absolute_path"),
		mapString(p.ToolInput, "path"),
		mapString(p.ToolCall.Args, "TargetFile"),
	)
	if file == "" {
		return
	}
	if fi, err := os.Stat(file); err != nil || fi.IsDir() {
		return
	}

	// This process does not source a login shell, so a Go toolchain installed
	// under the user's home has to be put on PATH explicitly.
	home, _ := os.UserHomeDir()
	os.Setenv("PATH", strings.Join([]string{
		filepath.Join(home, ".go", "bin"), filepath.Join(home, "go", "bin"),
		"/usr/local/go/bin", "/opt/homebrew/bin", os.Getenv("PATH"),
	}, string(os.PathListSeparator)))

	switch strings.ToLower(filepath.Ext(file)) {
	case ".go":
		if bin := trustedTool("goimports"); bin != "" {
			run(bin, "-w", file)
		} else if bin := trustedTool("gofmt"); bin != "" {
			run(bin, "-w", file)
		}
	case ".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs", ".css", ".scss",
		".json", ".md", ".yaml", ".yml", ".html":
		if bin := trustedTool("prettier"); bin != "" {
			run(bin, "--write", file)
		}
	}
}

// trustedTool resolves a formatter, but only outside the repository being
// edited.
//
// This hook fires automatically on every Edit and Write, with no operator
// action. An earlier version ran `npx --no-install prettier`, which resolves
// the repository's own `node_modules/.bin/prettier` — so opening a hostile repo
// and editing any markdown file executed that repo's code. A formatter is only
// safe to run if it came from the operator's environment, never from the tree
// under review.
func trustedTool(name string) string {
	bin, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	resolved, err := filepath.EvalSymlinks(bin)
	if err != nil {
		return ""
	}
	if underRepo(resolved) {
		return ""
	}
	return resolved
}

// underRepo reports whether a path lies inside the working directory's
// repository, or inside any node_modules — the two places a tree can put an
// executable that PATH might reach.
func underRepo(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"node_modules"+string(os.PathSeparator)) {
		return true
	}
	wd, err := os.Getwd()
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(wd, path)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = nil, nil
	_ = cmd.Run()
}

// correctionMarkers match the shapes a correction takes. Deliberately loose:
// a false positive costs one log line that /retro reads with judgment, while a
// missed correction is a lesson lost for good.
var correctionMarkers = regexp.MustCompile(`(?i)^no[,. !]|not what i|i (said|told you)|i asked (for|you)|that.?s (wrong|not right|incorrect)|don.?t do that|stop (doing|that)|undo th|revert th|why did you|you (didn.?t|ignored|missed|forgot)|(^|[^[:alnum:]])instead([^[:alnum:]]|$)`)

// CaptureFeedback appends likely corrections to the learning-loop log mined by
// /retro and /distill. It writes nothing to stdout: UserPromptSubmit stdout is
// injected into model context, and this must stay invisible.
func CaptureFeedback(p hookio.Payload) {
	prompt := strings.TrimSpace(p.Prompt)
	// Harness notifications and reminders arrive on this event too, wrapped in
	// XML. They are not the user correcting anything.
	if prompt == "" || strings.HasPrefix(prompt, "<") {
		return
	}
	if !correctionMarkers.MatchString(prompt) {
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	dir := filepath.Join(home, ".agents", "feedback")
	if os.MkdirAll(dir, 0o700) != nil {
		return
	}
	entry, err := json.Marshal(map[string]string{
		"ts":      time.Now().UTC().Format(time.RFC3339),
		"project": p.CWD,
		"session": p.SessionID,
		"source":  "auto-hook",
		"prompt":  truncate(prompt, 400),
	})
	if err != nil {
		return
	}
	f, err := os.OpenFile(filepath.Join(dir, "log.jsonl"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintln(f, string(entry))
}

// statusPayload is the session JSON the status line receives.
type statusPayload struct {
	Model struct {
		DisplayName string `json:"display_name"`
	} `json:"model"`
	Workspace struct {
		CurrentDir string `json:"current_dir"`
	} `json:"workspace"`
	CWD string `json:"cwd"`
}

// StatusLine renders "model · branch · gcp project" for the terminal chrome.
func StatusLine(raw []byte) string {
	var s statusPayload
	_ = json.Unmarshal(raw, &s)

	model := s.Model.DisplayName
	if model == "" {
		model = "?"
	}
	dir := firstString(s.Workspace.CurrentDir, s.CWD, ".")

	out := "⚡ " + model
	if branch, err := exec.Command("git", "-C", dir, "branch", "--show-current").Output(); err == nil {
		if b := strings.TrimSpace(string(branch)); b != "" {
			out += "  ⎇ " + b
		}
	}
	if p := gcloudProject(); p != "" {
		out += "  ☁ " + p
	}
	return out
}

// gcloudProject reads the active project from gcloud's config files directly.
// Shelling out to gcloud would cost hundreds of milliseconds on every render.
func gcloudProject() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	active, err := os.ReadFile(filepath.Join(home, ".config", "gcloud", "active_config"))
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(home, ".config", "gcloud", "configurations",
		"config_"+strings.TrimSpace(string(active))))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if rest, found := strings.CutPrefix(strings.TrimSpace(line), "project"); found {
			return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(rest), "="))
		}
	}
	return ""
}

func mapString(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	s, _ := m[key].(string)
	return s
}

func firstString(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
