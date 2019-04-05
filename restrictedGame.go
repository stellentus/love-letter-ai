package main

import (
	"fmt"
	"love-letter-ai/gamemaster"
	"love-letter-ai/montecarlo"
	"love-letter-ai/players"
)

const (
	rounds = 1000000000
	gamma  = 0.05
)

func main() {
	vf := montecarlo.EvenValueFunction()
	pl := players.RandomPlayer{}

	fmt.Println("Running...")
	for i := 0; i < rounds; i++ {
		vf.Update(&pl, gamma)
	}

	for i := 0; i < 20; i++ {
		states, _, winner, err := gamemaster.TraceOneGame(&pl, gamma)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Winner:", winner)
		for _, v := range states {
			fmt.Printf("    % 8d: %0.3f\n", v, vf[v])
		}
	}
}
