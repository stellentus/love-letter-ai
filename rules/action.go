package rules

type Action struct {
	// PlayRecent is true if the recently dealt card is played (otherwise, the old card is played).
	PlayRecent bool

	// TargetPlayer is set to the ID of the player targeted by the card, if applicable.
	TargetPlayer int

	// SelectedCard is set to the Card chosen by the action, if applicable.
	SelectedCard Card
}
