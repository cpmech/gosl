// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// CheckDerivT checks derivatives w.r.t to t for fixed coordinates x
func CheckDerivT(tst *testing.T, o T, t0, tf float64, xcte []float64, np int, tskip []float64, sktol, dtol, dtol2 float64, ver bool) {
	t := utl.LinSpace(t0, tf, np)
	for i := 0; i < np; i++ {
		g := o.G(t[i], xcte)
		h := o.H(t[i], xcte)
		skip := false
		for _, val := range tskip {
			if math.Abs(val-t[i]) < sktol {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		chk.DerivScaSca(tst, io.Sf("G(%10f)", t[i]), dtol, g, t[i], 1e-3, ver, func(τ float64) (float64, error) {
			return o.F(τ, xcte), nil
		})
		chk.DerivScaSca(tst, io.Sf("H(%10f)", t[i]), dtol2, h, t[i], 1e-3, ver, func(τ float64) (float64, error) {
			return o.G(τ, xcte), nil
		})
	}
}

// CheckDerivX checks derivatives w.r.t to x for fixed t
func CheckDerivX(tst *testing.T, o T, tcte float64, xmin, xmax []float64, np int, xskip [][]float64, sktol, dtol float64, ver bool) {
	ndim := len(xmin)
	dx := make([]float64, ndim)
	for i := 0; i < ndim; i++ {
		dx[i] = (xmax[i] - xmin[i]) / float64(np-1)
	}
	x := make([]float64, ndim)
	g := make([]float64, ndim)
	nz := 1
	if ndim == 3 {
		nz = np
	}
	for k := 0; k < nz; k++ {
		if ndim == 3 {
			x[2] = xmin[2] + float64(k)*dx[2]
		}
		for j := 0; j < np; j++ {
			x[1] = xmin[1] + float64(j)*dx[1]
			for i := 0; i < np; i++ {
				x[0] = xmin[0] + float64(i)*dx[0]
				skip := false
				for _, val := range xskip {
					for l := 0; l < ndim; l++ {
						if math.Abs(val[l]-x[l]) < sktol {
							skip = true
							break
						}
					}
				}
				if skip {
					continue
				}
				o.Grad(g, tcte, x)
				chk.DerivScaVec(tst, io.Sf("dFdX(t,%10v)", x), dtol, g, x, 1e-3, ver, func(xVec []float64) (float64, error) {
					return o.F(tcte, xVec), nil
				})
			}
		}
	}
}
