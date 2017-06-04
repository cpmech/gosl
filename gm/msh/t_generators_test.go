// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
)

func TestGen01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Gen01. 2D ring")

	r, R := 1.0, 3.0
	mesh, err := GenRing2d(4, 5, r, R, math.Pi/4.0)
	if err != nil {
		tst.Errorf("%v", err)
		return
	}

	tm, err := mesh.GetTagMaps()
	if err != nil {
		tst.Errorf("%v", err)
		return
	}

	vset := tm.VertTag2verts[4]
	vset = append(vset, tm.VertTag2verts[41]...)
	vset = append(vset, tm.VertTag2verts[43]...)
	chk.Int(tst, "len(vset)", len(vset), 11)

	for _, v := range vset {
		rm := math.Sqrt(v.X[0]*v.X[0] + v.X[1]*v.X[1])
		chk.Scalar(tst, "r", 1e-15, r, rm)
	}

	vset = tm.VertTag2verts[2]
	vset = append(vset, tm.VertTag2verts[21]...)
	vset = append(vset, tm.VertTag2verts[23]...)
	chk.Int(tst, "len(vset)", len(vset), 11)

	for _, v := range vset {
		Rm := math.Sqrt(v.X[0]*v.X[0] + v.X[1]*v.X[1])
		chk.Scalar(tst, "R", 1e-15, R, Rm)
	}

	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		args := NewArgs()
		args.WithEdges = true
		args.WithCells = true
		args.WithVerts = true
		args.WithIdsCells = true
		args.WithIdsVerts = true
		args.WithTagsVerts = true
		args.WithTagsEdges = true
		mesh.Draw(args)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/gm", "ring2d")
	}
}
