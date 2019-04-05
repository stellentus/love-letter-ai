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

	fmt.Println("Running vs random...")
	pl.TrainWithPlayerPolicy(rounds, &players.RandomPlayer{})
	printTraces(20, pl)
	fightRandom(1000, pl)

	for j := 0; j < 5; j++ {
		fmt.Printf("Running vs self %d...\n", j+1)
		pl.TrainWithSelfPolicy(rounds)
		printTraces(20, pl)
		fightRandom(1000, pl)
	}
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
	fmt.Printf("MC is 0; ")
	fightPlayers(n, []players.Player{
		pl,
		&players.RandomPlayer{},
	})
	fmt.Printf("MC is 1; ")
	fightPlayers(n, []players.Player{
		&players.RandomPlayer{},
		pl,
	})
}

func fightPlayers(n int, pls []players.Player) {
	// Now fight vs Random
	gm, err := gamemaster.New(pls)
	if err != nil {
		panic(err)
	}
	winner, err := gm.PlaySeries(n)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Player %d won %d-%d\n", winner, gm.Wins[winner], gm.Wins[(winner+1)%2])
}
