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
	{rules.Guard, rules.Guard, rules.Guard, 0, rules.Deck{}, rules.Action{PlayRecent: true, TargetPlayerOffset: 0, SelectedCard: rules.Princess}, 15 << 26, "zero state full action"},
	{rules.Princess, rules.Countess, rules.King, -15, rules.Deck{}, rules.Action{}, 16317, "full state zero action"},
	{rules.Princess, rules.Countess, rules.King, -15, rules.Deck{}, rules.Action{PlayRecent: true, TargetPlayerOffset: 0, SelectedCard: rules.Princess}, (15 << 26) + 16317, "simple deck full action"},
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
	for i, val := range expected {
		expected[i] = val << 26
	}
	assert.EqualValues(t, expected, AllActionStates(0))
}

func TestAllActionListNonzeroState(t *testing.T) {
	expected := [16]int{(0 << 26) + 19827491, (1 << 26) + 19827491, (2 << 26) + 19827491, (3 << 26) + 19827491, (4 << 26) + 19827491, (5 << 26) + 19827491, (6 << 26) + 19827491, (7 << 26) + 19827491, (8 << 26) + 19827491, (9 << 26) + 19827491, (10 << 26) + 19827491, (11 << 26) + 19827491, (12 << 26) + 19827491, (13 << 26) + 19827491, (14 << 26) + 19827491, (15 << 26) + 19827491}
	assert.EqualValues(t, expected, AllActionStates(19827491))
}
