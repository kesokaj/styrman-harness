package gate

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kesokaj/styrman-harness/internal/hookio"
	"github.com/kesokaj/styrman-harness/internal/trail"
)

// TrailNotice describes what the project currently owes, as plain context for
// the agent to act on at its own pace.
//
// This is the informational half of the trail gate, and it is deliberately the
// half that runs most often. A gate that blocks costs a full model round-trip
// every time it is wrong — measured at eight forced continuations for a single
// mistaken block. A gate that informs costs a few tokens and is never wrong in
// a way that burns time.
func TrailNotice(dir string) string {
	proj, ok := trail.Resolve(dir)
	if !ok || !proj.Opted() {
		return ""
	}
	cfg, err := trail.LoadConfig(proj.Root)
	if err != nil {
		return "harness: " + err.Error()
	}
	rep := proj.Owed(cfg, trail.Today())
	if len(rep.Due) == 0 {
		return ""
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Definition of Done — %d source file(s) changed, %d line(s). Owed before this work is committed:",
		len(rep.Changed), rep.Lines)
	for _, due := range rep.Due {
		fmt.Fprintf(&b, "\n  - %s: %s", due.Row, due.Reason)
	}
	for _, skip := range rep.Skipped {
		fmt.Fprintf(&b, "\n  - %s: not required — %s", skip.Row, skip.Reason)
	}
	b.WriteString("\nRun `harness dispatch` to satisfy the owed rows, or say why they are not needed.")
	return b.String()
}

// NoticeOnce is TrailNotice, but only when the answer has changed since the
// last time this session heard it. Repeating an unchanged notice on every
// prompt would spend tokens to say nothing new.
func NoticeOnce(p hookio.Payload) string {
	notice := TrailNotice(p.Dir())
	key := safeName(p.SessionID)
	if notice == "" || key == "" {
		return notice
	}
	dir, ok := markerDir("harness-notice")
	if !ok {
		return notice
	}
	marker := filepath.Join(dir, key)
	sum := fmt.Sprintf("%x", sha256.Sum256([]byte(notice)))
	if prev, err := os.ReadFile(marker); err == nil && string(prev) == sum {
		return ""
	}
	_ = os.WriteFile(marker, []byte(sum), 0o600)
	return notice
}

// StopDecision is the Claude Code Stop hook.
//
// It blocks at one boundary only: a commit exists that the changelog trail
// does not cover. Uncommitted work in progress is reported by TrailNotice
// instead — blocking at every turn end is what made the old gate feel like it
// never let go, and the DoD is a commit contract, not a turn contract.
func StopDecision(p hookio.Payload) hookio.Decision {
	if p.StopHookActive {
		return hookio.Allow() // already blocked once for this stop
	}
	proj, ok := trail.Resolve(p.Dir())
	if !ok || !proj.Opted() {
		return hookio.Allow()
	}
	covered, head, newest := proj.CoversHead()
	if covered {
		return hookio.Allow()
	}
	if newest == "" {
		newest = "none"
	}
	return hookio.Block(fmt.Sprintf(
		"Trail behind code: HEAD was committed %s but the newest changelog entry is %s. "+
			"Dispatch the changelog agent for the committed work — pass today's date %s — or say why it is not needed.",
		head, newest, trail.Today()))
}

// AgyStopDecision is the Antigravity Stop hook. Antigravity spells "do not
// stop yet" as continue, sends no stop_hook_active, and caps repeats itself —
// so the loop guard here is a per-conversation marker: nudge once, never more.
func AgyStopDecision(p hookio.Payload) hookio.Decision {
	if !p.Idle() {
		return hookio.Allow() // pending work of its own; never nudge
	}
	proj, ok := trail.Resolve(p.Dir())
	if !ok || !proj.Opted() {
		return hookio.Allow()
	}
	covered, head, newest := proj.CoversHead()
	if covered {
		return hookio.Allow()
	}

	// Without a conversation id there is no way to nudge exactly once, so do
	// not nudge at all rather than risk a loop.
	conv := safeName(p.ConversationID)
	if conv == "" {
		return hookio.Allow()
	}
	dir, ok := markerDir("harness-agy-stop")
	if !ok {
		return hookio.Allow()
	}
	marker := filepath.Join(dir, conv)
	if fileExists(marker) || !touch(marker) {
		return hookio.Allow()
	}

	if newest == "" {
		newest = "none"
	}
	return hookio.Continue(fmt.Sprintf(
		"Trail behind code: HEAD was committed %s but the newest changelog entry is %s. "+
			"Record the committed work — today is %s.", head, newest, trail.Today()))
}
