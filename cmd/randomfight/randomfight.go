package main

import (
	"flag"
	"fmt"
	"love-letter-ai/gamemaster"
	"love-letter-ai/players"
	"math/rand"
	"time"
)

var nTest = flag.Int("n", 1000, "Number of games played in each test against random")

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	pls := []players.Player{
		&players.RandomPlayer{},
		&players.RandomPlayer{},
	}
	gm, err := gamemaster.New(pls)
	if err != nil {
		panic(err)
	}
	winner, err := gm.PlaySeries(*nTest)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Player %d won %d-%d\n", winner, gm.Wins[winner], gm.Wins[(winner+1)%2])
}
