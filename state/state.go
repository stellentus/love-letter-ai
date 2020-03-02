package state

import "love-letter-ai/rules"

// SpaceMagnitude represents the number of possible states considered.
// Note that some of these states are impossible to achieve in real gameplay.
const SpaceMagnitude = rules.DeckSpaceMagnitude * 32 * 512 // deck*score*hand bits equals the size of the statespace

// Index returns a unique number between 0 and SpaceMagnitude-1 for the given state.
func Index(seenCards rules.Deck, recent, old, opponent rules.Card, scoreDelta int) int {
	return (((seenCards.AsInt() << 5) + scoreValue(scoreDelta)) << 9) + handValue(recent, old, opponent)
}

// expects state value in lower 5 bits
// a delta greater than 15 is just 15
func scoreFromValue(state int) int {
	scoreDelta := state & 0xF
	if state&0x10 == 0x10 {
		// Sign bit was set
		scoreDelta *= -1
	}
	return scoreDelta
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

// 9 bits
func handValue(recent, old, opponent rules.Card) int {
	recent--
	old--
	opponent--
	return int((((old << 3) + recent) << 3) + opponent)
}
