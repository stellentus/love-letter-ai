package rules

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	_, err := NewGame(2)
	assert.NoError(t, err)
}

func TestPlayingPrince(t *testing.T) {
	rand.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Countess
	state.ActivePlayerCard = Guard

	// Play the Prince on the other player
	err := state.PlayCard(Action{
		PlayRecent:   false,
		TargetPlayer: 1,
		SelectedCard: None,
	})
	assert.NoError(t, err)

	assert.Equal(t, 1, state.ActivePlayer)         // It's the next player's turn
	assert.Equal(t, Guard, state.CardInHand[1])    // The other player had to discard and ended up drawing a Guard (based on the seed)
	assert.Equal(t, Guard, state.CardInHand[0])    // Our guard was moved into hand
	assert.Equal(t, Guard, state.ActivePlayerCard) // The other player now gets a turn and drew a Guard (based on the seed)
	assert.Equal(t, state.Deck[Guard], 2)          // So only two guards remain
}

func TestPlayingGuard(t *testing.T) {
	rand.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Countess
	state.ActivePlayerCard = Guard

	// Play the Prince on the other player
	err := state.PlayCard(Action{
		PlayRecent:   true,
		TargetPlayer: 1,
		SelectedCard: Handmaid, // incorrect guess
	})
	assert.NoError(t, err)

	assert.Equal(t, 1, state.ActivePlayer)         // It's the next player's turn
	assert.Equal(t, Countess, state.CardInHand[1]) // The other player did not discard the Countess
	assert.Equal(t, Prince, state.CardInHand[0])   // Our prince remained in hand
	assert.Equal(t, Guard, state.ActivePlayerCard) // The other player now gets a turn and drew a Guard (based on the seed)
	assert.Equal(t, state.Deck[Guard], 3)          // So only three guards remain
}
