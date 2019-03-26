package rules

// SimpleState provides a simplified state. It is only valid for 2 players.
type SimpleState struct {
	// discards is all of the cards discarded so far (unsorted, unattributed)
	discards Deck

	// highCard is the current player's high card
	highCard Card

	// lowCard is the current player's low card
	lowCard Card

	// opponentCard is the card most recently played by the opponent
	opponentCard Card

	// scoreDiff is the current player's score lead compared to the opponent
	scoreDiff int
}

// SimpleState converts a Gamestate to a SimpleState
func (state *Gamestate) AsSimpleState() SimpleState {
	simple := SimpleState{}

	simple.discards = state.Faceup.AsDeck()
	for _, val := range state.PlayerHistory {
		simple.discards.AddStack(val)
	}
	if state.ActiveCardIsHighest() {
		simple.highCard = state.ActivePlayerCard
		simple.lowCard = state.CardInHand[state.ActivePlayer]
	} else {
		simple.lowCard = state.ActivePlayerCard
		simple.highCard = state.CardInHand[state.ActivePlayer]
	}

	// Figure out opponent's ID
	opponent := 0
	if state.ActivePlayer == 0 {
		opponent = 1
	}
	if len(state.PlayerHistory[opponent]) > 0 {
		// Get opponent's last played card. (Note if the opponent played the Prince, then this will show the discarded card, instead.)
		simple.opponentCard = state.PlayerHistory[opponent][len(state.PlayerHistory[opponent])-1]
	} // else default to None
	simple.scoreDiff = state.PlayerHistory[state.ActivePlayer].Score() - state.PlayerHistory[opponent].Score()

	return simple
}
