// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// PlotX plots X values
//   args -- maps int(y) values to plot arguments
func (o *DataMatrix) PlotX(iFeature, jFeature int, args map[int]*plt.A) {
	a := &plt.A{C: "orange", M: ".", Ls: "none", NoClip: true}
	i, j := 1+iFeature, 1+jFeature
	if args == nil {
		xivals := make([]float64, o.nSamples)
		xjvals := make([]float64, o.nSamples)
		for k := 0; k < o.nSamples; k++ {
			xivals[k] = o.xMat.Get(k, i)
			xjvals[k] = o.xMat.Get(k, j)
		}
		plt.Plot(xivals, xjvals, a)
	} else {
		var aa *plt.A
		for k := 0; k < o.nSamples; k++ {
			aa = args[int(o.yVec[k])]
			if aa == nil {
				aa = a
			}
			aa.NoClip = true
			plt.PlotOne(o.xMat.Get(k, i), o.xMat.Get(k, j), aa)
		}
	}
	plt.Gll(io.Sf("$x_{%d}$", iFeature), io.Sf("$x_{%d}$", jFeature), nil)
}

// PlotModel plots regression model @ iFeature and all other values set to MinX
func (o *DataMatrix) PlotModel(reg Regression, iFeature, npts int, args *plt.A) {
	if npts < 3 {
		npts = 11
	}
	if args == nil {
		args = &plt.A{C: plt.C(0, 0), NoClip: true}
	}
	i := 1 + iFeature
	xVec := o.newXvec(true, false, false)
	if !o.statOk {
		o.stat()
	}
	dxi := (o.maxX[iFeature] - o.minX[iFeature]) / float64(npts-1)
	if len(o.horiz) != npts {
		o.horiz = make([]float64, npts)
		o.vert = make([]float64, npts)
	}
	for k := 0; k < npts; k++ {
		o.horiz[k] = o.minX[iFeature] + float64(k)*dxi
		xVec[i] = o.horiz[k]
		o.vert[k] = reg.Model(xVec)
	}
	plt.Plot(o.horiz, o.vert, args)
}

// PlotContModel plots contour of model
//  xmin -- minimum x[iFeature,jFeature] value for meshgrid; may be nil
//  xmax -- maximum x[iFeature,jFeature] value for meshgrid; may be nil
func (o *DataMatrix) PlotContModel(reg Regression, iFeature, jFeature, npi, npj int, mapper DataMapper, xmin, xmax []float64, filled bool, args *plt.A) {
	if len(o.xxi) != npj {
		o.xxi = utl.Alloc(npj, npi)
		o.xxj = utl.Alloc(npj, npi)
		o.zz = utl.Alloc(npj, npi)
	}
	i, j := 1+iFeature, 1+jFeature
	xVec := o.newXvec(true, false, false)
	if !o.statOk {
		o.stat()
	}
	if xmin == nil {
		xmin = o.minX
	}
	if xmax == nil {
		xmax = o.maxX
	}
	var xRaw []float64
	if mapper != nil {
		xRaw = make([]float64, mapper.NumOriginalFeatures())
	}
	dxi := (xmax[iFeature] - xmin[iFeature]) / float64(npi-1)
	dxj := (xmax[jFeature] - xmin[jFeature]) / float64(npj-1)
	for r := 0; r < npj; r++ {
		for c := 0; c < npi; c++ {
			o.xxi[r][c] = xmin[iFeature] + float64(c)*dxi
			o.xxj[r][c] = xmin[jFeature] + float64(r)*dxj
			if mapper == nil {
				xVec[i] = o.xxi[r][c]
				xVec[j] = o.xxj[r][c]
			} else {
				xRaw[iFeature] = o.xxi[r][c]
				xRaw[jFeature] = o.xxj[r][c]
				mapper.Map(xVec, xRaw)
			}
			o.zz[r][c] = reg.Model(xVec)
		}
	}
	if filled {
		plt.ContourF(o.xxi, o.xxj, o.zz, args)
	} else {
		if args == nil {
			args = &plt.A{Colors: []string{"k"}, Levels: []float64{0}}
		}
		plt.ContourL(o.xxi, o.xxj, o.zz, args)
	}
	plt.Gll(io.Sf("$x_{%d}$", iFeature), io.Sf("$x_{%d}$", jFeature), nil)
}

// PlotContCost plots contour of cost
func (o *DataMatrix) PlotContCost(reg Regression, iPrm, jPrm, npi, npj int, minθ, maxθ []float64, args *plt.A) {
	if len(o.thi) != npj {
		o.thi = utl.Alloc(npj, npi)
		o.thj = utl.Alloc(npj, npi)
		o.cc = utl.Alloc(npj, npi)
	}
	θcpy, bcpy := reg.GetParams()
	dthi := (maxθ[iPrm] - minθ[iPrm]) / float64(npi-1)
	dthj := (maxθ[jPrm] - minθ[jPrm]) / float64(npj-1)
	for r := 0; r < npj; r++ {
		for c := 0; c < npi; c++ {
			o.thi[r][c] = minθ[iPrm] + float64(c)*dthi
			o.thj[r][c] = minθ[jPrm] + float64(r)*dthj
			reg.SetTheta(iPrm, o.thi[r][c])
			reg.SetTheta(jPrm, o.thj[r][c])
			o.cc[r][c] = reg.Cost(o)
		}
	}
	plt.ContourF(o.thi, o.thj, o.cc, args)
	plt.Gll(io.Sf("$\\theta_{%d}$", iPrm), io.Sf("$\\theta_{%d}$", jPrm), nil)
	reg.SetParams(θcpy, bcpy)
}
