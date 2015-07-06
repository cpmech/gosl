// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
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

// Check checks derivatives w.r.t to t for fixed coordinates x
func CheckT(tst *testing.T, o Func, t0, tf float64, xcte []float64, np int, tskip []float64, sktol, dtol, dtol2 float64, ver bool) {
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
