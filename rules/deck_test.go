package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDrawingCard(t *testing.T) {
	deck := Deck{Guard: 1}
	card := deck.Draw()
	assert.EqualValues(t, card, Guard)
}

func TestDrawingEmpty(t *testing.T) {
	deck := Deck{}
	card := deck.Draw()
	assert.EqualValues(t, card, None)
}

func TestDrawingCardUntilEmpty(t *testing.T) {
	deck := Deck{Guard: 1}
	card := deck.Draw()
	card = deck.Draw()
	assert.EqualValues(t, card, None)
}
