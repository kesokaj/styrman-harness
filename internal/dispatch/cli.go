package dispatch

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// cliTool is one vendor's headless entry point. Both CLIs can run a named
// agent non-interactively and report structured usage, which is what makes
// deterministic dispatch possible without a model in the loop.
type cliTool struct {
	name string
	args func(agent, prompt string) []string
}

// runOutput is the part of a headless run this package cares about.
type runOutput struct {
	Result  string
	CostUSD float64
	IsError bool
}

// resolveCLI picks the CLI to dispatch with.
func resolveCLI(preferred string) (cliTool, error) {
	claude := cliTool{
		name: "claude",
		args: func(agent, prompt string) []string {
			return []string{"-p", prompt, "--agent", agent, "--output-format", "json"}
		},
	}
	agy := cliTool{
		name: "agy",
		args: func(agent, prompt string) []string {
			return []string{"-p", prompt, "--agent", agent, "--output-format", "json"}
		},
	}

	switch preferred {
	case "claude":
		return claude, requireBinary("claude")
	case "agy":
		return agy, requireBinary("agy")
	case "":
		if requireBinary("claude") == nil {
			return claude, nil
		}
		if requireBinary("agy") == nil {
			return agy, nil
		}
		return cliTool{}, fmt.Errorf("neither claude nor agy is on PATH")
	default:
		return cliTool{}, fmt.Errorf("unknown vendor %q (expected claude or agy)", preferred)
	}
}

func requireBinary(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("%s not on PATH", name)
	}
	return nil
}

// run executes one agent headlessly and parses what it cost.
func (c cliTool) run(ctx context.Context, dir, agent, prompt string) (runOutput, error) {
	cmd := exec.CommandContext(ctx, c.name, c.args(agent, prompt)...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if ctx.Err() != nil {
			return runOutput{}, ctx.Err()
		}
		return runOutput{}, fmt.Errorf("%s: %w", c.name, err)
	}
	return parseRun(out), nil
}

// parseRun reads the headless result envelope. Field names differ between the
// vendors and may gain fields, so anything missing is simply zero rather than
// an error — a dispatch must not fail because a cost field was renamed.
func parseRun(data []byte) runOutput {
	var env struct {
		Result   string  `json:"result"`
		IsError  bool    `json:"is_error"`
		CostUSD  float64 `json:"total_cost_usd"`
		Response string  `json:"response"`
	}
	if err := json.Unmarshal(data, &env); err != nil {
		return runOutput{Result: strings.TrimSpace(string(data))}
	}
	result := env.Result
	if result == "" {
		result = env.Response
	}
	return runOutput{Result: strings.TrimSpace(result), CostUSD: env.CostUSD, IsError: env.IsError}
}
