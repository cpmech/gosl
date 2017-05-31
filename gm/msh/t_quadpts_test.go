// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"math"
	"strings"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/plt"
)

func TestQuadpts01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("quadpts01. quadrature points for 'lin'")

	linPoints := map[int][][]float64{
		1: [][]float64{
			{0, 0, 0, 2},
		},
		2: [][]float64{
			{-0.5773502691896257, 0, 0, 1},
			{+0.5773502691896257, 0, 0, 1},
		},
		3: [][]float64{
			{-0.7745966692414834, 0, 0, 0.5555555555555556},
			{+0.0000000000000000, 0, 0, 0.8888888888888888},
			{+0.7745966692414834, 0, 0, 0.5555555555555556},
		},
		4: [][]float64{
			{-0.8611363115940526, 0, 0, 0.3478548451374538},
			{-0.3399810435848562, 0, 0, 0.6521451548625462},
			{+0.3399810435848562, 0, 0, 0.6521451548625462},
			{+0.8611363115940526, 0, 0, 0.3478548451374538},
		},
		5: [][]float64{
			{-0.9061798459386640, 0, 0, 0.2369268850561891},
			{-0.5384693101056831, 0, 0, 0.4786286704993665},
			{+0.0000000000000000, 0, 0, 0.5688888888888889},
			{+0.5384693101056831, 0, 0, 0.4786286704993665},
			{+0.9061798459386640, 0, 0, 0.2369268850561891},
		},
	}

	for n, qpts := range linPoints {
		pts := linIntPointsSet.Find("LE", n) // Gauss-Legendre
		if pts == nil {
			tst.Errorf("cannot find set of %d points for lin elements\n", n)
			return
		}
		io.Pfblue2("\n----------------------------------- lin %d -----------------------------------\n", n)
		for i := 0; i < n; i++ {
			x := qpts[i][:1]
			w := qpts[i][3]
			io.Pforan("x=%19v w=%19v\n", pts.Points[i].X, pts.Points[i].W)
			io.Pfgrey("  %19v w=%19v\n", x, w)
			chk.Vector(tst, io.Sf("lin: %d: x%d", n, i), 1e-15, x, pts.Points[i].X)
			chk.Scalar(tst, io.Sf("lin: %d: w%d", n, i), 1e-15, w, pts.Points[i].W)
		}
	}
}

func TestQuadpts02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("quadpts02. quadrature points for 'qua'")

	quaPoints := map[string][][]float64{
		"LE_4": [][]float64{
			{-0.5773502691896257, -0.5773502691896257, 0, 1},
			{+0.5773502691896257, -0.5773502691896257, 0, 1},
			{-0.5773502691896257, +0.5773502691896257, 0, 1},
			{+0.5773502691896257, +0.5773502691896257, 0, 1},
		},
		"LE_9": [][]float64{
			{-0.7745966692414834, -0.7745966692414834, 0, 25.0 / 81.0},
			{+0.0000000000000000, -0.7745966692414834, 0, 40.0 / 81.0},
			{+0.7745966692414834, -0.7745966692414834, 0, 25.0 / 81.0},
			{-0.7745966692414834, +0.0000000000000000, 0, 40.0 / 81.0},
			{+0.0000000000000000, +0.0000000000000000, 0, 64.0 / 81.0},
			{+0.7745966692414834, +0.0000000000000000, 0, 40.0 / 81.0},
			{-0.7745966692414834, +0.7745966692414834, 0, 25.0 / 81.0},
			{+0.0000000000000000, +0.7745966692414834, 0, 40.0 / 81.0},
			{+0.7745966692414834, +0.7745966692414834, 0, 25.0 / 81.0},
		},
		"W5corner_5": [][]float64{
			{-1, -1, 0, 1.0 / 3.0},
			{+1, -1, 0, 1.0 / 3.0},
			{+0, +0, 0, 8.0 / 3.0},
			{-1, +1, 0, 1.0 / 3.0},
			{+1, +1, 0, 1.0 / 3.0},
		},
		"W4stable_5": [][]float64{
			{-0.5776391, -0.5776391, 0, 0.999},
			{+0.5776391, -0.5776391, 0, 0.999},
			{+0.0000000, +0.0000000, 0, 0.004},
			{-0.5776391, +0.5776391, 0, 0.999},
			{+0.5776391, +0.5776391, 0, 0.999},
		},
	}

	for key, qpts := range quaPoints {
		res := strings.Split(key, "_")
		rule, n := res[0], io.Atoi(res[1])
		pts := quaIntPointsSet.Find(rule, n) // Gauss-Legendre
		if pts == nil {
			tst.Errorf("cannot find set of %d points for qua elements\n", n)
			return
		}
		io.Pfblue2("\n----------------------------------- qua %d -----------------------------------\n", n)
		for i := 0; i < n; i++ {
			x := qpts[i][:2]
			w := qpts[i][3]
			io.Pforan("x=%19v w=%19v\n", pts.Points[i].X, pts.Points[i].W)
			io.Pfgrey("  %19v w=%19v\n", x, w)
			chk.Vector(tst, io.Sf("qua: %d: x%d", n, i), 1e-15, x, pts.Points[i].X)
			chk.Scalar(tst, io.Sf("qua: %d: w%d", n, i), 1e-15, w, pts.Points[i].W)
		}
	}

	io.Pfblue2("\n-------------------------------- qua 5 W5 corner ---------------------------\n")
	w5crn := NewIntPoints("W5", 2, 5, []*fun.P{&fun.P{N: "w0", V: 8.0 / 3.0}})
	tmppt := quaPoints["W5corner_5"]
	for i := 0; i < 5; i++ {
		x := tmppt[i][:2]
		w := tmppt[i][3]
		io.Pforan("x=%19v w=%19v\n", w5crn.Points[i].X, w5crn.Points[i].W)
		io.Pfgrey("  %19v w=%19v\n", x, w)
		chk.Vector(tst, io.Sf("qua: 5: x%d", i), 1e-15, w5crn.Points[i].X, x)
		chk.Scalar(tst, io.Sf("qua: 5: w%d", i), 1e-15, w5crn.Points[i].W, w)
	}

	io.Pfblue2("\n-------------------------------- qua 8 W8 fixed ----------------------------\n")
	w8fixA := NewIntPoints("W8", 2, 8, []*fun.P{&fun.P{N: "wb", V: 40.0 / 49.0}})
	w8fixB := quaIntPointsSet.Find("W8fixed", 8)
	for i := 0; i < 8; i++ {
		io.Pforan("x=%19v w=%19v\n", w8fixA.Points[i].X, w8fixA.Points[i].W)
		io.Pfgrey("  %19v   %19v\n", w8fixB.Points[i].X, w8fixB.Points[i].W)
		chk.Vector(tst, io.Sf("qua: 8: x%d", i), 1e-15, w8fixA.Points[i].X, w8fixB.Points[i].X)
		chk.Scalar(tst, io.Sf("qua: 8: w%d", i), 1e-15, w8fixA.Points[i].W, w8fixB.Points[i].W)
	}

	if chk.Verbose {
		plt.Reset(true, nil)
		w5crn := quaIntPointsSet.Find("W5corner", 5)
		w5crn.Draw(nil, nil)
		w4sta := quaIntPointsSet.Find("W4stable", 5)
		w4sta.Draw([]float64{2.5, 0}, nil)
		w8fixA.Draw([]float64{0, 2.5}, nil)
		w8fixB.Draw([]float64{2.5, 2.5}, nil)
		plt.Equal()
		plt.AxisRange(-2, 4.5, -2, 4.5)
		plt.HideAllBorders()
		plt.Gll("x", "y", nil)
		plt.Save("/tmp/gosl", "quadpts02")
	}
}

func TestQuadpts03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("quadpts03. quadrature points for 'hex'")

	hexPoints := map[int][][]float64{
		8: [][]float64{
			{-0.5773502691896257, -0.5773502691896257, -0.5773502691896257, 1},
			{+0.5773502691896257, -0.5773502691896257, -0.5773502691896257, 1},
			{-0.5773502691896257, +0.5773502691896257, -0.5773502691896257, 1},
			{+0.5773502691896257, +0.5773502691896257, -0.5773502691896257, 1},
			{-0.5773502691896257, -0.5773502691896257, +0.5773502691896257, 1},
			{+0.5773502691896257, -0.5773502691896257, +0.5773502691896257, 1},
			{-0.5773502691896257, +0.5773502691896257, +0.5773502691896257, 1},
			{+0.5773502691896257, +0.5773502691896257, +0.5773502691896257, 1},
		},
		27: [][]float64{
			{-0.774596669241483, -0.774596669241483, -0.774596669241483, 0.171467764060357},
			{+0.000000000000000, -0.774596669241483, -0.774596669241483, 0.274348422496571},
			{+0.774596669241483, -0.774596669241483, -0.774596669241483, 0.171467764060357},
			{-0.774596669241483, +0.000000000000000, -0.774596669241483, 0.274348422496571},
			{+0.000000000000000, +0.000000000000000, -0.774596669241483, 0.438957475994513},
			{+0.774596669241483, +0.000000000000000, -0.774596669241483, 0.274348422496571},
			{-0.774596669241483, +0.774596669241483, -0.774596669241483, 0.171467764060357},
			{+0.000000000000000, +0.774596669241483, -0.774596669241483, 0.274348422496571},
			{+0.774596669241483, +0.774596669241483, -0.774596669241483, 0.171467764060357},
			{-0.774596669241483, -0.774596669241483, +0.000000000000000, 0.274348422496571},
			{+0.000000000000000, -0.774596669241483, +0.000000000000000, 0.438957475994513},
			{+0.774596669241483, -0.774596669241483, +0.000000000000000, 0.274348422496571},
			{-0.774596669241483, +0.000000000000000, +0.000000000000000, 0.438957475994513},
			{+0.000000000000000, +0.000000000000000, +0.000000000000000, 0.702331961591221},
			{+0.774596669241483, +0.000000000000000, +0.000000000000000, 0.438957475994513},
			{-0.774596669241483, +0.774596669241483, +0.000000000000000, 0.274348422496571},
			{+0.000000000000000, +0.774596669241483, +0.000000000000000, 0.438957475994513},
			{+0.774596669241483, +0.774596669241483, +0.000000000000000, 0.274348422496571},
			{-0.774596669241483, -0.774596669241483, +0.774596669241483, 0.171467764060357},
			{+0.000000000000000, -0.774596669241483, +0.774596669241483, 0.274348422496571},
			{+0.774596669241483, -0.774596669241483, +0.774596669241483, 0.171467764060357},
			{-0.774596669241483, +0.000000000000000, +0.774596669241483, 0.274348422496571},
			{+0.000000000000000, +0.000000000000000, +0.774596669241483, 0.438957475994513},
			{+0.774596669241483, +0.000000000000000, +0.774596669241483, 0.274348422496571},
			{-0.774596669241483, +0.774596669241483, +0.774596669241483, 0.171467764060357},
			{+0.000000000000000, +0.774596669241483, +0.774596669241483, 0.274348422496571},
			{+0.774596669241483, +0.774596669241483, +0.774596669241483, 0.171467764060357},
		},
	}

	for n, qpts := range hexPoints {
		pts := hexIntPointsSet.Find("LE", n) // Gauss-Legendre
		if pts == nil {
			tst.Errorf("cannot find set of %d points for hex elements\n", n)
			return
		}
		io.Pfblue2("\n----------------------------------- hex %d -----------------------------------\n", n)
		for i := 0; i < n; i++ {
			x := qpts[i][:3]
			w := qpts[i][3]
			io.Pforan("x=%19v w=%19v\n", pts.Points[i].X, pts.Points[i].W)
			io.Pfgrey("  %19v w=%19v\n", x, w)
			chk.Vector(tst, io.Sf("hex: %d: x%d", n, i), 1e-15, x, pts.Points[i].X)
			chk.Scalar(tst, io.Sf("hex: %d: w%d", n, i), 1e-14, w, pts.Points[i].W)
		}
	}
}

func TestQuadpts10(tst *testing.T) {

	//verbose()
	chk.PrintTitle("quadpts02. quadrature points")

	degreeMax := 5
	glX := make([][]float64, degreeMax+1)
	glW := make([][]float64, degreeMax+1)
	for n := 1; n <= degreeMax; n++ {
		glX[n], glW[n] = num.GaussLegendreXW(-1, 1, n)
	}

	for name, allPts := range IntPointsOld {

		io.PfYel("\n--------------------------------- %-6s---------------------------------\n", name)

		switch name {
		case "lin":
			for n, pts := range allPts {
				x := make([]float64, n)
				w := make([]float64, n)
				for i := 0; i < n; i++ {
					x[i] = pts[i][0]
					w[i] = pts[i][3]
				}
				io.Pf("\nx = %v\n", x)
				io.Pfgreen("    %v\n", glX[n])
				io.Pf("w = %v\n", w)
				io.Pfgreen("    %v\n", glW[n])
				chk.Vector(tst, io.Sf("lin:%d x", n), 1e-15, x, glX[n])
				chk.Vector(tst, io.Sf("lin:%d w", n), 1e-15, w, glW[n])
			}

		case "qua":
			for n, pts := range allPts {
				io.Pl()
				n1d := int(math.Sqrt(float64(n)))
				x1d := glX[n1d]
				w1d := glW[n1d]
				for j := 0; j < n1d; j++ {
					for i := 0; i < n1d; i++ {
						m := i + n1d*j
						x := pts[m][:2]
						v := pts[m][3]
						y := []float64{x1d[i], x1d[j]}
						w := w1d[i] * w1d[j]
						io.Pf("  %d%d x = %23v  w = %23v\n", i, j, x, v)
						io.Pfgreen("         %23v      %23v\n", y, w)
						chk.Vector(tst, "x", 1e-15, x, y)
						chk.Scalar(tst, "w", 1e-15, v, w)
					}
				}
			}

		case "hex":
			for n, pts := range allPts {
				if n == 14 {
					continue
				}
				io.Pl()
				n1d := int(math.Floor(math.Pow(float64(n), 1.0/3.0) + 0.5))
				x1d := glX[n1d]
				w1d := glW[n1d]
				for k := 0; k < n1d; k++ {
					for j := 0; j < n1d; j++ {
						for i := 0; i < n1d; i++ {
							m := i + n1d*j + (n1d*n1d)*k
							x := pts[m][:3]
							v := pts[m][3]
							y := []float64{x1d[i], x1d[j], x1d[k]}
							w := w1d[i] * w1d[j] * w1d[k]
							io.Pf("%d%d x=%18v w=%18v\n", i, j, x, v)
							io.Pfgreen("     %18v   %18v\n", y, w)
							chk.Vector(tst, "x", 1e-15, x, y)
							chk.Scalar(tst, "w", 1e-14, v, w)
						}
					}
				}
			}
		}

	}
}
