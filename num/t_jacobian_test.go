// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
)

func sin(x float64) float64 { return math.Sin(x) }
func cos(x float64) float64 { return math.Cos(x) }

func TestJacobian01a(tst *testing.T) {

	// verbose()
	chk.PrintTitle("TestJacobian 01a (sparse)")

	ffcn := func(fx, x la.Vector) {
		fx[0] = math.Pow(x[0], 3.0) + x[1] - 1.0
		fx[1] = -x[0] + math.Pow(x[1], 3.0) + 1.0
	}
	Jfcn := func(dfdx *la.Triplet, x la.Vector) {
		dfdx.Start()
		dfdx.Put(0, 0, 3.0*x[0]*x[0])
		dfdx.Put(0, 1, 1.0)
		dfdx.Put(1, 0, -1.0)
		dfdx.Put(1, 1, 3.0*x[1]*x[1])
	}
	x := []float64{0.5, 0.5}
	CompareJac(tst, ffcn, Jfcn, x, 1e-7)
}

func TestJacobian02a(tst *testing.T) {

	// verbose()
	chk.PrintTitle("TestJacobian 02a (sparse)")

	ffcn := func(fx, x la.Vector) {
		fx[0] = 2.0*x[0] - x[1] + sin(x[2]) - cos(x[3]) - x[5]*x[5] - 1.0      // 0
		fx[1] = -x[0] + 2.0*x[1] + cos(x[2]) - sin(x[3]) + x[5] - 1.0          // 1
		fx[2] = x[0] + 3.0*x[1] + sin(x[3]) - cos(x[4]) - x[5]*x[5] - 1.0      // 2
		fx[3] = 2.0*x[0] + 4.0*x[1] + cos(x[3]) - cos(x[4]) + x[5] - 1.0       // 3
		fx[4] = x[0] + 5.0*x[1] - sin(x[2]) + sin(x[4]) - x[5]*x[5]*x[5] - 1.0 // 4
		fx[5] = x[0] + 6.0*x[1] - cos(x[2]) + cos(x[4]) + x[5] - 1.0           // 5
	}
	Jfcn := func(dfdx *la.Triplet, x la.Vector) {
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
	}
	x := []float64{5.0, 5.0, math.Pi, math.Pi, math.Pi, 5.0}
	CompareJac(tst, ffcn, Jfcn, x, 1e-6)
}

func TestJacobian03(tst *testing.T) {

	// verbose()
	chk.PrintTitle("TestJacobian 03 (dense)")

	ffcn := func(fx, x la.Vector) {
		fx[0] = math.Pow(x[0], 3.0) + x[1] - 1.0
		fx[1] = -x[0] + math.Pow(x[1], 3.0) + 1.0
	}
	Jfcn := func(dfdx *la.Matrix, x la.Vector) {
		dfdx.Set(0, 0, 3.0*x[0]*x[0])
		dfdx.Set(0, 1, 1.0)
		dfdx.Set(1, 0, -1.0)
		dfdx.Set(1, 1, 3.0*x[1]*x[1])
	}
	x := []float64{0.5, 0.5}
	CompareJacDense(tst, ffcn, Jfcn, x, 1e-7)
}
