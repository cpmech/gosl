// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

const (
	T_NURBS_SAVE = false
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

func do_plot_nurbs_basis(b *Nurbs, la, lb int) {
	npts := 21
	plt.SetForEps(1.2, 500)

	plt.Subplot(3, 2, 1)
	b.DrawCtrl2D(false)
	b.DrawElems2D(npts, false, "", "")
	t0 := time.Now()
	b.PlotBasis(la, "", 11, 0) // 0 => CalcBasis
	io.Pfcyan("time elapsed (calcbasis) = %v\n", time.Now().Sub(t0))
	plt.Equal()

	plt.Subplot(3, 2, 2)
	b.DrawCtrl2D(false)
	b.DrawElems2D(npts, false, "", "")
	b.PlotBasis(lb, "", 11, 0) // 0 => CalcBasis
	plt.Equal()

	plt.Subplot(3, 2, 3)
	b.DrawCtrl2D(false)
	b.DrawElems2D(npts, false, "", "")
	b.PlotBasis(la, "", 11, 1) // 1 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(3, 2, 4)
	b.DrawCtrl2D(false)
	b.DrawElems2D(npts, false, "", "")
	b.PlotBasis(lb, "", 11, 1) // 1 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(3, 2, 5)
	b.DrawCtrl2D(false)
	b.DrawElems2D(npts, false, "", "")
	t0 = time.Now()
	b.PlotBasis(la, "", 11, 2) // 2 => RecursiveBasis
	io.Pfcyan("time elapsed (recursive) = %v\n", time.Now().Sub(t0))
	plt.Equal()

	plt.Subplot(3, 2, 6)
	b.DrawCtrl2D(false)
	b.DrawElems2D(npts, false, "", "")
	b.PlotBasis(lb, "", 11, 2) // 2 => RecursiveBasis
	plt.Equal()
}

func do_plot_nurbs_derivs(b *Nurbs, la, lb int) {
	np := 11
	plt.SetForEps(1.5, 500)

	plt.Subplot(4, 2, 1)
	t0 := time.Now()
	b.PlotDeriv(la, 0, "", np, 0) // 0 => CalcBasisAndDerivs
	io.Pfcyan("time elapsed (calcbasis) = %v\n", time.Now().Sub(t0))
	plt.Equal()

	plt.Subplot(4, 2, 2)
	t0 = time.Now()
	b.PlotDeriv(la, 0, "", np, 1) // 1 => NumericalDeriv
	io.Pfcyan("time elapsed (numerical) = %v\n", time.Now().Sub(t0))
	plt.Equal()

	plt.Subplot(4, 2, 3)
	b.PlotDeriv(la, 1, "", np, 0) // 0 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(4, 2, 4)
	b.PlotDeriv(la, 1, "", np, 1) // 0 => NumericalDeriv
	plt.Equal()

	plt.Subplot(4, 2, 5)
	b.PlotDeriv(lb, 0, "", np, 0) // 0 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(4, 2, 6)
	b.PlotDeriv(lb, 0, "", np, 1) // 0 => NumericalDeriv
	plt.Equal()

	plt.Subplot(4, 2, 7)
	b.PlotDeriv(lb, 1, "", np, 0) // 0 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(4, 2, 8)
	b.PlotDeriv(lb, 1, "", np, 1) // 0 => NumericalDeriv
	plt.Equal()
}

func do_plot_nurbs(b *Nurbs) {
	plt.SetForEps(0.75, 400)
	b.DrawCtrl2D(true)
	b.DrawElems2D(21, true, "", "")
	plt.Equal()
}

func do_plot_nurbs_refined(b, c *Nurbs) {
	plt.SetForEps(1.5, 400)
	plt.Subplot(3, 1, 1)
	b.DrawCtrl2D(true)
	b.DrawElems2D(21, true, "", "")
	plt.Equal()

	plt.Subplot(3, 1, 2)
	c.DrawCtrl2D(true)
	c.DrawElems2D(21, true, "", "")
	plt.Equal()

	plt.Subplot(3, 1, 3)
	b.DrawElems2D(21, true, ", lw=3", "")
	c.DrawElems2D(21, true, ", color='red', marker='+', markevery=10", "color='magenta', size=8, va='bottom'")
	plt.Equal()
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

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//verbose()
	chk.PrintTitle("nurbs01")

	b := get_nurbs_A()
	elems := b.Elements()
	enodes := b.Enodes()
	io.PfYel("enodes = %v\n", enodes)
	chk.Ints(tst, "elem[0]", elems[0], []int{2, 3, 1, 2})
	chk.Ints(tst, "elem[1]", elems[1], []int{3, 4, 1, 2})
	chk.Ints(tst, "elem[2]", elems[2], []int{4, 5, 1, 2})
	chk.Ints(tst, "enodes[0]", enodes[0], []int{0, 1, 2, 5, 6, 7})
	chk.Ints(tst, "enodes[1]", enodes[1], []int{1, 2, 3, 6, 7, 8})
	chk.Ints(tst, "enodes[2]", enodes[2], []int{2, 3, 4, 7, 8, 9})

	if T_NURBS_SAVE {
		do_plot_nurbs_basis(b, 0, 7)
		plt.SaveD("/tmp/gosl", "t_nurbs01.eps")
	}
}

func Test_nurbs02(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//verbose()
	chk.PrintTitle("nurbs02")

	b := get_nurbs_A()
	do_check_derivs(tst, b, 11, 1e-5, false)

	if T_NURBS_SAVE {
		do_plot_nurbs_derivs(b, 0, 7)
		plt.SaveD("/tmp/gosl", "t_nurbs02.eps")
	}
}

func Test_nurbs03(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//verbose()
	chk.PrintTitle("nurbs03")

	b := get_nurbs_B()

	if T_NURBS_SAVE {
		la := 0 + 0*b.n[0]
		lb := 2 + 1*b.n[0]
		do_plot_nurbs_basis(b, la, lb)
		plt.SaveD("/tmp/gosl", "t_nurbs03.eps")
	}
}

func Test_nurbs04(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//verbose()
	chk.PrintTitle("nurbs04")

	b := get_nurbs_B()
	do_check_derivs(tst, b, 11, 1e-5, false)

	if T_NURBS_SAVE {
		la := 0 + 0*b.n[0]
		lb := 2 + 1*b.n[0]
		do_plot_nurbs_derivs(b, la, lb)
		plt.SaveD("/tmp/gosl", "t_nurbs04.eps")
	}
}

func Test_nurbs05(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

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

	if T_NURBS_SAVE {
		plt.SetForEps(0.7, 300)
		b.DrawCtrl2D(true)
		b.DrawElems2D(21, true, "", "")
		plt.Equal()
		plt.SaveD("/tmp/gosl", "t_nurbs05.eps")
	}
}

func Test_nurbs06(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//verbose()
	chk.PrintTitle("nurbs06")

	b := get_nurbs_C()
	elems := b.Elements()
	enodes := b.Enodes()
	chk.Ints(tst, "elem[0]", elems[0], []int{3, 4})
	chk.Ints(tst, "elem[1]", elems[1], []int{4, 5})
	chk.Ints(tst, "elem[2]", elems[2], []int{5, 6})
	chk.Ints(tst, "enodes[0]", enodes[0], []int{0, 1, 2, 3})
	chk.Ints(tst, "enodes[1]", enodes[1], []int{1, 2, 3, 4})
	chk.Ints(tst, "enodes[2]", enodes[2], []int{2, 3, 4, 5})
	c := b.Krefine([][]float64{
		//{0.15},
		{0.15, 0.5, 0.85},
	})

	if T_NURBS_SAVE {
		do_plot_nurbs_refined(b, c)
		plt.SaveD("/tmp/gosl", "t_nurbs06.eps")
	}
}

func Test_nurbs07(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//verbose()
	chk.PrintTitle("nurbs07")

	b := get_nurbs_A()
	c := b.Krefine([][]float64{
		{0.5, 1.5, 2.5},
		{0.5},
	})

	if T_NURBS_SAVE {
		do_plot_nurbs_refined(b, c)
		plt.SaveD("/tmp/gosl", "t_nurbs07.eps")
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

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//verbose()
	chk.PrintTitle("nurbs08")

	a := get_nurbs_B()
	a_el := a.Elements()
	a_en := a.Enodes()
	chk.Ints(tst, "a_el[0]", a_el[0], []int{2, 3, 2, 3})
	chk.Ints(tst, "a_el[1]", a_el[1], []int{3, 4, 2, 3})
	chk.Ints(tst, "a_en[0]", a_en[0], []int{0, 1, 2, 4, 5, 6, 8, 9, 10})
	chk.Ints(tst, "a_en[1]", a_en[1], []int{1, 2, 3, 5, 6, 7, 9, 10, 11})

	b := a.KrefineN(2, false)
	b_el := b.Elements()
	b_en := b.Enodes()
	chk.Ints(tst, "b_el[0]", b_el[0], []int{2, 3, 2, 3})
	chk.Ints(tst, "b_el[1]", b_el[1], []int{3, 4, 2, 3})
	chk.Ints(tst, "b_el[2]", b_el[2], []int{4, 5, 2, 3})
	chk.Ints(tst, "b_el[3]", b_el[3], []int{5, 6, 2, 3})
	chk.Ints(tst, "b_el[4]", b_el[4], []int{2, 3, 3, 4})
	chk.Ints(tst, "b_el[5]", b_el[5], []int{3, 4, 3, 4})
	chk.Ints(tst, "b_el[6]", b_el[6], []int{4, 5, 3, 4})
	chk.Ints(tst, "b_el[7]", b_el[7], []int{5, 6, 3, 4})
	chk.Ints(tst, "b_en[0]", b_en[0], []int{0, 1, 2, 6, 7, 8, 12, 13, 14})
	chk.Ints(tst, "b_en[1]", b_en[1], []int{1, 2, 3, 7, 8, 9, 13, 14, 15})
	chk.Ints(tst, "b_en[2]", b_en[2], []int{2, 3, 4, 8, 9, 10, 14, 15, 16})
	chk.Ints(tst, "b_en[3]", b_en[3], []int{3, 4, 5, 9, 10, 11, 15, 16, 17})
	chk.Ints(tst, "b_en[4]", b_en[4], []int{6, 7, 8, 12, 13, 14, 18, 19, 20})
	chk.Ints(tst, "b_en[5]", b_en[5], []int{7, 8, 9, 13, 14, 15, 19, 20, 21})
	chk.Ints(tst, "b_en[6]", b_en[6], []int{8, 9, 10, 14, 15, 16, 20, 21, 22})
	chk.Ints(tst, "b_en[7]", b_en[7], []int{9, 10, 11, 15, 16, 17, 21, 22, 23})

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
	utl.CompareDeep4(tst, "Read: Q", a.Q, B[0].Q)
	utl.CheckIntMat(tst, "Read: l2i", a.l2i, B[0].l2i)

	if T_NURBS_SAVE {
		do_plot_nurbs_refined(a, c)
		plt.SaveD("/tmp/gosl", "t_nurbs08.eps")
		do_plot_nurbs(B[0])
		plt.SaveD("/tmp/gosl", "t_nurbs08_read.eps")
	}
}
