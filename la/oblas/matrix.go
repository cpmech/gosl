// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oblas

// Matrix implements a column-major representation of a matrix by using a linear array that can be passed to Fortran code
//
//  NOTE: the functions related to Matrix do not check for the limits of indices and dimensions. Panic may occur then.
//
//  Example:
//             _      _
//            |  0  3  |
//        M = |  1  4  |
//            |_ 2  5 _|(m x n)
//
//     data[i+j*m] = M[i][j]
//
type Matrix struct {
	m, n int       // dimensions
	data []float64 // data array. column-major => Fortran
}

// NewMatrix allocates a new Matrix
func NewMatrix(m, n int) (o *Matrix) {
	o = new(Matrix)
	o.m, o.n = m, n
	o.data = make([]float64, m*n)
	return
}

// SetFromMat sets matrix with data from a nested slice structure; i.e. form a given la.Mat.
func (o *Matrix) SetFromMat(a [][]float64) {
	k := 0
	for j := 0; j < o.n; j++ {
		for i := 0; i < o.m; i++ {
			o.data[k] = a[i][j]
			k += 1
		}
	}
}

// Sets value
func (o *Matrix) Set(i, j int, val float64) {
	o.data[i+j*o.m] = val // col-major
}

// GetMat returns nested slice representation; i.e. a la.Mat structure
func (o Matrix) GetMat() (M [][]float64) {
	M = make([][]float64, o.m)
	for i := 0; i < o.m; i++ {
		M[i] = make([]float64, o.n)
		for j := 0; j < o.n; j++ {
			M[i][j] = o.data[i+j*o.m]
		}
	}
	return
}

// complex ///////////////////////////////////////////////////////////////////////////////////////

// MatrixC implements a column-major representation of a matrix of complex numbers by using a linear
// array that can be passed to Fortran code.
//
//  NOTE: the functions related to MatrixC do not check for the limits of indices and dimensions.
//        Panic may occur then.
//
//  Example:
//             _            _
//            |  0+0i  3+3i  |
//        M = |  1+1i  4+4i  |
//            |_ 2+2i  5+5i _|(m x n)
//
//     data[i+j*m] = M[i][j]
//
type MatrixC struct {
	m, n int          // dimensions
	data []complex128 // data array. column-major => Fortran
}

// NewMatrixC allocates a new MatrixC
func NewMatrixC(m, n int) (o *MatrixC) {
	o = new(MatrixC)
	o.m, o.n = m, n
	o.data = make([]complex128, m*n)
	return
}

// SetFromMat sets matrix with data from a nested slice structure; i.e. form a given la.Mat.
func (o *MatrixC) SetFromMat(a [][]complex128) {
	k := 0
	for j := 0; j < o.n; j++ {
		for i := 0; i < o.m; i++ {
			o.data[k] = a[i][j]
			k += 1
		}
	}
}

// Sets value
func (o *MatrixC) Set(i, j int, val complex128) {
	o.data[i+j*o.m] = val // col-major
}

// GetMat returns nested slice representation; i.e. a la.Mat structure
func (o MatrixC) GetMat() (M [][]complex128) {
	M = make([][]complex128, o.m)
	for i := 0; i < o.m; i++ {
		M[i] = make([]complex128, o.n)
		for j := 0; j < o.n; j++ {
			M[i][j] = o.data[i+j*o.m]
		}
	}
	return
}
