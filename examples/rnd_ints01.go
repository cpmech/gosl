// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

const NSAMPLES = 100000

func main() {

	// initialise seed with fixed number; use 0 to use current time
	rnd.Init(1234)

	// allocate slice for integers
	nints := 10
	vals := make([]int, NSAMPLES)

	// using the rnd.Int function
	t0 := time.Now()
	for i := 0; i < NSAMPLES; i++ {
		vals[i] = rnd.Int(0, nints-1)
	}
	io.Pf("time elapsed = %v\n", time.Now().Sub(t0))

	// text histogram
	hist := rnd.IntHistogram{Stations: utl.IntRange(nints + 1)}
	hist.Count(vals, true)
	io.Pf(rnd.TextHist(hist.GenLabels("%d"), hist.Counts, 60))

	// using the rnd.Ints function
	t0 = time.Now()
	rnd.Ints(vals, 0, nints-1)
	io.Pf("time elapsed = %v\n", time.Now().Sub(t0))

	// text histogram
	hist.Count(vals, true)
	io.Pf(rnd.TextHist(hist.GenLabels("%d"), hist.Counts, 60))
}
