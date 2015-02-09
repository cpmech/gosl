// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"
	"testing"

	"code.google.com/p/gosl/utl"
)

func TestVector01(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()
	utl.TTitle("TestVector 01")

	utl.Pf("func VecFill(v []float64, s float64)\n")
	v := make([]float64, 5)
	VecFill(v, 666)
	PrintVec("v", v, "%5g", false)
	utl.CheckVector(tst, "v", 1e-17, v, []float64{666, 666, 666, 666, 666})

	utl.Pf("\nfunc VecFillC(v []complex128, s complex128)\n")
	vc := make([]complex128, 5)
	VecFillC(vc, 666+666i)
	PrintVecC("vc", vc, "(%2g +", "%4gi) ", false)
	utl.CheckVectorC(tst, "vc", 1e-17, vc, []complex128{666 + 666i, 666 + 666i, 666 + 666i, 666 + 666i, 666 + 666i})

	utl.Pf("\nfunc VecNorm(v []float64) (nrm float64)\n")
	PrintVec("v", v, "%5g", false)
	nrm := VecNorm(v)
	utl.Pf("norm(v) = %23.15e\n", nrm)
	utl.CheckScalar(tst, "norm(v)", 1e-17, nrm, 1.489221273014860e+03)

	utl.Pf("\nfunc VecNormDiff(u, v []float64) (nrm float64)\n")
	u := []float64{333, 333, 333, 333, 333}
	PrintVec("u", u, "%5g", false)
	PrintVec("v", v, "%5g", false)
	nrm = VecNormDiff(u, v)
	utl.Pf("norm(u-v) = %23.15e\n", nrm)
	utl.CheckScalar(tst, "norm(u-v)", 1e-17, nrm, math.Sqrt(5.0*333.0*333.0))

	utl.Pf("\nfunc VecDot(u, v []float64) (res float64)\n")
	u = []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	PrintVec("u", u, "%5g", false)
	PrintVec("v", v, "%5g", false)
	udotv := VecDot(u, v)
	utl.Pf("u dot v = %v\n", udotv)
	utl.CheckScalar(tst, "u dot v", 1e-12, udotv, 999)

	utl.Pf("\nfunc VecCopy(a []float64, alp float64, b []float64)\n")
	a := make([]float64, len(u))
	VecCopy(a, 1, u)
	PrintVec("u     ", u, "%5g", false)
	PrintVec("a := u", a, "%5g", false)
	utl.CheckVector(tst, "a", 1e-17, a, []float64{0.1, 0.2, 0.3, 0.4, 0.5})

	utl.Pf("\nfunc VecAdd(a []float64, alp float64, b []float64)\n")
	b := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	PrintVec("b        ", b, "%5g", false)
	VecAdd(b, 10, b) // b += 10.0*b
	PrintVec("b += 10*b", b, "%5g", false)
	utl.CheckVector(tst, "b", 1e-17, b, []float64{11, 22, 33, 44, 55})

	utl.Pf("\nfunc VecAdd2(u []float64, alp float64, a []float64, bet float64, b []float64)\n")
	PrintVec("a", a, "%7g", false)
	PrintVec("b", b, "%7g", false)
	c := make([]float64, len(a))
	VecAdd2(c, 1, a, 10, b) // c = 1.0*a + 10.0*b
	PrintVec("c = 1*a+10*b", c, "%7g", false)
	utl.CheckVector(tst, "c", 1e-17, c, []float64{110.1, 220.2, 330.3, 440.4, 550.5})

	utl.Pf("\nfunc VecMin(v []float64) (min float64)\n")
	PrintVec("a", a, "%5g", false)
	mina := VecMin(a)
	utl.Pf("min(a) = %v\n", mina)
	utl.CheckScalar(tst, "min(a)", 1e-17, mina, 0.1)

	utl.Pf("\nfunc VecMax(v []float64) (max float64)\n")
	PrintVec("a", a, "%5g", false)
	maxa := VecMax(a)
	utl.Pf("max(a) = %v\n", maxa)
	utl.CheckScalar(tst, "max(a)", 1e-17, maxa, 0.5)

	utl.Pf("\nfunc VecMinMax(v []float64) (min, max float64)\n")
	PrintVec("a", a, "%5g", false)
	min2a, max2a := VecMinMax(a)
	utl.Pf("min(a) = %v\n", min2a)
	utl.Pf("max(a) = %v\n", max2a)
	utl.CheckScalar(tst, "min(a)", 1e-17, min2a, 0.1)
	utl.CheckScalar(tst, "max(a)", 1e-17, max2a, 0.5)

	utl.Pf("\nfunc VecLargest(u []float64, den float64) (largest float64)\n")
	PrintVec("b     ", b, "%5g", false)
	bdiv11 := []float64{b[0] / 11.0, b[1] / 11.0, b[2] / 11.0, b[3] / 11.0, b[4] / 11.0}
	PrintVec("b / 11", bdiv11, "%5g", false)
	maxbdiv11 := VecLargest(b, 11)
	utl.Pf("max(b/11) = %v\n", maxbdiv11)
	utl.CheckScalar(tst, "max(b/11)", 1e-17, maxbdiv11, 5)

	utl.Pf("\nfunc VecMaxDiff(a, b []float64) (maxdiff float64)\n")
	amb1 := []float64{a[0] - b[0], a[1] - b[1], a[2] - b[2], a[3] - b[3], a[4] - b[4]}
	amb2 := make([]float64, len(a))
	VecAdd2(amb2, 1, a, -1, b)
	PrintVec("a  ", a, "%7g", false)
	PrintVec("b  ", b, "%7g", false)
	PrintVec("a-b", amb1, "%7g", false)
	PrintVec("a-b", amb2, "%7g", false)
	maxdiffab := VecMaxDiff(a, b)
	utl.Pf("maxdiff(a,b) = max(abs(a-b)) = %v\n", maxdiffab)
	utl.CheckVector(tst, "amb1 == amb2", 1e-17, amb1, amb2)
	utl.CheckScalar(tst, "maxdiff(a,b)", 1e-17, maxdiffab, 54.5)

	utl.Pf("\nfunc VecMaxDiffC(a, b []complex128) (maxdiff float64)\n")
	az := []complex128{complex(a[0], 1), complex(a[1], 3), complex(a[2], 0.5), complex(a[3], 1), complex(a[4], 0)}
	bz := []complex128{complex(b[0], 1), complex(b[1], 6), complex(b[2], 0.8), complex(b[3], -3), complex(b[4], 1)}
	ambz := []complex128{az[0] - bz[0], az[1] - bz[1], az[2] - bz[2], az[3] - bz[3], az[4] - bz[4]}
	PrintVecC("az   ", az, "(%5g +", "%4gi) ", false)
	PrintVecC("bz   ", bz, "(%5g +", "%4gi) ", false)
	PrintVecC("az-bz", ambz, "(%5g +", "%4gi) ", false)
	maxdiffabz := VecMaxDiffC(az, bz)
	utl.Pf("maxdiff(az,bz) = %v\n", maxdiffabz)
	utl.CheckScalar(tst, "maxdiff(az,bz)", 1e-17, maxdiffabz, 54.5)

	utl.Pf("\nfunc VecScale(res []float64, Atol, Rtol float64, v []float64)\n")
	scal1 := make([]float64, len(a))
	VecScale(scal1, 0.5, 0.1, amb1)
	PrintVec("a-b            ", amb1, "%7g", false)
	PrintVec("0.5 + 0.1*(a-b)", scal1, "%7g", false)
	utl.CheckVector(tst, "0.5 + 0.1*(a-b)", 1e-15, scal1, []float64{-0.59, -1.68, -2.77, -3.86, -4.95})

	utl.Pf("\nfunc VecScaleAbs(res []float64, Atol, Rtol float64, v []float64)\n")
	scal2 := make([]float64, len(a))
	VecScaleAbs(scal2, 0.5, 0.1, amb1)
	PrintVec("a-b            ", amb1, "%7g", false)
	PrintVec("0.5 + 0.1*|a-b|", scal2, "%7g", false)
	utl.CheckVector(tst, "0.5 + 0.1*|a-b|", 1e-15, scal2, []float64{1.59, 2.68, 3.77, 4.86, 5.95})

	utl.Pf("\nfunc VecRms(u []float64) (rms float64)\n")
	PrintVec("v", v, "%5g", false)
	rms := VecRms(v)
	utl.Pf("rms(v) = %23.15e\n", rms)
	utl.CheckScalar(tst, "rms(v)", 1e-17, rms, 666.0)

	utl.Pf("func VecRmsErr(u []float64, Atol, Rtol float64, v []float64) (rms float64)\n")
	PrintVec("v", v, "%5g", false)
	rmserr := VecRmsErr(v, 0, 1, v)
	utl.Pf("rmserr(v,v) = %23.15e\n", rmserr)
	utl.CheckScalar(tst, "rmserr(v,v,0,1)", 1e-17, rmserr, 1)

	utl.Pf("func VecRmsError(u, w []float64, Atol, Rtol float64, v []float64) (rms float64)\n")
	PrintVec("v", v, "%5g", false)
	w := []float64{333, 333, 333, 333, 333}
	rmserr = VecRmsError(v, w, 0, 1, v)
	utl.Pf("rmserr(v,w,v) = %23.15e\n", rmserr)
	utl.CheckScalar(tst, "rmserr(v,w,0,1,v)", 1e-17, rmserr, 0.5)
}
