package main

import (
	"flag"
	"fmt"
	"strings"
)

type SimulationConfig struct {
	NumAgents             int     `json:"num_agents"`
	NumIslands            int     `json:"num_islands"`
	InitialEnergy         int     `json:"initial_energy"`
	EnergyTransfer        int     `json:"energy_transfer"`
	ReproductionThreshold int     `json:"reproduction_threshold"`
	ChildEnergy           int     `json:"child_energy"`
	DeathThreshold        int     `json:"death_threshold"`
	SameIslandProbability float64 `json:"same_island_probability"`
}

func DefaultSimulationConfig() SimulationConfig {
	return SimulationConfig{
		NumAgents:             10_000,
		NumIslands:            100,
		InitialEnergy:         10,
		EnergyTransfer:        1,
		ReproductionThreshold: 15,
		ChildEnergy:           5,
		DeathThreshold:        0,
		SameIslandProbability: 0.999,
	}
}

func BindSimulationFlags(fs *flag.FlagSet, cfg *SimulationConfig) {
	fs.IntVar(&cfg.NumAgents, "num-agents", cfg.NumAgents, "Number of agents")
	fs.IntVar(&cfg.NumIslands, "num-islands", cfg.NumIslands, "Number of islands")
	fs.IntVar(&cfg.InitialEnergy, "initial-energy", cfg.InitialEnergy, "Initial energy of each agent")
	fs.IntVar(&cfg.EnergyTransfer, "energy-transfer", cfg.EnergyTransfer, "Energy transferred after a meeting")
	fs.IntVar(&cfg.ReproductionThreshold, "reproduction-threshold", cfg.ReproductionThreshold, "Energy threshold for reproduction")
	fs.IntVar(&cfg.ChildEnergy, "child-energy", cfg.ChildEnergy, "Initial energy given to a child")
	fs.IntVar(&cfg.DeathThreshold, "death-threshold", cfg.DeathThreshold, "Energy threshold at or below which an agent dies")
	fs.Float64Var(
		&cfg.SameIslandProbability,
		"same-island-probability",
		cfg.SameIslandProbability,
		"Probability of selecting the agent's own island",
	)
}

func (c SimulationConfig) Validate() error {
	var errs []string

	if c.NumAgents <= 0 {
		errs = append(errs, "numAgents must be > 0")
	}
	if c.NumIslands <= 0 {
		errs = append(errs, "numIslands must be > 0")
	}
	if c.InitialEnergy < 0 {
		errs = append(errs, "initialEnergy must be >= 0")
	}
	if c.EnergyTransfer <= 0 {
		errs = append(errs, "energyTransfer must be > 0")
	}
	if c.ReproductionThreshold < 0 {
		errs = append(errs, "reproductionThreshold must be >= 0")
	}
	if c.ChildEnergy < 0 {
		errs = append(errs, "childEnergy must be >= 0")
	}
	if c.SameIslandProbability < 0.0 || c.SameIslandProbability > 1.0 {
		errs = append(errs, "sameIslandProbability must be between 0 and 1")
	}
	if c.NumIslands > c.NumAgents {
		errs = append(errs, "numIslands should not be greater than numAgents")
	}
	if c.ChildEnergy > c.ReproductionThreshold {
		errs = append(errs, "childEnergy should not be greater than reproductionThreshold")
	}

	if len(errs) > 0 {
		return fmt.Errorf("invalid simulation config:\n- %s", strings.Join(errs, "\n- "))
	}

	return nil
}
