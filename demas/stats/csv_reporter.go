package stats

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	interval      time.Duration
}

func NewCSVReporter(store *Store, runInfo RunInfo, interval time.Duration) (*CSVReporter, error) {
	if err := os.MkdirAll(runInfo.OutputDir, 0o755); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	globalPath := filepath.Join(runInfo.OutputDir, "global_stats.csv")
	islandsPath := filepath.Join(runInfo.OutputDir, "island_stats.csv")

	globalFile, err := os.Create(globalPath)
	if err != nil {
		return nil, fmt.Errorf("create global stats file: %w", err)
	}

	islandsFile, err := os.Create(islandsPath)
	if err != nil {
		globalFile.Close()
		return nil, fmt.Errorf("create island stats file: %w", err)
	}

	globalCSV := csv.NewWriter(globalFile)
	islandsCSV := csv.NewWriter(islandsFile)

	if err := writeCSVRow(globalCSV, []string{
		"run_id",
		"timestamp",
		"elapsed_ms",
		"snapshot_index",
		"total_meetings",
		"same_island_meetings",
		"cross_island_meetings",
		"same_island_ratio",
		"best_score",
		"total_agents",
	}); err != nil {
		globalFile.Close()
		islandsFile.Close()
		return nil, fmt.Errorf("write global CSV header: %w", err)
	}

	if err := writeCSVRow(islandsCSV, []string{
		"run_id",
		"timestamp",
		"elapsed_ms",
		"snapshot_index",
		"island_id",
		"agents",
		"meetings",
	}); err != nil {
		globalFile.Close()
		islandsFile.Close()
		return nil, fmt.Errorf("write island CSV header: %w", err)
	}

	return &CSVReporter{
		Store:         store,
		RunInfo:       runInfo,
		globalFile:    globalFile,
		globalCSV:     globalCSV,
		islandsFile:   islandsFile,
		islandsCSV:    islandsCSV,
		startedAt:     time.Now(),
		snapshotIndex: 0,
		interval:      interval,
	}, nil
}

func (r *CSVReporter) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	defer func() {
		if err := r.close(); err != nil {
			fmt.Fprintf(os.Stderr, "csv reporter close error: %v\n", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			if err := r.writeSnapshot(r.Store.Snapshot()); err != nil {
				fmt.Fprintf(os.Stderr, "csv reporter final snapshot error: %v\n", err)
			}
			return

		case <-ticker.C:
			if err := r.writeSnapshot(r.Store.Snapshot()); err != nil {
				fmt.Fprintf(os.Stderr, "csv reporter periodic snapshot error: %v\n", err)
			}
		}
	}
}

func (r *CSVReporter) writeSnapshot(s Snapshot) error {
	now := time.Now()
	timestamp := now.Format(time.RFC3339Nano)
	elapsedMS := now.Sub(r.startedAt).Milliseconds()

	totalAgents := 0
	for _, count := range s.AgentsPerIsland {
		totalAgents += count
	}

	if err := writeCSVRow(r.globalCSV, []string{
		r.RunInfo.RunID,
		timestamp,
		strconv.FormatInt(elapsedMS, 10),
		strconv.Itoa(r.snapshotIndex),
		strconv.Itoa(s.TotalMeetings),
		strconv.Itoa(s.SameIslandMeetings),
		strconv.Itoa(s.CrossIslandMeetings),
		fmt.Sprintf("%.8f", s.SameIslandRatio()),
		fmt.Sprintf("%.8f", s.BestScore),
		strconv.Itoa(totalAgents),
	}); err != nil {
		return fmt.Errorf("write global snapshot row: %w", err)
	}

	for islandID, agents := range s.AgentsPerIsland {
		meetings := s.MeetingsPerIsland[islandID]

		if err := writeCSVRow(r.islandsCSV, []string{
			r.RunInfo.RunID,
			timestamp,
			strconv.FormatInt(elapsedMS, 10),
			strconv.Itoa(r.snapshotIndex),
			strconv.Itoa(islandID),
			strconv.Itoa(agents),
			strconv.Itoa(meetings),
		}); err != nil {
			return fmt.Errorf("write island snapshot row for island %d: %w", islandID, err)
		}
	}

	r.snapshotIndex++
	return nil
}

func (r *CSVReporter) close() error {
	var errs []error

	r.globalCSV.Flush()
	if err := r.globalCSV.Error(); err != nil {
		errs = append(errs, fmt.Errorf("flush global csv: %w", err))
	}

	r.islandsCSV.Flush()
	if err := r.islandsCSV.Error(); err != nil {
		errs = append(errs, fmt.Errorf("flush islands csv: %w", err))
	}

	if err := r.globalFile.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close global csv file: %w", err))
	}

	if err := r.islandsFile.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close islands csv file: %w", err))
	}

	if len(errs) == 0 {
		return nil
	}

	return joinErrors(errs)
}

func writeCSVRow(w *csv.Writer, row []string) error {
	if err := w.Write(row); err != nil {
		return err
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func joinErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}

	var msg strings.Builder
	msg.WriteString(errs[0].Error())
	for i := 1; i < len(errs); i++ {
		msg.WriteString("; " + errs[i].Error())
	}
	return fmt.Errorf("%s", msg.String())
}

