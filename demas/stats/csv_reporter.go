package stats

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type CSVReporter struct {
	Store   *Store
	RunInfo RunInfo

	globalFile  *os.File
	globalCSV   *csv.Writer
	islandsFile *os.File
	islandsCSV  *csv.Writer

	startedAt     time.Time
	snapshotIndex int

	interval time.Duration
}

func NewCSVReporter(store *Store, runInfo RunInfo, interval time.Duration) (*CSVReporter, error) {
	if err := os.MkdirAll(runInfo.OutputDir, 0o755); err != nil {
		return nil, err
	}

	globalPath := filepath.Join(runInfo.OutputDir, "global_stats.csv")
	islandsPath := filepath.Join(runInfo.OutputDir, "island_stats.csv")

	globalFile, err := os.Create(globalPath)
	if err != nil {
		return nil, err
	}

	islandsFile, err := os.Create(islandsPath)
	if err != nil {
		globalFile.Close()
		return nil, err
	}

	globalCSV := csv.NewWriter(globalFile)
	islandsCSV := csv.NewWriter(islandsFile)

	// header global
	_ = globalCSV.Write([]string{
		"run_id",
		"timestamp",
		"elapsed_ms",
		"snapshot_index",
		"total_meetings",
		"same_island",
		"cross_island",
		"same_ratio",
		"best_score",
		"total_agents",
	})

	// header islands
	_ = islandsCSV.Write([]string{
		"run_id",
		"timestamp",
		"elapsed_ms",
		"snapshot_index",
		"island_id",
		"agents",
		"meetings",
	})

	globalCSV.Flush()
	islandsCSV.Flush()

	return &CSVReporter{
		Store:         store,
		RunInfo:       runInfo,
		globalFile:    globalFile,
		globalCSV:     globalCSV,
		islandsFile:   islandsFile,
		islandsCSV:    islandsCSV,
		startedAt:     time.Now(),
		interval:      interval,
		snapshotIndex: 0,
	}, nil
}

func (r *CSVReporter) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	defer r.close()

	for {
		select {
		case <-ctx.Done():
			r.writeSnapshot(r.Store.Snapshot())
			return
		case <-ticker.C:
			r.writeSnapshot(r.Store.Snapshot())
		}
	}
}

func (r *CSVReporter) writeSnapshot(s Snapshot) {
	now := time.Now()
	ts := now.Format(time.RFC3339Nano)
	elapsed := now.Sub(r.startedAt).Milliseconds()

	totalAgents := 0
	for _, v := range s.AgentsPerIsland {
		totalAgents += v
	}

	// global row
	_ = r.globalCSV.Write([]string{
		r.RunInfo.RunID,
		ts,
		strconv.FormatInt(elapsed, 10),
		strconv.Itoa(r.snapshotIndex),
		strconv.Itoa(s.TotalMeetings),
		strconv.Itoa(s.SameIslandMeetings),
		strconv.Itoa(s.CrossIslandMeetings),
		fmt.Sprintf("%.6f", s.SameIslandRatio()),
		fmt.Sprintf("%.6f", s.BestScore),
		strconv.Itoa(totalAgents),
	})

	// per island rows
	for id, agents := range s.AgentsPerIsland {
		meetings := s.MeetingsPerIsland[id]

		_ = r.islandsCSV.Write([]string{
			r.RunInfo.RunID,
			ts,
			strconv.FormatInt(elapsed, 10),
			strconv.Itoa(r.snapshotIndex),
			strconv.Itoa(id),
			strconv.Itoa(agents),
			strconv.Itoa(meetings),
		})
	}

	r.globalCSV.Flush()
	r.islandsCSV.Flush()

	r.snapshotIndex++
}

func (r *CSVReporter) close() {
	r.globalCSV.Flush()
	r.islandsCSV.Flush()
	_ = r.globalFile.Close()
	_ = r.islandsFile.Close()
}
