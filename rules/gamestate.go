package rules

import (
	"errors"
)

type Gamestate struct {
	// Deck includes all remaining cards.
	// Note this includes the one card that is always dealt face-down in a real game, so here a game should end with one card in this deck.
	Deck

	// Faceup is a Stack of all cards that were dealt face-up to no one. Order is unimportant.
	Faceup Stack

	// Discards contains a Stack for each player, showing their face-up cards.
	// Note the top card here might not be the most recently played card if a Prince was played against this player.
	Discards []Stack

	// LastPlay contains the last card played by each player. This will often be the last card in the player's discard stack, but not always.
	LastPlay Stack

	// KnownCards contains a Stack for each player, with a Stack of their knowledge of opponents' cards.
	// Index first by player about whom you want to know, then by the index of the player who might know something.
	// A card of 'None' means no knowledge.
	KnownCards []Stack

	// ActivePlayer is the id of the active player.
	ActivePlayer int

	// EliminatedPlayers is true if a given player has been eliminated.
	EliminatedPlayers []bool

	// CardInHand contains the single card in each player's hand. (Only the active player has a second card, which is separate below.)
	// This is NOT public information.
	CardInHand Stack

	// ActivePlayerCard is the active player's second card.
	// This is NOT public information.
	ActivePlayerCard Card

	// GameEnded is true if the game is over.
	GameEnded bool

	// Winner is the id of the winning player. It is only valid once a player has won.
	Winner int
}

// NewGame deals out a new game for the specified number of players.
// This always assumes that player 0 is the starting player.
func NewGame(playerCount int) (Gamestate, error) {
	state := Gamestate{
		Deck:             DefaultDeck(),
		Discards:         make([]Stack, playerCount),
		ActivePlayer:     0,
		CardInHand:       make([]Card, playerCount),
		LastPlay:         make([]Card, playerCount),
		KnownCards:       make([]Stack, playerCount),
		ActivePlayerCard: None,
	}

	if playerCount == 2 {
		// Draw 3 cards face up
		for i := 0; i < 3; i++ {
			state.Faceup = append(state.Faceup, state.Deck.Draw())
		}
	} else if playerCount != 3 && playerCount != 4 {
		return Gamestate{}, errors.New("Only games with 2, 3, or 4 players are supported")
	}

	for i := range state.CardInHand {
		state.CardInHand[i] = state.Deck.Draw()
	}
	for i := range state.KnownCards {
		state.KnownCards[i] = make([]Card, playerCount)
	}
	state.ActivePlayerCard = state.Deck.Draw()

	return state, nil
}

// NewSimpleGame deals out a new game for 2 players with a simplified deck.
// This always assumes that player 0 is the starting player.
// It draws 4 cards from the previous deck:
//	* 2 for player 0's hand
//	* 1 for player 1's hand
//	* 1 for player 1's last play
func NewSimpleGame(deck Deck) Gamestate {
	playerCount := 2

	state := Gamestate{
		Deck:             deck,
		Discards:         make([]Stack, playerCount),
		ActivePlayer:     0,
		CardInHand:       make([]Card, playerCount),
		LastPlay:         make([]Card, playerCount),
		KnownCards:       make([]Stack, playerCount),
		ActivePlayerCard: None,
		Faceup:           []Card{},
	}

	state.LastPlay[1] = state.Deck.Draw()

	for i := range state.CardInHand {
		state.CardInHand[i] = state.Deck.Draw()
	}
	for i := range state.KnownCards {
		state.KnownCards[i] = make([]Card, playerCount)
	}
	state.ActivePlayerCard = state.Deck.Draw()

	return state
}

func (state *Gamestate) AllDiscards() Deck {
	discards := state.Faceup.AsDeck()
	for _, val := range state.Discards {
		discards.AddStack(val)
	}
	return discards
}

func (state *Gamestate) eliminatePlayer(player int) {
	state.EliminatedPlayers[player] = true
	state.Discards[state.ActivePlayer] = append(state.Discards[state.ActivePlayer], state.CardInHand[player])
	state.CardInHand[player] = None
	for i := range state.KnownCards[player] {
		state.KnownCards[player][i] = None
	}

	pInGame := 0
	remainingPlayer := 0
	for pid, isIn := range state.EliminatedPlayers {
		if isIn {
			pInGame += 1
			remainingPlayer = pid
		}
	}
	if pInGame == 1 {
		state.Winner = remainingPlayer
		state.GameEnded = true
	}
}

func (state *Gamestate) clearKnownCard(player int, card Card) {
	// Range through the list of my known cards to reset it if I discard the known card
	for i, val := range state.KnownCards[player] {
		if card == val {
			state.KnownCards[player][i] = None
		}
	}
}

// PlayCard takes the provided action. Of course only the active player should call this at any time.
func (state *Gamestate) PlayCard(action Action) error {
	if state.GameEnded {
		return errors.New("The game has already ended")
	}

	// If the card to be played isn't the recent card, swap them to make the rest of this function easier.
	// Since that card will be discarded this turn, it doesn't matter that we do this.
	if !action.PlayRecent {
		card := state.ActivePlayerCard
		state.ActivePlayerCard = state.CardInHand[state.ActivePlayer]
		state.CardInHand[state.ActivePlayer] = card
	}

	// If the retained card is the Countess, make sure that's allowed
	if state.CardInHand[state.ActivePlayer] == Countess {
		if state.ActivePlayerCard == King || state.ActivePlayerCard == Prince {
			// Automatically eliminated for cheating. This is not the same as the rules, which simply forbid this.
			state.eliminatePlayer(state.ActivePlayer)
			return nil
		}
	}

	state.clearKnownCard(state.ActivePlayer, state.ActivePlayerCard)
	state.Discards[state.ActivePlayer] = append(state.Discards[state.ActivePlayer], state.ActivePlayerCard)
	state.LastPlay[state.ActivePlayer] = state.ActivePlayerCard

	switch state.ActivePlayerCard {
	case Guard:
		if !(action.TargetPlayer >= 0 && action.TargetPlayer <= 1 && action.TargetPlayer != state.ActivePlayer) {
			return errors.New("You must target a valid player")
		}
		if state.LastPlay[action.TargetPlayer] == Handmaid {
			break
		}
		targetCard := state.CardInHand[action.TargetPlayer]
		if targetCard == action.SelectedCard && targetCard != Guard {
			state.eliminatePlayer(action.TargetPlayer)
		}
		// Note we don't store this history, which a real player would rely upon. e.g. if I guess 4 and it's wrong, do I guess 4 again the next turn when no Handmaids have shown up? This bot would do that.
	case Priest:
		if !(action.TargetPlayer >= 0 && action.TargetPlayer <= 1 && action.TargetPlayer != state.ActivePlayer) {
			return errors.New("You must target a valid player")
		}
		if state.LastPlay[action.TargetPlayer] == Handmaid {
			break
		}
		state.KnownCards[action.TargetPlayer][state.ActivePlayer] = state.CardInHand[action.TargetPlayer]
	case Baron:
		if !(action.TargetPlayer >= 0 && action.TargetPlayer <= 1 && action.TargetPlayer != state.ActivePlayer) {
			return errors.New("You must target a valid player")
		}
		if state.LastPlay[action.TargetPlayer] == Handmaid {
			break
		}
		// Compare cards. Eliminate low. Tie does nothing
		targetValue := int(state.CardInHand[action.TargetPlayer])
		activeValue := int(state.CardInHand[state.ActivePlayer])
		switch {
		case targetValue < activeValue:
			state.eliminatePlayer(action.TargetPlayer)
		case targetValue > activeValue:
			state.eliminatePlayer(state.ActivePlayer)
		}
	case Handmaid:
		// Do nothing
	case Prince:
		if !(action.TargetPlayer >= 0 && action.TargetPlayer <= 1) {
			return errors.New("You must target a valid player")
		}
		if state.LastPlay[action.TargetPlayer] == Handmaid {
			// If you target someone invalid, default to self.
			// The game rules say that if everyone else has a Handmaid, you must target yourself, so this is a good default.
			action.TargetPlayer = state.ActivePlayer
		}
		targetCard := state.CardInHand[action.TargetPlayer]
		state.clearKnownCard(action.TargetPlayer, targetCard)
		state.Discards[action.TargetPlayer] = append(state.Discards[action.TargetPlayer], targetCard)
		state.CardInHand[action.TargetPlayer] = state.Deck.Draw()

		if targetCard == Princess {
			state.eliminatePlayer(action.TargetPlayer)
		}
	case King:
		if !(action.TargetPlayer >= 0 && action.TargetPlayer <= 1) {
			return errors.New("You must target a valid player")
		}
		if state.LastPlay[action.TargetPlayer] == Handmaid {
			break
		}
		// Trade hands
		targetCard := state.CardInHand[action.TargetPlayer]
		activeCard := state.CardInHand[state.ActivePlayer]
		state.CardInHand[action.TargetPlayer] = activeCard
		state.CardInHand[state.ActivePlayer] = targetCard
		// Update knowledge
		for i := range state.KnownCards[state.ActivePlayer] {
			if activeCard == state.KnownCards[state.ActivePlayer][i] {
				state.KnownCards[state.ActivePlayer][i] = targetCard
			}
			if targetCard == state.KnownCards[action.TargetPlayer][i] {
				state.KnownCards[action.TargetPlayer][i] = activeCard
			}
		}
		state.KnownCards[state.ActivePlayer][action.TargetPlayer] = targetCard
		state.KnownCards[action.TargetPlayer][state.ActivePlayer] = activeCard
	case Countess:
		// Do nothing
	case Princess:
		// Idiot!
		state.eliminatePlayer(state.ActivePlayer)
	default:
		return errors.New("An invalid card was played")
	}

	if state.Deck.Size() > 1 {
		state.ActivePlayerCard = state.Deck.Draw()
		state.ActivePlayer++
	} else {
		state.triggerGameEnd()
	}

	return nil
}

func (state *Gamestate) triggerGameEnd() error {
	if state.Deck.Size() >= 1 {
		// The deck could be size 0 if a prince was played in the last round.
		return errors.New("The deck is too big")
	}

	tie := false
	maxCard := 0
	state.Winner = 0
	for pid, val := range state.CardInHand {
		if int(val) > maxCard {
			maxCard = int(val)
			state.Winner = pid
			tie = false
		} else if int(val) == maxCard {
			tie = true
		}
	}

	if tie {
		scores := make([]int, len(state.Discards))
		maxScore := 0
		for i := range scores {
			for _, val := range state.Discards[i] {
				scores[i] += int(val)
			}
			if scores[i] > maxScore {
				maxScore = scores[i]
				state.Winner = i
				tie = false
			} else if scores[i] == maxScore {
				tie = true // We don't deal with this
			}
		}
	}

	state.GameEnded = true

	return nil
}
