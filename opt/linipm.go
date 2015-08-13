// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package opt implements routines for solving optimisation problems
package opt

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

// LinIpm implements the interior-point mehtods for linear programming problems
//  Solve:
//          min cᵀx   s.t.   Aᵀx = b, x ≥ 0
//           x
//
//  or the dual problem:
//
//          max bᵀλ   s.t.   Aᵀλ + s = c, s ≥ 0
//           λ
type LinIpm struct {

	// problem
	A *la.Triplet // [Nl][nx]
	B []float64   // [Nl]
	C []float64   // [nx]

	// constants
	NmaxIt int // max number of iterations

	// dimensions
	Nx int // number of x
	Nl int // number of λ
	Ny int // number of y = nx + ns + nl = 2 * nx + nl

	// solution vector
	Y   []float64 // y := [x, λ, s]
	X   []float64 // subset of y
	L   []float64 // subset of y
	S   []float64 // subset of y
	Mdy []float64 // -Δy

	// affine solution
	R  []float64   // residual
	Rx []float64   // subset of R
	Rl []float64   // subset of R
	Rs []float64   // subset of R
	J  *la.Triplet // [ny][ny] Jacobian matrix

	// linear solver
	Lis la.LinSol // linear solver
}

// Clean cleans allocated memory
func (o *LinIpm) Clean() {
	o.Lis.Clean()
}

// Init initialises LinIpm
func (o *LinIpm) Init(A *la.Triplet, b, c []float64, prms fun.Prms) {

	// problem
	o.A, o.B, o.C = A, b, c

	// constants
	o.NmaxIt = 1
	for _, p := range prms {
		switch p.N {
		case "nmaxit":
			o.NmaxIt = int(p.V)
		}
	}

	// dimensions
	o.Nx = len(o.C)
	o.Nl = len(o.B)
	o.Ny = 2*o.Nx + o.Nl
	ix, jx := 0, o.Nx
	il, jl := o.Nx, o.Nx+o.Nl
	is, js := o.Nx+o.Nl, o.Ny

	// solution vector
	o.Y = make([]float64, o.Ny)
	o.X = o.Y[ix:jx]
	o.L = o.Y[il:jl]
	o.S = o.Y[is:js]
	o.Mdy = make([]float64, o.Ny)

	// affine solution
	o.R = make([]float64, o.Ny)
	o.Rx = o.R[ix:jx]
	o.Rl = o.R[il:jl]
	o.Rs = o.R[is:js]
	o.J = new(la.Triplet)
	nnz := 2*o.Nl*o.Nx + 3*o.Nx
	o.J.Init(o.Ny, o.Ny, nnz)

	// linear solver
	o.Lis = la.GetSolver("umfpack")
}

// Solve solves linear programming problem
func (o *LinIpm) Solve() (err error) {

	// constants for linear solver
	symmetric := false
	verbose := false
	timing := false

	// auxiliary
	ix, jx := 0, o.Nx
	is, js := o.Nx+o.Nl, o.Ny

	// control variables
	var μ, σ float64  // μ and σ
	var xrmin float64 // min{ x_i / (-Δx_i) } (x_ratio_minimum)
	var srmin float64 // min{ s_i / (-Δs_i) } (s_ratio_minimum)
	var αpa float64   // α_prime_affine
	var αda float64   // α_dual_affine
	var μaff float64  // μ_affine

	// perform iterations
	it := 0
	for it = 0; it < o.NmaxIt; it++ {

		// compute residual
		la.SpTriMatTrVecMul(o.Rx, o.A, o.L) // rx := Aᵀλ
		la.SpTriMatVecMul(o.Rl, o.A, o.X)   // rλ := A x
		for i := 0; i < o.Nx; i++ {
			o.Rx[i] += o.S[i] - o.C[i]
			o.Rs[i] = o.X[i] * o.S[i]
		}
		for i := 0; i < o.Nl; i++ {
			o.Rl[i] -= o.B[i]
		}

		// assemble Jacobian
		o.J.Start()
		o.J.PutMatAndMatT(o.A)
		for i := 0; i < o.Nx; i++ {
			o.J.Put(i, is+i, 1.0)
			o.J.Put(is+i, i, o.S[i])
			o.J.Put(is+i, is+i, o.X[i])
		}

		// solve linear system
		if it == 0 {
			err = o.Lis.InitR(o.J, symmetric, verbose, timing)
			if err != nil {
				return
			}
		}
		err = o.Lis.Fact()
		if err != nil {
			return
		}
		err = o.Lis.SolveR(o.Mdy, o.R, false) // mdy := inv(J) * R
		if err != nil {
			return
		}

		// control variables
		mdx := o.Mdy[ix:jx] // -Δx
		mds := o.Mdy[is:js] // -Δs
		xrmin, srmin, μ = 0, 0, 0
		for i := 0; i < o.Nx; i++ {
			if mdx[i] > 0 {
				xrmin = min(xrmin, o.X[i]/mdx[i])
			}
			if mds[i] > 0 {
				srmin = min(srmin, o.S[i]/mds[i])
			}
			μ += o.X[i] * o.S[i]
		}
		μ /= float64(o.Nx)

		// compute σ
		αpa = min(1, xrmin)
		αda = min(1, srmin)
		μaff = 0
		for i := 0; i < o.Nx; i++ {
			μaff += (o.X[i] - αpa*mdx[i]) * (o.S[i] - αda*mds[i])
		}
		σ = math.Pow(μaff/μ, 3)
		io.Pforan("σ = %v\n", σ)
	}

	// check convergence
	if it == o.NmaxIt {
		err = chk.Err("iterations dit not converge")
	}
	return
}
