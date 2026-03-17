package main

import (
	"fmt"
	"math/rand/v2"
)

func EMAS(
	cities []City,
	populationSize int,
	iterations int,
	initialEnergy int,
	reproductionEnergy int,
	childEnergy int,
	deathEnergy int,
) float64 {
	problem := NewProblem(cities)
	agents := make([]Agent, populationSize)

	for i := range populationSize {
		agents[i] = NewRandomAgent(problem, initialEnergy)
	}

	best := agents[0]
	for _, agent := range agents[1:] {
		if agent.Cost < best.Cost {
			best = agent
		}
	}

	for iter := range iterations {
		rand.Shuffle(len(agents), func(i, j int) {
			agents[i], agents[j] = agents[j], agents[i]
		})

		for i := 0; i+1 < len(agents); i += 2 {
			Meet(&agents[i], &agents[i+1])
		}

		newAgents := make([]Agent, 0, len(agents)*2)

		for i := range agents {
			if agents[i].Energy >= reproductionEnergy {
				childTour := MutateInsert(agents[i].Tour)
				newAgents = append(newAgents, NewAgent(problem, childTour, childEnergy))
				agents[i].Energy -= childEnergy
			}
		}

		for _, agent := range agents {
			if agent.Energy > deathEnergy {
				newAgents = append(newAgents, agent)
			}
		}

		agents = newAgents

		for _, agent := range agents {
			if agent.Cost < best.Cost {
				best = agent
			}
		}
		if iter%200 == 0 {
			fmt.Printf("Iter %d: %f\n", iter, best.Cost)
			fmt.Printf("Population %d: \n", len(agents))
		}

	}
	return best.Cost
}

func Meet(a, b *Agent) {
	if a.Cost < b.Cost {
		a.Energy++
		b.Energy--
	} else if b.Cost < a.Cost {
		a.Energy--
		b.Energy++
	}
}

