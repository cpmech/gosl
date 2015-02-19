// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_grid2D(tst *testing.T) {

	chk.PrintTitle("grid2D")

	var g Grid2D
	//g.Init(6.0, 4.0, 7, 5)
	g.Init(6.0, 4.0, 21, 15)
	if false {
		fxy := func(x, y float64) float64 { return x*x + y*y }
		g.Draw("/tmp/gosl", "fig_t_grid2d_draw", true)
		g.Contour("/tmp/gosl", "fig_t_grid2d_contour", fxy, nil, 11, true)
	}
}
