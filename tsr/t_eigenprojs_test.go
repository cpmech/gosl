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

func Test_eigenp00(tst *testing.T) {

	//verbose()
	chk.PrintTitle("eigenp00")

	ncp := 6
	P := M_AllocEigenprojs(ncp)
	chk.Vector(tst, "P0", 1e-17, P[0], []float64{0, 0, 0, 0, 0, 0})
	chk.Vector(tst, "P1", 1e-17, P[1], []float64{0, 0, 0, 0, 0, 0})
	chk.Vector(tst, "P2", 1e-17, P[2], []float64{0, 0, 0, 0, 0, 0})
}

func Test_eigenp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("eigenp01")

	// constants
	tolP := 1e-14        // eigenprojectors
	tolS := 1e-13        // spectral decomposition
	toldP := 1e-9        // derivatives of eigenprojectors
	ver := chk.Verbose   // check P verbose
	verdP := chk.Verbose // check dPda verbose

	// run test
	nd := test_nd
	for idxA := 0; idxA < len(test_nd); idxA++ {
		//for idxA := 10; idxA < 11; idxA++ {
		//for idxA := 11; idxA < 12; idxA++ {
		//for idxA := 12; idxA < 13; idxA++ {

		// tensor and eigenvalues
		A := test_AA[idxA]
		a := M_Alloc2(nd[idxA])
		Ten2Man(a, A)
		io.PfYel("\n\ntst # %d ###################################################################################\n", idxA)
		io.Pfblue2("a = %v\n", a)
		io.Pfblue2("λ = %v\n", test_λ[idxA])

		// check eigenprojectors
		io.Pforan("\neigenprojectors\n")
		λsorted := checkEigenprojs(tst, a, tolP, tolS, ver)
		io.Pfyel("λsorted = %v\n", λsorted)
		λchk := utl.GetSorted(test_λ[idxA])
		chk.Vector(tst, "λchk", 1e-12, λsorted, λchk)

		// check derivatives of eigenprojectors
		io.Pforan("\nderivatives\n")
		checkEigenprojsDerivs(tst, a, toldP, verdP, EV_ZERO)

	}
}

func Test_eigenp02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("eigenp02")

	// constants
	tolP := 1e-14        // eigenprojectors
	tolS := 1e-14        // spectral decomposition
	toldP := 1e-14       // derivatives of eigenprojectors
	ver := chk.Verbose   // check P verbose
	verdP := chk.Verbose // check dPda verbose

	// set tensor
	ϵ := 1e-10
	s := 1.0
	a := []float64{1, 1, 2, ϵ, 0, 0}
	for i := 0; i < len(a); i++ {
		a[i] *= s
	}

	// run test
	io.Pforan("\neigenprojectors\n")
	checkEigenprojs(tst, a, tolP, tolS, ver)

	io.Pforan("\nderivatives\n")
	checkEigenprojsDerivs(tst, a, toldP, verdP, EV_ZERO)
}

func Test_eigenp03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("eigenp03")

	// constants
	tolP := 1e-14        // eigenprojectors
	tolS := 1e-13        // spectral decomposition
	toldP := 1e-14       // derivatives of eigenprojectors
	ver := chk.Verbose   // check P verbose
	verdP := chk.Verbose // check dPda verbose

	// run
	s := 100.0
	for _, ϵ := range []float64{1e-7, 1e-2, 1e-2} {

		// set tensor
		a := []float64{s, s + ϵ/4.0, s + 3.0*ϵ/4.0, 0, SQ3 * ϵ / 4.0, 0}
		io.PfYel("\n\ntst ####################################################################################\n")
		io.Pfblue2("ϵ = %v\n", ϵ)
		io.Pfblue2("a = %v\n", a)

		// run test
		io.Pforan("\neigenprojectors (num)\n")
		checkEigenprojs(tst, a, tolP, tolS, ver)

		// check derivatives
		io.Pforan("\nderivatives\n")
		checkEigenprojsDerivs(tst, a, toldP, verdP, EV_ZERO)
	}
}

func Test_eigenp04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("eigenp04")

	// constants
	tolP := 1e-14        // eigenprojectors
	tolS := 1e-14        // spectral decomposition
	toldP := 1e-6        // derivatives of eigenprojectors
	ver := chk.Verbose   // check P verbose
	verdP := chk.Verbose // check dPda verbose

	// tensor main value
	//s := 400.0
	//s := 0.0
	s := 1.0

	// run for a number of δ
	for _, i := range []int{3, 5, 7, 9} {

		// set tensor
		δ := math.Pow(10.0, -float64(i))
		a := []float64{s, s + δ, s - δ, 0, 0, 0}
		io.PfYel("\n\ntst ####################################################################################\n")
		io.Pfblue2("δ = %v\n", δ)
		io.Pfblue2("a = %v\n", a)

		// run test
		io.Pforan("\neigenprojectors\n")
		checkEigenprojs(tst, a, tolP, tolS, ver)

		// check derivatives
		io.Pforan("\nderivatives\n")
		checkEigenprojsDerivs(tst, a, toldP, verdP, EV_ZERO)
	}
}

// auxiliary /////////////////////////////////////////////////////////////////////////////////////////

// checkEigenprojs checks eigen projectors
func checkEigenprojs(tst *testing.T, a []float64, tolP, tolS float64, ver bool) (λsorted []float64) {

	// compute eigenvalues and eigenprojectors
	ncp := len(a)
	λ := make([]float64, 3)
	P := la.MatAlloc(3, ncp)
	err := M_EigenValsProjsNum(P, λ, a)
	if err != nil {
		chk.Panic("eigenprojs.go: CheckEigenprojs failed:\n %v", err.Error())
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

	// sort eigenvalues
	λsorted = make([]float64, 3)
	I := []int{0, 1, 2}
	I, λsorted, _, _, err = utl.SortQuadruples(I, λ, nil, nil, "x")
	if err != nil {
		chk.Panic("%v", err)
	}
	return
}

// checkEigenprojsDerivs checks the derivatives of eigen projectors w.r.t defining tensor
func checkEigenprojsDerivs(tst *testing.T, a []float64, tol float64, ver bool, zero float64) {

	// variables
	ncp := len(a)
	λ := make([]float64, 3)
	P := la.MatAlloc(3, ncp)

	// compute derivatives of eigenprojectors
	dPda := utl.Deep3alloc(3, ncp, ncp)
	err := M_EigenProjsDerivAuto(dPda, a, λ, P)
	if err != nil {
		tst.Errorf("M_EigenProjsDerivAuto failed:\n%v\n", err)
		return
	}

	// check
	for k := 0; k < 3; k++ {
		chk.DerivVecVec(tst, io.Sf("dP%d/da", k), tol, dPda[k], a, 1e-6, ver, func(f, x []float64) error {
			err := M_EigenValsProjsNum(P, λ, x)
			copy(f, P[k])
			return err
		})
	}
}
