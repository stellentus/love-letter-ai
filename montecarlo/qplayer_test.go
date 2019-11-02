package montecarlo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLoadSave(t *testing.T) {
	path := "temp-qplayer-test-file.dat"
	epsilon := float32(0.125) // Integer can be represented exactly
	tableSize := 100
	qp := newTestQPlayer(epsilon, tableSize)

	// To avoid using much RAM, only fill the first portion (but all will be written)
	for i := range qp.qf {
		// Fill with arbitrary data
		qp.qf[i].sum = uint16(i)
		qp.qf[i].count = uint16((i * 1293) % 289)
	}

	err := qp.SaveToFile(path)
	defer os.Remove(path)
	assert.NoError(t, err)

	qp2 := newTestQPlayer(epsilon, tableSize)
	err = qp2.LoadFromFile(path)
	assert.NoError(t, err)

	assert.Equal(t, epsilon, qp2.epsilon, "Epsilon didn't save/load the same")
	assert.Equal(t, qp.qf, qp2.qf, "Table didn't save/load the same")
}

func newTestQPlayer(epsilon float32, size int) *QPlayer {
	return &QPlayer{
		qf:      make([]Value, size, size),
		epsilon: epsilon,
	}
}
