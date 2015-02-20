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
	b := 0.5
	dtol := 1e-7
	dver := chk.Verbose
	dtol2 := 1e-6

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
		m := SmpUnitDirector(n, σ, b)
		SmpDirector(N, σ, b)

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
		dNdσ := la.MatAlloc(3, 3)
		SmpDirectorDeriv1(dNdσ, σ, b)
		io.Pfpink("\ndNdσ = %v\n", dNdσ)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, σ[j] = σ[j], x
					SmpDirector(N_tmp, σ, b)
					σ[j] = tmp
					return N_tmp[i]
				}, σ[j], 1e-6)
				chk.AnaNum(tst, io.Sf("dN/dσ[%d][%d]", i, j), dtol, dNdσ[i][j], dnum, dver)
			}
		}

		// dm/dσ
		n_tmp := make([]float64, 3)
		dmdσ := make([]float64, 3)
		mm := SmpNormDirectorDeriv1(dmdσ, σ, b)
		io.Pfpink("\ndmdσ = %v\n", dmdσ)
		chk.Scalar(tst, "m", 1e-17, m, mm)
		dtol_tmp := dtol
		if idxA == 5 {
			dtol = 1e-6
		}
		for j := 0; j < 3; j++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, σ[j] = σ[j], x
				m := SmpUnitDirector(n_tmp, σ, b)
				σ[j] = tmp
				return m
			}, σ[j], 1e-6)
			chk.AnaNum(tst, io.Sf("dm/dσ[%d]", j), dtol, dmdσ[j], dnum, dver)
		}
		dtol = dtol_tmp

		// dn/dσ
		dndσ := la.MatAlloc(3, 3)
		SmpUnitDirectorDeriv1(dndσ, σ, n, b, m, dmdσ)
		io.Pfpink("\ndndσ = %v\n", dndσ)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, σ[j] = σ[j], x
					SmpUnitDirector(n_tmp, σ, b)
					σ[j] = tmp
					return n_tmp[i]
				}, σ[j], 1e-6)
				chk.AnaNum(tst, io.Sf("dn/dσ[%d][%d]", i, j), dtol, dndσ[i][j], dnum, dver)
			}
		}

		// d²m/dσdσ
		dmdσ_tmp := make([]float64, 3)
		d2mdσdσ := la.MatAlloc(3, 3)
		SmpNormDirectorDeriv2(d2mdσdσ, σ, b, m, dmdσ)
		io.Pfpink("\nd2mdσdσ = %v\n", d2mdσdσ)
		tol_tmp := dtol2
		if idxA == 5 {
			dtol2 = 1e-3
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, σ[j] = σ[j], x
					SmpNormDirectorDeriv1(dmdσ_tmp, σ, b)
					σ[j] = tmp
					return dmdσ_tmp[i]
				}, σ[j], 1e-6)
				chk.AnaNum(tst, io.Sf("d2m/dσ[%d]dσ[%d]", i, j), dtol2, d2mdσdσ[i][j], dnum, dver)
			}
		}
		dtol2 = tol_tmp

		// d²N/dσdσ
		io.Pfpink("\nd²N/dσdσ\n")
		dNdσ_tmp := la.MatAlloc(3, 3)
		tol_tmp = dtol2
		if idxA == 5 {
			dtol2 = 1e-4
		}
		if idxA == 7 {
			dtol2 = 1e-5
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				for k := 0; k < 3; k++ {
					dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
						tmp, σ[k] = σ[k], x
						SmpDirectorDeriv1(dNdσ_tmp, σ, b)
						σ[k] = tmp
						return dNdσ_tmp[i][j]
					}, σ[k], 1e-6)
					dana := SmpDirectorDeriv2(i, j, k, σ, b)
					chk.AnaNum(tst, io.Sf("d²N[%d]/dσ[%d]dσ[%d]", i, j, k), dtol2, dana, dnum, dver)
				}
			}
		}
		dtol2 = tol_tmp

		// d²n/dσdσ
		io.Pfpink("\nd²n/dσdσ\n")
		dndσ_tmp := la.MatAlloc(3, 3)
		d2ndσdσ := utl.Deep3alloc(3, 3, 3)
		SmpUnitDirectorDeriv2(d2ndσdσ, σ, n, dmdσ, b, m, d2mdσdσ, dndσ)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				for k := 0; k < 3; k++ {
					dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
						tmp, σ[k] = σ[k], x
						SmpUnitDirector(n_tmp, σ, b)
						m_tmp := SmpNormDirectorDeriv1(dmdσ_tmp, σ, b)
						SmpUnitDirectorDeriv1(dndσ_tmp, σ, n_tmp, b, m_tmp, dmdσ_tmp)
						σ[k] = tmp
						return dndσ_tmp[i][j]
					}, σ[k], 1e-6)
					chk.AnaNum(tst, io.Sf("d²n[%d]/dσ[%d]dσ[%d]", i, j, k), dtol2, d2ndσdσ[i][j][k], dnum, dver)
				}
			}
		}

		// SMP derivs
		io.Pfpink("\nSMP derivs\n")
		d2mdσdσ_ := la.MatAlloc(3, 3)
		dndσ_ := la.MatAlloc(3, 3)
		dmdσ_ := make([]float64, 3)
		n_ := make([]float64, 3)
		m_ := SmpDerivs(d2mdσdσ_, dndσ_, dmdσ_, n_, σ, b)
		chk.Scalar(tst, "m_", 1e-14, m_, m)
		chk.Vector(tst, "n_", 1e-15, n_, n)
		chk.Vector(tst, "dmdσ_", 1e-15, dmdσ_, dmdσ)
		chk.Matrix(tst, "dndσ_", 1e-13, dndσ_, dndσ)
		chk.Matrix(tst, "d2mdσdσ_", 1e-13, d2mdσdσ_, d2mdσdσ)
	}
}

func Test_geninvs02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("geninvs02")

	b := 0.5
	tol := 1e-10
	dtol := 1e-7
	dtol2 := 1e-7
	dver := chk.Verbose

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
		d2mdσdσ := la.MatAlloc(3, 3)
		dndσ := la.MatAlloc(3, 3)
		dmdσ := make([]float64, 3)
		n := make([]float64, 3)
		m := SmpDerivs(d2mdσdσ, dndσ, dmdσ, n, σ, b)

		// SMP invariants
		p, q, err := GenInvs(σ, n, 1)
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
		chk.Scalar(tst, "p", 1e-14, p, norm_proj)
		chk.Scalar(tst, "q", 1e-13, q, q_)

		// dt/dσ
		var tmp float64
		n_tmp := make([]float64, 3)
		tvec_tmp := make([]float64, 3)
		dtdσ := la.MatAlloc(3, 3)
		GenTvecDeriv1(dtdσ, σ, n, dndσ)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, σ[j] = σ[j], x
					SmpUnitDirector(n_tmp, σ, b)
					GenTvec(tvec_tmp, σ, n_tmp)
					σ[j] = tmp
					return tvec_tmp[i]
				}, σ[j], 1e-6)
				chk.AnaNum(tst, io.Sf("dt/dσ[%d][%d]", i, j), dtol, dtdσ[i][j], dnum, dver)
			}
		}

		// d²t/dσdσ
		io.Pfpink("\nd²t/dσdσ\n")
		d2mdσdσ_tmp := la.MatAlloc(3, 3)
		dndσ_tmp := la.MatAlloc(3, 3)
		dmdσ_tmp := make([]float64, 3)
		dtdσ_tmp := la.MatAlloc(3, 3)
		d2ndσdσ := utl.Deep3alloc(3, 3, 3)
		SmpUnitDirectorDeriv2(d2ndσdσ, σ, n, dmdσ, b, m, d2mdσdσ, dndσ)
		dtol2_tmp := dtol2
		if idxA == 5 {
			dtol2 = 1e-6
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				for k := 0; k < 3; k++ {
					dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
						tmp, σ[k] = σ[k], x
						SmpDerivs(d2mdσdσ_tmp, dndσ_tmp, dmdσ_tmp, n_tmp, σ, b)
						GenTvecDeriv1(dtdσ_tmp, σ, n_tmp, dndσ_tmp)
						σ[k] = tmp
						return dtdσ_tmp[i][j]
					}, σ[k], 1e-6)
					dana := GenTvecDeriv2(i, j, k, σ, dndσ, d2ndσdσ[i][j][k])
					chk.AnaNum(tst, io.Sf("d²t[%d]/dσ[%d]dσ[%d]", i, j, k), dtol2, dana, dnum, dver)
				}
			}
		}
		dtol2 = dtol2_tmp

		// first order derivatives
		dpdσ := make([]float64, 3)
		dqdσ := make([]float64, 3)
		p_, q_, err := GenInvsDeriv1(dpdσ, dqdσ, σ, n, dndσ, 1)
		if err != nil {
			chk.Panic("%v", err)
		}
		chk.Scalar(tst, "p", 1e-17, p, p_)
		chk.Scalar(tst, "q", 1e-17, q, q_)
		var ptmp, qtmp float64
		io.Pfpink("\ndp/dσ\n")
		for j := 0; j < 3; j++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, σ[j] = σ[j], x
				SmpUnitDirector(n_tmp, σ, b)
				ptmp, _, err = GenInvs(σ, n_tmp, 1)
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
				SmpUnitDirector(n_tmp, σ, b)
				_, qtmp, err = GenInvs(σ, n_tmp, 1)
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

		// second order derivatives
		dpdσ_tmp := make([]float64, 3)
		dqdσ_tmp := make([]float64, 3)
		d2pdσdσ := la.MatAlloc(3, 3)
		d2qdσdσ := la.MatAlloc(3, 3)
		GenInvsDeriv2(d2pdσdσ, d2qdσdσ, σ, n, dpdσ, dqdσ, p, q, dndσ, d2ndσdσ, 1)
		io.Pfpink("\nd²p/dσdσ\n")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, σ[j] = σ[j], x
					SmpDerivs(d2mdσdσ_tmp, dndσ_tmp, dmdσ_tmp, n_tmp, σ, b)
					GenInvsDeriv1(dpdσ_tmp, dqdσ_tmp, σ, n_tmp, dndσ_tmp, 1)
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
				//dnum, _ := num.DerivForward(func(x float64, args ...interface{}) (res float64) {
				//dnum, _ := num.DerivBackward(func(x float64, args ...interface{}) (res float64) {
				dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
					tmp, σ[j] = σ[j], x
					SmpDerivs(d2mdσdσ_tmp, dndσ_tmp, dmdσ_tmp, n_tmp, σ, b)
					GenInvsDeriv1(dpdσ_tmp, dqdσ_tmp, σ, n_tmp, dndσ_tmp, 1)
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
