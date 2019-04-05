package gamemaster

import (
	"fmt"
	"love-letter-ai/players"
	"love-letter-ai/rules"
)

type Trace struct {
	States       []int
	ActionStates []int
	Returns      []int
	Winner       int
}

// TraceOneGame returns the states for one gameplay played by the provided player pl.
func TraceOneGame(pl players.Player) (Trace, error) {
	sg, err := rules.NewGame(2)
	if err != nil {
		return Trace{}, err
	}

	tr := Trace{States: make([]int, 0, 15)}

	for !sg.GameEnded {
		s := players.NewSimpleState(sg)
		if s.OpponentCard == 0 {
			s.OpponentCard++
		}

		action := pl.PlayCard(s)
		sa, ss := s.AsIntWithAction(action)
		if ss < 0 || sa < 0 {
			return Trace{}, fmt.Errorf("Negative state was calculated: %d %d", ss, sa)
		}
		tr.States = append(tr.States, ss)
		tr.ActionStates = append(tr.ActionStates, sa)
		if err := sg.PlayCard(action); err != nil {
			return Trace{}, fmt.Errorf("Game failed: %+v", sg)
		}
	}

	numPlays := len(tr.States)
	tr.Returns = make([]int, numPlays)
	for i := numPlays - 1; i >= 0; i-- {
		// Leave ret[i] as zero unless this was the winner
		if i%2 == sg.Winner {
			tr.Returns[i] = 1
		}
	}
	tr.Winner = sg.Winner

	return tr, nil
}
