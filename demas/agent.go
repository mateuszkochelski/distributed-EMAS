package main

import (
	"context"
	"math/rand/v2"

	"github.com/mateuszkochelski/tsp-emas/stats"
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
	EventCh       chan stats.Event

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

func (a *Agent[S]) emitEvent(event stats.Event) {
	select {
	case a.EventCh <- event:
	default:
	}
}

func (a *Agent[S]) resolveMeeting(reply Message) {
	enemyScore := reply.Score

	if a.Problem.IsBetter(a.Score, enemyScore) {
		a.Energy += a.Config.EnergyTransfer
	} else if a.Problem.IsBetter(enemyScore, a.Score) {
		a.Energy -= a.Config.EnergyTransfer
	}

	a.emitEvent(stats.Event{
		Score:           a.Score,
		SameIsland:      a.PrimaryIsland == reply.PrimaryIslandId,
		PrimaryIslandID: a.PrimaryIsland,
		EventType:       stats.EventMeeting,
	})
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

			a.emitEvent(stats.Event{
				Score:           a.Score,
				PrimaryIslandID: a.PrimaryIsland,
				EventType:       stats.EventBorn,
			})

			go child.Run(ctx)
		} else if a.Energy <= a.Config.DeathThreshold {
			a.emitEvent(stats.Event{
				Score:           a.Score,
				PrimaryIslandID: a.PrimaryIsland,
				EventType:       stats.EventDead,
			})
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
		EventCh:       parent.EventCh,
		Solution:      childSolution,
		Score:         childScore,
		Problem:       parent.Problem,
		Config:        parent.Config,
	}
}

func NewAgent[S any](
	id string,
	islands []chan Message,
	eventCh chan stats.Event,
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
		EventCh:       eventCh,
		Solution:      solution,
		Score:         problem.Evaluate(solution),
		Problem:       problem,
		Config:        config,
	}
}
