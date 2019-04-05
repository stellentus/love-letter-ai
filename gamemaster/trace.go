package gamemaster

import (
	"fmt"
	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/state"
)

// TraceOneGame returns the states for one gameplay played by the provided player pl.
// The return values are a slice of states, the index of the winning player, and an error.
func TraceOneGame(pl players.Player) ([]int, int, error) {
	sg, err := rules.NewGame(2)
	if err != nil {
		return nil, 0, err
	}

	states := make([]int, 0, 15)
	for !sg.GameEnded {
		s := sg.AsSimpleState()
		if s.OpponentCard == 0 {
			s.OpponentCard++
		}

		ss := state.Index(s.Discards, s.RecentDraw, s.OldCard, s.OpponentCard, s.ScoreDiff)
		if ss < 0 {
			return nil, 0, fmt.Errorf("Negative state was calculated: %d", ss)
		}
		states = append(states, ss)
		if err := sg.PlayCard(pl.PlayCard(s, sg.ActivePlayer)); err != nil {
			fmt.Printf("Game failed: %+v\n", sg)
			panic(err)
		}
	}

	return states, sg.Winner, nil
}
