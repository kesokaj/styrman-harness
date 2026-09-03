// Package harness installs the styrman harness (global context, skills,
// subagents) into Claude Code and the Antigravity CLI, and reports drift
// between the repo and what is deployed.
package harness

import (
	"errors"
	"regexp"
	"strings"
)

// Frontmatter is the YAML header of a SKILL.md or agent file, flattened to
// strings. Block scalars are folded onto one line, lists keep their bracketed
// source form, and nested maps are dropped. That is everything the installer
// needs, and it keeps the parser dependency-free.
type Frontmatter map[string]string

var keyLine = regexp.MustCompile(`^([A-Za-z0-9_-]+):(?:[ \t]+(.*))?$`)

// Split separates a leading `---` block from the body that follows it.
func Split(text string) (front, body string, err error) {
	first, rest, ok := strings.Cut(text, "\n")
	if !ok || strings.TrimRight(first, "\r") != "---" {
		return "", "", errors.New("no YAML frontmatter block")
	}
	lines := strings.SplitAfter(rest, "\n")
	for i, l := range lines {
		if strings.TrimRight(l, "\r\n") == "---" {
			return strings.Join(lines[:i], ""), strings.Join(lines[i+1:], ""), nil
		}
	}
	return "", "", errors.New("unterminated YAML frontmatter block")
}

// Parse reads the top-level keys of a frontmatter block.
func Parse(front string) Frontmatter {
	fm := Frontmatter{}
	lines := strings.Split(front, "\n")
	for i := 0; i < len(lines); i++ {
		m := keyLine.FindStringSubmatch(strings.TrimRight(lines[i], "\r"))
		if m == nil {
			continue
		}
		key, val := m[1], strings.TrimSpace(m[2])

		// Indented or blank lines after a key belong to it: a block scalar,
		// a nested map, or a plain scalar that wraps.
		var cont []string
		j := i + 1
		for ; j < len(lines); j++ {
			l := strings.TrimRight(lines[j], "\r")
			if l != "" && !strings.HasPrefix(l, " ") && !strings.HasPrefix(l, "\t") {
				break
			}
			if t := strings.TrimSpace(l); t != "" {
				cont = append(cont, t)
			}
		}

		switch {
		case isBlockIndicator(val):
			fm[key] = strings.Join(cont, " ")
		case val == "" && len(cont) > 0 && strings.Contains(cont[0], ": "):
			// A nested map such as `metadata:`. Not needed; skipped whole.
		case len(cont) > 0:
			fm[key] = strings.Join(append([]string{unquote(val)}, cont...), " ")
		default:
			fm[key] = unquote(val)
		}
		i = j - 1
	}
	return fm
}

func isBlockIndicator(v string) bool {
	switch v {
	case "|", "|-", "|+", ">", ">-", ">+":
		return true
	}
	return false
}

func unquote(v string) string {
	if len(v) >= 2 {
		switch {
		case v[0] == '"' && v[len(v)-1] == '"':
			inner := v[1 : len(v)-1]
			inner = strings.ReplaceAll(inner, `\"`, `"`)
			return strings.ReplaceAll(inner, `\\`, `\`)
		case v[0] == '\'' && v[len(v)-1] == '\'':
			return strings.ReplaceAll(v[1:len(v)-1], "''", "'")
		}
	}
	return v
}

// SplitList reads a YAML list written either as `[a, b]`, `a, b`, or `a b`.
func SplitList(s string) []string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	sep := ","
	if !strings.Contains(s, ",") {
		sep = " "
	}
	var out []string
	for _, p := range strings.Split(s, sep) {
		if p = unquote(strings.TrimSpace(p)); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// yamlQuote renders a string as a double-quoted YAML scalar so that colons,
// hashes, and quotes inside a description can never change its meaning.
func yamlQuote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}
