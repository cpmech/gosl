// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/utl"
)

// Callbacks
type Cb_isofun_f func(p, q float64, args ...interface{}) float64
type Cb_isofun_g func(p, q float64, args ...interface{}) (dfdp, dfdq float64)
type Cb_isofun_h func(p, q float64, args ...interface{}) (d2fdp2, d2fdq2, d2fdpdq float64)

// IsoFun handles isotropic functions and their derivatives
type IsoFun struct {
	// smp invs parameters
	a, b, β, ϵ float64 // smp parameters
	shift      float64 // shift eigenvalues; e.g. to account for cohesion
	ncp        int     // number of stress components
	// callback functions
	ffcn Cb_isofun_f // callback function
	gfcn Cb_isofun_g // callback function
	hfcn Cb_isofun_h // callback function
	// eigenvalues/projectors
	FixRep bool        // do run perturbation code to fix repeated eigenvalues
	HasRep bool        // has repeated eigenvalues
	Pert   float64     // perturbation values
	EvTol  float64     // tolerance to detect repeated eigenvalues
	Zero   float64     // minimum λ to be considered zero
	L, Ls  []float64   // eigenvalues and shifted eigenvalues
	P      [][]float64 // eigenprojectors
	// smp variables
	m          float64   // norm of SMP director
	N, n       []float64 // SMP director and unit director
	Frmp, Grmp []float64 // ramp function values
	p, q       float64   // invariants
	// derivatives
	dNdλ       []float64   // derivatives
	dndλ       [][]float64 // derivatives
	dpdλ, dqdλ []float64   // derivatives
	dfdp, dfdq float64     // derivatives of f w.r.t invariants
	Dfdλ       []float64   // first derivative of f
	// for second derivatives
	Acpy    []float64     // copy of tensor 'A'
	d2ndλdλ [][][]float64 // derivatives
	d2pdλdλ [][]float64   // derivatives
	d2qdλdλ [][]float64   // derivatives
	Dgdλ    [][]float64   // second derivative of f
	d2fdp2  float64       // derivatives of f
	d2fdq2  float64       // derivatives of f
	d2fdpdq float64       // derivatives of f
	dPdA    [][][]float64 // derivatives of eigenprojectors
}

// Accessors
func (o *IsoFun) Get_bsmp() float64      { return o.b }
func (o *IsoFun) Set_bsmp(b float64)     { o.b = b }
func (o *IsoFun) Get_pq() (p, q float64) { return o.p, o.q }

// Init initialises the isotropic function structure
func (o *IsoFun) Init(a, b, β, ϵ, shift float64, ncp int, ffcn Cb_isofun_f, gfcn Cb_isofun_g, hfcn Cb_isofun_h) {

	// smp invs parameters and number of stress components
	o.a, o.b, o.β, o.ϵ = a, b, β, ϵ
	o.shift = shift
	o.ncp = ncp

	// callback functions
	o.ffcn, o.gfcn, o.hfcn = ffcn, gfcn, hfcn

	// do fix repeated eigenvalues?
	o.FixRep = true

	// eigenvalues/projectors
	o.Pert = EV_PERT
	o.EvTol = EV_EVTOL
	o.Zero = EV_ZERO
	o.L = make([]float64, 3)
	o.Ls = make([]float64, 3)
	o.P = la.MatAlloc(3, ncp)

	// smp variables and derivatives
	o.N, o.n = make([]float64, 3), make([]float64, 3)
	o.Frmp, o.Grmp = make([]float64, 3), make([]float64, 3)
	o.dNdλ, o.dndλ = make([]float64, 3), la.MatAlloc(3, 3)

	// derivatives
	o.dpdλ, o.dqdλ = make([]float64, 3), make([]float64, 3)
	o.Dfdλ = make([]float64, 3)

	// for second derivatives
	o.Acpy = make([]float64, ncp)
	o.d2ndλdλ = utl.Deep3alloc(3, 3, 3)
	o.d2pdλdλ = la.MatAlloc(3, 3)
	o.d2qdλdλ = la.MatAlloc(3, 3)
	o.Dgdλ = la.MatAlloc(3, 3)
	o.dPdA = utl.Deep3alloc(3, ncp, ncp)
}

// functions w.r.t principal values //////////////////////////////////////////////////////////////

// Fp evaluates the isotropic function @ λ (principal values)
func (o *IsoFun) Fp(λ []float64, args ...interface{}) (res float64, err error) {

	// apply shift
	o.apply_shift(o.Ls, λ)

	// SMP director and unit director
	o.m = NewSmpDirector(o.N, o.Ls, o.a, o.b, o.β, o.ϵ)
	NewSmpUnitDirector(o.n, o.m, o.N)

	// SMP invariants
	o.p, o.q, err = GenInvs(o.Ls, o.n, o.a)
	if err != nil {
		return
	}

	// function evaluation
	res = o.ffcn(o.p, o.q, args...)
	return
}

// Gp computes f and Dfdλ @ λ (principal values)
// Notes:
//  1) λ is input => shifted and copied into to internal Ls for use in Hp
//  2) output is stored in Dfdλ
func (o *IsoFun) Gp(λ []float64, args ...interface{}) (fval float64, err error) {

	// apply shift
	o.apply_shift(o.Ls, λ)

	// SMP director, unit director and derivatives
	o.m = NewSmpDerivs1(o.dndλ, o.dNdλ, o.N, o.Frmp, o.Grmp, o.Ls, o.a, o.b, o.β, o.ϵ)
	NewSmpUnitDirector(o.n, o.m, o.N)

	// SMP invariants and derivatives
	o.p, o.q, err = GenInvsDeriv1(o.dpdλ, o.dqdλ, o.Ls, o.n, o.dndλ, o.a)
	if err != nil {
		return
	}

	// function evaluation
	fval = o.ffcn(o.p, o.q, args...)
	o.dfdp, o.dfdq = o.gfcn(o.p, o.q, args...)

	// derivatives of f w.r.t eigenvalues
	o.Dfdλ[0] = o.dfdp*o.dpdλ[0] + o.dfdq*o.dqdλ[0]
	o.Dfdλ[1] = o.dfdp*o.dpdλ[1] + o.dfdq*o.dqdλ[1]
	o.Dfdλ[2] = o.dfdp*o.dpdλ[2] + o.dfdq*o.dqdλ[2]

	//o.DebugOutput(true)
	return
}

// HafterGp computes d2fdλdλ == dgdλ after Gp was called
// Notes:
//  1) output is stored in Dgdλ
func (o *IsoFun) HafterGp(args ...interface{}) (err error) {

	// SMP director second derivatives
	NewSmpDerivs2(o.d2ndλdλ, o.Ls, o.a, o.b, o.β, o.ϵ, o.m, o.N, o.Frmp, o.Grmp, o.dNdλ, o.dndλ)

	// SMP invariants and derivatives
	GenInvsDeriv2(o.d2pdλdλ, o.d2qdλdλ, o.Ls, o.n, o.dpdλ, o.dqdλ, o.p, o.q, o.dndλ, o.d2ndλdλ, o.a)

	// function evaluation
	o.d2fdp2, o.d2fdq2, o.d2fdpdq = o.hfcn(o.p, o.q, args...)

	// derivatives of g w.r.t eigenvalues
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			o.Dgdλ[i][j] = o.dpdλ[i]*(o.d2fdp2*o.dpdλ[j]+o.d2fdpdq*o.dqdλ[j]) + o.dfdp*o.d2pdλdλ[i][j] +
				o.dqdλ[i]*(o.d2fdpdq*o.dpdλ[j]+o.d2fdq2*o.dqdλ[j]) + o.dfdq*o.d2qdλdλ[i][j]
		}
	}

	//o.DebugOutput(true)
	return
}

// functions w.r.t full tensor (Mandel) //////////////////////////////////////////////////////////

// ApplyPert computes eigenvalues of A and applies a perturbation to avoid repeated eigenvalues
//  Notes:
//   1) 'A' is modified
//   2) eigenvalues are stored in L
//   3) HasRep indicates that there were repeated eigenvalues
func (o *IsoFun) ApplyPert(A []float64) (err error) {
	o.HasRep, err = M_FixZeroOrRepeated(o.L, A, o.Pert, o.EvTol, o.Zero)
	return
}

// Fa evaluates the isotropic function @ A (a second order tensor in Mandel's basis)
func (o *IsoFun) Fa(A []float64, args ...interface{}) (res float64, err error) {

	// eigenvalues
	copy(o.Acpy, A)
	if o.FixRep {
		err = o.ApplyPert(o.Acpy)
	} else {
		err = M_EigenValsNum(o.L, o.Acpy)
	}
	if err != nil {
		return
	}

	// F w.r.t eigenvalues
	res, err = o.Fp(o.L, args...)
	return
}

// Ga computes f and dfdA @ A (a second order tensor in Mandel's basis)
//  Notes:
//   1) eigenvalues are stored in L and eigenprojectors are stored in P
//   2) Dfdλ is calculated and is availabe for external use
func (o *IsoFun) Ga(dfdA, A []float64, args ...interface{}) (fval float64, err error) {

	// fix repeated
	copy(o.Acpy, A)
	if o.FixRep {
		err = o.ApplyPert(o.Acpy)
		if err != nil {
			return
		}
	}

	// eigenvalues and eigenprojectors
	err = M_EigenValsProjsNum(o.P, o.L, o.Acpy)
	if err != nil {
		return
	}

	// G w.r.t eigenvalues => Dfdλ
	fval, err = o.Gp(o.L, args...)
	if err != nil {
		return
	}

	// dfdA
	for i := 0; i < o.ncp; i++ {
		dfdA[i] = 0
		for k := 0; k < 3; k++ {
			dfdA[i] += o.Dfdλ[k] * o.P[k][i]
		}
	}
	return
}

// HafterGa computes d2fdada after G is called
//  Notes:
//   1) Dgdλ==d2fdλdλ is calculated and is availabe for external use
func (o *IsoFun) HafterGa(d2fdAdA [][]float64, args ...interface{}) (err error) {

	// derivatives of eigenprojectors
	err = M_EigenProjsDeriv(o.dPdA, o.Acpy, o.L, o.P, o.Zero)
	if err != nil {
		return
	}

	// H w.r.t eigenvalues => Dgdλ
	err = o.HafterGp(args...)
	if err != nil {
		return
	}

	// d2fdAdA
	for i := 0; i < o.ncp; i++ {
		for j := 0; j < o.ncp; j++ {
			d2fdAdA[i][j] = o.Dgdλ[0][0]*o.P[0][i]*o.P[0][j] + o.Dgdλ[0][1]*o.P[0][i]*o.P[1][j] + o.Dgdλ[0][2]*o.P[0][i]*o.P[2][j] +
				o.Dgdλ[1][0]*o.P[1][i]*o.P[0][j] + o.Dgdλ[1][1]*o.P[1][i]*o.P[1][j] + o.Dgdλ[1][2]*o.P[1][i]*o.P[2][j] +
				o.Dgdλ[2][0]*o.P[2][i]*o.P[0][j] + o.Dgdλ[2][1]*o.P[2][i]*o.P[1][j] + o.Dgdλ[2][2]*o.P[2][i]*o.P[2][j] +
				+o.Dfdλ[0]*o.dPdA[0][i][j] + o.Dfdλ[1]*o.dPdA[1][i][j] + o.Dfdλ[2]*o.dPdA[2][i][j]
		}
	}
	return
}

// Get_derivs_afterHa returns the derivatives of SMP invariants after a call to HafterGa
func (o *IsoFun) Get_derivs_afterHa(dpdσ, dqdσ []float64) {
	for i := 0; i < o.ncp; i++ {
		dpdσ[i] = o.dpdλ[0]*o.P[0][i] + o.dpdλ[1]*o.P[1][i] + o.dpdλ[2]*o.P[2][i]
		dqdσ[i] = o.dqdλ[0]*o.P[0][i] + o.dqdλ[1]*o.P[1][i] + o.dqdλ[2]*o.P[2][i]
	}
}

// apply_shift applies shift of eigenvalues
func (o *IsoFun) apply_shift(λshifted, λoriginal []float64) {
	λshifted[0] = λoriginal[0] + o.a*o.shift
	λshifted[1] = λoriginal[1] + o.a*o.shift
	λshifted[2] = λoriginal[2] + o.a*o.shift
}

// auxiliary /////////////////////////////////////////////////////////////////////////////////////

// FindIntersect find point on surface using Newton's method
func (o *IsoFun) FindIntersect(p0, k float64, Δλ []float64, usek, debug bool, args ...interface{}) (λ_at_int []float64) {

	// for nonlinear solver
	prms := map[string]float64{
		"atol":    1e-15,
		"rtol":    1e-15,
		"ftol":    1e-15,
		"lSearch": 0.0,
		"maxIt":   20,
	}
	var intersect_ffcn num.Cb_f
	var intersect_Jfcn num.Cb_J
	var xsol []float64
	var neq int

	// auxiliary variables
	A := make([]float64, o.ncp)
	dfdA := make([]float64, o.ncp)
	var err error

	// k-path
	if usek {

		// nonlinear problem
		a0, a1 := k-3.0, 2.0*k+3.0
		intersect_ffcn = func(fx, x []float64) error { // x[0]==σ1, x[1]==σ3
			A[0], A[1], A[2] = x[0], x[1], x[2]
			fx[0], err = o.Fa(A, args...)
			if err != nil {
				chk.Panic(_isofun_err1)
			}
			fx[1] = a0*x[0] + a1*x[1] + 3.0*k*p0
			fx[2] = a0*x[0] + a1*x[2] + 3.0*k*p0
			return nil
		}

		// Jacobian of nonlinear problem
		intersect_Jfcn = func(dfdx *la.Triplet, x []float64) error {
			A[0], A[1], A[2] = x[0], x[1], x[2]
			_, err = o.Ga(dfdA, A, args...)
			if err != nil {
				chk.Panic(_isofun_err2)
			}
			dfdx.Start()
			dfdx.Put(0, 0, o.dfdp*o.dpdλ[0]+o.dfdq*o.dqdλ[0])
			dfdx.Put(0, 1, o.dfdp*o.dpdλ[1]+o.dfdq*o.dqdλ[1])
			dfdx.Put(0, 2, o.dfdp*o.dpdλ[2]+o.dfdq*o.dqdλ[2])
			dfdx.Put(1, 0, a0)
			dfdx.Put(1, 1, a1)
			dfdx.Put(1, 2, 0.0)
			dfdx.Put(2, 0, a0)
			dfdx.Put(2, 1, 0.0)
			dfdx.Put(2, 2, a1)
			return nil
		}

		// trial solution
		neq = 3 // number of equations
		xsol = []float64{-p0, -p0 * 1.1, -p0 * 1.1}

		// delta-path
	} else {

		// check
		if len(Δλ) != 3 {
			chk.Panic(_isofun_err3)
		}

		// nonlinear problem
		intersect_ffcn = func(fx, x []float64) error { // x[0]==m
			m := x[0]
			A[0] = -p0 + m*Δλ[0]
			A[1] = -p0 + m*Δλ[1]
			A[2] = -p0 + m*Δλ[2]
			fx[0], err = o.Fa(A, args...)
			if err != nil {
				chk.Panic(_isofun_err1)
			}
			return nil
		}

		// Jacobian of nonlinear problem
		intersect_Jfcn = func(dfdx *la.Triplet, x []float64) error {
			m := x[0]
			A[0] = -p0 + m*Δλ[0]
			A[1] = -p0 + m*Δλ[1]
			A[2] = -p0 + m*Δλ[2]
			_, err = o.Ga(dfdA, A, args...)
			if err != nil {
				chk.Panic(_isofun_err2)
			}
			dfdx.Start()
			dfdx.Put(0, 0, o.dfdp*(o.dpdλ[0]*Δλ[0]+o.dpdλ[1]*Δλ[1]+o.dpdλ[2]*Δλ[2])+
				o.dfdq*(o.dqdλ[0]*Δλ[0]+o.dqdλ[1]*Δλ[1]+o.dqdλ[2]*Δλ[2]))
			return nil
		}

		// trial solution
		neq = 1 // number of equations
		xsol = []float64{0.5}
	}

	// solve
	jacfcn := intersect_Jfcn
	numjac := false
	if debug {
		f_xsol := make([]float64, neq)
		intersect_ffcn(f_xsol, xsol)
		io.Pforan("f_xsol (before) = %+v\n", f_xsol)
	}
	var nls num.NlSolver
	nls.Init(neq, intersect_ffcn, jacfcn, nil, false, numjac, prms)
	defer nls.Clean()
	err = nls.Solve(xsol, !debug)
	if err != nil {
		chk.Panic(_isofun_err10, err.Error())
	}
	if debug {
		f_xsol := make([]float64, neq)
		intersect_ffcn(f_xsol, xsol)
		io.Pforan("f_xsol (after)  = %+v\n", f_xsol)
	}

	// check Jacobian
	if debug {
		_, err = nls.CheckJ(xsol, 1.5e-6, false, true)
		if err != nil {
			chk.Panic("%v", err)
		}
	}

	// set solution
	λ_at_int = make([]float64, 3)
	if !usek {
		m := xsol[0]
		λ_at_int[0] = -p0 + m*Δλ[0]
		λ_at_int[1] = -p0 + m*Δλ[1]
		λ_at_int[2] = -p0 + m*Δλ[2]
	} else {
		λ_at_int[0] = xsol[0]
		λ_at_int[1] = xsol[1]
		λ_at_int[2] = xsol[2]
	}
	return
}

// CheckGrads check df/dA and d²f/dA²
func (o *IsoFun) CheckGrads(A []float64, tol, tol2 float64, ver bool) {

	// compute derivatives
	dfdA := make([]float64, len(A))
	_, err := o.Ga(dfdA, A)
	if err != nil {
		chk.Panic(_isofun_err4, err)
	}
	d2fdAdA := la.MatAlloc(len(A), len(A))
	err = o.HafterGa(d2fdAdA)
	if err != nil {
		chk.Panic(_isofun_err6, err)
	}

	// df/dA
	if ver {
		io.Pfpink("\ndfdA . . . . . . . . \n")
	}
	var fval, tmp float64
	has_error := false
	for j := 0; j < len(A); j++ {
		dnum := num.DerivFwd(func(x float64, args ...interface{}) (res float64) {
			tmp, A[j] = A[j], x
			fval, err = o.Fa(A)
			if err != nil {
				chk.Panic(_isofun_err5, err)
			}
			A[j] = tmp
			return fval
		}, A[j], 1e-6)
		err := chk.PrintAnaNum(io.Sf("df/dA[%d]", j), tol, dfdA[j], dnum, ver)
		if err != nil {
			has_error = true
		}
	}
	if has_error {
		chk.Panic(_isofun_err8)
	}

	// d2f/dAdA
	if ver {
		io.Pfpink("\nd²f/dAdA . . . . . . . . \n")
	}
	dfdA_tmp := make([]float64, len(A))
	for i := 0; i < len(A); i++ {
		for j := 0; j < len(A); j++ {
			dnum, _ := num.DerivCentral(func(x float64, args ...interface{}) (res float64) {
				tmp, A[j] = A[j], x
				_, err = o.Ga(dfdA_tmp, A)
				if err != nil {
					chk.Panic(_isofun_err7, err)
				}
				A[j] = tmp
				return dfdA_tmp[i]
			}, A[j], 1e-6)
			err := chk.PrintAnaNum(io.Sf("d²f/dA[%d]dA[%d]", i, j), tol2, d2fdAdA[i][j], dnum, ver)
			if err != nil {
				has_error = true
			}
		}
	}
	if has_error {
		chk.Panic(_isofun_err9)
	}
}

// String returns information of this object
func (o *IsoFun) String() (l string) {
	l = io.Sf("  a=%v b=%v β=%v ϵ=%v\n", o.a, o.b, o.β, o.ϵ)
	//l += io.Sf(" λ    = %v\n", o.L)
	//l += io.Sf(" p, q = %v, %v\n", o.p, o.q)
	return
}

// DebugOutput outputs debug information
func (o *IsoFun) DebugOutput(princOnly bool) {
	io.Pfblue("a, b, β, ϵ = %v, %v, %v, %v\n", o.a, o.b, o.β, o.ϵ)
	io.Pfcyan2("L       = %v\n", o.L)
	io.Pfcyan2("Ls      = %v\n", o.Ls)
	io.Pfblue2("m       = %v\n", o.m)
	io.Pfblue2("N       = %v\n", o.N)
	io.Pfblue2("n       = %v\n", o.n)
	io.Pfblue("Frmp    = %v\n", o.Frmp)
	io.Pfblue("Grmp    = %v\n", o.Grmp)
	io.Pfblue("dNdλ    = %v\n", o.dNdλ)
	io.Pfblue("dndλ    = %v\n", o.dndλ)
	io.Pfblue("d2ndλdλ = %v\n", o.d2ndλdλ)
	io.Pfblue("d2pdλdλ = %v\n", o.d2pdλdλ)
	io.Pfblue("d2qdλdλ = %v\n", o.d2qdλdλ)
	io.Pfblue2("p       = %v\n", o.p)
	io.Pfblue2("q       = %v\n", o.q)
	io.Pfblue2("dpdλ    = %v\n", o.dpdλ)
	io.Pfblue2("dqdλ    = %v\n", o.dqdλ)
	io.Pfblue2("dfdp    = %v\n", o.dfdp)
	io.Pfblue2("dfdq    = %v\n", o.dfdq)
	io.Pfblue2("Dfdλ    = %v\n", o.Dfdλ)
	io.Pfblue("d2fdp2  = %v\n", o.d2fdp2)
	io.Pfblue("d2fdq2  = %v\n", o.d2fdq2)
	io.Pfblue("d2fdpdq = %v\n", o.d2fdpdq)
	io.Pfblue("Dgdλ    = %v\n", o.Dgdλ)
	if !princOnly {
		io.Pf("Acpy = %+#v\n", o.Acpy)
		io.Pf("P[0] = %v\n", o.P[0])
		io.Pf("P[1] = %v\n", o.P[1])
		io.Pf("P[2] = %v\n", o.P[2])
		la.PrintMat("dPdA[0]", o.dPdA[0], "%15e", false)
		la.PrintMat("dPdA[1]", o.dPdA[1], "%15e", false)
		la.PrintMat("dPdA[2]", o.dPdA[2], "%15e", false)
	}
}

// error messages
var (
	_isofun_err1  = "isofun.go: FindIntersect: SMPinvs failed in NlSolve for fx"
	_isofun_err2  = "isofun.go: FindIntersect: SMPderivs1 failed in NlSolve for dfdx"
	_isofun_err3  = "isofun.go: FindIntersect: if usek=false, Δλ must be given"
	_isofun_err4  = "isofun.go: CheckGrads: cannot compute G(A):\n%v"
	_isofun_err5  = "isofun.go: CheckGrads: DerivCentral: cannot compute F(A):\n%v"
	_isofun_err6  = "isofun.go: CheckGrads: cannot compute H(A):\n%v"
	_isofun_err7  = "isofun.go: CheckGrads: DerivCentral: cannot compute G(A):\n%v"
	_isofun_err8  = "isofun.go: CheckGrads: df/dA failed"
	_isofun_err9  = "isofun.go: CheckGrads: d²f/dA² failed"
	_isofun_err10 = "isofun.go: FindIntersect: nonlinear soler failed:\n %v"
)
