package trail

import (
	"path"
	"strings"
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
		if ok, _ := path.Match(glob, rel); ok {
			return true
		}
		// A trailing /** also covers everything below the prefix.
		if prefix, found := strings.CutSuffix(glob, "/**"); found && strings.HasPrefix(rel, prefix+"/") {
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
