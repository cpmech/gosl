// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dbf implements a database of f(t,{x}) functions (e.g. time-space functions). The
// functions in this package are accompanied by derivatives w.r.t t and gradients w.r.t {x}. For
// instance: g(t,{x}) = df/dt, h(t,{x}) = dg/dt, and grad = df/d{x}
package dbf

import (
	"os"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// T defines the interface for t-x functions; i.e. f(t, {x})
type T interface {
	Init(prms Params) error                   // initialise function parameters
	F(t float64, x []float64) float64         // y = F(t, x)
	G(t float64, x []float64) float64         // ∂y/∂t_cteX = G(t, x)
	H(t float64, x []float64) float64         // ∂²y/∂t²_cteX = H(t, x)
	Grad(v []float64, t float64, x []float64) // ∇F = ∂y/∂x = Grad(t, x)
}

// allocators maps function type to function allocator
var allocators = map[string]func() T{} // type => function allocator

// New allocates function by name
func New(name string, prms Params) (T, error) {
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
//  fnkey       -- filename key (without extension)
//  args{F,G,H} -- if any is "", the corresponding plot is not created
func PlotT(o T, dirout, fnkey string, t0, tf float64, xcte []float64, np int) {

	// variables
	t := utl.LinSpace(t0, tf, np)
	var f, g, h []float64
	withF := true
	withG := true
	withH := true
	nrow := 0

	// y-values
	if withF {
		f = make([]float64, np)
		for i := 0; i < np; i++ {
			f[i] = o.F(t[i], xcte)
		}
		nrow++
	}
	if withG {
		g = make([]float64, np)
		for i := 0; i < np; i++ {
			g[i] = o.G(t[i], xcte)
		}
		nrow++
	}
	if withH {
		h = make([]float64, np)
		for i := 0; i < np; i++ {
			h[i] = o.H(t[i], xcte)
		}
		nrow++
	}
	if nrow == 0 {
		chk.Panic("one of args{F,G,H} must be provided")
	}

	// labels
	labelT := "$t$"
	labelX := ""
	for _, x := range xcte {
		labelX += io.Sf(",%g", x)
	}
	labelF := "$f(t" + labelX + ")$"
	labelG := "$g(t" + labelX + ")=\\frac{\\mathrm{d}f}{\\mathrm{d}t}$"
	labelH := "$h(t" + labelX + ")=\\frac{\\mathrm{d}^2f}{\\mathrm{d}t^2}$"

	// plot F
	pidx := 1
	if withF {
		if nrow > 1 {
			plt.Subplot(nrow, 1, pidx)
		}
		plt.Plot(t, f, nil)
		plt.Gll(labelT, labelF, nil)
		pidx++
	}

	// plot G
	if withG {
		if nrow > 1 {
			plt.Subplot(nrow, 1, pidx)
		}
		plt.Plot(t, g, nil)
		plt.Gll(labelT, labelG, nil)
		pidx++
	}

	// plot H
	if withH {
		if nrow > 1 {
			plt.Subplot(nrow, 1, pidx)
		}
		plt.Plot(t, h, nil)
		plt.Gll(labelT, labelH, nil)
	}

	// save figure
	if fnkey != "" {
		plt.Save(dirout, fnkey)
	}
}

// PlotX plots F and the gradient of F, Gx and Gy, for varying x and fixed t
//  hlZero  -- highlight F(t,x) = 0
//  axEqual -- use axis['equal']
func PlotX(o T, dirout, fnkey string, tcte float64, xmin, xmax []float64, np int) {
	withGrad := true
	hlZero := true
	axEqual := true
	if len(xmin) == 3 {
		chk.Panic("PlotX works in 2D only")
	}
	X, Y := utl.MeshGrid2d(xmin[0], xmax[0], xmin[1], xmax[1], np, np)
	F := utl.Alloc(np, np)
	var Gx, Gy [][]float64
	nrow := 1
	if withGrad {
		Gx = utl.Alloc(np, np)
		Gy = utl.Alloc(np, np)
		nrow++
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
	figsize := &plt.A{Prop: 1, WidthPt: 600, Dpi: 200}
	os.MkdirAll(dirout, 0777)
	if withGrad {
		figsize.Prop = 2
		plt.Reset(true, figsize)
		plt.Subplot(nrow, 1, 1)
		plt.Title("F(t,x)", nil)
	} else {
		plt.Reset(true, figsize)
	}
	plt.ContourF(X, Y, F, nil)
	if hlZero {
		plt.ContourL(X, Y, F, &plt.A{Levels: []float64{0}, Lw: 2, Colors: []string{"yellow"}})
	}
	if axEqual {
		plt.Equal()
	}
	plt.Gll("x", "y", nil)
	if withGrad {
		plt.Subplot(2, 1, 2)
		plt.Title("gradient", nil)
		plt.Quiver(X, Y, Gx, Gy, nil)
		if axEqual {
			plt.Equal()
		}
		plt.Gll("x", "y", nil)
	}
	if fnkey != "" {
		plt.Save(dirout, fnkey)
	}
}
