package montecarlo

import (
	"love-letter-ai/gamemaster"
	"love-letter-ai/players"
	"love-letter-ai/state"
)

type Value struct {
	// The maximum visits to a state isn't likely to be above 4 billion, so this is fine.
	sum   uint32
	count uint32
}

type ValueFunction [state.SpaceMagnitude]Value
type Action [state.SpaceMagnitude]uint8

func (vf *ValueFunction) Update(pl players.Player) {
	tr, err := gamemaster.TraceOneGame(pl)
	if err != nil {
		panic(err.Error())
	}

	for i, s := range tr.States {
		vf[s].sum += uint32(tr.Returns[i])
		vf[s].count++
	}
}

func (vf *ValueFunction) Value(state int) float32 {
	return float32(vf[state].sum) / float32(vf[state].count)
}
