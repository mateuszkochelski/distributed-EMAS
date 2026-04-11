package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

func appendCSV(path string, headers []string, row []string) error {
	fileExists := true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fileExists = false
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if !fileExists {
		if err := w.Write(headers); err != nil {
			return err
		}
	}

	if err := w.Write(row); err != nil {
		return err
	}

	return w.Error()
}

type ContinuousAgent struct {
	X      []float64
	Energy int
	Cost   float64
}

func cloneVector(x []float64) []float64 {
	out := make([]float64, len(x))
	copy(out, x)
	return out
}

func randomVector(dim int, lower, upper float64) []float64 {
	x := make([]float64, dim)
	width := upper - lower
	for i := range x {
		x[i] = lower + rand.Float64()*width
	}
	return x
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func mutateGaussian(parent []float64, sigma, lower, upper float64) []float64 {
	child := cloneVector(parent)
	if len(child) == 0 {
		return child
	}

	k := 1 + rand.IntN(maxInt(1, len(child)/5))
	for range k {
		i := rand.IntN(len(child))
		child[i] += rand.NormFloat64() * sigma
		child[i] = clamp(child[i], lower, upper)
	}

	return child
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sphere(x []float64) float64 {
	var sum float64
	for _, xi := range x {
		sum += xi * xi
	}
	return sum
}

func rastrigin(x []float64) float64 {
	n := float64(len(x))
	sum := 10.0 * n
	for _, xi := range x {
		sum += xi*xi - 10.0*math.Cos(2.0*math.Pi*xi)
	}
	return sum
}

func rosenbrock(x []float64) float64 {
	if len(x) < 2 {
		return 0
	}

	var sum float64
	for i := 0; i < len(x)-1; i++ {
		a := x[i+1] - x[i]*x[i]
		b := 1.0 - x[i]
		sum += 100.0*a*a + b*b
	}
	return sum
}

func selectObjective(name string) (func([]float64) float64, string) {
	switch strings.ToLower(name) {
	case "sphere":
		return sphere, "sphere"
	case "rosenbrock":
		return rosenbrock, "rosenbrock"
	default:
		return rastrigin, "rastrigin"
	}
}

func runContinuousEMAS(
	dim int,
	lower float64,
	upper float64,
	sigma float64,
	populationSize int,
	iterations int,
	initialEnergy int,
	reproductionEnergy int,
	childEnergy int,
	deathEnergy int,
	objective func([]float64) float64,
) float64 {
	agents := make([]ContinuousAgent, populationSize)
	for i := range populationSize {
		x := randomVector(dim, lower, upper)
		agents[i] = ContinuousAgent{
			X:      x,
			Energy: initialEnergy,
			Cost:   objective(x),
		}
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
			if agents[i].Cost < agents[i+1].Cost {
				agents[i].Energy++
				agents[i+1].Energy--
			} else if agents[i+1].Cost < agents[i].Cost {
				agents[i].Energy--
				agents[i+1].Energy++
			}
		}

		newAgents := make([]ContinuousAgent, 0, len(agents)*2)
		for i := range agents {
			if agents[i].Energy >= reproductionEnergy {
				childX := mutateGaussian(agents[i].X, sigma, lower, upper)
				newAgents = append(newAgents, ContinuousAgent{
					X:      childX,
					Energy: childEnergy,
					Cost:   objective(childX),
				})
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
			fmt.Printf("Population %d\n", len(agents))
		}
	}

	return best.Cost
}

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
	mode := flag.String("mode", "continuous", "optimization mode: continuous or tsp")
	objectiveName := flag.String("objective", "rastrigin", "continuous objective: rastrigin, sphere, rosenbrock")
	dim := flag.Int("dim", 30, "continuous problem dimension")
	lower := flag.Float64("lower", -5.12, "continuous lower bound")
	upper := flag.Float64("upper", 5.12, "continuous upper bound")
	sigma := flag.Float64("sigma", 0.2, "continuous mutation sigma")
	csvPath := flag.String("csv", "", "optional output CSV file path")
	runs := flag.Int("runs", 1, "number of independent runs")

	populationSize := flag.Int("pop", 50, "population size")
	iterations := flag.Int("iter", 100000, "number of iterations")
	initialEnergy := flag.Int("e0", 10, "initial energy")
	reproductionEnergy := flag.Int("erepro", 15, "reproduction threshold")
	childEnergy := flag.Int("echild", 5, "child energy")
	deathEnergy := flag.Int("edeath", 0, "death threshold")
	flag.Parse()

	selectedMode := strings.ToLower(*mode)
	if selectedMode == "continuous" {
		objective, objectiveLabel := selectObjective(*objectiveName)
		for run := 1; run <= *runs; run++ {
			result := runContinuousEMAS(
				*dim,
				*lower,
				*upper,
				*sigma,
				*populationSize,
				*iterations,
				*initialEnergy,
				*reproductionEnergy,
				*childEnergy,
				*deathEnergy,
				objective,
			)
			fmt.Printf("Run %d: Continuous EMAS (%s) result %f\n", run, objectiveLabel, result)

			if *csvPath != "" {
				err := appendCSV(
					*csvPath,
					[]string{"timestamp", "mode", "objective", "run", "result", "dim", "lower", "upper", "sigma", "pop", "iter", "e0", "erepro", "echild", "edeath"},
					[]string{
						time.Now().Format(time.RFC3339),
						"continuous",
						objectiveLabel,
						strconv.Itoa(run),
						fmt.Sprintf("%.12f", result),
						strconv.Itoa(*dim),
						fmt.Sprintf("%.6f", *lower),
						fmt.Sprintf("%.6f", *upper),
						fmt.Sprintf("%.6f", *sigma),
						strconv.Itoa(*populationSize),
						strconv.Itoa(*iterations),
						strconv.Itoa(*initialEnergy),
						strconv.Itoa(*reproductionEnergy),
						strconv.Itoa(*childEnergy),
						strconv.Itoa(*deathEnergy),
					},
				)
				if err != nil {
					fmt.Printf("failed to write CSV: %v\n", err)
				}
			}
		}
		return
	}
	for run := 1; run <= *runs; run++ {
		cities := make([]City, 100)
		for i := range len(cities) {
			cities[i] = City{
				X: -500.0 + rand.Float64()*1000,
				Y: -500.0 + rand.Float64()*1000,
			}
		}

		nearestNeighbourTour := NearestNeighbourTSP(NewProblem(cities))
		emasResult := EMAS(
			cities,
			*populationSize,
			*iterations,
			*initialEnergy,
			*reproductionEnergy,
			*childEnergy,
			*deathEnergy,
		)

		fmt.Printf("Run %d: EMAS result %f\n", run, emasResult)
		fmt.Printf("Run %d: nearestNeighbour result %f\n", run, nearestNeighbourTour)

		if *csvPath != "" {
			err := appendCSV(
				*csvPath,
				[]string{"timestamp", "mode", "run", "emas_result", "nearest_neighbour", "cities", "pop", "iter", "e0", "erepro", "echild", "edeath"},
				[]string{
					time.Now().Format(time.RFC3339),
					"tsp",
					strconv.Itoa(run),
					fmt.Sprintf("%.12f", emasResult),
					fmt.Sprintf("%.12f", nearestNeighbourTour),
					strconv.Itoa(len(cities)),
					strconv.Itoa(*populationSize),
					strconv.Itoa(*iterations),
					strconv.Itoa(*initialEnergy),
					strconv.Itoa(*reproductionEnergy),
					strconv.Itoa(*childEnergy),
					strconv.Itoa(*deathEnergy),
				},
			)
			if err != nil {
				fmt.Printf("failed to write CSV: %v\n", err)
			}
		}
	}

}
