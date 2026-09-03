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

// commitBoundary reports whether this is a commit the trail has not recorded,
// and the scope an entry must declare to record it.
//
// The gate no longer computes a scope of its own. It used to: coverage was
// decided with one function and the remediation message printed another, so on
// a dirty tree the block named a scope the gate would never accept — doing
// exactly what it said left it blocking. Two functions that must agree will
// eventually not. Asking Owed, which is also what `harness dispatch` satisfies,
// removes the possibility rather than re-aligning it.
func commitBoundary(dir string) (blocked bool, scope, head, newest string) {
	proj, ok := trail.Resolve(dir)
	if !ok || !proj.Opted() {
		return false, "", "", ""
	}
	// Uncommitted work is not this gate's business — TrailNotice reports it.
	// The Definition of Done is a commit contract, not a turn contract.
	if len(proj.DirtySource()) > 0 {
		return false, "", "", ""
	}
	cfg, err := trail.LoadConfig(proj.Root)
	if err != nil {
		return false, "", "", ""
	}
	rep := proj.Owed(cfg, trail.Today())
	for _, due := range rep.Due {
		if due.Row == trail.RowChangelog {
			return true, rep.Scope, proj.HeadDate(), proj.NewestChangeDate()
		}
	}
	return false, "", "", ""
}

// StopDecision is the Claude Code Stop hook.
//
// It blocks at one boundary only: a commit the changelog does not record.
// Uncommitted work in progress is reported by TrailNotice instead — blocking at
// every turn end is what made the old gate feel like it never let go.
func StopDecision(p hookio.Payload) hookio.Decision {
	if p.StopHookActive {
		return hookio.Allow() // already blocked once for this stop
	}
	blocked, scope, head, newest := commitBoundary(p.Dir())
	if !blocked {
		return hookio.Allow()
	}
	if newest == "" {
		newest = "none"
	}
	return hookio.Block(fmt.Sprintf(
		"Trail behind code: HEAD was committed %s and no changelog entry covers it "+
			"(newest entry: %s). Run `harness dispatch`, or write the entry yourself with "+
			"frontmatter `owner: changelog`, `date: %s` and `scope: %s` — an entry without "+
			"that exact scope covers nothing. Or say why it is not needed.",
		head, newest, trail.Today(), scope))
}

// AgyStopDecision is the Antigravity Stop hook. Antigravity spells "do not
// stop yet" as continue, sends no stop_hook_active, and caps repeats itself —
// so the loop guard here is a per-conversation marker: nudge once, never more.
func AgyStopDecision(p hookio.Payload) hookio.Decision {
	if !p.Idle() {
		return hookio.Allow() // pending work of its own; never nudge
	}
	blocked, scope, head, newest := commitBoundary(p.Dir())
	if !blocked {
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
		"Trail behind code: HEAD was committed %s and no changelog entry covers it "+
			"(newest entry: %s). Record the committed work — today is %s, scope %s.",
		head, newest, trail.Today(), scope))
}
