package eval

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds evaluation system settings and GCP integration parameters.
type Config struct {
	ProjectID   string `json:"projectId"`
	Region      string `json:"region"`
	GCSBucket   string `json:"gcsBucket"`
	VertexModel string `json:"vertexModel"`
	LocalLedger string `json:"localLedger"`
}

// DefaultConfig returns baseline configuration defaults.
func DefaultConfig() Config {
	return Config{
		ProjectID:   os.Getenv("GCP_PROJECT"),
		Region:      firstNonEmpty(os.Getenv("GCP_REGION"), "europe-west1"),
		GCSBucket:   os.Getenv("GCS_BUCKET"),
		VertexModel: firstNonEmpty(os.Getenv("VERTEX_MODEL"), "gemini-1.5-pro"),
		LocalLedger: ".agents/docs/evals/eval-ledger.jsonl",
	}
}

// LoadConfig loads configuration from the repo root or returns defaults.
func LoadConfig(repoRoot string) Config {
	cfg := DefaultConfig()
	path := filepath.Join(repoRoot, ".agents", "eval-config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	var loaded Config
	if err := json.Unmarshal(data, &loaded); err == nil {
		if loaded.ProjectID != "" {
			cfg.ProjectID = loaded.ProjectID
		}
		if loaded.Region != "" {
			cfg.Region = loaded.Region
		}
		if loaded.GCSBucket != "" {
			cfg.GCSBucket = loaded.GCSBucket
		}
		if loaded.VertexModel != "" {
			cfg.VertexModel = loaded.VertexModel
		}
		if loaded.LocalLedger != "" {
			cfg.LocalLedger = loaded.LocalLedger
		}
	}
	return cfg
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
