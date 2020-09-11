// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/utl"
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
		chk.Float64(tst, io.Sf("error @ (%+.8f,%+.8f) == 0?", x[0], x[1]), 1e-15, e, 0)
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
			chk.Float64(tst, io.Sf("error @ (%+.8f,%+.8f) == 0?", x[0], x[1]), 1e-15, e, 0)
		}
	}

	// original curve
	curve := FactoryNurbs.Curve2dCircle(xc, yc, r)
	checkCircle(curve)

	// refine NURBS
	refined := curve.Krefine([][]float64{{0.5, 1.5, 2.5, 3.5}})
	checkCircle(refined)
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
				chk.Float64(tst, io.Sf("error @ (%+.8f,%+8f,%+8f) == 0?", x[0], x[1], x[2]), tol, e, 0)
			}
		}
	}

	// surface
	surf := FactoryNurbs.Surf3dTorus(xc, yc, zc, r, R)
	checkTorus(surf)
}
