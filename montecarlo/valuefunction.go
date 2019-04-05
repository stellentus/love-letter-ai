package montecarlo

import (
	"love-letter-ai/gamemaster"
	"love-letter-ai/players"
	"love-letter-ai/state"
)

type Value struct {
	// The maximum visits to a state isn't likely to be much above 65k, so let's save some RAM.
	// It's necessary to deal with overflow, but we don't care about loss of precision from occasionally dividing by two.
	sum   uint16
	count uint16
}

type ValueFunction [state.SpaceMagnitude]Value
type Action [state.SpaceMagnitude]uint8

func (vf *ValueFunction) Update(pl players.Player) {
	tr, err := gamemaster.TraceOneGame(pl)
	if err != nil {
		panic(err.Error())
	}

	for _, si := range tr.StateInfos {
		s := si.State
		if si.Won {
			vf[s].sum++
		}
		vf[s].count++

		// Check for overflow
		if vf[s].count == 0xFFFF {
			vf[s].sum /= 2
			vf[s].count /= 2
		}
	}
}

func (vf *ValueFunction) Value(state int) float32 {
	return float32(vf[state].sum) / float32(vf[state].count)
}
