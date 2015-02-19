// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
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

const (
	//SAVEISOPLOT = true
	SAVEISOPLOT = false
)

func Test_isofun01(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose() = false
	chk.PrintTitle("isofun01")

	a, b, β, ϵ := 1.0, 0.5, 2.0, 1e-3
	shift := 0.0

	φ := 30.0
	μ := NewSmpCalcμ(φ, a, b, β, ϵ)
	io.Pforan("μ = %v\n", μ)

	simpleform := true
	notfcrit := false

	dver := true
	dtol := 1e-5
	dtol2 := 1e-6

	ffcn := func(p, q float64, args ...interface{}) float64 {
		if notfcrit {
			return p*p + q*q
		}
		if simpleform {
			return q - μ*p
		}
		return q/p - μ
	}

	gfcn := func(p, q float64, args ...interface{}) (dfdp, dfdq float64) {
		if notfcrit {
			dfdp = 2.0 * p
			dfdq = 2.0 * q
			return
		}
		if simpleform {
			dfdp = -μ
			dfdq = 1.0
			return
		}
		dfdp = -q / (p * p)
		dfdq = 1.0 / p
		return
	}

	hfcn := func(p, q float64, args ...interface{}) (d2fdp2, d2fdq2, d2fdpdq float64) {
		if notfcrit {
			d2fdp2 = 2.0
			d2fdq2 = 2.0
			d2fdpdq = 0.0
			return
		}
		if simpleform {
			return
		}
		d2fdp2 = 2.0 * q / (p * p * p)
		d2fdq2 = 0.0
		d2fdpdq = -1.0 / (p * p)
		return
	}

	nd := test_nd
	for idxA := 0; idxA < len(test_nd); idxA++ {
		//for idxA := 10; idxA < 11; idxA++ {
		//for idxA := 10; idxA < len(test_nd); idxA++ {

		// change tolerances
		dtol_, dtol2_ := dtol, dtol2
		switch idxA {
		case 12:
			dtol, dtol2 = 1.72e-5, 1e-4
		}

		// tensor
		AA := test_AA[idxA]
		A := M_Alloc2(nd[idxA])
		Ten2Man(A, AA)
		io.PfYel("\n\ntst # %d ###################################################################################\n", idxA)
		io.Pfblue2("A = %v\n", A)

		// isotropic function
		var o IsoFun
		o.Init(a, b, β, ϵ, shift, len(A), ffcn, gfcn, hfcn)

		// function evaluation and shifted eigenvalues
		fval, err := o.Fa(A)
		if err != nil {
			chk.Panic("cannot compute F(A):\n%v", err)
		}
		if o.HasRep {
			copy(A, o.Acpy)
			io.Pfyel("A(pert) = %v\n", A)
			io.Pfyel("λ(pert) = %v\n", o.L)
		}
		io.Pforan("p, q = %v, %v\n", o.p, o.q)
		io.Pforan("f(A) = %v\n", fval)

		// check gradients
		o.CheckGrads(A, dtol, dtol2, dver)

		// recover tolerances
		dtol, dtol2 = dtol_, dtol2_
	}
}

func Test_isofun02(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose() = false
	chk.PrintTitle("isofun02")

	// SMP director parameters
	a, b, β, ϵ := -1.0, 0.5, 2.0, 1e-3
	shift := 0.0

	// failure crit parameters and number of stress components
	φ, ncp := 30.0, 4

	// q/p coefficient
	μ := NewSmpCalcμ(φ, a, b, β, ϵ)
	io.Pforan("μ = %v\n", μ)

	// yield func coefficients
	r := 2.0
	qy := 0.0
	α := math.Atan(μ)
	sα := math.Sin(α)
	cα := math.Cos(α)

	// yield function
	ffcn := func(p, q float64, args ...interface{}) float64 {
		//io.Pfgrey("p, q = %v, %v\n", p, q)
		pc := -qy/μ + r/sα
		R := math.Sqrt(q*q + (p-pc)*(p-pc))
		pd := pc - R*sα
		if p < pd {
			return R - r
		}
		return q*cα - (p-pc)*sα - r
	}

	// first order derivative
	gfcn := func(p, q float64, args ...interface{}) (dfdp, dfdq float64) {
		return
	}

	// second order derivative
	hfcn := func(p, q float64, args ...interface{}) (d2fdp2, d2fdq2, d2fdpdq float64) {
		return
	}

	// isotropic functions
	var o IsoFun
	o.Init(a, b, β, ϵ, shift, ncp, ffcn, gfcn, hfcn)

	// plot
	σcCte := 10.0
	M := Phi2M(φ, "oct")
	rmin, rmax := 0.0, 1.28*M*σcCte
	nr, nα := 31, 81
	//nr,   nα   := 31, 1001
	npolarc := true
	simplec := false
	only0 := false
	grads := false
	showpts := false
	ferr := 10.0
	if SAVEISOPLOT {
		PlotOct("fig_isofun02.png", σcCte, rmin, rmax, nr, nα, φ, o.Fa, o.Ga,
			npolarc, simplec, only0, grads, showpts, true, true, ferr)
	}
}

func Test_isofun03(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose() = false
	chk.PrintTitle("isofun03")

	// SMP director parameters
	a, b, β, ϵ := -1.0, 0.0, 1.0, 1e-3
	shift := 0.0

	// failure crit parameters
	φ := 20.0
	M := NewSmpCalcμ(φ, a, b, β, ϵ)

	// functions
	ffcn := func(p, q float64, args ...interface{}) float64 {
		return M*M*p*p + q*q - 1.0
	}
	gfcn := func(p, q float64, args ...interface{}) (dfdp, dfdq float64) {
		dfdp = M * M * 2.0 * p
		dfdq = 2.0 * q
		return
	}
	hfcn := func(p, q float64, args ...interface{}) (d2fdp2, d2fdq2, d2fdpdq float64) {
		d2fdp2 = M * M * 2.0
		d2fdq2 = 2.0
		d2fdpdq = 0
		return
	}

	// isotropic functions
	ncp := 6
	var o IsoFun
	o.Init(a, b, β, ϵ, shift, ncp, ffcn, gfcn, hfcn)

	// find point on surface using Newton's method
	p0, k := 1.0, 0.8
	//usek  := true
	usek := false
	debug := false
	Δλ := []float64{-2, -3, -4}
	//Δλ    := []float64{-0.2, -0.3, -4}
	//Δλ    := []float64{-1, 0, 0}
	λ := o.FindIntersect(p0, k, Δλ, usek, debug)

	// change λ
	//λ = []float64{-1,-1,-1}

	// check gradients @ intersection
	dtol, dtol2, ver := 1e-7, 1e-9, true
	A := make([]float64, ncp)
	A[0], A[1], A[2] = λ[0], λ[1], λ[2]
	f_at_A, _ := o.Fa(A)
	io.Pforan("λ==A = %+v\n", λ)
	io.Pforan("f(A) = %+v\n", f_at_A)
	o.CheckGrads(A, dtol, dtol2, ver)

	// plot
	σcCte := 10.0
	Moct := Phi2M(φ, "oct")
	rmin, rmax := 0.0, 1.28*Moct*σcCte
	nr, nα := 31, 81
	npolarc := true
	simplec := false
	only0 := false
	grads := false
	showpts := false
	ferr := 10.0
	if SAVEISOPLOT {
		PlotOct("fig_isofun03.png", σcCte, rmin, rmax, nr, nα, φ, o.Fa, o.Ga,
			npolarc, simplec, only0, grads, showpts, true, true, ferr)
	}
}

func Test_isofun04(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose() = false
	chk.PrintTitle("isofun04")

	// constants
	dtol, dtol2, ver := 1e-6, 1e-8, true

	// SMP director parameters
	a, b, β, ϵ := -1.0, 1.0, 1.0, 1e-3
	shift := 2.0

	// failure crit parameters
	φ := 20.0
	M := NewSmpCalcμ(φ, a, b, β, ϵ)

	// for rounding
	r := 0.3
	α := math.Atan(M)
	sα := math.Sin(α)
	cα := math.Cos(α)

	// functions
	auxvars := func(p, q float64) (pc, R, pd float64) {
		pc = r / sα
		R = math.Sqrt(q*q + (p-pc)*(p-pc))
		pd = pc - R*sα
		return
	}
	ffcn := func(p, q float64, args ...interface{}) float64 {
		pc, R, pd := auxvars(p, q)
		if p < pd {
			return R - r
		}
		return q*cα - (p-pc)*sα - r
	}
	gfcn := func(p, q float64, args ...interface{}) (dfdp, dfdq float64) {
		pc, R, pd := auxvars(p, q)
		if p < pd {
			dfdp = (p - pc) / R
			dfdq = q / R
			return
		}
		dfdp = -sα
		dfdq = cα
		return
	}
	hfcn := func(p, q float64, args ...interface{}) (d2fdp2, d2fdq2, d2fdpdq float64) {
		pc, R, pd := auxvars(p, q)
		if p < pd {
			R3 := R * R * R
			d2fdp2 = 1.0/R - (p-pc)*(p-pc)/R3
			d2fdq2 = 1.0/R - q*q/R3
			d2fdpdq = -(p - pc) * q / R3
			return
		}
		return
	}

	// isotropic functions
	ncp := 6
	var o IsoFun
	o.Init(a, b, β, ϵ, shift, ncp, ffcn, gfcn, hfcn)

	// find point on surface using Newton's method
	p0, k := 3.0, 0.8
	usek := false
	debug := true
	idx := 3
	var Δλ []float64
	switch idx {
	case 1:
		usek = true
	case 2:
		Δλ = []float64{-2, -3, -4}
	case 3:
		Δλ = []float64{4, 4, 4}
		dtol2 = 1e-7
	}
	λ := o.FindIntersect(p0, k, Δλ, usek, debug)
	A := make([]float64, ncp)
	A[0], A[1], A[2] = λ[0], λ[1], λ[2]

	// handle repeated eigenvalues
	f_at_A, _ := o.Fa(A)
	if o.HasRep {
		copy(A, o.Acpy)
		copy(λ, o.L)
	}
	f_at_λ, _ := o.Fp(λ)
	io.Pforan("λ = %v\n", λ)
	io.Pforan("A = %v\n", A)
	io.Pforan("f(λ) = %v\n", f_at_λ)
	io.Pforan("f(A) = %v\n", f_at_A)

	// check gradients @ intersection
	o.CheckGrads(A, dtol, dtol2, ver)

	// plot
	σcCte := 10.0
	Moct := Phi2M(φ, "oct")
	rmin, rmax := 0.0, 1.3*Moct*σcCte
	nr, nα := 61, 81
	npolarc := true
	simplec := false
	only0 := false
	grads := false
	showpts := false
	ferr := 10.0
	if SAVEISOPLOT {
		PlotOct("fig_isofun04.png", σcCte, rmin, rmax, nr, nα, φ, o.Fa, o.Ga,
			npolarc, simplec, only0, grads, showpts, true, true, ferr)
	}
}
