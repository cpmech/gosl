// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

func Test_smpinvs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("smpinvs01")

	a, b, β, ϵ := -1.0, 0.5, 1e-3, 1e-3

	L := []float64{-8.0, -8.0, -8.0}
	N := make([]float64, 3)
	n := make([]float64, 3)
	m := SmpDirector(N, L, a, b, β, ϵ)
	SmpUnitDirector(n, m, N)
	io.Pforan("L = %v\n", L)
	io.Pforan("N = %v\n", N)
	io.Pforan("m = %v\n", m)
	io.Pforan("n = %v\n", n)
	chk.Vector(tst, "n", 1e-15, n, []float64{a / SQ3, a / SQ3, a / SQ3})

	p, q, err := GenInvs(L, n, a)
	if err != nil {
		chk.Panic("GenInvs failed:\n%v", err.Error())
	}
	io.Pforan("p = %v\n", p)
	io.Pforan("q = %v\n", q)
	if q < 0.0 || q > GENINVSQEPS {
		chk.Panic("q=%g is incorrect", q)
	}
	if math.Abs(p-a*L[0]) > 1e-14 {
		chk.Panic("p=%g is incorrect. err = %g", p, math.Abs(p-a*L[0]))
	}
}

func Test_smpinvs02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("smpinvs02")

	// coefficients for smp invariants
	smp_a := -1.0
	smp_b := 0.5
	smp_β := 1e-1 // derivative values become too high with
	smp_ϵ := 1e-1 // small β and ϵ @ zero

	// constants for checking derivatives
	dver := chk.Verbose
	dtol := 1e-9
	dtol2 := 1e-8

	// run tests
	nd := test_nd
	for idxA := 0; idxA < len(test_nd); idxA++ {
		//for idxA := 0; idxA < 1; idxA++ {
		//for idxA := 10; idxA < 11; idxA++ {

		// tensor and eigenvalues
		A := test_AA[idxA]
		a := M_Alloc2(nd[idxA])
		Ten2Man(a, A)
		L := make([]float64, 3)
		M_EigenValsNum(L, a)

		// SMP director
		N := make([]float64, 3)
		n := make([]float64, 3)
		m := SmpDirector(N, L, smp_a, smp_b, smp_β, smp_ϵ)
		SmpUnitDirector(n, m, N)

		// output
		io.PfYel("\n\ntst # %d ###################################################################################\n", idxA)
		io.Pforan("L = %v\n", L)
		io.Pforan("N = %v\n", N)
		io.Pforan("m = %v\n", m)
		io.Pfpink("n = %v\n", n)
		chk.Vector(tst, "L", 1e-12, L, test_λ[idxA])
		chk.Scalar(tst, "norm(n)==1", 1e-15, la.VecNorm(n), 1)
		chk.Scalar(tst, "m=norm(N)", 1e-14, m, la.VecNorm(N))

		// dN/dL
		var tmp float64
		N_tmp := make([]float64, 3)
		dNdL := make([]float64, 3)
		SmpDirectorDeriv1(dNdL, L, smp_a, smp_b, smp_β, smp_ϵ)
		io.Pfpink("\ndNdL = %v\n", dNdL)
		for i := 0; i < 3; i++ {
			chk.DerivScaSca(tst, io.Sf("dN%d/dL%d", i, i), dtol, dNdL[i], L[i], 1e-1, dver, func(x float64) (float64, error) {
				tmp, L[i] = L[i], x
				SmpDirector(N_tmp, L, smp_a, smp_b, smp_β, smp_ϵ)
				L[i] = tmp
				return N_tmp[i], nil
			})
		}

		// dm/dL
		dmdL := make([]float64, 3)
		SmpNormDirectorDeriv1(dmdL, m, N, dNdL)
		io.Pfpink("\ndmdL = %v\n", dmdL)
		chk.DerivScaVec(tst, "dm/dL", dtol, dmdL, L, 1e-1, dver, func(x []float64) (float64, error) {
			return SmpDirector(N_tmp, x, smp_a, smp_b, smp_β, smp_ϵ), nil
		})

		// dn/dL
		dndL := la.MatAlloc(3, 3)
		SmpUnitDirectorDeriv1(dndL, m, N, dNdL, dmdL)
		io.Pfpink("\ndndL = %v\n", dndL)
		chk.DerivVecVec(tst, "dn/dL", dtol, dndL, L, 1e-1, dver, func(f, x []float64) error {
			m_tmp := SmpDirector(N_tmp, x, smp_a, smp_b, smp_β, smp_ϵ)
			SmpUnitDirector(f, m_tmp, N_tmp) // f := n
			return nil
		})

		// change tolerance
		dtol2_tmp := dtol2
		if idxA == 10 || idxA == 11 {
			dtol2 = 1e-6
		}

		// d²m/dLdL
		dNdL_tmp := make([]float64, 3)
		dmdL_tmp := make([]float64, 3)
		d2NdL2 := make([]float64, 3)
		d2mdLdL := la.MatAlloc(3, 3)
		SmpDirectorDeriv2(d2NdL2, L, smp_a, smp_b, smp_β, smp_ϵ)
		SmpNormDirectorDeriv2(d2mdLdL, L, smp_a, smp_b, smp_β, smp_ϵ, m, N, dNdL, d2NdL2, dmdL)
		io.Pfpink("\nd2mdLdL = %v\n", d2mdLdL)
		chk.DerivVecVec(tst, "d2m/dLdL", dtol2, d2mdLdL, L, 1e-6, dver, func(f, x []float64) error {
			m_tmp := SmpDirector(N_tmp, x, smp_a, smp_b, smp_β, smp_ϵ)
			SmpDirectorDeriv1(dNdL_tmp, x, smp_a, smp_b, smp_β, smp_ϵ)
			SmpNormDirectorDeriv1(f, m_tmp, N_tmp, dNdL_tmp) // f := dmdL
			return nil
		})

		// d²N/dLdL
		io.Pfpink("\nd²N/dLdL\n")
		for i := 0; i < 3; i++ {
			chk.DerivScaSca(tst, io.Sf("d²N%d/dL%ddL%d", i, i, i), dtol2, d2NdL2[i], L[i], 1e-6, dver, func(x float64) (float64, error) {
				tmp, L[i] = L[i], x
				SmpDirectorDeriv1(dNdL_tmp, L, smp_a, smp_b, smp_β, smp_ϵ)
				L[i] = tmp
				return dNdL_tmp[i], nil
			})
		}

		// d²n/dLdL
		io.Pfpink("\nd²n/dLdL\n")
		dndL_tmp := la.MatAlloc(3, 3)
		d2ndLdL := utl.Deep3alloc(3, 3, 3)
		SmpUnitDirectorDeriv2(d2ndLdL, m, N, dNdL, d2NdL2, dmdL, n, d2mdLdL, dndL)
		for i := 0; i < 3; i++ {
			chk.DerivVecVec(tst, io.Sf("d²n%d/dLdL", i), dtol2, d2ndLdL[i], L, 1e-6, dver, func(f, x []float64) error {
				m_tmp := SmpDirector(N_tmp, x, smp_a, smp_b, smp_β, smp_ϵ)
				SmpDirectorDeriv1(dNdL_tmp, x, smp_a, smp_b, smp_β, smp_ϵ)
				SmpNormDirectorDeriv1(dmdL_tmp, m_tmp, N_tmp, dNdL_tmp)
				SmpUnitDirectorDeriv1(dndL_tmp, m_tmp, N_tmp, dNdL_tmp, dmdL_tmp)
				copy(f, dndL_tmp[i])
				return nil
			})
		}

		// recover tolerance
		dtol2 = dtol2_tmp

		// SMP derivs
		//if false {
		if true {
			io.Pfpink("\nSMP derivs\n")
			dndL_ := la.MatAlloc(3, 3)
			dNdL_ := make([]float64, 3)
			d2ndLdL_ := utl.Deep3alloc(3, 3, 3)
			N_ := make([]float64, 3)
			F_ := make([]float64, 3)
			G_ := make([]float64, 3)
			m_ := SmpDerivs1(dndL_, dNdL_, N_, F_, G_, L, smp_a, smp_b, smp_β, smp_ϵ)
			SmpDerivs2(d2ndLdL_, L, smp_a, smp_b, smp_β, smp_ϵ, m_, N_, F_, G_, dNdL_, dndL_)
			chk.Scalar(tst, "m_", 1e-14, m_, m)
			chk.Vector(tst, "N_", 1e-15, N_, N)
			chk.Vector(tst, "dNdL_", 1e-15, dNdL_, dNdL)
			chk.Matrix(tst, "dndL_", 1e-13, dndL_, dndL)
			chk.Deep3(tst, "d2ndLdL_", 1e-11, d2ndLdL_, d2ndLdL)
		}
	}
}
