package main

import (
	"flag"
	"fmt"
	"love-letter-ai/gamemaster"
	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/td"
	"os"
	"path/filepath"
)

var loadPath = flag.String("load", "", "Path to the file to load weights")
var savePath = flag.String("save", "", "Path to the file to save weights")
var gamma = flag.Float64("gamma", 1, "Value of the starting gamma")
var epsilon = flag.Float64("epsilon", 0.3, "Value of the starting epsilon")
var epsilonDecay = flag.Float64("epsilondecay", 0.7, "Factor for scaling epsilon after each training epoch")
var epsilonDecayPeriod = flag.Int("epsilondecayperiod", 100, "Number of training epochs between each epsilon adjustment")
var alpha = flag.Float64("*alpha", 0.3, "Value of the starting alpha")
var alphaDecay = flag.Float64("alphadecay", 0.995, "Factor for scaling alpha after each training epoch")
var nEpochs = flag.Int("epochs", 1000, "Number of epochs")
var nTraces = flag.Int("traces", 50, "Number of game traces to print after each epoch")
var nGames = flag.Int("games", 1000000000, "Number of games per training epoch")
var nTest = flag.Int("n", 10000, "Number of games played in each test against random")

func main() {
	flag.Parse()

	sar := td.NewTD(float32(*alpha), float32(*gamma))
	var err error

	if *loadPath != "" {
		err = sar.LoadFromFile(*loadPath)
		if err != nil {
			// Okay, no file, print a warning and keep going
			fmt.Println("WARNING: Could not find the file you wanted to load, so proceeding with newly initialized SARSA")
			sar = td.NewTD(float32(*alpha), float32(*gamma))
		} else {
			fmt.Println("The weights were loaded from '" + *loadPath + "'")
		}
	}

	if *savePath != "" {
		if _, err := os.Stat(filepath.Dir(*savePath)); os.IsNotExist(err) {
			panic("The path you plan to save at is a non-existent directory")
		}
		fmt.Println("The final weights will be saved at '" + *savePath + "'")
	}

	pls := []players.TrainingPlayer{
		sar.ExpectedSarsaLearner(),
		sar.ExpectedSarsaLearner(),
	}

	for j := 0; j < *nEpochs; j++ {
		fmt.Printf("Running vs self %d...\n", j+1)
		players.Train(pls, *nGames, *epsilon)

		fightRandom(*nTest, sar)

		*alpha *= *alphaDecay
		sar.Alpha = float32(*alpha)

		if (j % *epsilonDecayPeriod) == 0 {
			*epsilon *= *epsilonDecay
		}
	}

	fmt.Printf("\n\nPlaying greedily...\n")
	printTraces(*nTraces, sar)
	fightRandom(*nTest, sar)

	if *savePath != "" {
		err := sar.SaveToFile(*savePath)
		if err != nil {
			panic(err)
		}
	}
}

func printTraces(n int, sar *td.TD) {
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
	fmt.Printf("Sarsa win rates: %2.1f%%,", fightPlayers(n, []players.Player{
		pl,
		&players.RandomPlayer{},
	}))
	fmt.Printf(" %2.1f%%\n", 100.0-fightPlayers(n, []players.Player{
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
