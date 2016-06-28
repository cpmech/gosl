// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_nurbs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs01")

	// NURBS
	b := FactoryNurbs2dStrip()
	elems := b.Elements()
	nbasis := b.GetElemNumBasis()
	io.Pforan("nbasis = %v\n", nbasis)
	chk.IntAssert(nbasis, 6) // orders := (2,1) => nbasis = (2+1)*(1+1) = 6

	// check basis and elements
	chk.Ints(tst, "elem[0]", elems[0], []int{2, 3, 1, 2})
	chk.Ints(tst, "elem[1]", elems[1], []int{3, 4, 1, 2})
	chk.Ints(tst, "elem[2]", elems[2], []int{4, 5, 1, 2})
	chk.Ints(tst, "ibasis0", b.IndBasis(elems[0]), []int{0, 1, 2, 5, 6, 7})
	chk.Ints(tst, "ibasis1", b.IndBasis(elems[1]), []int{1, 2, 3, 6, 7, 8})
	chk.Ints(tst, "ibasis2", b.IndBasis(elems[2]), []int{2, 3, 4, 7, 8, 9})
	chk.IntAssert(b.GetElemNumBasis(), len(b.IndBasis(elems[0])))

	// check derivatives
	b.CheckDerivs(tst, 11, 1e-5, false)

	// check spans
	solE := [][]int{{2, 3, 1, 2}, {3, 4, 1, 2}, {4, 5, 1, 2}}
	solL := [][]int{{0, 1, 2, 5, 6, 7}, {1, 2, 3, 6, 7, 8}, {2, 3, 4, 7, 8, 9}}
	for k, e := range elems {
		L := b.IndBasis(e)
		io.Pforan("e=%v: L=%v\n", e, L)
		chk.Ints(tst, "span", e, solE[k])
		chk.Ints(tst, "L", L, solL[k])
	}

	// check indices along curve
	io.Pf("\n------------ indices along curve -------------\n")
	chk.Ints(tst, "l0s2a0", b.IndsAlongCurve(0, 2, 0), []int{0, 1, 2})
	chk.Ints(tst, "l0s3a0", b.IndsAlongCurve(0, 3, 0), []int{1, 2, 3})
	chk.Ints(tst, "l0s4a0", b.IndsAlongCurve(0, 4, 0), []int{2, 3, 4})
	chk.Ints(tst, "l0s2a1", b.IndsAlongCurve(0, 2, 1), []int{5, 6, 7})
	chk.Ints(tst, "l0s3a1", b.IndsAlongCurve(0, 3, 1), []int{6, 7, 8})
	chk.Ints(tst, "l0s4a1", b.IndsAlongCurve(0, 4, 1), []int{7, 8, 9})
	chk.Ints(tst, "l1s1a0", b.IndsAlongCurve(1, 1, 0), []int{0, 5})
	chk.Ints(tst, "l1s1a1", b.IndsAlongCurve(1, 1, 1), []int{1, 6})
	chk.Ints(tst, "l1s1a2", b.IndsAlongCurve(1, 1, 2), []int{2, 7})
	chk.Ints(tst, "l1s1a3", b.IndsAlongCurve(1, 1, 3), []int{3, 8})
	chk.Ints(tst, "l1s1a4", b.IndsAlongCurve(1, 1, 4), []int{4, 9})

	// extract surfaces and check
	io.Pf("\n------------ extract surfaces -------------\n")
	surfs := b.ExtractSurfaces()
	chk.Deep4(tst, "surf0: Q", 1e-15, surfs[0].Q, [][][][]float64{
		{{{0, 0, 0, 0.8}}},         // 0
		{{{0, 0.4 * 0.9, 0, 0.9}}}, // 5
	})
	chk.Deep4(tst, "surf1: Q", 1e-15, surfs[1].Q, [][][][]float64{
		{{{1.0 * 1.1, 0.1 * 1.1, 0, 1.1}}}, // 4
		{{{1.0 * 0.5, 0.5 * 0.5, 0, 0.5}}}, // 9
	})
	chk.Deep4(tst, "surf2: Q", 1e-15, surfs[2].Q, [][][][]float64{
		{{{0.00 * 0.80, 0.00 * 0.80, 0, 0.80}}}, // 0
		{{{0.25 * 1.00, 0.15 * 1.00, 0, 1.00}}}, // 1
		{{{0.50 * 0.70, 0.00 * 0.70, 0, 0.70}}}, // 2
		{{{0.75 * 1.20, 0.00 * 1.20, 0, 1.20}}}, // 3
		{{{1.00 * 1.10, 0.10 * 1.10, 0, 1.10}}}, // 4
	})
	chk.Deep4(tst, "surf3: Q", 1e-15, surfs[3].Q, [][][][]float64{
		{{{0.00 * 0.90, 0.40 * 0.90, 0, 0.90}}}, // 5
		{{{0.25 * 0.60, 0.55 * 0.60, 0, 0.60}}}, // 6
		{{{0.50 * 1.50, 0.40 * 1.50, 0, 1.50}}}, // 7
		{{{0.75 * 1.40, 0.40 * 1.40, 0, 1.40}}}, // 8
		{{{1.00 * 0.50, 0.50 * 0.50, 0, 0.50}}}, // 9
	})

	io.Pf("\n------------ elem bry local inds -----------\n")
	elembryinds := b.ElemBryLocalInds()
	io.Pforan("elembryinds = %v\n", elembryinds)
	chk.IntMat(tst, "elembryinds", elembryinds, [][]int{
		{0, 1, 2},
		{2, 5},
		{3, 4, 5},
		{0, 3},
	})

	// refine NURBS
	c := b.Krefine([][]float64{
		{0.5, 1.5, 2.5},
		{0.5},
	})
	elems = c.Elements()
	chk.IntAssert(c.GetElemNumBasis(), len(c.IndBasis(elems[0])))

	// check refined elements: round 1
	io.Pf("\n------------ refined -------------\n")
	chk.Ints(tst, "elem[0]", elems[0], []int{2, 3, 1, 2})
	chk.Ints(tst, "elem[1]", elems[1], []int{3, 4, 1, 2})
	chk.Ints(tst, "elem[2]", elems[2], []int{4, 5, 1, 2})
	chk.Ints(tst, "elem[3]", elems[3], []int{5, 6, 1, 2})
	chk.Ints(tst, "elem[4]", elems[4], []int{6, 7, 1, 2})
	chk.Ints(tst, "elem[5]", elems[5], []int{7, 8, 1, 2})

	// check refined elements: round 2
	chk.Ints(tst, "elem[ 6]", elems[6], []int{2, 3, 2, 3})
	chk.Ints(tst, "elem[ 7]", elems[7], []int{3, 4, 2, 3})
	chk.Ints(tst, "elem[ 8]", elems[8], []int{4, 5, 2, 3})
	chk.Ints(tst, "elem[ 9]", elems[9], []int{5, 6, 2, 3})
	chk.Ints(tst, "elem[10]", elems[10], []int{6, 7, 2, 3})
	chk.Ints(tst, "elem[11]", elems[11], []int{7, 8, 2, 3})

	// check refined basis: round 1
	chk.Ints(tst, "basis0", c.IndBasis(elems[0]), []int{0, 1, 2, 8, 9, 10})
	chk.Ints(tst, "basis1", c.IndBasis(elems[1]), []int{1, 2, 3, 9, 10, 11})
	chk.Ints(tst, "basis2", c.IndBasis(elems[2]), []int{2, 3, 4, 10, 11, 12})
	chk.Ints(tst, "basis3", c.IndBasis(elems[3]), []int{3, 4, 5, 11, 12, 13})
	chk.Ints(tst, "basis4", c.IndBasis(elems[4]), []int{4, 5, 6, 12, 13, 14})
	chk.Ints(tst, "basis5", c.IndBasis(elems[5]), []int{5, 6, 7, 13, 14, 15})

	// check refined basis: round 2
	chk.Ints(tst, "basis6", c.IndBasis(elems[6]), []int{8, 9, 10, 16, 17, 18})
	chk.Ints(tst, "basis7", c.IndBasis(elems[7]), []int{9, 10, 11, 17, 18, 19})
	chk.Ints(tst, "basis8", c.IndBasis(elems[8]), []int{10, 11, 12, 18, 19, 20})
	chk.Ints(tst, "basis9", c.IndBasis(elems[9]), []int{11, 12, 13, 19, 20, 21})
	chk.Ints(tst, "basis10", c.IndBasis(elems[10]), []int{12, 13, 14, 20, 21, 22})
	chk.Ints(tst, "basis11", c.IndBasis(elems[11]), []int{13, 14, 15, 21, 22, 23})

	io.Pf("\n------------ refined: inds along curve -------------\n")
	chk.Ints(tst, "l0s2a0", c.IndsAlongCurve(0, 2, 0), []int{0, 1, 2})
	chk.Ints(tst, "l0s7a0", c.IndsAlongCurve(0, 7, 0), []int{5, 6, 7})
	chk.Ints(tst, "l0s3a2", c.IndsAlongCurve(0, 3, 2), []int{17, 18, 19})
	chk.Ints(tst, "l0s7a2", c.IndsAlongCurve(0, 7, 2), []int{21, 22, 23})
	chk.Ints(tst, "l1s1a0", c.IndsAlongCurve(1, 1, 0), []int{0, 8})
	chk.Ints(tst, "l1s1a0", c.IndsAlongCurve(1, 2, 0), []int{8, 16})
	chk.Ints(tst, "l1s2a7", c.IndsAlongCurve(1, 1, 7), []int{7, 15})
	chk.Ints(tst, "l1s2a7", c.IndsAlongCurve(1, 2, 7), []int{15, 23})

	// plot
	if chk.Verbose {
		io.Pf("\n------------ plot -------------\n")
		PlotNurbs("/tmp/gosl/gm", "nurbs01a.png", b, 41, true, func() {
			colors := []string{"#e78005", "#3944db", "#983cc3", "#529174"}
			for k, surf := range surfs {
				surf.DrawCtrl2d(true, io.Sf(", ls='none', marker='o', ms=4, zorder=20, color='%s'", colors[k]), "va='top', size=7")
				surf.DrawElems2d(21, true, io.Sf(", ls='none', marker='o', ms=4, zorder=20, color='%s'", colors[k]), "va='top', size=7")
			}
		})
		PlotNurbsBasis("/tmp/gosl/gm", "nurbs01b.png", b, 0, 7)
		PlotNurbsDerivs("/tmp/gosl/gm", "nurbs01c.png", b, 0, 7)
		PlotTwoNurbs("/tmp/gosl/gm", "nurbs01d.png", b, c, 41, true, nil)
	}
}

func Test_nurbs02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs02")

	// NURBS
	b := FactoryNurbs2dPlateHole()
	elems := b.Elements()
	nbasis := b.GetElemNumBasis()
	io.Pforan("nbasis = %v\n", nbasis)
	chk.IntAssert(nbasis, 9) // orders := (2,2) => nbasis = (2+1)*(2+1) = 9

	// check basis and elements
	chk.Ints(tst, "elem[0]", elems[0], []int{2, 3, 2, 3})
	chk.Ints(tst, "elem[1]", elems[1], []int{3, 4, 2, 3})
	chk.Ints(tst, "ibasis0", b.IndBasis(elems[0]), []int{0, 1, 2, 4, 5, 6, 8, 9, 10})
	chk.Ints(tst, "ibasis1", b.IndBasis(elems[1]), []int{1, 2, 3, 5, 6, 7, 9, 10, 11})
	chk.IntAssert(b.GetElemNumBasis(), len(b.IndBasis(elems[0])))

	// check derivatives
	b.CheckDerivs(tst, 11, 1e-5, false)

	// refine NURBS
	c := b.KrefineN(2, false)
	elems = c.Elements()
	chk.IntAssert(c.GetElemNumBasis(), len(c.IndBasis(elems[0])))

	// check refined elements
	io.Pf("\n------------ refined -------------\n")
	chk.Ints(tst, "elem[0]", elems[0], []int{2, 3, 2, 3})
	chk.Ints(tst, "elem[1]", elems[1], []int{3, 4, 2, 3})
	chk.Ints(tst, "elem[2]", elems[2], []int{4, 5, 2, 3})
	chk.Ints(tst, "elem[3]", elems[3], []int{5, 6, 2, 3})
	chk.Ints(tst, "elem[4]", elems[4], []int{2, 3, 3, 4})
	chk.Ints(tst, "elem[5]", elems[5], []int{3, 4, 3, 4})
	chk.Ints(tst, "elem[6]", elems[6], []int{4, 5, 3, 4})
	chk.Ints(tst, "elem[7]", elems[7], []int{5, 6, 3, 4})

	// check refined basis
	chk.Ints(tst, "ibasis0", c.IndBasis(elems[0]), []int{0, 1, 2, 6, 7, 8, 12, 13, 14})
	chk.Ints(tst, "ibasis1", c.IndBasis(elems[1]), []int{1, 2, 3, 7, 8, 9, 13, 14, 15})
	chk.Ints(tst, "ibasis2", c.IndBasis(elems[2]), []int{2, 3, 4, 8, 9, 10, 14, 15, 16})
	chk.Ints(tst, "ibasis3", c.IndBasis(elems[3]), []int{3, 4, 5, 9, 10, 11, 15, 16, 17})
	chk.Ints(tst, "ibasis4", c.IndBasis(elems[4]), []int{6, 7, 8, 12, 13, 14, 18, 19, 20})
	chk.Ints(tst, "ibasis5", c.IndBasis(elems[5]), []int{7, 8, 9, 13, 14, 15, 19, 20, 21})
	chk.Ints(tst, "ibasis6", c.IndBasis(elems[6]), []int{8, 9, 10, 14, 15, 16, 20, 21, 22})
	chk.Ints(tst, "ibasis7", c.IndBasis(elems[7]), []int{9, 10, 11, 15, 16, 17, 21, 22, 23})

	// plot
	if chk.Verbose {
		io.Pf("\n------------ plot -------------\n")
		la := 0 + 0*b.n[0]
		lb := 2 + 1*b.n[0]
		PlotNurbs("/tmp/gosl/gm", "nurbs02a.png", b, 41, true, nil)
		PlotNurbsBasis("/tmp/gosl/gm", "nurbs02b.png", b, la, lb)
		PlotNurbsDerivs("/tmp/gosl/gm", "nurbs02c.png", b, la, lb)
		PlotTwoNurbs("/tmp/gosl/gm", "nurbs02d.png", b, c, 41, true, nil)
	}
}

func Test_nurbs03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs03")

	// NURBS
	b := FactoryNurbs1dCurveA()
	elems := b.Elements()
	nbasis := b.GetElemNumBasis()
	io.Pforan("nbasis = %v\n", nbasis)
	chk.IntAssert(nbasis, 4) // orders := (3,) => nbasis = (3+1) = 4

	// check basis and elements
	chk.Ints(tst, "elem[0]", elems[0], []int{3, 4})
	chk.Ints(tst, "elem[1]", elems[1], []int{4, 5})
	chk.Ints(tst, "elem[2]", elems[2], []int{5, 6})
	chk.Ints(tst, "ibasis0", b.IndBasis(elems[0]), []int{0, 1, 2, 3})
	chk.Ints(tst, "ibasis1", b.IndBasis(elems[1]), []int{1, 2, 3, 4})
	chk.Ints(tst, "ibasis2", b.IndBasis(elems[2]), []int{2, 3, 4, 5})

	// refine NURBS
	c := b.Krefine([][]float64{
		{0.15, 0.5, 0.85},
	})

	// plot
	if chk.Verbose {
		PlotNurbs("/tmp/gosl/gm", "nurbs03a.png", b, 41, true, nil)
		PlotTwoNurbs("/tmp/gosl/gm", "nurbs03b.png", b, c, 41, true, nil)
	}
}

func Test_nurbs04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs04")

	// NURBS
	a := FactoryNurbs2dPlateHole()
	b := a.KrefineN(2, false)
	c := a.KrefineN(4, false)

	// tags
	a_vt := tag_verts(a)
	a_ct := map[string]int{
		"0_0": -1,
		"0_1": -2,
	}
	b_vt := tag_verts(b)
	c_vt := tag_verts(c)

	// write .msh files
	WriteMshD("/tmp/gosl/gm", "m_nurbs04a", []*Nurbs{a}, a_vt, a_ct)
	WriteMshD("/tmp/gosl/gm", "m_nurbs04b", []*Nurbs{b}, b_vt, nil)
	WriteMshD("/tmp/gosl/gm", "m_nurbs04c", []*Nurbs{c}, c_vt, nil)

	// read .msh file back and check
	a_read := ReadMsh("/tmp/gosl/gm/m_nurbs04a")[0]
	chk.IntAssert(a_read.gnd, a.gnd)
	chk.Ints(tst, "p", a.p, a_read.p)
	chk.Ints(tst, "n", a.n, a_read.n)
	chk.Deep4(tst, "Q", 1.0e-17, a.Q, a_read.Q)
	chk.IntMat(tst, "l2i", a.l2i, a_read.l2i)

	// plot
	if chk.Verbose {
		PlotNurbs("/tmp/gosl/gm", "nurbs04a.png", a_read, 41, true, nil)
		PlotTwoNurbs("/tmp/gosl/gm", "nurbs04b.png", a, b, 41, true, nil)
		PlotTwoNurbs("/tmp/gosl/gm", "nurbs04c.png", a, c, 41, true, nil)
	}
}
