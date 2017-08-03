// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// Waterfall draws parallel lines @ t along x with height = z. z[len(t)][len(x)]
func Waterfall(X, T []float64, Z [][]float64, argsLine, argsFace *A) {
	if argsLine == nil {
		argsLine = &A{C: "k"}
	}
	if argsFace == nil {
		argsFace = &A{Fc: "w", Ec: "k", Closed: false}
	}
	createAxes3d()
	uid := genUID()
	sx := io.Sf("X%d", uid)
	sz := io.Sf("Z%d", uid)
	genArray(&bufferPy, sx, X)
	genMat(&bufferPy, sz, Z)
	nx := len(X)
	nt := len(T)
	tt := make([]float64, nx)
	P := utl.Alloc(nx, 3)
	xmin, xmax, tmin, tmax, zmin, zmax := X[0], X[0], T[0], T[0], Z[0][0], Z[0][0]
	for i := nt - 1; i >= 0; i-- {
		t := T[i]
		utl.Fill(tt, t)
		uid = genUID()
		st := io.Sf("T%d", uid)
		genArray(&bufferPy, st, tt)
		for j, x := range X {
			P[j][0] = x
			P[j][1] = t
			P[j][2] = Z[i][j]
			zmin = utl.Min(zmin, Z[i][j])
			zmax = utl.Max(zmax, Z[i][j])
		}
		tmin = utl.Min(tmin, t)
		tmax = utl.Max(tmax, t)
		Polygon3d(P, argsFace)
	}
	for _, x := range X {
		xmin = utl.Min(xmin, x)
		xmax = utl.Max(xmax, x)
	}
	AxisRange3d(xmin, xmax, tmin, tmax, zmin, zmax)
}
