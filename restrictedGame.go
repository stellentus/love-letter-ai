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
		updateValueFunction(&vf, &pl)
	}

	for i := 0; i < 20; i++ {
		states, winner, err := gamemaster.TraceOneGame(&pl)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Winner:", winner)
		for _, v := range states {
			fmt.Printf("    % 8d: %0.3f\n", v, vf[v])
		}
	}
}

func updateValueFunction(vf *montecarlo.ValueFunction, pl players.Player) {
	states, winner, err := gamemaster.TraceOneGame(pl)
	if err != nil {
		panic(err.Error())
	}

	p1v, p2v := float32(1.0), float32(0.0)
	if winner != 0 {
		p1v, p2v = 0.0, 1.0
	}

	for _, s := range states {
		vf[s] += (p1v - vf[s]) * gamma
		p1v, p2v = p2v, p1v
	}
}
