// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

// HashPoint returns a unique id of a point
func HashPoint(x, xmin, xdel []float64, tol float64) int {
	if tol < 1e-16 {
		chk.Panic("HashPoint: minimum tolerance must be 1e-16. %v is invalid", tol)
	}
	coefs := []float64{11, 101, 1001}
	n := utl.Imin(len(x), 3)
	var hash, xbar float64
	for i := 0; i < n; i++ {
		if xdel[i] > 0 {
			xbar = (x[i] - xmin[i]) / xdel[i]
			if xbar < 0 {
				xbar = 0
			}
			if xbar > 1 {
				xbar = 1
			}
			//io.Pfblue2("xbar=%20v  xbar/tol=%22v  hash=%22v\n", xbar, xbar/tol, (xbar/tol)*coefs[i])
			hash += (xbar / tol) * coefs[i]
		}
	}
	return int(hash)
}
