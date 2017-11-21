// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/la"
)

// facTfinite defines a structure to implement transfinite mappings
type facTfinite struct{}

// FactoryTfinite generates transfinite mappings
var FactoryTfinite = facTfinite{}

// 2D surfaces ////////////////////////////////////////////////////////////////////////////////////

// Surf2dQuad generates a transfinite mapping of a quadrilateral
//
//   A,B,C,D -- the four corners (counter-clockwise order)
//
//              Γ[3](r)
//             D───────C
//             │       │
//      Γ[0](s)│       │Γ[1](s)
//             │       │
//             A───────B
//              Γ[2](r)
//
func (o facTfinite) Surf2dQuad(A, B, C, D []float64) (surf *Transfinite) {

	u0 := []float64{D[0] - A[0], D[1] - A[1]}
	u1 := []float64{C[0] - B[0], C[1] - B[1]}
	u2 := []float64{B[0] - A[0], B[1] - A[1]}
	u3 := []float64{C[0] - D[0], C[1] - D[1]}

	surf = NewTransfinite2d([]fun.Vs{

		// Γ[0](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			x[0] = A[0] + (1.0+s)*u0[0]/2.0
			x[1] = A[1] + (1.0+s)*u0[1]/2.0
		},

		// Γ[1](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			x[0] = B[0] + (1.0+s)*u1[0]/2.0
			x[1] = B[1] + (1.0+s)*u1[1]/2.0
		},

		// Γ[2](r)
		func(x la.Vector, r float64) { // r ϵ [-1,+1]
			x[0] = A[0] + (1.0+r)*u2[0]/2.0
			x[1] = A[1] + (1.0+r)*u2[1]/2.0
		},

		// Γ[3](r)
		func(x la.Vector, r float64) { // r ϵ [-1,+1]
			x[0] = D[0] + (1.0+r)*u3[0]/2.0
			x[1] = D[1] + (1.0+r)*u3[1]/2.0
		},

		// first order derivatives

	}, []fun.Vs{

		// dΓ[0]/ds
		func(dxds la.Vector, s float64) {
			dxds[0] = u0[0] / 2.0
			dxds[1] = u0[1] / 2.0
		},

		// dΓ[1]/ds
		func(dxds la.Vector, s float64) {
			dxds[0] = u1[0] / 2.0
			dxds[1] = u1[1] / 2.0
		},

		// dΓ[2]/dr
		func(dxdr la.Vector, r float64) {
			dxdr[0] = u2[0] / 2.0
			dxdr[1] = u2[1] / 2.0
		},

		// dΓ[3]/dr
		func(dxdr la.Vector, r float64) {
			dxdr[0] = u3[0] / 2.0
			dxdr[1] = u3[1] / 2.0
		},
	}, nil)
	return
}

// Surf2dQuarterRing generates a transfinite mapping of a quarter of a ring centered @ (0,0)
//   a -- inner radius
//   b -- outer radius
func (o facTfinite) Surf2dQuarterRing(a, b float64) (surf *Transfinite) {

	π := math.Pi
	surf = NewTransfinite2d([]fun.Vs{

		// B[0](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			θ := π * (1.0 + s) / 4.0
			x[0] = a * math.Cos(θ)
			x[1] = a * math.Sin(θ)
		},

		// B[1](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			θ := π * (1.0 + s) / 4.0
			x[0] = b * math.Cos(θ)
			x[1] = b * math.Sin(θ)
		},

		// B[2](r)
		func(x la.Vector, r float64) { // r ϵ [-1,+1]
			x[0] = a + 0.5*(1.0+r)*(b-a)
			x[1] = 0.0
		},

		// B[3](r)
		func(x la.Vector, r float64) { // r ϵ [-1,+1]
			x[0] = 0.0
			x[1] = a + 0.5*(1.0+r)*(b-a)
		},

		// first order derivatives

	}, []fun.Vs{

		// dB[0]/ds
		func(dxds la.Vector, s float64) {
			θ := π * (1.0 + s) / 4.0
			dxds[0] = -a * math.Sin(θ) * π / 4.0
			dxds[1] = +a * math.Cos(θ) * π / 4.0
		},

		// dB[1]/ds
		func(dxds la.Vector, s float64) {
			θ := π * (1.0 + s) / 4.0
			dxds[0] = -b * math.Sin(θ) * π / 4.0
			dxds[1] = +b * math.Cos(θ) * π / 4.0
		},

		// dB[2]/dr
		func(dxdr la.Vector, r float64) {
			dxdr[0] = 0.5 * (b - a)
			dxdr[1] = 0.0
		},

		// dB[3]/dr
		func(dxdr la.Vector, r float64) {
			dxdr[0] = 0.0
			dxdr[1] = 0.5 * (b - a)
		},

		// second order derivatives

	}, []fun.Vs{

		// d²B[0]/ds²
		func(ddxdss la.Vector, s float64) {
			θ := π * (1.0 + s) / 4.0
			ddxdss[0] = -a * math.Cos(θ) * π * π / 16.0
			ddxdss[1] = -a * math.Sin(θ) * π * π / 16.0
		},

		// d²B[1]/ds²
		func(ddxdss la.Vector, s float64) {
			θ := π * (1.0 + s) / 4.0
			ddxdss[0] = -b * math.Cos(θ) * π * π / 16.0
			ddxdss[1] = -b * math.Sin(θ) * π * π / 16.0
		},

		// d²B[2]/dr²
		func(ddxdrr la.Vector, r float64) {
			ddxdrr[0] = 0.0
			ddxdrr[1] = 0.0
		},

		// d²B[3]/dr²
		func(ddxdrr la.Vector, r float64) {
			ddxdrr[0] = 0.0
			ddxdrr[1] = 0.0
		},
	})
	return
}

// Surf2dQuarterPerfLozenge generates a transfinite mapping of a quarter of a perforated lozenge
// (diamond shape) centered @ (0,0)
//   a -- inner radius
//   b -- diagonal of lozenge (diamond)
func (o facTfinite) Surf2dQuarterPerfLozenge(a, b float64) (surf *Transfinite) {

	π := math.Pi
	surf = NewTransfinite2d([]fun.Vs{

		// B[0](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			θ := π * (1.0 + s) / 4.0
			x[0] = a * math.Cos(θ)
			x[1] = a * math.Sin(θ)
		},

		// B[1](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			x[0] = b * 0.5 * (1.0 - s)
			x[1] = b * 0.5 * (1.0 + s)
		},

		// B[2](r)
		func(x la.Vector, r float64) { // r ϵ [-1,+1]
			x[0] = a + 0.5*(1.0+r)*(b-a)
			x[1] = 0.0
		},

		// B[3](r)
		func(x la.Vector, r float64) { // r ϵ [-1,+1]
			x[0] = 0.0
			x[1] = a + 0.5*(1.0+r)*(b-a)
		},

		// first order derivatives

	}, []fun.Vs{

		// dB[0]/ds
		func(dxds la.Vector, s float64) {
			θ := π * (1.0 + s) / 4.0
			dxds[0] = -a * math.Sin(θ) * π / 4.0
			dxds[1] = +a * math.Cos(θ) * π / 4.0
		},

		// dB[1]/ds
		func(dxds la.Vector, s float64) {
			dxds[0] = -b * 0.5
			dxds[1] = +b * 0.5
		},

		// dB[2]/dr
		func(dxdr la.Vector, r float64) {
			dxdr[0] = 0.5 * (b - a)
			dxdr[1] = 0.0
		},

		// dB[3]/dr
		func(dxdr la.Vector, r float64) {
			dxdr[0] = 0.0
			dxdr[1] = 0.5 * (b - a)
		},

		// second order derivatives

	}, []fun.Vs{

		// d²B[0]/ds²
		func(ddxdss la.Vector, s float64) {
			θ := π * (1.0 + s) / 4.0
			ddxdss[0] = -a * math.Cos(θ) * π * π / 16.0
			ddxdss[1] = -a * math.Sin(θ) * π * π / 16.0
		},

		// d²B[1]/ds²
		func(ddxdss la.Vector, s float64) {
			ddxdss[0] = 0.0
			ddxdss[1] = 0.0
		},

		// d²B[2]/dr²
		func(ddxdrr la.Vector, r float64) {
			ddxdrr[0] = 0.0
			ddxdrr[1] = 0.0
		},

		// d²B[3]/dr²
		func(ddxdrr la.Vector, r float64) {
			ddxdrr[0] = 0.0
			ddxdrr[1] = 0.0
		},
	})
	return
}

// solids /////////////////////////////////////////////////////////////////////////////////////////

// SolidCube generates a transfinite mapping of a cube
func (o facTfinite) SolidCube(lx, ly, lz float64) (solid *Transfinite) {

	solid = NewTransfinite3d([]fun.Vss{

		// B[0](s,t)
		func(x la.Vector, s, t float64) {
			x[0] = 0.0
			x[1] = (1.0 + s) * ly / 2.0
			x[2] = (1.0 + t) * lz / 2.0
		},

		// B[1](s,t)
		func(x la.Vector, s, t float64) {
			x[0] = lx
			x[1] = (1.0 + s) * ly / 2.0
			x[2] = (1.0 + t) * lz / 2.0
		},

		// B[2](r,t)
		func(x la.Vector, r, t float64) {
			x[0] = (1.0 + r) * lx / 2.0
			x[1] = 0.0
			x[2] = (1.0 + t) * lz / 2.0
		},

		// B[3](r,t)
		func(x la.Vector, r, t float64) {
			x[0] = (1.0 + r) * lx / 2.0
			x[1] = ly
			x[2] = (1.0 + t) * lz / 2.0
		},

		// B[4](r,s)
		func(x la.Vector, r, s float64) {
			x[0] = (1.0 + r) * lx / 2.0
			x[1] = (1.0 + s) * ly / 2.0
			x[2] = 0.0
		},

		// B[5](r,s)
		func(x la.Vector, r, s float64) {
			x[0] = (1.0 + r) * lx / 2.0
			x[1] = (1.0 + s) * ly / 2.0
			x[2] = lz
		},

		// first order derivatives

	}, []fun.Vvss{

		// Bd[0](s,t)
		func(dxds, dxdt la.Vector, s, t float64) {
			dxds[0] = 0.0
			dxds[1] = ly / 2.0
			dxds[2] = 0.0
			dxdt[0] = 0.0
			dxdt[1] = 0.0
			dxdt[2] = lz / 2.0
		},

		// Bd[1](s,t)
		func(dxds, dxdt la.Vector, s, t float64) {
			dxds[0] = 0.0
			dxds[1] = ly / 2.0
			dxds[2] = 0.0
			dxdt[0] = 0.0
			dxdt[1] = 0.0
			dxdt[2] = lz / 2.0
		},

		// Bd[2](r,t)
		func(dxdr, dxdt la.Vector, r, t float64) {
			dxdr[0] = lx / 2.0
			dxdr[1] = 0.0
			dxdr[2] = 0.0
			dxdt[0] = 0.0
			dxdt[1] = 0.0
			dxdt[2] = lz / 2.0
		},

		// Bd[3](r,t)
		func(dxdr, dxdt la.Vector, r, t float64) {
			dxdr[0] = lx / 2.0
			dxdr[1] = 0.0
			dxdr[2] = 0.0
			dxdt[0] = 0.0
			dxdt[1] = 0.0
			dxdt[2] = lz / 2.0
		},

		// Bd[4](r,s)
		func(dxdr, dxds la.Vector, r, s float64) {
			dxdr[0] = lx / 2.0
			dxdr[1] = 0.0
			dxdr[2] = 0.0
			dxds[0] = 0.0
			dxds[1] = ly / 2.0
			dxds[2] = 0.0
		},

		// Bd[5](r,s)
		func(dxdr, dxds la.Vector, r, s float64) {
			dxdr[0] = lx / 2.0
			dxdr[1] = 0.0
			dxdr[2] = 0.0
			dxds[0] = 0.0
			dxds[1] = ly / 2.0
			dxds[2] = 0.0
		},
	}, nil)
	return
}

// SolidQuarterRing generates a transfinite mapping of a quarter of a 3d ring centered @ (0,0)
//   a -- inner radius
//   b -- outer radius
//   h -- thickness
func (o facTfinite) SolidQuarterRing(a, b, h float64) (solid *Transfinite) {

	π := math.Pi
	tmp := la.NewVector(2)
	x2d := la.NewVector(2)
	u2d := la.NewVector(2)
	dxdr2d := la.NewVector(2)
	dxds2d := la.NewVector(2)
	ddxdrr2d := la.NewVector(2)
	ddxdss2d := la.NewVector(2)
	ddxdrs2d := la.NewVector(2)

	surf := FactoryTfinite.Surf2dQuarterRing(a, b)

	solid = NewTransfinite3d([]fun.Vss{

		// B[0](s,t)
		func(x la.Vector, s, t float64) {
			u2d[0], u2d[1] = s, t
			surf.Point(x2d, u2d)
			x[0] = 0.0
			x[1] = x2d[0]
			x[2] = x2d[1]
		},

		// B[1](s,t)
		func(x la.Vector, s, t float64) {
			u2d[0], u2d[1] = s, t
			surf.Point(x2d, u2d)
			x[0] = h
			x[1] = x2d[0]
			x[2] = x2d[1]
		},

		// B[2](r,t)
		func(x la.Vector, r, t float64) {
			θ := (1.0 + t) * π / 4.0
			x[0] = (1.0 + r) * h / 2.0
			x[1] = a * math.Cos(θ)
			x[2] = a * math.Sin(θ)
		},

		// B[3](r,t)
		func(x la.Vector, r, t float64) {
			θ := (1.0 + t) * π / 4.0
			x[0] = (1.0 + r) * h / 2.0
			x[1] = b * math.Cos(θ)
			x[2] = b * math.Sin(θ)
		},

		// B[4](r,s)
		func(x la.Vector, r, s float64) {
			x[0] = (1.0 + r) * h / 2.0
			x[1] = a + (1.0+s)*(b-a)/2.0
			x[2] = 0.0
		},

		// B[5](r,s)
		func(x la.Vector, r, s float64) {
			x[0] = (1.0 + r) * h / 2.0
			x[1] = 0.0
			x[2] = a + (1.0+s)*(b-a)/2.0
		},

		// first order derivatives

	}, []fun.Vvss{

		// Bd[0](s,t)
		func(dxds, dxdt la.Vector, s, t float64) {
			u2d[0], u2d[1] = s, t
			surf.PointAndDerivs(tmp, dxdr2d, dxds2d, nil, nil, nil, nil, nil, nil, nil, u2d)
			dxds[0] = 0.0
			dxds[1] = dxdr2d[0]
			dxds[2] = dxdr2d[1]
			dxdt[0] = 0.0
			dxdt[1] = dxds2d[0]
			dxdt[2] = dxds2d[1]
		},

		// Bd[1](s,t)
		func(dxds, dxdt la.Vector, s, t float64) {
			u2d[0], u2d[1] = s, t
			surf.PointAndDerivs(tmp, dxdr2d, dxds2d, nil, nil, nil, nil, nil, nil, nil, u2d)
			dxds[0] = 0.0
			dxds[1] = dxdr2d[0]
			dxds[2] = dxdr2d[1]
			dxdt[0] = 0.0
			dxdt[1] = dxds2d[0]
			dxdt[2] = dxds2d[1]
		},

		// Bd[2](r,t)
		func(dxdr, dxdt la.Vector, r, t float64) {
			θ := (1.0 + t) * π / 4.0
			dxdr[0] = h / 2.0
			dxdr[1] = 0.0
			dxdr[2] = 0.0
			dxdt[0] = 0.0
			dxdt[1] = -a * math.Sin(θ) * π / 4.0
			dxdt[2] = +a * math.Cos(θ) * π / 4.0
		},

		// Bd[3](r,t)
		func(dxdr, dxdt la.Vector, r, t float64) {
			θ := (1.0 + t) * π / 4.0
			dxdr[0] = h / 2.0
			dxdr[1] = 0.0
			dxdr[2] = 0.0
			dxdt[0] = 0.0
			dxdt[1] = -b * math.Sin(θ) * π / 4.0
			dxdt[2] = +b * math.Cos(θ) * π / 4.0
		},

		// Bd[4](r,s)
		func(dxdr, dxds la.Vector, r, s float64) {
			dxdr[0] = h / 2.0
			dxdr[1] = 0.0
			dxdr[2] = 0.0
			dxds[0] = 0.0
			dxds[1] = (b - a) / 2.0
			dxds[2] = 0.0
		},

		// Bd[5](r,s)
		func(dxdr, dxds la.Vector, r, s float64) {
			dxdr[0] = h / 2.0
			dxdr[1] = 0.0
			dxdr[2] = 0.0
			dxds[0] = 0.0
			dxds[1] = 0.0
			dxds[2] = (b - a) / 2.0
		},

		// second order derivatives

	}, []fun.Vvvss{

		// Bdd[0](s,t)
		func(ddxdss, ddxdtt, ddxdst la.Vector, s, t float64) {
			u2d[0], u2d[1] = s, t
			surf.PointAndDerivs(tmp, tmp, tmp, nil, ddxdrr2d, ddxdss2d, tmp, ddxdrs2d, tmp, tmp, u2d)
			ddxdss[0] = 0.0
			ddxdss[1] = ddxdrr2d[0]
			ddxdss[2] = ddxdrr2d[1]

			ddxdtt[0] = 0.0
			ddxdtt[1] = ddxdss2d[0]
			ddxdtt[2] = ddxdss2d[1]

			ddxdst[0] = 0.0
			ddxdst[1] = ddxdrs2d[0]
			ddxdst[2] = ddxdrs2d[1]
		},

		// Bdd[1](s,t)
		func(ddxdss, ddxdtt, ddxdst la.Vector, s, t float64) {
			u2d[0], u2d[1] = s, t
			surf.PointAndDerivs(tmp, tmp, tmp, nil, ddxdrr2d, ddxdss2d, tmp, ddxdrs2d, tmp, tmp, u2d)
			ddxdss[0] = 0.0
			ddxdss[1] = ddxdrr2d[0]
			ddxdss[2] = ddxdrr2d[1]

			ddxdtt[0] = 0.0
			ddxdtt[1] = ddxdss2d[0]
			ddxdtt[2] = ddxdss2d[1]

			ddxdst[0] = 0.0
			ddxdst[1] = ddxdrs2d[0]
			ddxdst[2] = ddxdrs2d[1]
		},

		// Bdd[2](r,t)
		func(ddxdrr, ddxdtt, ddxdrt la.Vector, r, t float64) {
			θ := (1.0 + t) * π / 4.0
			ddxdrr[0] = 0.0
			ddxdrr[1] = 0.0
			ddxdrr[2] = 0.0

			ddxdtt[0] = 0.0
			ddxdtt[1] = -a * math.Cos(θ) * π * π / 16.0
			ddxdtt[2] = -a * math.Sin(θ) * π * π / 16.0

			ddxdrt[0] = 0.0
			ddxdrt[1] = 0.0
			ddxdrt[2] = 0.0
		},

		// Bdd[3](r,t)
		func(ddxdrr, ddxdtt, ddxdrt la.Vector, r, t float64) {
			θ := (1.0 + t) * π / 4.0
			ddxdrr[0] = 0.0
			ddxdrr[1] = 0.0
			ddxdrr[2] = 0.0

			ddxdtt[0] = 0.0
			ddxdtt[1] = -b * math.Cos(θ) * π * π / 16.0
			ddxdtt[2] = -b * math.Sin(θ) * π * π / 16.0

			ddxdrt[0] = 0.0
			ddxdrt[1] = 0.0
			ddxdrt[2] = 0.0
		},

		// Bdd[4](r,s)
		func(ddxdrr, ddxdss, ddxdrs la.Vector, r, s float64) {
			ddxdrr[0] = 0.0
			ddxdrr[1] = 0.0
			ddxdrr[2] = 0.0

			ddxdss[0] = 0.0
			ddxdss[1] = 0.0
			ddxdss[2] = 0.0

			ddxdrs[0] = 0.0
			ddxdrs[1] = 0.0
			ddxdrs[2] = 0.0
		},

		// Bdd[5](r,s)
		func(ddxdrr, ddxdss, ddxdrs la.Vector, r, s float64) {
			ddxdrr[0] = 0.0
			ddxdrr[1] = 0.0
			ddxdrr[2] = 0.0

			ddxdss[0] = 0.0
			ddxdss[1] = 0.0
			ddxdss[2] = 0.0

			ddxdrs[0] = 0.0
			ddxdrs[1] = 0.0
			ddxdrs[2] = 0.0
		},
	})
	return
}
