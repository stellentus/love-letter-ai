package rules

type Action struct {
	// PlayHighest is true if the highest card is played (otherwise, lowest is played).
	PlayHighest bool

	// TargetPlayer is set to the ID of the player targeted by the card, if applicable.
	TargetPlayer int

	// SelectedCard is set to the Card chosen by the action, if applicable.
	SelectedCard Card
}
