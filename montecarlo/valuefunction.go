package montecarlo

type ValueFunction [StateSpaceMagnitude]float32
type Action [StateSpaceMagnitude]uint8

func EvenValueFunction() ValueFunction {
	vf := ValueFunction{}
	for i := range vf {
		vf[i] = 0.5
	}
	return vf
}
