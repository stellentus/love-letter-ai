package main

import (
	"fmt"
	"love-letter-ai/montecarlo"
	"love-letter-ai/players"
	"love-letter-ai/rules"
)

const (
	rounds = 1000000000
	gamma  = 0.05
)

func main() {
	vf := montecarlo.EvenValueFunction()

	fmt.Println("Running...")
	for i := 0; i < rounds; i++ {
		updateValueFunction(&vf)
	}

	for i := 0; i < 20; i++ {
		states, winner := traceGame()
		fmt.Println("Winner:", winner)
		for _, v := range states {
			fmt.Printf("    % 8d: %0.3f\n", v, vf[v])
		}
	}
}

func traceGame() ([]int, int) {
	sg, err := rules.NewGame(2)
	if err != nil {
		panic(err)
	}
	p := players.RandomPlayer{}

	states := make([]int, 0, 15)
	for !sg.GameEnded {
		s := sg.AsSimpleState()
		if s.OpponentCard == 0 {
			s.OpponentCard++
		}

		ss := montecarlo.IndexOfState(s.Discards, s.RecentDraw, s.OldCard, s.OpponentCard, s.ScoreDiff)
		if ss < 0 {
			panic(fmt.Sprintf("Negative state was calculated: %d", ss))
		}
		states = append(states, ss)
		if err := sg.PlayCard(p.PlayCard(s, sg.ActivePlayer)); err != nil {
			fmt.Printf("Game failed: %+v\n", sg)
			panic(err)
		}
	}

	return states, sg.Winner
}

func updateValueFunction(vf *montecarlo.ValueFunction) {
	states, winner := traceGame()

	p1v, p2v := float32(1.0), float32(0.0)
	if winner != 0 {
		p1v, p2v = 0.0, 1.0
	}

	for _, s := range states {
		vf[s] += (p1v - vf[s]) * gamma
		p1v, p2v = p2v, p1v
	}
}
