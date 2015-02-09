// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"github.com/cpmech/gosl/utl"
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
func (t *Triplet) Init(m, n, max int) {
	t.m, t.n, t.pos, t.max = m, n, 0, max
	t.i = make([]int, max)
	t.j = make([]int, max)
	t.x = make([]float64, max)
}

// Put inserts an element to a pre-allocated (with Init) triplet matrix
func (t *Triplet) Put(i, j int, x float64) {
	if t.pos >= t.max {
		utl.Panic(_sparsemat_err1, t.pos, t.max)
	}
	t.i[t.pos], t.j[t.pos], t.x[t.pos] = i, j, x
	t.pos++
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
		utl.Panic(_sparsemat_err3, a.m, a.n, o.m, o.n)
	}
	for k := 0; k < a.pos; k++ {
		o.Put(a.n+a.i[k], a.j[k], a.x[k]) // puts a
		o.Put(a.j[k], a.n+a.i[k], a.x[k]) // puts at
	}
}

// Start (re-)starts the insertion index within "t" in order to allow (re-)insertion of
// items using the Put method
func (t *Triplet) Start() {
	t.pos = 0
}

// Len returns the number of items just inserted in the triplet
func (t *Triplet) Len() int {
	return t.pos
}

// Max returns the maximum number of items that can be inserted in the triplet
func (t *Triplet) Max() int {
	return t.max
}

// TripletC is the equivalent to Triplet but with values stored as pairs of
// float64 representing complex numbers
type TripletC struct {
	m, n     int       // matrix dimension (rows, columns)
	pos, max int       // current position and max number of entries allowed (non-zeros, including repetitions)
	i, j     []int     // indices for each x values (size=max)
	x        []float64 // x,z split: values for each i, j (size=max)
	z        []float64 // x,z split: complex values
	xz       []float64 // xz monolithic: real and complex values: {real,cmplx,r,c,r,c,...}
}

// CCMatrixC is the equivalent to CCMatrix but with values stored as paris of
// float64 representing compressed numbers
type CCMatrixC struct {
	m, n int       // matrix dimension (rows, columns)
	nnz  int       // number of non-zeros
	p, i []int     // pointers and row indices (len(p)=n+1, len(i)=nnz)
	x    []float64 // values (len(x)=nnz)
	z    []float64 // complex values (len(z)=nnz)
}

// Init allocates all memory required to hold a sparse matrix in triplet (complex) form
func (t *TripletC) Init(m, n, max int, xzmonolithic bool) {
	t.m, t.n, t.pos, t.max = m, n, 0, max
	t.i = make([]int, max)
	t.j = make([]int, max)
	if xzmonolithic {
		t.xz = make([]float64, 2*max)
	} else {
		t.x = make([]float64, max)
		t.z = make([]float64, max)
	}
}

// Put inserts an element to a pre-allocated (with Init) triplet (complex) matrix
func (t *TripletC) Put(i, j int, x, z float64) {
	if t.pos >= t.max {
		utl.Panic(_sparsemat_err2, t.pos, t.max)
	}
	t.i[t.pos], t.j[t.pos] = i, j
	if t.xz != nil {
		t.xz[t.pos*2], t.xz[t.pos*2+1] = x, z
	} else {
		t.x[t.pos], t.z[t.pos] = x, z
	}
	t.pos++
}

// Start (re-)starts the insertion index within "t" in order to allow (re-)insertion of
// items using the Put method
func (t *TripletC) Start() {
	t.pos = 0
}

// Len returns the number of items just inserted in the complex triplet
func (t *TripletC) Len() int {
	return t.pos
}

// Max returns the maximum number of items that can be inserted in the complex triplet
func (t *TripletC) Max() int {
	return t.max
}

// ToDense converts a column-compressed matrix to dense form
func (a *CCMatrix) ToDense() [][]float64 {
	r := make([][]float64, a.m)
	for i := 0; i < a.m; i++ {
		r[i] = make([]float64, a.n)
	}
	for j := 0; j < a.n; j++ {
		for p := a.p[j]; p < a.p[j+1]; p++ {
			r[a.i[p]][j] = a.x[p]
		}
	}
	return r
}

// ToDense converts a column-compressed matrix (complex) to dense form
func (a *CCMatrixC) ToDense() [][]complex128 {
	r := make([][]complex128, a.m)
	for i := 0; i < a.m; i++ {
		r[i] = make([]complex128, a.n)
	}
	for j := 0; j < a.n; j++ {
		for p := a.p[j]; p < a.p[j+1]; p++ {
			r[a.i[p]][j] = complex(a.x[p], a.z[p])
		}
	}
	return r
}

// error messages
var (
	_sparsemat_err1 = "sparsemat.go: la.Triplet.Put: cannot put item because max number of items has been exceeded (pos = %d, max = %d)"
	_sparsemat_err2 = "sparsemat.go: la.TripletC.Put: cannot put item because max number of items has been exceeded (pos = %d, max = %d)"
	_sparsemat_err3 = "sparsemat.go: la.Triplet.PutMatAndMatT: cannot put larger matrix into sparse matrix.\nb := [[.. at] [a ..]] with len(a)=(%d,%d) and len(b)=(%d,%d)"
)
