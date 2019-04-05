package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	descr     string
	action    Action
	out       int
	converted Action
}{
	// An empty action should be the same as selecting a Guard.
	{"Empty action", Action{}, 0, Action{SelectedCard: Guard}},
	{"Target Guard is empty", Action{SelectedCard: Guard}, 0, Action{SelectedCard: Guard}},

	// Note these two states give the same output; that's okay because these occur when the played card is different, so these actions show up in different contexts.
	// Also note that anytime there is a SelectedCard, the TargetPlayerOffset ought to be 1.
	{"Target 1 equals Selected Priest", Action{PlayRecent: false, TargetPlayerOffset: 1}, 2, Action{PlayRecent: false, TargetPlayerOffset: 1, SelectedCard: Priest}},
	{"Selected Priest equals Target 1", Action{PlayRecent: false, SelectedCard: Priest}, 2, Action{PlayRecent: false, TargetPlayerOffset: 1, SelectedCard: Priest}},

	// Test a few different actions
	{"Maximal action", Action{PlayRecent: true, SelectedCard: Princess}, 15, Action{PlayRecent: true, TargetPlayerOffset: 1, SelectedCard: Princess}},
	{"Target player", Action{PlayRecent: true, TargetPlayerOffset: 1}, 3, Action{PlayRecent: true, TargetPlayerOffset: 1, SelectedCard: Priest}},
	{"Nearly maximal", Action{PlayRecent: false, SelectedCard: Princess}, 14, Action{PlayRecent: false, TargetPlayerOffset: 1, SelectedCard: Princess}},

	// When TargetPlayerOffset and SelectedCard are both set, ignore TargetPlayerOffset for conversion (since this is a guard targeting a specific card).
	{"Conflict input", Action{PlayRecent: true, TargetPlayerOffset: 1, SelectedCard: Princess}, 15, Action{PlayRecent: true, TargetPlayerOffset: 1, SelectedCard: Princess}},
}

func TestHandValues(t *testing.T) {
	for _, test := range tests {
		assert.EqualValues(t, test.out, test.action.AsInt(), "To int: "+test.descr)
	}
}

func TestActionConversion(t *testing.T) {
	for _, test := range tests {
		assert.EqualValues(t, test.converted, ActionFromInt(test.action.AsInt()), "Convert: "+test.descr)
	}
}
