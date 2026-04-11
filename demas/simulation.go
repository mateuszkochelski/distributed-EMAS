package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"time"

	continuous "github.com/mateuszkochelski/tsp-emas/continuous"
	"github.com/mateuszkochelski/tsp-emas/stats"
	tsp "github.com/mateuszkochelski/tsp-emas/tsp"
)

func randomID(nBytes int) string {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

func buildRunInfo(config *SimulationConfig, problemName string, problemConfig any) stats.RunInfo {
	now := time.Now()
	runID := fmt.Sprintf(
		"%s-%s-%s",
		problemName,
		now.Format("20060102-150405"),
		randomID(6),
	)

	return stats.RunInfo{
		RunID:            runID,
		Mode:             problemName,
		StartedAt:        now,
		OutputDir:        filepath.Join("runs", runID),
		SimulationConfig: *config,
		ProblemConfig:    problemConfig,
	}
}

func setupStats(
	ctx context.Context,
	agentsPerIsland map[int]int,
	config *SimulationConfig,
	problemName string,
	problemConfig any,
	eventCh chan stats.Event,
	maximum float64,
	isBetter func(newScore, currentBest float64) bool,
) error {
	store := stats.NewStore(agentsPerIsland, maximum, isBetter)
	collector := stats.NewCollector(eventCh, store)
	consoleReporter := stats.NewConsoleReporter(store, 1_000_000, 1*time.Second)

	runInfo := buildRunInfo(config, problemName, problemConfig)

	if err := stats.WriteRunInfo(runInfo); err != nil {
		return fmt.Errorf("failed to create run info file: %w", err)
	}

	csvReporter, err := stats.NewCSVReporter(store, runInfo, 1*time.Second)
	if err != nil {
		return fmt.Errorf("failed to create CSV reporter: %w", err)
	}

	go collector.Run(ctx)
	go consoleReporter.Run(ctx)
	go csvReporter.Run(ctx)

	return nil
}

func runSimulation[S any](
	ctx context.Context,
	problem Problem[S],
	config *SimulationConfig,
	problemName string,
	problemConfig any,
) {
	eventCh := make(chan stats.Event, 1_000_000)
	islands := make([]chan Message, config.NumIslands)

	bufferSize := maxInt(int(math.Sqrt(float64(config.NumAgents/config.NumIslands))), 1)

	for i := range islands {
		islands[i] = make(chan Message, bufferSize)
	}

	agentsPerIsland := make(map[int]int)
	agents := make([]Agent[S], config.NumAgents)

	for i := range agents {
		agents[i] = NewAgent(strconv.Itoa(i), islands, eventCh, problem, config)
		agentsPerIsland[agents[i].PrimaryIsland]++
	}

	if err := setupStats(ctx, agentsPerIsland, config, problemName, problemConfig, eventCh, problem.Maximum(), problem.IsBetter); err != nil {
		fmt.Println(err)
		return
	}

	for i := range agents {
		go agents[i].Run(ctx)
	}
}

func runTSP(ctx context.Context, simCfg *SimulationConfig, tspCfg *tsp.Config) error {
	if err := simCfg.Validate(); err != nil {
		return fmt.Errorf("invalid simulation config: %w", err)
	}
	if err := tspCfg.Validate(); err != nil {
		return fmt.Errorf("invalid TSP config: %w", err)
	}

	var problem tsp.Problem
	var err error

	switch tspCfg.Source {
	case tsp.SourceRandom:
		problem = tsp.NewRandomProblem(tspCfg.NumCities, tspCfg.MinCoord, tspCfg.MaxCoord)
	case tsp.SourceFile:
		problem, err = tsp.LoadFromFile(tspCfg.InputFile)
		if err != nil {
			return fmt.Errorf("failed to load TSP problem: %w", err)
		}
	default:
		return fmt.Errorf("unknown tsp source: %s", tspCfg.Source)
	}

	runSimulation(ctx, problem, simCfg, "tsp", *tspCfg)
	return nil
}

func runContinuous(ctx context.Context, simCfg *SimulationConfig, continuousCfg *continuous.Config) error {
	if err := simCfg.Validate(); err != nil {
		return fmt.Errorf("invalid simulation config: %w", err)
	}
	if err := continuousCfg.Validate(); err != nil {
		return fmt.Errorf("invalid continuous config: %w", err)
	}

	problem := continuous.NewProblem(*continuousCfg)
	runSimulation(ctx, problem, simCfg, "continuous", *continuousCfg)
	return nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
