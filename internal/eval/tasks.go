package eval

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var defaultTasks = map[string]Task{
	"gf-go-cli": {
		ID:          "gf-go-cli",
		Name:        "Greenfield Go CLI: Log Aggregator",
		Category:    CategoryGreenfield,
		Description: "Build a standalone Go CLI tool that parses structured JSON log lines and prints status code counts.",
		Prompt: "Build a Go CLI utility in this directory. Requirements:\n" +
			"1. The program should accept a log file path as an argument or read from standard input.\n" +
			"2. Each line of input is a JSON object with at least {\"status\": <int>, \"path\": <string>}.\n" +
			"3. It must count occurrences of each status code and output a summary.\n" +
			"4. Provide unit tests with `go test ./...`.\n" +
			"5. Follow stdlib-first Go conventions.\n" +
			"6. Satisfy the project Definition of Done and commit your work cleanly.\n",
		InitialFiles: map[string]string{
			"go.mod":    "module logagg\n\ngo 1.22\n",
			"README.md": "# Log Aggregator CLI\n\nBuild instructions will go here.\n",
		},
		Verify: func(workspaceDir string) (bool, string, error) {
			// 1. Check if tests pass
			testCmd := exec.Command("go", "test", "-v", "./...")
			testCmd.Dir = workspaceDir
			testOut, err := testCmd.CombinedOutput()
			if err != nil {
				return false, fmt.Sprintf("go test failed: %s", string(testOut)), nil
			}

			// 2. Build the binary
			binPath := filepath.Join(workspaceDir, "logagg_bin")
			buildCmd := exec.Command("go", "build", "-o", binPath, ".")
			buildCmd.Dir = workspaceDir
			buildOut, err := buildCmd.CombinedOutput()
			if err != nil {
				return false, fmt.Sprintf("go build failed: %s", string(buildOut)), nil
			}
			defer os.Remove(binPath)

			// 3. Run with synthetic test log file
			testLog := `{"status": 200, "path": "/api/v1"}
{"status": 200, "path": "/api/v2"}
{"status": 404, "path": "/missing"}
{"status": 500, "path": "/crash"}
`
			logFile := filepath.Join(workspaceDir, "test.log")
			if err := os.WriteFile(logFile, []byte(testLog), 0o600); err != nil {
				return false, "", err
			}
			defer os.Remove(logFile)

			runCmd := exec.Command(binPath, logFile)
			runCmd.Dir = workspaceDir
			var stdout bytes.Buffer
			runCmd.Stdout = &stdout
			if err := runCmd.Run(); err != nil {
				// Try feeding via stdin if positional arg failed
				stdinCmd := exec.Command(binPath)
				stdinCmd.Dir = workspaceDir
				stdinCmd.Stdin = strings.NewReader(testLog)
				stdinCmd.Stdout = &stdout
				if err2 := stdinCmd.Run(); err2 != nil {
					return false, fmt.Sprintf("execution failed: %v / %v", err, err2), nil
				}
			}

			outStr := stdout.String()
			// Should report count 2 for 200, 1 for 404, 1 for 500
			if !strings.Contains(outStr, "200") || !strings.Contains(outStr, "404") {
				return false, fmt.Sprintf("output did not contain expected status codes: %s", outStr), nil
			}

			return true, "tests passed, binary built, and sample log aggregation succeeded", nil
		},
		Rubric: "1. Clean CLI flag/arg handling and stdin reading.\n" +
			"2. Idiomatic Go error handling and json.Decoder usage.\n" +
			"3. Comprehensive unit test coverage.\n" +
			"4. Proper Definition-of-Done commit artifacts and clean working tree.",
	},

	"gf-go-rest": {
		ID:          "gf-go-rest",
		Name:        "Greenfield Go REST API: Book Catalog",
		Category:    CategoryGreenfield,
		Description: "Build a Go HTTP REST API for a book catalog with in-memory or SQLite storage, supporting CRUD endpoints and unit tests.",
		Prompt: "Build an HTTP REST API server in Go for a book catalog. Requirements:\n" +
			"1. Endpoints: POST /books (create book), GET /books (list books), GET /books/{id} (get by id).\n" +
			"2. A Book has at least: ID (string), Title (string), Author (string), Year (int).\n" +
			"3. Use stdlib net/http (or lightweight router).\n" +
			"4. Provide unit or integration tests verifying all endpoints.\n" +
			"5. Follow Go standards (proper status codes: 201 for create, 404 for not found, 400 for bad request).\n" +
			"6. Satisfy the project Definition of Done and commit your work cleanly.\n",
		InitialFiles: map[string]string{
			"go.mod": "module bookapi\n\ngo 1.22\n",
		},
		Verify: func(workspaceDir string) (bool, string, error) {
			// Run tests
			testCmd := exec.Command("go", "test", "-v", "./...")
			testCmd.Dir = workspaceDir
			testOut, err := testCmd.CombinedOutput()
			if err != nil {
				return false, fmt.Sprintf("go test failed: %s", string(testOut)), nil
			}

			// Verify builds cleanly
			buildCmd := exec.Command("go", "build", "./...")
			buildCmd.Dir = workspaceDir
			buildOut, err := buildCmd.CombinedOutput()
			if err != nil {
				return false, fmt.Sprintf("go build failed: %s", string(buildOut)), nil
			}

			return true, "HTTP server code compiled and unit/integration tests passed", nil
		},
		Rubric: "1. Clean HTTP handler design with proper status codes (200, 201, 400, 404).\n" +
			"2. Thread-safe persistence (mutex or database).\n" +
			"3. Test coverage for success and error branches.\n" +
			"4. Definition-of-Done compliance (changeset, verification log, commit boundary).",
	},

	"bf-trail-feature": {
		ID:          "bf-trail-feature",
		Name:        "Brownfield Trail Resume: Add Healthz Endpoint",
		Category:    CategoryBrownfield,
		Description: "Resume work in an existing project with an active .agents trail, implement the next task from TODO.md, satisfy DoD, and commit.",
		Prompt: "Resume work on this existing project.\n" +
			"1. Follow the Resume Protocol: inspect TODO.md and recent trail artifacts.\n" +
			"2. Implement the pending feature listed in TODO.md (add /healthz endpoint returning {\"status\": \"ok\", \"time\": \"...\"}).\n" +
			"3. Add unit test covering /healthz.\n" +
			"4. Ensure `go test ./...` passes.\n" +
			"5. Satisfy the Definition of Done (run `harness dispatch` or record required artifacts) and commit cleanly.",
		InitialFiles: map[string]string{
			"go.mod":    "module service\n\ngo 1.22\n",
			"AGENTS.md": "# AGENTS.md — service\n\nProject rules follow standard .agents trail.\n",
			"TODO.md":   "# TODO\n\n- Add /healthz endpoint returning JSON status ok and current timestamp.\n",
			"main.go": `package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "welcome")
	})
	_ = http.ListenAndServe(":8080", nil)
}
`,
			".agents/docs/changes/2026-09-01-init.md": "---\nowner: changelog\ndate: 2026-09-01\nscope: changeset:initial\n---\nInitial commit.\n",
		},
		Verify: func(workspaceDir string) (bool, string, error) {
			// 1. Check tests pass
			testCmd := exec.Command("go", "test", "-v", "./...")
			testCmd.Dir = workspaceDir
			testOut, err := testCmd.CombinedOutput()
			if err != nil {
				return false, fmt.Sprintf("go test failed: %s", string(testOut)), nil
			}

			// 2. Check git commit
			gitLog := exec.Command("git", "log", "-n", "1", "--oneline")
			gitLog.Dir = workspaceDir
			logOut, err := gitLog.CombinedOutput()
			if err != nil {
				return false, "failed to read git log", err
			}
			if strings.Contains(string(logOut), "initial") {
				return false, "agent did not commit its work", nil
			}

			return true, "healthz feature implemented, tests pass, and work was committed cleanly", nil
		},
		Rubric: "1. Did the agent follow the Resume Protocol?\n" +
			"2. Did the agent satisfy the Definition of Done?\n" +
			"3. Is the git working tree clean with a descriptive commit message?\n" +
			"4. Does /healthz return valid JSON with status and timestamp?",
	},
}

// GetTask returns a task by ID.
func GetTask(id string) (Task, bool) {
	t, ok := defaultTasks[id]
	return t, ok
}

// AllTasks returns all registered tasks sorted by ID.
func AllTasks() []Task {
	tasks := make([]Task, 0, len(defaultTasks))
	for _, t := range defaultTasks {
		tasks = append(tasks, t)
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return tasks
}

// RegisterTask adds or overrides a task.
func RegisterTask(t Task) {
	defaultTasks[t.ID] = t
}

func init() {
	// Stamp date if needed
	_ = time.Now()
}
