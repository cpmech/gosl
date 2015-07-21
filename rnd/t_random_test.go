// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"sort"
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

const NSAMPLES = 1000

func Test_GOint01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GOint01. integers")

	Init(1234)

	nints := 10
	vals := make([]int, NSAMPLES)

	// using Int
	t0 := time.Now()
	for i := 0; i < NSAMPLES; i++ {
		vals[i] = Int(0, nints-1)
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist := IntHistogram{Stations: utl.IntRange(nints + 1)}
	hist.Count(vals)
	io.Pfyel(TextHist(hist.GenLabels("%d"), hist.Counts, 60))

	// using Ints
	t0 = time.Now()
	Ints(vals, 0, nints-1)
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist.Count(vals)
	io.Pfcyan(TextHist(hist.GenLabels("%d"), hist.Counts, 60))
}

func Test_MTint01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MTint01. integers (Mersenne Twister)")

	Init(1234)

	nints := 10
	vals := make([]int, NSAMPLES)

	// using MTint
	t0 := time.Now()
	for i := 0; i < NSAMPLES; i++ {
		vals[i] = MTint(0, nints-1)
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist := IntHistogram{Stations: utl.IntRange(nints + 1)}
	hist.Count(vals)
	io.Pfyel(TextHist(hist.GenLabels("%d"), hist.Counts, 60))

	// using MTints
	t0 = time.Now()
	MTints(vals, 0, nints-1)
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist.Count(vals)
	io.Pfcyan(TextHist(hist.GenLabels("%d"), hist.Counts, 60))
}

func Test_GOflt01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GOflt01. float64")

	Init(1234)

	xmin := 10.0
	xmax := 20.0
	vals := make([]float64, NSAMPLES)

	// using Float64
	t0 := time.Now()
	for i := 0; i < NSAMPLES; i++ {
		vals[i] = Float64(xmin, xmax)
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist := Histogram{Stations: []float64{10, 12.5, 15, 17.5, 20}}
	hist.Count(vals)
	io.Pfpink(TextHist(hist.GenLabels("%4g"), hist.Counts, 60))

	// using Float64s
	t0 = time.Now()
	Float64s(vals, xmin, xmax)
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist.Count(vals)
	io.Pfblue2(TextHist(hist.GenLabels("%4g"), hist.Counts, 60))
}

func Test_MTflt01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MTflt01. float64 (Mersenne Twister)")

	Init(1234)

	xmin := 10.0
	xmax := 20.0
	vals := make([]float64, NSAMPLES)

	// using MTfloat64
	t0 := time.Now()
	for i := 0; i < NSAMPLES; i++ {
		vals[i] = MTfloat64(xmin, xmax)
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist := Histogram{Stations: []float64{10, 12.5, 15, 17.5, 20}}
	hist.Count(vals)
	io.Pfpink(TextHist(hist.GenLabels("%4g"), hist.Counts, 60))

	// using MTfloat64s
	t0 = time.Now()
	MTfloat64s(vals, xmin, xmax)
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist.Count(vals)
	io.Pfblue2(TextHist(hist.GenLabels("%4g"), hist.Counts, 60))
}

func Test_flip01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flip01")

	Init(1234)

	p := 0.5
	nsamples := 100
	ntrue := 0
	nfalse := 0
	for i := 0; i < nsamples; i++ {
		if FlipCoin(p) {
			ntrue++
		} else {
			nfalse++
		}
	}

	io.Pforan("ntrue  = %v (42)\n", ntrue)
	io.Pforan("nfalse = %v (58)\n", nfalse)
}

func Test_GOshuffleInts01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GOshuffleInts01")

	Init(0)

	n := 10
	nums := utl.IntRange(n)
	io.Pfgreen("before = %v\n", nums)
	IntShuffle(nums)
	io.Pfcyan("after  = %v\n", nums)

	sort.Ints(nums)
	io.Pforan("sorted = %v\n", nums)
	chk.Ints(tst, "nums", nums, utl.IntRange(n))

	shufled := IntGetShuffled(nums)
	io.Pfyel("shufled = %v\n", shufled)
	sort.Ints(shufled)
	chk.Ints(tst, "shufled", shufled, utl.IntRange(n))
}

func Test_MTshuffleInts01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MTshuffleInts01. Mersenne Twister")

	Init(0)

	n := 10
	nums := utl.IntRange(n)
	io.Pfgreen("before = %v\n", nums)
	MTintShuffle(nums)
	io.Pfcyan("after  = %v\n", nums)

	sort.Ints(nums)
	io.Pforan("sorted = %v\n", nums)
	chk.Ints(tst, "nums", nums, utl.IntRange(n))
}
