package fun

import (
	"math/rand"
	"testing"
)

var (
	benchmarkingXX []float64
	benchmarkingYY []float64
	benchmarkingF  []float64

	// Testing for scaling
	benchmarkingXX2 []float64
	benchmarkingYY2 []float64
	benchmarkingF2  []float64

	multInterpRes float64
)

func init() {
	N := 1 << 8
	M := 1 << 9

	benchmarkingXX = make([]float64, N)
	benchmarkingYY = make([]float64, N)
	benchmarkingF = make([]float64, N*N)

	// Testing for scaling
	benchmarkingXX2 = make([]float64, M)
	benchmarkingYY2 = make([]float64, M)
	benchmarkingF2 = make([]float64, M*M)

	for i := 0; i < N; i++ {
		benchmarkingXX[i] = 100 * float64(i) / float64(N)
		benchmarkingYY[i] = 100 * float64(i) / float64(N)
	}

	// Testing for scaling
	for i := 0; i < M; i++ {
		benchmarkingXX2[i] = 100 * float64(i) / float64(M)
		benchmarkingYY2[i] = 100 * float64(i) / float64(M)
	}

	for i := 0; i < (N * N); i++ {
		benchmarkingF[i] = 5000 * rand.Float64()
	}

	for i := 0; i < (M * M); i++ {
		benchmarkingF2[i] = 5000 * rand.Float64()
	}

}

func BenchmarkMultInterpSpeed(b *testing.B) {
	o := NewBiLinear(benchmarkingF2, benchmarkingXX2, benchmarkingYY2)
	var res float64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res = o.P(rand.Float64()*100, rand.Float64()*100)
	}
	multInterpRes = res
}

func BenchmarkMultInterp(b *testing.B) {
	o := NewBiLinear(benchmarkingF, benchmarkingXX, benchmarkingYY)
	var res float64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res = interpoRunBiLinearSearch(o)
	}
	multInterpRes = res
}

func BenchmarkMultInterpScale(b *testing.B) {
	o := NewBiLinear(benchmarkingF2, benchmarkingXX2, benchmarkingYY2)
	var res float64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res = interpoRunBiLinearSearch(o)
	}
	multInterpRes = res
}

func BenchmarkMultInterpNoHunt(b *testing.B) {
	o := NewBiLinear(benchmarkingF, benchmarkingXX, benchmarkingYY)
	o.SetDisableHunt(true)
	var res float64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res = interpoRunBiLinearSearch(o)
	}
	multInterpRes = res
}

func interpoRunBiLinearSearch(o *BiLinear) (res float64) {
	Mseq := 1000
	Mrnd := 100
	for j := 0; j < Mseq; j++ {
		x := 100 * float64(j) / float64(Mseq)
		y := 100 * float64(j) / float64(Mseq)
		res = o.P(x, y)
	}
	for j := 0; j < Mrnd; j++ {
		x := 100 * rand.Float64()
		y := 100 * rand.Float64()
		res = o.P(x, y)
	}
	return
}
