// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ml implements functions to develop Machine Learning algorithms
package ml

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Regression defines the functions required to perform regression computations
type Regression interface {
	Model(xVec, theta la.Vector) float64     // model equation where xVec[1+nFeatures] (augmented vector) and theta[1+nFeatures]
	Cost(data *RegData) float64              // computes cost
	Deriv(dCdTheta la.Vector, data *RegData) // computes dCdθ for given data len(dCdθ) = 1+nFeatures
}

// RegData holds data for regression computations
type RegData struct {

	// input
	mData    int        // m: number of data points
	nParams  int        // n: number of features + 1
	xMat     *la.Matrix // [mData][nParams] matrix with the first column being filled with ones
	yVec     la.Vector  // [mData] y-data
	lVec     la.Vector  // [mData] l = X⋅θ (linear model)
	thetaVec la.Vector  // [nParams] parameters θ

	// control
	statOk bool // indicates that Stat is OK; otherwise Stat() must be called because x,y changed

	// stat
	minX, maxX  []float64 // [nFeatures] range of x
	sumX, meanX []float64 // [nFeatures] sum and mean of x
	sigX        []float64 // [nFeatures] standard deviations of x
	delX        []float64 // [nFeatures] difference: maxX - minX
	minY, maxY  float64   // range of y
	sumY, meanY float64   // sum and mean of y
	sigY        float64   // standard deviation of y
	delY        float64   // difference: maxY - minY

	// plotting
	horiz, vert  []float64   // [npts] horizontal and vertical arrays for curves
	thi, thj, cc [][]float64 // [npts][npts] meshgrid over θ
	xxi, xxj, zz [][]float64 // [npts][npts] meshgrid over x
}

// NewRegData returns a new structure to hold Regression Data
func NewRegData(nData, nFeatures int) (o *RegData) {

	// main
	o = new(RegData)
	o.mData = nData
	o.nParams = nFeatures + 1
	o.xMat = la.NewMatrix(o.mData, o.nParams)
	o.yVec = la.NewVector(o.mData)
	o.lVec = la.NewVector(o.mData)
	o.thetaVec = la.NewVector(o.nParams)
	o.xMat.SetCol(0, 1.0)

	// stat
	o.minX = make([]float64, nFeatures)
	o.maxX = make([]float64, nFeatures)
	o.sumX = make([]float64, nFeatures)
	o.meanX = make([]float64, nFeatures)
	o.sigX = make([]float64, nFeatures)
	o.delX = make([]float64, nFeatures)
	return
}

// NewRegDataTable sets X and Y values given table
//  xyRawTable -- [nData][nFeatures+1] table with x and y raw values,
//                where the last column contains y-values
func NewRegDataTable(xyRawTable [][]float64) (o *RegData) {
	nData := len(xyRawTable)
	if nData < 1 {
		chk.Panic("at leat one row of data in table must be provided\n")
	}
	nFeatures := len(xyRawTable[0]) - 1
	o = NewRegData(nData, nFeatures)
	for i := 0; i < nData; i++ {
		for j := 0; j < nFeatures; j++ {
			o.SetX(i, j, xyRawTable[i][j])
		}
		o.SetY(i, xyRawTable[i][nFeatures])
	}
	o.stat()
	return
}

// Ndata returns the number of data points
func (o *RegData) Ndata() int {
	return o.mData
}

// Nparams returns the number of parameters = len(θ)
func (o *RegData) Nparams() int {
	return o.nParams
}

// Nfeatures returns the number of features = number of parameters - 1
func (o *RegData) Nfeatures() int {
	return o.nParams - 1
}

// SetX sets x-value
func (o *RegData) SetX(iData, jFeature int, value float64) {
	o.xMat.Set(iData, 1+jFeature, value) // 1+j maps to augmented array
	o.statOk = false
}

// GetXvals get all x-values corresponding to feature iFeature
func (o *RegData) GetXvals(iFeature int) (xValues []float64) {
	return o.xMat.GetCol(1 + iFeature) // 1+j maps to augmented array
}

// SetY sets y-value
func (o *RegData) SetY(iData int, value float64) {
	o.yVec[iData] = value
	o.statOk = false
}

// GetYvals returns all y-values
func (o *RegData) GetYvals() (yValues []float64) {
	yValues = make([]float64, o.mData)
	copy(yValues, o.yVec)
	return
}

// Normalize normalizes x values
func (o *RegData) Normalize(useMinMax bool) {
	if !o.statOk {
		o.stat()
	}
	den := o.sigX
	if useMinMax {
		den = o.delX
	}
	for i := 0; i < o.mData; i++ {
		for j := 1; j < o.nParams; j++ {
			jFeature := j - 1
			o.xMat.Set(i, j, (o.xMat.Get(i, j)-o.meanX[jFeature])/den[jFeature])
		}
	}
	o.statOk = false
}

// plotting //////////////////////////////////////////////////////////////////////////////////////

// PlotX plots X values
//   args -- maps int(y) values to plot arguments
func (o *RegData) PlotX(iFeature, jFeature int, args map[int]*plt.A) {
	a := &plt.A{C: "orange", M: ".", Ls: "none", NoClip: true}
	i, j := 1+iFeature, 1+jFeature
	if args == nil {
		xivals := make([]float64, o.mData)
		xjvals := make([]float64, o.mData)
		for k := 0; k < o.mData; k++ {
			xivals[k] = o.xMat.Get(k, i)
			xjvals[k] = o.xMat.Get(k, j)
		}
		plt.Plot(xivals, xjvals, a)
	} else {
		var aa *plt.A
		for k := 0; k < o.mData; k++ {
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
func (o *RegData) PlotModel(reg Regression, iFeature, npts int, args *plt.A) {
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
		o.vert[k] = reg.Model(xVec, o.thetaVec)
	}
	plt.Plot(o.horiz, o.vert, args)
}

// PlotContModel plots contour of model
//  xmin -- minimum x[iFeature,jFeature] value for meshgrid; may be nil
//  xmax -- maximum x[iFeature,jFeature] value for meshgrid; may be nil
func (o *RegData) PlotContModel(reg Regression, iFeature, jFeature, npi, npj int, xmin, xmax []float64, filled bool, args *plt.A) {
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
	dxi := (xmax[iFeature] - xmin[iFeature]) / float64(npi-1)
	dxj := (xmax[jFeature] - xmin[jFeature]) / float64(npj-1)
	for r := 0; r < npj; r++ {
		for c := 0; c < npi; c++ {
			o.xxi[r][c] = xmin[iFeature] + float64(c)*dxi
			o.xxj[r][c] = xmin[jFeature] + float64(r)*dxj
			xVec[i] = o.xxi[r][c]
			xVec[j] = o.xxj[r][c]
			o.zz[r][c] = reg.Model(xVec, o.thetaVec)
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
func (o *RegData) PlotContCost(reg Regression, iPrm, jPrm, npi, npj int, minθ, maxθ []float64, args *plt.A) {
	if len(o.thi) != npj {
		o.thi = utl.Alloc(npj, npi)
		o.thj = utl.Alloc(npj, npi)
		o.cc = utl.Alloc(npj, npi)
	}
	θcpy := o.thetaVec.GetCopy()
	dthi := (maxθ[iPrm] - minθ[iPrm]) / float64(npi-1)
	dthj := (maxθ[jPrm] - minθ[jPrm]) / float64(npj-1)
	for r := 0; r < npj; r++ {
		for c := 0; c < npi; c++ {
			o.thi[r][c] = minθ[iPrm] + float64(c)*dthi
			o.thj[r][c] = minθ[jPrm] + float64(r)*dthj
			o.thetaVec[iPrm] = o.thi[r][c]
			o.thetaVec[jPrm] = o.thj[r][c]
			o.cc[r][c] = reg.Cost(o)
		}
	}
	o.thetaVec.Apply(1, θcpy)
	plt.ContourF(o.thi, o.thj, o.cc, args)
	plt.Gll(io.Sf("$\\theta_{%d}$", iPrm), io.Sf("$\\theta_{%d}$", jPrm), nil)
}

// auxiliary /////////////////////////////////////////////////////////////////////////////////////

// stat computes statistics
func (o *RegData) stat() {
	for J := 0; J < o.Nfeatures(); J++ {
		j := 1 + J
		o.minX[J] = o.xMat.Get(0, j)
		o.maxX[J] = o.minX[J]
		o.sumX[J] = 0.0
		for i := 0; i < o.mData; i++ {
			xval := o.xMat.Get(i, j)
			o.minX[J] = utl.Min(o.minX[J], xval)
			o.maxX[J] = utl.Max(o.maxX[J], xval)
			o.sumX[J] += xval
		}
		o.meanX[J] = o.sumX[J] / float64(o.mData)
		o.sigX[J] = rnd.StatDevFirst(o.xMat.Col(j), o.meanX[J], true)
		o.delX[J] = o.maxX[J] - o.minX[J]
	}
	o.minY = o.yVec[0]
	o.maxY = o.minY
	o.sumY = 0.0
	for i := 0; i < o.mData; i++ {
		o.minY = utl.Min(o.minY, o.yVec[i])
		o.maxY = utl.Max(o.maxY, o.yVec[i])
		o.sumY += o.yVec[i]
	}
	o.meanY = o.sumY / float64(o.mData)
	o.sigY = rnd.StatDevFirst(o.yVec, o.meanY, true)
	o.delY = o.maxY - o.minY
	o.statOk = true
}

// newXvec creates a new xVec[1+nFeatures] vector
func (o *RegData) newXvec(initWithMin, initWithMax, initWithMean bool) (xVec la.Vector) {
	xVec = la.NewVector(o.nParams)
	if initWithMin || initWithMax || initWithMax {
		if !o.statOk {
			o.stat()
		}
	}
	if initWithMin {
		copy(xVec, o.minX)
	}
	if initWithMax {
		copy(xVec, o.maxX)
	}
	if initWithMean {
		copy(xVec, o.meanX)
	}
	xVec[0] = 1.0
	return
}
