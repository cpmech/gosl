// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gm (geometry and meshes) implements auxiliary functions for
// handling geometry and mesh structures
package gm

import (
	"strings"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// Draw2d draws curve and control points
// option =  0 : use CalcBasis
//           1 : use RecursiveBasis
func (o *Bspline) Draw2d(curveArgs, ctrlArgs string, npts, option int) {
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
	qx := make([]float64, o.NumBasis())
	qy := make([]float64, o.NumBasis())
	for i := 0; i < o.NumBasis(); i++ {
		qx[i], qy[i] = o.Q[i][0], o.Q[i][1]
	}
	lbls := []string{"Nonly", "recN"}
	plt.Plot(xx, yy, io.Sf("'k-', clip_on=0, label=r'%s'", lbls[option])+curveArgs)
	plt.Plot(qx, qy, "'r-', clip_on=0, label=r'ctrl', marker='.'"+ctrlArgs)
	plt.Gll("$x$", "$y$", "leg=1, leg_out=1, leg_ncol=2, leg_hlen=1.5, leg_fsz=7")
}

func (o *Bspline) Draw3d(curveArgs, ctrlArgs string, npts int, first bool) {
	t := utl.LinSpace(o.tmin, o.tmax, npts)
	x := make([]float64, npts)
	y := make([]float64, npts)
	z := make([]float64, npts)
	for i, t := range t {
		C := o.Point(t, 0)
		x[i], y[i], z[i] = C[0], C[1], C[2]
	}
	plt.Plot3dLine(x, y, z, first, "")
}

// PlotBasis plots basis functions in I
// option =  0 : use CalcBasis
//           1 : use CalcBasisAndDerivs
//           2 : use RecursiveBasis
func (o *Bspline) PlotBasis(args string, npts, option int) {
	nmks := 10
	tt := utl.LinSpace(o.tmin, o.tmax, npts)
	I := utl.IntRange(o.NumBasis())
	f := make([]float64, len(tt))
	lbls := []string{"Nonly", "N\\&dN", "recN"}
	var cmd string
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
		if strings.Contains(args, "marker") {
			cmd = io.Sf("label=r'%s:%d', color=GetClr(%d, 2) %s", lbls[option], i, i, args)
		} else {
			cmd = io.Sf("label=r'%s:%d', marker=(None if %d %%2 == 0 else GetMrk(%d/2,1)), markevery=(%d-1)/%d, clip_on=0, color=GetClr(%d, 2) %s", lbls[option], i, i, i, npts, nmks, i, args)
		}
		plt.Plot(tt, f, cmd)
	}
	plt.Gll("$t$", io.Sf("$N_{i,%d}$", o.p), io.Sf("leg=1, leg_out=1, leg_ncol=%d, leg_hlen=1.5, leg_fsz=7", o.NumBasis()))
	o.plt_ticks_spans()
}

// PlotDerivs plots derivatives of basis functions in I
// option =  0 : use CalcBasisAndDerivs
//           1 : use NumericalDeriv
func (o *Bspline) PlotDerivs(args string, npts, option int) {
	nmks := 10
	tt := utl.LinSpace(o.tmin, o.tmax, npts)
	I := utl.IntRange(o.NumBasis())
	f := make([]float64, len(tt))
	lbls := []string{"N\\&dN", "numD"}
	var cmd string
	for _, i := range I {
		for j, t := range tt {
			switch option {
			case 0:
				o.CalcBasisAndDerivs(t)
				f[j] = o.GetDeriv(i)
			case 1:
				f[j] = o.NumericalDeriv(t, i)
			}
		}
		if strings.Contains(args, "marker") {
			cmd = io.Sf("label=r'%s:%d', color=GetClr(%d, 2) %s", lbls[option], i, i, args)
		} else {
			cmd = io.Sf("label=r'%s:%d', marker=(None if %d %%2 == 0 else GetMrk(%d/2,1)), markevery=(%d-1)/%d, clip_on=0, color=GetClr(%d, 2) %s", lbls[option], i, i, i, npts, nmks, i, args)
		}
		plt.Plot(tt, f, cmd)
	}
	plt.Gll("$t$", io.Sf(`$\frac{\mathrm{d}N_{i,%d}}{\mathrm{d}t}$`, o.p), io.Sf("leg=1, leg_out=1, leg_ncol=%d, leg_hlen=1.5, leg_fsz=7", o.NumBasis()))
	o.plt_ticks_spans()
}

// plt_ticks_spans adds ticks indicating spans
func (o *Bspline) plt_ticks_spans() {
	lbls := make(map[float64]string, 0)
	for i, t := range o.T {
		if _, ok := lbls[t]; !ok {
			lbls[t] = io.Sf("'[%d", i)
		} else {
			lbls[t] += io.Sf(",%d", i)
		}
	}
	for t, l := range lbls {
		plt.AnnotateXlabels(t, io.Sf("%s]'", l), "")
	}
}
