// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/utl"
)

// M_PrincValsNum returns the (sorted, ascending) eigenvalues of tensor 'a' (2nd order symmetric tensor in Mandel's basis)
// using Jacobi rotation.
func M_PrincValsNum(a []float64) (λ0, λ1, λ2 float64, err error) {
	Q := Alloc2()
	A := Alloc2()
	v := make([]float64, 3)
	Man2Ten(A, a)
	_, err = la.Jacobi(Q, v, A)
	if err != nil {
		return
	}
	λ0, λ1, λ2 = v[0], v[1], v[2]
	utl.DblSort3(&λ0, &λ1, &λ2)
	return
}

// M_EigenValsNum returns the eigenvalues of tensor 'a' (2nd order symmetric tensor in Mandel's basis)
// using Jacobi rotation
func M_EigenValsNum(λ, a []float64) (err error) {
	Q := Alloc2()
	A := Alloc2()
	Man2Ten(A, a)
	_, err = la.Jacobi(Q, λ, A)
	return
}

// M_EigenValsVecs returns the eigenvalues and eigenvectors of tensor 'a' (2nd order symmetric tensor in Mandel's basis)
// using Jacobi rotation.
func M_EigenValsVecsNum(Q [][]float64, λ, a []float64) (err error) {
	A := Alloc2()
	Man2Ten(A, a)
	_, err = la.Jacobi(Q, λ, A)
	return
}

// M_AllocEigenprojs allocates new eigenprojectors P[3][ncp].
//  P[0] = {P0_0, P0_1, P0_2, P0_3, P0_4, P0_5}
//  P[1] = {P1_0, P1_1, P1_2, P1_3, P1_4, P1_5}
//  P[2] = {P2_0, P2_1, P2_2, P2_3, P2_4, P2_5}
func M_AllocEigenprojs(ncp int) (P [][]float64) {
	P = make([][]float64, 3)
	for k := 0; k < 3; k++ {
		P[k] = make([]float64, ncp)
	}
	return P
}

// M_EigenValsProjsNum computes the eigenvalues and eigenprojectors of tensor 'a' (2nd order symmetric tensor in Mandel's basis)
// using Jacobi rotation.
func M_EigenValsProjsNum(P [][]float64, λ, a []float64) (err error) {
	Q := Alloc2()
	A := Alloc2()
	Man2Ten(A, a)
	_, err = la.Jacobi(Q, λ, A)
	if err != nil {
		return
	}
	P[0][0] = Q[0][0] * Q[0][0]
	P[0][1] = Q[1][0] * Q[1][0]
	P[0][2] = Q[2][0] * Q[2][0]
	P[0][3] = Q[0][0] * Q[1][0] * SQ2
	P[1][0] = Q[0][1] * Q[0][1]
	P[1][1] = Q[1][1] * Q[1][1]
	P[1][2] = Q[2][1] * Q[2][1]
	P[1][3] = Q[0][1] * Q[1][1] * SQ2
	P[2][0] = Q[0][2] * Q[0][2]
	P[2][1] = Q[1][2] * Q[1][2]
	P[2][2] = Q[2][2] * Q[2][2]
	P[2][3] = Q[0][2] * Q[1][2] * SQ2
	if len(a) == 6 {
		P[0][4] = Q[1][0] * Q[2][0] * SQ2
		P[0][5] = Q[2][0] * Q[0][0] * SQ2
		P[1][4] = Q[1][1] * Q[2][1] * SQ2
		P[1][5] = Q[2][1] * Q[0][1] * SQ2
		P[2][4] = Q[1][2] * Q[2][2] * SQ2
		P[2][5] = Q[2][2] * Q[0][2] * SQ2
	}
	return
}

// M_FixZeroOrRepeated computes the eigenvalues λ and applies pertubation to tensor 'a'
// if there are zero or repeated components
//  Note: 'a' is modified
func M_FixZeroOrRepeated(λ, a []float64, pert, evtol, zero float64) (haspert bool, err error) {

	// debug
	if EV_DEBUG {
		io.PfGreen("\nM_FixZeroOrRepeated:\n")
		io.Pfgreen(" a = %#v\n", a)
	}

	// define: compute eigenvalues
	calceigen := func() (e error) {
		ee := M_EigenValsNum(λ, a)
		if ee != nil {
			return chk.Err(_eigenprojs_err9, ee)
		}
		return
	}

	// define: fix any zero
	fix_any_zero := func() (recompute bool) {
		for i := 0; i < 3; i++ {
			if math.Abs(λ[i]) < zero {
				a[i] -= pert
				recompute = true
				if EV_DEBUG {
					io.Pfgreen(" λ%d zero. a(after)=%v\n", i, a)
				}
			}
		}
		return
	}

	// define: doadd
	doadd := func(i, j int, val, cf float64) {
		if math.Abs(a[i]+val-a[j]) > cf {
			a[i] += val
		}
	}

	// define: fix repeated
	fix_repeated := func() (recompute bool) {
		fac := max(1, la.VecLargest(a, 1))
		cf := evtol * fac
		pt := pert + cf
		if EV_DEBUG {
			io.Pfgreen(" cf=%v, pert=%v, pt=%v\n", cf, pert, pt)
		}
		if math.Abs(λ[0]-λ[1]) < cf && math.Abs(λ[1]-λ[2]) < cf && math.Abs(λ[2]-λ[0]) < cf {
			a[0] -= pt * 2.0
			a[1] -= pt
			a[2] += pt * 2.0
			recompute = true
			if EV_DEBUG {
				io.Pfgreen(" λ0 ≈ λ1 ≈ λ2 (diff=%v, %v, %v)\n  λ=%v\n  a(after)=%v\n", math.Abs(λ[0]-λ[1]), math.Abs(λ[1]-λ[2]), math.Abs(λ[2]-λ[0]), λ, a)
			}
		} else {
			if math.Abs(λ[0]-λ[1]) < cf {
				doadd(0, 2, -pt, cf)
				doadd(1, 2, +pt, cf)
				recompute = true
				if EV_DEBUG {
					io.Pfgreen(" λ0 ≈ λ1 (diff=%v)\n  λ=%v\n  a(after)=%v\n", math.Abs(λ[0]-λ[1]), λ, a)
				}
			}
			if math.Abs(λ[1]-λ[2]) < cf {
				doadd(1, 0, -pt, cf)
				doadd(2, 0, +pt, cf)
				recompute = true
				if EV_DEBUG {
					io.Pfgreen(" λ1 ≈ λ2 (diff=%v)\n  λ=%v\n  a(after)=%v\n", math.Abs(λ[1]-λ[2]), λ, a)
				}
			}
			if math.Abs(λ[2]-λ[0]) < cf {
				doadd(2, 1, -pt, cf)
				doadd(0, 1, +pt, cf)
				recompute = true
				if EV_DEBUG {
					io.Pfgreen(" λ2 ≈ λ0 (diff=%v)\n  λ=%v\n  a(after)=%v\n", math.Abs(λ[2]-λ[0]), λ, a)
				}
			}
		}
		return
	}

	// compute eigenvalues
	err = calceigen()
	if err != nil {
		return
	}

	// fix any zero
	recompute := fix_any_zero()
	if recompute {
		haspert = true
		err = calceigen()
		if err != nil {
			return
		}
	}

	// fix repeated (first call)
	if EV_DEBUG {
		io.PfGreen(" first call to fix_repeated\n")
	}
	recompute = fix_repeated()
	if recompute {
		haspert = true
		err = calceigen()
		if err != nil {
			return
		}

		// fix repeated (second call)
		if EV_DEBUG {
			io.PfGreen(" second call to fix_repeated\n")
		}
		recompute = fix_repeated()
		if recompute {
			err = calceigen()
			if err != nil {
				return
			}

			// check
			if EV_DEBUG {
				io.PfGreen(" last call to fix_repeated\n")
			}
			recompute = fix_repeated()
			if recompute {
				return true, chk.Err(_eigenprojs_err11, a, λ)
			}
		}
	}
	return
}

// M_EigenValsAna returns the eigenvalues of tensor 'a' (2nd order symmetric tensor in Mandel's basis)
// using the analytical formula
func M_EigenValsAna(λ, a []float64, evtol, zero float64) {
	// check for diagonal tensor
	isdiagonal := true
	for i := 3; i < len(a); i++ {
		if math.Abs(a[i]) > zero {
			isdiagonal = false
			break
		}
	}
	// copy eigenvalues
	if isdiagonal {
		λ[0], λ[1], λ[2] = a[0], a[1], a[2]
		return
	}
	// compute eigenvalues
	I1, I2, I3 := M_CharInvs(a)
	d := I1*I1 - 3.0*I2
	if d < 0 {
		d = 0
		if EV_DEBUG {
			io.Pfgrey("eigenprojs.go: M_EigenValsAna: d = %g < 0\n", d)
		}
	}
	fac := max(1, la.VecLargest(a, 1))
	sd := math.Sqrt(d)
	if sd < evtol*fac {
		chk.Panic("eigenprojs.go: M_EigenValsAna: all repeated eigenvalues\n  a=%v\n  λ=%v\n", a, λ)
	}
	nm := 2.0*I1*I1*I1 - 9.0*I1*I2 + 27.0*I3
	if math.Abs(nm) < evtol {
		nm = 0
		if EV_DEBUG {
			io.Pfgrey("eigenprojs.go: M_EigenValsAna: |nm| = %g < %g\n", math.Abs(nm), evtol)
		}
	}
	t := 0.5 * nm / (sd * sd * sd)
	if math.Abs(t) > 1.0 {
		t = fun.Sign(t)
		if EV_DEBUG {
			io.Pfgrey("eigenprojs.go: M_EigenValsAna: |t| = %g > 1\n", math.Abs(t))
		}
	}
	θ := math.Acos(t)
	for k := 0; k < 3; k++ {
		λ[k] = (I1 + 2.0*sd*math.Cos((θ+2.0*math.Pi*(1.0+float64(k)))/3.0)) / 3.0
	}
	if EV_DEBUG {
		io.Pfgrey("eigenprojs.go: M_EigenValsAna:\n  a = %v\n  λ = %v\n", a, λ)
	}
	return
}

// M_EigenProjsAna computes the eigenprojectors of tensor 'a' (2nd order symmetric tensor in Mandel's basis)
// using the analytical formula.
func M_EigenProjsAna(P [][]float64, a, λ []float64, zero float64) (err error) {

	// check eigenvalues
	if math.Abs(λ[0]) < zero || math.Abs(λ[1]) < zero || math.Abs(λ[2]) < zero {
		return chk.Err(_eigenprojs_err1, "M_EigenProjsAna", λ, a, zero)
	}

	// inverse tensor
	ai := make([]float64, len(a))
	_, err = M_Inv(ai, a, MINDET)
	if err != nil {
		return chk.Err(_eigenprojs_err7, "M_EigenProjsAna", err.Error())
	}

	// characteristic invariants
	I1 := λ[0] + λ[1] + λ[2]
	I3 := λ[0] * λ[1] * λ[2]

	// eigenvalues and eigenprojectors
	var dn, cf float64
	for k := 0; k < 3; k++ {
		dn = 2.0*λ[k]*λ[k] - λ[k]*I1 + I3/λ[k]
		cf = λ[k] / dn
		if math.Abs(dn) < EV_DNMIN {
			return chk.Err(_eigenprojs_err4, dn, EV_DNMIN, a, λ, I1, I3, cf)
		}
		for i := 0; i < len(a); i++ {
			P[k][i] = cf * (a[i] + (λ[k]-I1)*Im[i] + (I3/λ[k])*ai[i])
		}
	}
	return
}

// M_EigenProjsDeriv returns the derivatives of the eigenprojectors w.r.t its defining tensor
// using the analytical formula.
//  Input:  a      -- (perturbed) tensor 'a' (in Mandel basis)
//          λ      -- (shifted) eigenvalues of 'a'
//          P      -- eigenprojectors of 'a'
//          zero   -- tolerance to assume zero eigenvalues
//          usefdm -- use finite difference method instead of analytical formula
//  Output: dPda   -- the derivatives of P w.r.t 'a'
func M_EigenProjsDeriv(dPda [][][]float64, a, λ []float64, P [][]float64, zero float64) (err error) {

	// check eigenvalues
	if math.Abs(λ[0]) < zero || math.Abs(λ[1]) < zero || math.Abs(λ[2]) < zero {
		return chk.Err(_eigenprojs_err5, "M_EigenProjsDeriv", λ, zero)
	}

	// derivative of inverse tensor
	ai := make([]float64, len(a))
	_, err = M_Inv(ai, a, MINDET)
	if err != nil {
		return chk.Err(_eigenprojs_err7, "M_EigenProjsDeriv", err.Error())
	}
	M_InvDeriv(dPda[0], ai) // dPda[0] := daida

	// characteristic invariants
	I1 := λ[0] + λ[1] + λ[2]
	I3 := λ[0] * λ[1] * λ[2]

	// alpha coefficients
	α0 := 2.0*λ[0]*λ[0] - I1*λ[0] + I3/λ[0]
	α1 := 2.0*λ[1]*λ[1] - I1*λ[1] + I3/λ[1]
	α2 := 2.0*λ[2]*λ[2] - I1*λ[2] + I3/λ[2]
	if math.Abs(α0) < EV_ALPMIN || math.Abs(α1) < EV_ALPMIN || math.Abs(α2) < EV_ALPMIN {
		return chk.Err(_eigenprojs_err2, α0, α1, α2, λ)
	}

	// compute derivatives
	ncp := len(ai)
	var daida_ij float64
	for i := 0; i < ncp; i++ {
		for j := 0; j < ncp; j++ {
			daida_ij = dPda[0][i][j]
			dPda[0][i][j] = (λ[0]*IIm[i][j] + I3*daida_ij + ((I3/(λ[0]*λ[0])-λ[0])*P[0][i]*P[0][j] + (I3/(λ[1]*λ[1])-λ[0])*P[1][i]*P[1][j] + (I3/(λ[2]*λ[2])-λ[0])*P[2][i]*P[2][j])) / α0
			dPda[1][i][j] = (λ[1]*IIm[i][j] + I3*daida_ij + ((I3/(λ[0]*λ[0])-λ[1])*P[0][i]*P[0][j] + (I3/(λ[1]*λ[1])-λ[1])*P[1][i]*P[1][j] + (I3/(λ[2]*λ[2])-λ[1])*P[2][i]*P[2][j])) / α1
			dPda[2][i][j] = (λ[2]*IIm[i][j] + I3*daida_ij + ((I3/(λ[0]*λ[0])-λ[2])*P[0][i]*P[0][j] + (I3/(λ[1]*λ[1])-λ[2])*P[1][i]*P[1][j] + (I3/(λ[2]*λ[2])-λ[2])*P[2][i]*P[2][j])) / α2
		}
	}
	return
}

// M_EigenProjsDerivNum returns the derivatives of the eigenprojectors w.r.t its defining tensor
// using the finite differences method.
//  Input:
//    a -- tensor in Mandel basis
//    h -- step size for finite differences
//  Output:
//    dPda -- derivatives [3][ncp][ncp]
func M_EigenProjsDerivNum(dPda [][][]float64, a []float64, h float64) (err error) {
	ncp := len(a)
	λ := make([]float64, 3)
	P := la.MatAlloc(3, ncp)
	Q := Alloc2()
	A := Alloc2()
	q2p := func(k int) {
		switch k {
		case 0:
			P[0][0] = Q[0][0] * Q[0][0]
			P[0][1] = Q[1][0] * Q[1][0]
			P[0][2] = Q[2][0] * Q[2][0]
			P[0][3] = Q[0][0] * Q[1][0] * SQ2
			if ncp == 6 {
				P[0][4] = Q[1][0] * Q[2][0] * SQ2
				P[0][5] = Q[2][0] * Q[0][0] * SQ2
			}
		case 1:
			P[1][0] = Q[0][1] * Q[0][1]
			P[1][1] = Q[1][1] * Q[1][1]
			P[1][2] = Q[2][1] * Q[2][1]
			P[1][3] = Q[0][1] * Q[1][1] * SQ2
			if ncp == 6 {
				P[1][4] = Q[1][1] * Q[2][1] * SQ2
				P[1][5] = Q[2][1] * Q[0][1] * SQ2
			}
		case 2:
			P[2][0] = Q[0][2] * Q[0][2]
			P[2][1] = Q[1][2] * Q[1][2]
			P[2][2] = Q[2][2] * Q[2][2]
			P[2][3] = Q[0][2] * Q[1][2] * SQ2
			if ncp == 6 {
				P[2][4] = Q[1][2] * Q[2][2] * SQ2
				P[2][5] = Q[2][2] * Q[0][2] * SQ2
			}
		}
	}
	var tmp float64
	failed := false
	for k := 0; k < 3; k++ {
		for i := 0; i < ncp; i++ {
			for j := 0; j < ncp; j++ {
				dPda[k][i][j], _ = num.DerivCentral(func(x float64, args ...interface{}) float64 {
					tmp, a[j] = a[j], x
					defer func() { a[j] = tmp }()
					Man2Ten(A, a)
					_, err = la.Jacobi(Q, λ, A)
					if err != nil {
						failed = true
						return 0
					}
					q2p(k)
					return P[k][i]
				}, a[j], h)
				if failed {
					return
				}
			}
		}
	}
	return
}

// auxiliary /////////////////////////////////////////////////////////////////////////////////////////

func CheckEigenprojs(a []float64, usenum bool, tolP, tolS float64, ver bool, evtol, zero float64, dosort bool) (λsorted []float64, Psorted [][]float64) {

	// compute eigenvalues and eigenprojectors
	ncp := len(a)
	λ := make([]float64, 3)
	P := la.MatAlloc(3, ncp)
	if usenum {
		err := M_EigenValsProjsNum(P, λ, a)
		if err != nil {
			chk.Panic("eigenprojs.go: CheckEigenprojs failed:\n %v", err.Error())
		}
	} else {
		M_EigenValsAna(λ, a, evtol, zero)
		err := M_EigenProjsAna(P, a, λ, zero)
		if err != nil {
			chk.Panic("eigenprojs.go: CheckEigenprojs failed:\n %v", err.Error())
		}
	}

	// print projectors
	if ver {
		la.PrintVec("P0", P[0], "%14.6e", false)
		la.PrintVec("P1", P[1], "%14.6e", false)
		la.PrintVec("P2", P[2], "%14.6e", false)
	}

	// check P dot P
	PdotP := make([]float64, ncp)
	Z := make([]float64, ncp)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			err := M_Dot(PdotP, P[i], P[j], 1e-14)
			if err != nil {
				chk.Panic("%v", err)
			}
			if i == j {
				diff := la.VecMaxDiff(PdotP, P[i])
				if diff > tolP {
					chk.Panic("eigenprojs.go: CheckEigenprojs failed: P%d dot P%d != P%d (diff=%g)", i, j, i, diff)
				} else if ver {
					io.Pf("P%d dot P%d == P%d [1;32mOK[0m (diff=%g)\n", i, j, i, diff)
				}
			} else {
				diff := la.VecMaxDiff(PdotP, Z)
				if diff > tolP {
					chk.Panic("eigenprojs.go: CheckEigenprojs failed: P%d dot P%d !=  0 (diff=%g)", i, j, diff)
				} else if ver {
					io.Pf("P%d dot P%d ==  0 [1;32mOK[0m (diff=%g)\n", i, j, diff)
				}
			}
		}
	}

	// check sum of eigenprojectors
	sumP := make([]float64, ncp)
	for k := 0; k < 3; k++ {
		for i := 0; i < ncp; i++ {
			sumP[i] += P[k][i]
		}
	}
	diff := la.VecMaxDiff(sumP, Im[:ncp])
	if diff > tolP {
		chk.Panic("eigenprojs.go: CheckEigenprojs failed: sumP != I (diff=%g)", diff)
	} else if ver {
		io.Pf("sum(P) [1;32mOK[0m (diff=%g)\n", diff)
	}

	// check spectral decomposition
	as := make([]float64, len(a))
	for k := 0; k < 3; k++ {
		for i := 0; i < len(a); i++ {
			as[i] += λ[k] * P[k][i]
		}
	}
	diff = la.VecMaxDiff(as, a)
	if diff > tolS {
		chk.Panic("eigenprojs.go: CheckEigenprojs failed: a(spectral) != a (diff=%g)", diff)
	} else if ver {
		io.Pf("a(spectral) == a [1;32mOK[0m (diff=%g)\n", diff)
	}

	// sort eigenvalues and eigenprojectors
	if dosort {
		λsorted = make([]float64, 3)
		Psorted = la.MatAlloc(3, ncp)
		I := []int{0, 1, 2}
		var err error
		I, λsorted, _, _, err = utl.SortQuadruples(I, λ, nil, nil, "x")
		if err != nil {
			chk.Panic("%v", err)
		}
		for i, k := range I {
			copy(Psorted[i], P[k])
		}
	}
	return
}

func CheckEigenprojsDerivs(usenum bool, a []float64, tol float64, ver bool, evtol, zero float64) (dPda, dPda_num [][][]float64) {

	// compute eigenvalues and eigenprojectors
	ncp := len(a)
	λ := make([]float64, 3)
	P := la.MatAlloc(3, ncp)
	docalc := func() {
		if usenum {
			err := M_EigenValsProjsNum(P, λ, a)
			if err != nil {
				chk.Panic("eigenprojs.go: CheckEigenprojsDerivs failed:\n %v", err.Error())
			}
		} else {
			M_EigenValsAna(λ, a, evtol, zero)
			//io.Pforan("a = %v\n", a)
			//io.Pforan("λ = %v\n", λ)
			err := M_EigenProjsAna(P, a, λ, zero)
			if err != nil {
				chk.Panic("eigenprojs.go: CheckEigenprojsDerivs failed:\n %v", err.Error())
			}
		}
	}

	// compute derivatives of eigenprojectors
	docalc()
	dPda = utl.Deep3alloc(3, ncp, ncp)
	err := M_EigenProjsDeriv(dPda, a, λ, P, zero)
	if err != nil {
		chk.Panic("%v", err)
	}

	// compute numerical derivatives of eigenprojectors
	dPda_num = utl.Deep3alloc(3, ncp, ncp)
	err = M_EigenProjsDerivNum(dPda_num, a, 1e-6)
	if err != nil {
		chk.Panic("%v", err)
	}

	// check
	var tmp float64
	has_error := false
	for k := 0; k < 3; k++ {
		for i := 0; i < ncp; i++ {
			for j := 0; j < ncp; j++ {
				dnum := num.DerivFwd(func(x float64, args ...interface{}) (res float64) {
					//dnum := num.DerivCen(func(x float64, args ...interface{}) (res float64) {
					tmp, a[j] = a[j], x
					docalc()
					a[j] = tmp
					return P[k][i]
				}, a[j], 1e-6)
				err := chk.PrintAnaNum(io.Sf("dP%d[%d]/da[%d]", k, i, j), tol, dPda[k][i][j], dnum, ver)
				if err != nil {
					has_error = true
				}
			}
		}
		if ver {
			io.Pf("\n")
		}
	}
	if has_error {
		chk.Panic(_eigenprojs_err8)
	}
	return
}

// print_eigenvecs prints the eigenvectors in matrix Q
func print_eigenvecs(Q [][]float64) {
	io.Pforan("Q0 = [%v %v %v]\n", Q[0][0], Q[1][0], Q[2][0])
	io.Pforan("Q1 = [%v %v %v]\n", Q[0][1], Q[1][1], Q[2][1])
	io.Pforan("Q2 = [%v %v %v]\n", Q[0][2], Q[1][2], Q[2][2])
}

// error messages
var (
	_eigenprojs_err1  = "eigenprojs.go: %s: λ=%v of tensor a=%v is too small (zero=%v)"
	_eigenprojs_err2  = "eigenprojs.go: M_EigenProjsDeriv: α=[%v %v %v] coefficients must be non-zero (λ=%v)"
	_eigenprojs_err3  = "eigenprojs.go: %s: cannot handle repeated eigenvalues λ=%v of tensor a=%v (tol=%v)"
	_eigenprojs_err4  = "eigenprojs.go: M_EigenProjsAna: cannot compute eigenprojectors since:\n  den = 2*λ² - λ*I1 + I3/λ = %g < %g\n  a=%v\n  λ=%v\n  I1=%v  I3=%v\n  cf=%v"
	_eigenprojs_err5  = "eigenprojs.go: %s: λ=%v of tensor is too small (zero=%v)"
	_eigenprojs_err6  = "eigenprojs.go: %s: there are still repeated eigenvalues after perturbation\n  A_perturbed = %v\n  L = %v"
	_eigenprojs_err7  = "eigenprojs.go: %s: cannot compute inverse tensor:\n %v"
	_eigenprojs_err8  = "eigenprojs.go: CheckEigenprojsDerivs failed"
	_eigenprojs_err9  = "eigenprojs.go: M_FixZeroOrRepeated failed:\n %v"
	_eigenprojs_err10 = "eigenprojs.go: M_FixZeroOrRepeated: inconsistency with a and λ found (not all zero components):\n  a=%v\n  λ=%v\n"
	_eigenprojs_err11 = "eigenprojs.go: M_FixZeroOrRepeated failed to fix repeated values:\n  a=%v\n  λ=%v\n"
)
