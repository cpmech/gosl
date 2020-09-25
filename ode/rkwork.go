// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import "gosl/la"

// rkwork holds the Runge-Kuta "workspace" variables
type rkwork struct {

	// workspace
	nstg int         // number of stages
	ndim int         // dimension of y vector
	u    la.Vector   // u[stg] = x + h*c[stg]
	v    []la.Vector // v[stg][dim] = ya[dim] + h*sum(a[stg][j]*f[j][dim], j, nstg)
	f    []la.Vector // f[stg][dim] = f(u[stg], v[stg][dim])

	// step data
	rs    float64   // stiffness ratio ρ = ‖ k[s] - k[s-1] ‖ / ‖ v[s] - v[s-1] ‖
	h     float64   // current stepsize
	hPrev float64   // previous stepsize
	first bool      // first step
	f0    la.Vector // f(x,y) before step
	scal  la.Vector // scal = Atol + Rtol*abs(y)

	// step control data
	reuseJacAndDecOnce bool    // reuse current Jacobian and current decomposition
	reuseJacOnce       bool    // reuse last Jacobian (only)
	jacIsOK            bool    // Jacobian is OK
	nit                int     // current number of iterations
	eta                float64 // eta tolerance
	theta              float64 // theta variable
	dvfac              float64 // divergence factor
	diverg             bool    // flag diverging step
	reject             bool    // reject step

	// error control
	rerr     float64 // relative error
	rerrPrev float64 // previous relative error

	// stiffness detection
	stiffYes int // counter of "stiff" steps
	stiffNot int // counter of not "stiff" steps
}

// newRKwork returns a new structure
func newRKwork(nstg, ndim int) (o *rkwork) {

	// workspace
	o = new(rkwork)
	o.nstg = nstg
	o.ndim = ndim
	o.u = la.NewVector(o.nstg)
	o.v = make([]la.Vector, o.nstg)
	o.f = make([]la.Vector, o.nstg)
	for i := 0; i < o.nstg; i++ {
		o.v[i] = la.NewVector(o.ndim)
		o.f[i] = la.NewVector(o.ndim)
	}

	// step data
	o.f0 = la.NewVector(o.ndim)
	o.scal = la.NewVector(o.ndim)
	return
}
