// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/gm/rw"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func test_rwstep01(tst *testing.T) {

	verbose()
	chk.PrintTitle("rwstep01")

	buf, err := io.ReadFile("rw/data/beadpanel.step")
	if err != nil {
		tst.Errorf("cannot read file:\n%v", err)
		return
	}
	dat := string(buf)

	var stp rw.STEP
	err = stp.ParseDATA(dat)
	if err != nil {
		tst.Errorf("Parse filed:\n%v", err)
		return
	}

	var bsplines []*Bspline

	for _, scurve := range stp.Scurves {
		curve := stp.BScurves[scurve.Curve_3d]
		if curve == nil {
			continue
		}

		// collect vertices
		nv := len(curve.Control_points_list)
		verts := utl.DblsAlloc(nv, 4)
		for i, key := range curve.Control_points_list {
			if p, ok := stp.Points[key]; ok {
				for j := 0; j < 3; j++ {
					verts[i][j] = p.Coordinates[j]
				}
				verts[i][3] = 1.0
			} else {
				chk.Panic("cannot find point %q", key)
			}
		}

		// collect knots
		nk := 0
		for _, m := range curve.Knot_multiplicities {
			nk += m
		}
		knots := make([]float64, nk)
		k := 0
		for i, u := range curve.Knots {
			m := curve.Knot_multiplicities[i]
			for j := 0; j < m; j++ {
				knots[k] = u
				k++
			}
		}

		// create B-spline
		bsp := new(Bspline)
		bsp.Init(knots, curve.Degree)
		bsp.SetControl(verts)
		bsplines = append(bsplines, bsp)
	}

	if true {
		io.Pforan("n = %v\n", len(bsplines))
		for i, bsp := range bsplines {
			bsp.Draw3d("", "", 21, i == 0)
		}
		//plt.Show()
	}
}
