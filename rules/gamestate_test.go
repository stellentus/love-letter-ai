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

func TestPlayingPrincessIdiot(t *testing.T) {
	rand.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Princess
	state.CardInHand[1] = Countess
	state.ActivePlayerCard = Guard

	// Play the Princess
	err := state.PlayCard(Action{
		PlayRecent:         false,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	})
	assert.NoError(t, err)

	assert.True(t, state.GameEnded)
	assert.Equal(t, 1, state.Winner)
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
		PlayRecent:         false,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	})
	assert.NoError(t, err)

	assert.Equal(t, 1, state.ActivePlayer)         // It's the next player's turn
	assert.Equal(t, Guard, state.CardInHand[1])    // The other player had to discard and ended up drawing a Guard (based on the seed)
	assert.Equal(t, Guard, state.CardInHand[0])    // Our guard was moved into hand
	assert.Equal(t, Guard, state.ActivePlayerCard) // The other player now gets a turn and drew a Guard (based on the seed)
	assert.Equal(t, state.Deck[Guard], 2)          // So only two guards remain
}

func TestPlayingGuardBadGuess(t *testing.T) {
	rand.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Countess
	state.ActivePlayerCard = Guard

	// Play the Guard on the other player
	err := state.PlayCard(Action{
		PlayRecent:         true,
		TargetPlayerOffset: 1,
		SelectedCard:       Handmaid, // incorrect guess
	})
	assert.NoError(t, err)

	assert.Equal(t, 1, state.ActivePlayer)         // It's the next player's turn
	assert.Equal(t, Countess, state.CardInHand[1]) // The other player did not discard the Countess
	assert.Equal(t, Prince, state.CardInHand[0])   // Our prince remained in hand
	assert.Equal(t, Guard, state.ActivePlayerCard) // The other player now gets a turn and drew a Guard (based on the seed)
	assert.Equal(t, state.Deck[Guard], 3)          // So only three guards remain
}

func TestPlayingGuardGoodGuess(t *testing.T) {
	rand.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Countess
	state.ActivePlayerCard = Guard

	// Play the Guard on the other player
	err := state.PlayCard(Action{
		PlayRecent:         true,
		TargetPlayerOffset: 1,
		SelectedCard:       Countess, // correct guess
	})
	assert.NoError(t, err)

	assert.True(t, state.GameEnded)
	assert.Equal(t, 0, state.Winner)
}

func TestPlayingPrinceOnPrincess(t *testing.T) {
	rand.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Princess
	state.ActivePlayerCard = Guard

	// Play the Prince on the other player
	err := state.PlayCard(Action{
		PlayRecent:         false,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	})
	assert.NoError(t, err)

	assert.True(t, state.GameEnded)
	assert.Equal(t, 0, state.Winner)
}

func TestPlayingPrinceWithCountess(t *testing.T) {
	rand.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Princess
	state.ActivePlayerCard = Countess

	// Play the Prince on the other player, even though we MUST play the countess
	err := state.PlayCard(Action{
		PlayRecent:         false,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	})
	assert.NoError(t, err)

	// ...so we lose
	assert.True(t, state.GameEnded)
	assert.Equal(t, 1, state.Winner)
}

func TestPlayingCountessCorrectly(t *testing.T) {
	rand.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Princess
	state.ActivePlayerCard = Countess

	// Play the Prince on the other player, even though we MUST play the countess
	err := state.PlayCard(Action{
		PlayRecent:         true,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	})
	assert.NoError(t, err)

	// .. and the game goes on
	assert.False(t, state.GameEnded)
	assert.Equal(t, 1, state.ActivePlayer)         // It's the next player's turn
	assert.Equal(t, Princess, state.CardInHand[1]) // The other player still has the Princess
	assert.Equal(t, Prince, state.CardInHand[0])   // Our prince remained in hand
	assert.Equal(t, Guard, state.ActivePlayerCard) // The other player now gets a turn and drew a Guard (based on the seed)
	assert.Equal(t, state.Deck[Guard], 3)          // So only three guards remain
}

func TestPlayingBaronWithCountessVsGuard(t *testing.T) {
	rand.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  1,
	}, 2)

	state.CardInHand[0] = Baron
	state.CardInHand[1] = Guard
	state.ActivePlayerCard = Countess

	err := state.PlayCard(Action{
		PlayRecent:         false,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	})
	assert.NoError(t, err)

	assert.True(t, state.GameEnded)
	assert.Equal(t, 0, state.Winner)
}

func TestPlayingBaronWithPriestVsKing(t *testing.T) {
	rand.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 1,
		Baron:  1,
	}, 2)

	state.CardInHand[0] = Priest
	state.CardInHand[1] = King
	state.ActivePlayerCard = Baron

	err := state.PlayCard(Action{
		PlayRecent:         true,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	})
	assert.NoError(t, err)

	assert.True(t, state.GameEnded)
	assert.Equal(t, 1, state.Winner)
}

func TestAlmostEmptyDeck(t *testing.T) {
	rand.Seed(0)
	state := newGame(Deck{
		Priest: 1,
		Prince: 1,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Princess
	state.ActivePlayerCard = Guard
	state.ActivePlayer = 1

	err := state.PlayCard(Action{
		PlayRecent:         true,
		TargetPlayerOffset: 1,
		SelectedCard:       Handmaid,
	})
	assert.NoError(t, err)

	assert.False(t, state.GameEnded)
	assert.Equal(t, 0, state.ActivePlayer)          // It's the next player's turn
	assert.Equal(t, Princess, state.CardInHand[1])  // The other player still has the Princess
	assert.Equal(t, Prince, state.CardInHand[0])    // Our prince remained in hand
	assert.Equal(t, Priest, state.ActivePlayerCard) // The other player now gets a turn and drew a Priest (based on the seed)
}

func TestLastPlay(t *testing.T) {
	rand.Seed(0)
	state := newGame(Deck{
		Prince: 1,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Princess
	state.ActivePlayerCard = Guard
	state.ActivePlayer = 1

	err := state.PlayCard(Action{
		PlayRecent:         true,
		TargetPlayerOffset: 1,
		SelectedCard:       Handmaid,
	})
	assert.NoError(t, err)

	assert.True(t, state.GameEnded)
	assert.Equal(t, 1, state.Winner)
}
