package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	_, err := NewGame(2, r)
	assert.NoError(t, err)
}

const newGameToken = "2.2302.[443].[[]-[]].[00].[[00]-[00]].0.00.[78].1"

func TestGameToToken(t *testing.T) {
	r.Seed(0)
	game, _ := NewGame(2, r)
	assert.Equal(t, newGameToken, game.Token())
}

func TestGameFromToken(t *testing.T) {
	game := &Gamestate{}
	game.FromToken(newGameToken)
	// Now we test if it converts back again, since we don't have equality for a game
	assert.Equal(t, newGameToken, game.Token())
}

func TestPlayingPrincessIdiot(t *testing.T) {
	r.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Princess
	state.CardInHand[1] = Countess
	state.ActivePlayerCard = Guard

	// Play the Princess
	state.PlayCard(Action{
		PlayRecent:         false,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	}, r)

	assert.True(t, state.GameEnded)
	assert.Equal(t, 1, state.Winner)
}

func TestPlayingPrince(t *testing.T) {
	r.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Countess
	state.ActivePlayerCard = Guard

	// Play the Prince on the other player
	state.PlayCard(Action{
		PlayRecent:         false,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	}, r)

	assert.Equal(t, 1, state.ActivePlayer)         // It's the next player's turn
	assert.Equal(t, Guard, state.CardInHand[1])    // The other player had to discard and ended up drawing a Guard (based on the seed)
	assert.Equal(t, Guard, state.CardInHand[0])    // Our guard was moved into hand
	assert.Equal(t, Guard, state.ActivePlayerCard) // The other player now gets a turn and drew a Guard (based on the seed)
	assert.Equal(t, state.Deck[Guard], 2)          // So only two guards remain
}

func TestPlayingGuardBadGuess(t *testing.T) {
	r.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Countess
	state.ActivePlayerCard = Guard

	// Play the Guard on the other player
	state.PlayCard(Action{
		PlayRecent:         true,
		TargetPlayerOffset: 1,
		SelectedCard:       Handmaid, // incorrect guess
	}, r)

	assert.Equal(t, 1, state.ActivePlayer)         // It's the next player's turn
	assert.Equal(t, Countess, state.CardInHand[1]) // The other player did not discard the Countess
	assert.Equal(t, Prince, state.CardInHand[0])   // Our prince remained in hand
	assert.Equal(t, Guard, state.ActivePlayerCard) // The other player now gets a turn and drew a Guard (based on the seed)
	assert.Equal(t, state.Deck[Guard], 3)          // So only three guards remain
}

func TestPlayingGuardGoodGuess(t *testing.T) {
	r.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Countess
	state.ActivePlayerCard = Guard

	// Play the Guard on the other player
	state.PlayCard(Action{
		PlayRecent:         true,
		TargetPlayerOffset: 1,
		SelectedCard:       Countess, // correct guess
	}, r)

	assert.True(t, state.GameEnded)
	assert.Equal(t, 0, state.Winner)
}

func TestPlayingPrinceOnPrincess(t *testing.T) {
	r.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Princess
	state.ActivePlayerCard = Guard

	// Play the Prince on the other player
	state.PlayCard(Action{
		PlayRecent:         false,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	}, r)

	assert.True(t, state.GameEnded)
	assert.Equal(t, 0, state.Winner)
}

func TestPlayingPrinceWithCountess(t *testing.T) {
	r.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Princess
	state.ActivePlayerCard = Countess

	// Play the Prince on the other player, even though we MUST play the countess
	state.PlayCard(Action{
		PlayRecent:         false,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	}, r)

	// ...so we lose
	assert.True(t, state.GameEnded)
	assert.Equal(t, 1, state.Winner)
}

func TestPlayingCountessCorrectly(t *testing.T) {
	r.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  2,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Princess
	state.ActivePlayerCard = Countess

	// Play the Prince on the other player, even though we MUST play the countess
	state.PlayCard(Action{
		PlayRecent:         true,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	}, r)

	// .. and the game goes on
	assert.False(t, state.GameEnded)
	assert.Equal(t, 1, state.ActivePlayer)         // It's the next player's turn
	assert.Equal(t, Princess, state.CardInHand[1]) // The other player still has the Princess
	assert.Equal(t, Prince, state.CardInHand[0])   // Our prince remained in hand
	assert.Equal(t, Guard, state.ActivePlayerCard) // The other player now gets a turn and drew a Guard (based on the seed)
	assert.Equal(t, state.Deck[Guard], 3)          // So only three guards remain
}

func TestPlayingBaronWithCountessVsGuard(t *testing.T) {
	r.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 2,
		Baron:  1,
	}, 2)

	state.CardInHand[0] = Baron
	state.CardInHand[1] = Guard
	state.ActivePlayerCard = Countess

	state.PlayCard(Action{
		PlayRecent:         false,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	}, r)

	assert.True(t, state.GameEnded)
	assert.Equal(t, 0, state.Winner)
}

func TestPlayingBaronWithPriestVsKing(t *testing.T) {
	r.Seed(0)
	state := newGame(Deck{
		Guard:  4,
		Priest: 1,
		Baron:  1,
	}, 2)

	state.CardInHand[0] = Priest
	state.CardInHand[1] = King
	state.ActivePlayerCard = Baron

	state.PlayCard(Action{
		PlayRecent:         true,
		TargetPlayerOffset: 1,
		SelectedCard:       None,
	}, r)

	assert.True(t, state.GameEnded)
	assert.Equal(t, 1, state.Winner)
}

func TestAlmostEmptyDeck(t *testing.T) {
	r.Seed(0)
	state := newGame(Deck{
		Priest: 1,
		Prince: 1,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Princess
	state.ActivePlayerCard = Guard
	state.ActivePlayer = 1

	state.PlayCard(Action{
		PlayRecent:         true,
		TargetPlayerOffset: 1,
		SelectedCard:       Handmaid,
	}, r)

	assert.False(t, state.GameEnded)
	assert.Equal(t, 0, state.ActivePlayer)          // It's the next player's turn
	assert.Equal(t, Princess, state.CardInHand[1])  // The other player still has the Princess
	assert.Equal(t, Prince, state.CardInHand[0])    // Our prince remained in hand
	assert.Equal(t, Priest, state.ActivePlayerCard) // The other player now gets a turn and drew a Priest (based on the seed)
}

func TestLastPlay(t *testing.T) {
	r.Seed(0)
	state := newGame(Deck{
		Prince: 1,
	}, 2)

	state.CardInHand[0] = Prince
	state.CardInHand[1] = Princess
	state.ActivePlayerCard = Guard
	state.ActivePlayer = 1

	state.PlayCard(Action{
		PlayRecent:         true,
		TargetPlayerOffset: 1,
		SelectedCard:       Handmaid,
	}, r)

	assert.True(t, state.GameEnded)
	assert.Equal(t, 1, state.Winner)
}
