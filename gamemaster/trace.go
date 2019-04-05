package gamemaster

import (
	"fmt"
	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/state"
)

// TraceOneGame returns the states for one gameplay played by the provided player pl.
// The return values are a slice of states, a slice of returns, the index of the winning player, and an error.
func TraceOneGame(pl players.Player, gamma float32) ([]int, []float32, int, error) {
	sg, err := rules.NewGame(2)
	if err != nil {
		return nil, nil, 0, err
	}

	states := make([]int, 0, 15)
	for !sg.GameEnded {
		s := sg.AsSimpleState()
		if s.OpponentCard == 0 {
			s.OpponentCard++
		}

		ss := state.Index(s.Discards, s.RecentDraw, s.OldCard, s.OpponentCard, s.ScoreDiff)
		if ss < 0 {
			return nil, nil, 0, fmt.Errorf("Negative state was calculated: %d", ss)
		}
		states = append(states, ss)
		if err := sg.PlayCard(pl.PlayCard(s, sg.ActivePlayer)); err != nil {
			fmt.Printf("Game failed: %+v\n", sg)
			panic(err)
		}
	}

	numPlays := len(states)
	rets := make([]float32, numPlays)
	thisRet := float32(1.0)
	for i := numPlays - 1; i >= 0; i-- {
		// Leave ret[i] as zero unless this was the winner
		if i%2 == sg.Winner {
			rets[i] = thisRet
			thisRet *= gamma // only scale by gamma for actions taken by this player
		}
	}

	return states, rets, sg.Winner, nil
}
