// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!darwin

package ode

import "github.com/cpmech/gosl/mpi"

func (o *Solver) init_mpi() {
	if mpi.IsOn() {
		o.root = (mpi.Rank() == 0)
		if mpi.Size() > 1 {
			o.Distr = true
		}
	}
}
