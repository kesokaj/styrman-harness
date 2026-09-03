// Package hookio handles the wire format both vendors use for hooks: JSON on
// stdin, a JSON decision on stdout, and an exit status that must never break
// the session.
//
// The two vendors spell the same concepts differently. Claude Code blocks by
// emitting {"decision":"block"}; Antigravity's Stop hook re-enters the loop by
// emitting {"decision":"continue"}. Everything else in the harness speaks in
// terms of Decision and lets this package do the translation.
package hookio

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Payload is the union of the fields the harness reads from either vendor's
// hook JSON. Unknown fields are ignored; absent fields stay zero.
type Payload struct {
	// Claude Code
	CWD              string         `json:"cwd"`
	SessionID        string         `json:"session_id"`
	StopHookActive   bool           `json:"stop_hook_active"`
	AgentType        string         `json:"agent_type"`
	AgentID          string         `json:"agent_id"`
	LastAssistantMsg string         `json:"last_assistant_message"`
	HookEventName    string         `json:"hook_event_name"`
	Prompt           string         `json:"prompt"`
	ToolInput        map[string]any `json:"tool_input"`
	ToolName         string         `json:"tool_name"`

	// Antigravity
	WorkspacePaths []string `json:"workspacePaths"`
	ConversationID string   `json:"conversationId"`
	InvocationNum  int      `json:"invocationNum"`
	FullyIdle      *bool    `json:"fullyIdle"`
	ToolCall       struct {
		Args map[string]any `json:"args"`
	} `json:"toolCall"`
}

// Dir is the working directory the hook should act on: the payload's own
// directory when it names one, otherwise the process's.
func (p Payload) Dir() string {
	if p.CWD != "" && isDir(p.CWD) {
		return p.CWD
	}
	if len(p.WorkspacePaths) > 0 && isDir(p.WorkspacePaths[0]) {
		return p.WorkspacePaths[0]
	}
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

// Idle reports whether an Antigravity Stop payload is fully idle. A stop that
// still has pending work of its own is never nudged.
func (p Payload) Idle() bool { return p.FullyIdle == nil || *p.FullyIdle }

func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// Read parses the hook payload from stdin.
//
// Hook stdin is a closed pipe in practice, but a hand-run with an open, silent
// terminal must not hang forever — so a read that produces nothing within the
// cap yields an empty payload rather than blocking. Malformed JSON is not an
// error either: a gate that crashes on an unfamiliar payload is worse than one
// that declines to act.
func Read() Payload {
	var p Payload
	fi, err := os.Stdin.Stat()
	if err != nil || fi.Mode()&os.ModeCharDevice != 0 {
		return p // a terminal, not a pipe: nothing to read
	}

	done := make(chan []byte, 1)
	go func() {
		data, _ := io.ReadAll(os.Stdin)
		done <- data
	}()

	select {
	case data := <-done:
		_ = json.Unmarshal(data, &p)
	case <-time.After(2 * time.Second):
	}
	return p
}

// Decision is a hook's answer. The zero value allows the action.
type Decision struct {
	Verb   string `json:"decision,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// Allow lets the action proceed with no comment.
func Allow() Decision { return Decision{} }

// Block stops a Claude Code action and tells the agent why.
func Block(reason string) Decision { return Decision{Verb: "block", Reason: reason} }

// Continue is Antigravity's inverse spelling of "do not stop yet".
func Continue(reason string) Decision { return Decision{Verb: "continue", Reason: reason} }

// Deny refuses a Claude Code PreToolUse tool call.
func Deny(reason string) Decision { return Decision{Verb: "deny", Reason: reason} }

// Emit writes the decision and exits successfully. A hook must never break the
// session, so the exit status is always zero and blocking travels in the JSON.
func Emit(d Decision) {
	if d.Verb == "" {
		fmt.Println("{}")
		return
	}
	out, err := json.Marshal(d)
	if err != nil {
		fmt.Println("{}")
		return
	}
	fmt.Println(string(out))
}

// EmitContext prints text that the vendor injects into model context.
// SessionStart and UserPromptSubmit stdout is context, not a decision.
func EmitContext(text string) {
	if text != "" {
		fmt.Println(text)
	}
}

// EmitAgyContext wraps context text in Antigravity's PreInvocation shape.
func EmitAgyContext(text string) {
	if text == "" {
		fmt.Println("{}")
		return
	}
	payload := map[string]any{
		"injectSteps": []map[string]any{{"ephemeralMessage": text}},
	}
	out, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("{}")
		return
	}
	fmt.Println(string(out))
}
