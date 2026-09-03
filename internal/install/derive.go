package install

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// modelMap translates Claude model aliases to Antigravity tiers. Anything
// unlisted is omitted, which agy reads as "inherit".
var modelMap = map[string]string{"opus": "pro", "haiku": "flash"}

// toolMap translates Claude tool names to agy ones. An unmapped tool is an
// error rather than a guess: Antigravity's documented behaviour for an
// unrecognised tool name is that the subagent hangs during dispatch, so a
// silent pass-through would trade a build failure for a mystery stall.
var toolMap = map[string]string{
	"Read":  "view_file",
	"Grep":  "grep_search",
	"Glob":  "find_by_name",
	"Bash":  "run_command",
	"Write": "write_to_file",
	"Edit":  "replace_file_content",
}

// AgyToolNames lists the agy tool names this harness depends on, so doctor can
// check them against the installed CLI.
func AgyToolNames() []string {
	names := make([]string, 0, len(toolMap))
	for _, v := range toolMap {
		names = append(names, v)
	}
	return names
}

const agySkillsDir = "~/.gemini/config/skills/"

// Derive is the `harness derive` subcommand.
func Derive(args []string) error {
	if len(args) != 2 {
		return errors.New("usage: harness derive <claude-agent.md> <out.md>")
	}
	return DeriveFile(args[0], args[1])
}

// DeriveFile writes the Antigravity flavour of a Claude Code agent definition.
// Claude Code is the source of truth; this file is generated and should never
// be hand-edited.
func DeriveFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	front, body, err := splitFrontmatter(string(data))
	if err != nil {
		return err
	}

	name, description := field(front, "name"), field(front, "description")
	if name == "" || description == "" {
		return errors.New("frontmatter must contain single-line name and description")
	}

	out := []string{
		"name: " + name,
		"description: " + description,
		"subagent: true",
		"mainAgent: false",
	}
	if tier, ok := modelMap[field(front, "model")]; ok {
		out = append(out, "model: "+tier)
	}
	if tools := field(front, "tools"); tools != "" {
		mapped, err := mapTools(tools)
		if err != nil {
			return err
		}
		out = append(out, "tools: ["+strings.Join(mapped, ", ")+"]")
	}
	if skills := field(front, "skills"); skills != "" {
		var paths []string
		for _, s := range splitList(skills) {
			paths = append(paths, agySkillsDir+s)
		}
		out = append(out, "skills: ["+strings.Join(paths, ", ")+"]")
	}
	// Claude-only keys (maxTurns, effort, memory, permissionMode, isolation)
	// are dropped: agy has no equivalent and would reject or ignore them.

	body = strings.ReplaceAll(body, "~/.claude/skills/", agySkillsDir)
	content := "---\n" + strings.Join(out, "\n") + "\n---\n" + body

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, []byte(content), 0o644)
}

func mapTools(list string) ([]string, error) {
	var mapped []string
	for _, tool := range splitList(list) {
		agy, ok := toolMap[tool]
		if !ok {
			return nil, fmt.Errorf("tool %q has no agy mapping — add it to toolMap", tool)
		}
		mapped = append(mapped, agy)
	}
	return mapped, nil
}

// splitFrontmatter separates a leading YAML frontmatter block from the body.
func splitFrontmatter(text string) (front, body string, err error) {
	if !strings.HasPrefix(text, "---\n") {
		return "", "", errors.New("no YAML frontmatter block found")
	}
	rest := text[len("---\n"):]
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		return "", "", errors.New("unterminated YAML frontmatter block")
	}
	return rest[:end], rest[end+len("\n---\n"):], nil
}

// field reads a single-line scalar from a frontmatter block.
func field(front, key string) string {
	for _, line := range strings.Split(front, "\n") {
		name, value, found := strings.Cut(line, ":")
		if !found || strings.TrimSpace(name) != key {
			continue
		}
		return strings.TrimSpace(value)
	}
	return ""
}

func splitList(s string) []string {
	s = strings.Trim(strings.TrimSpace(s), "[]")
	var out []string
	for _, part := range strings.Split(s, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}
