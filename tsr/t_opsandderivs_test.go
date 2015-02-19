// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/num"
)

var (
	//              0  1  2  3  4  5  6  7  8  9 10 11 12
	test_nd = []int{2, 2, 3, 3, 3, 2, 2, 3, 3, 3, 3, 3, 3}
	test_AA = [][][]float64{
		{
			{1, 2, 0}, // 0
			{2, -2, 0},
			{0, 0, -2},
		},
		{
			{-100, 33, 0}, // 1
			{33, -200, 0},
			{0, 0, 150},
		},
		{
			{1, 2, 4}, // 2
			{2, -2, 3},
			{4, 3, -2},
		},
		{
			{-100, -10, 20}, // 3
			{-10, -200, 15},
			{20, 15, -300},
		},
		{
			{-100, 0, -10}, // 4
			{0, -200, 0},
			{-10, 0, 100},
		},
		{
			{0.13, 1.2, 0}, // 5
			{1.2, -20, 0},
			{0, 0, -28},
		},
		{
			{-10, 3.3, 0}, // 6
			{3.3, -2, 0},
			{0, 0, 1.5},
		},
		{
			{0.1, 0.2, 0.8}, // 7
			{0.2, -1.3, 0.3},
			{0.8, 0.3, -0.2},
		},
		{
			{-10, -1, 2}, // 8
			{-1, -20, 1},
			{2, 1, -30},
		},
		{
			{-10, 0, -1}, // 9
			{0, -20, 0},
			{-1, 0, 10},
		},
		{
			{0, 0, 0}, // 10
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{1, 0, 0}, // 11
			{0, 1, 0},
			{0, 0, 0},
		},
		{
			{1, 0, 0}, // 12
			{0, 1, 0},
			{0, 0, 1},
		},
	}
	test_λ = [][]float64{
		{2, -3, -2}, // 0
		{-9.009173679699937e+01, -2.099082632030006e+02, 1.500000000000000e+02},  // 1
		{5.340450317010319e+00, -2.810439391134889e+00, -5.530010925875433e+00},  // 2
		{-9.731191510676280e+01, -1.982745793278015e+02, -3.044135055654359e+02}, // 3 swap(1,2)
		{-1.004987562112089e+02, -2.000000000000000e+02, 1.004987562112089e+02},  // 4 swap(1,2)
		{2.012826026112755e-01, -2.007128260261128e+01, -2.800000000000000e+01},  // 5
		{-1.118555686498567e+01, -8.144431350143311e-01, 1.500000000000000e+00},  // 6
		{8.204027629949283e-01, -1.376514253428709e+00, -8.438885095662192e-01},  // 7 swap(1,2)
		{-9.723289905003320e+00, -1.996226364398301e+01, -3.031444645101367e+01}, // 8 swap(1,2)
		{-1.004987562112089e+01, -2.000000000000000e+01, 1.004987562112089e+01},  // 9 swap(1,2)
		{0, 0, 0},
		{1, 1, 0},
		{1, 1, 1},
	}
	test_Q = [][][]float64{
		{
			{8.944271909999159e-01, -4.472135954999579e-01, 0.000000000000000e+00}, // 0
			{4.472135954999579e-01, 8.944271909999159e-01, 0.000000000000000e+00},
			{0.000000000000000e+00, 0.000000000000000e+00, 1.000000000000000e+00},
		},
		{
			{9.577602535503666e-01, -2.875678993196866e-01, 0.000000000000000e+00}, // 1
			{2.875678993196866e-01, 9.577602535503666e-01, 0.000000000000000e+00},
			{0.000000000000000e+00, 0.000000000000000e+00, 1.000000000000000e+00},
		},
		{
			{7.117223868913104e-01, -6.134981578280129e-01, -3.421567686592956e-01}, // 2
			{4.230940036436777e-01, 7.632060732740258e-01, -4.883727611143006e-01},
			{5.607519131295937e-01, 2.028213501071249e-01, 8.027582399840089e-01},
		},
		{
			{9.922973478006673e-01, 6.795178308599809e-02, -1.035786113210405e-01}, // 3 swap(1,2)
			{-8.322931230548185e-02, 9.850011890023485e-01, -1.511474089659375e-01},
			{9.175431935837298e-02, 1.586039496336501e-01, 9.830696476037130e-01},
		},
		{
			{9.987585269247989e-01, 0.000000000000000e+00, -4.981370188015976e-02}, // 4 swap(1,2)
			{0.000000000000000e+00, 1.000000000000000e+00, 0.000000000000000e+00},
			{4.981370188015976e-02, 0.000000000000000e+00, 9.987585269247989e-01},
		},
		{
			{9.982403466593390e-01, -5.929764161788244e-02, 0.000000000000000e+00}, // 5
			{5.929764161788244e-02, 9.982403466593390e-01, 0.000000000000000e+00},
			{0.000000000000000e+00, 0.000000000000000e+00, 1.000000000000000e+00},
		},
		{
			{9.411092600065299e-01, 3.381025890613104e-01, 0.000000000000000e+00}, // 6
			{-3.381025890613104e-01, 9.411092600065299e-01, 0.000000000000000e+00},
			{0.000000000000000e+00, 0.000000000000000e+00, 1.000000000000000e+00},
		},
		{
			{7.530915723725707e-01, 4.145581192313702e-03, -6.579026506847245e-01}, // 7 swap(1,2)
			{1.612764587047917e-01, 9.683071396355756e-01, 1.907123152772568e-01},
			{6.378424472305450e-01, -2.497280470578905e-01, 7.285553616737726e-01},
		},
		{
			{9.917737775527090e-01, 7.655540877078122e-02, -1.025867610696840e-01}, // 8 swap(1,2)
			{-8.740740608531795e-02, 9.905342925071262e-01, -1.058383707775363e-01},
			{9.351320505824934e-02, 1.139345634798570e-01, 9.890774467777477e-01},
		},
		{
			{9.987585269247989e-01, 0.000000000000000e+00, -4.981370188015975e-02}, // 9 swap(1,2)
			{0.000000000000000e+00, 1.000000000000000e+00, 0.000000000000000e+00},
			{4.981370188015975e-02, 0.000000000000000e+00, 9.987585269247989e-01},
		},
		{
			{1, 0, 0}, // 10
			{0, 1, 0},
			{0, 0, 1},
		},
		{
			{1, 0, 0}, // 11
			{0, 1, 0},
			{0, 0, 1},
		},
		{
			{1, 0, 0}, // 12
			{0, 1, 0},
			{0, 0, 1},
		},
	}
	//                  0    1   2  3    4    5  6    7  8   9  10 11 12
	test_σc = []float64{4, 300, 10, 0, 200, 1.0, 3, 1.6, 0, 20, 0, 0, 0}
)

func Test_Ts(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	chk.PrintTitle("Ts")

	nd := test_nd
	for m := 0; m < len(test_nd)-3; m++ {
		//for m := 0; m < 3; m++ {
		A := test_AA[m]
		a := M_Alloc2(nd[m])
		Ten2Man(a, A)
		s := M_Dev(a)
		Ts := M_Alloc4(nd[m])
		M_Ts(Ts, s)
		s2 := M_Alloc2(nd[m])
		ds2ds := M_Alloc4(nd[m])
		M_Sq(s2, s)
		M_SqDeriv(ds2ds, s)
		Ts_ := M_Alloc4(nd[m])
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a); j++ {
				for k := 0; k < len(a); k++ {
					for l := 0; l < len(a); l++ {
						Ts_[i][j] += Psd[i][k] * ds2ds[k][l] * Psd[l][j]
					}
				}
			}
		}
		chk.Matrix(tst, "Ts", 1e-13, Ts, Ts_)
	}
}

func Test_ops01(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	chk.PrintTitle("ops01")

	// basic derivatives
	dver := false
	dtol := 1e-5

	// invariants derivatives
	dveri := false
	dtoli1 := []float64{1e-6, 1e-6, 1e-6, 1e-6, 1e-6, 1e-6, 1e-6, 1e-6, 1e-6, 1e-6}
	dtoli2 := []float64{1e-6, 1e-5, 1e-6, 1e-4, 1e-5, 1e-6, 1e-6, 1e-6, 1e-6, 1e-6}
	dtoli3 := []float64{1e-6, 1e-3, 1e-6, 1e-3, 1e-3, 1e-6, 1e-6, 1e-6, 1e-5, 1e-6}

	// lode derivatives
	dverw := true
	dtolw := 1e-8

	nd := test_nd
	for m := 0; m < len(test_nd)-3; m++ {
		//for m := 0; m < 3; m++ {
		A := test_AA[m]
		a := M_Alloc2(nd[m])
		Ten2Man(a, A)
		trA := Tr(A)
		tra := M_Tr(a)
		detA := Det(A)
		deta := M_Det(a)
		devA := Dev(A)
		deva := M_Dev(a)
		devA_ := Alloc2()
		a2 := M_Alloc2(nd[m])
		A2 := Alloc2()
		A2_ := Alloc2()
		trDevA := Tr(devA)
		deva__ := M_Alloc2(nd[m])
		devA__ := Alloc2()
		s2 := M_Alloc2(nd[m])
		M_Sq(a2, a)
		M_Sq(s2, deva)
		Man2Ten(A2, a2)
		Man2Ten(devA_, deva)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				for k := 0; k < 3; k++ {
					A2_[i][j] += A[i][k] * A[k][j]
				}
			}
		}
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a); j++ {
				deva__[i] += Psd[i][j] * a[j]
			}
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				for k := 0; k < 3; k++ {
					for l := 0; l < 3; l++ {
						devA__[i][j] += M2TT(Psd, i, j, k, l) * A[k][l]
					}
				}
			}
		}
		// check basic
		if math.Abs(trA-tra) > 1e-17 {
			chk.Panic("tra failed. diff = %g", trA-tra)
		}
		if math.Abs(detA-deta) > 1e-14 {
			chk.Panic("detA failed. diff = %g", detA-deta)
		}
		if math.Abs(trDevA) > 1e-13 {
			chk.Panic("trDevA failed. error = %g", trDevA)
		}
		chk.Matrix(tst, "devA", 1e-13, devA, devA_)
		chk.Matrix(tst, "devA", 1e-13, devA, devA__)
		chk.Vector(tst, "devA", 1e-13, deva, deva__)
		chk.Matrix(tst, "A²", 1e-11, A2, A2_)
		// check tr(s2)
		io.Pfblue2("tr(s²) = %v\n", M_Tr(s2))
		if M_Tr(s2) < 1 {
			chk.Panic("Tr(s2) failed")
		}
		// check derivatives
		da2da := M_Alloc4(nd[m])
		a2tmp := M_Alloc2(nd[m]) // a2tmp == a²
		M_SqDeriv(da2da, a)
		var tmp float64
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a); j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, a[j] = a[j], x
					M_Sq(a2tmp, a)
					a[j] = tmp
					return a2tmp[i]
				}, a[j], 1e-6)
				chk.AnaNum(tst, io.Sf("da²/da[%d][%d]", i, j), dtol, da2da[i][j], dnum, dver)
			}
		}
		// characteristic invariants
		I1, I2, I3 := M_CharInvs(a)
		I2a := 0.5 * (tra*tra - M_Tr(a2))
		I1_, I2_, I3_, dI1da, dI2da, dI3da := M_CharInvsAndDerivs(a)
		if math.Abs(I1-tra) > 1e-17 {
			chk.Panic("I1 failed (a). error = %v", I1-tra)
		}
		if math.Abs(I2-I2a) > 1e-12 {
			chk.Panic("I2 failed (a). error = %v (I2=%v, I2_=%v)", I2-I2a, I2, I2a)
		}
		if math.Abs(I3-deta) > 1e-17 {
			chk.Panic("I3 failed (a). error = %v", I3-deta)
		}
		if math.Abs(I1-I1_) > 1e-17 {
			chk.Panic("I1 failed (b). error = %v", I1-I1_)
		}
		if math.Abs(I2-I2_) > 1e-17 {
			chk.Panic("I2 failed (b). error = %v", I2-I2_)
		}
		if math.Abs(I3-I3_) > 1e-17 {
			chk.Panic("I3 failed (b). error = %v", I3-I3_)
		}
		// dI1da
		for j := 0; j < len(a); j++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, a[j] = a[j], x
				i1, _, _ := M_CharInvs(a)
				a[j] = tmp
				return i1
			}, a[j], 1e-6)
			chk.AnaNum(tst, io.Sf("dI1/da[%d]", j), dtoli1[m], dI1da[j], dnum, dveri)
		}
		// dI2da
		for j := 0; j < len(a); j++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, a[j] = a[j], x
				_, i2, _ := M_CharInvs(a)
				a[j] = tmp
				return i2
			}, a[j], 1e-6)
			chk.AnaNum(tst, io.Sf("dI2/da[%d]", j), dtoli2[m], dI2da[j], dnum, dveri)
		}
		// dI3da
		for j := 0; j < len(a); j++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, a[j] = a[j], x
				_, _, i3 := M_CharInvs(a)
				a[j] = tmp
				return i3
			}, a[j], 1e-6)
			chk.AnaNum(tst, io.Sf("dI3/da[%d]", j), dtoli3[m], dI3da[j], dnum, dveri)
		}
		// dDet(a)/da
		DdetaDa := make([]float64, len(a))
		M_DetDeriv(DdetaDa, a)
		for j := 0; j < len(a); j++ {
			chk.AnaNum(tst, io.Sf("dDet(a)/da[%d]", j), dtoli3[m], dI3da[j], DdetaDa[j], dveri)
		}
		// lode angle
		if true {
			s := M_Alloc2(nd[m])
			dwda := M_Alloc2(nd[m])
			dwda_ := M_Alloc2(nd[m])
			d2wdada := M_Alloc4(nd[m])
			p, q, w := M_pqws(s, a)
			M_LodeDeriv1(dwda, a, s, p, q, w)
			M_LodeDeriv2(d2wdada, dwda_, a, s, p, q, w)
			chk.Vector(tst, "s", 1e-13, deva, s)
			chk.Vector(tst, "dwda", 1e-13, dwda, dwda_)
			// dwda
			for j := 0; j < len(a); j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, a[j] = a[j], x
					res = M_w(a)
					a[j] = tmp
					return res
				}, a[j], 1e-6)
				chk.AnaNum(tst, io.Sf("dw/da[%d]", j), dtolw, dwda[j], dnum, dverw)
			}
			// d2wdada
			s_tmp := M_Alloc2(nd[m])
			dwda_tmp := M_Alloc2(nd[m])
			for i := 0; i < len(a); i++ {
				for j := 0; j < len(a); j++ {
					dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
						tmp, a[j] = a[j], x
						p_tmp, q_tmp, w_tmp := M_pqws(s_tmp, a)
						M_LodeDeriv1(dwda_tmp, a, s_tmp, p_tmp, q_tmp, w_tmp)
						a[j] = tmp
						return dwda_tmp[i]
					}, a[j], 1e-6)
					chk.AnaNum(tst, io.Sf("d2w/dada[%d][%d]", i, j), dtolw, d2wdada[i][j], dnum, dverw)
				}
			}
		}
	}
}

func Test_ops02(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	chk.PrintTitle("ops02")

	nd := []int{2, 2, 3, 3, 3}
	AA := [][][]float64{
		{
			{1, 2, 0},
			{2, -2, 0},
			{0, 0, -2},
		},
		{
			{-100, 33, 0},
			{33, -200, 0},
			{0, 0, 150},
		},
		{
			{1, 2, 4},
			{2, -2, 3},
			{4, 3, -2},
		},
		{
			{-100, -10, 20},
			{-10, -200, 15},
			{20, 15, -300},
		},
		{
			{-100, 0, -10},
			{0, -200, 0},
			{-10, 0, 100},
		},
	}
	BB := [][][]float64{
		{
			{0.13, 1.2, 0},
			{1.2, -20, 0},
			{0, 0, -28},
		},
		{
			{-10, 3.3, 0},
			{3.3, -2, 0},
			{0, 0, 1.5},
		},
		{
			{0.1, 0.2, 0.8},
			{0.2, -1.3, 0.3},
			{0.8, 0.3, -0.2},
		},
		{
			{-10, -1, 2},
			{-1, -20, 1},
			{2, 1, -30},
		},
		{
			{-10, 3, -1},
			{3, -20, 1},
			{-1, 1, 10},
		},
	}

	nonsymTol := 1e-15

	for m := 0; m < len(nd); m++ {

		// tensors
		A := AA[m]
		B := BB[m]
		a := M_Alloc2(nd[m])
		b := M_Alloc2(nd[m])
		Ten2Man(a, A)
		Ten2Man(b, B)
		io.PfYel("\n\ntst # %d ###################################################################################\n", m)
		io.Pfblue2("a = %v\n", a)
		io.Pfblue2("b = %v\n", b)

		// dyadic
		c := M_Dy(a, b)
		c_ := M_Alloc4(nd[m])
		c__ := M_Alloc4(nd[m])
		M_DyAdd(c_, 1, a, b)
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a); j++ {
				for k := 0; k < len(a); k++ {
					c__[i][j] = a[i] * b[j]
				}
			}
		}
		chk.Matrix(tst, "a dy b", 1e-12, c, c_)
		chk.Matrix(tst, "a dy b", 1e-12, c, c__)

		// dot product
		d := M_Alloc2(nd[m])
		D := Alloc2()
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				for k := 0; k < 3; k++ {
					D[i][j] += A[i][k] * B[k][j]
				}
			}
		}
		err := M_Dot(d, a, b, nonsymTol)
		for i := 0; i < 3; i++ {
			chk.Scalar(tst, io.Sf("a_dot_b[%d][%d]", i, i), 1e-15, D[i][i], d[i])
		}
		/*
		   for k := 0; k < 2*nd[m]; k++ {
		       I, J := M2Ti[k], M2Tj[k]
		       cf   := 1.0
		       if k > 2 {
		           cf = 1.0 / SQ2
		       }
		       io.Pforan("%v %v\n", D[I][J], d[k] * cf)
		       chk.Scalar(tst, io.Sf("a_dot_b[%d][%d]",I,J), 1e-15, D[I][J], d[k] * cf)
		   }
		*/
		if err == nil {
			chk.Panic("dot product failed: error should be non-nil, since the result is expected to be non-symmetric")
		}

		// dot product (square tensor)
		a2 := M_Alloc2(nd[m])
		M_Sq(a2, a)
		aa := M_Alloc2(nd[m])
		tol_tmp := nonsymTol
		if m == 3 || m == 4 {
			nonsymTol = 1e-12
		}
		err = M_Dot(aa, a, a, nonsymTol)
		io.Pforan("a2 = %v\n", a2)
		io.Pforan("aa = %v\n", aa)
		if err != nil {
			chk.Panic("%v", err)
		}
		chk.Vector(tst, "a2", 1e-15, a2, aa)
		nonsymTol = tol_tmp
	}
}

func Test_ops03(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	chk.PrintTitle("ops03")

	nonsymTol := 1e-15

	dtol := 1e-9
	dver := true

	nd := test_nd
	for idxA := 0; idxA < len(test_nd)-3; idxA++ {
		//for idxA := 0; idxA < 1; idxA++ {

		// tensor and eigenvalues
		A := test_AA[idxA]
		a := M_Alloc2(nd[idxA])
		Ten2Man(a, A)
		io.PfYel("\n\ntst # %d ###################################################################################\n", idxA)
		io.Pfblue2("a = %v\n", a)

		// inverse
		Ai := Alloc2()
		ai := M_Alloc2(nd[idxA])
		detA, err := Inv(Ai, A)
		if err != nil {
			chk.Panic("%v", err)
		}
		deta_ := M_Det(a)
		deta, err := M_Inv(ai, a, MINDET)
		if err != nil {
			chk.Panic("%v", err)
		}
		Ai_ := Alloc2()
		Man2Ten(Ai_, ai)
		aia := M_Alloc2(nd[idxA])
		err = M_Dot(aia, ai, a, nonsymTol)
		if err != nil {
			chk.Panic("%v", err)
		}
		chk.Scalar(tst, "detA", 1e-14, detA, deta)
		chk.Scalar(tst, "deta", 1e-14, deta, deta_)
		chk.Matrix(tst, "Ai", 1e-14, Ai, Ai_)
		chk.Vector(tst, "ai*a", 1e-15, aia, Im[:2*nd[idxA]])
		io.Pforan("ai*a = %v\n", aia)

		// derivative of inverse
		dtol_tmp := dtol
		if idxA == 5 {
			dtol = 1e-8
		}
		var tmp float64
		ai_tmp := M_Alloc2(nd[idxA])
		daida := M_Alloc4(nd[idxA])
		M_InvDeriv(daida, ai)
		io.Pforan("ai = %v\n", ai)
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a); j++ {
				//dnum, _ := num.DerivForward(func(x float64, args ...interface{}) (res float64) {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, a[j] = a[j], x
					_, err := M_Inv(ai_tmp, a, MINDET)
					a[j] = tmp
					if err != nil {
						chk.Panic("daida failed:\n%v", err)
					}
					return ai_tmp[i]
				}, a[j], 1e-6)
				chk.AnaNum(tst, io.Sf("dai/da[%d][%d]", i, j), dtol, daida[i][j], dnum, dver)
			}
		}
		dtol = dtol_tmp
	}
}
