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

func TestAsInt(t *testing.T) {
	assert.EqualValues(t, 0, Deck{}.AsInt(), "Empty")
	assert.EqualValues(t, 1, Deck{Guard: 1}.AsInt(), "One guard")
	assert.EqualValues(t, 3, Deck{Guard: 3}.AsInt(), "Three guards")
	assert.EqualValues(t, 8, Deck{Guard: 2, Priest: 1}.AsInt(), "Two guards and a priest walk into a bar...")
	assert.EqualValues(t, DeckSpaceMagnitude-1, DefaultDeck().AsInt(), "Full")
}
