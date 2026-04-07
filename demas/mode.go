package main

type Mode string

const (
	ModeTSP        Mode = "tsp"
	ModeContinuous Mode = "continuous"
)

var AllModes = []Mode{
	ModeTSP,
	ModeContinuous,
}
