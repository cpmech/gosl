// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"gosl/chk"
	"gosl/io"
	"gosl/plt"
	"gosl/utl"
)

// Draw2d draws curve and control points
// option =  0 : use CalcBasis
//           1 : use RecursiveBasis
func (o *Bspline) Draw2d(npts, option int, withCtrl bool, argsCurve, argsCtrl *plt.A) {
	if !o.okQ {
		chk.Panic("Q must be set before calling this method")
	}
	tt := utl.LinSpace(o.tmin, o.tmax, npts)
	xx := make([]float64, npts)
	yy := make([]float64, npts)
	for i, t := range tt {
		C := o.Point(t, option)
		xx[i], yy[i] = C[0], C[1]
	}
	aCurve := argsCurve
	if aCurve == nil {
		lbls := []string{"non-recursive", "recursive"}
		aCurve = &plt.A{C: "k", Ls: "-", L: lbls[option], NoClip: true}
	}
	plt.Plot(xx, yy, aCurve)
	if withCtrl {
		qx := make([]float64, o.NumBasis())
		qy := make([]float64, o.NumBasis())
		for i := 0; i < o.NumBasis(); i++ {
			qx[i], qy[i] = o.Q[i][0], o.Q[i][1]
		}
		aCtrl := argsCtrl
		if aCtrl == nil {
			aCtrl = &plt.A{C: "r", Ls: "-", L: "ctrl", M: ".", NoClip: true}
		}
		plt.Plot(qx, qy, aCtrl)
	}
}

// Draw3d draws bspline in 3D
func (o *Bspline) Draw3d(npts int) {
	t := utl.LinSpace(o.tmin, o.tmax, npts)
	x := make([]float64, npts)
	y := make([]float64, npts)
	z := make([]float64, npts)
	for i, t := range t {
		C := o.Point(t, 0)
		x[i], y[i], z[i] = C[0], C[1], C[2]
	}
	plt.Plot3dLine(x, y, z, nil)
}

// PlotBasis plots basis functions in I
// option =  0 : use CalcBasis
//           1 : use CalcBasisAndDerivs
//           2 : use RecursiveBasis
func (o *Bspline) PlotBasis(npts, option int) {
	tt := utl.LinSpace(o.tmin, o.tmax, npts)
	I := utl.IntRange(o.NumBasis())
	f := make([]float64, len(tt))
	for _, i := range I {
		for j, t := range tt {
			switch option {
			case 0:
				o.CalcBasis(t)
				f[j] = o.GetBasis(i)
			case 1:
				o.CalcBasisAndDerivs(t)
				f[j] = o.GetBasis(i)
			case 2:
				f[j] = o.RecursiveBasis(t, i)
			}
		}
		/* TODO
		if strings.Contains(args, "marker") {
			cmd = io.Sf("label=r'%s:%d', color=GetClr(%d, 2) %s", lbls[option], i, i, args)
		} else {
			cmd = io.Sf("label=r'%s:%d', marker=(None if %d %%2 == 0 else GetMrk(%d/2,1)), markevery=(%d-1)/%d, clip_on=0, color=GetClr(%d, 2) %s", lbls[option], i, i, i, npts, nmks, i, args)
		}
		plt.Plot(tt, f, cmd)
		*/
		plt.Plot(tt, f, nil)
	}
	plt.Gll("$x$", io.Sf("$N_{i,%d}$", o.p), &plt.A{LegOut: true, LegNcol: o.NumBasis(), LegHlen: 1.5, FszLeg: 7})
	o.pltTicksSpans()
}

// PlotDerivs plots derivatives of basis functions in I
func (o *Bspline) PlotDerivs(npts int) {
	tt := utl.LinSpace(o.tmin, o.tmax, npts)
	I := utl.IntRange(o.NumBasis())
	f := make([]float64, len(tt))
	for _, i := range I {
		for j, t := range tt {
			o.CalcBasisAndDerivs(t)
			f[j] = o.GetDeriv(i)
		}
		/* TODO
		if strings.Contains(args, "marker") {
			cmd = io.Sf("label=r'%s:%d', color=GetClr(%d, 2) %s", lbls[option], i, i, args)
		} else {
			cmd = io.Sf("label=r'%s:%d', marker=(None if %d %%2 == 0 else GetMrk(%d/2,1)), markevery=(%d-1)/%d, clip_on=0, color=GetClr(%d, 2) %s", lbls[option], i, i, i, npts, nmks, i, args)
		}
		*/
		plt.Plot(tt, f, nil)
	}
	plt.Gll("$t$", io.Sf(`$\frac{\mathrm{d}N_{i,%d}}{\mathrm{d}t}$`, o.p), &plt.A{LegOut: true, LegNcol: o.NumBasis(), LegHlen: 1.5, FszLeg: 7})
	o.pltTicksSpans()
}

// pltTicksSpans adds ticks indicating spans
func (o *Bspline) pltTicksSpans() {
	lbls := make(map[float64]string, 0)
	for i, t := range o.T {
		if _, ok := lbls[t]; !ok {
			lbls[t] = io.Sf("[%d", i)
		} else {
			lbls[t] += io.Sf(",%d", i)
		}
	}
	for t, l := range lbls {
		plt.AnnotateXlabels(t, io.Sf("%s]", l), nil)
	}
}
