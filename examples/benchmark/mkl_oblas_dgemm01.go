// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math/rand"
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la/mkl"
	"github.com/cpmech/gosl/la/oblas"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// set number of threads
	mkl.SetNumThreads(1)
	oblas.SetNumThreads(1)

	// run small values
	mValues := utl.IntRange3(2, 66, 2)
	nSamples := 1000
	bench("mkl-oblas-dgemm01a", nSamples, mValues)

	io.Pl()

	// run larger values
	mValues = utl.IntRange3(16, 1424, 64)
	nSamples = 100
	bench("mkl-oblas-dgemm01b", nSamples, mValues)
}

func bench(fnkey string, nSamples int, mValues []int) {

	// constants
	α, β := 1.0, 0.0
	_, mMax := utl.IntMinMax(mValues)

	// Dgemm: allocate matrices
	a := make([]float64, mMax*mMax)
	b := make([]float64, mMax*mMax)
	c := make([]float64, mMax*mMax)

	// Dgemm: generate random matrices
	for j := 0; j < mMax; j++ {
		for i := 0; i < mMax; i++ {
			a[i+j*mMax] = rand.Float64() - 0.5
			b[i+j*mMax] = rand.Float64() - 0.5
			c[i+j*mMax] = rand.Float64() - 0.5
		}
	}

	// Dgemm: run first to "warm-up"
	mkl.Dgemm(false, false, 2, 2, 2, α, a, 2, b, 2, β, c, 2)
	mkl.Dgemm(false, false, 4, 4, 4, α, a, 4, b, 4, β, c, 4)
	mkl.Dgemm(false, false, 8, 8, 8, α, a, 8, b, 8, β, c, 8)
	oblas.Dgemm(false, false, 2, 2, 2, α, a, 2, b, 2, β, c, 2)
	oblas.Dgemm(false, false, 4, 4, 4, α, a, 4, b, 4, β, c, 4)
	oblas.Dgemm(false, false, 8, 8, 8, α, a, 8, b, 8, β, c, 8)

	// data for plotting
	idx, mklIdx := 0, 0
	npts := len(mValues) + 1
	xx, yy := make([]float64, npts), make([]float64, npts)
	mklX, mklY := make([]float64, npts), make([]float64, npts)

	// header
	io.Pf("   size   ┃     OpenBLAS dgemm     (Dt) ┃            MKL           (mklDt) ┃ Dt/mklDt\n")
	io.Pf("━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━\n")

	// run all sizes
	for _, m := range mValues {

		// oblas: run benchmark
		t0 := time.Now()
		for l := 0; l < nSamples; l++ {
			oblas.Dgemm(false, false, m, m, m, α, a, m, b, m, β, c, m)
		}

		// oblas: compute MFlops
		dt := time.Now().Sub(t0) / time.Duration(nSamples)
		dtMicros := float64(dt.Nanoseconds()) * 1e-3
		mflops := 2.0 * float64(m) * float64(m) * float64(m) / dtMicros
		gflops := mflops * 1e-3

		// oblas: set data for plot
		xx[idx] = float64(m)
		yy[idx] = gflops
		idx++

		// ------------------------------- MKL

		// mkl: run benchmark
		mklT0 := time.Now()
		for l := 0; l < nSamples; l++ {
			mkl.Dgemm(false, false, m, m, m, α, a, m, b, m, β, c, m)
		}

		// mkl: compute MFlops
		mklDt := time.Now().Sub(mklT0) / time.Duration(nSamples)
		mklDtMicros := float64(mklDt.Nanoseconds()) * 1e-3
		mklMflops := 2.0 * float64(m) * float64(m) * float64(m) / mklDtMicros
		mklGflops := mklMflops * 1e-3

		// mkl: set data for plot
		mklX[mklIdx] = float64(m)
		mklY[mklIdx] = mklGflops
		mklIdx++

		// ------------------------------- message

		// print message
		io.Pf("%4d×%4d ┃ %5.2f GFlops (%12v) ┃ mkl: %5.2f GFlops (%12v) ┃ %.3f \n", m, m, gflops, dt, mklGflops, mklDt, dtMicros/mklDtMicros)
	}

	// plot
	if true {
		plt.Reset(true, &plt.A{WidthPt: 450})
		plt.Plot(xx[:idx], yy[:idx], &plt.A{C: "#1549bd", M: ".", L: "OpenBLAS"})
		plt.Plot(mklX[:mklIdx], mklY[:mklIdx], &plt.A{C: "r", M: ".", L: "MKL"})
		plt.Title(io.Sf("OpenBLAS $versus$ MKL MatMul. Single-threaded. nSamples=%d", nSamples), &plt.A{Fsz: 9})
		plt.SetTicksXlist(xx)
		plt.SetTicksRotationX(45)
		plt.SetYnticks(16)
		plt.Gll("$m$", "$GFlops$", &plt.A{LegOut: false, LegNcol: 2})
		plt.Save("/tmp/gosl/bench", fnkey)
	}
}
