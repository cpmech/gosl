// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/ode"
	"github.com/cpmech/gosl/plt"
)

func main() {

	mpi.Start(false)
	defer func() {
		mpi.Stop(false)
	}()

	if mpi.Rank() == 0 {
		chk.PrintTitle("ode04: Hairer-Wanner VII-p376 Transistor Amplifier\n")
	}
	if mpi.Size() != 3 {
		chk.Panic(">> error: this test requires 3 MPI processors\n")
		return
	}

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
	w := make([]float64, 8) // workspace
	fcn := func(f []float64, dx, x float64, y []float64, args ...interface{}) error {
		UET := UE * math.Sin(W*x)
		FAC1 := BETA * (math.Exp((y[3]-y[2])/UF) - 1.0)
		FAC2 := BETA * (math.Exp((y[6]-y[5])/UF) - 1.0)
		la.VecFill(f, 0)
		switch mpi.Rank() {
		case 0:
			f[0] = y[0] / R9
		case 1:
			f[1] = (y[1]-UB)/R8 + ALPHA*FAC1
			f[2] = y[2]/R7 - FAC1
		case 2:
			f[3] = y[3]/R5 + (y[3]-UB)/R6 + (1.0-ALPHA)*FAC1
			f[4] = (y[4]-UB)/R4 + ALPHA*FAC2
			f[5] = y[5]/R3 - FAC2
			f[6] = y[6]/R1 + (y[6]-UB)/R2 + (1.0-ALPHA)*FAC2
			f[7] = (y[7] - UET) / R0
		}
		mpi.AllReduceSum(f, w)
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
		switch mpi.Rank() {
		case 0:
			dfdy.Put(2+0-NU, 0, 1.0/R9)
			dfdy.Put(2+1-NU, 1, 1.0/R8)
			dfdy.Put(1+2-NU, 2, -ALPHA*FAC14)
			dfdy.Put(0+3-NU, 3, ALPHA*FAC14)
			dfdy.Put(2+2-NU, 2, 1.0/R7+FAC14)
		case 1:
			dfdy.Put(1+3-NU, 3, -FAC14)
			dfdy.Put(2+3-NU, 3, 1.0/R5+1.0/R6+(1.0-ALPHA)*FAC14)
			dfdy.Put(3+2-NU, 2, -(1.0-ALPHA)*FAC14)
			dfdy.Put(2+4-NU, 4, 1.0/R4)
			dfdy.Put(1+5-NU, 5, -ALPHA*FAC27)
		case 2:
			dfdy.Put(0+6-NU, 6, ALPHA*FAC27)
			dfdy.Put(2+5-NU, 5, 1.0/R3+FAC27)
			dfdy.Put(1+6-NU, 6, -FAC27)
			dfdy.Put(2+6-NU, 6, 1.0/R1+1.0/R2+(1.0-ALPHA)*FAC27)
			dfdy.Put(3+5-NU, 5, -(1.0-ALPHA)*FAC27)
			dfdy.Put(2+7-NU, 7, 1.0/R0)
		}
		return nil
	}

	// matrix "M"
	c1, c2, c3, c4, c5 := 1.0e-6, 2.0e-6, 3.0e-6, 4.0e-6, 5.0e-6
	var M la.Triplet
	M.Init(8, 8, 14)
	M.Start()
	NU := 1
	switch mpi.Rank() {
	case 0:
		M.Put(1+0-NU, 0, -c5)
		M.Put(0+1-NU, 1, c5)
		M.Put(2+0-NU, 0, c5)
		M.Put(1+1-NU, 1, -c5)
		M.Put(1+2-NU, 2, -c4)
		M.Put(1+3-NU, 3, -c3)
	case 1:
		M.Put(0+4-NU, 4, c3)
		M.Put(2+3-NU, 3, c3)
		M.Put(1+4-NU, 4, -c3)
	case 2:
		M.Put(1+5-NU, 5, -c2)
		M.Put(1+6-NU, 6, -c1)
		M.Put(0+7-NU, 7, c1)
		M.Put(2+6-NU, 6, c1)
		M.Put(1+7-NU, 7, -c1)
	}

	// flags
	silent := false
	fixstp := false
	//method := "Dopri5"
	method := "Radau5"
	ndim := len(ya)
	numjac := false

	// structure to hold numerical results
	res := ode.Results{Method: method}

	// ODE solver
	var osol ode.Solver
	osol.Pll = true

	// solve problem
	if numjac {
		osol.Init(method, ndim, fcn, nil, &M, ode.SimpleOutput, silent)
	} else {
		osol.Init(method, ndim, fcn, jac, &M, ode.SimpleOutput, silent)
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

	// plot
	if mpi.Rank() == 0 {
		io.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))
		plt.SetForEps(2.0, 400)
		args := "'b-', marker='.', lw=1, clip_on=0"
		ode.Plot("/tmp/gosl/ode", "hwamplifier_mpi.eps", &res, nil, xa, xb, "", args, func() {
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
