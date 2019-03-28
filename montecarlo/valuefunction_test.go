package montecarlo

import (
	"testing"

	"bellstone.ca/go/love-letter-ai/rules"
	"github.com/stretchr/testify/assert"
)

func BenchmarkIndexOfState(b *testing.B) {
	seenCards := rules.DefaultDeck()
	high, low, opponent := rules.Priest, rules.Baron, rules.Handmaid
	scoreDelta := -8
	for n := 0; n < b.N; n++ {
		IndexOfState(seenCards, high, low, opponent, scoreDelta)
	}
}

func TestZeroState(t *testing.T) {
	seenCards := rules.DefaultDeck()
	for i := range seenCards {
		seenCards[i] = 0
	}
	high, low, opponent := rules.Guard, rules.Guard, rules.Guard
	scoreDelta := 0
	assert.EqualValues(t, 0, IndexOfState(seenCards, high, low, opponent, scoreDelta))
}

func TestFullStateSimpleDeck(t *testing.T) {
	high, low, opponent := rules.Princess, rules.Countess, rules.King
	scoreDelta := -15
	assert.EqualValues(t, 16317, IndexOfState(rules.Deck{}, high, low, opponent, scoreDelta))
}

func TestFullState(t *testing.T) {
	seenCards := rules.DefaultDeck()
	high, low, opponent := rules.Princess, rules.Princess, rules.Princess
	scoreDelta := -15
	assert.EqualValues(t, StateSpaceMagnitude-1, IndexOfState(seenCards, high, low, opponent, scoreDelta))
}

func TestScoreDelta(t *testing.T) {
	tests := []struct{ in, out int }{
		{0, 0},
		{3, 3},
		{15, 15},
		{16, 15},
		{2356, 15},
		{-3, 19},
		{-15, 31},
		{-16, 31},
		{-2356, 31},
	}
	for _, test := range tests {
		assert.EqualValues(t, test.out, scoreValue(test.in))
	}
}

func TestDeckValues(t *testing.T) {
	assert.EqualValues(t, 0, deckValue(rules.Deck{}), "Empty")
	assert.EqualValues(t, 3, deckValue(rules.Deck{rules.Guard: 3}), "Three guards")
	assert.EqualValues(t, 8, deckValue(rules.Deck{rules.Guard: 2, rules.Priest: 1}), "Two guards and a priest walk into a bar...")
	assert.EqualValues(t, deckSpaceMagnitude-1, deckValue(rules.DefaultDeck()), "Full")
}

func TestHandValues(t *testing.T) {
	assert.EqualValues(t, 0, handValue(rules.Guard, rules.Guard, rules.Guard), "Lowest")
	assert.EqualValues(t, 73, handValue(rules.Priest, rules.Priest, rules.Priest), "Priests")
	assert.EqualValues(t, 511, handValue(rules.Princess, rules.Princess, rules.Princess), "Max theoretically")
}
