package gamemaster

import (
	"fmt"
	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/state"
)

type Trace struct {
	States  []int
	Returns []float32
	Winner  int
}

// TraceOneGame returns the states for one gameplay played by the provided player pl.
func TraceOneGame(pl players.Player, gamma float32) (Trace, error) {
	sg, err := rules.NewGame(2)
	if err != nil {
		return Trace{}, err
	}

	tr := Trace{States: make([]int, 0, 15)}

	for !sg.GameEnded {
		s := sg.AsSimpleState()
		if s.OpponentCard == 0 {
			s.OpponentCard++
		}

		ss := state.Index(s.Discards, s.RecentDraw, s.OldCard, s.OpponentCard, s.ScoreDiff)
		if ss < 0 {
			return Trace{}, fmt.Errorf("Negative state was calculated: %d", ss)
		}
		tr.States = append(tr.States, ss)
		if err := sg.PlayCard(pl.PlayCard(s, sg.ActivePlayer)); err != nil {
			return Trace{}, fmt.Errorf("Game failed: %+v", sg)
		}
	}

	numPlays := len(tr.States)
	tr.Returns = make([]float32, numPlays)
	thisRet := float32(1.0)
	for i := numPlays - 1; i >= 0; i-- {
		// Leave ret[i] as zero unless this was the winner
		if i%2 == sg.Winner {
			tr.Returns[i] = thisRet
			thisRet *= gamma // only scale by gamma for actions taken by this player
		}
	}
	tr.Winner = sg.Winner

	return tr, nil
}
