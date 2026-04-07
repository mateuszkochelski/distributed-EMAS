package main

import (
	"context"
	"fmt"
	"time"
)

const (
	RunDuration  = 300 * time.Second
	ShutdownWait = 5 * time.Second
)

func main() {
	cfg, err := ParseCLI()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch cfg.Mode {
	case ModeTSP:
		if err := runTSP(ctx, &cfg.Simulation, cfg.TSP); err != nil {
			fmt.Println(err)
			return
		}
	case ModeContinuous:
		if err := runContinuous(ctx, &cfg.Simulation, cfg.Continuous); err != nil {
			fmt.Println(err)
			return
		}
	}

	time.Sleep(RunDuration)
	cancel()
	time.Sleep(ShutdownWait)
}
