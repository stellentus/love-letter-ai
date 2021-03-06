package state

import "love-letter-ai/rules"

// Simple provides a simplified state. It is only valid for 2 players.
type Simple struct {
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

// Simple converts a rules.Gamestate to a Simple
func NewSimple(gs rules.Gamestate) Simple {
	simple := Simple{}

	simple.Discards = gs.AllDiscards()
	simple.RecentDraw = gs.ActivePlayerCard
	simple.OldCard = gs.CardInHand[gs.ActivePlayer]

	// Figure out opponent's ID
	opponent := 0
	if gs.ActivePlayer == 0 {
		opponent = 1
	}
	if len(gs.Discards[opponent]) > 0 {
		// Get opponent's last played card. (If a Prince was played on the opponent, this will still show the last played card.)
		simple.OpponentCard = gs.LastPlay[opponent]
	} else {
		// else default to Princess, since that can't ever actually be the opponent's last card
		simple.OpponentCard = rules.Princess
	}
	simple.ScoreDiff = gs.Discards[gs.ActivePlayer].Score() - gs.Discards[opponent].Score()

	return simple
}

// AsIndex converts the simple state into an array index.
func (ss Simple) AsIndex() int {
	return Index(ss.Discards, ss.RecentDraw, ss.OldCard, ss.OpponentCard, ss.ScoreDiff)
}

// AsIndexWithAction converts the simple state and action into an array index.
func (ss Simple) AsIndexWithAction(act rules.Action) (int, int) {
	return Indices(ss.Discards, ss.RecentDraw, ss.OldCard, ss.OpponentCard, ss.ScoreDiff, act)
}
