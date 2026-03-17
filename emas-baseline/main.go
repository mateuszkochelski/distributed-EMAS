package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

func FindNearestCity(cities []City, visited []bool, id int) int {
	n := len(cities)
	bestDist := math.Inf(1)
	bestCity := -1
	for i := range n {
		if visited[i] {
			continue
		}
		if newDist := Distance(cities[i], cities[id]); newDist < bestDist {
			bestDist = newDist
			bestCity = i
		}
	}

	return bestCity
}

func NearestNeighbourTSP(problem Problem) float64 {
	n := len(problem.Cities)
	start := rand.IntN(n)
	tour := []int{start}
	visited := make([]bool, n)
	visited[start] = true

	for range n - 1 {
		nearestNeighbour := FindNearestCity(problem.Cities, visited, tour[len(tour)-1])
		visited[nearestNeighbour] = true
		tour = append(tour, nearestNeighbour)
	}

	return TourLength(problem, tour)
}

func main() {

	cities := make([]City, 100)

	for i := range len(cities) {
		cities[i] = City{
			X: -500.0 + rand.Float64()*1000,
			Y: -500.0 + rand.Float64()*1000,
		}
	}

	// randomTour := RandomTSP(cities)
	nearestNeighbourTour := NearestNeighbourTSP(NewProblem(cities))

	populationSize := 50
	iterations := 100000
	initialEnergy := 10
	reproductionEnergy := 15
	childEnergy := 5
	deathEnergy := 0

	emasResult := EMAS(
		cities,
		populationSize,
		iterations,
		initialEnergy,
		reproductionEnergy,
		childEnergy,
		deathEnergy,
	)

	fmt.Printf("EMAS result %f\n", emasResult)
	fmt.Printf("nearestNeighbour result %f\n", nearestNeighbourTour)

}
