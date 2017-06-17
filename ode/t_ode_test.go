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
	ya := []float64{0.0}
	ndim := len(ya)
	y := make([]float64, ndim)

	fcn := func(f []float64, dx, x float64, y []float64) error {
		f[0] = lam*y[0] - lam*math.Cos(x)
		return nil
	}

	jac := func(dfdy *la.Triplet, dx, x float64, y []float64) error {
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
	var FwEuler Solver
	FwEuler.Init("FwEuler", ndim, fcn, jac, nil, nil)
	FwEuler.SaveXY = true
	FwEuler.Solve(y, xa, xb, dx, true)
	chk.Int(tst, "number of F evaluations ", FwEuler.Nfeval, 40)
	chk.Int(tst, "number of J evaluations ", FwEuler.Njeval, 0)
	chk.Int(tst, "total number of steps   ", FwEuler.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", FwEuler.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", FwEuler.Nrejected, 0)
	chk.Int(tst, "number of decompositions", FwEuler.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", FwEuler.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", FwEuler.Nitmax, 0)
	chk.Int(tst, "IdxSave", FwEuler.IdxSave, FwEuler.Nsteps+1)

	// BwEuler
	io.Pforan(". . . BwEuler . . . \n")
	copy(y, ya)
	var BwEuler Solver
	//BwEuler.Init("BwEuler", ndim, fcn, nil, nil, out)
	BwEuler.Init("BwEuler", ndim, fcn, jac, nil, nil)
	BwEuler.SaveXY = true
	BwEuler.Solve(y, xa, xb, dx, true)
	chk.Int(tst, "number of F evaluations ", BwEuler.Nfeval, 80)
	chk.Int(tst, "number of J evaluations ", BwEuler.Njeval, 40)
	chk.Int(tst, "total number of steps   ", BwEuler.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", BwEuler.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", BwEuler.Nrejected, 0)
	chk.Int(tst, "number of decompositions", BwEuler.Ndecomp, 40)
	chk.Int(tst, "number of lin solutions ", BwEuler.Nlinsol, 40)
	chk.Int(tst, "max number of iterations", BwEuler.Nitmax, 2)
	chk.Int(tst, "IdxSave", BwEuler.IdxSave, BwEuler.Nsteps+1)

	// MoEuler
	io.Pforan(". . . MoEuler . . . \n")
	copy(y, ya)
	var MoEuler Solver
	MoEuler.Init("MoEuler", ndim, fcn, jac, nil, nil)
	MoEuler.SaveXY = true
	MoEuler.Solve(y, xa, xb, xb-xa, false)
	chk.Int(tst, "number of F evaluations ", MoEuler.Nfeval, 379)
	chk.Int(tst, "number of J evaluations ", MoEuler.Njeval, 0)
	chk.Int(tst, "total number of steps   ", MoEuler.Nsteps, 189)
	chk.Int(tst, "number of accepted steps", MoEuler.Naccepted, 189)
	chk.Int(tst, "number of rejected steps", MoEuler.Nrejected, 0)
	chk.Int(tst, "number of decompositions", MoEuler.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", MoEuler.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", MoEuler.Nitmax, 0)
	chk.Int(tst, "IdxSave", MoEuler.IdxSave, MoEuler.Naccepted+1)

	// Dopri5
	io.Pforan(". . . Dopri5 . . . \n")
	copy(y, ya)
	var Dopri5 Solver
	Dopri5.Init("Dopri5", ndim, fcn, jac, nil, nil)
	Dopri5.SaveXY = true
	Dopri5.Solve(y, xa, xb, xb-xa, false)
	chk.Int(tst, "number of F evaluations ", Dopri5.Nfeval, 1132)
	chk.Int(tst, "number of J evaluations ", Dopri5.Njeval, 0)
	chk.Int(tst, "total number of steps   ", Dopri5.Nsteps, 172)
	chk.Int(tst, "number of accepted steps", Dopri5.Naccepted, 99)
	chk.Int(tst, "number of rejected steps", Dopri5.Nrejected, 73)
	chk.Int(tst, "number of decompositions", Dopri5.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", Dopri5.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", Dopri5.Nitmax, 0)
	chk.Int(tst, "IdxSave", Dopri5.IdxSave, Dopri5.Naccepted+1)

	// Radau5
	io.Pforan(". . . Radau5 . . . \n")
	copy(y, ya)
	var Radau5 Solver
	Radau5.Init("Radau5", ndim, fcn, jac, nil, nil)
	Radau5.SaveXY = true
	Radau5.Solve(y, xa, xb, xb-xa, false)
	chk.Int(tst, "number of F evaluations ", Radau5.Nfeval, 66)
	chk.Int(tst, "number of J evaluations ", Radau5.Njeval, 1)
	chk.Int(tst, "total number of steps   ", Radau5.Nsteps, 15)
	chk.Int(tst, "number of accepted steps", Radau5.Naccepted, 15)
	chk.Int(tst, "number of rejected steps", Radau5.Nrejected, 0)
	chk.Int(tst, "number of decompositions", Radau5.Ndecomp, 13)
	chk.Int(tst, "number of lin solutions ", Radau5.Nlinsol, 17)
	chk.Int(tst, "max number of iterations", Radau5.Nitmax, 2)
	chk.Int(tst, "IdxSave", Radau5.IdxSave, Radau5.Naccepted+1)

	if chk.Verbose {
		X := utl.LinSpace(xa, xb, 101)
		Y := make([]float64, len(X))
		for i := 0; i < len(X); i++ {
			Y[i] = -lam * (math.Sin(X[i]) - lam*math.Cos(X[i]) + lam*math.Exp(lam*X[i])) / (lam*lam + 1.0)
		}
		a, b, c, d, e := FwEuler.IdxSave, BwEuler.IdxSave, MoEuler.IdxSave, Dopri5.IdxSave, Radau5.IdxSave
		plt.Reset(false, nil)
		plt.Plot(X, Y, &plt.A{C: "grey", Ls: "-", Lw: 10, L: "solution"})
		plt.Plot(FwEuler.Xvalues[:a], FwEuler.Yvalues[0][:a], &plt.A{C: "k", M: ".", Ls: ":", L: "FwEuler"})
		plt.Plot(BwEuler.Xvalues[:b], BwEuler.Yvalues[0][:b], &plt.A{C: "r", M: ".", Ls: ":", L: "BwEuler"})
		plt.Plot(MoEuler.Xvalues[:c], MoEuler.Yvalues[0][:c], &plt.A{C: "c", M: "+", Ls: ":", L: "MoEuler"})
		plt.Plot(Dopri5.Xvalues[:d], Dopri5.Yvalues[0][:d], &plt.A{C: "m", M: ".", Ls: "--", L: "Dopri5"})
		plt.Plot(Radau5.Xvalues[:e], Radau5.Yvalues[0][:e], &plt.A{C: "b", M: "o", Ls: "-", L: "Radau5"})
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
	fcn := func(f []float64, dx, x float64, y []float64) error {
		f[0] = y[1]
		f[1] = ((1.0-y[0]*y[0])*y[1] - y[0]) / eps
		return nil
	}
	jac := func(dfdy *la.Triplet, dx, x float64, y []float64) error {
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
	//method := "Dopri5"
	method := "Radau5"
	numjac := false
	xa, xb := 0.0, 2.0
	ya := []float64{2.0, -0.6}
	ndim := len(ya)

	// allocate ODE object
	var o Solver
	if numjac {
		o.Init(method, ndim, fcn, nil, nil, nil)
	} else {
		o.Init(method, ndim, fcn, jac, nil, nil)
	}

	// tolerances and initial step size
	rtol := 1e-4
	atol := rtol
	o.IniH = 1.0e-4
	o.SetTol(atol, rtol)
	//o.NmaxSS = 2
	o.SaveXY = true

	// solve problem
	y := make([]float64, ndim)
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
			plt.Subplot(ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(o.Xvalues[:s], o.Yvalues[j][:s], &plt.A{C: "r", M: ".", Ms: 2, Ls: "none", L: labelB})
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(ndim+1, 1, ndim+1)
		plt.Plot(o.Xvalues[1:s], o.Hvalues[1:s], &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl", "vdpolA")
	}
}

// Hairer-Wanner VII-p3 Eq.(1.4) Robertson's Equation
func Test_ode03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode03: Hairer-Wanner VII-p3 Eq.(1.4) Robertson's Equation")

	fcn := func(f []float64, dx, x float64, y []float64) error {
		f[0] = -0.04*y[0] + 1.0e4*y[1]*y[2]
		f[1] = 0.04*y[0] - 1.0e4*y[1]*y[2] - 3.0e7*y[1]*y[1]
		f[2] = 3.0e7 * y[1] * y[1]
		return nil
	}
	jac := func(dfdy *la.Triplet, dx, x float64, y []float64) error {
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
	//method := "Dopri5"
	method := "Radau5"
	xa, xb := 0.0, 0.3
	ya := []float64{1.0, 0.0, 0.0}
	ndim := len(ya)

	// allocate ODE object
	var o Solver
	o.Init(method, ndim, fcn, jac, nil, nil)
	o.SaveXY = true

	// tolerances and initial step size
	rtol := 1e-2
	atol := rtol * 1e-6
	o.SetTol(atol, rtol)
	o.IniH = 1.0e-6

	// solve problem
	y := make([]float64, ndim)
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
			plt.Subplot(ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(o.Xvalues[:s], o.Yvalues[j][:s], &plt.A{C: "r", M: ".", Ms: 2, Ls: "none", L: labelB})
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(ndim+1, 1, ndim+1)
		plt.Plot(o.Xvalues[1:s], o.Hvalues[1:s], &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl", "rober")
	}
}

func Test_ode04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode04: Hairer-Wanner VII-p376 Transistor Amplifier\n")
	// NOTE: from E Hairer's website, not the system in the book

	// data
	UE, UB, UF, ALPHA, BETA := 0.1, 6.0, 0.026, 0.99, 1.0e-6
	R0, R1, R2, R3, R4, R5 := 1000.0, 9000.0, 9000.0, 9000.0, 9000.0, 9000.0
	R6, R7, R8, R9 := 9000.0, 9000.0, 9000.0, 9000.0
	W := 2.0 * 3.141592654 * 100.0

	// initial values
	xa := 0.0
	ya := []float64{0.0,
		UB,
		UB / (R6/R5 + 1.0),
		UB / (R6/R5 + 1.0),
		UB,
		UB / (R2/R1 + 1.0),
		UB / (R2/R1 + 1.0),
		0.0}

	// endpoint of integration
	xb := 0.05
	//xb = 0.0123 // OK
	//xb = 0.01235 // !OK

	// right-hand side of the amplifier problem
	fcn := func(f []float64, dx, x float64, y []float64) error {
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
	jac := func(dfdy *la.Triplet, dx, x float64, y []float64) error {
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
	//method := "Dopri5"
	method := "Radau5"
	ndim := len(ya)
	numjac := false

	// ODE solver
	var o Solver
	o.Pll = true
	o.SaveXY = true

	if numjac {
		o.Init(method, ndim, fcn, nil, &M, nil)
	} else {
		o.Init(method, ndim, fcn, jac, &M, nil)
	}
	o.IniH = 1.0e-6 // initial step size

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
	if false {
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
			plt.Subplot(ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(o.Xvalues[:s], o.Yvalues[j][:s], &plt.A{C: "r", M: ".", Ms: 1, Ls: "none", L: labelB})
			plt.AxisXmax(0.05)
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(ndim+1, 1, ndim+1)
		plt.Plot(o.Xvalues[1:s], o.Hvalues[1:s], &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.AxisXmax(0.05)
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl", "hwamplifier")
	}
}
