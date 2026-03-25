package main

type Problem[S any] interface {
	NewRandomSolution() S
	MutateSolution(S) S
	Evaluate(S) float64
	IsBetter(newScore, oldScore float64) bool
}
