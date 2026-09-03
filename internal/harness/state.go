package harness

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// State is the install manifest. It is the only thing uninstall trusts: a
// path not listed here is not styrman's to delete, whatever directory it
// sits in.
type State struct {
	Repo        string   `json:"repo"`
	Commit      string   `json:"commit"`
	InstalledAt string   `json:"installedAt"`
	Files       []string `json:"files"` // relative to $HOME, sorted
}

func statePath(home string) string {
	return filepath.Join(home, ".config", "styrman", "state.json")
}

// LoadState reads the manifest. A missing file is not an error: it means
// nothing has been installed into this home yet.
func LoadState(home string) (State, bool, error) {
	data, err := os.ReadFile(statePath(home))
	if errors.Is(err, os.ErrNotExist) {
		return State{}, false, nil
	}
	if err != nil {
		return State{}, false, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, false, fmt.Errorf("parsing %s: %w", statePath(home), err)
	}
	return s, true, nil
}

func SaveState(home string, s State) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return writeAtomic(statePath(home), append(data, '\n'), 0o644)
}

func RemoveState(home string) error {
	p := statePath(home)
	if err := os.Remove(p); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	// The directory is styrman's own; drop it when empty.
	_ = os.Remove(filepath.Dir(p))
	return nil
}
