package players

import (
	"love-letter-ai/rules"
	"love-letter-ai/state"
	"math/rand"
)

type Player interface {
	PlayCard(state.Simple) rules.Action
}

type TrainingPlayer interface {
	Player

	// GreedyAction returns the greedy action for the given state, along with the corresponding state-action.
	// The action may be nil if nothing has trained yet.
	GreedyAction(state int) (*rules.Action, int)
}

// EpsilonGreedyAction provides a suggested action for the provided state.
// If it hasn't learned anything for this state, it plays randomly.
// It will also choose a random action with probability Epsilon. This isn't exactly
// Epsilon-greedy because it doesn't subtract the probability of the greedy action.
func EpsilonGreedyAction(pl TrainingPlayer, st state.Simple, epsilon float32) (rules.Action, int) {
	sNoAct := st.AsIndex()
	act, sa := pl.GreedyAction(sNoAct)
	if act == nil || rand.Float32() < epsilon {
		action := (&RandomPlayer{}).PlayCard(st)
		return action, state.IndexWithAction(sNoAct, action)
	}
	return *act, sa
}
