// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"

	"math/rand"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_bins01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bins01. save and recovery")

	var bins Bins
	bins.Init([]float64{0, 0, 0}, []float64{10, 10, 10}, []int{100, 100, 100})

	// fill bins structure
	maxit := 1000 // number of entries
	X := make([]float64, maxit)
	Y := make([]float64, maxit)
	Z := make([]float64, maxit)
	ID := make([]int, maxit)
	for k := 0; k < maxit; k++ {
		x := rand.Float64() * 10
		y := rand.Float64() * 10
		z := rand.Float64() * 10
		X[k] = x
		Y[k] = y
		Z[k] = z
		ID[k] = k
		bins.Append([]float64{x, y, z}, k, nil)
	}

	// getting ids from bins
	IDchk := make([]int, maxit)
	for k := 0; k < maxit; k++ {
		x := X[k]
		y := Y[k]
		z := Z[k]
		id, sqDist := bins.FindClosest([]float64{x, y, z})
		IDchk[k] = id
		if sqDist > 0 {
			tst.Errorf("sqDist is incorrect: %g", sqDist)
			return
		}
	}
	chk.Ints(tst, "check ids", ID, IDchk)

	// plot
	if chk.Verbose {

		// draw
		plt.Reset(false, nil)
		bins.Draw(true, false, false, false, nil, nil, nil, nil, nil)
		plt.Default3dView(bins.Xmin[0], bins.Xmax[0], bins.Xmin[1], bins.Xmax[1], bins.Xmin[2], bins.Xmax[2], true)
		if false {
			plt.ShowSave("/tmp/gosl/gm", "t_bins01")
		} else {
			plt.Save("/tmp/gosl/gm", "t_bins01")
		}
	}
}

func Test_bins02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bins02. find closest")

	// bins
	var bins Bins
	bins.Init([]float64{0, 0}, []float64{4, 3}, []int{5, 5})

	// append points
	X := []float64{0.5, 1.0, 2.0, 2.0, 2.1, 3.0, 2.1, 2.2}
	Y := []float64{0.0, 0.5, 0.5, 1.0, 2.0, 2.0, 2.1, 2.1}
	for i := 0; i < len(X); i++ {
		bins.Append([]float64{X[i], Y[i]}, i, nil)
	}

	// check
	io.Pf(bins.Summary())
	io.Pf("\n")
	chk.Int(tst, "Ndim", bins.Ndim, 2)
	chk.Array(tst, "Xmin", 1e-15, bins.Xmin, []float64{0, 0})
	chk.Array(tst, "Xmax", 1e-15, bins.Xmax, []float64{4, 3})
	chk.Array(tst, "Xdel", 1e-15, bins.Xdel, []float64{4, 3})
	chk.Array(tst, "Size", 1e-15, bins.Size, []float64{4.0 / 5.0, 3.0 / 5.0})
	chk.Ints(tst, "Npts", bins.Npts, []int{6, 6})
	chk.Int(tst, "Nall", len(bins.All), 6*6) // there are ghost bins along each direction
	chk.Int(tst, "Nactive", bins.Nactive(), 6)
	chk.Int(tst, "Nentries", bins.Nentries(), 8)

	// find closest
	id, sqDist := bins.FindClosest([]float64{2.2, 2.0})
	io.Pforan("\nid=%v  sqDist=%v\n", id, sqDist)
	chk.Int(tst, "closest 4: id", id, 4)
	chk.Float64(tst, "closest 4: sqDist", 1e-15, sqDist, 0.1*0.1)

	// find closest again
	id, sqDist = bins.FindClosest([]float64{2.2, 2.01})
	io.Pforan("\nid=%v  sqDist=%v\n", id, sqDist)
	chk.Int(tst, "closest 7: id", id, 7)
	chk.Float64(tst, "closest 7: sqDist", 1e-15, sqDist, math.Pow(0.1-0.01, 2))

	// append more points
	nextID := bins.Nentries()
	tolerance := 1e-2
	currentID, ex := bins.FindClosestAndAppend(&nextID, []float64{1.0, 1.5}, nil, tolerance, nil)
	io.Pf("\n")
	if ex {
		tst.Errorf("existent flag is incorrect")
		return
	}
	chk.Int(tst, "currentId 8", currentID, 8)
	chk.Int(tst, "nextId 9", nextID, 9)
	chk.Int(tst, "Nactive", bins.Nactive(), 7)
	chk.Int(tst, "Nentries", bins.Nentries(), 9)

	// add point: repeated, no change
	io.Pf("\n")
	currentID, ex = bins.FindClosestAndAppend(&nextID, []float64{1.0, 1.5}, nil, tolerance, nil)
	if !ex {
		tst.Errorf("existent flag is incorrect")
		return
	}
	chk.Int(tst, "currentId 8", currentID, 8)
	chk.Int(tst, "nextId 9", nextID, 9)
	chk.Int(tst, "Nactive", bins.Nactive(), 7)
	chk.Int(tst, "Nentries", bins.Nentries(), 9)

	// add point: very close
	io.Pf("\n")
	tolerance = 0.1
	currentID, ex = bins.FindClosestAndAppend(&nextID, []float64{1.0, 1.59999}, nil, tolerance, nil)
	if !ex {
		tst.Errorf("existent flag is incorrect")
		return
	}
	chk.Int(tst, "currentId 8", currentID, 8)
	chk.Int(tst, "nextId 9", nextID, 9)
	chk.Int(tst, "Nactive", bins.Nactive(), 7)
	chk.Int(tst, "Nentries", bins.Nentries(), 9)

	// add point: new
	io.Pf("\n")
	currentID, ex = bins.FindClosestAndAppend(&nextID, []float64{1.0, 1.6}, nil, tolerance, nil)
	if ex {
		tst.Errorf("existent flag is incorrect")
		return
	}
	chk.Int(tst, "currentId 9", currentID, 9)
	chk.Int(tst, "nextId 10", nextID, 10)
	chk.Int(tst, "Nactive", bins.Nactive(), 7)
	chk.Int(tst, "Nentries", bins.Nentries(), 10)

	// check entries
	io.Pf("\n")
	entries := map[int][]int{0: {0}, 1: {1}, 2: {2}, 8: {3}, 13: {8, 9}, 20: {4, 6, 7}, 21: {5}}
	checkBinsEntries(tst, bins.All, entries)

	// draw
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 500})
		bins.Draw(true, true, true, true, nil, nil, nil, nil, nil)
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/gm", "t_bins02")
	}
}

func Test_bins03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bins03. find along line (2D)")

	// bins
	var bins Bins
	bins.Init([]float64{-0.2, -0.2}, []float64{0.8, 1.8}, []int{5, 5})
	io.Pf(bins.Summary())

	// check
	io.Pf("\n")
	chk.Int(tst, "Ndim", bins.Ndim, 2)
	chk.Array(tst, "Xmin", 1e-15, bins.Xmin, []float64{-0.2, -0.2})
	chk.Array(tst, "Xmax", 1e-15, bins.Xmax, []float64{0.8, 1.8})
	chk.Array(tst, "Xdel", 1e-15, bins.Xdel, []float64{1, 2})
	chk.Array(tst, "Size", 1e-15, bins.Size, []float64{1.0 / 5.0, 2.0 / 5.0})
	chk.Ints(tst, "Npts", bins.Npts, []int{6, 6})
	chk.Int(tst, "Nall", len(bins.All), 6*6) // there are ghost bins along each direction
	chk.Int(tst, "Nactive", bins.Nactive(), 0)
	chk.Int(tst, "Nentries", bins.Nentries(), 0)

	// fill bins structure
	maxit := 5 // number of entries
	ID := make([]int, maxit)
	for k := 0; k < maxit; k++ {
		pos := float64(k) / float64(maxit)
		ID[k] = k
		x := []float64{pos, 2*pos + 0.2}
		if k == 2 {
			io.Pf("\n------------------------------------------------------------------------------\n")
			io.Pf("  With x = %v,\n", x)
			io.Pf("  The expression (x - xmin) / size results in:\n")
			io.Pf("    %v (my: 3.0000000000000004)\n", (x[0]-bins.Xmin[0])/bins.Size[0])
			io.Pf("    %v (my: 2.9999999999999996)\n", (x[1]-bins.Xmin[1])/bins.Size[1])
			io.Pf("  Therefore, bin # 15 will be selected instead of bin # 21.\n")
			io.Pf("\n")
			io.Pf("  This is OK, but other systems may have slightly different rounding errors.\n")
			io.Pf("\n")
			io.Pf("  To avoid results in different systems, x[1] is subtracted by a small value\n")
			io.Pf("  in order to make sure bin # 15 is selected.\n")
			io.Pf("\n")
			δ := 1e-15
			x[1] -= δ
			io.Pf("  The small value is = %v leading to x = %v.\n", δ, x)
			io.Pf("\n")
			io.Pf("  Now, (x - xmin) / size results in:\n")
			io.Pf("    %v (my: 3.0000000000000004)\n", (x[0]-bins.Xmin[0])/bins.Size[0])
			io.Pf("    %v (my: 2.9999999999999973)\n", (x[1]-bins.Xmin[1])/bins.Size[1])
			io.Pf("  which will induce x falling within bin # 15.\n")
			io.Pf("------------------------------------------------------------------------------\n")
		}
		if k == 3 {
			io.Pf("\n------------------------------------------------------------------------------\n")
			io.Pf("  With x = %v,\n", x)
			io.Pf("  The expression (x - xmin) / size results in:\n")
			io.Pf("    %18v (my: 4)\n", (x[0]-bins.Xmin[0])/bins.Size[0])
			io.Pf("    %v (my: 3.9999999999999996)\n", (x[1]-bins.Xmin[1])/bins.Size[1])
			io.Pf("  Therefore, bin # 22 will be selected instead of bin # 28.\n")
			io.Pf("\n")
			δ := 1e-15
			x[1] += δ
			io.Pf("  A small value = %v is now added to x[1] leading:\n", δ)
			io.Pf("  x = %v.\n", x)
			io.Pf("\n")
			io.Pf("  Now, (x - xmin) / size results in:\n")
			io.Pf("    %17v (my: 4)\n", (x[0]-bins.Xmin[0])/bins.Size[0])
			io.Pf("    %v (my: 4.000000000000002)\n", (x[1]-bins.Xmin[1])/bins.Size[1])
			io.Pf("  which will induce x falling within bin # 28.\n")
			io.Pf("------------------------------------------------------------------------------\n")
		}
		bins.Append(x, ID[k], nil)
	}

	// message
	io.Pf("\n")
	for _, bin := range bins.All {
		if bin != nil {
			io.Pf("%v\n", bin)
		}
	}

	// check entries
	io.Pf("\n")
	entries := map[int][]int{7: {0}, 14: {1}, 15: {2}, 28: {3}, 35: {4}}
	checkBinsEntries(tst, bins.All, entries)
	chk.Int(tst, "Nactive", bins.Nactive(), 5)
	chk.Int(tst, "Nentries", bins.Nentries(), 5)

	// add more points to bins
	for i := 0; i < 5; i++ {
		bins.Append([]float64{float64(i) * 0.1, 1.8}, 100+i, nil)
	}

	// find points along diagonal
	io.Pf("\n")
	ids := bins.FindAlongSegment([]float64{0.0, 0.2}, []float64{0.8, 1.8}, 1e-8)
	io.Pf("ids along diagonal = %v\n", ids)
	chk.Ints(tst, "ids along diagonal ", ids, ID)

	// find additional points
	io.Pf("\n")
	ids = bins.FindAlongSegment([]float64{-0.2, 1.8}, []float64{0.8, 1.8}, 1e-8)
	io.Pf("ids along top edge = %v\n", ids)
	chk.Ints(tst, "ids along top edge", ids, []int{100, 101, 102, 103, 104, 4})

	// draw
	if chk.Verbose {
		selBins := map[int]bool{8: true, 9: true, 10: true}
		plt.Reset(true, &plt.A{WidthPt: 500})
		bins.Draw(true, true, true, true, nil, nil, nil, nil, selBins)
		plt.SetXnticks(15)
		plt.SetYnticks(12)
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/gm", "t_bins03")
	}
}

func Test_bins04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bins04. find along line (3D)")

	// bins
	var bins Bins
	bins.Init([]float64{0, 0, 0}, []float64{10, 10, 10}, []int{10, 10, 10})

	// fill bins structure
	maxit := 10 // number of entries
	ID := make([]int, maxit)
	for k := 0; k < maxit; k++ {
		x := float64(k) / float64(maxit) * 10
		ID[k] = k * 11
		bins.Append([]float64{x, x, x}, ID[k], nil)
	}

	// find points along along space diagonal
	ids := bins.FindAlongSegment([]float64{0, 0, 0}, []float64{10, 10, 10}, 0.0000001)
	io.Pforan("ids along space diagonal = %v\n", ids)
	chk.Ints(tst, "ids along space diagonal", ID, ids)

	// draw
	if chk.Verbose {
		argsGrid := &plt.A{C: "#427ce5", Lw: 0.1}
		plt.Reset(true, &plt.A{WidthPt: 500})
		bins.Draw(true, true, false, false, nil, argsGrid, nil, nil, nil)
		plt.DefaultTriad(10.1)
		plt.Default3dView(bins.Xmin[0], bins.Xmax[0], bins.Xmin[1], bins.Xmax[1], bins.Xmin[2], bins.Xmax[2], true)
		plt.Save("/tmp/gosl/gm", "t_bins04")
	}
}

func Test_bins05a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bins05a. find along line (2D)")

	// bins
	var bins Bins
	bins.Init([]float64{0, 0}, []float64{1, 2}, []int{10, 10})

	// add points
	points := [][]float64{
		{0.21132486540518713, 0.21132486540518713},
		{0.7886751345948129, 0.21132486540518713},
		{0.21132486540518713, 0.7886751345948129},
		{0.7886751345948129, 0.7886751345948129},
		{0.21132486540518713, 1.2113248654051871},
		{0.7886751345948129, 1.2113248654051871},
		{0.21132486540518713, 1.788675134594813},
		{0.7886751345948129, 1.788675134594813},
	}
	for i := 0; i < 8; i++ {
		bins.Append(points[i], i, nil)
	}
	io.Pf("bins = %v\n", bins)

	// find points
	x := 0.7886751345948129
	ids := bins.FindAlongSegment([]float64{x, 0}, []float64{x, 2}, 1.e-15)
	io.Pf("ids = %v\n", ids)
	chk.Ints(tst, "ids", []int{1, 3, 5, 7}, ids)

	// draw
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 500})
		bins.Draw(true, true, true, true, nil, nil, nil, nil, nil)
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "t_bins05a")
	}
}

func Test_bins05b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bins05b. find along line (2D). unequal ndiv")

	// bins
	var bins Bins
	bins.Init([]float64{0, 0}, []float64{1, 2}, []int{5, 10})

	// add points
	points := [][]float64{
		{0.21132486540518713, 0.21132486540518713},
		{0.7886751345948129, 0.21132486540518713},
		{0.21132486540518713, 0.7886751345948129},
		{0.7886751345948129, 0.7886751345948129},
		{0.21132486540518713, 1.2113248654051871},
		{0.7886751345948129, 1.2113248654051871},
		{0.21132486540518713, 1.788675134594813},
		{0.7886751345948129, 1.788675134594813},
	}
	for i := 0; i < 8; i++ {
		bins.Append(points[i], i, nil)
	}
	io.Pf("bins = %v\n", bins)

	// find points
	x := 0.7886751345948129
	ids := bins.FindAlongSegment([]float64{x, 0}, []float64{x, 2}, 1.e-15)
	io.Pf("ids = %v\n", ids)
	chk.Ints(tst, "ids", []int{1, 3, 5, 7}, ids)

	// draw
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 500})
		bins.Draw(true, true, true, true, nil, nil, nil, nil, nil)
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "t_bins05b")
	}
}

func Test_bins06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bins06. 3D grid")

	// bins
	var bins Bins
	bins.Init([]float64{5, 5, 5}, []float64{10, 10, 10}, []int{2, 2, 2})
	io.Pfpink(bins.Summary())

	// check
	chk.Int(tst, "Ndim", bins.Ndim, 3)
	chk.Array(tst, "Xmin", 1e-15, bins.Xmin, []float64{5, 5, 5})
	chk.Array(tst, "Xmax", 1e-15, bins.Xmax, []float64{10, 10, 10})
	chk.Array(tst, "Xdel", 1e-15, bins.Xdel, []float64{5, 5, 5})
	chk.Array(tst, "Size", 1e-15, bins.Size, []float64{2.5, 2.5, 2.5})
	chk.Ints(tst, "Npts", bins.Npts, []int{3, 3, 3})
	chk.Int(tst, "Nall", len(bins.All), 27) // there are extra bins along each direction
	chk.Int(tst, "Nactive", bins.Nactive(), 0)
	chk.Int(tst, "Nentries", bins.Nentries(), 0)

	// append
	bins.Append([]float64{9, 7, 6}, 1, nil)
	bins.Append([]float64{8, 5, 6}, 2, nil)
	bins.Append([]float64{7, 7, 5}, 3, nil)
	bins.Append([]float64{5, 7, 6}, 4, nil)
	bins.Append([]float64{5, 5, 5}, 5, nil)
	bins.Append([]float64{10, 10, 10}, 6, nil) // this one goes to a ghost bin
	bins.Append([]float64{5, 5, 10}, 7, nil)   // this one goes to a ghost bin too

	// check again
	chk.Int(tst, "Nactive", bins.Nactive(), 4)
	chk.Int(tst, "Nentries", bins.Nentries(), 7)
	chk.Int(tst, "N0", len(bins.All[0].Entries), 3)
	chk.Int(tst, "N1", len(bins.All[1].Entries), 2)
	for i := 2; i < 25; i++ {
		if i == 18 { // this contain one entry
			continue
		}
		if bins.All[i] != nil {
			tst.Errorf("bin # %d should be empty", i)
			return
		}
	}
	chk.Int(tst, "N26", len(bins.All[26].Entries), 1)

	// plot
	if chk.Verbose {
		plt.Reset(false, nil)
		bins.Draw(true, true, true, true, nil, nil, nil, nil, nil)
		plt.Default3dView(bins.Xmin[0], bins.Xmax[0], bins.Xmin[1], bins.Xmax[1], bins.Xmin[2], bins.Xmax[2], true)
		if false {
			plt.ShowSave("/tmp/gosl/gm", "t_bins06")
		} else {
			plt.Save("/tmp/gosl/gm", "t_bins06")
		}
	}
}

// auxiliary /////////////////////////////////////////////////////////////////////////////////////

// entries is a map with the ids of each entry in each bin: maps binId => entries ids
func checkBinsEntries(tst *testing.T, bins []*Bin, entries map[int][]int) {
	for idx, bin := range bins {
		txt := io.Sf("N%d", idx)
		if e, ok := entries[idx]; ok {
			if bin == nil {
				tst.Errorf("bin " + txt + " should not be nil\n")
				return
			}
			chk.Int(tst, txt, len(bin.Entries), len(e))
			ee := make([]int, len(bin.Entries))
			for k, entry := range bin.Entries {
				ee[k] = entry.ID
			}
			chk.Ints(tst, txt, ee, e)
		} else {
			if bin != nil {
				tst.Errorf("bin " + txt + " should be nil\n")
				return
			}
		}
	}
}
