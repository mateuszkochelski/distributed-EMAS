package main

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"
	"time"
)

const (
	NumAgents             = 10_000
	NumIslands            = 100
	InitialEnergy         = 10
	EnergyTransfer        = 1
	ReproductionThreshold = 15
	ChildEnergy           = 5
	DeathThreshold        = 0
	SameIslandProbability = 0.9999
)

type Log struct {
	From            string
	Msg             string
	Score           float64
	sameIsland      bool
	PrimaryIslandId int
}

type Message struct {
	From       string
	Score      float64
	Body       int
	ResponseCh chan Message
}

type LogCollector struct {
	LogCh chan Log
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
			counts[log.PrimaryIslandId]++
			if log.Score < minimum_distance {
				minimum_distance = log.Score
				fmt.Println("new minimum:", minimum_distance, "discovered on island: ", log.PrimaryIslandId)
			}
			if log.sameIsland {
				sameIslandLogs++
			} else {
				differentIslandLogs++
			}
			if (sameIslandLogs+differentIslandLogs)%1_000_000 == 0 {
				fmt.Println("Meetings: ", (sameIslandLogs+differentIslandLogs)/1_000_000, " mln")
				fmt.Println("SAME ISLANDS: ", sameIslandLogs, " DIFFERENT ISLANDS: ", differentIslandLogs)
				fmt.Println("PERCENTAGE: ", float64(sameIslandLogs)/float64(sameIslandLogs+differentIslandLogs))
			}
		}
	}
}

type Agent[S any] struct {
	ID string
	// Tour          []int
	Energy        int
	IslandsChs    []chan Message
	InboxCh       chan Message
	PrimaryIsland int
	LogCh         chan Log

	Solution S
	Score    float64
	Problem  Problem[S]
}

func (a *Agent[S]) selectTargetIslandCh() chan Message {
	if rand.Float64() < SameIslandProbability || len(a.IslandsChs) == 1 {
		return a.IslandsChs[a.PrimaryIsland]
	}

	targetIdx := rand.IntN(len(a.IslandsChs))
	for targetIdx == a.PrimaryIsland {
		targetIdx = rand.IntN(len(a.IslandsChs))
	}

	return a.IslandsChs[targetIdx]
}

func (a *Agent[S]) newHandshakeMessage() Message {
	return Message{
		From:       a.ID,
		ResponseCh: a.InboxCh,
	}
}

func (a *Agent[S]) CreateMeetingMessage(score float64) Message {
	return Message{
		From:       a.ID,
		Score:      score,
		Body:       a.PrimaryIsland,
		ResponseCh: a.InboxCh,
	}
}

func (a *Agent[S]) runMeetingAsInitiator(ctx context.Context, sessionCh chan Message) {
	sessionCh <- a.CreateMeetingMessage(a.Score)
	msg := <-a.InboxCh

	enemyScore := msg.Score
	if a.Problem.IsBetter(a.Score, enemyScore) {
		a.Energy += EnergyTransfer
	} else if a.Problem.IsBetter(enemyScore, a.Score) {
		a.Energy -= EnergyTransfer
	}
}

func (a *Agent[S]) runMeetingAsResponder(ctx context.Context, sessionCh chan Message) {
	msg := <-sessionCh
	enemyScore := msg.Score

	msg.ResponseCh <- a.CreateMeetingMessage(a.Score)

	if a.Problem.IsBetter(a.Score, enemyScore) {
		a.Energy += EnergyTransfer
	} else if a.Problem.IsBetter(enemyScore, a.Score) {
		a.Energy -= EnergyTransfer
	}

	a.LogCh <- Log{
		Score:           a.Score,
		sameIsland:      a.PrimaryIsland == msg.Body,
		PrimaryIslandId: a.PrimaryIsland,
	}
}

func (a *Agent[S]) Run(ctx context.Context) {
	for {
		target := a.selectTargetIslandCh()
		ping := a.newHandshakeMessage()

		if a.Energy >= ReproductionThreshold {
			child := NewChildAgent(a, "CHILD", ChildEnergy)
			a.Energy -= ChildEnergy
			go child.Run(ctx)
		} else if a.Energy <= DeathThreshold {
			return
		}

		select {
		case <-ctx.Done():
			return

		case target <- ping:
			response := <-a.InboxCh
			a.runMeetingAsInitiator(ctx, response.ResponseCh)

		case incoming := <-a.IslandsChs[a.PrimaryIsland]:
			pong := a.newHandshakeMessage()
			incoming.ResponseCh <- pong
			a.runMeetingAsResponder(ctx, pong.ResponseCh)
		}
	}
}

func NewChildAgent[S any](parent *Agent[S], id string, childEnergy int) Agent[S] {
	childSolution := parent.Problem.MutateSolution(parent.Solution)
	childScore := parent.Problem.Evaluate(childSolution)

	return Agent[S]{
		ID:            id,
		Energy:        childEnergy,
		InboxCh:       make(chan Message),
		IslandsChs:    parent.IslandsChs,
		PrimaryIsland: parent.PrimaryIsland,
		LogCh:         parent.LogCh,
		Solution:      childSolution,
		Score:         childScore,
		Problem:       parent.Problem,
	}
}

func NewAgent[S any](id string, islands []chan Message, logCh chan Log, initialEnergy int, problem Problem[S]) Agent[S] {
	solution := problem.NewRandomSolution()
	return Agent[S]{
		ID:            id,
		Energy:        initialEnergy,
		InboxCh:       make(chan Message),
		IslandsChs:    islands,
		PrimaryIsland: rand.IntN(len(islands)),
		LogCh:         logCh,

		Solution: solution,
		Score:    problem.Evaluate(solution),
		Problem:  problem,
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	logChan := make(chan Log, 1000)
	logger := LogCollector{
		LogCh: logChan,
	}
	go logger.Run(ctx)

	cities := GenerateRandomCities(100, -100, 100)
	problem := NewTSPProblem(cities)

	islands := make([]chan Message, NumIslands)
	for i := range NumIslands {
		islands[i] = make(chan Message, int64(math.Sqrt(NumAgents/NumIslands)))
	}

	agents := make([]Agent[TSPSolution], NumAgents)
	for i := range NumAgents {
		agents[i] = NewAgent(strconv.Itoa(i), islands, logChan, InitialEnergy, problem)
	}

	for i := range NumAgents {
		go agents[i].Run(ctx)
	}

	time.Sleep(120 * time.Second)
	cancel()

	fmt.Println("simple heuristic: ", NearestNeighbourTSP(problem, 0))
	time.Sleep(5 * time.Second)
}
