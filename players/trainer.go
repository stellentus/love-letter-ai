package players

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"

	"love-letter-ai/rules"
	"love-letter-ai/state"
)

type TrainingPlayer interface {
	Player

	// GreedyAction returns the greedy action for the given state, along with the corresponding state-action.
	// The action may be nil if nothing has trained yet.
	GreedyAction(state int) (*rules.Action, int)

	// UpdateQ is called to update the player's state-action values.
	UpdateQ(gameEnded bool, lastQ, sa int, reward float32)

	// Finalize is called at the end in case any cleanup is necessary.
	Finalize()
}

type trainer struct {
	// tp is the player model being trained
	tp TrainingPlayer

	// qStates is a slice of state-action ints representing the states seen
	// (and actions chosen) so far by this player.
	qStates []int
}

const (
	unsetState       = state.ActionSpaceMagnitude
	winReward        = 100
	stupidReward     = -100
	forfeitWinReward = 1 // Only a minor benefit for winning because the other player was an idiot
	noReward         = 0
	lossReward       = -0.1 // The penalty for losing is minor since it might not have been the player's fault
	chunkSize        = 1000
)

var (
	Runners = runtime.GOMAXPROCS(0)
	Output  = true
)

func Train(pls []TrainingPlayer, episodes int, epsilon float64) {
	trs := make([]trainer, len(pls))
	for i, pl := range pls {
		trs[i] = trainer{tp: pl}
	}

	wg := sync.WaitGroup{}
	in := make(chan int)
	out := make(chan int)

	for i := 0; i < Runners; i++ {
		wg.Add(1)
		go func() {
			r := rand.New(rand.NewSource(int64(i)))
			for games := range in {
				templateSG, err := rules.NewGame(2, r)
				if err != nil {
					panic(err.Error())
				}
				for i := 0; i < games; i++ {
					sg := templateSG.Copy()

					trs[0].qStates = make([]int, 0, 8) // I think maximum number of turns is 6, but whatever
					trs[1].qStates = make([]int, 0, 8)

					for !sg.GameEnded {
						action, err := trs[sg.ActivePlayer].learningAction(sg, epsilon, r)
						if err != nil {
							panic(err.Error())
						}

						sg.PlayCard(action, r)
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
				out <- games
			}
			wg.Done()
		}()
	}

	epPrintMod := episodes / 100000
	if epPrintMod < 1 {
		epPrintMod = 1
	}

	go func() { status(episodes, epPrintMod, out); wg.Done() }()

	for episodes > 0 {
		if episodes > chunkSize {
			in <- chunkSize
		} else {
			in <- episodes
		}
		episodes -= chunkSize
	}
	close(in)

	wg.Wait()
	wg.Add(1)
	close(out)
	wg.Wait()

	trs[0].tp.Finalize()
	trs[1].tp.Finalize()

	if Output {
		fmt.Fprintln(os.Stderr, "\r100.0% complete")
	}
}

func status(episodes, epPrintMod int, ch chan int) {
	count := 0
	current := 0
	for i := range ch {
		count += i
		current += i
		if current >= epPrintMod && Output {
			fmt.Fprintf(os.Stderr, "\r%2.2f%% complete", float32(count)/float32(episodes)*100)
			current -= epPrintMod
		}
	}
}

// epsilonGreedyAction provides a suggested action for the provided state.
// If it hasn't learned anything for this state, it plays randomly.
// It will also choose a random action with probability Epsilon. This isn't exactly
// Epsilon-greedy because it doesn't subtract the probability of the greedy action.
func epsilonGreedyAction(pl TrainingPlayer, st state.Simple, epsilon float64, r *rand.Rand) (rules.Action, int) {
	sNoAct := st.AsIndex()
	act, sa := pl.GreedyAction(sNoAct)
	if act == nil || r.Float64() < epsilon {
		action := (&RandomPlayer{}).PlayCardRand(st, r)
		return action, state.IndexWithAction(sNoAct, action)
	}
	return *act, sa
}

// learningAction provides a suggested action for the provided state.
// However, it also assumes it's being called for each play in a game so it can update the policy.
func (tr *trainer) learningAction(game rules.Gamestate, epsilon float64, r *rand.Rand) (rules.Action, error) {
	action, sa := epsilonGreedyAction(tr.tp, state.NewSimple(game), epsilon, r)
	tr.updateQ(game.GameEnded, sa, noReward)
	return action, nil
}

func (tr *trainer) updateQ(gameEnded bool, sa int, reward float32) {
	tr.qStates = append(tr.qStates, sa)

	numStates := len(tr.qStates)
	// Now save the update
	if numStates > 1 {
		tr.tp.UpdateQ(gameEnded, tr.qStates[numStates-2], sa, reward)
	}
}
