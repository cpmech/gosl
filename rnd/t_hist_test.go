// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_hist01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("hist01")

	lims := []float64{0, 1, 2, 3, 4, 5}
	hist := Histogram{Stations: lims}

	idx := hist.FindBin(-3.3)
	chk.IntAssert(idx, -1)

	idx = hist.FindBin(7.0)
	chk.IntAssert(idx, -1)

	for i, x := range lims {
		idx = hist.FindBin(x)
		io.Pforan("x=%g idx=%d\n", x, idx)
		if i < len(lims)-1 {
			chk.IntAssert(idx, i)
		} else {
			chk.IntAssert(idx, -1)
		}
	}

	idx = hist.FindBin(0.5)
	chk.IntAssert(idx, 0)

	idx = hist.FindBin(1.5)
	chk.IntAssert(idx, 1)

	idx = hist.FindBin(2.5)
	chk.IntAssert(idx, 2)

	idx = hist.FindBin(3.99999999999999)
	chk.IntAssert(idx, 3)

	idx = hist.FindBin(4.999999)
	chk.IntAssert(idx, 4)

	hist.Count([]float64{
		0, 0.1, 0.2, 0.3, 0.9, // 5
		1, 1, 1, 1.2, 1.3, 1.4, 1.5, 1.99, // 8
		2, 2.5, // 2
		3, 3.5, // 2
		4.1, 4.5, 4.9, // 3
		-3, -2, -1,
		5, 6, 7, 8,
	}, true)
	io.Pforan("counts = %v\n", hist.Counts)
	chk.Ints(tst, "counts", hist.Counts, []int{5, 8, 2, 2, 3})

	labels := hist.GenLabels("%g")
	io.Pforan("labels = %v\n", labels)
}

func Test_hist02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("hist02")

	lims := []int{0, 1, 2, 3, 4, 5}
	hist := IntHistogram{Stations: lims}

	idx := hist.FindBin(-3)
	chk.IntAssert(idx, -1)

	idx = hist.FindBin(7)
	chk.IntAssert(idx, -1)

	for i, x := range lims {
		idx = hist.FindBin(x)
		io.Pforan("x=%v idx=%v\n", x, idx)
		if i < len(lims)-1 {
			chk.IntAssert(idx, i)
		} else {
			chk.IntAssert(idx, -1)
		}
	}

	hist.Count([]int{
		0, 0, 0, 0, 0, // 5
		1, 1, 1, 1, 1, 1, 1, 1, // 8
		2, 2, // 2
		3, 3, // 2
		4, 4, 4, // 3
		-3, -2, -1,
		5, 6, 7, 8,
	}, true)
	io.Pforan("counts = %v\n", hist.Counts)
	chk.Ints(tst, "counts", hist.Counts, []int{5, 8, 2, 2, 3})

	labels := hist.GenLabels("%d")
	io.Pforan("labels = %v\n", labels)

	if chk.Verbose {
		plt.Reset(true, nil)
		hist.Plot(true, nil, nil)
		plt.Save("/tmp/gosl/rnd", "hist02")
	}
}
