package state

import "love-letter-ai/rules"

const (
	// ActionSpaceMagnitude represents the number of possible action-states considered.
	// Note that some of these states are impossible to achieve in real gameplay.
	ActionSpaceMagnitude = 16 * SpaceMagnitude // 16 actions
)

// ActionIndex returns a unique number between 0 and ActionSpaceMagnitude-1 for the given state.
func ActionIndex(seenCards rules.Deck, recent, old, opponent rules.Card, scoreDelta int, act rules.Action) int {
	return actionIndex(Index(seenCards, recent, old, opponent, scoreDelta), act)
}

// Indices returns the ActionIndex and the state Index.
func Indices(seenCards rules.Deck, recent, old, opponent rules.Card, scoreDelta int, act rules.Action) (int, int) {
	stateIndex := Index(seenCards, recent, old, opponent, scoreDelta)
	return actionIndex(stateIndex, act), stateIndex
}

func actionIndex(stateIndex int, act rules.Action) int {
	return (stateIndex << 4) + act.AsInt()
}
