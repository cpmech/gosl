// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fdm"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

const (
	pi = math.Pi
)

func sin(x float64) float64 { return math.Sin(x) }
func cos(x float64) float64 { return math.Cos(x) }

func TestJacobian01a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestJacobian 01a")

	ffcn := func(fx, x []float64) error {
		fx[0] = math.Pow(x[0], 3.0) + x[1] - 1.0
		fx[1] = -x[0] + math.Pow(x[1], 3.0) + 1.0
		return nil
	}
	Jfcn := func(dfdx *la.Triplet, x []float64) error {
		dfdx.Start()
		dfdx.Put(0, 0, 3.0*x[0]*x[0])
		dfdx.Put(0, 1, 1.0)
		dfdx.Put(1, 0, -1.0)
		dfdx.Put(1, 1, 3.0*x[1]*x[1])
		return nil
	}
	x := []float64{0.5, 0.5}
	CompareJac(tst, ffcn, Jfcn, x, 1e-8)
}

func TestJacobian02a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestJacobian 02a")

	ffcn := func(fx, x []float64) error {
		fx[0] = 2.0*x[0] - x[1] + sin(x[2]) - cos(x[3]) - x[5]*x[5] - 1.0      // 0
		fx[1] = -x[0] + 2.0*x[1] + cos(x[2]) - sin(x[3]) + x[5] - 1.0          // 1
		fx[2] = x[0] + 3.0*x[1] + sin(x[3]) - cos(x[4]) - x[5]*x[5] - 1.0      // 2
		fx[3] = 2.0*x[0] + 4.0*x[1] + cos(x[3]) - cos(x[4]) + x[5] - 1.0       // 3
		fx[4] = x[0] + 5.0*x[1] - sin(x[2]) + sin(x[4]) - x[5]*x[5]*x[5] - 1.0 // 4
		fx[5] = x[0] + 6.0*x[1] - cos(x[2]) + cos(x[4]) + x[5] - 1.0           // 5
		return nil
	}
	Jfcn := func(dfdx *la.Triplet, x []float64) error {
		dfdx.Start()
		dfdx.Put(0, 0, 2.0)
		dfdx.Put(0, 1, -1.0)
		dfdx.Put(0, 2, cos(x[2]))
		dfdx.Put(0, 3, sin(x[3]))
		dfdx.Put(0, 5, -2.0*x[5])
		dfdx.Put(1, 0, -1.0)
		dfdx.Put(1, 1, 2.0)
		dfdx.Put(1, 2, -sin(x[2]))
		dfdx.Put(1, 3, -cos(x[3]))
		dfdx.Put(1, 5, 1.0)
		dfdx.Put(2, 0, 1.0)
		dfdx.Put(2, 1, 3.0)
		dfdx.Put(2, 3, cos(x[3]))
		dfdx.Put(2, 4, sin(x[4]))
		dfdx.Put(2, 5, -2.0*x[5])
		dfdx.Put(3, 0, 2.0)
		dfdx.Put(3, 1, 4.0)
		dfdx.Put(3, 3, -sin(x[3]))
		dfdx.Put(3, 4, sin(x[4]))
		dfdx.Put(3, 5, 1.0)
		dfdx.Put(4, 0, 1.0)
		dfdx.Put(4, 1, 5.0)
		dfdx.Put(4, 2, -cos(x[2]))
		dfdx.Put(4, 4, cos(x[4]))
		dfdx.Put(4, 5, -3.0*x[5]*x[5])
		dfdx.Put(5, 0, 1.0)
		dfdx.Put(5, 1, 6.0)
		dfdx.Put(5, 2, sin(x[2]))
		dfdx.Put(5, 4, -sin(x[4]))
		dfdx.Put(5, 5, 1.0)
		return nil
	}
	x := []float64{5.0, 5.0, pi, pi, pi, 5.0}
	CompareJac(tst, ffcn, Jfcn, x, 1e-6)
}

func TestJacobian03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestJacobian 03")

	// grid
	var g fdm.Grid2D
	//g.Init(1.0, 1.0, 4, 4)
	g.Init(1.0, 1.0, 6, 6)
	//g.Init(1.0, 1.0, 11, 11)

	// equations numbering
	var e fdm.Equations
	peq := utl.IntUnique(g.L, g.R, g.B, g.T)
	e.Init(g.N, peq)

	// K11 and K12
	var K11, K12 la.Triplet
	fdm.InitK11andK12(&K11, &K12, &e)

	// assembly
	F1 := make([]float64, e.N1)
	fdm.AssemblePoisson2d(&K11, &K12, F1, 1, 1, nil, &g, &e)

	// prescribed values
	U2 := make([]float64, e.N2)
	for _, eq := range g.L {
		U2[e.FR2[eq]] = 50.0
	}
	for _, eq := range g.R {
		U2[e.FR2[eq]] = 0.0
	}
	for _, eq := range g.B {
		U2[e.FR2[eq]] = 0.0
	}
	for _, eq := range g.T {
		U2[e.FR2[eq]] = 50.0
	}

	// functions
	k11 := K11.ToMatrix(nil)
	k12 := K12.ToMatrix(nil)
	ffcn := func(fU1, U1 []float64) error { // K11*U1 + K12*U2 - F1
		la.VecCopy(fU1, -1, F1)            // fU1 := (-F1)
		la.SpMatVecMulAdd(fU1, 1, k11, U1) // fU1 += K11*U1
		la.SpMatVecMulAdd(fU1, 1, k12, U2) // fU1 += K12*U2
		return nil
	}
	Jfcn := func(dfU1dU1 *la.Triplet, U1 []float64) error {
		fdm.AssemblePoisson2d(dfU1dU1, &K12, F1, 1, 1, nil, &g, &e)
		return nil
	}
	U1 := make([]float64, e.N1)
	CompareJac(tst, ffcn, Jfcn, U1, 0.0075)

	print_jac := false
	if print_jac {
		W1 := make([]float64, e.N1)
		fU1 := make([]float64, e.N1)
		ffcn(fU1, U1)
		var Jnum la.Triplet
		Jnum.Init(e.N1, e.N1, e.N1*e.N1)
		Jacobian(&Jnum, ffcn, U1, fU1, W1)
		la.PrintMat("K11 ", K11.ToMatrix(nil).ToDense(), "%g ", false)
		la.PrintMat("Jnum", Jnum.ToMatrix(nil).ToDense(), "%g ", false)
	}

	test_ffcn := false
	if test_ffcn {
		Uc := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 50.0, 25.0, 325.0 / 22.0, 100.0 / 11.0, 50.0 / 11.0,
			0.0, 50.0, 775.0 / 22.0, 25.0, 375.0 / 22.0, 100.0 / 11.0, 0.0, 50.0, 450.0 / 11.0, 725.0 / 22.0,
			25.0, 325.0 / 22.0, 0.0, 50.0, 500.0 / 11.0, 450.0 / 11.0, 775.0 / 22.0, 25.0, 0.0, 50.0, 50.0,
			50.0, 50.0, 50.0, 50.0,
		}
		for i := 0; i < e.N1; i++ {
			U1[i] = Uc[e.RF1[i]]
		}
		fU1 := make([]float64, e.N1)
		min, max := la.VecMinMax(fU1)
		io.Pf("min/max fU1 = %v\n", min, max)
	}
}
