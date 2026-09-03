package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kesokaj/styrman-harness/internal/hookio"
	"github.com/kesokaj/styrman-harness/internal/trail"
)

// ownerRows maps a subagent's type to the DoD row it owns. Agents outside this
// map own no dated artifact and are never held.
var ownerRows = map[string]trail.Row{
	"reviewer":  trail.RowReviewer,
	"security":  trail.RowSecurity,
	"verifier":  trail.RowVerifier,
	"changelog": trail.RowChangelog,
	"explorer":  trail.RowExplorer,
}

// maxHolds caps how many times one subagent may be held. Claude Code stops a
// looping hook on its own, but only after eight forced continuations — enough
// to burn a dollar and degrade the agent into refusing outright. Holding once
// is a reminder; holding repeatedly is just damage.
const maxHolds = 1

// SubagentStopDecision holds a Definition-of-Done owner that stopped without
// writing its artifact, so the artifact gets written rather than reported in
// chat and lost.
//
// Every guard here fails open. The loop protection rests on stop_hook_active,
// which both vendors send but neither documents; the hold counter is the
// second guard, so an unfamiliar payload costs one wasted turn rather than
// eight.
func SubagentStopDecision(p hookio.Payload) hookio.Decision {
	row, owned := ownerRows[p.AgentType]
	if !owned {
		return hookio.Allow()
	}

	dir, markersOK := markerDir("harness-subagent")
	marker := ""
	if markersOK {
		marker = filepath.Join(dir, safeName(p.AgentID)+"."+p.AgentType)
	}
	release := func() hookio.Decision {
		if marker != "" {
			_ = os.Remove(marker)
		}
		return hookio.Allow()
	}

	if p.StopHookActive {
		return release()
	}
	proj, ok := trail.Resolve(p.Dir())
	if !ok || !proj.Opted() {
		return release()
	}
	if hasFreshArtifact(proj, row, marker) {
		return release()
	}
	// Second guard: never hold the same agent more than maxHolds times, even
	// if stop_hook_active never arrives.
	if marker == "" || bump(marker+".holds") > maxHolds {
		return release()
	}

	return hookio.Block(fmt.Sprintf(
		"Definition of Done: the %s agent stopped without writing its artifact. Write %s/%s/%s-<slug>.md "+
			"with frontmatter `owner: %s` and `date: %s`, then name its path in your final message and stop. "+
			"A refusal, a no-op, or an empty result is still a report with a verdict.",
		p.AgentType, proj.DocsRel(), row.Class(), trail.Today(), row, trail.Today()))
}

// SubagentStartMark stamps the start of a subagent run, so its artifact can be
// told apart from one a sibling wrote earlier in the same day.
func SubagentStartMark(p hookio.Payload) {
	if _, owned := ownerRows[p.AgentType]; !owned {
		return
	}
	dir, ok := markerDir("harness-subagent")
	if !ok {
		return
	}
	base := filepath.Join(dir, safeName(p.AgentID)+"."+p.AgentType)
	touch(base)
	_ = os.Remove(base + ".holds")
}

// hasFreshArtifact reports whether the row's artifact for today exists and was
// written during this run. Without a start marker — a resumed agent never gets
// SubagentStart — anything from the last hour counts.
func hasFreshArtifact(proj trail.Project, row trail.Row, marker string) bool {
	since := time.Now().Add(-time.Hour)
	if marker != "" {
		if fi, err := os.Stat(marker); err == nil {
			since = fi.ModTime()
		}
	}
	today := trail.Today()
	for _, a := range proj.Artifacts(row.Class()) {
		if a.Owner != row || a.Date != today {
			continue
		}
		fi, err := os.Stat(a.Path)
		if err != nil {
			continue
		}
		// Not-before rather than after: an artifact written inside the
		// marker's own second still counts.
		if !fi.ModTime().Before(since) {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
