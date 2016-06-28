// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tri

import "github.com/cpmech/gosl/plt"

func Draw(V [][]float64, C [][]int, style *plt.Fmt) {
	if style == nil {
		style = &plt.Fmt{C: "b", M: "o", Ms: 2}
	}
	type edgeType struct{ A, B int }
	drawnEdges := make(map[edgeType]bool)
	for _, cell := range C {
		for i := 0; i < 3; i++ {
			a, b := cell[i], cell[(i+1)%3]
			edge := edgeType{a, b}
			if b < a {
				edge.A, edge.B = edge.B, edge.A
			}
			if _, found := drawnEdges[edge]; !found {
				x := []float64{V[a][0], V[b][0]}
				y := []float64{V[a][1], V[b][1]}
				plt.Plot(x, y, style.GetArgs(""))
				drawnEdges[edge] = true
			}
		}
	}
}
