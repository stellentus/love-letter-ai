package players

import (
	"fmt"
	"love-letter-ai/rules"
	"love-letter-ai/state"
	"math/rand"
	"os"
)

type TrainingPlayer interface {
	Player

	// GreedyAction returns the greedy action for the given state, along with the corresponding state-action.
	// The action may be nil if nothing has trained yet.
	GreedyAction(state int) (*rules.Action, int)

	// UpdateQ is called to update the player's state-action values.
	UpdateQ(gameEnded bool, lastQ, sa int, reward float32)
}

type trainer struct {
	// tp is the player model being trained
	tp TrainingPlayer

	// lastQ was the last state-action value.
	lastQ int
}

const (
	unsetState       = state.ActionSpaceMagnitude
	winReward        = 100
	stupidReward     = -100
	forfeitWinReward = 1 // Only a minor benefit for winning because the other player was an idiot
	noReward         = 0
	lossReward       = -0.1 // The penalty for losing is minor since it might not have been the player's fault
)

func Train(pls []TrainingPlayer, episodes int, epsilon float32) {
	templateSG, err := rules.NewGame(2)
	if err != nil {
		panic(err.Error())
	}

	trs := make([]trainer, len(pls))
	for i, pl := range pls {
		trs[i] = trainer{tp: pl}
	}

	epPrintMod := episodes / 100000
	if epPrintMod < 1 {
		epPrintMod = 1
	}

	for i := 0; i < episodes; i++ {
		if (i % epPrintMod) == 0 {
			fmt.Fprintf(os.Stderr, "\r%2.2f%% complete", float32(i)/float32(episodes)*100)
		}
		if (i % 100) == 0 {
			// Every so often, start from a new starting state
			templateSG, err = rules.NewGame(2)
			if err != nil {
				panic(err.Error())
			}
		}
		sg := templateSG.Copy()

		trs[0].lastQ = unsetState
		trs[1].lastQ = unsetState

		for !sg.GameEnded {
			action, err := trs[sg.ActivePlayer].learningAction(sg, epsilon)
			if err != nil {
				panic(err.Error())
			}

			sg.PlayCard(action)
		}

		// Now allow both players to update based on the end of the game.
		sa, _ := state.NewSimple(sg).AsIndexWithAction(rules.Action{})
		if sa < 0 {
			panic(fmt.Sprintf("Negative state was calculated: %d", sa))
		}
		if sg.LossWasStupid {
			// This only happens if the play is something that will ALWAYS lose the game, so incur a huge penalty
			trs[(sg.Winner+1)%2].updateQ(sg.GameEnded, sa, stupidReward)
			trs[sg.Winner].updateQ(sg.GameEnded, sa, forfeitWinReward)
		} else {
			trs[(sg.Winner+1)%2].updateQ(sg.GameEnded, sa, lossReward)
			trs[sg.Winner].updateQ(sg.GameEnded, sa, winReward)
		}
	}
	fmt.Fprintln(os.Stderr, "\r100.0% complete")
}

// epsilonGreedyAction provides a suggested action for the provided state.
// If it hasn't learned anything for this state, it plays randomly.
// It will also choose a random action with probability Epsilon. This isn't exactly
// Epsilon-greedy because it doesn't subtract the probability of the greedy action.
func epsilonGreedyAction(pl TrainingPlayer, st state.Simple, epsilon float32) (rules.Action, int) {
	sNoAct := st.AsIndex()
	act, sa := pl.GreedyAction(sNoAct)
	if act == nil || rand.Float32() < epsilon {
		action := (&RandomPlayer{}).PlayCard(st)
		return action, state.IndexWithAction(sNoAct, action)
	}
	return *act, sa
}

// learningAction provides a suggested action for the provided state.
// However, it also assumes it's being called for each play in a game so it can update the policy.
func (tr *trainer) learningAction(game rules.Gamestate, epsilon float32) (rules.Action, error) {
	action, sa := epsilonGreedyAction(tr.tp, state.NewSimple(game), epsilon)
	tr.updateQ(game.GameEnded, sa, noReward)
	return action, nil
}

func (tr *trainer) updateQ(gameEnded bool, sa int, reward float32) {
	// Now save the update
	if tr.lastQ != unsetState {
		tr.tp.UpdateQ(gameEnded, tr.lastQ, sa, reward)
	}
	tr.lastQ = sa
}
