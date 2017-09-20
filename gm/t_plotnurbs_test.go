// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func Test_plotnurbs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plotnurbs01. 2D NURBS curve in 2D")

	if chk.Verbose {

		// configurations
		ndim := 2
		npts := 21
		withIds := true
		argsIds := &plt.A{C: "r", Fsz: 10}
		argsIdsC := &plt.A{C: "k", Fsz: 10}

		// curve
		circle := false
		var curve *Nurbs
		if circle {
			xc, yc, r := 0.5, 0.5, 0.5
			curve = FactoryNurbs.Curve2dCircle(xc, yc, r)
		} else {
			curve = FactoryNurbs.Curve2dExample2()
		}

		// plot
		plt.Reset(false, nil)
		curve.DrawCtrl(ndim, false, nil, argsIdsC)
		curve.DrawElems(ndim, npts, withIds, nil, argsIds)
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl", "t_plotnurbs01")
	}
}

func Test_plotnurbs02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plotnurbs02. 2D NURBS surface in 2D")

	if chk.Verbose {

		// configurations
		ndim := 2
		npts := 21
		withIds := true
		argsIds := &plt.A{C: "r", Fsz: 10}
		argsIdsC := &plt.A{C: "k", Fsz: 10}

		// surface
		surf := FactoryNurbs.Surf2dExample1()

		// plot
		plt.Reset(false, nil)
		surf.DrawCtrl(ndim, false, nil, argsIdsC)
		surf.DrawElems(ndim, npts, withIds, nil, argsIds)
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl", "t_plotnurbs02")
	}
}

func Test_plotnurbs03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plotnurbs03. 2D NURBS curve in 3D")

	if chk.Verbose {

		// configurations
		ndim := 3
		npts := 21
		withIds := false

		// curve
		circle := true
		var curve *Nurbs
		if circle {
			xc, yc, r := 0.5, 0.5, 0.5
			curve = FactoryNurbs.Curve2dCircle(xc, yc, r)
		} else {
			curve = FactoryNurbs.Curve2dExample2()
		}

		// plot
		plt.Reset(false, nil)
		curve.DrawCtrl(ndim, withIds, nil, nil)
		curve.DrawElems(ndim, npts, withIds, nil, nil)
		plt.Save("/tmp/gosl", "t_plotnurbs03")
	}
}

func Test_plotnurbs04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plotnurbs04. 2D NURBS surface in 3D")

	if chk.Verbose {

		// configurations
		ndim := 3
		npts := 21
		nu, nv := 21, 11
		withIds := true
		argsIds := &plt.A{C: "r", Fsz: 10}
		argsIdsC := &plt.A{C: "k", Fsz: 10}

		// surface
		surf := FactoryNurbs.Surf2dExample1()

		// plot
		plt.Reset(false, nil)
		surf.DrawCtrl(ndim, false, nil, argsIdsC)
		surf.DrawElems(ndim, npts, withIds, nil, argsIds)
		surf.DrawSurface(ndim, nu, nv, false, true, nil, &plt.A{C: "orange", Lw: 0.7})
		plt.Save("/tmp/gosl", "t_plotnurbs04")
	}
}

func Test_plotnurbs05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plotnurbs05. 3D NURBS curve in 3D")

	if chk.Verbose {

		// curve
		verts := [][]float64{
			{0, 0, 0, 1},
			{1, 0, 0, 1},
			{1, 1, 0, 1},
			{1, 1, 1, 1},
		}
		knots := [][]float64{{0, 0, 0, 0.5, 1, 1, 1}}
		curve := NewNurbs(1, []int{2}, knots)
		err := curve.SetControl(verts, utl.IntRange(len(verts)))
		if err != nil {
			tst.Errorf("%v", err)
			return
		}

		// configurations
		ndim := 3
		npts := 21
		withIds := true

		// plot
		plt.Reset(false, nil)
		plt.Triad(1.1, "x", "y", "z", &plt.A{C: "orange"}, &plt.A{C: "green"})
		plt.Default3dView(-0.1, 1.1, -0.1, 1.1, -0.1, 1.1, true)
		curve.DrawCtrl(ndim, false, nil, nil)
		curve.DrawElems(ndim, npts, withIds, nil, nil)
		plt.Save("/tmp/gosl", "t_plotnurbs05")
	}
}

func Test_plotnurbs06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plotnurbs06. 3D NURBS plane in 3D")

	if chk.Verbose {

		// surface
		verts := [][]float64{
			{0.0, 0.0, 0.0, 1.0},
			{0.5, 0.0, 0.0, 1.0},
			{1.0, 0.0, 0.0, 1.0},
			{0.0, 0.5, 0.5, 1.0},
			{0.5, 0.5, 0.5, 1.0},
			{1.0, 0.5, 0.5, 1.0},
			{0.0, 1.0, 1.0, 1.0},
			{0.5, 1.0, 1.0, 1.0},
			{1.0, 1.0, 1.0, 1.0},
		}
		knots := [][]float64{
			{0, 0, 0, 1, 1, 1},
			{0, 0, 0, 1, 1, 1},
		}
		curve := NewNurbs(2, []int{2, 2}, knots)
		err := curve.SetControl(verts, utl.IntRange(len(verts)))
		if err != nil {
			tst.Errorf("%v", err)
			return
		}

		// configurations
		ndim := 3
		npts := 11
		nu, nv := 21, 21
		withIds := true

		// plot
		plt.Reset(false, nil)
		plt.Triad(1.1, "x", "y", "z", &plt.A{C: "orange"}, &plt.A{C: "green"})
		plt.Scale3d(-0.1, 1.1, -0.1, 1.1, -0.1, 1.1, true)
		curve.DrawCtrl(ndim, false, nil, nil)
		curve.DrawElems(ndim, npts, withIds, nil, nil)
		curve.DrawSurface(3, nu, nv, false, true, nil, &plt.A{C: "#2782c8", Lw: 0.5})
		plt.Save("/tmp/gosl", "t_plotnurbs06")
	}
}

func Test_plotnurbs07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plotnurbs07. 3D NURBS cylinder in 3D")

	if chk.Verbose {

		// surface
		xc, yc, zc, r, h := 0.5, 0.5, 0.0, 0.25, 1.0
		curve := FactoryNurbs.Surf3dCylinder(xc, yc, zc, r, h)

		// configurations
		ndim := 3
		npts := 11
		nu, nv := 21, 21
		withIds := false

		// plot
		plt.Reset(false, nil)
		plt.Triad(1.1, "x", "y", "z", &plt.A{C: "orange"}, &plt.A{C: "green"})
		plt.Default3dView(-0.1, 1.1, -0.1, 1.1, -0.1, 1.1, true)
		curve.DrawCtrl(ndim, false, nil, nil)
		curve.DrawElems(ndim, npts, withIds, nil, nil)
		curve.DrawSurface(3, nu, nv, true, true, nil, &plt.A{C: "#2782c8", Lw: 0.5})
		plt.Save("/tmp/gosl", "t_plotnurbs07")
	}
}

func Test_plotnurbs08(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plotnurbs08. 3D NURBS donut in 3D")

	if chk.Verbose {

		// surface
		xc, yc, zc, r, R := 0.0, 0.0, 0.0, 2.0, 4.0
		curve := FactoryNurbs.Surf3dTorus(xc, yc, zc, r, R)

		// configurations
		ndim := 3
		npts := 11
		nu, nv := 18, 41
		withIds := false

		// plot
		plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
		plt.Triad(1.1, "x", "y", "z", &plt.A{C: "orange"}, &plt.A{C: "green"})
		curve.DrawCtrl(ndim, false, &plt.A{C: "grey", Lw: 0.5}, nil)
		if false {
			curve.DrawElems(ndim, npts, withIds, nil, nil)
		}
		if true {
			curve.DrawSurface(3, nu, nv, true, false, &plt.A{CmapIdx: 3}, &plt.A{C: "#2782c8", Lw: 0.5})
		}
		plt.Default3dView(-6.1, 6.1, -6.1, 6.1, -6.1, 6.1, true)
		if false {
			plt.ShowSave("/tmp/gosl", "t_plotnurbs08")
		} else {
			plt.Save("/tmp/gosl", "t_plotnurbs08")
		}
	}
}
