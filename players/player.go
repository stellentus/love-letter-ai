package players

import (
	"math/rand"

	"love-letter-ai/rules"
	"love-letter-ai/state"
)

type Player interface {
	PlayCard(state.Simple) rules.Action
	PlayCardRand(state.Simple, *rand.Rand) rules.Action
}
