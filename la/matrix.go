// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package la implements functions and structure for Linear Algebra computations. It defines a
// Vector and Matrix types for computations with dense data and also a Triplet and CCMatrix for
// sparse data.
package la

import (
	"math"
	"strings"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la/oblas"
	"github.com/cpmech/gosl/utl"
)

// Matrix implements a column-major representation of a matrix by using a linear array that can be passed to Fortran code
//
//  NOTE: the functions related to Matrix do not check for the limits of indices and dimensions.
//        Panic may occur then.
//
//  Example:
//             _      _
//            |  0  3  |
//        A = |  1  4  |
//            |_ 2  5 _|(m x n)
//
//     data[i+j*m] = A[i][j]
//
type Matrix struct {
	M, N int       // dimensions
	Data []float64 // data array. column-major => Fortran
}

// NewMatrix allocates a new (empty) Matrix with given (m,n) (row/col sizes)
func NewMatrix(m, n int) (o *Matrix) {
	o = new(Matrix)
	o.M, o.N = m, n
	o.Data = make([]float64, m*n)
	return
}

// NewMatrixDeep2 allocates a new Matrix from given (Deep2) nested slice.
// NOTE: make sure to have at least 1x1 item
func NewMatrixDeep2(a [][]float64) (o *Matrix) {
	o = new(Matrix)
	o.M, o.N = len(a), len(a[0])
	o.Data = make([]float64, o.M*o.N)
	o.SetFromDeep2(a)
	return
}

// SetFromDeep2 sets matrix with data from a nested slice (Deep2) structure
func (o *Matrix) SetFromDeep2(a [][]float64) {
	k := 0
	for j := 0; j < o.N; j++ {
		for i := 0; i < o.M; i++ {
			o.Data[k] = a[i][j]
			k++
		}
	}
}

// SetDiag sets diagonal matrix with diagonal components equal to val
func (o *Matrix) SetDiag(val float64) {
	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			if i == j {
				o.Data[i+j*o.M] = val
			} else {
				o.Data[i+j*o.M] = 0
			}
		}
	}
}

// Set sets value
func (o *Matrix) Set(i, j int, val float64) {
	o.Data[i+j*o.M] = val // col-major
}

// Get gets value
func (o *Matrix) Get(i, j int) float64 {
	return o.Data[i+j*o.M] // col-major
}

// GetDeep2 returns nested slice representation
func (o *Matrix) GetDeep2() (M [][]float64) {
	M = make([][]float64, o.M)
	for i := 0; i < o.M; i++ {
		M[i] = make([]float64, o.N)
		for j := 0; j < o.N; j++ {
			M[i][j] = o.Data[i+j*o.M]
		}
	}
	return
}

// GetCopy returns a copy of this matrix
func (o *Matrix) GetCopy() (clone *Matrix) {
	clone = NewMatrix(o.M, o.N)
	copy(clone.Data, o.Data)
	return
}

// GetTranspose returns the tranpose matrix
func (o *Matrix) GetTranspose() (tran *Matrix) {
	tran = NewMatrix(o.N, o.M)
	for i := 0; i < o.N; i++ {
		for j := 0; j < o.M; j++ {
			tran.Set(i, j, o.Get(j, i))
		}
	}
	return
}

// GetComplex returns a complex version of this matrix
func (o *Matrix) GetComplex() (b *MatrixC) {
	b = NewMatrixC(o.M, o.N)
	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			b.Set(i, j, complex(o.Get(i, j), 0))
		}
	}
	return
}

// CopyInto copies the scaled components of this matrix into another one (result)
//  result := α * this   ⇒   result[ij] := α * this[ij]
func (o *Matrix) CopyInto(result *Matrix, α float64) {
	for k := 0; k < o.M*o.N; k++ {
		result.Data[k] = α * o.Data[k]
	}
}

// Add adds value to (i,j) location
func (o *Matrix) Add(i, j int, val float64) {
	o.Data[i+j*o.M] += val // col-major
}

// Fill fills this matrix with a single number val
//  aij = val
func (o *Matrix) Fill(val float64) {
	for k := 0; k < o.M*o.N; k++ {
		o.Data[k] = val
	}
}

// ClearRC clear rows and columns and set diagonal components
//                 _         _                                     _         _
//  Example:      |  1 2 3 4  |                                   |  1 2 3 4  |
//            A = |  5 6 7 8  |  ⇒  clear([1,2], [], 1.0)  ⇒  A = |  0 1 0 0  |
//                |_ 4 3 2 1 _|                                   |_ 0 0 1 0 _|
//
func (o *Matrix) ClearRC(rows, cols []int, diag float64) {
	for _, r := range rows {
		for j := 0; j < o.N; j++ {
			if r == j {
				o.Set(r, j, diag)
			} else {
				o.Set(r, j, 0.0)
			}
		}
	}
	for _, c := range cols {
		for i := 0; i < o.M; i++ {
			if i == c {
				o.Set(i, c, diag)
			} else {
				o.Set(i, c, 0.0)
			}
		}
	}
}

// ClearBry clears boundaries
//                 _       _                          _       _
//  Example:      |  1 2 3  |                        |  1 0 0  |
//            A = |  4 5 6  |  ⇒  clear(1.0)  ⇒  A = |  0 5 0  |
//                |_ 7 8 9 _|                        |_ 0 0 1 _|
//
func (o *Matrix) ClearBry(diag float64) {
	o.ClearRC([]int{0, o.M - 1}, []int{0, o.N - 1}, diag)
}

// MaxDiff returns the maximum difference between the components of this and another matrix
func (o *Matrix) MaxDiff(another *Matrix) (maxdiff float64) {
	maxdiff = math.Abs(o.Data[0] - another.Data[0])
	for k := 1; k < o.M*o.N; k++ {
		diff := math.Abs(o.Data[k] - another.Data[k])
		if diff > maxdiff {
			maxdiff = diff
		}
	}
	return
}

// Largest returns the largest component |a[ij]| of this matrix, normalised by den
//   largest := |a[ij]| / den
func (o *Matrix) Largest(den float64) (largest float64) {
	largest = math.Abs(o.Data[0])
	for k := 1; k < o.M*o.N; k++ {
		tmp := math.Abs(o.Data[k])
		if tmp > largest {
			largest = tmp
		}
	}
	return largest / den
}

// Col access column j of this matrix. No copies are made since the internal data are in
// col-major format already.
// NOTE: this method can be used to modify the columns; e.g. with o.Col(0)[0] = 123
func (o *Matrix) Col(j int) Vector {
	return o.Data[j*o.M : (j+1)*o.M]
}

// GetRow returns row i of this matrix
func (o *Matrix) GetRow(i int) (row Vector) {
	row = make([]float64, o.N)
	for j := 0; j < o.N; j++ {
		row[j] = o.Data[i+j*o.M]
	}
	return
}

// GetCol returns column j of this matrix
func (o *Matrix) GetCol(j int) (col Vector) {
	col = make([]float64, o.M)
	copy(col, o.Data[j*o.M:(j+1)*o.M])
	return
}

// NormFrob returns the Frobenious norm of this matrix
//  nrm := ‖a‖_F = sqrt(Σ_i Σ_j a[ij]⋅a[ij]) = ‖a‖_2
func (o *Matrix) NormFrob() (nrm float64) {
	for k := 0; k < o.M*o.N; k++ {
		nrm += o.Data[k] * o.Data[k]
	}
	return math.Sqrt(nrm)
}

// NormInf returns the infinite norm of this matrix
//  nrm := ‖a‖_∞ = max_i ( Σ_j a[ij] )
func (o *Matrix) NormInf() (nrm float64) {
	for j := 0; j < o.N; j++ { // sum first row
		nrm += math.Abs(o.Data[j*o.M])
	}
	var sumrow float64
	for i := 1; i < o.M; i++ {
		sumrow = 0.0
		for j := 0; j < o.N; j++ { // sum the other rows
			sumrow += math.Abs(o.Data[i+j*o.M])
			if sumrow > nrm {
				nrm = sumrow
			}
		}
	}
	return
}

// Apply sets this matrix with the scaled components of another matrix
//  this := α * another   ⇒   this[i] := α * another[i]
//  NOTE: "another" may be "this"
func (o Matrix) Apply(α float64, another *Matrix) {
	for k := 0; k < o.M*o.N; k++ {
		o.Data[k] = α * another.Data[k]
	}
}

// Det computes the determinant of matrix using the LU factorization
//   NOTE: this method may fail due to overflow...
func (o *Matrix) Det() (det float64, err error) {
	if o.M != o.N {
		err = chk.Err("matrix must be square to compute determinant. %d × %d is invalid\n", o.M, o.N)
		return
	}
	ai := make([]float64, len(o.Data))
	copy(ai, o.Data)
	ipiv := make([]int32, utl.Imin(o.M, o.N))
	err = oblas.Dgetrf(o.M, o.N, ai, o.M, ipiv) // NOTE: ipiv are 1-based indices
	if err != nil {
		return
	}
	det = 1.0
	for i := 0; i < o.M; i++ {
		if ipiv[i]-1 == int32(i) { // NOTE: ipiv are 1-based indices
			det = +det * ai[i+i*o.M]
		} else {
			det = -det * ai[i+i*o.M]
		}
	}
	return
}

// Print prints matrix (without commas or brackets)
func (o *Matrix) Print(nfmt string) (l string) {
	if nfmt == "" {
		nfmt = "%g "
	}
	for i := 0; i < o.M; i++ {
		if i > 0 {
			l += "\n"
		}
		for j := 0; j < o.N; j++ {
			l += io.Sf(nfmt, o.Get(i, j))
		}
	}
	return
}

// PrintGo prints matrix in Go format
func (o *Matrix) PrintGo(nfmt string) (l string) {
	if nfmt == "" {
		nfmt = "%10g"
	}
	l = "[][]float64{\n"
	for i := 0; i < o.M; i++ {
		l += "    {"
		for j := 0; j < o.N; j++ {
			if j > 0 {
				l += ","
			}
			l += io.Sf(nfmt, o.Get(i, j))
		}
		l += "},\n"
	}
	l += "}"
	return
}

// PrintPy prints matrix in Python format
func (o *Matrix) PrintPy(nfmt string) (l string) {
	if nfmt == "" {
		nfmt = "%10g"
	}
	l = "np.matrix([\n"
	for i := 0; i < o.M; i++ {
		l += "    ["
		for j := 0; j < o.N; j++ {
			if j > 0 {
				l += ","
			}
			l += io.Sf(nfmt, o.Get(i, j))
		}
		l += "],\n"
	}
	l += "], dtype=float)"
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
//        A = |  1+1i  4+4i  |
//            |_ 2+2i  5+5i _|(m x n)
//
//     data[i+j*m] = A[i][j]
//
type MatrixC struct {
	M, N int          // dimensions
	Data []complex128 // data array. column-major => Fortran
}

// NewMatrixC allocates a new (empty) MatrixC with given (m,n) (row/col sizes)
func NewMatrixC(m, n int) (o *MatrixC) {
	o = new(MatrixC)
	o.M, o.N = m, n
	o.Data = make([]complex128, m*n)
	return
}

// NewMatrixDeep2c allocates a new MatrixC from given (Deep2c) nested slice.
// NOTE: make sure to have at least 1x1 items
func NewMatrixDeep2c(a [][]complex128) (o *MatrixC) {
	o = new(MatrixC)
	o.M, o.N = len(a), len(a[0])
	o.Data = make([]complex128, o.M*o.N)
	o.SetFromDeep2c(a)
	return
}

// SetFromDeep2c sets matrix with data from a nested slice (Deep2c) structure
func (o *MatrixC) SetFromDeep2c(a [][]complex128) {
	k := 0
	for j := 0; j < o.N; j++ {
		for i := 0; i < o.M; i++ {
			o.Data[k] = a[i][j]
			k++
		}
	}
}

// Set sets value
func (o *MatrixC) Set(i, j int, val complex128) {
	o.Data[i+j*o.M] = val // col-major
}

// Get gets value
func (o *MatrixC) Get(i, j int) complex128 {
	return o.Data[i+j*o.M] // col-major
}

// GetDeep2 returns nested slice representation
func (o *MatrixC) GetDeep2() (M [][]complex128) {
	M = make([][]complex128, o.M)
	for i := 0; i < o.M; i++ {
		M[i] = make([]complex128, o.N)
		for j := 0; j < o.N; j++ {
			M[i][j] = o.Data[i+j*o.M]
		}
	}
	return
}

// GetCopy returns a copy of this matrix
func (o *MatrixC) GetCopy() (clone *MatrixC) {
	clone = NewMatrixC(o.M, o.N)
	copy(clone.Data, o.Data)
	return
}

// GetTranspose returns the tranpose matrix
func (o *MatrixC) GetTranspose() (tran *MatrixC) {
	tran = NewMatrixC(o.N, o.M)
	for i := 0; i < o.N; i++ {
		for j := 0; j < o.M; j++ {
			tran.Set(i, j, o.Get(j, i))
		}
	}
	return
}

// Add adds value to (i,j) location
func (o *MatrixC) Add(i, j int, val complex128) {
	o.Data[i+j*o.M] += val // col-major
}

// Fill fills this matrix with a single number val
//  aij = val
func (o *MatrixC) Fill(val complex128) {
	for k := 0; k < o.M*o.N; k++ {
		o.Data[k] = val
	}
}

// Col access column j of this matrix. No copies are made since the internal data are in
// col-major format already.
// NOTE: this method can be used to modify the columns; e.g. with o.Col(0)[0] = 123
func (o *MatrixC) Col(j int) VectorC {
	return o.Data[j*o.M : (j+1)*o.M]
}

// GetRow returns row i of this matrix
func (o *MatrixC) GetRow(i int) (row VectorC) {
	row = make([]complex128, o.N)
	for j := 0; j < o.N; j++ {
		row[j] = o.Data[i+j*o.M]
	}
	return
}

// GetCol returns column j of this matrix
func (o *MatrixC) GetCol(j int) (col VectorC) {
	col = make([]complex128, o.M)
	copy(col, o.Data[j*o.M:(j+1)*o.M])
	return
}

// Apply sets this matrix with the scaled components of another matrix
//  this := α * another   ⇒   this[i] := α * another[i]
//  NOTE: "another" may be "this"
func (o MatrixC) Apply(α complex128, another *MatrixC) {
	for k := 0; k < o.M*o.N; k++ {
		o.Data[k] = α * another.Data[k]
	}
}

// Print prints matrix (without commas or brackets).
// NOTE: if non-empty, nfmtI must have '+' e.g. %+g
func (o *MatrixC) Print(nfmtR, nfmtI string) (l string) {
	if nfmtR == "" {
		nfmtR = "%g"
	}
	if nfmtI == "" {
		nfmtI = "%+g"
	}
	if !strings.ContainsRune(nfmtI, '+') {
		nfmtI = strings.Replace(nfmtI, "%", "%+", -1)
	}
	for i := 0; i < o.M; i++ {
		if i > 0 {
			l += "\n"
		}
		for j := 0; j < o.N; j++ {
			if j > 0 {
				l += " "
			}
			l += io.Sf(nfmtR, real(o.Get(i, j))) + io.Sf(nfmtI, imag(o.Get(i, j))) + "i"
		}
	}
	return
}

// PrintGo prints matrix in Go format
// NOTE: if non-empty, nfmtI must have '+' e.g. %+g
func (o *MatrixC) PrintGo(nfmtR, nfmtI string) (l string) {
	if nfmtR == "" {
		nfmtR = "%g"
	}
	if nfmtI == "" {
		nfmtI = "%+g"
	}
	if !strings.ContainsRune(nfmtI, '+') {
		nfmtI = strings.Replace(nfmtI, "%", "%+", -1)
	}
	l = "[][]complex128{\n"
	for i := 0; i < o.M; i++ {
		l += "    {"
		for j := 0; j < o.N; j++ {
			if j > 0 {
				l += ","
			}
			l += io.Sf(nfmtR, real(o.Get(i, j))) + io.Sf(nfmtI, imag(o.Get(i, j))) + "i"
		}
		l += "},\n"
	}
	l += "}"
	return
}

// PrintPy prints matrix in Python format
// NOTE: if non-empty, nfmtI must have '+' e.g. %+g
func (o *MatrixC) PrintPy(nfmtR, nfmtI string) (l string) {
	if nfmtR == "" {
		nfmtR = "%g"
	}
	if nfmtI == "" {
		nfmtI = "%+g"
	}
	if !strings.ContainsRune(nfmtI, '+') {
		nfmtI = strings.Replace(nfmtI, "%", "%+", -1)
	}
	l = "np.matrix([\n"
	for i := 0; i < o.M; i++ {
		l += "    ["
		for j := 0; j < o.N; j++ {
			if j > 0 {
				l += ","
			}
			l += io.Sf(nfmtR, real(o.Get(i, j))) + io.Sf(nfmtI, imag(o.Get(i, j))) + "j"
		}
		l += "],\n"
	}
	l += "], dtype=complex)"
	return
}
