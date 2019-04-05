package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandValues(t *testing.T) {
	tests := []struct {
		descr  string
		action Action
		out    int
	}{
		{"Empty action", Action{}, 0},                             // This should be the same as selecting Guard.
		{"Target Guard is empty", Action{SelectedCard: Guard}, 0}, // This is the true minimal input, but 'None' must also work
		{"Maximal action", Action{PlayRecent: true, SelectedCard: Princess}, 15},
		{"Maximal target", Action{PlayRecent: true, TargetPlayer: 1}, 3},
		{"Conflict input", Action{PlayRecent: true, TargetPlayer: 1, SelectedCard: Princess}, 3}, // when TargetPlayer and SelectedCard are both set, ignore SelectedCard
		{"Maximal target", Action{PlayRecent: false, SelectedCard: Princess}, 14},
		{"Target 1 equals Selected Priest", Action{PlayRecent: false, TargetPlayer: 1}, 2},
		{"Selected Priest equals Target 1", Action{PlayRecent: false, SelectedCard: Priest}, 2}, // note two states give the same output; that's okay because these occur when the played card is different, so these actions show up in different contexts
	}
	for _, test := range tests {
		assert.EqualValues(t, test.out, test.action.AsInt(), test.descr)
	}
}
