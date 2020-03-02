package state

import (
	"testing"

	"love-letter-ai/rules"

	"github.com/stretchr/testify/assert"
)

func BenchmarkIndex(b *testing.B) {
	seenCards := rules.DefaultDeck()
	high, low, opponent := rules.Priest, rules.Baron, rules.Handmaid
	scoreDelta := -8
	for n := 0; n < b.N; n++ {
		Index(seenCards, high, low, opponent, scoreDelta)
	}
}

func TestZeroState(t *testing.T) {
	seenCards := rules.DefaultDeck()
	for i := range seenCards {
		seenCards[i] = 0
	}
	high, low, opponent := rules.Guard, rules.Guard, rules.Guard
	scoreDelta := 0
	assert.EqualValues(t, 0, Index(seenCards, high, low, opponent, scoreDelta))
}

func TestMinimalDeck(t *testing.T) {
	high, low, opponent := rules.Guard, rules.Guard, rules.Guard
	scoreDelta := 0
	deck := rules.Deck{rules.Guard: 1}
	assert.EqualValues(t, 1<<(5+9), Index(deck, high, low, opponent, scoreDelta))
}

func TestFullStateSimpleDeck(t *testing.T) {
	high, low, opponent := rules.Princess, rules.Countess, rules.King
	scoreDelta := -15
	assert.EqualValues(t, 16317, Index(rules.Deck{}, high, low, opponent, scoreDelta))
}

func TestFullState(t *testing.T) {
	seenCards := rules.DefaultDeck()
	high, low, opponent := rules.Princess, rules.Princess, rules.Princess
	scoreDelta := -15
	assert.EqualValues(t, SpaceMagnitude-1, Index(seenCards, high, low, opponent, scoreDelta))
}

var scoreDeltaTests = []struct{ score, state int }{
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

func TestScoreDelta(t *testing.T) {
	for _, test := range scoreDeltaTests {
		assert.EqualValues(t, test.state, scoreValue(test.score))
	}
}

func TestReverseScoreDelta(t *testing.T) {
	for _, test := range scoreDeltaTests {
		score := test.score
		// Reversing has smaller range
		if score > 15 {
			score = 15
		}
		if score < -15 {
			score = -15
		}
		assert.EqualValues(t, score, scoreFromValue(test.state))
	}
}

var handValuesTests = []struct {
	recent, old, opponent rules.Card
	state                 int
	msg                   string
}{
	{rules.Guard, rules.Guard, rules.Guard, 0, "Lowest"},
	{rules.Priest, rules.Priest, rules.Priest, 73, "Priests"},
	{rules.Princess, rules.Princess, rules.Princess, 511, "Max theoretically"},
}

func TestHandValues(t *testing.T) {
	for _, test := range handValuesTests {
		assert.EqualValues(t, test.state, handValue(test.recent, test.old, test.opponent), test.msg)
	}
}

func TestReverseHandValues(t *testing.T) {
	for _, test := range handValuesTests {
		recent, old, opponent := handFromValue(test.state)
		assert.EqualValues(t, test.recent, recent, "Reverse recent with "+test.msg)
		assert.EqualValues(t, test.old, old, "Reverse old with "+test.msg)
		assert.EqualValues(t, test.opponent, opponent, "Reverse opponent with "+test.msg)
	}
}
