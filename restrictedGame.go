package main

import (
	"fmt"
	"love-letter-ai/montecarlo"
	"love-letter-ai/players"
	"love-letter-ai/rules"
)

const (
	rounds = 1000000
	gamma  = 0.02
)

var simpleDeck = rules.Deck{
	rules.Guard:    3,
	rules.Handmaid: 2,
	rules.Prince:   2,
	rules.Princess: 1,
}

func main() {
	vf := montecarlo.EvenValueFunction()

	fmt.Println("Running...")
	for i := 0; i < rounds; i++ {
		updateValueFunction(&vf)
	}
	for _, v := range []int{
		24,
		895011,
		3564548,
		256,
		25280,
		33056,
	} {
		if v < 16384 {
			fmt.Println()
		}
		fmt.Printf("% 8d: %0.3f\n", v, vf[v])
	}
}

func updateValueFunction(vf *montecarlo.ValueFunction) {
	sg := rules.NewSimpleGame(simpleDeck)
	p := players.RandomPlayer{}

	states := make([]int, 0, 8)
	for !sg.GameEnded {
		s := sg.AsSimpleState()
		if s.OpponentCard == 0 {
			s.OpponentCard++
		}

		ss := montecarlo.IndexOfState(s.Discards, s.RecentDraw, s.OldCard, s.OpponentCard, s.ScoreDiff)
		if ss < 0 {
			panic("Huh")
		}
		states = append(states, ss)
		if err := sg.PlayCard(p.PlayCard(s, sg.ActivePlayer)); err != nil {
			fmt.Printf("Game failed: %+v\n", sg)
			panic(err)
		}
	}

	p1v, p2v := float32(1.0), float32(0.0)
	if sg.Winner != 0 {
		p1v, p2v = 0.0, 1.0
	}

	for _, s := range states {
		vf[s] += (p1v - vf[s]) * gamma
		p1v, p2v = p2v, p1v
	}
}
