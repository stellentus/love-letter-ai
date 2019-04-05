package state

import "love-letter-ai/rules"

const (
	// ActionSpaceMagnitude represents the number of possible action-states considered.
	// Note that some of these states are impossible to achieve in real gameplay.
	ActionSpaceMagnitude = 16 * SpaceMagnitude // 16 actions
)

// ActionIndex returns a unique number between 0 and ActionSpaceMagnitude-1 for the given state.
func ActionIndex(seenCards rules.Deck, recent, old, opponent rules.Card, scoreDelta int, act rules.Action) int {
	return (Index(seenCards, recent, old, opponent, scoreDelta) << 4) + act.AsInt()
}
