package td

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLoadSave(t *testing.T) {
	path := "temp-td-test-file.dat"
	alpha := float32(0.5)  // Integer can be represented exactly
	gamma := float32(0.25) // Integer can be represented exactly
	tableSize := 100
	td := newTestTDlayer(alpha, gamma, tableSize)

	// To avoid using much RAM, only fill the first portion (but all will be written)
	for i := range td.qf {
		// Fill with arbitrary data
		td.qf[i] = float32(i*208284) / 7282
	}

	err := td.SaveToFile(path)
	defer os.Remove(path)
	assert.NoError(t, err)

	td2 := newTestTDlayer(alpha, gamma, tableSize)
	err = td2.LoadFromFile(path)
	assert.NoError(t, err)

	assert.Equal(t, alpha, td2.Alpha, "Alpha didn't save/load the same")
	assert.Equal(t, gamma, td2.Gamma, "Gamma didn't save/load the same")
	assert.Equal(t, td.qf, td2.qf, "Table didn't save/load the same")
}

func newTestTDlayer(alpha, gamma float32, size int) *TD {
	return &TD{
		qf:    make([]float32, size, size),
		Alpha: alpha,
		Gamma: gamma,
	}
}
