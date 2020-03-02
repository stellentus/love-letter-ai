package state

import (
	"testing"

	"love-letter-ai/rules"

	"github.com/stretchr/testify/assert"
)

func BenchmarkIndex(b *testing.B) {
	seenCards := rules.DefaultDeck()
	recent, old, opponent := rules.Priest, rules.Baron, rules.Handmaid
	scoreDelta := -8
	for n := 0; n < b.N; n++ {
		Index(seenCards, recent, old, opponent, scoreDelta)
	}
}

var entireStateTests = []struct {
	recent, old, opponent rules.Card
	scoreDelta            int
	seenCards             rules.Deck
	state                 int
	msg                   string
}{
	{rules.Guard, rules.Guard, rules.Guard, 0, rules.Deck{}, 0, "zero state"},
	{rules.Guard, rules.Guard, rules.Guard, 0, rules.Deck{rules.Guard: 1}, 1 << (5 + 9), "minimal deck"},
	{rules.Princess, rules.Countess, rules.King, -15, rules.Deck{}, 16317, "full state simple deck"},
	{rules.Princess, rules.Princess, rules.Princess, -15, rules.DefaultDeck(), largestPossibleStateValue, "full state"},
}

func TestState(t *testing.T) {
	for _, test := range entireStateTests {
		assert.EqualValues(t, test.state, Index(test.seenCards, test.recent, test.old, test.opponent, test.scoreDelta), "State for "+test.msg)
	}
}

func TestStateInversion(t *testing.T) {
	for _, test := range entireStateTests {
		seenCards, recent, old, opponent, scoreDelta := FromIndex(test.state)

		assert.EqualValues(t, test.seenCards, seenCards, "State inversion for "+test.msg)
		assert.EqualValues(t, test.recent, recent, "State inversion for "+test.msg)
		assert.EqualValues(t, test.old, old, "State inversion for "+test.msg)
		assert.EqualValues(t, test.opponent, opponent, "State inversion for "+test.msg)
		assert.EqualValues(t, test.scoreDelta, scoreDelta, "State inversion for "+test.msg)
	}
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
