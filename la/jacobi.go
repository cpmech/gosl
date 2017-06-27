// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Jacobi performs the Jacobi transformation of a symmetric matrix to find its eigenvectors and
// eigenvalues.
//
// The Jacobi method consists of a sequence of orthogonal similarity transformations. Each
// transformation (a Jacobi rotation) is just a plane rotation designed to annihilate one of the
// off-diagonal matrix elements. Successive transformations undo previously set zeros, but the
// off-diagonal elements nevertheless get smaller and smaller. Accumulating the product of the
// transformations as you go gives the matrix of eigenvectors (Q), while the elements of the final
// diagonal matrix (A) are the eigenvalues.
//
// The Jacobi method is absolutely foolproof for all real symmetric matrices.
//
//         A = Q ⋅ L ⋅ Qᵀ
//
//   Input:
//    A -- matrix to compute eigenvalues (SYMMETRIC and SQUARE)
//   Output:
//    A -- modified
//    Q -- matrix which columns are the eigenvectors
//    v -- vector with the eigenvalues
//
//   NOTE: for matrices of order greater than about 10, say, the algorithm is slower,
//         by a significant constant factor, than the QR method.
//
func Jacobi(Q *Matrix, v Vector, A *Matrix) (err error) {

	// constants
	tol := 1e-15
	nItMax := 20

	// auxiliary variables
	var j, p, q int
	var θ, τ, t, sm, s, h, g, c float64

	// auxiliary variables
	n := A.M
	b := NewVector(n)
	z := NewVector(n) // this vector will accumulate terms of the form tapq as in equation (11.1.14).

	// initialize Q to the identity matrix.
	for p = 0; p < n; p++ {
		for q = 0; q < n; q++ {
			Q.Set(p, q, 0.0)
		}
		Q.Set(p, p, 1.0)
	}

	// initialize b and v to the diagonal of A
	for p = 0; p < n; p++ {
		b[p] = A.Get(p, p)
		v[p] = A.Get(p, p)
		z[p] = 0.0
	}

	// perform iterations
	for it := 0; it < nItMax; it++ {

		// sum off-diagonal elements.
		sm = 0.0
		for p = 0; p < n-1; p++ {
			for q = p + 1; q < n; q++ {
				sm += math.Abs(A.Get(p, q))
			}
		}

		// exit point
		if sm < tol {
			return
		}

		// rotations
		for p = 0; p < n-1; p++ {
			for q = p + 1; q < n; q++ {
				h = v[q] - v[p]
				if math.Abs(h) <= tol {
					t = 1.0
				} else {
					θ = 0.5 * h / (A.Get(p, q))
					t = 1.0 / (math.Abs(θ) + math.Sqrt(1.0+θ*θ))
					if θ < 0.0 {
						t = -t
					}
				}
				c = 1.0 / math.Sqrt(1.0+t*t)
				s = t * c
				τ = s / (1.0 + c)
				h = t * A.Get(p, q)
				z[p] -= h
				z[q] += h
				v[p] -= h
				v[q] += h
				A.Set(p, q, 0.0)
				for j = 0; j < p; j++ { // case of rotations 0 <= j < p.
					g = A.Get(j, p)
					h = A.Get(j, q)
					A.Set(j, p, g-s*(h+g*τ))
					A.Set(j, q, h+s*(g-h*τ))
				}
				for j = p + 1; j < q; j++ { // case of rotations p < j < q.
					g = A.Get(p, j)
					h = A.Get(j, q)
					A.Set(p, j, g-s*(h+g*τ))
					A.Set(j, q, h+s*(g-h*τ))
				}
				for j = q + 1; j < n; j++ { //case of rotations q < j < n.
					g = A.Get(p, j)
					h = A.Get(q, j)
					A.Set(p, j, g-s*(h+g*τ))
					A.Set(q, j, h+s*(g-h*τ))
				}
				for j = 0; j < n; j++ { // Q matrix
					g = Q.Get(j, p)
					h = Q.Get(j, q)
					Q.Set(j, p, g-s*(h+g*τ))
					Q.Set(j, q, h+s*(g-h*τ))
				}
			}
		}
		for p = 0; p < n; p++ {
			b[p] += z[p]
			z[p] = 0.0  // reinitialize z.
			v[p] = b[p] // update v with the sum of tapq,
		}
	}

	err = chk.Err("Jacobi rotation dit not converge after %d iterations", nItMax+1)
	return
}
