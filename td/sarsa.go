package td

import (
	"fmt"
	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/state"
	"os"
)

type UpdateQFunc func(gameEnded bool, sa int, reward float32)

type sarsaLearner struct {
	// td is the backing data
	td *TD

	// lastQ was the last state-action value.
	lastQ int

	// updateQFunc is the function that updates TD's qf
	UpdateQFunc
}

const (
	unsetState       = state.ActionSpaceMagnitude
	winReward        = 100
	stupidReward     = -100
	forfeitWinReward = 1 // Only a minor benefit for winning because the other player was an idiot
	noReward         = 0
	lossReward       = -0.1 // The penalty for losing is minor since it might not have been the player's fault
)

func TrainSarsa(td *TD, episodes int) {
	pls := []*sarsaLearner{
		&sarsaLearner{td: td, UpdateQFunc: nil},
		&sarsaLearner{td: td, UpdateQFunc: nil},
	}

	templateSG, err := rules.NewGame(2)
	if err != nil {
		panic(err.Error())
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

		pls[0].lastQ = unsetState
		pls[1].lastQ = unsetState

		for !sg.GameEnded {
			action, err := pls[sg.ActivePlayer].learningAction(sg)
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
			pls[(sg.Winner+1)%2].updateQ(sg.GameEnded, sa, stupidReward)
			pls[sg.Winner].updateQ(sg.GameEnded, sa, forfeitWinReward)
		} else {
			pls[(sg.Winner+1)%2].updateQ(sg.GameEnded, sa, lossReward)
			pls[sg.Winner].updateQ(sg.GameEnded, sa, winReward)
		}
	}
	fmt.Fprintln(os.Stderr, "\r100.0% complete")
}

// learningAction provides a suggested action for the provided state.
// However, it also assumes it's being called for each play in a game so it can update the policy.
func (sl *sarsaLearner) learningAction(game rules.Gamestate) (rules.Action, error) {
	action, sa := players.EpsilonGreedyAction(sl.td, state.NewSimple(game), sl.td.Epsilon)
	sl.updateQ(game.GameEnded, sa, noReward)
	return action, nil
}

func (sl *sarsaLearner) updateQ(gameEnded bool, sa int, reward float32) {
	// Now save the update
	if sl.lastQ != unsetState {
		thisValue := float32(0) // If game ended, the value of the new state is 0 because it's a terminal state
		if !gameEnded {
			thisValue = sl.td.Gamma * sl.td.qf[sa]
		}
		sl.td.qf[sl.lastQ] += sl.td.Alpha * (reward + thisValue - sl.td.qf[sl.lastQ])
	}
	sl.lastQ = sa
}
