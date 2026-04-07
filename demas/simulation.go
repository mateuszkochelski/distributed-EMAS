package main

import (
	"context"
	"fmt"

	continuous "github.com/mateuszkochelski/tsp-emas/continuous"
	tsp "github.com/mateuszkochelski/tsp-emas/tsp"
)

func runTSP(ctx context.Context, simCfg *SimulationConfig, tspCfg *tsp.Config) error {
	if err := simCfg.Validate(); err != nil {
		return fmt.Errorf("invalid simulation config: %w", err)
	}
	if err := tspCfg.Validate(); err != nil {
		return fmt.Errorf("invalid tsp config: %w", err)
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

	runSimulation(ctx, problem, simCfg)
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
	runSimulation(ctx, problem, simCfg)
	return nil
}
