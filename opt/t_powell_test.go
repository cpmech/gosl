// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func runPowellTest(tst *testing.T, fnkey string, p *Problem, x0 la.Vector, tolf, tolx, α float64) (sol *Powell) {

	// solve using Gradient-Descent
	xmin := x0.GetCopy()
	sol = NewPowell(p)
	sol.UseHist = true
	fmin := sol.Min(xmin, nil)

	// check
	name := "Powell"
	io.Pforan("%s: NumIter = %v\n", name, sol.NumIter)
	io.Pf("%s: NumFeval = %v\n", name, sol.NumFeval)
	chk.Float64(tst, io.Sf("%s: fmin", name), tolf, fmin, p.Fref)
	chk.Array(tst, io.Sf("%s: xmin", name), tolx, xmin, p.Xref)
	io.Pl()
	return
}

func TestPowell01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Powell01. Very simple bi-dimensional optimization")

	// problem
	p := Factory.SimpleParaboloid()

	// initial point
	x0 := la.NewVectorSlice([]float64{1, 1})

	// run test
	α := 0.6
	runPowellTest(tst, "powell01", p, x0, 1e-10, 1e-5, α)
}

func TestPowell02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Powell02. quadratic optimization in 2D")

	// problem
	p := Factory.SimpleQuadratic2d()

	// initial point
	x0 := la.NewVectorSlice([]float64{1.5, -0.75})

	// run test
	α := 0.2
	runPowellTest(tst, "powell02", p, x0, 1e-13, 1e-6, α)
}

func TestPowell03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Powell02. quadratic optimization in 3D")

	// problem
	p := Factory.SimpleQuadratic3d()

	// initial point
	x0 := la.NewVectorSlice([]float64{1, 2, 3})

	// run test
	α := 0.1
	runPowellTest(tst, "powell03", p, x0, 1e-10, 1e-5, α)
}
