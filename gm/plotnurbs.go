// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// basis and derivatives ///////////////////////////////////////////////////////////////////////////

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
		xx := utl.Alloc(npts, npts)
		yy := utl.Alloc(npts, npts)
		zz := utl.Alloc(npts, npts)
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
}

// PlotDeriv2d plots derivative dR[i][j][k]du[d] (2D only)
func (o *Nurbs) PlotDeriv2d(l, d int, npts int) {
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
			o.CalcBasisAndDerivs(uvec)
			o.GetDerivL(gvec, l)
			G[m] = gvec[0]
		}
		plt.Plot(U, G, &plt.A{NoClip: true})
		plt.Gll("$u$", io.Sf("$G_%d$", l), nil)
	// surface
	case 2:
		xx := utl.Alloc(npts, npts)
		yy := utl.Alloc(npts, npts)
		zz := utl.Alloc(npts, npts)
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
				o.CalcBasisAndDerivs(u)
				o.GetDerivL(drdu, l)
				zz[m][n] = drdu[d]
			}
		}
		plt.ContourF(xx, yy, zz, &plt.A{NoCbar: true})
	}
}

// global functions ////////////////////////////////////////////////////////////////////////////////

// PlotNurbs plots a NURBS
//  ndim -- 2 or 3 => to plot in a 2D or 3D space
func PlotNurbs(dirout, fnkey string, b *Nurbs, ndim, npts int, withIds, withCtrl bool, argsElems, argsCtrl, argsIds *plt.A, extra func()) {
	if withCtrl {
		b.DrawCtrl(ndim, withIds, argsCtrl, argsIds)
	}
	b.DrawElems(ndim, npts, withIds, argsElems, argsIds)
	if extra != nil {
		extra()
	}
	plt.Save(dirout, fnkey)
}

// PlotNurbsBasis2d plots basis functions la and lb
//   First row == CalcBasis
//   Second row == CalcBasisAndDerivs
//   Third row == RecursiveBasis
func PlotNurbsBasis2d(dirout, fnkey string, b *Nurbs, la, lb int, withElems, withCtrl bool, argsElems, argsCtrl *plt.A, extra func(idxSubplot int)) {

	// configuration
	ndim := 2
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
		b.DrawCtrl(ndim, false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems(ndim, npts, false, argsElems, nil)
	}
	b.PlotBasis2d(la, ndiv, 0) // 0 => CalcBasis
	if extra != nil {
		extra(1)
	}

	plt.Subplot(3, 2, 2)
	if withCtrl {
		b.DrawCtrl(ndim, false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems(ndim, npts, false, argsElems, nil)
	}
	b.PlotBasis2d(lb, ndiv, 0) // 0 => CalcBasis
	if extra != nil {
		extra(2)
	}

	// second row -----------------------

	plt.Subplot(3, 2, 3)
	if withCtrl {
		b.DrawCtrl(ndim, false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems(ndim, npts, false, argsElems, nil)
	}
	b.PlotBasis2d(la, ndiv, 1) // 1 => CalcBasisAndDerivs
	if extra != nil {
		extra(3)
	}

	plt.Subplot(3, 2, 4)
	if withCtrl {
		b.DrawCtrl(ndim, false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems(ndim, npts, false, argsElems, nil)
	}
	b.PlotBasis2d(lb, ndiv, 1) // 1 => CalcBasisAndDerivs
	if extra != nil {
		extra(4)
	}

	// third row ------------------------

	plt.Subplot(3, 2, 5)
	if withCtrl {
		b.DrawCtrl(ndim, false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems(ndim, npts, false, argsElems, nil)
	}
	b.PlotBasis2d(la, ndiv, 2) // 2 => RecursiveBasis
	if extra != nil {
		extra(5)
	}

	plt.Subplot(3, 2, 6)
	if withCtrl {
		b.DrawCtrl(ndim, false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems(ndim, npts, false, argsElems, nil)
	}
	b.PlotBasis2d(lb, ndiv, 2) // 2 => RecursiveBasis
	if extra != nil {
		extra(6)
	}

	plt.Save(dirout, fnkey)
}

// PlotNurbsDerivs2d plots derivatives of basis functions la and lb in 2D space
//  Left column == Analytical
//  Right column == Numerical
func PlotNurbsDerivs2d(dirout, fnkey string, b *Nurbs, la, lb int, withElems, withCtrl bool, argsElems, argsCtrl *plt.A, extra func(idxSubplot int)) {

	// configuration
	ndim := 2
	npts := 41
	ndiv := 41
	if b.gnd == 1 {
		ndiv = 101
	}
	if argsElems == nil {
		argsElems = &plt.A{C: "green", Ls: "-", NoClip: true}
	}

	plt.Subplot(2, 2, 1)
	b.PlotDeriv2d(la, 0, ndiv)
	if withCtrl {
		b.DrawCtrl(ndim, false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems(ndim, npts, false, argsElems, nil)
	}
	if extra != nil {
		extra(1)
	}

	plt.Subplot(2, 2, 2)
	b.PlotDeriv2d(la, 1, ndiv)
	if withCtrl {
		b.DrawCtrl(ndim, false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems(ndim, npts, false, argsElems, nil)
	}
	if extra != nil {
		extra(2)
	}

	plt.Subplot(2, 2, 3)
	b.PlotDeriv2d(lb, 0, ndiv)
	if withCtrl {
		b.DrawCtrl(ndim, false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems(ndim, npts, false, argsElems, nil)
	}
	if extra != nil {
		extra(3)
	}

	plt.Subplot(2, 2, 4)
	b.PlotDeriv2d(lb, 1, ndiv)
	if withCtrl {
		b.DrawCtrl(ndim, false, argsCtrl, nil)
	}
	if withElems {
		b.DrawElems(ndim, npts, false, argsElems, nil)
	}
	if extra != nil {
		extra(4)
	}

	plt.Save(dirout, fnkey)
}

// draw NURBS methods //////////////////////////////////////////////////////////////////////////////

// DrawCtrl draws control net in 2D or 3D space
//  ndim -- 2 or 3 => to plot in a 2D or 3D space
func (o *Nurbs) DrawCtrl(ndim int, withIds bool, args, argsIds *plt.A) {
	if args == nil {
		args = &plt.A{C: "k", M: ".", Ls: "--", NoClip: true}
	}
	if argsIds == nil {
		argsIds = &plt.A{C: "r", Fsz: 7, NoClip: true}
	}
	if o.gnd == 3 {
		ndim = 3
	}
	switch o.gnd {
	// curve
	case 1:
		xa := make([]float64, o.n[0])
		ya := make([]float64, o.n[0])
		za := make([]float64, o.n[0])
		j, k := 0, 0
		for i := 0; i < o.n[0]; i++ {
			xa[i] = o.Q[i][j][k][0] / o.Q[i][j][k][3]
			ya[i] = o.Q[i][j][k][1] / o.Q[i][j][k][3]
			za[i] = o.Q[i][j][k][2] / o.Q[i][j][k][3]
		}
		if ndim == 3 {
			plt.Plot3dLine(xa, ya, za, args)
		} else {
			plt.Plot(xa, ya, args)
		}
		if withIds {
			for i := 0; i < o.n[0]; i++ {
				x := o.Q[i][j][k][0] / o.Q[i][j][k][3]
				y := o.Q[i][j][k][1] / o.Q[i][j][k][3]
				z := o.Q[i][j][k][2] / o.Q[i][j][k][3]
				if ndim == 3 {
					plt.Text3d(x, y, z, io.Sf("%d", i), argsIds)
				} else {
					plt.Text(x, y, io.Sf("%d", i), argsIds)
				}
			}
		}
	// surface
	case 2:
		xa := make([]float64, o.n[1])
		ya := make([]float64, o.n[1])
		za := make([]float64, o.n[1])
		k := 0
		for i := 0; i < o.n[0]; i++ {
			for j := 0; j < o.n[1]; j++ {
				xa[j] = o.Q[i][j][k][0] / o.Q[i][j][k][3]
				ya[j] = o.Q[i][j][k][1] / o.Q[i][j][k][3]
				za[j] = o.Q[i][j][k][2] / o.Q[i][j][k][3]
			}
			if ndim == 3 {
				plt.Plot3dLine(xa, ya, za, args)
			} else {
				plt.Plot(xa, ya, args)
			}
		}
		xb := make([]float64, o.n[0])
		yb := make([]float64, o.n[0])
		zb := make([]float64, o.n[0])
		for j := 0; j < o.n[1]; j++ {
			for i := 0; i < o.n[0]; i++ {
				xb[i] = o.Q[i][j][k][0] / o.Q[i][j][k][3]
				yb[i] = o.Q[i][j][k][1] / o.Q[i][j][k][3]
				zb[i] = o.Q[i][j][k][2] / o.Q[i][j][k][3]
			}
			if ndim == 3 {
				plt.Plot3dLine(xb, yb, zb, args)
			} else {
				plt.Plot(xb, yb, args)
			}
		}
		if withIds {
			for i := 0; i < o.n[0]; i++ {
				for j := 0; j < o.n[1]; j++ {
					x := o.Q[i][j][k][0] / o.Q[i][j][k][3]
					y := o.Q[i][j][k][1] / o.Q[i][j][k][3]
					z := o.Q[i][j][k][2] / o.Q[i][j][k][3]
					l := i + j*o.n[0]
					if ndim == 3 {
						plt.Text3d(x, y, z, io.Sf("%d", l), argsIds)
					} else {
						plt.Text(x, y, io.Sf("%d", l), argsIds)
					}
				}
			}
		}
	// solid
	case 3:
		xx := make([]float64, 2) // min,max
		yy := make([]float64, 2) // min,max
		zz := make([]float64, 2) // min,max
		for i := 0; i < o.n[0]; i++ {
			for j := 0; j < o.n[1]; j++ {
				for k := 0; k < o.n[2]; k++ {
					x := o.Q[i][j][k][0] / o.Q[i][j][k][3]
					y := o.Q[i][j][k][1] / o.Q[i][j][k][3]
					z := o.Q[i][j][k][2] / o.Q[i][j][k][3]
					if i > 0 {
						xp := o.Q[i-1][j][k][0] / o.Q[i-1][j][k][3]
						yp := o.Q[i-1][j][k][1] / o.Q[i-1][j][k][3]
						zp := o.Q[i-1][j][k][2] / o.Q[i-1][j][k][3]
						xx[0], xx[1] = xp, x
						yy[0], yy[1] = yp, y
						zz[0], zz[1] = zp, z
						plt.Plot3dLine(xx, yy, zz, args)
					}
				}
			}
		}
		for j := 0; j < o.n[1]; j++ {
			for i := 0; i < o.n[0]; i++ {
				for k := 0; k < o.n[2]; k++ {
					x := o.Q[i][j][k][0] / o.Q[i][j][k][3]
					y := o.Q[i][j][k][1] / o.Q[i][j][k][3]
					z := o.Q[i][j][k][2] / o.Q[i][j][k][3]
					if j > 0 {
						xp := o.Q[i][j-1][k][0] / o.Q[i][j-1][k][3]
						yp := o.Q[i][j-1][k][1] / o.Q[i][j-1][k][3]
						zp := o.Q[i][j-1][k][2] / o.Q[i][j-1][k][3]
						xx[0], xx[1] = xp, x
						yy[0], yy[1] = yp, y
						zz[0], zz[1] = zp, z
						plt.Plot3dLine(xx, yy, zz, args)
					}
				}
			}
		}
		for k := 0; k < o.n[2]; k++ {
			for j := 0; j < o.n[1]; j++ {
				for i := 0; i < o.n[0]; i++ {
					x := o.Q[i][j][k][0] / o.Q[i][j][k][3]
					y := o.Q[i][j][k][1] / o.Q[i][j][k][3]
					z := o.Q[i][j][k][2] / o.Q[i][j][k][3]
					if k > 0 {
						xp := o.Q[i][j][k-1][0] / o.Q[i][j][k-1][3]
						yp := o.Q[i][j][k-1][1] / o.Q[i][j][k-1][3]
						zp := o.Q[i][j][k-1][2] / o.Q[i][j][k-1][3]
						xx[0], xx[1] = xp, x
						yy[0], yy[1] = yp, y
						zz[0], zz[1] = zp, z
						plt.Plot3dLine(xx, yy, zz, args)
					}
					if withIds {
						l := i + j*o.n[0] + k*o.n[0]*o.n[1]
						plt.Text3d(x, y, z, io.Sf("%d", l), argsIds)
					}
				}
			}
		}
	}
}

// DrawEdge draws edge from tmin to tmax in 2D or 3D space
//  ndim -- 2 or 3 => to plot in a 2D or 3D space
func (o *Nurbs) DrawEdge(ndim int, tmin, tmax, cteA, cteB float64, along, npts int, args *plt.A) {
	if args == nil {
		args = &plt.A{C: "b", Ls: "-", NoClip: true}
	}
	tt := utl.LinSpace(tmin, tmax, npts)
	xx := make([]float64, npts)
	yy := make([]float64, npts)
	zz := make([]float64, npts)
	x := make([]float64, 3)
	u := make([]float64, 3)
	for i, t := range tt {
		switch along {
		case 0:
			u[0], u[1], u[2] = t, cteA, cteB
		case 1:
			u[0], u[1], u[2] = cteA, t, cteB
		case 2:
			u[0], u[1], u[2] = cteA, cteB, t
		}
		o.Point(x, u, ndim)
		xx[i], yy[i], zz[i] = x[0], x[1], x[2]
	}
	if ndim == 2 {
		plt.Plot(xx, yy, args)
	} else {
		plt.Plot3dLine(xx, yy, zz, args)
	}
}

// DrawElem draws element (non-zero span) in 2D or 3D space
//  ndim -- 2 or 3 => to plot in a 2D or 3D space
func (o *Nurbs) DrawElem(ndim int, span []int, npts int, withIds bool, args, argsIds *plt.A) {
	if argsIds == nil {
		argsIds = &plt.A{C: "r", Fsz: 7}
	}
	c := make([]float64, 3)
	if o.gnd == 3 {
		ndim = 3
	}
	switch o.gnd {
	// curve
	case 1:
		umin, umax := o.b[0].T[span[0]], o.b[0].T[span[1]]
		o.DrawEdge(ndim, umin, umax, 0.0, 0.0, 0, npts, args)
		if withIds {
			o.Point(c, []float64{umin}, ndim)
			drawElemID(c, ndim, span[0], -1, -1, argsIds)
			o.Point(c, []float64{umax}, ndim)
			drawElemID(c, ndim, span[1], -1, -1, argsIds)
		}
	// surface
	case 2:
		umin, umax := o.b[0].T[span[0]], o.b[0].T[span[1]]
		vmin, vmax := o.b[1].T[span[2]], o.b[1].T[span[3]]
		o.DrawEdge(ndim, umin, umax, vmin, 0.0, 0, npts, args)
		o.DrawEdge(ndim, umin, umax, vmax, 0.0, 0, npts, args)
		o.DrawEdge(ndim, vmin, vmax, umin, 0.0, 1, npts, args)
		o.DrawEdge(ndim, vmin, vmax, umax, 0.0, 1, npts, args)
		if withIds {
			o.Point(c, []float64{umin, vmin}, ndim)
			drawElemID(c, ndim, span[0], span[2], -1, argsIds)
			o.Point(c, []float64{umin, vmax}, ndim)
			drawElemID(c, ndim, span[0], span[3], -1, argsIds)
			o.Point(c, []float64{umax, vmin}, ndim)
			drawElemID(c, ndim, span[1], span[2], -1, argsIds)
			o.Point(c, []float64{umax, vmax}, ndim)
			drawElemID(c, ndim, span[1], span[3], -1, argsIds)
		}
		// solid
	case 3:
		umin, umax := o.b[0].T[span[0]], o.b[0].T[span[1]]
		vmin, vmax := o.b[1].T[span[2]], o.b[1].T[span[3]]
		wmin, wmax := o.b[2].T[span[4]], o.b[2].T[span[5]]
		o.DrawEdge(ndim, umin, umax, vmin, wmin, 0, npts, args)
		o.DrawEdge(ndim, umin, umax, vmax, wmin, 0, npts, args)
		o.DrawEdge(ndim, vmin, vmax, umin, wmin, 1, npts, args)
		o.DrawEdge(ndim, vmin, vmax, umax, wmin, 1, npts, args)
		o.DrawEdge(ndim, umin, umax, vmin, wmax, 0, npts, args)
		o.DrawEdge(ndim, umin, umax, vmax, wmax, 0, npts, args)
		o.DrawEdge(ndim, vmin, vmax, umin, wmax, 1, npts, args)
		o.DrawEdge(ndim, vmin, vmax, umax, wmax, 1, npts, args)
		o.DrawEdge(ndim, wmin, wmax, umin, vmin, 2, npts, args)
		o.DrawEdge(ndim, wmin, wmax, umax, vmin, 2, npts, args)
		o.DrawEdge(ndim, wmin, wmax, umin, vmax, 2, npts, args)
		o.DrawEdge(ndim, wmin, wmax, umax, vmax, 2, npts, args)
		if withIds {
			o.Point(c, []float64{umin, vmin, wmin}, ndim)
			drawElemID(c, ndim, span[0], span[2], span[4], argsIds)
			o.Point(c, []float64{umin, vmax, wmin}, ndim)
			drawElemID(c, ndim, span[0], span[3], span[4], argsIds)
			o.Point(c, []float64{umax, vmin, wmin}, ndim)
			drawElemID(c, ndim, span[1], span[2], span[4], argsIds)
			o.Point(c, []float64{umax, vmax, wmin}, ndim)
			drawElemID(c, ndim, span[1], span[3], span[4], argsIds)
			o.Point(c, []float64{umin, vmin, wmax}, ndim)
			drawElemID(c, ndim, span[0], span[2], span[5], argsIds)
			o.Point(c, []float64{umin, vmax, wmax}, ndim)
			drawElemID(c, ndim, span[0], span[3], span[5], argsIds)
			o.Point(c, []float64{umax, vmin, wmax}, ndim)
			drawElemID(c, ndim, span[1], span[2], span[5], argsIds)
			o.Point(c, []float64{umax, vmax, wmax}, ndim)
			drawElemID(c, ndim, span[1], span[3], span[5], argsIds)
		}
	}
}

// DrawElems draws all elements (non-zero spans) in 2D or 3D space
//  ndim -- 2 or 3 => to plot in a 2D or 3D space
func (o *Nurbs) DrawElems(ndim int, npts int, withIds bool, args, argsIds *plt.A) {
	elems := o.Elements()
	for _, e := range elems {
		o.DrawElem(ndim, e, npts, withIds, args, argsIds)
	}
}

// DrawSurface draws surface
//  ndim -- 2 or 3 => to plot in a 2D or 3D space
func (o *Nurbs) DrawSurface(ndim int, nu, nv int, withSurf, withWire bool, argsSurf, argsWire *plt.A) {
	if o.gnd != 2 {
		return
	}
	X := make([][]float64, nu)
	Y := make([][]float64, nu)
	var Z [][]float64
	if ndim == 3 {
		Z = make([][]float64, nu)
	}
	du0 := (o.b[0].tmax - o.b[0].tmin) / float64(nu-1)
	du1 := (o.b[1].tmax - o.b[1].tmin) / float64(nv-1)
	x := []float64{0, 0, 0}
	u := []float64{0, 0}
	for m := 0; m < nu; m++ {
		X[m] = make([]float64, nv)
		Y[m] = make([]float64, nv)
		if ndim == 3 {
			Z[m] = make([]float64, nv)
		}
		u[0] = o.b[0].tmin + float64(m)*du0
		for n := 0; n < nv; n++ {
			u[1] = o.b[1].tmin + float64(n)*du1
			o.Point(x, u, ndim)
			X[m][n] = x[0]
			Y[m][n] = x[1]
			if ndim == 3 {
				Z[m][n] = x[2]
			}
		}
	}
	if ndim == 2 {
		if withWire {
			plt.Grid2d(X, Y, false, argsWire, nil)
		}
	} else {
		if withSurf {
			plt.Surface(X, Y, Z, argsSurf)
		}
		if withWire {
			plt.Wireframe(X, Y, Z, argsWire)
		}
	}
}

// DrawSolid draws wireframe representing solid NURBS
func (o *Nurbs) DrawSolid(nu, nv, nw int, args *plt.A) {
	if o.gnd != 3 {
		return
	}
	du := (o.b[0].tmax - o.b[0].tmin) / float64(nu-1)
	dv := (o.b[1].tmax - o.b[1].tmin) / float64(nv-1)
	dw := (o.b[2].tmax - o.b[2].tmin) / float64(nw-1)
	u := make([]float64, 3)
	x := make([]float64, 3)
	x0 := make([]float64, nu)
	y0 := make([]float64, nu)
	z0 := make([]float64, nu)
	x1 := make([]float64, nv)
	y1 := make([]float64, nv)
	z1 := make([]float64, nv)
	x2 := make([]float64, nw)
	y2 := make([]float64, nw)
	z2 := make([]float64, nw)

	// draw 0-lines
	for k := 0; k < nw; k++ {
		u[2] = o.b[2].tmin + float64(k)*dw
		for j := 0; j < nv; j++ {
			u[1] = o.b[1].tmin + float64(j)*dv
			for i := 0; i < nu; i++ {
				u[0] = o.b[0].tmin + float64(i)*du
				o.Point(x, u, 3)
				x0[i] = x[0]
				y0[i] = x[1]
				z0[i] = x[2]
			}
			plt.Plot3dLine(x0, y0, z0, args)
		}
	}

	// draw 1-lines
	for k := 0; k < nw; k++ {
		u[2] = o.b[2].tmin + float64(k)*dw
		for i := 0; i < nu; i++ {
			u[0] = o.b[0].tmin + float64(i)*du
			for j := 0; j < nv; j++ {
				u[1] = o.b[1].tmin + float64(j)*dv
				o.Point(x, u, 3)
				x1[j] = x[0]
				y1[j] = x[1]
				z1[j] = x[2]
			}
			plt.Plot3dLine(x1, y1, z1, args)
		}
	}

	// draw 2-lines
	for j := 0; j < nv; j++ {
		u[1] = o.b[1].tmin + float64(j)*dv
		for i := 0; i < nu; i++ {
			u[0] = o.b[0].tmin + float64(i)*du
			for k := 0; k < nw; k++ {
				u[2] = o.b[2].tmin + float64(k)*dw
				o.Point(x, u, 3)
				x2[k] = x[0]
				y2[k] = x[1]
				z2[k] = x[2]
			}
			plt.Plot3dLine(x2, y2, z2, args)
		}
	}
}

// DrawVectors3d draws tangent vectors of 3D surface
func (o *Nurbs) DrawVectors3d(nu, nv int, sf float64, argsU, argsV *plt.A) {
	if o.gnd != 2 {
		chk.Panic("method works with surfaces only\n")
	}
	if argsU == nil {
		argsU = &plt.A{C: plt.C(0, 0)}
	}
	if argsV == nil {
		argsV = &plt.A{C: plt.C(1, 0)}
	}
	du0 := (o.b[0].tmax - o.b[0].tmin) / float64(nu-1)
	du1 := (o.b[1].tmax - o.b[1].tmin) / float64(nv-1)
	u := la.NewVector(2)
	c := la.NewVector(3)
	dCdu := la.NewMatrix(3, o.gnd)
	for n := 0; n < nv; n++ {
		u[1] = o.b[1].tmin + float64(n)*du1
		for m := 0; m < nu; m++ {
			u[0] = o.b[0].tmin + float64(m)*du0
			o.PointAndFirstDerivs(dCdu, c, u, 3)
			plt.DrawArrow3d(c, dCdu.Col(0), true, sf, argsU)
			plt.DrawArrow3d(c, dCdu.Col(1), true, sf, argsV)
		}
	}
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// drawElemID draws element id
func drawElemID(c []float64, ndim, i, j, k int, args *plt.A) {
	txt := io.Sf("(%d", i)
	if j >= 0 {
		txt += io.Sf(",%d", j)
	}
	if k >= 0 {
		txt += io.Sf(",%d", k)
	}
	txt += ")"
	if ndim == 2 {
		plt.Text(c[0], c[1], txt, args)
	} else {
		plt.Text3d(c[0], c[1], c[2], txt, args)
	}
}

// plotTwoNurbs2d plots two NURBS in 2d space
func plotTwoNurbs2d(dirout, fnkey string, a, b *Nurbs, labelA, labelB string, extra func()) {
	str := "curve: "
	if a.gnd > 1 {
		str = "elems: "
	}
	argsCtrlA := &plt.A{C: "k", Ls: "--", L: "control: " + labelA, NoClip: true}
	argsCtrlB := &plt.A{C: "green", L: "control: " + labelB, NoClip: true}
	argsElemsA := &plt.A{C: "b", L: str + labelA, NoClip: true}
	argsElemsB := &plt.A{C: "orange", Ls: "none", M: "*", Me: 20, L: str + labelB, NoClip: true}
	argsIdsCtrlA := &plt.A{C: "k", Fsz: 7, NoClip: true}
	argsIdsCtrlB := &plt.A{C: "green", Fsz: 7, NoClip: true}
	argsIdsA := &plt.A{C: "r", Fsz: 7, NoClip: true}
	ndim := 2
	npts := 41
	a.DrawCtrl(ndim, true, argsCtrlA, argsIdsCtrlA)
	a.DrawElems(ndim, npts, true, argsElemsA, argsIdsA)
	b.DrawCtrl(ndim, true, argsCtrlB, argsIdsCtrlB)
	b.DrawElems(ndim, npts, false, argsElemsB, nil)
	if extra != nil {
		extra()
	}
	plt.LegendX([]*plt.A{argsCtrlA, argsCtrlB, argsElemsA, argsElemsB}, &plt.A{LegOut: true, LegNcol: 2})
	plt.Save(dirout, fnkey)
}
