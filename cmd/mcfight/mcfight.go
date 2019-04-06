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
	epsilon := float32(0.3)
	pl := montecarlo.NewQPlayer(epsilon)

	fmt.Println("Running vs random...")
	pl.TrainWithPlayerPolicy(rounds, &players.RandomPlayer{})
	printTraces(20, pl)
	fightRandom(1000, pl)

	for j := 0; j < 5; j++ {
		epsilon *= 0.7
		pl.SetEpsilon(epsilon)
		fmt.Printf("Running vs self %d...\n", j+1)
		pl.TrainWithSelfPolicy(rounds)
		printTraces(20, pl)
		fightRandom(1000, pl)
	}

	pl.SetEpsilon(0.0)
	fmt.Printf("Playing greedily...\n")
	printTraces(20, pl)
	fightRandom(1000, pl)
}

func printTraces(n int, pl *montecarlo.QPlayer) {
	for i := 0; i < n; i++ {
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

func fightRandom(n int, pl *montecarlo.QPlayer) {
	fmt.Printf("MC is 0 with a win rate of % 2.1f%%\n", fightPlayers(n, []players.Player{
		pl,
		&players.RandomPlayer{},
	}))
	fmt.Printf("MC is 1 with a win rate of % 2.1f%%\n", 100.0-fightPlayers(n, []players.Player{
		&players.RandomPlayer{},
		pl,
	}))
}

func fightPlayers(n int, pls []players.Player) float32 {
	// Now fight vs Random
	gm, err := gamemaster.New(pls)
	if err != nil {
		panic(err)
	}
	wins, err := gm.PlayStatistics(n)
	if err != nil {
		panic(err)
	}

	return float32(wins) / float32(n) * 100.0
}
