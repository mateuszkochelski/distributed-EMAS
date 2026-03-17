package main

import (
	"math"
	"math/rand/v2"
)

type City struct {
	X float64
	Y float64
}

type Problem struct {
	Cities []City
	Dist   []float64
	N      int
}

func (p Problem) D(i, j int) float64 {
	return p.Dist[i*p.N+j]
}

func Distance(a, b City) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func BuildDistanceMatrix1D(cities []City) []float64 {
	n := len(cities)
	dist := make([]float64, n*n)

	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			d := Distance(cities[i], cities[j])
			dist[i*n+j] = d
			dist[j*n+i] = d
		}
	}

	return dist
}

func RandomTSP(cities []City) []int {
	n := len(cities)
	tour := make([]int, n)

	for i := range tour {
		tour[i] = i
	}

	rand.Shuffle(n, func(i, j int) {
		tour[i], tour[j] = tour[j], tour[i]
	})

	return tour
}

func NewProblem(cities []City) Problem {
	return Problem{
		Cities: cities,
		Dist:   BuildDistanceMatrix1D(cities),
		N:      len(cities),
	}
}

func TourLength(problem Problem, tour []int) float64 {
	n := len(tour)
	if n == 0 {
		return 0
	}

	dist := problem.Dist
	size := problem.N

	var total float64
	for i := 0; i < n-1; i++ {
		total += dist[tour[i]*size+tour[i+1]]
	}
	total += dist[tour[n-1]*size+tour[0]]

	return total
}
