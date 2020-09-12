// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package rnd

import (
	"sort"
	"testing"
	"time"

	"gosl/chk"
	"gosl/io"
	"gosl/utl"
)

func Test_MTint01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MTint01. integers (Mersenne Twister)")

	MTinit(1234)

	nints := 10
	vals := make([]int, Nsamples)

	// using MTint
	t0 := time.Now()
	for i := 0; i < Nsamples; i++ {
		vals[i] = MTint(0, nints-1)
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist := IntHistogram{Stations: utl.IntRange(nints + 1)}
	hist.Count(vals, true)
	io.Pfyel(TextHist(hist.GenLabels("%d"), hist.Counts, 60))

	// using MTints
	t0 = time.Now()
	MTints(vals, 0, nints-1)
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist.Count(vals, true)
	io.Pfcyan(TextHist(hist.GenLabels("%d"), hist.Counts, 60))
}

func Test_MTflt01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MTflt01. float64 (Mersenne Twister)")

	MTinit(1234)

	xmin := 10.0
	xmax := 20.0
	vals := make([]float64, Nsamples)

	// using MTfloat64
	t0 := time.Now()
	for i := 0; i < Nsamples; i++ {
		vals[i] = MTfloat64(xmin, xmax)
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist := Histogram{Stations: []float64{10, 12.5, 15, 17.5, 20}}
	hist.Count(vals, true)
	io.Pfpink(TextHist(hist.GenLabels("%4g"), hist.Counts, 60))

	// using MTfloat64s
	t0 = time.Now()
	MTfloat64s(vals, xmin, xmax)
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist.Count(vals, true)
	io.Pfblue2(TextHist(hist.GenLabels("%4g"), hist.Counts, 60))
}

func Test_MTshuffleInts01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MTshuffleInts01. Mersenne Twister")

	MTinit(0)

	n := 10
	nums := utl.IntRange(n)
	io.Pfgreen("before = %v\n", nums)
	MTintShuffle(nums)
	io.Pfcyan("after  = %v\n", nums)

	sort.Ints(nums)
	io.Pforan("sorted = %v\n", nums)
	chk.Ints(tst, "nums", nums, utl.IntRange(n))
}
