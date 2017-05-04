// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/ode"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// functions
type TimeSpace func(t, x float64) float64
type TimeOnly func(t float64) float64
type SpaceOnly func(x float64) float64

// Diffusion1d problem data
type Diffusion1d struct {

	// problem data
	Kx float64   // kx: diffusion coefficient
	Xf float64   // Xf: maximum x
	U  TimeSpace // u: true solution
	S  TimeSpace // s: source function
	U0 SpaceOnly // u0: initial values function
	UL TimeOnly  // uL: boundary condition function at left (x=0)
	UR TimeOnly  // uR :boundary condition function at right (x=xf)

	// meshgrid for plotting
	xx [][]float64 // x values
	tt [][]float64 // t values
	uu [][]float64 // u values
	ww [][]float64 // true u values
	ee [][]float64 // error values
}

// Problem returns problem by index
func Problem(problem int) (p *Diffusion1d) {

	p = new(Diffusion1d)

	switch problem {

	case 1:

		// constants
		l := -0.1     // λ
		pi := math.Pi // π

		// data
		p.Kx = 1.0
		p.Xf = 1.0

		// functions
		p.U = func(t, x float64) float64 { return math.Exp(l*t) * math.Sin(pi*x) }
		p.S = func(t, x float64) float64 { return (l + p.Kx*pi*pi) * p.U(t, x) }
		p.U0 = func(x float64) float64 { return math.Sin(pi * x) }
		p.UL = func(t float64) float64 { return 0 }
		p.UR = func(t float64) float64 { return 0 }

	default:
		chk.Panic("problem index = %d is not available", problem)
	}
	return
}

// Version1 returns the ODE function that uses L matrix
//   Input:
//     m -- problem dimension; e.g. nx - 2
//     dx -- grid spacing
//     printL -- print L matrix
//   Output:
//     fcn -- f(t,Ub) function in dUbdt = f(t,Ub)
//     jac -- dfdUb function
func (o Diffusion1d) Version1(m int, dx float64, printL bool) (fcn ode.Cb_fcn, jac ode.Cb_jac) {

	// auxiliary coefficient
	p := o.Kx / (dx * dx)

	// L matrix
	var Ltri la.Triplet
	Ltri.Init(m, m, 3*m-2)
	for i := 0; i < m; i++ {
		if i > 0 {
			Ltri.Put(i, i-1, p)
		}
		if i < m-1 {
			Ltri.Put(i, i+1, p)
		}
		Ltri.Put(i, i, -2*p)
	}
	L := Ltri.ToMatrix(nil)

	// print L matrix
	if printL {
		la.PrintMat("L", L.ToDense(), "%9.2f", false)
	}

	// f(t,Ub) function
	fcn = func(f []float64, dt, t float64, Ub []float64) (err error) {
		la.VecApplyFunc(f, func(i int, _ float64) (res float64) {
			x := float64(i+1) * dx
			res = o.S(t, x)
			if i == 1 {
				res += p * o.UL(t)
			}
			if i == m-1 {
				res += p * o.UR(t)
			}
			return
		})
		la.SpMatVecMulAdd(f, 1, L, Ub) // f += L*U
		return
	}
	return
}

// Error computes the error
func (o Diffusion1d) Error(dx float64, sol *ode.Solver) (emax float64) {
	nt := sol.IdxSave
	nx := len(sol.Yvalues) + 2 // +2 => left and right points
	for i := 0; i < nt; i++ {
		t := sol.Xvalues[i]
		for j := 1; j < nx-1; j++ {
			x := float64(j) * dx
			u := sol.Yvalues[j-1][i]
			e := math.Abs(o.U(t, x) - u)
			emax = utl.Max(emax, e)
		}
	}
	return
}

// Results computes the meshgrid values and error
//  nlt -- number of lines along t
//  nlx -- number of lines along x
func (o *Diffusion1d) Results(dx float64, sol *ode.Solver, nlt, nlx int) {

	// total number of points along t and x
	nt := sol.IdxSave
	nx := len(sol.Yvalues) + 2      // +2 => left and right points
	ii := utl.GetStrides(nt-1, nlt) // -1 => because we want indices
	jj := utl.GetStrides(nx-1, nlx) // -1 => because we want indices

	// meshgrid for plotting
	mt := len(ii)
	mx := len(jj)
	o.xx = la.MatAlloc(mt, mx)
	o.tt = la.MatAlloc(mt, mx)
	o.uu = la.MatAlloc(mt, mx)
	o.ww = la.MatAlloc(mt, mx)
	o.ee = la.MatAlloc(mt, mx)
	for i, I := range ii {
		for j, J := range jj {
			o.xx[i][j] = float64(J) * dx
			o.tt[i][j] = sol.Xvalues[I]
			if J == 0 {
				o.uu[i][j] = o.UL(o.tt[i][j])
			} else if J == nx-1 {
				o.uu[i][j] = o.UR(o.tt[i][j])
			} else {
				o.uu[i][j] = sol.Yvalues[J-1][I]
			}
			o.ww[i][j] = o.U(o.tt[i][j], o.xx[i][j])
			o.ee[i][j] = math.Abs(o.ww[i][j] - o.uu[i][j])
		}
	}
	return
}

// PlotRes3d plots results in 3D
func (o Diffusion1d) PlotRes3d(dirout, fnkey string) {
	plt.Reset(false, nil)
	plt.Wireframe(o.tt, o.xx, o.uu, &plt.A{C: "red", CmapIdx: 0})
	plt.Wireframe(o.tt, o.xx, o.ww, &plt.A{C: "blue"})
	plt.SetLabels3d("$t$", "$x$", "$u$", nil)
	plt.Save(dirout, fnkey)
}

// PlotErr3d plots error in 3D
func (o Diffusion1d) PlotErr3d(dirout, fnkey string) {
	plt.Reset(false, nil)
	plt.Surface(o.tt, o.xx, o.ee, &plt.A{})
	plt.SetLabels3d("$t$", "$x$", "$e$", nil)
	plt.Save(dirout, fnkey)
}
