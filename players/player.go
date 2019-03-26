package players

import "love-letter-ai/rules"

type Player interface {
	PlayCard(rules.SimpleState, int) rules.Action
}
