// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// data
	x := utl.LinSpace(0.0, 100.0, 11)
	y1 := make([]float64, len(x))
	y2 := make([]float64, len(x))
	y3 := make([]float64, len(x))
	y4 := make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		y1[i] = x[i] * x[i]
		y2[i] = x[i]
		y3[i] = x[i] * 100
		y4[i] = x[i] * 2
	}

	// clear figure
	plt.Reset(false, nil)

	// plot curve on main figure
	plt.Plot(x, y1, &plt.A{L: "curve on old"})

	// plot curve on zoom window
	old, new := plt.ZoomWindow(0.25, 0.5, 0.3, 0.3, nil)
	plt.Plot(x, y2, &plt.A{C: "r", L: "curve on new"})

	// activate main figure
	plt.Sca(old)
	plt.Plot(x, y3, &plt.A{C: "orange", L: "curve ond old again"})
	plt.Gll("x", "y", &plt.A{LegLoc: "lower right"})

	// activate zoom window
	plt.Sca(new)
	plt.Plot(x, y4, &plt.A{C: "cyan", L: "curve ond new again"})
	plt.Gll("xnew", "ynew", nil)

	err := plt.Save("/tmp/gosl", "plt_zoomwindow01")
	if err != nil {
		io.PfRed("save failed:\n%v\n", err)
	}
}
