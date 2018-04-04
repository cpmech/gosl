// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/la"
)

// Problem holds the functions defining an optimization problem
type Problem struct {
	Ndim int       // dimension of x == len(x)
	Ffcn fun.Sv    // objective function f({x})
	Gfcn fun.Vv    // gradient function df/d{x}|(x)
	Hfcn fun.Mv    // Hessian function d²f/d{x}d{x}|(x)
	Fref float64   // known solution fmin = f({x})
	Xref la.Vector // known solution {x} @ min
}

// NewQuadraticProblem returns a quadratic optimization problem such that f(x) = xᵀ A x
func NewQuadraticProblem(Amat [][]float64) (p *Problem) {

	// new problem
	p = new(Problem)
	p.Ndim = len(Amat)

	// objective function f({x})
	A := la.NewMatrixDeep2(Amat)
	tmp := la.NewVector(A.M)
	p.Ffcn = func(x la.Vector) float64 {
		la.MatVecMul(tmp, 1, A, x)
		return la.VecDot(x, tmp) // xᵀ A x
	}

	// gradient function df/d{x}|(x)
	At := A.GetTranspose()
	AtPlusA := la.NewMatrix(A.M, A.M)
	la.MatAdd(AtPlusA, 1, At, 1, A)
	p.Gfcn = func(g, x la.Vector) {
		la.MatVecMul(g, 1, AtPlusA, x) // g := (Aᵀ+A)⋅x
	}

	// Hessian function d²f/d{x}d{x}!(x)
	p.Hfcn = func(h *la.Matrix, x la.Vector) {
		AtPlusA.CopyInto(h, 1) // H := Aᵀ + A
	}

	// solution
	p.Fref = 0.0
	p.Xref = la.NewVector(A.M) // xmin := 0.0
	return
}
