package main

import "math/rand/v2"

func CloneTour(tour []int) []int {
	cloned := make([]int, len(tour))
	copy(cloned, tour)
	return cloned
}

func MutateSwap(tour []int) []int {
	child := CloneTour(tour)

	i := rand.IntN(len(child))
	j := rand.IntN(len(child))

	child[i], child[j] = child[j], child[i]

	return child
}

func MutateInsert(tour []int) []int {
	if len(tour) < 2 {
		return CloneTour(tour)
	}

	from := rand.IntN(len(tour))
	to := rand.IntN(len(tour))

	for to == from {
		to = rand.IntN(len(tour))
	}

	moved := tour[from]

	without := make([]int, 0, len(tour)-1)
	without = append(without, tour[:from]...)
	without = append(without, tour[from+1:]...)

	if to > from {
		to--
	}

	child := make([]int, 0, len(tour))
	child = append(child, without[:to]...)
	child = append(child, moved)
	child = append(child, without[to:]...)

	return child
}
