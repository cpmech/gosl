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
	"github.com/cpmech/gosl/utl"
)

// Hairer-Wanner VII-p2 Eq.(1.1)
func Test_ode01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode01: Hairer-Wanner VII-p2 Eq.(1.1)")

	lam := -50.0
	xa, xb := 0.0, 1.5
	ya := la.Vector([]float64{0.0})
	ndim := len(ya)
	y := la.NewVector(ndim)

	fcn := func(f la.Vector, dx, x float64, y la.Vector) error {
		f[0] = lam*y[0] - lam*math.Cos(x)
		return nil
	}

	jac := func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
		if dfdy.Max() == 0 {
			dfdy.Init(1, 1, 1)
		}
		dfdy.Start()
		dfdy.Put(0, 0, lam)
		return nil
	}

	// FwEuler
	io.Pforan(". . . FwEuler . . . \n")
	dx := 1.875 / 50.0
	copy(y, ya)
	sol1 := NewSolver(FwEulerKind, ndim, fcn, jac, nil, nil)
	sol1.SaveXY = true
	sol1.Solve(y, xa, xb, dx, true)
	chk.Int(tst, "number of F evaluations ", sol1.Nfeval, 40)
	chk.Int(tst, "number of J evaluations ", sol1.Njeval, 0)
	chk.Int(tst, "total number of steps   ", sol1.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", sol1.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", sol1.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol1.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", sol1.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", sol1.Nitmax, 0)
	chk.Int(tst, "IdxSave", sol1.IdxSave, sol1.Nsteps+1)

	// BwEuler
	io.Pforan(". . . BwEuler . . . \n")
	copy(y, ya)
	sol2 := NewSolver(BwEulerKind, ndim, fcn, jac, nil, nil)
	sol2.SaveXY = true
	sol2.Solve(y, xa, xb, dx, true)
	chk.Int(tst, "number of F evaluations ", sol2.Nfeval, 80)
	chk.Int(tst, "number of J evaluations ", sol2.Njeval, 40)
	chk.Int(tst, "total number of steps   ", sol2.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", sol2.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", sol2.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol2.Ndecomp, 40)
	chk.Int(tst, "number of lin solutions ", sol2.Nlinsol, 40)
	chk.Int(tst, "max number of iterations", sol2.Nitmax, 2)
	chk.Int(tst, "IdxSave", sol2.IdxSave, sol2.Nsteps+1)

	// MoEuler
	io.Pforan(". . . MoEuler . . . \n")
	copy(y, ya)
	sol3 := NewSolver(MoEulerKind, ndim, fcn, jac, nil, nil)
	sol3.SaveXY = true
	sol3.Solve(y, xa, xb, xb-xa, false)
	chk.Int(tst, "number of F evaluations ", sol3.Nfeval, 379)
	chk.Int(tst, "number of J evaluations ", sol3.Njeval, 0)
	chk.Int(tst, "total number of steps   ", sol3.Nsteps, 189)
	chk.Int(tst, "number of accepted steps", sol3.Naccepted, 189)
	chk.Int(tst, "number of rejected steps", sol3.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol3.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", sol3.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", sol3.Nitmax, 0)
	chk.Int(tst, "IdxSave", sol3.IdxSave, sol3.Naccepted+1)

	// DoPri5
	io.Pforan(". . . DoPri5 . . . \n")
	copy(y, ya)
	sol4 := NewSolver(DoPri5kind, ndim, fcn, jac, nil, nil)
	sol4.SaveXY = true
	sol4.Solve(y, xa, xb, xb-xa, false)
	chk.Int(tst, "number of F evaluations ", sol4.Nfeval, 1132)
	chk.Int(tst, "number of J evaluations ", sol4.Njeval, 0)
	chk.Int(tst, "total number of steps   ", sol4.Nsteps, 172)
	chk.Int(tst, "number of accepted steps", sol4.Naccepted, 99)
	chk.Int(tst, "number of rejected steps", sol4.Nrejected, 73)
	chk.Int(tst, "number of decompositions", sol4.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", sol4.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", sol4.Nitmax, 0)
	chk.Int(tst, "IdxSave", sol4.IdxSave, sol4.Naccepted+1)

	// Radau5
	io.Pforan(". . . Radau5 . . . \n")
	copy(y, ya)
	sol5 := NewSolver(Radau5kind, ndim, fcn, jac, nil, nil)
	sol5.SaveXY = true
	sol5.Solve(y, xa, xb, xb-xa, false)
	chk.Int(tst, "number of F evaluations ", sol5.Nfeval, 66)
	chk.Int(tst, "number of J evaluations ", sol5.Njeval, 1)
	chk.Int(tst, "total number of steps   ", sol5.Nsteps, 15)
	chk.Int(tst, "number of accepted steps", sol5.Naccepted, 15)
	chk.Int(tst, "number of rejected steps", sol5.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol5.Ndecomp, 13)
	chk.Int(tst, "number of lin solutions ", sol5.Nlinsol, 17)
	chk.Int(tst, "max number of iterations", sol5.Nitmax, 2)
	chk.Int(tst, "IdxSave", sol5.IdxSave, sol5.Naccepted+1)

	if chk.Verbose {
		X := utl.LinSpace(xa, xb, 101)
		Y := make([]float64, len(X))
		for i := 0; i < len(X); i++ {
			Y[i] = -lam * (math.Sin(X[i]) - lam*math.Cos(X[i]) + lam*math.Exp(lam*X[i])) / (lam*lam + 1.0)
		}
		a, b, c, d, e := sol1.IdxSave, sol2.IdxSave, sol3.IdxSave, sol4.IdxSave, sol5.IdxSave
		Ya, Yb, Yc, Yd, Ye := sol1.ExtractTimeSeries(0), sol2.ExtractTimeSeries(0), sol3.ExtractTimeSeries(0), sol4.ExtractTimeSeries(0), sol5.ExtractTimeSeries(0)
		plt.Reset(false, nil)
		plt.Plot(X, Y, &plt.A{C: "grey", Ls: "-", Lw: 10, L: "solution"})
		plt.Plot(sol1.Xvalues[:a], Ya, &plt.A{C: "k", M: ".", Ls: ":", L: "FwEuler"})
		plt.Plot(sol2.Xvalues[:b], Yb, &plt.A{C: "r", M: ".", Ls: ":", L: "BwEuler"})
		plt.Plot(sol3.Xvalues[:c], Yc, &plt.A{C: "c", M: "+", Ls: ":", L: "MoEuler"})
		plt.Plot(sol4.Xvalues[:d], Yd, &plt.A{C: "m", M: ".", Ls: "--", L: "Dopri5"})
		plt.Plot(sol5.Xvalues[:e], Ye, &plt.A{C: "b", M: "o", Ls: "-", L: "Radau5"})
		plt.Gll("$x$", "$y$", nil)
		plt.Save("/tmp/gosl/ode", "ode1")
	}
}

// Hairer-Wanner VII-p5 Eq.(1.5) Van der Pol's Equation
func Test_ode02(tst *testing.T) {

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

	// method and flags
	fixstp := false
	method := Radau5kind
	numjac := false
	xa, xb := 0.0, 2.0
	ya := la.Vector([]float64{2.0, -0.6})
	ndim := len(ya)

	// allocate ODE object
	var o *Solver
	if numjac {
		o = NewSolver(method, ndim, fcn, nil, nil, nil)
	} else {
		o = NewSolver(method, ndim, fcn, jac, nil, nil)
	}

	// tolerances and initial step size
	rtol := 1e-4
	atol := rtol
	o.IniH = 1.0e-4
	o.SetTol(atol, rtol)
	//o.NmaxSS = 2
	o.SaveXY = true

	// solve problem
	y := la.NewVector(ndim)
	copy(y, ya)
	t0 := time.Now()
	if fixstp {
		o.Solve(y, xa, xb, 0.05, fixstp)
	} else {
		o.Solve(y, xa, xb, xb-xa, fixstp)
	}
	chk.Int(tst, "number of F evaluations ", o.Nfeval, 2233)
	chk.Int(tst, "number of J evaluations ", o.Njeval, 160)
	chk.Int(tst, "total number of steps   ", o.Nsteps, 280)
	chk.Int(tst, "number of accepted steps", o.Naccepted, 241)
	chk.Int(tst, "number of rejected steps", o.Nrejected, 7)
	chk.Int(tst, "number of decompositions", o.Ndecomp, 251)
	chk.Int(tst, "number of lin solutions ", o.Nlinsol, 663)
	chk.Int(tst, "max number of iterations", o.Nitmax, 6)
	io.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5, FszXtck: 6, FszYtck: 6})
		_, T, err := io.ReadTable("data/vdpol_radau5_for.dat")
		if err != nil {
			chk.Panic("%v", err)
		}
		s := o.IdxSave
		for j := 0; j < ndim; j++ {
			labelA, labelB := "", ""
			if j == 2 {
				labelA, labelB = "reference", "gosl"
			}
			Yj := o.ExtractTimeSeries(j)
			plt.Subplot(ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(o.Xvalues[:s], Yj, &plt.A{C: "r", M: ".", Ms: 2, Ls: "none", L: labelB})
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(ndim+1, 1, ndim+1)
		plt.Plot(o.Xvalues[1:s], o.Hvalues[1:s], &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl/ode", "vdpolA")
	}
}

// Hairer-Wanner VII-p3 Eq.(1.4) Robertson's Equation
func Test_ode03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode03: Hairer-Wanner VII-p3 Eq.(1.4) Robertson's Equation")

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

	// data
	fixstp := false
	method := Radau5kind
	xa, xb := 0.0, 0.3
	ya := la.Vector([]float64{1.0, 0.0, 0.0})
	ndim := len(ya)

	// allocate ODE object
	o := NewSolver(method, ndim, fcn, jac, nil, nil)
	o.SaveXY = true

	// tolerances and initial step size
	rtol := 1e-2
	atol := rtol * 1e-6
	o.SetTol(atol, rtol)
	o.IniH = 1.0e-6

	// solve problem
	y := la.NewVector(ndim)
	copy(y, ya)
	if fixstp {
		o.Solve(y, xa, xb, 0.01, fixstp)
	} else {
		o.Solve(y, xa, xb, xb-xa, fixstp)
	}
	chk.Int(tst, "number of F evaluations ", o.Nfeval, 87)
	chk.Int(tst, "number of J evaluations ", o.Njeval, 8)
	chk.Int(tst, "total number of steps   ", o.Nsteps, 17)
	chk.Int(tst, "number of accepted steps", o.Naccepted, 15)
	chk.Int(tst, "number of rejected steps", o.Nrejected, 1)
	chk.Int(tst, "number of decompositions", o.Ndecomp, 15)
	chk.Int(tst, "number of lin solutions ", o.Nlinsol, 24)
	chk.Int(tst, "max number of iterations", o.Nitmax, 2)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5, FszXtck: 6, FszYtck: 6})
		_, T, err := io.ReadTable("data/rober_radau5_cpp.dat")
		if err != nil {
			chk.Panic("%v", err)
		}
		s := o.IdxSave
		for j := 0; j < ndim; j++ {
			labelA, labelB := "", ""
			if j == 2 {
				labelA, labelB = "reference", "gosl"
			}
			Yj := o.ExtractTimeSeries(j)
			plt.Subplot(ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(o.Xvalues[:s], Yj, &plt.A{C: "r", M: ".", Ms: 2, Ls: "none", L: labelB})
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(ndim+1, 1, ndim+1)
		plt.Plot(o.Xvalues[1:s], o.Hvalues[1:s], &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl/ode", "rober")
	}
}

func Test_ode04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode04: Hairer-Wanner VII-p376 Transistor Amplifier")
	// NOTE: from E Hairer's website, not the system in the book

	// data
	UE, UB, UF, ALPHA, BETA := 0.1, 6.0, 0.026, 0.99, 1.0e-6
	R0, R1, R2, R3, R4, R5 := 1000.0, 9000.0, 9000.0, 9000.0, 9000.0, 9000.0
	R6, R7, R8, R9 := 9000.0, 9000.0, 9000.0, 9000.0
	W := 2.0 * 3.141592654 * 100.0

	// initial values
	xa := 0.0
	ya := la.Vector([]float64{0.0,
		UB,
		UB / (R6/R5 + 1.0),
		UB / (R6/R5 + 1.0),
		UB,
		UB / (R2/R1 + 1.0),
		UB / (R2/R1 + 1.0),
		0.0,
	})

	// endpoint of integration
	xb := 0.05
	//xb = 0.0123 // OK
	//xb = 0.01235 // !OK

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
	var M la.Triplet
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

	// flags
	fixstp := false
	method := Radau5kind
	ndim := len(ya)
	numjac := false

	// ODE solver
	var o *Solver
	if numjac {
		o = NewSolver(method, ndim, fcn, nil, &M, nil)
	} else {
		o = NewSolver(method, ndim, fcn, jac, &M, nil)
	}
	o.IniH = 1.0e-6 // initial step size
	o.SaveXY = true

	// set tolerances
	atol, rtol := 1e-11, 1e-5
	o.SetTol(atol, rtol)

	// run
	t0 := time.Now()
	if fixstp {
		o.Solve(ya, xa, xb, 0.01, fixstp)
	} else {
		o.Solve(ya, xa, xb, xb-xa, fixstp)
	}
	if true {
		chk.Int(tst, "number of F evaluations ", o.Nfeval, 2599)
		chk.Int(tst, "number of J evaluations ", o.Njeval, 216)
		chk.Int(tst, "total number of steps   ", o.Nsteps, 275)
		chk.Int(tst, "number of accepted steps", o.Naccepted, 219)
		chk.Int(tst, "number of rejected steps", o.Nrejected, 20)
		chk.Int(tst, "number of decompositions", o.Ndecomp, 274)
		chk.Int(tst, "number of lin solutions ", o.Nlinsol, 792)
		chk.Int(tst, "max number of iterations", o.Nitmax, 6)
	}
	io.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 450, Dpi: 150, Prop: 1.8, FszXtck: 6, FszYtck: 6})
		_, T, err := io.ReadTable("data/radau5_hwamplifier.dat")
		if err != nil {
			chk.Panic("%v", err)
		}
		s := o.IdxSave
		for j := 0; j < ndim; j++ {
			labelA, labelB := "", ""
			if j == 4 {
				labelA, labelB = "reference", "gosl"
			}
			Yj := o.ExtractTimeSeries(j)
			plt.Subplot(ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(o.Xvalues[:s], Yj, &plt.A{C: "r", M: ".", Ms: 1, Ls: "none", L: labelB})
			plt.AxisXmax(0.05)
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(ndim+1, 1, ndim+1)
		plt.Plot(o.Xvalues[1:s], o.Hvalues[1:s], &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.AxisXmax(0.05)
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl/ode", "hwamplifier")
	}
}
