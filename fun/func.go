// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// fun (functions) implements a number of y(t,x) functions and their first and second order derivatives
package fun

import (
	"os"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// Func defines the interface for y(t, x) functions
type Func interface {
	Init(prms Prms) error                     // initialise function parameters
	F(t float64, x []float64) float64         // y = F(t, x)
	G(t float64, x []float64) float64         // ∂y/∂t_cteX = G(t, x)
	H(t float64, x []float64) float64         // ∂²y/∂t²_cteX = H(t, x)
	Grad(v []float64, t float64, x []float64) // ∇F = ∂y/∂x = Grad(t, x)
}

// allocators maps function type to function allocator
var allocators = map[string]func() Func{} // type => function allocator

// New allocates function by name
func New(name string, prms Prms) (Func, error) {
	if name == "zero" {
		return &Zero, nil
	}
	allocator, ok := allocators[name]
	if !ok {
		return nil, chk.Err("cannot find function named %q", name)
	}
	o := allocator()
	err := o.Init(prms)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// PlotT plots F, G and H for varying t and fixed coordinates x
func PlotT(o Func, dirout, fname string, t0, tf float64, xcte []float64, np int, args string, withG, withH, save, show bool, extra func()) {
	t := utl.LinSpace(t0, tf, np)
	y := make([]float64, np)
	for i := 0; i < np; i++ {
		y[i] = o.F(t[i], xcte)
	}
	var g, h []float64
	nrow := 1
	if withG {
		g = make([]float64, np)
		for i := 0; i < np; i++ {
			g[i] = o.G(t[i], xcte)
		}
		nrow += 1
	}
	if withH {
		h = make([]float64, np)
		for i := 0; i < np; i++ {
			h[i] = o.H(t[i], xcte)
		}
		nrow += 1
	}
	os.MkdirAll(dirout, 0777)
	if withG || withH {
		plt.Subplot(nrow, 1, 1)
	}
	plt.Plot(t, y, args)
	if extra != nil {
		extra()
	}
	plt.Gll("t", "y", "")
	pidx := 2
	if withG {
		plt.Subplot(nrow, 1, pidx)
		plt.Plot(t, g, args)
		plt.Gll("t", "dy/dt", "")
		pidx += 1
	}
	if withH {
		plt.Subplot(nrow, 1, pidx)
		plt.Plot(t, h, args)
		plt.Gll("t", "d2y/dt2", "")
	}
	if save {
		plt.Save(dirout + "/" + fname)
	}
	if show {
		plt.Show()
	}
}

// PlotX plots F and the gradient of F, Gx and Gy, for varying x and fixed t
//  hlZero  -- highlight F(t,x) = 0
//  axEqual -- use axis['equal']
func PlotX(o Func, dirout, fname string, tcte float64, xmin, xmax []float64, np int, args string, withGrad, hlZero, axEqual, save, show bool, extra func()) {
	if len(xmin) == 3 {
		chk.Panic("PlotX works in 2D only")
	}
	X, Y := utl.MeshGrid2D(xmin[0], xmax[0], xmin[1], xmax[1], np, np)
	F := la.MatAlloc(np, np)
	var Gx, Gy [][]float64
	nrow := 1
	if withGrad {
		Gx = la.MatAlloc(np, np)
		Gy = la.MatAlloc(np, np)
		nrow += 1
	}
	x := make([]float64, 2)
	g := make([]float64, 2)
	for i := 0; i < np; i++ {
		for j := 0; j < np; j++ {
			x[0], x[1] = X[i][j], Y[i][j]
			F[i][j] = o.F(tcte, x)
			if withGrad {
				o.Grad(g, tcte, x)
				Gx[i][j] = g[0]
				Gy[i][j] = g[1]
			}
		}
	}
	prop, wid, dpi := 1.0, 600.0, 200
	os.MkdirAll(dirout, 0777)
	if withGrad {
		prop = 2
		plt.SetForPng(prop, wid, dpi)
		plt.Subplot(nrow, 1, 1)
		plt.Title("F(t,x)", "")
	} else {
		plt.SetForPng(prop, wid, dpi)
	}
	plt.Contour(X, Y, F, args)
	if hlZero {
		plt.ContourSimple(X, Y, F, "levels=[0], linewidths=[2], colors=['yellow']")
	}
	if axEqual {
		plt.Equal()
	}
	if extra != nil {
		extra()
	}
	plt.Gll("x", "y", "")
	if withGrad {
		plt.Subplot(2, 1, 2)
		plt.Title("gradient", "")
		plt.Quiver(X, Y, Gx, Gy, args)
		if axEqual {
			plt.Equal()
		}
		plt.Gll("x", "y", "")
	}
	if save {
		plt.Save(dirout + "/" + fname)
	}
	if show {
		plt.Show()
	}
}
