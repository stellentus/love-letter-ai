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

	if state.ActivePlayer == 0 {
		simple.opponentCard = state.CardInHand[1]
		simple.scoreDiff = state.PlayerHistory[0].Score() - state.PlayerHistory[1].Score()
	} else {
		simple.opponentCard = state.CardInHand[0]
		simple.scoreDiff = state.PlayerHistory[1].Score() - state.PlayerHistory[0].Score()
	}

	return simple
}
