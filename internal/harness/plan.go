package harness

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// Deployment roots, relative to $HOME. Claude Code reads its user-level
// customizations from ~/.claude; the Antigravity CLI scans ~/.gemini/config.
const (
	ClaudeRoot = ".claude"
	AgyRoot    = ".gemini/config"
	BinPath    = ".local/bin/sh-agent"
)

// Supported vendor identifiers.
const (
	VendorAll         = "all"
	VendorAntigravity = "antigravity"
	VendorClaude      = "claude"
)

// NormalizeVendor parses and standardizes vendor names. It accepts "all",
// "antigravity" (or "agy", "gemini"), and "claude". An empty string defaults to "all".
func NormalizeVendor(v string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "", "all":
		return VendorAll, nil
	case "antigravity", "agy", "gemini":
		return VendorAntigravity, nil
	case "claude", "claudecode":
		return VendorClaude, nil
	default:
		return "", fmt.Errorf("unknown vendor %q: must be all, antigravity (or agy), or claude", v)
	}
}

// File is one path the install produces, relative to $HOME.
type File struct {
	Rel  string
	Data []byte
	Mode fs.FileMode

	// KeepForeign marks a file the operator may already own. If it exists
	// without the sh-agent marker the installer leaves it alone and says so,
	// rather than overwriting a hand-written CLAUDE.md.
	KeepForeign bool
}

// Plan computes every file the install writes. Each repo directory maps to
// one place per vendor, so a file dropped into the repo deploys without a
// code change:
//
//	context/*.md          → .claude/AGENTS.md (+ CLAUDE.md stub), .gemini/config/rules/AGENTS.md
//	skills/<name>/**      → .claude/skills/<name>/**,   .gemini/config/skills/<name>/**
//	agents/<name>.md      → .claude/agents/<name>.md,   .gemini/config/agents/<name>/agent.md (derived)
//	commands/<name>.md    → .claude/commands/<name>.md, .gemini/config/workflows/<name>.md
//	workflows/<name>.md   → same as commands/
//	rules/<name>.md       → .claude/rules/<name>.md,    .gemini/config/rules/<name>.md
//
// src is the checkout, plain copy, or embedded files to read.
// vendor selects which targets to generate ("all", "antigravity", or "claude").
// selfExe is the running binary, copied to ~/.local/bin so `sh-agent update`
// works from any directory; pass nil to skip it.
func Plan(src Source, selfExe []byte, vendor string) ([]File, error) {
	v, err := NormalizeVendor(vendor)
	if err != nil {
		return nil, err
	}
	wantClaude := v == VendorAll || v == VendorClaude
	wantAgy := v == VendorAll || v == VendorAntigravity

	inv, err := LoadInventory(src)
	if err != nil {
		return nil, err
	}
	context, err := RenderContext(src, inv)
	if err != nil {
		return nil, err
	}

	var files []File
	if wantClaude {
		files = append(files,
			File{Rel: filepath.Join(ClaudeRoot, "AGENTS.md"), Data: context, Mode: 0o644},
			File{Rel: filepath.Join(ClaudeRoot, "CLAUDE.md"), Data: []byte(claudeStub), Mode: 0o644, KeepForeign: true},
		)
	}
	if wantAgy {
		files = append(files,
			File{Rel: filepath.Join(AgyRoot, "rules", "AGENTS.md"), Data: context, Mode: 0o644},
		)
	}

	for _, s := range inv.Skills {
		err := fs.WalkDir(src.FS, s.Dir, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			// A link is refused, not followed: reading it would deploy the
			// TARGET'S bytes under the LINK'S mode (d.Info reports the link, and
			// `|0o600` on a link's 0777 ships a world-writable file into
			// ~/.claude) — an arbitrary-file read into every consumer home
			// (security review 2026-09-14). A skill ships regular files only.
			if d.Type()&fs.ModeSymlink != 0 {
				return fmt.Errorf("%s is a symlink — a skill ships regular files only", src.path(p))
			}
			if d.IsDir() || d.Name() == ".DS_Store" {
				return nil
			}
			// A secret-shaped file is refused rather than copied into both
			// vendor directories on every machine. The repo's .gitignore
			// keeps these out of git, which means review never sees them and
			// they can sit in a skill directory unnoticed; install is the one
			// step that would spread them (security review 2026-09-16). The
			// names match the context's Git rule.
			if secretShaped(d.Name()) {
				return fmt.Errorf("%s is a secret-shaped file — install refuses to deploy it; move it out of the skill", src.path(p))
			}
			rel := filepath.FromSlash(strings.TrimPrefix(p, s.Dir+"/"))
			info, err := d.Info()
			if err != nil {
				return err
			}
			data, err := fs.ReadFile(src.FS, p)
			if err != nil {
				return err
			}
			// Mask off group and other WRITE before deploying. A repo file
			// that is 0664 or 0666 — a loose umask, a copied archive —
			// would otherwise land group-writable in ~/.claude, where its
			// contents are instructions every session reads: another local
			// account could rewrite the harness without touching the repo
			// (security review of v1.0.0). Read and execute bits are kept.
			mode := info.Mode().Perm()&0o755 | 0o600
			// Both CLIs get their own copy: neither is documented to read the
			// other's directory, and a copy has no dangling-symlink failure mode.
			if wantClaude {
				files = append(files, File{Rel: filepath.Join(ClaudeRoot, "skills", s.Name, rel), Data: data, Mode: mode})
			}
			if wantAgy {
				files = append(files, File{Rel: filepath.Join(AgyRoot, "skills", s.Name, rel), Data: data, Mode: mode})
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("skill %s: %w", s.Name, err)
		}
	}

	for _, a := range inv.Agents {
		if wantClaude {
			files = append(files, File{Rel: filepath.Join(ClaudeRoot, "agents", a.Name+".md"), Data: a.Raw, Mode: 0o644})
		}
		if wantAgy {
			derived, err := DeriveAgy(a)
			if err != nil {
				return nil, err
			}
			files = append(files, File{Rel: filepath.Join(AgyRoot, "agents", a.Name, "agent.md"), Data: derived, Mode: 0o644})
		}
	}

	// A command is a markdown prompt invoked as /name on both vendors: Claude
	// Code calls it a command, Antigravity a workflow. Copied verbatim; both
	// ignore frontmatter keys they do not know.
	for _, c := range inv.Commands {
		if wantClaude {
			files = append(files, File{Rel: filepath.Join(ClaudeRoot, "commands", c.Name+".md"), Data: c.Raw, Mode: 0o644})
		}
		if wantAgy {
			files = append(files, File{Rel: filepath.Join(AgyRoot, "workflows", c.Name+".md"), Data: c.Raw, Mode: 0o644})
		}
	}

	for _, r := range inv.Rules {
		if wantClaude {
			files = append(files, File{Rel: filepath.Join(ClaudeRoot, "rules", r.Name+".md"), Data: r.Raw, Mode: 0o644})
		}
		if wantAgy {
			files = append(files, File{Rel: filepath.Join(AgyRoot, "rules", r.Name+".md"), Data: r.Raw, Mode: 0o644})
		}
	}

	if selfExe != nil {
		files = append(files, File{Rel: BinPath, Data: selfExe, Mode: 0o755})
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Rel < files[j].Rel })
	for i := 1; i < len(files); i++ {
		if files[i].Rel == files[i-1].Rel {
			return nil, fmt.Errorf("two sources map to ~/%s (a rule may not be named AGENTS.md)", files[i].Rel)
		}
	}
	return files, nil
}

// IsOurs reports whether deployed content was written by sh-agent.
func IsOurs(data []byte) bool {
	return strings.HasPrefix(string(data), Marker)
}

// secretShaped reports whether a filename matches the credential shapes the
// context's Git rule stops a commit for. install deploys a skill's files to
// every vendor directory on every machine, so one that slipped past
// .gitignore into a skill directory would be spread rather than caught.
// SecretShaped reports whether a file name looks like a secret (.env, a
// key, a token file). The search index skips such files too.
func SecretShaped(name string) bool { return secretShaped(name) }

func secretShaped(name string) bool {
	for _, pattern := range []string{
		".env", ".env.*", "*.env", ".envrc", ".netrc", ".npmrc", ".pypirc",
		"*.pem", "*.key", "*.p12", "*.pfx", "*.keystore", "*.jks",
		"*.tfvars", "*.tfstate", "*.tfstate.backup",
		"*credentials*.json", "*-key.json", "id_rsa", "id_ed25519", "id_ecdsa", "id_dsa",
		"*oauth-token*", // the Antigravity CLI keeps its login in antigravity-oauth-token
	} {
		if ok, _ := filepath.Match(pattern, name); ok {
			return true
		}
	}
	return false
}
