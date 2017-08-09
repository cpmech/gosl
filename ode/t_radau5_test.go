// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func TestRadau501a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Radau501a. Eq11 (analytical Jacobian)")

	// problem
	p := ProbHwEq11()

	// configuration
	conf, err := NewConfig("radau5", "", nil)
	status(tst, err)
	conf.SaveXY = true

	// solver
	sol, err := NewSolver(conf, p.Ndim, p.Fcn, p.Jac, nil, nil)
	status(tst, err)
	defer sol.Free()

	// solve ODE
	err = sol.Solve(p.Y, 0.0, p.Xf)
	status(tst, err)

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
	chk.Float64(tst, "yFin", 2.88898538383e-5, p.Y[0], p.Yana(p.Xf))

	// plot
	if chk.Verbose {
		plt.Reset(true, nil)
		p.Plot("Radau5,Jana", sol.Out, 101, true, nil, nil)
		plt.Save("/tmp/gosl/ode", "radau501a")
	}
}

func TestRadau502(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Radau502: Van der Pol's Equation")

	// problem
	p := ProbVanDerPol()
	p.Y[1] = -0.66 // for some reason the previous reference code was using -0.6
	//p.Xf = 0.6

	// configuration
	conf, err := NewConfig("radau5", "", nil)
	status(tst, err)
	conf.SaveXY = true
	conf.IniH = 1e-6
	conf.SetTol(1e-4, 1e-4)
	conf.ContDx = 0.2
	conf.ContNmax = int(math.Ceil(p.Xf/conf.ContDx)) + 1

	// output function
	out := func(istep int, h, x float64, y la.Vector) error {
		//io.Pf("x=%5.2f  y=%18.10e%18.10e  nstep=%d\n", x, y[0], y[1], istep)
		return nil
	}

	// allocate ODE object
	sol, err := NewSolver(conf, p.Ndim, p.Fcn, p.Jac, nil, out)
	status(tst, err)
	defer sol.Free()

	// solve problem
	err = sol.Solve(p.Y, 0, p.Xf)
	status(tst, err)

	// compare with fortran code
	_, d, err := io.ReadTable("data/dr1_radau5.cmp")
	status(tst, err)
	chk.Int(tst, "ContIdx", sol.Out.ContIdx, len(d["x"]))
	for i := 0; i < sol.Out.ContIdx; i++ {
		for j := 0; j < p.Ndim; j++ {
			key := io.Sf("y%d", j)
			chk.AnaNum(tst, key, 1e-11, sol.Out.ContY[i][j], d[key][i], chk.Verbose)
		}
	}
	for i := 0; i < sol.Out.ContIdx; i++ {
		chk.Int(tst, "stp", sol.Out.ContStp[i], int(d["stp"][i]))
	}

	// check
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 2218)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 161)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 275)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 238)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 8)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 248)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 660)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 6)
}
