// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux

package ode

func (o *Solver) init_mpi() {
}

func radau5_step_mpi(o *Solver, y0 []float64, x0 float64, args ...interface{}) (rerr float64, err error) {
	return
}
