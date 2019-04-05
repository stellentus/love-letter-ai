package rules

type Action struct {
	// PlayRecent is true if the recently dealt card is played (otherwise, the old card is played).
	PlayRecent bool

	// TargetPlayer is set to the ID of the player targeted by the card, if applicable.
	TargetPlayer int

	// SelectedCard is set to the Card chosen by the action, if applicable.
	SelectedCard Card
}

// AsInt converts an action to an integer to be used for indexing.
// This integer only uses 4 bits.
func (act Action) AsInt() int {
	retVal := 0
	if act.PlayRecent {
		retVal = 1
	}
	if act.TargetPlayer != 0 {
		retVal += 2 * act.TargetPlayer
	} else if act.SelectedCard != None {
		retVal += 2 * (int(act.SelectedCard) - 1)
	}
	return retVal
}
