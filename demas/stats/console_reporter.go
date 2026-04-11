package stats

import (
	"context"
	"fmt"
	"time"
)

type ConsoleReporter struct {
	Store              *Store
	PrintEveryMeetings int
	PrintInterval      time.Duration
}

func NewConsoleReporter(store *Store, printEveryMeetings int, printInterval time.Duration) *ConsoleReporter {
	return &ConsoleReporter{
		Store:              store,
		PrintEveryMeetings: printEveryMeetings,
		PrintInterval:      printInterval,
	}
}

func (r *ConsoleReporter) Run(ctx context.Context) {
	ticker := time.NewTicker(r.PrintInterval)
	defer ticker.Stop()

	lastPrintedMeetings := 0
	lastBestScore := r.Store.BestScore

	for {
		select {
		case <-ctx.Done():
			snapshot := r.Store.Snapshot()
			r.printFinal(snapshot)
			return

		case <-ticker.C:
			snapshot := r.Store.Snapshot()

			if snapshot.BestScore < lastBestScore {
				lastBestScore = snapshot.BestScore
				fmt.Println("new minimum:", snapshot.BestScore)
			}

			if r.PrintEveryMeetings > 0 && snapshot.TotalMeetings-lastPrintedMeetings >= r.PrintEveryMeetings {
				lastPrintedMeetings = snapshot.TotalMeetings
				r.printPeriodic(snapshot)
			}
		}
	}
}

func (r *ConsoleReporter) printPeriodic(snapshot Snapshot) {
	fmt.Printf(
		"meetings=%d same=%d cross=%d ratio=%.6f best=%.6f\n",
		snapshot.TotalMeetings,
		snapshot.SameIslandMeetings,
		snapshot.CrossIslandMeetings,
		snapshot.SameIslandRatio(),
		snapshot.BestScore,
	)

	// for islandID, count := range snapshot.AgentsPerIsland {
	// 	fmt.Printf("%d -> %d\n", islandID, count)
	// }
}

func (r *ConsoleReporter) printFinal(snapshot Snapshot) {
	fmt.Println("final meetings per island:")
	for islandID, count := range snapshot.MeetingsPerIsland {
		fmt.Printf("%d -> %d\n", islandID, count)
	}
}
