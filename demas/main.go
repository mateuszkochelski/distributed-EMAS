package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ShutdownWait = 5 * time.Second

func main() {
	cfg, err := ParseCLI()
	if err != nil {
		fmt.Println(err)
		return
	}

	baseCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithTimeout(baseCtx, cfg.RunDuration)
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

	<-ctx.Done()
	time.Sleep(ShutdownWait)
}

