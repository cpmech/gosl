// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math/rand"
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la/oblas"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func MatMul(c [][]float64, α float64, a, b [][]float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b[0]); j++ {
			c[i][j] = 0.0
			for k := 0; k < len(a[0]); k++ {
				c[i][j] += α * a[i][k] * b[k][j]
			}
		}
	}
}

func main() {

	// set number of threads
	oblas.SetNumThreads(1)

	// run small values
	mValues := utl.IntRange3(2, 66, 2)
	nSamples := 1000
	bench("oblas-dgemm01a", nSamples, mValues)

	io.Pl()

	// run larger values
	mValues = utl.IntRange3(16, 1424, 64)
	nSamples = 10
	bench("oblas-dgemm01b", nSamples, mValues)
}

func bench(fnkey string, nSamples int, mValues []int) {

	// constants
	α, β := 1.0, 0.0
	_, mMax := utl.IntMinMax(mValues)

	// Dgemm: allocate matrices
	a := oblas.NewMatrixMN(mMax, mMax)
	b := oblas.NewMatrixMN(mMax, mMax)
	c := oblas.NewMatrixMN(mMax, mMax)

	// Dgemm: generate random matrices
	for j := 0; j < mMax; j++ {
		for i := 0; i < mMax; i++ {
			a.Set(i, j, rand.Float64()-0.5)
			b.Set(i, j, rand.Float64()-0.5)
			c.Set(i, j, rand.Float64()-0.5)
		}
	}

	// Dgemm: run first to "warm-up"
	oblas.Dgemm(false, false, 2, 2, 2, α, a, 2, b, 2, β, c, 2)
	oblas.Dgemm(false, false, 4, 4, 4, α, a, 4, b, 4, β, c, 4)
	oblas.Dgemm(false, false, 8, 8, 8, α, a, 8, b, 8, β, c, 8)

	// data for plotting
	idx, naiveIdx := 0, 0
	npts := len(mValues) + 1
	xx, yy := make([]float64, npts), make([]float64, npts)
	naiveX, naiveY := make([]float64, npts), make([]float64, npts)

	// header
	io.Pf("   size   ┃     OpenBLAS dgemm     (Dt) ┃          naïve           (naiveDt) ┃ naiveDt/Dt\n")
	io.Pf("━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━\n")

	// run all sizes
	for _, m := range mValues {

		// Dgemm: run benchmark
		t0 := time.Now()
		for l := 0; l < nSamples; l++ {
			oblas.Dgemm(false, false, m, m, m, α, a, m, b, m, β, c, m)
		}

		// Dgemm: compute MFlops
		dt := time.Now().Sub(t0) / time.Duration(nSamples)
		dtMicros := float64(dt.Nanoseconds()) * 1e-3
		mflops := 2.0 * float64(m) * float64(m) * float64(m) / dtMicros
		gflops := mflops * 1e-3

		// Dgemm: set data for plot
		xx[idx] = float64(m)
		yy[idx] = gflops
		idx++

		// ------------------------------- Naïve

		var naiveDt time.Duration
		var naiveGflops float64
		if m <= 720 {

			// Naive: allocate matrices
			A := make([][]float64, m)
			B := make([][]float64, m)
			C := make([][]float64, m)

			// Naive: generate random matrices
			for i := 0; i < m; i++ {
				A[i] = make([]float64, m)
				B[i] = make([]float64, m)
				C[i] = make([]float64, m)
				for j := 0; j < m; j++ {
					A[i][j] = rand.Float64() - 0.5
					B[i][j] = rand.Float64() - 0.5
					C[i][j] = rand.Float64() - 0.5
				}
			}

			// Naive: run benchmark
			naiveT0 := time.Now()
			for l := 0; l < nSamples; l++ {
				MatMul(C, α, A, B)
			}

			// Naive: compute MFlops
			naiveDt = time.Now().Sub(naiveT0) / time.Duration(nSamples)
			naiveDtMicros := float64(naiveDt.Nanoseconds()) * 1e-3
			naiveMflops := 2.0 * float64(m) * float64(m) * float64(m) / naiveDtMicros
			naiveGflops = naiveMflops * 1e-3

			// Naive: set data for plot
			naiveX[naiveIdx] = float64(m)
			naiveY[naiveIdx] = naiveGflops
			naiveIdx++

			// ------------------------------- message

			// print message
			io.Pf("%4d×%4d ┃ %5.2f GFlops (%12v) ┃ naive: %5.2f GFlops (%12v) ┃ %.3f \n", m, m, gflops, dt, naiveGflops, naiveDt, naiveDtMicros/dtMicros)

		} else {

			// ------------------------------- message

			// print message
			io.Pf("%4d×%4d ┃ %5.2f GFlops (%12v) ┃ naive: N/A                         ┃ N/A\n", m, m, gflops, dt)

		}
	}

	// plot
	if true {
		plt.Reset(true, &plt.A{WidthPt: 450})
		plt.Plot(xx[:idx], yy[:idx], &plt.A{C: "#1549bd", M: ".", L: "dgemm"})
		plt.Plot(naiveX[:naiveIdx], naiveY[:naiveIdx], &plt.A{C: "r", M: ".", L: "naive"})
		plt.Title(io.Sf("OpenBLAS $versus$ Naive MatMul. Single-threaded. nSamples=%d", nSamples), &plt.A{Fsz: 9})
		plt.SetTicksXlist(xx)
		plt.SetTicksRotationX(45)
		plt.SetYnticks(16)
		plt.Gll("$m$", "$GFlops$", &plt.A{LegOut: false, LegNcol: 2})
		plt.Save("/tmp/gosl/bench", fnkey)
	}
}
