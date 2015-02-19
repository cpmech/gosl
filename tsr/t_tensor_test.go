// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"fmt"
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_tsr01(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	chk.PrintTitle("tsr01")

	F := [][]float64{
		{1, 1.5, 0},
		{0, 1.5, 0},
		{0, 0.0, 1},
	}
	Fi := Alloc2()       // inverse of F
	C := Alloc2()        // right Cauchy-Green deformation tensor
	b := Alloc2()        // left Cauchy-Green deformation tensor
	E := Alloc2()        // Green strain
	e := Alloc2()        // Almansi strain
	epf := Alloc2()      // e as push-forward of E
	Epb := Alloc2()      // E as pull-back of e
	Ipf := Alloc2()      // I as push-forward of C
	Cpb := Alloc2()      // C as a pull-back of I
	Fi[0][0] = 666       // noise
	C[0][0] = 666        // noise
	b[0][0] = 666        // noise
	E[0][0] = 666        // noise
	e[0][0] = 666        // noise
	epf[0][0] = 666      // noise
	Epb[0][0] = 666      // noise
	Ipf[0][0] = 666      // noise
	Cpb[0][0] = 666      // noise
	J, err := Inv(Fi, F) // Fi  := inv(F)
	if err != nil {
		chk.Panic("%v", err)
	}
	RightCauchyGreenDef(C, F)  // C   := Ft * F
	LeftCauchyGreenDef(b, F)   // b   := F * Ft
	GreenStrain(E, F)          // E   := 0.5 * (Ft * F - I)
	AlmansiStrain(e, Fi)       // e   := 0.5 * (I - Fit * F)
	PushForward(epf, E, F, Fi) // epf := push-forward(E)
	PullBack(Epb, e, F, Fi)    // Epb := pull-back(e)
	PushForward(Ipf, C, F, Fi) // Ipf := push-forward(C)
	PullBack(Cpb, It, F, Fi)   // Cpb := pull-back(I)
	detC, detb := Det(C), Det(b)
	io.Pf("F   = %v\n", F)
	io.Pf("Fi  = %v\n", Fi)
	io.Pf("C   = %v\n", C)
	io.Pf("b   = %v\n", b)
	io.Pf("E   = %v\n", E)
	io.Pf("e   = %v\n", e)
	io.Pf("epf = %v\n", epf)
	io.Pf("Epb = %v\n", Epb)
	io.Pf("Ipf = %v\n", Ipf)
	io.Pf("Cpb = %v\n", Cpb)
	io.Pf("det(F)=%v, det(C)=%v, det(b)=%v\n", J, detC, detb)
	chk.Matrix(tst, "Fi", 1.0e-17, Fi, [][]float64{{1, -1, 0}, {0, 2.0 / 3.0, 0}, {0, 0, 1}})
	chk.Matrix(tst, "C", 1.0e-17, C, [][]float64{{1, 1.5, 0}, {1.5, 4.5, 0}, {0, 0, 1}})
	chk.Matrix(tst, "b", 1.0e-17, b, [][]float64{{3.25, 2.25, 0}, {2.25, 2.25, 0}, {0, 0, 1}})
	chk.Matrix(tst, "E", 1.0e-17, E, [][]float64{{0, 0.75, 0}, {0.75, 1.75, 0}, {0, 0, 0}})
	chk.Matrix(tst, "e", 1.0e-17, e, [][]float64{{0, 0.5, 0}, {0.5, -2.0 / 9.0, 0}, {0, 0, 0}})
	chk.Matrix(tst, "epf", 1.0e-15, epf, [][]float64{{0, 0.5, 0}, {0.5, -2.0 / 9.0, 0}, {0, 0, 0}})
	chk.Matrix(tst, "Epb", 1.0e-17, Epb, [][]float64{{0, 0.75, 0}, {0.75, 1.75, 0}, {0, 0, 0}})
	chk.Matrix(tst, "Ipf", 1.0e-17, Ipf, It)
	chk.Matrix(tst, "Cpb", 1.0e-17, Cpb, [][]float64{{1, 1.5, 0}, {1.5, 4.5, 0}, {0, 0, 1}})
	chk.Scalar(tst, "det(F)", 1.0e-17, J, 1.5)
	chk.Scalar(tst, "det(C)", 1.0e-17, detC, 1.5*1.5)
	chk.Scalar(tst, "det(b)", 1.0e-17, detb, 1.5*1.5)
}

func Test_tsr02(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	chk.PrintTitle("tsr02")

	F := [][]float64{
		{2, 8.0 / 3.0, 0},
		{0, 2, 0},
		{0, 0.0, 1},
	}
	Fi := Alloc2()
	C := Alloc2()
	b := Alloc2()
	J, err := Inv(Fi, F)
	if err != nil {
		chk.Panic("%v", err)
	}
	RightCauchyGreenDef(C, F)
	LeftCauchyGreenDef(b, F)
	io.Pf("F = %v\n", F)
	io.Pf("C = %v\n", C)
	io.Pf("b = %v\n", b)
	chk.Scalar(tst, "J", 1.0e-17, J, 4.0)
	chk.Matrix(tst, "C", 1.0e-17, C, [][]float64{{36.0 / 9.0, 48.0 / 9.0, 0}, {48.0 / 9.0, 100.0 / 9.0, 0}, {0, 0, 1}})
	chk.Matrix(tst, "b", 1.0e-17, b, [][]float64{{100.0 / 9.0, 48.0 / 9.0, 0}, {48.0 / 9.0, 36.0 / 9.0, 0}, {0, 0, 1}})

	λ, μ := 2.0, 3.0
	σ := Alloc2()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			σ[i][j] = (λ/J)*math.Log(J)*It[i][j] + (μ/J)*(b[i][j]-It[i][j])
		}
	}

	P := Alloc2()
	S := Alloc2()
	σfromP := Alloc2()
	σfromS := Alloc2()
	CauchyToPK1(P, σ, F, Fi, J)
	CauchyToPK2(S, σ, F, Fi, J)
	PK1ToCauchy(σfromP, P, F, Fi, J)
	PK2ToCauchy(σfromS, S, F, Fi, J)

	io.Pf("σ = %v\n", σ)
	io.Pf("P = %v\n", P)
	io.Pf("S = %v\n", S)
	chk.Matrix(tst, "σfromP", 1.0e-17, σfromP, σ)
	chk.Matrix(tst, "σfromS", 1.0e-14, σfromS, σ)
}

func Test_tsr03(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	chk.PrintTitle("tsr03")

	a := [][]float64{
		{4.0, 1.0 / SQ2, 0},
		{1.0 / SQ2, 5.0, 0},
		{0, 0, 6.0},
	}
	am := make([]float64, 4)
	aa := Alloc2()
	Ten2Man(am, a)
	Man2Ten(aa, am)
	io.Pf("a  = %v\n", a)
	io.Pf("am = %v\n", am)
	io.Pf("aa = %v\n", aa)
	chk.Matrix(tst, "aa", 1.0e-15, aa, a)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			chk.Scalar(tst, fmt.Sprintf("am[%d][%d]", i, j), 1.0e-15, M2T(am, i, j), a[i][j])
		}
	}

	b := [][]float64{
		{4.0, 1.0 / SQ2, 3.0 / SQ2},
		{1.0 / SQ2, 5.0, 2.0 / SQ2},
		{3.0 / SQ2, 2.0 / SQ2, 6.0},
	}
	bm := make([]float64, 6)
	bb := Alloc2()
	Ten2Man(bm, b)
	Man2Ten(bb, bm)
	io.Pf("b  = %v\n", b)
	io.Pf("bm = %v\n", bm)
	io.Pf("bb = %v\n", bb)
	chk.Matrix(tst, "bb", 1.0e-15, bb, b)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			chk.Scalar(tst, fmt.Sprintf("bm[%d][%d]", i, j), 1.0e-15, M2T(bm, i, j), b[i][j])
		}
	}
}

func Test_tsr04(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	chk.PrintTitle("tsr04")

	a := [][]float64{
		{8.0, 1.0 / SQ2, 2.0 / SQ2},
		{1.0 / SQ2, -5.0, 0.5 / SQ2},
		{2.0 / SQ2, 0.5 / SQ2, 7.0},
	}

	am := make([]float64, 6)
	Ten2Man(am, a)
	chk.Vector(tst, "Ten2Man", 1e-17, am, []float64{8, -5, 7, 1, 0.5, 2})

	amdyam := make([][]float64, 6)
	for i := 0; i < 6; i++ {
		amdyam[i] = make([]float64, 6)
		for j := 0; j < 6; j++ {
			amdyam[i][j] = am[i] * am[j]
		}
	}

	adya := Alloc4()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				for l := 0; l < 3; l++ {
					adya[i][j][k][l] = a[i][j] * a[k][l]
				}
			}
		}
	}

	//io.Pforan("adya = %v\n", adya)
	//io.Pforan("amdyam = %v\n", amdyam)

	var err float64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				for l := 0; l < 3; l++ {
					err += math.Abs(adya[i][j][k][l] - M2TT(amdyam, i, j, k, l))
				}
			}
		}
	}
	if err > 1e-13 {
		chk.Panic("M2TT failed")
	}
}
