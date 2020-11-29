// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestEigen01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eigen01")

	A := NewMatrixDeep2([][]float64{
		{2, 0, 0},
		{0, 2, 0},
		{0, 0, 2},
	})

	w := NewVectorC(A.M)
	EigenVal(w, A, true)
	chk.ArrayC(tst, "w", 1e-17, w, []complex128{2, 2, 2})

	v := NewMatrixC(A.M, A.M)
	EigenVecR(v, w, A, true)

	io.Pf("v = \n")
	io.Pf("%v\n", v.Print("", ""))

	chk.Deep2c(tst, "v", 1e-17, v.GetDeep2(), [][]complex128{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	})

	CheckEigenVecR(tst, A, w, v, 1e-17)
}

func TestEigen02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eigen02")

	A := NewMatrixDeep2([][]float64{
		{2, 0, 0},
		{0, 3, 4},
		{0, 4, 9},
	})

	w := NewVectorC(A.M)
	EigenVal(w, A, true)
	io.Pforan("w = %v\n", w)
	chk.ArrayC(tst, "w", 1e-17, w, []complex128{11, 1, 2})

	v := NewMatrixC(A.M, A.M)
	EigenVecR(v, w, A, true)

	io.Pf("v = \n")
	io.Pf("%v\n", v.Print("%12.8f", "%12.8f"))

	os3 := complex(1.0/math.Sqrt(5.0), 0)
	io.Pforan("os3 = %v\n", os3)
	io.Pforan("2*os3 = %v\n", 2*os3)
	chk.Deep2c(tst, "v", 1e-17, v.GetDeep2(), [][]complex128{
		{0 * os3, +0 * os3, 1},
		{1 * os3, +2 * os3, 0},
		{2 * os3, -1 * os3, 0},
	})

	CheckEigenVecR(tst, A, w, v, 1e-15)
}

func TestEigen03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eigen03")

	A := NewMatrixDeep2([][]float64{
		{1, 2, 3},
		{2, 3, 2},
		{3, 2, 2},
	})

	w := NewVectorC(A.M)
	EigenVal(w, A, true)
	io.Pforan("w = %v\n", w)
	chk.ArrayC(tst, "w", 1e-14, w, []complex128{6.69537390404459476e+00, -1.55809924785903786e+00, 8.62725343814443657e-01})

	v := NewMatrixC(A.M, A.M)
	EigenVecR(v, w, A, true)

	io.Pf("v = \n")
	io.Pf("%v\n", v.Print("%12.8f", "%12.8f"))

	chk.Deep2c(tst, "v", 1e-15, v.GetDeep2(), [][]complex128{
		{-5.26633230856907386e-01, -7.81993314738381295e-01, +3.33382506832158143e-01},
		{-6.07084171793832561e-01, +7.14394870018381645e-02, -7.91419742017035133e-01},
		{-5.95068272145819699e-01, +6.19179178753124115e-01, +5.12358171676802088e-01},
	})

	CheckEigenVecR(tst, A, w, v, 1e-14)
}

func TestEigen04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eigen04")

	A := NewMatrixDeep2([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	w := NewVectorC(A.M)
	v := NewMatrixC(A.M, A.M)
	EigenVecR(v, w, A, true)

	io.Pforan("w = %v\n", w)
	chk.ArrayC(tst, "w", 1.43e-14, w, []complex128{16.116843969807043, -1.116843969807043, 0.0})

	io.Pf("v = \n")
	io.Pf("%v\n", v.Print("%12.8f", "%12.8f"))

	chk.Deep2c(tst, "v", 1e-15, v.GetDeep2(), [][]complex128{
		{-0.231970687246286, -0.785830238742067, +0.408248290463864},
		{-0.525322093301234, -0.086751339256628, -0.816496580927726},
		{-0.818673499356181, +0.612327560228810, +0.408248290463863},
	})

	CheckEigenVecR(tst, A, w, v, 1e-14)
}

func TestEigen05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eigen05. Rutis3 matrix")

	A := NewMatrixDeep2([][]float64{
		{+4.0, +0.0, +5.0, 3.0},
		{-5.0, +4.0, -3.0, 0.0},
		{+0.0, -3.0, +4.0, 5.0},
		{+3.0, -5.0, +0.0, 4.0},
	})

	u := NewMatrixC(A.M, A.M)
	v := NewMatrixC(A.M, A.M)
	w := NewVectorC(A.M)
	EigenVecLR(u, v, w, A, true)

	io.Pforan("w = %v\n", w)
	chk.ArrayC(tst, "w", 1.6e-14, w, []complex128{12.0, 1.0 + 5.0i, 1.0 - 5.0i, 2.0})

	io.Pf("u = \n")
	io.Pf("%v\n", u.Print("%10.6f", "%10.6f"))

	uRef := [][]complex128{
		{+0.5, -0.5 + 0.0i, -0.5 - 0.0i, +0.5},
		{-0.5, +0.0 - 0.5i, +0.0 + 0.5i, +0.5},
		{+0.5, +0.0 - 0.5i, +0.0 + 0.5i, -0.5},
		{+0.5, +0.5 + 0.0i, +0.5 + 0.0i, +0.5},
	}
	chk.Deep2c(tst, "u", 1e-15, u.GetDeep2(), uRef)

	io.Pl()
	io.Pf("v = \n")
	io.Pf("%v\n", v.Print("%10.6f", "%10.6f"))

	vRef := [][]complex128{
		{+0.5, +0.0 + 0.5i, +0.0 - 0.5i, +0.5},
		{-0.5, -0.5 + 0.0i, -0.5 + 0.0i, +0.5},
		{+0.5, -0.5 + 0.0i, -0.5 - 0.0i, -0.5},
		{+0.5, -0.0 - 0.5i, -0.0 + 0.5i, +0.5},
	}
	chk.Deep2c(tst, "v", 1e-15, v.GetDeep2(), vRef)

	CheckEigenVecL(tst, A, w, u, 1.16e-14)
	CheckEigenVecR(tst, A, w, v, 1.16e-14)

	// compute left eigenvector again
	u2 := NewMatrixC(A.M, A.M)
	w2 := NewVectorC(A.M)
	EigenVecL(u2, w2, A, true)
	chk.Deep2c(tst, "u2", 1e-15, u2.GetDeep2(), uRef)

	// compute right eigenvector again
	v3 := NewMatrixC(A.M, A.M)
	w3 := NewVectorC(A.M)
	EigenVecR(v3, w3, A, true)
	chk.Deep2c(tst, "v3", 1e-15, v3.GetDeep2(), vRef)
}
