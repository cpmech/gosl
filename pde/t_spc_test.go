// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestSpc01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Spc01. Auu matrix after homogeneous bcs")

	// lagrange interpolators (5x5 grid)
	l := fun.NewLagIntSet(2, []int{4, 4}, []string{"cgl", "cgl"})

	// grid
	g := new(gm.Grid)
	g.RectSet2d(l[0].X, l[1].X)

	// solver
	p := dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}
	s := NewSpcLaplacian(p, l, g, nil)

	// essential boundary conditions
	s.AddBc(true, 10, 0.0, nil) // left
	s.AddBc(true, 11, 0.0, nil) // right
	s.AddBc(true, 20, 0.0, nil) // bottom
	s.AddBc(true, 21, 0.0, nil) // top

	// assemble
	s.Assemble(false)
	Duu := s.Eqs.Auu.ToDense()
	io.Pf("%v\n", Duu.Print("%7.2f"))

	// check
	zzz := 0.0
	chk.Deep2(tst, "Auu", 1e-14, Duu.GetDeep2(), [][]float64{
		{-28, +6., -2. /* */, +6., zzz, zzz /* */, -2., zzz, zzz},
		{+4., -20, +4. /* */, zzz, +6., zzz /* */, zzz, -2., zzz},
		{-2., +6., -28 /* */, zzz, zzz, +6. /* */, zzz, zzz, -2.},

		{+4., zzz, zzz /* */, -20, +6., -2. /* */, +4., zzz, zzz},
		{zzz, +4., zzz /* */, +4., -12, +4. /* */, zzz, +4., zzz},
		{zzz, zzz, +4. /* */, -2., +6., -20 /* */, zzz, zzz, +4.},

		{-2., zzz, zzz /* */, +6., zzz, zzz /* */, -28, +6., -2.},
		{zzz, -2., zzz /* */, zzz, +6., zzz /* */, +4., -20, +4.},
		{zzz, zzz, -2. /* */, zzz, zzz, +6. /* */, -2., +6., -28},
	})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		gp := gm.GridPlotter{G: g, WithVids: true}
		gp.Draw()
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/pde", "spc01")
	}
}

func TestSpc02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Spc02. simple Dirichlet problem (unif-grid / Laplace)")

	// solve problem
	//    ∂²u     ∂²u
	//    ———  +  ——— = 0    with   u(x,0)=1   u(3,y)=2   u(x,3)=2   u(0,y)=1
	//    ∂x²     ∂y²               (bottom)   (right)    (top)      (left)

	// Lagrange interpolators
	lis := fun.NewLagIntSet(2, []int{3, 3}, []string{"uni", "uni"})

	// grid
	g := new(gm.Grid)
	g.RectSet2dU([]float64{0, 0}, []float64{2, 2}, lis[0].X, lis[1].X)

	// solver
	p := dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}
	s := NewSpcLaplacian(p, lis, g, nil)

	// essential boundary conditions
	s.AddBc(true, 10, 1.0, nil) // left
	s.AddBc(true, 11, 2.0, nil) // right
	s.AddBc(true, 20, 1.0, nil) // bottom
	s.AddBc(true, 21, 2.0, nil) // top

	// set equationa and assemble A matrix
	reactions := true
	s.Assemble(reactions)

	// check equations (must be after Assemble)
	chk.Int(tst, "number of equations == number of nodes", s.Eqs.N, 16)
	chk.Ints(tst, "UtoF", s.Eqs.UtoF, []int{5, 6, 9, 10})
	chk.Ints(tst, "KtoF", s.Eqs.KtoF, []int{0, 1, 2, 3, 4, 7, 8, 11, 12, 13, 14, 15})

	// check Duu
	Duu := s.Eqs.Auu.ToDense()
	io.Pf("Auu =\n%v\n", Duu.Print("%7.2f"))
	chk.Deep2(tst, "Auu", 1e-17, Duu.GetDeep2(), [][]float64{
		{-9.00, +2.25, +2.25, +0.00}, // 0 ⇒ node (1,1) 5
		{+2.25, -9.00, +0.00, +2.25}, // 1 ⇒ node (2,1) 6
		{+2.25, +0.00, -9.00, +2.25}, // 2 ⇒ node (1,2) 9
		{+0.00, +2.25, +2.25, -9.00}, // 3 ⇒ node (2,2) 10
	})

	// solve
	u, f := s.SolveSteady(reactions)
	chk.Array(tst, "xk", 1e-17, s.Eqs.Xk, []float64{1, 1, 1, 1, 1, 2, 1, 2, 2, 2, 2, 2})
	chk.Array(tst, "u", 1e-15, u, []float64{1, 1, 1, 1, 1, 1.25, 1.5, 2, 1, 1.5, 1.75, 2, 2, 2, 2, 2})

	// check f
	sFull := NewSpcLaplacian(p, lis, g, nil)
	sFull.Assemble(false)
	K := sFull.Eqs.Auu.ToMatrix(nil)
	Fref := la.NewVector(g.Size())
	la.SpMatVecMul(Fref, 1.0, K, u)
	chk.Array(tst, "f", 1e-14, f, Fref)

	// get results over grid
	uu := g.MapMeshgrid2d(u)
	chk.Deep2(tst, "uu", 1e-15, uu, [][]float64{
		{1.00, 1.00, 1.00, 1.00},
		{1.00, 1.25, 1.50, 2.00},
		{1.00, 1.50, 1.75, 2.00},
		{2.00, 2.00, 2.00, 2.00},
	})

	// plot
	if chk.Verbose {
		gp := gm.GridPlotter{G: g, WithVids: true}
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		gp.Draw()
		plt.ContourL(gp.X2d, gp.Y2d, uu, nil)
		plt.Gll("$x$", "$y$", nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/pde", "spc02")
	}
}

func TestSpc03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Spc03. Trefethen's p16")

	// solve problem
	//    ∂²u     ∂²u
	//    ———  +  ——— = 10 sin(8x⋅(y-1))    with   homogeneous BCs
	//    ∂x²     ∂y²

	// Lagrange interpolators
	lis := fun.NewLagIntSet(2, []int{4, 4}, []string{"cgl", "cgl"})

	// grid
	g := new(gm.Grid)
	g.RectSet2dU([]float64{-1, -1}, []float64{+1, +1}, lis[0].X, lis[1].X)

	// solver
	p := dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}
	source := func(x la.Vector, t float64) float64 {
		return 10.0 * math.Sin(8.0*x[0]*(x[1]-1.0))
	}
	s := NewSpcLaplacian(p, lis, g, source)

	// essential boundary conditions
	s.AddBc(true, 10, 0.0, nil) // left
	s.AddBc(true, 11, 0.0, nil) // right
	s.AddBc(true, 20, 0.0, nil) // bottom
	s.AddBc(true, 21, 0.0, nil) // top

	// set equationa and assemble A matrix
	reactions := false
	s.Assemble(reactions)

	// solve
	u, _ := s.SolveSteady(reactions)

	// check
	uu := g.MapMeshgrid2d(u)
	chk.Deep2(tst, "uu", 1e-14, uu, [][]float64{
		{0, +0.000000000000000, +0, +0.000000000000000, 0},
		{0, +0.181363633964132, +0, -0.181363633964131, 0},
		{0, +0.292713394079481, +0, -0.292713394079479, 0},
		{0, -0.329593843114906, +0, +0.329593843114906, 0},
		{0, +0.000000000000000, +0, +0.000000000000000, 0},
	})

	// draw surface
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 600})

		// data for dense output
		npts := 101
		xdense := make([]float64, npts)
		ydense := make([]float64, npts)
		zdense := make([]float64, npts)

		// draw lines along x
		nx, ny := g.Npts(0), g.Npts(1)
		X := lis[0].X
		Y := make([]float64, ny)
		Z := make([]float64, nx)
		for j := 0; j < ny; j++ {
			utl.Fill(Y, lis[1].X[j])
			for i := 0; i < nx; i++ {
				Z[i] = uu[j][i]
			}
			lis[0].U = Z
			for k := 0; k < npts; k++ {
				xdense[k] = -1.0 + 2.0*float64(k)/float64(npts-1)
				ydense[k] = lis[1].X[j]
				zdense[k] = lis[0].I(xdense[k])
			}
			plt.Plot3dLine(xdense, ydense, zdense, &plt.A{C: "gray"})
			plt.Plot3dPoints(X, Y, Z, &plt.A{C: plt.C(0, 0), M: ".", Mec: plt.C(0, 0)})
		}

		// draw lines along y
		X = make([]float64, nx)
		Y = lis[1].X
		Z = make([]float64, ny)
		for i := 0; i < nx; i++ {
			utl.Fill(X, lis[0].X[i])
			for j := 0; j < ny; j++ {
				Z[j] = uu[j][i]
			}
			lis[1].U = Z
			for k := 0; k < npts; k++ {
				xdense[k] = lis[0].X[i]
				ydense[k] = -1.0 + 2.0*float64(k)/float64(npts-1)
				zdense[k] = lis[1].I(ydense[k])
			}
			plt.Plot3dLine(xdense, ydense, zdense, &plt.A{C: "gray"})
			plt.Plot3dPoints(X, Y, Z, &plt.A{C: plt.C(2, 0), M: "+", Mec: plt.C(0, 0)})
		}

		plt.AxisRange3d(-1, 1, -1, 1, -0.5, 0.5)
		plt.Camera(45, 220, nil)
		plt.Save("/tmp/gosl/pde", "spc03")
	}
}
