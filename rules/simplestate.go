package rules

// SimpleState provides a simplified state. It is only valid for 2 players.
type SimpleState struct {
	// Discards is all of the cards discarded so far (unsorted, unattributed)
	Discards Deck

	// RecentDraw is the card the current player just drew
	RecentDraw Card

	// OldCard is the card the current player already had
	OldCard Card

	// OpponentCard is the card most recently played by the opponent
	OpponentCard Card

	// ScoreDiff is the current player's score lead compared to the opponent
	ScoreDiff int
}

// SimpleState converts a Gamestate to a SimpleState
func (state *Gamestate) AsSimpleState() SimpleState {
	simple := SimpleState{}

	simple.Discards = state.Faceup.AsDeck()
	for _, val := range state.Discards {
		simple.Discards.AddStack(val)
	}
	simple.RecentDraw = state.ActivePlayerCard
	simple.OldCard = state.CardInHand[state.ActivePlayer]

	// Figure out opponent's ID
	opponent := 0
	if state.ActivePlayer == 0 {
		opponent = 1
	}
	if len(state.Discards[opponent]) > 0 {
		// Get opponent's last played card. (Note if the opponent played the Prince, then this will show the discarded card, instead.)
		simple.OpponentCard = state.Discards[opponent][len(state.Discards[opponent])-1]
	} // else default to None
	simple.ScoreDiff = state.Discards[state.ActivePlayer].Score() - state.Discards[opponent].Score()

	return simple
}
