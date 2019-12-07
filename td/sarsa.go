package td

import "love-letter-ai/players"

func (td TD) SarsaLearner() players.TrainingPlayer {
	return sarsaLearner{TD: td}
}

type sarsaLearner struct{ TD }

func (sl sarsaLearner) UpdateQ(gameEnded bool, lastQ, sa int, reward float32) {
	thisValue := float32(0) // If game ended, the value of the new state is 0 because it's a terminal state
	if !gameEnded {
		thisValue = sl.Gamma * sl.qf[sa]
	}
	sl.qf[lastQ] += sl.Alpha * (reward + thisValue - sl.qf[lastQ])
}
