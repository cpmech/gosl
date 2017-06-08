// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"testing"

	"github.com/cpmech/gosl/fun/fftw"
)

var benchmarkingData []float64

func init() {
	N := 1 << 20 // 2²⁰ = 1,048,576
	benchmarkingData = make([]float64, N*2)
	for i := 0; i < N; i++ {
		ii := float64(i * 2)
		benchmarkingData[i*2] = ii + 1
		benchmarkingData[i*2+1] = ii + 2
	}
}

func BenchmarkFFT01(b *testing.B) {
	inverse, inplace, measure := false, false, false
	plan, _ := fftw.NewPlan1d(benchmarkingData, 0, inverse, inplace, measure)
	defer plan.Free()
	for i := 0; i < b.N; i++ {
		plan.Execute()
	}
}

func BenchmarkFourierTransLL01(b *testing.B) {
	inverse := false
	x := make([]float64, len(benchmarkingData))
	copy(x, benchmarkingData)
	for i := 0; i < b.N; i++ {
		FourierTransLL(x, inverse)
	}
}
