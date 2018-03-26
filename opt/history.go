// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// History holds history of optmization using directiors; e.g. for Debugging
type History struct {
	Ndim  int         // dimension of x-vector
	HistX []la.Vector // [it] history of x-values (position)
	HistN []la.Vector // [it] history of n-values (direction)
	HistF []float64   // [it] history of f-values
	HistI []float64   // [it] index of iteration
	ffcn  fun.Sv      // f({x}) function
}

// NewHistory returns new object
func NewHistory(nMaxIt int, f0 float64, x0 la.Vector, ffcn fun.Sv) (o *History) {
	o = new(History)
	o.Ndim = len(x0)
	o.HistX = append(o.HistX, x0.GetCopy())
	o.HistN = append(o.HistN, nil)
	o.HistF = append(o.HistF, f0)
	o.HistI = append(o.HistI, 0)
	o.ffcn = ffcn
	return
}

// Append appends new x and n vectors, and updates F and I arrays
func (o *History) Append(fx float64, x, n la.Vector) {
	o.HistX = append(o.HistX, x.GetCopy())
	o.HistN = append(o.HistN, n.GetCopy())
	o.HistF = append(o.HistF, fx)
	o.HistI = append(o.HistI, float64(len(o.HistI)))
}

// Plot plots history
func (o *History) Plot(iDim, jDim int, ximin, ximax, xjmin, xjmax float64, npts int) {

	// contour
	plt.Subplot(2, 1, 1)
	xvec := la.NewVector(2)
	xx, yy, zz := utl.MeshGrid2dF(ximin, ximax, xjmin, xjmax, npts, npts, func(r, s float64) float64 {
		xvec[0], xvec[1] = r, s
		return o.ffcn(xvec)
	})
	plt.ContourF(xx, yy, zz, nil)

	// trajectory
	for k := 0; k < len(o.HistX)-1; k++ {
		x := o.HistX[k]
		n := o.HistN[1+k]
		if n.Norm() > 1e-10 {
			plt.PlotOne(x[0], x[1], &plt.A{C: "y", M: "o", Z: 10, NoClip: true})
			plt.DrawArrow2d(x, n, true, 1, nil)
		}
	}

	// final point
	l := len(o.HistX) - 1
	plt.PlotOne(o.HistX[l][0], o.HistX[l][1], &plt.A{C: "y", M: "*", Ms: 10, Z: 10, NoClip: true})

	// convergence
	plt.Subplot(2, 1, 2)
	plt.Plot(o.HistI, o.HistF, &plt.A{C: plt.C(2, 0), M: ".", Ls: "-", Lw: 2, NoClip: true})
	plt.Text(o.HistI[0], o.HistF[0], io.Sf("%.3f", o.HistF[0]), &plt.A{C: plt.C(0, 0), Fsz: 7, Ha: "left", Va: "top", NoClip: true})
	plt.Text(o.HistI[l], o.HistF[l], io.Sf("%.3f", o.HistF[l]), &plt.A{C: plt.C(0, 0), Fsz: 7, Ha: "right", Va: "bottom", NoClip: true})
	plt.Gll("$iteration$", "$f(x)$", nil)
	plt.HideTRborders()
}
