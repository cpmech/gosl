// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func doplot(maxN int) {
	l := 21
	N := utl.LinSpace(1, float64(maxN), l)
	Y1 := make([]float64, l)
	Y2 := make([]float64, l)
	Y3 := make([]float64, l)
	Y4 := make([]float64, l)
	Y5 := make([]float64, l)
	Y6 := make([]float64, l)
	for i, n := range N {
		x := float64(n)
		Y1[i] = 1
		Y2[i] = x
		Y3[i] = math.Log2(x)
		Y4[i] = x * math.Log2(x)
		Y5[i] = math.Pow(x, 2)
		Y6[i] = math.Pow(2, x)
	}
	plt.Plot(N, Y1, &plt.A{NoClip: true, M: plt.M(0, 0), Ms: 5, L: "$O(1)$"})
	plt.Plot(N, Y2, &plt.A{NoClip: true, M: plt.M(1, 0), Ms: 5, L: "$O(n)$"})
	plt.Plot(N, Y3, &plt.A{NoClip: true, M: plt.M(2, 0), Ms: 5, L: "$O(\\log{n})$"})
	plt.Plot(N, Y4, &plt.A{NoClip: true, M: plt.M(3, 0), Ms: 5, L: "$O(n \\cdot\\log{n})$"})
	plt.Plot(N, Y5, &plt.A{NoClip: true, M: plt.M(4, 0), Ms: 5, L: "$O(n^2)$"})
	plt.Plot(N, Y6, &plt.A{NoClip: true, M: plt.M(5, 0), Ms: 3, L: "$O(2^n)$"})
	plt.SetNoXtickLabels()
	plt.SetNoYtickLabels()
	plt.HideTRborders()
}

func main() {
	plt.Reset(true, nil)
	doplot(5)
	plt.Title("Big-O complexity", nil)
	plt.Gll("elements", "operations", &plt.A{LegNcol: 3})
	plt.ZoomWindow(0.15, 0.38, 0.4, 0.4, nil)
	doplot(10)
	plt.Save("/tmp/gosl", "plt_bigo")
}
