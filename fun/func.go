// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// fun (functions) implements a number of y(t,x) functions and their first and second order derivatives
package fun

import (
	"os"

	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// Func defines the interface for y(t, x) functions
type Func interface {
	Init(prms Prms)
	F(t float64, x []float64) float64         // y = F(t, x)
	G(t float64, x []float64) float64         // ∂y/∂t_cteX = G(t, x)
	H(t float64, x []float64) float64         // ∂²y/∂t²_cteX = H(t, x)
	Grad(v []float64, t float64, x []float64) // ∇F = ∂y/∂x = Grad(t, x)
}

// allocators maps function type to function allocator
var allocators = map[string]func() Func{} // type => function allocator

// New allocates function by name
func New(name string, prms Prms) Func {
	if name == "zero" {
		return &Zero
	}
	allocator, ok := allocators[name]
	if !ok {
		utl.Panic("cannot find function named %s", name)
	}
	o := allocator()
	o.Init(prms)
	return o
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
