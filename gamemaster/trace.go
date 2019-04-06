package gamemaster

import (
	"fmt"
	"love-letter-ai/players"
	"love-letter-ai/rules"
)

type StateInfo struct {
	State       int
	ActionState int
	Won         bool
}

type Trace struct {
	StateInfos []StateInfo
	FinalState rules.FinalState
	Winner     int
}

// TraceOneGame returns the states for one gameplay played by the provided player pl.
func TraceOneGame(pl players.Player) (Trace, error) {
	sg, err := rules.NewGame(2)
	if err != nil {
		return Trace{}, err
	}

	tr := Trace{StateInfos: make([]StateInfo, 0, 15)}

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
		tr.StateInfos = append(tr.StateInfos, StateInfo{State: ss, ActionState: sa})
		if err := sg.PlayCard(action); err != nil {
			return Trace{}, fmt.Errorf("Game failed: %+v", sg)
		}
	}

	numPlays := len(tr.StateInfos)
	for i := numPlays - 1; i >= 0; i-- {
		// Leave ret[i] as zero unless this was the winner
		if i%2 == sg.Winner {
			tr.StateInfos[i].Won = true
		}
	}
	tr.Winner = sg.Winner
	tr.FinalState = sg.FinalState

	return tr, nil
}

// TraceIntoChannels calculates the states for one gameplay played by the provided player pl.
// It sends the results into the provided channels. The slice of channels should be of length 8.
func TraceIntoChannels(pl players.Player, chs []chan StateInfo) error {
	tr, err := TraceOneGame(pl)
	if err != nil {
		return err
	}

	for _, si := range tr.StateInfos {
		chs[si.State>>(12+5+9-3)] <- si
	}

	return nil
}
