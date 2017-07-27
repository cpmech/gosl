// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import "github.com/cpmech/gosl/la"

func fweuler_accept(o *Solver, y la.Vector) {
}

// forward-Euler
func fweuler_step(o *Solver, y la.Vector, x float64) (rerr float64, err error) {
	o.Nfeval += 1
	err = o.fcn(o.f[0], o.h, x, y)
	if err != nil {
		return
	}
	for i := 0; i < o.ndim; i++ {
		y[i] += o.h * o.f[0][i]
	}
	return 1e+20, err // must not be used with automatic substepping
}
