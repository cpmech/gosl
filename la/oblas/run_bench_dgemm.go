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
)

func main() {

	// set number of threads
	oblas.SetNumThreads(1)

	// constants
	nSamples := 10
	mMin := 8
	mMax := 1064
	mStp := 48
	useN := false
	nVal := 1
	α, β := 1.0, 0.0

	// allocate matrices
	a := oblas.NewMatrixMN(mMax, mMax)
	b := oblas.NewMatrixMN(mMax, mMax)
	c := oblas.NewMatrixMN(mMax, mMax)

	// generate random matrices
	for j := 0; j < mMax; j++ {
		for i := 0; i < mMax; i++ {
			a.Set(i, j, rand.Float64()-0.5)
			b.Set(i, j, rand.Float64()-0.5)
			c.Set(i, j, rand.Float64()-0.5)
		}
	}

	// run first to "warm-up"
	oblas.Dgemm(false, false, 2, 2, 2, α, a, 2, b, 2, β, c, 2)
	oblas.Dgemm(false, false, 4, 4, 4, α, a, 4, b, 4, β, c, 4)
	oblas.Dgemm(false, false, 8, 8, 8, α, a, 8, b, 8, β, c, 8)

	// data for plotting
	idx := 0
	npts := (mMax-mMin)/mStp + 1
	xx, yy := make([]float64, npts), make([]float64, npts)

	// run all sizes
	for m := mMin; m <= mMax; m += mStp {

		// set n
		n := m
		if useN {
			n = nVal
		}

		// run benchmark
		t0 := time.Now()
		for l := 0; l < nSamples; l++ {
			oblas.Dgemm(false, false, m, n, m, α, a, m, b, m, β, c, m)
		}

		// compute MFlops
		dt := time.Now().Sub(t0) / time.Duration(nSamples)
		dtMicros := float64(dt.Nanoseconds()) * 1e-3
		mflops := 2.0 * float64(m) * float64(m) * float64(n) / dtMicros
		gflops := mflops * 1e-3

		// print message
		io.Pf("%4d × %4d: %10.2f GFlops   %v\n", m, n, gflops, dt)

		// set data for plot
		xx[idx] = float64(m)
		yy[idx] = gflops
		idx++
	}

	// plot
	if true {
		plt.Reset(true, nil)
		plt.Plot(xx, yy, &plt.A{M: "."})
		plt.SetTicksXlist(xx)
		plt.SetTicksRotationX(45)
		plt.SetYnticks(16)
		plt.Gll("$m$", "$GFlops$", nil)
		plt.Save("/tmp/gosl/oblas", "benchDgemm")
	}
}
