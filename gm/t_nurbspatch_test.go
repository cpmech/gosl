// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"path/filepath"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_npatch01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("npatch01. NURBS patch")

	// NURBS patch
	binsNdiv := 3
	tolerance := 1e-10
	patch, err := NewNurbsPatch(binsNdiv, tolerance,
		FactoryNurbs.Surf2dRectangleQL(-1, -0.5, 2, 1),
		FactoryNurbs.Surf2dRectangleQL(1, 0.5, 2, 1),
		FactoryNurbs.Surf2dRectangleQL(0.5, 1.7, 2, 0.8),
	)
	if err != nil {
		tst.Errorf("%v", err)
		return
	}

	// check bins
	io.Pf("%v\n", patch.Bins)
	chk.Vector(tst, "Xmin", 1e-15, patch.Bins.Xmin, []float64{-1, -0.5})
	chk.Vector(tst, "Xmax", 1e-15, patch.Bins.Xmax, []float64{3, 2.5})
	chk.Vector(tst, "Xdel", 1e-15, patch.Bins.Xdel, []float64{3 + 1, 2.5 + 0.5})
	chk.Vector(tst, "Size", 1e-15, patch.Bins.Size, []float64{4.0 / 3.0, 3.0 / 3.0})
	chk.Ints(tst, "Npts", patch.Bins.Npts, []int{4, 4})
	chk.Int(tst, "Nall", len(patch.Bins.All), 4*4)    // there are ghost bins along each direction
	chk.Int(tst, "Nactive", patch.Bins.Nactive(), 11) // mind the ghost bins
	chk.Int(tst, "Nentries", patch.Bins.Nentries(), 17)

	// check number of entries in bins
	io.Pf("\n")
	entries := map[int][]int{ // maps idx of bin to ids of entries
		0:  []int{0, 1},
		1:  []int{2},
		4:  []int{3, 4},
		5:  []int{5},
		6:  []int{6},
		7:  []int{7},
		9:  []int{8, 11, 12},
		10: []int{9, 13},
		11: []int{10},
		13: []int{14, 15},
		14: []int{16},
	}
	checkBinsEntries(tst, patch.Bins.All, entries)

	// check exchange data
	chk.Int(tst, "0: Id", patch.ExchangeData[0].Id, 0)
	chk.Int(tst, "1: Id", patch.ExchangeData[1].Id, 1)
	chk.Int(tst, "0: Gnd", patch.ExchangeData[0].Gnd, 2)
	chk.Int(tst, "1: Gnd", patch.ExchangeData[1].Gnd, 2)
	chk.Ints(tst, "0: Ords", patch.ExchangeData[0].Ords, []int{2, 1})
	chk.Ints(tst, "1: Ords", patch.ExchangeData[1].Ords, []int{2, 1})
	chk.Matrix(tst, "0: Knots", 1e-15, patch.ExchangeData[0].Knots, [][]float64{{0, 0, 0, 1, 1, 1}, {0, 0, 1, 1}})
	chk.Matrix(tst, "1: Knots", 1e-15, patch.ExchangeData[1].Knots, [][]float64{{0, 0, 0, 1, 1, 1}, {0, 0, 1, 1}})
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

	verbose()
	chk.PrintTitle("npatch02. write NURBS patch to json file")

	// NURBS patch
	binsNdiv := 3
	tolerance := 1e-10
	patch, err := NewNurbsPatch(binsNdiv, tolerance,
		FactoryNurbs.Surf2dRectangleQL(-1, -0.5, 2, 1),
		FactoryNurbs.Surf2dRectangleQL(1, 0.5, 2, 1),
		FactoryNurbs.Surf2dRectangleQL(0.5, 1.7, 2, 0.8),
		FactoryNurbs.Curve2dCircle(0, 1.7, 0.5),
		FactoryNurbs.Curve2dCircle(3, 1.2, 0.5), // will fail if xc==3.0
	)
	if err != nil {
		tst.Errorf("%v", err)
		return
	}

	// check number of entries in bins
	io.Pf("\n")
	entries := map[int][]int{ // maps idx of bin to ids of entries
		0:  []int{0, 1},
		1:  []int{2},
		4:  []int{3, 4, 21, 22},
		5:  []int{5, 23},
		6:  []int{6, 7, 28, 29, 30},
		7:  []int{24, 31},
		8:  []int{18, 19, 20},
		9:  []int{8, 11, 12, 17},
		10: []int{9, 10, 13, 26},
		11: []int{25},
		13: []int{14, 15},
		14: []int{16},
	}
	checkBinsEntries(tst, patch.Bins.All, entries)

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
	pp, err := NewNurbsPatchFromFile(filepath.Join("/tmp/gosl/gm", "npatch02.json"), binsNdiv, tolerance)
	if err != nil {
		tst.Errorf("%v", err)
		return
	}

	// check read dataa
	io.Pf("\n")
	chk.Int(tst, "0: Id", pp.ExchangeData[0].Id, 0)
	chk.Int(tst, "1: Id", pp.ExchangeData[1].Id, 1)
	chk.Int(tst, "0: Gnd", pp.ExchangeData[0].Gnd, 2)
	chk.Int(tst, "1: Gnd", pp.ExchangeData[1].Gnd, 2)
	chk.Ints(tst, "0: Ords", pp.ExchangeData[0].Ords, []int{2, 1})
	chk.Ints(tst, "1: Ords", pp.ExchangeData[1].Ords, []int{2, 1})
	chk.Matrix(tst, "0: Knots", 1e-15, pp.ExchangeData[0].Knots, [][]float64{{0, 0, 0, 1, 1, 1}, {0, 0, 1, 1}})
	chk.Matrix(tst, "1: Knots", 1e-15, pp.ExchangeData[1].Knots, [][]float64{{0, 0, 0, 1, 1, 1}, {0, 0, 1, 1}})
	chk.Ints(tst, "0: Ctrls", pp.ExchangeData[0].Ctrls, []int{0, 1, 2, 3, 4, 5})
	chk.Ints(tst, "1: Ctrls", pp.ExchangeData[1].Ctrls, []int{5, 6, 7, 8, 9, 10})
	chk.Ints(tst, "0: Ctrls", pp.ExchangeData[0].Ctrls, []int{0, 1, 2, 3, 4, 5})
	chk.Ints(tst, "1: Ctrls", pp.ExchangeData[1].Ctrls, []int{5, 6, 7, 8, 9, 10})
	chk.Ints(tst, "2: Ctrls", pp.ExchangeData[2].Ctrls, []int{11, 12, 13, 14, 15, 16})
	chk.Ints(tst, "3: Ctrls", pp.ExchangeData[3].Ctrls, []int{11, 17, 18, 19, 20, 21, 22, 23, 11})
	chk.Ints(tst, "4: Ctrls", pp.ExchangeData[4].Ctrls, []int{24, 25, 26, 27, 28, 29, 30, 31, 24})

	// plot
	if chk.Verbose {
		ndim := 2
		plt.Reset(false, nil)
		for _, surf := range patch.Entities {
			surf.DrawCtrl(ndim, false, nil, nil)
			surf.DrawElems(ndim, 11, false, nil, nil)
		}
		patch.Bins.Draw(true, true, true, true, nil, nil, nil, nil, nil)
		plt.Equal()
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/gm", "t_npatch02")
	}
}
