// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

func TestFdm01a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Fdm01a. Full Auu matrix.")

	// 3x3 grid ⇒ 9 equations
	g := new(gm.Grid)
	g.RectGenUniform([]float64{0, 0}, []float64{2, 2}, []int{3, 3})

	// operator
	s := NewFdmLaplacian(utl.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}, g, nil)

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

func TestFdm01b(tst *testing.T) {
	//verbose()
	chk.PrintTitle("Fdm01b. panic on params")
	defer chk.RecoverTstPanicIsOK(tst)
	NewFdmLaplacian(nil, nil, nil)
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
	p := utl.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}
	s := NewFdmLaplacian(p, g, nil)

	// essential boundary conditions
	s.AddEbc(10, 1.0, nil) // left
	s.AddEbc(11, 2.0, nil) // right
	s.AddEbc(20, 1.0, nil) // bottom
	s.AddEbc(21, 2.0, nil) // top

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
	sNoreact.AddEbc(10, 1.0, nil) // left
	sNoreact.AddEbc(11, 2.0, nil) // right
	sNoreact.AddEbc(20, 1.0, nil) // bottom
	sNoreact.AddEbc(21, 2.0, nil) // top
	sNoreact.Assemble(reactions)
	uNoreact, fNoreact := sNoreact.SolveSteady(reactions)
	io.Pf("uNoreact = %v\n", uNoreact)
	io.Pf("fNoreact = %v\n", fNoreact)
	if fNoreact != nil {
		tst.Errorf("fNoreact should be nil\n")
		return
	}
	chk.Array(tst, "uNoreact", 1e-15, uNoreact, u)
}

func TestFdm03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Fdm03. ZZ p1342 Example 2")

	// 5x5 grid ⇒ 25 equations
	g := new(gm.Grid)
	g.RectGenUniform([]float64{0, 0}, []float64{1, 1}, []int{5, 5})

	// solver
	p := utl.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}
	s := NewFdmLaplacian(p, g, func(X la.Vector, t float64) float64 {
		x, y := X[0], X[1]
		xx, yy := x*x, y*x
		xxx, yyy := xx*x, yy*y
		return 14.0*yyy - (16.0-12.0*x)*yy - (-42.0*xx+54.0*x-2.0)*y + 4.0*xxx - 16.0*xx + 12.0*x
	})

	// homogeneous boundary conditions
	s.SetHbc()

	// set equations and assemble A matrix
	reactions := false
	s.Assemble(reactions)

	// solve
	u, _ := s.SolveSteady(reactions)

	// check
	ana := func(X []float64) float64 {
		x, y := X[0], X[1]
		return x * (1.0 - x) * y * (1.0 - y) * (1.0 + 2.0*x + 7.0*y)
	}
	for n := 0; n < g.Npts(1); n++ {
		for m := 0; m < g.Npts(0); m++ {
			chk.AnaNum(tst, "u", 0.021, u[g.IndexMNPtoI(m, n, 0)], ana(g.X(m, n, 0)), chk.Verbose)
		}
	}
}
