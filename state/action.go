package state

import "love-letter-ai/rules"

const (
	// ActionSpaceMagnitude represents the number of possible action-states considered.
	// Note that some of these states are impossible to achieve in real gameplay.
	ActionSpaceMagnitude = 16 * SpaceMagnitude // 16 actions
)

// ActionIndex returns a unique number between 0 and ActionSpaceMagnitude-1 for the given state.
func ActionIndex(seenCards rules.Deck, recent, old, opponent rules.Card, scoreDelta int, act rules.Action) int {
	return IndexWithAction(Index(seenCards, recent, old, opponent, scoreDelta), act)
}

// Indices returns the ActionIndex and the state Index.
func Indices(seenCards rules.Deck, recent, old, opponent rules.Card, scoreDelta int, act rules.Action) (int, int) {
	stateIndex := Index(seenCards, recent, old, opponent, scoreDelta)
	return IndexWithAction(stateIndex, act), stateIndex
}

func IndexWithAction(stateIndex int, act rules.Action) int {
	return (stateIndex << 4) + act.AsInt()
}

func IndexWithoutAction(actionIndex int) int {
	return actionIndex >> 4
}

func ActionFromIndex(actionIndex int) int {
	return actionIndex&0xF
}

// AllActionStates returns all possible ActionStates for a given state.
// It does not eliminate actions that are impossible for the given state, so it always returns 16 options.
func AllActionStates(state int) [16]int {
	states := [16]int{}
	baseState := state << 4
	for i := 0; i < 16; i++ {
		states[i] = baseState + i
	}
	return states
}
