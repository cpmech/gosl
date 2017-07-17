// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestMatrix01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Matrix01. (real) new and print")

	A := [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 0, -1, -2},
	}

	a := NewMatrixSlice(A)
	chk.Array(tst, "A to a", 1e-15, a.Data, []float64{1, 5, 9, 2, 6, 0, 3, 7, -1, 4, 8, -2})

	chk.Float64(tst, "Get(0,0)", 1e-17, a.Get(0, 0), 1)
	chk.Float64(tst, "Get(0,1)", 1e-17, a.Get(0, 1), 2)
	chk.Float64(tst, "Get(0,2)", 1e-17, a.Get(0, 2), 3)
	chk.Float64(tst, "Get(0,3)", 1e-17, a.Get(0, 3), 4)

	chk.Float64(tst, "Get(1,0)", 1e-17, a.Get(1, 0), 5)
	chk.Float64(tst, "Get(1,1)", 1e-17, a.Get(1, 1), 6)
	chk.Float64(tst, "Get(1,2)", 1e-17, a.Get(1, 2), 7)
	chk.Float64(tst, "Get(1,3)", 1e-17, a.Get(1, 3), 8)

	chk.Float64(tst, "Get(2,0)", 1e-17, a.Get(2, 0), 9)
	chk.Float64(tst, "Get(2,1)", 1e-17, a.Get(2, 1), 0)
	chk.Float64(tst, "Get(2,2)", 1e-17, a.Get(2, 2), -1)
	chk.Float64(tst, "Get(2,3)", 1e-17, a.Get(2, 3), -2)

	Aback := a.GetSlice()
	chk.Matrix(tst, "a to A", 1e-15, Aback, A)

	l := a.Print("")
	chk.String(tst, l, "1 2 3 4 \n5 6 7 8 \n9 0 -1 -2 ")

	l = a.PrintGo("%2g")
	lCorrect := "[][]float64{\n    { 1, 2, 3, 4},\n    { 5, 6, 7, 8},\n    { 9, 0,-1,-2},\n}"
	chk.String(tst, l, lCorrect)

	l = a.PrintPy("%2g")
	lCorrect = "np.matrix([\n    [ 1, 2, 3, 4],\n    [ 5, 6, 7, 8],\n    [ 9, 0,-1,-2],\n], dtype=float)"
	chk.String(tst, l, lCorrect)

	a.Add(1, 2, -7)
	chk.Float64(tst, "Get(1,2)", 1e-17, a.Get(1, 2), 0)

	b := a.GetCopy()
	if b.M != a.M {
		tst.Errorf("b.M should be equal to a.M\n")
		return
	}
	if b.N != a.N {
		tst.Errorf("b.N should be equal to a.N\n")
		return
	}
	chk.Array(tst, "b.Data", 1e-17, b.Data, a.Data)
}

func TestMatrix02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Matrix02. (complex) new and print")

	A := [][]complex128{
		{1 + 0.1i, 2, 3, 4 - 0.4i},
		{5 + 0.5i, 6, 7, 8 - 0.8i},
		{9 + 0.9i, 0, -1, -2 + 1i},
	}

	a := NewMatrixSliceC(A)
	chk.ArrayC(tst, "A to a", 1e-15, a.Data, []complex128{1 + 0.1i, 5 + 0.5i, 9 + 0.9i, 2, 6, 0, 3, 7, -1, 4 - 0.4i, 8 - 0.8i, -2 + 1i})

	chk.Complex128(tst, "Get(0,0)", 1e-17, a.Get(0, 0), 1+0.1i)
	chk.Complex128(tst, "Get(0,1)", 1e-17, a.Get(0, 1), 2)
	chk.Complex128(tst, "Get(0,2)", 1e-17, a.Get(0, 2), 3)
	chk.Complex128(tst, "Get(0,3)", 1e-17, a.Get(0, 3), 4-0.4i)

	chk.Complex128(tst, "Get(1,0)", 1e-17, a.Get(1, 0), 5+0.5i)
	chk.Complex128(tst, "Get(1,1)", 1e-17, a.Get(1, 1), 6)
	chk.Complex128(tst, "Get(1,2)", 1e-17, a.Get(1, 2), 7)
	chk.Complex128(tst, "Get(1,3)", 1e-17, a.Get(1, 3), 8-0.8i)

	chk.Complex128(tst, "Get(2,0)", 1e-17, a.Get(2, 0), 9+0.9i)
	chk.Complex128(tst, "Get(2,1)", 1e-17, a.Get(2, 1), 0)
	chk.Complex128(tst, "Get(2,2)", 1e-17, a.Get(2, 2), -1)
	chk.Complex128(tst, "Get(2,3)", 1e-17, a.Get(2, 3), -2+1i)

	Aback := a.GetSlice()
	chk.MatrixC(tst, "a to A", 1e-15, Aback, A)

	l := a.Print("%g", "")
	chk.String(tst, l, "1+0.1i, 2+0i, 3+0i, 4-0.4i\n5+0.5i, 6+0i, 7+0i, 8-0.8i\n9+0.9i, 0+0i, -1+0i, -2+1i")

	l = a.PrintGo("%2g", "%+4.1f")
	lCorrect := "[][]complex128{\n    { 1+0.1i, 2+0.0i, 3+0.0i, 4-0.4i},\n    { 5+0.5i, 6+0.0i, 7+0.0i, 8-0.8i},\n    { 9+0.9i, 0+0.0i,-1+0.0i,-2+1.0i},\n}"
	chk.String(tst, l, lCorrect)

	l = a.PrintPy("%2g", "%4.1f")
	lCorrect = "np.matrix([\n    [ 1+0.1j, 2+0.0j, 3+0.0j, 4-0.4j],\n    [ 5+0.5j, 6+0.0j, 7+0.0j, 8-0.8j],\n    [ 9+0.9j, 0+0.0j,-1+0.0j,-2+1.0j],\n], dtype=complex)"
	chk.String(tst, l, lCorrect)

	a.Add(1, 3, -8+0.8i)
	chk.Complex128(tst, "Get(1,3)", 1e-17, a.Get(1, 3), 0)

	b := a.GetCopy()
	if b.M != a.M {
		tst.Errorf("b.M should be equal to a.M\n")
		return
	}
	if b.N != a.N {
		tst.Errorf("b.N should be equal to a.N\n")
		return
	}
	chk.ArrayC(tst, "b.Data", 1e-17, b.Data, a.Data)
}

func TestMatrix03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Matrix03. (real) matrix methods")

	// Set
	a := NewMatrix(3, 5)
	a.Set(0, 0, 1)
	a.Set(0, 1, 2)
	a.Set(0, 2, 3)
	a.Set(0, 3, 4)
	a.Set(0, 4, 5)
	a.Set(1, 0, 0.1)
	a.Set(1, 1, 0.2)
	a.Set(1, 2, 0.3)
	a.Set(1, 3, 0.4)
	a.Set(1, 4, 0.5)
	a.Set(2, 0, 10)
	a.Set(2, 1, 20)
	a.Set(2, 2, 30)
	a.Set(2, 3, 40)
	a.Set(2, 4, 50)
	chk.Matrix(tst, "a", 1e-17, a.GetSlice(), [][]float64{
		{1, 2, 3, 4, 5},
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{10, 20, 30, 40, 50},
	})

	// GetCopy
	aclone := a.GetCopy()
	chk.Array(tst, "aclone", 1e-17, a.Data, aclone.Data)

	// Fill
	b := NewMatrix(5, 3)
	b.Fill(2)
	chk.Matrix(tst, "b", 1e-17, b.GetSlice(), [][]float64{
		{2, 2, 2},
		{2, 2, 2},
		{2, 2, 2},
		{2, 2, 2},
		{2, 2, 2},
	})

	// Scale
	c := NewMatrix(5, 3)
	c.Fill(2)
	c.Apply(1.0/4.0, c)
	chk.Matrix(tst, "c := c/4", 1e-17, c.GetSlice(), [][]float64{
		{0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5},
	})

	// MatCopy
	d := NewMatrix(3, 5)
	a.CopyInto(d, 1)
	chk.Matrix(tst, "d", 1e-17, d.GetSlice(), [][]float64{
		{1, 2, 3, 4, 5},
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{10, 20, 30, 40, 50},
	})

	// SetDiag
	e := NewMatrix(3, 3)
	e.SetDiag(1)
	chk.Matrix(tst, "e", 1e-17, e.GetSlice(), [][]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	})

	// MaxDiff
	f := NewMatrixSlice([][]float64{
		{1.1, 2.2, 3.3, 4.4, 5.5},
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{1, 2, 3, 4, 5},
	})
	maxdiff := a.MaxDiff(f)
	chk.Float64(tst, "MaxDiff", 1e-17, maxdiff, 45)

	// Largest
	largest := a.Largest(1)
	chk.Float64(tst, "Largest", 1e-17, largest, 50)

	// GetRow
	row0 := a.GetRow(0)
	row1 := a.GetRow(1)
	row2 := a.GetRow(2)
	chk.Array(tst, "GetRow(0)", 1e-17, row0, []float64{1, 2, 3, 4, 5})
	chk.Array(tst, "GetRow(1)", 1e-17, row1, []float64{0.1, 0.2, 0.3, 0.4, 0.5})
	chk.Array(tst, "GetRow(2)", 1e-17, row2, []float64{10, 20, 30, 40, 50})

	// GetCol
	col0 := a.GetCol(0)
	col2 := a.GetCol(2)
	col4 := a.GetCol(4)
	chk.Array(tst, "GetCol(0)", 1e-17, col0, []float64{1, 0.1, 10})
	chk.Array(tst, "GetCol(2)", 1e-17, col2, []float64{3, 0.3, 30})
	chk.Array(tst, "GetCol(4)", 1e-17, col4, []float64{5, 0.5, 50})

	// NormFrob
	A := NewMatrixSlice([][]float64{
		{-3, 5, 7},
		{+2, 6, 4},
		{+0, 2, 8},
	})
	normFA := A.NormFrob()
	chk.Float64(tst, "NormFrob", 1e-17, normFA, 1.438749456993816e+01)

	// NormInf
	normIA := A.NormInf()
	chk.Float64(tst, "NormInf", 1e-17, normIA, 15.0)
}

func TestMatrix04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Matrix04. (complex) matrix methods")

	// Set
	a := NewMatrixC(3, 5)
	a.Set(0, 0, 1+1i)
	a.Set(0, 1, 2+2i)
	a.Set(0, 2, 3+3i)
	a.Set(0, 3, 4+4i)
	a.Set(0, 4, 5+5i)
	a.Set(1, 0, 0.1)
	a.Set(1, 1, 0.2)
	a.Set(1, 2, 0.3)
	a.Set(1, 3, 0.4)
	a.Set(1, 4, 0.5)
	a.Set(2, 0, 10)
	a.Set(2, 1, 20)
	a.Set(2, 2, 30)
	a.Set(2, 3, 40)
	a.Set(2, 4, 50)
	chk.MatrixC(tst, "a", 1e-17, a.GetSlice(), [][]complex128{
		{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i, 5 + 5i},
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{10, 20, 30, 40, 50},
	})

	// GetCopy
	aclone := a.GetCopy()
	chk.ArrayC(tst, "aclone", 1e-17, a.Data, aclone.Data)

	// Fill
	b := NewMatrixC(5, 3)
	b.Fill(2 - 1i)
	chk.MatrixC(tst, "b", 1e-17, b.GetSlice(), [][]complex128{
		{2 - 1i, 2 - 1i, 2 - 1i},
		{2 - 1i, 2 - 1i, 2 - 1i},
		{2 - 1i, 2 - 1i, 2 - 1i},
		{2 - 1i, 2 - 1i, 2 - 1i},
		{2 - 1i, 2 - 1i, 2 - 1i},
	})

	// Scale
	c := NewMatrixC(5, 3)
	c.Fill(2 + 2i)
	c.Apply(1.0/4.0, c)
	chk.MatrixC(tst, "c := c/4", 1e-17, c.GetSlice(), [][]complex128{
		{0.5 + 0.5i, 0.5 + 0.5i, 0.5 + 0.5i},
		{0.5 + 0.5i, 0.5 + 0.5i, 0.5 + 0.5i},
		{0.5 + 0.5i, 0.5 + 0.5i, 0.5 + 0.5i},
		{0.5 + 0.5i, 0.5 + 0.5i, 0.5 + 0.5i},
		{0.5 + 0.5i, 0.5 + 0.5i, 0.5 + 0.5i},
	})
}
