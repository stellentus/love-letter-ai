package players

import (
	"math/rand"

	"love-letter-ai/rules"
	"love-letter-ai/state"
)

// This AI will play random, but with the constraint that it won't break a rule or do something that will obviously eliminate it (i.e. discard the Princess).
// It might still make incredibly stupid plays (like look at an opponent's hand one round and then use the Baron to compare values, even though it should know that it will lose).

type RandomPlayer struct{}

func (rp *RandomPlayer) PlayCard(state state.Simple) rules.Action {
	action := rules.Action{
		PlayRecent:         rand.Int31n(1) == 0,
		SelectedCard:       rules.Card(rand.Int31n(int32(rules.Princess))),
		TargetPlayerOffset: int(rand.Int31n(1) + 1), // This assumes two players
	}

	return playAction(state, action)
}

func (rp *RandomPlayer) PlayCardRand(state state.Simple, r *rand.Rand) rules.Action {
	action := rules.Action{
		PlayRecent:         r.Int31n(1) == 0,
		SelectedCard:       rules.Card(r.Int31n(int32(rules.Princess))),
		TargetPlayerOffset: int(r.Int31n(1) + 1), // This assumes two players
	}

	return playAction(state, action)
}

func playAction(state state.Simple, action rules.Action) rules.Action {
	playedCard := state.RecentDraw
	otherCard := state.OldCard
	if !action.PlayRecent {
		playedCard = state.OldCard
		otherCard = state.RecentDraw
	}

	// This assumes 2 players
	otherPlayerOffset := 1

	if playedCard == rules.Princess {
		action.PlayRecent = !action.PlayRecent
		otherCard, playedCard = playedCard, otherCard
	}

	switch playedCard {
	case rules.Guard:
		action.TargetPlayerOffset = otherPlayerOffset
	case rules.Priest:
		action.TargetPlayerOffset = otherPlayerOffset
	case rules.Baron:
		action.TargetPlayerOffset = otherPlayerOffset
	case rules.Prince:
		if otherCard == rules.Princess {
			action.TargetPlayerOffset = otherPlayerOffset
		} else if otherCard == rules.Countess {
			// Nope, must play Countess
			action.PlayRecent = !action.PlayRecent
		} else {
			// Okay, leave random offset
		}
	case rules.King:
		if otherCard == rules.Countess {
			// Nope, must play Countess
			action.PlayRecent = !action.PlayRecent
		} else {
			action.TargetPlayerOffset = otherPlayerOffset
		}
	default:
		// This is an error
	}

	return action
}
