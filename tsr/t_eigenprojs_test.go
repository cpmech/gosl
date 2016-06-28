// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
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
		λsorted := CheckEigenprojs(a, tolP, tolS, ver)
		io.Pfyel("λsorted = %v\n", λsorted)
		λchk := utl.DblGetSorted(test_λ[idxA])
		chk.Vector(tst, "λchk", 1e-12, λsorted, λchk)

		// check derivatives of eigenprojectors
		io.Pforan("\nderivatives\n")
		CheckEigenprojsDerivs(a, toldP, verdP, EV_ZERO)

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
	CheckEigenprojs(a, tolP, tolS, ver)

	io.Pforan("\nderivatives\n")
	CheckEigenprojsDerivs(a, toldP, verdP, EV_ZERO)
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
		CheckEigenprojs(a, tolP, tolS, ver)

		// check derivatives
		io.Pforan("\nderivatives\n")
		CheckEigenprojsDerivs(a, toldP, verdP, EV_ZERO)
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
		CheckEigenprojs(a, tolP, tolS, ver)

		// check derivatives
		io.Pforan("\nderivatives\n")
		CheckEigenprojsDerivs(a, toldP, verdP, EV_ZERO)
	}
}
