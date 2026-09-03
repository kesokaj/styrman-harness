package trail

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"sync"
)

// sourceExts are the file types whose change owes the DoD rows.
var sourceExts = map[string]bool{
	".go": true, ".ts": true, ".tsx": true, ".js": true, ".mjs": true,
	".cjs": true, ".jsx": true, ".sh": true, ".py": true, ".tf": true,
	".tfvars": true, ".sql": true,
}

// deployDirs are directories whose YAML/JSON manifests ship — a change to one
// deploys something, so it owes the same rows as code.
var deployDirs = map[string]bool{
	".github": true, ".circleci": true, ".gitlab": true, "k8s": true,
	"kubernetes": true, "helm": true, "charts": true, "overlays": true,
	"base": true, "deploy": true, "deployments": true, "manifests": true,
	"infra": true,
}

// sourceNames are extensionless source files, build definitions, and the
// dependency manifests that form the security row's surface.
var sourceNames = map[string]bool{
	"Makefile": true, "makefile": true, "Justfile": true, "justfile": true,
	"cloudbuild.yml": true, "cloudbuild.yaml": true, "cloudbuild.json": true,
	"skaffold.yml": true, "skaffold.yaml": true, "Chart.yaml": true,
	".gitlab-ci.yml": true, "kustomization.yaml": true,
	"go.mod": true, "go.sum": true, "package.json": true,
	"package-lock.json": true, "pnpm-lock.yaml": true, "yarn.lock": true,
	"pyproject.toml": true, "poetry.lock": true, "uv.lock": true,
	"Cargo.toml": true, "Cargo.lock": true, "Gemfile": true, "Gemfile.lock": true,
}

// IsSourcePath reports whether changing this path owes the DoD rows.
//
// Markdown and anything under the docs root never trigger it by default, and
// neither does plain config: a settings.json or editor toml tweak does not owe
// a review, a verification and a changelog. Code, IaC, container builds,
// deployment and CI manifests, and dependency manifests do.
//
// A project whose product is prose — a harness whose deliverable is a
// constitution and a set of skills — declares its own source in
// `.agents/harness.json` under `sourceGlobs`, which is consulted before the
// markdown exclusion. Without that escape hatch, rewriting every agent
// definition in this very repo would report that nothing is owed.
func (p Project) IsSourcePath(rel string) bool {
	rel = strings.TrimPrefix(toSlash(rel), "./")
	if rel == ".agents" || strings.HasPrefix(rel, ".agents/") {
		return false
	}
	if docsRel := toSlash(p.DocsRel()); docsRel != "." && strings.HasPrefix(rel, docsRel+"/") {
		return false
	}
	for _, glob := range p.SourceGlobs {
		if matchGlob(glob, rel) {
			return true
		}
	}
	if strings.HasSuffix(rel, ".md") {
		return false
	}

	base := path.Base(rel)
	if sourceExts[path.Ext(rel)] {
		return true
	}
	if strings.HasSuffix(rel, ".tf.json") {
		return true
	}
	if sourceNames[base] {
		return true
	}
	if strings.HasPrefix(base, "Dockerfile") || strings.HasPrefix(base, "Containerfile") ||
		strings.HasSuffix(base, ".mk") {
		return true
	}
	if strings.Contains(base, "compose") &&
		(strings.HasSuffix(base, ".yml") || strings.HasSuffix(base, ".yaml")) {
		return true
	}
	if strings.HasPrefix(base, "values") && strings.HasSuffix(base, ".yaml") {
		return true
	}
	if strings.HasPrefix(base, "requirements") && strings.HasSuffix(base, ".txt") {
		return true
	}

	// Deployment and CI manifests, at any depth.
	if ext := path.Ext(rel); ext == ".yaml" || ext == ".yml" || ext == ".json" {
		for _, seg := range strings.Split(path.Dir(rel), "/") {
			if deployDirs[seg] {
				return true
			}
		}
	}
	return false
}

// globCache memoises compiled source globs; IsSourcePath runs once per changed
// file per gate invocation, and the pattern set is tiny and fixed per project.
var globCache sync.Map // string -> *regexp.Regexp

// matchGlob reports whether a path matches a glob, supporting `**` for "any
// number of path segments".
//
// Go's path.Match deliberately does not implement `**` — its `*` never crosses
// a separator, so `**/*.go` matches exactly one directory level. Since the
// documented way to declare a project's source is patterns like `**/*.go` and
// `claude/agents/**`, the matcher has to mean what the documentation says.
func matchGlob(glob, name string) bool {
	re, ok := globCache.Load(glob)
	if !ok {
		compiled, err := compileGlob(glob)
		if err != nil {
			return false
		}
		re, _ = globCache.LoadOrStore(glob, compiled)
	}
	return re.(*regexp.Regexp).MatchString(name)
}

// compileGlob translates a glob into an anchored regexp.
//
// It walks runes, not bytes: quoting one byte of a multi-byte character would
// produce a pattern that can never match, and a gate that silently matches
// nothing is a gate that silently switches itself off. Character classes are
// passed through so the syntax stays a superset of path.Match's rather than
// quietly turning `[a-z]*.go` into a search for a literal bracket.
func compileGlob(glob string) (*regexp.Regexp, error) {
	var b strings.Builder
	b.WriteString("^")

	runes := []rune(glob)
	for i := 0; i < len(runes); i++ {
		rest := string(runes[i:])
		switch {
		case strings.HasPrefix(rest, "**/"):
			// Any number of leading segments, including none.
			b.WriteString("(?:[^/]+/)*")
			i += 2
		case strings.HasPrefix(rest, "/**"):
			// Everything beneath this prefix.
			b.WriteString("(?:/.*)?")
			i += 2
		case strings.HasPrefix(rest, "**"):
			b.WriteString(".*")
			i++
		case runes[i] == '*':
			b.WriteString("[^/]*")
		case runes[i] == '?':
			b.WriteString("[^/]")
		case runes[i] == '[':
			class, next, err := compileClass(runes, i)
			if err != nil {
				return nil, err
			}
			b.WriteString(class)
			i = next
		default:
			b.WriteString(regexp.QuoteMeta(string(runes[i])))
		}
	}

	b.WriteString("$")
	return regexp.Compile(b.String())
}

// compileClass translates a `[...]` character class, returning the regexp
// fragment and the index of its closing bracket. An unterminated class is an
// error rather than a literal: silently reinterpreting a malformed pattern is
// how a gate ends up matching something nobody asked for.
func compileClass(runes []rune, start int) (string, int, error) {
	var b strings.Builder
	b.WriteRune('[')

	i := start + 1
	if i < len(runes) && (runes[i] == '!' || runes[i] == '^') {
		b.WriteRune('^')
		i++
	}
	// A ']' immediately after the opener is a literal, per glob convention.
	if i < len(runes) && runes[i] == ']' {
		b.WriteString(`\]`)
		i++
	}
	for ; i < len(runes); i++ {
		switch runes[i] {
		case ']':
			b.WriteRune(']')
			return b.String(), i, nil
		case '\\', '^':
			b.WriteRune('\\')
			b.WriteRune(runes[i])
		default:
			b.WriteRune(runes[i])
		}
	}
	return "", 0, fmt.Errorf("unterminated character class in glob %q", string(runes))
}

// securityPatterns name the surfaces that make the security row due: auth,
// crypto, secrets, input parsing, and dependency manifests.
var securityPatterns = []string{
	"auth", "crypt", "secret", "credential", "token", "session", "jwt", "oidc",
	"password", "login", "permission", "policy", "sanitiz", "validat", "parse",
}

// OwesSecurity reports whether any changed path touches the security row's
// surface. Dependency manifests always do — a bumped version is the classic
// supply-chain change.
func (p Project) OwesSecurity(changed []string) bool {
	for _, rel := range changed {
		lower := strings.ToLower(rel)
		if sourceNames[path.Base(rel)] && isDependencyManifest(path.Base(rel)) {
			return true
		}
		for _, pat := range securityPatterns {
			if strings.Contains(lower, pat) {
				return true
			}
		}
	}
	return false
}

// TouchesDependencyManifest reports whether any changed path is a dependency manifest.
func (p Project) TouchesDependencyManifest(changed []string) bool {
	for _, rel := range changed {
		if sourceNames[path.Base(rel)] && isDependencyManifest(path.Base(rel)) {
			return true
		}
	}
	return false
}

func isDependencyManifest(base string) bool {
	switch base {
	case "go.mod", "go.sum", "package.json", "package-lock.json", "pnpm-lock.yaml",
		"yarn.lock", "pyproject.toml", "poetry.lock", "uv.lock", "Cargo.toml",
		"Cargo.lock", "Gemfile", "Gemfile.lock":
		return true
	}
	return false
}

// toSlash normalises separators without importing path/filepath into
// the pure-string helpers above; git always reports forward slashes.
func toSlash(s string) string { return strings.ReplaceAll(s, "\\", "/") }
