package montecarlo

import (
	"testing"
)

func BenchmarkEvenValueFunction(b *testing.B) {
	for n := 0; n < b.N; n++ {
		EvenValueFunction()
	}
}
