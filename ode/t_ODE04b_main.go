// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/ode"
)

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

func main() {

	mpi.Start(false)
	defer func() {
		mpi.Stop(false)
	}()

	if mpi.Rank() == 0 {
		chk.PrintTitle("Test ODE 04b (MPI)")
		io.Pfcyan("Hairer-Wanner VII-p376 Transistor Amplifier (MPI)\n")
		io.Pfcyan("(from E Hairer's website, not the system in the book)\n")
	}
	if mpi.Size() != 3 {
		chk.Panic(">> error: this test requires 3 MPI processors\n")
		return
	}

	// RIGHT-HAND SIDE OF THE AMPLIFIER PROBLEM
	w := make([]float64, 8) // workspace
	fcn := func(f []float64, x float64, y []float64, args ...interface{}) error {
		d := args[0].(*HWtransData)
		UET := d.UE * math.Sin(d.W*x)
		FAC1 := d.BETA * (math.Exp((y[3]-y[2])/d.UF) - 1.0)
		FAC2 := d.BETA * (math.Exp((y[6]-y[5])/d.UF) - 1.0)
		la.VecFill(f, 0)
		switch mpi.Rank() {
		case 0:
			f[0] = y[0] / d.R9
		case 1:
			f[1] = (y[1]-d.UB)/d.R8 + d.ALPHA*FAC1
			f[2] = y[2]/d.R7 - FAC1
		case 2:
			f[3] = y[3]/d.R5 + (y[3]-d.UB)/d.R6 + (1.0-d.ALPHA)*FAC1
			f[4] = (y[4]-d.UB)/d.R4 + d.ALPHA*FAC2
			f[5] = y[5]/d.R3 - FAC2
			f[6] = y[6]/d.R1 + (y[6]-d.UB)/d.R2 + (1.0-d.ALPHA)*FAC2
			f[7] = (y[7] - UET) / d.R0
		}
		mpi.AllReduceSum(f, w)
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
		switch mpi.Rank() {
		case 0:
			dfdy.Put(2+0-NU, 0, 1.0/d.R9)
			dfdy.Put(2+1-NU, 1, 1.0/d.R8)
			dfdy.Put(1+2-NU, 2, -d.ALPHA*FAC14)
			dfdy.Put(0+3-NU, 3, d.ALPHA*FAC14)
			dfdy.Put(2+2-NU, 2, 1.0/d.R7+FAC14)
		case 1:
			dfdy.Put(1+3-NU, 3, -FAC14)
			dfdy.Put(2+3-NU, 3, 1.0/d.R5+1.0/d.R6+(1.0-d.ALPHA)*FAC14)
			dfdy.Put(3+2-NU, 2, -(1.0-d.ALPHA)*FAC14)
			dfdy.Put(2+4-NU, 4, 1.0/d.R4)
			dfdy.Put(1+5-NU, 5, -d.ALPHA*FAC27)
		case 2:
			dfdy.Put(0+6-NU, 6, d.ALPHA*FAC27)
			dfdy.Put(2+5-NU, 5, 1.0/d.R3+FAC27)
			dfdy.Put(1+6-NU, 6, -FAC27)
			dfdy.Put(2+6-NU, 6, 1.0/d.R1+1.0/d.R2+(1.0-d.ALPHA)*FAC27)
			dfdy.Put(3+5-NU, 5, -(1.0-d.ALPHA)*FAC27)
			dfdy.Put(2+7-NU, 7, 1.0/d.R0)
		}
		return nil
	}

	// MATRIX "M"
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

	// WRITE FILE FUNCTION
	idxstp := 1
	var b bytes.Buffer
	out := func(first bool, dx, x float64, y []float64, args ...interface{}) error {
		if mpi.Rank() == 0 {
			if first {
				fmt.Fprintf(&b, "%6s%23s%23s%23s%23s%23s%23s%23s%23s%23s\n", "ns", "x", "y0", "y1", "y2", "y3", "y4", "y5", "y6", "y7")
			}
			fmt.Fprintf(&b, "%6d%23.15E", idxstp, x)
			for j := 0; j < len(y); j++ {
				fmt.Fprintf(&b, "%23.15E", y[j])
			}
			fmt.Fprintf(&b, "\n")
			idxstp += 1
		}
		return nil
	}
	defer func() {
		if mpi.Rank() == 0 {
			io.WriteFileD("/tmp/gosl", "hwamplifierB.res", &b)
		}
	}()

	// INITIAL DATA
	D, xa, xb, ya := HWtransIni()

	// SET ODE SOLVER
	silent := false
	fixstp := false
	//method := "Dopri5"
	method := "Radau5"
	ndim := len(ya)
	//numjac := true
	numjac := false
	var osol ode.ODE

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
	if mpi.Rank() == 0 {
		io.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))
	}
}
