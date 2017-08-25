// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oblas

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestSplitJoinComplex(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SplitJoinComplex")

	v := []complex128{1 + 0.1i, 2 + 0.2i, 3 - 0.3i}
	vr, vi := GetSplitComplex(v)
	chk.Array(tst, "vr", 1e-17, vr, []float64{1, 2, 3})
	chk.Array(tst, "vi", 1e-17, vi, []float64{0.1, 0.2, -0.3})

	u := GetJoinComplex(vr, vi)
	chk.ArrayC(tst, "u=v", 1e-17, u, v)

	wr, wi := make([]float64, len(v)), make([]float64, len(v))
	SplitComplex(wr, wi, v)
	chk.Array(tst, "wr", 1e-17, wr, []float64{1, 2, 3})
	chk.Array(tst, "wi", 1e-17, wi, []float64{0.1, 0.2, -0.3})

	w := make([]complex128, len(v))
	JoinComplex(w, wr, wi)
	chk.ArrayC(tst, "w=v", 1e-17, w, v)
}

func TestConversions01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Conversions01. real")

	A := [][]float64{
		{1, 2, +3, +4},
		{5, 6, +7, +8},
		{9, 0, -1, -2},
	}
	m, n := len(A), len(A[0])

	a := SliceToColMajor(A)

	row1 := ExtractRow(1, m, n, a)
	io.Pf("row1 = %v\n", row1)
	chk.Array(tst, "row1", 1e-17, row1, []float64{5, 6, 7, 8})

	col2 := ExtractCol(2, m, n, a)
	io.Pf("col2 = %v\n", col2)
	chk.Array(tst, "col2", 1e-17, col2, []float64{3, 7, -1})
}

func TestConversions02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Conversions02. complex")

	A := [][]complex128{
		{1 + 0.1i, 2, +3, +4 - 0.4i},
		{5 + 0.5i, 6, +7, +8 - 0.8i},
		{9 + 0.9i, 0, -1, -2 + 1.0i},
	}
	m, n := len(A), len(A[0])

	a := SliceToColMajorC(A)

	row1 := ExtractRowC(1, m, n, a)
	io.Pf("row1 = %v\n", row1)
	chk.ArrayC(tst, "row1", 1e-17, row1, []complex128{5 + 0.5i, 6, 7, 8 - 0.8i})

	col2 := ExtractColC(2, m, n, a)
	io.Pf("col2 = %v\n", col2)
	chk.ArrayC(tst, "col2", 1e-17, col2, []complex128{3, 7, -1})
}
