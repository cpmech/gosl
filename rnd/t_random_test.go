// Copyright 2016 The Gosl Authors. All rights reserved.
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

const Nsamples = 10

func Test_GOint01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GOint01. integers")

	Init(1234)

	nints := 10
	vals := make([]int, Nsamples)

	// using Int
	t0 := time.Now()
	for i := 0; i < Nsamples; i++ {
		vals[i] = Int(0, nints-1)
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist := IntHistogram{Stations: utl.IntRange(nints + 1)}
	hist.Count(vals, true)
	io.Pfyel(TextHist(hist.GenLabels("%d"), hist.Counts, 60))

	// using Ints
	t0 = time.Now()
	Ints(vals, 0, nints-1)
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist.Count(vals, true)
	io.Pfcyan(TextHist(hist.GenLabels("%d"), hist.Counts, 60))
}

func Test_GOflt01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GOflt01. float64")

	Init(1234)

	xmin := 10.0
	xmax := 20.0
	vals := make([]float64, Nsamples)

	// using Float64
	t0 := time.Now()
	for i := 0; i < Nsamples; i++ {
		vals[i] = Float64(xmin, xmax)
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist := Histogram{Stations: []float64{10, 12.5, 15, 17.5, 20}}
	hist.Count(vals, true)
	io.Pfpink(TextHist(hist.GenLabels("%4g"), hist.Counts, 60))

	// using Float64s
	t0 = time.Now()
	Float64s(vals, xmin, xmax)
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	hist.Count(vals, true)
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

func Test_getunique01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("getunique01")

	Init(0)

	nsel := 5 // number of selections
	size := 10
	nums := utl.IntRange(size)
	hist := IntHistogram{Stations: utl.IntRange(size + 5)}
	sel := IntGetUnique(nums, nsel)
	io.Pfgreen("nums = %v\n", nums)
	io.Pfcyan("sel  = %v\n", sel)
	for i := 0; i < Nsamples; i++ {
		sel := IntGetUnique(nums, nsel)
		checkRepeated(sel)
		hist.Count(sel, false)
		//io.Pfgrey("sel  = %v\n", sel)
	}

	io.Pf(TextHist(hist.GenLabels("%d"), hist.Counts, 60))
}

func Test_getunique02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("getunique02")

	Init(0)

	nsel := 5 // number of selections
	start := 2
	endp1 := 10
	hist := IntHistogram{Stations: utl.IntRange(endp1 + 3)}
	sel := IntGetUniqueN(start, endp1, nsel)
	io.Pfcyan("sel  = %v\n", sel)
	for i := 0; i < Nsamples; i++ {
		sel := IntGetUniqueN(start, endp1, nsel)
		checkRepeated(sel)
		hist.Count(sel, false)
		//io.Pfgrey("sel  = %v\n", sel)
	}

	io.Pf(TextHist(hist.GenLabels("%d"), hist.Counts, 60))
}

func Test_groups01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("groups01")

	Init(0)

	ng := 3            // number of groups
	nints := 12        // number of integers
	size := nints / ng // groups size
	ints := utl.IntRange(nints)
	groups := utl.IntAlloc(ng, size)
	hists := make([]*IntHistogram, ng)
	for i := 0; i < ng; i++ {
		hists[i] = &IntHistogram{Stations: utl.IntRange(nints + 1)}
	}
	IntGetGroups(groups, ints)
	io.Pfcyan("groups = %v\n", groups)
	for i := 0; i < Nsamples; i++ {
		IntGetGroups(groups, ints)
		for j := 0; j < ng; j++ {
			checkRepeated(groups[j])
			hists[j].Count(groups[j], false)
		}
	}
	for i := 0; i < ng; i++ {
		io.Pf("\n")
		io.Pf(TextHist(hists[i].GenLabels("%d"), hists[i].Counts, 60))
	}
}

func checkRepeated(v []int) {
	for i := 1; i < len(v); i++ {
		if v[i] == v[i-1] {
			chk.Panic("there are repeated entries in v = %v", v)
		}
	}
}
