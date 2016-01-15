// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
	"github.com/cpmech/gosl/vtk"
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
	L, Ls []float64   // eigenvalues and shifted eigenvalues
	P     [][]float64 // eigenprojectors

	// smp variables
	m          float64   // norm of SMP director
	N, n       []float64 // SMP director and unit director
	Frmp, Grmp []float64 // ramp function values
	p, q       float64   // invariants

	// derivatives
	dNdL       []float64   // derivatives
	dndL       [][]float64 // derivatives
	dpdL, dqdL []float64   // derivatives
	dfdp, dfdq float64     // derivatives of f w.r.t invariants
	DfdL       []float64   // first derivative of f

	// for second derivatives
	Acpy    []float64     // copy of tensor 'A'
	d2ndLdL [][][]float64 // derivatives
	d2pdLdL [][]float64   // derivatives
	d2qdLdL [][]float64   // derivatives
	DgdL    [][]float64   // second derivative of f
	d2fdp2  float64       // derivatives of f
	d2fdq2  float64       // derivatives of f
	d2fdpdq float64       // derivatives of f
	dPdA    [][][]float64 // derivatives of eigenprojectors
}

// Accessors
func (o *IsoFun) Get_bsmp() float64      { return o.b }
func (o *IsoFun) Set_bsmp(b float64)     { o.b = b }
func (o *IsoFun) Get_pq() (p, q float64) { return o.p, o.q }

// SetPrms set parameters
func (o *IsoFun) SetPrms(a, b, β, ϵ, shift float64, ffcn Cb_isofun_f, gfcn Cb_isofun_g, hfcn Cb_isofun_h) {
	o.a, o.b, o.β, o.ϵ = a, b, β, ϵ
	o.shift = shift
	o.ffcn, o.gfcn, o.hfcn = ffcn, gfcn, hfcn // TODO: check why this is neccessary to refresh plot?
}

// Init initialises the isotropic function structure
func (o *IsoFun) Init(a, b, β, ϵ, shift float64, ncp int, ffcn Cb_isofun_f, gfcn Cb_isofun_g, hfcn Cb_isofun_h) {

	// smp invs parameters and number of stress components
	o.a, o.b, o.β, o.ϵ = a, b, β, ϵ
	o.shift = shift
	o.ncp = ncp

	// callback functions
	o.ffcn, o.gfcn, o.hfcn = ffcn, gfcn, hfcn

	// eigenvalues/projectors
	o.L = make([]float64, 3)
	o.Ls = make([]float64, 3)
	o.P = la.MatAlloc(3, ncp)

	// smp variables and derivatives
	o.N, o.n = make([]float64, 3), make([]float64, 3)
	o.Frmp, o.Grmp = make([]float64, 3), make([]float64, 3)
	o.dNdL, o.dndL = make([]float64, 3), la.MatAlloc(3, 3)

	// derivatives
	o.dpdL, o.dqdL = make([]float64, 3), make([]float64, 3)
	o.DfdL = make([]float64, 3)

	// for second derivatives
	o.Acpy = make([]float64, ncp)
	o.d2ndLdL = utl.Deep3alloc(3, 3, 3)
	o.d2pdLdL = la.MatAlloc(3, 3)
	o.d2qdLdL = la.MatAlloc(3, 3)
	o.DgdL = la.MatAlloc(3, 3)
	o.dPdA = utl.Deep3alloc(3, ncp, ncp)
}

// functions w.r.t principal values //////////////////////////////////////////////////////////////

// Fp evaluates the isotropic function @ L (principal values)
func (o *IsoFun) Fp(L []float64, args ...interface{}) (res float64, err error) {

	// apply shift
	o.apply_shift(o.Ls, L)

	// SMP director and unit director
	o.m = SmpDirector(o.N, o.Ls, o.a, o.b, o.β, o.ϵ)
	SmpUnitDirector(o.n, o.m, o.N)

	// SMP invariants
	o.p, o.q, err = GenInvs(o.Ls, o.n, o.a)
	if err != nil {
		return
	}

	// function evaluation
	res = o.ffcn(o.p, o.q, args...)
	return
}

// Gp computes f and DfdL @ L (principal values)
// Notes:
//  1) L is input => shifted and copied into to internal Ls for use in Hp
//  2) output is stored in DfdL
func (o *IsoFun) Gp(L []float64, args ...interface{}) (fval float64, err error) {

	// apply shift
	o.apply_shift(o.Ls, L)

	// SMP director, unit director and derivatives
	o.m = SmpDerivs1(o.dndL, o.dNdL, o.N, o.Frmp, o.Grmp, o.Ls, o.a, o.b, o.β, o.ϵ)
	SmpUnitDirector(o.n, o.m, o.N)

	// SMP invariants and derivatives
	o.p, o.q, err = GenInvsDeriv1(o.dpdL, o.dqdL, o.Ls, o.n, o.dndL, o.a)
	if err != nil {
		return
	}

	// function evaluation
	fval = o.ffcn(o.p, o.q, args...)
	o.dfdp, o.dfdq = o.gfcn(o.p, o.q, args...)

	// derivatives of f w.r.t eigenvalues
	o.DfdL[0] = o.dfdp*o.dpdL[0] + o.dfdq*o.dqdL[0]
	o.DfdL[1] = o.dfdp*o.dpdL[1] + o.dfdq*o.dqdL[1]
	o.DfdL[2] = o.dfdp*o.dpdL[2] + o.dfdq*o.dqdL[2]

	//o.DebugOutput(true)
	return
}

// HafterGp computes d2fdLdL == dgdL after Gp was called
// Notes:
//  1) output is stored in DgdL
func (o *IsoFun) HafterGp(args ...interface{}) (err error) {

	// SMP director second derivatives
	SmpDerivs2(o.d2ndLdL, o.Ls, o.a, o.b, o.β, o.ϵ, o.m, o.N, o.Frmp, o.Grmp, o.dNdL, o.dndL)

	// SMP invariants and derivatives
	GenInvsDeriv2(o.d2pdLdL, o.d2qdLdL, o.Ls, o.n, o.dpdL, o.dqdL, o.p, o.q, o.dndL, o.d2ndLdL, o.a)

	// function evaluation
	o.d2fdp2, o.d2fdq2, o.d2fdpdq = o.hfcn(o.p, o.q, args...)

	// derivatives of g w.r.t eigenvalues
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			o.DgdL[i][j] =
				o.dpdL[i]*(o.d2fdp2*o.dpdL[j]+o.d2fdpdq*o.dqdL[j]) + o.dfdp*o.d2pdLdL[i][j] +
					o.dqdL[i]*(o.d2fdpdq*o.dpdL[j]+o.d2fdq2*o.dqdL[j]) + o.dfdq*o.d2qdLdL[i][j]
		}
	}

	//o.DebugOutput(true)
	return
}

// functions w.r.t full tensor (Mandel) //////////////////////////////////////////////////////////

// Fa evaluates the isotropic function @ A (a second order tensor in Mandel's basis)
func (o *IsoFun) Fa(A []float64, args ...interface{}) (res float64, err error) {

	// eigenvalues
	chk.IntAssert(len(A), o.ncp)
	err = M_EigenValsNum(o.L, A)
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
//   2) DfdL is calculated and is availabe for external use
func (o *IsoFun) Ga(dfdA, A []float64, args ...interface{}) (fval float64, err error) {

	// eigenvalues and eigenprojectors
	chk.IntAssert(len(A), o.ncp)
	copy(o.Acpy, A)
	err = M_EigenValsProjsNum(o.P, o.L, o.Acpy)
	if err != nil {
		return
	}

	// G w.r.t eigenvalues => DfdL
	fval, err = o.Gp(o.L, args...)
	if err != nil {
		return
	}

	// dfdA
	for i := 0; i < o.ncp; i++ {
		dfdA[i] = 0
		for k := 0; k < 3; k++ {
			dfdA[i] += o.DfdL[k] * o.P[k][i]
		}
	}
	return
}

// HafterGa computes d2fdada after G is called
//  Input:
//    dfdA -- df/da
//    args -- extra arguments
//  Output:
//    d2fdAdA -- d²f/dAdA
//  Notes:
//    1) DgdL==d2fdLdL is calculated and is availabe for external use
func (o *IsoFun) HafterGa(d2fdAdA [][]float64, args ...interface{}) (err error) {

	// derivatives of eigenprojectors
	err = M_EigenProjsDerivAuto(o.dPdA, o.Acpy, o.L, o.P)
	if err != nil {
		return
	}

	// H w.r.t eigenvalues => DgdL
	err = o.HafterGp(args...)
	if err != nil {
		return
	}

	// d2fdAdA
	for i := 0; i < o.ncp; i++ {
		for j := 0; j < o.ncp; j++ {
			d2fdAdA[i][j] = o.DgdL[0][0]*o.P[0][i]*o.P[0][j] + o.DgdL[0][1]*o.P[0][i]*o.P[1][j] + o.DgdL[0][2]*o.P[0][i]*o.P[2][j] +
				o.DgdL[1][0]*o.P[1][i]*o.P[0][j] + o.DgdL[1][1]*o.P[1][i]*o.P[1][j] + o.DgdL[1][2]*o.P[1][i]*o.P[2][j] +
				o.DgdL[2][0]*o.P[2][i]*o.P[0][j] + o.DgdL[2][1]*o.P[2][i]*o.P[1][j] + o.DgdL[2][2]*o.P[2][i]*o.P[2][j] +
				+o.DfdL[0]*o.dPdA[0][i][j] + o.DfdL[1]*o.dPdA[1][i][j] + o.DfdL[2]*o.dPdA[2][i][j]
		}
	}
	/*
		for i := 0; i < o.ncp; i++ {
			for j := 0; j < o.ncp; j++ {
				d2fdAdA[i][j] = 0
				for k := 0; k < 3; k++ {
					for l := 0; l < 3; l++ {
						d2fdAdA[i][j] += o.DgdL[k][l] * o.P[k][i] * o.P[l][j]
					}
					d2fdAdA[i][j] += o.DfdL[k] * o.dPdA[k][i][j]
				}
			}
		}
	*/
	return
}

// Get_derivs_afterHa returns the derivatives of SMP invariants after a call to HafterGa
func (o *IsoFun) Get_derivs_afterHa(dpdσ, dqdσ []float64) {
	for i := 0; i < o.ncp; i++ {
		dpdσ[i] = o.dpdL[0]*o.P[0][i] + o.dpdL[1]*o.P[1][i] + o.dpdL[2]*o.P[2][i]
		dqdσ[i] = o.dqdL[0]*o.P[0][i] + o.dqdL[1]*o.P[1][i] + o.dqdL[2]*o.P[2][i]
	}
}

// apply_shift applies shift of eigenvalues
func (o *IsoFun) apply_shift(Lshifted, Loriginal []float64) {
	Lshifted[0] = Loriginal[0] + o.a*o.shift
	Lshifted[1] = Loriginal[1] + o.a*o.shift
	Lshifted[2] = Loriginal[2] + o.a*o.shift
}

// auxiliary /////////////////////////////////////////////////////////////////////////////////////

// Plots Ffcn
func (o *IsoFun) PlotFfcn(dirout, fname string, pmin, pmax float64, np int, pq_point []float64, args_contour, args_point string, extra_before, extra_after func(), args ...interface{}) {
	if extra_before != nil {
		extra_before()
	}
	x, y := utl.MeshGrid2D(pmin, pmax, pmin, pmax, np, np)
	z := la.MatAlloc(np, np)
	for i := 0; i < np; i++ {
		for j := 0; j < np; j++ {
			z[i][j] = o.ffcn(x[i][j], y[i][j], args...)
		}
	}
	plt.Contour(x, y, z, args_contour)
	plt.ContourSimple(x, y, z, false, 8, "levels=[0], colors=['yellow'], linewidths=[2]")
	if pq_point != nil {
		plt.PlotOne(pq_point[0], pq_point[1], args_point)
	}
	plt.Gll("p", "q", "")
	if extra_after != nil {
		extra_after()
	}
	plt.SaveD(dirout, fname)
}

// FindIntersect find point on surface using Newton's method
func (o *IsoFun) FindIntersect(p0, k float64, ΔL []float64, usek, debug bool, args ...interface{}) (L_at_int []float64) {

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
				chk.Panic("Fa failed:\n%v", err)
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
				chk.Panic("Ga failed:%v\n", err)
			}
			dfdx.Start()
			dfdx.Put(0, 0, o.dfdp*o.dpdL[0]+o.dfdq*o.dqdL[0])
			dfdx.Put(0, 1, o.dfdp*o.dpdL[1]+o.dfdq*o.dqdL[1])
			dfdx.Put(0, 2, o.dfdp*o.dpdL[2]+o.dfdq*o.dqdL[2])
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
		if len(ΔL) != 3 {
			chk.Panic("if usek=false, ΔL must be given")
		}

		// nonlinear problem
		intersect_ffcn = func(fx, x []float64) error { // x[0]==m
			m := x[0]
			A[0] = -p0 + m*ΔL[0]
			A[1] = -p0 + m*ΔL[1]
			A[2] = -p0 + m*ΔL[2]
			fx[0], err = o.Fa(A, args...)
			if err != nil {
				chk.Panic("Fa failed:%v\n", err)
			}
			return nil
		}

		// Jacobian of nonlinear problem
		intersect_Jfcn = func(dfdx *la.Triplet, x []float64) error {
			m := x[0]
			A[0] = -p0 + m*ΔL[0]
			A[1] = -p0 + m*ΔL[1]
			A[2] = -p0 + m*ΔL[2]
			_, err = o.Ga(dfdA, A, args...)
			if err != nil {
				chk.Panic("Ga failed:%v\n", err)
			}
			dfdx.Start()
			dfdx.Put(0, 0, o.dfdp*(o.dpdL[0]*ΔL[0]+o.dpdL[1]*ΔL[1]+o.dpdL[2]*ΔL[2])+
				o.dfdq*(o.dqdL[0]*ΔL[0]+o.dqdL[1]*ΔL[1]+o.dqdL[2]*ΔL[2]))
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
		chk.Panic("nonlinear solver failed:\n%v", err)
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
	L_at_int = make([]float64, 3)
	if !usek {
		m := xsol[0]
		L_at_int[0] = -p0 + m*ΔL[0]
		L_at_int[1] = -p0 + m*ΔL[1]
		L_at_int[2] = -p0 + m*ΔL[2]
	} else {
		L_at_int[0] = xsol[0]
		L_at_int[1] = xsol[1]
		L_at_int[2] = xsol[2]
	}
	return
}

// View visualises the isosurface with VTK
func (o IsoFun) View(l float64, L []float64, grads bool, gradsFtol float64, extraconf func(is *vtk.IsoSurf), args ...interface{}) {

	// scene
	scn := vtk.NewScene()
	scn.AxesLen = l
	scn.Reverse = o.a < 0

	// stress point
	var sph *vtk.Sphere
	if L != nil {
		sph = vtk.NewSphere()
		sph.Cen = L
		sph.R = 0.2
		sph.Color = []float64{1, 1, 0, 1}
		sph.AddTo(scn)
	}

	// callback function
	fcn := func(x []float64) (f, vx, vy, vz float64) {
		o.L[0], o.L[1], o.L[2] = x[0], x[1], x[2]
		f, err := o.Fp(o.L, args...)
		if err != nil {
			f = 1e+5
			return
		}
		if grads {
			_, err = o.Gp(o.L, args...)
			if err != nil {
				return
			}
			if math.Abs(f) < gradsFtol {
				vx, vy, vz = o.DfdL[0], o.DfdL[1], o.DfdL[2]
			}
		}
		return
	}

	// isosurfaces
	isf1 := vtk.NewIsoSurf(fcn)
	isf2 := vtk.NewIsoSurf(fcn)
	isf1.AddTo(scn)
	isf2.AddTo(scn)

	// set isosurfaces
	npq := 61
	npt := 101
	setisof := func(isf *vtk.IsoSurf) {
		isf.OctRotate = true
		if isf.OctRotate {
			isf.Limits = []float64{-l, l, 0, l, 0, 360}
			isf.Ndiv = []int{npq, npq, npt}
		} else {
			isf.Limits = []float64{-l, l, -l, l, -l, l}
			isf.Ndiv = []int{npq, npq, npq}
		}
		isf.Nlevels = 1
		isf.Frange = []float64{0, 0}
	}
	setisof(isf1)
	isf1.CmapNclrs = 24
	isf1.CmapType = "warm"
	if extraconf != nil {
		extraconf(isf1)
	}
	setisof(isf2)
	isf2.CmapNclrs = 0
	isf2.Color = []float64{0, 0, 0, 1}
	isf2.ShowWire = true

	// run visualisation
	scn.Run()
}

// String returns information of this object
func (o *IsoFun) String() (l string) {
	l = io.Sf("  a=%v b=%v β=%v ϵ=%v\n", o.a, o.b, o.β, o.ϵ)
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
	io.Pfblue("dNdL    = %v\n", o.dNdL)
	io.Pfblue("dndL    = %v\n", o.dndL)
	io.Pfblue("d2ndLdL = %v\n", o.d2ndLdL)
	io.Pfblue("d2pdLdL = %v\n", o.d2pdLdL)
	io.Pfblue("d2qdLdL = %v\n", o.d2qdLdL)
	io.Pfblue2("p       = %v\n", o.p)
	io.Pfblue2("q       = %v\n", o.q)
	io.Pfblue2("dpdL    = %v\n", o.dpdL)
	io.Pfblue2("dqdL    = %v\n", o.dqdL)
	io.Pfblue2("dfdp    = %v\n", o.dfdp)
	io.Pfblue2("dfdq    = %v\n", o.dfdq)
	io.Pfblue2("DfdL    = %v\n", o.DfdL)
	io.Pfblue("d2fdp2  = %v\n", o.d2fdp2)
	io.Pfblue("d2fdq2  = %v\n", o.d2fdq2)
	io.Pfblue("d2fdpdq = %v\n", o.d2fdpdq)
	io.Pfblue("DgdL    = %v\n", o.DgdL)
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
