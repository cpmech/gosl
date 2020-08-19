// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"path/filepath"
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/plt"
)

func Test_npatch01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("npatch01. NURBS patch")

	// NURBS patch
	binsNdiv := 3
	tolerance := 1e-10
	patch := NewNurbsPatch(binsNdiv, tolerance,
		FactoryNurbs.Surf2dRectangleQL(-1, -0.5, 2, 1),
		FactoryNurbs.Surf2dRectangleQL(1, 0.5, 2, 1),
		FactoryNurbs.Surf2dRectangleQL(0.5, 1.7, 2, 0.8),
	)

	// check bins
	io.Pforan("entries = %v\n", patch.Bins.All)
	io.Pf("%v\n", patch.Bins)
	chk.Array(tst, "Xmin", 1e-15, patch.Bins.Xmin, []float64{-1, -0.5})
	chk.Array(tst, "Xmax", 1e-15, patch.Bins.Xmax, []float64{3, 2.5})
	chk.Array(tst, "Xdel", 1e-15, patch.Bins.Xdel, []float64{3 + 1, 2.5 + 0.5})
	chk.Array(tst, "Size", 1e-15, patch.Bins.Size, []float64{4.0 / 3.0, 3.0 / 3.0})
	chk.Ints(tst, "Ndiv", patch.Bins.Ndiv, []int{3, 3})
	chk.Int(tst, "Nall", len(patch.Bins.All), 3*3)
	chk.Int(tst, "Nactive", patch.Bins.Nactive(), 7)
	chk.Int(tst, "Nentries", patch.Bins.Nentries(), 17)

	// check number of entries in bins
	io.Pf("\n")
	entries := map[int][]int{ // maps idx of bin to ids of entries
		0: {0, 1},
		1: {2},
		3: {3, 4},
		4: {5},
		5: {6, 7},
		7: {8, 11, 12, 14, 15},
		8: {9, 10, 13, 16},
	}
	checkBinsEntries(tst, patch.Bins.All, entries)

	// check exchange data
	chk.Int(tst, "0: Id", patch.ExchangeData[0].ID, 0)
	chk.Int(tst, "1: Id", patch.ExchangeData[1].ID, 1)
	chk.Int(tst, "0: Gnd", patch.ExchangeData[0].Gnd, 2)
	chk.Int(tst, "1: Gnd", patch.ExchangeData[1].Gnd, 2)
	chk.Ints(tst, "0: Ords", patch.ExchangeData[0].Ords, []int{2, 1})
	chk.Ints(tst, "1: Ords", patch.ExchangeData[1].Ords, []int{2, 1})
	chk.Deep2(tst, "0: Knots", 1e-15, patch.ExchangeData[0].Knots, [][]float64{{0, 0, 0, 1, 1, 1}, {0, 0, 1, 1}})
	chk.Deep2(tst, "1: Knots", 1e-15, patch.ExchangeData[1].Knots, [][]float64{{0, 0, 0, 1, 1, 1}, {0, 0, 1, 1}})
	chk.Ints(tst, "0: Ctrls", patch.ExchangeData[0].Ctrls, []int{0, 1, 2, 3, 4, 5})
	chk.Ints(tst, "1: Ctrls", patch.ExchangeData[1].Ctrls, []int{5, 6, 7, 8, 9, 10})

	// plot
	if chk.Verbose {
		ndim := 2
		plt.Reset(false, nil)
		for _, surf := range patch.Entities {
			surf.DrawCtrl(ndim, false, nil, nil)
			surf.DrawElems(ndim, 3, false, nil, nil)
		}
		patch.Bins.Draw(true, true, true, true, nil, nil, nil, nil, nil)
		plt.Equal()
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/gm", "t_npatch01")
	}
}

func Test_npatch02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("npatch02. write NURBS patch to json file")

	// NURBS patch
	binsNdiv := 3
	tolerance := 1e-10
	patch := NewNurbsPatch(binsNdiv, tolerance,
		FactoryNurbs.Surf2dRectangleQL(-1, -0.5, 2, 1),
		FactoryNurbs.Surf2dRectangleQL(1, 0.5, 2, 1),
		FactoryNurbs.Surf2dRectangleQL(0.5, 1.7, 2, 0.8),
		FactoryNurbs.Curve2dCircle(0, 1.7, 0.5),
		FactoryNurbs.Curve2dCircle(3, 1.2, 0.5), // will fail if xc==3.0
	)

	// check number of entries in bins
	io.Pf("\n")
	entries := map[int][]int{ // maps idx of bin to ids of entries
		0: {0, 1},
		1: {2},
		3: {3, 4, 21, 22},
		4: {5, 23},
		5: {6, 7, 24, 28, 29, 30, 31},
		6: {18, 19, 20},
		7: {8, 11, 12, 14, 15, 17},
		8: {9, 10, 13, 16, 25, 26, 27},
	}
	checkBinsEntries(tst, patch.Bins.All, entries)

	// check control points
	for i, cp := range patch.ControlPoints {
		if i != cp.ID {
			tst.Errorf("control point id is incorrect")
			return
		}
	}

	// check exchange data
	chk.Ints(tst, "0: Ctrls", patch.ExchangeData[0].Ctrls, []int{0, 1, 2, 3, 4, 5})
	chk.Ints(tst, "1: Ctrls", patch.ExchangeData[1].Ctrls, []int{5, 6, 7, 8, 9, 10})
	chk.Ints(tst, "2: Ctrls", patch.ExchangeData[2].Ctrls, []int{11, 12, 13, 14, 15, 16})
	chk.Ints(tst, "3: Ctrls", patch.ExchangeData[3].Ctrls, []int{11, 17, 18, 19, 20, 21, 22, 23, 11})
	chk.Ints(tst, "4: Ctrls", patch.ExchangeData[4].Ctrls, []int{24, 25, 26, 27, 28, 29, 30, 31, 24})
	//                                                                        ^
	//                                                                        |
	//                                      the fourth id is NOT 13 because the weight is different

	// write file
	patch.Write("/tmp/gosl/gm", "npatch02")

	// read file back
	pp := NewNurbsPatchFromFile(filepath.Join("/tmp/gosl/gm", "npatch02.json"), binsNdiv, tolerance)

	// check read control points
	for i, cp := range pp.ControlPoints {
		if i != cp.ID {
			tst.Errorf("read control point id is incorrect")
			return
		}
	}

	// check read data
	io.Pf("\n")
	chk.Int(tst, "0: Id", pp.ExchangeData[0].ID, 0)
	chk.Int(tst, "1: Id", pp.ExchangeData[1].ID, 1)
	chk.Int(tst, "0: Gnd", pp.ExchangeData[0].Gnd, 2)
	chk.Int(tst, "1: Gnd", pp.ExchangeData[1].Gnd, 2)
	chk.Ints(tst, "0: Ords", pp.ExchangeData[0].Ords, []int{2, 1})
	chk.Ints(tst, "1: Ords", pp.ExchangeData[1].Ords, []int{2, 1})
	chk.Deep2(tst, "0: Knots", 1e-15, pp.ExchangeData[0].Knots, [][]float64{{0, 0, 0, 1, 1, 1}, {0, 0, 1, 1}})
	chk.Deep2(tst, "1: Knots", 1e-15, pp.ExchangeData[1].Knots, [][]float64{{0, 0, 0, 1, 1, 1}, {0, 0, 1, 1}})
	chk.Ints(tst, "0: Ctrls", pp.ExchangeData[0].Ctrls, []int{0, 1, 2, 3, 4, 5})
	chk.Ints(tst, "1: Ctrls", pp.ExchangeData[1].Ctrls, []int{5, 6, 7, 8, 9, 10})
	chk.Ints(tst, "0: Ctrls", pp.ExchangeData[0].Ctrls, []int{0, 1, 2, 3, 4, 5})
	chk.Ints(tst, "1: Ctrls", pp.ExchangeData[1].Ctrls, []int{5, 6, 7, 8, 9, 10})
	chk.Ints(tst, "2: Ctrls", pp.ExchangeData[2].Ctrls, []int{11, 12, 13, 14, 15, 16})
	chk.Ints(tst, "3: Ctrls", pp.ExchangeData[3].Ctrls, []int{11, 17, 18, 19, 20, 21, 22, 23, 11})
	chk.Ints(tst, "4: Ctrls", pp.ExchangeData[4].Ctrls, []int{24, 25, 26, 27, 28, 29, 30, 31, 24})

	// check read bins
	io.Pf("\n")
	checkBinsEntries(tst, pp.Bins.All, entries)

	// check read nurbs
	io.Pf("\n")
	chk.Int(tst, "0: Gnd", pp.Entities[0].Gnd(), 2)
	chk.Int(tst, "1: Gnd", pp.Entities[1].Gnd(), 2)
	chk.Int(tst, "2: Gnd", pp.Entities[2].Gnd(), 2)
	chk.Int(tst, "3: Gnd", pp.Entities[3].Gnd(), 1)
	chk.Int(tst, "4: Gnd", pp.Entities[4].Gnd(), 1)
	chk.Ints(tst, "0: Ord", []int{pp.Entities[0].Ord(0), pp.Entities[0].Ord(1)}, []int{2, 1})
	chk.Ints(tst, "1: Ord", []int{pp.Entities[1].Ord(0), pp.Entities[1].Ord(1)}, []int{2, 1})
	chk.Ints(tst, "2: Ord", []int{pp.Entities[2].Ord(0), pp.Entities[2].Ord(1)}, []int{2, 1})
	chk.Ints(tst, "3: Ord", []int{pp.Entities[3].Ord(0)}, []int{2})
	chk.Ints(tst, "4: Ord", []int{pp.Entities[4].Ord(0)}, []int{2})
	chk.Ints(tst, "0: Nbasis", pp.Entities[0].n, []int{3, 2, 1})
	chk.Ints(tst, "1: Nbasis", pp.Entities[1].n, []int{3, 2, 1})
	chk.Ints(tst, "2: Nbasis", pp.Entities[2].n, []int{3, 2, 1})
	chk.Ints(tst, "3: Nbasis", pp.Entities[3].n, []int{9, 1, 1})
	chk.Ints(tst, "4: Nbasis", pp.Entities[4].n, []int{9, 1, 1})

	// plot
	if chk.Verbose {
		ndim := 2
		plt.Reset(false, nil)

		// original patch
		for _, surf := range patch.Entities {
			surf.DrawCtrl(ndim, false, nil, nil)
			surf.DrawElems(ndim, 11, false, &plt.A{C: "blue", Z: 10}, nil)
		}
		patch.Bins.Draw(true, true, true, true, nil, nil, nil, nil, nil)

		// read patch
		for _, surf := range pp.Entities {
			//surf.DrawCtrl(ndim, false, &plt.A{C: "g"}, nil)
			surf.DrawElems(ndim, 11, false, &plt.A{C: "#f1e09a", Lw: 10, NoClip: true}, nil)
		}
		argsEntry := &plt.A{C: "k", M: "o", Ms: 10, Void: true, NoClip: true}
		pp.Bins.Draw(true, false, false, false, argsEntry, nil, nil, nil, nil)

		// save
		plt.Equal()
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/gm", "t_npatch02")
	}
}
