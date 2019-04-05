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
		if (i % 100000) == 0 {
			fmt.Printf("% 2.2f%% complete\r", float32(i)/rounds*100)
		}
		vf.Update(&pl)
	}
	fmt.Println("100.0%% complete")

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
