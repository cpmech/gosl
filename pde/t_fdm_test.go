// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func TestFdm01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Fdm01. Full Auu matrix.")

	// 3x3 grid ⇒ 9 equations
	g := new(gm.Grid)
	g.RectGenUniform([]float64{0, 0}, []float64{2, 2}, []int{3, 3})

	// operator
	s := NewFdmLaplacian(dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}, g, nil)

	// assemble
	s.Assemble(false)
	Duu := s.Eqs.Auu.ToDense()
	io.Pf("%v\n", Duu.Print("%4g"))

	// check
	chk.Deep2(tst, "Auu", 1e-17, Duu.GetDeep2(), [][]float64{
		{-4, +2, +0, +2, +0, +0, +0, +0, +0}, // 0
		{+1, -4, +1, +0, +2, +0, +0, +0, +0}, // 1
		{+0, +2, -4, +0, +0, +2, +0, +0, +0}, // 2
		{+1, +0, +0, -4, +2, +0, +1, +0, +0}, // 3
		{+0, +1, +0, +1, -4, +1, +0, +1, +0}, // 4
		{+0, +0, +1, +0, +2, -4, +0, +0, +1}, // 5
		{+0, +0, +0, +2, +0, +0, -4, +2, +0}, // 6
		{+0, +0, +0, +0, +2, +0, +1, -4, +1}, // 7
		{+0, +0, +0, +0, +0, +2, +0, +2, -4}, // 8
	})
}

func TestFdm02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Fdm02. simple Dirichlet problem (Laplace)")

	// solve problem
	//    ∂²u     ∂²u
	//    ———  +  ——— = 0    with   u(x,0)=1   u(3,y)=2   u(x,3)=2   u(0,y)=1
	//    ∂x²     ∂y²               (bottom)   (right)    (top)      (left)

	// 4x4 grid ⇒ 16 equations
	g := new(gm.Grid)
	g.RectGenUniform([]float64{0, 0}, []float64{3, 3}, []int{4, 4})

	// solver
	p := dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}
	s := NewFdmLaplacian(p, g, nil)

	// essential boundary conditions
	s.AddBc(true, 10, 1.0, nil) // left
	s.AddBc(true, 11, 2.0, nil) // right
	s.AddBc(true, 20, 1.0, nil) // bottom
	s.AddBc(true, 21, 2.0, nil) // top

	// set equations and assemble A matrix
	reactions := true
	s.Assemble(reactions)

	// check equations (must be after Assemble)
	chk.Int(tst, "number of equations == number of nodes", s.Eqs.N, 16)
	chk.Ints(tst, "UtoF", s.Eqs.UtoF, []int{5, 6, 9, 10})
	chk.Ints(tst, "KtoF", s.Eqs.KtoF, []int{0, 1, 2, 3, 4, 7, 8, 11, 12, 13, 14, 15})

	// check Duu
	Duu := s.Eqs.Auu.ToDense()
	io.Pf("Auu =\n%v\n", Duu.Print("%4g"))
	chk.Deep2(tst, "Auu", 1e-17, Duu.GetDeep2(), [][]float64{
		{-4, +1, +1, +0}, // 0 ⇒ node (1,1) 5
		{+1, -4, +0, +1}, // 1 ⇒ node (2,1) 6
		{+1, +0, -4, +1}, // 2 ⇒ node (1,2) 9
		{+0, +1, +1, -4}, // 3 ⇒ node (2,2) 10
	})

	// solve
	u, f := s.SolveSteady(reactions)
	io.Pf("u = %v\n", u)
	io.Pf("f = %v\n", f)
	chk.Array(tst, "u", 1e-15, u, []float64{1, 1, 1, 1, 1, 1.25, 1.5, 2, 1, 1.5, 1.75, 2, 2, 2, 2, 2})

	// check f
	sFull := NewFdmLaplacian(p, g, nil)
	sFull.Assemble(false)
	K := sFull.Eqs.Auu.ToMatrix(nil)
	Fref := la.NewVector(g.Size())
	la.SpMatVecMul(Fref, 1.0, K, u)
	chk.Array(tst, "f", 1e-15, f, Fref)

	// get results over grid
	uu := g.MapMeshgrid2d(u)
	chk.Deep2(tst, "uu", 1e-15, uu, [][]float64{
		{1.00, 1.00, 1.00, 1.00},
		{1.00, 1.25, 1.50, 2.00},
		{1.00, 1.50, 1.75, 2.00},
		{2.00, 2.00, 2.00, 2.00},
	})

	// solve again without reactions
	reactions = false
	sNoreact := NewFdmLaplacian(p, g, nil)
	sNoreact.AddBc(true, 10, 1.0, nil) // left
	sNoreact.AddBc(true, 11, 2.0, nil) // right
	sNoreact.AddBc(true, 20, 1.0, nil) // bottom
	sNoreact.AddBc(true, 21, 2.0, nil) // top
	sNoreact.Assemble(reactions)
	uNoreact, fNoreact := sNoreact.SolveSteady(reactions)
	io.Pf("uNoreact = %v\n", uNoreact)
	io.Pf("fNoreact = %v\n", fNoreact)
	if fNoreact != nil {
		tst.Errorf("fNoreact should be nil\n")
		return
	}
	chk.Array(tst, "uNoreact", 1e-15, uNoreact, u)

	// plot
	if chk.Verbose {
		gp := gm.GridPlotter{G: g, WithVids: true}
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		gp.Draw()
		plt.ContourL(gp.X2d, gp.Y2d, uu, nil)
		plt.Gll("$x$", "$y$", nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/pde", "fdm02")
	}
}
