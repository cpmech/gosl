// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func get_nurbs_A() (b *Nurbs) {
	verts := [][]float64{
		{0.00, 0.00, 0, 0.80}, // 0
		{0.25, 0.15, 0, 1.00}, // 1
		{0.50, 0.00, 0, 0.70}, // 2
		{0.75, 0.00, 0, 1.20}, // 3
		{1.00, 0.10, 0, 1.10}, // 4
		{0.00, 0.40, 0, 0.90}, // 5
		{0.25, 0.55, 0, 0.60}, // 6
		{0.50, 0.40, 0, 1.50}, // 7
		{0.75, 0.40, 0, 1.40}, // 8
		{1.00, 0.50, 0, 0.50}, // 9
	}
	knots := [][]float64{
		{0, 0, 0, 1, 2, 3, 3, 3},
		{0, 0, 1, 1},
	}
	b = new(Nurbs)
	b.Init(2, []int{2, 1}, knots)
	b.SetControl(verts, utl.IntRange(len(verts)))
	return
}

func get_nurbs_B() (b *Nurbs) {
	verts := [][]float64{
		{-1.000000000000000e+00, 0.000000000000000e+00, 0, 1.000000000000000e+00},
		{-1.000000000000000e+00, 4.142135623730951e-01, 0, 8.535533905932737e-01},
		{-4.142135623730951e-01, 1.000000000000000e+00, 0, 8.535533905932737e-01},
		{0.000000000000000e+00, 1.000000000000000e+00, 0, 1.000000000000000e+00},
		{-2.500000000000000e+00, 0.000000000000000e+00, 0, 1.000000000000000e+00},
		{-2.500000000000000e+00, 7.500000000000000e-01, 0, 1.000000000000000e+00},
		{-7.500000000000000e-01, 2.500000000000000e+00, 0, 1.000000000000000e+00},
		{0.000000000000000e+00, 2.500000000000000e+00, 0, 1.000000000000000e+00},
		{-4.000000000000000e+00, 0.000000000000000e+00, 0, 1.000000000000000e+00},
		{-4.000000000000000e+00, 4.000000000000000e+00, 0, 1.000000000000000e+00},
		{0.000000000000000e+00, 4.000000000000000e+00, 0, 1.000000000000000e+00},
	}
	knots := [][]float64{
		{0, 0, 0, 0.5, 1, 1, 1},
		{0, 0, 0, 1, 1, 1},
	}
	b = new(Nurbs)
	b.Init(2, []int{2, 2}, knots)
	b.SetControl(verts, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 10})
	return
}

func get_nurbs_C() (b *Nurbs) {
	verts := [][]float64{
		{0.0, 0.0, 0, 1},
		{1.0, 0.2, 0, 1},
		{0.5, 1.5, 0, 1},
		{2.5, 2.0, 0, 1},
		{2.0, 0.4, 0, 1},
		{3.0, 0.0, 0, 1},
	}
	knots := [][]float64{
		{0, 0, 0, 0, 0.3, 0.7, 1, 1, 1, 1},
	}
	b = new(Nurbs)
	b.Init(1, []int{3}, knots)
	b.SetControl(verts, utl.IntRange(len(verts)))
	return
}

func do_check_derivs(tst *testing.T, b *Nurbs, nn int, tol float64, ver bool) {
	dana := make([]float64, 2)
	dnum := make([]float64, 2)
	for _, u := range utl.LinSpace(b.b[0].tmin, b.b[0].tmax, nn) {
		for _, v := range utl.LinSpace(b.b[1].tmin, b.b[1].tmax, nn) {
			uu := []float64{u, v}
			b.CalcBasisAndDerivs(uu)
			for i := 0; i < b.n[0]; i++ {
				for j := 0; j < b.n[1]; j++ {
					l := i + j*b.n[0]
					b.GetDerivL(dana, l)
					b.NumericalDeriv(dnum, uu, l)
					chk.AnaNum(tst, io.Sf("dR[%d][%d][0](%g,%g)", i, j, uu[0], uu[1]), tol, dana[0], dnum[0], ver)
					chk.AnaNum(tst, io.Sf("dR[%d][%d][1](%g,%g)", i, j, uu[0], uu[1]), tol, dana[1], dnum[1], ver)
				}
			}
		}
	}
}

func Test_nurbs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs01")

	b := get_nurbs_A()
	elems := b.Elements()
	nbasis := b.GetElemNumBasis()
	io.Pforan("nbasis = %v\n", nbasis)
	chk.IntAssert(nbasis, 6) // orders := (2,1) => nbasis = (2+1)*(1+1)
	chk.Ints(tst, "elem[0]", elems[0], []int{2, 3, 1, 2})
	chk.Ints(tst, "elem[1]", elems[1], []int{3, 4, 1, 2})
	chk.Ints(tst, "elem[2]", elems[2], []int{4, 5, 1, 2})
	chk.Ints(tst, "ibasis0", b.IndBasis(elems[0]), []int{0, 1, 2, 5, 6, 7})
	chk.Ints(tst, "ibasis1", b.IndBasis(elems[1]), []int{1, 2, 3, 6, 7, 8})
	chk.Ints(tst, "ibasis2", b.IndBasis(elems[2]), []int{2, 3, 4, 7, 8, 9})
	chk.IntAssert(b.GetElemNumBasis(), len(b.IndBasis(elems[0])))

	if chk.Verbose {
		PlotNurbsBasis("/tmp/gosl", "t_nurbs01", b, 0, 7)
	}
}

func Test_nurbs02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs02")

	b := get_nurbs_A()
	do_check_derivs(tst, b, 11, 1e-5, false)

	if chk.Verbose {
		PlotNurbsDerivs("/tmp/gosl", "t_nurbs02", b, 0, 7)
	}
}

func Test_nurbs03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs03")

	b := get_nurbs_B()
	nbasis := b.GetElemNumBasis()
	io.Pforan("nbasis = %v\n", nbasis)
	chk.IntAssert(nbasis, 9) // orders := (2,1) => nbasis = (2+1)*(2+1)

	if chk.Verbose {
		la := 0 + 0*b.n[0]
		lb := 2 + 1*b.n[0]
		PlotNurbsBasis("/tmp/gosl", "t_nurbs03", b, la, lb)
	}
}

func Test_nurbs04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs04")

	b := get_nurbs_B()
	do_check_derivs(tst, b, 11, 1e-5, false)

	if chk.Verbose {
		la := 0 + 0*b.n[0]
		lb := 2 + 1*b.n[0]
		PlotNurbsDerivs("/tmp/gosl", "t_nurbs04", b, la, lb)
	}
}

func Test_nurbs05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs05")

	b := get_nurbs_A()
	elems := b.Elements()
	solE := [][]int{{2, 3, 1, 2}, {3, 4, 1, 2}, {4, 5, 1, 2}}
	solL := [][]int{{0, 1, 2, 5, 6, 7}, {1, 2, 3, 6, 7, 8}, {2, 3, 4, 7, 8, 9}}
	for k, e := range elems {
		L := b.IndBasis(e)
		io.Pforan("e=%v: L=%v\n", e, L)
		chk.Ints(tst, "span", e, solE[k])
		chk.Ints(tst, "L", L, solL[k])
	}

	if chk.Verbose {
		PlotNurbs("/tmp/gosl", "t_nurbs05", b)
	}
}

func Test_nurbs06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs06")

	b := get_nurbs_C()
	elems := b.Elements()
	chk.Ints(tst, "elem[0]", elems[0], []int{3, 4})
	chk.Ints(tst, "elem[1]", elems[1], []int{4, 5})
	chk.Ints(tst, "elem[2]", elems[2], []int{5, 6})
	chk.Ints(tst, "ibasis0", b.IndBasis(elems[0]), []int{0, 1, 2, 3})
	chk.Ints(tst, "ibasis1", b.IndBasis(elems[1]), []int{1, 2, 3, 4})
	chk.Ints(tst, "ibasis2", b.IndBasis(elems[2]), []int{2, 3, 4, 5})

	if chk.Verbose {
		c := b.Krefine([][]float64{
			//{0.15},
			{0.15, 0.5, 0.85},
		})
		PlotNurbsRefined("/tmp/gosl", "t_nurbs06", b, c)
	}
}

func Test_nurbs07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs07")

	b := get_nurbs_A()
	c := b.Krefine([][]float64{
		{0.5, 1.5, 2.5},
		{0.5},
	})

	elems := c.Elements()

	chk.Ints(tst, "elem[0]", elems[0], []int{2, 3, 1, 2})
	chk.Ints(tst, "elem[1]", elems[1], []int{3, 4, 1, 2})
	chk.Ints(tst, "elem[2]", elems[2], []int{4, 5, 1, 2})
	chk.Ints(tst, "elem[3]", elems[3], []int{5, 6, 1, 2})
	chk.Ints(tst, "elem[4]", elems[4], []int{6, 7, 1, 2})
	chk.Ints(tst, "elem[5]", elems[5], []int{7, 8, 1, 2})

	chk.Ints(tst, "elem[ 6]", elems[6], []int{2, 3, 2, 3})
	chk.Ints(tst, "elem[ 7]", elems[7], []int{3, 4, 2, 3})
	chk.Ints(tst, "elem[ 8]", elems[8], []int{4, 5, 2, 3})
	chk.Ints(tst, "elem[ 9]", elems[9], []int{5, 6, 2, 3})
	chk.Ints(tst, "elem[10]", elems[10], []int{6, 7, 2, 3})
	chk.Ints(tst, "elem[11]", elems[11], []int{7, 8, 2, 3})

	chk.Ints(tst, "basis0", c.IndBasis(elems[0]), []int{0, 1, 2, 8, 9, 10})
	chk.Ints(tst, "basis1", c.IndBasis(elems[1]), []int{1, 2, 3, 9, 10, 11})
	chk.Ints(tst, "basis2", c.IndBasis(elems[2]), []int{2, 3, 4, 10, 11, 12})
	chk.Ints(tst, "basis3", c.IndBasis(elems[3]), []int{3, 4, 5, 11, 12, 13})
	chk.Ints(tst, "basis4", c.IndBasis(elems[4]), []int{4, 5, 6, 12, 13, 14})
	chk.Ints(tst, "basis5", c.IndBasis(elems[5]), []int{5, 6, 7, 13, 14, 15})

	chk.Ints(tst, "basis6", c.IndBasis(elems[6]), []int{8, 9, 10, 16, 17, 18})
	chk.Ints(tst, "basis7", c.IndBasis(elems[7]), []int{9, 10, 11, 17, 18, 19})
	chk.Ints(tst, "basis8", c.IndBasis(elems[8]), []int{10, 11, 12, 18, 19, 20})
	chk.Ints(tst, "basis9", c.IndBasis(elems[9]), []int{11, 12, 13, 19, 20, 21})
	chk.Ints(tst, "basis10", c.IndBasis(elems[10]), []int{12, 13, 14, 20, 21, 22})
	chk.Ints(tst, "basis11", c.IndBasis(elems[11]), []int{13, 14, 15, 21, 22, 23})

	if chk.Verbose {
		PlotNurbsRefined("/tmp/gosl", "t_nurbs07", b, c)
	}
}

func tag_verts(b *Nurbs) (vt map[int]int) {
	vt = make(map[int]int)
	n0, n1 := b.NumBasis(0), b.NumBasis(1)
	for j := 0; j < n1; j++ {
		for i := 0; i < n0; i++ {
			x := b.GetQ(i, j, 0)
			if math.Abs(x[0]) < 1e-7 { // right
				vt[HashPoint(x[0], x[1], x[2])] = -1
			}
			if math.Abs(x[1]) < 1e-7 { // bottom
				vt[HashPoint(x[0], x[1], x[2])] = -2
			}
			if math.Abs(x[0]+4.0) < 1e-7 { // left
				vt[HashPoint(x[0], x[1], x[2])] = -3
			}
			if math.Abs(x[0]+4.0) < 1e-7 && math.Abs(x[1]) < 1e-7 { // left-bottom
				vt[HashPoint(x[0], x[1], x[2])] = -4
			}
		}
	}
	return
}

func Test_nurbs08(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs08")

	a := get_nurbs_B()
	a_el := a.Elements()
	chk.Ints(tst, "a_el[0]", a_el[0], []int{2, 3, 2, 3})
	chk.Ints(tst, "a_el[1]", a_el[1], []int{3, 4, 2, 3})
	chk.Ints(tst, "a_ibasis0", a.IndBasis(a_el[0]), []int{0, 1, 2, 4, 5, 6, 8, 9, 10})
	chk.Ints(tst, "a_ibasis1", a.IndBasis(a_el[1]), []int{1, 2, 3, 5, 6, 7, 9, 10, 11})
	chk.IntAssert(a.GetElemNumBasis(), len(a.IndBasis(a_el[0])))

	b := a.KrefineN(2, false)
	b_el := b.Elements()
	chk.Ints(tst, "b_el[0]", b_el[0], []int{2, 3, 2, 3})
	chk.Ints(tst, "b_el[1]", b_el[1], []int{3, 4, 2, 3})
	chk.Ints(tst, "b_el[2]", b_el[2], []int{4, 5, 2, 3})
	chk.Ints(tst, "b_el[3]", b_el[3], []int{5, 6, 2, 3})
	chk.Ints(tst, "b_el[4]", b_el[4], []int{2, 3, 3, 4})
	chk.Ints(tst, "b_el[5]", b_el[5], []int{3, 4, 3, 4})
	chk.Ints(tst, "b_el[6]", b_el[6], []int{4, 5, 3, 4})
	chk.Ints(tst, "b_el[7]", b_el[7], []int{5, 6, 3, 4})
	chk.Ints(tst, "b_ibasis0", b.IndBasis(b_el[0]), []int{0, 1, 2, 6, 7, 8, 12, 13, 14})
	chk.Ints(tst, "b_ibasis1", b.IndBasis(b_el[1]), []int{1, 2, 3, 7, 8, 9, 13, 14, 15})
	chk.Ints(tst, "b_ibasis2", b.IndBasis(b_el[2]), []int{2, 3, 4, 8, 9, 10, 14, 15, 16})
	chk.Ints(tst, "b_ibasis3", b.IndBasis(b_el[3]), []int{3, 4, 5, 9, 10, 11, 15, 16, 17})
	chk.Ints(tst, "b_ibasis4", b.IndBasis(b_el[4]), []int{6, 7, 8, 12, 13, 14, 18, 19, 20})
	chk.Ints(tst, "b_ibasis5", b.IndBasis(b_el[5]), []int{7, 8, 9, 13, 14, 15, 19, 20, 21})
	chk.Ints(tst, "b_ibasis6", b.IndBasis(b_el[6]), []int{8, 9, 10, 14, 15, 16, 20, 21, 22})
	chk.Ints(tst, "b_ibasis7", b.IndBasis(b_el[7]), []int{9, 10, 11, 15, 16, 17, 21, 22, 23})
	chk.IntAssert(b.GetElemNumBasis(), len(b.IndBasis(b_el[0])))

	c := a.KrefineN(4, false)

	a_vt := tag_verts(a)
	a_ct := map[string]int{
		"0_0": -1,
		"0_1": -2,
	}
	b_vt := tag_verts(b)
	c_vt := tag_verts(c)

	WriteMshD("/tmp/gosl", "m_nurbs08a", []*Nurbs{a}, a_vt, a_ct)
	WriteMshD("/tmp/gosl", "m_nurbs08b", []*Nurbs{b}, b_vt, nil)
	WriteMshD("/tmp/gosl", "m_nurbs08c", []*Nurbs{c}, c_vt, nil)

	B := ReadMsh("/tmp/gosl/m_nurbs08a")

	if a.gnd != B[0].gnd {
		chk.Panic("Read: gnd is wrong")
	}
	chk.Ints(tst, "Read: p", a.p, B[0].p)
	chk.Ints(tst, "Read: n", a.n, B[0].n)
	chk.Deep4(tst, "Read: Q", 1.0e-17, a.Q, B[0].Q)
	chk.IntMat(tst, "Read: l2i", a.l2i, B[0].l2i)

	if chk.Verbose {
		PlotNurbs("/tmp/gosl", "t_nurbs08_read", B[0])
		PlotNurbsRefined("/tmp/gosl", "t_nurbs08_ref2", a, b)
		PlotNurbsRefined("/tmp/gosl", "t_nurbs08_ref4", a, c)
	}
}

func Test_nurbs09(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs09")

	b := get_nurbs_A()
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

	if chk.Verbose {
		I := [][]int{{2, 2, 1}, {2, 2, 2}, {2, 2, 3}, {2, 2, 4}}
		plt.SetForEps(1, 500)
		for k, surf := range surfs {
			plt.SubplotI(I[k])
			b.DrawElems2D(41, false, ", ls=':', lw=0.5", "")
			surf.DrawCtrl2D(true)
			surf.DrawElems2D(41, true, ", zorder=20, color='#e78005', lw=2", "")
			plt.Equal()
		}
		plt.SaveD("/tmp/gosl", "test_nurbs09.eps")
	}
}

func Test_nurbs10(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nurbs10")

	if chk.Verbose {
		b := get_nurbs_B()
		surfs := b.ExtractSurfaces()
		I := [][]int{{2, 2, 1}, {2, 2, 2}, {2, 2, 3}, {2, 2, 4}}
		plt.SetForEps(1, 500)
		for k, surf := range surfs {
			plt.SubplotI(I[k])
			b.DrawElems2D(41, false, ", ls=':', lw=0.5", "")
			surf.DrawCtrl2D(true)
			surf.DrawElems2D(41, true, ", zorder=20, color='#e78005', lw=2", "")
			plt.Equal()
		}
		plt.SaveD("/tmp/gosl", "test_nurbs10.eps")
	}
}
