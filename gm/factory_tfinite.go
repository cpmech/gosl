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

// Surf2dQuarterRing generates a transfinite mapping of a quarter of a ring centered @ (0,0)
//   a -- inner radius
//   b -- outer radius
func (o facTfinite) Surf2dQuarterRing(a, b float64) (surf *Transfinite) {

	π := math.Pi
	surf = NewTransfinite2d(2, []fun.Vs{

		// B[0](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			θ := π * (s + 1) / 4.0
			x[0] = a * math.Cos(θ)
			x[1] = a * math.Sin(θ)
		},

		// B[1](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			θ := π * (s + 1) / 4.0
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
	}, []fun.Vs{

		// dB[0]/ds
		func(dxds la.Vector, s float64) {
			θ := π * (s + 1) / 4.0
			dθds := π / 4.0
			dxds[0] = -a * math.Sin(θ) * dθds
			dxds[1] = +a * math.Cos(θ) * dθds
		},

		// dB[1]/ds
		func(dxds la.Vector, s float64) {
			θ := π * (s + 1) / 4.0
			dθds := π / 4.0
			dxds[0] = -b * math.Sin(θ) * dθds
			dxds[1] = +b * math.Cos(θ) * dθds
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
	})
	return
}

// Surf2dQuarterPerfLozenge generates a transfinite mapping of a quarter of a perforated lozenge
// (diamond shape) centered @ (0,0)
//   a -- inner radius
//   b -- diagonal of lozenge (diamond)
func (o facTfinite) Surf2dQuarterPerfLozenge(a, b float64) (surf *Transfinite) {

	π := math.Pi
	surf = NewTransfinite2d(2, []fun.Vs{

		// B[0](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			θ := π * (s + 1) / 4.0
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
	}, []fun.Vs{

		// dB[0]/ds
		func(dxds la.Vector, s float64) {
			θ := π * (s + 1) / 4.0
			dθds := π / 4.0
			dxds[0] = -a * math.Sin(θ) * dθds
			dxds[1] = +a * math.Cos(θ) * dθds
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
	})
	return
}

// 3D surfaces ////////////////////////////////////////////////////////////////////////////////////

// Surf3dCube generates a transfinite mapping of a cube
func (o facTfinite) Surf3dCube(lx, ly, lz float64) (solid *Transfinite) {

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
	}, []fun.Mss{

		// dB[0]/dst
		func(dxdst *la.Matrix, s, t float64) {
			dxdst.Set(0, 0, 0.0)    // dx[0]/ds
			dxdst.Set(0, 1, 0.0)    // dx[0]/dt
			dxdst.Set(1, 0, ly/2.0) // dx[1]/ds
			dxdst.Set(1, 1, 0.0)    // dx[1]/dt
			dxdst.Set(2, 0, 0.0)    // dx[2]/ds
			dxdst.Set(2, 1, lz/2.0) // dx[2]/dt
		},

		// dB[1]/dst
		func(dxdst *la.Matrix, s, t float64) {
			dxdst.Set(0, 0, 0.0)    // dx[0]/ds
			dxdst.Set(0, 1, 0.0)    // dx[0]/dt
			dxdst.Set(1, 0, ly/2.0) // dx[1]/ds
			dxdst.Set(1, 1, 0.0)    // dx[1]/dt
			dxdst.Set(2, 0, 0.0)    // dx[2]/ds
			dxdst.Set(2, 1, lz/2.0) // dx[2]/dt
		},

		// dB[2]/drt
		func(dxdrt *la.Matrix, r, t float64) {
			dxdrt.Set(0, 0, lx/2.0) // dx[0]/dr
			dxdrt.Set(0, 1, 0.0)    // dx[0]/dt
			dxdrt.Set(1, 0, 0.0)    // dx[1]/dr
			dxdrt.Set(1, 1, 0.0)    // dx[1]/dt
			dxdrt.Set(2, 0, 0.0)    // dx[2]/dr
			dxdrt.Set(2, 1, lz/2.0) // dx[2]/dt
		},

		// dB[3]/drt
		func(dxdrt *la.Matrix, r, t float64) {
			dxdrt.Set(0, 0, lx/2.0) // dx[0]/dr
			dxdrt.Set(0, 1, 0.0)    // dx[0]/dt
			dxdrt.Set(1, 0, 0.0)    // dx[1]/dr
			dxdrt.Set(1, 1, 0.0)    // dx[1]/dt
			dxdrt.Set(2, 0, 0.0)    // dx[2]/dr
			dxdrt.Set(2, 1, lz/2.0) // dx[2]/dt
		},

		// dB[4]/drs
		func(dxdrs *la.Matrix, r, s float64) {
			dxdrs.Set(0, 0, lx/2.0) // dx[0]/dr
			dxdrs.Set(0, 1, 0.0)    // dx[0]/ds
			dxdrs.Set(1, 0, 0.0)    // dx[1]/dr
			dxdrs.Set(1, 1, ly/2.0) // dx[1]/ds
			dxdrs.Set(2, 0, 0.0)    // dx[2]/dr
			dxdrs.Set(2, 1, 0.0)    // dx[2]/ds
		},

		// dB[5]/drs
		func(dxdrs *la.Matrix, r, s float64) {
			dxdrs.Set(0, 0, lx/2.0) // dx[0]/dr
			dxdrs.Set(0, 1, 0.0)    // dx[0]/ds
			dxdrs.Set(1, 0, 0.0)    // dx[1]/dr
			dxdrs.Set(1, 1, ly/2.0) // dx[1]/ds
			dxdrs.Set(2, 0, 0.0)    // dx[2]/dr
			dxdrs.Set(2, 1, 0.0)    // dx[2]/ds
		},
	})
	return
}

// Surf3dQuarterRing generates a transfinite mapping of a quarter of a 3d ring centered @ (0,0)
//   a -- inner radius
//   b -- outer radius
//   h -- thickness
func (o facTfinite) Surf3dQuarterRing(a, b, h float64) (solid *Transfinite) {

	π := math.Pi
	x2d := la.NewVector(2)
	u2d := la.NewVector(2)
	dxdu2d := la.NewMatrix(2, 2)
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
	}, []fun.Mss{

		// dB[0]/dst
		func(dxdst *la.Matrix, s, t float64) {
			u2d[0], u2d[1] = s, t
			surf.Derivs(dxdu2d, u2d)
			dxdst.Set(0, 0, 0.0)              // dx[0]/ds
			dxdst.Set(0, 1, 0.0)              // dx[0]/dt
			dxdst.Set(1, 0, dxdu2d.Get(0, 0)) // dx[1]/ds
			dxdst.Set(1, 1, dxdu2d.Get(0, 1)) // dx[1]/dt
			dxdst.Set(2, 0, dxdu2d.Get(1, 0)) // dx[2]/ds
			dxdst.Set(2, 1, dxdu2d.Get(1, 1)) // dx[2]/dt
		},

		// dB[1]/dst
		func(dxdst *la.Matrix, s, t float64) {
			u2d[0], u2d[1] = s, t
			surf.Derivs(dxdu2d, u2d)
			dxdst.Set(0, 0, 0.0)              // dx[0]/ds
			dxdst.Set(0, 1, 0.0)              // dx[0]/dt
			dxdst.Set(1, 0, dxdu2d.Get(0, 0)) // dx[1]/ds
			dxdst.Set(1, 1, dxdu2d.Get(0, 1)) // dx[1]/dt
			dxdst.Set(2, 0, dxdu2d.Get(1, 0)) // dx[2]/ds
			dxdst.Set(2, 1, dxdu2d.Get(1, 1)) // dx[2]/dt
		},

		// dB[2]/drt
		func(dxdrt *la.Matrix, r, t float64) {
			θ := (1.0 + t) * π / 4.0
			dxdrt.Set(0, 0, h/2.0)                // dx[0]/dr
			dxdrt.Set(0, 1, 0.0)                  // dx[0]/dt
			dxdrt.Set(1, 0, 0.0)                  // dx[1]/dr
			dxdrt.Set(1, 1, -a*math.Sin(θ)*π/4.0) // dx[1]/dt
			dxdrt.Set(2, 0, 0.0)                  // dx[2]/dr
			dxdrt.Set(2, 1, a*math.Cos(θ)*π/4.0)  // dx[2]/dt
		},

		// dB[3]/drt
		func(dxdrt *la.Matrix, r, t float64) {
			θ := (1.0 + t) * π / 4.0
			dxdrt.Set(0, 0, h/2.0)                // dx[0]/dr
			dxdrt.Set(0, 1, 0.0)                  // dx[0]/dt
			dxdrt.Set(1, 0, 0.0)                  // dx[1]/dr
			dxdrt.Set(1, 1, -b*math.Sin(θ)*π/4.0) // dx[1]/dt
			dxdrt.Set(2, 0, 0.0)                  // dx[2]/dr
			dxdrt.Set(2, 1, b*math.Cos(θ)*π/4.0)  // dx[2]/dt
		},

		// dB[4]/drs
		func(dxdrs *la.Matrix, r, s float64) {
			dxdrs.Set(0, 0, h/2.0)     // dx[0]/dr
			dxdrs.Set(0, 1, 0.0)       // dx[0]/ds
			dxdrs.Set(1, 0, 0.0)       // dx[1]/dr
			dxdrs.Set(1, 1, (b-a)/2.0) // dx[1]/ds
			dxdrs.Set(2, 0, 0.0)       // dx[2]/dr
			dxdrs.Set(2, 1, 0.0)       // dx[2]/ds
		},

		// dB[5]/drs
		func(dxdrs *la.Matrix, r, s float64) {
			dxdrs.Set(0, 0, h/2.0)     // dx[0]/dr
			dxdrs.Set(0, 1, 0.0)       // dx[0]/ds
			dxdrs.Set(1, 0, 0.0)       // dx[1]/dr
			dxdrs.Set(1, 1, 0.0)       // dx[1]/ds
			dxdrs.Set(2, 0, 0.0)       // dx[2]/dr
			dxdrs.Set(2, 1, (b-a)/2.0) // dx[2]/ds
		},
	})
	return
}