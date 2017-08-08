// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

// Hairer-Wanner VII-p2 Eq.(1.1)
func TestOde01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode01: Hairer-Wanner VII-p2 Eq.(1.1)")

	// problem
	p := ProbHwEq11()

	// FwEuler
	io.Pforan("\n. . . FwEuler . . . \n")
	stat1, out1, err := p.Solve(FwEulerKind, true, false)
	status(tst, err)
	chk.Int(tst, "number of F evaluations ", stat1.Nfeval, 40)
	chk.Int(tst, "number of J evaluations ", stat1.Njeval, 0)
	chk.Int(tst, "total number of steps   ", stat1.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", stat1.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", stat1.Nrejected, 0)
	chk.Int(tst, "number of decompositions", stat1.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", stat1.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", stat1.Nitmax, 0)

	// BwEuler
	io.Pforan("\n. . . BwEuler . . . \n")
	stat2, out2, err := p.Solve(BwEulerKind, true, false)
	status(tst, err)
	chk.Int(tst, "number of F evaluations ", stat2.Nfeval, 80)
	chk.Int(tst, "number of J evaluations ", stat2.Njeval, 40)
	chk.Int(tst, "total number of steps   ", stat2.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", stat2.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", stat2.Nrejected, 0)
	chk.Int(tst, "number of decompositions", stat2.Ndecomp, 40)
	chk.Int(tst, "number of lin solutions ", stat2.Nlinsol, 40)
	chk.Int(tst, "max number of iterations", stat2.Nitmax, 2)

	// MoEuler
	io.Pforan("\n. . . MoEuler . . . \n")
	stat3, out3, err := p.Solve(MoEulerKind, false, false)
	status(tst, err)
	chk.Int(tst, "number of F evaluations ", stat3.Nfeval, 379)
	chk.Int(tst, "number of J evaluations ", stat3.Njeval, 0)
	chk.Int(tst, "total number of steps   ", stat3.Nsteps, 189)
	chk.Int(tst, "number of accepted steps", stat3.Naccepted, 189)
	chk.Int(tst, "number of rejected steps", stat3.Nrejected, 0)
	chk.Int(tst, "number of decompositions", stat3.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", stat3.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", stat3.Nitmax, 0)

	// DoPri5
	io.Pforan("\n. . . DoPri5 . . . \n")
	stat4, out4, err := p.Solve(DoPri5kind, false, false)
	status(tst, err)
	chk.Int(tst, "number of F evaluations ", stat4.Nfeval, 1132)
	chk.Int(tst, "number of J evaluations ", stat4.Njeval, 0)
	chk.Int(tst, "total number of steps   ", stat4.Nsteps, 172)
	chk.Int(tst, "number of accepted steps", stat4.Naccepted, 99)
	chk.Int(tst, "number of rejected steps", stat4.Nrejected, 73)
	chk.Int(tst, "number of decompositions", stat4.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", stat4.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", stat4.Nitmax, 0)

	// Radau5
	io.Pforan("\n. . . Radau5 . . . \n")
	stat5, out5, err := p.Solve(Radau5kind, false, false)
	status(tst, err)
	chk.Int(tst, "number of F evaluations ", stat5.Nfeval, 66)
	chk.Int(tst, "number of J evaluations ", stat5.Njeval, 1)
	chk.Int(tst, "total number of steps   ", stat5.Nsteps, 15)
	chk.Int(tst, "number of accepted steps", stat5.Naccepted, 15)
	chk.Int(tst, "number of rejected steps", stat5.Nrejected, 0)
	chk.Int(tst, "number of decompositions", stat5.Ndecomp, 13)
	chk.Int(tst, "number of lin solutions ", stat5.Nlinsol, 17)
	chk.Int(tst, "max number of iterations", stat5.Nitmax, 2)

	// plot
	if chk.Verbose {
		npts := 201
		plt.Reset(true, nil)
		p.Plot("FwEuler", out1, npts, true, nil, &plt.A{C: "k", M: ".", Ls: ":"})
		p.Plot("BwEuler", out2, npts, false, nil, &plt.A{C: "r", M: ".", Ls: ":"})
		p.Plot("MoEuler", out3, npts, false, nil, &plt.A{C: "c", M: "+", Ls: ":"})
		p.Plot("Dopri5", out4, npts, false, nil, &plt.A{C: "m", M: "^", Ls: "--", Ms: 3})
		p.Plot("Radau5", out5, npts, false, nil, &plt.A{C: "b", M: "s", Ls: "-", Ms: 3})
		plt.Gll("$x$", "$y$", nil)
		plt.Save("/tmp/gosl/ode", "ode1")
	}
}

// Hairer-Wanner VII-p5 Eq.(1.5) Van der Pol's Equation
func TestOde02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode02: Hairer-Wanner VII-p5 Eq.(1.5) Van der Pol's Equation")

	// problem definition
	eps := 1.0e-6
	fcn := func(f la.Vector, dx, x float64, y la.Vector) error {
		f[0] = y[1]
		f[1] = ((1.0-y[0]*y[0])*y[1] - y[0]) / eps
		return nil
	}
	jac := func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
		if dfdy.Max() == 0 {
			dfdy.Init(2, 2, 4)
		}
		dfdy.Start()
		dfdy.Put(0, 0, 0.0)
		dfdy.Put(0, 1, 1.0)
		dfdy.Put(1, 0, (-2.0*y[0]*y[1]-1.0)/eps)
		dfdy.Put(1, 1, (1.0-y[0]*y[0])/eps)
		return nil
	}
	xf := 2.0
	y := la.Vector([]float64{2.0, -0.6})
	ndim := len(y)

	// configuration
	conf, err := NewConfig(Radau5kind, "", nil, nil)
	status(tst, err)
	conf.SaveXY = true

	// allocate ODE object
	sol, err := NewSolver(conf, ndim, fcn, jac, nil, nil)
	status(tst, err)
	defer sol.Free()

	// solve problem
	err = sol.Solve(y, 0, xf)
	status(tst, err)

	// check
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 2233)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 160)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 280)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 241)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 7)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 251)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 663)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 6)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5, FszXtck: 6, FszYtck: 6})
		_, T, err := io.ReadTable("data/vdpol_radau5_for.dat")
		if err != nil {
			chk.Panic("%v", err)
		}
		s := sol.Out.IdxSave
		X := sol.Out.Xvalues[:s]
		for j := 0; j < ndim; j++ {
			labelA, labelB := "", ""
			if j == 2 {
				labelA, labelB = "reference", "gosl"
			}
			Yj := sol.Out.ExtractTimeSeries(j)
			plt.Subplot(ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(X, Yj, &plt.A{C: "r", M: ".", Ms: 2, Ls: "none", L: labelB})
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(ndim+1, 1, ndim+1)
		plt.Plot(X, sol.Out.Hvalues[:s], &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl/ode", "vdpolA")
	}
}

// Hairer-Wanner VII-p3 Eq.(1.4) Robertson's Equation
func TestOde03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode03: Hairer-Wanner VII-p3 Eq.(1.4) Robertson's Equation")

	// problem definition
	fcn := func(f la.Vector, dx, x float64, y la.Vector) error {
		f[0] = -0.04*y[0] + 1.0e4*y[1]*y[2]
		f[1] = 0.04*y[0] - 1.0e4*y[1]*y[2] - 3.0e7*y[1]*y[1]
		f[2] = 3.0e7 * y[1] * y[1]
		return nil
	}
	jac := func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
		if dfdy.Max() == 0 {
			dfdy.Init(3, 3, 9)
		}
		dfdy.Start()
		dfdy.Put(0, 0, -0.04)
		dfdy.Put(0, 1, 1.0e4*y[2])
		dfdy.Put(0, 2, 1.0e4*y[1])
		dfdy.Put(1, 0, 0.04)
		dfdy.Put(1, 1, -1.0e4*y[2]-6.0e7*y[1])
		dfdy.Put(1, 2, -1.0e4*y[1])
		dfdy.Put(2, 0, 0.0)
		dfdy.Put(2, 1, 6.0e7*y[1])
		dfdy.Put(2, 2, 0.0)
		return nil
	}
	xf := 0.3
	y := la.Vector([]float64{1.0, 0.0, 0.0})
	ndim := len(y)

	// configuration
	conf, err := NewConfig(Radau5kind, "", nil, nil)
	status(tst, err)
	conf.SaveXY = true

	// tolerances and initial step size
	rtol := 1e-2
	atol := rtol * 1e-6
	conf.SetTol(atol, rtol)
	conf.IniH = 1.0e-6

	// allocate ODE object
	sol, err := NewSolver(conf, ndim, fcn, jac, nil, nil)
	status(tst, err)
	defer sol.Free()

	// solve problem
	err = sol.Solve(y, 0.0, xf)
	status(tst, err)

	// check
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 87)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 8)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 17)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 15)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 1)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 15)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 24)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 2)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5, FszXtck: 6, FszYtck: 6})
		_, T, err := io.ReadTable("data/rober_radau5_cpp.dat")
		if err != nil {
			chk.Panic("%v", err)
		}
		s := sol.Out.IdxSave
		X := sol.Out.Xvalues[:s]
		for j := 0; j < ndim; j++ {
			labelA, labelB := "", ""
			if j == 2 {
				labelA, labelB = "reference", "gosl"
			}
			Yj := sol.Out.ExtractTimeSeries(j)
			plt.Subplot(ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(X, Yj, &plt.A{C: "r", M: ".", Ms: 2, Ls: "none", L: labelB})
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(ndim+1, 1, ndim+1)
		plt.Plot(X, sol.Out.Hvalues[:s], &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl/ode", "rober")
	}
}

func TestOde04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode04: Hairer-Wanner VII-p376 Transistor Amplifier")
	// NOTE: from E Hairer's website, not the system in the book

	// data
	UE, UB, UF, ALPHA, BETA := 0.1, 6.0, 0.026, 0.99, 1.0e-6
	R0, R1, R2, R3, R4, R5 := 1000.0, 9000.0, 9000.0, 9000.0, 9000.0, 9000.0
	R6, R7, R8, R9 := 9000.0, 9000.0, 9000.0, 9000.0
	W := 2.0 * 3.141592654 * 100.0

	// initial values
	y := la.Vector([]float64{0.0,
		UB,
		UB / (R6/R5 + 1.0),
		UB / (R6/R5 + 1.0),
		UB,
		UB / (R2/R1 + 1.0),
		UB / (R2/R1 + 1.0),
		0.0,
	})

	// endpoint of integration
	xf := 0.05
	//xf = 0.0123 // OK
	//xf = 0.01235 // !OK

	// right-hand side of the amplifier problem
	fcn := func(f la.Vector, dx, x float64, y la.Vector) error {
		UET := UE * math.Sin(W*x)
		FAC1 := BETA * (math.Exp((y[3]-y[2])/UF) - 1.0)
		FAC2 := BETA * (math.Exp((y[6]-y[5])/UF) - 1.0)
		f[0] = y[0] / R9
		f[1] = (y[1]-UB)/R8 + ALPHA*FAC1
		f[2] = y[2]/R7 - FAC1
		f[3] = y[3]/R5 + (y[3]-UB)/R6 + (1.0-ALPHA)*FAC1
		f[4] = (y[4]-UB)/R4 + ALPHA*FAC2
		f[5] = y[5]/R3 - FAC2
		f[6] = y[6]/R1 + (y[6]-UB)/R2 + (1.0-ALPHA)*FAC2
		f[7] = (y[7] - UET) / R0
		return nil
	}

	// Jacobian of the amplifier problem
	jac := func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
		FAC14 := BETA * math.Exp((y[3]-y[2])/UF) / UF
		FAC27 := BETA * math.Exp((y[6]-y[5])/UF) / UF
		if dfdy.Max() == 0 {
			dfdy.Init(8, 8, 16)
		}
		NU := 2
		dfdy.Start()
		dfdy.Put(2+0-NU, 0, 1.0/R9)
		dfdy.Put(2+1-NU, 1, 1.0/R8)
		dfdy.Put(1+2-NU, 2, -ALPHA*FAC14)
		dfdy.Put(0+3-NU, 3, ALPHA*FAC14)
		dfdy.Put(2+2-NU, 2, 1.0/R7+FAC14)
		dfdy.Put(1+3-NU, 3, -FAC14)
		dfdy.Put(2+3-NU, 3, 1.0/R5+1.0/R6+(1.0-ALPHA)*FAC14)
		dfdy.Put(3+2-NU, 2, -(1.0-ALPHA)*FAC14)
		dfdy.Put(2+4-NU, 4, 1.0/R4)
		dfdy.Put(1+5-NU, 5, -ALPHA*FAC27)
		dfdy.Put(0+6-NU, 6, ALPHA*FAC27)
		dfdy.Put(2+5-NU, 5, 1.0/R3+FAC27)
		dfdy.Put(1+6-NU, 6, -FAC27)
		dfdy.Put(2+6-NU, 6, 1.0/R1+1.0/R2+(1.0-ALPHA)*FAC27)
		dfdy.Put(3+5-NU, 5, -(1.0-ALPHA)*FAC27)
		dfdy.Put(2+7-NU, 7, 1.0/R0)
		return nil
	}

	// matrix "M"
	c1, c2, c3, c4, c5 := 1.0e-6, 2.0e-6, 3.0e-6, 4.0e-6, 5.0e-6
	M := new(la.Triplet)
	M.Init(8, 8, 14)
	M.Start()
	NU := 1
	M.Put(1+0-NU, 0, -c5)
	M.Put(0+1-NU, 1, c5)
	M.Put(2+0-NU, 0, c5)
	M.Put(1+1-NU, 1, -c5)
	M.Put(1+2-NU, 2, -c4)
	M.Put(1+3-NU, 3, -c3)
	M.Put(0+4-NU, 4, c3)
	M.Put(2+3-NU, 3, c3)
	M.Put(1+4-NU, 4, -c3)
	M.Put(1+5-NU, 5, -c2)
	M.Put(1+6-NU, 6, -c1)
	M.Put(0+7-NU, 7, c1)
	M.Put(2+6-NU, 6, c1)
	M.Put(1+7-NU, 7, -c1)

	// configurations
	conf, err := NewConfig(Radau5kind, "", nil, nil)
	status(tst, err)
	conf.SaveXY = true
	conf.IniH = 1.0e-6 // initial step size

	// set tolerances
	atol, rtol := 1e-11, 1e-5
	conf.SetTol(atol, rtol)

	// ODE solver
	ndim := len(y)
	sol, err := NewSolver(conf, ndim, fcn, jac, M, nil)
	status(tst, err)
	defer sol.Free()

	// run
	t0 := time.Now()
	err = sol.Solve(y, 0.0, xf)
	status(tst, err)
	io.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))

	// check
	if true {
		chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 2599)
		chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 216)
		chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 275)
		chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 219)
		chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 20)
		chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 274)
		chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 792)
		chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 6)
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 450, Dpi: 150, Prop: 1.8, FszXtck: 6, FszYtck: 6})
		_, T, err := io.ReadTable("data/radau5_hwamplifier.dat")
		if err != nil {
			chk.Panic("%v", err)
		}
		s := sol.Out.IdxSave
		X := sol.Out.Xvalues[:s]
		for j := 0; j < ndim; j++ {
			labelA, labelB := "", ""
			if j == 4 {
				labelA, labelB = "reference", "gosl"
			}
			Yj := sol.Out.ExtractTimeSeries(j)
			plt.Subplot(ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(X, Yj, &plt.A{C: "r", M: ".", Ms: 1, Ls: "none", L: labelB})
			plt.AxisXmax(0.05)
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(ndim+1, 1, ndim+1)
		plt.Plot(X, sol.Out.Hvalues[:s], &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.AxisXmax(0.05)
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl/ode", "hwamplifier")
	}
}
