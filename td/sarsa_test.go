package td

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLoadSave(t *testing.T) {
	path := "temp-td-test-file.dat"
	epsilon := float32(0.125) // Integer can be represented exactly
	alpha := float32(0.5)     // Integer can be represented exactly
	gamma := float32(0.25)    // Integer can be represented exactly
	tableSize := 100
	sarsa := newTestSarsalayer(epsilon, alpha, gamma, tableSize)

	// To avoid using much RAM, only fill the first portion (but all will be written)
	for i := range sarsa.qf {
		// Fill with arbitrary data
		sarsa.qf[i] = float32(i*208284) / 7282
	}

	err := sarsa.SaveToFile(path)
	defer os.Remove(path)
	assert.NoError(t, err)

	sarsa2 := newTestSarsalayer(epsilon, alpha, gamma, tableSize)
	err = sarsa2.LoadFromFile(path)
	assert.NoError(t, err)

	assert.Equal(t, epsilon, sarsa2.Epsilon, "Epsilon didn't save/load the same")
	assert.Equal(t, alpha, sarsa2.Alpha, "Alpha didn't save/load the same")
	assert.Equal(t, gamma, sarsa2.Gamma, "Gamma didn't save/load the same")
	assert.Equal(t, sarsa.qf, sarsa2.qf, "Table didn't save/load the same")
}

func newTestSarsalayer(epsilon, alpha, gamma float32, size int) *Sarsa {
	return &Sarsa{
		qf:      make([]float32, size, size),
		Epsilon: epsilon,
		Alpha:   alpha,
		Gamma:   gamma,
	}
}
