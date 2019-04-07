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

const unsetState = state.SpaceMagnitude

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

	for i := 0; i < episodes; i++ {
		if (i % 100000) == 0 {
			fmt.Printf("\r%2.2f%% complete", float32(i)/float32(episodes)*100)
		}

		sg, err := rules.NewGame(2)
		if err != nil {
			panic(err.Error())
		}

		pls[0].lastQ = unsetState
		pls[1].lastQ = unsetState

		s := players.NewSimpleState(sg)
		if s.OpponentCard == 0 {
			s.OpponentCard++
		}

		for !sg.GameEnded {
			action, err := pls[sg.ActivePlayer].learningAction(s, sg)
			if err != nil {
				panic(err.Error())
			}

			sg.PlayCard(action)
		}

		// Now allow both players to update based on the end of the game.
		sa, _ := s.AsIntWithAction(rules.Action{})
		if sa < 0 {
			panic(fmt.Sprintf("Negative state was calculated: %d", sa))
		}
		pls[0].updateLearning(sg.GameEnded, sa)
		pls[1].updateLearning(sg.GameEnded, sa)
	}
	fmt.Println("\r100.0% complete")
}

// PlayCard provides a suggested action for the provided state.
// If it hasn't learned anything for this state, it plays randomly.
func (sl *sarsaLearner) PlayCard(state players.SimpleState) rules.Action {
	act := sl.sarsa.greedyAction(state.AsInt())
	if act == nil {
		return (&players.RandomPlayer{}).PlayCard(state)
	}
	return *act
}

// learningAction provides a suggested action for the provided state.
// However, it also assumes it's being called for each play in a game so it can update the policy.
func (sl *sarsaLearner) learningAction(state players.SimpleState, game rules.Gamestate) (rules.Action, error) {
	action := sl.PlayCard(state)

	// Calculate the new state
	sa, _ := state.AsIntWithAction(action)
	if sa < 0 {
		return action, fmt.Errorf("Negative state was calculated: %d", sa)
	}

	sl.updateLearning(game.GameEnded, sa)

	return action, nil
}

func (sl *sarsaLearner) updateLearning(gameEnded bool, sa int) {
	// Now save the update
	if sl.lastQ != unsetState {
		reward := float32(0)
		if gameEnded {
			reward = 1.0
		}

		sl.sarsa.qf[sl.lastQ] += sl.sarsa.Alpha * (reward + sl.sarsa.Gamma*sl.sarsa.qf[sa] - sl.sarsa.qf[sl.lastQ])
	}
	sl.lastQ = sa
}

// epsilonGreedyAction provides a suggested action for the provided state.
// If it hasn't learned anything for this state, it plays randomly.
// It will also choose a random action with probability Epsilon. This isn't exactly
// Epsilon-greedy because it doesn't subtract the probability of the greedy action.
func (sarsa Sarsa) epsilonGreedyAction(state players.SimpleState) rules.Action {
	act := sarsa.greedyAction(state.AsInt())
	if act == nil || rand.Float32() < sarsa.Epsilon {
		return (&players.RandomPlayer{}).PlayCard(state)
	}
	return *act
}

// greedyAction returns the greedy action for the given state. (Note the argument should be a state, not an action-state.)
// Ties are broken by choosing the first option (i.e. arbitrarily in a deterministic way).
func (sarsa Sarsa) greedyAction(st int) *rules.Action {
	bestActs := []int{}
	bestActValue := float32(0)
	for act, actState := range state.AllActionStates(st) {
		thisVal := sarsa.qf[actState]
		if thisVal > bestActValue {
			bestActValue = thisVal
			bestActs = []int{act}
		} else if thisVal == bestActValue {
			bestActs = append(bestActs, act)
		}
	}
	if len(bestActs) > 0 {
		bestAct := rules.ActionFromInt(bestActs[0])
		return &bestAct
	} else {
		return nil
	}
}
