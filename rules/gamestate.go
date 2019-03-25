package rules

import (
	"errors"
)

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

func (state *Gamestate) ActiveCardIsHighest() bool {
	return int(state.ActivePlayerCard) > int(state.CardInHand[state.ActivePlayer])
}

func (state *Gamestate) Win(activeWins bool) {
	if activeWins {
		state.Winner = state.ActivePlayer
	} else {
		if state.ActivePlayer == 0 {
			state.Winner = 1
		} else {
			state.Winner = 0
		}
	}
	state.GameEnded = true
}

func (state *Gamestate) TopCardForPlayer(player int) Card {
	return state.PlayerHistory[player][len(state.PlayerHistory[player])-1]
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
			// Automatically lose for cheating
			state.Win(false)
			return nil
		}
	}

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
			state.Win(true)
		}
		// Note we don't store this history, which a real player would rely upon. e.g. if I guess 4 and it's wrong, do I guess 4 again the next turn when no Handmaids have shown up? This bot would do that.
	case Priest:
		if !(targetPlayer >= 0 && targetPlayer <= 1 && targetPlayer != state.ActivePlayer) {
			return errors.New("You must target a valid player")
		}
		if state.TopCardForPlayer(targetPlayer) == Handmaid {
			break
		}
		// TODO: Store this knowledge, keeping track of whether the other player still has the observed card.
		// TODO: Also store whether the other player has seen my current card.
	case Baron:
		if !(targetPlayer >= 0 && targetPlayer <= 1 && targetPlayer != state.ActivePlayer) {
			return errors.New("You must target a valid player")
		}
		if state.TopCardForPlayer(targetPlayer) == Handmaid {
			break
		}
		// Compare cards. High wins. Tie does nothing
		targetValue := int(state.CardInHand[targetPlayer])
		activeValue := int(state.CardInHand[state.ActivePlayer])
		switch {
		case targetValue < activeValue:
			state.Win(true)
		case targetValue > activeValue:
			state.Win(false)
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
		state.PlayerHistory[targetPlayer] = append(state.PlayerHistory[targetPlayer], targetCard)
		state.CardInHand[targetPlayer] = state.Deck.Draw()

		if targetCard == Princess {
			state.Win(targetPlayer != state.ActivePlayer)
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
		state.CardInHand[targetPlayer] = state.CardInHand[state.ActivePlayer]
		state.CardInHand[state.ActivePlayer] = targetCard
	case Countess:
		// Do nothing
	case Princess:
		// Idiot!
		state.Win(false)
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
	if state.Deck.Size() != 1 {
		return errors.New("The deck size is incorrect for game end")
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
