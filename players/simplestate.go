package players

import (
	"love-letter-ai/rules"
	"love-letter-ai/state"
)

// SimpleState provides a simplified state. It is only valid for 2 players.
type SimpleState struct {
	// Discards is all of the cards discarded so far (unsorted, unattributed)
	Discards rules.Deck

	// RecentDraw is the card the current player just drew
	RecentDraw rules.Card

	// OldCard is the card the current player already had
	OldCard rules.Card

	// OpponentCard is the card most recently played by the opponent
	OpponentCard rules.Card

	// ScoreDiff is the current player's score lead compared to the opponent
	ScoreDiff int
}

// SimpleState converts a rules.Gamestate to a SimpleState
func NewSimpleState(state rules.Gamestate) SimpleState {
	simple := SimpleState{}

	simple.Discards = state.AllDiscards()
	simple.RecentDraw = state.ActivePlayerCard
	simple.OldCard = state.CardInHand[state.ActivePlayer]

	// Figure out opponent's ID
	opponent := 0
	if state.ActivePlayer == 0 {
		opponent = 1
	}
	if len(state.Discards[opponent]) > 0 {
		// Get opponent's last played card. (If a Prince was played on the opponent, this will still show the last played card.)
		simple.OpponentCard = state.LastPlay[opponent]
	} // else default to None
	simple.ScoreDiff = state.Discards[state.ActivePlayer].Score() - state.Discards[opponent].Score()

	return simple
}

func (ss SimpleState) AsInt() int {
	return state.Index(ss.Discards, ss.RecentDraw, ss.OldCard, ss.OpponentCard, ss.ScoreDiff)
}

func (ss SimpleState) AsIntWithAction(act rules.Action) (int, int) {
	return state.Indices(ss.Discards, ss.RecentDraw, ss.OldCard, ss.OpponentCard, ss.ScoreDiff, act)
}
