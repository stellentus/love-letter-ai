package players

import (
	"math/rand"

	"love-letter-ai/rules"
)

// This AI will play random, but with the constraint that it won't break a rule or do something that will obviously eliminate it (i.e. discard the Princess).
// It might still make incredibly stupid plays (like look at an opponent's hand one round and then use the Baron to compare values, even though it should know that it will lose).

type RandomPlayer struct{}

func (rp *RandomPlayer) PlayCard(state rules.SimpleState, myID int) rules.Action {
	action := rules.Action{
		PlayRecent:   rand.Int31n(1) == 0,
		TargetPlayer: int(rand.Int31n(1) + 1),
		SelectedCard: rules.Card(rand.Int31n(int32(rules.Princess))),
	}

	playedCard := state.RecentDraw
	otherCard := state.OldCard
	if !action.PlayRecent {
		playedCard = state.OldCard
		otherCard = state.RecentDraw
	}

	// This assumes 2 players
	otherPlayer := 0
	if myID == 0 {
		otherPlayer = 1
	}

	if playedCard == rules.Princess {
		action.PlayRecent = !action.PlayRecent
		otherCard, playedCard = playedCard, otherCard
	}

	switch playedCard {
	case rules.Guard:
		action.TargetPlayer = otherPlayer
	case rules.Priest:
		action.TargetPlayer = otherPlayer
	case rules.Baron:
		action.TargetPlayer = otherPlayer
	case rules.Prince:
		if otherCard == rules.Princess && action.TargetPlayer == myID {
			action.TargetPlayer = otherPlayer
		} else if otherCard == rules.Countess {
			// Nope, must play Countess
			action.PlayRecent = !action.PlayRecent
		}
	case rules.King:
		if otherCard == rules.Countess {
			// Nope, must play Countess
			action.PlayRecent = !action.PlayRecent
		}
	default:
		// This is an error
	}

	return action
}
