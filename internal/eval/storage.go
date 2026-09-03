package eval

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Storage handles persisting evaluation records locally and to GCP Cloud Storage.
type Storage struct {
	RepoRoot string
	Config   Config
}

// NewStorage creates a storage manager.
func NewStorage(repoRoot string, cfg Config) *Storage {
	return &Storage{
		RepoRoot: repoRoot,
		Config:   cfg,
	}
}

// SaveResult writes the result to the local JSONL ledger and optionally syncs to GCS.
func (s *Storage) SaveResult(ctx context.Context, res EvalRunResult) error {
	// 1. Local append
	ledgerPath := filepath.Join(s.RepoRoot, s.Config.LocalLedger)
	if err := os.MkdirAll(filepath.Dir(ledgerPath), 0o755); err != nil {
		return err
	}

	data, err := json.Marshal(res)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(ledgerPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := fmt.Fprintln(f, string(data)); err != nil {
		return err
	}

	// 2. GCS Sync (if bucket configured)
	if s.Config.GCSBucket != "" {
		_ = s.uploadToGCS(ctx, res)
	}

	return nil
}

func (s *Storage) uploadToGCS(ctx context.Context, res EvalRunResult) error {
	// Retrieve access token
	tokenCmd := exec.CommandContext(ctx, "gcloud", "auth", "application-default", "print-access-token")
	tokenOut, err := tokenCmd.Output()
	if err != nil {
		tokenCmd = exec.CommandContext(ctx, "gcloud", "auth", "print-access-token")
		tokenOut, err = tokenCmd.Output()
		if err != nil {
			return err
		}
	}
	token := strings.TrimSpace(string(tokenOut))

	data, _ := json.MarshalIndent(res, "", "  ")
	objectName := fmt.Sprintf("evals/%s/%s.json", res.RunID, res.TaskID)
	url := fmt.Sprintf("https://storage.googleapis.com/upload/storage/v1/b/%s/o?uploadType=media&name=%s",
		s.Config.GCSBucket, objectName)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ReadLedger returns all results from the local ledger file.
func (s *Storage) ReadLedger() ([]EvalRunResult, error) {
	ledgerPath := filepath.Join(s.RepoRoot, s.Config.LocalLedger)
	data, err := os.ReadFile(ledgerPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var results []EvalRunResult
	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		var r EvalRunResult
		if err := json.Unmarshal([]byte(l), &r); err == nil {
			results = append(results, r)
		}
	}
	return results, nil
}
