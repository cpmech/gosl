// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func TestPowell01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Powell01. Simple bi-dimensional optimization")

	// function f({x})
	ffcn := func(x la.Vector) float64 {
		return x[0]*x[0] + x[1]*x[1] - 0.5
	}

	// initial point
	x0 := la.NewVectorSlice([]float64{1, 1})

	// solver
	solver := NewPowell(2, ffcn)
	solver.History = true
	fmin := solver.Min(x0, false)
	chk.Float64(tst, "fmin", 1e-15, fmin, -0.5)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.7})
		solver.Hist.Plot(0, 1, -1.1, 1.1, -1.1, 1.1, 41)
		plt.Save("/tmp/gosl/opt", "powell01")
	}
}
