package main

type Agent struct {
	Tour   []int
	Energy int
	Cost   float64
}

func NewAgent(problem Problem, tour []int, energy int) Agent {
	return Agent{
		Tour:   tour,
		Energy: energy,
		Cost:   TourLength(problem, tour),
	}
}

func NewRandomAgent(problem Problem, energy int) Agent {
	tour := RandomTSP(problem.Cities)
	cost := TourLength(problem, tour)

	return Agent{
		Tour:   tour,
		Energy: energy,
		Cost:   cost,
	}
}
