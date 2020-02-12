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

func TestZeroStateFullActionToState(t *testing.T) {
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
	actual := IndexWithoutAction(ActionIndex(seenCards, high, low, opponent, scoreDelta, action))
	expect := Index(seenCards, high, low, opponent, scoreDelta)
	assert.EqualValues(t, actual, expect)
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

func TestFullActionStateSimpleDeckToState(t *testing.T) {
	high, low, opponent := rules.Princess, rules.Countess, rules.King
	scoreDelta := -15
	action := rules.Action{
		PlayRecent:         true,
		TargetPlayerOffset: 0,
		SelectedCard:       rules.Princess,
	}
	actual := IndexWithoutAction(ActionIndex(rules.Deck{}, high, low, opponent, scoreDelta, action))
	expect := Index(rules.Deck{}, high, low, opponent, scoreDelta)
	assert.EqualValues(t, actual, expect)
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

func TestAllActionList(t *testing.T) {
	expected := [16]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	assert.EqualValues(t, expected, AllActionStates(0))
}

func TestAllActionListNonzeroState(t *testing.T) {
	expected := [16]int{0 + (19827491 << 4), 1 + (19827491 << 4), 2 + (19827491 << 4), 3 + (19827491 << 4), 4 + (19827491 << 4), 5 + (19827491 << 4), 6 + (19827491 << 4), 7 + (19827491 << 4), 8 + (19827491 << 4), 9 + (19827491 << 4), 10 + (19827491 << 4), 11 + (19827491 << 4), 12 + (19827491 << 4), 13 + (19827491 << 4), 14 + (19827491 << 4), 15 + (19827491 << 4)}
	assert.EqualValues(t, expected, AllActionStates(19827491))
}
