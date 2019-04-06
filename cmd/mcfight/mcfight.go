package main

import (
	"fmt"
	"love-letter-ai/gamemaster"
	"love-letter-ai/montecarlo"
	"love-letter-ai/players"
	"love-letter-ai/rules"
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
	fists := make([]rules.FinalState, 0, n)
	for i := 0; i < n; i++ {
		tr, err := gamemaster.TraceOneGame(&players.RandomPlayer{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Game %d winner: %d\n", i, tr.Winner)
		for _, v := range tr.StateInfos {
			fmt.Printf("    %08X: %0.3f\n", v.ActionState, pl.Value(v.ActionState))
		}
		fists = append(fists, tr.FinalState)
	}
	fmt.Println("Game | Discard | InHand | Opponent | Deck | Won? ")
	fmt.Println("-----|---------|--------|----------|------|-------")
	for i, fist := range fists {
		fmt.Printf(" %3d | %d       | %d      | %d        | %2d   | %t \n", i, fist.LastDiscard, fist.LastInHand, fist.OpponentInHand, fist.RemainingDeck, fist.DiscardWon)
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
