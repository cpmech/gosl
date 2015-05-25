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

	tol := 1e-15
	dtol := 1e-6
	dver := chk.Verbose
	dtol2 := 1e-6

	smp_a := -1.0
	smp_b := 0.5
	smp_β := 2.0
	smp_ϵ := 1e-3

	nd := test_nd
	for idxA := 0; idxA < len(test_nd)-3; idxA++ {
		//for idxA := 0; idxA < 1; idxA++ {

		// tensor and eigenvalues
		A := test_AA[idxA]
		a := M_Alloc2(nd[idxA])
		Ten2Man(a, A)
		λ := make([]float64, 3)
		M_EigenValsNum(λ, a)

		// shifted eigenvalues
		σc := 0.0
		for j := 0; j < 3; j++ {
			if λ[j] >= σc {
				σc = λ[j] * 1.01
			}
		}
		σ := make([]float64, 3)
		err := ShiftedEigenvs(σ, λ, σc, tol)
		if err != nil {
			chk.Panic("%v\n", err)
		}

		// SMP director
		N := make([]float64, 3)
		n := make([]float64, 3)
		m := SmpDirector(N, σ, smp_a, smp_b, smp_β, smp_ϵ)
		SmpUnitDirector(n, m, N)

		// output
		io.PfYel("\n\ntst # %d ###################################################################################\n", idxA)
		io.Pfblue2("σ  = %v\n", σ)
		io.Pforan("λ = %v\n", λ)
		io.Pforan("N = %v\n", N)
		io.Pforan("m = %v\n", m)
		io.Pfpink("n = %v\n", n)
		chk.Vector(tst, "λ", 1e-12, λ, test_λ[idxA])
		chk.Scalar(tst, "norm(n)==1", 1e-15, la.VecNorm(n), 1)
		chk.Scalar(tst, "m=norm(N)", 1e-14, m, la.VecNorm(N))

		// dN/dσ
		var tmp float64
		N_tmp := make([]float64, 3)
		dNdσ := make([]float64, 3)
		SmpDirectorDeriv1(dNdσ, σ, smp_a, smp_b, smp_β, smp_ϵ)
		io.Pfpink("\ndNdσ = %v\n", dNdσ)
		dtol_tmp := dtol
		if idxA == 0 || idxA == 2 {
			dtol = 1e-5
		}
		for i := 0; i < 3; i++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, σ[i] = σ[i], x
				SmpDirector(N_tmp, σ, smp_a, smp_b, smp_β, smp_ϵ)
				σ[i] = tmp
				return N_tmp[i]
			}, σ[i], 1e-6)
			chk.AnaNum(tst, io.Sf("dN/dσ[%d][%d]", i, i), dtol, dNdσ[i], dnum, dver)
		}
		dtol = dtol_tmp

		// dm/dσ
		n_tmp := make([]float64, 3)
		dmdσ := make([]float64, 3)
		SmpNormDirectorDeriv1(dmdσ, m, N, dNdσ)
		io.Pfpink("\ndmdσ = %v\n", dmdσ)
		dtol_tmp = dtol
		if idxA == 0 || idxA == 2 {
			dtol = 1e-5
		}
		for j := 0; j < 3; j++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, σ[j] = σ[j], x
				m_tmp := SmpDirector(N_tmp, σ, smp_a, smp_b, smp_β, smp_ϵ)
				σ[j] = tmp
				return m_tmp
			}, σ[j], 1e-6)
			chk.AnaNum(tst, io.Sf("dm/dσ[%d]", j), dtol, dmdσ[j], dnum, dver)
		}
		dtol = dtol_tmp

		// dn/dσ
		dndσ := la.MatAlloc(3, 3)
		SmpUnitDirectorDeriv1(dndσ, m, N, dNdσ, dmdσ)
		io.Pfpink("\ndndσ = %v\n", dndσ)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, σ[j] = σ[j], x
					m_tmp := SmpDirector(N_tmp, σ, smp_a, smp_b, smp_β, smp_ϵ)
					SmpUnitDirector(n_tmp, m_tmp, N_tmp)
					σ[j] = tmp
					return n_tmp[i]
				}, σ[j], 1e-6)
				chk.AnaNum(tst, io.Sf("dn/dσ[%d][%d]", i, j), dtol, dndσ[i][j], dnum, dver)
			}
		}

		// d²m/dσdσ
		dNdσ_tmp := make([]float64, 3)
		dmdσ_tmp := make([]float64, 3)
		d2Ndσ2 := make([]float64, 3)
		d2mdσdσ := la.MatAlloc(3, 3)
		SmpDirectorDeriv2(d2Ndσ2, σ, smp_a, smp_b, smp_β, smp_ϵ)
		SmpNormDirectorDeriv2(d2mdσdσ, σ, smp_a, smp_b, smp_β, smp_ϵ, m, N, dNdσ, d2Ndσ2, dmdσ)
		io.Pfpink("\nd2mdσdσ = %v\n", d2mdσdσ)
		tol_tmp := dtol2
		if idxA == 0 {
			dtol2 = 1e-5
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, σ[j] = σ[j], x
					m_tmp := SmpDirector(N_tmp, σ, smp_a, smp_b, smp_β, smp_ϵ)
					SmpDirectorDeriv1(dNdσ_tmp, σ, smp_a, smp_b, smp_β, smp_ϵ)
					SmpNormDirectorDeriv1(dmdσ_tmp, m_tmp, N_tmp, dNdσ_tmp)
					σ[j] = tmp
					return dmdσ_tmp[i]
				}, σ[j], 1e-6)
				chk.AnaNum(tst, io.Sf("d2m/dσ[%d]dσ[%d]", i, j), dtol2, d2mdσdσ[i][j], dnum, dver)
			}
		}
		dtol2 = tol_tmp

		// d²N/dσdσ
		io.Pfpink("\nd²N/dσdσ\n")
		tol_tmp = dtol2
		if idxA == 0 || idxA == 7 {
			dtol2 = 1e-5
		}
		for i := 0; i < 3; i++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, σ[i] = σ[i], x
				SmpDirectorDeriv1(dNdσ_tmp, σ, smp_a, smp_b, smp_β, smp_ϵ)
				σ[i] = tmp
				return dNdσ_tmp[i]
			}, σ[i], 1e-6)
			chk.AnaNum(tst, io.Sf("d²N[%d]/dσ[%d]dσ[%d]", i, i, i), dtol2, d2Ndσ2[i], dnum, dver)
		}
		dtol2 = tol_tmp

		// d²n/dσdσ
		io.Pfpink("\nd²n/dσdσ\n")
		dndσ_tmp := la.MatAlloc(3, 3)
		d2ndσdσ := utl.Deep3alloc(3, 3, 3)
		SmpUnitDirectorDeriv2(d2ndσdσ, m, N, dNdσ, d2Ndσ2, dmdσ, n, d2mdσdσ, dndσ)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				for k := 0; k < 3; k++ {
					dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
						tmp, σ[k] = σ[k], x
						m_tmp := SmpDirector(N_tmp, σ, smp_a, smp_b, smp_β, smp_ϵ)
						SmpDirectorDeriv1(dNdσ_tmp, σ, smp_a, smp_b, smp_β, smp_ϵ)
						SmpNormDirectorDeriv1(dmdσ_tmp, m_tmp, N_tmp, dNdσ_tmp)
						SmpUnitDirectorDeriv1(dndσ_tmp, m_tmp, N_tmp, dNdσ_tmp, dmdσ_tmp)
						σ[k] = tmp
						return dndσ_tmp[i][j]
					}, σ[k], 1e-6)
					chk.AnaNum(tst, io.Sf("d²n[%d]/dσ[%d]dσ[%d]", i, j, k), dtol2, d2ndσdσ[i][j][k], dnum, dver)
				}
			}
		}

		// SMP derivs
		io.Pfpink("\nSMP derivs\n")
		dndσ_ := la.MatAlloc(3, 3)
		dNdσ_ := make([]float64, 3)
		d2ndσdσ_ := utl.Deep3alloc(3, 3, 3)
		N_ := make([]float64, 3)
		F_ := make([]float64, 3)
		G_ := make([]float64, 3)
		m_ := SmpDerivs1(dndσ_, dNdσ_, N_, F_, G_, σ, smp_a, smp_b, smp_β, smp_ϵ)
		SmpDerivs2(d2ndσdσ_, σ, smp_a, smp_b, smp_β, smp_ϵ, m_, N_, F_, G_, dNdσ_, dndσ_)
		chk.Scalar(tst, "m_", 1e-14, m_, m)
		chk.Vector(tst, "N_", 1e-15, N_, N)
		chk.Vector(tst, "dNdσ_", 1e-15, dNdσ_, dNdσ)
		chk.Matrix(tst, "dndσ_", 1e-13, dndσ_, dndσ)
		chk.Deep3(tst, "d2ndσdσ_", 1e-13, d2ndσdσ_, d2ndσdσ)
	}
}

func Test_geninvs02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("geninvs02")

	tol := 1e-10
	dtol := 1e-6
	dtol2 := 1e-6
	dver := chk.Verbose

	smp_a := -1.0
	smp_b := 0.5
	smp_β := 2.0
	smp_ϵ := 1e-3

	nd := test_nd
	for idxA := 0; idxA < len(test_nd)-3; idxA++ {
		//for idxA := 0; idxA < 1; idxA++ {

		// tensor and eigenvalues
		A := test_AA[idxA]
		a := M_Alloc2(nd[idxA])
		Ten2Man(a, A)
		λ := make([]float64, 3)
		M_EigenValsNum(λ, a)
		σc := 0.0
		cf := 1.1 // 1.01
		for j := 0; j < 3; j++ {
			if λ[j] >= σc {
				σc = λ[j] * cf
			}
		}

		// shifted eigenvalues
		σ := make([]float64, 3)
		err := ShiftedEigenvs(σ, λ, σc, tol)
		if err != nil {
			chk.Panic("%v\n", err)
		}

		// SMP derivs and SMP director
		dndσ := la.MatAlloc(3, 3)
		dNdσ := make([]float64, 3)
		d2ndσdσ := utl.Deep3alloc(3, 3, 3)
		N := make([]float64, 3)
		F := make([]float64, 3)
		G := make([]float64, 3)
		m := SmpDerivs1(dndσ, dNdσ, N, F, G, σ, smp_a, smp_b, smp_β, smp_ϵ)
		SmpDerivs2(d2ndσdσ, σ, smp_a, smp_b, smp_β, smp_ϵ, m, N, F, G, dNdσ, dndσ)
		n := make([]float64, 3)
		SmpUnitDirector(n, m, N)

		// SMP invariants
		p, q, err := GenInvs(σ, n, smp_a)
		if err != nil {
			chk.Panic("SmpInvs failed:\n%v", err)
		}

		// output
		io.PfYel("\n\ntst # %d ###################################################################################\n", idxA)
		io.Pfblue2("σ = %v\n", σ)
		io.Pforan("n = %v\n", n)
		io.Pforan("p = %v\n", p)
		io.Pforan("q = %v\n", q)

		// check invariants
		tvec := make([]float64, 3)
		GenTvec(tvec, σ, n)
		proj := make([]float64, 3) // projection of tvec along n
		tdn := la.VecDot(tvec, n)  // tvec dot n
		for i := 0; i < 3; i++ {
			proj[i] = tdn * n[i]
		}
		norm_proj := la.VecNorm(proj)
		norm_tvec := la.VecNorm(tvec)
		q_ := math.Sqrt(norm_tvec*norm_tvec - norm_proj*norm_proj)
		io.Pforan("proj = %v\n", proj)
		io.Pforan("norm(proj) = %v == p\n", norm_proj)
		chk.Scalar(tst, "p", 1e-14, p, smp_a*norm_proj)
		chk.Scalar(tst, "q", 1e-13, q, q_)

		// dt/dσ
		var tmp float64
		N_tmp := make([]float64, 3)
		n_tmp := make([]float64, 3)
		tvec_tmp := make([]float64, 3)
		dtdσ := la.MatAlloc(3, 3)
		GenTvecDeriv1(dtdσ, σ, n, dndσ)
		dtol_tmp := dtol
		if idxA == 4 {
			dtol = 1e-5
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, σ[j] = σ[j], x
					m_tmp := SmpDirector(N_tmp, σ, smp_a, smp_b, smp_β, smp_ϵ)
					SmpUnitDirector(n_tmp, m_tmp, N_tmp)
					GenTvec(tvec_tmp, σ, n_tmp)
					σ[j] = tmp
					return tvec_tmp[i]
				}, σ[j], 1e-6)
				chk.AnaNum(tst, io.Sf("dt/dσ[%d][%d]", i, j), dtol, dtdσ[i][j], dnum, dver)
			}
		}
		dtol = dtol_tmp

		// d²t/dσdσ
		io.Pfpink("\nd²t/dσdσ\n")
		dNdσ_tmp := make([]float64, 3)
		dndσ_tmp := la.MatAlloc(3, 3)
		dtdσ_tmp := la.MatAlloc(3, 3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				for k := 0; k < 3; k++ {
					dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
						tmp, σ[k] = σ[k], x
						m_tmp := SmpDerivs1(dndσ_tmp, dNdσ_tmp, N_tmp, F, G, σ, smp_a, smp_b, smp_β, smp_ϵ)
						SmpUnitDirector(n_tmp, m_tmp, N_tmp)
						GenTvecDeriv1(dtdσ_tmp, σ, n_tmp, dndσ_tmp)
						σ[k] = tmp
						return dtdσ_tmp[i][j]
					}, σ[k], 1e-6)
					dana := GenTvecDeriv2(i, j, k, σ, dndσ, d2ndσdσ[i][j][k])
					chk.AnaNum(tst, io.Sf("d²t[%d]/dσ[%d]dσ[%d]", i, j, k), dtol2, dana, dnum, dver)
				}
			}
		}

		// first order derivatives
		dpdσ := make([]float64, 3)
		dqdσ := make([]float64, 3)
		p_, q_, err := GenInvsDeriv1(dpdσ, dqdσ, σ, n, dndσ, smp_a)
		if err != nil {
			chk.Panic("%v", err)
		}
		chk.Scalar(tst, "p", 1e-17, p, p_)
		chk.Scalar(tst, "q", 1e-17, q, q_)
		var ptmp, qtmp float64
		io.Pfpink("\ndp/dσ\n")
		dtol_tmp = dtol
		if idxA == 4 {
			dtol = 1e-5
		}
		for j := 0; j < 3; j++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, σ[j] = σ[j], x
				m_tmp := SmpDirector(N_tmp, σ, smp_a, smp_b, smp_β, smp_ϵ)
				SmpUnitDirector(n_tmp, m_tmp, N_tmp)
				ptmp, _, err = GenInvs(σ, n_tmp, smp_a)
				if err != nil {
					chk.Panic("DerivCentral: SmpInvs failed:\n%v", err)
				}
				if σ[0] < 1e-14 || σ[1] < 1e-14 || σ[2] < 1e-14 {
					chk.Panic("dp/dσ failed: σ=%v must be all greater than %v", σ, 1e-14)
				}
				σ[j] = tmp
				return ptmp
			}, σ[j], 1e-6)
			chk.AnaNum(tst, io.Sf("dp/dσ[%d]", j), dtol, dpdσ[j], dnum, dver)
		}
		io.Pfpink("\ndq/dσ\n")
		for j := 0; j < 3; j++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, σ[j] = σ[j], x
				m_tmp := SmpDirector(N_tmp, σ, smp_a, smp_b, smp_β, smp_ϵ)
				SmpUnitDirector(n_tmp, m_tmp, N_tmp)
				_, qtmp, err = GenInvs(σ, n_tmp, smp_a)
				if err != nil {
					chk.Panic("DerivCentral: SmpInvs failed:\n%v", err)
				}
				if σ[0] < 1e-14 || σ[1] < 1e-14 || σ[2] < 1e-14 {
					chk.Panic("dq/dσ failed: σ=%v must be all greater than %v", σ, 1e-14)
				}
				σ[j] = tmp
				return qtmp
			}, σ[j], 1e-6)
			chk.AnaNum(tst, io.Sf("dq/dσ[%d]", j), dtol, dqdσ[j], dnum, dver)
		}
		dtol = dtol_tmp

		// second order derivatives
		dpdσ_tmp := make([]float64, 3)
		dqdσ_tmp := make([]float64, 3)
		d2pdσdσ := la.MatAlloc(3, 3)
		d2qdσdσ := la.MatAlloc(3, 3)
		GenInvsDeriv2(d2pdσdσ, d2qdσdσ, σ, n, dpdσ, dqdσ, p, q, dndσ, d2ndσdσ, smp_a)
		io.Pfpink("\nd²p/dσdσ\n")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, σ[j] = σ[j], x
					m_tmp := SmpDerivs1(dndσ_tmp, dNdσ_tmp, N_tmp, F, G, σ, smp_a, smp_b, smp_β, smp_ϵ)
					SmpUnitDirector(n_tmp, m_tmp, N_tmp)
					GenInvsDeriv1(dpdσ_tmp, dqdσ_tmp, σ, n_tmp, dndσ_tmp, smp_a)
					if σ[0] < 1e-14 || σ[1] < 1e-14 || σ[2] < 1e-14 {
						chk.Panic("d²p/dσdσdp/dσ failed: σ=%v must be all greater than %v", σ, 1e-14)
					}
					σ[j] = tmp
					return dpdσ_tmp[i]
				}, σ[j], 1e-6)
				chk.AnaNum(tst, io.Sf("d²p/dσ[%d][%d]", i, j), dtol2, d2pdσdσ[i][j], dnum, dver)
			}
		}
		io.Pfpink("\nd²q/dσdσ\n")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, σ[j] = σ[j], x
					m_tmp := SmpDerivs1(dndσ_tmp, dNdσ_tmp, N_tmp, F, G, σ, smp_a, smp_b, smp_β, smp_ϵ)
					SmpUnitDirector(n_tmp, m_tmp, N_tmp)
					GenInvsDeriv1(dpdσ_tmp, dqdσ_tmp, σ, n_tmp, dndσ_tmp, smp_a)
					if σ[0] < 1e-14 || σ[1] < 1e-14 || σ[2] < 1e-14 {
						chk.Panic("d²q/dσdσdp/dσ failed: σ=%v must be all greater than %v", σ, 1e-14)
					}
					σ[j] = tmp
					return dqdσ_tmp[i]
				}, σ[j], 1e-6)
				chk.AnaNum(tst, io.Sf("d²q/dσ[%d][%d]", i, j), dtol2, d2qdσdσ[i][j], dnum, dver)
			}
		}
	}
}
