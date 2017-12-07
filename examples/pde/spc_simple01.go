// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/pde"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

const (
	π    = math.Pi
	nsum = 100
)

func ana(x []float64) float64 {
	var sum float64
	for j := 0; j < nsum; j++ {
		c := 2.0*float64(j) + 1.0
		cp := c * π
		op := fun.NegOnePowN(j)
		sum += op * math.Sin(cp*x[0]) * math.Sinh(cp*(1.0-x[1])) / (c * c * math.Sinh(cp))
	}
	return 4.0 * sum / (π * π)
}

func run(N int, doplot bool) (maxerr float64) {

	// Lagrange interpolators
	lis := fun.NewLagIntSet(2, []int{N, N}, []string{"cgl", "cgl"})

	// grid
	g := new(gm.Grid)
	g.RectSet2dU([]float64{0, 0}, []float64{1, 1}, lis[0].X, lis[1].X)

	// solver
	p := dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}
	s := pde.NewSpcLaplacian(p, lis, g, nil)

	// essential boundary conditions
	s.AddEbc(10, 0.0, nil) // left
	s.AddEbc(11, 0.0, nil) // right
	s.AddEbc(20, 0.0, func(x la.Vector, t float64) float64 {
		if x[0] < 0.5 {
			return x[0]
		}
		return 1.0 - x[0]
	}) // bottom
	s.AddEbc(21, 0.0, nil) // top

	// solve
	reactions := false
	s.Assemble(reactions)
	u, _ := s.SolveSteady(reactions)

	// error
	e := make([]float64, len(u))
	for n := 0; n < g.Npts(1); n++ {
		for m := 0; m < g.Npts(0); m++ {
			I := g.IndexMNPtoI(m, n, 0)
			x := g.X(m, n, 0)
			e[I] = math.Abs(u[I] - ana(x))
			maxerr = utl.Max(maxerr, e[I])
		}
	}

	// plot
	if doplot {
		uu := g.MapMeshgrid2d(u)
		ee := g.MapMeshgrid2d(e)
		gp := gm.GridPlotter{G: g, WithVids: false, ArgsEdges: &plt.A{C: "grey", Lw: 0.5}}

		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		gp.Draw()
		plt.ContourF(gp.X2d, gp.Y2d, uu, nil)
		plt.Equal()
		plt.Gll("$x$", "$y$", nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/pde", io.Sf("spc_simple01_contour_N%d", N))

		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		gp.Draw()
		plt.ContourF(gp.X2d, gp.Y2d, ee, nil)
		plt.Equal()
		plt.Gll("$x$", "$y$", nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/pde", io.Sf("spc_simple01_error_N%d", N))

		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.Surface(gp.X2d, gp.Y2d, uu, &plt.A{})
		plt.Default3dView(0, 1, 0, 1, 0, 1, true)
		plt.Save("/tmp/gosl/pde", io.Sf("spc_simple01_surface_N%d", N))
	}
	return
}

func main() {

	run(10, true)

	Nvalues := []int{10, 16, 22, 26, 32, 64}
	nn := make([]float64, len(Nvalues))
	ee := make([]float64, len(Nvalues))
	for i, N := range Nvalues {
		io.Pforan("running with N = %v\n", N)
		nn[i] = float64(N)
		ee[i] = run(N, false)
	}
	io.Pf("nn = %v\n", nn)
	io.Pf("ee = %v\n", ee)

	plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
	plt.Plot(nn, ee, &plt.A{C: plt.C(0, 0), Lw: 1.5})
	plt.SetYlog()
	plt.Gll("$N$", "$max(error)$", nil)
	plt.Save("/tmp/gosl/pde", "spc_simple01_error")
}
