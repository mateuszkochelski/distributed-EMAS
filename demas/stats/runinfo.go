package stats

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type RunInfo struct {
	RunID            string    `json:"run_id"`
	Mode             string    `json:"mode"`
	StartedAt        time.Time `json:"started_at"`
	OutputDir        string    `json:"output_dir"`
	SimulationConfig any       `json:"simulation_config"`
	ProblemConfig    any       `json:"problem_config"`
}

func WriteRunInfo(runInfo RunInfo) error {
	if err := os.MkdirAll(runInfo.OutputDir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(runInfo.OutputDir, "run.json")

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	return enc.Encode(runInfo)
}
