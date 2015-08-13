// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import "github.com/cpmech/gosl/chk"

// --------------------------------------------------------------------------------------------------
// matrix-vector ------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------

// SpMatVecMul returns the (sparse) matrix-vector multiplication (scaled):
//  v := α * a * u  =>  vi = α * aij * uj
//  NOTE: dense vector v will be first initialised with zeros
func SpMatVecMul(v []float64, α float64, a *CCMatrix, u []float64) {
	VecFill(v, 0)
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			v[a.i[k]] += α * a.x[k] * u[j]
		}
	}
}

// SpMatVecMulAdd returns the (sparse) matrix-vector multiplication with addition (scaled):
//  v += α * a * u  =>  vi += α * aij * uj
func SpMatVecMulAdd(v []float64, α float64, a *CCMatrix, u []float64) {
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			v[a.i[k]] += α * a.x[k] * u[j]
		}
	}
}

// SpMatVecMulAddX returns the (sparse) matrix-vector multiplication with addition (scaled/extended):
//  v += a * (α*u + β*w)  =>  vi += aij * (α*uj + β*wj)
func SpMatVecMulAddX(v []float64, a *CCMatrix, α float64, u []float64, β float64, w []float64) {
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			v[a.i[k]] += a.x[k] * (α*u[j] + β*w[j])
		}
	}
}

// SpMatTrVecMult returns the (sparse) matrix-vector multiplication with "a" transposed (scaled):
//  v := α * transp(a) * u  =>  vj = α * aij * ui
//  NOTE: dense vector v will be first initialised with zeros
func SpMatTrVecMul(v []float64, α float64, a *CCMatrix, u []float64) {
	VecFill(v, 0)
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			v[j] += α * a.x[k] * u[a.i[k]]
		}
	}
}

// SpMatTrVecMulAdd returns the (sparse) matrix-vector multiplication with addition and "a" transposed (scaled):
//  v += α * transp(a) * u  =>  vj += α * aij * ui
func SpMatTrVecMulAdd(v []float64, α float64, a *CCMatrix, u []float64) {
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			v[j] += α * a.x[k] * u[a.i[k]]
		}
	}
}

// SpMatVecMulC returns the (sparse/complex) matrix-vector multiplication (scaled):
//  v := α * a * u  =>  vi = α * aij * uj
//  NOTE: dense vector v will be first initialised with zeros
func SpMatVecMulC(v []complex128, α float64, a *CCMatrixC, u []complex128) {
	VecFillC(v, 0)
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			v[a.i[k]] += complex(α*a.x[k], α*a.z[k]) * u[j]
		}
	}
}

// SpMatVecMulAddC returns the (sparse/complex) matrix-vector multiplication with addition (scaled):
//  v += α * a * u  =>  vi += α * aij * uj
func SpMatVecMulAddC(v []complex128, α float64, a *CCMatrixC, u []complex128) {
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			v[a.i[k]] += complex(α*a.x[k], α*a.z[k]) * u[j]
		}
	}
}

// SpMatTrVecMultC returns the (sparse/complex) matrix-vector multiplication with "a" transposed (scaled):
//  v := α * transp(a) * u  =>  vj = α * aij * ui
//  NOTE: dense vector v will be first initialised with zeros
func SpMatTrVecMulC(v []complex128, α float64, a *CCMatrixC, u []complex128) {
	VecFillC(v, 0)
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			v[j] += complex(α*a.x[k], α*a.z[k]) * u[a.i[k]]
		}
	}
}

// SpMatTrVecMulAddC returns the (sparse/complex) matrix-vector multiplication with addition and "a" transposed (scaled):
//  v += α * transp(a) * u  =>  vj += α * aij * ui
func SpMatTrVecMulAddC(v []complex128, α float64, a *CCMatrixC, u []complex128) {
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			v[j] += complex(α*a.x[k], α*a.z[k]) * u[a.i[k]]
		}
	}
}

// --------------------------------------------------------------------------------------------------
// auxiliary ----------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------

// SpInitSimilar initialises another matrix "b" with the same structure (Ap, Ai) of
// sparse matrix "a". The values Ax are not copied though.
func SpInitSimilar(b *CCMatrix, a *CCMatrix) {
	b.m, b.n, b.nnz = a.m, a.n, a.nnz
	b.p = make([]int, a.n+1)
	b.i = make([]int, a.nnz)
	b.x = make([]float64, a.nnz)
	for j := 0; j < a.n+1; j++ {
		b.p[j] = a.p[j]
	}
	for k := 0; k < a.nnz; k++ {
		b.i[k] = a.i[k]
	}
}

// SpInitSimilarR2C initialises another matrix "b" (complex) with the same structure (Ap, Ai) of
// sparse matrix "a" (real). The values Ax are not copied though (Bx and Bz are not set).
func SpInitSimilarR2C(b *CCMatrixC, a *CCMatrix) {
	b.m, b.n, b.nnz = a.m, a.n, a.nnz
	b.p = make([]int, a.n+1)
	b.i = make([]int, a.nnz)
	b.x = make([]float64, a.nnz)
	b.z = make([]float64, a.nnz)
	for j := 0; j < a.n+1; j++ {
		b.p[j] = a.p[j]
	}
	for k := 0; k < a.nnz; k++ {
		b.i[k] = a.i[k]
	}
}

// SpMatAddI adds an identity matrix I to "a", scaled by α and β according to:
//  r := α*a + β*I
func SpMatAddI(r *CCMatrix, α float64, a *CCMatrix, β float64) {
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			if a.i[k] == j {
				r.x[k] = α*a.x[k] + β
			} else {
				r.x[k] = α * a.x[k]
			}
		}
	}
}

// SpCheckDiag checks if all elements on the diagonal of "a" are present.
//  OUTPUT:
//   ok -- true if all diagonal elements are present;
//         otherwise, ok == false if any diagonal element is missing.
func SpCheckDiag(a *CCMatrix) bool {
	rowok := make([]bool, a.m)
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			if a.i[k] == j {
				rowok[j] = true
			}
		}
	}
	for i := 0; i < a.m; i++ {
		if !rowok[i] {
			return false
		}
	}
	return true
}

// SpInitRc initialises two complex sparse matrices (residual correction) according to:
//  Real:     γ      *I - J
//  Complex: (α + βi)*I - J
//  NOTE: "a" must include all diagonal elements
func SpInitRc(rcmatR *CCMatrix, rcmatC *CCMatrixC, a *CCMatrix) {
	rcmatR.m, rcmatR.n, rcmatR.nnz = a.m, a.n, a.nnz
	rcmatR.p = make([]int, a.n+1)
	rcmatR.i = make([]int, a.nnz)
	rcmatR.x = make([]float64, a.nnz)
	rcmatC.m, rcmatC.n, rcmatC.nnz = a.m, a.n, a.nnz
	rcmatC.p = make([]int, a.n+1)
	rcmatC.i = make([]int, a.nnz)
	rcmatC.x = make([]float64, a.nnz)
	rcmatC.z = make([]float64, a.nnz)
	for j := 0; j < a.n+1; j++ {
		rcmatR.p[j] = a.p[j]
		rcmatC.p[j] = a.p[j]
	}
	for k := 0; k < a.nnz; k++ {
		rcmatR.i[k] = a.i[k]
		rcmatC.i[k] = a.i[k]
	}
}

// SpSetRc sets the values within two complex sparse matrices (residual correction) according to:
//  Real:     γ      *I - J
//  Complex: (α + βi)*I - J
//  NOTE: "a" must include all diagonal elements
func SpSetRc(rcmatR *CCMatrix, rcmatC *CCMatrixC, α, β, γ float64, a *CCMatrix) {
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			if a.i[k] == j {
				rcmatR.x[k] = γ - a.x[k]
				rcmatC.x[k] = α - a.x[k]
				rcmatC.z[k] = β
			} else {
				rcmatR.x[k] = -a.x[k]
				rcmatC.x[k] = -a.x[k]
				rcmatC.z[k] = 0.0
			}
		}
	}
}

// SpTriSetDiag sets a (n x n) real triplet with diagonal values 'v'
func SpTriSetDiag(a *Triplet, n int, v float64) {
	a.Init(n, n, n)
	for i := 0; i < n; i++ {
		a.Put(i, i, v)
	}
}

// --------------------------------------------------------------------------------------------------
// matrix-matrix ------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------

// SpAllocMatAddMat allocates a matrix 'c' to hold the result of the addition of 'a' and 'b'.
// It also allocates the mapping arrays a2c and b2c, where:
//  a2c maps 'k' in 'a' to 'k' in 'c': len(a2c) = a.nnz
//  b2c maps 'k' in 'b' to 'k' in 'c': len(b2c) = b.nnz
func SpAllocMatAddMat(a, b *CCMatrix) (c *CCMatrix, a2c, b2c []int) {
	if a.m != b.m || a.n != b.n {
		chk.Panic("matrices 'a' (%dx%d) and 'b' (%dx%d) must have the same dimensions", a.m, a.n, b.m, b.n)
	}
	// number of nonzeros in 'c'
	var i, j, k, nnz int
	r2a := make([]int, a.m) // maps a row index to the corresponding k index of 'a'
	r2b := make([]int, a.m) // maps a row index to the corresponding k index of 'b'
	exact := true
	if exact {
		for j = 0; j < a.n; j++ {
			for i = 0; i < a.m; i++ {
				r2a[i], r2b[i] = -1, -1
			}
			for k = a.p[j]; k < a.p[j+1]; k++ {
				r2a[a.i[k]] = k
			}
			for k = b.p[j]; k < b.p[j+1]; k++ {
				r2b[b.i[k]] = k
			}
			for i = 0; i < a.m; i++ {
				if r2a[i] > -1 || r2b[i] > -1 {
					nnz += 1
				}
			}
		}
	} else {
		nnz = a.nnz + b.nnz
	}
	// allocate c, a2c, and b2c
	c = new(CCMatrix)
	c.m, c.n, c.nnz = a.m, a.n, nnz
	c.x = make([]float64, nnz)
	c.i = make([]int, nnz)
	c.p = make([]int, c.n+1)
	a2c = make([]int, a.nnz)
	b2c = make([]int, b.nnz)
	nnz = 0 // == k of 'c'
	for j = 0; j < a.n; j++ {
		for i = 0; i < a.m; i++ {
			r2a[i], r2b[i] = -1, -1
		}
		for k = a.p[j]; k < a.p[j+1]; k++ {
			r2a[a.i[k]] = k
		}
		for k = b.p[j]; k < b.p[j+1]; k++ {
			r2b[b.i[k]] = k
		}
		for i = 0; i < a.m; i++ {
			if r2a[i] > -1 || r2b[i] > -1 {
				if r2a[i] > -1 {
					a2c[r2a[i]] = nnz
				}
				if r2b[i] > -1 {
					b2c[r2b[i]] = nnz
				}
				c.i[nnz] = i
				nnz += 1
			}
		}
		c.p[j+1] = nnz
	}
	return
}

// SpMatAddMat adds two sparse matrices. The 'c' matrix matrix and the 'a2c' and 'b2c' arrays
// must be pre-allocated by SpAllocMatAddMat. The result is:
//  c := α*a + β*b
//  NOTE: this routine does not check for the correct sizes, since this is expect to be
//        done by SpAllocMatAddMat
func SpMatAddMat(c *CCMatrix, α float64, a *CCMatrix, β float64, b *CCMatrix, a2c, b2c []int) {
	VecFill(c.x, 0)
	for k := 0; k < a.nnz; k++ {
		c.x[a2c[k]] += α * a.x[k]
	}
	for k := 0; k < b.nnz; k++ {
		c.x[b2c[k]] += β * b.x[k]
	}
}

// SpMatAddMatC adds two real sparse matrices with two sets of coefficients in such a way that
// one real matrix (cR) and another complex matrix (cC) are obtained. The results are:
//    cR :=      γ*a + μ*b
//    cC := (α+βi)*a + μ*b
//  NOTE: the structure of cR and cC are the same and can be allocated with SpAllocMatAddMat,
//        followed by one call to SpInitSimilarR2C. For example:
//            cR, a2c, b2c := SpAllocMatAddMat(a, b)
//            SpInitSimilarR2C(cC, cR)
func SpMatAddMatC(cC *CCMatrixC, cR *CCMatrix, α, β, γ float64, a *CCMatrix, μ float64, b *CCMatrix, a2c, b2c []int) {
	for k := 0; k < cR.nnz; k++ {
		cR.x[k], cC.x[k], cC.z[k] = 0, 0, 0
	}
	for k := 0; k < a.nnz; k++ {
		cR.x[a2c[k]] += γ * a.x[k]
		cC.x[a2c[k]] += α * a.x[k]
		cC.z[a2c[k]] += β * a.x[k]
	}
	for k := 0; k < b.nnz; k++ {
		cR.x[b2c[k]] += μ * b.x[k]
		cC.x[b2c[k]] += μ * b.x[k]
	}
}

// SpTriAdd adds two matrices in Triplet format:
//   c := α*a + β*b
//   NOTE: the output 'c' triplet must be able to hold all nonzeros of 'a' and 'b'
//         actually the 'c' triplet is just expanded
func SpTriAdd(c *Triplet, α float64, a *Triplet, β float64, b *Triplet) {
	c.Start()
	for k := 0; k < a.pos; k++ {
		c.Put(a.i[k], a.j[k], α*a.x[k])
	}
	for k := 0; k < b.pos; k++ {
		c.Put(b.i[k], b.j[k], β*b.x[k])
	}
}

// SpTriAddR2C adds two real matrices in Triplet format generating a complex triplet
// accordint to:
//   c := (α+βi)*a + μ*b
//   NOTE: the output 'c' triplet must be able to hold all nonzeros of 'a' and 'b'
//         actually the 'c' triplet is just expanded
func SpTriAddR2C(c *TripletC, α, β float64, a *Triplet, μ float64, b *Triplet) {
	c.Start()
	for k := 0; k < a.pos; k++ {
		c.Put(a.i[k], a.j[k], α*a.x[k], β*a.x[k])
	}
	for k := 0; k < b.pos; k++ {
		c.Put(b.i[k], b.j[k], μ*b.x[k], 0)
	}
}

// SpTriMatVecMul returns the matrix-vector multiplication with matrix a in
// triplet format and two dense vectors x and y
//  y := a * x    or    y_i := a_ij * x_j
func SpTriMatVecMul(y []float64, a *Triplet, x []float64) {
	if len(y) != a.m {
		chk.Panic("length of vector y must be equal to %d. y_(%d × 1). a_(%d × %d)", a.m, len(y), a.m, a.n)
	}
	if len(x) != a.n {
		chk.Panic("length of vector x must be equal to %d. x_(%d × 1). a_(%d × %d)", a.n, len(x), a.m, a.n)
	}
	for i := 0; i < len(y); i++ {
		y[i] = 0
	}
	for k := 0; k < a.pos; k++ {
		y[a.i[k]] += a.x[k] * x[a.j[k]]
	}
}

// SpTriMatTrVecMul returns the matrix-vector multiplication with transposed matrix a in
// triplet format and two dense vectors x and y
//  y := transpose(a) * x    or    y_I := a_JI * x_J    or     y_j := a_ij * x_i
func SpTriMatTrVecMul(y []float64, a *Triplet, x []float64) {
	if len(y) != a.n {
		chk.Panic("length of vector y must be equal to %d. y_(%d × 1). a_(%d × %d)", a.n, len(y), a.m, a.n)
	}
	if len(x) != a.m {
		chk.Panic("length of vector x must be equal to %d. x_(%d × 1). a_(%d × %d)", a.m, len(x), a.m, a.n)
	}
	for j := 0; j < len(y); j++ {
		y[j] = 0
	}
	for k := 0; k < a.pos; k++ {
		y[a.j[k]] += a.x[k] * x[a.i[k]]
	}
}
