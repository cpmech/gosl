// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// surface
	xc, yc, zc, r, R := 0.0, 0.0, 0.0, 2.0, 4.0
	curve := gm.FactoryNurbs.Surf3dTorus(xc, yc, zc, r, R)

	// configurations
	ndim := 3
	nu, nv := 18, 41

	// plot
	plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
	plt.Triad(1.1, "x", "y", "z", &plt.A{C: "orange"}, &plt.A{C: "green"})
	curve.DrawCtrl(ndim, false, &plt.A{C: "grey", Lw: 0.5}, nil)
	curve.DrawSurface(3, nu, nv, true, false, &plt.A{CmapIdx: 3, Rstride: 1, Cstride: 1}, &plt.A{C: "#2782c8", Lw: 0.5})
	plt.Default3dView(-6.1, 6.1, -6.1, 6.1, -6.1, 6.1, true)
	plt.Save("/tmp/gosl", "gm_nurbs03")
}
