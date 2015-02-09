// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"bytes"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

const (
	TEST_PLOT = false
)

// Hairer-Wanner VII-p2 Eq.(1.1)
func TestODE01(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	//utl.UseColors = false
	utl.TTitle("Test ODE 01")
	utl.Pfcyan("Hairer-Wanner VII-p2 Eq.(1.1) and work/correctness analysis\n")

	lam := -50.0
	silent := false
	xa, xb := 0.0, 1.5
	ya := []float64{0.0}
	ndim := len(ya)
	y := make([]float64, ndim)

	fcn := func(f []float64, x float64, y []float64, args ...interface{}) error {
		f[0] = lam*y[0] - lam*math.Cos(x)
		return nil
	}

	jac := func(dfdy *la.Triplet, x float64, y []float64, args ...interface{}) error {
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
			utl.Panic("cannot add more than %d points in slice", MAXN)
		}
		(*X)[i], (*Y)[i] = x, y[0]
		*k = i + 1
		return nil
	}

	// FwEuler
	utl.Pforan(". . . FwEuler . . . \n")
	dx := 1.875 / 50.0
	copy(y, ya)
	k_FwEuler, X_FwEuler, Y_FwEuler := 0, make([]float64, MAXN), make([]float64, MAXN)
	var FwEuler ODE
	FwEuler.Init("FwEuler", ndim, fcn, jac, nil, out, silent)
	FwEuler.Solve(y, xa, xb, dx, true, &k_FwEuler, &X_FwEuler, &Y_FwEuler)

	// BwEuler
	utl.Pforan(". . . BwEuler . . . \n")
	copy(y, ya)
	k_BwEuler, X_BwEuler, Y_BwEuler := 0, make([]float64, MAXN), make([]float64, MAXN)
	var BwEuler ODE
	//BwEuler.Init("BwEuler", ndim, fcn, nil, nil, out, silent)
	BwEuler.Init("BwEuler", ndim, fcn, jac, nil, out, silent)
	BwEuler.Solve(y, xa, xb, dx, true, &k_BwEuler, &X_BwEuler, &Y_BwEuler)

	// MoEuler
	utl.Pforan(". . . MoEuler . . . \n")
	copy(y, ya)
	k_MoEuler, X_MoEuler, Y_MoEuler := 0, make([]float64, MAXN), make([]float64, MAXN)
	var MoEuler ODE
	MoEuler.Init("MoEuler", ndim, fcn, jac, nil, out, silent)
	MoEuler.Solve(y, xa, xb, xb-xa, false, &k_MoEuler, &X_MoEuler, &Y_MoEuler)

	// Dopri5
	utl.Pforan(". . . Dopri5 . . . \n")
	copy(y, ya)
	k_Dopri5, X_Dopri5, Y_Dopri5 := 0, make([]float64, MAXN), make([]float64, MAXN)
	var Dopri5 ODE
	Dopri5.Init("Dopri5", ndim, fcn, jac, nil, out, silent)
	Dopri5.Solve(y, xa, xb, xb-xa, false, &k_Dopri5, &X_Dopri5, &Y_Dopri5)

	// Radau5
	utl.Pforan(". . . Radau5 . . . \n")
	copy(y, ya)
	k_Radau5, X_Radau5, Y_Radau5 := 0, make([]float64, MAXN), make([]float64, MAXN)
	var Radau5 ODE
	Radau5.Init("Radau5", ndim, fcn, jac, nil, out, silent)
	Radau5.Solve(y, xa, xb, xb-xa, false, &k_Radau5, &X_Radau5, &Y_Radau5)

	if TEST_PLOT {
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
func TestODE02a(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	//utl.UseColors = false
	utl.TTitle("Test ODE 02a")
	utl.Pfcyan("Hairer-Wanner VII-p5 Eq.(1.5) Van der Pol's Equation\n")

	eps := 1.0e-6
	fcn := func(f []float64, x float64, y []float64, args ...interface{}) error {
		f[0] = y[1]
		f[1] = ((1.0-y[0]*y[0])*y[1] - y[0]) / eps
		return nil
	}
	jac := func(dfdy *la.Triplet, x float64, y []float64, args ...interface{}) error {
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

	// data
	silent := false
	fixstp := false
	//method := "Dopri5"
	method := "Radau5"
	xa, xb := 0.0, 2.0
	ya := []float64{2.0, -0.6}
	ndim := len(ya)

	// output
	var b bytes.Buffer
	out := func(first bool, dx, x float64, y []float64, args ...interface{}) error {
		if first {
			fmt.Fprintf(&b, "%23s %23s %23s %23s\n", "dx", "x", "y0", "y1")
		}
		fmt.Fprintf(&b, "%23.15E %23.15E %23.15E %23.15E\n", dx, x, y[0], y[1])
		return nil
	}
	defer func() {
		extra := "d2 = Read('data/vdpol_radau5_for.dat')\n" +
			"subplot(3,1,1)\n" +
			"plot(d2['x'],d2['y0'],'k+',label='res',ms=10)\n" +
			"subplot(3,1,2)\n" +
			"plot(d2['x'],d2['y1'],'k+',label='res',ms=10)\n"
		Plot("/tmp/gosl", "vdpolA", method, &b, []int{0, 1}, ndim, nil, xa, xb, true, false, extra)
	}()

	// one run
	var o ODE
	numjac := true
	if numjac {
		o.Init(method, ndim, fcn, nil, nil, out, silent)
	} else {
		o.Init(method, ndim, fcn, jac, nil, out, silent)
	}

	// tolerances and initial step size
	rtol := 1e-4
	atol := rtol
	o.SetTol(atol, rtol)
	o.IniH = 1.0e-4

	//o.NmaxSS = 2

	y := make([]float64, ndim)
	copy(y, ya)
	t0 := time.Now()
	if fixstp {
		o.Solve(y, xa, xb, 0.05, fixstp)
	} else {
		o.Solve(y, xa, xb, xb-xa, fixstp)
	}
	utl.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))
}

// Hairer-Wanner VII-p3 Eq.(1.4) Robertson's Equation
func TestODE03(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	//utl.UseColors = false
	utl.TTitle("Test ODE 03")
	utl.Pfcyan("Hairer-Wanner VII-p3 Eq.(1.4) Robertson's Equation\n")

	fcn := func(f []float64, x float64, y []float64, args ...interface{}) error {
		f[0] = -0.04*y[0] + 1.0e4*y[1]*y[2]
		f[1] = 0.04*y[0] - 1.0e4*y[1]*y[2] - 3.0e7*y[1]*y[1]
		f[2] = 3.0e7 * y[1] * y[1]
		return nil
	}
	jac := func(dfdy *la.Triplet, x float64, y []float64, args ...interface{}) error {
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

	// output
	var b bytes.Buffer
	out := func(first bool, dx, x float64, y []float64, args ...interface{}) error {
		if first {
			fmt.Fprintf(&b, "%23s %23s %23s %23s %23s\n", "dx", "x", "y0", "y1", "y2")
		}
		fmt.Fprintf(&b, "%23.15E %23.15E %23.15E %23.15E %23.15E\n", dx, x, y[0], y[1], y[2])
		return nil
	}
	defer func() {
		extra := "d2 = Read('data/rober_radau5_cpp.dat')\n" +
			"subplot(4,1,1)\n" +
			"plot(d2['x'],d2['y0'],'k+',label='res',ms=10)\n" +
			"subplot(4,1,2)\n" +
			"plot(d2['x'],d2['y1'],'k+',label='res',ms=10)\n" +
			"subplot(4,1,3)\n" +
			"plot(d2['x'],d2['y2'],'k+',label='res',ms=10)\n"
		Plot("/tmp/gosl", "rober", method, &b, []int{0, 1, 2}, ndim, nil, xa, xb, true, false, extra)
	}()

	// one run
	var o ODE
	o.Init(method, ndim, fcn, jac, nil, out, silent)
	rtol := 1e-2
	atol := rtol * 1e-6
	o.SetTol(atol, rtol)
	o.IniH = 1.0e-6
	y := make([]float64, ndim)
	copy(y, ya)
	if fixstp {
		o.Solve(y, xa, xb, 0.01, fixstp)
	} else {
		o.Solve(y, xa, xb, xb-xa, fixstp)
	}
}

// DATA STRUCTURE FOR HW TRANSISTOR PROBLEM
type HWtransData struct {
	UE, UB, UF, ALPHA, BETA                float64
	R0, R1, R2, R3, R4, R5, R6, R7, R8, R9 float64
	W                                      float64
}

// INITIAL DATA FOR THE AMPLIFIER PROBLEM
func HWtransIni() (D HWtransData, xa, xb float64, ya []float64) {
	// DATA
	D.UE, D.UB, D.UF, D.ALPHA, D.BETA = 0.1, 6.0, 0.026, 0.99, 1.0e-6
	D.R0, D.R1, D.R2, D.R3, D.R4, D.R5 = 1000.0, 9000.0, 9000.0, 9000.0, 9000.0, 9000.0
	D.R6, D.R7, D.R8, D.R9 = 9000.0, 9000.0, 9000.0, 9000.0
	D.W = 2.0 * 3.141592654 * 100.0

	// INITIAL VALUES
	xa = 0.0
	ya = []float64{0.0,
		D.UB,
		D.UB / (D.R6/D.R5 + 1.0),
		D.UB / (D.R6/D.R5 + 1.0),
		D.UB,
		D.UB / (D.R2/D.R1 + 1.0),
		D.UB / (D.R2/D.R1 + 1.0),
		0.0}

	// ENDPOINT OF INTEGRATION
	xb = 0.05
	//xb = 0.0123 // OK
	//xb = 0.01235 // !OK
	return
}

// Hairer-Wanner VII-p376 Transistor Amplifier
// (from E Hairer's website, not the system in the book)
func TestODE04a(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	//utl.UseColors = false
	utl.TTitle("Test ODE 04a")
	utl.Pfcyan("Hairer-Wanner VII-p376 Transistor Amplifier\n")
	utl.Pfcyan("(from E Hairer's website, not the system in the book)\n")

	// RIGHT-HAND SIDE OF THE AMPLIFIER PROBLEM
	fcn := func(f []float64, x float64, y []float64, args ...interface{}) error {
		d := args[0].(*HWtransData)
		UET := d.UE * math.Sin(d.W*x)
		FAC1 := d.BETA * (math.Exp((y[3]-y[2])/d.UF) - 1.0)
		FAC2 := d.BETA * (math.Exp((y[6]-y[5])/d.UF) - 1.0)
		f[0] = y[0] / d.R9
		f[1] = (y[1]-d.UB)/d.R8 + d.ALPHA*FAC1
		f[2] = y[2]/d.R7 - FAC1
		f[3] = y[3]/d.R5 + (y[3]-d.UB)/d.R6 + (1.0-d.ALPHA)*FAC1
		f[4] = (y[4]-d.UB)/d.R4 + d.ALPHA*FAC2
		f[5] = y[5]/d.R3 - FAC2
		f[6] = y[6]/d.R1 + (y[6]-d.UB)/d.R2 + (1.0-d.ALPHA)*FAC2
		f[7] = (y[7] - UET) / d.R0
		return nil
	}

	// JACOBIAN OF THE AMPLIFIER PROBLEM
	jac := func(dfdy *la.Triplet, x float64, y []float64, args ...interface{}) error {
		d := args[0].(*HWtransData)
		FAC14 := d.BETA * math.Exp((y[3]-y[2])/d.UF) / d.UF
		FAC27 := d.BETA * math.Exp((y[6]-y[5])/d.UF) / d.UF
		if dfdy.Max() == 0 {
			dfdy.Init(8, 8, 16)
		}
		NU := 2
		dfdy.Start()
		dfdy.Put(2+0-NU, 0, 1.0/d.R9)
		dfdy.Put(2+1-NU, 1, 1.0/d.R8)
		dfdy.Put(1+2-NU, 2, -d.ALPHA*FAC14)
		dfdy.Put(0+3-NU, 3, d.ALPHA*FAC14)
		dfdy.Put(2+2-NU, 2, 1.0/d.R7+FAC14)
		dfdy.Put(1+3-NU, 3, -FAC14)
		dfdy.Put(2+3-NU, 3, 1.0/d.R5+1.0/d.R6+(1.0-d.ALPHA)*FAC14)
		dfdy.Put(3+2-NU, 2, -(1.0-d.ALPHA)*FAC14)
		dfdy.Put(2+4-NU, 4, 1.0/d.R4)
		dfdy.Put(1+5-NU, 5, -d.ALPHA*FAC27)
		dfdy.Put(0+6-NU, 6, d.ALPHA*FAC27)
		dfdy.Put(2+5-NU, 5, 1.0/d.R3+FAC27)
		dfdy.Put(1+6-NU, 6, -FAC27)
		dfdy.Put(2+6-NU, 6, 1.0/d.R1+1.0/d.R2+(1.0-d.ALPHA)*FAC27)
		dfdy.Put(3+5-NU, 5, -(1.0-d.ALPHA)*FAC27)
		dfdy.Put(2+7-NU, 7, 1.0/d.R0)
		return nil
	}

	// MATRIX "M"
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

	// WRITE FILE FUNCTION
	idxstp := 1
	var b bytes.Buffer
	out := func(first bool, dx, x float64, y []float64, args ...interface{}) error {
		if first {
			fmt.Fprintf(&b, "%6s%23s%23s%23s%23s%23s%23s%23s%23s%23s\n", "ns", "x", "y0", "y1", "y2", "y3", "y4", "y5", "y6", "y7")
		}
		fmt.Fprintf(&b, "%6d%23.15E", idxstp, x)
		for j := 0; j < len(y); j++ {
			fmt.Fprintf(&b, "%23.15E", y[j])
		}
		fmt.Fprintf(&b, "\n")
		idxstp += 1
		return nil
	}
	defer func() {
		utl.WriteFileD("/tmp/gosl", "hwamplifierA.res", &b)
	}()

	// INITIAL DATA
	D, xa, xb, ya := HWtransIni()

	// SET ODE SOLVER
	silent := false
	fixstp := false
	//method := "Dopri5"
	method := "Radau5"
	ndim := len(ya)
	numjac := false
	var osol ODE

	osol.Pll = true

	if numjac {
		osol.Init(method, ndim, fcn, nil, &M, out, silent)
	} else {
		osol.Init(method, ndim, fcn, jac, &M, out, silent)
	}
	osol.IniH = 1.0e-6 // initial step size

	// SET TOLERANCES
	atol, rtol := 1e-11, 1e-5
	osol.SetTol(atol, rtol)

	// RUN
	t0 := time.Now()
	if fixstp {
		osol.Solve(ya, xa, xb, 0.01, fixstp, &D)
	} else {
		osol.Solve(ya, xa, xb, xb-xa, fixstp, &D)
	}
	utl.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))
}
