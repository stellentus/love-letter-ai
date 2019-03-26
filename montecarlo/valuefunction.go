package montecarlo

import "love-letter-ai/rules"

const (
	// StateSpaceMagnitude represents the number of possible states considered.
	// Note that some of these states are impossible to achieve in real gameplay.
	StateSpaceMagnitude = deckSpaceMagnitude * 32 * 512 // deck*score*hand bits equals the size of the statespace
	deckSpaceMagnitude  = 6 * 3 * 3 * 3 * 3 * 2 * 2 * 2
)

type ValueFunction [StateSpaceMagnitude]float32

// Index state returns a unique number between 0 and StateSpaceMagnitude-1 for the given state.
func IndexOfState(seenCards rules.Deck, high, low, opponent rules.Card, scoreDelta int) int {
	return (((deckValue(seenCards) << 5) + scoreValue(scoreDelta)) << 9) + handValue(high, low, opponent)
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
func handValue(high, low, opponent rules.Card) int {
	high--
	low--
	opponent--
	return int((((low << 3) + high) << 3) + opponent)
}
