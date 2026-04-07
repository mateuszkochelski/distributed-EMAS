package continuous

import (
	"math"
	"math/rand/v2"
)

type Solution struct {
	X []float64
}

func (s Solution) Clone() Solution {
	out := make([]float64, len(s.X))
	copy(out, s.X)
	return Solution{X: out}
}

type Problem struct {
	Config Config
}

func NewProblem(cfg Config) Problem {
	return Problem{Config: cfg}
}

func (p Problem) NewRandomSolution() Solution {
	x := make([]float64, p.Config.Dim)
	width := p.Config.Upper - p.Config.Lower

	for i := range x {
		x[i] = p.Config.Lower + rand.Float64()*width
	}

	return Solution{X: x}
}

func (p Problem) MutateSolution(s Solution) Solution {
	child := s.Clone()
	if p.Config.Dim == 0 {
		return child
	}

	k := 1 + rand.IntN(maxInt(1, p.Config.Dim/5))
	for range k {
		i := rand.IntN(p.Config.Dim)
		child.X[i] += rand.NormFloat64() * p.Config.Sigma
		child.X[i] = clamp(child.X[i], p.Config.Lower, p.Config.Upper)
	}

	return child
}

func (p Problem) Evaluate(s Solution) float64 {
	switch p.Config.Objective {
	case ObjectiveSphere:
		return Sphere(s.X)
	case ObjectiveRastrigin:
		return Rastrigin(s.X)
	case ObjectiveRosenbrock:
		return Rosenbrock(s.X)
	default:
		panic("unsupported objective")
	}
}

func (p Problem) IsBetter(newScore, oldScore float64) bool {
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
