package main

import (
	"flag"
	"fmt"
	"love-letter-ai/gamemaster"
	"love-letter-ai/montecarlo"
	"love-letter-ai/players"
	"love-letter-ai/rules"
)

var loadPath = flag.String("load", "res/weights/mc.dat", "Path to the file to load weights")
var savePath = flag.String("save", "res/weights/mc.dat", "Path to the file to save weights")
var epsilon = flag.Float64("epsilon", 0.3, "Value of the starting epsilon")
var epsilonDecay = flag.Float64("epsilondecay", 0.7, "Factor for scaling epsilon after each training epoch")
var nEpochs = flag.Int("epochs", 5, "Number of epochs")
var nTraces = flag.Int("traces", 20, "Number of game traces to print after each epoch")
var nGames = flag.Int("games", 1000000000, "Number of games per training epoch")
var nTest = flag.Int("n", 1000, "Number of games played in each test against random")

func main() {
	flag.Parse()

	pl := montecarlo.NewQPlayer(float32(*epsilon))

	fmt.Println("Running vs random...")
	pl.TrainWithPlayerPolicy(*nGames, &players.RandomPlayer{})
	printTraces(*nTraces, pl)
	fightRandom(*nTest, pl)

	for j := 0; j < *nEpochs; j++ {
		*epsilon *= *epsilonDecay
		pl.SetEpsilon(float32(*epsilon))
		fmt.Printf("Running vs self %d...\n", j+1)
		pl.TrainWithSelfPolicy(*nGames)
		printTraces(*nTraces, pl)
		fightRandom(*nTest, pl)
	}

	pl.SetEpsilon(0.0)
	fmt.Printf("\n\nPlaying greedily...\n")
	printTraces(*nTraces, pl)
	fightRandom(*nTest, pl)

	if *savePath != "" {
		err := pl.SaveToFile(*savePath)
		if err != nil {
			panic(err)
		}
	}
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
