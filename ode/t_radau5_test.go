// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"encoding/json"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func TestRadau501a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Radau501a. Eq11 (analytical Jacobian)")

	// problem
	p := ProbHwEq11()

	// configuration
	conf := NewConfig("radau5", "", nil)
	conf.SetStepOut(true, nil)

	// solver
	sol := NewSolver(p.Ndim, conf, p.Fcn, p.Jac, nil)
	defer sol.Free()

	// solve ODE
	sol.Solve(p.Y, 0.0, p.Xf)

	// check Stat
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 66)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 1)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 15)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 15)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 13)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 17)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 2)

	// check results
	chk.Float64(tst, "yFin", 2.889e-5, p.Y[0], p.CalcYana(0, p.Xf))
}

func TestRadau502(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Radau502: Van der Pol's Equation")

	// problem
	p := ProbVanDerPol(0, false)
	p.Y[1] = -0.66 // for some reason, the previous reference code was using -0.6
	///////////////// -0.66 is the value from Hairer's website code
	//p.Xf = 0.6

	// configuration
	conf := NewConfig("radau5", "", nil)
	conf.SetStepOut(true, nil)
	conf.IniH = 1e-6
	conf.SetTols(1e-4, 1e-4)

	// dense output function
	ss := make([]int, 11)
	xx := make([]float64, 11)
	yy0 := make([]float64, 11)
	yy1 := make([]float64, 11)
	iout := 0
	io.Pf("\n%5s%7s%23s%23s\n", "s", "x", "y0", "y1")
	conf.SetDenseOut(true, 0.2, p.Xf, func(istep int, h, x float64, y la.Vector, xout float64, yout la.Vector) (stop bool) {
		io.Pf("%5d%7.3f%23.15e%23.15e\n", istep, x, y[0], y[1])
		ss[iout] = istep
		xx[iout] = xout
		yy0[iout] = yout[0]
		yy1[iout] = yout[1]
		iout++
		return
	})

	// allocate ODE object
	sol := NewSolver(p.Ndim, conf, p.Fcn, p.Jac, nil)
	defer sol.Free()

	// solve problem
	sol.Solve(p.Y, 0, p.Xf)

	// check
	io.Pl()
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 2218)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 161)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 275)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 238)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 8)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 248)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 660)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 6)

	// compare with fortran code
	d := readRefData("data/dr1_radau5.cmp")
	chk.Ints(tst, "S", sol.Out.GetDenseS(), d.S)
	chk.Array(tst, "X", 1e-15, sol.Out.GetDenseX(), d.X)
	chk.Deep2(tst, "Y", 1e-11, sol.Out.GetDenseYtableT(), d.Y)

	// check saved output
	chk.Ints(tst, "ss", ss, d.S)
	chk.Array(tst, "xx", 1e-15, xx, d.X)
	chk.Array(tst, "yy0", 1e-11, yy0, d.Y[0])
	chk.Array(tst, "yy1", 1e-11, yy1, d.Y[1])
}

type refData struct {
	S []int       // [nout]
	X []float64   // [nout]
	Y [][]float64 // [dim][nout]
}

func readRefData(fn string) (o *refData) {
	o = new(refData)
	b := io.ReadFile(fn)
	err := json.Unmarshal(b, o)
	if err != nil {
		chk.Panic("%v\n", err)
	}
	return
}
