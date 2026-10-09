package metrics

import (
	"path"
	"regexp"
	"strings"
)

// subcommanded are the programs whose first non-flag word is kept as a
// subcommand (spec § Normalised verb).
var subcommanded = map[string]bool{
	"git": true, "go": true, "gh": true, "docker": true, "gcloud": true, "kubectl": true,
	"npm": true, "npx": true, "sh-agent": true, "terraform": true,
}

// gitGlobalArg are git's global options that take the next word as their
// value (`git -C dir`, `git -c k=v`); the others are single words.
var gitGlobalArg = map[string]bool{"-C": true, "-c": true, "--git-dir": true, "--work-tree": true, "--namespace": true, "--exec-path": true, "--config-env": true}

// skipGitGlobals drops git's global options and their values from the front
// of its arguments, so the subcommand behind them is found (`git -C dir add
// -A` is `git add -A`).
func skipGitGlobals(rest []string) []string {
	for len(rest) > 0 && strings.HasPrefix(rest[0], "-") {
		name, _, hasValue := strings.Cut(rest[0], "=")
		rest = rest[1:]
		if gitGlobalArg[name] && !hasValue && len(rest) > 0 {
			rest = rest[1:]
		}
	}
	return rest
}

// shAgentSub are sh-agent's second-level verbs, kept after the subcommand.
var shAgentSub = map[string]bool{"ingest": true, "show": true, "report": true, "export": true, "add": true}

// goStyle programs take long flags behind a single dash (`go test -race
// -count=1`) and never glue a value to a flag, so the whole name is kept.
var goStyle = map[string]bool{"go": true, "gofmt": true, "govulncheck": true, "find": true}

// dropped are wrapper words skipped before the program.
var dropped = map[string]bool{"sudo": true, "env": true, "command": true, "time": true, "nohup": true, "exec": true}

// wrapperArg are wrapper flags that take a value (`env -u NAME`,
// `sudo -u USER`, `env -C DIR`): the value is skipped with the flag.
var wrapperArg = map[string]bool{"-u": true, "--unset": true, "-g": true, "-C": true, "--chdir": true}

// keywordsLead open a compound command and are followed by one (`do git
// add`); keywordsEnd close one and carry none; keywordsSkip open a segment
// that is not a command (`for f in a b`).
var (
	keywordsLead = map[string]bool{"if": true, "then": true, "else": true, "elif": true, "do": true, "while": true, "until": true, "!": true}
	keywordsEnd  = map[string]bool{"done": true, "fi": true, "esac": true, "in": true}
	keywordsSkip = map[string]bool{"for": true, "case": true, "select": true, "function": true}
)

// programName is what a command name looks like. Anything else in command
// position (a glob, an expansion, a quoted path) is content and yields no
// verb.
var (
	programName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_.+-]*$`)
	longFlag    = regexp.MustCompile(`^--[A-Za-z0-9][A-Za-z0-9_-]*$`)
	letters     = regexp.MustCompile(`^[A-Za-z]+$`)
	// redirection is a word that opens with one (`>f`, `2>>log`, `<in`,
	// `&>f`); redirectionOnly is one whose target is the next word.
	redirection     = regexp.MustCompile(`^(?:[0-9]*|&)[<>]`)
	redirectionOnly = regexp.MustCompile(`^(?:[0-9]*|&)(?:>>?|<<?|>&|<&)$`)
)

// Verb reduces a shell command line to its shape: one string per command in
// it, "program [subcommand] [flag names…]". No positional word, no flag value
// and no comment is kept, except a bare "." after `git add` (D22). Content in
// the shell constructs the tokenizer models — quotes, here-documents, ${…},
// $((…)), case patterns, comments — never reaches a verb; a word in command
// position must look like a command name.
func Verb(cmd string) []string {
	var out []string
	for _, seg := range segments(stripHeredocs(cmd)) {
		if v := segmentVerb(seg); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// tokenizer splits a command line into segments of words.
type tokenizer struct {
	segs        [][]string
	words       []string
	cur         strings.Builder
	inWord      bool
	parenDepth  int  // open ( and $( subshells
	caseDepth   int  // open case … esac
	wantPattern bool // inside a case, before the next pattern's )
}

func (t *tokenizer) flushWord() {
	if !t.inWord {
		return
	}
	w := t.cur.String()
	t.cur.Reset()
	t.inWord = false
	if len(t.words) == 0 {
		switch w {
		case "case":
			t.caseDepth++
		case "esac":
			if t.caseDepth > 0 {
				t.caseDepth--
			}
			t.wantPattern = false
		}
	}
	t.words = append(t.words, w)
	if w == "in" && len(t.words) >= 3 && t.words[0] == "case" {
		t.words = nil // the case header carries no command
		t.wantPattern = true
	}
}

func (t *tokenizer) flushSeg() {
	t.flushWord()
	if len(t.words) > 0 && !t.wantPattern {
		t.segs = append(t.segs, t.words)
	}
	t.words = nil
}

// segments splits cmd at ; && || | & newlines and subshell parentheses,
// outside quotes and expansions, and returns each segment's words.
func segments(cmd string) [][]string {
	t := &tokenizer{}
	rs := []rune(cmd)
	at := func(i int, r rune) bool { return i < len(rs) && rs[i] == r }
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		switch {
		case r == '\'' || r == '"':
			// a quoted run is part of the current word; its content is never split
			q := r
			t.inWord = true
			for i++; i < len(rs) && rs[i] != q; i++ {
				if q == '"' && rs[i] == '\\' && i+1 < len(rs) {
					i++
				}
				t.cur.WriteRune(rs[i])
			}
		case r == '\\' && i+1 < len(rs):
			i++
			if rs[i] != '\n' {
				t.cur.WriteRune(rs[i])
				t.inWord = true
			}
		case r == '$' && at(i+1, '{'):
			i = t.swallow(rs, i+1, '{', '}')
		case r == '$' && at(i+1, '(') && at(i+2, '('):
			i = t.swallow(rs, i+1, '(', ')')
		case r == '$' && at(i+1, '('):
			t.flushSeg()
			t.parenDepth++
			i++
		case r == '#' && !t.inWord:
			for i+1 < len(rs) && rs[i+1] != '\n' {
				i++
			}
		case r == '&' && (i > 0 && (rs[i-1] == '>' || rs[i-1] == '<') || at(i+1, '>')):
			t.cur.WriteRune(r) // a redirection (2>&1, &>f), never a separator
			t.inWord = true
		case r == ';' && at(i+1, ';'):
			i++
			t.flushSeg()
			if t.caseDepth > 0 {
				t.wantPattern = true
			}
		case r == '|' && t.wantPattern:
			t.flushWord() // an alternative in a case pattern
		case r == ';' || r == '\n' || r == '|' || r == '&' || r == '`':
			if (r == '|' || r == '&') && at(i+1, r) {
				i++
			}
			t.flushSeg()
		case r == '(':
			t.flushSeg()
			t.parenDepth++
		case r == ')':
			switch {
			case t.wantPattern || t.parenDepth == 0:
				// the end of a case pattern: its words are content
				t.flushWord()
				t.words = nil
				t.wantPattern = false
			default:
				t.parenDepth--
				t.flushSeg()
			}
		case (r == '{' || r == '}') && !t.inWord:
			t.flushSeg()
		case r == ' ' || r == '\t':
			t.flushWord()
		default:
			t.cur.WriteRune(r)
			t.inWord = true
		}
	}
	t.flushSeg()
	return t.segs
}

// swallow appends the balanced open…close run starting at rs[i] to the
// current word, unparsed, and returns the index of its last rune.
func (t *tokenizer) swallow(rs []rune, i int, open, close rune) int {
	t.inWord = true
	t.cur.WriteRune('$')
	depth := 0
	for ; i < len(rs); i++ {
		t.cur.WriteRune(rs[i])
		switch rs[i] {
		case open:
			depth++
		case close:
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return i
}

func segmentVerb(words []string) string {
	i := 0
	for i < len(words) {
		w := words[i]
		switch {
		case keywordsLead[w], dropped[w], isAssignment(w):
			i++
		case redirection.MatchString(w):
			i++ // a leading redirection: >file, 2>err, or > and its target
			if redirectionOnly.MatchString(w) && i < len(words) {
				i++
			}
		case strings.HasPrefix(w, "-") && i > 0 && wrapper(words[:i]):
			i++ // a wrapper's own flag (env -u X, sudo -n)
			if wrapperArg[w] && i < len(words) {
				i++ // and the value it takes
			}
		default:
			return verbFrom(words[i:])
		}
	}
	return ""
}

func verbFrom(words []string) string {
	if keywordsEnd[words[0]] || keywordsSkip[words[0]] {
		return ""
	}
	prog := path.Base(words[0])
	rest := words[1:]
	if prog == "go" && len(rest) >= 2 && rest[0] == "run" {
		prog, rest = path.Base(strings.TrimSuffix(rest[1], "/")), rest[2:]
	}
	if !programName.MatchString(prog) {
		return ""
	}
	parts := []string{prog}
	if prog == "git" {
		rest = skipGitGlobals(rest)
	}
	if subcommanded[prog] && len(rest) > 0 && programName.MatchString(rest[0]) {
		parts = append(parts, rest[0])
		rest = rest[1:]
	}
	// sh-agent's own second-level verbs (`metrics show`, `memory add`) are
	// identifiers from a fixed list; any other second word stays a positional.
	if prog == "sh-agent" && len(parts) == 2 && len(rest) > 0 && shAgentSub[rest[0]] {
		parts = append(parts, rest[0])
		rest = rest[1:]
	}
	gitAdd := prog == "git" && len(parts) == 2 && parts[1] == "add"
	for _, w := range rest {
		if f := flagName(prog, w); f != "" {
			parts = append(parts, f)
		} else if gitAdd && w == "." {
			parts = append(parts, ".")
		}
	}
	return strings.Join(parts, " ")
}

// flagName is the name of the flag w, or "" when w is no flag or its name
// cannot be told from a value glued to it. `--name=v` keeps `--name`; a
// go-style `-name=v` keeps `-name`; a short-flag group of at most three
// letters (`-la`, `-rnw`) is kept whole; anything longer or with a digit
// keeps its first letter only (`-c1-200`, `-psecret`, `-xzvf` → `-c`, `-p`,
// `-x`), since a short option can carry its value glued on.
func flagName(prog, w string) string {
	if strings.HasPrefix(w, "--") {
		name, _, _ := strings.Cut(w, "=")
		if longFlag.MatchString(name) {
			return name
		}
		return ""
	}
	if len(w) < 2 || w[0] != '-' {
		return ""
	}
	name, _, hasValue := strings.Cut(w[1:], "=")
	switch {
	case goStyle[prog] && letters.MatchString(name):
		return "-" + name
	case !letters.MatchString(w[1:2]):
		return "" // -20, -1: a value, not a flag name
	case !hasValue && len(w) <= 4 && letters.MatchString(w[1:]):
		return w
	}
	return w[:2]
}

// wrapper reports whether the words before a flag are a wrapper and its own
// flags, so the flag belongs to the wrapper and not to the program.
func wrapper(before []string) bool {
	for _, w := range before {
		if !dropped[w] && !strings.HasPrefix(w, "-") && !isAssignment(w) {
			return false
		}
	}
	return true
}

// stripHeredocs removes every here-document body from cmd, so that the
// lines of a file written with `cat <<EOF` never become commands (they are
// content, and the store keeps none).
func stripHeredocs(cmd string) string {
	lines := strings.Split(cmd, "\n")
	var out []string
	var pending []string // delimiters opened on the current line, in order
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		out = append(out, line)
		pending = append(pending, heredocDelims(line)...)
		for len(pending) > 0 {
			d := pending[0]
			pending = pending[1:]
			for i+1 < len(lines) {
				i++
				if strings.TrimLeft(lines[i], "\t") == d {
					break
				}
			}
		}
	}
	return strings.Join(out, "\n")
}

// heredocDelims returns the delimiters of the here-documents line opens:
// `<<WORD`, `<<-WORD`, `<<'WORD'`, `<<"WORD"`, `<<\WORD`; `<<<` is a
// here-string and opens none.
func heredocDelims(line string) []string {
	var ds []string
	for {
		j := strings.Index(line, "<<")
		if j < 0 {
			return ds
		}
		line = line[j+2:]
		if strings.HasPrefix(line, "<") {
			line = line[1:]
			continue
		}
		line = strings.TrimPrefix(line, "-")
		line = strings.TrimLeft(line, " \t")
		q := ""
		switch {
		case strings.HasPrefix(line, "'") || strings.HasPrefix(line, "\""):
			q, line = line[:1], line[1:]
		case strings.HasPrefix(line, "\\"):
			line = line[1:] // <<\EOF: a quoted delimiter, the same here-doc
		}
		end := strings.IndexFunc(line, func(r rune) bool {
			return !(r == '_' || r == '-' || r == '.' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9')
		})
		if end < 0 {
			end = len(line)
		}
		if end > 0 {
			ds = append(ds, line[:end])
		}
		line = strings.TrimPrefix(line[end:], q)
	}
}

func isAssignment(w string) bool {
	name, _, ok := strings.Cut(w, "=")
	if !ok || name == "" {
		return false
	}
	for i, r := range name {
		if !(r == '_' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || i > 0 && r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}
