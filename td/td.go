package td

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/state"
	"math/rand"
	"os"
)

type TD struct {
	qf      []float32
	Epsilon float32
	Alpha   float32
	Gamma   float32
}

func NewTD(epsilon, alpha, gamma float32) *TD {
	sar := &TD{
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

func (sarsa TD) Value(actState int) float32 {
	return sarsa.qf[actState]
}

// PlayCard provides a suggested action for the provided state.
// If it hasn't learned anything for this state, it plays randomly.
func (sar TD) PlayCard(state state.Simple) rules.Action {
	act, _ := sar.greedyAction(state.AsIndex())
	if act == nil {
		return (&players.RandomPlayer{}).PlayCard(state)
	}
	return *act
}

// epsilonGreedyAction provides a suggested action for the provided state.
// If it hasn't learned anything for this state, it plays randomly.
// It will also choose a random action with probability Epsilon. This isn't exactly
// Epsilon-greedy because it doesn't subtract the probability of the greedy action.
func (sarsa TD) epsilonGreedyAction(st state.Simple) (rules.Action, int) {
	sNoAct := st.AsIndex()
	act, sa := sarsa.greedyAction(sNoAct)
	if act == nil || rand.Float32() < sarsa.Epsilon {
		action := (&players.RandomPlayer{}).PlayCard(st)
		return action, state.IndexWithAction(sNoAct, action)
	}
	return *act, sa
}

// greedyAction returns the greedy action for the given state. (Note the argument should be a state, not an action-state.)
// Ties are broken by choosing the first option (i.e. arbitrarily in a deterministic way).
func (sarsa TD) greedyAction(st int) (*rules.Action, int) {
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

type fileHeader struct {
	Version              uint32
	Epsilon              float32
	Alpha                float32
	Gamma                float32
	ActionSpaceMagnitude uint64
}

func (sarsa TD) SaveToFile(path string) error {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}

	length := len(sarsa.qf)
	writer := bufio.NewWriter(file)

	err = binary.Write(writer, binary.BigEndian, fileHeader{
		Version:              1,
		Epsilon:              sarsa.Epsilon,
		Alpha:                sarsa.Alpha,
		Gamma:                sarsa.Gamma,
		ActionSpaceMagnitude: uint64(length),
	})
	if err != nil {
		return err
	}

	for i, val := range sarsa.qf {
		if i%(length/100) == 0 {
			fmt.Printf("\rSaving %2d%%", i*100/length)
		}
		if err := binary.Write(writer, binary.BigEndian, val); err != nil {
			return err
		}
	}
	fmt.Printf("\rSaved 100%%\n")

	return writer.Flush()
}

func (sarsa *TD) LoadFromFile(path string) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}

	length := len(sarsa.qf)

	reader := bufio.NewReader(file)
	header := &fileHeader{}
	if err = binary.Read(reader, binary.BigEndian, header); err != nil {
		return err
	}
	if header.Version != 1 {
		return fmt.Errorf("Cannot load SARSA weights from version not 1 (%d)", header.Version)
	}
	if int(header.ActionSpaceMagnitude) != length {
		return fmt.Errorf("Cannot load SARSA weights from file size not %d (%d)", state.ActionSpaceMagnitude, header.ActionSpaceMagnitude)
	}
	sarsa.Epsilon = header.Epsilon
	sarsa.Alpha = header.Alpha
	sarsa.Gamma = header.Gamma

	for i := range sarsa.qf {
		if i%(length/100) == 0 {
			fmt.Printf("\rLoading %2d%%", i*100/length)
		}
		if err := binary.Read(reader, binary.BigEndian, &sarsa.qf[i]); err != nil {
			return nil
		}
	}
	fmt.Printf("\rLoaded 100%%\n")

	return nil
}
