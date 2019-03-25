package rules

// Stack is an ordered stack of cards.
type Stack []Card

type Gamestate struct {
	// Deck includes all remaining cards INCLUDING any face-down cards which (in a real game) are pre-dealt face-down and are unknown.
	Deck

	// Faceup is a Stack of all cards that were dealt face-up to no one. Order is unimportant.
	Faceup Stack

	// PlayerHistory contains a Stack for each player, showing their face-up cards.
	PlayerHistory []Stack

	// CardInHand contains the single card in each player's hand. (Only the active player has a second card, which is separate below.)
	// This is NOT public information.
	CardInHand Stack

	// ActivePlayerCard is the active player's second card.
	// This is NOT public information.
	ActivePlayerCard Card
}
