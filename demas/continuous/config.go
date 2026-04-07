package continuous

import (
	"flag"
	"fmt"
)

type Objective string

const (
	ObjectiveSphere     Objective = "sphere"
	ObjectiveRastrigin  Objective = "rastrigin"
	ObjectiveRosenbrock Objective = "rosenbrock"
)

type Config struct {
	Dim       int
	Lower     float64
	Upper     float64
	Sigma     float64
	Objective Objective
}

func DefaultConfig() Config {
	return Config{
		Dim:       30,
		Lower:     -5.12,
		Upper:     5.12,
		Sigma:     0.2,
		Objective: ObjectiveRastrigin,
	}
}

func BindFlags(fs *flag.FlagSet, cfg *Config) {
	fs.IntVar(&cfg.Dim, "cont-dim", cfg.Dim, "Dimension of the continuous problem")
	fs.Float64Var(&cfg.Lower, "cont-lower", cfg.Lower, "Lower bound")
	fs.Float64Var(&cfg.Upper, "cont-upper", cfg.Upper, "Upper bound")
	fs.Float64Var(&cfg.Sigma, "cont-sigma", cfg.Sigma, "Mutation sigma")
	fs.StringVar((*string)(&cfg.Objective), "cont-objective", string(cfg.Objective), "Objective: sphere, rastrigin, rosenbrock")
}

func (c Config) Validate() error {
	if c.Dim <= 0 {
		return fmt.Errorf("dim must be > 0")
	}
	if c.Lower >= c.Upper {
		return fmt.Errorf("lower must be < upper")
	}
	if c.Sigma <= 0 {
		return fmt.Errorf("sigma must be > 0")
	}

	switch c.Objective {
	case ObjectiveSphere, ObjectiveRastrigin, ObjectiveRosenbrock:
		return nil
	default:
		return fmt.Errorf("invalid objective %q", c.Objective)
	}
}
