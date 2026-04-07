package main

import (
	"context"
	"fmt"
	"math"
)

type Log struct {
	From            string
	Msg             string
	Score           float64
	SameIsland      bool
	PrimaryIslandId int
	EventType       string
}

type LogCollector struct {
	AgentsPerIsland map[int]int
	LogCh           chan Log
}

func (lc *LogCollector) Run(ctx context.Context) {
	sameIslandLogs := 0
	differentIslandLogs := 0
	minimum_distance := math.Inf(1)
	counts := make(map[int]int)
	for {
		select {
		case <-ctx.Done():
			for k, v := range counts {
				fmt.Printf("%d -> %d\n", k, v)
			}
			return

		case log := <-lc.LogCh:
			switch log.EventType {
			case "BORN":
				lc.AgentsPerIsland[log.PrimaryIslandId]++
			case "DEAD":
				lc.AgentsPerIsland[log.PrimaryIslandId]--
			case "MEETING":
				counts[log.PrimaryIslandId]++
				if log.Score < minimum_distance {
					minimum_distance = log.Score
					fmt.Println("new minimum:", minimum_distance, "discovered on island: ", log.PrimaryIslandId)
				}
				if log.SameIsland {
					sameIslandLogs++
				} else {
					differentIslandLogs++
				}
				if (sameIslandLogs+differentIslandLogs)%1_000_000 == 0 {
					// fmt.Println("Meetings: ", (sameIslandLogs+differentIslandLogs)/1_000_000, " mln")
					// fmt.Println("SAME ISLANDS: ", sameIslandLogs, " DIFFERENT ISLANDS: ", differentIslandLogs)
					// fmt.Println("PERCENTAGE: ", float64(sameIslandLogs)/float64(sameIslandLogs+differentIslandLogs))
					//
					for k, v := range lc.AgentsPerIsland {
						fmt.Printf("%d -> %d\n", k, v)
					}

				}

			}
		}
	}
}
