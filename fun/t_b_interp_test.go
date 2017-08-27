// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math/rand"
	"testing"

	"github.com/cpmech/gosl/plt"
)

var (
	benchmarkingX   []float64
	benchmarkingY   []float64
	benchmarkingRes float64
)

func init() {
	N := 1 << 10
	//N := 1 << 20 // 2²⁰ = 1,048,576
	benchmarkingX = make([]float64, N)
	benchmarkingY = make([]float64, N)
	for i := 0; i < N; i++ {
		benchmarkingX[i] = 100 * float64(i) / float64(N)
		benchmarkingY[i] = -5.0 + 10.0*rand.Float64()
	}
	if false {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.Plot(benchmarkingX, benchmarkingY, &plt.A{C: "r", Ls: "-", M: ".", NoClip: true})
		plt.Gll("x", "y", nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/fun", "benchInterp01")
	}
}

func BenchmarkInterp(b *testing.B) {
	o, _ := NewInterpolator("lin", 1, benchmarkingX, benchmarkingY)
	var res float64
	for i := 0; i < b.N; i++ {
		res = interpoRunSearch(o)
	}
	benchmarkingRes = res
}

func BenchmarkInterpNoHunt(b *testing.B) {
	o, _ := NewInterpolator("lin", 1, benchmarkingX, benchmarkingY)
	o.DisableHunt = true
	var res float64
	for i := 0; i < b.N; i++ {
		res = interpoRunSearch(o)
	}
	benchmarkingRes = res
}

func interpoRunSearch(o *Interpolator) (res float64) {
	Mseq := 1000
	Mrnd := 100
	for j := 0; j < Mseq; j++ {
		x := 100 * float64(j) / float64(Mseq)
		res = o.P(x)
	}
	for j := 0; j < Mrnd; j++ {
		x := 100 * rand.Float64()
		res = o.P(x)
	}
	return
}
