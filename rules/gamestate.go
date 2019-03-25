package rules

import "errors"

// Stack is an ordered stack of cards.
type Stack []Card

type Gamestate struct {
	// Deck includes all remaining cards.
	// Note this includes the one card that is always dealt face-down in a real game, so here a game should end with one card in this deck.
	Deck

	// Faceup is a Stack of all cards that were dealt face-up to no one. Order is unimportant.
	Faceup Stack

	// PlayerHistory contains a Stack for each player, showing their face-up cards.
	PlayerHistory []Stack

	// ActivePlayer is the id of the active player.
	ActivePlayer int

	// CardInHand contains the single card in each player's hand. (Only the active player has a second card, which is separate below.)
	// This is NOT public information.
	CardInHand Stack

	// ActivePlayerCard is the active player's second card.
	// This is NOT public information.
	ActivePlayerCard Card
}

// NewGame deals out a new game for the specified number of players.
// This always assumes that player 0 is the starting player.
func NewGame(playerCount int) (Gamestate, error) {
	state := Gamestate{
		Deck:             DefaultDeck(),
		PlayerHistory:    make([]Stack, playerCount),
		ActivePlayer:     0,
		CardInHand:       make([]Card, playerCount),
		ActivePlayerCard: None,
	}

	switch playerCount {
	case 2:
		// Draw 3 cards face up
		for i := 0; i < 3; i++ {
			state.Faceup = append(state.Faceup, state.Deck.Draw())
		}
	default:
		return Gamestate{}, errors.New("Only 2-player games are supported")
	}

	for i := range state.CardInHand {
		state.CardInHand[i] = state.Deck.Draw()
	}
	state.ActivePlayerCard = state.Deck.Draw()

	return state, nil
}
