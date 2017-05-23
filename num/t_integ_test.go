// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

func Test_trapz01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("trapz01")

	x := []float64{4, 6, 8}
	y := []float64{1, 2, 3}
	A := QuadDiscreteTrapzXY(x, y)
	chk.Scalar(tst, "A", 1e-17, A, 8)
}

func Test_trapz02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("trapz02")

	y := func(x float64) float64 {
		return math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0))
	}

	n := 11
	x := utl.LinSpace(0, 1, n)
	A := QuadDiscreteTrapzXF(x, y)
	A_ := QuadDiscreteTrapzRF(0, 1, n, y)
	io.Pforan("A  = %v\n", A)
	io.Pforan("A_ = %v\n", A_)
	Acor := 1.08306090851465 // right value is Acor := 1.08268158558
	chk.Scalar(tst, "A", 1e-15, A, Acor)
	chk.Scalar(tst, "A_", 1e-15, A_, Acor)
}

func Test_2dinteg01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("2dinteg01. volume of box")

	x := utl.LinSpace(-1, 1, 7)
	y := utl.LinSpace(-2, 2, 9)
	m, n := len(x), len(y)
	f := la.MatAlloc(m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			f[i][j] = 3.0
		}
	}
	dx, dy := x[1]-x[0], y[1]-y[0]
	Vt := QuadDiscreteTrapz2d(dx, dy, f)
	Vs := Simps2D(dx, dy, f)
	io.Pforan("Vt = %v\n", Vt)
	io.Pforan("Vs = %v\n", Vs)
	chk.Scalar(tst, "Vt", 1e-14, Vt, 24.0)
	chk.Scalar(tst, "Vs", 1e-14, Vs, 24.0)
}

func Test_2dinteg02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("2dinteg02. bidimensional integral")

	// Γ(1/4, 1)
	gamma_1div4_1 := 0.2462555291934987088744974330686081384629028737277219

	x := utl.LinSpace(0, 1, 11)
	y := utl.LinSpace(0, 1, 11)
	m, n := len(x), len(y)
	f := la.MatAlloc(m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			f[i][j] = 8.0 * math.Exp(-math.Pow(x[i], 2)-math.Pow(y[j], 4))
		}
	}
	dx, dy := x[1]-x[0], y[1]-y[0]
	Vt := QuadDiscreteTrapz2d(dx, dy, f)
	Vs := Simps2D(dx, dy, f)
	Vc := math.Sqrt(math.Pi) * math.Erf(1) * (math.Gamma(1.0/4.0) - gamma_1div4_1)
	io.Pforan("Vt = %v\n", Vt)
	io.Pforan("Vs = %v\n", Vs)
	io.Pfgreen("Vc = %v\n", Vc)
	chk.Scalar(tst, "Vt", 0.0114830435645548, Vt, Vc)
	chk.Scalar(tst, "Vs", 1e-4, Vs, Vc)

}

func Test_2dinteg03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("2dinteg03. ∫∫(1+8xy)dxdy")

	x := utl.LinSpace(0, 3, 5)
	y := utl.LinSpace(1, 2, 5)
	m, n := len(x), len(y)
	f := la.MatAlloc(m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			f[i][j] = 1.0 + 8.0*x[i]*y[j]
		}
	}
	dx, dy := x[1]-x[0], y[1]-y[0]
	Vt := QuadDiscreteTrapz2d(dx, dy, f)
	Vs := Simps2D(dx, dy, f)
	io.Pforan("Vt = %v\n", Vt)
	io.Pforan("Vs = %v\n", Vs)
	chk.Scalar(tst, "Vt", 1e-13, Vt, 57.0)
	chk.Scalar(tst, "Vs", 1e-13, Vs, 57.0)
}

func Test_2dinteg04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("2dinteg04. ∫∫y・sin(x)")

	x := utl.LinSpace(0, math.Pi/2.0, 11)
	y := utl.LinSpace(0, 1, 11)
	m, n := len(x), len(y)
	f := la.MatAlloc(m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			f[i][j] = math.Sin(x[i]) * y[j]
		}
	}
	dx, dy := x[1]-x[0], y[1]-y[0]
	Vt := QuadDiscreteTrapz2d(dx, dy, f)
	Vs := Simps2D(dx, dy, f)
	io.Pforan("Vt = %v\n", Vt)
	io.Pforan("Vs = %v\n", Vs)
	chk.Scalar(tst, "Vt", 0.00103, Vt, 0.5)
	chk.Scalar(tst, "Vs", 1e-5, Vs, 0.5)
}
