// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func Test_DiscSimpson01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("DiscSimpson01. Simple function")

	y := func(x float64) (res float64, err error) {
		res = math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0))
		return
	}

	n := 1000
	A, _ := QuadDiscreteSimpsonRF(0, 1, n, y)
	io.Pforan("A  = %v\n", A)
	Acor := 1.08268158558
	chk.Float64(tst, "A", 1e-11, A, Acor)
}

func Test_DiscTrapz01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("DiscTrapz01. Discrete trapezoidal XY")

	x := []float64{4, 6, 8}
	y := []float64{1, 2, 3}
	A := QuadDiscreteTrapzXY(x, y)
	chk.Float64(tst, "A", 1e-17, A, 8)
}

func Test_DiscTrapz02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("DiscTrapz02. Discrete trapezoidal XF and RF")

	y := func(x float64) (res float64, err error) {
		res = math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0))
		return
	}

	n := 11
	x := utl.LinSpace(0, 1, n)
	A, _ := QuadDiscreteTrapzXF(x, y)
	A1, _ := QuadDiscreteTrapzRF(0, 1, n, y)
	io.Pforan("A  = %v\n", A)
	io.Pforan("A1 = %v\n", A1)
	Acor := 1.08306090851465 // right value is Acor := 1.08268158558
	chk.Float64(tst, "A", 1e-15, A, Acor)
	chk.Float64(tst, "A1", 1e-15, A1, Acor)
}

func Test_Disc2dInteg01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Disc2dInteg01. Discrete 2D. Volume of box")

	x := utl.LinSpace(-1, 1, 7)
	y := utl.LinSpace(-2, 2, 9)
	m, n := len(x), len(y)
	f := utl.Alloc(m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			f[i][j] = 3.0
		}
	}
	dx, dy := x[1]-x[0], y[1]-y[0]
	Vt := QuadDiscreteTrapz2d(dx, dy, f)
	Vs := QuadDiscreteSimps2d(dx, dy, f)
	io.Pforan("Vt = %v\n", Vt)
	io.Pforan("Vs = %v\n", Vs)
	chk.Float64(tst, "Vt", 1e-14, Vt, 24.0)
	chk.Float64(tst, "Vs", 1e-14, Vs, 24.0)
}

func Test_Disc2dInteg02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Disc2dInteg02. Discrete 2D. Exp function")

	// Γ(1/4, 1)
	gamma1div4o1 := 0.2462555291934987088744974330686081384629028737277219

	x := utl.LinSpace(0, 1, 11)
	y := utl.LinSpace(0, 1, 11)
	m, n := len(x), len(y)
	f := utl.Alloc(m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			f[i][j] = 8.0 * math.Exp(-math.Pow(x[i], 2)-math.Pow(y[j], 4))
		}
	}
	dx, dy := x[1]-x[0], y[1]-y[0]
	Vt := QuadDiscreteTrapz2d(dx, dy, f)
	Vs := QuadDiscreteSimps2d(dx, dy, f)
	Vc := math.Sqrt(math.Pi) * math.Erf(1) * (math.Gamma(1.0/4.0) - gamma1div4o1)
	io.Pforan("Vt = %v\n", Vt)
	io.Pforan("Vs = %v\n", Vs)
	io.Pfgreen("Vc = %v\n", Vc)
	chk.Float64(tst, "Vt", 0.0114830435645548, Vt, Vc)
	chk.Float64(tst, "Vs", 1e-4, Vs, Vc)

}

func Test_Disc2dInteg03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Disc2dInteg03. Discrete 2D. ∫∫(1+8xy)dxdy")

	x := utl.LinSpace(0, 3, 5)
	y := utl.LinSpace(1, 2, 5)
	m, n := len(x), len(y)
	f := utl.Alloc(m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			f[i][j] = 1.0 + 8.0*x[i]*y[j]
		}
	}
	dx, dy := x[1]-x[0], y[1]-y[0]
	Vt := QuadDiscreteTrapz2d(dx, dy, f)
	Vs := QuadDiscreteSimps2d(dx, dy, f)
	io.Pforan("Vt = %v\n", Vt)
	io.Pforan("Vs = %v\n", Vs)
	chk.Float64(tst, "Vt", 1e-13, Vt, 57.0)
	chk.Float64(tst, "Vs", 1e-13, Vs, 57.0)
}

func Test_Disc2dInteg04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Disc2dInteg04. Discrete 2D. ∫∫y・sin(x)")

	x := utl.LinSpace(0, math.Pi/2.0, 11)
	y := utl.LinSpace(0, 1, 11)
	m, n := len(x), len(y)
	f := utl.Alloc(m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			f[i][j] = math.Sin(x[i]) * y[j]
		}
	}
	dx, dy := x[1]-x[0], y[1]-y[0]
	Vt := QuadDiscreteTrapz2d(dx, dy, f)
	Vs := QuadDiscreteSimps2d(dx, dy, f)
	io.Pforan("Vt = %v\n", Vt)
	io.Pforan("Vs = %v\n", Vs)
	chk.Float64(tst, "Vt", 0.00103, Vt, 0.5)
	chk.Float64(tst, "Vs", 1e-5, Vs, 0.5)
}
