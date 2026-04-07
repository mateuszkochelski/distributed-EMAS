package tsp

import (
	"flag"
	"fmt"
)

type Source string

const (
	SourceRandom Source = "random"
	SourceFile   Source = "file"
)

type Config struct {
	Source    Source
	NumCities int
	MinCoord  float64
	MaxCoord  float64
	InputFile string
}

func DefaultConfig() Config {
	return Config{
		Source:    SourceRandom,
		NumCities: 100,
		MinCoord:  -100,
		MaxCoord:  100,
		InputFile: "",
	}
}

func BindFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar((*string)(&cfg.Source), "tsp-source", string(cfg.Source), "TSP source: random or file")
	fs.IntVar(&cfg.NumCities, "tsp-num-cities", cfg.NumCities, "Number of cities for random TSP")
	fs.Float64Var(&cfg.MinCoord, "tsp-min-coord", cfg.MinCoord, "Minimum coordinate for random TSP")
	fs.Float64Var(&cfg.MaxCoord, "tsp-max-coord", cfg.MaxCoord, "Maximum coordinate for random TSP")
	fs.StringVar(&cfg.InputFile, "tsp-input-file", cfg.InputFile, "Path to TSPLIB input file")
}

func (c Config) Validate() error {
	switch c.Source {
	case SourceRandom:
		if c.NumCities <= 1 {
			return fmt.Errorf("cities must be > 1")
		}
		if c.MinCoord >= c.MaxCoord {
			return fmt.Errorf("min-coord must be < max-coord")
		}
	case SourceFile:
		if c.InputFile == "" {
			return fmt.Errorf("input-file is required when tsp-source=file")
		}
	default:
		return fmt.Errorf("invalid tsp source %q", c.Source)
	}

	return nil
}

