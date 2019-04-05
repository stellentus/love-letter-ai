package montecarlo

import (
	"love-letter-ai/gamemaster"
	"love-letter-ai/players"
	"love-letter-ai/state"
)

type Value struct {
	sum   float32
	count int32
}

type ValueFunction [state.SpaceMagnitude]Value
type Action [state.SpaceMagnitude]uint8

func (vf *ValueFunction) Update(pl players.Player, gamma float32) {
	states, rets, _, err := gamemaster.TraceOneGame(pl, gamma)
	if err != nil {
		panic(err.Error())
	}

	for i, s := range states {
		vf[s].sum += rets[i]
		vf[s].count++
	}
}

func (vf *ValueFunction) Value(state int) float32 {
	return vf[state].sum / float32(vf[state].count)
}
