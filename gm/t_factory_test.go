// Copyright 2016 The Gosl Authors. All rights reserved.
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

func Test_factory01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("factory01. quarter circle curve in 2d")

	// geometry
	xc, yc, r := 0.5, 1.0, 2.0
	curve := FactoryNurbs.Curve2dQuarterCircle(xc, yc, r)

	// function to check circle
	npcheck := 5
	U := utl.LinSpace(curve.b[0].tmin, curve.b[0].tmax, npcheck)
	x := make([]float64, 2)
	for _, u := range U {
		curve.Point(x, []float64{u}, 2)
		e := math.Sqrt(math.Pow(x[0]-xc, 2)+math.Pow(x[1]-yc, 2)) - r
		chk.Scalar(tst, io.Sf("error @ (%+.8f,%+.8f) == 0?", x[0], x[1]), 1e-15, e, 0)
	}

	// plot
	if chk.Verbose {
		extra := func() {
			plt.Circle(xc, yc, r, &plt.A{C: "#478275", Lw: 1})
		}
		ndim := 2
		argsCurve := &plt.A{C: "orange", M: "+", Mec: "k", Lw: 4, L: "curve", NoClip: true}
		argsCtrl := &plt.A{C: "k", M: ".", Ls: "--", L: "control", NoClip: true}
		argsIds := &plt.A{C: "b", Fsz: 10}
		plt.Reset(false, nil)
		plt.Equal()
		plt.HideAllBorders()
		plt.PlotOne(xc, yc, &plt.A{C: "k", M: "+", Ms: 20})
		PlotNurbs("/tmp/gosl", "t_factory01", curve, ndim, 11, true, true, argsCurve, argsCtrl, argsIds, extra)
	}
}

func Test_factory02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("factory02. circle curve in 2d")

	// geometry
	xc, yc, r := 0.5, 0.25, 1.75

	// function to check circle
	npcheck := 11
	checkCircle := func(nurbs *Nurbs) {
		U := utl.LinSpace(nurbs.b[0].tmin, nurbs.b[0].tmax, npcheck)
		x := make([]float64, 2)
		for _, u := range U {
			nurbs.Point(x, []float64{u}, 2)
			e := math.Sqrt(math.Pow(x[0]-xc, 2)+math.Pow(x[1]-yc, 2)) - r
			chk.Scalar(tst, io.Sf("error @ (%+.8f,%+.8f) == 0?", x[0], x[1]), 1e-15, e, 0)
		}
	}

	// original curve
	curve := FactoryNurbs.Curve2dCircle(xc, yc, r)
	checkCircle(curve)

	// refine NURBS
	refined := curve.Krefine([][]float64{{0.5, 1.5, 2.5, 3.5}})
	checkCircle(refined)

	// plot
	if chk.Verbose {

		argsIdsA := &plt.A{C: "b", Fsz: 10}
		argsCtrlA := &plt.A{C: "k", M: ".", Ls: "--", L: "control", NoClip: true}
		argsCurveA := &plt.A{C: "orange", M: "+", Mec: "k", Lw: 4, L: "curve", NoClip: true}

		argsIdsB := &plt.A{C: "green", Fsz: 7}
		argsCtrlB := &plt.A{C: "green", L: "refined: control"}
		argsElemsB := &plt.A{C: "orange", Ls: "none", M: "*", Me: 20, L: "refined: curve"}

		ndim := 2
		np := 11
		extra := func() {
			plt.Circle(xc, yc, r, &plt.A{C: "#478275", Lw: 1})
			refined.DrawCtrl(ndim, true, argsCtrlB, argsIdsB)
			refined.DrawElems(ndim, np, false, argsElemsB, nil)
		}

		plt.Reset(false, nil)
		plt.Equal()
		plt.HideAllBorders()
		plt.AxisRange(-3, 3, -3, 3)
		PlotNurbs("/tmp/gosl", "t_factory02", curve, ndim, np, true, true, argsCurveA, argsCtrlA, argsIdsA, extra)
	}
}

func Test_factory03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("factory04. toroidal surface")

	// geometry
	xc, yc, zc, r, R := 0.0, 0.0, 0.0, 1.0, 5.0

	// function to check toroidal surface
	tol := 1e-14
	npcheck := 9
	checkTorus := func(nurbs *Nurbs) {
		U := utl.LinSpace(nurbs.b[0].tmin, nurbs.b[0].tmax, npcheck)
		V := utl.LinSpace(nurbs.b[1].tmin, nurbs.b[1].tmax, npcheck)
		x := make([]float64, 3)
		for _, u := range U {
			for _, v := range V {
				nurbs.Point(x, []float64{u, v}, 3)
				e := math.Pow(R-math.Sqrt(math.Pow(x[0]-xc, 2)+math.Pow(x[1]-yc, 2)), 2) +
					math.Pow(x[2]-zc, 2) - r*r
				chk.Scalar(tst, io.Sf("error @ (%+.8f,%+8f,%+8f) == 0?", x[0], x[1], x[2]), tol, e, 0)
			}
		}
	}

	// surface
	surf := FactoryNurbs.Surf3dTorus(xc, yc, zc, r, R)
	checkTorus(surf)
}
