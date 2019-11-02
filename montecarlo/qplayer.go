package montecarlo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"love-letter-ai/gamemaster"
	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/state"
	"math/rand"
	"os"
)

type QPlayer struct {
	qf      []Value
	epsilon float32
}

func NewQPlayer(epsilon float32) *QPlayer {
	return &QPlayer{
		qf:      make([]Value, state.ActionSpaceMagnitude, state.ActionSpaceMagnitude),
		epsilon: epsilon,
	}
}

func (qp *QPlayer) SetEpsilon(epsilon float32) {
	qp.epsilon = epsilon
}

func (qp *QPlayer) TrainWithPlayerPolicy(episodes int, pl players.Player) {
	for i := 0; i < episodes; i++ {
		if (i % 100000) == 0 {
			fmt.Printf("\r%2.2f%% complete", float32(i)/float32(episodes)*100)
		}

		tr, err := gamemaster.TraceOneGame(pl)
		if err != nil {
			panic(err.Error())
		}

		for _, si := range tr.StateInfos {
			qp.SaveState(si)
		}
	}
	fmt.Println("\r100.0% complete")
}

func (qp *QPlayer) TrainWithSelfPolicy(episodes int) {
	qp.TrainWithPlayerPolicy(episodes, qp)
}

// PlayCard provides a suggested action for the provided state.
// If it hasn't learned anything for this state, it plays randomly.
// It will also choose a random action with probability epsilon. This isn't exactly
// epsilon-greedy because it doesn't subtract the probability of the greedy action.
func (qp *QPlayer) PlayCard(state players.SimpleState) rules.Action {
	act := qp.policy(state.AsInt())
	if act == nil || rand.Float32() < qp.epsilon {
		return (&players.RandomPlayer{}).PlayCard(state)
	}
	return *act
}

func (qp QPlayer) Value(st int) float32 {
	cnt := float32(qp.qf[st].count)
	if cnt == 0.0 {
		return 0.0
	}
	sum := float32(qp.qf[st].sum)
	if sum == 0.0 {
		return 0.0
	}
	return sum / cnt
}

// policy returns the greedy action for the given state. (Note the argument should be a state, not an action-state.)
// Ties are broken by choosing the first option (i.e. arbitrarily in a deterministic way).
func (qp QPlayer) policy(st int) *rules.Action {
	bestActs := []int{}
	bestActValue := float32(0)
	for act, actState := range state.AllActionStates(st) {
		thisVal := qp.Value(actState)
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

func (qp *QPlayer) SaveState(si gamemaster.StateInfo) {
	s := si.ActionState
	if si.Won {
		qp.qf[s].sum++
	}
	qp.qf[s].count++

	// Check for overflow
	if qp.qf[s].count == 0xFFFF {
		qp.qf[s].sum /= 2
		qp.qf[s].count /= 2
	}
}

type fileHeader struct {
	version              uint32
	epsilon              float32
	ActionSpaceMagnitude uint64
}

func (qp QPlayer) SaveToFile(path string) error {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)

	var buf bytes.Buffer
	err = binary.Write(&buf, binary.BigEndian, fileHeader{
		version:              1,
		epsilon:              qp.epsilon,
		ActionSpaceMagnitude: state.ActionSpaceMagnitude,
	})
	if err != nil {
		return err
	}

	if _, err := writer.Write(buf.Bytes()); err != nil {
		return err
	}

	length := len(qp.qf)
	by := make([]byte, 4)
	for i, val := range qp.qf {
		if i%(length/100) == 0 {
			fmt.Printf("\rSaving %2d%%", i*100/length)
		}

		by[0] = byte(val.sum & 0xFF)
		by[1] = byte(val.sum >> 8)
		by[2] = byte(val.count & 0xFF)
		by[3] = byte(val.count >> 8)

		if _, err := writer.Write(by); err != nil {
			return err
		}
	}
	fmt.Printf("\rSaved 100%%\n")

	return writer.Flush()
}
