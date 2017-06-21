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
	nSamples := 1
	m, n := 10000, 6000
	kMin := 64
	kMax := 384
	kStp := 16
	α, β := 1.0, 0.0

	// allocate matrices
	a := oblas.NewMatrixMN(m, kMax)
	b := oblas.NewMatrixMN(kMax, n)
	c := oblas.NewMatrixMN(m, n)

	// generate random matrices
	for i := 0; i < m; i++ {
		for j := 0; j < kMax; j++ {
			a.Set(i, j, rand.Float64()-0.5)
		}
		for j := 0; j < n; j++ {
			c.Set(i, j, rand.Float64()-0.5)
		}
	}
	for i := 0; i < kMax; i++ {
		for j := 0; j < n; j++ {
			b.Set(i, j, rand.Float64()-0.5)
		}
	}

	// run first to "warm-up"
	oblas.Dgemm(false, false, 2, 2, 2, α, a, 2, b, 2, β, c, 2)
	oblas.Dgemm(false, false, 4, 4, 4, α, a, 4, b, 4, β, c, 4)
	oblas.Dgemm(false, false, 8, 8, 8, α, a, 8, b, 8, β, c, 8)

	// data for plotting
	idx := 0
	npts := (kMax-kMin)/kStp + 1
	xx, yy := make([]float64, npts), make([]float64, npts)

	// run all sizes
	for k := kMin; k <= kMax; k += kStp {

		// run benchmark
		lda, ldb, ldc := m, k, m
		t0 := time.Now()
		for l := 0; l < nSamples; l++ {
			oblas.Dgemm(false, false, m, n, k, α, a, lda, b, ldb, β, c, ldc)
		}

		// compute MFlops
		dt := time.Now().Sub(t0) / time.Duration(nSamples)
		dtMicros := float64(dt.Nanoseconds()) * 1e-3
		mflops := 2.0 * float64(m) * float64(k) * float64(n) / dtMicros
		gflops := mflops * 1e-3

		// print message
		io.Pf("%4d × %4d k=%3d: %10.2f GFlops   %v\n", m, n, k, gflops, dt)

		// set data for plot
		xx[idx] = float64(k)
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
		plt.Gll("$k$", "$GFlops$", nil)
		plt.Save("/tmp/gosl/oblas", "benchDgemmBig")
	}
}
