package main

import (
	"context"
	"math"
	"math/rand/v2"
	"strconv"
)

type Message struct {
	From            string
	Score           float64
	Body            int
	ResponseCh      chan Message
	PrimaryIslandId int
}

type Agent[S any] struct {
	ID            string
	Energy        int
	IslandsChs    []chan Message
	InboxCh       chan Message
	PrimaryIsland int
	LogCh         chan Log

	Solution S
	Score    float64
	Problem  Problem[S]
	Config   *SimulationConfig
}

func (a *Agent[S]) selectTargetIslandCh() chan Message {
	if rand.Float64() < a.Config.SameIslandProbability || len(a.IslandsChs) == 1 {
		return a.IslandsChs[a.PrimaryIsland]
	}

	targetIdx := rand.IntN(len(a.IslandsChs))
	for targetIdx == a.PrimaryIsland {
		targetIdx = rand.IntN(len(a.IslandsChs))
	}

	return a.IslandsChs[targetIdx]
}

func (a *Agent[S]) CreateMeetingMessage(score float64) Message {
	return Message{
		From:            a.ID,
		Score:           score,
		Body:            a.PrimaryIsland,
		ResponseCh:      a.InboxCh,
		PrimaryIslandId: a.PrimaryIsland,
	}
}

func (a *Agent[S]) resolveMeeting(reply Message) {
	enemyScore := reply.Score

	if a.Problem.IsBetter(a.Score, enemyScore) {
		a.Energy += a.Config.EnergyTransfer
	} else if a.Problem.IsBetter(enemyScore, a.Score) {
		a.Energy -= a.Config.EnergyTransfer
	}

	a.LogCh <- Log{
		Score:           a.Score,
		SameIsland:      a.PrimaryIsland == reply.PrimaryIslandId,
		PrimaryIslandId: a.PrimaryIsland,
		EventType:       "MEETING",
	}
}

func (a *Agent[S]) runMeetingAsResponder(incoming Message) {
	incoming.ResponseCh <- a.CreateMeetingMessage(a.Score)
	a.resolveMeeting(incoming)
}

func (a *Agent[S]) Run(ctx context.Context) {
	for {
		target := a.selectTargetIslandCh()

		if a.Energy >= a.Config.ReproductionThreshold {
			child := NewChildAgent(a, "CHILD")
			a.Energy -= a.Config.ChildEnergy

			a.LogCh <- Log{
				Score:           a.Score,
				PrimaryIslandId: a.PrimaryIsland,
				EventType:       "BORN",
			}

			go child.Run(ctx)
		} else if a.Energy <= a.Config.DeathThreshold {
			a.LogCh <- Log{
				Score:           a.Score,
				PrimaryIslandId: a.PrimaryIsland,
				EventType:       "DEAD",
			}
			return
		}

		select {
		case <-ctx.Done():
			return

		case target <- a.CreateMeetingMessage(a.Score):
			reply := <-a.InboxCh
			a.resolveMeeting(reply)

		case incoming := <-a.IslandsChs[a.PrimaryIsland]:
			a.runMeetingAsResponder(incoming)
		}
	}
}

func NewChildAgent[S any](parent *Agent[S], id string) Agent[S] {
	childSolution := parent.Problem.MutateSolution(parent.Solution)
	childScore := parent.Problem.Evaluate(childSolution)

	return Agent[S]{
		ID:            id,
		Energy:        parent.Config.ChildEnergy,
		InboxCh:       make(chan Message),
		IslandsChs:    parent.IslandsChs,
		PrimaryIsland: parent.PrimaryIsland,
		LogCh:         parent.LogCh,
		Solution:      childSolution,
		Score:         childScore,
		Problem:       parent.Problem,
		Config:        parent.Config,
	}
}

func NewAgent[S any](
	id string,
	islands []chan Message,
	logCh chan Log,
	problem Problem[S],
	config *SimulationConfig,
) Agent[S] {
	solution := problem.NewRandomSolution()

	return Agent[S]{
		ID:            id,
		Energy:        config.InitialEnergy,
		InboxCh:       make(chan Message),
		IslandsChs:    islands,
		PrimaryIsland: rand.IntN(len(islands)),
		LogCh:         logCh,
		Solution:      solution,
		Score:         problem.Evaluate(solution),
		Problem:       problem,
		Config:        config,
	}
}

func runSimulation[S any](ctx context.Context, problem Problem[S], config *SimulationConfig) {
	logChan := make(chan Log, 1024)
	islands := make([]chan Message, config.NumIslands)

	bufferSize := max(int(math.Sqrt(float64(config.NumAgents/config.NumIslands))), 1)

	for i := range islands {
		islands[i] = make(chan Message, bufferSize)
	}

	agentsPerIsland := make(map[int]int)
	agents := make([]Agent[S], config.NumAgents)

	for i := range agents {
		agents[i] = NewAgent(strconv.Itoa(i), islands, logChan, problem, config)
		agentsPerIsland[agents[i].PrimaryIsland]++
	}

	logger := LogCollector{
		LogCh:           logChan,
		AgentsPerIsland: agentsPerIsland,
	}

	go logger.Run(ctx)

	for i := range agents {
		go agents[i].Run(ctx)
	}
}
