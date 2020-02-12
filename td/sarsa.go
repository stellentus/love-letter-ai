package td

import (
	"math/rand"

	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/state"
)

func (td TD) SarsaLearner() players.TrainingPlayer {
	return sarsaLearner{TD: td}
}

type sarsaLearner struct{ TD }

func (lrn sarsaLearner) Finalize() {}

func (sl sarsaLearner) UpdateQ(gameEnded bool, qStates []int, reward float32) {
	lastQ, sa := qStates[len(qStates)-2], qStates[len(qStates)-1]

	thisValue := float32(0) // If game ended, the value of the new state is 0 because it's a terminal state
	if !gameEnded {
		thisValue = sl.Gamma * sl.qf[sa]
	}
	sl.qf[lastQ] += sl.Alpha * (reward + thisValue - sl.qf[lastQ])
}

func (td TD) QLearner() players.TrainingPlayer {
	return qLearner{TD: td}
}

type qLearner struct{ TD }

func (lrn qLearner) Finalize() {}

func (lrn qLearner) UpdateQ(gameEnded bool, qStates []int, reward float32) {
	lastQ, sa := qStates[len(qStates)-2], qStates[len(qStates)-1]

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

func (td TD) DoubleQLearner() players.TrainingPlayer {
	return doubleQLearner{
		td: []TD{td, *NewTD(td.Alpha, td.Gamma)},
	}
}

type doubleQLearner struct {
	td          []TD
	selectFirst bool
}

func (lrn doubleQLearner) Finalize() {
	for i := range lrn.td[0].qf {
		lrn.td[0].qf[i] += lrn.td[1].qf[i]
	}
}

func (lrn doubleQLearner) randTD() int                 { return rand.Int() % 2 }
func (lrn doubleQLearner) randTDRand(r *rand.Rand) int { return r.Int() % 2 }

func (lrn doubleQLearner) PlayCard(st state.Simple) rules.Action {
	return lrn.td[lrn.randTD()].PlayCard(st)
}

func (lrn doubleQLearner) PlayCardRand(st state.Simple, r *rand.Rand) rules.Action {
	return lrn.td[lrn.randTDRand(r)].PlayCard(st)
}

func (lrn doubleQLearner) GreedyAction(state int) (*rules.Action, int) {
	return lrn.td[lrn.randTD()].GreedyAction(state)
}

func (lrn doubleQLearner) UpdateQ(gameEnded bool, qStates []int, reward float32) {
	lastQ, sa := qStates[len(qStates)-2], qStates[len(qStates)-1]

	pick := lrn.randTD()
	lrn.updateQ(lrn.td[pick], lrn.td[(pick+1)%2], gameEnded, lastQ, sa, reward)
}

func (lrn doubleQLearner) updateQ(a, b TD, gameEnded bool, lastQ, sa int, reward float32) {
	// The expected value is the greedy policy.
	thisValue := float32(0) // If game ended, the value of the new state is 0 because it's a terminal state
	if !gameEnded {
		st := state.IndexWithoutAction(sa)
		act, greedySA := a.GreedyAction(st)
		if act == nil {
			// We don't have enough data to know what's greedy. I'm not sure if this is common or impossible.
			greedySA = sa
		}
		thisValue = a.Gamma * b.qf[greedySA]
	}
	a.qf[lastQ] += a.Alpha * (reward + thisValue - a.qf[lastQ])
}
