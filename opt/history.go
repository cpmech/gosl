// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"math"

	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// History holds history of optmization using directiors; e.g. for Debugging
type History struct {

	// data
	Ndim  int         // dimension of x-vector
	HistX []la.Vector // [it] history of x-values (position)
	HistU []la.Vector // [it] history of u-values (direction)
	HistF []float64   // [it] history of f-values
	HistI []float64   // [it] index of iteration

	// configuration
	NptsI   int       // number of points for contour
	NptsJ   int       // number of points for contour
	RangeXi []float64 // {ximin, ximax} [may be nil for default]
	RangeXj []float64 // {xjmin, xjmax} [may be nil for default]
	GapXi   float64   // expand {ximin, ximax}
	GapXj   float64   // expand {ximin, ximax}

	// internal
	ffcn fun.Sv // f({x}) function
}

// NewHistory returns new object
func NewHistory(nMaxIt int, f0 float64, x0 la.Vector, ffcn fun.Sv) (o *History) {
	o = new(History)
	o.Ndim = len(x0)
	o.HistX = append(o.HistX, x0.GetCopy())
	o.HistU = append(o.HistU, nil)
	o.HistF = append(o.HistF, f0)
	o.HistI = append(o.HistI, 0)
	o.NptsI = 41
	o.NptsJ = 41
	o.ffcn = ffcn
	return
}

// Append appends new x and u vectors, and updates F and I arrays
func (o *History) Append(fx float64, x, u la.Vector) {
	o.HistX = append(o.HistX, x.GetCopy())
	o.HistU = append(o.HistU, u.GetCopy())
	o.HistF = append(o.HistF, fx)
	o.HistI = append(o.HistI, float64(len(o.HistI)))
}

// PlotX plots contour and trajectory of x-points
func (o *History) PlotX(iDim, jDim int, xref la.Vector) {

	// i-range
	var ximin, ximax float64
	if len(o.RangeXi) == 2 {
		ximin, ximax = o.RangeXi[0], o.RangeXi[1]
	} else {
		ximin, ximax = math.MaxFloat64, math.SmallestNonzeroFloat64
		for _, x := range o.HistX {
			ximin = utl.Min(ximin, x[iDim])
			ximax = utl.Max(ximax, x[iDim])
		}
	}

	// j-range
	var xjmin, xjmax float64
	if len(o.RangeXj) == 2 {
		xjmin, xjmax = o.RangeXj[0], o.RangeXj[1]
	} else {
		xjmin, xjmax = math.MaxFloat64, math.SmallestNonzeroFloat64
		for _, x := range o.HistX {
			xjmin = utl.Min(xjmin, x[jDim])
			xjmax = utl.Max(xjmax, x[jDim])
		}
	}

	// use gap
	ximin -= o.GapXi
	ximax += o.GapXi
	xjmin -= o.GapXj
	xjmax += o.GapXj

	// contour
	xvec := xref.GetCopy()
	xx, yy, zz := utl.MeshGrid2dF(ximin, ximax, xjmin, xjmax, o.NptsI, o.NptsJ, func(r, s float64) float64 {
		xvec[iDim], xvec[jDim] = r, s
		return o.ffcn(xvec)
	})
	plt.ContourF(xx, yy, zz, nil)

	// trajectory
	x2d := la.NewVector(2)
	u2d := la.NewVector(2)
	for k := 0; k < len(o.HistX)-1; k++ {
		x := o.HistX[k]
		u := o.HistU[1+k]
		x2d[0], x2d[1] = x[iDim], x[jDim]
		u2d[0], u2d[1] = u[iDim], u[jDim]
		if u.Norm() > 1e-10 {
			plt.PlotOne(x2d[0], x2d[1], &plt.A{C: "y", M: "o", Z: 10, NoClip: true})
			plt.DrawArrow2d(x2d, u2d, false, 1, nil)
		}
	}

	// final point
	l := len(o.HistX) - 1
	plt.PlotOne(o.HistX[l][iDim], o.HistX[l][jDim], &plt.A{C: "y", M: "*", Ms: 10, Z: 10, NoClip: true})

	// labels
	plt.SetLabels(io.Sf("$x_{%d}$", iDim), io.Sf("$x_{%d}$", jDim), nil)
}

// PlotF plots f-values along iterations
func (o *History) PlotF() {
	l := len(o.HistI) - 1
	plt.Plot(o.HistI, o.HistF, &plt.A{C: plt.C(2, 0), M: ".", Ls: "-", Lw: 2, NoClip: true})
	plt.Text(o.HistI[0], o.HistF[0], io.Sf("%.3f", o.HistF[0]), &plt.A{C: plt.C(0, 0), Fsz: 7, Ha: "left", Va: "top", NoClip: true})
	plt.Text(o.HistI[l], o.HistF[l], io.Sf("%.3f", o.HistF[l]), &plt.A{C: plt.C(0, 0), Fsz: 7, Ha: "right", Va: "bottom", NoClip: true})
	plt.Gll("$iteration$", "$f(x)$", nil)
	plt.HideTRborders()
}
