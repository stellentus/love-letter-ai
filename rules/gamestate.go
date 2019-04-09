package rules

import (
	"errors"
)

type Gamestate struct {
	// NumPlayers is the number of players in the game.
	NumPlayers int

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

	// FinalState stores some state at the time the game ended. It's only set once GameEnded is true.
	FinalState

	// LossWasStupid is set to true if the play was something that would never, ever, be a good idea.
	// This is only meaningful in 2-player.
	LossWasStupid bool
}

type FinalState struct {
	// LastDiscard was the card discarded by the last play of the game.
	LastDiscard Card

	// LastInHand was the card that wasn't discarded in the last play.
	LastInHand Card

	// OpponentInHand was the card the opponent held at the start of the last play.
	OpponentInHand Card

	// RemainingDeck is the number of cards remaining in the deck at the end of the game.
	RemainingDeck int

	// DiscardWon is true if the active player won.
	DiscardWon bool
}

// NewGame deals out a new game for the specified number of players.
// This always assumes that player 0 is the starting player.
func NewGame(playerCount int) (Gamestate, error) {
	state := newGame(DefaultDeck(), playerCount)

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
	state.ActivePlayerCard = state.Deck.Draw()

	return state, nil
}

func (game Gamestate) Copy() Gamestate {
	gs := Gamestate{
		NumPlayers:       game.NumPlayers,
		Deck:             game.Deck.Copy(),
		Faceup:           game.Faceup.Copy(),
		LastPlay:         game.LastPlay.Copy(),
		ActivePlayer:     game.ActivePlayer,
		CardInHand:       game.CardInHand.Copy(),
		ActivePlayerCard: game.ActivePlayerCard,
		GameEnded:        game.GameEnded,
		Winner:           game.Winner,
		FinalState:       game.FinalState,
		LossWasStupid:    game.LossWasStupid,
	}

	gs.Discards = make([]Stack, 0, len(game.Discards))
	for _, val := range game.Discards {
		gs.Discards = append(gs.Discards, val.Copy())
	}

	gs.KnownCards = make([]Stack, 0, len(game.KnownCards))
	for _, val := range game.KnownCards {
		gs.KnownCards = append(gs.KnownCards, val.Copy())
	}

	gs.EliminatedPlayers = make([]bool, 0, len(game.EliminatedPlayers))
	for _, val := range game.EliminatedPlayers {
		gs.EliminatedPlayers = append(gs.EliminatedPlayers, val)
	}

	return gs
}

// NewSimpleGame deals out a new game for 2 players with a simplified deck.
// This always assumes that player 0 is the starting player.
// It draws 4 cards from the previous deck:
//	* 2 for player 0's hand
//	* 1 for player 1's hand
//	* 1 for player 1's last play
func NewSimpleGame(deck Deck) Gamestate {
	state := newGame(deck, 2)

	state.LastPlay[1] = state.Deck.Draw()

	for i := range state.CardInHand {
		state.CardInHand[i] = state.Deck.Draw()
	}
	state.ActivePlayerCard = state.Deck.Draw()

	return state
}

func newGame(deck Deck, playerCount int) Gamestate {
	state := Gamestate{
		NumPlayers:        playerCount,
		Deck:              deck,
		Faceup:            []Card{},
		Discards:          make([]Stack, playerCount),
		LastPlay:          make([]Card, playerCount),
		KnownCards:        make([]Stack, playerCount),
		ActivePlayer:      0,
		EliminatedPlayers: make([]bool, playerCount),
		CardInHand:        make([]Card, playerCount),
		ActivePlayerCard:  None,
		GameEnded:         false,
		Winner:            -1,
	}
	for i := range state.KnownCards {
		state.KnownCards[i] = make([]Card, playerCount)
	}
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

	pInGame := 0
	remainingPlayer := 0
	for pid, isElim := range state.EliminatedPlayers {
		if !isElim {
			pInGame += 1
			remainingPlayer = pid
		}
	}
	if pInGame == 1 {
		state.Winner = remainingPlayer
		state.GameEnded = true
	}

	state.updateFinalState()

	if state.CardInHand[player] != None {
		state.Discards[state.ActivePlayer] = append(state.Discards[state.ActivePlayer], state.CardInHand[player])
		state.CardInHand[player] = None
	}
	for i := range state.KnownCards[player] {
		state.KnownCards[player][i] = None
	}
}

func (state *Gamestate) updateFinalState() {
	state.FinalState = FinalState{
		LastDiscard:    state.ActivePlayerCard,
		LastInHand:     state.CardInHand[state.ActivePlayer],
		OpponentInHand: state.CardInHand[(state.ActivePlayer+1)%state.NumPlayers],
		RemainingDeck:  state.Deck.Size(),
		DiscardWon:     state.Winner == state.ActivePlayer,
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
func (state *Gamestate) PlayCard(action Action) {
	if state.GameEnded {
		return
	}

	// If the card to be played isn't the recent card, swap them to make the rest of this function easier.
	// Since that card will be discarded this turn, it doesn't matter that we do this.
	if !action.PlayRecent {
		card := state.ActivePlayerCard
		state.ActivePlayerCard = state.CardInHand[state.ActivePlayer]
		state.CardInHand[state.ActivePlayer] = card
	}

	state.clearKnownCard(state.ActivePlayer, state.ActivePlayerCard)
	state.Discards[state.ActivePlayer] = append(state.Discards[state.ActivePlayer], state.ActivePlayerCard)
	state.LastPlay[state.ActivePlayer] = state.ActivePlayerCard

	// If the retained card is the Countess, make sure that's allowed
	if state.CardInHand[state.ActivePlayer] == Countess {
		if state.ActivePlayerCard == King || state.ActivePlayerCard == Prince {
			// Automatically eliminated for cheating. This is not the same as the rules, which simply forbid this.
			state.eliminatePlayer(state.ActivePlayer)
			state.LossWasStupid = true
			return
		}
	}

	switch state.ActivePlayerCard {
	case Guard:
		if !(action.TargetPlayerOffset > 0 && action.TargetPlayerOffset < state.NumPlayers) {
			// You must target a valid player with a Guard
			state.eliminatePlayer(state.ActivePlayer)
			state.LossWasStupid = true
			break
		}
		targetPlayer := state.getTargetIDFromOffset(action.TargetPlayerOffset)
		if state.LastPlay[targetPlayer] == Handmaid {
			break
		}
		targetCard := state.CardInHand[targetPlayer]
		if targetCard == action.SelectedCard && targetCard != Guard {
			state.eliminatePlayer(targetPlayer)
		}
		// Note we don't store this history, which a real player would rely upon. e.g. if I guess 4 and it's wrong, do I guess 4 again the next turn when no Handmaids have shown up? This bot would do that.
	case Priest:
		if !(action.TargetPlayerOffset > 0 && action.TargetPlayerOffset < state.NumPlayers) {
			// You must target a valid player with a Priest
			state.eliminatePlayer(state.ActivePlayer)
			state.LossWasStupid = true
			break
		}
		targetPlayer := state.getTargetIDFromOffset(action.TargetPlayerOffset)
		if state.LastPlay[targetPlayer] == Handmaid {
			break
		}
		state.KnownCards[targetPlayer][state.ActivePlayer] = state.CardInHand[targetPlayer]
	case Baron:
		if !(action.TargetPlayerOffset > 0 && action.TargetPlayerOffset < state.NumPlayers) {
			// You must target a valid player with a Baron
			state.eliminatePlayer(state.ActivePlayer)
			state.LossWasStupid = true
			break
		}
		targetPlayer := state.getTargetIDFromOffset(action.TargetPlayerOffset)
		if state.LastPlay[targetPlayer] == Handmaid {
			break
		}
		// Compare cards. Eliminate low. Tie does nothing
		targetValue := int(state.CardInHand[targetPlayer])
		activeValue := int(state.CardInHand[state.ActivePlayer])
		switch {
		case targetValue < activeValue:
			state.eliminatePlayer(targetPlayer)
		case targetValue > activeValue:
			state.eliminatePlayer(state.ActivePlayer)
		}
	case Handmaid:
		// Do nothing
	case Prince:
		if !(action.TargetPlayerOffset >= 0 && action.TargetPlayerOffset < state.NumPlayers) {
			// You must target a valid player with a Prince
			state.eliminatePlayer(state.ActivePlayer)
			state.LossWasStupid = true
			break
		}
		targetPlayer := state.getTargetIDFromOffset(action.TargetPlayerOffset)
		if state.LastPlay[targetPlayer] == Handmaid {
			// If you target someone invalid, default to self.
			// The game rules say that if everyone else has a Handmaid, you must target yourself, so this is a good default.
			targetPlayer = state.ActivePlayer
		}
		targetCard := state.CardInHand[targetPlayer]
		if targetCard == Princess {
			// Do this first to update FinalState.
			state.eliminatePlayer(targetPlayer)
		} else {
			state.Discards[targetPlayer] = append(state.Discards[targetPlayer], targetCard)
			state.CardInHand[targetPlayer] = state.Deck.Draw()
		}
		state.clearKnownCard(targetPlayer, targetCard)
		if targetCard == Princess && targetPlayer == state.ActivePlayer {
			// This was stupid UNLESS the other player has a Handmaid, in which case this is okay.
			if state.NumPlayers == 2 && state.LastPlay[(state.ActivePlayer+1)%2] != Handmaid {
				state.LossWasStupid = true
			}
		}
	case King:
		if !(action.TargetPlayerOffset > 0 && action.TargetPlayerOffset < state.NumPlayers) {
			// You must target a valid player with a King
			state.eliminatePlayer(state.ActivePlayer)
			state.LossWasStupid = true
			break
		}
		targetPlayer := state.getTargetIDFromOffset(action.TargetPlayerOffset)
		if state.LastPlay[targetPlayer] == Handmaid {
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
		state.eliminatePlayer(state.ActivePlayer)
		state.LossWasStupid = true
	default:
		// An invalid card was played
		state.eliminatePlayer(state.ActivePlayer)
		state.LossWasStupid = true
	}

	if state.Deck.Size() > 1 {
		state.ActivePlayerCard = state.Deck.Draw()
		state.incrementPlayerTurn()
	} else {
		state.triggerGameEnd()
	}
}

func (state *Gamestate) getTargetIDFromOffset(offset int) int {
	return (state.ActivePlayer + offset) % state.NumPlayers
}

// incrementPlayerTurn increments the player turn. It assumes there are at least 2 active players
func (state *Gamestate) incrementPlayerTurn() {
	// Increment with rollover
	state.ActivePlayer = (state.ActivePlayer + 1) % state.NumPlayers

	// Skip past eliminated players
	for state.EliminatedPlayers[state.ActivePlayer] {
		state.ActivePlayer = (state.ActivePlayer + 1) % state.NumPlayers
	}
}

func (state *Gamestate) triggerGameEnd() {
	// It's an error if the deck size is > 1, but test code in this module could confirm that never happens.
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

	state.updateFinalState()
}
