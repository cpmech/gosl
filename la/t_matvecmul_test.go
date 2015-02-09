// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"code.google.com/p/gosl/utl"
)

func Test_mvMul01(tst *testing.T) {

	tsprev := utl.Tsilent
	defer func() {
		utl.Tsilent = tsprev
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("mvMul01. MatrixVector multiplication")

	a := [][]float64{
		{10.0, 20.0, 30.0, 40.0, 50.0},
		{1.0, 20.0, 3.0, 40.0, 5.0},
		{10.0, 2.0, 30.0, 4.0, 50.0},
	}
	u := []float64{0.5, 0.4, 0.3, 0.2, 0.1}
	r := []float64{0.5, 0.4, 0.3}
	z := []float64{1000, 1000, 1000}
	w := []float64{1000, 2000, 3000, 4000, 5000}
	au := make([]float64, 3)
	atr := make([]float64, 5)
	au_cor := []float64{35, 17.9, 20.6}
	atr_cor := []float64{8.4, 18.6, 25.2, 37.2, 42}
	zpau_cor := []float64{1035, 1017.9, 1020.6}
	wpar_cor := []float64{1008.4, 2018.6, 3025.2, 4037.2, 5042}
	MatVecMul(au, 1, a, u)     // au  = 1*a*u
	MatTrVecMul(atr, 1, a, r)  // atr = 1*transp(a)*r
	MatVecMulAdd(z, 1, a, u)   // z  += 1*a*u
	MatTrVecMulAdd(w, 1, a, r) // w  += 1*transp(a)*r
	utl.CheckVector(tst, "au", 1.0e-17, au, au_cor)
	utl.CheckVector(tst, "atr", 1.0e-17, atr, atr_cor)
	utl.CheckVector(tst, "zpau", 1.0e-12, z, zpau_cor)
	utl.CheckVector(tst, "wpar", 1.0e-12, w, wpar_cor)
}

func Test_mmMul01(tst *testing.T) {

	tsprev := utl.Tsilent
	defer func() {
		utl.Tsilent = tsprev
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("mmMul01. MatrixMatrix multiplication")

	a := [][]float64{
		{1.0, 2.0, 3.0},
		{0.5, 0.75, 1.5},
	}
	b := [][]float64{
		{0.1, 0.5, 0.5, 0.75},
		{0.2, 2.0, 2.0, 2.0},
		{0.3, 0.5, 0.5, 0.5},
	}
	c := MatAlloc(2, 4)
	MatMul(c, 1, a, b) // c := 1*a*b
	utl.Pf("a = %v\n", a)
	utl.Pf("b = %v\n", b)
	utl.Pf("c = %v\n", c)

	ccor := [][]float64{
		{1.4, 6.0, 6.0, 6.25},
		{0.65, 2.5, 2.5, 2.625},
	}
	utl.CheckMatrix(tst, "c", 1.0e-15, c, ccor)
}

func Test_mmMul02(tst *testing.T) {

	tsprev := utl.Tsilent
	defer func() {
		utl.Tsilent = tsprev
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("mmMul02. MatrixMatrix multiplication")

	a := [][]float64{
		{0.1, 0.2, 0.3},
		{1.0, 0.2, 0.3},
		{2.0, 0.2, 0.3},
		{3.0, 0.2, 0.3},
	}

	b := [][]float64{
		{10.0, 20.0, 30.0, 40.0, 50.0},
		{10.0, 20.0, 30.0, 40.0, 50.0},
		{10.0, 20.0, 30.0, 40.0, 50.0},
	}

	c := [][]float64{
		{0.1, 0.2, 0.3, 0.4},
		{0.1, 0.2, 0.3, 0.4},
		{0.1, 0.2, 0.3, 0.4},
		{0.1, 0.2, 0.3, 0.4},
		{0.1, 0.2, 0.3, 0.4},
	}

	e := [][]float64{
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{1.0, 0.2, 0.3, 0.4, 0.5},
		{2.0, 0.2, 0.3, 0.4, 0.5},
	}

	ab := MatAlloc(len(a), len(b[0]))
	abc := MatAlloc(len(a), len(c[0]))
	etbc := MatAlloc(len(e[0]), len(c[0]))

	detbc := [][]float64{
		{1000.0, 1000.0, 1000.0, 1000.0},
		{1000.0, 1000.0, 1000.0, 1000.0},
		{1000.0, 1000.0, 1000.0, 1000.0},
		{1000.0, 1000.0, 1000.0, 1000.0},
		{1000.0, 1000.0, 1000.0, 1000.0},
	}

	MatMul(ab, 0.5, a, b)
	MatMul3(abc, 0.5, a, b, c)
	MatTrMul3(etbc, 0.5, e, b, c)
	MatTrMulAdd3(detbc, 0.5, e, b, c)
	PrintMat("a", a, "%g ", false)
	PrintMat("b", b, "%g ", false)
	PrintMat("c", c, "%g ", false)
	PrintMat("ab", ab, "%g ", false)
	PrintMat("abc", abc, "%g ", false)
	PrintMat("etbc", etbc, "%g ", false)
	PrintMat("detbc", detbc, "%g ", false)

	ab_cor := [][]float64{
		{3.0, 6, 9.0, 12, 15.0},
		{7.5, 15, 22.5, 30, 37.5},
		{12.5, 25, 37.5, 50, 62.5},
		{17.5, 35, 52.5, 70, 87.5},
	}
	abc_cor := [][]float64{
		{4.5, 9.0, 13.5, 18},
		{11.25, 22.5, 33.75, 45},
		{18.75, 37.5, 56.25, 75},
		{26.25, 52.5, 78.75, 105},
	}
	etbc_cor := [][]float64{
		{23.25, 46.5, 69.75, 93},
		{4.5, 9.0, 13.5, 18},
		{6.75, 13.5, 20.25, 27},
		{9.0, 18.0, 27.0, 36},
		{11.25, 22.5, 33.75, 45},
	}
	detbc_cor := [][]float64{
		{1023.25, 1046.5, 1069.75, 1093},
		{1004.5, 1009.0, 1013.5, 1018},
		{1006.75, 1013.5, 1020.25, 1027},
		{1009.0, 1018.0, 1027.0, 1036},
		{1011.25, 1022.5, 1033.75, 1045},
	}

	utl.CheckMatrix(tst, "ab", 1.0e-17, ab, ab_cor)
	utl.CheckMatrix(tst, "abc", 1.0e-17, abc, abc_cor)
	utl.CheckMatrix(tst, "etbc", 1.0e-13, etbc, etbc_cor)
	utl.CheckMatrix(tst, "detbc", 1.0e-13, detbc, detbc_cor)
}
