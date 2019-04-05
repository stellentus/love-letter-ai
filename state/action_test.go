package state

import (
	"testing"

	"love-letter-ai/rules"

	"github.com/stretchr/testify/assert"
)

func BenchmarkActionIndex(b *testing.B) {
	seenCards := rules.DefaultDeck()
	high, low, opponent := rules.Priest, rules.Baron, rules.Handmaid
	scoreDelta := -8
	action := rules.Action{
		PlayRecent:         true,
		TargetPlayerOffset: 0,
		SelectedCard:       rules.Prince,
	}
	for n := 0; n < b.N; n++ {
		ActionIndex(seenCards, high, low, opponent, scoreDelta, action)
	}
}

func TestZeroActionState(t *testing.T) {
	seenCards := rules.DefaultDeck()
	for i := range seenCards {
		seenCards[i] = 0
	}
	high, low, opponent := rules.Guard, rules.Guard, rules.Guard
	scoreDelta := 0
	action := rules.Action{}
	assert.EqualValues(t, 0, ActionIndex(seenCards, high, low, opponent, scoreDelta, action))
}

func TestZeroStateFullAction(t *testing.T) {
	seenCards := rules.DefaultDeck()
	for i := range seenCards {
		seenCards[i] = 0
	}
	high, low, opponent := rules.Guard, rules.Guard, rules.Guard
	scoreDelta := 0
	action := rules.Action{
		PlayRecent:         true,
		TargetPlayerOffset: 0,
		SelectedCard:       rules.Princess,
	}
	assert.EqualValues(t, 15, ActionIndex(seenCards, high, low, opponent, scoreDelta, action))
}

func TestFullStateZeroAction(t *testing.T) {
	seenCards := rules.DefaultDeck()
	for i := range seenCards {
		seenCards[i] = 0
	}
	high, low, opponent := rules.Princess, rules.Countess, rules.King
	scoreDelta := -15
	action := rules.Action{}
	assert.EqualValues(t, 16317<<4, ActionIndex(seenCards, high, low, opponent, scoreDelta, action))
}

func TestFullActionStateSimpleDeck(t *testing.T) {
	high, low, opponent := rules.Princess, rules.Countess, rules.King
	scoreDelta := -15
	action := rules.Action{
		PlayRecent:         true,
		TargetPlayerOffset: 0,
		SelectedCard:       rules.Princess,
	}
	assert.EqualValues(t, 15+(16317<<4), ActionIndex(rules.Deck{}, high, low, opponent, scoreDelta, action))
}

func TestFullActionState(t *testing.T) {
	seenCards := rules.DefaultDeck()
	high, low, opponent := rules.Princess, rules.Princess, rules.Princess
	scoreDelta := -15
	action := rules.Action{
		PlayRecent:         true,
		TargetPlayerOffset: 0,
		SelectedCard:       rules.Princess,
	}
	assert.EqualValues(t, ActionSpaceMagnitude-1, ActionIndex(seenCards, high, low, opponent, scoreDelta, action))
}
