// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"bytes"
	"gosl/mpi"
	"math"
	"math/cmplx"
	"strings"

	"gosl/chk"
	"gosl/io"
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

// NewTriplet returns a new Triplet. This is a wrapper to new(Triplet) followed by Init()
func NewTriplet(m, n, max int) (o *Triplet) {
	o = new(Triplet)
	o.Init(m, n, max)
	return
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

// Start (re)starts index for inserting items using the Put command
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

// Size returns the row/column size of the matrix represented by the Triplet
func (o *Triplet) Size() (m, n int) {
	return o.m, o.n
}

// ToDense returns the dense matrix corresponding to this Triplet
func (o *Triplet) ToDense() (a *Matrix) {
	a = NewMatrix(o.m, o.n)
	for k := 0; k < o.max; k++ {
		a.Add(o.i[k], o.j[k], o.x[k])
	}
	return
}

// ReadSmat reads a SMAT file or a MatrixMarket file
//
//  About the .smat file:
//   - lines starting with the percent mark (%) are ignored (they are comments)
//   - the first non-comment line contains the number of rows, columns, and non-zero entries
//   - the following lines contain the indices of row, column, and the non-zero entry
//
//  Example of .smat file (0-based indices):
//
//     % this is a comment
//     % ---------------------
//     % m n nnz
//     %  i j x
//     %   ...
//     %  i j x
//     % ---------------------
//        5  5  8
//          0     0   1.000e+00
//          1     1   1.050e+01
//          2     2   1.500e-02
//          0     3   6.000e+00
//          3     1   2.505e+02
//          3     3  -2.800e+02
//          3     4   3.332e+01
//          4     4   1.200e+01
//
//  Example of MatrixMarket file (1-based indices):
//
//     %%MatrixMarket matrix coordinate real general
//     %=================================================================================
//     %
//     % This ASCII file represents a sparse MxN matrix with L
//     % nonzeros in the following Matrix Market format:
//     %
//     % Reference: https://math.nist.gov/MatrixMarket/formats.html
//     %
//     % +----------------------------------------------+
//     % |%%MatrixMarket matrix coordinate real general | <--- header line
//     % |%                                             | <--+
//     % |% comments                                    |    |-- 0 or more comment lines
//     % |%                                             | <--+
//     % |    M  N  L                                   | <--- rows, columns, entries
//     % |    I1  J1  A(I1, J1)                         | <--+
//     % |    I2  J2  A(I2, J2)                         |    |
//     % |    I3  J3  A(I3, J3)                         |    |-- L lines
//     % |        . . .                                 |    |
//     % |    IL JL  A(IL, JL)                          | <--+
//     % +----------------------------------------------+
//     %
//     % Indices are 1-based, i.e. A(1,1) is the first element.
//     %
//     %=================================================================================
//       5  5  8
//         1     1   1.000e+00
//         2     2   1.050e+01
//         3     3   1.500e-02
//         1     4   6.000e+00
//         4     2   2.505e+02
//         4     4  -2.800e+02
//         4     5   3.332e+01
//         5     5   1.200e+01
//
//  NOTE: this function can only read a "coordinate" type MatrixMarket at the moment
//
//  Input:
//   filename -- filename
//   mirrorIfSym -- do mirror off-diagonal entries is the matrix symmetric;
//      i.e. for each i != j, set both A(i,j) = A(j,i) with the value just read from file.
//      This option helps when you want all non-zero values of a symmetric matrix,
//      noting that the MatrixMarket format stores only one "triangular side" of
//      the matrix and the diagonal.
//   comm -- [may be nil]. If the communication is given, this function will
//      distribute the number of non-zeros (near) equally to all processors.
//
//  Output:
//   symmetric -- [MatrixMarket only] return true if the MatrixMarket header has "symmetric"
//
func (o *Triplet) ReadSmat(filename string, mirrorIfSym bool, comm *mpi.Communicator) (symmetric bool) {
	deltaIndex := 0
	initialized := false
	id, sz := 0, 1
	if comm != nil {
		id, sz = comm.Rank(), comm.Size()
	}
	start, endp1 := 0, 0
	indexNnz := 0
	io.ReadLines(filename, func(idx int, line string) (stop bool) {
		if strings.HasPrefix(line, "%%MatrixMarket") {
			info := strings.Fields(line)
			if info[1] != "matrix" {
				chk.Panic("can only read \"matrix\" MatrixMarket at the moment")
			}
			if info[2] != "coordinate" {
				chk.Panic("can only read \"coordinate\" MatrixMarket at the moment")
			}
			if info[3] != "real" {
				chk.Panic("the given MatrixMarket file must have the word \"real\" in the header")
			}
			if info[4] != "general" && info[4] != "symmetric" && info[4] != "unsymmetric" {
				chk.Panic("this function only works with \"general\", \"symmetric\" and \"unsymmetric\" MatrixMarket files")
			}
			if info[4] == "symmetric" {
				symmetric = true
			}
			deltaIndex = 1
			return
		}
		if strings.HasPrefix(line, "%") {
			return
		}
		r := strings.Fields(line)
		if !initialized {
			if len(r) != 3 {
				chk.Panic("the number of columns in the line with dimensions must be 3 (m,n,nnz)\n")
			}
			m, n, nnz := io.Atoi(r[0]), io.Atoi(r[1]), io.Atoi(r[2])
			start, endp1 = (id*nnz)/sz, ((id+1)*nnz)/sz
			if symmetric {
				o.Init(m, n, (endp1-start)*2) // assuming that the diagonal is all-zeros (for safety)
			} else {
				o.Init(m, n, endp1-start)
			}
			initialized = true
		} else {
			if len(r) != 3 {
				chk.Panic("the number of columns in the data lines must be 3 (i,j,x)\n")
			}
			i, j, x := io.Atoi(r[0]), io.Atoi(r[1]), io.Atof(r[2])
			if indexNnz >= start && indexNnz < endp1 {
				o.Put(i-deltaIndex, j-deltaIndex, x)
				if symmetric && mirrorIfSym && i != j {
					o.Put(j-deltaIndex, i-deltaIndex, x)
				}
			}
			indexNnz++
		}
		return
	})
	return
}

// WriteSmat writes a SMAT file (that can be visualised with vismatrix) or a MatrixMarket file
//
//   For more information, see:
//
//           func (o *Triplet) ReadSmat()
//
//  dirout -- directory (to be created if not empty) where the file is saved
//  fnkey -- filename without extension (we add .smat or .mtx if matrixMarket == true)
//  tol -- tolerance to ignore near-zero values. only save values such that |value| > tol
//  format -- format for real numbers; e.g. "%23.15g" [default is "%g"]
//  matrixMarket -- save according to the matrixMarket file format (1-based indices + header)
//  enforceSymmetry -- [MatrixMarket only] ignore upper band of the matrix and save only the lower band + main diagonal
//
//  NOTE: This function converts the Triplet into CCMatrix (returned)
//        because there may be repeated entries (added)
//
func (o *Triplet) WriteSmat(dirout, fnkey string, tol float64, format string, matrixMarket, enforceSymmetry bool) (cmat *CCMatrix) {
	cmat = o.ToMatrix(nil)
	cmat.WriteSmat(dirout, fnkey, tol, format, matrixMarket, enforceSymmetry)
	return
}

// WriteSmat writes a SMAT file (that can be visualised with vismatrix) or a MatrixMarket file
//
//   For more information, see:
//
//           func (o *Triplet) ReadSmat()
//
//  dirout -- directory (to be created if not empty) where the file is saved
//  fnkey -- filename without extension (we add .smat or .mtx if matrixMarket == true)
//  tol -- tolerance to ignore near-zero values. only save values such that |value| > tol
//  format -- format for real numbers; e.g. "%23.15g" [default is "%g"]
//  matrixMarket -- save according to the matrixMarket file format (1-based indices + header)
//  enforceSymmetry -- [MatrixMarket only] ignore upper band of the matrix and save only the lower band + main diagonal
//
func (o *CCMatrix) WriteSmat(dirout, fnkey string, tol float64, format string, matrixMarket, enforceSymmetry bool) {
	fmtVal := "%g"
	if format != "" {
		fmtVal = format
	}
	fmtStr := "%d %d " + fmtVal + "\n"
	deltaIndex := 0
	if matrixMarket {
		deltaIndex = 1
	}
	var bfa, bfb bytes.Buffer
	var nnz int
	for j := 0; j < o.n; j++ {
		for p := o.p[j]; p < o.p[j+1]; p++ {
			if math.Abs(o.x[p]) > tol {
				row, col := o.i[p]+deltaIndex, j+deltaIndex
				if matrixMarket && enforceSymmetry {
					if col > row {
						continue
					}
				}
				io.Ff(&bfb, fmtStr, row, col, o.x[p])
				nnz++
			}
		}
	}
	ext := ".smat"
	if matrixMarket {
		ext = ".mtx"
		kind := "general"
		if enforceSymmetry {
			kind = "symmetric"
		}
		header := "%%%%MatrixMarket matrix coordinate real " + kind + "\n"
		io.Ff(&bfa, header)
	}
	io.Ff(&bfa, "%d %d %d\n", o.m, o.n, nnz)
	io.WriteFileVD(dirout, fnkey+ext, &bfa, &bfb)
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

// NewTripletC returns a new TripletC. This is a wrapper to new(TripletC) followed by Init()
func NewTripletC(m, n, max int) (o *TripletC) {
	o = new(TripletC)
	o.Init(m, n, max)
	return
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

// Start (re)starts index for inserting items using the Put command
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

// ToDense returns the dense matrix corresponding to this Triplet
func (o *TripletC) ToDense() (a *MatrixC) {
	a = NewMatrixC(o.m, o.n)
	for k := 0; k < o.max; k++ {
		a.Add(o.i[k], o.j[k], o.x[k])
	}
	return
}

// ReadSmat reads a SMAT file or a MatrixMarket file
//
//  About the .smat file:
//   - lines starting with the percent mark (%) are ignored (they are comments)
//   - the first non-comment line contains the number of rows, columns, and non-zero entries
//   - the following lines contain the indices of row, column, and the non-zero entry
//
//  Example of .smat file (0-based indices):
//
//     % this is a comment
//     % ---------------------
//     % m n nnz
//     %  i j real(x) imag(x)
//     %   ...
//     %  i j real(x) imag(x)
//     % ---------------------
//        5  5  8
//          0     0   1.000e+00  0.0
//          1     1   1.050e+01  0.0
//          2     2   1.500e-02  0.1
//          0     3   6.000e+00  0.1
//          3     1   2.505e+02  0.0
//          3     3  -2.800e+02  0.0
//          3     4   3.332e+01  0.2
//          4     4   1.200e+01  0.2
//
//  Example of MatrixMarket file (1-based indices):
//
//     %%MatrixMarket matrix coordinate complex general
//     %=================================================================================
//     %
//     % This ASCII file represents a sparse MxN matrix with L
//     % nonzeros in the following Matrix Market format:
//     %
//     % Reference: https://math.nist.gov/MatrixMarket/formats.html
//     %
//     % +-------------------------------------------------+
//     % |%%MatrixMarket matrix coordinate complex general | <--- header line
//     % |%                                                | <--+
//     % |% comments                                       |    |-- 0 or more comment lines
//     % |%                                                | <--+
//     % |    M  N  L                                      | <--- rows, columns, entries
//     % |    I1  J1  real(A(I1, J1)) imag(A(I1, J1))      | <--+
//     % |    I2  J2  real(A(I2, J2)) imag(A(I2, J2))      |    |
//     % |    I3  J3  real(A(I3, J3)) imag(A(I3, J3))      |    |-- L lines
//     % |        . . .                                    |    |
//     % |    IL JL  real(A(IL, JL)) imag(A(IL, JL))       | <--+
//     % +-------------------------------------------------+
//     %
//     % Indices are 1-based, i.e. A(1,1) is the first element.
//     %
//     %=================================================================================
//       5  5  8
//         1     1   1.000e+00  0.0
//         2     2   1.050e+01  0.0
//         3     3   1.500e-02  0.1
//         1     4   6.000e+00  0.1
//         4     2   2.505e+02  0.0
//         4     4  -2.800e+02  0.0
//         4     5   3.332e+01  0.2
//         5     5   1.200e+01  0.2
//
//  NOTE: this function can only read a "coordinate" type MatrixMarket at the moment
//
//  Input:
//   filename -- filename
//   mirrorIfSym -- do mirror off-diagonal entries is the matrix symmetric;
//      i.e. for each i != j, set both A(i,j) = A(j,i) with the value just read from file.
//      This option helps when you want all non-zero values of a symmetric matrix,
//      noting that the MatrixMarket format stores only one "triangular side" of
//      the matrix and the diagonal.
//   comm -- [may be nil]. If the communication is given, this function will
//      distribute the number of non-zeros (near) equally to all processors.
//
//  Output:
//   symmetric -- [MatrixMarket only] return true if the MatrixMarket header has "symmetric"
//
func (o *TripletC) ReadSmat(filename string, mirrorIfSym bool, comm *mpi.Communicator) (symmetric bool) {
	deltaIndex := 0
	initialized := false
	id, sz := 0, 1
	if comm != nil {
		id, sz = comm.Rank(), comm.Size()
	}
	start, endp1 := 0, 0
	indexNnz := 0
	io.ReadLines(filename, func(idx int, line string) (stop bool) {
		if strings.HasPrefix(line, "%%MatrixMarket") {
			info := strings.Fields(line)
			if info[1] != "matrix" {
				chk.Panic("can only read \"matrix\" MatrixMarket at the moment")
			}
			if info[2] != "coordinate" {
				chk.Panic("can only read \"coordinate\" MatrixMarket at the moment")
			}
			if info[3] != "complex" {
				chk.Panic("the given MatrixMarket file must have the word \"complex\" in the header")
			}
			if info[4] != "general" && info[4] != "symmetric" && info[4] != "unsymmetric" {
				chk.Panic("this function only works with \"general\", \"symmetric\" and \"unsymmetric\" MatrixMarket files")
			}
			if info[4] == "symmetric" {
				symmetric = true
			}
			deltaIndex = 1
			return
		}
		if strings.HasPrefix(line, "%") {
			return
		}
		r := strings.Fields(line)
		if !initialized {
			if len(r) != 3 {
				chk.Panic("number of columns in header must be 3 (m,n,nnz)\n")
			}
			m, n, nnz := io.Atoi(r[0]), io.Atoi(r[1]), io.Atoi(r[2])
			start, endp1 = (id*nnz)/sz, ((id+1)*nnz)/sz
			if symmetric {
				o.Init(m, n, (endp1-start)*2) // assuming that the diagonal is all-zeros (for safety)
			} else {
				o.Init(m, n, endp1-start)
			}
			initialized = true
		} else {
			if len(r) != 4 {
				chk.Panic("number of columns in data lines must be 4 (i,j,xReal,xImag)\n")
			}
			i, j, x := io.Atoi(r[0]), io.Atoi(r[1]), complex(io.Atof(r[2]), io.Atof(r[3]))
			if indexNnz >= start && indexNnz < endp1 {
				o.Put(i-deltaIndex, j-deltaIndex, x)
				if symmetric && mirrorIfSym && i != j {
					o.Put(j-deltaIndex, i-deltaIndex, x)
				}
			}
			indexNnz++
		}
		return
	})
	return
}

// WriteSmat writes a SMAT file (that can be visualised with vismatrix) or a MatrixMarket file
//
//   For more information, see:
//
//           func (o *TripletC) ReadSmat()
//
//  dirout -- directory (to be created if not empty) where the file is saved
//  fnkey -- filename without extension (we add .smat or .mtx if matrixMarket == true)
//  tol -- tolerance to ignore near-zero values. only save values such that |real(value)| > tol OR |imag(value)| > tol
//  format -- format for numbers; e.g. "%23.15g" [default is "%g"]
//  matrixMarket -- save according to the matrixMarket file format (1-based indices + header)
//  enforceSymmetry -- [MatrixMarket only] ignore upper band of the matrix and save only the lower band + main diagonal
//  norm -- writes a different matrix (real) such that the entries are the abs(entry) [modulus matrix]
//
//  NOTE: This function converts the Triplet into CCMatrixC (returned)
//        because there may be repeated entries (added)
//
func (o *TripletC) WriteSmat(dirout, fnkey string, tol float64, format string, matrixMarket, enforceSymmetry, norm bool) (cmat *CCMatrixC) {
	cmat = o.ToMatrix(nil)
	cmat.WriteSmat(dirout, fnkey, tol, format, matrixMarket, enforceSymmetry, norm)
	return
}

// WriteSmat writes a SMAT file (that can be visualised with vismatrix) or a MatrixMarket file
//
//   For more information, see:
//
//           func (o *TripletC) ReadSmat()
//
//  dirout -- directory (to be created if not empty) where the file is saved
//  fnkey -- filename without extension (we add .smat or .mtx if matrixMarket == true)
//  tol -- tolerance to ignore near-zero values. only save values such that |real(value)| > tol OR |imag(value)| > tol
//  format -- format for numbers; e.g. "%23.15g" [default is "%g"]
//  matrixMarket -- save according to the matrixMarket file format (1-based indices + header)
//  enforceSymmetry -- [MatrixMarket only] ignore upper band of the matrix and save only the lower band + main diagonal
//  norm -- writes a different matrix (real) such that the entries are the abs(entry) [modulus matrix]
//
func (o *CCMatrixC) WriteSmat(dirout, fnkey string, tol float64, format string, matrixMarket, enforceSymmetry, norm bool) {
	fmtVal := "%g"
	if format != "" {
		fmtVal = format
	}
	fmtStr := "%d %d " + fmtVal + " " + fmtVal + "\n"
	deltaIndex := 0
	if matrixMarket {
		deltaIndex = 1
	}
	var bfa, bfb bytes.Buffer
	var nnz int
	dataType := "complex"
	if norm {
		dataType = "real"
		fmtStr = "%d %d " + fmtVal + "\n"
	}
	for j := 0; j < o.n; j++ {
		for p := o.p[j]; p < o.p[j+1]; p++ {
			if math.Abs(real(o.x[p])) > tol || math.Abs(imag(o.x[p])) > tol {
				row, col := o.i[p]+deltaIndex, j+deltaIndex
				if matrixMarket && enforceSymmetry {
					if col > row {
						continue
					}
				}
				if norm {
					io.Ff(&bfb, fmtStr, row, col, cmplx.Abs(o.x[p]))
				} else {
					io.Ff(&bfb, fmtStr, row, col, real(o.x[p]), imag(o.x[p]))
				}
				nnz++
			}
		}
	}
	ext := ".smat"
	if matrixMarket {
		ext = ".mtx"
		kind := "general"
		if enforceSymmetry {
			kind = "symmetric"
		}
		header := "%%%%MatrixMarket matrix coordinate " + dataType + " " + kind + "\n"
		io.Ff(&bfa, header)
	}
	io.Ff(&bfa, "%d %d %d\n", o.m, o.n, nnz)
	io.WriteFileVD(dirout, fnkey+ext, &bfa, &bfb)
}

// ToDense converts a column-compressed matrix (complex) to dense form
func (o *CCMatrixC) ToDense() (res *MatrixC) {
	res = NewMatrixC(o.m, o.n)
	for j := 0; j < o.n; j++ {
		for p := o.p[j]; p < o.p[j+1]; p++ {
			res.Set(o.i[p], j, o.x[p])
		}
	}
	return
}
