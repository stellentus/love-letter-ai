package state

import "love-letter-ai/rules"

const (
	// SpaceMagnitude represents the number of possible states considered.
	// Note that some of these states are impossible to achieve in real gameplay.
	SpaceMagnitude     = deckSpaceMagnitude * 32 * 512 // deck*score*hand bits equals the size of the statespace
	deckSpaceMagnitude = 6 * 3 * 3 * 3 * 3 * 2 * 2 * 2
)

// Index returns a unique number between 0 and SpaceMagnitude-1 for the given state.
func Index(seenCards rules.Deck, recent, old, opponent rules.Card, scoreDelta int) int {
	return (((deckValue(seenCards) << 5) + scoreValue(scoreDelta)) << 9) + handValue(recent, old, opponent)
}

// 5 bits
func scoreValue(scoreDelta int) int {
	value := scoreDelta
	if scoreDelta < 0 {
		value = -scoreDelta
	}
	if value > 15 {
		value = 15
	}
	if scoreDelta < 0 {
		value |= 16
	}
	return value
}

// 12 bits
func deckValue(sc rules.Deck) int {
	return sc[rules.Guard] +
		6*(sc[rules.Priest]+
			3*(sc[rules.Baron]+
				3*(sc[rules.Handmaid]+
					3*(sc[rules.Prince]+
						3*(sc[rules.King]*4+
							sc[rules.Countess]*2+
							sc[rules.Princess])))))
}

// 9 bits
func handValue(recent, old, opponent rules.Card) int {
	recent--
	old--
	opponent--
	return int((((old << 3) + recent) << 3) + opponent)
}
