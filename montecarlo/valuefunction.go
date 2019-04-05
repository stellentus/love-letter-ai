package montecarlo

import (
	"love-letter-ai/gamemaster"
	"love-letter-ai/players"
	"love-letter-ai/state"
)

type ValueFunction [state.SpaceMagnitude]float32
type Action [state.SpaceMagnitude]uint8

func EvenValueFunction() ValueFunction {
	vf := ValueFunction{}
	for i := range vf {
		vf[i] = 0.5
	}
	return vf
}

func (vf *ValueFunction) Update(pl players.Player, gamma, valueScale float32) {
	states, rets, _, err := gamemaster.TraceOneGame(pl, gamma)
	if err != nil {
		panic(err.Error())
	}

	for i, s := range states {
		vf[s] += (rets[i] - vf[s]) * valueScale
	}
}
