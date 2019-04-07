package main

import (
	"fmt"
	"love-letter-ai/gamemaster"
	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/td"
)

const (
	rounds  = 1000000000
	gamma   = 1
	epsilon = 0.3
)

func main() {
	alpha := float32(0.3)
	sar := td.NewSarsa(epsilon, alpha, gamma)
	pl := sar.NewPlayer()

	for j := 0; j < 6; j++ {
		fmt.Printf("Running vs self %d...\n", j+1)
		sar.Train(rounds)

		printTraces(20, pl, sar)
		fightRandom(1000, pl)

		alpha *= 0.7
		sar.Alpha = alpha
	}

	fmt.Printf("\n\nPlaying greedily...\n")
	printTraces(50, pl, sar)
	fightRandom(1000, pl)
}

func printTraces(n int, pl players.Player, sar *td.Sarsa) {
	fists := make([]rules.FinalState, 0, n)
	for i := 0; i < n; i++ {
		tr, err := gamemaster.TraceOneGame(&players.RandomPlayer{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Game %d winner: %d\n", i, tr.Winner)
		for _, v := range tr.StateInfos {
			fmt.Printf("    %08X: %0.3f\n", v.ActionState, sar.Value(v.ActionState))
		}
		fists = append(fists, tr.FinalState)
	}
	fmt.Println("Game | Discard | InHand | Opponent | Deck | Won? ")
	fmt.Println("-----|---------|--------|----------|------|-------")
	for i, fist := range fists {
		fmt.Printf(" %3d | %d       | %d      | %d        | %2d   | %t \n", i, fist.LastDiscard, fist.LastInHand, fist.OpponentInHand, fist.RemainingDeck, fist.DiscardWon)
	}
}

func fightRandom(n int, pl players.Player) {
	fmt.Printf("MC playing 1st has a win rate of %2.1f%%\n", fightPlayers(n, []players.Player{
		pl,
		&players.RandomPlayer{},
	}))
	fmt.Printf("MC playing 2nd has a win rate of %2.1f%%\n", 100.0-fightPlayers(n, []players.Player{
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