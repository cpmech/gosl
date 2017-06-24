// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la/oblas"
)

// MatAlloc allocates a matrix with size m,n:
//  a := 0  =>  aij = 0
func MatAlloc(m, n int) (mat [][]float64) {
	mat = make([][]float64, m)
	for k := 0; k < m; k++ {
		mat[k] = make([]float64, n)
	}
	return mat
}

// MatClone allocates and clones a matrix
//  b := a
func MatClone(a [][]float64) (b [][]float64) {
	b = make([][]float64, len(a))
	for i := 0; i < len(a); i++ {
		b[i] = make([]float64, len(a[i]))
		for j := 0; j < len(a[i]); j++ {
			b[i][j] = a[i][j]
		}
	}
	return
}

// MatFill fills a matrix with a single number s:
//  aij := s
func MatFill(a [][]float64, s float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			a[i][j] = s
		}
	}
}

// MatScale scales a matrix by a scalar s:
//  a := α * a  =>  aij := α * aij
func MatScale(a [][]float64, α float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			a[i][j] *= α
		}
	}
}

// MatCopy copies to matrix "a", another matrix "b" (scaled):
//  a := α * b  =>  aij := α * bij
func MatCopy(a [][]float64, α float64, b [][]float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			a[i][j] = α * b[i][j]
		}
	}
}

// MatSetDiag sets a diagonal matrix where the diagonal components are equal to s
func MatSetDiag(a [][]float64, s float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			a[i][j] = 0.0
			if i == j {
				a[i][j] = s
			}
		}
	}
}

// MatMaxDiff returns the maximum difference between components of two matrices
func MatMaxDiff(a, b [][]float64) (maxdiff float64) {
	maxdiff = math.Abs(a[0][0] - b[0][0])
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			diff := math.Abs(a[i][j] - b[i][j])
			if diff > maxdiff {
				maxdiff = diff
			}
		}
	}
	return
}

// MatLargest returns the largest component (abs(a_ij)) of a matrix, normalised by den
func MatLargest(a [][]float64, den float64) (largest float64) {
	largest = math.Abs(a[0][0])
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			tmp := math.Abs(a[i][j])
			if tmp > largest {
				largest = tmp
			}
		}
	}
	return
}

// MatGetCol returns a column of matrix a into vector col:
//  col := a[:][j]
func MatGetCol(j int, a [][]float64) (col []float64) {
	col = make([]float64, len(a))
	for i := 0; i < len(a); i++ {
		col[i] = a[i][j]
	}
	return
}

// MatNormF returns the Frobenious norm of a matrix:
//  ||A||_F := sqrt(sum_i sum_j aij*aij)
func MatNormF(a [][]float64) (res float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			res += a[i][j] * a[i][j]
		}
	}
	return math.Sqrt(res)
}

// MatNormI returns the infinite norm of a matrix:
//  ||A||_∞ := max_i sum_j aij
func MatNormI(a [][]float64) (res float64) {
	if len(a) < 1 {
		return
	}
	for j := 0; j < len(a[0]); j++ {
		res += math.Abs(a[0][j])
	}
	var sumrow float64
	for i := 0; i < len(a); i++ {
		sumrow = 0.0
		for j := 0; j < len(a[i]); j++ {
			sumrow += math.Abs(a[i][j])
			if sumrow > res {
				res = sumrow
			}
		}
	}
	return
}

// --------------------------------------------------------------------------------------------------
// matrix-vector / matrix-matrix --------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------

// MatVecMul returns the matrix-vector multiplication (scaled):
//  v := α * a * u  =>  vi = α * aij * uj
//  NOTE: not efficient implementation => use for small matrices
func MatVecMul(v []float64, α float64, a [][]float64, u []float64) {
	for i := 0; i < len(a); i++ {
		v[i] = 0.0
		for j := 0; j < len(u); j++ {
			v[i] += α * a[i][j] * u[j]
		}
	}
}

// MatVecMulAdd returns the matrix-vector multiplication with addition (scaled):
//  v += α * a * u  =>  vi += α * aij * uj
//  NOTE: not efficient implementation => use for small matrices
func MatVecMulAdd(v []float64, α float64, a [][]float64, u []float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(u); j++ {
			v[i] += α * a[i][j] * u[j]
		}
	}
}

// MatVecMulAddC (complex) returns the matrix-vector multiplication with addition (scaled):
//  v += α * a * u  =>  vi += α * aij * uj
//  NOTE: not efficient implementation => use for small matrices
func MatVecMulAddC(v []complex128, α complex128, a [][]complex128, u []complex128) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(u); j++ {
			v[i] += α * a[i][j] * u[j]
		}
	}
}

// MatTrVecMult returns the matrix-vector multiplication with "a" transposed (scaled):
//  v := α * transp(a) * u  =>  vi = α * aji * uj
//  NOTE: not efficient implementation => use for small matrices
func MatTrVecMul(v []float64, α float64, a [][]float64, u []float64) {
	for i := 0; i < len(a[0]); i++ {
		v[i] = 0.0
		for j := 0; j < len(u); j++ {
			v[i] += α * a[j][i] * u[j]
		}
	}
}

// MatTrVecMulAdd returns the matrix-vector multiplication with addition and "a" transposed (scaled):
//  v += α * transp(a) * u  =>  vi += α * aji * uj
//  NOTE: not efficient implementation => use for small matrices
func MatTrVecMulAdd(v []float64, α float64, a [][]float64, u []float64) {
	for i := 0; i < len(a[0]); i++ {
		for j := 0; j < len(u); j++ {
			v[i] += α * a[j][i] * u[j]
		}
	}
}

// MatVecMulCopyAdd returns the matrix-vector multiplication with copy and addition (scaled):
//  w = α*v + β*a*u  =>  wi = α*vi + β*aij*uj
//  NOTE: not efficient implementation => use for small matrices
func MatVecMulCopyAdd(w []float64, α float64, v []float64, β float64, a [][]float64, u []float64) {
	for i := 0; i < len(a); i++ {
		w[i] = α * v[i]
		for j := 0; j < len(u); j++ {
			w[i] += β * a[i][j] * u[j]
		}
	}
}

// MatMul returns the matrix multiplication (scaled):
//  c := α * a * b  =>  cij := α * aik * bkj
//  NOTE: not efficient implementation => use for small matrices
func MatMul(c [][]float64, α float64, a, b [][]float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b[0]); j++ {
			c[i][j] = 0.0
			for k := 0; k < len(a[0]); k++ {
				c[i][j] += α * a[i][k] * b[k][j]
			}
		}
	}
}

// MatMulNew returns the matrix multiplication (scaled):
//  c := α * a * b  =>  cij := α * aik * bkj
//  NOTE: not efficient implementation => use for small matrices
func MatMulNew(c *Matrix, α float64, a, b *Matrix) {
	for i := 0; i < c.M*c.N; i++ {
		c.Data[i] = 0
	}
	err := oblas.Dgemm(false, false, a.M, b.N, a.N, α, a.Data, a.M, b.Data, b.M, 1, c.Data, c.M)
	if err != nil {
		chk.Panic("%v\n", err)
	}
}

// MatMul3 returns the triple matrix multiplication:
//  d := α * a * b * c  =>  dij := α * aik * bkl * clj
//  NOTE: not efficient implementation => use for small matrices
func MatMul3(d [][]float64, α float64, a, b, c [][]float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(c[0]); j++ {
			d[i][j] = 0.0
			for k := 0; k < len(a[0]); k++ {
				for l := 0; l < len(b[0]); l++ {
					d[i][j] += α * a[i][k] * b[k][l] * c[l][j]
				}
			}
		}
	}
}

// MatTrMul3 returns the triple matrix multiplication with tranposed "a":
//  d := α * trans(a) * b * c  =>  dij := α * aki * bkl * clj
//  NOTE: not efficient implementation => use for small matrices
func MatTrMul3(d [][]float64, α float64, a, b, c [][]float64) {
	for i := 0; i < len(a[0]); i++ {
		for j := 0; j < len(c[0]); j++ {
			d[i][j] = 0.0
			for k := 0; k < len(a); k++ {
				for l := 0; l < len(b[0]); l++ {
					d[i][j] += α * a[k][i] * b[k][l] * c[l][j]
				}
			}
		}
	}
}

// MatTrMulAdd3 returns the triple matrix multiplication with addition and tranposed "a":
//  d += α * trans(a) * b * c  =>  dij += α * aki * bkl * clj
//  NOTE: not efficient implementation => use for small matrices
func MatTrMulAdd3(d [][]float64, α float64, a, b, c [][]float64) {
	for i := 0; i < len(a[0]); i++ {
		for j := 0; j < len(c[0]); j++ {
			for k := 0; k < len(a); k++ {
				for l := 0; l < len(b[0]); l++ {
					d[i][j] += α * a[k][i] * b[k][l] * c[l][j]
				}
			}
		}
	}
}

// VecOuterAdd returns the outer product between two vectors, with addition (scaled)
//  aij += α * u[i] * v[j]
//  NOTE: not efficient implementation => use for small matrices
func VecOuterAdd(a [][]float64, α float64, u, v []float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			a[i][j] += α * u[i] * v[j]
		}
	}
}
