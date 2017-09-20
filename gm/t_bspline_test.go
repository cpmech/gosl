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

	T1 := []float64{0, 0, 0, 1, 1, 1}
	s1 := NewBspline(T1, 2)
	s1.SetControl([][]float64{{0, 0}, {0.5, 1}, {1, 0}})

	T2 := []float64{0, 0, 0, 0.5, 1, 1, 1}
	s2 := NewBspline(T2, 2)
	s2.SetControl([][]float64{{0, 0}, {0.25, 0.5}, {0.75, 0.5}, {1, 0}})

	if chk.Verbose {

		argsRec := &plt.A{C: "k", M: "+", Me: 15, Ls: "none", L: "recursive"}

		npts := 201
		plt.Reset(true, &plt.A{Prop: 1.5})
		plt.SplotGap(0.2, 0.4)

		plt.Subplot(3, 2, 1)
		s1.Draw2d(npts, 0, true, nil, nil)      // 0 => CalcBasis
		s1.Draw2d(npts, 1, false, argsRec, nil) // 1 => RecursiveBasis

		plt.Subplot(3, 2, 2)
		plt.SetAxis(0, 1, 0, 1)
		s2.Draw2d(npts, 0, true, nil, nil)      // 0 => CalcBasis
		s2.Draw2d(npts, 1, false, argsRec, nil) // 1 => RecursiveBasis

		plt.Subplot(3, 2, 3)
		s1.PlotBasis(npts, 0) // 0 => CalcBasis
		s1.PlotBasis(npts, 1) // 1 => CalcBasisAndDerivs
		s1.PlotBasis(npts, 2) // 2 => RecursiveBasis

		plt.Subplot(3, 2, 4)
		s2.PlotBasis(npts, 0) // 0 => CalcBasis
		s2.PlotBasis(npts, 1) // 1 => CalcBasisAndDerivs
		s2.PlotBasis(npts, 2) // 2 => RecursiveBasis

		plt.Subplot(3, 2, 5)
		s1.PlotDerivs(npts) // 0 => CalcBasisAndDerivs

		plt.Subplot(3, 2, 6)
		s2.PlotDerivs(npts) // 0 => CalcBasisAndDerivs

		plt.Save("/tmp/gosl", "bspline01")
	}
}

func Test_bspline02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bspline02")

	//               0 1 2 3 4 5 6 7 8 9 10
	T := []float64{0, 0, 0, 1, 2, 3, 4, 4, 5, 5, 5}
	sol := []int{2, 2, 3, 3, 4, 4, 5, 5, 7, 7, 7}
	s := NewBspline(T, 2)
	s.SetControl([][]float64{{0, 0}, {0.5, 1}, {1, 0}, {1.5, 0}, {2, 1}, {2.5, 1}, {3, 0.5}, {3.5, 0}})

	tt := utl.LinSpace(0, 5, 11)
	for k, t := range tt {
		span := s.findSpan(t)
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
		chk.Array(tst, "Point", tol, pa, pb)
	}

	if chk.Verbose {
		argsRec := &plt.A{C: "k", M: "+", Me: 15, Ls: "none", L: "recursive"}
		npts := 201
		plt.Reset(false, nil)
		s.Draw2d(npts, 0, true, nil, nil)      // 0 => CalcBasis
		s.Draw2d(npts, 1, false, argsRec, nil) // 1 => RecursiveBasis
		plt.Plot(xx, yy, &plt.A{C: "b", L: "check"})
		plt.Gll("x", "y", nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl", "bspline02")
	}
}

func Test_bspline03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bspline03")

	//             0 1 2 3 4 5 6 7 8 9 10
	T := []float64{0, 0, 0, 1, 2, 3, 4, 4, 5, 5, 5}
	s := NewBspline(T, 2)
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

	δ := 1e-1 // used to avoid using central differences @ boundaries of t in [0,5]
	ver := chk.Verbose
	tol := 1e-5
	tt := utl.LinSpace(δ, 5.0-δ, 11)
	anad := make([]float64, s.NumBasis())
	for _, t := range tt {
		for i := 0; i < s.NumBasis(); i++ {
			s.CalcBasisAndDerivs(t)
			anad[i] = s.GetDeriv(i)
			chk.DerivScaSca(tst, io.Sf("i=%d t=%.3f", i, t), tol, anad[i], t, 1e-1, ver, func(x float64) float64 {
				return s.RecursiveBasis(x, i)
			})
		}
	}

	if chk.Verbose {

		argsRec := &plt.A{C: "k", M: "+", Me: 15, Ls: "none", L: "recursive"}

		npts := 201
		plt.Reset(true, &plt.A{Prop: 1.5})
		plt.SplotGap(0, 0.3)

		plt.Subplot(3, 1, 1)
		s.Draw2d(npts, 0, true, nil, nil)      // 0 => CalcBasis
		s.Draw2d(npts, 1, false, argsRec, nil) // 1 => RecursiveBasis

		plt.Subplot(3, 1, 2)
		s.PlotBasis(npts, 0) // 0 => CalcBasis
		s.PlotBasis(npts, 1) // 1 => CalcBasisAndDerivs
		s.PlotBasis(npts, 2) // 2 => RecursiveBasis

		plt.Subplot(3, 1, 3)
		s.PlotDerivs(npts) // 0 => CalcBasisAndDerivs

		plt.Save("/tmp/gosl", "bspline03")
	}
}
