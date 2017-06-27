// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/la"
)

// Jacobian computes Jacobian (sparse) matrix
//      Calculates (with N=n-1):
//          df0dx0, df0dx1, df0dx2, ... df0dxN
//          df1dx0, df1dx1, df1dx2, ... df1dxN
//               . . . . . . . . . . . . .
//          dfNdx0, dfNdx1, dfNdx2, ... dfNdxN
//  INPUT:
//      ffcn : f(x) function
//      x    : station where dfdx has to be calculated
//      fx   : f @ x
//      w    : workspace with size == n == len(x)
//  RETURNS:
//      J : dfdx @ x [must be pre-allocated]
func Jacobian(J *la.Triplet, ffcn fun.Vv, x, fx, w []float64) (err error) {
	ndim := len(x)
	start, endp1 := 0, ndim
	if J.Max() == 0 {
		J.Init(ndim, ndim, ndim*ndim)
	}
	J.Start()
	var df float64
	for col := 0; col < ndim; col++ {
		xsafe := x[col]
		delta := math.Sqrt(MACHEPS * max(1e-5, math.Abs(xsafe)))
		x[col] = xsafe + delta
		err = ffcn(w, x) // w := f(x+δx[col])
		if err != nil {
			return
		}
		for row := start; row < endp1; row++ {
			df = w[row] - fx[row]
			J.Put(row, col, df/delta)
		}
		x[col] = xsafe
	}
	return
}

// CompareJac compares Jacobian matrix (e.g. for testing)
func CompareJac(tst *testing.T, ffcn fun.Vv, Jfcn fun.Tv, x []float64, tol float64) {

	// numerical
	n := len(x)
	fx := make([]float64, n)
	w := make([]float64, n) // workspace
	ffcn(fx, x)
	var Jnum la.Triplet
	Jnum.Init(n, n, n*n)
	Jacobian(&Jnum, ffcn, x, fx, w)

	// analytical
	var Jana la.Triplet
	Jana.Init(n, n, n*n)
	Jfcn(&Jana, x)

	// compare
	chk.Vector(tst, "Jacobian matrix", tol, Jnum.GetDenseMatrix().Data, Jana.GetDenseMatrix().Data)
}
