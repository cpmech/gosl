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
	surf = NewTransfinite(2, []fun.Vs{

		// B[0](r)
		func(x la.Vector, r float64) { // r ϵ [-1,+1]
			x[0] = a + 0.5*(1.0+r)*(b-a)
			x[1] = 0.0
		},

		// B[1](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			θ := π * (s + 1) / 4.0
			x[0] = b * math.Cos(θ)
			x[1] = b * math.Sin(θ)
		},

		// B[2](r)
		func(x la.Vector, r float64) { // r ϵ [-1,+1]
			x[0] = 0.0
			x[1] = a + 0.5*(1.0+r)*(b-a)
		},

		// B[3](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			θ := π * (s + 1) / 4.0
			x[0] = a * math.Cos(θ)
			x[1] = a * math.Sin(θ)
		},
	}, []fun.Vs{

		// dB[0]/dr
		func(dxdr la.Vector, r float64) {
			dxdr[0] = 0.5 * (b - a)
			dxdr[1] = 0.0
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
			dxdr[0] = 0.0
			dxdr[1] = 0.5 * (b - a)
		},

		// dB[3]/ds
		func(dxds la.Vector, s float64) {
			θ := π * (s + 1) / 4.0
			dθds := π / 4.0
			dxds[0] = -a * math.Sin(θ) * dθds
			dxds[1] = +a * math.Cos(θ) * dθds
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
	surf = NewTransfinite(2, []fun.Vs{

		// B[0](r)
		func(x la.Vector, r float64) { // r ϵ [-1,+1]
			x[0] = a + 0.5*(1.0+r)*(b-a)
			x[1] = 0.0
		},

		// B[1](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			x[0] = b * 0.5 * (1.0 - s)
			x[1] = b * 0.5 * (1.0 + s)
		},

		// B[2](r)
		func(x la.Vector, r float64) { // r ϵ [-1,+1]
			x[0] = 0.0
			x[1] = a + 0.5*(1.0+r)*(b-a)
		},

		// B[3](s)
		func(x la.Vector, s float64) { // s ϵ [-1,+1]
			θ := π * (s + 1) / 4.0
			x[0] = a * math.Cos(θ)
			x[1] = a * math.Sin(θ)
		},
	}, []fun.Vs{

		// dB[0]/dr
		func(dxdr la.Vector, r float64) {
			dxdr[0] = 0.5 * (b - a)
			dxdr[1] = 0.0
		},

		// dB[1]/ds
		func(dxds la.Vector, s float64) {
			dxds[0] = -b * 0.5
			dxds[1] = +b * 0.5
		},

		// dB[2]/dr
		func(dxdr la.Vector, r float64) {
			dxdr[0] = 0.0
			dxdr[1] = 0.5 * (b - a)
		},

		// dB[3]/ds
		func(dxds la.Vector, s float64) {
			θ := π * (s + 1) / 4.0
			dθds := π / 4.0
			dxds[0] = -a * math.Sin(θ) * dθds
			dxds[1] = +a * math.Cos(θ) * dθds
		},
	})
	return
}
