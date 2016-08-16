// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// NURBS methods ///////////////////////////////////////////////////////////////////////////////////

// DrawCtrl2d draws control net
func (o *Nurbs) DrawCtrl2d(ids bool, args, idargs string) {
	if len(idargs) == 0 {
		idargs = "color='r', size=7"
	}
	switch o.gnd {
	// curve
	case 1:
		xa := make([]float64, o.n[0])
		ya := make([]float64, o.n[0])
		j, k := 0, 0
		for i := 0; i < o.n[0]; i++ {
			xa[i] = o.Q[i][j][k][0] / o.Q[i][j][k][3]
			ya[i] = o.Q[i][j][k][1] / o.Q[i][j][k][3]
		}
		plt.Plot(xa, ya, "'k.--', clip_on=0"+args)
		if ids {
			for i := 0; i < o.n[0]; i++ {
				x := o.Q[i][j][k][0] / o.Q[i][j][k][3]
				y := o.Q[i][j][k][1] / o.Q[i][j][k][3]
				plt.Text(x, y, io.Sf("%d", i), idargs)
			}
		}
	// surface
	case 2:
		xa := make([]float64, o.n[1])
		ya := make([]float64, o.n[1])
		k := 0
		for i := 0; i < o.n[0]; i++ {
			for j := 0; j < o.n[1]; j++ {
				xa[j] = o.Q[i][j][k][0] / o.Q[i][j][k][3]
				ya[j] = o.Q[i][j][k][1] / o.Q[i][j][k][3]
			}
			plt.Plot(xa, ya, "'k.--', clip_on=0"+args)
		}
		xb := make([]float64, o.n[0])
		yb := make([]float64, o.n[0])
		for j := 0; j < o.n[1]; j++ {
			for i := 0; i < o.n[0]; i++ {
				xb[i] = o.Q[i][j][k][0] / o.Q[i][j][k][3]
				yb[i] = o.Q[i][j][k][1] / o.Q[i][j][k][3]
			}
			plt.Plot(xb, yb, "'k.--', clip_on=0"+args)
		}
		if ids {
			for i := 0; i < o.n[0]; i++ {
				for j := 0; j < o.n[1]; j++ {
					x := o.Q[i][j][k][0] / o.Q[i][j][k][3]
					y := o.Q[i][j][k][1] / o.Q[i][j][k][3]
					l := i + j*o.n[0]
					plt.Text(x, y, io.Sf("%d", l), idargs)
				}
			}
		}
	}
}

// DrawElem2d draws element (non-zero span)
func (o *Nurbs) DrawElem2d(span []int, npts int, ids bool, args, idargs string) {
	if len(idargs) == 0 {
		idargs = "color='b', va='top', size=7"
	}
	switch o.gnd {
	// curve
	case 1:
		umin, umax := o.b[0].T[span[0]], o.b[0].T[span[1]]
		o.DrawEdge2d(umin, umax, 0.0, 0, npts, args)
		if ids {
			c := o.Point([]float64{umin})
			plt.Text(c[0], c[1], io.Sf("(%d)", span[0]), idargs)
			if span[1] == o.b[0].m-o.p[0]-1 {
				c := o.Point([]float64{umax})
				plt.Text(c[0], c[1], io.Sf("(%d)", span[1]), idargs)
			}
		}
	// surface
	case 2:
		umin, umax := o.b[0].T[span[0]], o.b[0].T[span[1]]
		vmin, vmax := o.b[1].T[span[2]], o.b[1].T[span[3]]
		o.DrawEdge2d(umin, umax, vmin, 0, npts, args)
		o.DrawEdge2d(umin, umax, vmax, 0, npts, args)
		o.DrawEdge2d(vmin, vmax, umin, 1, npts, args)
		o.DrawEdge2d(vmin, vmax, umax, 1, npts, args)
		if ids {
			c := o.Point([]float64{umin, vmin})
			plt.Text(c[0], c[1], io.Sf("(%d,%d)", span[0], span[2]), idargs)
			c = o.Point([]float64{umin, vmax})
			plt.Text(c[0], c[1], io.Sf("(%d,%d)", span[0], span[3]), idargs)
			c = o.Point([]float64{umax, vmin})
			plt.Text(c[0], c[1], io.Sf("(%d,%d)", span[1], span[2]), idargs)
			c = o.Point([]float64{umax, vmax})
			plt.Text(c[0], c[1], io.Sf("(%d,%d)", span[1], span[3]), idargs)
		}
	}
}

// DrawElems2d draws all elements (non-zero spans)
func (o *Nurbs) DrawElems2d(npts int, ids bool, args, idargs string) {
	elems := o.Elements()
	for _, e := range elems {
		o.DrawElem2d(e, npts, ids, args, idargs)
	}
}

// PlotBasis plots basis function (2D only)
// option =  0 : use CalcBasis
//           1 : use CalcBasisAndDerivs
//           2 : use RecursiveBasis
func (o *Nurbs) PlotBasis(l int, args string, npts, option int) {
	lbls := []string{"CalcBasis function", "CalcBasisAndDerivs function", "RecursiveBasis function"}
	switch o.gnd {
	// curve
	case 1:
		xx := make([]float64, npts)
		yy := make([]float64, npts)
		du0 := (o.b[0].tmax - o.b[0].tmin) / float64(npts-1)
		for m := 0; m < npts; m++ {
			u0 := o.b[0].tmin + float64(m)*du0
			u := []float64{u0}
			x := o.Point(u)
			xx[m] = x[0]
			switch option {
			case 0:
				o.CalcBasis(u)
				yy[m] = o.GetBasisL(l)
			case 1:
				o.CalcBasisAndDerivs(u)
				yy[m] = o.GetBasisL(l)
			case 2:
				yy[m] = o.RecursiveBasis(u, l)
			}
		}
		plt.Plot(xx, yy, args)
	// surface
	case 2:
		xx := la.MatAlloc(npts, npts)
		yy := la.MatAlloc(npts, npts)
		zz := la.MatAlloc(npts, npts)
		du0 := (o.b[0].tmax - o.b[0].tmin) / float64(npts-1)
		du1 := (o.b[1].tmax - o.b[1].tmin) / float64(npts-1)
		for m := 0; m < npts; m++ {
			u0 := o.b[0].tmin + float64(m)*du0
			for n := 0; n < npts; n++ {
				u1 := o.b[1].tmin + float64(n)*du1
				u := []float64{u0, u1}
				x := o.Point(u)
				xx[m][n] = x[0]
				yy[m][n] = x[1]
				switch option {
				case 0:
					o.CalcBasis(u)
					zz[m][n] = o.GetBasisL(l)
				case 1:
					o.CalcBasisAndDerivs(u)
					zz[m][n] = o.GetBasisL(l)
				case 2:
					zz[m][n] = o.RecursiveBasis(u, l)
				}
			}
		}
		plt.Contour(xx, yy, zz, "fsz=7")
	}
	plt.Title(io.Sf("%s @ %d", lbls[option], l), "size=7")
}

// PlotDeriv plots derivative dR[i][j][k]du[d] (2D only)
// option =  0 : use CalcBasisAndDerivs
//           1 : use NumericalDeriv
func (o *Nurbs) PlotDeriv(l, d int, args string, npts, option int) {
	lbls := []string{"CalcBasisAndDerivs function", "NumericalDeriv function"}
	switch o.gnd {
	// curve
	case 1:
	// surface
	case 2:
		xx := la.MatAlloc(npts, npts)
		yy := la.MatAlloc(npts, npts)
		zz := la.MatAlloc(npts, npts)
		du0 := (o.b[0].tmax - o.b[0].tmin) / float64(npts-1)
		du1 := (o.b[1].tmax - o.b[1].tmin) / float64(npts-1)
		drdu := make([]float64, 2)
		for m := 0; m < npts; m++ {
			u0 := o.b[0].tmin + float64(m)*du0
			for n := 0; n < npts; n++ {
				u1 := o.b[1].tmin + float64(n)*du1
				u := []float64{u0, u1}
				x := o.Point(u)
				xx[m][n] = x[0]
				yy[m][n] = x[1]
				switch option {
				case 0:
					o.CalcBasisAndDerivs(u)
					o.GetDerivL(drdu, l)
				case 1:
					o.NumericalDeriv(drdu, u, l)
				}
				zz[m][n] = drdu[d]
			}
		}
		plt.Title(io.Sf("%s @ %d,%d", lbls[option], l, d), "size=7")
		plt.Contour(xx, yy, zz, "fsz=7")
	}
}

// DrawEdge2d draws and edge from tmin to tmax
func (o *Nurbs) DrawEdge2d(tmin, tmax, cte float64, along, npts int, args string) {
	tt := utl.LinSpace(tmin, tmax, npts)
	xx := make([]float64, npts)
	yy := make([]float64, npts)
	u := make([]float64, 2)
	for i, t := range tt {
		if along == 0 {
			u[0], u[1] = t, cte
		} else {
			u[0], u[1] = cte, t
		}
		x := o.Point(u)
		xx[i], yy[i] = x[0], x[1]
	}
	plt.Plot(xx, yy, "'k-', clip_on=0"+args)
}

// global functions ////////////////////////////////////////////////////////////////////////////////

// PlotNurbs plots a NURBS
func PlotNurbs(dirout, fn string, b *Nurbs, npts int, ids bool, extra func()) {
	plt.Reset()
	if io.FnExt(fn) == ".eps" {
		plt.SetForEps(1.0, 500)
	} else {
		plt.SetForPng(1.0, 500, 150)
	}
	b.DrawCtrl2d(ids, "", "")
	b.DrawElems2d(npts, ids, "", "")
	if extra != nil {
		extra()
	}
	plt.Equal()
	plt.SaveD(dirout, fn)
}

// PlotTwoNurbs plots two NURBS for comparison
func PlotTwoNurbs(dirout, fn string, b, c *Nurbs, npts int, ids bool, extra func()) {
	plt.Reset()
	if io.FnExt(fn) == ".eps" {
		plt.SetForEps(1.5, 500)
	} else {
		plt.SetForPng(1.5, 500, 150)
	}

	plt.Subplot(3, 1, 1)
	b.DrawCtrl2d(ids, "", "")
	b.DrawElems2d(npts, ids, "", "")
	if extra != nil {
		extra()
	}
	plt.Equal()

	plt.Subplot(3, 1, 2)
	c.DrawCtrl2d(ids, "", "")
	c.DrawElems2d(npts, ids, "", "")
	plt.Equal()

	plt.Subplot(3, 1, 3)
	b.DrawElems2d(npts, ids, ", lw=3", "")
	c.DrawElems2d(npts, ids, ", color='red', marker='+', markevery=10", "color='green', size=7, va='bottom'")
	plt.Equal()

	plt.SaveD(dirout, fn)
}

// PlotNurbsBasis plots basis functions la and lb
func PlotNurbsBasis(dirout, fn string, b *Nurbs, la, lb int) {
	npts := 41
	plt.Reset()
	if io.FnExt(fn) == ".eps" {
		plt.SetForEps(1.5, 500)
	} else {
		plt.SetForPng(1.5, 600, 150)
	}

	plt.Subplot(3, 2, 1)
	b.DrawCtrl2d(false, "", "")
	b.DrawElems2d(npts, false, "", "")
	t0 := time.Now()
	b.PlotBasis(la, "", 11, 0) // 0 => CalcBasis
	io.Pfcyan("time elapsed (calcbasis) = %v\n", time.Now().Sub(t0))
	plt.Equal()

	plt.Subplot(3, 2, 2)
	b.DrawCtrl2d(false, "", "")
	b.DrawElems2d(npts, false, "", "")
	b.PlotBasis(lb, "", 11, 0) // 0 => CalcBasis
	plt.Equal()

	plt.Subplot(3, 2, 3)
	b.DrawCtrl2d(false, "", "")
	b.DrawElems2d(npts, false, "", "")
	b.PlotBasis(la, "", 11, 1) // 1 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(3, 2, 4)
	b.DrawCtrl2d(false, "", "")
	b.DrawElems2d(npts, false, "", "")
	b.PlotBasis(lb, "", 11, 1) // 1 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(3, 2, 5)
	b.DrawCtrl2d(false, "", "")
	b.DrawElems2d(npts, false, "", "")
	t0 = time.Now()
	b.PlotBasis(la, "", 11, 2) // 2 => RecursiveBasis
	io.Pfcyan("time elapsed (recursive) = %v\n", time.Now().Sub(t0))
	plt.Equal()

	plt.Subplot(3, 2, 6)
	b.DrawCtrl2d(false, "", "")
	b.DrawElems2d(npts, false, "", "")
	b.PlotBasis(lb, "", 11, 2) // 2 => RecursiveBasis
	plt.Equal()

	plt.SaveD(dirout, fn)
}

// PlotNurbsDerivs plots derivatives of basis functions la and lb
func PlotNurbsDerivs(dirout, fn string, b *Nurbs, la, lb int) {
	npts := 41
	plt.Reset()
	if io.FnExt(fn) == ".eps" {
		plt.SetForEps(1.5, 500)
	} else {
		plt.SetForPng(1.5, 600, 150)
	}

	plt.Subplot(4, 2, 1)
	t0 := time.Now()
	b.PlotDeriv(la, 0, "", npts, 0) // 0 => CalcBasisAndDerivs
	io.Pfcyan("time elapsed (calcbasis) = %v\n", time.Now().Sub(t0))
	plt.Equal()

	plt.Subplot(4, 2, 2)
	t0 = time.Now()
	b.PlotDeriv(la, 0, "", npts, 1) // 1 => NumericalDeriv
	io.Pfcyan("time elapsed (numerical) = %v\n", time.Now().Sub(t0))
	plt.Equal()

	plt.Subplot(4, 2, 3)
	b.PlotDeriv(la, 1, "", npts, 0) // 0 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(4, 2, 4)
	b.PlotDeriv(la, 1, "", npts, 1) // 0 => NumericalDeriv
	plt.Equal()

	plt.Subplot(4, 2, 5)
	b.PlotDeriv(lb, 0, "", npts, 0) // 0 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(4, 2, 6)
	b.PlotDeriv(lb, 0, "", npts, 1) // 0 => NumericalDeriv
	plt.Equal()

	plt.Subplot(4, 2, 7)
	b.PlotDeriv(lb, 1, "", npts, 0) // 0 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(4, 2, 8)
	b.PlotDeriv(lb, 1, "", npts, 1) // 0 => NumericalDeriv
	plt.Equal()

	plt.SaveD(dirout, fn)
}
