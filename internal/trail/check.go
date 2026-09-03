package trail

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Freshness is the human-readable trail report used by /preflight, /audit and
// /handoff: does the documentation trail actually cover the code?
type Freshness struct {
	Lines []string
	Stale bool
}

func (f *Freshness) note(format string, args ...any) {
	f.Lines = append(f.Lines, fmt.Sprintf(format, args...))
}

func (f *Freshness) fail(format string, args ...any) {
	f.Stale = true
	f.note("STALE: "+format, args...)
}

// CheckFreshness reports whether the trail covers the code.
//
// A project that never adopted the contract passes: the gate binds projects
// that keep a trail, and failing a foreign clone for not having one would make
// the check useless in exactly the repos it should stay quiet in.
func (p Project) CheckFreshness() Freshness {
	var f Freshness
	f.note("Trail freshness — %s", p.Root)
	f.note("docs root: %s/", p.DocsRel())

	if !p.Opted() {
		f.note("note: no trail exists yet — %s/changes/ not found (contract not adopted)", p.DocsRel())
		f.note("RESULT: OK (nothing to check; run /init-project to adopt the contract)")
		return f
	}
	// Source classification is project-configurable, so the freshness report
	// must see the same set of files the dispatcher does.
	if cfg, err := LoadConfig(p.Root); err == nil {
		p.SourceGlobs = cfg.SourceGlobs
	}

	covered, head, newest := p.CoversHead()
	switch {
	case head == "":
		f.note("commits: none yet — HEAD coverage vacuously OK")
	case newest == "":
		f.fail("repo has commits (HEAD %s) but no changelog artifacts", head)
	case !covered:
		f.fail("newest changelog entry (%s) is older than HEAD (%s)", newest, head)
	default:
		f.note("OK: newest changelog entry (%s) covers HEAD (%s)", newest, head)
	}

	today := Today()
	if p.HasArtifact(RowChangelog, today) {
		f.note("OK: changelog entry for today (%s) present", today)
	} else if changed := p.DirtySource(); len(changed) > 0 {
		f.fail("%d uncommitted source file(s) but no changelog entry dated %s", len(changed), today)
	} else {
		f.note("note: no changelog entry dated today (%s); no uncommitted source changes either", today)
	}

	f.checkDates(p, today)
	f.checkAppendOnly(p)
	f.checkRootDocs(p)

	if f.Stale {
		f.note("RESULT: STALE — the documentation trail does not cover the code.")
	} else {
		f.note("RESULT: FRESH")
	}
	return f
}

// checkDates rejects artifacts dated in the future. A date from model memory
// rather than the system clock is the usual cause, and it silently breaks every
// date comparison above.
func (f *Freshness) checkDates(p Project, today string) {
	for _, class := range []string{"changes", "reviews", "testing", "plans"} {
		for _, a := range p.Artifacts(class) {
			if a.Date > today {
				f.fail("future-dated artifact: %s declares %s (today is %s)",
					filepath.Base(a.Path), a.Date, today)
			}
		}
	}
}

// checkAppendOnly enforces that committed changelog entries and ADR text are
// never rewritten. A rewritten record is a stale trail, mechanically.
func (f *Freshness) checkAppendOnly(p Project) {
	out, err := git(p.Root, "diff", "--no-ext-diff", "--no-textconv", "--name-status",
		"HEAD", "--", filepath.Join(p.Docs, "changes"))
	if err == nil && out != "" {
		var touched []string
		for _, line := range strings.Split(out, "\n") {
			fields := strings.Fields(line)
			if len(fields) >= 2 && strings.ContainsAny(fields[0][:1], "MDR") {
				touched = append(touched, fields[len(fields)-1])
			}
		}
		if len(touched) > 0 {
			sort.Strings(touched)
			f.fail("committed changelog entry modified, renamed or deleted (append-only): %s",
				strings.Join(touched, ", "))
		}
	}

	adr := filepath.Join(p.Docs, "architecture", "DECISIONS.md")
	if _, err := os.Stat(adr); err != nil {
		return
	}
	diff, err := git(p.Root, "diff", "--no-ext-diff", "--no-textconv", "HEAD", "--", adr)
	if err != nil {
		return
	}
	for _, line := range strings.Split(diff, "\n") {
		// A removed line starts with one '-'; '---' is the file header, and a
		// removed "- item" bullet shows as "-- item".
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			f.fail("lines removed from DECISIONS.md (append-only: amend by adding, never by rewriting)")
			return
		}
	}
}

// checkRootDocs reports lag in the planning documents. Informational only — a
// stale ROADMAP is worth mentioning, but it is not a broken trail.
func (f *Freshness) checkRootDocs(p Project) {
	headEpoch := p.headEpoch()
	for _, name := range []string{"ROADMAP.md", "TODO.md"} {
		path := filepath.Join(p.Root, name)
		fi, err := os.Stat(path)
		switch {
		case err != nil:
			f.note("note: %s missing at repo root", name)
		case headEpoch > 0 && fi.ModTime().Unix() < headEpoch:
			f.note("note: %s not touched since HEAD commit", name)
		default:
			f.note("OK: %s current relative to HEAD", name)
		}
	}
}

func (p Project) headEpoch() int64 {
	out, err := git(p.Root, "log", "-1", "--format=%ct")
	if err != nil {
		return 0
	}
	var epoch int64
	if _, err := fmt.Sscanf(out, "%d", &epoch); err != nil {
		return 0
	}
	return epoch
}
