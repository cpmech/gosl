// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux

package ode

func (o *ODE) init_mpi() {
}

func radau5_step_mpi(o *ODE, y0 []float64, x0 float64, args ...interface{}) (rerr float64, err error) {
return
}
