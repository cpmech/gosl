// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
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
	silent := false
	xa, xb := 0.0, 1.5
	ya := []float64{0.0}
	ndim := len(ya)
	y := make([]float64, ndim)

	fcn := func(f []float64, dx, x float64, y []float64, args ...interface{}) error {
		f[0] = lam*y[0] - lam*math.Cos(x)
		return nil
	}

	jac := func(dfdy *la.Triplet, dx, x float64, y []float64, args ...interface{}) error {
		if dfdy.Max() == 0 {
			dfdy.Init(1, 1, 1)
		}
		dfdy.Start()
		dfdy.Put(0, 0, lam)
		return nil
	}

	MAXN := 1000 // maximum number of points for output

	out := func(first bool, dx, x float64, y []float64, args ...interface{}) error {
		k := args[0].(*int)
		X := args[1].(*[]float64)
		Y := args[2].(*[]float64)
		i := *k
		if i >= MAXN-1 {
			chk.Panic("cannot add more than %d points in slice", MAXN)
		}
		(*X)[i], (*Y)[i] = x, y[0]
		*k = i + 1
		return nil
	}

	// FwEuler
	io.Pforan(". . . FwEuler . . . \n")
	dx := 1.875 / 50.0
	copy(y, ya)
	k_FwEuler, X_FwEuler, Y_FwEuler := 0, make([]float64, MAXN), make([]float64, MAXN)
	var FwEuler Solver
	FwEuler.Init("FwEuler", ndim, fcn, jac, nil, out, silent)
	FwEuler.Solve(y, xa, xb, dx, true, &k_FwEuler, &X_FwEuler, &Y_FwEuler)

	// BwEuler
	io.Pforan(". . . BwEuler . . . \n")
	copy(y, ya)
	k_BwEuler, X_BwEuler, Y_BwEuler := 0, make([]float64, MAXN), make([]float64, MAXN)
	var BwEuler Solver
	//BwEuler.Init("BwEuler", ndim, fcn, nil, nil, out, silent)
	BwEuler.Init("BwEuler", ndim, fcn, jac, nil, out, silent)
	BwEuler.Solve(y, xa, xb, dx, true, &k_BwEuler, &X_BwEuler, &Y_BwEuler)

	// MoEuler
	io.Pforan(". . . MoEuler . . . \n")
	copy(y, ya)
	k_MoEuler, X_MoEuler, Y_MoEuler := 0, make([]float64, MAXN), make([]float64, MAXN)
	var MoEuler Solver
	MoEuler.Init("MoEuler", ndim, fcn, jac, nil, out, silent)
	MoEuler.Solve(y, xa, xb, xb-xa, false, &k_MoEuler, &X_MoEuler, &Y_MoEuler)

	// Dopri5
	io.Pforan(". . . Dopri5 . . . \n")
	copy(y, ya)
	k_Dopri5, X_Dopri5, Y_Dopri5 := 0, make([]float64, MAXN), make([]float64, MAXN)
	var Dopri5 Solver
	Dopri5.Init("Dopri5", ndim, fcn, jac, nil, out, silent)
	Dopri5.Solve(y, xa, xb, xb-xa, false, &k_Dopri5, &X_Dopri5, &Y_Dopri5)

	// Radau5
	io.Pforan(". . . Radau5 . . . \n")
	copy(y, ya)
	k_Radau5, X_Radau5, Y_Radau5 := 0, make([]float64, MAXN), make([]float64, MAXN)
	var Radau5 Solver
	Radau5.Init("Radau5", ndim, fcn, jac, nil, out, silent)
	Radau5.Solve(y, xa, xb, xb-xa, false, &k_Radau5, &X_Radau5, &Y_Radau5)

	if chk.Verbose {
		X := utl.LinSpace(xa, xb, 101)
		Y := make([]float64, len(X))
		for i := 0; i < len(X); i++ {
			Y[i] = -lam * (math.Sin(X[i]) - lam*math.Cos(X[i]) + lam*math.Exp(lam*X[i])) / (lam*lam + 1.0)
		}
		plt.SetForEps(0.75, 500)
		plt.Plot(X, Y, "'y.-', label='solution', lw=5")
		plt.Plot(X_FwEuler[:k_FwEuler], Y_FwEuler[:k_FwEuler], "'k.:',  label='FwEuler'")
		plt.Plot(X_BwEuler[:k_BwEuler], Y_BwEuler[:k_BwEuler], "'r.:',  label='BwEuler'")
		plt.Plot(X_MoEuler[:k_MoEuler], Y_MoEuler[:k_MoEuler], "'c+:',  label='MoEuler'")
		plt.Plot(X_Dopri5[:k_Dopri5], Y_Dopri5[:k_Dopri5], "'m.--', label='Dopri5'")
		plt.Plot(X_Radau5[:k_Radau5], Y_Radau5[:k_Radau5], "'bo-',  label='Radau5'")
		plt.Gll("$x$", "$y$", "")
		plt.SaveD("/tmp/gosl", "ode1.eps")
	}

	// work/correctness analysis
	//WcAnalysis("/tmp/gosl", "ode1", method, fcn, jac, nil, ycfcn, ya, xa, xb, []float64{5}, false)
}

// Hairer-Wanner VII-p5 Eq.(1.5) Van der Pol's Equation
func Test_ode02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode02: Hairer-Wanner VII-p5 Eq.(1.5) Van der Pol's Equation")

	// problem definition
	eps := 1.0e-6
	fcn := func(f []float64, dx, x float64, y []float64, args ...interface{}) error {
		f[0] = y[1]
		f[1] = ((1.0-y[0]*y[0])*y[1] - y[0]) / eps
		return nil
	}
	jac := func(dfdy *la.Triplet, dx, x float64, y []float64, args ...interface{}) error {
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
	silent := false
	fixstp := false
	//method := "Dopri5"
	method := "Radau5"
	numjac := false
	xa, xb := 0.0, 2.0
	ya := []float64{2.0, -0.6}
	ndim := len(ya)

	// structure to hold numerical results
	res := Results{Method: method}

	// allocate ODE object
	var o Solver
	if numjac {
		o.Init(method, ndim, fcn, nil, nil, SimpleOutput, silent)
	} else {
		o.Init(method, ndim, fcn, jac, nil, SimpleOutput, silent)
	}

	// tolerances and initial step size
	rtol := 1e-4
	atol := rtol
	o.IniH = 1.0e-4
	o.SetTol(atol, rtol)
	//o.NmaxSS = 2

	// solve problem
	y := make([]float64, ndim)
	copy(y, ya)
	t0 := time.Now()
	if fixstp {
		o.Solve(y, xa, xb, 0.05, fixstp, &res)
	} else {
		o.Solve(y, xa, xb, xb-xa, fixstp, &res)
	}
	io.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))

	// plot
	if chk.Verbose {
		plt.SetForEps(1.5, 400)
		args := "'b-', marker='.', lw=1, ms=4, clip_on=0"
		Plot("/tmp/gosl/ode", "vdpolA.eps", &res, nil, xa, xb, "", args, func() {
			_, T, err := io.ReadTable("data/vdpol_radau5_for.dat")
			if err != nil {
				chk.Panic("%v", err)
			}
			plt.Subplot(3, 1, 1)
			plt.Plot(T["x"], T["y0"], "'k+',label='reference',ms=7")
			plt.Subplot(3, 1, 2)
			plt.Plot(T["x"], T["y1"], "'k+',label='reference',ms=7")
		})
	}
}

// Hairer-Wanner VII-p3 Eq.(1.4) Robertson's Equation
func Test_ode03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode03: Hairer-Wanner VII-p3 Eq.(1.4) Robertson's Equation")

	fcn := func(f []float64, dx, x float64, y []float64, args ...interface{}) error {
		f[0] = -0.04*y[0] + 1.0e4*y[1]*y[2]
		f[1] = 0.04*y[0] - 1.0e4*y[1]*y[2] - 3.0e7*y[1]*y[1]
		f[2] = 3.0e7 * y[1] * y[1]
		return nil
	}
	jac := func(dfdy *la.Triplet, dx, x float64, y []float64, args ...interface{}) error {
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
	silent := false
	fixstp := false
	//method := "Dopri5"
	method := "Radau5"
	xa, xb := 0.0, 0.3
	ya := []float64{1.0, 0.0, 0.0}
	ndim := len(ya)

	// structure to hold numerical results
	res := Results{Method: method}

	// allocate ODE object
	var o Solver
	o.Init(method, ndim, fcn, jac, nil, SimpleOutput, silent)

	// tolerances and initial step size
	rtol := 1e-2
	atol := rtol * 1e-6
	o.SetTol(atol, rtol)
	o.IniH = 1.0e-6

	// solve problem
	y := make([]float64, ndim)
	copy(y, ya)
	if fixstp {
		o.Solve(y, xa, xb, 0.01, fixstp, &res)
	} else {
		o.Solve(y, xa, xb, xb-xa, fixstp, &res)
	}

	// plot
	if chk.Verbose {
		plt.SetForEps(1.5, 400)
		args := "'b-', marker='.', lw=1, clip_on=0"
		Plot("/tmp/gosl/ode", "rober.eps", &res, nil, xa, xb, "", args, func() {
			_, T, err := io.ReadTable("data/rober_radau5_cpp.dat")
			if err != nil {
				chk.Panic("%v", err)
			}
			plt.Subplot(4, 1, 1)
			plt.Plot(T["x"], T["y0"], "'k+',label='reference',ms=10")
			plt.Subplot(4, 1, 2)
			plt.Plot(T["x"], T["y1"], "'k+',label='reference',ms=10")
			plt.Subplot(4, 1, 3)
			plt.Plot(T["x"], T["y2"], "'k+',label='reference',ms=10")
		})
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
	fcn := func(f []float64, dx, x float64, y []float64, args ...interface{}) error {
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
	jac := func(dfdy *la.Triplet, dx, x float64, y []float64, args ...interface{}) error {
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
	silent := false
	fixstp := false
	//method := "Dopri5"
	method := "Radau5"
	ndim := len(ya)
	numjac := false

	// structure to hold numerical results
	res := Results{Method: method}

	// ODE solver
	var osol Solver
	osol.Pll = true

	if numjac {
		osol.Init(method, ndim, fcn, nil, &M, SimpleOutput, silent)
	} else {
		osol.Init(method, ndim, fcn, jac, &M, SimpleOutput, silent)
	}
	osol.IniH = 1.0e-6 // initial step size

	// set tolerances
	atol, rtol := 1e-11, 1e-5
	osol.SetTol(atol, rtol)

	// run
	t0 := time.Now()
	if fixstp {
		osol.Solve(ya, xa, xb, 0.01, fixstp, &res)
	} else {
		osol.Solve(ya, xa, xb, xb-xa, fixstp, &res)
	}
	io.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))

	// plot
	if chk.Verbose {
		plt.SetForEps(2.0, 400)
		args := "'b-', marker='.', lw=1, clip_on=0"
		Plot("/tmp/gosl/ode", "hwamplifier.eps", &res, nil, xa, xb, "", args, func() {
			_, T, err := io.ReadTable("data/radau5_hwamplifier.dat")
			if err != nil {
				chk.Panic("%v", err)
			}
			for j := 0; j < ndim; j++ {
				plt.Subplot(ndim+1, 1, j+1)
				plt.Plot(T["x"], T[io.Sf("y%d", j)], "'k+',label='reference',ms=10")
			}
		})
	}
}
