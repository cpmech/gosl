// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

func Test_args01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("args01")

	var a A

	// plot and basic options
	a.C = "red"
	a.M = "o"
	a.Ls = "--"
	a.Lw = 1.2
	a.Ms = -1
	a.L = "gosl"
	a.Me = 2
	a.Z = 123
	a.Mec = "blue"
	a.Mew = 0.3
	a.Void = true
	a.NoClip = true

	// shapes
	a.Ha = "center"
	a.Va = "center"
	a.Fc = "magenta"
	a.Ec = "yellow"

	// text and extra arguments
	a.Fsz = 7

	l := a.String(false, false)
	chk.String(tst, l, "color='red',marker='o',ls='--',lw=1.2,label='gosl',markevery=2,zorder=123,markeredgecolor='blue',mew=0.3,markerfacecolor='none',clip_on=0,facecolor='magenta',edgecolor='yellow',ha='center',va='center',fontsize=7")
}

func Test_args02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("args02")

	a := &A{
		Colors:  []string{"red", "tan", "lime"},
		Type:    "bar",
		Stacked: true,
		NoFill:  true,
		Nbins:   10,
		Normed:  true,
	}

	l := a.String(true, false)
	chk.String(tst, l, "color=['red','tan','lime'],histtype='bar',stacked=1,fill=0,bins=10,normed=1")
}

func Test_nlevels01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nlevels01")

	nlevels := 3
	Z := [][]float64{
		{0, 1}, {0, 1},
		{0, 1}, {0, 1},
	}

	l := getContourLevels(nlevels, Z)
	chk.String(tst, l, "[0,0.5,1]")
}

func Test_plot01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot01")

	if chk.Verbose {

		x := utl.LinSpace(0.0, 1.0, 11)
		y := make([]float64, len(x))
		for i := 0; i < len(x); i++ {
			y[i] = x[i] * x[i]
		}

		Reset(false, nil)
		SetFontSizes(&A{Fsz: 20, FszLbl: 20, FszXtck: 10, FszYtck: 10})
		Plot(x, y, &A{L: "first", C: "r", M: "o", Ls: "-", Lw: 2, NoClip: true})
		Plot(y, x, &A{L: "second", C: "b", M: ".", Ls: ":", Lw: 40})
		Text(0.2, 0.8, "HERE", &A{Fsz: 20, Ha: "center", Va: "center"})
		SetTicksX(0.1, 0.01, "%.3f")
		SetTicksY(0.2, 0.1, "%.2f")
		HideBorders(&A{HideR: true, HideT: true})
		Gll(`$\varepsilon$`, `$\sigma$`, &A{
			LegOut:  true,
			LegNcol: 3,
			FszLeg:  14,
			HideR:   true,
		})

		err := Save("/tmp/gosl", "t_plot01")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot02")

	if chk.Verbose {

		Reset(true, &A{Eps: true, WidthPt: 380})
		ReplaceAxes(0, 0, 1, 1, 0.04, 0.04, "the x", "the y", &A{Style: "->"}, &A{})
		Arrow(0, 0, 1, 1, &A{C: "orange"})
		AxHline(0, &A{C: "red"})
		AxVline(0, &A{C: "blue"})
		Annotate(0, 0, "TEST", &A{C: "green", FigFraction: true})
		AnnotateXlabels(0, "HERE", &A{Fsz: 10})
		SupTitle("suptitle goes here", &A{C: "red"})
		Title("title goes here", &A{C: "cyan"})
		Text(0.5, 0.5, "TEXT", &A{C: "orange", Va: "top"})
		Cross(0.5, 0.5, nil)
		PlotOne(0, 0, &A{M: "*"})

		err := Save("/tmp/gosl", "t_plot02")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot03")

	if chk.Verbose {

		// grid size
		xmin, xmax, N := -math.Pi/2.0+0.1, math.Pi/2.0-0.1, 21

		// mesh grid
		X, Y, F := utl.MeshGrid2dF(xmin, xmax, xmin, xmax, N, N, func(x, y float64) (z float64) {
			m := math.Pow(math.Cos(x), 2.0) + math.Pow(math.Cos(y), 2.0)
			z = -math.Pow(m, 2.0)
			return
		})

		// configuration
		a := &A{
			NumFmt:  "%.1f",
			Lw:      0.8,
			CbarLbl: "NICE",
			SelectC: "yellow",
			SelectV: -2.5,
			Nlevels: 10,
		}

		Reset(true, nil)
		Equal()
		ContourF(X, Y, F, a)

		err := Save("/tmp/gosl", "t_plot03")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot04")

	if chk.Verbose {

		// grid size
		xmin, xmax, N := -math.Pi/2.0+0.1, math.Pi/2.0-0.1, 21

		// mesh grid
		X, Y, F, U, V := utl.MeshGrid2dFG(xmin, xmax, xmin, xmax, N, N, func(x, y float64) (z, u, v float64) {
			m := math.Pow(math.Cos(x), 2.0) + math.Pow(math.Cos(y), 2.0)
			z = -math.Pow(m, 2.0)
			u = 4.0 * math.Cos(x) * math.Sin(x) * m
			v = 4.0 * math.Cos(y) * math.Sin(y) * m
			return
		})

		// configuration
		a := &A{
			Colors:   []string{"cyan", "blue", "yellow", "green"},
			Levels:   []float64{-4, -3, -2, -1, 0},
			NumFmt:   "%.1f",
			NoLines:  true,
			NoLabels: true,
			NoInline: true,
			NoCbar:   true,
			Lw:       1.5,
			SelectC:  "white",
			SelectV:  -2.5,
		}

		b := &A{
			CmapIdx: 4,
			SelectC: "black",
			SelectV: -2.5,
		}

		Reset(true, nil)
		Equal()
		ContourF(X, Y, F, a)
		ContourL(X, Y, F, b)
		Quiver(X, Y, U, V, nil)

		err := Save("/tmp/gosl", "t_plot04")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot05")

	if chk.Verbose {

		X := [][]float64{
			{1, 1, 1, 2, 2, 2, 2, 2, 3, 3, 4, 5, 6}, // first series
			{-1, -1, 0, 1, 2, 3},                    // second series
			{5, 6, 7, 8},                            // third series
		}

		L := []string{
			"first",
			"second",
			"third",
		}

		a := &A{
			Colors:  []string{"red", "tan", "lime"},
			Ec:      "black",
			Lw:      0.5,
			Type:    "bar",
			Stacked: true,
		}

		Reset(true, nil)
		Hist(X, L, a)
		Gll("series", "count", nil)

		err := Save("/tmp/gosl", "t_plot05")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot06")

	if chk.Verbose {

		x := []float64{0, 1, 1, 1}
		y := []float64{0, 0, 1, 1}
		z := []float64{0, 0, 0, 1}

		np := 5

		X, Y, Z := utl.MeshGrid2dF(0, 1, 0, 1, np, np, func(x, y float64) float64 {
			return x + y
		})

		U, V, W := utl.MeshGrid2dF(0, 1, 0, 1, np, np, func(u, v float64) float64 {
			return u*u + v*v
		})

		Reset(true, nil)

		Triad(1.3, true, nil, nil)

		Plot3dLine(x, y, z, nil)
		Plot3dPoints(x, y, z, nil)
		Surface(U, V, W, &A{CmapIdx: 4, Rstride: 1, Cstride: 1})

		Wireframe(X, Y, Z, &A{C: "orange", Lw: 0.4})

		elev, azim := 30.0, 20.0
		Camera(elev, azim, nil)
		AxDist(10.5)
		Scale3d(0, 1.5, 0, 1.5, 0, 1.5, true)

		//err := ShowSave("/tmp/gosl", "t_plot06")
		err := Save("/tmp/gosl", "t_plot06")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot07")

	if chk.Verbose {

		// point on plane
		p := []float64{0.5, 0.5, 0.5}

		// normal vector
		n := []float64{1, 0, 1}

		// limits and divisions
		xmin, xmax := 0.0, 1.0
		ymin, ymax := 0.0, 1.0
		nu, nv := 5, 5

		// draw
		Reset(true, nil)
		Triad(1.0, true, &A{C: "orange"}, &A{C: "red"})
		PlaneZ(p, n, xmin, xmax, ymin, ymax, nu, nv, true, nil)
		Default3dView(-0.1, 1.1, -0.1, 1.1, -0.1, 1.1, true)

		// save
		//err := ShowSave("/tmp/gosl", "t_plot07")
		err := Save("/tmp/gosl", "t_plot07")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot08(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot08")

	if chk.Verbose {

		// draw
		Reset(true, nil)
		Triad(1.0, true, &A{C: "orange"}, &A{C: "red"})

		alpha, height, radius := 15.0, 1.0, 0.5
		nu, nv := 7, 11
		if true {
			c1 := []float64{0.5, -0.5, 0}
			c2 := []float64{0, 1, 0}
			c3 := []float64{0, 0, 0}
			CylinderZ(c1, radius, height, nu, nv, &A{C: "grey"})
			ConeZ(c2, alpha, height, nu, nv, &A{C: "green"})
			ConeDiag(c3, alpha, height, nu, nv, nil)
			Diag3d(1, nil)
		}

		centre := []float64{0.7, 0.7, 0.0}
		radius, amin, amax := 0.3, 0.0, 180.0
		nu, nv = 21, 5
		cup := false
		Hemisphere(centre, radius, amin, amax, nu, nv, cup, &A{C: "k", Surf: true, Wire: false})

		Default3dView(-0.1, 1.1, -0.1, 1.1, -0.1, 1.1, true)

		// save
		err := Save("/tmp/gosl", "t_plot08")
		//err := ShowSave("/tmp/gosl", "t_plot08")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot09(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot09")

	if chk.Verbose {

		// draw
		Reset(true, nil)
		Triad(1.0, true, &A{C: "orange"}, &A{C: "orange"})

		// centre and radii
		c := []float64{0, 0, 0}
		r := []float64{1, 1, 1}
		a := []float64{2, 2, 2}

		// cup
		//alpMin, alpMax := -180.0, 180.0
		//etaMin, etaMax := -90.0, 0.0

		// hemisphere
		//alpMin, alpMax := -180.0, 180.0
		//etaMin, etaMax := 0.0, 90.0

		// sphere
		alpMin, alpMax := -180.0, 180.0
		etaMin, etaMax := -90.0, 90.0

		// rounded cube
		//a = []float64{10, 10, 10}

		// star
		//a = []float64{0.5, 0.5, 0.5}

		// blob
		a = []float64{2.0, 1.5, 1.0}

		// divisions
		nalp, neta := 30, 30

		// generate
		Superquadric(c, r, a, alpMin, alpMax, etaMin, etaMax, nalp, neta, &A{Surf: true})

		//Default3dView(-2.1, 2.1, -2.1, 2.1, -2.1, 2.1, true)
		Default3dView(-1.1, 1.1, -1.1, 1.1, -1.1, 1.1, true)

		// save
		err := Save("/tmp/gosl", "t_plot09")
		//err := ShowSave("/tmp/gosl", "t_plot09")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}
