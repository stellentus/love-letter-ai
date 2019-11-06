package players

import (
	"love-letter-ai/rules"
	"love-letter-ai/state"
)

type Player interface {
	PlayCard(state.Simple) rules.Action
}
