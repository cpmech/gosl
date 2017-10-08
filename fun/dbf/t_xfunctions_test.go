// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
)

func Test_xfun01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("xfun01. 2D halo => circle")

	o := New("halo", []*P{
		{N: "r", V: 0.5},
		{N: "xc", V: 0.5},
		{N: "yc", V: 0.5},
	})

	tcte := 0.0
	xmin := []float64{-1, -1}
	xmax := []float64{2, 2}
	np := 21
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotX(o, "/tmp/gosl", "t_halo", tcte, xmin, xmax, np)
	}

	np = 4
	sktol := 1e-10
	dtol := 1e-10
	ver := chk.Verbose
	CheckDerivX(tst, o, tcte, xmin, xmax, np, nil, sktol, dtol, ver)
}

func Test_xfun02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("xfun02. 2D circle distance")

	xc := []float64{0.5, 0.5}
	o := New("cdist", []*P{
		{N: "r", V: 0.5},
		{N: "xc", V: xc[0]},
		{N: "yc", V: xc[1]},
	})

	tcte := 0.0
	xmin := []float64{-1, -1}
	xmax := []float64{2, 2}
	np := 21
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotX(o, "/tmp/gosl", "t_cdist", tcte, xmin, xmax, np)
	}

	np = 5
	sktol := 1e-10
	xskip := [][]float64{xc}
	dtol := 1e-10
	ver := chk.Verbose
	CheckDerivX(tst, o, tcte, xmin, xmax, np, xskip, sktol, dtol, ver)
}

func Test_xfun03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("xfun03. xpoly2: 2nd order polynomial with x coordinates")

	o := New("xpoly2", []*P{
		{N: "a0", V: 1.5}, {N: "a1", V: 0.5}, {N: "a2", V: -1.5},
		{N: "b0", V: -1.5}, {N: "b1", V: -0.5}, {N: "b2", V: 1.5},
		{N: "c01", V: 2.0}, {N: "c12", V: -2.0}, {N: "c20", V: 1.0},
		//&Prm{N: "2D", V: 1},
	})

	tcte := 0.0
	xmin := []float64{-1, -1, -1}
	xmax := []float64{2, 2, 2}
	np := 21
	if chk.Verbose && len(xmin) == 2 {
		plt.Reset(false, nil)
		PlotX(o, "/tmp/gosl", "t_xpoly2", tcte, xmin, xmax, np)
	}

	np = 3
	sktol := 1e-10
	xskip := [][]float64{}
	dtol := 1e-10
	ver := chk.Verbose
	CheckDerivX(tst, o, tcte, xmin, xmax, np, xskip, sktol, dtol, ver)
}

func Test_xfun04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("xfun04. xpoly1: 1st order polynomial with x coordinates")

	o := New("xpoly2", []*P{
		{N: "a0", V: 1.5}, {N: "a1", V: 0.5}, {N: "a2", V: -1.5},
		{N: "2D", V: 1},
	})

	tcte := 0.0
	xmin := []float64{-1, -1, -1}
	xmax := []float64{2, 2, 2}
	np := 21
	if chk.Verbose && len(xmin) == 2 {
		plt.Reset(false, nil)
		PlotX(o, "/tmp/gosl", "t_xpoly2", tcte, xmin, xmax, np)
	}

	np = 3
	sktol := 1e-10
	xskip := [][]float64{}
	dtol := 1e-10
	ver := chk.Verbose
	CheckDerivX(tst, o, tcte, xmin, xmax, np, xskip, sktol, dtol, ver)
}
