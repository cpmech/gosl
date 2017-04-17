// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/utl"
)

// CheckDerivT checks derivatives w.r.t to t for fixed coordinates x
func CheckDerivT(tst *testing.T, o TimeSpace, t0, tf float64, xcte []float64, np int, tskip []float64, sktol, dtol, dtol2 float64, ver bool) {
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
		dnum := num.DerivCen(func(t float64, args ...interface{}) (res float64) {
			return o.F(t, xcte)
		}, t[i])
		chk.AnaNum(tst, io.Sf("G(%10f)", t[i]), dtol, g, dnum, ver)
		dnum2 := num.DerivCen(func(t float64, args ...interface{}) (res float64) {
			return o.G(t, xcte)
		}, t[i])
		chk.AnaNum(tst, io.Sf("H(%10f)", t[i]), dtol2, h, dnum2, ver)
	}
}

// CheckDerivX checks derivatives w.r.t to x for fixed t
func CheckDerivX(tst *testing.T, o TimeSpace, tcte float64, xmin, xmax []float64, np int, xskip [][]float64, sktol, dtol float64, ver bool) {
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
	xtmp := make([]float64, ndim)
	for k := 0; k < nz; k++ {
		if ndim == 3 {
			x[2] = xmin[2] + float64(k)*dx[2]
		}
		for j := 0; j < np; j++ {
			x[1] = xmin[1] + float64(j)*dx[1]
			for i := 0; i < np; i++ {
				x[0] = xmin[0] + float64(i)*dx[0]
				o.Grad(g, tcte, x)
				for l := 0; l < ndim; l++ {
					skip := false
					for _, val := range xskip {
						if math.Abs(val[l]-x[l]) < sktol {
							skip = true
							break
						}
					}
					if skip {
						continue
					}
					dnum := num.DerivCen(func(s float64, args ...interface{}) (res float64) {
						copy(xtmp, x)
						xtmp[l] = s
						return o.F(tcte, xtmp)
					}, x[l])
					chk.AnaNum(tst, io.Sf("dFdX(t,%10v)[%d]", x, l), dtol, g[l], dnum, ver)
				}
			}
		}
	}
}
