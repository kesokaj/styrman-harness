package eval

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// TelemetryData holds extracted metrics from an agent run.
type TelemetryData struct {
	Turns        int
	ToolCalls    map[string]int
	TotalTokens  int
	InputTokens  int
	OutputTokens int
	CacheTokens  int
	CostUSD      float64
}

// ParseTranscriptJSONL inspects transcript.jsonl (from Antigravity or Claude logs)
// to count turns, tool calls, and tokens.
func ParseTranscriptJSONL(transcriptPath string) (TelemetryData, error) {
	data := TelemetryData{
		ToolCalls: make(map[string]int),
	}

	f, err := os.Open(transcriptPath)
	if err != nil {
		return data, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// Some transcript lines can be large
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var step struct {
			Type      string `json:"type"`
			ToolCalls []struct {
				Function struct {
					Name string `json:"name"`
				} `json:"function"`
				ToolName string `json:"tool_name"`
				Name     string `json:"name"`
			} `json:"tool_calls"`
			Usage struct {
				InputTokens  int `json:"input_tokens"`
				OutputTokens int `json:"output_tokens"`
				TotalTokens  int `json:"total_tokens"`
				CacheRead    int `json:"cache_read_input_tokens"`
				CacheWrite   int `json:"cache_creation_input_tokens"`
			} `json:"usage"`
		}

		if err := json.Unmarshal(line, &step); err != nil {
			continue
		}

		if step.Type == "PLANNER_RESPONSE" || step.Type == "assistant" {
			data.Turns++
		}

		for _, tc := range step.ToolCalls {
			name := tc.Function.Name
			if name == "" {
				name = tc.ToolName
			}
			if name == "" {
				name = tc.Name
			}
			if name != "" {
				data.ToolCalls[name]++
			}
		}

		if step.Usage.TotalTokens > 0 {
			data.InputTokens += step.Usage.InputTokens
			data.OutputTokens += step.Usage.OutputTokens
			data.TotalTokens += step.Usage.TotalTokens
			data.CacheTokens += (step.Usage.CacheRead + step.Usage.CacheWrite)
		}
	}

	return data, scanner.Err()
}

// ParseEnvelopeJSON reads the CLI's standard JSON output envelope.
func ParseEnvelopeJSON(raw []byte) TelemetryData {
	data := TelemetryData{
		ToolCalls: make(map[string]int),
	}

	// 1. Try Antigravity envelope
	var agyEnv struct {
		ConversationID  string  `json:"conversation_id"`
		Status          string  `json:"status"`
		DurationSeconds float64 `json:"duration_seconds"`
		NumTurns        int     `json:"num_turns"`
		Usage           struct {
			InputTokens     int `json:"input_tokens"`
			OutputTokens    int `json:"output_tokens"`
			ThinkingTokens  int `json:"thinking_tokens"`
			CacheReadTokens int `json:"cache_read_tokens"`
			TotalTokens     int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(raw, &agyEnv); err == nil && agyEnv.ConversationID != "" {
		data.Turns = agyEnv.NumTurns
		data.TotalTokens = agyEnv.Usage.TotalTokens
		data.InputTokens = agyEnv.Usage.InputTokens
		data.OutputTokens = agyEnv.Usage.OutputTokens
		data.CacheTokens = agyEnv.Usage.CacheReadTokens

		// Approximate cost (Gemini Flash rates: ~$0.10/1M input, ~$0.40/1M output)
		inputCost := float64(data.InputTokens) * 0.0000001
		outputCost := float64(data.OutputTokens) * 0.0000004
		data.CostUSD = inputCost + outputCost

		// Locate transcript using conversation ID
		if transcriptPath := FindAgyTranscript(agyEnv.ConversationID); transcriptPath != "" {
			if tData, err := ParseTranscriptJSONL(transcriptPath); err == nil {
				if len(tData.ToolCalls) > 0 {
					data.ToolCalls = tData.ToolCalls
				}
				if tData.Turns > data.Turns {
					data.Turns = tData.Turns
				}
			}
		}
		return data
	}

	// 2. Try Claude Code envelope
	var env struct {
		CostUSD      float64        `json:"total_cost_usd"`
		TotalTokens  int            `json:"total_tokens"`
		InputTokens  int            `json:"input_tokens"`
		OutputTokens int            `json:"output_tokens"`
		Turns        int            `json:"turns"`
		ToolCalls    map[string]int `json:"tool_calls"`
	}

	if err := json.Unmarshal(raw, &env); err == nil {
		data.CostUSD = env.CostUSD
		data.TotalTokens = env.TotalTokens
		data.InputTokens = env.InputTokens
		data.OutputTokens = env.OutputTokens
		data.Turns = env.Turns
		if env.ToolCalls != nil {
			data.ToolCalls = env.ToolCalls
		}
	}

	return data
}

// FindAgyTranscript locates the transcript file for an Antigravity conversation.
func FindAgyTranscript(convID string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	candidates := []string{
		filepath.Join(home, ".gemini", "antigravity-cli", "brain", convID, ".system_generated", "logs", "transcript.jsonl"),
		filepath.Join(home, ".gemini", "antigravity", "brain", convID, ".system_generated", "logs", "transcript.jsonl"),
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			return c
		}
	}
	return ""
}

// FindTranscriptInWorkspace searches for transcript or log files in the workspace.
func FindTranscriptInWorkspace(workspaceDir string) string {
	candidates := []string{
		filepath.Join(workspaceDir, ".agents", "logs", "transcript.jsonl"),
		filepath.Join(workspaceDir, "transcript.jsonl"),
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			return c
		}
	}

	// Search recursively in .system_generated
	var found string
	_ = filepath.Walk(workspaceDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, "transcript.jsonl") {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	return found
}
