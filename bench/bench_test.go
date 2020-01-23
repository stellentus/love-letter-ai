package bench

import (
	"testing"

	"love-letter-ai/players"
	"love-letter-ai/td"
)

func BenchmarkQLearner1(b *testing.B) {
	players.Runners = 1
	sar := td.NewTD(0.3, 1)
	pls := []players.TrainingPlayer{
		sar.QLearner(),
		sar.QLearner(),
	}

	b.ResetTimer()
	players.Train(pls, b.N, 1)
}

func BenchmarkQLearner2(b *testing.B) {
	players.Runners = 2
	sar := td.NewTD(0.3, 1)
	pls := []players.TrainingPlayer{
		sar.QLearner(),
		sar.QLearner(),
	}

	b.ResetTimer()
	players.Train(pls, b.N, 1)
}

func BenchmarkQLearner4(b *testing.B) {
	players.Runners = 4
	sar := td.NewTD(0.3, 1)
	pls := []players.TrainingPlayer{
		sar.QLearner(),
		sar.QLearner(),
	}

	b.ResetTimer()
	players.Train(pls, b.N, 1)
}

func BenchmarkQLearner8(b *testing.B) {
	players.Runners = 8
	sar := td.NewTD(0.3, 1)
	pls := []players.TrainingPlayer{
		sar.QLearner(),
		sar.QLearner(),
	}

	b.ResetTimer()
	players.Train(pls, b.N, 1)
}

func BenchmarkQLearner16(b *testing.B) {
	players.Runners = 16
	sar := td.NewTD(0.3, 1)
	pls := []players.TrainingPlayer{
		sar.QLearner(),
		sar.QLearner(),
	}

	b.ResetTimer()
	players.Train(pls, b.N, 1)
}

func BenchmarkQLearner32(b *testing.B) {
	players.Runners = 32
	sar := td.NewTD(0.3, 1)
	pls := []players.TrainingPlayer{
		sar.QLearner(),
		sar.QLearner(),
	}

	b.ResetTimer()
	players.Train(pls, b.N, 1)
}
