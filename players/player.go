package players

import "love-letter-ai/rules"

type Player interface {
	PlayCard(SimpleState) rules.Action
}
