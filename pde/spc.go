// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/la"
)

// SpcLaplacian implements the Spectral Collocation (SPC) Laplacian operator (2D or 3D)
//
//                 ∂²φ        ∂²φ        ∂²φ         ∂φ      ∂φ
//    L{φ} = ∇²φ = ——— g¹¹ +  ——— g²² +  ———— 2g¹² - —— L¹ - —— L²
//                 ∂a²        ∂b²        ∂a∂b        ∂a      ∂b
//
//    with a=u[0]=r and b=u[1]=s
//
type SpcLaplacian struct {
	LagInt   fun.LagIntSet  // Lagrange interpolators [ndim]
	Grid     *gm.Grid       // grid
	Source   fun.Svs        // source term function s({x},t)
	EssenBcs *BoundaryConds // essential boundary conditions
	NaturBcs *BoundaryConds // natural boundary conditions
	Eqs      *la.Equations  // equations
	bcsReady bool           // boundary conditions are set
}

// NewSpcLaplacian creates a new SPC Laplacian operator with given parameters
//  NOTE: params is not used at the moment
func NewSpcLaplacian(params dbf.Params, lis fun.LagIntSet, grid *gm.Grid, source fun.Svs) (o *SpcLaplacian) {
	o = new(SpcLaplacian)
	o.LagInt = lis
	o.Grid = grid
	o.Source = source
	o.EssenBcs = NewBoundaryCondsGrid(grid, 1) // 1:maxNdof
	o.NaturBcs = NewBoundaryCondsGrid(grid, 1) // 1:maxNdof
	o.bcsReady = false
	return
}

// AddEbc adds essential boundary condition given tag of edge or face
//   tag    -- edge or face tag in grid
//   cvalue -- constant value [optional]; or
//   fvalue -- function value [optional]
func (o *SpcLaplacian) AddEbc(tag int, cvalue float64, fvalue fun.Svs) {
	o.bcsReady = false
	o.EssenBcs.AddUsingTag(tag, 0, cvalue, fvalue)
}

// SetHbc sets homogeneous boundary conditions; i.e. all boundaries with zero EBC
func (o *SpcLaplacian) SetHbc() {
	if o.Grid.Ndim() == 2 {
		o.AddEbc(10, 0.0, nil)
		o.AddEbc(11, 0.0, nil)
		o.AddEbc(20, 0.0, nil)
		o.AddEbc(21, 0.0, nil)
		return
	}
	o.AddEbc(100, 0.0, nil)
	o.AddEbc(101, 0.0, nil)
	o.AddEbc(200, 0.0, nil)
	o.AddEbc(201, 0.0, nil)
	o.AddEbc(300, 0.0, nil)
	o.AddEbc(301, 0.0, nil)
}

// AddNbc adds natural boundary condition given tag of edge or face
//   tag    -- edge or face tag in grid
//   cvalue -- constant value [optional]; or
//   fvalue -- function value [optional]
func (o *SpcLaplacian) AddNbc(tag int, cvalue float64, fvalue fun.Svs) {
	o.bcsReady = false
	o.NaturBcs.AddUsingTag(tag, 0, cvalue, fvalue)
}

// Assemble assembles operator into A matrix from [A] ⋅ {u} = {b}
//  reactions -- prepare for computation of RHS
func (o *SpcLaplacian) Assemble(reactions bool) {
	nx := o.LagInt[0].N + 1
	ny := o.LagInt[1].N + 1
	if !o.bcsReady {
		nnz := (nx * nx) * (ny * ny)
		o.Eqs = la.NewEquations(o.Grid.Size(), o.EssenBcs.Nodes())
		o.Eqs.Alloc([]int{nnz, nnz, nnz, nnz}, reactions, true) // TODO: optimise nnz
		for _, li := range o.LagInt {
			li.CalcD2() // also calculates D1
		}
		o.bcsReady = true
	}
	ι := func(m, n int) int { return o.Grid.IndexMNPtoI(m, n, 0) }
	δ := func(m, n int) float64 {
		if m == n {
			return 1
		}
		return 0
	}
	g := func(i, j, m, n int) float64 { return o.Grid.ContraMatrix(m, n, 0).Get(i, j) }
	L := func(i, m, n int) float64 { return o.Grid.Lcoeff(m, n, 0, i) }
	D1a := func(m, n int) float64 { return o.LagInt[0].D1.Get(m, n) }
	D1b := func(m, n int) float64 { return o.LagInt[1].D1.Get(m, n) }
	D2a := func(m, n int) float64 { return o.LagInt[0].D2.Get(m, n) }
	D2b := func(m, n int) float64 { return o.LagInt[1].D2.Get(m, n) }
	N := la.NewVector(o.Grid.Ndim())
	o.Eqs.Start()
	if o.Grid.Ndim() == 2 {
		for p := 0; p < nx; p++ {
			for q := 0; q < ny; q++ {
				I := ι(p, q)
				if o.NaturBcs.Has(I) {
					o.NaturBcs.NormalGrid(N, I)
					α := la.VecDot(N, o.Grid.ContraBasis(p, q, 0, 0)) // n⋅g^0
					β := la.VecDot(N, o.Grid.ContraBasis(p, q, 0, 1)) // n⋅g^1
					for m := 0; m < nx; m++ {
						for n := 0; n < ny; n++ {
							J := ι(m, n)
							o.Eqs.Put(I, J, D1a(p, m)*δ(q, n)*α+δ(p, m)*D1b(q, n)*β)
						}
					}
				} else {
					for m := 0; m < nx; m++ {
						for n := 0; n < ny; n++ {
							J := ι(m, n)
							o.Eqs.Put(I, J, 0+
								D2a(p, m)*δ(q, n)*g(0, 0, p, q)+
								δ(p, m)*D2b(q, n)*g(1, 1, p, q)+
								D1a(p, m)*D1b(q, n)*2.0*g(0, 1, p, q)+
								-D1a(p, m)*δ(q, n)*L(0, p, q)+
								-δ(p, m)*D1b(q, n)*L(1, p, q))
						}
					}
				}
			}
		}
		return
	}
	chk.Panic("TODO: Implement Assemble() in 3D\n")
}

// SolveSteady solves steady problem
//   Solves: [K]⋅{u} = {f} represented by [A]⋅{x} = {b}
func (o *SpcLaplacian) SolveSteady(reactions bool) (u, f []float64) {
	o.Eqs.SolveOnce(o.calcXk, o.calcBu)
	u = make([]float64, o.Grid.Size())
	o.Eqs.JoinVector(u, o.Eqs.Xu, o.Eqs.Xk)
	if reactions {
		f = make([]float64, o.Grid.Size())
		if o.Eqs.Nk > 0 { // need to calc Bu again because it was modified
			for i, I := range o.Eqs.UtoF {
				o.Eqs.Bu[i] = o.calcBu(I, 0)
			}
		}
		o.Eqs.JoinVector(f, o.Eqs.Bu, o.Eqs.Bk)
	}
	return
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// calcXk calculates known {u} values (CalcXk in la.Equations)
//  I -- node number
//  t -- time
func (o *SpcLaplacian) calcXk(I int, t float64) float64 {
	_, val, available := o.EssenBcs.Value(I, 0, t)
	if available {
		return val
	}
	return 0
}

// calcBu calculates RHS vector (e.g. source) corresponding to known values of {u} (CalcBu in la.Equations)
//  I -- node number
//  t -- time
func (o *SpcLaplacian) calcBu(I int, t float64) float64 {
	if o.Source != nil {
		return o.Source(o.Grid.Node(I), t)
	}
	return 0
}
