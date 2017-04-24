// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// NURBS methods ///////////////////////////////////////////////////////////////////////////////////

// DrawCtrl2d draws control net
func (o *Nurbs) DrawCtrl2d(withIds bool, args, argsIds *plt.A) {
	if args == nil {
		args = &plt.A{C: "k", M: ".", Ls: "--", NoClip: true}
	}
	if argsIds == nil {
		argsIds = &plt.A{C: "r", Fsz: 7, NoClip: true}
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
		plt.Plot(xa, ya, args)
		if withIds {
			for i := 0; i < o.n[0]; i++ {
				x := o.Q[i][j][k][0] / o.Q[i][j][k][3]
				y := o.Q[i][j][k][1] / o.Q[i][j][k][3]
				plt.Text(x, y, io.Sf("%d", i), argsIds)
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
			plt.Plot(xa, ya, args)
		}
		xb := make([]float64, o.n[0])
		yb := make([]float64, o.n[0])
		for j := 0; j < o.n[1]; j++ {
			for i := 0; i < o.n[0]; i++ {
				xb[i] = o.Q[i][j][k][0] / o.Q[i][j][k][3]
				yb[i] = o.Q[i][j][k][1] / o.Q[i][j][k][3]
			}
			plt.Plot(xb, yb, args)
		}
		if withIds {
			for i := 0; i < o.n[0]; i++ {
				for j := 0; j < o.n[1]; j++ {
					x := o.Q[i][j][k][0] / o.Q[i][j][k][3]
					y := o.Q[i][j][k][1] / o.Q[i][j][k][3]
					l := i + j*o.n[0]
					plt.Text(x, y, io.Sf("%d", l), argsIds)
				}
			}
		}
	}
}

// DrawEdge2d draws and edge from tmin to tmax
func (o *Nurbs) DrawEdge2d(tmin, tmax, cte float64, along, npts int, args *plt.A) {
	if args == nil {
		args = &plt.A{C: "b", Ls: "-", NoClip: true}
	}
	tt := utl.LinSpace(tmin, tmax, npts)
	xx := make([]float64, npts)
	yy := make([]float64, npts)
	x := make([]float64, 2)
	u := make([]float64, 2)
	for i, t := range tt {
		if along == 0 {
			u[0], u[1] = t, cte
		} else {
			u[0], u[1] = cte, t
		}
		o.Point(x, u, 2)
		xx[i], yy[i] = x[0], x[1]
	}
	plt.Plot(xx, yy, args)
}

// DrawElem2d draws element (non-zero span)
func (o *Nurbs) DrawElem2d(span []int, npts int, withIds bool, args, argsIds *plt.A) {
	if argsIds == nil {
		argsIds = &plt.A{C: "r", Fsz: 7}
	}
	c := make([]float64, 2)
	switch o.gnd {
	// curve
	case 1:
		umin, umax := o.b[0].T[span[0]], o.b[0].T[span[1]]
		o.DrawEdge2d(umin, umax, 0.0, 0, npts, args)
		if withIds {
			o.Point(c, []float64{umin}, 2)
			plt.Text(c[0], c[1], io.Sf("(%d)", span[0]), argsIds)
			if span[1] == o.b[0].m-o.p[0]-1 {
				o.Point(c, []float64{umax}, 2)
				plt.Text(c[0], c[1], io.Sf("(%d)", span[1]), argsIds)
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
		if withIds {
			o.Point(c, []float64{umin, vmin}, 2)
			plt.Text(c[0], c[1], io.Sf("(%d,%d)", span[0], span[2]), argsIds)
			o.Point(c, []float64{umin, vmax}, 2)
			plt.Text(c[0], c[1], io.Sf("(%d,%d)", span[0], span[3]), argsIds)
			o.Point(c, []float64{umax, vmin}, 2)
			plt.Text(c[0], c[1], io.Sf("(%d,%d)", span[1], span[2]), argsIds)
			o.Point(c, []float64{umax, vmax}, 2)
			plt.Text(c[0], c[1], io.Sf("(%d,%d)", span[1], span[3]), argsIds)
		}
	}
}

// DrawElems2d draws all elements (non-zero spans)
func (o *Nurbs) DrawElems2d(npts int, withIds bool, args, argsIds *plt.A) {
	elems := o.Elements()
	for _, e := range elems {
		o.DrawElem2d(e, npts, withIds, args, argsIds)
	}
}

// PlotBasis2d plots basis function (2D only)
// option =  0 : use CalcBasis
//           1 : use CalcBasisAndDerivs
//           2 : use RecursiveBasis
func (o *Nurbs) PlotBasis2d(l int, npts, option int) {
	x := make([]float64, o.gnd)
	switch o.gnd {
	// curve
	case 1:
		U := make([]float64, npts)
		S := make([]float64, npts)
		du := (o.b[0].tmax - o.b[0].tmin) / float64(npts-1)
		uvec := []float64{0}
		for m := 0; m < npts; m++ {
			U[m] = o.b[0].tmin + float64(m)*du
			uvec[0] = U[m]
			switch option {
			case 0:
				o.CalcBasis(uvec)
				S[m] = o.GetBasisL(l)
			case 1:
				o.CalcBasisAndDerivs(uvec)
				S[m] = o.GetBasisL(l)
			case 2:
				S[m] = o.RecursiveBasis(uvec, l)
			}
		}
		plt.Plot(U, S, nil)
		plt.Gll("$u$", io.Sf("$S_%d$", l), nil)
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
				o.Point(x, u, 2)
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
		plt.ContourF(xx, yy, zz, &plt.A{NoCbar: true})
	}
	if false {
		lbls := []string{"CalcBasis function", "CalcBasisAndDerivs function", "RecursiveBasis function"}
		plt.Title(io.Sf("%s @ %d", lbls[option], l), nil)
	}
}

// PlotDeriv2d plots derivative dR[i][j][k]du[d] (2D only)
// option =  0 : use CalcBasisAndDerivs
//           1 : use NumericalDeriv
func (o *Nurbs) PlotDeriv2d(l, d int, npts, option int) {
	x := make([]float64, o.gnd)
	switch o.gnd {
	// curve
	case 1:
		U := make([]float64, npts)
		G := make([]float64, npts)
		du := (o.b[0].tmax - o.b[0].tmin) / float64(npts-1)
		uvec := []float64{0}
		gvec := []float64{0}
		for m := 0; m < npts; m++ {
			U[m] = o.b[0].tmin + float64(m)*du
			uvec[0] = U[m]
			switch option {
			case 0:
				o.CalcBasisAndDerivs(uvec)
				o.GetDerivL(gvec, l)
			case 1:
				o.NumericalDeriv(gvec, uvec, l)
			}
			G[m] = gvec[0]
		}
		plt.Plot(U, G, &plt.A{NoClip: true})
		plt.Gll("$u$", io.Sf("$G_%d$", l), nil)
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
				o.Point(x, u, 2)
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
		plt.ContourF(xx, yy, zz, &plt.A{NoCbar: true})
	}
	if false {
		lbls := []string{"CalcBasisAndDerivs function", "NumericalDeriv function"}
		plt.Title(io.Sf("%s @ %d,%d", lbls[option], l, d), nil)
	}
}

// global functions ////////////////////////////////////////////////////////////////////////////////

// PlotNurbs plots a NURBS
func PlotNurbs2d(dirout, fnkey string, b *Nurbs, npts int, withIds, withCtrl bool, argsElems, argsCtrl, argsIds *plt.A, extra func()) {
	if withCtrl {
		b.DrawCtrl2d(withIds, argsCtrl, argsIds)
	}
	b.DrawElems2d(npts, withIds, argsElems, argsIds)
	if extra != nil {
		extra()
	}
	plt.Save(dirout, fnkey)
}

// PlotNurbsBasis plots basis functions la and lb
//   First row == CalcBasis
//   Second row == CalcBasisAndDerivs
//   Third row == RecursiveBasis
func PlotNurbsBasis2d(dirout, fnkey string, b *Nurbs, la, lb int, withElems, withCtrl bool, argsElems, argsCtrl *plt.A, extra func(idxSubplot int)) {

	// configuration
	npts := 41
	ndiv := 11
	if b.gnd == 1 {
		ndiv = 101
	}
	if argsElems == nil {
		argsElems = &plt.A{C: "yellow", Ls: "-", NoClip: true}
	}

	// first row -----------------------

	plt.Subplot(3, 2, 1)
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	b.PlotBasis2d(la, ndiv, 0) // 0 => CalcBasis
	if extra != nil {
		extra(1)
	}

	plt.Subplot(3, 2, 2)
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	b.PlotBasis2d(lb, ndiv, 0) // 0 => CalcBasis
	if extra != nil {
		extra(2)
	}

	// second row -----------------------

	plt.Subplot(3, 2, 3)
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	b.PlotBasis2d(la, ndiv, 1) // 1 => CalcBasisAndDerivs
	if extra != nil {
		extra(3)
	}

	plt.Subplot(3, 2, 4)
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	b.PlotBasis2d(lb, ndiv, 1) // 1 => CalcBasisAndDerivs
	if extra != nil {
		extra(4)
	}

	// third row ------------------------

	plt.Subplot(3, 2, 5)
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	b.PlotBasis2d(la, ndiv, 2) // 2 => RecursiveBasis
	if extra != nil {
		extra(5)
	}

	plt.Subplot(3, 2, 6)
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	b.PlotBasis2d(lb, ndiv, 2) // 2 => RecursiveBasis
	if extra != nil {
		extra(6)
	}

	plt.Save(dirout, fnkey)
}

// PlotNurbsDerivs2d plots derivatives of basis functions la and lb
//  Left column == Analytical
//  Right column == Numerical
func PlotNurbsDerivs2d(dirout, fnkey string, b *Nurbs, la, lb int, withElems, withCtrl bool, argsElems, argsCtrl *plt.A, extra func(idxSubplot int)) {

	// configuration
	npts := 41
	ndiv := 41
	if b.gnd == 1 {
		ndiv = 101
	}
	if argsElems == nil {
		argsElems = &plt.A{C: "green", Ls: "-", NoClip: true}
	}

	plt.Subplot(4, 2, 1)
	b.PlotDeriv2d(la, 0, ndiv, 0) // 0 => CalcBasisAndDerivs
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	if extra != nil {
		extra(1)
	}

	plt.Subplot(4, 2, 2)
	b.PlotDeriv2d(la, 0, ndiv, 1) // 1 => NumericalDeriv
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	if extra != nil {
		extra(2)
	}

	plt.Subplot(4, 2, 3)
	b.PlotDeriv2d(la, 1, ndiv, 0) // 0 => CalcBasisAndDerivs
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	if extra != nil {
		extra(3)
	}

	plt.Subplot(4, 2, 4)
	b.PlotDeriv2d(la, 1, ndiv, 1) // 1 => NumericalDeriv
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	if extra != nil {
		extra(4)
	}

	plt.Subplot(4, 2, 5)
	b.PlotDeriv2d(lb, 0, ndiv, 0) // 0 => CalcBasisAndDerivs
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	if extra != nil {
		extra(5)
	}

	plt.Subplot(4, 2, 6)
	b.PlotDeriv2d(lb, 0, ndiv, 1) // 1 => NumericalDeriv
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	if extra != nil {
		extra(6)
	}

	plt.Subplot(4, 2, 7)
	b.PlotDeriv2d(lb, 1, ndiv, 0) // 0 => CalcBasisAndDerivs
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	if extra != nil {
		extra(7)
	}

	plt.Subplot(4, 2, 8)
	b.PlotDeriv2d(lb, 1, ndiv, 1) // 1 => NumericalDeriv
	if withCtrl {
		b.DrawCtrl2d(false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems2d(npts, false, argsElems, nil)
	}
	if extra != nil {
		extra(8)
	}

	plt.Save(dirout, fnkey)
}
