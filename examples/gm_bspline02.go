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

	T1 := []float64{0, 0, 0, 1, 1, 1}
	s1 := gm.NewBspline(T1, 2)
	s1.SetControl([][]float64{{0, 0}, {0.5, 1}, {1, 0}})

	T2 := []float64{0, 0, 0, 0.5, 1, 1, 1}
	s2 := gm.NewBspline(T2, 2)
	s2.SetControl([][]float64{{0, 0}, {0.25, 0.5}, {0.75, 0.5}, {1, 0}})

	npts := 201
	plt.Reset(true, &plt.A{Prop: 1.5})
	plt.SplotGap(0.2, 0.4)

	plt.Subplot(3, 2, 1)
	s1.Draw2d(npts, 0, true, nil, nil) // 0 => CalcBasis
	plt.HideAllBorders()

	plt.Subplot(3, 2, 2)
	plt.SetAxis(0, 1, 0, 1)
	s2.Draw2d(npts, 0, true, nil, nil) // 0 => CalcBasis
	plt.HideAllBorders()

	plt.Subplot(3, 2, 3)
	s1.PlotBasis(npts, 0) // 0 => CalcBasis

	plt.Subplot(3, 2, 4)
	s2.PlotBasis(npts, 0) // 0 => CalcBasis

	plt.Subplot(3, 2, 5)
	s1.PlotDerivs(npts)

	plt.Subplot(3, 2, 6)
	s2.PlotDerivs(npts)

	plt.Save("/tmp/gosl", "gm_bspline02")
}
