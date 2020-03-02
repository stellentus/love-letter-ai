package rules

type Action struct {
	// PlayRecent is true if the recently dealt card is played (otherwise, the old card is played).
	PlayRecent bool

	// TargetPlayerOffset is set (targetPlayerID - myID + numPlayers)%numPlayers.
	TargetPlayerOffset int

	// SelectedCard is set to the Card chosen by the action, if applicable.
	SelectedCard Card
}

// AsInt converts an action to an integer to be used for indexing.
// This integer only uses 4 bits and is only valid for 2 players.
// TargetPlayer is instead encoded as a bool: targetSelf. That never conflicts with SelectedCard.
func (act Action) AsInt() int {
	retVal := 0
	if act.PlayRecent {
		retVal = 1
	}
	if act.SelectedCard != None && act.SelectedCard != Guard {
		retVal += 2 * (int(act.SelectedCard) - 1)
	} else if act.TargetPlayerOffset > 0 {
		retVal += 2 * act.TargetPlayerOffset
	}
	return retVal
}

// ActionFromInt reverses action.AsInt, but only for the 4 bits that compose the action. Other bits are ignored.
// This only works for a 2-player game.
func ActionFromInt(st int) Action {
	act := Action{}
	if st%2 == 1 {
		act.PlayRecent = true
	}

	st = (st & 0xF) >> 1
	// Now st is the TargetPlayerOffset or SelectedCard. We don't know which, but if *any* card was selected,
	// then the other player must be targeted (since the played card is a guard), so for 2 players the offset is 1.
	if st > 0 {
		act.TargetPlayerOffset = 1
	}
	act.SelectedCard = Card(st + 1)
	return act
}

func (card Card) PossibleActions(isRecent bool) []Action {
	switch card {
	case Guard:
		acts := make([]Action, 0, 8)
		for card := Guard + 1; card <= Princess; card++ {
			acts = append(acts, Action{
				PlayRecent:         isRecent,
				TargetPlayerOffset: 1,
				SelectedCard:       card,
			})
		}
		return acts
	case Prince:
		return []Action{
			{
				PlayRecent:         isRecent,
				TargetPlayerOffset: 0,
			},
			{
				PlayRecent:         isRecent,
				TargetPlayerOffset: 1,
			},
		}
	case King:
		return []Action{{
			PlayRecent:         isRecent,
			TargetPlayerOffset: 1,
		}}
	}

	return []Action{{PlayRecent: isRecent}}
}
