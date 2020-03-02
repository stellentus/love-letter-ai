package state

import "love-letter-ai/rules"

// SpaceMagnitude represents the number of possible states considered.
// Note that some of these states are impossible to achieve in real gameplay.
const SpaceMagnitude = rules.DeckSpaceMagnitude * 32 * 512 // deck*score*hand bits equals the size of the statespace

// TerminalState represents a state that can't regularly be reached.
// The value of `FromIndex(TerminalState)` is that the active player is holding 2 princess cards, the opponent also
// has a princess, the entire deck is in the discard pile, and the score delta is -15. So, obviously impossible.
const TerminalState = SpaceMagnitude - 1

// Index returns a unique number between 0 and SpaceMagnitude-1 for the given state.
func Index(seenCards rules.Deck, recent, old, opponent rules.Card, scoreDelta int) int {
	return (((seenCards.AsInt() << 5) + scoreValue(scoreDelta)) << 9) + handValue(recent, old, opponent)
}

func FromIndex(st int) (seenCards rules.Deck, recent, old, opponent rules.Card, scoreDelta int) {
	recent, old, opponent = handFromValue(st & 0x1FF)
	st >>= 9

	scoreDelta = scoreFromValue(st)
	st >>= 5

	seenCards.FromInt(st)

	return
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

func handFromValue(state int) (recent, old, opponent rules.Card) {
	opponent = rules.Card(state&0x7) + 1
	recent = rules.Card((state>>3)&0x7) + 1
	old = rules.Card((state>>6)&0x7) + 1
	return
}

// 9 bits
func handValue(recent, old, opponent rules.Card) int {
	recent--
	old--
	opponent--
	return int((((old << 3) + recent) << 3) + opponent)
}
