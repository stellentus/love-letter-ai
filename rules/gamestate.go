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

	// PlayerHistory contains a Stack for each player, showing their face-up cards.
	PlayerHistory []Stack

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
		PlayerHistory:    make([]Stack, playerCount),
		ActivePlayer:     0,
		CardInHand:       make([]Card, playerCount),
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

func (state *Gamestate) EliminatePlayer(player int) {
	state.EliminatedPlayers[player] = true
	state.PlayerHistory[state.ActivePlayer] = append(state.PlayerHistory[state.ActivePlayer], state.CardInHand[player])
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

func (state *Gamestate) ActiveCardIsHighest() bool {
	return int(state.ActivePlayerCard) > int(state.CardInHand[state.ActivePlayer])
}

func (state *Gamestate) TopCardForPlayer(player int) Card {
	return state.PlayerHistory[player][len(state.PlayerHistory[player])-1]
}

func (state *Gamestate) ClearKnownCard(player int, card Card) {
	// Range through the list of my known cards to reset it if I discard the known card
	for i, val := range state.KnownCards[player] {
		if card == val {
			state.KnownCards[player][i] = None
		}
	}
}

func (state *Gamestate) PlayCard(playHighest bool, targetPlayer int, selectedCard Card) error {
	if state.GameEnded {
		return errors.New("The game has already ended")
	}

	// If the card to be played isn't the "active" card, swap them to make the rest of this function easier
	if state.ActiveCardIsHighest() != playHighest {
		card := state.ActivePlayerCard
		state.ActivePlayerCard = state.CardInHand[state.ActivePlayer]
		state.CardInHand[state.ActivePlayer] = card
	}

	// If the retained card is the Countess, make sure that's allowed
	if state.CardInHand[state.ActivePlayer] == Countess {
		if state.ActivePlayerCard == King || state.ActivePlayerCard == Prince {
			// Automatically eliminated for cheating
			state.EliminatePlayer(state.ActivePlayer)
			return nil
		}
	}

	state.ClearKnownCard(state.ActivePlayer, state.ActivePlayerCard)
	state.PlayerHistory[state.ActivePlayer] = append(state.PlayerHistory[state.ActivePlayer], state.ActivePlayerCard)

	switch state.ActivePlayerCard {
	case Guard:
		if !(targetPlayer >= 0 && targetPlayer <= 1 && targetPlayer != state.ActivePlayer) {
			return errors.New("You must target a valid player")
		}
		if state.TopCardForPlayer(targetPlayer) == Handmaid {
			break
		}
		targetCard := state.CardInHand[targetPlayer]
		if targetCard == selectedCard && targetCard != Guard {
			state.EliminatePlayer(targetPlayer)
		}
		// Note we don't store this history, which a real player would rely upon. e.g. if I guess 4 and it's wrong, do I guess 4 again the next turn when no Handmaids have shown up? This bot would do that.
	case Priest:
		if !(targetPlayer >= 0 && targetPlayer <= 1 && targetPlayer != state.ActivePlayer) {
			return errors.New("You must target a valid player")
		}
		if state.TopCardForPlayer(targetPlayer) == Handmaid {
			break
		}
		state.KnownCards[targetPlayer][state.ActivePlayer] = state.CardInHand[targetPlayer]
	case Baron:
		if !(targetPlayer >= 0 && targetPlayer <= 1 && targetPlayer != state.ActivePlayer) {
			return errors.New("You must target a valid player")
		}
		if state.TopCardForPlayer(targetPlayer) == Handmaid {
			break
		}
		// Compare cards. Eliminate low. Tie does nothing
		targetValue := int(state.CardInHand[targetPlayer])
		activeValue := int(state.CardInHand[state.ActivePlayer])
		switch {
		case targetValue < activeValue:
			state.EliminatePlayer(targetPlayer)
		case targetValue > activeValue:
			state.EliminatePlayer(state.ActivePlayer)
		}
	case Handmaid:
		// Do nothing
	case Prince:
		if !(targetPlayer >= 0 && targetPlayer <= 1) {
			return errors.New("You must target a valid player")
		}
		if state.TopCardForPlayer(targetPlayer) == Handmaid {
			break
		}
		targetCard := state.CardInHand[targetPlayer]
		state.ClearKnownCard(targetPlayer, targetCard)
		state.PlayerHistory[targetPlayer] = append(state.PlayerHistory[targetPlayer], targetCard)
		state.CardInHand[targetPlayer] = state.Deck.Draw()

		if targetCard == Princess {
			state.EliminatePlayer(targetPlayer)
		}
	case King:
		if !(targetPlayer >= 0 && targetPlayer <= 1) {
			return errors.New("You must target a valid player")
		}
		if state.TopCardForPlayer(targetPlayer) == Handmaid {
			break
		}
		// Trade hands
		targetCard := state.CardInHand[targetPlayer]
		activeCard := state.CardInHand[state.ActivePlayer]
		state.CardInHand[targetPlayer] = activeCard
		state.CardInHand[state.ActivePlayer] = targetCard
		// Update knowledge
		for i := range state.KnownCards[state.ActivePlayer] {
			if activeCard == state.KnownCards[state.ActivePlayer][i] {
				state.KnownCards[state.ActivePlayer][i] = targetCard
			}
			if targetCard == state.KnownCards[targetPlayer][i] {
				state.KnownCards[targetPlayer][i] = activeCard
			}
		}
		state.KnownCards[state.ActivePlayer][targetPlayer] = targetCard
		state.KnownCards[targetPlayer][state.ActivePlayer] = activeCard
	case Countess:
		// Do nothing
	case Princess:
		// Idiot!
		state.EliminatePlayer(state.ActivePlayer)
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
		scores := make([]int, len(state.PlayerHistory))
		maxScore := 0
		for i := range scores {
			for _, val := range state.PlayerHistory[i] {
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

// SimpleState provides a simplified state. It is only valid for 2 players.
func (state *Gamestate) SimpleState() (discards Deck, highCard, lowCard, opponentCard Card, scoreDiff int) {
	discards = state.Faceup.AsDeck()
	for _, val := range state.PlayerHistory {
		discards.AddStack(val)
	}
	if state.ActiveCardIsHighest() {
		highCard = state.ActivePlayerCard
		lowCard = state.CardInHand[state.ActivePlayer]
	} else {
		lowCard = state.ActivePlayerCard
		highCard = state.CardInHand[state.ActivePlayer]
	}

	if state.ActivePlayer == 0 {
		opponentCard = state.CardInHand[1]
		scoreDiff = state.PlayerHistory[0].Score() - state.PlayerHistory[1].Score()
	} else {
		opponentCard = state.CardInHand[0]
		scoreDiff = state.PlayerHistory[1].Score() - state.PlayerHistory[0].Score()
	}

	return
}
