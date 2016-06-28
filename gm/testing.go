// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// CheckDerivs compares analytical with numerical derivatives
func (o *Nurbs) CheckDerivs(tst *testing.T, nn int, tol float64, ver bool) {
	dana := make([]float64, 2)
	dnum := make([]float64, 2)
	for _, u := range utl.LinSpace(o.b[0].tmin, o.b[0].tmax, nn) {
		for _, v := range utl.LinSpace(o.b[1].tmin, o.b[1].tmax, nn) {
			uu := []float64{u, v}
			o.CalcBasisAndDerivs(uu)
			for i := 0; i < o.n[0]; i++ {
				for j := 0; j < o.n[1]; j++ {
					l := i + j*o.n[0]
					o.GetDerivL(dana, l)
					o.NumericalDeriv(dnum, uu, l)
					chk.AnaNum(tst, io.Sf("dR[%d][%d][0](%g,%g)", i, j, uu[0], uu[1]), tol, dana[0], dnum[0], ver)
					chk.AnaNum(tst, io.Sf("dR[%d][%d][1](%g,%g)", i, j, uu[0], uu[1]), tol, dana[1], dnum[1], ver)
				}
			}
		}
	}
}
