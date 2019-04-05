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

func (vf *ValueFunction) Update(pl players.Player, gamma float32) {
	states, winner, err := gamemaster.TraceOneGame(pl)
	if err != nil {
		panic(err.Error())
	}

	p1v, p2v := float32(1.0), float32(0.0)
	if winner != 0 {
		p1v, p2v = 0.0, 1.0
	}

	for _, s := range states {
		vf[s] += (p1v - vf[s]) * gamma
		p1v, p2v = p2v, p1v
	}
}
