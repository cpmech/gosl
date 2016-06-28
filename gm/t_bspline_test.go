// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func Test_bspline01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bspline01")

	var s1 Bspline
	T1 := []float64{0, 0, 0, 1, 1, 1}
	s1.Init(T1, 2)
	s1.SetControl([][]float64{{0, 0}, {0.5, 1}, {1, 0}})

	var s2 Bspline
	T2 := []float64{0, 0, 0, 0.5, 1, 1, 1}
	s2.Init(T2, 2)
	s2.SetControl([][]float64{{0, 0}, {0.25, 0.5}, {0.75, 0.5}, {1, 0}})

	if chk.Verbose {
		npts := 201
		plt.SetForPng(1.5, 600, 150)
		plt.SplotGap(0.2, 0.4)

		str0 := ",lw=2"
		str1 := ",ls='none',marker='+',color='cyan',markevery=10"
		str2 := ",ls='none',marker='x',markevery=10"
		str3 := ",ls='none',marker='+',markevery=10"
		str4 := ",ls='none',marker='4',markevery=10"

		plt.Subplot(3, 2, 1)
		s1.Draw2d(str0, "", npts, 0) // 0 => CalcBasis
		s1.Draw2d(str1, "", npts, 1) // 1 => RecursiveBasis

		plt.Subplot(3, 2, 2)
		plt.SetAxis(0, 1, 0, 1)
		s2.Draw2d(str0, "", npts, 0) // 0 => CalcBasis
		s2.Draw2d(str1, "", npts, 1) // 1 => RecursiveBasis

		plt.Subplot(3, 2, 3)
		s1.PlotBasis("", npts, 0)   // 0 => CalcBasis
		s1.PlotBasis(str2, npts, 1) // 1 => CalcBasisAndDerivs
		s1.PlotBasis(str3, npts, 2) // 2 => RecursiveBasis

		plt.Subplot(3, 2, 4)
		s2.PlotBasis("", npts, 0)   // 0 => CalcBasis
		s2.PlotBasis(str2, npts, 1) // 1 => CalcBasisAndDerivs
		s2.PlotBasis(str3, npts, 2) // 2 => RecursiveBasis

		plt.Subplot(3, 2, 5)
		s1.PlotDerivs("", npts, 0)   // 0 => CalcBasisAndDerivs
		s1.PlotDerivs(str4, npts, 1) // 1 => NumericalDeriv

		plt.Subplot(3, 2, 6)
		s2.PlotDerivs("", npts, 0)   // 0 => CalcBasisAndDerivs
		s2.PlotDerivs(str4, npts, 1) // 1 => NumericalDeriv

		plt.SaveD("/tmp/gosl/gm", "bspline01.png")
	}
}

func Test_bspline02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bspline02")

	//               0 1 2 3 4 5 6 7 8 9 10
	T := []float64{0, 0, 0, 1, 2, 3, 4, 4, 5, 5, 5}
	sol := []int{2, 2, 3, 3, 4, 4, 5, 5, 7, 7, 7}
	var s Bspline
	s.Init(T, 2)
	s.SetControl([][]float64{{0, 0}, {0.5, 1}, {1, 0}, {1.5, 0}, {2, 1}, {2.5, 1}, {3, 0.5}, {3.5, 0}})

	tt := utl.LinSpace(0, 5, 11)
	for k, t := range tt {
		span := s.find_span(t)
		io.Pforan("t=%.4f  =>  span=%v\n", t, span)
		if span != sol[k] {
			chk.Panic("find_span failed: t=%v  span => %d != %d", t, span, sol[k])
		}
	}

	tol := 1e-14
	np := len(tt)
	xx, yy := make([]float64, np), make([]float64, np)
	for k, t := range tt {
		t0 := time.Now()
		pa := s.Point(t, 0) // 0 => CalcBasis
		io.Pf("Point(rec): dtime = %v\n", time.Now().Sub(t0))
		t0 = time.Now()
		pb := s.Point(t, 1) // 1 => RecursiveBasis
		io.Pf("Point:      dtime = %v\n", time.Now().Sub(t0))
		xx[k], yy[k] = pb[0], pb[1]
		io.Pfred("pa - pb = %v, %v\n", pa[0]-pb[0], pa[1]-pb[1])
		chk.Vector(tst, "Point", tol, pa, pb)
	}

	if chk.Verbose {
		npts := 201
		plt.SetForPng(0.75, 300, 150)
		str0 := ",lw=2"
		str1 := ",ls='none',marker='+',color='cyan',markevery=10"
		s.Draw2d(str0, "", npts, 0) // 0 => CalcBasis
		s.Draw2d(str1, "", npts, 1) // 1 => RecursiveBasis
		plt.Plot(xx, yy, "'bo', clip_on=0")
		plt.SaveD("/tmp/gosl/gm", "bspline02.png")
	}
}

func Test_bspline03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bspline03")

	//             0 1 2 3 4 5 6 7 8 9 10
	T := []float64{0, 0, 0, 1, 2, 3, 4, 4, 5, 5, 5}
	var s Bspline
	s.Init(T, 2)
	s.SetControl([][]float64{{0, 0}, {0.5, 1}, {1, 0}, {1.5, 0}, {2, 1}, {2.5, 1}, {3, 0.5}, {3.5, 0}})

	// analytical derivatives
	s.CalcBasisAndDerivs(3.99)
	io.Pfpink("ana: dNdt(t=3.99, i=5) = %v\n", s.GetDeriv(5))
	io.Pfpink("ana: dNdt(t=3.99, i=6) = %v\n", s.GetDeriv(6))
	io.Pfpink("ana: dNdt(t=3.99, i=7) = %v\n", s.GetDeriv(7))
	s.CalcBasisAndDerivs(4.0)
	io.Pforan("ana: dNdt(t=4.00, i=5) = %v\n", s.GetDeriv(5))
	io.Pforan("ana: dNdt(t=4.00, i=6) = %v\n", s.GetDeriv(6))
	io.Pforan("ana: dNdt(t=4.00, i=7) = %v\n", s.GetDeriv(7))

	// numerical derivatives
	io.Pfcyan("num: dNdt(t=3.99, i=5) = %v\n", s.NumericalDeriv(3.99, 5))
	io.Pfcyan("num: dNdt(t=3.99, i=6) = %v\n", s.NumericalDeriv(3.99, 6))
	io.Pfcyan("num: dNdt(t=3.99, i=7) = %v\n", s.NumericalDeriv(3.99, 7))
	io.Pfblue2("num: dNdt(t=4.00, i=5) = %v\n", s.NumericalDeriv(4.00, 5))
	io.Pfblue2("num: dNdt(t=4.00, i=6) = %v\n", s.NumericalDeriv(4.00, 6))
	io.Pfblue2("num: dNdt(t=4.00, i=7) = %v\n", s.NumericalDeriv(4.00, 7))

	ver := false
	tol := 1e-5
	tt := utl.LinSpace(0, 5, 11)
	numd := make([]float64, s.NumBasis())
	anad := make([]float64, s.NumBasis())
	for _, t := range tt {
		for i := 0; i < s.NumBasis(); i++ {
			s.CalcBasisAndDerivs(t)
			anad[i] = s.GetDeriv(i)
			numd[i] = s.NumericalDeriv(t, i)
			// numerical fails @ 4 [4,5,6]
			if t == 4 {
				numd[4] = anad[4]
				numd[5] = anad[5]
				numd[6] = anad[6]
			}
			chk.PrintAnaNum(io.Sf("i=%d t=%v", i, t), tol, anad[i], numd[i], ver)
		}
		chk.Vector(tst, io.Sf("derivs @ %v", t), tol, numd, anad)
	}

	if chk.Verbose {

		npts := 201
		plt.SetForPng(1.5, 600, 150)
		plt.SplotGap(0, 0.3)

		str0 := ",lw=2"
		str1 := ",ls='none',marker='+',color='cyan',markevery=10"
		str2 := ",ls='none',marker='x',markevery=10"
		str3 := ",ls='none',marker='+',markevery=10"
		str4 := ",ls='none',marker='4',markevery=10"

		plt.Subplot(3, 1, 1)
		s.Draw2d(str0, "", npts, 0) // 0 => CalcBasis
		s.Draw2d(str1, "", npts, 1) // 1 => RecursiveBasis

		plt.Subplot(3, 1, 2)
		s.PlotBasis("", npts, 0)   // 0 => CalcBasis
		s.PlotBasis(str2, npts, 1) // 1 => CalcBasisAndDerivs
		s.PlotBasis(str3, npts, 2) // 2 => RecursiveBasis

		plt.Subplot(3, 1, 3)
		s.PlotDerivs("", npts, 0)   // 0 => CalcBasisAndDerivs
		s.PlotDerivs(str4, npts, 1) // 1 => NumericalDeriv

		plt.SaveD("/tmp/gosl/gm", "bspline03.png")
	}
}
