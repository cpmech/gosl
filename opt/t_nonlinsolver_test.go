// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func TestNLS01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("NLS01. NonLinSolver")

	// problem and initial point
	p := Factory.SimpleParaboloid()
	x := la.NewVectorSlice([]float64{1, 1})

	for _, kind := range []string{"conjgrad", "powell", "graddesc"} {
		io.Pf(">>>>>>>>>>>>>>>>>>> running %q <<<<<<<<<<<<<<<<<<<<\n", kind)
		sol := GetNonLinSolver(kind, p)
		fmin := sol.Min(x, nil)
		io.Pf("fmin = %v (%v)\n", fmin, p.Fref)
		chk.Float64(tst, "fmin", 1e-10, fmin, p.Fref)
		chk.Array(tst, "xmin", 1e-10, x, p.Xref)
	}
}
