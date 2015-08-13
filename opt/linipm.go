// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package opt implements routines for solving optimisation problems
package opt

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

// LinIPM implements the interior-point mehtods for linear programming problems
//  Solve:
//          min cᵀx   s.t.   Aᵀx = b, x ≥ 0
//           x
//
//  or the dual problem:
//
//          max bᵀλ   s.t.   Aᵀλ + s = c, s ≥ 0
//           λ
type LinIPM struct {
	A      *la.Triplet // [nλ][nx]
	b      []float64   // [nλ]
	c      []float64   // [nx]
	R      []float64   // residual
	y      []float64   // auxiliary vector: y := [x, λ, s]
	mdy    []float64   // -Δy
	J      *la.Triplet // [ny][ny] Jacobian matrix
	nx     int         // number of x
	nλ     int         // number of λ
	ny     int         // number of y == nx + ns + nλ == 2*nx + nλ
	ix, jx int         // start and end+1 in R corresponding to rx
	il, jl int         // start and end+1 in R corresponding to rl
	is, js int         // start and end+1 in R corresponding to rs
	x      []float64   // subset of y
	λ      []float64   // subset of y
	s      []float64   // subset of y
	rx     []float64   // subset of y
	rλ     []float64   // subset of y
	rs     []float64   // subset of y
	NmaxIt int         // max number of iterations
}

func (o *LinIPM) Init(A *la.Triplet, b, c []float64) {
	o.A, o.b, o.c = A, b, c
	o.nx = len(o.c)
	o.nλ = len(o.b)
	o.ny = 2*o.nx + o.nλ
	o.R = make([]float64, o.ny)
	o.y = make([]float64, o.ny)
	o.mdy = make([]float64, o.ny)
	o.ix, o.jx = 0, o.nx
	o.il, o.jl = o.nx, o.nx+o.nλ
	o.is, o.js = o.nx+o.nλ, o.ny
	o.J = new(la.Triplet)
	nnz := 2*o.nλ*o.nx + 3*o.nx
	o.J.Init(o.ny, o.ny, nnz)
	o.x = o.y[o.ix:o.jx]
	o.λ = o.y[o.il:o.jl]
	o.s = o.y[o.is:o.js]
	o.rx = o.R[o.ix:o.jx]
	o.rλ = o.R[o.il:o.jl]
	o.rs = o.R[o.is:o.js]
	o.NmaxIt = 1
}

func (o *LinIPM) Solve() (err error) {

	// allocate solver
	lis := la.GetSolver("umfpack")
	defer lis.Clean()
	symmetric := false
	verbose := false
	timing := false

	it := 0
	for it = 0; it < o.NmaxIt; it++ {
		μ := o.CalcMu()
		o.CalcJ()
		if it == 0 {
			err = lis.InitR(o.J, symmetric, verbose, timing)
			if err != nil {
				return
			}
		}
		err = lis.Fact()
		if err != nil {
			return
		}
		err = lis.SolveR(o.mdy, o.R, false) // mdy := inv(J) * R
		if err != nil {
			return
		}
		mdx := o.mdy[o.ix:o.jx] // -Δx
		mds := o.mdy[o.is:o.js] // -Δy
		xdivdx_min := 0.0
		sdivds_min := 0.0
		for i := 0; i < o.nx; i++ {
			if mdx[i] > 0 {
				xdivdx_min = min(xdivdx_min, o.x[i]/mdx[i])
			}
			if mds[i] > 0 {
				sdivds_min = min(sdivds_min, o.s[i]/mds[i])
			}
		}
		αpri := min(1, xdivdx_min)
		αdua := min(1, sdivds_min)
		μaff := 0.0
		for i := 0; i < o.nx; i++ {
			μaff += (o.x[i] - αpri*mdx[i]) * (o.s[i] - αdua*mds[i])
		}
		σ := math.Pow(μaff/μ/float64(o.nx), 3.0)
		io.Pforan("σ = %v\n", σ)
	}

	if it == o.NmaxIt {
		err = chk.Err("iterations dit not converge")
	}
	return
}

func (o *LinIPM) CalcR(y []float64) {
	la.SpTriMatTrVecMul(o.rx, o.A, o.λ) // rx := Aᵀλ
	la.SpTriMatVecMul(o.rλ, o.A, o.x)   // rλ := A x
	for i := 0; i < o.nx; i++ {
		o.rx[i] += o.s[i] - o.c[i]
		o.rs[i] = o.x[i] * o.s[i]
	}
	for i := 0; i < o.nλ; i++ {
		o.rλ[i] -= o.b[i]
	}
}

func (o *LinIPM) CalcMu() (μ float64) {
	for i := 0; i < o.nx; i++ {
		μ += o.x[i] * o.s[i]
	}
	return μ / float64(o.nx)
}

func (o *LinIPM) CalcJ() {
	o.J.Start()
	o.J.PutMatAndMatT(o.A)
	for i := 0; i < o.nx; i++ {
		o.J.Put(i, o.is+i, 1)
		o.J.Put(o.is+i, i, o.s[i])
		o.J.Put(o.is+i, o.is+i, o.x[i])
	}
}
