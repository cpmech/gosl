// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/fun/fftw"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
)

func ScaledMflops(N int, dt time.Duration) float64 {
	n := float64(N)
	dtMicroseconds := float64(dt.Nanoseconds()) * 1e-3
	return 5.0 * n * math.Log2(n) / dtMicroseconds
}

func main() {

	// set seed
	rnd.Init(0)

	// size
	//N := 1 << 4 // 2⁴ = 512
	//N := 1 << 10 // 2¹⁰ = 1024
	//N := 1 << 11 // 2¹⁰ = 1024
	//N := 1 << 18 // 2¹⁸ = 262144
	N := 1 << 20 // 2²⁰ = 1,048,576
	//N := 1 << 25 // 2²⁵ = 33,554,432
	io.Pf("len(complex(x)):   N = %v\n", N)

	// allocate array
	t0 := time.Now()
	data0 := make([]float64, N*2)
	io.Pf("alloc data (N*2)   dt = %v\n", time.Now().Sub(t0))

	// random array
	t0 = time.Now()
	rnd.Float64s(data0, -0.5, 0.5)
	io.Pf("generate data:     dt = %v\n", time.Now().Sub(t0))

	// copy data
	data1 := make([]float64, len(data0))
	t0 = time.Now()
	copy(data1, data0)
	io.Pf("copy data:         dt = %v\n", time.Now().Sub(t0))

	// benchmark slow method
	io.Pl()
	var X []complex128
	if N < 30000 {
		x := la.RCpairsToComplex(data0)
		t0 = time.Now()
		X = fun.DftSlow(x)
		io.Pf("DftSlow            dt = %v\n", time.Now().Sub(t0))
	}

	// benchmark Go FFT code
	inverse := false
	t0 = time.Now()
	err := fun.FourierTransLL(data0, inverse)
	if err != nil {
		chk.Panic("FourierTransLL failed:\n%v\n", err)
	}
	tf := time.Now()
	io.Pf("FourierTransLL     dt = %v\n", tf.Sub(t0))
	io.Pf("FourierTransLL mflops = %.1f\n", ScaledMflops(N, tf.Sub(t0)))

	// allocate plan
	t0 = time.Now()
	inplace, measure := true, false
	plan, err := fftw.NewPlan1d(data1, 0, inverse, inplace, measure)
	if err != nil {
		chk.Panic("fftw failed:\n%v\n", err)
	}
	defer plan.Free()
	io.Pl()
	io.Pf("FFTW: allocate     dt = %v\n", time.Now().Sub(t0))

	// benchmark FFTW
	t1 := time.Now()
	plan.Execute()
	tf = time.Now()
	io.Pf("FFTW: execute      dt = %v\n", tf.Sub(t1))
	io.Pf("FFTW: total        dt = %v\n", tf.Sub(t0))
	io.Pf("FFTW:          mflops = %.1f\n", ScaledMflops(N, tf.Sub(t1)))

	// check
	if true && len(X) > 0 {
		io.Pl()
		tst := new(testing.T)
		Y := la.RCpairsToComplex(data0)
		t0 = time.Now()
		chk.VectorC(tst, "data0", 1e-10, X, Y)
		io.Pf("error check:       dt = %v\n", time.Now().Sub(t0))
		Y = la.RCpairsToComplex(data1)
		chk.VectorC(tst, "data1", 1e-10, X, Y)
	}
}
