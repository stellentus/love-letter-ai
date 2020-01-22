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
	qf    []float32
	Alpha float32
	Gamma float32
}

func NewTD(alpha, gamma float32) *TD {
	sar := &TD{
		qf:    make([]float32, state.ActionSpaceMagnitude, state.ActionSpaceMagnitude),
		Alpha: alpha,
		Gamma: gamma,
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
	act, _ := sar.GreedyAction(state.AsIndex())
	if act == nil {
		return (&players.RandomPlayer{}).PlayCard(state)
	}
	return *act
}

// PlayCard provides a suggested action for the provided state.
// If it hasn't learned anything for this state, it plays randomly.
func (sar TD) PlayCardRand(state state.Simple, r *rand.Rand) rules.Action {
	act, _ := sar.GreedyAction(state.AsIndex())
	if act == nil {
		return (&players.RandomPlayer{}).PlayCardRand(state, r)
	}
	return *act
}

// greedyAction returns the greedy action for the given state. (Note the argument should be a state, not an action-state.)
// Ties are broken by choosing the first option (i.e. arbitrarily in a deterministic way).
func (sarsa TD) GreedyAction(st int) (*rules.Action, int) {
	bestAct := 0
	bestActValue := float32(0)
	bestActState := 0
	for act, actState := range state.AllActionStates(st) {
		thisVal := sarsa.qf[actState]
		if thisVal > bestActValue {
			bestActValue = thisVal
			bestAct = act
			bestActState = actState
		}
	}
	if bestAct != 0 {
		act := rules.ActionFromInt(bestAct)
		return &act, bestActState
	} else {
		return nil, 0
	}
}

type fileHeader struct {
	Version              uint32
	Alpha                float32
	Gamma                float32
	ActionSpaceMagnitude uint64
}

const currentFileFormatVersion = 2

func (sarsa TD) SaveToFile(path string) error {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}

	length := len(sarsa.qf)
	writer := bufio.NewWriter(file)

	err = binary.Write(writer, binary.BigEndian, fileHeader{
		Version:              currentFileFormatVersion,
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
	if header.Version != currentFileFormatVersion {
		return fmt.Errorf("Cannot load SARSA weights from version not %d (%d)", currentFileFormatVersion, header.Version)
	}
	if int(header.ActionSpaceMagnitude) != length {
		return fmt.Errorf("Cannot load SARSA weights from file size not %d (%d)", state.ActionSpaceMagnitude, header.ActionSpaceMagnitude)
	}
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
