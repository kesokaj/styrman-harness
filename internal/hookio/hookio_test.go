package hookio

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestAgyContextJSON covers the Antigravity wrapping, which until now was the
// one vendor path with no test at all: the plain Claude Code path is exercised
// constantly, so a break in the agy shape would have surfaced only on the
// operator's machine.
func TestAgyContextJSON(t *testing.T) {
	t.Run("empty injects nothing", func(t *testing.T) {
		if got := AgyContextJSON(""); got != "{}" {
			t.Errorf("got %q, want {}", got)
		}
	})

	t.Run("text is wrapped in injectSteps", func(t *testing.T) {
		out := AgyContextJSON("Resume Protocol: line one\nline two")

		var payload struct {
			InjectSteps []struct {
				EphemeralMessage string `json:"ephemeralMessage"`
			} `json:"injectSteps"`
		}
		if err := json.Unmarshal([]byte(out), &payload); err != nil {
			t.Fatalf("not valid JSON: %v — %s", err, out)
		}
		if len(payload.InjectSteps) != 1 {
			t.Fatalf("got %d inject steps, want 1", len(payload.InjectSteps))
		}
		if !strings.Contains(payload.InjectSteps[0].EphemeralMessage, "line two") {
			t.Error("multi-line text was truncated or mangled")
		}
	})

	t.Run("control characters stay escaped", func(t *testing.T) {
		out := AgyContextJSON("has \"quotes\" and\ttabs")
		if !json.Valid([]byte(out)) {
			t.Errorf("payload is not valid JSON: %s", out)
		}
	})
}

// TestDecisionShapes pins the two vendors' opposite spellings. Claude Code
// blocks with "block"; Antigravity re-enters its loop with "continue". Swapping
// them would make each gate do the opposite of what it intends.
func TestDecisionShapes(t *testing.T) {
	cases := map[string]struct {
		d        Decision
		wantVerb string
	}{
		"allow":    {Allow(), ""},
		"block":    {Block("why"), "block"},
		"continue": {Continue("why"), "continue"},
		"deny":     {Deny("why"), "deny"},
	}
	for name, tc := range cases {
		if tc.d.Verb != tc.wantVerb {
			t.Errorf("%s: verb = %q, want %q", name, tc.d.Verb, tc.wantVerb)
		}
	}
}

func TestPayloadDirPrefersItsOwnCWD(t *testing.T) {
	dir := t.TempDir()

	if got := (Payload{CWD: dir}).Dir(); got != dir {
		t.Errorf("Claude cwd ignored: got %q", got)
	}
	if got := (Payload{WorkspacePaths: []string{dir}}).Dir(); got != dir {
		t.Errorf("Antigravity workspacePaths ignored: got %q", got)
	}
	// A path that does not exist must fall back rather than be trusted.
	if got := (Payload{CWD: "/nonexistent/xyzzy"}).Dir(); got == "/nonexistent/xyzzy" {
		t.Error("a nonexistent cwd was used as-is")
	}
}

func TestPayloadIdleDefaultsTrue(t *testing.T) {
	if !(Payload{}).Idle() {
		t.Error("a payload without fullyIdle must be treated as idle")
	}
	no := false
	if (Payload{FullyIdle: &no}).Idle() {
		t.Error("fullyIdle:false was read as idle — jq's // treats false like null, and this must not")
	}
}
