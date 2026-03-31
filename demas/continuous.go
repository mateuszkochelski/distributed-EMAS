package main

import (
	"math"
	"math/rand/v2"
)

type ContinuousSolution struct {
	X []float64
}

func (s ContinuousSolution) Clone() ContinuousSolution {
	out := make([]float64, len(s.X))
	copy(out, s.X)
	return ContinuousSolution{X: out}
}

type ObjectiveFunc func(x []float64) float64

type ContinuousProblem struct {
	Dim       int
	Lower     float64
	Upper     float64
	Sigma     float64
	Objective ObjectiveFunc
}

func (p ContinuousProblem) NewRandomSolution() ContinuousSolution {
	x := make([]float64, p.Dim)
	width := p.Upper - p.Lower
	for i := range x {
		x[i] = p.Lower + rand.Float64()*width
	}
	return ContinuousSolution{X: x}
}

func (p ContinuousProblem) MutateSolution(s ContinuousSolution) ContinuousSolution {
	child := s.Clone()
	if p.Dim == 0 {
		return child
	}

	k := 1 + rand.IntN(maxInt(1, p.Dim/5))
	for range k {
		i := rand.IntN(p.Dim)
		child.X[i] += rand.NormFloat64() * p.Sigma
		child.X[i] = clamp(child.X[i], p.Lower, p.Upper)
	}

	return child
}

func (p ContinuousProblem) Evaluate(s ContinuousSolution) float64 {
	return p.Objective(s.X)
}

func (p ContinuousProblem) IsBetter(newScore, oldScore float64) bool {
	return newScore < oldScore
}

func Sphere(x []float64) float64 {
	var sum float64
	for _, xi := range x {
		sum += xi * xi
	}
	return sum
}

func Rastrigin(x []float64) float64 {
	n := float64(len(x))
	sum := 10.0 * n
	for _, xi := range x {
		sum += xi*xi - 10.0*math.Cos(2.0*math.Pi*xi)
	}
	return sum
}

func Rosenbrock(x []float64) float64 {
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

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
