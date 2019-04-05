package main

import (
	"fmt"
	"love-letter-ai/gamemaster"
	"love-letter-ai/montecarlo"
	"love-letter-ai/players"
)

const (
	rounds = 1000000000
)

func main() {
	vf := montecarlo.ValueFunction{}
	pl := players.RandomPlayer{}

	fmt.Println("Running...")
	for i := 0; i < rounds; i++ {
		vf.Update(&pl)
	}

	for i := 0; i < 20; i++ {
		tr, err := gamemaster.TraceOneGame(&pl)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Winner:", tr.Winner)
		for _, v := range tr.States {
			fmt.Printf("    % 8d: %0.3f\n", v, vf.Value(v))
		}
	}
}
