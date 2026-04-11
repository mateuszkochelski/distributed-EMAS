package tsp

import (
	"math"
	"math/rand/v2"
)

type Solution struct {
	Tour []int
}

func (s Solution) Clone() Solution {
	tourCopy := make([]int, len(s.Tour))
	copy(tourCopy, s.Tour)

	return Solution{
		Tour: tourCopy,
	}
}

type City struct {
	X float64
	Y float64
}

func GenerateRandomCities(n int, minCoord, maxCoord float64) []City {
	if n <= 0 {
		return nil
	}

	cities := make([]City, n)
	width := maxCoord - minCoord

	for i := range cities {
		cities[i] = City{
			X: minCoord + rand.Float64()*width,
			Y: minCoord + rand.Float64()*width,
		}
	}

	return cities
}

type Problem struct {
	Cities []City
	Dist   []float64
}

func Distance(a, b City) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func BuildDistanceMatrix1D(cities []City) []float64 {
	n := len(cities)
	dist := make([]float64, n*n)

	for i := range n {
		for j := i; j < n; j++ {
			d := Distance(cities[i], cities[j])
			dist[i*n+j] = d
			dist[j*n+i] = d
		}
	}

	return dist
}

func NewProblem(cities []City) Problem {
	return Problem{
		Cities: cities,
		Dist:   BuildDistanceMatrix1D(cities),
	}
}

func NewRandomProblem(numCities int, minCoord, maxCoord float64) Problem {
	cities := GenerateRandomCities(numCities, minCoord, maxCoord)
	return NewProblem(cities)
}

func (p Problem) distance(i, j int) float64 {
	return p.Dist[i*len(p.Cities)+j]
}

func (p Problem) NewRandomSolution() Solution {
	return Solution{
		Tour: rand.Perm(len(p.Cities)),
	}
}

func (p Problem) MutateSolution(s Solution) Solution {
	mutated := s.Clone()
	n := len(mutated.Tour)

	if n < 2 {
		return mutated
	}

	from := rand.IntN(n)
	moved := mutated.Tour[from]

	without := make([]int, 0, n-1)
	without = append(without, mutated.Tour[:from]...)
	without = append(without, mutated.Tour[from+1:]...)

	to := rand.IntN(len(without) + 1)

	child := make([]int, 0, n)
	child = append(child, without[:to]...)
	child = append(child, moved)
	child = append(child, without[to:]...)

	return Solution{
		Tour: child,
	}
}

func (p Problem) Evaluate(s Solution) float64 {
	n := len(s.Tour)
	if n < 2 {
		return 0
	}

	var total float64
	for i := range n {
		from := s.Tour[i]
		to := s.Tour[(i+1)%n]
		total += p.distance(from, to)
	}

	return total
}

func NearestNeighbourTSP(problem Problem, start int) float64 {
	n := len(problem.Cities)
	if start < 0 || start >= n {
		start = 0
	}

	visited := make([]bool, n)
	tour := make([]int, 0, n)

	current := start
	tour = append(tour, current)
	visited[current] = true

	for len(tour) < n {
		bestCity := -1
		bestDist := math.Inf(1)

		for candidate := range n {
			if visited[candidate] {
				continue
			}

			dist := problem.distance(current, candidate)
			if dist < bestDist {
				bestDist = dist
				bestCity = candidate
			}
		}

		visited[bestCity] = true
		tour = append(tour, bestCity)
		current = bestCity
	}
	return problem.Evaluate(Solution{
		Tour: tour,
	})
}

func (p Problem) IsBetter(newScore, oldScore float64) bool {
	return newScore < oldScore
}

func (p Problem) Maximum() float64 {
	return math.Inf(1)
}
