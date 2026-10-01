package harness

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// The one vendor setting the harness owns (R4, ADR-025): Claude Code's commit
// and PR attribution, switched off in the operator's ~/.claude/settings.json
// so no commit carries the trailer the context's Git rule forbids. Every
// function here is advisory: no state of that file fails an install or a
// status, and nothing from the file is echoed.

// AttributionSet writes the one vendor setting the harness owns (R4,
// ADR-025): Claude Code's attribution off, so no commit or PR carries the
// trailer the context's Git rule forbids. Only an absent key is written — a
// present value of any kind is the operator's and AttributionHint says when
// it is on. The write is a textual insertion before the object's closing
// brace, re-parsed before it lands, because encoding/json would reorder the
// operator's keys; every byte up to the object's last value stays, and the
// file keeps its own line endings. No state of the file
// fails the install: the deployment already succeeded and this file is
// advisory. dryRun prints what would be written and touches nothing.
func AttributionSet(home string, dryRun bool, out io.Writer) {
	const fragment = `"attribution": {"commit": "", "pr": ""}`
	rel := filepath.Join(ClaudeRoot, "settings.json")
	path := filepath.Join(home, rel)
	// A linked directory above the file (a dotfiles checkout's ~/.claude) is
	// never written through either: the same rule every installed path follows
	// (underLink, BUG-129).
	if underLink(home, rel) {
		fmt.Fprintf(out, "settings: a directory above %s is a symlink; add the key to its target by hand\n", path)
		return
	}
	// Lstat, not Stat: a settings.json that is a symlink (a dotfiles checkout)
	// is never written through — a rename would replace the link with a file.
	info, err := os.Lstat(path)
	switch {
	case err == nil && info.Mode()&os.ModeSymlink != 0:
		fmt.Fprintf(out, "settings: %s is a symlink; add the key to its target by hand\n", path)
		return
	case errors.Is(err, os.ErrNotExist):
		if dryRun {
			fmt.Fprintf(out, "[dry-run] would create %s with %s\n", path, fragment)
			return
		}
		content := []byte("{\n  " + fragment + "\n}\n")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			settingsUnwritable(out, path, err)
			return
		}
		// O_EXCL: the Lstat above is advisory; whatever appears at the path
		// between the two calls (a planted link, a FIFO) makes this fail
		// instead of being written through or blocked on.
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			settingsUnwritable(out, path, err)
			return
		}
		if _, err := f.Write(content); err != nil {
			f.Close()
			settingsUnwritable(out, path, err)
			return
		}
		if err := f.Sync(); err != nil {
			f.Close()
			settingsUnwritable(out, path, err)
			return
		}
		if err := f.Close(); err != nil {
			settingsUnwritable(out, path, err)
			return
		}
		fmt.Fprintf(out, "settings: %s created with %s — Claude Code's commit and PR trailers are off\n", path, fragment)
		return
	case err != nil, !info.Mode().IsRegular():
		return // AttributionHint names it, and ReadSettings bounds the size
	}
	data, problem := ReadSettings(path)
	if problem != "" {
		return // AttributionHint names it
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil || top == nil {
		// malformed, not an object, or the bare value null (a nil map with no
		// error): attributionHint names each with one line
		return
	}
	if _, present := top["attribution"]; present {
		return
	}
	// Insert after the last byte of the object's body, keeping the operator's
	// indentation style (two spaces, the file's own newline before the brace)
	// and whatever follows the closing brace.
	trimmed := bytes.TrimRight(data, " \t\r\n")
	closing := bytes.LastIndexByte(trimmed, '}')
	opening := bytes.IndexByte(trimmed, '{')
	if closing <= opening {
		fmt.Fprintf(out, "settings: %s could not be edited safely (the result would not parse); add the key by hand\n", path)
		return
	}
	nl := "\n"
	if bytes.Contains(data, []byte("\r\n")) {
		nl = "\r\n" // keep the file's own line endings
	}
	body := bytes.TrimRight(trimmed[:closing], " \t\r\n")
	insert := "," + nl + "  " + fragment + nl
	if len(bytes.TrimSpace(trimmed[opening+1:closing])) == 0 {
		body = trimmed[:opening+1]
		insert = nl + "  " + fragment + nl
	}
	result := make([]byte, 0, len(data)+len(insert))
	result = append(result, body...)
	result = append(result, insert...)
	result = append(result, trimmed[closing:]...)
	result = append(result, data[len(trimmed):]...)
	if !json.Valid(result) {
		fmt.Fprintf(out, "settings: %s could not be edited safely (the result would not parse); add the key by hand\n", path)
		return
	}
	if dryRun {
		fmt.Fprintf(out, "[dry-run] would add %s to %s\n", fragment, path)
		return
	}
	if err := WriteAtomic(path, result, info.Mode().Perm()); err != nil {
		settingsUnwritable(out, path, err)
		return
	}
	fmt.Fprintf(out, "settings: %s — %s added; other keys untouched\n", path, fragment)
}

// settingsUnwritable is the one sentence for every way the settings file
// could not be written.
func settingsUnwritable(out io.Writer, path string, err error) {
	fmt.Fprintf(out, "settings: %s could not be written (%v); add the key by hand\n", path, err)
}

// ReadSettings reads the operator's settings file, bounded to a regular file
// of at most one megabyte; problem names why it could not be read, without a
// byte of its contents.
func ReadSettings(path string) (data []byte, problem string) {
	info, err := os.Stat(path)
	switch {
	case err != nil:
		return nil, "could not be read"
	case !info.Mode().IsRegular():
		return nil, "is not a regular file"
	case info.Size() > 1<<20:
		// One bound, here, so install and status say the same thing about
		// a file this tool will not rewrite.
		return nil, "is larger than 1 MiB, more than this tool will read or rewrite"
	}
	data, err = ReadCapped(path, 1<<20) // the one capped reader (S7)
	if err != nil {
		return nil, "could not be read"
	}
	return data, ""
}

// AttributionHint says when Claude Code will still add its Co-Authored-By
// trailer. The context's Git rule forbids the trailer, but a rule alone does
// not hold in a headless session; the CLI's own setting is what stops it, and
// settings.json is the operator's file, so install and status only say so. prefix is each command's note
// style ("note: " for install, "note       " for status). Advisory only: no
// state of that file may fail install or status, and nothing from it is
// echoed.
func AttributionHint(home, prefix string, out io.Writer) {
	path := filepath.Join(home, ClaudeRoot, "settings.json")
	off, problem := attributionOff(path)
	switch {
	case problem != "":
		fmt.Fprintf(out, "%s%s %s, so the attribution setting could not be checked; Claude Code adds a Co-Authored-By trailer to every commit unless it sets \"attribution\": {\"commit\": \"\"} (context § Git)\n", prefix, path, problem)
	case !off:
		fmt.Fprintf(out, "%s%s does not set \"attribution\": {\"commit\": \"\"}; Claude Code adds a Co-Authored-By trailer to every commit until it does (context § Git)\n", prefix, path)
	}
}

// attributionOff reads the operator's settings file, bounded to a regular
// file of at most one megabyte. off is true when attribution.commit is "" or
// false; problem names why the file could not be judged, without a byte of
// its contents (a JSON error string carries the input it choked on).
func attributionOff(path string) (off bool, problem string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false, ""
	}
	data, problem := ReadSettings(path)
	if problem != "" {
		return false, problem
	}
	var settings struct {
		Attribution struct {
			Commit any `json:"commit"`
		} `json:"attribution"`
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		if json.Valid(data) {
			return false, "is not a JSON object"
		}
		return false, "is not valid JSON"
	}
	switch v := settings.Attribution.Commit.(type) {
	case string:
		return v == "", ""
	case bool:
		return !v, ""
	}
	return false, ""
}
