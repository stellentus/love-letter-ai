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

var deckIntTests = []struct {
	descr string
	asInt int
	deck  Deck
}{
	{"Empty", 0, Deck{}},
	{"One guard", 1, Deck{Guard: 1}},
	{"One King", 1944, Deck{King: 1}},
	{"One Countess", 972, Deck{Countess: 1}},
	{"One Princess", 486, Deck{Princess: 1}},
	{"Three guards", 3, Deck{Guard: 3}},
	{"Two guards and a priest walk into a bar...", 8, Deck{Guard: 2, Priest: 1}},
	{"Full", DeckSpaceMagnitude - 1, DefaultDeck()},
}

func TestAsInt(t *testing.T) {
	for _, test := range deckIntTests {
		assert.EqualValues(t, test.asInt, test.deck.AsInt(), "To int: "+test.descr)
	}
}

func TestFromInt(t *testing.T) {
	for _, test := range deckIntTests {
		deck := Deck{}
		deck.FromInt(test.asInt)
		assert.EqualValues(t, test.deck, deck, "From int: "+test.descr)
	}
}
