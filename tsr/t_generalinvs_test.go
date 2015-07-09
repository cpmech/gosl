// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/utl"
)

func Test_geninvs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("geninvs01")

	// coefficients for smp invariants
	smp_a := -1.0
	smp_b := 0.5
	smp_β := 1e-1 // derivative values become too high with
	smp_ϵ := 1e-1 // small β and ϵ @ zero

	// constants for checking derivatives
	dver := chk.Verbose
	dtol := 1e-6
	dtol2 := 1e-6

	// run tests
	nd := test_nd
	for idxA := 0; idxA < len(test_nd); idxA++ {
		//for idxA := 10; idxA < 11; idxA++ {

		// tensor and eigenvalues
		A := test_AA[idxA]
		a := M_Alloc2(nd[idxA])
		Ten2Man(a, A)
		L := make([]float64, 3)
		M_EigenValsNum(L, a)

		// SMP derivs and SMP director
		dndL := la.MatAlloc(3, 3)
		dNdL := make([]float64, 3)
		d2ndLdL := utl.Deep3alloc(3, 3, 3)
		N := make([]float64, 3)
		F := make([]float64, 3)
		G := make([]float64, 3)
		m := SmpDerivs1(dndL, dNdL, N, F, G, L, smp_a, smp_b, smp_β, smp_ϵ)
		SmpDerivs2(d2ndLdL, L, smp_a, smp_b, smp_β, smp_ϵ, m, N, F, G, dNdL, dndL)
		n := make([]float64, 3)
		SmpUnitDirector(n, m, N)

		// SMP invariants
		p, q, err := GenInvs(L, n, smp_a)
		if err != nil {
			chk.Panic("SmpInvs failed:\n%v", err)
		}

		// output
		io.PfYel("\n\ntst # %d ###################################################################################\n", idxA)
		io.Pfblue2("L = %v\n", L)
		io.Pforan("n = %v\n", n)
		io.Pforan("p = %v\n", p)
		io.Pforan("q = %v\n", q)

		// check invariants
		tvec := make([]float64, 3)
		GenTvec(tvec, L, n)
		proj := make([]float64, 3) // projection of tvec along n
		tdn := la.VecDot(tvec, n)  // tvec dot n
		for i := 0; i < 3; i++ {
			proj[i] = tdn * n[i]
		}
		norm_proj := la.VecNorm(proj)
		norm_tvec := la.VecNorm(tvec)
		q_ := GENINVSQEPS + math.Sqrt(norm_tvec*norm_tvec-norm_proj*norm_proj)
		io.Pforan("proj = %v\n", proj)
		io.Pforan("norm(proj) = %v == p\n", norm_proj)
		chk.Scalar(tst, "p", 1e-14, math.Abs(p), norm_proj)
		chk.Scalar(tst, "q", 1e-13, q, q_)

		// dt/dL
		var tmp float64
		N_tmp := make([]float64, 3)
		n_tmp := make([]float64, 3)
		tvec_tmp := make([]float64, 3)
		dtdL := la.MatAlloc(3, 3)
		GenTvecDeriv1(dtdL, L, n, dndL)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, L[j] = L[j], x
					m_tmp := SmpDirector(N_tmp, L, smp_a, smp_b, smp_β, smp_ϵ)
					SmpUnitDirector(n_tmp, m_tmp, N_tmp)
					GenTvec(tvec_tmp, L, n_tmp)
					L[j] = tmp
					return tvec_tmp[i]
				}, L[j], 1e-6)
				chk.AnaNum(tst, io.Sf("dt/dL[%d][%d]", i, j), dtol, dtdL[i][j], dnum, dver)
			}
		}

		// d²t/dLdL
		io.Pfpink("\nd²t/dLdL\n")
		dNdL_tmp := make([]float64, 3)
		dndL_tmp := la.MatAlloc(3, 3)
		dtdL_tmp := la.MatAlloc(3, 3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				for k := 0; k < 3; k++ {
					dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
						tmp, L[k] = L[k], x
						m_tmp := SmpDerivs1(dndL_tmp, dNdL_tmp, N_tmp, F, G, L, smp_a, smp_b, smp_β, smp_ϵ)
						SmpUnitDirector(n_tmp, m_tmp, N_tmp)
						GenTvecDeriv1(dtdL_tmp, L, n_tmp, dndL_tmp)
						L[k] = tmp
						return dtdL_tmp[i][j]
					}, L[k], 1e-6)
					dana := GenTvecDeriv2(i, j, k, L, dndL, d2ndLdL[i][j][k])
					chk.AnaNum(tst, io.Sf("d²t[%d]/dL[%d]dL[%d]", i, j, k), dtol2, dana, dnum, dver)
				}
			}
		}

		// change tolerance
		dtol_tmp := dtol
		switch idxA {
		case 5, 11:
			dtol = 1e-5
		case 12:
			dtol = 0.0013
		}

		// first order derivatives
		dpdL := make([]float64, 3)
		dqdL := make([]float64, 3)
		p_, q_, err := GenInvsDeriv1(dpdL, dqdL, L, n, dndL, smp_a)
		if err != nil {
			chk.Panic("%v", err)
		}
		chk.Scalar(tst, "p", 1e-17, p, p_)
		chk.Scalar(tst, "q", 1e-17, q, q_)
		var ptmp, qtmp float64
		io.Pfpink("\ndp/dL\n")
		for j := 0; j < 3; j++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, L[j] = L[j], x
				m_tmp := SmpDirector(N_tmp, L, smp_a, smp_b, smp_β, smp_ϵ)
				SmpUnitDirector(n_tmp, m_tmp, N_tmp)
				ptmp, _, err = GenInvs(L, n_tmp, smp_a)
				if err != nil {
					chk.Panic("DerivCentral: SmpInvs failed:\n%v", err)
				}
				L[j] = tmp
				return ptmp
			}, L[j], 1e-6)
			chk.AnaNum(tst, io.Sf("dp/dL[%d]", j), dtol, dpdL[j], dnum, dver)
		}
		io.Pfpink("\ndq/dL\n")
		for j := 0; j < 3; j++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, L[j] = L[j], x
				m_tmp := SmpDirector(N_tmp, L, smp_a, smp_b, smp_β, smp_ϵ)
				SmpUnitDirector(n_tmp, m_tmp, N_tmp)
				_, qtmp, err = GenInvs(L, n_tmp, smp_a)
				if err != nil {
					chk.Panic("DerivCentral: SmpInvs failed:\n%v", err)
				}
				L[j] = tmp
				return qtmp
			}, L[j], 1e-6)
			chk.AnaNum(tst, io.Sf("dq/dL[%d]", j), dtol, dqdL[j], dnum, dver)
		}

		// recover tolerance
		dtol = dtol_tmp

		// change tolerance
		io.Pforan("dtol2 = %v\n", dtol2)
		dtol2_tmp := dtol2
		switch idxA {
		case 5:
			dtol2 = 1e-5
		case 10:
			dtol2 = 0.72
		case 11:
			dtol2 = 1e-5
		case 12:
			dtol2 = 544
		}

		// second order derivatives
		dpdL_tmp := make([]float64, 3)
		dqdL_tmp := make([]float64, 3)
		d2pdLdL := la.MatAlloc(3, 3)
		d2qdLdL := la.MatAlloc(3, 3)
		GenInvsDeriv2(d2pdLdL, d2qdLdL, L, n, dpdL, dqdL, p, q, dndL, d2ndLdL, smp_a)
		io.Pfpink("\nd²p/dLdL\n")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, L[j] = L[j], x
					m_tmp := SmpDerivs1(dndL_tmp, dNdL_tmp, N_tmp, F, G, L, smp_a, smp_b, smp_β, smp_ϵ)
					SmpUnitDirector(n_tmp, m_tmp, N_tmp)
					GenInvsDeriv1(dpdL_tmp, dqdL_tmp, L, n_tmp, dndL_tmp, smp_a)
					L[j] = tmp
					return dpdL_tmp[i]
				}, L[j], 1e-6)
				chk.AnaNum(tst, io.Sf("d²p/dL[%d][%d]", i, j), dtol2, d2pdLdL[i][j], dnum, dver)
			}
		}
		io.Pfpink("\nd²q/dLdL\n")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, L[j] = L[j], x
					m_tmp := SmpDerivs1(dndL_tmp, dNdL_tmp, N_tmp, F, G, L, smp_a, smp_b, smp_β, smp_ϵ)
					SmpUnitDirector(n_tmp, m_tmp, N_tmp)
					GenInvsDeriv1(dpdL_tmp, dqdL_tmp, L, n_tmp, dndL_tmp, smp_a)
					L[j] = tmp
					return dqdL_tmp[i]
				}, L[j], 1e-6)
				chk.AnaNum(tst, io.Sf("d²q/dL[%d][%d]", i, j), dtol2, d2qdLdL[i][j], dnum, dver)
			}
		}

		// recover tolerance
		dtol2 = dtol2_tmp
	}
}
