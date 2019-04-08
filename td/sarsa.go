package td

import (
	"fmt"
	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/state"
	"math/rand"
)

type Sarsa struct {
	qf      []float32
	Epsilon float32
	Alpha   float32
	Gamma   float32
}

type sarsaLearner struct {
	// sarsa is the backing data
	sarsa *Sarsa

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

func NewSarsa(epsilon, alpha, gamma float32) *Sarsa {
	sar := &Sarsa{
		qf:      make([]float32, state.ActionSpaceMagnitude, state.ActionSpaceMagnitude),
		Epsilon: epsilon,
		Alpha:   alpha,
		Gamma:   gamma,
	}
	for i := range sar.qf {
		sar.qf[i] = 0.5
	}
	return sar
}

func (sarsa *Sarsa) NewPlayer() players.Player {
	return &sarsaLearner{
		sarsa: sarsa,
		lastQ: unsetState,
	}
}

func (sarsa Sarsa) Value(actState int) float32 {
	return sarsa.qf[actState]
}

func (sarsa *Sarsa) Train(episodes int) {
	pls := []*sarsaLearner{
		&sarsaLearner{sarsa: sarsa},
		&sarsaLearner{sarsa: sarsa},
	}

	templateSG, err := rules.NewGame(2)
	if err != nil {
		panic(err.Error())
	}

	for i := 0; i < episodes; i++ {
		if (i % 100000) == 0 {
			fmt.Printf("\r%2.2f%% complete", float32(i)/float32(episodes)*100)
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
		sa, _ := players.NewSimpleState(sg).AsIntWithAction(rules.Action{})
		if sa < 0 {
			panic(fmt.Sprintf("Negative state was calculated: %d", sa))
		}
		if sg.LossWasStupid {
			// This only happens if the play is something that will ALWAYS lose the game, so incur a huge penalty
			pls[(sg.Winner+1)%2].updateLearning(sg.GameEnded, sa, stupidReward)
			pls[sg.Winner].updateLearning(sg.GameEnded, sa, forfeitWinReward)
		} else {
			pls[(sg.Winner+1)%2].updateLearning(sg.GameEnded, sa, lossReward)
			pls[sg.Winner].updateLearning(sg.GameEnded, sa, winReward)
		}
	}
	fmt.Println("\r100.0% complete")
}

// PlayCard provides a suggested action for the provided state.
// If it hasn't learned anything for this state, it plays randomly.
func (sl *sarsaLearner) PlayCard(state players.SimpleState) rules.Action {
	act, _ := sl.sarsa.greedyAction(state.AsInt())
	if act == nil {
		return (&players.RandomPlayer{}).PlayCard(state)
	}
	return *act
}

// learningAction provides a suggested action for the provided state.
// However, it also assumes it's being called for each play in a game so it can update the policy.
func (sl *sarsaLearner) learningAction(game rules.Gamestate) (rules.Action, error) {
	action, sa := sl.sarsa.epsilonGreedyAction(players.NewSimpleState(game))
	sl.updateLearning(game.GameEnded, sa, noReward)
	return action, nil
}

func (sl *sarsaLearner) updateLearning(gameEnded bool, sa int, reward float32) {
	// Now save the update
	if sl.lastQ != unsetState {
		thisValue := float32(0) // If game ended, the value of the new state is 0 because it's a terminal state
		if !gameEnded {
			thisValue = sl.sarsa.Gamma * sl.sarsa.qf[sa]
		}
		sl.sarsa.qf[sl.lastQ] += sl.sarsa.Alpha * (reward + thisValue - sl.sarsa.qf[sl.lastQ])
	}
	sl.lastQ = sa
}

// epsilonGreedyAction provides a suggested action for the provided state.
// If it hasn't learned anything for this state, it plays randomly.
// It will also choose a random action with probability Epsilon. This isn't exactly
// Epsilon-greedy because it doesn't subtract the probability of the greedy action.
func (sarsa Sarsa) epsilonGreedyAction(st players.SimpleState) (rules.Action, int) {
	sNoAct := st.AsInt()
	act, sa := sarsa.greedyAction(sNoAct)
	if act == nil || rand.Float32() < sarsa.Epsilon {
		action := (&players.RandomPlayer{}).PlayCard(st)
		return action, state.IndexWithAction(sNoAct, action)
	}
	return *act, sa
}

// greedyAction returns the greedy action for the given state. (Note the argument should be a state, not an action-state.)
// Ties are broken by choosing the first option (i.e. arbitrarily in a deterministic way).
func (sarsa Sarsa) greedyAction(st int) (*rules.Action, int) {
	bestActs := []int{}
	bestActValue := float32(0)
	bestActState := 0
	for act, actState := range state.AllActionStates(st) {
		thisVal := sarsa.qf[actState]
		if thisVal > bestActValue {
			bestActValue = thisVal
			bestActs = []int{act}
			bestActState = actState
		} else if thisVal == bestActValue {
			bestActs = append(bestActs, act)
		}
	}
	if len(bestActs) > 0 {
		bestAct := rules.ActionFromInt(bestActs[0])
		return &bestAct, bestActState
	} else {
		return nil, 0
	}
}
