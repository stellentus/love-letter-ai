package main

import (
	"fmt"
	"love-letter-ai/gamemaster"
	"love-letter-ai/montecarlo"
	"love-letter-ai/players"
)

const (
	rounds  = 1000000000
	epsilon = 0.05
)

func main() {
	pl := montecarlo.NewQPlayer(epsilon)
	fmt.Println("Running...")
	pl.TrainWithPlayerPolicy(rounds, &players.RandomPlayer{})

	for i := 0; i < 20; i++ {
		tr, err := gamemaster.TraceOneGame(&players.RandomPlayer{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Winner:", tr.Winner)
		for _, v := range tr.StateInfos {
			fmt.Printf("    % 8d: %0.3f\n", v.ActionState, pl.Value(v.ActionState))
		}
	}
}
