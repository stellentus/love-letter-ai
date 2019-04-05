package montecarlo

import (
	"fmt"
	"love-letter-ai/gamemaster"
	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/state"
	"math/rand"
)

type qPlayer struct {
	qf      []Value
	epsilon float32
}

func NewQPlayer(epsilon float32) *qPlayer {
	return &qPlayer{
		qf:      make([]Value, state.ActionSpaceMagnitude, state.ActionSpaceMagnitude),
		epsilon: epsilon,
	}
}

func (qp *qPlayer) TrainWithPlayerPolicy(episodes int, pl players.Player) {
	for i := 0; i < episodes; i++ {
		if (i % 100000) == 0 {
			fmt.Printf("% 2.2f%% complete\r", float32(i)/float32(episodes)*100)
		}

		tr, err := gamemaster.TraceOneGame(pl)
		if err != nil {
			panic(err.Error())
		}

		for _, si := range tr.StateInfos {
			qp.SaveState(si)
		}
	}
	fmt.Println("100.0% complete")
}

func (qp *qPlayer) TrainWithSelfPolicy(episodes int) {
	qp.TrainWithPlayerPolicy(episodes, qp)
}

// PlayCard provides a suggested action for the provided state.
// If it hasn't learned anything for this state, it plays randomly.
// It will also choose a random action with probability epsilon. This isn't exactly
// epsilon-greedy because it doesn't subtract the probability of the greedy action.
func (qp *qPlayer) PlayCard(state players.SimpleState) rules.Action {
	act := qp.policy(state.AsInt())
	if act == nil || rand.Float32() < qp.epsilon {
		return (&players.RandomPlayer{}).PlayCard(state)
	}
	return *act
}

func (qp qPlayer) Value(st int) float32 {
	return float32(qp.qf[st].sum) / float32(qp.qf[st].count)
}

// policy returns the greedy action for the given state. (Note the argument should be a state, not an action-state.)
// Ties are broken by choosing the first option (i.e. arbitrarily in a deterministic way).
func (qp qPlayer) policy(st int) *rules.Action {
	bestActs := []int{}
	bestActValue := float32(0)
	for act, actState := range state.AllActionStates(st) {
		thisVal := qp.Value(actState)
		if thisVal > bestActValue {
			bestActValue = thisVal
			bestActs = []int{act}
		} else if thisVal == bestActValue {
			bestActs = append(bestActs, act)
		}
	}
	if len(bestActs) > 0 {
		bestAct := rules.ActionFromInt(bestActs[0])
		return &bestAct
	} else {
		return nil
	}
}

func (qp *qPlayer) SaveState(si gamemaster.StateInfo) {
	s := si.ActionState
	if si.Won {
		qp.qf[s].sum++
	}
	qp.qf[s].count++

	// Check for overflow
	if qp.qf[s].count == 0xFFFF {
		qp.qf[s].sum /= 2
		qp.qf[s].count /= 2
	}
}
