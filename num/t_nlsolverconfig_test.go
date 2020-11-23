// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"testing"

	"gosl/chk"
)

func TestNlsConfig01(tst *testing.T) {

	// verbose()
	chk.PrintTitle("NlsConfig01")

	c := NewNlSolverConfig()

	// flags
	res := []bool{
		c.Verbose,
		c.ConstantJacobian,
		c.LineSearch,
		c.EnforceConvRate,
		c.useDenseSolver,
		c.hasJacobianFunction,
		c.LinSolConfig.Verbose,
	}
	correct := []bool{false, false, false, false, false, false, false}
	chk.Bools(tst, "flags", res, correct)

	// tolerances
	chk.Float64(tst, "atol", 1e-15, c.atol, 1e-8)
	chk.Float64(tst, "rtol", 1e-15, c.rtol, 1e-8)
	chk.Float64(tst, "ftol", 1e-15, c.ftol, 1e-9)
	chk.Float64(tst, "fnewt", 1e-15, c.fnewt, 0.0001)
}
