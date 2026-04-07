package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	continuous "github.com/mateuszkochelski/tsp-emas/continuous"
	tsp "github.com/mateuszkochelski/tsp-emas/tsp"
)

type CLIConfig struct {
	Mode       Mode
	Simulation SimulationConfig
	TSP        *tsp.Config
	Continuous *continuous.Config
}

func parseMode(args []string) (Mode, error) {
	mode := ModeTSP

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "-mode" {
			if i+1 >= len(args) {
				return "", fmt.Errorf("missing value for -mode")
			}
			mode = Mode(args[i+1])
			i++
			continue
		}

		const prefix = "-mode="
		if after, ok := strings.CutPrefix(arg, prefix); ok {
			mode = Mode(after)
		}
	}

	if isValidMode(mode) {
		return mode, nil
	} else {
		return "", fmt.Errorf("invalid mode %q, expected one of: %q, %q", mode, ModeTSP, ModeContinuous)
	}
}

func ParseCLI() (CLIConfig, error) {
	mode, err := parseMode(os.Args[1:])
	if err != nil {
		return CLIConfig{}, err
	}

	simCfg := DefaultSimulationConfig()

	switch mode {
	case ModeTSP:
		tspCfg := tsp.DefaultConfig()

		fs := flag.NewFlagSet("tsp", flag.ContinueOnError)
		fs.SetOutput(os.Stderr)

		fs.String("mode", string(ModeTSP), "Problem mode: tsp or continuous")
		BindSimulationFlags(fs, &simCfg)
		tsp.BindFlags(fs, &tspCfg)

		if err := fs.Parse(os.Args[1:]); err != nil {
			return CLIConfig{}, err
		}
		if err := simCfg.Validate(); err != nil {
			return CLIConfig{}, fmt.Errorf("invalid simulation config: %w", err)
		}
		if err := tspCfg.Validate(); err != nil {
			return CLIConfig{}, fmt.Errorf("invalid tsp config: %w", err)
		}

		return CLIConfig{
			Mode:       mode,
			Simulation: simCfg,
			TSP:        &tspCfg,
		}, nil

	case ModeContinuous:
		continuousCfg := continuous.DefaultConfig()

		fs := flag.NewFlagSet("continuous", flag.ContinueOnError)
		fs.SetOutput(os.Stderr)

		fs.String("mode", string(ModeContinuous), "Problem mode: tsp or continuous")
		BindSimulationFlags(fs, &simCfg)
		continuous.BindFlags(fs, &continuousCfg)

		if err := fs.Parse(os.Args[1:]); err != nil {
			return CLIConfig{}, err
		}
		if err := simCfg.Validate(); err != nil {
			return CLIConfig{}, fmt.Errorf("invalid simulation config: %w", err)
		}
		if err := continuousCfg.Validate(); err != nil {
			return CLIConfig{}, fmt.Errorf("invalid continuous config: %w", err)
		}

		return CLIConfig{
			Mode:       mode,
			Simulation: simCfg,
			Continuous: &continuousCfg,
		}, nil
	}

	return CLIConfig{}, fmt.Errorf("unsupported mode %q", mode)
}
