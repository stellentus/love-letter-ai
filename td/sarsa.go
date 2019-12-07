package td

import "love-letter-ai/players"
import "love-letter-ai/state"

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

func (td TD) ExpectedSarsaLearner() players.TrainingPlayer {
	return expectedSarsaLearner{TD: td}
}

type expectedSarsaLearner struct{ TD }

func (lrn expectedSarsaLearner) UpdateQ(gameEnded bool, lastQ, sa int, reward float32) {
	// The expected value is the greedy policy.
	thisValue := float32(0) // If game ended, the value of the new state is 0 because it's a terminal state
	if !gameEnded {
		st := state.IndexWithoutAction(sa)
		act, greedySA := lrn.GreedyAction(st)
		if act == nil {
			// We don't have enough data to know what's greedy. I'm not sure if this is common or impossible.
			greedySA = sa
		}
		thisValue = lrn.Gamma * lrn.qf[greedySA]
	}
	lrn.qf[lastQ] += lrn.Alpha * (reward + thisValue - lrn.qf[lastQ])
}
