package state

import (
	"testing"

	"love-letter-ai/rules"

	"github.com/stretchr/testify/assert"
)

func BenchmarkActionIndex(b *testing.B) {
	seenCards := rules.DefaultDeck()
	recent, old, opponent := rules.Priest, rules.Baron, rules.Handmaid
	scoreDelta := -8
	action := rules.Action{
		PlayRecent:         true,
		TargetPlayerOffset: 0,
		SelectedCard:       rules.Prince,
	}
	for n := 0; n < b.N; n++ {
		ActionIndex(seenCards, recent, old, opponent, scoreDelta, action)
	}
}

var entireStateActionTests = []struct {
	recent, old, opponent rules.Card
	scoreDelta            int
	seenCards             rules.Deck
	action                rules.Action
	state                 int
	msg                   string
}{
	{rules.Guard, rules.Guard, rules.Guard, 0, rules.Deck{}, rules.Action{}, 0, "zero state zero action"},
	{rules.Guard, rules.Guard, rules.Guard, 0, rules.Deck{}, rules.Action{PlayRecent: true, TargetPlayerOffset: 0, SelectedCard: rules.Princess}, 15, "zero state full action"},
	{rules.Princess, rules.Countess, rules.King, -15, rules.Deck{}, rules.Action{}, 16317 << 4, "full state zero action"},
	{rules.Princess, rules.Countess, rules.King, -15, rules.Deck{}, rules.Action{PlayRecent: true, TargetPlayerOffset: 0, SelectedCard: rules.Princess}, 15 + (16317 << 4), "simple deck full action"},
	{rules.Princess, rules.Princess, rules.Princess, -15, rules.DefaultDeck(), rules.Action{PlayRecent: true, TargetPlayerOffset: 0, SelectedCard: rules.Princess}, ActionSpaceMagnitude - 1, "full state full action"},
}

func TestActionState(t *testing.T) {
	for _, test := range entireStateActionTests {
		assert.EqualValues(t, test.state, ActionIndex(test.seenCards, test.recent, test.old, test.opponent, test.scoreDelta, test.action), "Action state for "+test.msg)
	}
}

func TestActionStateToState(t *testing.T) {
	for _, test := range entireStateActionTests {
		actual := IndexWithoutAction(ActionIndex(test.seenCards, test.recent, test.old, test.opponent, test.scoreDelta, test.action))
		expect := Index(test.seenCards, test.recent, test.old, test.opponent, test.scoreDelta)
		assert.EqualValues(t, actual, expect, "Action state to state for "+test.msg)
	}
}

func TestAllActionList(t *testing.T) {
	expected := [16]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	assert.EqualValues(t, expected, AllActionStates(0))
}

func TestAllActionListNonzeroState(t *testing.T) {
	expected := [16]int{0 + (19827491 << 4), 1 + (19827491 << 4), 2 + (19827491 << 4), 3 + (19827491 << 4), 4 + (19827491 << 4), 5 + (19827491 << 4), 6 + (19827491 << 4), 7 + (19827491 << 4), 8 + (19827491 << 4), 9 + (19827491 << 4), 10 + (19827491 << 4), 11 + (19827491 << 4), 12 + (19827491 << 4), 13 + (19827491 << 4), 14 + (19827491 << 4), 15 + (19827491 << 4)}
	assert.EqualValues(t, expected, AllActionStates(19827491))
}
