// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"bytes"
	"math"
	"strings"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// Triplet is a simple representation of a sparse matrix, where the indices and values
// of this matrix are stored directly.
type Triplet struct {
	m, n     int       // matrix dimension (rows, columns)
	pos, max int       // current position and max number of entries allowed (non-zeros, including repetitions)
	i, j     []int     // indices for each x values (size=max)
	x        []float64 // values for each i, j (size=max)
}

// CCMatrix represents a sparse matrix using the so-called "column-compressed format".
type CCMatrix struct {
	m, n int       // matrix dimension (rows, columns)
	nnz  int       // number of non-zeros
	p, i []int     // pointers and row indices (len(p)=n+1, len(i)=nnz)
	x    []float64 // values (len(x)=nnz)
}

// Init allocates all memory required to hold a sparse matrix in triplet form
func (o *Triplet) Init(m, n, max int) {
	o.m, o.n, o.pos, o.max = m, n, 0, max
	o.i = make([]int, max)
	o.j = make([]int, max)
	o.x = make([]float64, max)
}

// Put inserts an element to a pre-allocated (with Init) triplet matrix
func (o *Triplet) Put(i, j int, x float64) {
	if o.pos >= o.max {
		chk.Panic("cannot put item because max number of items has been exceeded (pos = %d, max = %d)", o.pos, o.max)
	}
	o.i[o.pos], o.j[o.pos], o.x[o.pos] = i, j, x
	o.pos++
}

// PutMatAndMatT adds the content of a matrix "a" and its transpose "at" to triplet "o"
// ex:    0   1   2   3   4   5
//      [... ... ... a00 a10 ...] 0
//      [... ... ... a01 a11 ...] 1
//      [... ... ... a02 a12 ...] 2      [. at  .]
//      [a00 a01 a02 ... ... ...] 3  =>  [a  .  .]
//      [a10 a11 a12 ... ... ...] 4      [.  .  .]
//      [... ... ... ... ... ...] 5
func (o *Triplet) PutMatAndMatT(a *Triplet) {
	if a.n+a.m > o.m || a.n+a.m > o.n {
		chk.Panic("cannot put larger matrix into sparse matrix.\nb := [[.. at] [a ..]] with len(a)=(%d,%d) and len(b)=(%d,%d)", a.m, a.n, o.m, o.n)
	}
	for k := 0; k < a.pos; k++ {
		o.Put(a.n+a.i[k], a.j[k], a.x[k]) // puts a
		o.Put(a.j[k], a.n+a.i[k], a.x[k]) // puts at
	}
}

// PutCCMatAndMatT adds the content of a compressed-column matrix "a" and its transpose "at" to triplet "o"
// ex:    0   1   2   3   4   5
//      [... ... ... a00 a10 ...] 0
//      [... ... ... a01 a11 ...] 1
//      [... ... ... a02 a12 ...] 2      [. at  .]
//      [a00 a01 a02 ... ... ...] 3  =>  [a  .  .]
//      [a10 a11 a12 ... ... ...] 4      [.  .  .]
//      [... ... ... ... ... ...] 5
func (o *Triplet) PutCCMatAndMatT(a *CCMatrix) {
	if a.n+a.m > o.m || a.n+a.m > o.n {
		chk.Panic("cannot put larger matrix into sparse matrix.\nb := [[.. at] [a ..]] with len(a)=(%d,%d) and len(b)=(%d,%d)", a.m, a.n, o.m, o.n)
	}
	for j := 0; j < a.n; j++ {
		for k := a.p[j]; k < a.p[j+1]; k++ {
			o.Put(a.n+a.i[k], j, a.x[k]) // puts a
			o.Put(j, a.n+a.i[k], a.x[k]) // puts at
		}
	}
}

// Start (re-)starts the insertion index within "o" in order to allow (re-)insertion of
// items using the Put method
func (o *Triplet) Start() {
	o.pos = 0
}

// Len returns the number of items just inserted in the triplet
func (o *Triplet) Len() int {
	return o.pos
}

// Max returns the maximum number of items that can be inserted in the triplet
func (o *Triplet) Max() int {
	return o.max
}

// GetDenseMatrix returns the dense matrix corresponding to this Triplet
func (o *Triplet) GetDenseMatrix() (a *Matrix) {
	a = NewMatrix(o.m, o.n)
	for k := 0; k < o.max; k++ {
		a.Add(o.i[k], o.j[k], o.x[k])
	}
	return
}

// WriteSmat writes a ".smat" file that can be visualised with vismatrix
//  tol -- tolerance to skip zero values
func (o *Triplet) WriteSmat(dirout, fnkey string, tol float64) {
	var bfa, bfb bytes.Buffer
	var nnz int
	for k := 0; k < o.pos; k++ {
		if math.Abs(o.x[k]) > tol {
			io.Ff(&bfb, "  %d  %d  %23.15e\n", o.i[k], o.j[k], o.x[k])
			nnz++
		}
	}
	io.Ff(&bfa, "%d  %d  %d\n", o.m, o.n, nnz)
	io.WriteFileD(dirout, fnkey+".smat", &bfa, &bfb)
}

// ReadSmat reads ".smat" file back
func (o *Triplet) ReadSmat(filename string) (err error) {
	var e error
	io.ReadLines(filename, func(idx int, line string) (stop bool) {
		r := strings.Fields(line)
		if len(r) != 3 {
			e = chk.Err("number of columns must be 3\n")
			return true // stop
		}
		if idx == 0 {
			m, n, nnz := io.Atoi(r[0]), io.Atoi(r[1]), io.Atoi(r[2])
			o.Init(m, n, nnz)
		} else {
			i, j, x := io.Atoi(r[0]), io.Atoi(r[1]), io.Atof(r[2])
			o.Put(i, j, x)
		}
		return
	})
	return
}

// ToDense converts a column-compressed matrix to dense form
func (o *CCMatrix) ToDense() (res *Matrix) {
	res = NewMatrix(o.m, o.n)
	for j := 0; j < o.n; j++ {
		for p := o.p[j]; p < o.p[j+1]; p++ {
			res.Set(o.i[p], j, o.x[p])
		}
	}
	return
}

// Set sets column-compressed matrix directly
func (o *CCMatrix) Set(m, n int, Ap, Ai []int, Ax []float64) {
	if len(Ap)-1 != n {
		chk.Panic("len(Ap) must be equal to n. %d != %d", len(Ap), n)
	}
	nnz := len(Ai)
	if len(Ax) != nnz {
		chk.Panic("len(Ax) must be equal to len(Ai) == nnz. %d != %d", len(Ax), nnz)
	}
	if Ap[n] != nnz {
		chk.Panic("last item in Ap must be equal to nnz. %d != %d", Ap[n], nnz)
	}
	o.m, o.n, o.nnz = m, n, nnz
	o.p, o.i, o.x = Ap, Ai, Ax
}

// complex /////////////////////////////////////////////////////////////////////////////////////////

// TripletC is a simple representation of a sparse matrix, where the indices and values
// of this matrix are stored directly. (complex version)
type TripletC struct {
	m, n     int          // matrix dimension (rows, columns)
	pos, max int          // current position and max number of entries allowed (non-zeros, including repetitions)
	i, j     []int        // indices for each x values (size=max)
	x        []complex128 // values for each i, j (size=max)
}

// CCMatrixC represents a sparse matrix using the so-called "column-compressed format".
// (complex version)
type CCMatrixC struct {
	m, n int          // matrix dimension (rows, columns)
	nnz  int          // number of non-zeros
	p, i []int        // pointers and row indices (len(p)=n+1, len(i)=nnz)
	x    []complex128 // values (len(x)=nnz)
}

// Init allocates all memory required to hold a sparse matrix in triplet (complex) form
func (o *TripletC) Init(m, n, max int) {
	o.m, o.n, o.pos, o.max = m, n, 0, max
	o.i = make([]int, max)
	o.j = make([]int, max)
	o.x = make([]complex128, max)
}

// Put inserts an element to a pre-allocated (with Init) triplet (complex) matrix
func (o *TripletC) Put(i, j int, x complex128) {
	if o.pos >= o.max {
		chk.Panic("cannot put item because max number of items has been exceeded (pos = %d, max = %d)", o.pos, o.max)
	}
	o.i[o.pos], o.j[o.pos], o.x[o.pos] = i, j, x
	o.pos++
}

// Start (re-)starts the insertion index within "o" in order to allow (re-)insertion of
// items using the Put method
func (o *TripletC) Start() {
	o.pos = 0
}

// Len returns the number of items just inserted in the complex triplet
func (o *TripletC) Len() int {
	return o.pos
}

// Max returns the maximum number of items that can be inserted in the complex triplet
func (o *TripletC) Max() int {
	return o.max
}

// GetDenseMatrix returns the dense matrix corresponding to this Triplet
func (o *TripletC) GetDenseMatrix() (a *MatrixC) {
	a = NewMatrixC(o.m, o.n)
	for k := 0; k < o.max; k++ {
		a.Add(o.i[k], o.j[k], o.x[k])
	}
	return
}

// ToDense converts a column-compressed matrix (complex) to dense form
func (a *CCMatrixC) ToDense() (res *MatrixC) {
	res = NewMatrixC(a.m, a.n)
	for j := 0; j < a.n; j++ {
		for p := a.p[j]; p < a.p[j+1]; p++ {
			res.Set(a.i[p], j, a.x[p])
		}
	}
	return
}
