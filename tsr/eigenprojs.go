// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"

	"github.com/cpmech/gosl/chk"
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
	utl.Sort3(&λ0, &λ1, &λ2)
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

// M_EigenProjsDerivAna returns the derivatives of the eigenprojectors w.r.t its defining tensor
// using the analytical formula.
//  Input:
//    a -- (perturbed) tensor 'a' (in Mandel basis)
//    λ -- eigenvalues of 'a'
//    P -- eigenprojectors of 'a'
//  Output:
//    dPda -- the derivatives of P w.r.t 'a'
func M_EigenProjsDerivAna(dPda [][][]float64, a, λ []float64, P [][]float64) (err error) {

	// check eigenvalues
	if math.Abs(λ[0]) < EV_ZERO || math.Abs(λ[1]) < EV_ZERO || math.Abs(λ[2]) < EV_ZERO {
		return chk.Err(_eigenprojs_err5, "M_EigenProjsDeriv", λ, EV_ZERO)
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

// M_EigenProjsDerivAuto computes the derivatives of the eigenprojectors of tensor a
// w.r.t. to itself by automatically calling the numerical or the analytical formulae
// depending on whether the eigenvalues are zero/repeated or not
//
//  Note: this function should work for non-perturbed tensors with zero/repeated eigenvalues.
//
//  Input:
//    a -- tensor 'a' (in Mandel basis)
//    λ -- eigenvalues of 'a'
//    P -- eigenprojectors of 'a'
//  Output:
//    dPda -- the derivatives of P w.r.t 'a'
func M_EigenProjsDerivAuto(dPda [][][]float64, a, λ []float64, P [][]float64) (err error) {
	if math.Abs(λ[0]) < EV_ZERO || math.Abs(λ[1]) < EV_ZERO || math.Abs(λ[2]) < EV_ZERO {
		return M_EigenProjsDerivNum(dPda, a, 1e-6)
	}
	//io.Pf("%v < %v\n", math.Abs(λ[0]-λ[1]), EV_EQUAL*max(λ[0], λ[1]))
	//io.Pf("%v < %v\n", math.Abs(λ[1]-λ[2]), EV_EQUAL*max(λ[1], λ[2]))
	//io.Pf("%v < %v\n", math.Abs(λ[2]-λ[0]), EV_EQUAL*max(λ[2], λ[0]))
	if math.Abs(λ[0]-λ[1]) < EV_EQUAL*math.Abs(max(λ[0], λ[1])) ||
		math.Abs(λ[1]-λ[2]) < EV_EQUAL*math.Abs(max(λ[1], λ[2])) ||
		math.Abs(λ[2]-λ[0]) < EV_EQUAL*math.Abs(max(λ[2], λ[0])) {
		return M_EigenProjsDerivNum(dPda, a, 1e-6)
	}
	return M_EigenProjsDerivAna(dPda, a, λ, P)
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
	_eigenprojs_err2  = "eigenprojs.go: M_EigenProjsDerivAna: α=[%v %v %v] coefficients must be non-zero (λ=%v)"
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
