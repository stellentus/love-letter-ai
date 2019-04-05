package montecarlo

import "love-letter-ai/state"

type ValueFunction [state.SpaceMagnitude]float32
type Action [state.SpaceMagnitude]uint8

func EvenValueFunction() ValueFunction {
	vf := ValueFunction{}
	for i := range vf {
		vf[i] = 0.5
	}
	return vf
}
