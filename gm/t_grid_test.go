// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func TestGrid01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid01")

	g, err := NewUniformGrid([]float64{-6, -3}, []float64{6, 3}, []int{4, 3})
	status(tst, err)

	chk.Int(tst, "N", g.N, 20)
	chk.Int(tst, "nx", g.Npts[0], 5)
	chk.Int(tst, "ny", g.Npts[1], 4)

	chk.Float64(tst, "Lx", 1e-15, g.Xdel[0], 12.0)
	chk.Float64(tst, "Ly", 1e-15, g.Xdel[1], 6.0)
	chk.Float64(tst, "Dx", 1e-15, g.Size[0], 3.0)
	chk.Float64(tst, "Dy", 1e-15, g.Size[1], 2.0)

	chk.Ints(tst, "B", g.Edge[0], []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "R", g.Edge[1], []int{4, 9, 14, 19})
	chk.Ints(tst, "T", g.Edge[2], []int{15, 16, 17, 18, 19})
	chk.Ints(tst, "L", g.Edge[3], []int{0, 5, 10, 15})

	chk.Ints(tst, "Tag # 10: L", g.GetNodesWithTag(10), []int{0, 5, 10, 15})
	chk.Ints(tst, "Tag # 11: R", g.GetNodesWithTag(11), []int{4, 9, 14, 19})
	chk.Ints(tst, "Tag # 20: B", g.GetNodesWithTag(20), []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "Tag # 21: T", g.GetNodesWithTag(21), []int{15, 16, 17, 18, 19})

	// plot
	if chk.Verbose {
		X, Y, F, err := g.Eval2d(func(x la.Vector) (float64, error) {
			return x[0]*x[0] + x[1]*x[1], nil
		})
		status(tst, err)
		plt.Reset(true, &plt.A{WidthPt: 500})
		g.Draw(true, true, true, true, nil, nil, nil, nil, nil)
		plt.Grid(&plt.A{C: "grey"})
		plt.ContourL(X, Y, F, nil)
		plt.Equal()
		plt.HideAllBorders()
		plt.SetXnticks(12)
		plt.SetYnticks(12)
		err = plt.Save("/tmp/gosl/gm", "grid01")
		status(tst, err)
	}
}
