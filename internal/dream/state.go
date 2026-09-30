package dream

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// state is ~/.config/sh-agent/dream-state.json: each project's watermark and
// the last run (D9; never in search.db, which is a cache).
type state struct {
	Projects map[string]time.Time `json:"projects"`
	LastRun  lastRun              `json:"last_run"`
}

type lastRun struct {
	At        time.Time `json:"at"`
	Projects  int       `json:"projects"`
	Proposals int       `json:"proposals"`
	Error     string    `json:"error,omitempty"`
}

func statePath(home string) string {
	return filepath.Join(home, ".config", "sh-agent", "dream-state.json")
}

// readState reads the state; a missing or unreadable file is a first run.
func readState(home string) state {
	st := state{Projects: map[string]time.Time{}}
	if data, err := os.ReadFile(statePath(home)); err == nil {
		json.Unmarshal(data, &st)
	}
	if st.Projects == nil {
		st.Projects = map[string]time.Time{}
	}
	return st
}
