// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"math/rand"
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la/mkl"
	"github.com/cpmech/gosl/la/oblas"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// set number of threads
	mkl.SetNumThreads(1)
	oblas.SetNumThreads(1)

	// run small values
	mValues := utl.IntRange3(2, 66, 2)
	nSamples := 1000
	bench("mkl-dgemm01a", nSamples, mValues)

	io.Pl()

	// run larger values
	mValues = utl.IntRange3(16, 1424, 64)
	nSamples = 100
	bench("mkl-dgemm01b", nSamples, mValues)
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

	// data for plotting
	buf := new(bytes.Buffer)
	io.Ff(buf, "%4s %4s %23s %23s\n", "m", "n", "Gflops", "DtMicros")

	// header
	io.Pf("   size   ┃            MKL           (mklDt) ┃ mklDt (μs)\n")
	io.Pf("━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━\n")

	// run all sizes
	for _, m := range mValues {

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

		// print message
		io.Pf("%4d×%4d ┃ mkl: %5.2f GFlops (%12v) ┃ %.3f \n", m, m, mklGflops, mklDt, mklDtMicros)

		// save buffer
		io.Ff(buf, "%4d %4d %23.15e %23.15e\n", m, m, mklGflops, mklDtMicros)
	}

	// save file
	io.WriteFileVD("/tmp/gosl/", io.Sf("%s-%dsamples.res", fnkey, nSamples), buf)
}
