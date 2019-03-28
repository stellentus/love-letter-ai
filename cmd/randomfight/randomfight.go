package main

import (
	"fmt"
	"love-letter-ai/gamemaster"
	"love-letter-ai/players"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	pls := []players.Player{
		&players.RandomPlayer{},
		&players.RandomPlayer{},
	}
	gm, err := gamemaster.New(pls)
	if err != nil {
		panic(err)
	}
	winner, err := gm.PlaySeries(1000)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Player %d won %d-%d\n", winner, gm.Wins[winner], gm.Wins[(winner+1)%2])
}
