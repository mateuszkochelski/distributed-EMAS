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

func isValidMode(mode Mode) bool {
	for _, m := range AllModes {
		if mode == m {
			return true
		}
	}
	return false
}
