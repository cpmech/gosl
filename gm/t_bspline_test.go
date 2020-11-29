// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
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
}
