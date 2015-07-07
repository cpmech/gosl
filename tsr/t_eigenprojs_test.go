// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
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

func Test_eigenp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("eigenp01")

	cmpλ := 1e-12
	cmpP := 1e-12
	tolP := 1e-12
	tolS := 1e-12
	toldP := 1e-6
	toldP2 := 1e-7
	verdP := chk.Verbose
	ver := chk.Verbose

	// run test
	nd := test_nd
	for idxA := 0; idxA < len(test_nd); idxA++ {
		//for idxA := 10; idxA < 11; idxA++ {
		//for idxA := 11; idxA < 12; idxA++ {
		//for idxA := 12; idxA < 13; idxA++ {

		// fix tolerances
		cmpλ_, cmpP_, tolP_, tolS_, toldP_ := cmpλ, cmpP, tolP, tolS, toldP
		switch idxA {
		case 5:
			toldP = 1e-4
		case 10:
			toldP = 1e-4
		case 11:
			toldP = 0.00021
		case 12:
			tolP, tolS, cmpP, toldP = 1e-9, 1e-9, 1e-10, 0.017
		}

		// tensor and eigenvalues
		A := test_AA[idxA]
		a := M_Alloc2(nd[idxA])
		Ten2Man(a, A)
		io.PfYel("\n\ntst # %d ###################################################################################\n", idxA)
		io.Pfblue2("a = %v\n", a)
		io.Pfblue2("λ = %v\n", test_λ[idxA])

		// perturbation
		λper := make([]float64, 3)
		haspert, err := M_FixZeroOrRepeated(λper, a, EV_PERT, EV_EVTOL, EV_ZERO)
		if haspert {
			io.Pfyel("a(pert) = %v\n", a)
			io.Pfyel("λ(pert) = %v\n", λper)
		}
		if err != nil {
			chk.Panic("%v", err.Error())
		}

		// check analytical eigenprojectors
		io.Pforan("\nana\n")
		λana, Pana := CheckEigenprojs(a, false, tolP, tolS, ver, EV_EVTOL, EV_ZERO, true)
		io.Pfyel("λana = %v\n", λana)
		if !haspert {
			λchk := utl.DblGetSorted(test_λ[idxA])
			chk.Vector(tst, "λchk", 1e-12, λana, λchk)
		}

		// check numerical eigenprojectors
		io.Pforan("\nnum\n")
		λnum, Pnum := CheckEigenprojs(a, true, tolP, tolS, ver, EV_EVTOL, EV_ZERO, true)
		io.Pfyel("λnum = %v\n", λnum)

		// compare ana-num
		io.Pforan("\nana-num\n")
		chk.Vector(tst, "λana-λnum", cmpλ, λana, λnum)
		chk.Matrix(tst, "Pana-Pnum", cmpP, Pana, Pnum)

		// check derivatives of analytical eigenprojectors
		io.Pforan("\nderivatives (anaP)\n")
		CheckEigenprojsDerivs(false, a, toldP, verdP, EV_EVTOL, EV_ZERO)

		// check derivatives of numerical eigenprojectors
		io.Pforan("\nderivatives (numP)\n")
		dPda, dPda_num := CheckEigenprojsDerivs(true, a, toldP, verdP, EV_EVTOL, EV_ZERO)

		// compare derivatives
		chk.Deep3(tst, "dPda", toldP2, dPda, dPda_num)

		// restore tolerances
		cmpλ, cmpP, tolP, tolS, toldP = cmpλ_, cmpP_, tolP_, tolS_, toldP_
	}
}

func Test_eigenp02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("eigenp02")

	// constants
	tolP := 1e-12
	tolS := 1e-12
	toldP := 1e-4
	toldP2 := 1e-9
	verdP := chk.Verbose
	ver := chk.Verbose

	// set tensor
	ϵ := 1e-10
	s := 1.0
	a := []float64{1, 1, 2, ϵ, 0, 0}
	for i := 0; i < len(a); i++ {
		a[i] *= s
	}

	// eigenvalues
	λ := make([]float64, 3)
	err := M_EigenValsNum(λ, a)
	if err != nil {
		chk.Panic("%v", err.Error())
	}
	io.Pfblue2("a = %v\n", a)
	io.Pfblue2("λ = %v\n", λ)

	// perturbation
	if true {
		λper := make([]float64, 3)
		haspert, err := M_FixZeroOrRepeated(λper, a, EV_PERT, EV_EVTOL, EV_ZERO)
		if haspert {
			io.Pfyel("a(pert) = %v\n", a)
			io.Pfyel("λ(pert) = %v\n", λper)
		}
		if err != nil {
			chk.Panic("%v", err.Error())
		}
	}

	// run test
	io.Pforan("\neigenprojectors\n")
	CheckEigenprojs(a, false, tolP, tolS, ver, EV_EVTOL, EV_ZERO, true)

	io.Pforan("\nderivatives\n")
	dPda, dPda_num := CheckEigenprojsDerivs(false, a, toldP, verdP, EV_EVTOL, EV_ZERO)

	// compare derivatives
	chk.Deep3(tst, "dPda", toldP2, dPda, dPda_num)
}

func Test_eigenp03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("eigenp03")

	// constants
	tolP := 1e-15
	tolS := 1e-13
	toldP := 0.008
	toldP2 := 0.002
	verdP := chk.Verbose
	ver := chk.Verbose

	// set tensor
	ϵ := 1e-7
	//ϵ := 1e-2 // this causes problmes with EVTOL=1e-5
	//ϵ := 1e-1
	s := 100.0
	a := []float64{s, s + ϵ/4.0, s + 3.0*ϵ/4.0, 0, SQ3 * ϵ / 4.0, 0}

	// eigenvalues
	λ := make([]float64, 3)
	err := M_EigenValsNum(λ, a)
	if err != nil {
		chk.Panic("%v", err.Error())
	}
	io.Pfblue2("a = %v\n", a)
	io.Pfblue2("λ = %v\n", λ)

	// perturbation
	if true {
		λper := make([]float64, 3)
		haspert, err := M_FixZeroOrRepeated(λper, a, EV_PERT, EV_EVTOL, EV_ZERO)
		if haspert {
			io.Pfyel("a(pert) = %v\n", a)
			io.Pfyel("λ(pert) = %v\n", λper)
		}
		if err != nil {
			chk.Panic("%v", err.Error())
		}
	}

	// run test (ana)
	if false {
		io.Pforan("\neigenprojectors (ana)\n")
		CheckEigenprojs(a, false, tolP, tolS, ver, EV_EVTOL, EV_ZERO, true)

		io.Pforan("\nderivatives (anaP)\n")
		CheckEigenprojsDerivs(false, a, toldP, verdP, EV_EVTOL, EV_ZERO)
	}

	// run test
	if true {
		io.Pforan("\neigenprojectors (num)\n")
		CheckEigenprojs(a, true, tolP, tolS, ver, EV_EVTOL, EV_ZERO, true)

		io.Pforan("\nderivatives (numP)\n")
		dPda, dPda_num := CheckEigenprojsDerivs(true, a, toldP, verdP, EV_EVTOL, EV_ZERO)

		// compare derivatives
		chk.Deep3(tst, "dPda", toldP2, dPda, dPda_num)
	}
}

func Test_eigenp04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("eigenp04")

	// constants
	tolP := 1e-15
	tolS := 1e-15
	toldP := 1e-4
	toldP2 := 1e-6
	verdP := chk.Verbose
	ver := chk.Verbose

	// tensor main value
	//s := 400.0
	//s := 0.0
	s := 1.0

	// run for a number of δ
	for _, i := range []int{3, 5, 7, 9} {

		// fix tolerances
		toldP_ := toldP
		switch i {
		case 5:
			toldP = 0.1582
			toldP2 = 0.105
		}

		// noise
		δ := math.Pow(10.0, -float64(i))
		io.PfYel("\n\nδ = %v ##################################################################################\n", δ)

		// tensor and eigenvalues
		a := []float64{s, s + δ, s - δ, 0, 0, 0}
		λ := make([]float64, 3)
		err := M_EigenValsNum(λ, a)
		if err != nil {
			chk.Panic("%v", err.Error())
		}
		io.Pfblue2("a = %v\n", a)
		io.Pfblue2("λ = %v\n", λ)

		// perturbation
		if true {
			λper := make([]float64, 3)
			haspert, err := M_FixZeroOrRepeated(λper, a, EV_PERT, EV_EVTOL, EV_ZERO)
			if haspert {
				io.Pfyel("a(pert) = %v\n", a)
				io.Pfyel("λ(pert) = %v\n", λper)
			}
			if err != nil {
				chk.Panic("%v", err.Error())
			}
		}

		// run test (ana)
		if false {
			io.Pforan("\neigenprojectors (ana)\n")
			CheckEigenprojs(a, false, tolP, tolS, ver, EV_EVTOL, EV_ZERO, true)

			io.Pforan("\nderivatives (anaP)\n")
			CheckEigenprojsDerivs(false, a, toldP, verdP, EV_EVTOL, EV_ZERO)
		}

		// run test
		if true {
			io.Pforan("\neigenprojectors (num)\n")
			CheckEigenprojs(a, true, tolP, tolS, ver, EV_EVTOL, EV_ZERO, true)

			io.Pforan("\nderivatives (numP)\n")
			dPda, dPda_num := CheckEigenprojsDerivs(true, a, toldP, verdP, EV_EVTOL, EV_ZERO)

			// compare derivatives
			chk.Deep3(tst, "dPda", toldP2, dPda, dPda_num)
		}

		// restore tolerances
		toldP = toldP_
	}
}

func Test_eigenp05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("eigenp05")

	// constants
	tolP := 1e-14
	tolS := 1e-15
	toldP := 1e-6
	toldP2 := 1e-6
	verdP := chk.Verbose
	ver := chk.Verbose

	// tensor
	idx := 3
	var a []float64
	switch idx {
	case 1:
		a = []float64{-3.4520830204845048, -3.452083020484557, -2.4945423175657706, 17.714258863230363, -2.242430664369921, -2.242430664370262}
	case 2:
		a = []float64{19538.58315173556, 18558.787254622286, 18558.787254622286, 0}
		tolP = 1e-8
		tolS = 1e-4
		toldP = 0.0025
	case 3:
		a = []float64{0.9985833333333334, 0.9995833333333333, 0.9995833333333333, 0, 0, 0}
		tolP = 1e-9
		tolS = 1e-9
		toldP = 0.023
	}
	io.Pfblue2("a = %v\n", a)

	// eigenvalues
	λ := make([]float64, 3)
	haspert, err := M_FixZeroOrRepeated(λ, a, EV_PERT, EV_EVTOL, EV_ZERO)
	if err != nil {
		chk.Panic("%v", err.Error())
	}
	if haspert {
		io.Pfyel("a(pert) = %v\n", a)
	}
	io.Pfblue2("λ = %v\n", λ)

	// run test (ana)
	if true {
		io.Pforan("\neigenprojectors (ana)\n")
		CheckEigenprojs(a, false, tolP, tolS, ver, EV_EVTOL, EV_ZERO, true)

		io.Pforan("\nderivatives (anaP)\n")
		dPda, dPda_num := CheckEigenprojsDerivs(false, a, toldP, verdP, EV_EVTOL, EV_ZERO)

		// compare derivatives
		chk.Deep3(tst, "dPda", toldP2, dPda, dPda_num)
	}

	// run test
	if false {
		io.Pforan("\neigenprojectors (num)\n")
		CheckEigenprojs(a, true, tolP, tolS, ver, EV_EVTOL, EV_ZERO, true)

		io.Pforan("\nderivatives (numP)\n")
		dPda, dPda_num := CheckEigenprojsDerivs(true, a, toldP, verdP, EV_EVTOL, EV_ZERO)

		// compare derivatives
		chk.Deep3(tst, "dPda", toldP2, dPda, dPda_num)
	}
}

func Test_eigenp06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("eigenp06")

	ncp := 6
	P := M_AllocEigenprojs(ncp)
	chk.Vector(tst, "P0", 1e-17, P[0], []float64{0, 0, 0, 0, 0, 0})
	chk.Vector(tst, "P1", 1e-17, P[1], []float64{0, 0, 0, 0, 0, 0})
	chk.Vector(tst, "P2", 1e-17, P[2], []float64{0, 0, 0, 0, 0, 0})
}
